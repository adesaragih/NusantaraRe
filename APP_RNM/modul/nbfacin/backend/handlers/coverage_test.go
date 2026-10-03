package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type coverageTiruan struct{ kata *string }

func (c coverageTiruan) CariCoverage(_ context.Context, k string) ([]models.BarisCoverage, error) {
	*c.kata = k
	return []models.BarisCoverage{{ID: "100815", OldID: "UJI-OLD", Nama: "UJI COVERAGE"}}, nil
}

func (coverageTiruan) CoverageOtomatis(context.Context) ([]models.BarisCoverage, error) {
	return []models.BarisCoverage{{ID: "100815", OldID: "UJI-OLD", Nama: "UJI COVERAGE"}, {ID: "100828"}}, nil
}

// kasusPeriode - case UJI dengan periode polis setahun penuh (Prorate 100).
func kasusPeriode() *models.Kasus {
	return &models.Kasus{CaseID: "UJI-NB-1", InsuredName: "UJI",
		General: models.General{StartDateTime: "20250101T050000.000 GMT", EndDateTime: "20260101T050000.000 GMT"}}
}

// kunciCoverage - 24 kunci kontrak CoverageObjek (frontend/api.ts).
var kunciCoverage = []string{"coverage", "oldId", "coverageNote", "coverageBasis", "day", "tsi", "indemnity", "rate", "rateOjk",
	"firstLoss", "discountPercentage", "tsiLiability", "netRate", "limitOfLiability", "pctLol", "proRatePercent",
	"indemnityPercentage", "firstScale", "sublimit", "lostLimit", "emlPml", "discount", "premium", "conditions"}

