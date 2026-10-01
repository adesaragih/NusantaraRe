package services_test

// PLAN LIST + UNDERWRITING LIMIT (paket 6, tiket 06).

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

// produkBerplan - produk tersimpan dengan satu plan ber-R/I Rate dari Pega.
const produkBerplan = `{"ID":"100007","PRODUCTNAME":"LAMA","CEDING":"UJI CEDING","CEDINGID":"L0UJI","SOBNAME":"UJI SOB",
	"SOBID":"L0SOB","PlanList":[{"Plan":"UJI COVER","PlanID":"P1","Name":"UJI BIZ","Benefit":"UJI MANFAAT",
	"RIRATE":"UJI RATE","RIRATEID":"R1"}]}`

func layananPlan() (*services.Layanan, func(models.Produk) error) {
	l, g := layananMaster()
	g.Umum["100007"] = produkBerplan
	g.Plan = []models.JenisPlan{{ID: "P1", CoverName: "UJI COVER", Business: "UJI BIZ", Benefit: "UJI MANFAAT"},
		{ID: "P2", CoverName: "UJI COVER DUA", Business: "UJI BIZ 2", Benefit: "UJI MANFAAT 2"}}
	return l, func(m models.Produk) error {
		m.ID = "100007"
		_, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
		return err
	}
}

var planTersimpan = models.BarisPlan{Plan: "UJI COVER", PlanID: "P1", Name: "UJI BIZ", Benefit: "UJI MANFAAT",
	RIRate: "UJI RATE", RIRateID: "R1"}

func TestPlanTersimpanDipertahankanUtuh(t *testing.T) {
	_, simpan := layananPlan()
	m := produkMasuk()
	m.PlanList = []models.BarisPlan{planTersimpan}
	if err := simpan(m); err != nil {
		t.Errorf("plan dan R/I Rate tersimpan tidak perlu diverifikasi ulang: %v", err)
	}
}

func TestPlanPesanVerbatimProteksiPlanListLife(t *testing.T) {
	_, simpan := layananPlan()
	m := produkMasuk()
	m.PlanList = []models.BarisPlan{{}, planTersimpan, planTersimpan}
	pesan := services.Pesan(simpan(m))
	for _, w := range []string{"PLAN LIST row 1: Plan tidak boleh kosong", "PLAN LIST row 1: RI/RATE tidak boleh kosong",
		"PLAN LIST row 3: Plan tidak boleh sama"} {
		if !strings.Contains(pesan, w) {
			t.Errorf("tanpa %q: %s", w, pesan)
		}
	}
	if strings.Contains(pesan, "row 2: Plan tidak boleh sama") {
		t.Errorf("baris pertama nilai itu bukan duplikat: %s", pesan)
	}
}

func TestPlanRIRateBaruMenungguOQ(t *testing.T) {
	_, simpan := layananPlan()
	m := produkMasuk()
	baru := planTersimpan
	baru.RIRateID, baru.RIRate = "R9", "RATE LAIN"
	m.PlanList = []models.BarisPlan{baru}
	err := simpan(m)
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), "OQ-MPNL-03") {
		t.Errorf("R/I Rate baru tidak dapat diverifikasi sampai OQ-MPNL-03: %v", err)
	}
}

func TestPlanMasterDiverifikasiNamaDariMaster(t *testing.T) {
	l, simpan := layananPlan()
	m := produkMasuk()
	baris := planTersimpan
	baris.PlanID, baris.Plan, baris.Name, baris.Benefit = "P2", "DIKETIK", "DIKETIK", "DIKETIK"
	m.PlanList = []models.BarisPlan{baris}
	if err := simpan(m); err != nil {
		t.Fatal(err)
	}
	p, _ := l.AmbilProduk(context.Background(), pelakuUji, "100007")
	if got := p.PlanList[0]; got.Plan != "UJI COVER DUA" || got.Name != "UJI BIZ 2" || got.Benefit != "UJI MANFAAT 2" ||
		got.RIRateID != "R1" {
		t.Errorf("BrowseProductTypeLife_RD: Plan ← CoverName, Name ← Business, Benefit ← Benefit: %+v", got)
	}
	baris.PlanID = "P-TIDAK-ADA"
	m.PlanList = []models.BarisPlan{baris}
	if err := simpan(m); !strings.Contains(services.Pesan(err), `PLAN LIST row 1: Plan Name "P-TIDAK-ADA" is not in the master list`) {
		t.Errorf("plan di luar master: %v", err)
	}
}

func TestUnderwritingLimitAngkaDanRentang(t *testing.T) {
	_, simpan := layananPlan()
	m := produkMasuk()
	m.UnderwritingLimit = []models.BarisUWLimit{
		{Description: "UJI", Medical: "BEBAS", MinAge: "70", MaxAge: "22", MinInsured: "5", MaxInsured: "x"},
		{Description: "UJI 2", Medical: "FCL", MinAge: "18", MaxAge: "65", MinInsured: "0", MaxInsured: "1175000000.5"},
	}
	pesan := services.Pesan(simpan(m))
	for _, w := range []string{"UNDERWRITING LIMIT row 1: Min Age must not be greater than Max Age",
		`UNDERWRITING LIMIT row 1: Max Insured "x" is not a number`} {
		if !strings.Contains(pesan, w) {
			t.Errorf("tanpa %q: %s", w, pesan)
		}
	}
	if strings.Contains(pesan, "Medical") || strings.Contains(pesan, "row 2") {
		t.Errorf("Medical teks bebas (R16), baris sah tidak ditolak: %s", pesan)
	}
}

func TestDaftarUrutDipertahankanDanKosongTetapKosong(t *testing.T) {
	l, simpan := layananPlan()
	m := produkMasuk()
	m.UnderwritingLimit = []models.BarisUWLimit{{Description: "C"}, {Description: "A"}, {Description: "B"}}
	if err := simpan(m); err != nil {
		t.Fatal(err)
	}
	p, _ := l.AmbilProduk(context.Background(), pelakuUji, "100007")
	if len(p.PlanList) != 0 || p.UnderwritingLimit[0].Description != "C" || p.UnderwritingLimit[2].Description != "B" {
		t.Errorf("urutan baris dipertahankan, daftar kosong = kosong: %+v %+v", p.PlanList, p.UnderwritingLimit)
	}
}
