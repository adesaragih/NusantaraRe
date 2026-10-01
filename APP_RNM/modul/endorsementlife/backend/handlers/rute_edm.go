// Package handlers adalah pintu HTTP modul Endorsement Life.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON.
// Nol aturan dagang di sini; nol impor repository (penjaga
// `TestHandlersTidakMengimporRepository`).
//
//	GET  /api/endorsement-life/inbox?halaman=&ukuran=        kotak masuk `InboxEndorsementLife`
//	GET  /api/endorsement-life/kasus/{id}                    kepala `InputEDMLife`
//	GET  /api/endorsement-life/kasus/{id}/peserta?halaman=   grid peserta b11899 / b17500
//	GET  /api/endorsement-life/kasus/{id}/peserta/{pid}      rincian `PL_Detail_Sec` + `RetroDetailLife`
//	GET  /api/endorsement-life/kasus/{id}/polis-lama?halaman= popup `View Old Policy`
//	POST /api/endorsement-life/kelayakan                     `SetErrorBatalEndorsement_Act` (tanpa tulis)
//	POST /api/endorsement-life/kasus                         `Submit` b4226 → `MappingEDMLife`
//	POST /api/endorsement-life/kasus/{id}/simpan             `Save` b37202 → `SetPremi_EDM`
//	POST /api/endorsement-life/kasus/{id}/unggah             `Upload CSV` b8973 - tinjau, nol tulis (multipart `berkas`)
//	POST /api/endorsement-life/kasus/{id}/csv                `Add CSV Data` b10405 → `SaveCSVEDMLife`
//	POST /api/endorsement-life/kasus/{id}/putuskan           `Submit` b37494 / b38109 → `IsLifeAccepted` (Confirm/Decline)
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
	daftarkanTulis(pasang)
	daftarkanCSV(pasang)
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
	pasang("GET "+Prefix+"/kasus/{id}/peserta/{pid}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.RincianPeserta(r.Context(), p, r.PathValue("id"), r.PathValue("pid"))
		tulis(w, d, err)
	})
	pasang("GET "+Prefix+"/kasus/{id}/polis-lama", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.PolisLama(r.Context(), p, r.PathValue("id"), angkaKueri(r, "halaman"), angkaKueri(r, "ukuran"))
		tulis(w, d, err)
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
	case errors.Is(err, services.ErrKasusTidakAda), errors.Is(err, services.ErrPesertaTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	case errors.As(err, new(services.GalatKelayakan)), errors.Is(err, services.ErrMasukanTidakSah),
		errors.Is(err, services.ErrTanpaPeserta), errors.Is(err, services.ErrCSVKosong), errors.Is(err, services.ErrCSVRusak),
		errors.Is(err, services.ErrCSVTanpaAcuan), errors.Is(err, services.ErrCSVBukanPerubahanData):
		// 422: JSON-nya sah, isinya ditolak gerbang - pesan VERBATIM korpus, satu per baris.
		galat.Tulis(w, http.StatusUnprocessableEntity, services.Pesan(err))
	case errors.Is(err, services.ErrKasusTerbukaGanda), errors.Is(err, services.ErrSumberWarisanEDM),
		errors.Is(err, services.ErrKasusTertutup), errors.Is(err, services.ErrSudahDisimpan), errors.Is(err, services.ErrCSVTerkunci),
		errors.Is(err, services.ErrBelumDisimpan), errors.Is(err, services.ErrVersiBerubah):
		// 409: keadaan DATA menolak - kasus terbuka lain lahir bersamaan, atau
		// versi berjalan polis belum dapat disalin (OQ-EDM-016).
		galat.Tulis(w, http.StatusConflict, services.Pesan(err))
	case errors.Is(err, services.ErrArasapas):
		log.Printf("endorsement life: %v", err)
		galat.Tulis(w, http.StatusServiceUnavailable, services.Pesan(err))
	default:
		log.Printf("endorsement life: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the endorsement life request")
	}
	return true
}
