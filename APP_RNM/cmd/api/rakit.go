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
)

// pilihModulAktif menyaring modul terdaftar menurut MODUL_AKTIF.
//
// Kosong = SEMUA modul (bawaan). Nama yang tidak dikenal DITOLAK: salah ketik
// di env yang diam-diam mematikan satu modul akan terbaca "modul itu memang
// tidak ada", dan tidak seorang pun tahu sebabnya.
func pilihModulAktif(terdaftar []inti.Modul, diminta []string) ([]inti.Modul, error) {
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

// rakitMux menyusun mux aplikasi: rute milik aplikasi, lalu rute setiap modul
// aktif menurut urutan daftar.
func rakitMux(dasar *inti.Dasar, aktif []inti.Modul) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz(dasar))
	mux.HandleFunc("GET /api/modul-aktif", modulAktif(aktif))
	for _, m := range aktif {
		m.DaftarkanRute(mux)
	}
	return mux
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
func tungguPekerja(tutup context.Context, semua []inti.Pekerja, catat func(string)) {
	for _, p := range semua {
		select {
		case <-p.Selesai:
		case <-tutup.Done():
			if p.PesanTerlambat != "" {
				catat(p.PesanTerlambat)
			}
		}
	}
}
