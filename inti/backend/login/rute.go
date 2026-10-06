package login

// Rute `/api/auth/*` dan middleware sesi - dipasang cmd/api di depan SELURUH
// rute, sehingga `inti.PelakuDari` membaca pelaku dari sesi tanpa satu pun
// modul diubah.

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
)

// Rute melayani login, logout, profil, dan ganti sandi.
type Rute struct {
	layanan *Layanan
	// cookieAman - atribut Secure (SESI_COOKIE_SECURE, bawaan true).
	cookieAman bool
	// kelola - Kelola User (`DenganKelola`); nil = rute `/api/admin/*` 503.
	kelola *Kelola
}

// NewRute menyusunnya. `l` nil = login belum dapat dipakai (tanpa Oracle atau
// tanpa SESI_RAHASIA): rutenya tetap ada dan menjawab 503 yang menjelaskan diri.
func NewRute(l *Layanan, cookieAman bool) *Rute { return &Rute{layanan: l, cookieAman: cookieAman} }

// Pasang mendaftarkan rute `/api/auth/*` dan Kelola User `/api/admin/*`.
func (r *Rute) Pasang(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", r.masuk)
	mux.HandleFunc("POST /api/auth/logout", r.keluar)
	mux.HandleFunc("GET /api/auth/saya", r.saya)
	mux.HandleFunc("POST /api/auth/ganti-sandi", r.gantiSandi)
	r.pasangKelola(mux)
}

type kunciSesi struct{}

type sesiPermintaan struct {
	profil Profil
	token  Token
}

func sesiDari(ctx context.Context) (sesiPermintaan, bool) {
	s, ok := ctx.Value(kunciSesi{}).(sesiPermintaan)
	return s, ok
}

// buangSetCookieSesi menjamin SATU cookie sesi per jawaban: middleware sudah
// memasang cookie perpanjangan sebelum handler logout atau ganti sandi
// memasang miliknya sendiri.
func buangSetCookieSesi(w http.ResponseWriter) {
	h := w.Header()
	var sisa []string
	for _, v := range h.Values("Set-Cookie") {
		if !strings.HasPrefix(v, NamaCookie+"=") {
			sisa = append(sisa, v)
		}
	}
	h.Del("Set-Cookie")
	for _, v := range sisa {
		h.Add("Set-Cookie", v)
	}
}

func (r *Rute) pasangCookie(w http.ResponseWriter, t Token) {
	buangSetCookieSesi(w)
	http.SetCookie(w, &http.Cookie{
		Name: NamaCookie, Value: Tandatangani(r.layanan.rahasia, t), Path: "/",
		Expires: t.Habis, HttpOnly: true, Secure: r.cookieAman, SameSite: http.SameSiteStrictMode,
	})
}

func (r *Rute) hapusCookie(w http.ResponseWriter) {
	buangSetCookieSesi(w)
	http.SetCookie(w, &http.Cookie{
		Name: NamaCookie, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: r.cookieAman, SameSite: http.SameSiteStrictMode,
	})
}

// Middleware membaca cookie sesi dan menaruh pelakunya di context.
//
// ⛔ Cookie rusak, kedaluwarsa, atau dicabut DIBUANG dan permintaannya
// diteruskan TANPA pelaku - jalur beridentitas menolaknya sendiri. Galat
// Oracle tidak membuang cookie: sesinya mungkin masih sah.
//
// ⛔ Akun yang wajib ganti sandi mendapat sesi, tetapi BUKAN pelaku dan BUKAN
// pemegang menu: hanya `/api/auth/*` yang melayaninya sampai sandinya diganti.
//
// Menu akun (`inti.DenganAksesMenu`) dibaca ulang dari Oracle di SETIAP
// permintaan, bersama workbasket-nya: perubahan di Kelola User berlaku pada
// permintaan berikutnya, tanpa login ulang.
func (r *Rute) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		c, err := req.Cookie(NamaCookie)
		if err != nil || r.layanan == nil {
			next.ServeHTTP(w, req)
			return
		}
		p, tok, err := r.layanan.Sesi(req.Context(), c.Value)
		if err != nil {
			if errors.Is(err, ErrSesiTidakSah) {
				r.hapusCookie(w)
			} else {
				log.Printf("login: membaca sesi: %v", err)
			}
			next.ServeHTTP(w, req)
			return
		}
		r.pasangCookie(w, tok)
		ctx := context.WithValue(req.Context(), kunciSesi{}, sesiPermintaan{profil: p, token: tok})
		if !p.WajibGantiSandi {
			ctx = inti.DenganPelakuSesi(ctx, inti.Pelaku{AkunID: p.AkunID, Peran: p.Peran})
			ctx = inti.DenganAksesMenu(ctx, p.Menu)
			ctx = inti.DenganMenuLihat(ctx, p.MenuLihat)
		}
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

