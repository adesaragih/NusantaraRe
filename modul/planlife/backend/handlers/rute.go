// Package handlers memasang API modul Plan - `/api/plan-life` (pola benefitlife + dua autocomplete; keputusan work
// owner 08-10-2026 K1-K7). Nol aturan dagang di sini; nol impor repository.
//
//	GET  /api/plan-life?urut=&arah=&halaman=   grid `BrowseProductTypeLife_RD` b4325 (10 per halaman b4401)
//	GET  /api/plan-life/pilihan/business      autocomplete Business (`BrowseBusinessLife_RD` b1183; K4 GROUPPANEL 009)
//	GET  /api/plan-life/pilihan/benefit       autocomplete Benefit (`BrowseBenefitLife_RD` b1537)
//	POST /api/plan-life                       Save - Add (`SaveProductTypeLife_Act` b6422) -> 201
//	PUT  /api/plan-life/{id}                  Save sesudah Edit (`EditProductTypeLife_Act` b4045)
//
// TIDAK ada DELETE (XML tanpa Delete), unggah, ringkasan, maupun rincian. New (`NewProductTypeLife_act` b6686) murni
// layar. Gerbang menu `planlife` dan gerbang tulis View only dipasang `cmd/api`; layanan menolak lagi (403).
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
	"nusantarare/modul/planlife/backend/models"
	"nusantarare/modul/planlife/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/plan-life"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "planlife"

// batasBadan - form kecil (tiga teks).
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
		d, err := l.Daftar(r.Context(), models.Saringan{Urut: q.Get("urut"), Turun: strings.EqualFold(q.Get("arah"), "desc"),
			Halaman: halaman(r)})
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, d)
	})
	pasang("GET "+Prefix+"/pilihan/business", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.PilihanBusiness(r.Context())
		if jawabGalat(w, err, "membaca pilihan business") {
			return
		}
		galat.TulisJSON(w, d)
	})
	pasang("GET "+Prefix+"/pilihan/benefit", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.PilihanBenefit(r.Context())
		if jawabGalat(w, err, "membaca pilihan benefit") {
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
			galat.Tulis(w, http.StatusNotFound, "Plan not found")
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

// jawabGalat menulis galat; true = sudah dijawab. K3: ID dari sequence yang sudah ada = 409 berkalimat.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Plan not found")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrIDTerpakai):
		log.Printf("planlife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusConflict, "Plan: "+pesan(err, services.ErrIDTerpakai))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("planlife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable,
			"Plan: table or sequence (PRODUCT_TYPE_LIFE, BUSINESS, BENEFIT_LIFE, M_PRODUCT_TYPE_LIFE_SEQ) is not in this schema")
	default:
		log.Printf("planlife: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Plan: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
