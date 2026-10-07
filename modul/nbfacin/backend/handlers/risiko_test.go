package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type riskTiruan struct{ saring *models.SaringRisk }

func (r riskTiruan) CariRisk(_ context.Context, s models.SaringRisk, _, _ int) ([]models.RiskAddress, int, error) {
	*r.saring = s
	return []models.RiskAddress{{ID: "UJI-R1", Title: "JL.", Address: "UJI JALAN", NationName: "UJI NEGARA", ProvinceName: "UJI PROV",
		CityName: "UJI KOTA", DistrictName: "UJI KEC", TerritoryName: "UJI WIL", PostalCode: "00000"}}, 1, nil
}

// TestCariRisk - GET /api/nbfacin/risk-address (tiket 36): bentuk jawaban persis kontrak
// frontend, nama saringan kueri = api.ts, 400 tanpa saringan, 503 tanpa DB.
func TestCariRisk(t *testing.T) {
	var s models.SaringRisk
	svc := services.Baru(nil).DenganRisk(riskTiruan{&s})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/risk-address?address=a&zipCode=b&country=c&province=d&city=e&district=f&territory=g&halaman=1", "", "")
	mau := `{"baris":[{"id":"UJI-R1","title":"JL.","address":"UJI JALAN","nationName":"UJI NEGARA","provinceName":"UJI PROV",` +
		`"cityName":"UJI KOTA","districtName":"UJI KEC","territoryName":"UJI WIL","postalCode":"00000"}],"total":1,"halaman":1,"ukuran":10}`
	if kode != 200 || isi != mau {
		t.Fatalf("%d %s\nmau %s", kode, isi, mau)
	}
	if s != (models.SaringRisk{Address: "a", ZipCode: "b", Country: "c", Province: "d", City: "e", District: "f", Territory: "g"}) {
		t.Errorf("saringan %+v", s)
	}
	for _, u := range []struct {
		jalur string
		svc   *services.Service
		kode  int
	}{{"/api/nbfacin/risk-address", svc, 400}, {"/api/nbfacin/risk-address?city=x&halaman=satu", svc, 400},
		{"/api/nbfacin/risk-address?city=x", services.Baru(nil), 503}} {
		if kode, isi := minta(t, u.svc, "GET", u.jalur, "", ""); kode != u.kode || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s: %d %s", u.jalur, kode, isi)
		}
	}
}
