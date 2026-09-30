package handlers

// Pintu HTTP tombol portal `Input Offer` / `Input Premium` - GILIRAN-13
// paket 1 (butir bn).
//
//	POST /api/polis-life   badan {"flag":"0"|"1"}
//
// `Section/PremiumList.xml` `Input Offer` b3273 / `Input Premium` b3921 ->
// `CreateInputLife`. Nol aturan dagang di sini; lihat services/polis_kasus.go.

import (
	"encoding/json"
	"net/http"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/premiumlistlife/services"
)

// isiBuatKasusPolis - `FlagPolicy` tombolnya, VERBATIM.
type isiBuatKasusPolis struct {
	Flag string `json:"flag"`
}

// buatKasusPolis melayani POST /api/polis-life.
func buatKasusPolis(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.PunyaDatabase() {
			galat.Tulis(w, http.StatusServiceUnavailable, "database belum dikonfigurasi")
			return
		}
		var isi isiBuatKasusPolis
		if err := json.NewDecoder(r.Body).Decode(&isi); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON yang sah")
			return
		}
		hasil, err := svc.KasusPolis().
			DenganJejak(jejak.PerekamJejakOracle(svc)).
			Buat(r.Context(), inti.PelakuDari(r, stubPelaku), isi.Flag, time.Now())
		if jawabGalatPolis(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(hasil)
	}
}
