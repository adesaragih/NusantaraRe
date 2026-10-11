// Package handlers memasang API modul Endorsment Fac In - `/api/endorsmentfacin` (kontrak sesi utama 11-10-2026).
// Nol aturan dagang di sini; nol impor repository.
//
//	POST /api/endorsmentfacin/periksa   periksa sebelum Create (nol tulis): enam gerbang, jenis bisnis, tanggal
//	POST /api/endorsmentfacin/kasus     buat kasus EDM-: 201 {caseId,pesan} lahir · 200 {caseId:"",periksa} ditolak
//
// Galat: 400 badan rusak · 403 tanpa akun · 413 badan terlalu besar · 422 isian tidak sah · 503 tanpa Oracle atau
// tabel/sequence modul belum dibuat · 500 lain (rincian di log server).
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/endorsmentfacin"

// batasBadan - badan JSON form periksa; isian jauh di bawahnya.
const batasBadan = 16 << 10

// Router - mux berisi rute modul (Oracle).
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// RouterDengan - mux di atas layanan tertentu (uji dengan tiruan).
func RouterDengan(l *services.Layanan, adaDB, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	daftarkan(mux, func() *services.Layanan { return l }, func() bool { return adaDB }, stubPelaku)
	return mux
}

// DaftarkanRute memasang rute ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	daftarkan(mux, func() *services.Layanan { return services.LayananOracle(svc) }, svc.PunyaDatabase, stubPelaku)
}

func daftarkan(mux *http.ServeMux, layanan func() *services.Layanan, adaDB func() bool, stub bool) {
	mux.HandleFunc("POST "+Prefix+"/periksa", func(w http.ResponseWriter, r *http.Request) {
		if !adaDB() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		var isi models.MasukanPeriksa
		if !bacaBadan(w, r, &isi) {
			return
		}
		h, err := layanan().Periksa(r.Context(), isi)
		if jawabGalat(w, err, "memeriksa") {
			return
		}
		galat.TulisJSON(w, h)
	})
	mux.HandleFunc("POST "+Prefix+"/kasus", func(w http.ResponseWriter, r *http.Request) {
		if !adaDB() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
			return
		}
		var isi models.MasukanKasus
		if !bacaBadan(w, r, &isi) {
			return
		}
		akun := strings.TrimSpace(inti.PelakuDari(r, stub).AkunID)
		h, err := layanan().Buat(r.Context(), akun, isi)
		if jawabGalat(w, err, "membuat kasus") {
			return
		}
		if h.CaseID == "" {
			galat.TulisJSON(w, h) // 200: ditolak / tidak lengkap - hasil bisnis, bukan galat
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(h); err != nil {
			log.Printf("endorsmentfacin: menulis jawaban kasus %s: %v", h.CaseID, err)
		}
	})
}

// bacaBadan membaca badan JSON; galat = 400 / 413.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan))
	d.DisallowUnknownFields()
	if err := d.Decode(ke); err != nil {
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			galat.Tulis(w, http.StatusRequestEntityTooLarge, "request body is too large")
			return false
		}
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON for this form")
		return false
	}
	return true
}

// jawabGalat menulis galat; true = sudah dijawab.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, strings.TrimPrefix(err.Error(), services.ErrMasukanTidakSah.Error()+": "))
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, strings.TrimPrefix(err.Error(), services.ErrDilarang.Error()+": "))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("endorsmentfacin: %s: %v", apa, err)
		// Generik: objek yang hilang bisa milik modul (260/262) atau warisan (mis. JSON_POLIS); namanya ada di log.
		galat.Tulis(w, http.StatusServiceUnavailable,
			"Endorsment Fac In: a required table or sequence is not in this schema; details are in the server log")
	default:
		log.Printf("endorsmentfacin: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Endorsment Fac In: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
