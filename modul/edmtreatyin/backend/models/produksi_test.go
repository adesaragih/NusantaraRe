package models

// Uji murni Utility1 EDM `SaveJsonPolisTreatyInEDM_Act`: json_polis tanpa JSON (`SavePolisTreatyInEDM_SQL`),
// ACHIEVEMENT (`SetEDMAchivementValue`), TREATYINPRODUCTION (`InsetTreatyInProdAddendum_Act` +
// `InsertTreatyInProdEDMT_SQL`). Nilai harapan ditulis tangan dari XML korpus EDM; fixture UJI-. Nilai selisih
// NEGATIF (pembatalan `SetEDMTCancel`: Diff = 0 - Old) dan nilai generasi baru BERBEDA dari selisihnya, supaya
// sumber yang salah (nilai baru, daftar NB) langsung merah.

import (
	"reflect"
	"testing"
	"time"
)

const (
	ujiIDEDM   = "UJI-EDMT-1"
	ujiNoPolis = "UJI-QR.T01.10.2026.00001"
	ujiNoEDM   = "UJI-QR.T01.10.2026.00001/E01"
)

// halamanProduksiEDMUji - kasus Proporsional: nilai baru (PolicyTreatyIn.*) sengaja berbeda dari selisih.
func halamanProduksiEDMUji() *Halaman {
	h := HalamanBaru()
	for k, v := range map[string]string{
		"PolicyNo": ujiNoPolis, "EDMNo": ujiNoEDM, "CedingCoName": "UJI CEDING;", "SOBName": "UJI SOB",
		"InsuredName": "UJI INSURED", "TreatyGroupID": "UJI-TG", "TreatyGroupName": "UJI GRUP", "TreatyYear": "2026",
		"YearOfQuartal": "2027", "NoOffer": "UJI-OFFER", "TreatyType": "SURPLUS", "Quartal": "4", "CedingCo": "UJI-C1",
		"SOB": "UJI-S1", "InsuredID": "UJI-I1", "IsNewPolicyNonProp": "0", "ShareValue": "12.5", "OutstandingClaim": "3",
		"StatementType": "SOA", "StartDate": "2026-07-01", "EndDate": "2027-06-30", "Currency": "USD",
		"IDCurrency": "UJI-USD", "MarketingOfficer": "UJI NAMA MO", "ProductionDate": "2026-10-06 17:07:03",
		"BizCode": "", "ClaimType": "UJI-CT", "ClaimPaymentType": "UJI-CPT", "Claim": "9",
		// nilai generasi baru - TIDAK boleh sampai ke produksi jalur Prop
		"PremiOgp": "9999", "ResultOgp1": "9999", "Deduction1": "9999", "NetPremium": "9999", "DueTo": "1",
		"GuaranteeFund": "77", "LayerType": "UJI-LT-BARU", "Layer": "9", "BalanceDueTo": "9999",
	} {
		h.Setel(pt+k, v)
	}
	h.Setel(HalamanPolis+".QuotationData.ProportionalType", ProporsionalPenuh)
	h.Setel(HalamanPolis+".QuotationData.MOID", "UJI-MO")
	h.Setel(HalamanPolis+".QuotationData.BusinessCode", "UJI-B1")
	for k, v := range map[string]string{
		"PremiOgp": "-1000.5", "RiCommOgp": "30", "ResultOgp1": "-300.15", "OveriddingCommOgp": "2.5",
		"ResultOgp2": "-25", "ExcessLoss": "-1", "PremiOnp": "-500", "RiCommOnp": "10", "ResultOnp1": "-50",
		"OveriddingCommOnp": "1", "ResultOnp2": "-5", "Deduction1": "-7", "Deduction2": "-8", "BalanceDueTo": "-600",
		"BalanceBeforePPH": "-610", "BalanceBeforeTax": "-620", "PPHValue": "-2", "PPNValue": "-11",
		"NetPremium": "-640", "Claim": "-4",
	} {
		h.Setel(sd+k, v)
	}
	h.SetelDaftar(DaftarSelisihSpreading, []Baris{
		{"TreatyType": "UJI-JR1", "SharePercentage": "88.8889", "ClaimPercentage": "80", "PremiumSpreaded": "-15828.75", "ClaimSpreaded": "-1.5"},
		{"TreatyType": "UJI-JR2", "SharePercentage": "11.1111", "ClaimPercentage": "20", "PremiumSpreaded": "-989.2969", "ClaimSpreaded": "0"},
	})
	// spreading generasi baru TIGA baris - jalur Prop membaca daftar SELISIH (dua baris)
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-X"}, {"TreatyType": "UJI-Y"}, {"TreatyType": "UJI-Z"}})
	return h
}

