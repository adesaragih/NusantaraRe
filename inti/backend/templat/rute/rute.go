// Package rute memasang API Template Manager - `/api/templat/*` (keputusan work owner 04-10-2026).
//
// Rute milik APLIKASI, bukan modul: gerbang menu `cmd/api` tidak melindunginya, jadi setiap rute memeriksa
// sendiri.
//   - Kelola (daftar, riwayat, periksa, unggah, aktifkan, unduh versi tertentu): pemegang menu Template Manager.
//   - Unduh versi aktif: pemegang menu PEMILIK slot (mis. `bordereaux`) atau Template Manager.
//
// Tanpa sesi login: ditolak 401, kecuali AUTH_STUB menyala (pengembangan) - sama dengan rute modul.
package rute

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/templat"
)

// Prefix rute API Template Manager.
const Prefix = "/api/templat"

// batasBadan - multipart: berkas (maks. templat.BatasUkuran) + catatan + kepala multipart.
const batasBadan = templat.BatasUkuran + 64<<10

type rute struct {
	l    *templat.Layanan
	stub bool
}

// Pasang mendaftarkan rute ke mux. l nil = setiap rute menjawab 503.
func Pasang(mux *http.ServeMux, l *templat.Layanan, stubPelaku bool) {
	r := rute{l: l, stub: stubPelaku}
	mux.HandleFunc("GET "+Prefix, r.daftar)
	mux.HandleFunc("GET "+Prefix+"/{kode}/riwayat", r.riwayat)
	mux.HandleFunc("POST "+Prefix+"/{kode}/periksa", r.periksa)
	mux.HandleFunc("POST "+Prefix+"/{kode}/unggah", r.unggah)
	mux.HandleFunc("POST "+Prefix+"/{kode}/aktifkan", r.aktifkan)
	mux.HandleFunc("GET "+Prefix+"/{kode}/unduh", r.unduh)
}

// izin menjawab apakah permintaan memegang salah satu `kode` menu, atau menulis penolakannya.
func (r rute) izin(w http.ResponseWriter, req *http.Request, kode ...string) bool {
	if r.l == nil {
		galat.Tulis(w, http.StatusServiceUnavailable, "Template Manager is not available: the template catalog is not loaded")
		return false
	}
	menuAkun, sesi := inti.AksesMenuDari(req.Context())
	if !sesi {
		if r.stub {
			return true
		}
		galat.Tulis(w, http.StatusUnauthorized, "belum login atau sesi sudah berakhir")
		return false
	}
	for _, k := range kode {
		if inti.PunyaMenu(menuAkun, k) {
			return true
		}
	}
	galat.Tulis(w, http.StatusForbidden, "This account does not have access to this template")
	return false
}

func (r rute) kelola(w http.ResponseWriter, req *http.Request) bool {
	return r.izin(w, req, menu.KodeTemplateManager)
}

func jawabGalat(w http.ResponseWriter, err error, apa string) {
	var gp templat.GalatPeriksa
	switch {
	case errors.As(err, &gp):
		galat.Tulis(w, http.StatusUnprocessableEntity, gp.Error())
	case errors.Is(err, templat.ErrSlotTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Template not found")
	case errors.Is(err, templat.ErrVersiTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Template version not found")
	case errors.Is(err, templat.ErrBelumDimigrasi), errors.Is(err, templat.ErrTanpaDatabase):
		log.Printf("templat: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable, "Template Manager needs the database: run -migrate (inti 912)")
	default:
		log.Printf("templat: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Template Manager: "+apa+" failed; see the server log")
	}
}

func (r rute) daftar(w http.ResponseWriter, req *http.Request) {
	if !r.kelola(w, req) {
		return
	}
	d, err := r.l.Daftar(req.Context())
	if err != nil {
		jawabGalat(w, err, "reading templates")
		return
	}
	galat.TulisJSON(w, struct {
		Daftar []templat.RingkasanSlot `json:"daftar"`
	}{d})
}

func (r rute) riwayat(w http.ResponseWriter, req *http.Request) {
	if !r.kelola(w, req) {
		return
	}
	vs, err := r.l.Riwayat(req.Context(), req.PathValue("kode"))
	if err != nil {
		jawabGalat(w, err, "reading versions")
		return
	}
	galat.TulisJSON(w, struct {
		Versi []templat.Versi `json:"versi"`
	}{vs})
}

// bacaBerkas membaca bagian `berkas` (dan `catatan`) multipart.
func bacaBerkas(w http.ResponseWriter, req *http.Request) (nama string, isi []byte, catatan string, ok bool) {
	req.Body = http.MaxBytesReader(w, req.Body, batasBadan)
	if err := req.ParseMultipartForm(batasBadan); err != nil {
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			galat.Tulis(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("File is larger than %d KB", templat.BatasUkuran>>10))
		} else {
			galat.Tulis(w, http.StatusBadRequest, "The request must be a multipart form with field berkas")
		}
		return "", nil, "", false
	}
	f, kepala, err := req.FormFile("berkas")
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "Choose a file to upload")
		return "", nil, "", false
	}
	defer func() { _ = f.Close() }()
	isi, err = io.ReadAll(io.LimitReader(f, templat.BatasUkuran+1))
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "The file cannot be read")
		return "", nil, "", false
	}
	return kepala.Filename, isi, req.FormValue("catatan"), true
}

