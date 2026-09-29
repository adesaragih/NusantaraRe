package handlers

// Pintu HTTP unggahan CSV premium list - tiket 04 PremiumList Life.
//
//	POST /api/polis-life/{id}/unggah/tinjau   (multipart, bagian `berkas`)
//	POST /api/polis-life/{id}/unggah/simpan   (multipart, bagian `berkas`)
//
// ⛔ DUA RUTE, dan keduanya POST. Tinjauan pun POST sebab ia MENGIRIM berkas;
// GET tidak punya badan permintaan. Yang membedakan keduanya bukan metodenya
// melainkan akibatnya, dan akibat itu dinyatakan di namanya: `tinjau` tidak
// menyentuh apa pun, `simpan` mengganti isi.
//
// ⛔ TINJAU TIDAK MENYIMPAN APA PUN. Rute tunggal yang "menyimpan bila lolos"
// membuat orang kehilangan kesempatan melihat hasilnya lebih dahulu - yaitu
// tepat yang AC tiket ini minta.
//
// Nol aturan dagang di sini; seluruhnya di `services/polis_unggah.go` dan
// `models/polis_unggah.go`.

import (
	"encoding/json"
	"errors"
	"net/http"

	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/galat"
)

// ambilBerkasCSV mengambil bagian `berkas` dari sebuah multipart.
//
// Mengembalikan true bila permintaan SUDAH dijawab.
func ambilBerkasCSV(w http.ResponseWriter, r *http.Request) (multipartBerkas, bool) {
	if err := r.ParseMultipartForm(galat.BatasFormulir); err != nil {
		galat.Tulis(w, http.StatusBadRequest, "permintaan bukan multipart yang sah")
		return multipartBerkas{}, true
	}
	berkas, kepala, err := r.FormFile("berkas")
	if err != nil {
		// ⛔ Pesannya MENYEBUT nama bagiannya. "Berkas tidak ada" membuat
		// orang menebak apakah namanya `file`, `csv`, atau `berkas`.
		galat.Tulis(w, http.StatusBadRequest, "bagian `berkas` tidak ada di permintaan")
		return multipartBerkas{}, true
	}
	return multipartBerkas{isi: berkas, nama: kepala.Filename}, false
}

type multipartBerkas struct {
	isi  interface{ Read([]byte) (int, error) }
	nama string
}

// tinjauUnggahPolis melayani POST .../unggah/tinjau.
func tinjauUnggahPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, sudah := ambilBerkasCSV(w, r)
		if sudah {
			return
		}
		hasil, err := svc.UnggahPremiumList().Tinjau(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.PathValue("id"), b.isi)
		if jawabGalatUnggahPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// simpanUnggahPolis melayani POST .../unggah/simpan.
func simpanUnggahPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		b, sudah := ambilBerkasCSV(w, r)
		if sudah {
			return
		}
		hasil, err := svc.UnggahPremiumList().Simpan(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.PathValue("id"), b.isi)
		// ⛔ PENOLAKAN IKUT DIKIRIM, bukan hanya kodenya. 409 dengan badan
		// kosong menyuruh orang mengunggah ulang ke rute tinjau untuk
		// mengetahui apa yang salah - dua putaran untuk satu jawaban.
		if errors.Is(err, services.ErrUnggahanDitolak) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(hasil)
			return
		}
		if jawabGalatUnggahPolis(w, err) {
			return
		}
		galat.TulisJSON(w, hasil)
	}
}

// jawabGalatUnggahPolis menerjemahkan galat unggahan menjadi kode HTTP.
func jawabGalatUnggahPolis(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, services.ErrCSVKosong),
		errors.Is(err, services.ErrCSVJudulGanda),
		errors.Is(err, services.ErrCSVTerlaluBanyakBaris),
		errors.Is(err, services.ErrCSVKolomKurang):
		// 400: berkasnya yang tidak dapat dibaca sama sekali - berbeda dari
		// berkas yang terbaca tetapi isinya ditolak (409).
		galat.Tulis(w, http.StatusBadRequest, err.Error())
		return true
	}
	return jawabGalatPolis(w, err)
}
