package models

// Mesin tangga Komite - tiket 02. TANPA Oracle.

import (
	"errors"
	"testing"
)

// TestTanggaTigaTingkatSetujuSeluruhnya - AC 1, 5 spec.
func TestTanggaTigaTingkatSetujuSeluruhnya(t *testing.T) {
	for _, u := range []struct {
		count            int
		berlanjut, akhir bool
	}{
		{1, true, false},
		{2, true, false},
		{3, false, true}, // tingkat akhir: akseptasi, tangga selesai
	} {
		a, err := TerapkanKeputusanKomite(KeputusanKomiteSetuju, u.count, 3)
		if err != nil {
			t.Fatal(err)
		}
		if a.Berlanjut != u.berlanjut || a.AkseptasiAkhir != u.akhir || a.CountBaru != u.count+1 {
			t.Errorf("count %d: %+v", u.count, a)
		}
		if a.TolakAkhir {
			t.Errorf("count %d: Setuju menyalakan TolakAkhir", u.count)
		}
	}
}

// TestTolakMenghentikanDiTingkatManaPun - AC 6 spec.
//
// ⚠️ Langkah 5 (tulis penolakan ke baris klaim) HANYA di tingkat akhir
// (b8119). Tolak di tingkat tengah menghentikan tangga TANPA langkah 5 -
// itu yang tiket 05 ("jalur balik dua tingkat") putuskan.
func TestTolakMenghentikanDiTingkatManaPun(t *testing.T) {
	for count := 1; count <= 3; count++ {
		a, err := TerapkanKeputusanKomite(KeputusanKomiteTolak, count, 3)
		if err != nil {
			t.Fatal(err)
		}
		if a.Berlanjut || a.AkseptasiAkhir {
			t.Errorf("Tolak di tingkat %d berlanjut/akseptasi: %+v", count, a)
		}
		if a.TolakAkhir != (count == 3) {
			t.Errorf("Tolak di tingkat %d: TolakAkhir %v", count, a.TolakAkhir)
		}
	}
}

// TestKeputusanDiLuarEnumDitolakTerang - AC 35 spec.
func TestKeputusanDiLuarEnumDitolakTerang(t *testing.T) {
	for _, k := range []string{"", "0", "3", "Setuju", "01"} {
		if _, err := TerapkanKeputusanKomite(k, 1, 3); !errors.Is(err, ErrKeputusanKomiteTidakDikenal) {
			t.Errorf("keputusan %q: %v", k, err)
		}
	}
	if _, err := TerapkanKeputusanKomite("1", 4, 3); err == nil {
		t.Error("count melebihi loop diterima")
	}
}

// TestIsKomiteLoopVERBATIM - `.AcceptStatus = "1"` DAN `.KomiteCount <= .KomiteLoop`.
func TestIsKomiteLoopVERBATIM(t *testing.T) {
	if !TanggaBerlanjut("1", 3, 3) || TanggaBerlanjut("1", 4, 3) || TanggaBerlanjut("2", 1, 3) {
		t.Error("TanggaBerlanjut tidak sama dengan IsKomiteLoop")
	}
	// Kasus baru tiba lewat Start2 [Always]; sesudah berhenti, tidak lagi.
	if !KasusDiTangga("", 1, 3) || KasusDiTangga("2", 2, 3) || KasusDiTangga("1", 4, 3) {
		t.Error("KasusDiTangga salah menilai kasus baru / berhenti / selesai")
	}
	if KataKeputusanKomite("1") != "Setuju" || KataKeputusanKomite("2") != "Tolak" {
		t.Error("kata keputusan")
	}
}

// TestEskalasiHanyaNaikSatu - AC 11 spec Komite.
func TestEskalasiHanyaNaikSatu(t *testing.T) {
	if n, err := EskalasiNaik(1, 3); err != nil || n != 2 {
		t.Errorf("eskalasi 1/3 = %d, %v", n, err)
	}
	if _, err := EskalasiNaik(3, 3); !errors.Is(err, ErrEskalasiTanpaTingkatAtas) {
		t.Errorf("eskalasi dari tingkat akhir: %v", err)
	}
	if _, err := EskalasiNaik(4, 3); err == nil {
		t.Error("tangga tidak sah diterima")
	}
}
