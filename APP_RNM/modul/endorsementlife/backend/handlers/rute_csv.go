package handlers

// Rute unggah CSV (tiket 07): `Upload CSV` b8973 → `POST /kasus/{id}/unggah`
// (tinjau, nol tulis) dan `Add CSV Data` b10405 → `POST /kasus/{id}/csv`.
// Berkasnya bagian multipart `berkas` (pola unggahan PremiumList).

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/endorsementlife/backend/services"
)

// batasBerkasCSV - badan unggahan terbesar (256 MiB).
//
// ⛔ BUKAN batas baris (AC 36): ±600.000 peserta pada lebar 38 kolom. Ia
// membatasi badan yang tak berujung, bukan polis grup yang besar.
const batasBerkasCSV = 256 << 20

// bagianBerkas membuka bagian `berkas` SESUDAH identitas diperiksa; false bila
// permintaan sudah dijawab.
func bagianBerkas(w http.ResponseWriter, r *http.Request, p inti.Pelaku) (io.Reader, bool) {
	if jawabGalat(w, inti.WajibIdentitas(p)) {
		return nil, false
	}
	r.Body = http.MaxBytesReader(w, r.Body, batasBerkasCSV)
	mr, err := r.MultipartReader()
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "request body must be multipart/form-data with a file field \"berkas\"")
		return nil, false
	}
	for {
		part, err := mr.NextPart()
		if err != nil {
			jawabBadan(w, err, "multipart field \"berkas\" is missing")
			return nil, false
		}
		if part.FormName() == "berkas" {
			return part, true
		}
	}
}

// jawabBadan - 413 bila badan melewati batas, selain itu 400.
func jawabBadan(w http.ResponseWriter, err error, pesan string) {
	var besar *http.MaxBytesError
	if errors.As(err, &besar) {
		galat.Tulis(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("uploaded file is larger than %d MiB", batasBerkasCSV>>20))
		return
	}
	galat.Tulis(w, http.StatusBadRequest, pesan)
}

// jawabCSV - penolakan baris berbadan daftar pesannya (422); galat lain lewat jawabGalat.
func jawabCSV(w http.ResponseWriter, isi any, err error) {
	var gc services.GalatCSV
	if errors.As(err, &gc) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(struct {
			Galat string `json:"galat"`
			services.HasilPeriksaCSV
		}{services.Pesan(gc), gc.HasilPeriksaCSV})
		return
	}
	tulis(w, isi, err)
}

func daftarkanCSV(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/kasus/{id}/unggah", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		berkas, ok := bagianBerkas(w, r, p)
		if !ok {
			return
		}
		h, err := l.PeriksaCSV(r.Context(), p, r.PathValue("id"), berkas)
		var besar *http.MaxBytesError
		if errors.As(err, &besar) {
			jawabBadan(w, err, "")
			return
		}
		jawabCSV(w, h, err)
	})
	pasang("POST "+Prefix+"/kasus/{id}/csv", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		berkas, ok := bagianBerkas(w, r, p)
		if !ok {
			return
		}
		// Ditampung ke berkas sementara: dua lintasan (periksa, lalu sisip di
		// dalam transaksi) tanpa memuat seluruh berkas ke memori.
		f, err := os.CreateTemp("", "edm-csv-*.csv")
		if err != nil {
			log.Printf("endorsement life: temp file: %v", err)
			galat.Tulis(w, http.StatusInternalServerError, "failed to receive the uploaded file")
			return
		}
		defer func() {
			_ = f.Close()
			_ = os.Remove(f.Name())
		}()
		if _, err := io.Copy(f, berkas); err != nil {
			jawabBadan(w, err, "failed to read the uploaded file")
			return
		}
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			galat.Tulis(w, http.StatusInternalServerError, "failed to receive the uploaded file")
			return
		}
		h, err := l.TambahCSV(r.Context(), p, r.PathValue("id"), f)
		jawabCSV(w, h, err)
	})
}
