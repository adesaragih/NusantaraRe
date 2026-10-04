// Package handlers adalah pintu HTTP modul Accounts.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON. Nol aturan dagang di sini; nol impor
// repository. Rute modul ini hanya terbuka bagi pemegang menu `accounts` (perakit `cmd/api`).
//
//	GET  /api/accounts                 satu halaman daftar (?cari=&halaman=&ukuran=), terbaru dulu
//	GET  /api/accounts/pilihan         pilihan Group Business (BUSINESSGROUP)
//	GET  /api/accounts/organisasi      pilihan Insured Name (?cari=) - organisasi CLIENT
//	POST /api/accounts                 tambah (INSERT) - akun hanya diinput sekali: nol ubah, nol hapus
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/accounts/backend/models"
	"nusantarare/modul/accounts/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/accounts"

// batasBadan - badan JSON form; Description 4000 byte jauh di bawahnya.
const batasBadan = 64 << 10

// Router menyusun rute modul ini SAJA di atas Oracle.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// RouterDengan menyusun rute di atas Layanan yang sudah disusun (uji dengan gudang tiruan). `adaDB` false = setiap
// rute menjawab 503.
func RouterDengan(l *services.Layanan, adaDB, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	daftarkan(mux, func() *services.Layanan { return l }, func() bool { return adaDB }, stubPelaku)
	return mux
}

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	daftarkan(mux, func() *services.Layanan { return services.LayananOracle(svc) }, svc.PunyaDatabase, stubPelaku)
}

type rute func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku)

func daftarkan(mux *http.ServeMux, layanan func() *services.Layanan, adaDB func() bool, stub bool) {
	pasang := func(pola string, f rute) {
		mux.HandleFunc(pola, func(w http.ResponseWriter, r *http.Request) {
			if !adaDB() {
				galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
				return
			}
			f(w, r, layanan(), inti.PelakuDari(r, stub))
		})
	}
	pasang("GET "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		q := r.URL.Query()
		h, err := l.Daftar(r.Context(), q.Get("cari"), angka(q.Get("halaman")), angka(q.Get("ukuran")))
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("GET "+Prefix+"/pilihan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		p, err := l.Pilihan(r.Context())
		if jawabGalat(w, err, "membaca pilihan") {
			return
		}
		galat.TulisJSON(w, p)
	})
	pasang("GET "+Prefix+"/organisasi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		h, err := l.CariOrganisasi(r.Context(), r.URL.Query().Get("cari"))
		if jawabGalat(w, err, "mencari organisasi") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("POST "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi) {
			return
		}
		a, err := l.Tambah(r.Context(), p, isi)
		if jawabGalat(w, err, "menambah") {
			return
		}
		galat.TulisJSON(w, a)
	})
}

// angka - bilangan bulat parameter kueri; kosong atau rusak = 0 (layanan memakai bawaannya).
func angka(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return n
}

// bacaBadan membaca badan JSON; galat = 400.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan))
	d.DisallowUnknownFields()
	if err := d.Decode(ke); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON for this form")
		return false
	}
	return true
}

// pesan - kalimat untuk layar: bagian sesudah nama galat penanda.
func pesan(err, penanda error) string {
	return strings.TrimPrefix(err.Error(), penanda.Error()+": ")
}

// jawabGalat menulis galat; true = sudah dijawab.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaPelaku):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrGanda):
		galat.Tulis(w, http.StatusConflict, pesan(err, services.ErrGanda))
	default:
		log.Printf("accounts: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Accounts: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
