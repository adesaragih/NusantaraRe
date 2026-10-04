package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// akumulasiTiruan - mencatat jalur yang dipanggil.
type akumulasiTiruan struct{ jalur *string }

func (a akumulasiTiruan) CariAkumulasiRD(_ context.Context, s models.SaringAkumulasi) ([]models.BarisAkumulasi, error) {
	*a.jalur = "rd:" + s.ID + "|" + s.Note + "|" + s.PostalCode + "|" + s.ProvinceID + "|" + s.CZone + "|" + s.Keyword
	return []models.BarisAkumulasi{{ID: "UJI-1", AccumulationName: "UJI TIPE", Note: "UJI JALAN"}}, nil
}

func (a akumulasiTiruan) AkumulasiWilayah(_ context.Context, kecamatan bool, id string) ([]models.BarisAkumulasi, error) {
	*a.jalur = "kota:" + id
	if kecamatan {
		*a.jalur = "kecamatan:" + id
	}
	return nil, nil
}

func (a akumulasiTiruan) AkumulasiPolis(_ context.Context, noPolis string) ([]models.BarisAkumulasi, error) {
	*a.jalur = "polis:" + noPolis
	return nil, nil
}

func (a akumulasiTiruan) Saran(_ context.Context, jenis, kata, induk string) ([]models.SaranAkumulasi, error) {
	*a.jalur = jenis + ":" + kata + "|" + induk
	return []models.SaranAkumulasi{{ID: "", Label: "UJI AREA", Ekstra: "12345"}, {ID: "C1", Label: "UJI KOTA"}}, nil
}

// TestCariAkumulasi - tiket 46 butir 2: bentuk {baris:[{id,accumulationName,note}]}; pemilihan jalur
// GetDataAccumulation_act (postalCode mengalahkan jalur SQL; kecamatan > kota; polis + wilayah = kosong; syariahStatus
// diabaikan); nilai dipangkas; 400 / 503.
func TestCariAkumulasi(t *testing.T) {
	var jalur string
	svc := services.Baru(nil).DenganAkumulasi(akumulasiTiruan{&jalur})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi?id=%20UJI-1%20&note=jl&czone=Z&keyword=K&provinceId=P&syariahStatus=S", "", "")
	if kode != 200 || isi != `{"baris":[{"id":"UJI-1","accumulationName":"UJI TIPE","note":"UJI JALAN"}]}` || jalur != "rd:UJI-1|jl||P|Z|K" {
		t.Fatalf("%d %s jalur %q", kode, isi, jalur)
	}
	for kueri, mau := range map[string]string{
		"cityId=C1":                       "kota:C1",
		"cityId=C1&districtId=D1":         "kecamatan:D1",
		"districtId=D1":                   "kecamatan:D1",
		"policyNo=POL-1":                  "polis:POL-1",
		"cityId=C1&postalCode=12345":      "rd:||12345|||",
		"policyNo=POL-1&postalCode=12345": "rd:||12345|||",
		"syariahStatus=S":                 "rd:|||||",
		"":                                "rd:|||||",
	} {
		jalur = ""
		if kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi?"+kueri, "", ""); kode != 200 || jalur != mau || !strings.HasPrefix(isi, `{"baris":[`) {
			t.Errorf("%q: %d %s jalur %q, mau %q", kueri, kode, isi, jalur, mau)
		}
	}
	jalur = ""
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi?policyNo=POL-1&cityId=C1", "", ""); kode != 200 || isi != `{"baris":[]}` || jalur != "" {
		t.Errorf("polis + kota: %d %s %q", kode, isi, jalur)
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/akumulasi?note="+strings.Repeat("A", 256), "", ""); kode != 400 {
		t.Errorf("400: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/akumulasi", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestSaranAkumulasi - tiket 46 butir 3: bentuk {baris:[{id,label,ekstra?}]}; ketujuh jenis dilayani (A178 diganti DDL
// 04-10-2026); tak dikenal -> 400.
func TestSaranAkumulasi(t *testing.T) {
	var jalur string
	svc := services.Baru(nil).DenganAkumulasi(akumulasiTiruan{&jalur})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/akumulasi/saran/area?q=%20jl%20&induk=UJI%20KEC", "", "")
	if kode != 200 || isi != `{"baris":[{"id":"","label":"UJI AREA","ekstra":"12345"},{"id":"C1","label":"UJI KOTA"}]}` || jalur != "area:jl|UJI KEC" {
		t.Fatalf("%d %s %q", kode, isi, jalur)
	}
	for _, j := range []string{"nation", "province", "accumtype", "czone", "city", "district"} {
		if kode, _ := minta(t, svc, "GET", "/api/nbfacin/akumulasi/saran/"+j+"?q=x", "", ""); kode != 200 || jalur != j+":x|" {
			t.Errorf("%s: %d %q", j, kode, jalur)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/akumulasi/saran/x", "", ""); kode != 400 {
		t.Errorf("tak dikenal: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/akumulasi/saran/city", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}
