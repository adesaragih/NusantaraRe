package models

// Uji murni Utility1 `SaveJsonPolisTreatyIn_Act` (ekspor WO 06-10-2026): json_polis tanpa JSON,
// TREATYINPRODUCTION (`InsetTreatyInProd_Act` + `InsertTreatyInProd_SQL`), ACHIEVEMENT (`SetAchivementValue`).
// Nilai harapan dihitung tangan dari XML; fixture UJI-.

import (
	"testing"
	"time"
)

func halamanProduksiUji() *Halaman {
	h := HalamanBaru()
	for k, v := range map[string]string{
		"PolicyNo": "UJI-QR.T01.10.2026.00001", "CedingCoName": "UJI CEDING;", "SOBName": "UJI SOB", "InsuredName": "UJI INSURED",
		"TreatyGroupID": "UJI-TG", "TreatyGroupName": "UJI GRUP", "TreatyYear": "2026", "NoOffer": "UJI-OFFER",
		"TreatyType": "SURPLUS", "Quartal": "4", "CedingCo": "UJI-C1", "SOB": "UJI-S1", "InsuredID": "UJI-I1",
		"IsNewPolicyNonProp": "0", "YearOfQuartal": "2026", "ShareValue": "12.5", "OutstandingClaim": "3",
		"StatementType": "SOA", "StartDate": "2026-07-01", "EndDate": "2027-06-30", "Currency": "USD", "IDCurrency": "UJI-USD",
		"PremiOgp": "1000.5", "RiCommOgp": "30", "ResultOgp1": "300.15", "OveriddingCommOgp": "2.5", "ResultOgp2": "25",
		"ExcessLoss": "1", "PremiOnp": "500", "RiCommOnp": "10", "ResultOnp1": "50", "OveriddingCommOnp": "1", "ResultOnp2": "5",
		"Deduction1": "7", "Deduction2": "8", "BalanceDueTo": "600", "BalanceBeforePPH": "610", "BalanceBeforeTax": "620",
		"PPHValue": "2", "PPNValue": "11", "DueTo": "1", "Claim": "0", "NetPremium": "640", "ProductionDate": "2026-10-06 17:07:03",
		"BizCode": "",
	} {
		h.Setel(pt+k, v)
	}
	h.Setel(HalamanPolis+".QuotationData.ProportionalType", ProporsionalPenuh)
	h.Setel(HalamanPolis+".QuotationData.MOID", "UJI-MO")
	h.Setel(HalamanPolis+".QuotationData.BusinessCode", "UJI-B1")
	h.SetelDaftar(DaftarSpreading, []Baris{
		{"TreatyType": "UJI-JR1", "SharePercentage": "88.8889", "ClaimPercentage": "88.8889", "PremiumSpreaded": "15828.75", "ClaimSpreaded": "0"},
		{"TreatyType": "UJI-JR2", "SharePercentage": "11.1111", "ClaimPercentage": "11.1111", "PremiumSpreaded": "989.2969", "ClaimSpreaded": "0"},
	})
	return h
}

func TestPrasimpanPolisMenurutSaveJsonPolis(t *testing.T) {
	h := halamanProduksiUji()
	h.Setel(pt+"ClaimType", "UJI-LAMA")
	h.Setel(pt+"StatementDate", "2026-10-01 00:00:00")
	sekarang := time.Date(2026, 10, 6, 17, 7, 3, 0, time.Local)
	PrasimpanPolis(h, "UJI-NB-1", sekarang, 25)

	if got := h.Ambil(pt + "ProductionDate"); got != "2026-10-06 17:07:03" { // 4.1 sekarang; 4.2/4.3 tak berlaku
		t.Errorf("ProductionDate = %q", got)
	}
	if got := h.Ambil(pt + "IDNewBisnis"); got != "UJI-NB-1" { // 7
		t.Errorf("IDNewBisnis = %q", got)
	}
	// 8 - When tak dicentang: SELALU dihitung ulang; Claim 0, SalvageValue kosong -> ""
	if h.Ambil(pt+"ClaimType") != "" || h.Ambil(pt+"ClaimPaymentType") != "" {
		t.Errorf("ClaimType/ClaimPaymentType = %q/%q", h.Ambil(pt+"ClaimType"), h.Ambil(pt+"ClaimPaymentType"))
	}
	if got := h.Ambil(pt + "BizCode"); got != "UJI-B1" { // InsetTreatyInProd_Act 9.1
		t.Errorf("BizCode = %q", got)
	}

	h.Setel(pt+"SalvageValue", "5")
	PrasimpanPolis(h, "UJI-NB-1", time.Date(2026, 10, 27, 9, 0, 0, 0, time.Local), 25)
	if h.Ambil(pt+"ClaimType") != "SOA" || h.Ambil(pt+"ClaimPaymentType") != "Claim" {
		t.Errorf("ada SalvageValue: %q/%q", h.Ambil(pt+"ClaimType"), h.Ambil(pt+"ClaimPaymentType"))
	}
	if got := h.Ambil(pt + "ProductionDate"); got != "2026-11-01 05:00:00" { // 4.3 lewat closing
		t.Errorf("ProductionDate lewat closing = %q", got)
	}
}

