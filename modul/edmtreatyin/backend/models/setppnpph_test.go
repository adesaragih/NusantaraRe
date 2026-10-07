package models

// Uji seam 3 - `Activity/SetPPNPPH` dengan nilai dari XML (tiket 07, 13; AC 27,
// 28, 79). Masukan dan harapan DIHITUNG TANGAN dari `PropertiesValue`:
//
//	langkah 2   .PPHValue = 0 ; .PPNValue = 0
//	langkah 4   syarat: .FlagPPH=="true"                    benar -> lewati syarat berikut (jalan)
//	                    ListAgent.pxResults(1).STS_PKP == 1  salah -> lewati langkah
//	  .BrokerageFee           = @divide(2.5,100,8) * (.PremiOgp + .PremiOnp)
//	                          = 0,025 x (1000 + 600)                        = 40
//	  .BrokerageFeeSebenarnya = @if(.TypeTax=="Inclusive",
//	                                @divide(.Deduction1,@divide(102.2,100,8),8), .Deduction1)
//	                          Inclusive: 51,1 / 1,022                       = 50
//	                          lainnya  : 51,1
//	  .PPHValue = .BrokerageFeeSebenarnya * @divide(2,100,8)   = 50 x 0,02  = 1     | 51,1 x 0,02  = 1,022
//	  .PPNValue = .BrokerageFeeSebenarnya * @divide(2.2,100,8) = 50 x 0,022 = 1,1   | 51,1 x 0,022 = 1,1242

import "testing"

func TestSetPPNPPHNilaiXML(t *testing.T) {
	for _, tt := range []struct {
		nama, flag, pkp, typeTax string
		jalan                    bool
		fee, seb, pph, ppn       string
	}{
		{"FlagPPH true, agen bukan PKP", "true", "", TypeTaxInclusive, true, "40", "50", "1", "1.1"},
		{"FlagPPH kosong, agen PKP", "", "1", "Exclusive", true, "40", "51.1", "1.022", "1.1242"},
		{"FlagPPH false, agen PKP", "false", "1", TypeTaxInclusive, true, "40", "50", "1", "1.1"},
		{"FlagPPH true melewati syarat PKP", "true", "0", "", true, "40", "51.1", "1.022", "1.1242"},
		{"FlagPPH false, agen bukan PKP", "false", "0", TypeTaxInclusive, false, "7", "9", "0", "0"},
		{"agen tak ditemukan (pxResults kosong)", "", "", TypeTaxInclusive, false, "7", "9", "0", "0"},
		{"FlagPPH huruf besar bukan \"true\"", "TRUE", "", TypeTaxInclusive, false, "7", "9", "0", "0"},
	} {
		t.Run(tt.nama, func(t *testing.T) {
			h := HalamanBaru()
			set := func(m, v string) { h.Setel("PolicyTreatyIn."+m, v) }
			set("PremiOgp", "1000")
			set("PremiOnp", "600")
			set("Deduction1", "51.1")
			set("TypeTax", tt.typeTax)
			set("FlagPPH", tt.flag)
			// nilai lama: membuktikan langkah 2 menolkan PPH/PPN dan langkah 4
			// yang dilewati tidak menyentuh BrokerageFee*.
			set("BrokerageFee", "7")
			set("BrokerageFeeSebenarnya", "9")
			set("PPHValue", "3")
			set("PPNValue", "4")
			h.Setel(JalurStsPKP, tt.pkp)
			if err := SetPPNPPH(h); err != nil {
				t.Fatal(err)
			}
			samaAngka(t, h, "PolicyTreatyIn.BrokerageFee", tt.fee)
			samaAngka(t, h, "PolicyTreatyIn.BrokerageFeeSebenarnya", tt.seb)
			samaAngka(t, h, "PolicyTreatyIn.PPHValue", tt.pph)
			samaAngka(t, h, "PolicyTreatyIn.PPNValue", tt.ppn)
		})
	}
}
