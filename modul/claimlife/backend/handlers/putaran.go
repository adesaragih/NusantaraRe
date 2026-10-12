package handlers

// Pintu HTTP `Add` grid adjustment - putaran berikutnya (tiket 11).
//
// ⛔ Sejak GILIRAN-14 butir bp baris PERTAMA lahir saat Submit Register, tidak
// lagi di sini (ralat bo).
//
// Nol aturan dagang di sini. Yang memutuskan boleh atau tidak adalah services.
//
// Dibaca sesudah: handlers.go dan services/hasilkomite.go.

import (
	"errors"
	"net/http"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// tambahPutaran melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/putaran.
func tambahPutaran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		err := svc.Putaran().DenganJejak(jejak.PerekamJejakOracle(svc)).Tambah(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), time.Now())
		if jawabGalatPutaran(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatPutaran menerjemahkan galat `Add` ke kode HTTP; false = tidak
// ada galat, pemanggil meneruskan ke jalur berhasil.
func jawabGalatPutaran(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden,
			"Only ReasLifeSPV can add an adjustment row.")
	case errors.Is(err, services.ErrBukanPenolakan):
		// 409: bukan permintaan yang salah, melainkan bentrokan dengan
		// keadaan yang ada - baris terakhirnya belum ditolak.
		galat.Tulis(w, http.StatusConflict,
			"A new adjustment round can only start after the last row is rejected.")
	case errors.Is(err, services.ErrTahapTanpaAddAdjustment):
		// Padanan syarat tampil b18160 - di Pega tombolnya tidak ada.
		galat.Tulis(w, http.StatusConflict,
			"Adding an adjustment row is only available at the Claim Analyst stage.")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat.Tulis(w, http.StatusConflict, "The stage of this case is not recognized.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to add the adjustment row.")
	}
	return true
}