// halamanProduksiEDMNonPropUji - kasus XOL: dua induk selisih per mata uang (2 + 1 lapisan).
func halamanProduksiEDMNonPropUji() *Halaman {
	h := halamanProduksiEDMUji()
	h.Setel(pt+"IsNewPolicyNonProp", "1")
	h.Setel(pt+"TreatyType", "XOL")
	h.Setel(HalamanPolis+".QuotationData.ProportionalType", "NonProportional")
	h.SetelDaftar(jMaster+"Share", []Baris{{"SpreadingTypeIDXOL": "UJI-JR-XOL", "SpreadingTypeXOL": "UJI NOTA XOL"}})
	h.SetelDaftar(DaftarSelisihXOL, []Baris{
		{"Currency": "USD", "IDCurrency": "UJI-USD", "GrossPremi": "-1500", "Deduction": "-15", "NetPremi": "-1350"},
		{"Currency": "IDR", "IDCurrency": "UJI-IDR", "GrossPremi": "200", "Deduction": "2", "NetPremi": "180"},
	})
	h.SetelDaftar(JalurAnak(DaftarSelisihXOL, 1, "ValueList"), []Baris{
		{"LayerType": "UJI-LT", "Layer": "1", "LayerPartType": "UJI-PT", "LayerPart": "1", "Currency": "USD",
			"GrossPremi": "-1000", "DueToValue": "-800", "NetPremi": "-900", "Deduction": "-10", "PPHValue": "-2", "PPNValue": "-3"},
		{"LayerType": "UJI-LT", "Layer": "2", "LayerPartType": "UJI-PT", "LayerPart": "1", "Currency": "USD",
			"GrossPremi": "-500", "DueToValue": "-400", "NetPremi": "-450", "Deduction": "-5", "PPHValue": "-1", "PPNValue": "-1.1"},
	})
	h.SetelDaftar(JalurAnak(DaftarSelisihXOL, 2, "ValueList"), []Baris{
		{"LayerType": "UJI-LT", "Layer": "1", "LayerPartType": "UJI-PT", "LayerPart": "1", "Currency": "IDR",
			"GrossPremi": "200", "DueToValue": "160", "NetPremi": "180", "Deduction": "2", "PPHValue": "0.4", "PPNValue": "0.44"},
	})
	// TreatyXOLList generasi baru - NB membacanya, EDM tidak
	h.SetelDaftar(DaftarXOL, []Baris{{"Currency": "EUR", "IDCurrency": "UJI-EUR", "GrossPremi": "1"}})
	h.SetelDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL), []Baris{{"Currency": "EUR", "GrossPremi": "1"}})
	return h
}

func cocokBaris(t *testing.T, nama string, dapat, harap Baris) {
	t.Helper()
	for k, v := range harap {
		if g, ada := dapat[k]; !ada || g != v {
			t.Errorf("%s %s = %q (ada %v), harap %q", nama, k, g, ada, v)
		}
	}
	for k := range dapat {
		if _, ada := harap[k]; !ada {
			t.Errorf("%s memuat kunci tak terduga %s", nama, k)
		}
	}
}

func kunciKolom(ks []Kolom) map[string]bool {
	m := map[string]bool{}
	for _, k := range ks {
		m[k.Kolom] = true
	}
	return m
}