func (r rute) periksa(w http.ResponseWriter, req *http.Request) {
	if !r.kelola(w, req) {
		return
	}
	nama, isi, _, ok := bacaBerkas(w, req)
	if !ok {
		return
	}
	h, err := r.l.Periksa(req.Context(), req.PathValue("kode"), nama, isi)
	if err != nil {
		jawabGalat(w, err, "checking the file")
		return
	}
	galat.TulisJSON(w, h)
}

func (r rute) unggah(w http.ResponseWriter, req *http.Request) {
	if !r.kelola(w, req) {
		return
	}
	akun := inti.PelakuDari(req, r.stub).AkunID
	if strings.TrimSpace(akun) == "" {
		galat.Tulis(w, http.StatusUnauthorized, "The uploader has no identity (log in, or X-Pelaku in stub mode)")
		return
	}
	nama, isi, catatan, ok := bacaBerkas(w, req)
	if !ok {
		return
	}
	v, err := r.l.Unggah(req.Context(), req.PathValue("kode"), nama, isi, catatan, akun)
	if err != nil {
		jawabGalat(w, err, "saving the file")
		return
	}
	galat.TulisJSON(w, struct {
		Versi int `json:"versi"`
	}{v})
}

func (r rute) aktifkan(w http.ResponseWriter, req *http.Request) {
	if !r.kelola(w, req) {
		return
	}
	var badan struct {
		Versi *int `json:"versi"`
	}
	dek := json.NewDecoder(http.MaxBytesReader(w, req.Body, 1<<10))
	dek.DisallowUnknownFields()
	if err := dek.Decode(&badan); err != nil || badan.Versi == nil {
		galat.Tulis(w, http.StatusBadRequest, `The body must be {"versi": <number>}; 0 = built-in file`)
		return
	}
	if err := r.l.Aktifkan(req.Context(), req.PathValue("kode"), *badan.Versi); err != nil {
		jawabGalat(w, err, "activating the version")
		return
	}
	galat.TulisJSON(w, struct {
		Versi int `json:"versi"`
	}{*badan.Versi})
}

func (r rute) unduh(w http.ResponseWriter, req *http.Request) {
	if r.l == nil {
		r.izin(w, req)
		return
	}
	kode := req.PathValue("kode")
	slot, err := r.l.Slot(kode)
	if err != nil {
		if r.izin(w, req, menu.KodeTemplateManager) {
			jawabGalat(w, err, "reading the template")
		}
		return
	}
	versi := -1
	if teks := req.URL.Query().Get("versi"); teks != "" {
		n, e := strconv.Atoi(teks)
		if e != nil || n < 0 {
			galat.Tulis(w, http.StatusBadRequest, "versi must be a number >= 0")
			return
		}
		versi = n
	}
	// Versi tertentu (riwayat) hanya untuk pengelola; versi aktif juga untuk pemegang menu pemilik.
	boleh := []string{menu.KodeTemplateManager}
	if versi < 0 {
		boleh = append(boleh, slot.Menu)
	}
	if !r.izin(w, req, boleh...) {
		return
	}
	b, err := r.l.Unduh(req.Context(), kode, versi)
	if err != nil {
		jawabGalat(w, err, "downloading the template")
		return
	}
	w.Header().Set("Content-Type", b.Mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`,
		namaAscii(b.Nama), url.PathEscape(b.Nama)))
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b.Isi)
}

// namaAscii - nama cadangan untuk peramban lama: selain ASCII cetak, `"` dan `\` diganti `_`.
func namaAscii(s string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r > 0x7e || r == '"' || r == '\\' {
			return '_'
		}
		return r
	}, s)
}
