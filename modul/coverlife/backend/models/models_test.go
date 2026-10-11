package models

import (
	"errors"
	"strings"
	"testing"
)

// C2: rumus prosedur PEGA_M_COVER_LIFE - '1' || LPAD(seq, 5, '0'); nomor lebih dari 5 angka DITOLAK (LPAD memotong).
// Sequence DEV last_number 5 -> ID berikutnya 100005.
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

// C3: Cover wajib (b843 / b895), dipangkas, huruf TIDAK diubah, paling panjang 200 byte; Note opsional, hanya spasi =
// kosong, paling panjang 1000 byte, isi tidak diubah.
func TestNormalIsian(t *testing.T) {
	if got, err := NormalCover("  uji Kecelakaan  "); err != nil || got != "uji Kecelakaan" {
		t.Errorf("rapikan %q %v", got, err)
	}
	for _, s := range []string{"", "   ", "\t\n"} {
		if _, err := NormalCover(s); err == nil || err.Error() != "Cover is required" {
			t.Errorf("kosong %q: %v", s, err)
		}
	}
	if _, err := NormalCover(strings.Repeat("A", BatasCover)); err != nil {
		t.Errorf("tepat %d byte ditolak: %v", BatasCover, err)
	}
	if _, err := NormalCover(strings.Repeat("€", 67)); err == nil || !strings.Contains(err.Error(), "longer than 200") {
		t.Errorf("201 byte UTF-8: %v", err)
	}
	if got, err := NormalNote("  "); err != nil || got != "" {
		t.Errorf("note spasi %q %v", got, err)
	}
	if got, err := NormalNote(" uji catatan\nbaris dua "); err != nil || got != " uji catatan\nbaris dua " {
		t.Errorf("note apa adanya %q %v", got, err)
	}
	if _, err := NormalNote(strings.Repeat("A", BatasNote+1)); err == nil || !strings.Contains(err.Error(), "Note is longer than 1000") {
		t.Errorf("note panjang: %v", err)
	}
	if KunciTeks(" uji Jiwa ") != "UJI JIWA" {
		t.Errorf("kunci %q", KunciTeks(" uji Jiwa "))
	}
}

// pyPageSize 50 (b4101).
func TestUkuranHalaman(t *testing.T) {
	if UkuranHalaman != 50 {
		t.Errorf("ukuran %d", UkuranHalaman)
	}
}
