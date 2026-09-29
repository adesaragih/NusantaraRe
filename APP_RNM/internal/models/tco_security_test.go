package models

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestPeriksaSecurityTCO(t *testing.T) {
	sah := SecurityReinsurer{ReasID: "1000001", ReasSecurity: "UJI-AG1"}
	if err := PeriksaSecurityTCO(sah); err != nil {
		t.Fatalf("baris sah ditolak: %v", err)
	}
	for _, k := range []struct {
		nama string
		s    SecurityReinsurer
		mau  error
	}{
		{"tanpa reinsurer", SecurityReinsurer{ReasSecurity: "UJI-AG1"}, ErrSecurityTanpaReinsurer},
		{"nama kosong", SecurityReinsurer{ReasID: "1000001", ReasSecurity: "  "}, ErrSecurityKosong},
	} {
		err := PeriksaSecurityTCO(k.s)
		if !errors.Is(err, k.mau) {
			t.Errorf("%s: %v, mau %v", k.nama, err, k.mau)
		}
	}
	// Pesan menyebut label form b19648 - bukan nama kolom.
	if err := PeriksaSecurityTCO(SecurityReinsurer{ReasID: "1"}); err == nil || !strings.Contains(err.Error(), "Security Name") {
		t.Errorf("pesan tidak menyebut medan: %v", err)
	}
}

// AC 17: security dicatat BESERTA porsinya - `%Share` wajib dan 0..100,
// desimal persis sampai 8 angka.
func TestShareSecurityTCO(t *testing.T) {
	d, err := UraiShareSecurityTCO("33,33333333")
	if err != nil || d.Text('f') != "33.33333333" {
		t.Fatalf("share: %v %v", d, err)
	}
	for teks, mau := range map[string]error{
		"":         ErrPersenKosong,
		"100.0001": ErrPersenDiLuarRentang,
		"-1":       ErrPersenDiLuarRentang,
		"1,5.2":    ErrPersenBukanDesimal,
		"abc":      ErrPersenBukanDesimal,
	} {
		if _, err := UraiShareSecurityTCO(teks); !errors.Is(err, mau) {
			t.Errorf("%q: %v, mau %v", teks, err, mau)
		}
	}
	if _, err := UraiShareSecurityTCO(""); err == nil || !strings.Contains(err.Error(), LabelShareSecurityTCO) {
		t.Errorf("pesan tidak menyebut label %%Share: %v", err)
	}
}

// ADR-0003: persen security tidak pernah menjadi float di jalur modul ini.
func TestSecurityTanpaFloat(t *testing.T) {
	for _, berkas := range []string{"tco_security.go", "../repository/tco_security.go", "../services/tco_security.go"} {
		isi, err := os.ReadFile(berkas)
		if err != nil {
			t.Fatal(err)
		}
		for _, larang := range []string{"float32", "float64", "ParseFloat", "strconv.Format"} {
			if strings.Contains(string(isi), larang) {
				t.Errorf("%s memuat %s", berkas, larang)
			}
		}
	}
}
