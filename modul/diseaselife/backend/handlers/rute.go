// Package handlers memasang API modul Disease Life - `/api/disease-life` (pola causeoflosslife / benefitlife;
// keputusan work owner 08-10-2026 D1-D4, K5). Nol aturan dagang di sini; nol impor repository.
//
//	GET  /api/disease-life?icd=&disease=&urut=&arah=&halaman=   grid `BrowseDiseaseLife_RD` b5094 (10 per halaman b5169,
//	                                                            bawaan ID menurun b5012; saring b5148 dan halaman DI SERVER)
//	POST /api/disease-life                                      Save - Add (`AddToList_Act` b2121) -> 201
//	PUT  /api/disease-life/{id}                                 Save - Edit (`EditList_DT` b4852 lalu `AddToList_Act`)
//
// TIDAK ada DELETE (XML tanpa Delete: `pyGridDeleteActivityExists` false b5297), unggah, ringkasan, maupun rincian.
// Cancel (`NewData_DT` b2388) murni layar.
//
// Gerbang menu `diseaselife` dan gerbang tulis View only dipasang `cmd/api` (HakLihat tanpa pola bebas); layanan
// menolak lagi lewat Aktor.Penuh (403).
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
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/diseaselife/backend/models"
	"nusantarare/modul/diseaselife/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/disease-life"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "diseaselife"

// batasBadan - form kecil (ICD Code + Disease <= 1000 byte).
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

type rute func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor)

// aktorDari - pelaku sesi (atau stub) dan hak menunya (Full = bukan View only).
func aktorDari(r *http.Request, stub bool) services.Aktor {
	p := inti.PelakuDari(r, stub)
	return services.Aktor{AkunID: strings.TrimSpace(p.AkunID), Penuh: menu.BolehUbah(r.Context(), KodeMenu)}
}

func halaman(r *http.Request) int {
	n, err := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("halaman")))
	if err != nil || n < 1 {
		return 1
	}
	return n
}

func daftarkan(mux *http.ServeMux, layanan func() *services.Layanan, adaDB func() bool, stub bool) {
	pasang := func(pola string, f rute) {
		mux.HandleFunc(pola, func(w http.ResponseWriter, r *http.Request) {
			if !adaDB() {
				galat.Tulis(w, http.StatusServiceUnavailable, "database is not configured")
				return
			}
			f(w, r, layanan(), aktorDari(r, stub))
		})
	}
	pasang("GET "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		q := r.URL.Query()
		d, err := l.Daftar(r.Context(), models.Saringan{ICDCode: q.Get("icd"), Disease: q.Get("disease"),
			Urut: strings.ToLower(strings.TrimSpace(q.Get("urut"))), Naik: strings.EqualFold(q.Get("arah"), "asc"),
			Halaman: halaman(r)})
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, d)
	})
	simpan := func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor, id string) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi) {
			return
		}
		hasil, err := l.Simpan(r.Context(), a, id, isi)
		if jawabGalat(w, err, "menyimpan") {
			return
		}
		if id == "" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(hasil)
			return
		}
		galat.TulisJSON(w, hasil)
	}
	pasang("POST "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, "")
	})
	pasang("PUT "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			galat.Tulis(w, http.StatusNotFound, "Disease not found")
			return
		}
		simpan(w, r, l, a, id)
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

func pesan(err, penanda error) string { return strings.TrimPrefix(err.Error(), penanda.Error()+": ") }

// jawabGalat menulis galat; true = sudah dijawab. ID dari sequence yang sudah ada = 409 berkalimat, bukan 500.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Disease not found")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrIDTerpakai):
		log.Printf("diseaselife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusConflict, "Disease: "+pesan(err, services.ErrIDTerpakai))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("diseaselife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable,
			"Disease: table or sequence (DISEASE_LIFE, SEQ_DISEASE_LIFE) is not in this schema")
	default:
		log.Printf("diseaselife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Disease: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
