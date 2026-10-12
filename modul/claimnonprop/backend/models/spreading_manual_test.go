package models

import "testing"

// TestNormalkanSpreadingManualRelatifTerhadapRNM - spreading manual master (persen treaty, total = % RNM Share)
// menjadi persen relatif bagian RNM seperti spreading bernama (CLMNP-000017 / master 1002320, dokumen Pega 1002304).
func TestNormalkanSpreadingManualRelatifTerhadapRNM(t *testing.T) {
	ors := ShareXOL{SpreadingTotalPctXOL: "40", RNMShare: "40",
		SpreadingListXOL: []SpreadingMaster{{ReinsTypeID: "10007", ReinsTypeName: "ORS", Pct: "40"}}}
	NormalkanSpreadingManual(&ors, "")
	if ors.SpreadingTotalPctXOL != "100" || ors.SpreadingListXOL[0].Pct != "100" {
		t.Errorf("ORS 40 dari RNM 40 = total %q, ORS %q; harap 100 / 100", ors.SpreadingTotalPctXOL, ors.SpreadingListXOL[0].Pct)
	}

	// 1002304: 14 + 2,5 + 2,5 = 19 (RNM baris kosong -> RNM master).
	qs := ShareXOL{SpreadingTotalPctXOL: "19.0", SpreadingListXOL: []SpreadingMaster{
		{ReinsTypeName: "2026 QS 212M TRT", Pct: "14"}, {ReinsTypeName: "2026 QS 17M TRT AS", Pct: "2.5"},
		{ReinsTypeName: "2026 QS 17M TRT EB", Pct: "2.5"}}}
	NormalkanSpreadingManual(&qs, "19")
	if qs.SpreadingTotalPctXOL != "100" || !SamaAngka(qs.SpreadingListXOL[1].Pct, qs.SpreadingListXOL[2].Pct) {
		t.Errorf("total %q, rincian %+v", qs.SpreadingTotalPctXOL, qs.SpreadingListXOL)
	}
	var kal Kalkulator
	jumlah := kal.D("0")
	for _, sp := range qs.SpreadingListXOL {
		jumlah = kal.Tambah(jumlah, kal.Teks(".Pct", sp.Pct))
	}
	if f, err := jumlah.Float64(); kal.Galat() != nil || err != nil || f < 99.999999 || f > 100.000001 {
		t.Errorf("jumlah rincian relatif = %s, harap 100 (%v)", Teks(jumlah), kal.Galat())
	}

	// Spreading bernama tidak disentuh (sudah relatif).
	bernama := ShareXOL{SpreadingTypeIDXOL: "10241", SpreadingTypeXOL: "2024 QS 155M TRT", SpreadingTotalPctXOL: "100",
		RNMShare: "30", SpreadingListXOL: []SpreadingMaster{{ReinsTypeName: "QS (OR)", Pct: "40"}, {ReinsTypeName: "QS (R/I)", Pct: "60"}}}
	NormalkanSpreadingManual(&bernama, "30")
	if bernama.SpreadingTotalPctXOL != "100" || bernama.SpreadingListXOL[0].Pct != "40" {
		t.Errorf("spreading bernama berubah: %+v", bernama)
	}

	// % RNM 0 / tak terbaca -> dibiarkan.
	nol := ShareXOL{SpreadingTotalPctXOL: "40", SpreadingListXOL: []SpreadingMaster{{ReinsTypeName: "ORS", Pct: "40"}}}
	NormalkanSpreadingManual(&nol, "0")
	if nol.SpreadingListXOL[0].Pct != "40" {
		t.Errorf("RNM 0 mestinya dibiarkan: %+v", nol)
	}
}
