package handlers

// Saran Occupation sub-tab Surrounding Risk (tiket 38): GET /api/nbfacin/occupation?cari=.

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// barisOccupation - kontrak `BarisOccupation` frontend (tiket 38; kdRiskExposure tiket 40).
type barisOccupation struct {
	OldID          string `json:"oldId"`
	Name           string `json:"name"`
	KdRiskExposure string `json:"kdRiskExposure"`
}

// cariOccupation - GET /api/nbfacin/occupation?cari= (tiket 38/40), kontrak `BarisOccupation`.
func cariOccupation(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.CariOccupation(r.Context(), r.URL.Query().Get("cari"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisOccupation, 0, len(hasil))
		for _, b := range hasil {
			baris = append(baris, barisOccupation{OldID: b.OldID, Name: b.Name, KdRiskExposure: b.KdRiskExposure})
		}
		galat.TulisJSON(w, struct {
			Baris []barisOccupation `json:"baris"`
		}{baris})
	}
}
