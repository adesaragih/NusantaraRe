package models

import (
	"errors"
	"strings"
	"testing"
)

// K3: '1' || LPAD(seq, 5, '0'); nomor > 5 angka DITOLAK; hasil muat VARCHAR2(6).
func TestBentukIDRumusProsedur(t *testing.T) {
	for nomor, mau := range map[string]string{"1": "100001", "44": "100044", " 99999 ": "199999"} {
		if got, err := BentukID(nomor); err != nil || got != mau || len(got) > BatasID {
			t.Errorf("BentukID(%q) = %q, %v; mau %q", nomor, got, err, mau)
		}
	}
	for _, nomor := range []string{"123456", "", "x", "-1"} {
		if _, err := BentukID(nomor); !errors.Is(err, ErrIDTidakSah) {
			t.Errorf("BentukID(%q): %v", nomor, err)
		}
	}
}

// K5: ketiga medan wajib (Business / Benefit tidak terbukti boleh kosong), dipangkas, <= 200 byte, huruf tidak diubah;
// semua masalah sekaligus.
func TestPeriksaIsian(t *testing.T) {
	got, err := PeriksaIsian(Isian{CoverName: "  Uji Plan  ", Business: " uji biz ", Benefit: "UJI MANFAAT"})
	if err != nil || got != (Isian{CoverName: "Uji Plan", Business: "uji biz", Benefit: "UJI MANFAAT"}) {
		t.Errorf("rapikan %+v %v", got, err)
	}
	_, err = PeriksaIsian(Isian{CoverName: " ", Benefit: strings.Repeat("A", 201)})
	if err == nil || err.Error() != "Plan Name is required; Business is required; Benefit is longer than 200 characters" {
		t.Errorf("galat gabungan %v", err)
	}
	if _, err := PeriksaIsian(Isian{CoverName: strings.Repeat("A", 200), Business: "B", Benefit: "C"}); err != nil {
		t.Errorf("tepat 200 byte: %v", err)
	}
	if KunciTeks(" Uji Biz ") != "UJI BIZ" {
		t.Error("KunciTeks")
	}
}

func TestKonstanta(t *testing.T) {
	if UkuranHalaman != 10 || GrupBusinessLife != "009" || BatasIDMaster != 10 {
		t.Errorf("%d %s %d", UkuranHalaman, GrupBusinessLife, BatasIDMaster)
	}
}
