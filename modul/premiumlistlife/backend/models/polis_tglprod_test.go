package models

import (
	"errors"
	"testing"
	"time"
)

// TestTanggalProduksi - TGL_PROD JSON_POLIS mengikuti periode PL Number (06-10-2026).
func TestTanggalProduksi(t *testing.T) {
	wib := time.FixedZone("WIB", 7*60*60)
	kasus := []struct {
		saat    time.Time
		periode string
		mau     string
	}{
		// Periode = bulan transaksi → tanggal transaksi, tanpa jam.
		{time.Date(2026, 10, 6, 15, 30, 0, 0, wib), "10.2026", "2026-10-06"},
		// Lewat tanggal closing → periode bulan depan → tanggal 1.
		{time.Date(2026, 10, 28, 9, 0, 0, 0, wib), "11.2026", "2026-11-01"},
		// Desember → Januari tahun berikutnya.
		{time.Date(2026, 12, 30, 9, 0, 0, 0, wib), "01.2027", "2027-01-01"},
		// Hari dibaca di Jakarta: 31-10 20:00 UTC = 01-11 03:00 WIB.
		{time.Date(2026, 10, 31, 20, 0, 0, 0, time.UTC), "11.2026", "2026-11-01"},
		{time.Date(2026, 11, 1, 20, 0, 0, 0, time.UTC), "11.2026", "2026-11-02"},
	}
	for _, k := range kasus {
		got, err := TanggalProduksi(k.saat, k.periode)
		if err != nil {
			t.Fatalf("%v %s: %v", k.saat, k.periode, err)
		}
		if got.Format("2006-01-02") != k.mau || got.Hour() != 0 || got.Minute() != 0 {
			t.Errorf("%v %s = %v, mau %s", k.saat, k.periode, got, k.mau)
		}
	}
	for _, salah := range []string{"", "2026-10", "13.2026", "AB.2026"} {
		if _, err := TanggalProduksi(time.Now(), salah); !errors.Is(err, ErrTanggalProduksi) {
			t.Errorf("periode %q: %v", salah, err)
		}
	}
}
