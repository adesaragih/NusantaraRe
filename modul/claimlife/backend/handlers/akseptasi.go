package handlers

// Pintu HTTP akseptasi baris adjustment - audit A0.
//
// ⭐ Jalur ini LAHIR dari audit: `SaveAdjustment_Act` mengaksep di Claim Life
// sendiri, dan modul ini sebelumnya hanya mengenal jalur Komite.
//
// Nol aturan dagang di sini. Dibaca sesudah: handlers.go, services/akseptasi.go.

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/services"
)

// simpanAdjustment melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/akseptasi.
func simpanAdjustment(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		nomor, err := svc.Akseptasi().DenganJejak(jejak.PerekamJejakOracle(svc)).
			DenganPenerbit(services.PenerbitAkseptasiOracle(svc)).
			SimpanAdjustment(r.Context(), inti.PelakuDari(r, stubPelaku),
				r.PathValue("id"), r.PathValue("pesertaId"), time.Now())

		switch {
		case errors.Is(err, inti.ErrTanpaIdentitas):
			galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
			return
		case errors.Is(err, inti.ErrTanpaWewenang):
			// 403: yang boleh mengaksep adalah PEMEGANG TAHAP klaim ini -
			// Outstanding oleh ReasLifeAdmin, Medical Check oleh
			// ReasLifeMedicalAdvisor, Claim Analis oleh ReasLifeSPV.
			galat.Tulis(w, http.StatusForbidden,
				"Only the holder of this claim stage can accept its rows.")
			return
		case errors.Is(err, services.ErrPesertaTidakDipilih):
			galat.Tulis(w, http.StatusConflict,
				"The insured has not been selected for the claim.")
			return
		case errors.Is(err, services.ErrBarisSudahBernomorAkseptasi):
			galat.Tulis(w, http.StatusConflict, "The row already has an acceptance number.")
			return
		case errors.Is(err, services.ErrBarisBukanOutstanding):
			galat.Tulis(w, http.StatusConflict,
				"Only Outstanding rows can be accepted.")
			return
		case errors.Is(err, services.ErrTypeTidakDikenal):
			galat.Tulis(w, http.StatusUnprocessableEntity,
				"The claim type is not recognized; the acceptance number cannot be issued.")
			return
		case errors.Is(err, kontrak.ErrKodeBisnisBelumTersimpan):
			// 501: temuan audit A0 - kolomnya belum ada, dan nomor akseptasi
			// memuatnya. Bukan kerusakan, melainkan yang belum dibangun.
			galat.Tulis(w, http.StatusNotImplemented,
				"The claim business code is not saved yet; the acceptance number cannot be built.")
			return
		case errors.Is(err, services.ErrNomorAkseptasiBerganda):
			galat.Tulis(w, http.StatusConflict,
				"The issued acceptance number is already in use; please try again.")
			return
		case errors.Is(err, kontrak.ErrKasusSudahTertutup):
			galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
			return
		case errors.Is(err, galat.ErrPermintaanTidakSah):
			galat.Tulis(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat.Tulis(w, http.StatusInternalServerError, "Failed to accept the row.")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			NomorAkseptasi string `json:"nomorAkseptasi"`
		}{nomor})
	}
}
