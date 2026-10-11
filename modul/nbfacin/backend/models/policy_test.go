package models

import (
	"testing"
	"time"
)

// TestProdDateTime - SaveJsonPolicy langkah 5: sekarang; Begin date >= sehari penuh di depan -> Begin date; hari Jakarta >
// closing -> tanggal 01 bulan berikut 05:00 GMT (Desember -> Januari tahun berikut).
func TestProdDateTime(t *testing.T) {
	jkt := time.FixedZone("WIB", 7*3600)
	j := func(s string) time.Time {
		w, err := time.ParseInLocation("2006-01-02 15:04", s, jkt)
		if err != nil {
			t.Fatal(err)
		}
		return w
	}
	for nama, c := range map[string]struct {
		sekarang, mulai time.Time
		closing         int
		mau             string
	}{
		"sekarang":            {j("2026-10-07 22:04"), time.Time{}, 25, "2026-10-07 22:04"},
		"tanggal = closing":   {j("2026-10-25 09:00"), time.Time{}, 25, "2026-10-25 09:00"},
		"lewat closing":       {j("2026-10-26 09:00"), time.Time{}, 25, "2026-11-01 12:00"},
		"lewat closing malam": {j("2026-10-25 23:30"), time.Time{}, 25, "2026-10-25 23:30"}, // UTC masih tgl 25, Jakarta 25
		"Jakarta tgl 26 dini": {j("2026-10-26 01:00"), time.Time{}, 25, "2026-11-01 12:00"}, // UTC 25 18:00
		"Desember":            {j("2026-12-28 10:00"), time.Time{}, 25, "2027-01-01 12:00"},
		"mulai depan":         {j("2026-10-07 10:00"), j("2026-10-19 00:00"), 25, "2026-10-19 00:00"},
		"mulai < sehari":      {j("2026-10-07 10:00"), j("2026-10-08 09:59"), 25, "2026-10-07 10:00"},
		"mulai tepat sehari":  {j("2026-10-07 10:00"), j("2026-10-08 10:00"), 25, "2026-10-08 10:00"},
		"mulai lalu":          {j("2026-10-07 10:00"), j("2026-09-01 00:00"), 25, "2026-10-07 10:00"},
		"mulai lewat closing": {j("2026-10-07 10:00"), j("2026-10-28 00:00"), 25, "2026-11-01 12:00"},
	} {
		if got := ProdDateTime(c.sekarang, c.mulai, c.closing, jkt).In(jkt).Format("2006-01-02 15:04"); got != c.mau {
			t.Errorf("%s: %s mau %s", nama, got, c.mau)
		}
	}
	if got := ProdDateTime(j("2026-10-26 09:00"), time.Time{}, 25, jkt); got.Location() != time.UTC || got.Hour() != jamPeriodeUTC {
		t.Errorf("jam GMT: %v", got)
	}
}

// jamPeriodeUTC - "01T050000.000 GMT".
const jamPeriodeUTC = 5

// TestRakitNomorPolis - KODE + "F" + BusinessOldId + "." + MM.YYYY + "." + LPAD 5; lebih dari 5 digit tidak dipotong.
func TestRakitNomorPolis(t *testing.T) {
	j := JenisNomorPolis("RNM")
	for urut, mau := range map[int]string{7: "RNMFB1.10.2026.00007", 1234: "RNMFB1.10.2026.01234", 12345: "RNMFB1.10.2026.12345", 123456: "RNMFB1.10.2026.123456"} {
		if got := RakitNomorPolis(j, "B1", "10.2026", urut); got != mau {
			t.Errorf("%d: %q", urut, got)
		}
	}
	if j != "RNMF" || KelasNomorPolis != "ASM-FW-GISFW-Work-NB/RNW" || HariClosingBawaan != 25 {
		t.Errorf("konstanta %q", j)
	}
}
