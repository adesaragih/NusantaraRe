package main

// Perakitan mux aplikasi dari modul yang AKTIF - refactor bentuk B.
//
// Untuk apa berkas ini: `cmd/api` tidak mengenal isi modul. Ia menerima daftar
// modul terdaftar (`modul.Rakit`), menyaringnya menurut MODUL_AKTIF, lalu
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

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/menu"
)

// pilihModulAktif menyaring modul terdaftar menurut MODUL_AKTIF.
//
// Kosong = SEMUA modul (bawaan). Nama yang tidak dikenal DITOLAK: salah ketik
// di env yang diam-diam mematikan satu modul akan terbaca "modul itu memang
// tidak ada", dan tidak seorang pun tahu sebabnya.
//
// `namaLama` (`modul.NamaLama`) memetakan nama modul sebelum tabel nama modul
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
// (`menu.SaringMenuUntukPelaku`, hari ini meneruskan semua).
func rakitMux(dasar *inti.Dasar, terdaftar, aktif []inti.Modul, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(dasar))
	mux.HandleFunc("GET /api/modul-aktif", modulAktif(aktif))
	mux.HandleFunc("GET /api/menu", ruteMenu(dasar, aktif, stubPelaku))
	dipasang := map[string]bool{}
	for _, m := range aktif {
		m.DaftarkanRute(mux)
		dipasang[m.Nama()] = true
	}
	type ruteNonaktif struct {
		nama string
		mux  *http.ServeMux
	}
	var nonaktif []ruteNonaktif
	for _, m := range terdaftar {
		if dipasang[m.Nama()] {
			continue
		}
		// Rutenya didaftarkan ke mux TERPISAH yang hanya dipakai untuk
		// mengenali jalurnya - handler-nya tidak pernah dipanggil.
		kenal := http.NewServeMux()
		m.DaftarkanRute(kenal)
		nonaktif = append(nonaktif, ruteNonaktif{nama: m.Nama(), mux: kenal})
	}
	if len(nonaktif) == 0 {
		return mux
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pola := mux.Handler(r); pola == "" {
			for _, n := range nonaktif {
				if _, p := n.mux.Handler(r); p != "" {
					galat.Tulis(w, http.StatusNotFound,
						fmt.Sprintf("modul %s tidak aktif di proses ini (MODUL_AKTIF)", n.nama))
					return
				}
			}
		}
		mux.ServeHTTP(w, r)
	})
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
