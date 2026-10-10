package handlers_test

// Premi diubah (work owner 10-10-2026: "udah di ubah premi OGP, tapi deduction A tidak berubah otomatis, harus di
// triger dulu"): sel Premi Ogp / Onp mengirim Count*_Act(Data=Pct) untuk Deduction In A / B OGP dan ONP LALU
// CountOGPONP (frontend/medan.ts `PREMI`). Deduction1 / Deduction2 TETAP isian manual - tidak ikut premi.

import (
	"testing"

	"github.com/cockroachdb/apd/v3"
)

// urutanSelPremi - salinan `PREMI` frontend/medan.ts (diuji sama di medan.test.ts).
var urutanSelPremi = []map[string]string{
	{"aksi": "CountResult1", "param": "Pct"},
	{"aksi": "CountResult2Ogp", "param": "Pct"},
	{"aksi": "CountResult1Onp", "param": "Pct"},
	{"aksi": "CountResult2Onp", "param": "Pct"},
	{"aksi": "CountOGPONP"},
}

func samaNilai(a, b string) bool {
	x, _, e1 := apd.NewFromString(a)
	y, _, e2 := apd.NewFromString(b)
	return e1 == nil && e2 == nil && x.Cmp(y) == 0
}

func TestPremiDiubahDeductionInABIkutDeduction1Tetap(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("")
	// keadaan lama: premi 1000, Deduction In A 32,5% = 325; premi diubah jadi 2000
	h.Setel("PolicyTreatyIn.PremiOgp", "2000")
	h.Setel("PolicyTreatyIn.RiCommOgp", "32.5")
	h.Setel("PolicyTreatyIn.ResultOgp1", "325")
	h.Setel("PolicyTreatyIn.OveriddingCommOgp", "5")
	h.Setel("PolicyTreatyIn.ResultOgp2", "50")
	h.Setel("PolicyTreatyIn.Deduction1", "25")
	ly := u.hitung(id, map[string]any{"urutan": urutanSelPremi, "halaman": h})
	for j, harap := range map[string]string{
		"PolicyTreatyIn.ResultOgp1": "650", // Deduction In A (OGP) = 2000 x 32,5%
		"PolicyTreatyIn.ResultOgp2": "100", // Deduction In B (OGP) = 2000 x 5%
		"PolicyTreatyIn.Deduction1": "25",  // manual - tidak ikut premi
		// Net Premium = (2000 - 650) - 5% x 2000 - 25
		"PolicyTreatyIn.NetPremium": "1225",
	} {
		if got := ly.Halaman.Ambil(j); !samaNilai(got, harap) {
			t.Errorf("%s = %q, harap %s", j, got, harap)
		}
	}
}
