//go:build db

package repository_test

// Grid Rate of Exchange — pengenal mata uang untuk dropdown `.CurrencyID`,
// diadu dengan POOLDATA (BACA SAJA, diukur 7 Oktober 2026).

import "testing"

func TestKursTahunanMembawaPengenalMataUang(t *testing.T) {
	g, ctx := gudangBaca(t)
	kurs, err := g.BacaKursTahunan(ctx, "2025")
	if err != nil {
		t.Fatal(err)
	}
	if len(kurs) == 0 {
		t.Fatal("kurs 2025 kosong")
	}
	// `IDCURRENCY` = `CURRENCY.ID`: USD → 10001 (`BrowseCurrency_RD` `.ID`).
	for _, k := range kurs {
		if k.MataUangID == "" {
			t.Errorf("%s tanpa pengenal mata uang", k.MataUang)
		}
		if k.MataUang == "USD" && k.MataUangID != "10001" {
			t.Errorf("USD berpengenal %q, mau 10001", k.MataUangID)
		}
	}
}
