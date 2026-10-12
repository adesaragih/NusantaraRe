package handlers

// Pintu HTTP tombol `Save` layar Claim Analis - 12-10-2026.
//
//	PUT /api/klaim-life/{id}/analis
//
// CLAIM ANALYST (NOTES, b15194) per peserta. Nol aturan dagang di sini.
//
// Dibaca sesudah: services/analis.go.

import (
	"encoding/json"
	"errors"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// simpanAnalis melayani PUT /api/klaim-life/{id}/analis.
func simpanAnalis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		var masuk services.IsianAnalis
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanIsian)).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request body is not valid JSON.")
			return
		}
		if err := svc.SimpanAnalis().Simpan(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"), masuk); err != nil {
			jawabGalatAnalis(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatAnalis - satu terjemahan galat tombol Save catatan Claim Analyst.
func jawabGalatAnalis(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "Only the Claim Analyst can save the Claim Analyst notes.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrAnalisBukanClaimAnalis):
		galat.Tulis(w, http.StatusConflict, "Save is only available at the Claim Analyst stage.")
	case errors.Is(err, services.ErrCatatanAnalisGagalTulis):
		galat.Tulis(w, http.StatusConflict, "The Claim Analyst notes cannot be saved for this insured.")
	case errors.Is(err, services.ErrCatatanAnalisKepanjangan):
		galat.Tulis(w, http.StatusUnprocessableEntity, "The Claim Analyst notes cannot exceed 1000 characters.")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "The claim does not exist.")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat.Tulis(w, http.StatusConflict, "The stage of this case is not recognized.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to save the Claim Analyst notes.")
	}
}
