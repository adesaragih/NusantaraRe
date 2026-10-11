package services

import (
	"errors"
	"reflect"
	"testing"
)

func kolomAkumulasi(xs []BarisAkumulasi, f func(BarisAkumulasi) string) []string {
	out := []string{}
	for _, x := range xs {
		out = append(out, f(x))
	}
	return out
}

func akumulasiUji(t *testing.T, m MasukanAkumulasi) []BarisAkumulasi {
	t.Helper()
	h, err := HitungAkumulasi(m)
	if err != nil {
		t.Fatalf("HitungAkumulasi: %v", err)
	}
	return h.AccumulationList
}

// TestAkumulasiDariReportingPeriod - `TreatyInSetAccountReport`: banyak baris
// = selisih bulan Start–End (tab Reporting Period) ÷ interval + 1.
func TestAkumulasiDariReportingPeriod(t *testing.T) {
	for _, k := range []struct {
		periode string
		period  []string
		tanggal []string
	}{
		{"quarter", []string{"Q 1", "Q 2", "Q 3", "Q 4"}, []string{"20250101", "20250401", "20250701", "20251001"}},
		{"half", []string{"H 1", "H 2"}, []string{"20250101", "20250701"}},
	} {
		xs := akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiPeriode, AccumulationPeriod: k.periode, ReportingStart: "20250101", ReportingEnd: "20251231"})
		if got := kolomAkumulasi(xs, func(b BarisAkumulasi) string { return b.Period }); !reflect.DeepEqual(got, k.period) {
			t.Errorf("%s Period %v", k.periode, got)
		}
		if got := kolomAkumulasi(xs, func(b BarisAkumulasi) string { return b.ReportDate }); !reflect.DeepEqual(got, k.tanggal) {
			t.Errorf("%s ReportDate %v", k.periode, got)
		}
	}
	// Bentuk kabel layar (DD-MM-YYYY) diterima sama.
	xs := akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiPeriode, AccumulationPeriod: "month", ReportingStart: "31-01-2025", ReportingEnd: "31-03-2025"})
	// TempDate BERANTAI (`@addCalendar(local.TempDate, …)`): 31 Jan → 28 Feb
	// → 28 Mar, bukan 31 Mar - pola yang sama dengan TreatyInSetReport.
	if got := kolomAkumulasi(xs, func(b BarisAkumulasi) string { return b.ReportDate }); !reflect.DeepEqual(got, []string{"20250131", "20250228", "20250328"}) {
		t.Errorf("month (berantai dari TempDate): %v", got)
	}
}

// TestAkumulasiNoneDanOther - langkah 2 (none: keluar, daftar utuh) dan 6
// (other: daftar dibuang).
func TestAkumulasiNoneDanOther(t *testing.T) {
	lama := []BarisAkumulasi{{Period: "T 1", ReportDate: "20240101", SubDays: "5"}}
	if xs := akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiPeriode, AccumulationPeriod: "none", ReportingStart: "20250101", ReportingEnd: "20251231", AccumulationList: lama}); !reflect.DeepEqual(xs, lama) {
		t.Errorf("none: %+v", xs)
	}
	if xs := akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiPeriode, AccumulationPeriod: "other", AccumulationList: lama}); len(xs) != 0 || xs == nil {
		t.Errorf("other: %+v", xs)
	}
}

// TestAkumulasiTanpaReportingPeriod - Start kosong (tab Reporting Period
// belum diisi): satu baris tanpa tanggal; periode kosong: durasi TIDAK dibagi.
func TestAkumulasiTanpaReportingPeriod(t *testing.T) {
	xs := akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiPeriode, AccumulationPeriod: "quarter"})
	if !reflect.DeepEqual(xs, []BarisAkumulasi{{Period: "Q 1"}}) {
		t.Errorf("tanpa Start/End: %+v", xs)
	}
	xs = akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiPeriode, ReportingStart: "20250101", ReportingEnd: "20250401"})
	if got := kolomAkumulasi(xs, func(b BarisAkumulasi) string { return b.Period + "@" + b.ReportDate }); !reflect.DeepEqual(got, []string{"1@20250101", "2@20250101", "3@20250101"}) {
		t.Errorf("periode kosong: %v", got)
	}
}

// TestAkumulasiJatuhTempo - `TreatyInAccumulationSetSubDue`.
func TestAkumulasiJatuhTempo(t *testing.T) {
	xs := akumulasiUji(t, MasukanAkumulasi{Aksi: AksiAkumulasiJatuhTempo, AccumulationList: []BarisAkumulasi{
		{Period: "Q 1", ReportDate: "20250101", SubDays: "30"},
		{Period: "Q 2", ReportDate: "20250401", SubDays: "", SubDueDate: "lama"},
		{Period: "Q 3", ReportDate: "20250701", SubDays: "x", SubDueDate: "lama"},
		{Period: "Q 4", ReportDate: "01-10-2025", SubDays: "45"},
	}})
	if got := kolomAkumulasi(xs, func(b BarisAkumulasi) string { return b.SubDueDate }); !reflect.DeepEqual(got, []string{"20250131", "lama", "lama", "20251115"}) {
		t.Errorf("SubDueDate %v", got)
	}
}

func TestAkumulasiAksiTakDikenal(t *testing.T) {
	if _, err := HitungAkumulasi(MasukanAkumulasi{Aksi: "hapus"}); !errors.Is(err, ErrMasukanTidakSah) {
		t.Errorf("err %v", err)
	}
}