func TestJSONPolisTanpaJSON(t *testing.T) {
	b := JSONPolis(halamanProduksiUji(), "UJI-NB-1", "UJI-DH")
	harap := Baris{"IDPEGA": KunciInstans("UJI-NB-1"), "NOPOLIS": "UJI-QR.T01.10.2026.00001", "PRODKE": "0",
		"TGL_PROD": "2026-10-06 17:07:03", "USERNAME": "UJI-DH"}
	if len(b) != len(harap) {
		t.Fatalf("json_polis = %v", b)
	}
	for k, v := range harap {
		if b[k] != v {
			t.Errorf("%s = %q, harap %q", k, b[k], v)
		}
	}
	for _, k := range KolomJSONPolis {
		if k.Kolom == "DATA_JSON" {
			t.Error("DATA_JSON tidak boleh ditulis (perintah WO 06-10-2026)")
		}
	}
}

func TestProduksiProporsionalSatuBarisPerSpreading(t *testing.T) {
	if len(KolomProduksi) != 58 {
		t.Fatalf("InsertTreatyInProd_SQL 58 kolom, katalog %d", len(KolomProduksi))
	}
	h := halamanProduksiUji()
	PrasimpanPolis(h, "UJI-NB-1", time.Date(2026, 10, 6, 17, 7, 3, 0, time.Local), 25)
	rows, err := BarisProduksi(h, "UJI-NB-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("harap 2 baris (satu per spreading), dapat %d", len(rows))
	}
	r := rows[0]
	for kol, v := range map[string]string{
		"IDPEGA": KunciInstans("UJI-NB-1"), "NOPOLIS": "UJI-QR.T01.10.2026.00001", "BUSINESSCODE": "UJI-B1",
		"CEDINGCO": "UJI CEDING", "SOB": "UJI SOB", "TREATYGROUPID": "UJI-TG", "TREATYGROUP": "UJI GRUP",
		"DUE_TO": "DUE TO US", "STATEMENT_DATE": "2026-07-01", "BEGINDATE": "2026-07-01", "ENDDATE": "2027-06-30",
		"CURR_ID": "USD", "UW_YEAR": "2026", "PREMI_OGP": "1000.5", "PERCENT_RI_COMM_OGP": "30", "RI_COMM_OGP": "300.15",
		"PERCENT_OVERRIDING_COMM_OGP": "2.5", "OVERRIDING_COMM_OGP": "25", "CLAIM": "0", "EXCESS_LOSS": "1",
		"PREMI_ONP": "500", "DEDUCTION1": "7", "DEDUCTION2": "8", "NET_PREMIUM": "15828.75", "BALANCE_DUE_TO": "600",
		"MARKETINGOFFICERCODE": "UJI-MO", "NOOFFER": "UJI-OFFER", "JN_REAS": "UJI-JR1", "PCT_SHARE_PREMI": "88.8889",
		"PCT_SHARE_CLAIM": "88.8889", "TREATYTYPE": "SURPLUS", "LAYERTYPE": "0", "LAYER": "0", "LAYERPARTTYPE": "0",
		"LAYERPART": "0", "CLAIMTYPE": "", "QUARTER": "4", "PPHVALUE": "2", "PPNVALUE": "11", "BALANCE_BEFORE_PPH": "610",
		"BALANCE_BEFORE_TAX": "620", "CEDINGCOID": "UJI-C1", "SOBID": "UJI-S1", "INSUREDID": "UJI-I1", "JENIS_TREATY": "0",
		"QUARTER_YEAR": "2026", "GUARANTEE_FUND": "", "SHARE_VALUE": "12.5", "PROPORTIONALTYPE": ProporsionalPenuh,
		"OUTSTANDING_CLAIM": "3", "PROD_DATE": "2026-10-06 17:07:03", "STATEMENT_TYPE": "SOA",
	} {
		if r[kol] != v {
			t.Errorf("%s = %q, harap %q", kol, r[kol], v)
		}
	}
	if rows[1]["JN_REAS"] != "UJI-JR2" || rows[1]["NET_PREMIUM"] != "989.2969" || rows[1]["PCT_SHARE_PREMI"] != "11.1111" {
		t.Errorf("baris spreading kedua = %v", rows[1])
	}
	for _, b := range rows {
		for _, k := range KolomProduksi {
			if _, ada := b[k.Kolom]; !ada {
				t.Errorf("baris tanpa kolom %s", k.Kolom)
			}
		}
	}

	h.Setel(pt+"DueTo", "0")
	h.Setel(pt+"TreatyType", "XOL")
	rows, _ = BarisProduksi(h, "UJI-NB-1")
	if rows[0]["DUE_TO"] != "DUE TO YOU" || rows[0]["LAYERTYPE"] != "" { // 9.5; 9.3 layer polis (tak tersimpan)
		t.Errorf("DueTo 0 / XOL: %q %q", rows[0]["DUE_TO"], rows[0]["LAYERTYPE"])
	}
	h.Setel(pt+"DueTo", "")
	rows, _ = BarisProduksi(h, "UJI-NB-1")
	if rows[0]["DUE_TO"] != "" {
		t.Errorf("DueTo kosong = %q", rows[0]["DUE_TO"])
	}
}

