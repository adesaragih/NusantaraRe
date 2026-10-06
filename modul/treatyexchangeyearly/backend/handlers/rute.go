// Package handlers memasang API modul Treaty Exchange Yearly - `/api/treaty-exchange-yearly` (keputusan work owner
// 05-10-2026). Nol aturan dagang di sini; nol impor repository.
//
//	GET  /api/treaty-exchange-yearly?q=&tahun=   daftar (cari ID / kode / nama mata uang; saring tahun treaty)
//	GET  /api/treaty-exchange-yearly/pilihan     pilihan Currency (view CURRENCY) dan tahun treaty yang ada
//	POST /api/treaty-exchange-yearly             Add (ID = situs aktif + TREATYEXCHANGE_SEQ)
//	PUT  /api/treaty-exchange-yearly             Edit - baris dikenali `kunci` (ROWID) di badan, karena ID warisan tidak unik
//
// Tanpa hapus (keputusan work owner). Gerbang menu `treatyexchangeyearly` dan gerbang tulis View only dipasang `cmd/api`.
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
	"nusantarare/modul/treatyexchangeyearly/backend/models"
	"nusantarare/modul/treatyexchangeyearly/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/treaty-exchange-yearly"

// KodeMenu - KODE menu modul ini (`M_LOGIN_GO_MENU.MENU_KODE`, sama dengan nama modul).
const KodeMenu = "treatyexchangeyearly"

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
		d, err := l.Daftar(r.Context(), q.Get("q"), q.Get("tahun"))
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, struct {
			Daftar []models.Kurs `json:"daftar"`
		}{d})
	})
	pasang("GET "+Prefix+"/pilihan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ services.Aktor) {
		m, th, err := l.Pilihan(r.Context())
		if jawabGalat(w, err, "membaca pilihan") {
			return
		}
		galat.TulisJSON(w, struct {
			MataUang []models.MataUang `json:"mataUang"`
			Tahun    []string          `json:"tahun"`
		}{m, th})
	})
	simpan := func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor, edit bool) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi) {
			return
		}
		if edit != (strings.TrimSpace(isi.Kunci) != "") {
			galat.Tulis(w, http.StatusBadRequest, "kunci is required for Edit (PUT) and must be empty for Add (POST)")
			return
		}
		hasil, err := l.Simpan(r.Context(), a, isi)
		if jawabGalat(w, err, "menyimpan") {
			return
		}
		galat.TulisJSON(w, hasil)
	}
	pasang("POST "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, false)
	})
	pasang("PUT "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, a services.Aktor) {
		simpan(w, r, l, a, true)
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
		galat.Tulis(w, http.StatusNotFound, "Exchange rate not found")
	case errors.Is(err, services.ErrDilarang):
		galat.Tulis(w, http.StatusForbidden, pesan(err, services.ErrDilarang))
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrBelumAda):
		log.Printf("treatyexchangeyearly: %s: %v", apa, err)
		galat.Tulis(w, http.StatusServiceUnavailable, "Treaty Exchange Yearly: table, view, or sequence is not in this schema")
	default:
		log.Printf("treatyexchangeyearly: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Treaty Exchange Yearly: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
