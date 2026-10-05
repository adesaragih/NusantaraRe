// Package handlers adalah pintu HTTP modul Aggregate.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON. Nol aturan dagang di sini; nol impor
// repository. Rute modul ini hanya terbuka bagi pemegang menu `aggregate` (perakit `cmd/api`).
//
//	GET  /api/aggregate?q=&halaman=&ukuran=   daftar GridDasbordAgg (satu baris per kunci)
//	GET  /api/aggregate/ringkasan             chart RNM Value (USD) per Ceding, Treaty Type, Coverage, seluruh As At
//	GET  /api/aggregate/rincian?<kunci>       seluruh baris satu kunci (klik ganda)
//	POST /api/aggregate/hapus                 Delete satu kunci {kunci}
//	GET  /api/aggregate/master-treaty?q=      popup Master ID
//	POST /api/aggregate/pratinjau             Upload CSV {csv, masterTreaty}
//	POST /api/aggregate/simpan                Save {baris}
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
	"nusantarare/modul/aggregate/backend/models"
	"nusantarare/modul/aggregate/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/aggregate"

// Batas badan permintaan: berkas CSV (5.000 baris) dan grid pratinjau yang dikirim balik saat Save.
const (
	batasBadanKecil = 1 << 16
	batasPratinjau  = 8 << 20
	batasSimpan     = 16 << 20
)

// Router menyusun rute modul ini SAJA di atas Oracle.
func Router(svc *services.Service, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stubPelaku)
	return mux
}

// RouterDengan menyusun rute di atas Layanan yang sudah disusun (uji dengan gudang tiruan). `adaDB` false = setiap
// rute menjawab 503.
func RouterDengan(l *services.Layanan, adaDB, stubPelaku bool) http.Handler {
	mux := http.NewServeMux()
	daftarkan(mux, func() *services.Layanan { return l }, func() bool { return adaDB }, stubPelaku)
	return mux
}

// DaftarkanRute mendaftarkan seluruh rute modul ini ke mux bersama.
func DaftarkanRute(mux *http.ServeMux, svc *services.Service, stubPelaku bool) {
	daftarkan(mux, func() *services.Layanan { return services.LayananOracle(svc) }, svc.PunyaDatabase, stubPelaku)
}

type rute func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku)

func angka(r *http.Request, nama string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get(nama)))
	return n
}

// kunciKueri - kunci daftar dari parameter kueri.
func kunciKueri(r *http.Request) models.Kunci {
	q := r.URL.Query()
	return models.Kunci{TanggalInput: q.Get("tanggalInput"), CedingCode: q.Get("cedingCode"),
		CedingName: q.Get("cedingName"), TreatyType: q.Get("treatyType"), AsAt: q.Get("asAt"), UwYear: q.Get("uwYear")}
}

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
	pasang("GET "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		h, err := l.Daftar(r.Context(), r.URL.Query().Get("q"), angka(r, "halaman"), angka(r, "ukuran"))
		if jawabGalat(w, err, "membaca daftar") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("GET "+Prefix+"/ringkasan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		h, err := l.Ringkasan(r.Context())
		if jawabGalat(w, err, "membaca ringkasan") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("GET "+Prefix+"/rincian", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		b, err := l.Rincian(r.Context(), kunciKueri(r))
		if jawabGalat(w, err, "membaca rincian") {
			return
		}
		galat.TulisJSON(w, jawabanBaris{Baris: b})
	})
	pasang("POST "+Prefix+"/hapus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var k models.Kunci
		if !bacaBadan(w, r, &k, batasBadanKecil) {
			return
		}
		n, err := l.Hapus(r.Context(), p, k)
		if jawabGalat(w, err, "menghapus") {
			return
		}
		galat.TulisJSON(w, jawabanHapus{Dihapus: n})
	})
	pasang("GET "+Prefix+"/master-treaty", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		d, err := l.CariTreaty(r.Context(), r.URL.Query().Get("q"))
		if jawabGalat(w, err, "mencari master treaty") {
			return
		}
		galat.TulisJSON(w, jawabanTreaty{Daftar: d})
	})
	pasang("POST "+Prefix+"/pratinjau", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		var req services.PermintaanPratinjau
		if !bacaBadan(w, r, &req, batasPratinjau) {
			return
		}
		h, err := l.Pratinjau(r.Context(), req)
		if jawabGalat(w, err, "membaca CSV") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("POST "+Prefix+"/simpan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var req services.PermintaanSimpan
		if !bacaBadan(w, r, &req, batasSimpan) {
			return
		}
		h, err := l.Simpan(r.Context(), p, req)
		if jawabGalat(w, err, "menyimpan") {
			return
		}
		galat.TulisJSON(w, h)
	})
}

type jawabanBaris struct {
	Baris []models.Baris `json:"baris"`
}

type jawabanHapus struct {
	Dihapus int64 `json:"dihapus"`
}

type jawabanTreaty struct {
	Daftar []models.MasterTreaty `json:"daftar"`
}

// bacaBadan membaca badan JSON; galat = 400.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any, batas int64) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, batas))
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

// pesan - kalimat untuk layar: bagian sesudah nama galat penanda.
func pesan(err, penanda error) string {
	return strings.TrimPrefix(err.Error(), penanda.Error()+": ")
}

// jawabGalat menulis galat; true = sudah dijawab.
func jawabGalat(w http.ResponseWriter, err error, apa string) bool {
	var simpan services.GalatSimpan
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaPelaku):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.As(err, &simpan):
		galat.Tulis(w, http.StatusUnprocessableEntity, simpan.Error())
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Aggregate data not found")
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrBelumDimigrasi):
		galat.Tulis(w, http.StatusServiceUnavailable, "Aggregate is not migrated yet (migration 880)")
	default:
		log.Printf("aggregate: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Aggregate: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
