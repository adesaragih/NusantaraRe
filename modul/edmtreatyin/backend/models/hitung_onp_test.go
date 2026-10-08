package models

// Uji audit silang putaran 3 (bab 3.2 butir 1): aktivitas rantai uang yang
// sudah dibangun tetapi belum punya uji berharapan XML. Setiap harapan
// DIHITUNG TANGAN dari `PropertiesValue` langkah yang dikutip (korpus
// `Activity/<nama>.xml`), bukan dengan memanggil rumus kode.
//
// Kode syarat langkah Pega (INVENTARIS bab 6): 2 lanjut, 3 lewati langkah,
// 6 keluar aktivitas.

import "testing"

// halamanOnp - hanya sisi ONP terisi; FlagPPH kosong dan STS_PKP kosong
// sehingga `SetPPNPPH` langkah 4 dilewati (PPH/PPN 0), Deduction1 kosong
// sehingga `CountNetPremi_act` langkah 6 dilewati.
func halamanOnp() *Halaman {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisProporsional)
	h.Setel("PolicyTreatyIn.PremiOnp", "2000")
	return h
}

// CountResult1Onp_Act (Param.Data):
//
//	3  syarat [(.RiCommOnp<100 && isDouble)||.RiCommOnp=""] F->6; [Data=="Pct"] T->2 F->3
//	   .ResultOnp1 = .PremiOnp*(.RiCommOnp/@String.toDecimal("100.00000"))
//	4  [Data=="Amount"] .RiCommOnp = (.ResultOnp1/.PremiOnp)*@String.toDecimal("100.00000")
//	5  call CountNetPremi_act
//
// Pct: ResultOnp1 = 2000 x 10/100 = 200; NetPremium (langkah 4 CountNetPremi)
// = (0-0)+(2000-200) = 1800; BalanceDueTo = 1800; DueTo 1.
// Amount: RiCommOnp = 300/2000 x 100 = 15; NetPremium = 2000-300 = 1700.
func TestCountResult1OnpPctDanAmount(t *testing.T) {
	h := halamanOnp()
	h.Setel("PolicyTreatyIn.RiCommOnp", "10")
	if err := CountResult1Onp(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.ResultOnp1", "200")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "1800")
	samaAngka(t, h, "PolicyTreatyIn.BalanceDueTo", "1800")
	if h.Ambil("PolicyTreatyIn.DueTo") != "1" {
		t.Errorf("DueTo = %q, harap 1", h.Ambil("PolicyTreatyIn.DueTo"))
	}

	h = halamanOnp()
	h.Setel("PolicyTreatyIn.RiCommOnp", "5")
	h.Setel("PolicyTreatyIn.ResultOnp1", "300")
	if err := CountResult1Onp(h, "Amount"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.RiCommOnp", "15")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "1700")
}

// CountResult1Onp_Act langkah 1 (`.PremiOnp=="" || =="0"` -> 6), langkah 2
// (`isDouble||""` salah -> 6; `>100` -> pesan ErrorMsg1), langkah 3 (150 tidak
// < 100 -> 6): tiga jalan keluar TANPA CountNetPremi (NetPremium tetap).
func TestCountResult1OnpKeluarTanpaMenghitung(t *testing.T) {
	for _, tt := range []struct {
		nama, premi, ri string
		pesan           bool
	}{
		{"PremiOnp kosong", "", "10", false},
		{"PremiOnp 0", "0", "10", false},
		{"RiCommOnp bukan angka", "2000", "UJI", false},
		{"RiCommOnp > 100", "2000", "150", true},
		{"RiCommOnp tepat 100", "2000", "100", false},
	} {
		h := halamanOnp()
		h.Setel("PolicyTreatyIn.PremiOnp", tt.premi)
		h.Setel("PolicyTreatyIn.RiCommOnp", tt.ri)
		h.Setel("PolicyTreatyIn.ResultOnp1", "7")
		h.Setel("PolicyTreatyIn.NetPremium", "5")
		if err := CountResult1Onp(h, "Pct"); err != nil {
			t.Fatalf("%s: %v", tt.nama, err)
		}
		if h.Ambil("PolicyTreatyIn.ResultOnp1") != "7" || h.Ambil("PolicyTreatyIn.NetPremium") != "5" {
			t.Errorf("%s: harus keluar tanpa menghitung, ResultOnp1=%q NetPremium=%q", tt.nama,
				h.Ambil("PolicyTreatyIn.ResultOnp1"), h.Ambil("PolicyTreatyIn.NetPremium"))
		}
		got := h.Pesan["PolicyTreatyIn.RiCommOnp"]
		if tt.pesan && (len(got) != 1 || got[0] != PesanErrorMsg1) {
			t.Errorf("%s: harap pesan [ErrorMsg1] pada RiCommOnp, dapat %v", tt.nama, h.SemuaPesan())
		}
		if !tt.pesan && adaPesan(h) {
			t.Errorf("%s: harap tanpa pesan, dapat %v", tt.nama, h.SemuaPesan())
		}
	}
}

// CountResult2Onp_act (Param.Data):
//
//	2  [@String.isDouble(.OveriddingCommOnp)] F->6 - kosong KELUAR
//	3  [.OveriddingCommOnp<100 && isDouble] F->6; [Data=="Pct"]
//	   .ResultOnp2 = (.OveriddingCommOnp/@String.toDecimal("100.00000"))*.PremiOnp
//	4  [Data=="Amount"] .OveriddingCommOnp = (.ResultOnp2/.PremiOnp) *@String.toDecimal("100.00000")
//	5  call CountNetPremi_act
//
// Pct 5: ResultOnp2 = 5/100 x 2000 = 100; NetPremium = 2000 - (5 x 2000/100) = 1900.
// Amount: OveriddingCommOnp = 50/2000 x 100 = 2,5; NetPremium = 2000 - 2,5 x 2000/100 = 1950.
func TestCountResult2OnpPctAmountDanKosong(t *testing.T) {
	h := halamanOnp()
	h.Setel("PolicyTreatyIn.OveriddingCommOnp", "5")
	if err := CountResult2Onp(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.ResultOnp2", "100")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "1900")

	h = halamanOnp()
	h.Setel("PolicyTreatyIn.OveriddingCommOnp", "1")
	h.Setel("PolicyTreatyIn.ResultOnp2", "50")
	if err := CountResult2Onp(h, "Amount"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.OveriddingCommOnp", "2.5")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "1950")

	h = halamanOnp()
	h.Setel("PolicyTreatyIn.ResultOnp2", "9")
	h.Setel("PolicyTreatyIn.NetPremium", "5")
	if err := CountResult2Onp(h, "Pct"); err != nil {
		t.Fatal(err)
	}
	if h.Ambil("PolicyTreatyIn.ResultOnp2") != "9" || h.Ambil("PolicyTreatyIn.NetPremium") != "5" {
		t.Errorf("OveriddingCommOnp kosong: langkah 2 keluar, dapat ResultOnp2=%q NetPremium=%q",
			h.Ambil("PolicyTreatyIn.ResultOnp2"), h.Ambil("PolicyTreatyIn.NetPremium"))
	}
}

// CountResult2Ogp_act langkah 4 (Data "Amount"):
// .OveriddingCommOgp = (.ResultOgp2/.PremiOgp)*100 = 30/1000 x 100 = 3;
// NetPremium = 1000 - 0 - 3 x 1000/100 = 970.
func TestCountResult2OgpAmount(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisProporsional)
	h.Setel("PolicyTreatyIn.PremiOgp", "1000")
	h.Setel("PolicyTreatyIn.OveriddingCommOgp", "1")
	h.Setel("PolicyTreatyIn.ResultOgp2", "30")
	if err := CountResult2Ogp(h, "Amount"); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.OveriddingCommOgp", "3")
	samaAngka(t, h, "PolicyTreatyIn.NetPremium", "970")
}

// CountNetPremi_act langkah 3 (hanya NonProportional):
// .GrossClaim = @divide(.Claim,TreatyIn.RNMShare,4) * 100 = 30/40 x 100 = 75
// (langkah 2 ber-RNMShareP dilewati: 30/50 x 100 = 60 BUKAN hasilnya).
func TestCountNetPremiNonPropMemakaiRNMShare(t *testing.T) {
	h := HalamanBaru()
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", JenisNonProporsional)
	h.Setel("PolicyTreatyIn.Claim", "30")
	h.Setel("TreatyIn.RNMShare", "40")
	h.Setel("TreatyIn.RNMShareP", "50")
	if err := CountNetPremi(h); err != nil {
		t.Fatal(err)
	}
	samaAngka(t, h, "PolicyTreatyIn.GrossClaim", "75")
}
