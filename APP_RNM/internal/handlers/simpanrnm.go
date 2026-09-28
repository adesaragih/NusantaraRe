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
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "hanya pemegang tahap Outstanding Claim yang dapat menyimpan ke RNM")
	case errors.Is(err, services.ErrKasusSudahTertutup):
		galat(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, services.ErrSimpanRNMBukanOutstanding):
		galat(w, http.StatusConflict, "Save to RNM hanya tersedia pada tahap Outstanding Claim")
	case errors.Is(err, services.ErrTypeTidakDikenal),
		errors.Is(err, services.ErrBusinessCodeTidakDikenal),
		errors.Is(err, services.ErrAmbangSTNCBelumDiketahui),
		errors.Is(err, services.ErrTanggalTerimaPolisKosong),
		errors.Is(err, services.ErrKategoriWajibBelumDiketahui),
		errors.Is(err, services.ErrAmbangProdukTakDitemukan),
		errors.Is(err, services.ErrPolisNomorTakDitemukan):
		// 422: datanya yang belum lengkap untuk diperiksa, bukan aplikasinya
		// yang rusak - kalimatnya menyebut apa yang kurang.
		galat(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	default:
		galat(w, http.StatusInternalServerError, "gagal menyimpan ke RNM")
	}
	return true
}

// simpanRNM melayani POST /api/klaim-life/{id}/outstanding.
func simpanRNM(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := services.SimpanRNMOracle(svc).Simpan(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("id"), time.Now())
		if jawabGalatSimpanRNM(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
