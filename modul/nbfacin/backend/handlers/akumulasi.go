package handlers

// Popup Choose Accumulation Code (tiket 46): GET /api/nbfacin/akumulasi dan GET /api/nbfacin/akumulasi/saran/{jenis}.
// Kontrak frontend modul/nbfacin/frontend/api.ts.

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// barisAkumulasiKabel - kontrak `{ id, accumulationName, note }`.
type barisAkumulasiKabel struct {
	ID               string `json:"id"`
	AccumulationName string `json:"accumulationName"`
	Note             string `json:"note"`
}

// saranAkumulasiKabel - kontrak `{ id, label, ekstra? }`.
type saranAkumulasiKabel struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Ekstra string `json:"ekstra,omitempty"`
}

// cariAkumulasi - GET /api/nbfacin/akumulasi?id=&policyNo=&note=&postalCode=&syariahStatus=&provinceId=&cityId=
// &districtId=&czone=&keyword=.
func cariAkumulasi(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		d, err := svc.CariAkumulasi(r.Context(), models.SaringAkumulasi{ID: q.Get("id"), PolicyNo: q.Get("policyNo"),
			Note: q.Get("note"), PostalCode: q.Get("postalCode"), SyariahStatus: q.Get("syariahStatus"),
			ProvinceID: q.Get("provinceId"), CityID: q.Get("cityId"), DistrictID: q.Get("districtId"), CZone: q.Get("czone"),
			Keyword: q.Get("keyword")})
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisAkumulasiKabel, 0, len(d))
		for _, b := range d {
			baris = append(baris, barisAkumulasiKabel(b))
		}
		galat.TulisJSON(w, struct {
			Baris []barisAkumulasiKabel `json:"baris"`
		}{baris})
	}
}

// saranAkumulasi - GET /api/nbfacin/akumulasi/saran/{jenis}?q=&induk=.
func saranAkumulasi(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.SaranAkumulasi(r.Context(), r.PathValue("jenis"), r.URL.Query().Get("q"), r.URL.Query().Get("induk"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]saranAkumulasiKabel, 0, len(d))
		for _, b := range d {
			baris = append(baris, saranAkumulasiKabel(b))
		}
		galat.TulisJSON(w, struct {
			Baris []saranAkumulasiKabel `json:"baris"`
		}{baris})
	}
}