func TestProduksiNonPropPerSpreadingPerLayer(t *testing.T) {
	h := halamanProduksiUji()
	h.Setel(pt+"IsNewPolicyNonProp", "1")
	h.Setel(HalamanPolis+".QuotationData.ProportionalType", "NonProportional")
	h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-JR1", "SharePercentage": "40"}, {"TreatyType": "UJI-JR2", "SharePercentage": "60"}})
	h.SetelDaftar(DaftarXOL, []Baris{{"Currency": "USD", "IDCurrency": "UJI-USD", "GrossPremi": "1000", "Deduction": "10", "NetPremi": "900"}})
	h.SetelDaftar(JalurAnak(DaftarXOL, 1, AnakLayerXOL), []Baris{
		{"LayerType": "UJI-LT", "Layer": "1", "LayerPartType": "UJI-PT", "LayerPart": "1", "Currency": "USD",
			"GrossPremi": "1000", "DueToValue": "800", "NetPremi": "900", "Deduction": "10", "PPHValue": "2", "PPNValue": "3"},
		{"Layer": "2", "GrossPremi": "500", "DueToValue": "", "NetPremi": "450", "Deduction": "5"},
	})
	rows, err := BarisProduksi(h, "UJI-NB-2")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 { // 2 spreading x 1 XOL x 2 layer
		t.Fatalf("harap 4 baris, dapat %d", len(rows))
	}
	r := rows[0]
	for kol, v := range map[string]string{
		"JN_REAS": "UJI-JR1", "PCT_SHARE_PREMI": "40", "PCT_SHARE_CLAIM": "0", "CLAIM": "0", "DUE_TO": "DUE TO US",
		"PERCENT_RI_COMM_OGP": "0", "DEDUCTION2": "0", "LAYERTYPE": "UJI-LT", "LAYER": "1", "CURR_ID": "USD",
		"PREMI_OGP": "400", "BALANCE_DUE_TO": "320", "NET_PREMIUM": "360", "DEDUCTION1": "4", "PPHVALUE": "2",
		"PPNVALUE": "3", "PREMI_ONP": "", "BALANCE_BEFORE_PPH": "", "GUARANTEE_FUND": "", "JENIS_TREATY": "1",
	} {
		if got := normal(r[kol]); got != v {
			t.Errorf("%s = %q, harap %q", kol, r[kol], v)
		}
	}
	if got := normal(rows[1]["BALANCE_DUE_TO"]); got != "0" { // DueToValue kosong = 0 di aritmetika
		t.Errorf("layer 2 BALANCE_DUE_TO = %q", rows[1]["BALANCE_DUE_TO"])
	}
	if got := normal(rows[3]["PREMI_OGP"]); got != "300" || rows[3]["JN_REAS"] != "UJI-JR2" {
		t.Errorf("spreading 2 layer 2 = %q %q", rows[3]["PREMI_OGP"], rows[3]["JN_REAS"])
	}

	h.Setel(pt+"SharePercentage", "")
	h.SetelDaftar(DaftarSpreading, []Baris{{"SharePercentage": "UJI-BUKAN-ANGKA"}})
	if _, err := BarisProduksi(h, "UJI-NB-2"); err == nil {
		t.Error("SharePercentage bukan angka harus galat")
	}
}

