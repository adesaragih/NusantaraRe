package models

// Format unggah ceding - keputusan work owner 02-10-2026.

import (
	"slices"
	"testing"
)

func TestDesimalKomaKeTitikHanyaSatuKomaTanpaTitik(t *testing.T) {
	for masuk, mau := range map[string]string{
		"10,70944011": "10.70944011",
		" 0,5 ":       "0.5",
		"10000":       "10000",
		"1.5":         "1.5",
		"1.234,56":    "1.234,56", // pemisah ribuan: dibiarkan, lalu ditolak UangCSV
		"1,234,567":   "1,234,567",
	} {
		if got := DesimalKomaKeTitik(masuk); got != mau {
			t.Errorf("DesimalKomaKeTitik(%q) = %q, mau %q", masuk, got, mau)
		}
	}
}

func TestEmpatTanggalTidakWajibTetapiDiperiksaBentuknya(t *testing.T) {
	opsional := []string{"STNC", "START_DATE", "EFFECTIVE_DATE",
		"RETROCESSION_VALUATION_BEGIN_DATE", "RETROCESSION_VALUATION_EXPIRED_DATE"}
	for _, k := range opsional {
		if slices.Contains(KolomWajibUnggah(), k) {
			t.Errorf("%s masih kolom judul wajib", k)
		}
	}
	b := barisSah(1)
	for _, k := range opsional {
		delete(b.Nilai, k)
	}
	if h := ValidasiUnggah([]BarisUnggah{b}); !h.Lolos() {
		t.Errorf("baris tanpa empat tanggal opsional ditolak: %+v", h.Ditolak)
	}
	b.Nilai["EFFECTIVE_DATE"] = "2026-01-01"
	h := ValidasiUnggah([]BarisUnggah{b})
	if len(h.Ditolak) != 1 || h.Ditolak[0].Kolom != "EFFECTIVE_DATE" {
		t.Errorf("tanggal opsional berbentuk salah tidak ditolak tepat: %+v", h.Ditolak)
	}
}
