package services_test

// Uji terjemahan `SaveTreatyInDetail_Act` / `SaveTreatyInDetailEdm_Act`
// (perintah WO 9 Oktober 2026). Bentuk dokumen = bentuk pendaratan: teks,
// larik `[]any` berisi `map[string]any`. Data UJI-*.

import (
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func larik(xs ...map[string]any) []any {
	out := make([]any, len(xs))
	for i, x := range xs {
		out[i] = x
	}
	return out
}

func uangDetail(cur, v string) map[string]any { return map[string]any{"Currency": cur, "Value": v} }

func angkaDi(t *testing.T, b models.BarisDetailTreaty, kolom string) string {
	t.Helper()
	d, ada := b.Angka[kolom]
	if !ada {
		t.Fatalf("kolom angka %s tidak ada", kolom)
	}
	if d == nil {
		return ""
	}
	return d.Text('f')
}

func dokProp() map[string]any {
	return map[string]any{
		"ProportionType": "Proportional", "TreatyContractName": "UJI-KONTRAK", "TreatyYear": "2026",
		"Ceding": "UJI-CEDANT", "CedingID": "G0", "LeadingReinsSource": "UJI-SOB", "LeadingReinsSourceID": "B0",
		"BrokeragePercentP": "2,5", "RNMShareP": "15",
		"Limits": larik(map[string]any{
			"TreatyType": "QUOTA SHARE",
			"Detail": larik(
				map[string]any{
					"TreatyGroup": "PROPERTY", "TreatyGroupID": "10007",
					"RIOGR": "225", "RIONR": "", "RNMShare": "", "SpreadingTypeID": "10236", "SpreadingType": "UJI QS",
					"Surplus": "3", "QSPct": "50",
					"IOOLimitList":  larik(uangDetail("IDR", "1000"), uangDetail("USD", "10")),
					"RetentionList": larik(uangDetail("IDR", "200")),
					"CessionList":   larik(uangDetail("IDR", "800")),
					"RNMShareList":  larik(uangDetail("IDR", "120")),
					"SpreadingList": larik(
						map[string]any{"ReinsTypeID": "10004", "Pct": "60"},
						map[string]any{"ReinsTypeID": "10028", "Pct": "40"},
					),
					"EPIList": larik(uangDetail("IDR", "5000")),
				},
				// Detail kedua TANPA IOOLimitList → satu sisipan dengan isi
				// SaveData yang masih membawa nilai detail pertama.
				map[string]any{"TreatyGroup": "MARINE", "TreatyGroupID": "10008", "SpreadingTypeID": "10236"},
			),
		}),
	}
}

func TestDetailProporsionalMengikutiActivity(t *testing.T) {
	baris, err := services.SusunDetailTreatyIn(dokProp(), "UJI-1")
	if err != nil {
		t.Fatal(err)
	}
	// IOO IDR → 1 baris ber-EPI; IOO USD → retensi IDR tidak cocok (CEK
	// CURR) → nol; detail kedua → 1 sisipan dini.
	if len(baris) != 2 {
		t.Fatalf("baris %d, mau 2", len(baris))
	}
	b := baris[0]
	for kolom, mau := range map[string]string{
		"TREATYID": "UJI-1", "PROPORTIONTYPE": "Proportional", "SOB": "UJI-SOB", "CEDINGID": "G0",
		"TREATYTYPE": "QUOTA SHARE", "TREATYGROUP": "PROPERTY", "LIMITCURRENCY": "IDR",
		"RETENTIONCURRENCY": "IDR", "CESSIONCURRENCY": "IDR", "SHARECURRENCY": "IDR", "EPICURRENCY": "IDR",
		"SPREADINGTYPEID": "10236", "SPREADINGTYPE": "UJI QS", "CLASSOFBUSINESS": "",
	} {
		if b.Teks[kolom] != mau {
			t.Errorf("%s = %q, mau %q", kolom, b.Teks[kolom], mau)
		}
	}
	for kolom, mau := range map[string]string{
		"LIMITVALUE": "1000", "RETENTIONVALUE": "200", "CESSIONVALUE": "800", "SHAREVALUE": "120", "EPIVALUE": "5000",
		// `@divide(@toDecimal(225),10,2)`; kosong = `@toDecimal("")` = 0.
		"RIOGR": "22.50", "RIONR": "0",
		// `.RNMShare` kosong → `TreatyIn.RNMShareP`; ketikan "2,5" terbaca 2.5.
		"RNM_SHARE": "15", "BROKERAGE": "2.5",
		"QS_RI": "60", "QS_OR": "40", "SPL_LINE": "3", "QS_PCT": "50",
		// Bukan milik cabang ini.
		"MDPVALUE": "", "SPREAD_RNM_SHARE_VALUE": "",
	} {
		if got := angkaDi(t, b, kolom); got != mau {
			t.Errorf("%s = %q, mau %q", kolom, got, mau)
		}
	}
	// ⭐ SaveData tidak dikosongkan antarbaris — sisipan dini detail kedua
	// membawa mata uang/nilai IOO terakhir (USD) dan EPI detail pertama.
	k := baris[1]
	if k.Teks["TREATYGROUP"] != "MARINE" || k.Teks["LIMITCURRENCY"] != "USD" || angkaDi(t, k, "EPIVALUE") != "5000" {
		t.Errorf("sisipan dini %v / %v", k.Teks, k.Angka)
	}
}

func TestDetailProporsionalSpreadingLebihDariSatu(t *testing.T) {
	doc := dokProp()
	lim := doc["Limits"].([]any)[0].(map[string]any)
	det := lim["Detail"].([]any)[0].(map[string]any)
	lim["Detail"] = larik(det)
	det["SpreadingTypeID"] = ""
	det["EPIList"] = larik()
	det["SpreadingList"] = larik(
		map[string]any{"ReinsTypeID": "10195", "ReinsTypeName": "UJI SPL", "Pct": "30", "BreakDownSprdList": larik(
			map[string]any{"ReinsID": "10004", "SharePct": "70", "Amount": "100"},
			map[string]any{"ReinsID": "10028", "SharePct": "30", "Amount": "50,5"},
		)},
		map[string]any{"ReinsTypeID": "10196", "ReinsTypeName": "UJI SPL 2", "Pct": "70"},
	)
	baris, err := services.SusunDetailTreatyIn(doc, "UJI-2")
	if err != nil {
		t.Fatal(err)
	}
	// Satu sisipan per baris SpreadingList (EPIList kosong).
	if len(baris) != 2 {
		t.Fatalf("baris %d, mau 2", len(baris))
	}
	if baris[0].Teks["SPREADINGTYPEID"] != "10195" || angkaDi(t, baris[0], "SPREAD_RNM_SHARE_VALUE") != "150.5" ||
		angkaDi(t, baris[0], "SPREAD_RNM_SHARE_PCT") != "30" || angkaDi(t, baris[0], "QS_RI") != "70" {
		t.Errorf("spreading 1: %v %v", baris[0].Teks, baris[0].Angka)
	}
	// Baris kedua tanpa rincian: nilai disetel ulang ke 0, QS ikut terbawa.
	if baris[1].Teks["SPREADINGTYPE"] != "UJI SPL 2" || angkaDi(t, baris[1], "SPREAD_RNM_SHARE_VALUE") != "0" ||
		angkaDi(t, baris[1], "QS_OR") != "30" {
		t.Errorf("spreading 2: %v %v", baris[1].Teks, baris[1].Angka)
	}
}

func shareNP(sxol string) map[string]any {
	return map[string]any{
		"Layer": "1", "LayerPart": "1", "LayerType": "layer", "LayerPartType": "layer",
		"SpreadingTypeXOL": sxol, "SpreadingTypeIDXOL": "10260",
		"TreatyGroupList":  larik(map[string]any{"TreatyGroup": "PROPERTY", "TreatyGroupID": "10007"}),
		"RnmLimitList":     larik(uangDetail("IDR", "7000")),
		"GrossPremiumList": larik(uangDetail("IDR", "300"), uangDetail("USD", "3")),
		"NetPremiumList":   larik(uangDetail("IDR", "250")),
		"DeductionList": larik(
			map[string]any{"Comment": "Brokerage", "Currency": "IDR", "Deduction": "10"},
			map[string]any{"Comment": "Commission", "Currency": "IDR", "Deduction": "20"},
			map[string]any{"Comment": "Other", "Currency": "IDR", "Deduction": "5"},
			map[string]any{"Comment": "brokerage", "Currency": "USD", "Deduction": "99"},
		),
	}
}

func dokNP(sh map[string]any) map[string]any {
	return map[string]any{
		"ProportionType": "NonProportional", "RNMShare": "19", "BrokeragePercent": "10", "InstallmentNo": "4",
		"Limits": larik(map[string]any{
			"MDPPct": "85", "ROLPct": "5,1", "Deductible": "600", "Deductible2": "0", "AdjRate": "3",
			"PremiumEarnedList": larik(uangDetail("IDR", "1"), uangDetail("IDR", "2")),
			"MDPList":           larik(uangDetail("IDR", "9")),
		}),
		"Share": larik(sh),
	}
}

// ⚠️ MEMAKU EKSPOR APA ADANYA: cabang Treaty In NonProp "spreading type 1"
// tidak punya langkah yang mengisi `SaveData.LIMITCURRENCY`, jadi Gross
// Premium selalu Exit Iteration — nol baris. Bila WO memutuskan lain, uji
// inilah yang berubah.
func TestDetailNonPropTreatyInTypeSatuNolBaris(t *testing.T) {
	baris, err := services.SusunDetailTreatyIn(dokNP(shareNP("2026 UJI QS")), "UJI-3")
	if err != nil {
		t.Fatal(err)
	}
	if len(baris) != 0 {
		t.Fatalf("baris %d, mau 0 (ekspor [5.4.2] tanpa RNM Limit)", len(baris))
	}
}

func TestDetailNonPropTreatyInTypeLebihDariSatu(t *testing.T) {
	sh := shareNP("")
	sh["SpreadingListXOL"] = larik(map[string]any{
		"ReinsTypeID": "10270", "ReinsTypeName": "UJI XOL AS", "Pct": "14", "Value": "980",
		"BreakDownSprdListXOL": larik(map[string]any{"ReinsTypeID": "10004", "SharePct": "100"}),
		"RnmLimitList":         larik(uangDetail("IDR", "980")),
		"GrossPremiumList":     larik(uangDetail("IDR", "42"), uangDetail("USD", "1")),
		"NetPremiumList":       larik(uangDetail("IDR", "35"), uangDetail("USD", "1")),
	})
	baris, err := services.SusunDetailTreatyIn(dokNP(sh), "UJI-4")
	if err != nil {
		t.Fatal(err)
	}
	if len(baris) != 1 {
		t.Fatalf("baris %d, mau 1 (mata uang lain Exit Iteration)", len(baris))
	}
	b := baris[0]
	for kolom, mau := range map[string]string{
		"LIMITCURRENCY": "IDR", "MDPCURRENCY": "IDR", "SPREADINGTYPE": "UJI XOL AS", "SPREADINGTYPEID": "10270",
		"LAYER": "1", "INSTALLMENTNO": "4", "TREATYGROUPID": "10007",
	} {
		if b.Teks[kolom] != mau {
			t.Errorf("%s = %q, mau %q", kolom, b.Teks[kolom], mau)
		}
	}
	for kolom, mau := range map[string]string{
		// LIMITVALUE/SHAREVALUE/NETPREMI dari RnmLimitList & NetPremiumList
		// baris Share ([5.4.3.2]); MDP dari spreading ([5.4.3.3.2]).
		"LIMITVALUE": "7000", "SHAREVALUE": "7000", "NETPREMIVALUE": "250", "MDPVALUE": "42",
		// Brokerage → DEDUCTION1; selainnya → DEDUCTION2; mata uang lain dilewati.
		"DEDUCTION1": "10", "DEDUCTION2": "25",
		"SPREAD_RNM_SHARE_PCT": "14", "SPREAD_RNM_SHARE_VALUE": "980", "QS_RI": "100",
		"RNM_SHARE": "19", "BROKERAGE": "10", "MDP_PCT": "85", "ROL_PCT": "5.1", "DEDUCTIBLE": "600",
		"ADJ_RATE": "3", "PREMIUM_EARNED": "2", "MDP": "9",
	} {
		if got := angkaDi(t, b, kolom); got != mau {
			t.Errorf("%s = %q, mau %q", kolom, got, mau)
		}
	}
}

func TestDetailNonPropEDMMengisiBatasRNM(t *testing.T) {
	baris, err := services.SusunDetailTreatyInEDM(dokNP(shareNP("2026 UJI QS")), "UJI-5/R01")
	if err != nil {
		t.Fatal(err)
	}
	if len(baris) != 1 {
		t.Fatalf("baris %d, mau 1", len(baris))
	}
	b := baris[0]
	if b.Teks["TREATYID"] != "UJI-5/R01" || b.Teks["LIMITCURRENCY"] != "IDR" || b.Teks["SPREADINGTYPE"] != "2026 UJI QS" {
		t.Errorf("teks %v", b.Teks)
	}
	if angkaDi(t, b, "MDPVALUE") != "300" || angkaDi(t, b, "DEDUCTION1") != "10" || angkaDi(t, b, "DEDUCTION2") != "25" {
		t.Errorf("angka %v", b.Angka)
	}
	// Prosedur EDM berhenti di QS_PCT: kolom khusus Treaty In tidak ada.
	if _, ada := b.Angka["MDP_PCT"]; ada {
		t.Error("kolom MDP_PCT ikut di baris EDM")
	}
}

func TestDetailAngkaTakTerbacaDitolak(t *testing.T) {
	doc := dokProp()
	doc["Limits"].([]any)[0].(map[string]any)["Detail"].([]any)[0].(map[string]any)["RIOGR"] = "30%"
	_, err := services.SusunDetailTreatyIn(doc, "UJI-6")
	if !errors.Is(err, services.ErrTombolDitolak) {
		t.Fatalf("galat %v, mau ErrTombolDitolak", err)
	}
}
