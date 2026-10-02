package premium

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/utils"
)

// TestTotalFireLokasi - SumTotalTSIPremiGross_Act langkah 1: premi item = Σ premi
// coverage (1.3.2-1.3.3); daftar per mata uang lokasi dibangun ulang (1.2) dan
// dijumlahkan per item (1.3.4-1.3.5); Rate = Premium/TSI 20 desimal × 1000.
func TestTotalFireLokasi(t *testing.T) {
	got, err := TotalFireLokasi([]ItemProperti{
		{MataUang: "USD", TSIObjectItem: "1000", PremiCoverage: []string{"1.50", "0.25"}},
		{MataUang: "IDR", TSIObjectItem: "3", PremiCoverage: []string{"1"}},
		{MataUang: "USD", TSIObjectItem: "2000", PremiCoverage: []string{"0.75"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	for i, mau := range []string{"1.75", "1", "0.75"} {
		if s := utils.FormatDecimal(got.PremiItem[i].Amount); s != mau {
			t.Errorf("premi item %d: %s, mau %s", i, s, mau)
		}
	}
	// Urutan daftar = urutan kemunculan mata uang pertama.
	mau := []struct {
		uang, tsi, premi, rate string
		lewatDouble            bool
	}{
		// 2.50/3000 = 0.00083333333333333333 (20 desimal) × 1000.
		{"USD", "3000", "2.50", "0.83333333333333333000", true},
		// 1/3 = 0.33333333333333333333 (20 desimal) × 1000.
		{"IDR", "3", "1", "333.33333333333333333000", false},
	}
	if len(got.PerMataUang) != len(mau) {
		t.Fatalf("%d mata uang, mau %d", len(got.PerMataUang), len(mau))
	}
	for i, m := range mau {
		g := got.PerMataUang[i]
		if g.LewatDouble != m.lewatDouble {
			t.Errorf("%d: LewatDouble %v, mau %v", i, g.LewatDouble, m.lewatDouble)
		}
		if g.TSI.Currency != m.uang || utils.FormatDecimal(g.TSI.Amount) != m.tsi || utils.FormatDecimal(g.Premium.Amount) != m.premi || utils.FormatDecimal(g.Rate) != m.rate {
			t.Errorf("%d: %s %s %s %s, mau %+v", i, g.TSI.Currency, utils.FormatDecimal(g.TSI.Amount), utils.FormatDecimal(g.Premium.Amount), utils.FormatDecimal(g.Rate), m)
		}
	}
}

// TestTotalFireLokasiTepi - TSI 0 → Rate 0 (`@if(.TSI=0,0,…)`); seri pada @divide
// (bukan @Math.divide, modenya belum terverifikasi) ditolak; masukan tak sah ditolak.
func TestTotalFireLokasiTepi(t *testing.T) {
	got, err := TotalFireLokasi([]ItemProperti{{MataUang: "IDR", TSIObjectItem: "0", PremiCoverage: []string{"5"}}})
	if err != nil || !got.PerMataUang[0].Rate.IsZero() {
		t.Fatalf("TSI 0: %+v (%v)", got, err)
	}
	// 1 / 8 = 0.125 → ×10²⁰ tidak seri; 1/ (2×10²⁰) = 5e-21 → seri di desimal ke-21.
	if _, err := TotalFireLokasi([]ItemProperti{{MataUang: "IDR", TSIObjectItem: "200000000000000000000", PremiCoverage: []string{"1"}}}); !errors.Is(err, ErrSeriDivide) {
		t.Errorf("seri: galat %v", err)
	}
	for nama, it := range map[string]ItemProperti{
		"premi kosong": {MataUang: "IDR", TSIObjectItem: "1", PremiCoverage: []string{""}},
		"TSI kosong":   {MataUang: "IDR", TSIObjectItem: "", PremiCoverage: []string{"1"}},
	} {
		if _, err := TotalFireLokasi([]ItemProperti{it}); !errors.Is(err, ErrAngkaTakTerbaca) {
			t.Errorf("%s: galat %v", nama, err)
		}
	}
	if _, err := TotalFireLokasi([]ItemProperti{{TSIObjectItem: "1", PremiCoverage: []string{"1"}}}); !errors.Is(err, ErrMataUangKosong) {
		t.Errorf("mata uang kosong: galat %v", err)
	}
}
