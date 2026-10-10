package models

import (
	"errors"
	"strings"
	"testing"
)

// D1.1: ID = TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL) apa adanya (DEV 197586); bukan angka / lebih dari VARCHAR2(100) DITOLAK.
func TestBentukIDDariSequence(t *testing.T) {
	for nomor, mau := range map[string]string{"197586": "197586", " 1 ": "1", "100006": "100006"} {
		if got, err := BentukID(nomor); err != nil || got != mau {
			t.Errorf("BentukID(%q) = %q, %v; mau %q", nomor, got, err, mau)
		}
	}
	for _, nomor := range []string{"", "1a", "-1", strings.Repeat("9", BatasID+1)} {
		if _, err := BentukID(nomor); !errors.Is(err, ErrIDTidakSah) {
			t.Errorf("BentukID(%q) mau ErrIDTidakSah, dapat %v", nomor, err)
		}
	}
}

// D3: ICD Code dan Disease wajib, dipangkas, HURUF BESAR (`SetUpperCase_DT` dipanggil kedua medan b1265 / b1538),
// paling panjang lebar kolom (100 / 1000 byte).
func TestNormalIsian(t *testing.T) {
	if got, err := NormalICDCode("  uji01.9 "); err != nil || got != "UJI01.9" {
		t.Errorf("icd %q %v", got, err)
	}
	if got, err := NormalDisease(" uji Kolera akut "); err != nil || got != "UJI KOLERA AKUT" {
		t.Errorf("disease %q %v", got, err)
	}
	for _, s := range []string{"", "   ", "\t\n"} {
		if _, err := NormalICDCode(s); err == nil || err.Error() != "ICD Code is required" {
			t.Errorf("icd kosong %q: %v", s, err)
		}
		if _, err := NormalDisease(s); err == nil || err.Error() != "Disease is required" {
			t.Errorf("disease kosong %q: %v", s, err)
		}
	}
	if _, err := NormalICDCode(strings.Repeat("A", BatasICDCode)); err != nil {
		t.Errorf("icd tepat batas: %v", err)
	}
	if _, err := NormalICDCode(strings.Repeat("A", BatasICDCode+1)); err == nil || !strings.Contains(err.Error(), "longer than 100") {
		t.Errorf("icd lebih: %v", err)
	}
	if _, err := NormalDisease(strings.Repeat("A", BatasDisease)); err != nil {
		t.Errorf("disease tepat batas: %v", err)
	}
	if _, err := NormalDisease(strings.Repeat("A", BatasDisease+1)); err == nil || !strings.Contains(err.Error(), "longer than 1000") {
		t.Errorf("disease lebih: %v", err)
	}
	// Byte, bukan karakter: 334 huruf 3-byte = 1002 byte.
	if _, err := NormalDisease(strings.Repeat("€", 334)); err == nil {
		t.Error("1002 byte UTF-8 diterima")
	}
	if KunciTeks(" uji01 ") != "UJI01" {
		t.Errorf("kunci %q", KunciTeks(" uji01 "))
	}
}

// Saring dipangkas dan tidak pernah melewati BatasSaring byte, terpotong pada batas huruf.
func TestRapikanSaring(t *testing.T) {
	if RapikanSaring("  kol  ") != "kol" {
		t.Error("pangkas")
	}
	s := RapikanSaring(strings.Repeat("€", 400)) // 1200 byte
	if len(s) > BatasSaring || len(s)%3 != 0 {
		t.Errorf("panjang %d", len(s))
	}
}

// `pyPageSize` 10 (b5169); bawaan urut ID.
func TestUkuranHalamanDanUrut(t *testing.T) {
	if UkuranHalaman != 10 || UrutID != "id" || UrutICD != "icd" {
		t.Errorf("%d %s %s", UkuranHalaman, UrutID, UrutICD)
	}
}
