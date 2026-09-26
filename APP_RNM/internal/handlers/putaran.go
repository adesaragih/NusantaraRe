package handlers

// Pintu HTTP putaran adjustment berikutnya - tiket 11.
//
// Nol aturan dagang di sini. Yang memutuskan boleh atau tidak adalah services.
//
// Dibaca sesudah: handlers.go dan services/hasilkomite.go.

import (
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// tambahPutaran melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/putaran.
func tambahPutaran(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Putaran().Tambah(r.Context(), pelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), time.Now())

		switch {
		case errors.Is(err, services.ErrTanpaIdentitas):
			galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, services.ErrTanpaWewenang):
			galat(w, http.StatusForbidden,
				"hanya ReasLifeSPV yang dapat membuka putaran adjustment berikutnya")
			return
		case errors.Is(err, services.ErrBukanPenolakan):
			// 409: bukan permintaan yang salah, melainkan bentrokan dengan
			// keadaan yang ada - baris terakhirnya belum ditolak.
			galat(w, http.StatusConflict,
				"putaran berikutnya hanya lahir sesudah baris terakhir ditolak")
			return
		case errors.Is(err, services.ErrJejakBelumDiputuskan):
			galat(w, http.StatusNotImplemented,
				"jejak audit belum dapat direkam: tempatnya belum diputuskan work owner")
			return
		case errors.Is(err, services.ErrPermintaanTidakSah):
			galat(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat(w, http.StatusInternalServerError,
				"gagal membuka putaran adjustment berikutnya")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
