package handlers

// Lookup form objek Aneka (tiket 73 P1): GET /api/nbfacin/kapal?cari= dan GET /api/nbfacin/tipe-kendaraan - kontrak api.ts
// `BarisKapal` / `BarisTipeKendaraan`. Semua teks apa adanya (kolom sumber VARCHAR2).

import (
	"net/http"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbfacin/backend/services"
)

type kapalKabel struct {
	ID          string `json:"id"`
	Nama        string `json:"nama"`
	Flag        string `json:"flag"`
	Tahun       string `json:"tahun"`
	NRT         string `json:"nrt"`
	GRT         string `json:"grt"`
	DWT         string `json:"dwt"`
	Klasifikasi string `json:"klasifikasi"`
	Konstruksi  string `json:"konstruksi"`
}

type tipeKendaraanKabel struct {
	ID       string `json:"id"`
	TypeName string `json:"typeName"`
}

// cariKapal - GET /api/nbfacin/kapal?cari=.
func cariKapal(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.CariKapal(r.Context(), r.URL.Query().Get("cari"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]kapalKabel, 0, len(d))
		for _, k := range d {
			baris = append(baris, kapalKabel(k))
		}
		galat.TulisJSON(w, struct {
			Baris []kapalKabel `json:"baris"`
		}{baris})
	}
}

// daftarTipeKendaraan - GET /api/nbfacin/tipe-kendaraan.
func daftarTipeKendaraan(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.DaftarTipeKendaraan(r.Context())
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]tipeKendaraanKabel, 0, len(d))
		for _, k := range d {
			baris = append(baris, tipeKendaraanKabel(k))
		}
		galat.TulisJSON(w, struct {
			Baris []tipeKendaraanKabel `json:"baris"`
		}{baris})
	}
}

type okupasiAnekaPilihanKabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// cariOkupasiAneka - GET /api/nbfacin/kasus/{caseId}/occupation-aneka?q= (popup Choose Occupation Aneka). Pilih: id ->
// okupasi.occupationId, name -> okupasi.occupationName (objek-aneka).
func cariOkupasiAneka(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.CariOkupasiAneka(r.Context(), r.PathValue("caseId"), r.URL.Query().Get("q"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		baris := make([]okupasiAnekaPilihanKabel, 0, len(d))
		for _, k := range d {
			baris = append(baris, okupasiAnekaPilihanKabel(k))
		}
		galat.TulisJSON(w, struct {
			Baris []okupasiAnekaPilihanKabel `json:"baris"`
		}{baris})
	}
}

// cariCoverageAneka - GET /api/nbfacin/kasus/{caseId}/coverage-aneka?q= (Coverage Note tab Coverage Aneka) - bentuk BarisCoverage
// FIRE. Pilih: id -> coverage, nama -> coverageNote, oldId -> oldId.
func cariCoverageAneka(svc *services.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d, err := svc.CariCoverageAneka(r.Context(), r.PathValue("caseId"), r.URL.Query().Get("q"))
		if err != nil {
			tulisGalat(w, err)
			return
		}
		tulisBarisCoverage(w, d)
	}
}