// bacaJSON - hanya `application/json` (formulir lintas situs tidak dapat
// mengirimnya), badan maksimal 4 KiB.
func bacaJSON(w http.ResponseWriter, req *http.Request, tujuan any) bool {
	return bacaJSONBatas(w, req, tujuan, 4096)
}

// bacaJSONBatas - `bacaJSON` dengan batas badan `batas` byte.
func bacaJSONBatas(w http.ResponseWriter, req *http.Request, tujuan any, batas int64) bool {
	if !strings.HasPrefix(strings.ToLower(req.Header.Get("Content-Type")), "application/json") {
		galat.Tulis(w, http.StatusUnsupportedMediaType, "badan permintaan wajib application/json")
		return false
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, req.Body, batas)).Decode(tujuan); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan tidak dapat dibaca")
		return false
	}
	return true
}

func (r *Rute) siap(w http.ResponseWriter) bool {
	if r.layanan == nil {
		galat.Tulis(w, http.StatusServiceUnavailable,
			"login belum dapat dipakai: ORACLE_DSN dan SESI_RAHASIA wajib terisi")
		return false
	}
	return true
}

func (r *Rute) masuk(w http.ResponseWriter, req *http.Request) {
	if !r.siap(w) {
		return
	}
	var m struct{ Akun, Sandi string }
	if !bacaJSON(w, req, &m) {
		return
	}
	p, tok, err := r.layanan.Masuk(req.Context(), m.Akun, m.Sandi)
	switch {
	case errors.Is(err, ErrKredensial):
		galat.Tulis(w, http.StatusUnauthorized, "akun atau sandi salah")
	case errors.Is(err, ErrTerkunci):
		galat.Tulis(w, http.StatusLocked, "akun terkunci sementara sesudah 5 kali sandi salah; coba lagi 15 menit lagi")
	case errors.Is(err, ErrTanpaRahasia):
		galat.Tulis(w, http.StatusServiceUnavailable, "login belum dapat dipakai: SESI_RAHASIA belum disetel")
	case errors.Is(err, ErrMenuBelumDimigrasi):
		galat.Tulis(w, http.StatusServiceUnavailable, "login belum dapat dipakai: "+strings.TrimPrefix(err.Error(), "login: "))
	case err != nil:
		log.Printf("login: masuk: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "login gagal; rinciannya di log server")
	default:
		r.pasangCookie(w, tok)
		galat.TulisJSON(w, p)
	}
}

func (r *Rute) saya(w http.ResponseWriter, req *http.Request) {
	s, ok := sesiDari(req.Context())
	if !ok {
		galat.Tulis(w, http.StatusUnauthorized, "belum login atau sesi sudah berakhir")
		return
	}
	galat.TulisJSON(w, s.profil)
}

// keluar mencabut SELURUH sesi akun itu dan membuang cookie - idempoten.
func (r *Rute) keluar(w http.ResponseWriter, req *http.Request) {
	if s, ok := sesiDari(req.Context()); ok {
		if err := r.layanan.Keluar(req.Context(), s.profil.AkunID); err != nil {
			log.Printf("login: keluar: %v", err)
			galat.Tulis(w, http.StatusInternalServerError, "logout gagal; rinciannya di log server")
			return
		}
	}
	r.hapusCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

// gantiSandi - galat sandi dijawab 400, bukan 401: 401 dibaca layar sebagai
// "sesi berakhir".
func (r *Rute) gantiSandi(w http.ResponseWriter, req *http.Request) {
	s, ok := sesiDari(req.Context())
	if !ok {
		galat.Tulis(w, http.StatusUnauthorized, "belum login atau sesi sudah berakhir")
		return
	}
	var m struct{ SandiLama, SandiBaru string }
	if !bacaJSON(w, req, &m) {
		return
	}
	p, tok, err := r.layanan.GantiSandi(req.Context(), s.token, m.SandiLama, m.SandiBaru)
	switch {
	case errors.Is(err, ErrKredensial):
		galat.Tulis(w, http.StatusBadRequest, "sandi lama salah")
	case errors.Is(err, ErrSandiTerlaluPendek), errors.Is(err, ErrSandiTerlaluPanjang), errors.Is(err, ErrSandiSama):
		galat.Tulis(w, http.StatusBadRequest, strings.TrimPrefix(err.Error(), "login: "))
	case errors.Is(err, ErrSesiTidakSah):
		r.hapusCookie(w)
		galat.Tulis(w, http.StatusUnauthorized, "sesi sudah berakhir")
	case err != nil:
		log.Printf("login: ganti sandi: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "ganti sandi gagal; rinciannya di log server")
	default:
		r.pasangCookie(w, tok)
		galat.TulisJSON(w, p)
	}
}