func TestProduksiPrasimpanSaveJsonEDM(t *testing.T) {
	h := halamanProduksiEDMUji()
	h.Setel(pt+"ClaimType", "")
	h.Setel(pt+"ClaimPaymentType", "")
	h.Setel(pt+"Claim", "0")
	h.Setel(pt+"StatementDate", "2026-10-01 00:00:00")
	PrasimpanPolis(h, ujiIDEDM, time.Date(2026, 10, 6, 17, 7, 3, 0, time.Local), 25)
	if got := h.Ambil(pt + "ProductionDate"); got != "2026-10-06 17:07:03" { // 4.1 sekarang; 4.2/4.3 tak berlaku
		t.Errorf("ProductionDate = %q", got)
	}
	if got := h.Ambil(pt + "IDNewBisnis"); got != ujiIDEDM { // 7
		t.Errorf("IDNewBisnis = %q", got)
	}
	// 8: keduanya kosong -> dihitung; Claim 0, SalvageValue kosong -> ""
	if h.Ambil(pt+"ClaimType") != "" || h.Ambil(pt+"ClaimPaymentType") != "" {
		t.Errorf("8 kosong: %q/%q", h.Ambil(pt+"ClaimType"), h.Ambil(pt+"ClaimPaymentType"))
	}
	if got := h.Ambil(pt + "BizCode"); got != "UJI-B1" { // InsetTreatyInProdAddendum_Act 9.1
		t.Errorf("BizCode = %q", got)
	}
	// 4.3 lewat closing -> tanggal 1 bulan berikut 05:00 GMT (penyimpangan sadar jam 24 tetap)
	PrasimpanPolis(h, ujiIDEDM, time.Date(2026, 10, 27, 9, 0, 0, 0, time.Local), 25)
	if got := h.Ambil(pt + "ProductionDate"); got != "2026-11-01 05:00:00" {
		t.Errorf("ProductionDate lewat closing = %q", got)
	}
}

func TestProduksiClaimTypeLangkah8BersyaratEDM(t *testing.T) {
	for _, c := range []struct {
		nama, ct, cpt, claim, salvage, harapCT, harapCPT string
	}{
		// When AKTIF: keduanya terisi -> DIBIARKAN walau Claim terisi (NB menghitung ulang selalu)
		{"keduanya terisi", "UJI-CT", "UJI-CPT", "9", "", "UJI-CT", "UJI-CPT"},
		// satu kosong -> keduanya DITIMPA
		{"CPT kosong, salvage", "UJI-CT", "", "0", "5", "SOA", "Claim"},
		{"CT kosong, claim", "", "UJI-CPT", "9", "", "SOA", "Claim"},
		{"CT kosong, nihil", "", "UJI-CPT", "0", "0", "", ""},
	} {
		h := halamanProduksiEDMUji()
		h.Setel(pt+"ClaimType", c.ct)
		h.Setel(pt+"ClaimPaymentType", c.cpt)
		h.Setel(pt+"Claim", c.claim)
		h.Setel(pt+"SalvageValue", c.salvage)
		PrasimpanMedan(h, ujiIDEDM)
		if h.Ambil(pt+"ClaimType") != c.harapCT || h.Ambil(pt+"ClaimPaymentType") != c.harapCPT {
			t.Errorf("%s: %q/%q, harap %q/%q", c.nama, h.Ambil(pt+"ClaimType"), h.Ambil(pt+"ClaimPaymentType"), c.harapCT, c.harapCPT)
		}
		if h.Ambil(pt+"ProductionDate") != "2026-10-06 17:07:03" {
			t.Errorf("%s: PrasimpanMedan mengubah ProductionDate", c.nama)
		}
	}
	// NonProp: 9.1 di luar blok 10 - BizCode dibiarkan kosong
	h := halamanProduksiEDMNonPropUji()
	PrasimpanMedan(h, ujiIDEDM)
	if h.Ambil(pt+"BizCode") != "" {
		t.Errorf("NonProp BizCode = %q", h.Ambil(pt+"BizCode"))
	}
}

