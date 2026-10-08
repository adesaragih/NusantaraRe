package models_test

import (
	"strings"
	"testing"

	"nusantarare/modul/riratelife/backend/models"
)

// Form Rate Detail (View Detail.xml): CONTRACT b2141 dan RATE b2579 wajib; GENDER b1938 dan AGE b2400 tidak.
func TestPeriksaIsianRate(t *testing.T) {
	got, err := models.PeriksaIsianRate(models.IsianRate{Gender: " m ", Contract: "05", Age: "", Rate: "0.7500"})
	if err != nil || got != (models.IsianRate{Gender: "M", Contract: "5", Age: "", Rate: "0,7500"}) {
		t.Fatalf("sah %+v %v", got, err)
	}
	if got, err := models.PeriksaIsianRate(models.IsianRate{Contract: "0", Rate: "1,25"}); err != nil || got.Gender != "" || got.Rate != "1,25" {
		t.Errorf("gender kosong %+v %v", got, err)
	}
	for _, k := range []struct {
		isi models.IsianRate
		mau string
	}{
		{models.IsianRate{Rate: "1"}, "CONTRACT is required"},
		{models.IsianRate{Contract: "1"}, "RATE is required"},
		{models.IsianRate{Gender: "X", Contract: "1", Rate: "1"}, "GENDER must be U, M, or F"},
		{models.IsianRate{Contract: "121", Rate: "1"}, "CONTRACT must be a whole number from 0 to 120"},
		{models.IsianRate{Contract: "1", Age: "-1", Rate: "1"}, "AGE must be a whole number"},
		{models.IsianRate{Contract: "1", Rate: "1.000,5"}, "RATE must be a number"},
	} {
		if _, err := models.PeriksaIsianRate(k.isi); err == nil || !strings.Contains(err.Error(), k.mau) {
			t.Errorf("%+v: galat %v, mau %q", k.isi, err, k.mau)
		}
	}
}
