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

	"nusantarare/internal/services"
	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/kontrak"
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
		_ = json.NewEncoder(w).Encode(struct {
			Galat   string `json:"galat"`
			Langkah string `json:"langkah"`
		}{tolak.Pesan, tolak.Langkah})
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "hanya pemegang tahap Outstanding Claim yang dapat menyimpan ke RNM")
	case errors.Is(err, services.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, services.ErrSimpanRNMBukanOutstanding):
		galat.Tulis(w, http.StatusConflict, "Save to RNM hanya tersedia pada tahap Outstanding Claim")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "klaim tidak ada")
	case errors.Is(err, services.ErrTypeTidakDikenal),
		errors.Is(err, services.ErrBusinessCodeTidakDikenal),
		errors.Is(err, kontrak.ErrPolisNomorTakDitemukan):
		// 422: datanya yang belum lengkap untuk diperiksa, bukan aplikasinya
		// yang rusak - kalimatnya menyebut apa yang kurang.
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "gagal menyimpan ke RNM")
	}
	return true
}

// simpanRNM melayani POST /api/klaim-life/{id}/outstanding.
func simpanRNM(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
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
