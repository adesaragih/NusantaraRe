package models

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func tgl(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

func tahunWajar() TahunTreaty {
	return TahunTreaty{TreatyYear: "2026", TreatyGroupID: "10001",
		StartDate: tgl(2026, 1, 1), EndDate: tgl(2026, 12, 31)}
}

func TestPeriksaTahunTreatyWajarLolos(t *testing.T) {
	if err := PeriksaTahunTreaty(tahunWajar()); err != nil {
		t.Fatal(err)
	}
}

func TestPeriksaTahunTreatyGerbangExisting(t *testing.T) {
	kasus := []struct {
		ubah func(*TahunTreaty)
		mau  error
	}{
		{func(x *TahunTreaty) { x.TreatyGroupID = "  " }, ErrTahunTreatyGrupKosong},
		{func(x *TahunTreaty) { x.TreatyYear = "" }, ErrTahunTreatyTahunKosong},
		{func(x *TahunTreaty) { x.TreatyYear = "20A6" }, ErrTahunTreatyBukanAngka},
		{func(x *TahunTreaty) { x.TreatyYear = "2026.5" }, ErrTahunTreatyBukanAngka},
	}
	for _, k := range kasus {
		x := tahunWajar()
		k.ubah(&x)
		if err := PeriksaTahunTreaty(x); !errors.Is(err, k.mau) {
			t.Errorf("mau %v, dapat %v", k.mau, err)
		}
	}
}

// AC 9 - gerbang tambahan sadar: pesan menyebut KEDUA medan dan nilainya.
func TestPeriksaPeriodeTCOMenyebutMedan(t *testing.T) {
	err := PeriksaPeriodeTCO(tgl(2026, 12, 31), tgl(2026, 1, 1), "StartDate", "EndDate")
	if !errors.Is(err, ErrPeriodeTerbalik) {
		t.Fatalf("mau ErrPeriodeTerbalik, dapat %v", err)
	}
	for _, mau := range []string{"EndDate", "2026-01-01", "StartDate", "2026-12-31"} {
		if !strings.Contains(err.Error(), mau) {
			t.Errorf("pesan tanpa %q: %s", mau, err)
		}
	}
	// Sama hari sah; salah satu kosong bukan pelanggaran gerbang ini.
	if err := PeriksaPeriodeTCO(tgl(2026, 1, 1), tgl(2026, 1, 1), "a", "b"); err != nil {
		t.Errorf("sama hari ditolak: %v", err)
	}
	if err := PeriksaPeriodeTCO(time.Time{}, tgl(2026, 1, 1), "a", "b"); err != nil {
		t.Errorf("mulai kosong ditolak: %v", err)
	}
	x := tahunWajar()
	x.EndDate = tgl(2025, 12, 31)
	if err := PeriksaTahunTreaty(x); !errors.Is(err, ErrPeriodeTerbalik) {
		t.Errorf("tahun dengan periode terbalik lolos: %v", err)
	}
}
