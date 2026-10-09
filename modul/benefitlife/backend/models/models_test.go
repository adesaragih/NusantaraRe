package models

import (
	"errors"
	"strings"
	"testing"
)

// K3: rumus prosedur PEGA_M_BENEFIT_LIFE - '1' || LPAD(seq, 5, '0'); nomor lebih dari 5 angka DITOLAK (LPAD memotong).
func TestBentukIDRumusProsedur(t *testing.T) {
	for nomor, mau := range map[string]string{"1": "100001", "12": "100012", " 99999 ": "199999", "4": "100004"} {
		if got, err := BentukID(nomor); err != nil || got != mau {
			t.Errorf("BentukID(%q) = %q, %v; mau %q", nomor, got, err, mau)
		}
	}
	for _, nomor := range []string{"123456", "", "1a", "-1"} {
		if _, err := BentukID(nomor); !errors.Is(err, ErrIDTidakSah) {
			t.Errorf("BentukID(%q) mau ErrIDTidakSah, dapat %v", nomor, err)
		}
	}
	if id, _ := BentukID("99999"); len(id) > BatasID {
		t.Errorf("ID %s melewati VARCHAR2(%d)", id, BatasID)
	}
}

// K6: Benefit wajib (b1137), huruf besar (SetUpperCase_DT b1211), dipangkas, paling panjang 200 byte (lebar kolom).
func TestNormalBenefit(t *testing.T) {
	if got, err := NormalBenefit("  uji Rawat Inap  "); err != nil || got != "UJI RAWAT INAP" {
		t.Errorf("rapikan %q %v", got, err)
	}
	for _, s := range []string{"", "   ", "\t\n"} {
		if _, err := NormalBenefit(s); err == nil || err.Error() != "Benefit is required" {
			t.Errorf("kosong %q: %v", s, err)
		}
	}
	if _, err := NormalBenefit(strings.Repeat("A", BatasBenefit)); err != nil {
		t.Errorf("tepat %d byte ditolak: %v", BatasBenefit, err)
	}
	if _, err := NormalBenefit(strings.Repeat("A", BatasBenefit+1)); err == nil || !strings.Contains(err.Error(), "longer than 200") {
		t.Errorf("lebih dari %d byte: %v", BatasBenefit, err)
	}
	// Byte, bukan karakter: 67 huruf 3-byte = 201 byte.
	if _, err := NormalBenefit(strings.Repeat("€", 67)); err == nil {
		t.Error("201 byte UTF-8 diterima")
	}
}

// pyPageSize 10 (b4526).
func TestUkuranHalaman(t *testing.T) {
	if UkuranHalaman != 10 {
		t.Errorf("ukuran %d", UkuranHalaman)
	}
}
