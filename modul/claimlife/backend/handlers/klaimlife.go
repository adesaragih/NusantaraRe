package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/services"
)

// klaimLife melayani pembacaan satu klaim Life beserta seluruh barisnya.
//
// Tiket 01 AC-7: handlers tidak memanggil repository langsung. Berkas ini hanya
// menerjemahkan HTTP ke panggilan services dan sebaliknya - nol aturan dagang.
func klaimLife(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		layanan := svc.KlaimLife()
		if layanan == nil {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}

		klaim, err := layanan.Ambil(r.Context(), r.PathValue("id"))
		switch {
		case errors.Is(err, services.ErrKlaimTidakAda):
			galat.Tulis(w, http.StatusNotFound, "klaim tidak ada")
			return
		case errors.Is(err, services.ErrWajibIsi):
			galat.Tulis(w, http.StatusBadRequest, "pengenal klaim wajib diisi")
			return
		case err != nil:
			// Rincian galat tinggal di log server, tidak di badan jawaban. ⛔ 10-10-2026: dulu TIDAK dicatat sama sekali
			// (komentar ini berjanji, kodenya tidak) - laporan WO "gagal membaca klaim" tanpa jejak di log.
			log.Printf("claimlife: membaca klaim %q: %v", r.PathValue("id"), err)
			galat.Tulis(w, http.StatusInternalServerError, "Failed to load the claim.")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(klaim)
	}
}

// spreadingPolis melayani GET /api/klaim-life/{id}/peserta/{pesertaId}/spreading-polis - spreading PremiumList Life peserta
// itu (keputusan work owner 11-10-2026). Baca saja.
func spreadingPolis(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		layanan := svc.KlaimLife()
		if layanan == nil {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		hasil, err := layanan.SpreadingPolis(r.Context(), r.PathValue("id"), r.PathValue("pesertaId"))
		switch {
		case errors.Is(err, services.ErrKlaimTidakAda):
			galat.Tulis(w, http.StatusNotFound, "peserta klaim tidak ada")
			return
		case errors.Is(err, services.ErrWajibIsi):
			galat.Tulis(w, http.StatusBadRequest, "pengenal klaim dan peserta wajib diisi")
			return
		case err != nil:
			log.Printf("claimlife: membaca spreading polis klaim %q peserta %q: %v", r.PathValue("id"), r.PathValue("pesertaId"), err)
			galat.Tulis(w, http.StatusInternalServerError, "Failed to load the spreading.")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