func TestProduksiJSONPolisEDM(t *testing.T) {
	b := JSONPolis(halamanProduksiEDMUji(), ujiIDEDM, "UJI-DH")
	cocokBaris(t, "json_polis", b, Baris{
		"IDPEGA": KunciInstans(ujiIDEDM), "NOPOLIS": ujiNoPolis, "NOENDORS": ujiNoEDM,
		"PRODKE":   "", // TreatyInSearchProdKe - diisi penulis
		"TGL_PROD": "2026-10-06 17:07:03", "USERNAME": "UJI-DH",
	})
	var nama []string
	for _, k := range KolomJSONPolis {
		nama = append(nama, k.Kolom)
	}
	// urutan argumen PEGA_JSON_POLIS_TREATYIN tanpa json; DATA_JSON tidak ditulis (perintah WO 06-10-2026)
	if !reflect.DeepEqual(nama, []string{"IDPEGA", "NOPOLIS", "NOENDORS", "PRODKE", "TGL_PROD", "USERNAME"}) {
		t.Errorf("KolomJSONPolis = %v", nama)
	}
	if !reflect.DeepEqual(kunciKolom(KolomJSONPolis), map[string]bool{"IDPEGA": true, "NOPOLIS": true,
		"NOENDORS": true, "PRODKE": true, "TGL_PROD": true, "USERNAME": true}) {
		t.Error("kunci baris json_polis harus = katalog")
	}
}

// urutanInsertTreatyInProdEDMT - daftar kolom INSERT `RDBList/InsertTreatyInProdEDMT_SQL.xml`, disalin apa adanya.
var urutanInsertTreatyInProdEDMT = []string{"IDPEGA", "NOPOLIS", "BUSINESSCODE", "CEDINGCO", "SOB", "INSUREDNAME",
	"TREATYGROUP", "DUE_TO", "STATEMENT_DATE", "BEGINDATE", "ENDDATE", "CURR_ID", "UW_YEAR", "PREMI_OGP",
	"PERCENT_RI_COMM_OGP", "RI_COMM_OGP", "PERCENT_OVERRIDING_COMM_OGP", "OVERRIDING_COMM_OGP", "CLAIM", "EXCESS_LOSS",
	"PREMI_ONP", "PERCENT_RI_COMM_ONP", "RI_COMM_ONP", "PERCENT_OVERRIDING_COMM_ONP", "OVERRIDING_COMM_ONP",
	"DEDUCTION1", "DEDUCTION2", "NET_PREMIUM", "BALANCE_DUE_TO", "MARKETINGOFFICERCODE", "NOOFFER", "JN_REAS",
	"PCT_SHARE_PREMI", "PCT_SHARE_CLAIM", "TREATYTYPE", "LAYERTYPE", "LAYER", "LAYERPARTTYPE", "LAYERPART", "CLAIMTYPE",
	"CLAIMPAYMENTTYPE", "QUARTER", "NOENDORS", "PPHVALUE", "PPNVALUE", "BALANCE_BEFORE_PPH", "BALANCE_BEFORE_TAX",
	"CEDINGCOID", "SOBID", "INSUREDID", "JENIS_TREATY", "QUARTER_YEAR", "GUARANTEE_FUND", "SHARE_VALUE",
	"PROPORTIONALTYPE", "OUTSTANDING_CLAIM", "PROD_DATE", "STATEMENT_TYPE"}

func TestProduksiKatalogUrutanInsertTreatyInProdEDMT(t *testing.T) {
	var nama []string
	for _, k := range KolomProduksi {
		nama = append(nama, k.Kolom)
	}
	if len(nama) != 58 || !reflect.DeepEqual(nama, urutanInsertTreatyInProdEDMT) {
		t.Errorf("KolomProduksi (%d) = %v", len(nama), nama)
	}
	if kunciKolom(KolomProduksi)["TREATYGROUPID"] {
		t.Error("InsertTreatyInProdEDMT_SQL tidak memuat TREATYGROUPID")
	}
}

