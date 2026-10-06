package models_test

// Uji seam 3 - `InputPolicyTreatyInDetail_NonProp` dan langkah 18
// `InputPolicyTreatyInDetail_preACT`. Nilai harapan dihitung tangan dari
// langkah XML (INVENTARIS-XML.md bab 6) atas `masterUji` (nonprop_test.go).

import (
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestInputDetailNonProp(t *testing.T) {
	h := halamanXOL("true", "Inclusive")
	h.Setel("TreatyIn.Commencement", "2026-10-01 00:00:00")
	h.Setel("TreatyIn.Termination", "2027-09-30 00:00:00")
	h.SetelDaftar("PolicyTreatyIn.ListInstallment", []models.Baris{{}, {}, {}})
	h.SetelDaftar("PolicyTreatyIn.ListInstallment(3).InstallmentList", []models.Baris{{"Premium": "UJI-LAMA"}})
	h.TambahPesan("", "UJI-pesan-lama")
	if err := models.InputDetailNonProp(h, idUji); err != nil {
		t.Fatal(err)
	}
	pol := func(m string) string { return h.Ambil("PolicyTreatyIn." + m) }
	// 13
	for m, harap := range map[string]string{"DueTo": "1", "StartDate": "2026-10-01 00:00:00", "EndDate": "2027-09-30 00:00:00"} {
		if pol(m) != harap {
			t.Errorf("%s = %q, harap %q", m, pol(m), harap)
		}
	}
	// 16 (hanya IDR = PolicyTreatyIn.Currency): gross 1000+2000, net 900+1800,
	// deduction DeductionTotalList 306.6, share RnmLimitList 5000+7000.
	// 16.2.2 (FlagPPH, Inclusive, presisi 4): BFS 306.6/1.022 = 300 -> PPN 6.6, PPH 6.
	// 18 + 19: BalanceDueTo = netpremi + PPN + PPH.
	for m, harap := range map[string]string{
		"PremiOgp": "3000", "ShareValue": "12000", "Deduction1": "306.6", "Deduction2": "0", "NetPremium": "2700",
		"PPNValue": "6.6", "PPHValue": "6", "BalanceDueTo": "2712.6", "BalanceBeforePPH": "2706.6", "BalanceBeforeTax": "2700",
	} {
		sama(t, m, pol(m), harap)
	}
	// 20: satu ListInstallment per TreatyIn.Installment (baris lama dibuang, 18).
	ang := h.AmbilDaftar("PolicyTreatyIn.ListInstallment")
	if len(ang) != 2 {
		t.Fatalf("ListInstallment %d baris, harap 2", len(ang))
	}
	cekBaris(t, "Angsuran(1)", ang[0], map[string]string{"PaymentTotal": "2700", "InstallmentPercentage": "100", "Currency": "IDR", "IDCurrency": "UJI-ID-IDR"})
	cekBaris(t, "Angsuran(2)", ang[1], map[string]string{"PaymentTotal": "45", "InstallmentPercentage": "100", "Currency": "USD", "IDCurrency": "UJI-ID-USD"})
	r := h.AmbilDaftar("PolicyTreatyIn.ListInstallment(1).InstallmentList")
	if len(r) != 2 {
		t.Fatalf("InstallmentList(1) %d baris", len(r))
	}
	// 20.4.1 tanpa retro: Premium = .Amount; DueDate = .PaymentDate; InstallmentNo = .Installment.
	cekBaris(t, "Rinci(1,1)", r[0], map[string]string{"Premium": "1080", "Currency": "IDR", "IDCurrency": "UJI-ID-IDR",
		"InstallmentPercentage": "40", "InstallmentNo": "1", "DueDate": "2026-11-01"})
	cekBaris(t, "Rinci(1,2)", r[1], map[string]string{"Premium": "1620", "InstallmentPercentage": "60", "InstallmentNo": "2", "DueDate": "2027-02-01"})
	if len(h.AmbilDaftar("PolicyTreatyIn.ListInstallment(3).InstallmentList")) != 0 {
		t.Error("rincian baris angsuran lama ikut terhapus (langkah 18)")
	}
	// 22 (FacultativeShare 0) -> InsertToTreatyXOLList; 23 -> spreading.
	if n := len(h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")); n != 2 {
		t.Errorf("TreatyXOLList %d baris, harap 2", n)
	}
	cekBaris(t, "Spreading(1)", h.AmbilDaftar("PolicyTreatyIn.SpreadingRiskList")[0], map[string]string{
		"TreatyType": "UJI-SPR-ID", "PremiumSpreaded": "2700",
	})
	if len(h.SemuaPesan()) != 0 {
		t.Errorf("langkah 1 membersihkan pesan lama: %v", h.SemuaPesan())
	}
	if h.Ambil("PolicyTreatyIn.IsEDMInputOnNB") != "" {
		t.Error("EDMState kosong: IsEDMInputOnNB tidak diisi (10)")
	}
}

func TestInputDetailNonPropRetro(t *testing.T) {
	m := masterUji()
	m.Nilai["EDMState"] = "2"
	m.Daftar["FacultativeShareList"] = []models.Baris{{}}
	m.Daftar["FacultativeShareList(1).DeductionTotalList"] = []models.Baris{{"Currency": "IDR", "Value": "50"}, {"Currency": "USD", "Value": "9"}}
	m.Daftar["TotalFacShareDeductionNP"] = []models.Baris{{"Currency": "IDR", "Value": "400"}}
	m.Daftar["LimitFacShareSummaryList"] = []models.Baris{{"Note": "UJI", "Deductible": "40", "Deductible2": "4"}}
	m.Daftar["Share(1).SpreadingListXOL"] = []models.Baris{{"Pct": "30"}, {"Pct": "70"}}
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.Currency", "IDR")
	h.Setel("PolicyTreatyIn.FlagPPH", "true")
	h.Setel("PolicyTreatyIn.TypeTax", "Inclusive")
	h.Setel("PolicyTreatyIn.FlagRetroTreaty", "true")
	models.TerapkanMasterXOL(h, m)
	if err := models.InputDetailNonProp(h, idUji); err != nil {
		t.Fatal(err)
	}
	// 17: deduction 306.6 + retro IDR 50 = 356.6 (USD diabaikan).
	// 18: PremiOgp, NetPremium = deduction (retro); 19: BalanceDueTo tetap
	// local.netpremi 2700 + PPN 6.6 + PPH 6 (PPN dari 16.2.2, SEBELUM 17).
	for mm, harap := range map[string]string{
		"PremiOgp": "356.6", "NetPremium": "356.6", "Deduction1": "356.6", "BalanceDueTo": "2712.6",
	} {
		sama(t, mm, h.Ambil("PolicyTreatyIn."+mm), harap)
	}
	// 20.4.1 retro: Premium = deduction * @divide(.InstallmentPct,100,10): 356.6 x 0.4, x 0.6.
	r := h.AmbilDaftar("PolicyTreatyIn.ListInstallment(1).InstallmentList")
	sama(t, "Premium(1)", r[0]["Premium"], "142.64")
	sama(t, "Premium(2)", r[1]["Premium"], "213.96")
	// 7-8: tampilan master retro.
	sama(t, "TotalShareNetNP", h.AmbilDaftar("TreatyIn.TotalShareNetNP")[0]["Value"], "400")
	ls := h.AmbilDaftar("TreatyIn.LimitShareSummaryList")[0]
	cekBaris(t, "LimitShareSummaryList", ls, map[string]string{"NetPremi": "40", "NetPremi2": "4", "Deductible": "0", "Deductible2": "0"})
	// 9: isOR = @substring(Pct "30", 4, 6) bukan "OR" -> TotalSpreadedNetPremi =
	// 400 x Pct(2) 70 / 100 = 280; RI (isOR bukan "R/I") = 400 x 30 / 100 = 120.
	sama(t, "TotalSpreadedNetPremi", h.AmbilDaftar("TreatyIn.TotalSpreadedNetPremi")[0]["Value"], "280")
	sama(t, "TotalSpreadedNetPremiRI", h.AmbilDaftar("TreatyIn.TotalSpreadedNetPremiRI")[0]["Value"], "120")
	// 10: EDMState 2 -> TreatyMasterInEDM -> IsEDMInputOnNB.
	if h.Ambil("PolicyTreatyIn.IsEDMInputOnNB") != "true" {
		t.Error("IsEDMInputOnNB = true bila master EDM (10)")
	}
}

func TestPPNPPHLapisanXOL(t *testing.T) { // InputPolicyTreatyInDetail_preACT 18
	h := halamanXOL("true", "Inclusive")
	m := masterUji()
	m.Daftar["LimitShareSummaryList"] = []models.Baris{
		{"Limit": "1000", "Limit2": "0", "Deductible": "102.2", "Deductible2": "0", "NetPremi": "900", "NetPremi2": "0"},
		{"Limit": "0", "Limit2": "500", "Deductible": "0", "Deductible2": "10.22", "NetPremi": "0", "NetPremi2": "45"},
	}
	m.Daftar["TotalShareNetNP"] = []models.Baris{{"Currency": "IDR", "Value": "2700"}, {"Currency": "USD", "Value": "45"}}
	m.Daftar["TotalShareDeductionNP"] = []models.Baris{{"Currency": "IDR", "Value": "306.6"}, {"Currency": "USD", "Value": "10.22"}}
	models.TerapkanMasterXOL(h, m)
	if err := models.InputDetailNonProp(h, idUji); err != nil {
		t.Fatal(err)
	}
	if err := models.PPNPPHLapisanXOL(h); err != nil {
		t.Fatal(err)
	}
	ls := h.AmbilDaftar("TreatyIn.LimitShareSummaryList")
	// 18.1 (presisi 8): 102.2/1.022 = 100 -> PPN 2.2, PPH 2; 10.22/1.022 = 10 -> 0.22, 0.2.
	cekBaris(t, "Lapis(1)", ls[0], map[string]string{"PPNValue": "2.2", "NetPremiAfterPPN": "902.2", "PPHValue": "2", "NetPremiAfterPPH": "902"})
	cekBaris(t, "Lapis(2)", ls[1], map[string]string{"PPNValue2": "0.22", "NetPremiAfterPPN2": "45.22", "PPHValue2": "0.2", "NetPremiAfterPPH2": "45.2"})
	// 18.2: baris total yang mata uangnya = Local.curr / curr2 baris berjalan.
	tn := h.AmbilDaftar("TreatyIn.TotalShareNetNP")
	cekBaris(t, "TotalNet(IDR)", tn[0], map[string]string{"TotalNetPremiAfterPPN": "902.2", "TotalNetPremiAfterPPH": "902", "TotalNetPremiAfterTax": "904.2"})
	cekBaris(t, "TotalNet(USD)", tn[1], map[string]string{"TotalNetPremiAfterPPN": "45.22", "TotalNetPremiAfterPPH": "45.2", "TotalNetPremiAfterTax": "45.42"})
	td := h.AmbilDaftar("TreatyIn.TotalShareDeductionNP")
	cekBaris(t, "TotalDed(IDR)", td[0], map[string]string{"TotalPPNValue": "2.2", "TotalPPHValue": "2"})
	cekBaris(t, "TotalDed(USD)", td[1], map[string]string{"TotalPPNValue": "0.22", "TotalPPHValue": "0.2"})
	// 18.3.4.1: angsuran per mata uang.
	ang := h.AmbilDaftar("PolicyTreatyIn.ListInstallment")
	cekBaris(t, "Angsuran(IDR)", ang[0], map[string]string{"PPN": "2.2", "PPh": "2", "PaymentTotalAfterPPN": "2702.2", "PaymentTotalAfterTax": "904.2"})
	cekBaris(t, "Angsuran(USD)", ang[1], map[string]string{"PPN": "0.22", "PPh": "0.2", "PaymentTotalAfterPPN": "45.22", "PaymentTotalAfterTax": "45.42"})
	// 18.3.4.2: rincian x InstallmentPercentage/100.
	r := h.AmbilDaftar("PolicyTreatyIn.ListInstallment(1).InstallmentList")
	cekBaris(t, "Rinci(IDR,40%)", r[0], map[string]string{"PPN": "0.88", "PPh": "0.8", "PremiumAfterPPN": "1080.88", "PremiumAfterTax": "361.68"})
	cekBaris(t, "Rinci(IDR,60%)", r[1], map[string]string{"PPN": "1.32", "PPh": "1.2", "PremiumAfterPPN": "1621.32", "PremiumAfterTax": "542.52"})
	cekBaris(t, "Rinci(USD)", h.AmbilDaftar("PolicyTreatyIn.ListInstallment(2).InstallmentList")[0],
		map[string]string{"PPN": "0.22", "PPh": "0.2", "PremiumAfterPPN": "45.22", "PremiumAfterTax": "45.42"})
}

func TestPPNPPHLapisanXOLTanpaFlagPPH(t *testing.T) {
	h := halamanXOL("false", "Inclusive")
	m := masterUji()
	m.Daftar["LimitShareSummaryList"] = []models.Baris{{"Limit": "1", "Deductible": "102.2", "NetPremi": "900"}}
	models.TerapkanMasterXOL(h, m)
	if err := models.PPNPPHLapisanXOL(h); err != nil {
		t.Fatal(err)
	}
	if _, ada := h.AmbilDaftar("TreatyIn.LimitShareSummaryList")[0]["PPNValue"]; ada {
		t.Error("FlagPPH bukan true: langkah 18 dilewati")
	}
}