func TestCapaianProporsionalDanNonProp(t *testing.T) {
	h := halamanProduksiUji()
	c := BarisCapaian(h, "UJI-NB-1")
	if len(c) != 1 {
		t.Fatalf("proporsional: 1 baris, dapat %d", len(c))
	}
	for kol, v := range map[string]string{
		"IDPEGA": KunciInstans("UJI-NB-1"), "NOPOLIS": "", "NOOFFER": "UJI-OFFER", "SOBNAME": "UJI SOB",
		"TREATYGROUPNAME": "UJI GRUP", "TREATYTYPE": "SURPLUS", "QUARTER": "4", "QUARTERYEAR": "2026",
		"IDCURRENCY": "UJI-USD", "CURRENCY": "USD", "PREMIUM": "1000.5", "RICOMM": "300.15", "BROKERAGE": "7",
		"NETPREMIUM": "640", "PAIDCLAIM": "0", "OUTSTANDINGCLAIM": "3", "PROPORTIONALTYPE": ProporsionalPenuh,
	} {
		if c[0][kol] != v {
			t.Errorf("%s = %q, harap %q", kol, c[0][kol], v)
		}
	}

	h.Setel(HalamanPolis+".QuotationData.ProportionalType", "NonProportional")
	h.SetelDaftar(DaftarXOL, []Baris{
		{"Currency": "USD", "IDCurrency": "UJI-USD", "GrossPremi": "1000", "Deduction": "10", "NetPremi": "900"},
		{"Currency": "IDR", "IDCurrency": "UJI-IDR", "GrossPremi": "5", "Deduction": "1", "NetPremi": "4"},
	})
	c = BarisCapaian(h, "UJI-NB-1")
	if len(c) != 2 || c[1]["CURRENCY"] != "IDR" || c[0]["PREMIUM"] != "1000" || c[0]["RICOMM"] != "" ||
		c[0]["BROKERAGE"] != "10" || c[0]["NETPREMIUM"] != "900" || c[0]["PROPORTIONALTYPE"] != "NonProportional" {
		t.Errorf("non-proporsional per TreatyXOLList = %v", c)
	}
}

// normal membuang nol ekor desimal (`400.000...` -> `400`) supaya harapan ditulis seperti angka layar.
func normal(s string) string {
	d, err := AngkaTeks("uji", s)
	if err != nil || s == "" {
		return s
	}
	d.Reduce(d)
	return d.Text('f')
}

func TestPrasimpanMedanMembiarkanTanggalProduksi(t *testing.T) {
	h := halamanProduksiUji()
	h.Setel(pt+"ProductionDate", "2026-10-06 17:07:03")
	h.Setel(pt+"Claim", "9")
	PrasimpanMedan(h, "UJI-NB-1")
	if h.Ambil(pt+"ProductionDate") != "2026-10-06 17:07:03" {
		t.Errorf("ProductionDate berubah: %q", h.Ambil(pt+"ProductionDate"))
	}
	if h.Ambil(pt+"IDNewBisnis") != "UJI-NB-1" || h.Ambil(pt+"ClaimType") != "SOA" || h.Ambil(pt+"BizCode") != "UJI-B1" {
		t.Errorf("langkah 7/8/9.1: %q %q %q", h.Ambil(pt+"IDNewBisnis"), h.Ambil(pt+"ClaimType"), h.Ambil(pt+"BizCode"))
	}
}
