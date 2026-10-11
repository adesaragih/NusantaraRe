package handlers

// Copy To All Spreading per lokasi (tiket 73 P5, K73-40): POST …/spreading/salin-lokasi (FIRE) dan …/spreading-aneka/salin-lokasi
// (keluarga Aneka). Badan {percentShare, idxLokasi (1-based), template:[{treatyType, treatyName, sharePercentage, flagSpreading}]};
// jawaban = Tampilan salin dengan `template: []` (frontend memegang template per lokasi di state sendiri). Tidak menyimpan.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type templateLokasiKabel struct {
	TreatyType      string `json:"treatyType"`
	TreatyName      string `json:"treatyName"`
	SharePercentage string `json:"sharePercentage"`
	FlagSpreading   string `json:"flagSpreading"`
}

type badanSalinLokasi struct {
	PercentShare *string               `json:"percentShare"`
	IdxLokasi    *int                  `json:"idxLokasi"`
	Template     []templateLokasiKabel `json:"template"`
}

// uraiBadanSalinLokasi - JSON + percentShare + idxLokasi + template; galat 400 sudah ditulis bila ok=false.
func uraiBadanSalinLokasi(w http.ResponseWriter, r *http.Request) (ps *apd.Decimal, idx int, templat []models.TemplateSpreading, ok bool) {
	var b badanSalinLokasi
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanSpreading)).Decode(&b); err != nil || b.PercentShare == nil ||
		b.IdxLokasi == nil {
		galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus memuat percentShare dan idxLokasi")
		return nil, 0, nil, false
	}
	var masalah []string
	ps = services.UraiDesimalIsian("percentShare", *b.PercentShare, &masalah)
	templat = make([]models.TemplateSpreading, 0, len(b.Template))
	for i, t := range b.Template {
		templat = append(templat, models.TemplateSpreading{TreatyType: t.TreatyType, TreatyName: t.TreatyName, FlagSpreading: t.FlagSpreading,
			SharePercentage: services.UraiDesimalIsian(fmt.Sprintf("template[%d].sharePercentage", i), t.SharePercentage, &masalah)})
	}
	if len(masalah) > 0 {
		tulisGalat(w, fmt.Errorf("%w: %s", services.ErrMasukanSpreading, strings.Join(masalah, "; ")))
		return nil, 0, nil, false
	}
	return ps, *b.IdxLokasi, templat, true
}

// salinSpreadingLokasi - POST …/spreading/salin-lokasi (FIRE).
func salinSpreadingLokasi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		ps, idx, templat, ok := uraiBadanSalinLokasi(w, r)
		if !ok {
			return
		}
		d, err := svc.SalinSpreadingLokasi(r.Context(), pelaku, r.PathValue("caseId"), ps, idx, templat)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingKabel(d))
	}
}

// salinSpreadingAnekaLokasi - POST …/spreading-aneka/salin-lokasi (keluarga Aneka).
func salinSpreadingAnekaLokasi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pelaku := inti.PelakuDari(r, stubPelaku)
		if err := inti.WajibIdentitas(pelaku); err != nil {
			tulisGalat(w, err)
			return
		}
		ps, idx, templat, ok := uraiBadanSalinLokasi(w, r)
		if !ok {
			return
		}
		d, err := svc.SalinSpreadingAnekaLokasi(r.Context(), pelaku, r.PathValue("caseId"), ps, idx, templat)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, keSpreadingAnekaKabel(d))
	}
}
