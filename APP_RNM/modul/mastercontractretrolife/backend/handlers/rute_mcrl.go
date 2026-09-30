// Package handlers adalah pintu HTTP modul Master Contract Retro Life.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON.
// Nol aturan dagang di sini; nol impor repository (penjaga
// `TestHandlersTidakMengimporRepository`).
//
//	GET  /api/master-contract-retro-life/tahun                    grid tahun treaty (halaman awal)
//	GET  /api/master-contract-retro-life/tahun/{id}/kontrak       grid kontrak - tombol `ReinsType`
//	GET  /api/master-contract-retro-life/kontrak/{id}/reinsurer   grid reinsurer - tombol `Reinsurer List`
//	GET  /api/master-contract-retro-life/reinsurer/{id}/security  grid security - tombol `Security Reinsurer`
//	GET  /api/master-contract-retro-life/kontrak/{id}/business    grid business - tombol `Business List`
//	GET  /api/master-contract-retro-life/jenis-reasuransi         dropdown `REINS TYPE`
//	GET  /api/master-contract-retro-life/master-reinsurer?cari=   autocomplete `REINSURER NAME`
//	GET  /api/master-contract-retro-life/master-business?cari=    autocomplete `BUSINESS NAME`
//	GET  /api/master-contract-retro-life/ringkasan-rate?cari=     autocomplete `R/I RATE`
//	GET  /api/master-contract-retro-life/rate?idusedby=           section `ViewRate` (`Rate List`)
//	GET  /api/master-contract-retro-life/laporan/total-share-bukan-100?tahun=  tiket 11 (tanpa layar)
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/mastercontractretrolife/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/master-contract-retro-life"

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
}

func daftarkanBaca(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/tahun", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarTahun(r.Context(), p)
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/tahun/{id}/kontrak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		j, err := l.DaftarKontrak(r.Context(), p, r.PathValue("id"))
		tulis(w, j, err)
	})
	pasang("GET "+Prefix+"/kontrak/{id}/reinsurer", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		j, err := l.DaftarReinsurer(r.Context(), p, r.PathValue("id"))
		tulis(w, j, err)
	})
	pasang("GET "+Prefix+"/reinsurer/{id}/security", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		j, err := l.DaftarSecurity(r.Context(), p, r.PathValue("id"))
		tulis(w, j, err)
	})
	pasang("GET "+Prefix+"/kontrak/{id}/business", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		j, err := l.DaftarBusiness(r.Context(), p, r.PathValue("id"))
		tulis(w, j, err)
	})
	pasang("GET "+Prefix+"/jenis-reasuransi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.JenisReasuransi(r.Context(), p)
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/master-reinsurer", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.CariMasterReinsurer(r.Context(), p, r.URL.Query().Get("cari"))
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/master-business", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.CariMasterBusiness(r.Context(), p, r.URL.Query().Get("cari"))
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/ringkasan-rate", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.CariRingkasanRate(r.Context(), p, r.URL.Query().Get("cari"))
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/rate", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarRate(r.Context(), p, r.URL.Query().Get("idusedby"))
		tulisDaftar(w, d, err)
	})
	// Tiket 11 - kemampuan BARU tanpa padanan Pega: rute baca saja, nol layar (OQ-MCRL-07).
	pasang("GET "+Prefix+"/laporan/total-share-bukan-100", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.LaporanTotalShareBukan100(r.Context(), p, r.URL.Query().Get("tahun"))
		tulisDaftar(w, d, err)
	})
}

// jawabanDaftar adalah badan jawaban daftar sederhana.
type jawabanDaftar[T any] struct {
	Daftar []T `json:"daftar"`
	Total  int `json:"total"`
}

func tulisDaftar[T any](w http.ResponseWriter, d []T, err error) {
	if jawabGalat(w, err) {
		return
	}
	if d == nil {
		d = []T{}
	}
	galat.TulisJSON(w, jawabanDaftar[T]{Daftar: d, Total: len(d)})
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
// ⛔ K8: setiap galat sampai ke layar berkata-kata; galat tak terduga dicatat
// di log server (sebab aslinya) dan dijawab 500 berkalimat umum.
func jawabGalat(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "insufficient permission")
	case errors.Is(err, services.ErrTahunTidakAda), errors.Is(err, services.ErrKontrakTidakAda),
		errors.Is(err, services.ErrReinsurerTidakAda), errors.Is(err, services.ErrSecurityTidakAda),
		errors.Is(err, services.ErrBusinessTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	case errors.Is(err, services.ErrParameterWajib), errors.Is(err, services.ErrIDDariKlien):
		galat.Tulis(w, http.StatusBadRequest, services.Pesan(err))
	case errors.Is(err, services.ErrWajibIsi), errors.Is(err, services.ErrMasukanTidakSah):
		// 422: JSON-nya sah, isinya ditolak gerbang - pesan VERBATIM korpus
		// untuk wajib-isi, pesan yang menyebut medannya untuk nilai tak sah.
		// Pesan wajib-isi korpus tidak menyebut medannya; log server menyebutnya.
		var kosong services.GalatWajibIsi
		if errors.As(err, &kosong) {
			log.Printf("master contract retro life: required values empty: %s", strings.Join(kosong.Medan, ", "))
		}
		galat.Tulis(w, http.StatusUnprocessableEntity, services.Pesan(err))
	case errors.Is(err, services.ErrDampakBerubah):
		// 409: keadaan DATA berubah sejak pratinjau/popup - nol baris disentuh.
		galat.Tulis(w, http.StatusConflict, services.Pesan(err))
	case errors.Is(err, services.ErrRateMenungguPersetujuan):
		// 503 berkalimat: sumber tabel rate menunggu persetujuan (OQ-MCRL-13).
		galat.Tulis(w, http.StatusServiceUnavailable, services.Pesan(err))
	case errors.Is(err, services.ErrMasterTidakTerbaca):
		// 503: keadaan server - master rujukan tidak terbaca atau kosong -
		// dan pesannya MENYEBUT objeknya (ADR-0015).
		log.Printf("master contract retro life: %v", err)
		galat.Tulis(w, http.StatusServiceUnavailable, services.Pesan(err))
	default:
		log.Printf("master contract retro life: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the master contract retro life request")
	}
	return true
}
