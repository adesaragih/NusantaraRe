package handlers

// Popup Choose Class of Construction (tiket 40): GET /api/nbfacin/kasus/{caseId}/table-of-limit?category=.

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// barisTableOfLimit - kontrak `BarisTableOfLimit` frontend.
type barisTableOfLimit struct {
	Description string `json:"description"`
	PctLimit    string `json:"pctLimit"`
}

// cariTableOfLimit - GET /api/nbfacin/kasus/{caseId}/table-of-limit?category=: 409 bila kode bisnis case tidak dapat
// ditentukan.
func cariTableOfLimit(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.TableOfLimit(r.Context(), r.PathValue("caseId"), r.URL.Query().Get("category"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisTableOfLimit, 0, len(hasil))
		for _, b := range hasil {
			baris = append(baris, barisTableOfLimit{Description: b.Description, PctLimit: b.PctLimit})
		}
		galat.TulisJSON(w, struct {
			Baris []barisTableOfLimit `json:"baris"`
		}{baris})
	}
}
