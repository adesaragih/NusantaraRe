package handlers

// GET /api/nbfacin/risk-address?address=&zipCode=&country=&province=&city=&district=&territory=&halaman=
// - popup Choose Risk Address (tiket 36); kontrak `cariRiskAddress` frontend.

import (
	"net/http"
	"strconv"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type barisRisk struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Address       string `json:"address"`
	NationName    string `json:"nationName"`
	ProvinceName  string `json:"provinceName"`
	CityName      string `json:"cityName"`
	DistrictName  string `json:"districtName"`
	TerritoryName string `json:"territoryName"`
	PostalCode    string `json:"postalCode"`
}

func cariRisk(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		halaman := 1
		if h := q.Get("halaman"); h != "" {
			n, err := strconv.Atoi(h)
			if err != nil {
				n = 0 // tak sah -> 400 dari services
			}
			halaman = n
		}
		hasil, err := svc.CariRisk(r.Context(), models.SaringRisk{Address: q.Get("address"), ZipCode: q.Get("zipCode"),
			Country: q.Get("country"), Province: q.Get("province"), City: q.Get("city"), District: q.Get("district"),
			Territory: q.Get("territory")}, halaman)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisRisk, 0, len(hasil.Baris))
		for _, b := range hasil.Baris {
			baris = append(baris, barisRisk{ID: b.ID, Title: b.Title, Address: b.Address, NationName: b.NationName,
				ProvinceName: b.ProvinceName, CityName: b.CityName, DistrictName: b.DistrictName,
				TerritoryName: b.TerritoryName, PostalCode: b.PostalCode})
		}
		galat.TulisJSON(w, struct {
			Baris   []barisRisk `json:"baris"`
			Total   int         `json:"total"`
			Halaman int         `json:"halaman"`
			Ukuran  int         `json:"ukuran"`
		}{baris, hasil.Total, hasil.Halaman, hasil.Ukuran})
	}
}
