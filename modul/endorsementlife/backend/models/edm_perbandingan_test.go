package models_test

import (
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
)

// TestSusunPerbandingan - presisi penuh, tanda, dan mata uang terpisah.
func TestSusunPerbandingan(t *testing.T) {
	kol := []string{"PREMIUM"}
	b, err := models.SusunPerbandingan([]models.KelompokRekap{
		{Currency: "USD", Kelompok: models.KelompokLama, Cacah: 1, Nilai: map[string]string{"PREMIUM": "0.58"}},
		{Currency: "IDR", Kelompok: models.KelompokLama, Cacah: 10, Nilai: map[string]string{"PREMIUM": "5850038.75"}},
		{Currency: "IDR", Kelompok: models.KelompokHapus, Cacah: 2, Nilai: map[string]string{"PREMIUM": "-2047513.5625"}},
		{Currency: "IDR", Kelompok: models.KelompokBaru, Cacah: 3, Nilai: map[string]string{"PREMIUM": "100.0001"}},
	}, kol)
	if err != nil || len(b) != 2 || b[0].Currency != "IDR" || b[1].Currency != "USD" {
		t.Fatalf("%+v %v", b, err)
	}
	i := b[0]
	for _, k := range []struct {
		nama  string
		sel   models.SelPerbandingan
		cacah int
		nilai string
	}{
		{"before", i.Before, 12, "7897552.3125"},
		{"deleted", i.Deleted, -2, "-2047513.5625"},
		{"added", i.Added, 3, "100.0001"},
		{"after", i.After, 13, "5850138.7501"},
		{"difference", i.Difference, 1, "-2047413.5624"},
	} {
		if k.sel.Cacah != k.cacah || k.sel.Nilai["PREMIUM"] != k.nilai {
			t.Errorf("%s = %d %q; mau %d %q", k.nama, k.sel.Cacah, k.sel.Nilai["PREMIUM"], k.cacah, k.nilai)
		}
	}
	if u := b[1]; u.Before.Nilai["PREMIUM"] != "0.58" || u.Deleted.Nilai["PREMIUM"] != "0" || u.Difference.Nilai["PREMIUM"] != "0" {
		t.Errorf("USD %+v", u)
	}
	if _, err := models.SusunPerbandingan([]models.KelompokRekap{{Currency: "IDR", Kelompok: models.KelompokLama,
		Nilai: map[string]string{"PREMIUM": "1,5"}}}, kol); err == nil {
		t.Error("angka rusak diterima")
	}
	if b, err := models.SusunPerbandingan(nil, kol); err != nil || len(b) != 0 {
		t.Errorf("kosong %+v %v", b, err)
	}
}
