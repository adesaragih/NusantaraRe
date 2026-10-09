package models

import (
	"errors"
	"strings"
	"testing"
)

// K3: rumus prosedur PEGA_M_CAUSEOFLOSS_LIFE - '1' || LPAD(seq, 5, '0'); nomor lebih dari 5 angka DITOLAK (LPAD
// memotong). Sequence DEV last_number 5 -> ID berikutnya 100005.
func TestBentukIDRumusProsedur(t *testing.T) {
	for nomor, mau := range map[string]string{"1": "100001", "5": "100005", " 99999 ": "199999", "12": "100012"} {
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

// K4: Cause of Loss wajib (b996 / b1044), dipangkas, huruf TIDAK diubah (XML tanpa pengubah huruf), paling panjang 200
// byte (lebar kolom).
func TestNormalCauseOfLoss(t *testing.T) {
	if got, err := NormalCauseOfLoss("  uji Kecelakaan  "); err != nil || got != "uji Kecelakaan" {
		t.Errorf("rapikan %q %v", got, err)
	}
	for _, s := range []string{"", "   ", "\t\n"} {
		if _, err := NormalCauseOfLoss(s); err == nil || err.Error() != "Cause of Loss is required" {
			t.Errorf("kosong %q: %v", s, err)
		}
	}
	if _, err := NormalCauseOfLoss(strings.Repeat("A", BatasCauseOfLoss)); err != nil {
		t.Errorf("tepat %d byte ditolak: %v", BatasCauseOfLoss, err)
	}
	if _, err := NormalCauseOfLoss(strings.Repeat("A", BatasCauseOfLoss+1)); err == nil || !strings.Contains(err.Error(), "longer than 200") {
		t.Errorf("lebih dari %d byte: %v", BatasCauseOfLoss, err)
	}
	// Byte, bukan karakter: 67 huruf 3-byte = 201 byte.
	if _, err := NormalCauseOfLoss(strings.Repeat("€", 67)); err == nil {
		t.Error("201 byte UTF-8 diterima")
	}
	if KunciTeks(" uji Sakit ") != "UJI SAKIT" {
		t.Errorf("kunci %q", KunciTeks(" uji Sakit "))
	}
}

// pyPageSize 10 (b4057).
func TestUkuranHalaman(t *testing.T) {
	if UkuranHalaman != 10 {
		t.Errorf("ukuran %d", UkuranHalaman)
	}
}
