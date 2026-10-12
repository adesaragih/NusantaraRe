package handlers

// Pintu HTTP cabut peserta - OQ-M6 (GILIRAN-17). Nol aturan dagang di sini.
//
// Dibaca sesudah: services/cabutpeserta.go.

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

// cabutPeserta melayani POST /api/klaim-life/{id}/peserta/{pesertaId}/cabut.
func cabutPeserta(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		err := svc.Status().DenganJejak(jejak.PerekamJejakOracle(svc)).CabutPeserta(r.Context(),
			inti.PelakuDari(r, stubPelaku), r.PathValue("id"), r.PathValue("pesertaId"), time.Now())
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
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "Only ReasLifeAdmin can remove an insured.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrPesertaTidakDapatDicabut):
		galat.Tulis(w, http.StatusConflict,
			"An insured can only be removed at the Claim Admin stage before Save to RNM.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity, "The stage of this case is not recognized.")
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to remove the insured.")
	}
}
