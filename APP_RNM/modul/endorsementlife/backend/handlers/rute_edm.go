// Package handlers adalah pintu HTTP modul Endorsement Life.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON.
// Nol aturan dagang di sini; nol impor repository (penjaga
// `TestHandlersTidakMengimporRepository`).
//
//	GET  /api/endorsement-life/inbox?halaman=&ukuran=        kotak masuk `InboxEndorsementLife`
//	GET  /api/endorsement-life/kasus/{id}                    kepala `InputEDMLife`
//	GET  /api/endorsement-life/kasus/{id}/peserta?halaman=   grid peserta b11899 / b17500
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/endorsementlife/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/endorsement-life"

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

// angkaKueri membaca parameter kueri bilangan bulat; kosong/rusak = 0.
func angkaKueri(r *http.Request, nama string) int {
	n, err := strconv.Atoi(r.URL.Query().Get(nama))
	if err != nil {
		return 0
	}
	return n
}

func daftarkanBaca(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/inbox", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		h, err := l.Inbox(r.Context(), p, angkaKueri(r, "halaman"), angkaKueri(r, "ukuran"))
		tulis(w, h, err)
	})
	pasang("GET "+Prefix+"/kasus/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		k, err := l.BacaKasus(r.Context(), p, r.PathValue("id"))
		tulis(w, k, err)
	})
	pasang("GET "+Prefix+"/kasus/{id}/peserta", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		h, err := l.DaftarPeserta(r.Context(), p, r.PathValue("id"), angkaKueri(r, "halaman"), angkaKueri(r, "ukuran"))
		tulis(w, h, err)
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
//
// ⛔ Setiap galat sampai ke layar berkata-kata; galat tak terduga dicatat di
// log server (sebab aslinya) dan dijawab 500 berkalimat umum.
func jawabGalat(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "insufficient permission")
	case errors.Is(err, services.ErrKasusTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	default:
		log.Printf("endorsement life: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the endorsement life request")
	}
	return true
}
