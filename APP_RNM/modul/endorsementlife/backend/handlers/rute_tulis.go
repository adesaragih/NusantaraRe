package handlers

// Rute tulis modul Endorsement Life.
//
// ⛔ Identitas kasus baru tidak pernah datang dari klien (ADR-0006): `POST
// /kasus` menjawab pengenal `EDMLF-<n>` yang diterbitkan server.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
)

// batasBadan - badan JSON terbesar yang dibaca (64 KiB).
const batasBadan = 64 << 10

// batasBadanSimpan - `Save` membawa daftar ID peserta tercentang (32 heksa
// per ID); 4 MiB ≈ 100.000 ID.
const batasBadanSimpan = 4 << 20

// bacaBadan mengurai badan JSON SESUDAH identitas diperiksa (401 lebih dulu);
// terlalu besar = 413, rusak = 400. Mengembalikan true bila SUDAH dijawab.
func bacaBadan(w http.ResponseWriter, r *http.Request, p inti.Pelaku, ke any, batas int64) bool {
	if jawabGalat(w, inti.WajibIdentitas(p)) {
		return true
	}
	err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batas)).Decode(ke)
	var besar *http.MaxBytesError
	switch {
	case err == nil:
		return false
	case errors.As(err, &besar):
		galat.Tulis(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("request body is larger than %d KiB", batas>>10))
	default:
		galat.Tulis(w, http.StatusBadRequest, "request body is not valid JSON")
	}
	return true
}

func daftarkanTulis(pasang func(string, rute)) {
	pasang("POST "+Prefix+"/kelayakan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanKelayakan
		if bacaBadan(w, r, p, &m, batasBadan) {
			return
		}
		k, err := l.Kelayakan(r.Context(), p, m)
		tulis(w, k, err)
	})
	pasang("POST "+Prefix+"/kasus", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanKasus
		if bacaBadan(w, r, p, &m, batasBadan) {
			return
		}
		h, err := l.BuatKasus(r.Context(), p, m)
		if jawabGalat(w, err) {
			return
		}
		w.Header().Set("Location", Prefix+"/kasus/"+h.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(h)
	})
	pasang("POST "+Prefix+"/kasus/{id}/simpan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m models.PilihanHapus
		if bacaBadan(w, r, p, &m, batasBadanSimpan) {
			return
		}
		h, err := l.Simpan(r.Context(), p, r.PathValue("id"), m)
		tulis(w, h, err)
	})
	pasang("POST "+Prefix+"/kasus/{id}/putuskan", func(w http.ResponseWriter, r *http.Request, l *services.Layanan, p inti.Pelaku) {
		var m services.MasukanPutusan
		if bacaBadan(w, r, p, &m, batasBadan) {
			return
		}
		h, err := l.Putuskan(r.Context(), p, r.PathValue("id"), m)
		tulis(w, h, err)
	})
}
