package handlers

// Rute identitas kontrak — tiket 16, 17, 18, 19.
//
// Tiket 16 tidak punya rutenya sendiri: peringatan kunci alami ganda menumpang
// jawaban `POST /kontrak`, sebab penyimpanannya BERHASIL — ia bukan galat.

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatyin/backend/services"
)

func daftarkanIdentitas(pasang func(string, rute)) {
	// Tiket 17 — pencarian lewat nomor sistem lama.
	pasang("GET "+Prefix+"/kontrak/cari", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		k, err := l.CariLewatNomorWarisan(r.Context(), p, r.URL.Query().Get("nomorWarisan"))
		tulis(w, k, err)
	})

	// Tiket 18 — menyunting kepala kontrak; kelima ruas beku ditolak.
	pasang("PUT "+Prefix+"/kontrak/{id}", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenal(w, r)
		if !ok {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var m services.MasukanUbahKontrak
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		if jawabGalat(w, l.PerbaruiKontrak(r.Context(), p, id, m)) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// Tiket 19 — versi berikutnya; lapisan beku yang menyimpang ditolak.
	pasang("POST "+Prefix+"/kontrak/{id}/versi", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		id, ok := pengenal(w, r)
		if !ok {
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, batasBadanJSON)
		var m services.MasukanVersiTambahan
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "request body must be valid JSON")
			return
		}
		idVersi, err := l.TambahVersi(r.Context(), p, id, m)
		if jawabGalat(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			IDVersi int64 `json:"idVersi"`
		}{idVersi})
	})
}

// pengenal membaca `{id}` sebagai bilangan; false berarti permintaan SUDAH
// dijawab 400.
//
// Pengenal yang bukan angka adalah BENTUK permintaan yang salah, bukan kontrak
// yang tidak ada — dua pertanyaan, dua jawaban.
func pengenal(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		galat.Tulis(w, http.StatusBadRequest, "contract id must be a number")
		return 0, false
	}
	return id, true
}