// harapProp - baris pertama jalur Prop (blok 9) atas halamanProduksiEDMUji sesudah PrasimpanPolis, 58 kolom.
func harapProp() Baris {
	return Baris{
		"IDPEGA": KunciInstans(ujiIDEDM), "NOPOLIS": ujiNoPolis, "BUSINESSCODE": "UJI-B1", "CEDINGCO": "UJI CEDING",
		"SOB": "UJI SOB", "INSUREDNAME": "UJI INSURED", "TREATYGROUP": "UJI GRUP",
		"DUE_TO":         "", // TreatyDifference.DueTo tak pernah diisi rule EDM (PolicyTreatyIn.DueTo = 1 diabaikan)
		"STATEMENT_DATE": "2026-07-01", "BEGINDATE": "2026-07-01", "ENDDATE": "2027-06-30",
		"CURR_ID": "", // TreatyDifference.Currency (PolicyTreatyIn.Currency = USD diabaikan)
		"UW_YEAR": "2026", "PREMI_OGP": "-1000.5", "PERCENT_RI_COMM_OGP": "30", "RI_COMM_OGP": "-300.15",
		"PERCENT_OVERRIDING_COMM_OGP": "2.5", "OVERRIDING_COMM_OGP": "-25", "CLAIM": "-1.5", "EXCESS_LOSS": "-1",
		"PREMI_ONP": "-500", "PERCENT_RI_COMM_ONP": "10", "RI_COMM_ONP": "-50", "PERCENT_OVERRIDING_COMM_ONP": "1",
		"OVERRIDING_COMM_ONP": "-5", "DEDUCTION1": "-7", "DEDUCTION2": "-8", "NET_PREMIUM": "-15828.75",
		"BALANCE_DUE_TO": "-600", "MARKETINGOFFICERCODE": "UJI NAMA MO", "NOOFFER": "UJI-OFFER", "JN_REAS": "UJI-JR1",
		"PCT_SHARE_PREMI": "88.8889", "PCT_SHARE_CLAIM": "80", "TREATYTYPE": "SURPLUS", "LAYERTYPE": "0",
		"LAYER": "0", "LAYERPARTTYPE": "0", "LAYERPART": "0", "CLAIMTYPE": "UJI-CT", "CLAIMPAYMENTTYPE": "UJI-CPT",
		"QUARTER": "4", "NOENDORS": ujiNoEDM, "PPHVALUE": "-2", "PPNVALUE": "-11", "BALANCE_BEFORE_PPH": "-610",
		"BALANCE_BEFORE_TAX": "-620", "CEDINGCOID": "UJI-C1", "SOBID": "UJI-S1", "INSUREDID": "UJI-I1",
		"JENIS_TREATY": "0", "QUARTER_YEAR": "2027",
		"GUARANTEE_FUND": "", // TreatyDifference.GuaranteeFund (PolicyTreatyIn.GuaranteeFund = 77 diabaikan)
		"SHARE_VALUE":    "12.5", "PROPORTIONALTYPE": ProporsionalPenuh, "OUTSTANDING_CLAIM": "3",
		"PROD_DATE": "2026-10-06 17:07:03", "STATEMENT_TYPE": "SOA",
	}
}

