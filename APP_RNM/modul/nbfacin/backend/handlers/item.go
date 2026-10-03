package handlers

// Pilihan sub-tab Object Item (tiket 39): GET /api/nbfacin/jenis-item-objek dan GET /api/nbfacin/mata-uang.

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

// jenisItem - kontrak `JenisItem` frontend.
type jenisItem struct {
	Kode       string `json:"kode"`
	Nama       string `json:"nama"`
	Keterangan string `json:"keterangan"`
}

// daftarJenisItem - GET /api/nbfacin/jenis-item-objek.
func daftarJenisItem(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.DaftarJenisItem(r.Context())
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]jenisItem, 0, len(hasil))
		for _, j := range hasil {
			baris = append(baris, jenisItem{Kode: j.Kode, Nama: j.Nama, Keterangan: j.Keterangan})
		}
		galat.TulisJSON(w, struct {
			Baris []jenisItem `json:"baris"`
		}{baris})
	}
}

// daftarMataUang - GET /api/nbfacin/mata-uang: `{"baris": ["IDR", ...]}`.
func daftarMataUang(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hasil, err := svc.DaftarMataUang(r.Context())
		if err != nil {
			tulisGalat(w, err)
			return
		}
		galat.TulisJSON(w, struct {
			Baris []string `json:"baris"`
		}{append([]string{}, hasil...)})
	}
}
