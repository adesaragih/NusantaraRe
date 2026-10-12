package handlers

// Pintu HTTP tombol `Save` layar Medical Check - 11-10-2026.
//
//	PUT /api/klaim-life/{id}/medis
//
// `FlowAction/MedicalCheckClaimLife.xml` Save: RECOMMENDATION per peserta (`ClaimLifeDetailGCNM` b7323). Nol aturan
// dagang di sini.
//
// Dibaca sesudah: services/medis.go.

import (
	"encoding/json"
	"errors"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// simpanMedis melayani PUT /api/klaim-life/{id}/medis.
func simpanMedis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		var masuk services.IsianMedis
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanIsian)).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request body is not valid JSON.")
			return
		}
		if err := svc.SimpanMedis().Simpan(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"), masuk); err != nil {
			jawabGalatMedis(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatMedis - satu terjemahan galat tombol Save layar Medical Check.
func jawabGalatMedis(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "Only the Medical Advisor can save the Medical Check.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrMedisBukanMedicalCheck):
		galat.Tulis(w, http.StatusConflict, "Save is only available at the Medical Check stage.")
	case errors.Is(err, services.ErrRekomendasiTerkunci):
		galat.Tulis(w, http.StatusConflict, "The recommendation cannot be changed for an insured that has been decided.")
	case errors.Is(err, services.ErrRekomendasiKepanjangan):
		galat.Tulis(w, http.StatusUnprocessableEntity, "The recommendation cannot exceed 1000 characters.")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "The claim does not exist.")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat.Tulis(w, http.StatusConflict, "The stage of this case is not recognized.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to save the Medical Check.")
	}
}
