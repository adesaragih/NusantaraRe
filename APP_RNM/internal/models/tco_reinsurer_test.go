package models

import (
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

func des(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, _, err := apd.NewFromString(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// Normalisasi koma -> titik SEKALI di batas masukan (SetErrorMessageReinsurer
// langkah 1), rentang 0..100 (langkah 2-3, `SetErrorMessageBetween`), dan batas
// NUMBER(38,8).
func TestUraiPersenMasukTCO(t *testing.T) {
	sah := map[string]string{
		"12,5": "12.5", " 7.25 ": "7.25", "0": "0", "100": "100", "33,33333333": "33.33333333",
	}
	for masuk, mau := range sah {
		d, err := UraiPersenMasukTCO("PctShare", masuk)
		if err != nil || d.Text('f') != mau {
			t.Errorf("UraiPersenMasukTCO(%q) = %v %v, mau %s", masuk, d, err, mau)
		}
	}
	tolak := map[string]error{
		"":            ErrPersenKosong,
		"1.234,5":     ErrPersenBukanDesimal,
		"dua":         ErrPersenBukanDesimal,
		"100.0000001": ErrPersenDiLuarRentang,
		"-0.5":        ErrPersenDiLuarRentang,
		"1.123456789": ErrPersenBukanDesimal,
	}
	for masuk, mau := range tolak {
		_, err := UraiPersenMasukTCO("Ricomm", masuk)
		if !errors.Is(err, mau) {
			t.Errorf("UraiPersenMasukTCO(%q) = %v, mau %v", masuk, err, mau)
			continue
		}
		if !strings.Contains(err.Error(), "Ricomm") {
			t.Errorf("pesan %q tidak menyebut medannya", err)
		}
	}
}

// SaveTreatyReinsurerDetail1_Act langkah 7 + 9: tolak bila total baru > 100.
// Tiga sepertiga yang jumlahnya TEPAT 100 harus lolos - float64 tidak dapat
// menjaminnya (ADR-0003).
func TestPeriksaTotalShareTCO(t *testing.T) {
	lain := []*apd.Decimal{des(t, "33.33333333"), des(t, "33.33333333")}
	total, err := PeriksaTotalShareTCO(lain, des(t, "33.33333334"))
	if err != nil || total.Text('f') != "100.00000000" {
		t.Fatalf("tepat 100 ditolak: total %v, galat %v", total, err)
	}
	_, err = PeriksaTotalShareTCO(lain, des(t, "33.33333335"))
	if !errors.Is(err, ErrTotalShareMelebihi100) || !strings.Contains(err.Error(), "100.00000001") {
		t.Errorf("100,00000001 lolos atau pesan tanpa totalnya: %v", err)
	}
	// Kurang dari 100 BUKAN gerbang (hanya > 100 yang ditolak di Pega).
	if total, err := PeriksaTotalShareTCO(nil, des(t, "40")); err != nil || total.Text('f') != "40" {
		t.Errorf("40 ditolak: %v %v", total, err)
	}
	// Sel NULL warisan dihitung nol, bukan menggagalkan.
	if total, err := PeriksaTotalShareTCO([]*apd.Decimal{nil, des(t, "60")}, des(t, "40")); err != nil || total.Text('f') != "100" {
		t.Errorf("NULL warisan: %v %v", total, err)
	}
}

func TestTotalShareTCO(t *testing.T) {
	total, err := TotalShareTCO([]*apd.Decimal{des(t, "0.1"), des(t, "0.2"), nil})
	if err != nil || total.Text('f') != "0.3" {
		t.Errorf("0.1 + 0.2 = %v %v, mau tepat 0.3", total, err)
	}
	if total, _ := TotalShareTCO(nil); total.Text('f') != "0" {
		t.Errorf("kosong = %v", total)
	}
}

func TestPeriksaReinsurerTCO(t *testing.T) {
	r := ReinsurerTreaty{ReinsurerID: " ", PctShare: des(t, "10"), Ricomm: des(t, "5")}
	if err := PeriksaReinsurerTCO(r); !errors.Is(err, ErrReinsurerKosong) ||
		!strings.Contains(err.Error(), "Data tidak boleh kosong...!!!") {
		t.Errorf("reinsurer kosong: %v", err)
	}
	r.ReinsurerID = "UJI-R1"
	if err := PeriksaReinsurerTCO(r); err != nil {
		t.Errorf("wajar: %v", err)
	}
}