func TestProduksiPropKolomPerKolomSelisihNegatif(t *testing.T) {
	h := halamanProduksiEDMUji()
	PrasimpanPolis(h, ujiIDEDM, time.Date(2026, 10, 6, 17, 7, 3, 0, time.Local), 25)
	rows := BarisProduksi(h, ujiIDEDM)
	if len(rows) != 2 { // satu per TreatyDifference.SpreadingRiskList, bukan per spreading baru (3)
		t.Fatalf("harap 2 baris, dapat %d", len(rows))
	}
	harap := harapProp()
	if !reflect.DeepEqual(func() map[string]bool {
		m := map[string]bool{}
		for k := range harap {
			m[k] = true
		}
		return m
	}(), kunciKolom(KolomProduksi)) {
		t.Fatal("harapan Prop harus memuat tepat 58 kolom katalog")
	}
	cocokBaris(t, "Prop[0]", rows[0], harap)
	h2 := harapProp()
	h2["CLAIM"], h2["NET_PREMIUM"], h2["JN_REAS"], h2["PCT_SHARE_PREMI"], h2["PCT_SHARE_CLAIM"] = "0", "-989.2969", "UJI-JR2", "11.1111", "20"
	cocokBaris(t, "Prop[1]", rows[1], h2)

	// 9.3 TreatyType XOL: layer TreatyDifference (kosong), BUKAN layer polis baru
	h.Setel(pt+"TreatyType", "XOL")
	if r := BarisProduksi(h, ujiIDEDM)[0]; r["LAYERTYPE"] != "" || r["LAYER"] != "" || r["LAYERPARTTYPE"] != "" || r["LAYERPART"] != "" {
		t.Errorf("XOL Prop layer = %q %q %q %q", r["LAYERTYPE"], r["LAYER"], r["LAYERPARTTYPE"], r["LAYERPART"])
	}
	// 9.4 / 9.5 membaca TreatyDifference.DueTo bila suatu saat terisi
	for dueTo, teks := range map[string]string{"1": "DUE TO US", "0": "DUE TO YOU", "": "", "2": ""} {
		h.Setel(sd+"DueTo", dueTo)
		if got := BarisProduksi(h, ujiIDEDM)[0]["DUE_TO"]; got != teks {
			t.Errorf("TreatyDifference.DueTo %q -> DUE_TO %q, harap %q", dueTo, got, teks)
		}
	}
	h.SetelDaftar(DaftarSelisihSpreading, nil)
	if n := len(BarisProduksi(h, ujiIDEDM)); n != 0 {
		t.Errorf("tanpa spreading selisih nol baris, dapat %d", n)
	}
}

// harapNonProp - baris pertama jalur NonProp (blok 10) atas halamanProduksiEDMNonPropUji, 58 kolom.
func harapNonProp() Baris {
	b := harapProp()
	for k, v := range map[string]string{
		"BUSINESSCODE": "", "DUE_TO": "DUE TO US", "CURR_ID": "USD", "PREMI_OGP": "-1000", "PERCENT_RI_COMM_OGP": "0",
		"RI_COMM_OGP": "0", "PERCENT_OVERRIDING_COMM_OGP": "0", "OVERRIDING_COMM_OGP": "0", "CLAIM": "0",
		"EXCESS_LOSS": "0", "PREMI_ONP": "", "PERCENT_RI_COMM_ONP": "0", "RI_COMM_ONP": "0",
		"PERCENT_OVERRIDING_COMM_ONP": "0", "OVERRIDING_COMM_ONP": "0", "DEDUCTION1": "-10", "DEDUCTION2": "0",
		"NET_PREMIUM": "-900", "BALANCE_DUE_TO": "-800", "JN_REAS": "UJI-JR-XOL", "PCT_SHARE_PREMI": "100",
		"PCT_SHARE_CLAIM": "0", "TREATYTYPE": "XOL", "LAYERTYPE": "UJI-LT", "LAYER": "1", "LAYERPARTTYPE": "UJI-PT",
		"LAYERPART": "1", "PPHVALUE": "-2", "PPNVALUE": "-3", "BALANCE_BEFORE_PPH": "", "BALANCE_BEFORE_TAX": "",
		"JENIS_TREATY": "1", "GUARANTEE_FUND": "", "PROPORTIONALTYPE": "NonProportional",
	} {
		b[k] = v
	}
	return b
}

