package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type rwTiruan struct{}

func (rwTiruan) CariRW(_ context.Context, z string) ([]models.BarisRW, error) {
	return []models.BarisRW{{ZipCode: z + "45", TerritoryName: "UJI WIL", DistrictName: "UJI KEC", CityName: "UJI KOTA",
		ProvinceName: "UJI PROV", NationName: "UJI NEGARA"}}, nil
}

func (rwTiruan) SisipRiskAddress(context.Context, *db.Tx, models.AlamatBaru) (string, error) {
	return "UJI000000000001", nil
}

// TestRW - GET /api/nbfacin/rw dan POST /api/nbfacin/risk-address (tiket 37): bentuk jawaban
// persis kontrak api.ts; 201 {"id"}; 400/401/503; GET risk-address (tiket 36) tidak tertimpa.
func TestRW(t *testing.T) {
	svc := services.Baru(nil).DenganRW(rwTiruan{}).DenganTransaksi(tanpaOracle)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/rw?zipCode=123", "", "")
	mau := `{"baris":[{"zipCode":"12345","territoryName":"UJI WIL","districtName":"UJI KEC","cityName":"UJI KOTA","provinceName":"UJI PROV","nationName":"UJI NEGARA"}]}`
	if kode != 200 || isi != mau {
		t.Fatalf("%d %s\nmau %s", kode, isi, mau)
	}
	badan := `{"nationName":"","provinceName":"","districtName":"","cityName":"","territoryName":"","title":"DESA","address":"UJI JALAN","postalCode":"00000"}`
	if kode, isi := minta(t, svc, "POST", "/api/nbfacin/risk-address", badan, "UJI-USER"); kode != 201 || isi != `{"id":"UJI000000000001"}` {
		t.Errorf("simpan: %d %s", kode, isi)
	}
	for _, u := range []struct {
		metode, jalur, badan, pelaku string
		svc                          *services.Service
		kode                         int
	}{
		{"GET", "/api/nbfacin/rw?zipCode=12", "", "", svc, 400},
		{"GET", "/api/nbfacin/rw?zipCode=123", "", "", services.Baru(nil), 503},
		{"POST", "/api/nbfacin/risk-address", badan, "", svc, 401},
		{"POST", "/api/nbfacin/risk-address", `{"address":"x"}`, "UJI-USER", svc, 400},
		{"POST", "/api/nbfacin/risk-address", `{"address":`, "UJI-USER", svc, 400},
		{"POST", "/api/nbfacin/risk-address", badan, "UJI-USER", services.Baru(nil), 503},
		{"GET", "/api/nbfacin/risk-address", "", "", svc, 400}, // tiket 36 tetap: tanpa saringan
	} {
		if kode, isi := minta(t, u.svc, u.metode, u.jalur, u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s %s: %d %s, mau %d", u.metode, u.jalur, kode, isi, u.kode)
		}
	}
}
