package handlers

// Pintu HTTP tombol `Save` layar Outstanding - 10-10-2026.
//
//	PUT /api/klaim-life/{id}/outstanding/isian
//
// `Section/InputOSClaimLife.xml` b19413 (Save): Accept peserta (b17582) dan medan tab Register Claim (b3504, b3810,
// b4934). Nol aturan dagang di sini.
//
// Dibaca sesudah: services/isianoutstanding.go.

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/claimlife/backend/models"
	"nusantarare/modul/claimlife/backend/services"
)

// batasBadanIsian - isian layar ini kecil (tiga teks + satu penanda per peserta); badan yang lebih besar ditolak.
const batasBadanIsian = 1 << 20

// simpanIsianOutstanding melayani PUT /api/klaim-life/{id}/outstanding/isian.
func simpanIsianOutstanding(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		var masuk services.IsianOutstanding
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanIsian)).Decode(&masuk); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "The request body is not valid JSON.")
			return
		}
		err := svc.SimpanOutstanding().Simpan(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("id"),
			masuk, time.Now())
		if err != nil {
			jawabGalatIsianOutstanding(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// jawabGalatIsianOutstanding - satu terjemahan galat tombol Save layar Outstanding.
func jawabGalatIsianOutstanding(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, inti.ErrTanpaIdentitas):
		galat.Tulis(w, http.StatusUnauthorized, "The request has no user identity.")
	case errors.Is(err, inti.ErrTanpaWewenang):
		galat.Tulis(w, http.StatusForbidden, "Only the Claim Admin can save this screen.")
	case errors.Is(err, kontrak.ErrKasusSudahTertutup):
		galat.Tulis(w, http.StatusConflict, "The case is already closed and cannot be changed.")
	case errors.Is(err, services.ErrIsianBukanOutstanding):
		galat.Tulis(w, http.StatusConflict, "Save is only available at the Claim Admin stage.")
	case errors.Is(err, models.ErrClaimGrossMelebihiShare):
		// Kalimat b526 VERBATIM.
		galat.Tulis(w, http.StatusUnprocessableEntity, models.ErrClaimGrossMelebihiShare.Error())
	case errors.Is(err, services.ErrAdjustmentTerkunci):
		galat.Tulis(w, http.StatusConflict, "The adjustment cannot be changed after Save to RNM.")
	case errors.Is(err, services.ErrTanggalTerkunciSesudahSaveRNM):
		galat.Tulis(w, http.StatusConflict, "The insured dates are locked after Save to RNM.")
	case errors.Is(err, services.ErrDOLTidakSah):
		// 11-10-2026: kalimat jelas dengan nomor peserta dan rentang valuasi; "Invalid DOL" (ValidasiDOL_Act) tetap kode.
		var dol services.GalatDOL
		if errors.As(err, &dol) {
			galat.Tulis(w, http.StatusUnprocessableEntity, dol.Kalimat())
		} else {
			galat.Tulis(w, http.StatusUnprocessableEntity, err.Error())
		}
	case errors.Is(err, services.ErrCatatanAdminTerkunci):
		galat.Tulis(w, http.StatusConflict, "Admin Notes cannot be changed for an insured that has been decided.")
	case errors.Is(err, services.ErrKatastrofeTerkunci):
		galat.Tulis(w, http.StatusConflict, "Catastrophe and Catastrophe Note cannot be changed once the claim has a number.")
	case errors.Is(err, services.ErrKlaimTidakAda):
		galat.Tulis(w, http.StatusNotFound, "The claim does not exist.")
	case errors.Is(err, services.ErrTahapTidakDikenal):
		galat.Tulis(w, http.StatusConflict, "The stage of this case is not recognized.")
	case errors.Is(err, galat.ErrPermintaanTidakSah):
		galat.Tulis(w, http.StatusBadRequest, err.Error())
	default:
		galat.Tulis(w, http.StatusInternalServerError, "Failed to save the Claim Admin screen.")
	}
}

// pesertaPolis melayani GET /api/klaim-life/{id}/peserta-polis?halaman=N - View Detail tab Policy Details
// (DetailPolisLife b631, 11-10-2026). Baca saja.
func pesertaPolis(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "The database is not configured.")
			return
		}
		halaman := 1
		if t := r.URL.Query().Get("halaman"); t != "" {
			n, err := strconv.Atoi(t)
			if err != nil || n < 1 {
				galat.Tulis(w, http.StatusBadRequest, "The page must be a whole number of at least 1.")
				return
			}
			halaman = n
		}
		hasil, err := svc.SimpanOutstanding().PesertaPolis(r.Context(), r.PathValue("id"), halaman)
		switch {
		case errors.Is(err, services.ErrKlaimTidakAda):
			galat.Tulis(w, http.StatusNotFound, "The claim does not exist.")
			return
		case errors.Is(err, galat.ErrPermintaanTidakSah):
			galat.Tulis(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			log.Printf("claimlife: membaca peserta polis klaim %q: %v", r.PathValue("id"), err)
			galat.Tulis(w, http.StatusInternalServerError, "Failed to load the policy insured.")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