func TestProduksiNonPropPerLapisanTanpaPerkalianSpreading(t *testing.T) {
	h := halamanProduksiEDMNonPropUji()
	PrasimpanMedan(h, ujiIDEDM)
	rows := BarisProduksi(h, ujiIDEDM)
	if len(rows) != 3 { // 2 + 1 lapisan TreatyXOLDifferenceList; NB: x 3 baris spreading
		t.Fatalf("harap 3 baris, dapat %d", len(rows))
	}
	cocokBaris(t, "NonProp[0]", rows[0], harapNonProp())
	h1 := harapNonProp()
	h1["LAYER"], h1["PREMI_OGP"], h1["BALANCE_DUE_TO"], h1["NET_PREMIUM"] = "2", "-500", "-400", "-450"
	h1["DEDUCTION1"], h1["PPHVALUE"], h1["PPNVALUE"] = "-5", "-1", "-1.1"
	cocokBaris(t, "NonProp[1]", rows[1], h1)
	h2 := harapNonProp()
	h2["CURR_ID"], h2["PREMI_OGP"], h2["BALANCE_DUE_TO"], h2["NET_PREMIUM"] = "IDR", "200", "160", "180"
	h2["DEDUCTION1"], h2["PPHVALUE"], h2["PPNVALUE"] = "2", "0.4", "0.44"
	cocokBaris(t, "NonProp[2]", rows[2], h2)
	if NotaJenisReasXOL(h) != "" { // SpreadingTypeIDXOL terisi -> 10.3 dilewati
		t.Errorf("nota = %q", NotaJenisReasXOL(h))
	}

	// AdjPremi (EDMType 3) - jalur produksi sama persis
	adj := halamanProduksiEDMNonPropUji()
	adj.Setel(pt+"EDMType", "3")
	PrasimpanMedan(adj, ujiIDEDM)
	if !reflect.DeepEqual(BarisProduksi(adj, ujiIDEDM), rows) {
		t.Error("AdjPremi harus menghasilkan baris produksi yang sama")
	}
}

func TestProduksiJenisReasMenurutNotaNonProp(t *testing.T) {
	h := halamanProduksiEDMNonPropUji()
	h.SetelDaftar(jMaster+"Share", []Baris{{"SpreadingTypeIDXOL": "", "SpreadingTypeXOL": "UJI NOTA XOL"}})
	s, err := SusunSimpananPolis(h, ujiIDEDM, "UJI-DH")
	if err != nil {
		t.Fatal(err)
	}
	if s.NotaJenisReasXOL != "UJI NOTA XOL" {
		t.Fatalf("nota = %q", s.NotaJenisReasXOL)
	}
	asli := s.Produksi
	for _, r := range s.Produksi {
		if r["JN_REAS"] != "" {
			t.Errorf("sebelum dicari JN_REAS = %q", r["JN_REAS"])
		}
	}
	s.TerapkanJenisReasXOL("7") // 10.3.3
	for i, r := range s.Produksi {
		if r["JN_REAS"] != "7" {
			t.Errorf("baris %d JN_REAS = %q", i, r["JN_REAS"])
		}
	}
	if asli[0]["JN_REAS"] != "" {
		t.Error("TerapkanJenisReasXOL tidak boleh mengubah baris milik pemanggil")
	}
	// tak ditemukan -> "" (pxResults(1).CARI1 kosong)
	s.TerapkanJenisReasXOL("")
	if s.Produksi[0]["JN_REAS"] != "" {
		t.Errorf("tak ditemukan JN_REAS = %q", s.Produksi[0]["JN_REAS"])
	}

	// keduanya kosong / tanpa Share / jalur Prop -> tidak dicari, penerapan tak berefek
	for nama, hh := range map[string]*Halaman{
		"keduanya kosong": func() *Halaman {
			x := halamanProduksiEDMNonPropUji()
			x.SetelDaftar(jMaster+"Share", []Baris{{}})
			return x
		}(),
		"tanpa Share": func() *Halaman {
			x := halamanProduksiEDMNonPropUji()
			x.SetelDaftar(jMaster+"Share", nil)
			return x
		}(),
		"Prop": func() *Halaman {
			x := halamanProduksiEDMUji()
			x.SetelDaftar(jMaster+"Share", []Baris{{"SpreadingTypeXOL": "UJI NOTA XOL"}})
			return x
		}(),
	} {
		s, _ := SusunSimpananPolis(hh, ujiIDEDM, "UJI-DH")
		if s.NotaJenisReasXOL != "" {
			t.Errorf("%s: nota = %q", nama, s.NotaJenisReasXOL)
		}
		jr := s.Produksi[0]["JN_REAS"]
		s.TerapkanJenisReasXOL("7")
		if s.Produksi[0]["JN_REAS"] != jr {
			t.Errorf("%s: JN_REAS berubah %q -> %q", nama, jr, s.Produksi[0]["JN_REAS"])
		}
	}
}

