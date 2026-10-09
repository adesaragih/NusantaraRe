package main

// Perakitan mux aplikasi dari modul yang AKTIF - refactor bentuk B.
//
// Untuk apa berkas ini: `cmd/api` tidak mengenal isi modul. Ia menerima daftar
// modul terdaftar (`daftar.Rakit`), menyaringnya menurut MODUL_AKTIF, lalu
// hanya memasang modul yang lolos - rute dan pekerja latarnya. Modul yang
// nonaktif tidak punya rute sama sekali (jawabannya 404 dari mux), dan
// frontend menyembunyikan menunya dari GET /api/modul-aktif.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/daftar"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/login"
	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/templat"
	ruteTemplat "nusantarare/inti/backend/templat/rute"
)

// pilihModulAktif menyaring modul terdaftar menurut MODUL_AKTIF.
//
// Kosong = SEMUA modul (bawaan). Nama yang tidak dikenal DITOLAK: salah ketik
// di env yang diam-diam mematikan satu modul akan terbaca "modul itu memang
// tidak ada", dan tidak seorang pun tahu sebabnya.
//
// `namaLama` (`daftar.NamaLama`) memetakan nama modul sebelum tabel nama modul
// 30-09-2026 ke namanya kini: env yang masih memakainya DITOLAK dengan kalimat
// yang menyebut nama barunya - menerima diam-diam membuat dua nama untuk satu
// modul, menolak tanpa sebab membuat orang menebak.
func pilihModulAktif(terdaftar []inti.Modul, namaLama map[string]string, diminta []string) ([]inti.Modul, error) {
	if len(diminta) == 0 {
		return terdaftar, nil
	}
	dikenal := map[string]bool{}
	var semua []string
	for _, m := range terdaftar {
		dikenal[m.Nama()] = true
		semua = append(semua, m.Nama())
	}
	pilih := map[string]bool{}
	for _, n := range diminta {
		if baru, lama := namaLama[n]; lama && !dikenal[n] {
			return nil, fmt.Errorf("MODUL_AKTIF memakai nama modul lama %q; sejak 30-09-2026 namanya %q "+
				"(tabel nama modul, PANDUAN-DEPLOY-DAN-GIT-PER-MODUL.md)", n, baru)
		}
		if !dikenal[n] {
			sort.Strings(semua)
			return nil, fmt.Errorf("MODUL_AKTIF memuat modul yang tidak dikenal %q; yang dikenal: %s",
				n, strings.Join(semua, ", "))
		}
		pilih[n] = true
	}
	var aktif []inti.Modul
	for _, m := range terdaftar {
		if pilih[m.Nama()] {
			aktif = append(aktif, m)
		}
	}
	return aktif, nil
}

