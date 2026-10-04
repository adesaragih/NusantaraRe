// Package handlers memegang rute HTTP modul Treaty In Adjustment
// (`treatyinadjustment`).
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/treaty-in-adjustment"

// Router menyusun rute modul ini SAJA di atas Oracle - dipakai uji HTTP.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// RouterDengan menyusun rute di atas Layanan yang sudah disusun (uji dengan
// gudang tiruan). `adaDB` false = setiap rute menjawab 503.
func RouterDengan(l *services.Layanan, adaDB, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	daftarkan(mux, func() *services.Layanan { return l }, func() bool { return adaDB }, stubPelaku)
	return mux
}

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	daftarkan(mux, func() *services.Layanan { return services.LayananOracle(svc) }, svc.PunyaDatabase, stubPelaku)
}

// rute adalah satu penangan yang menerima Layanan dan Pelaku yang sudah jadi.
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
	daftarkanBaca(pasang)
}

// daftarkanBaca - dua jalur baca.
//
// ⛔ Ini SELURUH permukaan HTTP modul ini hari ini, dan itu disengaja. Papan
// tiket modul ini menyatakan `L-4` ("tidak ada spesifikasi layar di mana pun")
// dan melarang mengarang layar untuk memenuhi bentuk irisan tegak. Jalur SIMPAN
// - yang menegakkan materialitas (tiket 02) dan batas tanggal berlaku (tiket
// 03) - lahir bersama spesifikasinya, bukan sebelum.
func daftarkanBaca(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/kontrak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		k, err := l.DaftarKontrak(r.Context(), p)
		tulis(w, k, err)
	})
	pasang("GET "+Prefix+"/kontrak/{id}/versi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			// Pengenal yang bukan angka ditolak DI SINI: ia bentuk permintaan
			// yang salah, bukan kontrak yang tidak ada.
			galat.Tulis(w, http.StatusBadRequest, "contract id must be a number")
			return
		}
		v, err := l.RantaiVersi(r.Context(), p, id)
		tulis(w, v, err)
	})
}

func tulis(w http.ResponseWriter, isi any, err error) {
	if jawabGalat(w, err) {
		return
	}
	galat.TulisJSON(w, isi)
}

// jawabGalat menerjemahkan galat services menjadi kode HTTP berbadan
// `{"galat": ...}`. Mengembalikan true bila permintaan SUDAH dijawab.
func jawabGalat(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "insufficient permission")
	case errors.Is(err, services.ErrIDTidakSah):
		galat.Tulis(w, http.StatusBadRequest, services.Pesan(err))
	case errors.Is(err, services.ErrKontrakTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	default:
		// Galat tak terduga: sebab aslinya ke log server, kalimat umum ke layar.
		log.Printf("treaty in adjustment: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the treaty in adjustment request")
	}
	return true
}
