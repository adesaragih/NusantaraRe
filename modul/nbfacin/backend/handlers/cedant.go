package handlers

// Tab Inw Fac Cedant Panels kasus FIRE (tiket 49): GET / PUT /api/nbfacin/kasus/{caseId}/cedant. Kontrak frontend
// modul/nbfacin/frontend/api.ts `TampilanCedant` / `BarisCedant`.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// batasBadanCedant - badan PUT (paling banyak 100 baris cedant).
const batasBadanCedant = 1 << 20

type barisCedantKabel struct {
	CedingCo     string `json:"cedingCo"`
	CedingCoName string `json:"cedingCoName"`
	ShareCeding  string `json:"shareCeding"`
}

type tampilanCedantKabel struct {
	SobName         string             `json:"sobName"`
	ShareCedantType string             `json:"shareCedantType"`
	PercentShare    string             `json:"percentShare"`
	TotalTSIRNM     string             `json:"totalTsiRnm"`
	TotalPremiRNM   string             `json:"totalPremiRnm"`
	Wajib           bool               `json:"wajib"`
	Cedant          []barisCedantKabel `json:"cedant"`
	CedingUmum      []barisCedantKabel `json:"cedingUmum"`
}

func keBarisCedantKabel(d []models.BarisCedant) []barisCedantKabel {
	out := make([]barisCedantKabel, 0, len(d))
	for _, b := range d {
		out = append(out, barisCedantKabel{b.CedingCo, b.CedingCoName, teks(b.ShareCeding)})
	}
	return out
}

func keCedantKabel(d services.TampilanCedant) tampilanCedantKabel {
	return tampilanCedantKabel{SobName: d.SobName, ShareCedantType: d.ShareCedantType, PercentShare: teks(d.PercentShare),
		TotalTSIRNM: teks(d.TotalTSIRNM), TotalPremiRNM: teks(d.TotalPremiRNM), Wajib: d.Wajib,
		Cedant: keBarisCedantKabel(d.Cedant), CedingUmum: keBarisCedantKabel(d.CedingUmum)}
}

// bacaCedant - GET /api/nbfacin/kasus/{caseId}/cedant.
func bacaCedant(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.TampilanCedant(r.Context(), r.PathValue("caseId"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keCedantKabel(d))
	}
}

// simpanCedant - PUT …/cedant badan {shareCedantType, cedant: BarisCedant[]} (kunci asing -> 400): jawab baca ulang.
func simpanCedant(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			ShareCedantType *string            `json:"shareCedantType"`
			Cedant          []barisCedantKabel `json:"cedant"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanCedant))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil || b.ShareCedantType == nil || b.Cedant == nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus {\"shareCedantType\": teks, \"cedant\": [baris]}")
			return
		}
		var masalah []string
		baris := make([]models.BarisCedant, 0, len(b.Cedant))
		for i, c := range b.Cedant {
			baris = append(baris, models.BarisCedant{CedingCo: strings.TrimSpace(c.CedingCo), CedingCoName: c.CedingCoName,
				ShareCeding: services.UraiDesimalIsian(fmt.Sprintf("cedant[%d].shareCeding", i), strings.TrimSpace(c.ShareCeding), &masalah)})
		}
		if len(masalah) > 0 {
			tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanCedant, strings.Join(masalah, "; ")))
			return
		}
		d, err := svc.SimpanCedant(r.Context(), inti.PelakuDari(r, stubPelaku), r.PathValue("caseId"), *b.ShareCedantType, baris)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keCedantKabel(d))
	}
}