// TestCariCoverage - GET /api/nbfacin/coverage (tiket 43): bentuk persis BarisCoverage, kata di-trim, kosong = semua;
// tanpa identitas; 400 / 503.
func TestCariCoverage(t *testing.T) {
	var kata string
	svc := services.Baru(nil).DenganCoverage(coverageTiruan{&kata})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/coverage?cari=%20fire%20", "", ""); kode != 200 ||
		isi != `{"baris":[{"id":"100815","oldId":"UJI-OLD","nama":"UJI COVERAGE"}]}` || kata != "fire" {
		t.Errorf("%d %s kata %q", kode, isi, kata)
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/coverage", "", ""); kode != 200 || kata != "" {
		t.Errorf("kosong: %d kata %q", kode, kata)
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/coverage?cari="+strings.Repeat("A", 256), "", ""); kode != 400 {
		t.Errorf("400: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/coverage", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestCoverageOtomatis - GET /api/nbfacin/coverage-otomatis: urutan pembaca dipertahankan; 503.
func TestCoverageOtomatis(t *testing.T) {
	var kata string
	if kode, isi := minta(t, services.Baru(nil).DenganCoverage(coverageTiruan{&kata}), "GET", "/api/nbfacin/coverage-otomatis", "", ""); kode != 200 ||
		isi != `{"baris":[{"id":"100815","oldId":"UJI-OLD","nama":"UJI COVERAGE"},{"id":"100828","oldId":"","nama":""}]}` {
		t.Errorf("%d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/coverage-otomatis", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestHitungCoverageHandler - POST …/hitung-coverage: 24 kunci persis CoverageObjek; medan server (tsi, tsiLiability,
// proRatePercent) dari badan diabaikan; 400 desimal ber-jalur / basis 5 / mode; 409 periode kosong; 503.
func TestHitungCoverageHandler(t *testing.T) {
	svc := services.Baru(nil).DenganKasus(kasusTiruan{k: kasusPeriode()})
	jalur := "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage"
	cov := `{"coverage":"100815","coverageBasis":"1","rate":"1","indemnityPercentage":"100","tsi":"9","tsiLiability":"9","proRatePercent":"9"}`
	kode, isi := minta(t, svc, "POST", jalur, `{"coverage":`+cov+`,"tsi":"1000000000","mode":"percent"}`, "")
	var j map[string]any
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || len(j) != len(kunciCoverage) {
		t.Fatalf("%d %s", kode, isi)
	}
	for _, k := range kunciCoverage {
		if _, ada := j[k]; !ada {
			t.Errorf("kunci %q tidak ada (kontrak CoverageObjek)", k)
		}
	}
	if j["premium"] != "1000000" || j["tsi"] != "1000000000" || j["tsiLiability"] != "1000000000" || j["proRatePercent"] != "100" {
		t.Errorf("hasil: %s", isi)
	}
	for _, u := range []struct {
		nama, badan, pesan string
		svc                *services.Service
		kode               int
	}{
		{"JSON rusak", `{"coverage":`, "bukan JSON", svc, 400},
		{"rate koma", `{"coverage":{"coverageBasis":"1","rate":"1,5"},"tsi":"1","mode":"percent"}`, "coverage.rate harus angka", svc, 400},
		{"tsi minus", `{"coverage":{"coverageBasis":"1"},"tsi":"-1","mode":"percent"}`, "tsi harus angka", svc, 400},
		{"pctAdjustOther", `{"coverage":{"coverageBasis":"1"},"tsi":"1","mode":"percent","pctAdjustOther":"x"}`, "pctAdjustOther harus angka", svc, 400},
		{"basis 5", `{"coverage":{"coverageBasis":"5"},"tsi":"1","mode":"percent"}`, "Layering belum didukung", svc, 400},
		{"mode", `{"coverage":{"coverageBasis":"1"},"tsi":"1","mode":"x"}`, "mode harus", svc, 400},
		{"modeDiskon", `{"coverage":{"coverageBasis":"1"},"tsi":"1","mode":"percent","modeDiskon":"x"}`, "modeDiskon", svc, 400},
		{"periode kosong", `{"coverage":{"coverageBasis":"1"},"tsi":"1","mode":"percent"}`, "Begin / End",
			services.Baru(nil).DenganKasus(kasusTiruan{k: &models.Kasus{CaseID: "UJI-NB-1"}}), 409},
		{"tanpa DB", `{"coverage":{"coverageBasis":"1"},"tsi":"1","mode":"percent"}`, "galat", services.Baru(nil), 503},
	} {
		if kode, isi := minta(t, u.svc, "POST", jalur, u.badan, ""); kode != u.kode || !strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s, mau %d", u.nama, kode, isi, u.kode)
		}
	}
}

// TestObjekCoverage - PUT/GET objek dengan items[].coverages (tiket 43): premi dihitung ulang server, totalGrossPremi
// BACA-SAJA, totalNetRate pulang-pergi, totalPerCurrency dihitung; 400 ber-jalur coverages[k]; mata uang > 10 byte.
func TestObjekCoverage(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle).DenganPilihanItem(pilihanTiruan{}).
		DenganKasus(kasusTiruan{k: kasusPeriode()})
	jalur := "/api/nbfacin/kasus/UJI-NB-1/objek"
	cov := `{"coverage":"100815","coverageBasis":"1","rate":"1","indemnityPercentage":"100","premium":"5"}`
	badan := `{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","tsi":"1000000000","totalGrossPremi":"7","totalNetRate":"1.25",` +
		`"coverages":[` + cov + `]}],"totalPerCurrency":[{"currency":"X"}]}]}`
	kode, isi := minta(t, svc, "PUT", jalur, badan, "UJI-USER")
	var j struct {
		Baris []struct {
			Items []struct {
				Coverages       []map[string]any `json:"coverages"`
				TotalGrossPremi string           `json:"totalGrossPremi"`
				TotalNetRate    string           `json:"totalNetRate"`
			} `json:"items"`
			TotalPerCurrency []map[string]any `json:"totalPerCurrency"`
		} `json:"baris"`
	}
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || len(j.Baris) != 1 || len(j.Baris[0].Items) != 1 ||
		len(j.Baris[0].Items[0].Coverages) != 1 {
		t.Fatalf("%d %s", kode, isi)
	}
	it := j.Baris[0].Items[0]
	if len(it.Coverages[0]) != len(kunciCoverage) || it.Coverages[0]["premium"] != "1000000" || it.TotalGrossPremi != "1000000" ||
		it.TotalNetRate != "1.25" {
		t.Errorf("item: %s", isi)
	}
	tot := j.Baris[0].TotalPerCurrency
	if len(tot) != 1 || tot[0]["currency"] != "IDR" || tot[0]["tsi"] != "1000000000" || tot[0]["premium"] != "1000000" || tot[0]["rate"] != "1" {
		t.Errorf("totalPerCurrency: %s", isi)
	}
	for nama, u := range map[string]struct {
		badan, pesan string
	}{
		"rate koma": {`{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","coverages":[{"coverageBasis":"1"},{"coverageBasis":"1","rate":"1,5"}]}]}]}`,
			"baris[0].items[0].coverages[1].rate harus angka"},
		"basis 5":       {`{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","coverages":[{"coverageBasis":"5"}]}]}]}`, "Layering belum didukung"},
		"totalNetRate":  {`{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","totalNetRate":"x"}]}]}`, "baris[0].items[0].totalNetRate"},
		"mata uang 11B": {`{"baris":[{"objectType":"UJI","items":[{"currency":"IDRIDRIDRID","coverages":[{"coverageBasis":"1"}]}]}]}`, "currency paling banyak 10 byte"},
	} {
		if kode, isi := minta(t, svc, "PUT", jalur, u.badan, "UJI-USER"); kode != 400 || !strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
}
