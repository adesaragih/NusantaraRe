// Package handlers memasang API modul Treaty Group - `/api/treaty-group` (keputusan work owner 05-10-2026). Nol aturan
// dagang di sini; nol impor repository.
//
//	GET  /api/treaty-group?q=&ojk=     daftar (cari ID / nama / nama SOA / nama OJK; saring OJK), urut Order No
//	GET  /api/treaty-group/pilihan     pilihan OJK Business (TREATYGROUPOJK)
//	GET  /api/treaty-group/{id}        satu grup dan grup bisnis anaknya (tanpa SYARIAH)
//	POST /api/treaty-group             Add (ID = situs aktif + TREATYGROUP_SEQ)
//	PUT  /api/treaty-group/{id}        Edit
//
// Tanpa hapus (keputusan work owner). Gerbang menu `treatygroup` dan gerbang tulis View only dipasang `cmd/api`.
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
	"nusantarare/modul/treatygroup/backend/models"
	"nusantarare/modul/treatygroup/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/treaty-group"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "treatygroup"

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
		d, err := l.Daftar(r.Context(), q.Get("q"), q.Get("ojk"))
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.Grup `json:"daftar"`
		}{d})
	})
	pasang("GET "+Prefix+"/pilihan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		o, err := l.Pilihan(r.Context())
		if jawabGalat(w, err, "membaca pilihan") {
			return
		}
		galat.TulisJSON(w, struct {
			Ojk []models.Ojk `json:"ojk"`
		}{o})
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
		galat.Tulis(w, http.StatusNotFound, "Treaty group not found")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("treatygroup: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable, "Treaty Group: table or sequence TREATYGROUP is not in this schema")
	default:
		log.Printf("treatygroup: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Treaty Group: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
