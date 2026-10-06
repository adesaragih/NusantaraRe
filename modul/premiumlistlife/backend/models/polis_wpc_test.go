package models

import (
	"errors"
	"testing"
)

// Tabel work owner 03-10-2026, setiap bulan produksi.
func TestWPCPolisTabelWorkOwner(t *testing.T) {
	for _, k := range []struct{ tipe, periode, mau string }{
		{"QR", "01.2026", "2026-04-30"}, {"QR", "03.2026", "2026-04-30"},
		{"QP", "04.2026", "2026-07-31"}, {"QR", "06.2026", "2026-07-31"},
		{"QR", "07.2026", "2026-10-31"}, {"QP", "09.2026", "2026-10-31"},
		{"QR", "10.2026", "2027-01-31"}, {"QP", "12.2026", "2027-01-31"},
		{"TP", "02.2026", "2026-05-31"},
		{"TR", "05.2026", "2026-08-31"},
		{"TP", "08.2026", "2026-11-30"},
		{"TR", "11.2026", "2027-02-28"},
		{"TP", "10.2027", "2028-02-29"}, // kabisat
	} {
		got, err := WPCPolis(k.tipe, k.periode)
		if err != nil || got.Format("2006-01-02") != k.mau {
			t.Errorf("WPCPolis(%s, %s) = %s, %v - mau %s", k.tipe, k.periode, got.Format("2006-01-02"), err, k.mau)
		}
	}
}

func TestWPCPolisMenolakMasukanAsing(t *testing.T) {
	for _, k := range []struct{ tipe, periode string }{
		{"", "01.2026"}, {"XX", "01.2026"}, {"QR", ""}, {"QR", "13.2026"}, {"QR", "2026-01"},
	} {
		if _, err := WPCPolis(k.tipe, k.periode); !errors.Is(err, ErrWPCTakDapatDihitung) {
			t.Errorf("WPCPolis(%q, %q) galat %v", k.tipe, k.periode, err)
		}
	}
}