// rakitMux menyusun handler aplikasi: rute milik aplikasi, lalu rute setiap
// modul aktif menurut urutan daftar.
//
// Rute modul NONAKTIF dijawab 404 berbadan JSON `{galat}` yang menyebut
// modulnya - bukan teks `404 page not found` bawaan mux, yang frontend baca
// sebagai "jawaban bukan JSON, backend tidak terjangkau" dan menyuruh orang
// menyalakan ulang backend yang sedang berjalan (temuan /code-review paket
// 6-8). Saat semua modul aktif, handler yang dikembalikan ADALAH mux-nya:
// nol jawaban berubah.
//
// `GET /api/menu` (M_NAV_MENU, brief menu 30-09-2026) milik aplikasi: ia
// dipasang walau modul mana pun nonaktif, dan butir modul nonaktif tidak
// dikirimnya. `stubPelaku` = AUTH_STUB, diteruskan ke saringan per akun
// (`menu.SaringMenuUntukAkun`).
//
// `masuk` (login, keputusan work owner 01-10-2026) memasang `/api/auth/*` dan
// Kelola User `/api/admin/*`, dan membungkus SELURUH handler dengan
// middleware sesinya: pelaku hasil login dibaca `inti.PelakuDari` tanpa satu
// pun modul diubah. nil = tanpa login (uji) - dan tanpa gerbang menu.
//
// ⛔ GERBANG MENU (Kelola User, keputusan work owner 01-10-2026): rute milik
// modul AKTIF hanya dilayani bila akun yang login memegang menu modul itu
// (`M_LOGIN_GO_MENU`) - atau rute itu dipinjam modul yang dipegangnya
// (`ruteDipinjam`). Sesi tanpa menunya 403; tanpa sesi 401, kecuali
// AUTH_STUB=true (pengembangan, tanpa login). Pemilik rute dikenali dari
// mux terpisah per modul - nol modul diubah.
//
// ⛔ GERBANG TULIS (hak menu LIHAT, keputusan work owner 04-10-2026): `lihat` = modul yang mendaftar akses LIHAT
// (`daftar.HakLihat`) beserta pola rute yang dibebaskannya. Pemegang menu itu ber-hak LIHAT hanya dilayani GET/HEAD
// dan pola bebas; tulis lain 403. Modul yang tidak mendaftar tidak tersentuh.
//
// `aplikasi` - pemasang rute milik APLIKASI (bukan modul), mis. Template Manager: rutenya tidak melewati gerbang
// menu modul dan memeriksa aksesnya sendiri.
func rakitMux(dasar *inti.Dasar, terdaftar, aktif []inti.Modul, stubPelaku bool, masuk *login.Rute,
	lihat map[string][]string, aplikasi ...func(*http.ServeMux)) http.Handler {
	bungkus := func(h http.Handler) http.Handler { return h }
	mux := http.NewServeMux()
	if masuk != nil {
		masuk.Pasang(mux)
		bungkus = masuk.Middleware
	}
	mux.HandleFunc("GET /healthz", healthz(dasar))
	mux.HandleFunc("GET /api/modul-aktif", modulAktif(aktif))
	mux.HandleFunc("GET /api/menu", ruteMenu(dasar, aktif, stubPelaku))
	for _, pasang := range aplikasi {
		pasang(mux)
	}
	dipasang := map[string]bool{}
	var milikAktif []ruteModul
	for _, m := range aktif {
		m.DaftarkanRute(mux)
		dipasang[m.Nama()] = true
		if masuk != nil {
			milikAktif = append(milikAktif, kenali(m))
		}
	}
	var nonaktif []ruteModul
	for _, m := range terdaftar {
		if !dipasang[m.Nama()] {
			nonaktif = append(nonaktif, kenali(m))
		}
	}
	if len(nonaktif) == 0 && len(milikAktif) == 0 {
		return bungkus(mux)
	}
	return bungkus(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pola := mux.Handler(r)
		if pola == "" {
			for _, n := range nonaktif {
				if _, p := n.mux.Handler(r); p != "" {
					galat.Tulis(w, http.StatusNotFound,
						fmt.Sprintf("modul %s tidak aktif di proses ini (MODUL_AKTIF)", n.nama))
					return
				}
			}
		}
		if pemilik := pemilikPola(milikAktif, r, pola); pemilik != "" {
			if !izinMenu(w, r, stubPelaku, pemilik, ruteDipinjam[pola]) || !izinTulis(w, r, pemilik, pola, lihat) {
				return
			}
		}
		mux.ServeHTTP(w, r)
	}))
}

// ruteModul - rute satu modul, didaftarkan ke mux TERPISAH yang hanya dipakai
// untuk mengenali jalurnya; handler-nya tidak pernah dipanggil.
type ruteModul struct {
	nama string
	mux  *http.ServeMux
}

func kenali(m inti.Modul) ruteModul {
	kenal := http.NewServeMux()
	m.DaftarkanRute(kenal)
	return ruteModul{nama: m.Nama(), mux: kenal}
}

// pemilikPola - modul aktif yang mendaftarkan `pola` (pola yang dipilih mux
// aplikasi untuk permintaan ini); kosong = rute milik aplikasi.
func pemilikPola(milik []ruteModul, r *http.Request, pola string) string {
	if pola == "" {
		return ""
	}
	for _, m := range milik {
		if _, p := m.mux.Handler(r); p == pola {
			return m.nama
		}
	}
	return ""
}

