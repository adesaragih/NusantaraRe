package handlers

// Pintu HTTP penyerahan ke Komite - tiket 10; pop-up Send Claim to Committee (keputusan work owner 12-10-2026 opsi A).
//
//	GET  /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite  - pratinjau pop-up (nol tulisan)
//	POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite  - serahkan, badan {"remarks": "..."}
//
// Nol aturan dagang di sini: berkas ini menerjemahkan HTTP ke panggilan
// services dan sebaliknya. Yang memutuskan boleh atau tidak adalah services.
//
// Dibaca sesudah: handlers.go dan services/komite.go.

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

// penyerahanKomite merakit layanan penyerahan dengan ketergantungan Oracle-nya (satu tempat untuk GET dan POST).
func penyerahanKomite(svc *services.Service) *services.Penyerahan {
	return svc.Komite().DenganJejak(jejak.PerekamJejakOracle(svc)).
		DenganRoster(services.RosterKomiteOracle(svc)).
		DenganKasus(services.KasusKomiteOracle(svc)).
		DenganPenyalur(services.PenyalurClaimLifeOracle(svc))
}

// badanSerahkanKomite - badan POST penyerahan.
type badanSerahkanKomite struct {
	Remarks string `json:"remarks"`
}

// serahkanKomite melayani
// POST /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite.
func serahkanKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		var masuk badanSerahkanKomite
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanIsian)).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request body is not valid JSON.")
			return
		}
		err := penyerahanKomite(svc).Serahkan(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), r.PathValue("adjId"), masuk.Remarks, time.Now())
		if jawabGalatKomite(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// pratinjauKomite melayani
// GET /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}/komite - isi pop-up sebelum dikirim.
func pratinjauKomite(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		hasil, err := penyerahanKomite(svc).Pratinjau(r.Context(), inti.PelakuDari(r, stubPelaku),
			r.PathValue("id"), r.PathValue("pesertaId"), r.PathValue("adjId"), time.Now())
		if jawabGalatKomite(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(hasil)
	}
}

// jawabGalatKomite - satu terjemahan galat penyerahan / pratinjau ke Komite, kalimat Inggris (12-10-2026). false = tanpa
// galat (pemanggil melanjutkan).
func jawabGalatKomite(w http.ResponseWriter, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		// 403: identitasnya ada, perannya yang kurang untuk Type ini.
		galat.Tulis(w, http.StatusForbidden, "Your role is not allowed to send this claim type to the committee.")
	case errors.Is(err, services.ErrRemarksKomiteKosong):
		galat.Tulis(w, http.StatusUnprocessableEntity, "Remarks is required.")
	case errors.Is(err, services.ErrRemarksKomiteKepanjangan):
		galat.Tulis(w, http.StatusUnprocessableEntity, "Remarks cannot exceed 4000 characters.")
	case errors.Is(err, services.ErrTypeTidakDikenal):
		galat.Tulis(w, http.StatusUnprocessableEntity,
			"The claim type is not recognized; sending to the committee requires a valid type.")
	case errors.Is(err, services.ErrRekeningBelumLengkap):
		// 422: ⛔ Pesannya diteruskan APA ADANYA - teks sistem lama ("Name of bank cannot be empty").
		galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, services.ErrSerahkanBukanClaimAnalis):
		galat.Tulis(w, http.StatusConflict, "Send Claim to Committee is only available at the Claim Analyst stage.")
	case errors.Is(err, services.ErrBarisSudahDiserahkan):
		galat.Tulis(w, http.StatusConflict, "This row has already been sent to the committee.")
	case errors.Is(err, services.ErrBarisBukanOutstanding):
		galat.Tulis(w, http.StatusConflict, "Only Outstanding rows can be sent to the committee.")
	case errors.Is(err, services.ErrMataUangKlaimCampur):
		galat.Tulis(w, http.StatusUnprocessableEntity,
			"The rows of this claim use more than one currency; sending to the committee requires a single currency.")
	case errors.Is(err, services.ErrNilaiKlaimKomiteKosong):
		galat.Tulis(w, http.StatusUnprocessableEntity,
			"Claim RNM of this row is empty; the committee level cannot be determined.")
	case errors.Is(err, services.ErrRosterKomiteKosong):
		galat.Tulis(w, http.StatusUnprocessableEntity, "No committee level covers this claim amount.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat.Tulis(w, http.StatusConflict, "The stage of this case is not recognized.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to send the row to the committee.")
	}
	return true
}
