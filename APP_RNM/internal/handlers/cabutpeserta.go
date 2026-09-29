package handlers

// Pintu HTTP cabut peserta - OQ-M6 (GILIRAN-17). Nol aturan dagang di sini.
//
// Dibaca sesudah: services/cabutpeserta.go.

import (
	"errors"
	"net/http"
	"time"

	"nusantarare/internal/services"
)

// cabutPeserta melayani POST /api/klaim-life/{id}/peserta/{pesertaId}/cabut.
func cabutPeserta(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		err := svc.Status().DenganJejak(services.PerekamJejakOracle(svc)).CabutPeserta(r.Context(),
			pelakuDari(r, stubPelaku), r.PathValue("id"), r.PathValue("pesertaId"), time.Now())
		if err != nil {
			jawabGalatCabut(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatCabut - satu terjemahan galat cabut peserta.
func jawabGalatCabut(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, services.ErrTanpaIdentitas):
		galat(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
	case errors.Is(err, services.ErrTanpaWewenang):
		galat(w, http.StatusForbidden, "hanya ReasLifeAdmin yang dapat mencabut peserta")
	case errors.Is(err, services.ErrKasusSudahTertutup):
		galat(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
	case errors.Is(err, services.ErrPesertaTidakDapatDicabut):
		galat(w, http.StatusConflict,
			"peserta hanya dapat dicabut di tahap Outstanding Claim sebelum Save to RNM")
	case errors.Is(err, services.ErrPermintaanTidakSah):
		galat(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat(w, http.StatusUnprocessableEntity, "tahap kasus tidak dikenal")
	default:
		galat(w, http.StatusInternalServerError, "gagal mencabut peserta")
	}
}
