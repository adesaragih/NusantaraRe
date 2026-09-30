package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/services"
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
			// Rincian galat tinggal di log server, tidak di badan jawaban.
			galat.Tulis(w, http.StatusInternalServerError, "gagal membaca klaim")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(klaim)
	}
}
