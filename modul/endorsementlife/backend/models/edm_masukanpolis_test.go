package models_test

import (
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

func TestNomorPolisDariMasukan(t *testing.T) {
	for _, k := range []struct {
		masuk, polis, akhiran string
		ada                   bool
	}{
		{"RNML-QRL2.10.26.00042/01", "RNML-QRL2.10.26.00042", "01", true},
		{"  RNML-QRL2.10.26.00042/02  ", "RNML-QRL2.10.26.00042", "02", true},
		{"RNML-QRL2.10.26.00042", "RNML-QRL2.10.26.00042", "", false},
		{"RNML-QRL2/abc", "RNML-QRL2/abc", "", false},
		{"RNML-QRL2.10.26.00042/100", "RNML-QRL2.10.26.00042", "100", true},
		{"UJI/A/B/07", "UJI/A/B", "07", true},
		{"/01", "/01", "", false},
		{"", "", "", false},
	} {
		p, a, ada := models.NomorPolisDariMasukan(k.masuk)
		if p != k.polis || a != k.akhiran || ada != k.ada {
			t.Errorf("NomorPolisDariMasukan(%q) = %q %q %v", k.masuk, p, a, ada)
		}
	}
	if models.AkhiranVersi(0) != "" || models.AkhiranVersi(2) != "02" || models.AkhiranVersi(12) != "12" {
		t.Error("AkhiranVersi")
	}
}
