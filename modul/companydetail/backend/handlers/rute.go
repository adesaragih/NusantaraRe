// Package handlers adalah pintu HTTP modul Company Detail.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON. Nol aturan dagang di sini; nol impor
// repository. Rute modul ini hanya terbuka bagi pemegang menu `companydetail` (perakit `cmd/api`).
//
//	GET  /api/company-detail?q=&halaman=&ukuran=  satu halaman organisasi (CLIENT FLAG Org)
//	GET  /api/company-detail/pilihan              pilihan form: M_ENUMERASI dan NATION
//	GET  /api/company-detail/induk?q=&kecuali=    pilihan Parent organization
//	GET  /api/company-detail/periksa-nama?nama=&kecuali=  title di nama + organisasi bernama sama/mirip
//	GET  /api/company-detail/hak                  tombol yang boleh tampil (Copy Old: superadmin)
//	GET  /api/company-detail/lama                 popup Copy Old - superadmin
//	POST /api/company-detail/lama/salin           Process Copy {"ids": [...]} - superadmin
//	GET  /api/company-detail/{id}                 satu organisasi + PIC + alamat
//	POST /api/company-detail                      Create (nomor ORG dari SEQ_CLIENT_ORG)
//	PUT  /api/company-detail/{id}                 ubah - nol hapus organisasi
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
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/company-detail"

// batasBadan - badan JSON form: organisasi dengan puluhan PIC dan alamat jauh di bawahnya.
const batasBadan = 1 << 20

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
	pasang("GET "+Prefix+"/pilihan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		p, err := l.Pilihan(r.Context())
		if jawabGalat(w, err, "membaca pilihan") {
			return
		}
		galat.TulisJSON(w, p)
	})
	pasang("GET "+Prefix+"/induk", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		d, err := l.CariInduk(r.Context(), r.URL.Query().Get("q"), r.URL.Query().Get("kecuali"))
		if jawabGalat(w, err, "mencari parent organization") {
			return
		}
		galat.TulisJSON(w, jawabanInduk{Daftar: d})
	})
	pasang("GET "+Prefix+"/periksa-nama", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		h, err := l.Periksa(r.Context(), r.URL.Query().Get("nama"), r.URL.Query().Get("kecuali"))
		if jawabGalat(w, err, "memeriksa nama") {
			return
		}
		galat.TulisJSON(w, h)
	})
	pasang("GET "+Prefix+"/hak", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		galat.TulisJSON(w, l.HakAkun(r.Context()))
	})
	pasang("GET "+Prefix+"/lama", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarLama(r.Context(), p)
		if jawabGalat(w, err, "membaca organisasi lama") {
			return
		}
		galat.TulisJSON(w, jawabanLama{Daftar: d})
	})
	pasang("POST "+Prefix+"/lama/salin", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var b struct {
			IDs []string `json:"ids"`
		}
		if !bacaBadan(w, r, &b) {
			return
		}
		j, err := l.SalinLama(r.Context(), p, b.IDs)
		if jawabGalat(w, err, "Copy Old") {
			return
		}
		galat.TulisJSON(w, j)
	})
	pasang("GET "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, _ inti.Pelaku) {
		d, err := l.Ambil(r.Context(), r.PathValue("id"))
		if jawabGalat(w, err, "membaca organisasi") {
			return
		}
		galat.TulisJSON(w, d)
	})
	pasang("POST "+Prefix, func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi) {
			return
		}
		d, err := l.Tambah(r.Context(), p, isi)
		if jawabGalat(w, err, "membuat organisasi") {
			return
		}
		galat.TulisJSON(w, d)
	})
	pasang("PUT "+Prefix+"/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var isi models.Isian
		if !bacaBadan(w, r, &isi) {
			return
		}
		d, err := l.Ubah(r.Context(), p, r.PathValue("id"), isi)
		if jawabGalat(w, err, "mengubah organisasi") {
			return
		}
		galat.TulisJSON(w, d)
	})
}

// jawabanLama adalah badan jawaban popup Copy Old.
type jawabanLama struct {
	Daftar []models.OrgLama `json:"daftar"`
}

// jawabanInduk adalah badan jawaban pilihan Parent organization.
type jawabanInduk struct {
	Daftar []models.BarisDaftar `json:"daftar"`
}

// bacaBadan membaca badan JSON; galat = 400.
func bacaBadan(w http.ResponseWriter, r *http.Request, ke any) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadan))
	d.DisallowUnknownFields()
	if err := d.Decode(ke); err != nil {
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
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrTanpaPelaku):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, services.ErrBukanSuperadmin):
		galat.Tulis(w, http.StatusForbidden, "Copy Old is only for super admin (Kelola User holder)")
	case errors.Is(err, services.ErrTidakAda):
		galat.Tulis(w, http.StatusNotFound, "Organization not found")
	case errors.Is(err, services.ErrMasukanTidakSah):
		galat.Tulis(w, http.StatusUnprocessableEntity, pesan(err, services.ErrMasukanTidakSah))
	case errors.Is(err, services.ErrBelumDimigrasi):
		galat.Tulis(w, http.StatusServiceUnavailable, "Company Detail tables are not migrated yet (migrations 800-810)")
	default:
		log.Printf("company detail: %s: %v", apa, err)
		galat.Tulis(w, http.StatusInternalServerError, "Company Detail: "+apa+" gagal; rinciannya di log server")
	}
	return true
}
