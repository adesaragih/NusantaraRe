package handlers

// Popup Choose Accumulation Code (tiket 46): GET /api/nbfacin/akumulasi dan GET /api/nbfacin/akumulasi/saran/{jenis};
// form Add New: POST /api/nbfacin/akumulasi, GET /api/nbfacin/akumulasi/czone, GET /api/nbfacin/akumulasi/zipcode.
// Kontrak frontend modul/nbfacin/frontend/api.ts.

import (
	"encoding/json"
	"net/http"
	"strconv"

	inti "nusantarare/inti/backend"
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

// tambahAkumulasi - POST /api/nbfacin/akumulasi (form Add New) badan {accumulation, accumulationType, note, keyword,
// scopeArea, cZone, cZoneId, provinceId, zipCode, negara} (teks; kunci asing -> 400) -> 201 {id, note}.
func tambahAkumulasi(svc *services.Service, stubPelaku bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b struct {
			Accumulation     string `json:"accumulation"`
			AccumulationType string `json:"accumulationType"`
			Note             string `json:"note"`
			Keyword          string `json:"keyword"`
			ScopeArea        string `json:"scopeArea"`
			CZone            string `json:"cZone"`
			CZoneID          string `json:"cZoneId"`
			ProvinceID       string `json:"provinceId"`
			ZipCode          string `json:"zipCode"`
			Negara           string `json:"negara"`
		}
		dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&b); err != nil {
			galat.Tulis(w, http.StatusBadRequest, "badan permintaan harus objek JSON akumulasi bernilai teks")
			return
		}
		id, note, err := svc.TambahAkumulasi(r.Context(), inti.PelakuDari(r, stubPelaku), services.IsianAkumulasi(b))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			ID   string `json:"id"`
			Note string `json:"note"`
		}{id, note})
	}
}

// czoneZip - GET /api/nbfacin/akumulasi/czone?zip= -> {cZoneId, cZone} (kosong bila tidak ada).
func czoneZip(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := svc.CZoneZip(r.Context(), r.URL.Query().Get("zip"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			CZoneID string `json:"cZoneId"`
			CZone   string `json:"cZone"`
		}{c.ID, c.Code})
	}
}

// zipAkumulasiKabel - kontrak `{ zipCode, city, province, nation, nationInitial }`.
type zipAkumulasiKabel struct {
	ZipCode       string `json:"zipCode"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Nation        string `json:"nation"`
	NationInitial string `json:"nationInitial"`
}

// cariZipAkumulasi - GET /api/nbfacin/akumulasi/zipcode?provinceName=&q=&halaman= -> {baris, total, halaman, ukuran}.
func cariZipAkumulasi(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		nomor := 1
		if v := q.Get("halaman"); v != "" {
			n, err := strconv.Atoi(v)
			if err != nil {
				galat.Tulis(w, http.StatusBadRequest, "halaman harus bilangan bulat")
				return
			}
			nomor = n
		}
		h, err := svc.CariZipAkumulasi(r.Context(), q.Get("provinceName"), q.Get("q"), nomor)
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]zipAkumulasiKabel, 0, len(h.Baris))
		for _, b := range h.Baris {
			baris = append(baris, zipAkumulasiKabel(b))
		}
		galat.TulisJSON(w, struct {
			Baris   []zipAkumulasiKabel `json:"baris"`
			Total   int                 `json:"total"`
			Halaman int                 `json:"halaman"`
			Ukuran  int                 `json:"ukuran"`
		}{baris, h.Total, h.Nomor, h.Ukuran})
	}
}
