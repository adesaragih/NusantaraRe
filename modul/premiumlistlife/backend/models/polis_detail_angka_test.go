package models

// Kolom angka grid peserta - layar memberinya pemisah ribuan (03-10-2026).

import (
	"slices"
	"testing"
)

func TestNamaKolomGridAngkaHanyaKolomAngka(t *testing.T) {
	angka := NamaKolomGridAngka()
	for _, k := range []string{"SUM_INSURED", "GROSS_PREMIUM", "NET_PREMIUM", "CEDING_RETENTION"} {
		if !slices.Contains(angka, k) {
			t.Errorf("%s bukan kolom angka", k)
		}
	}
	for _, k := range []string{"CERTIFICATE_NO", "POLICY_HOLDER", "DOB", "BEGIN_DATE", "STNC", "NAME_OF_INSURED"} {
		if slices.Contains(angka, k) {
			t.Errorf("%s keliru dianggap kolom angka", k)
		}
	}
}