func TestProduksiCapaianPropDanNonPropEDM(t *testing.T) {
	h := halamanProduksiEDMUji()
	c := BarisCapaian(h, ujiIDEDM)
	if len(c) != 1 {
		t.Fatalf("Prop: 1 baris, dapat %d", len(c))
	}
	harap := Baris{
		"IDPEGA": KunciInstans(ujiIDEDM), "NOPOLIS": "", "NOOFFER": "UJI-OFFER", "SOBNAME": "UJI SOB",
		"TREATYGROUPNAME": "UJI GRUP", "TREATYTYPE": "SURPLUS", "QUARTER": "4",
		"QUARTERYEAR": "2026", // TreatyYear, bukan YearOfQuartal (2027)
		"IDCURRENCY":  "UJI-USD", "CURRENCY": "USD", "PREMIUM": "-1000.5", "RICOMM": "-300.15", "BROKERAGE": "-7",
		"NETPREMIUM": "-640", "PAIDCLAIM": "-4", "OUTSTANDINGCLAIM": "3", "PROPORTIONALTYPE": ProporsionalPenuh,
	}
	if !reflect.DeepEqual(func() map[string]bool {
		m := map[string]bool{}
		for k := range harap {
			m[k] = true
		}
		return m
	}(), kunciKolom(KolomCapaian)) {
		t.Fatal("harapan ACHIEVEMENT harus memuat tepat 17 kolom katalog")
	}
	cocokBaris(t, "capaian Prop", c[0], harap)

	n := halamanProduksiEDMNonPropUji()
	c = BarisCapaian(n, ujiIDEDM)
	if len(c) != 2 { // per induk TreatyXOLDifferenceList (USD, IDR); TreatyXOLList baru (EUR) diabaikan
		t.Fatalf("NonProp: 2 baris, dapat %d", len(c))
	}
	h0 := Baris{}
	for k, v := range harap {
		h0[k] = v
	}
	for k, v := range map[string]string{"IDCURRENCY": "UJI-USD", "CURRENCY": "USD", "PREMIUM": "-1500", "RICOMM": "",
		"BROKERAGE": "-15", "NETPREMIUM": "-1350", "PAIDCLAIM": "9", // PolicyTreatyIn.Claim, bukan selisih (-4)
		"TREATYTYPE": "XOL", "PROPORTIONALTYPE": "NonProportional"} {
		h0[k] = v
	}
	cocokBaris(t, "capaian NonProp[0]", c[0], h0)
	if c[1]["CURRENCY"] != "IDR" || c[1]["IDCURRENCY"] != "UJI-IDR" || c[1]["PREMIUM"] != "200" || c[1]["NETPREMIUM"] != "180" {
		t.Errorf("capaian NonProp[1] = %v", c[1])
	}
}

func TestProduksiSusunSimpananPolisEDM(t *testing.T) {
	h := halamanProduksiEDMUji()
	PrasimpanPolis(h, ujiIDEDM, time.Date(2026, 10, 6, 17, 7, 3, 0, time.Local), 25)
	s, err := SusunSimpananPolis(h, ujiIDEDM, "UJI-DH")
	if err != nil {
		t.Fatal(err)
	}
	if s.IDPega != KunciInstans(ujiIDEDM) || s.JSONPolis["NOENDORS"] != ujiNoEDM || len(s.Produksi) != 2 ||
		len(s.Capaian) != 1 || s.NotaJenisReasXOL != "" {
		t.Errorf("simpanan = %+v", s)
	}
	for _, r := range s.Produksi {
		if r["PROD_DATE"] != s.JSONPolis["TGL_PROD"] || r["NOENDORS"] != s.JSONPolis["NOENDORS"] {
			t.Errorf("PROD_DATE/NOENDORS produksi harus = json_polis: %q %q", r["PROD_DATE"], r["NOENDORS"])
		}
	}
}
