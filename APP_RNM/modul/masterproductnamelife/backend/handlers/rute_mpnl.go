// Package handlers adalah pintu HTTP modul Master Product Name Life.
//
// Tugasnya sempit: baca permintaan, panggil services, tulis jawaban JSON.
// Nol aturan dagang di sini; nol impor repository.
//
//	GET  /api/master-product-name-life/produk        grid `InboxProductName` (halaman awal)
//	GET  /api/master-product-name-life/produk/{id}   tombol `View` b74798
//	GET  /api/master-product-name-life/master/{jenis}?cari=&batas=  tujuh pemilih master (`Choose*`, PARITAS §4); `batas` = autocomplete
//	GET  /api/master-product-name-life/master-plan?cari=     autocomplete `Plan Name` (PLAN LIST)
//	GET  /api/master-product-name-life/rate?riRateId=        tombol `View Rate` - view `RATE_LIFE` (K1)
//	GET  /api/master-product-name-life/produk-lama           popup `Copy Old` (`rute_lama.go`, 03-10-2026)
//	POST /api/master-product-name-life/produk-lama/salin     `Process Copy`
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/master-product-name-life"

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
	daftarkanLampiran(pasang)
	daftarkanLama(pasang)
}

func daftarkanBaca(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/produk", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarProduk(r.Context(), p)
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/produk/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		pr, err := l.AmbilProduk(r.Context(), p, r.PathValue("id"))
		tulis(w, pr, err)
	})
	// Tujuh pemilih master (paket 2): `Choose*` → section → grid RD; juga autocomplete medan form.
	pasang("GET "+Prefix+"/master/{jenis}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		batas, ok := angkaKueri(w, r, "batas")
		if !ok {
			return
		}
		d, err := l.CariMaster(r.Context(), p, models.JenisMaster(r.PathValue("jenis")), r.URL.Query().Get("cari"), batas)
		tulisDaftar(w, d, err)
	})
	// Grid `PLAN LIST` (paket 6): autocomplete `Plan Name` dan tombol `View Rate`.
	pasang("GET "+Prefix+"/master-plan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.CariPlan(r.Context(), p, r.URL.Query().Get("cari"))
		tulisDaftar(w, d, err)
	})
	pasang("GET "+Prefix+"/rate", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarRate(r.Context(), p, r.URL.Query().Get("riRateId"))
		tulis(w, d, err)
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

// angkaKueri - parameter kueri bilangan bulat tak-negatif; kosong = 0; selain itu 400 dan false.
func angkaKueri(w http.ResponseWriter, r *http.Request, nama string) (int, bool) {
	t := r.URL.Query().Get(nama)
	if t == "" {
		return 0, true
	}
	n, err := strconv.Atoi(t)
	if err != nil || n < 0 {
		galat.Tulis(w, http.StatusBadRequest, nama+" must be a non-negative whole number")
		return 0, false
	}
	return n, true
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
	case errors.Is(err, services.ErrProdukTidakAda), errors.Is(err, services.ErrJenisMasterTidakDikenal),
		errors.Is(err, services.ErrSalinanTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	case errors.Is(err, services.ErrLampiranTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	case errors.Is(err, services.ErrNamaLampiranSudahAda), errors.Is(err, services.ErrLampiranBelumTerkirim),
		errors.Is(err, services.ErrLampiranSudahTerkirim), errors.Is(err, services.ErrBerkasSumberHilang),
		errors.Is(err, services.ErrBerkasTidakDiStub), errors.Is(err, services.ErrBerkasTidakDiStorage):
		// 409: keadaan lampiran menolak aksi - kalimat menyebut yang harus dilakukan.
		galat.Tulis(w, http.StatusConflict, services.Pesan(err))
	case errors.Is(err, services.ErrOfficeStub), errors.Is(err, services.ErrPenyimpananBelumDisetel),
		errors.Is(err, services.ErrStorageBelumSiap):
		// 503: berkas stub lokal tanpa URL untuk penampil kantor / folder antrean belum disetel / alamat, App, atau
		// token penyimpanan tidak tersedia.
		galat.Tulis(w, http.StatusServiceUnavailable, services.Pesan(err))
	case errors.Is(err, services.ErrStorageGagal):
		// 502: layanan penyimpanan nyata gagal - rinciannya di log layanan, tanpa alamat dan token.
		galat.Tulis(w, http.StatusBadGateway, services.Pesan(err))
	case errors.Is(err, services.ErrMasterTidakTerbaca):
		// 503: master rujukan tidak terbaca - pesannya MENYEBUT objeknya, sebab Oracle hanya di log.
		log.Printf("master product name life: %v", err)
		galat.Tulis(w, http.StatusServiceUnavailable, services.Pesan(err))
	case errors.Is(err, services.ErrIDDariKlien), errors.Is(err, services.ErrSalinanPadaUbah):
		galat.Tulis(w, http.StatusBadRequest, services.Pesan(err))
	case errors.Is(err, services.ErrMasukanTidakSah), errors.Is(err, services.ErrNilaiTidakMuat):
		// 422: JSON-nya sah, isinya ditolak - kalimat menyebut label medan VERBATIM; semua penolakan sekaligus.
		galat.Tulis(w, http.StatusUnprocessableEntity, services.Pesan(err))
	case errors.Is(err, services.ErrIdentitasGanda), errors.Is(err, services.ErrJSONRusak),
		errors.Is(err, services.ErrIdentitasMelampauiLebar), errors.Is(err, services.ErrIdentitasBentrok):
		// 500 berkalimat: keadaan DATA yang harus diperbaiki DBA, disebut terang.
		log.Printf("master product name life: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, services.Pesan(err))
	default:
		log.Printf("master product name life: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the master product name life request")
	}
	return true
}
