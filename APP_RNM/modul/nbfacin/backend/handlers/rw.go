package handlers

// Popup Add alamat risiko (tiket 37): GET /api/nbfacin/rw?zipCode= dan POST /api/nbfacin/risk-address.

import (
	"encoding/json"
	"net/http"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type barisRW struct {
	ZipCode       string `json:"zipCode"`
	TerritoryName string `json:"territoryName"`
	DistrictName  string `json:"districtName"`
	CityName      string `json:"cityName"`
	ProvinceName  string `json:"provinceName"`
	NationName    string `json:"nationName"`
}

type alamatBaru struct {
	NationName    string `json:"nationName"`
	ProvinceName  string `json:"provinceName"`
	DistrictName  string `json:"districtName"`
	CityName      string `json:"cityName"`
	TerritoryName string `json:"territoryName"`
	Title         string `json:"title"`
	Address       string `json:"address"`
	PostalCode    string `json:"postalCode"`
}

// batasBadanAlamat - delapan medan, terlebar 4000 bita.
const batasBadanAlamat = 64 << 10

// cariRW - GET /api/nbfacin/rw?zipCode=.
func cariRW(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.CariRW(r.Context(), r.URL.Query().Get("zipCode"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]barisRW, 0, len(hasil))
		for _, b := range hasil {
			baris = append(baris, barisRW{ZipCode: b.ZipCode, TerritoryName: b.TerritoryName, DistrictName: b.DistrictName,
				CityName: b.CityName, ProvinceName: b.ProvinceName, NationName: b.NationName})
		}
		galat.TulisJSON(w, struct {
			Baris []barisRW `json:"baris"`
		}{baris})
	}
}

// simpanAlamat - POST /api/nbfacin/risk-address: 201 {"id": "<ID baru>"}.
func simpanAlamat(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b alamatBaru
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, batasBadanAlamat)).Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan bukan JSON alamat risiko yang sah")
			return
		}
		id, err := svc.SimpanAlamatBaru(r.Context(), inti.PelakuDari(r, stubPelaku), models.AlamatBaru{NationName: b.NationName,
			ProvinceName: b.ProvinceName, DistrictName: b.DistrictName, CityName: b.CityName, TerritoryName: b.TerritoryName,
			Title: b.Title, Address: b.Address, PostalCode: b.PostalCode})
		if err != nil {
			tulisGalat(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			ID string `json:"id"`
		}{id})
	}
}