// ruteDipinjam - rute modul yang DIPANGGIL LAYAR modul lain, menurut pola mux
// pemiliknya: pemegang menu peminjam boleh memakainya walau tidak memegang
// menu pemiliknya. Daftarnya dijaga dua arah (`TestPanggilanLintasModulTerdaftar`).
var ruteDipinjam = map[string][]string{
	// Panel Data Polis Claim Life membaca polis versi berjalan dari PremiumList
	// Life (`modul/claimlife/frontend/api.ts`, butir av).
	"GET /api/polis-life/ringkas": {"claimlife"},
	// Tombol "+" Consultant / Adjuster Claim Prop menambah master adjuster lewat API modul Adjuster Consultant
	// (Pega MstAdjusterConsultant; keputusan work owner 08-10-2026, `modul/claimprop/frontend/api.ts`).
	"POST /api/adjuster-consultant": {"claimprop"},
	// Komite Claim Prop tanpa menu sendiri: tabel komite di inbox Claim Prop dan layar komitenya (jendela
	// `JENDELA_DIPINJAM` frontend/App.tsx) dipakai pemegang menu Claim Prop; siapa yang boleh memutus tetap dijaga
	// layanan komite (anggota workbasket tingkat berjalan). Keputusan work owner 09-10-2026.
	"GET /api/komite-claim-prop/kasus":                {"claimprop"},
	"GET /api/komite-claim-prop/kasus/{id}":           {"claimprop"},
	"POST /api/komite-claim-prop/kasus/{id}/putuskan": {"claimprop"},
	// Komite Claim Non Prop - pola yang sama atas inbox Claim Non Prop (perintah work owner 09-10-2026).
	"GET /api/komite-claim-non-prop/kasus":                {"claimnonprop"},
	"GET /api/komite-claim-non-prop/kasus/{id}":           {"claimnonprop"},
	"POST /api/komite-claim-non-prop/kasus/{id}/putuskan": {"claimnonprop"},
	// Komite Claim Life - menu komite dihapus (perintah work owner 09-10-2026, "anggap menu itu tidak pernah ada"):
	// tabel komite inbox Claim Life dan layar kasusnya (modul tanpa menu, `layar.ts`) dipakai pemegang menu Claim Life;
	// siapa yang boleh memutus tetap dijaga layanan komite (anggota tangga tingkat berjalan).
	"GET /api/komite":                 {"claimlife"},
	"GET /api/komite/laporan-harian":  {"claimlife"},
	"GET /api/komite/{id}":            {"claimlife"},
	"GET /api/komite/{id}/riwayat":    {"claimlife"},
	"POST /api/komite/{id}/keputusan": {"claimlife"},
	"POST /api/komite/{id}/eskalasi":  {"claimlife"},
}

// izinMenu menjawab apakah permintaan ini boleh memakai rute milik `pemilik`,
// atau menulis penolakannya.
func izinMenu(w http.ResponseWriter, r *http.Request, stubPelaku bool, pemilik string, peminjam []string) bool {
	kode, sesi := inti.AksesMenuDari(r.Context())
	if !sesi {
		if stubPelaku {
			return true
		}
		galat.Tulis(w, http.StatusUnauthorized, "belum login atau sesi sudah berakhir")
		return false
	}
	if inti.PunyaMenu(kode, pemilik) {
		return true
	}
	for _, p := range peminjam {
		if inti.PunyaMenu(kode, p) {
			return true
		}
	}
	galat.Tulis(w, http.StatusForbidden, fmt.Sprintf("akun ini tidak punya akses ke menu %s", pemilik))
	return false
}

// izinTulis menjawab apakah permintaan TULIS ke rute milik `pemilik` boleh lewat bagi akun ber-hak menu LIHAT, atau
// menulis penolakannya. GET/HEAD, modul yang tidak mendaftar (`lihat`), hak PENUH, dan pola yang dibebaskan modulnya
// selalu lewat. Superadmin TIDAK dikecualikan (keputusan work owner 05-10-2026 "ikuti B": View only berlaku juga untuk superadmin).
func izinTulis(w http.ResponseWriter, r *http.Request, pemilik, pola string, lihat map[string][]string) bool {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return true
	}
	bebas, daftar := lihat[pemilik]
	if !daftar || menu.BolehUbah(r.Context(), pemilik) {
		return true
	}
	for _, b := range bebas {
		if b == pola {
			return true
		}
	}
	galat.Tulis(w, http.StatusForbidden, fmt.Sprintf("akun ini hanya dapat melihat menu %s (View only)", pemilik))
	return false
}

// rakitTemplat menyusun Template Manager (keputusan work owner 04-10-2026) dari slot SEMUA modul terdaftar.
// Slot cacat atau kode ganda MENOLAK menyala - kesalahan modul terlihat saat proses dinyalakan. Tanpa Oracle,
// berkas bawaan tetap dapat diunduh; unggah menjawab 503.
func rakitTemplat(dasar *inti.Dasar, slot []templat.Slot, stubPelaku bool) (func(*http.ServeMux), error) {
	kat, err := templat.NewKatalog(slot)
	if err != nil {
		return nil, err
	}
	var g templat.Gudang
	if dasar.PunyaDatabase() {
		g = templat.NewGudangOracle(dasar.DB())
	}
	l := templat.NewLayanan(kat, g)
	return func(mux *http.ServeMux) { ruteTemplat.Pasang(mux, l, stubPelaku) }, nil
}

