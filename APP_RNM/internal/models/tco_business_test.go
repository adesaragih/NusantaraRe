package models

import (
	"errors"
	"testing"
)

func TestPeriksaBusinessTCO(t *testing.T) {
	if err := PeriksaBusinessTCO(BusinessTreaty{BizCode: "UJI-B1", IsActive: "1"}); err != nil {
		t.Errorf("wajar: %v", err)
	}
	if err := PeriksaBusinessTCO(BusinessTreaty{BizCode: "UJI-B1", IsActive: "0"}); err != nil {
		t.Errorf("nonaktif sah: %v", err)
	}
	if err := PeriksaBusinessTCO(BusinessTreaty{BizCode: " ", IsActive: "1"}); !errors.Is(err, ErrBusinessKodeKosong) {
		t.Errorf("kode kosong: %v", err)
	}
	for _, salah := range []string{"", "true", "Y", "2"} {
		if err := PeriksaBusinessTCO(BusinessTreaty{BizCode: "UJI-B1", IsActive: salah}); !errors.Is(err, ErrBusinessAktifTakSah) {
			t.Errorf("IsActive %q: %v", salah, err)
		}
	}
}

// Hilir membaca `isactive='1'` - nilai lain (termasuk NULL warisan) bukan aktif.
func TestBusinessAktifTCO(t *testing.T) {
	kasus := map[string]bool{"1": true, " 1 ": true, "0": false, "": false, "true": false}
	for v, mau := range kasus {
		if dapat := BusinessAktifTCO(BusinessTreaty{IsActive: v}); dapat != mau {
			t.Errorf("BusinessAktifTCO(%q) = %v, mau %v", v, dapat, mau)
		}
	}
}
