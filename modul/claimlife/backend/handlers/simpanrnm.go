package handlers

// Pintu HTTP `Save to RNM` - tiket 03, GILIRAN-11 paket 1.
//
//	POST /api/klaim-life/{id}/outstanding
//
// Tombol `Save to RNM` `Section/InputOSClaimLife.xml` b21102 ->
// `SaveOutStandingLife_Act` b21126. Nol aturan dagang di sini.
//
// Dibaca sesudah: services/simpanrnm.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// jawabGalatSimpanRNM menerjemahkan galat Save to RNM. true = sudah dijawab.
func jawabGalatSimpanRNM(w http.ResponseWriter, err error) bool {
	var tolak *services.PelanggaranRNM
	switch {
	case err == nil:
		return false
	case errors.As(err, &tolak):
		// ⭐ Kalimat XML apa adanya, dan langkah asalnya di medan sendiri.
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		// 11-10-2026: seluruh kekurangan kolom wajib ikut sebagai daftar terstruktur (layar mengelompokkan per peserta).
		kurang := tolak.Kekurangan
		if kurang == nil {
			kurang = []services.KekuranganRNM{}
		}
		_ = json.NewEncoder(w).Encode(struct {
			Galat      string                   `json:"galat"`
			Langkah    string                   `json:"langkah"`
			Kekurangan []services.KekuranganRNM `json:"kekurangan"`
		}{tolak.Pesan, tolak.Langkah, kurang})
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "Only the Claim Admin can Save to RNM.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrSimpanRNMBukanOutstanding):
		galat.Tulis(w, http.StatusConflict, "Save to RNM is only available at the Claim Admin stage.")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "The claim does not exist.")
	case errors.Is(err, services.ErrTypeTidakDikenal),
		errors.Is(err, services.ErrBusinessCodeTidakDikenal),
		errors.Is(err, kontrak.ErrPolisNomorTakDitemukan):
		// 422: datanya yang belum lengkap untuk diperiksa, bukan aplikasinya
		// yang rusak - kalimatnya menyebut apa yang kurang.
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to Save to RNM.")
	}
	return true
}

// simpanRNM melayani POST /api/klaim-life/{id}/outstanding.
func simpanRNM(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		hasil, err := services.SimpanRNMOracle(svc).Simpan(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.PathValue("id"), time.Now())
		if jawabGalatSimpanRNM(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