// rakitLogin menyusun rute login dan Kelola User. Tanpa Oracle atau tanpa
// SESI_RAHASIA keduanya tetap terpasang dan menjawab 503 yang menyebut sebabnya.
func rakitLogin(dasar *inti.Dasar, cfg config.Config) *login.Rute {
	if !dasar.PunyaDatabase() || cfg.SesiRahasia == "" {
		return login.NewRute(nil, cfg.SesiCookieAman)
	}
	gudang := login.NewGudangOracle(dasar.DB())
	l := login.NewLayanan(gudang, []byte(cfg.SesiRahasia))
	return login.NewRute(l, cfg.SesiCookieAman).DenganKelola(login.NewKelola(gudang, menu.NewPembaca(dasar.DB())).
		DenganHakLihat(kodeHakLihat(daftar.HakLihat())))
}

// jawabanModulAktif adalah badan GET /api/modul-aktif.
type jawabanModulAktif struct {
	Modul []string `json:"modul"`
}

// modulAktif menyebut modul yang dipasang proses ini - yang dibaca frontend
// untuk menampilkan menunya. Daftarnya tidak pernah kosong: kosong di env
// berarti semua.
func modulAktif(aktif []inti.Modul) http.HandlerFunc {
	jawab := jawabanModulAktif{Modul: []string{}}
	for _, m := range aktif {
		jawab.Modul = append(jawab.Modul, m.Nama())
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jawab)
	}
}

// ruteMenu merakit `GET /api/menu` di atas pembaca Oracle - atau tanpa
// pembaca bila proses berjalan tanpa database (jawabannya 503 bergalat).
func ruteMenu(dasar *inti.Dasar, aktif []inti.Modul, stubPelaku bool) http.HandlerFunc {
	var nama []string
	for _, m := range aktif {
		nama = append(nama, m.Nama())
	}
	// ⚠️ nil ANTARMUKA, bukan *menu.Pembaca bernilai nil: yang kedua tidak
	// sama dengan nil dan akan dipanggil.
	var pembaca menu.PembacaMenu
	if dasar.PunyaDatabase() {
		pembaca = menu.NewPembaca(dasar.DB())
	}
	return menu.Rute(pembaca, nama, stubPelaku)
}

type jawabanSehat struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// healthz menjawab tanpa menyentuh aturan dagang mana pun.
//
// Ia tetap menjawab 200 ketika Oracle belum dikonfigurasi: Fase 0 harus dapat
// dijalankan tanpa instance, dan keadaan database dilaporkan apa adanya di
// dalam badan jawaban, bukan disembunyikan.
//
// Refactor bentuk B (30-09-2026): dipindah apa adanya dari handlers Claim
// Life - ia milik aplikasi, bukan satu modul, dan tetap ada walau modul mana
// pun nonaktif.
func healthz(dasar *inti.Dasar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jawab := jawabanSehat{Status: "sehat", Database: "tidak dikonfigurasi"}

		if dasar.PunyaDatabase() {
			ctx, batal := context.WithTimeout(r.Context(), 3*time.Second)
			defer batal()
			if err := dasar.CekKesehatan(ctx); err != nil {
				jawab.Database = "tidak terjangkau"
			} else {
				jawab.Database = "terjangkau"
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jawab)
	}
}

// jalankanPekerja menyalakan pekerja latar setiap modul aktif.
func jalankanPekerja(ctx context.Context, aktif []inti.Modul) []inti.Pekerja {
	var semua []inti.Pekerja
	for _, m := range aktif {
		semua = append(semua, m.JalankanPekerja(ctx))
	}
	return semua
}

// tungguPekerja menunggu setiap pekerja berhenti sampai batas `tutup`, dan
// mencetak pesan modul yang pekerjanya belum berhenti.
//
// Pekerja tanpa kanal `Selesai` (nilai nol `inti.Pekerja{}`) dianggap sudah
// berhenti: menunggu kanal nil memakan seluruh batas penutupan tanpa satu
// baris log pun (temuan /code-review).
func tungguPekerja(tutup context.Context, semua []inti.Pekerja, catat func(string)) {
	for _, p := range semua {
		if p.Selesai == nil {
			continue
		}
		select {
		case <-p.Selesai:
		case <-tutup.Done():
			if p.PesanTerlambat != "" {
				catat(p.PesanTerlambat)
			}
		}
	}
}

// kodeHakLihat - KODE menu modul yang mendaftar akses LIHAT, urut.
func kodeHakLihat(lihat map[string][]string) []string {
	kode := make([]string, 0, len(lihat))
	for k := range lihat {
		kode = append(kode, k)
	}
	sort.Strings(kode)
	return kode
}
