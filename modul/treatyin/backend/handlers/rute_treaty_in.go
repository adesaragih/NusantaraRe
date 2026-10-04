// Package handlers memegang rute HTTP modul Treaty In (`treatyin`).
//
// Arah ketergantungan: handlers -> services -> repository. Handler tidak
// pernah memegang koneksi.
package handlers

import (
	"errors"
	"log"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// Prefix rute API modul ini (`MODUL.md`).
const Prefix = "/api/treaty-in"

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
	daftarkanAcuan(pasang)
	daftarkanKontrak(pasang)
	daftarkanIdentitas(pasang)
	daftarkanWarisan(pasang)
}

// daftarkanAcuan - jalur baca keenam tabel acuan (tiket 15).
//
// ⛔ Ini SELURUH permukaan HTTP modul ini hari ini, dan itu disengaja. Papan
// tiket Treaty In menyatakan `L-4`: "tidak ada spesifikasi layar di mana pun",
// dan melarang mengarang layar untuk memenuhi bentuk irisan tegak. Rute
// kontrak, versi, dan persetujuan LAHIR BERSAMA spesifikasinya, bukan sebelum.
func daftarkanAcuan(pasang func(string, rute)) {
	pasang("GET "+Prefix+"/acuan/{himpunan}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		d, err := l.DaftarAcuan(r.Context(), p, models.Himpunan(r.PathValue("himpunan")))
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
func jawabGalat(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "request without user identity is rejected")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "insufficient permission")
	case errors.Is(err, services.ErrMasukanTidakSah):
		// 422: JSON-nya sah, isinya ditolak gerbang - INV-29, INV-53, ruas
		// wajib. Pesannya memuat SELURUH pelanggaran, bukan yang pertama.
		galat.Tulis(w, http.StatusUnprocessableEntity, services.Pesan(err))
	case errors.Is(err, services.ErrPemulihanTidakSah):
		// 422: JSON-nya sah, isinya ditolak gerbang tiket 32 - INV-41 dan
		// nomor urut kembar. Pesannya menyebut BARIS KE BERAPA dan nilainya,
		// sebab daftar pemulihan dikirim sekaligus dan "ditolak" saja
		// membuat pengirimnya memeriksa seluruh baris satu per satu.
		galat.Tulis(w, http.StatusUnprocessableEntity, services.Pesan(err))
	case errors.Is(err, services.ErrRuasBekuBerubah), errors.Is(err, services.ErrVersiMenyimpang):
		// 422: JSON-nya sah, isinya ditolak gerbang identitas. INV-19 dan
		// ADR-0040: pesannya menyebut RUAS mana, bukan sekadar "ditolak".
		galat.Tulis(w, http.StatusUnprocessableEntity, services.Pesan(err))
	case errors.Is(err, services.ErrNomorUrutVersiGanda):
		// 409: keadaan DATA menolak - INV-04. Pesannya menyebut nomornya.
		galat.Tulis(w, http.StatusConflict, services.Pesan(err))
	case services.WarisanTidakAda(err):
		// 404: pengenalnya tidak menunjuk kontrak warisan mana pun.
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	case services.WarisanJSONRusak(err):
		// 500, dan pesannya MENYEBUT KONTRAKNYA. Dokumen warisan yang tidak
		// dapat diurai bukan kesalahan pemanggil — ia cacat data yang harus
		// dapat ditemukan, dan "invalid character" tanpa pengenal membuat
		// yang menyelidiki memeriksa 1.854 dokumen.
		log.Printf("treaty in: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, services.Pesan(err))
	case errors.Is(err, services.ErrKontrakTidakAda):
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	case errors.Is(err, services.ErrHimpunanTidakAda):
		// 404: himpunan yang diminta bukan salah satu dari enam. Pesannya
		// menyebut yang diminta - penolakan yang tidak menyebut apa yang
		// ditolak membuat pemanggilnya menebak.
		galat.Tulis(w, http.StatusNotFound, services.Pesan(err))
	default:
		// Galat tak terduga: sebab aslinya ke log server, kalimat umum ke layar.
		log.Printf("treaty in: %v", err)
		galat.Tulis(w, http.StatusInternalServerError, "failed to process the treaty in request")
	}
	return true
}
