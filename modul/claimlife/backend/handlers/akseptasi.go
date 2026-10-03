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
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		nomor, err := svc.Akseptasi().DenganJejak(jejak.PerekamJejakOracle(svc)).
			DenganPenerbit(services.PenerbitAkseptasiOracle(svc)).
			SimpanAdjustment(r.Context(), inti.PelakuDari(r, stubPelaku),
				r.PathValue("id"), r.PathValue("pesertaId"), time.Now())

		switch {
		case errors.Is(err, inti.ErrTanpaIdentitas):
			galat.Tulis(w, http.StatusUnauthorized, "permintaan tanpa identitas pelaku ditolak")
			return
		case errors.Is(err, inti.ErrTanpaWewenang):
			// 403: yang boleh mengaksep adalah PEMEGANG TAHAP klaim ini -
			// Outstanding oleh ReasLifeAdmin, Medical Check oleh
			// ReasLifeMedicalAdvisor, Claim Analis oleh ReasLifeSPV.
			galat.Tulis(w, http.StatusForbidden,
				"hanya pemegang tahap klaim ini yang dapat mengaksep barisnya")
			return
		case errors.Is(err, services.ErrPesertaTidakDipilih):
			galat.Tulis(w, http.StatusConflict,
				"peserta belum dipilih untuk diklaim")
			return
		case errors.Is(err, services.ErrBarisSudahBernomorAkseptasi):
			galat.Tulis(w, http.StatusConflict, "baris sudah punya nomor akseptasi")
			return
		case errors.Is(err, services.ErrBarisBukanOutstanding):
			galat.Tulis(w, http.StatusConflict,
				"hanya baris Outstanding yang dapat diaksep")
			return
		case errors.Is(err, services.ErrTypeTidakDikenal):
			galat.Tulis(w, http.StatusUnprocessableEntity,
				"Type klaim tidak dikenal; nomor akseptasi tidak dapat diterbitkan")
			return
		case errors.Is(err, kontrak.ErrKodeBisnisBelumTersimpan):
			// 501: temuan audit A0 - kolomnya belum ada, dan nomor akseptasi
			// memuatnya. Bukan kerusakan, melainkan yang belum dibangun.
			galat.Tulis(w, http.StatusNotImplemented,
				"kode bisnis klaim belum tersimpan; nomor akseptasi belum dapat dirakit")
			return
		case errors.Is(err, services.ErrNomorAkseptasiBerganda):
			galat.Tulis(w, http.StatusConflict,
				"nomor akseptasi yang terbit sudah dipakai; coba lagi")
			return
		case errors.Is(err, kontrak.ErrKasusSudahTertutup):
			galat.Tulis(w, http.StatusConflict, "kasus sudah ditutup dan tidak dapat diubah")
			return
		case errors.Is(err, galat.ErrPermintaanTidakSah):
			galat.Tulis(w, http.StatusBadRequest, err.Error())
			return
		case err != nil:
			galat.Tulis(w, http.StatusInternalServerError, "gagal mengaksep baris")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(struct {
			NomorAkseptasi string `json:"nomorAkseptasi"`
		}{nomor})
	}
}
