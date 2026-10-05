// Package handlers memasang API modul Treaty Description - `/api/treaty-description` (keputusan work owner
// 05-10-2026). Nol aturan dagang di sini; nol impor repository.
//
//	GET  /api/treaty-description?q=&xol=&status=   daftar (cari ID / nama; saring Non XOL / XOL dan Active / Inactive)
//	GET  /api/treaty-description/{id}               satu baris dan jumlah pemakaiannya
//	POST /api/treaty-description                    Add (ID = '1' + TREATY_DESCRIPTION_SEQ)
//	PUT  /api/treaty-description/{id}               Edit
//
// Tanpa hapus. Gerbang menu `treatydescription` dan gerbang tulis View only dipasang `cmd/api`.
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/treatydescription/backend/models"
	"nusantarare/modul/treatydescription/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/treaty-description"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "treatydescription"

// batasBadan - badan JSON form; isian jauh di bawahnya.
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
		d, err := l.Daftar(r.Context(), q.Get("q"), q.Get("xol"), q.Get("status"))
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.Desc `json:"daftar"`
		}{d})
	})
	pasang("GET "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		d, err := l.Buka(r.Context(), r.PathValue("id"))
		if jawabGalat(w, err, "membaca") {
			return
		}
		galat.TulisJSON(w, d)
	})
	simpan := func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor, id string) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi) {
			return
		}
		isi.ID = id
		hasil, err := l.Simpan(r.Context(), a, isi)
		if jawabGalat(w, err, "menyimpan") {
			return
		}
		galat.TulisJSON(w, hasil)
	}
	pasang("POST "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, "")
	})
	pasang("PUT "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, strings.TrimSpace(r.PathValue("id")))
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

// jawabGalat menulis galat; true = sudah dijawab.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Treaty description not found")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("treatydescription: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable,
			"Treaty Description: table TREATYDESC or sequence TREATY_DESCRIPTION_SEQ is not in this schema")
	default:
		log.Printf("treatydescription: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Treaty Description: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
