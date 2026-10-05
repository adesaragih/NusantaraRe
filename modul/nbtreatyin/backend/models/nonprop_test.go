package models_test

// Uji seam 3 (fungsi murni) jalur NB NonProporsional / XOL. Setiap nilai harapan
// DIHITUNG TANGAN dari langkah XML (INVENTARIS-XML.md bab 6) atas master fiktif
// berawalan UJI-; nomor langkah ditulis di samping tiap harapan.

import (
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// masterUji - dua layer (Share) dan dua angsuran per mata uang (IDR, USD).
//
//	Share(1) layer 1: Gross IDR 1000 / USD 50 · Net IDR 900 / USD 45 · Deduction IDR 102.2 / USD 10.22
//	Share(2) layer 2: Gross IDR 2000          · Net IDR 1800        · Deduction IDR 204.4
//	                  DeductionTotalList IDR 306.6 (pada Share(1)), RnmLimitList IDR 5000 / 7000
func masterUji() models.MasterXOL {
	b := func(kv ...string) models.Baris {
		r := models.Baris{}
		for i := 0; i+1 < len(kv); i += 2 {
			r[kv[i]] = kv[i+1]
		}
		return r
	}
	return models.MasterXOL{
		Nilai: map[string]string{"RNMShare": "10", "FacultativeShare": "0", "ProportionType": "NonProportional"},
		Daftar: map[string][]models.Baris{
			"Share": {
				b("LayerType", "UJI-LT", "Layer", "1", "LayerPartType", "UJI-PT", "LayerPart", "1",
					"SpreadingTypeXOL", "UJI-SPR", "SpreadingTypeIDXOL", "UJI-SPR-ID"),
				b("LayerType", "UJI-LT", "Layer", "2", "LayerPartType", "UJI-PT", "LayerPart", "1"),
			},
			"Share(1).GrossPremiumList":   {b("Currency", "IDR", "Value", "1000"), b("Currency", "USD", "Value", "50")},
			"Share(1).NetPremiumList":     {b("Currency", "IDR", "Value", "900"), b("Currency", "USD", "Value", "45")},
			"Share(1).DeductionList":      {b("Currency", "IDR", "Deduction", "102.2"), b("Currency", "USD", "Deduction", "10.22")},
			"Share(1).DeductionTotalList": {b("Currency", "IDR", "Value", "306.6")},
			"Share(1).RnmLimitList":       {b("Currency", "IDR", "Value", "5000")},
			"Share(2).GrossPremiumList":   {b("Currency", "IDR", "Value", "2000")},
			"Share(2).NetPremiumList":     {b("Currency", "IDR", "Value", "1800")},
			"Share(2).DeductionList":      {b("Currency", "IDR", "Deduction", "204.4")},
			"Share(2).RnmLimitList":       {b("Currency", "IDR", "Value", "7000")},
			"Installment": {
				b("Currency", "IDR", "AmountTotal", "2700", "PctTotal", "100"),
				b("Currency", "USD", "AmountTotal", "45", "PctTotal", "100"),
			},
			"Installment(1).InstallmentList": {
				b("Installment", "1", "PaymentDate", "2026-11-01", "InstallmentPct", "40", "Amount", "1080", "Currency", "IDR"),
				b("Installment", "2", "PaymentDate", "2027-02-01", "InstallmentPct", "60", "Amount", "1620", "Currency", "IDR"),
			},
			"Installment(2).InstallmentList": {
				b("Installment", "1", "PaymentDate", "2026-11-01", "InstallmentPct", "100", "Amount", "45", "Currency", "USD"),
			},
		},
	}
}

var idUji = models.IDMataUang{"IDR": "UJI-ID-IDR", "USD": "UJI-ID-USD"}

func halamanXOL(flagPPH, typeTax string) *models.Halaman {
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.FlagPPH", flagPPH)
	h.Setel("PolicyTreatyIn.TypeTax", typeTax)
	h.Setel("PolicyTreatyIn.Currency", "IDR")
	h.Setel("PolicyTreatyIn.IDCurrency", "UJI-ID-IDR")
	models.TerapkanMasterXOL(h, masterUji())
	return h
}

// sama membandingkan dua teks angka secara NUMERIK ("100.00000000" = "100");
// teks kosong hanya sama dengan teks kosong.
func sama(t *testing.T, apa, dapat, harap string) {
	t.Helper()
	if harap == "" || dapat == "" {
		if dapat != harap {
			t.Errorf("%s = %q, harap %q", apa, dapat, harap)
		}
		return
	}
	a, err1 := models.AngkaTeks(apa, dapat)
	b, err2 := models.AngkaTeks(apa, harap)
	if err1 != nil || err2 != nil {
		if dapat != harap {
			t.Errorf("%s = %q, harap %q", apa, dapat, harap)
		}
		return
	}
	if a.Cmp(b) != 0 {
		t.Errorf("%s = %s, harap %s", apa, a.Text('f'), b.Text('f'))
	}
}

func cekBaris(t *testing.T, nama string, b models.Baris, harap map[string]string) {
	t.Helper()
	for k, v := range harap {
		sama(t, nama+"."+k, b[k], v)
	}
}

func TestInsertToTreatyXOLList(t *testing.T) {
	h := halamanXOL("true", "Inclusive")
	h.SetelDaftar("PolicyTreatyIn.TreatyXOLList", []models.Baris{{"GrossPremi": "UJI-LAMA"}}) // langkah 1 menghapusnya
	if err := models.InsertToTreatyXOLList(h, idUji); err != nil {
		t.Fatal(err)
	}
	xol := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")
	if len(xol) != 2 {
		t.Fatalf("satu baris induk per TreatyIn.Installment (langkah 3), dapat %d", len(xol))
	}
	// IDR: 3.2.4 jumlah per layer; 3.5 + 3.6 (FlagPPH, Inclusive):
	// Deduction 102.2+204.4 = 306.6 -> BrokerageFeeSebenarnya 306.6/1.022 = 300,
	// PPH 300*0.02 = 6, PPN 300*0.022 = 6.6.
	// ⚠️ DueTo "" - `Page-New InputXOL` (3.2.1) menghapus CARI1 "DUE TO US" langkah 2.
	cekBaris(t, "XOL(1)", xol[0], map[string]string{
		"Currency": "IDR", "IDCurrency": "UJI-ID-IDR", "GrossPremi": "3000", "NetPremi": "2700",
		"DueToValue": "2700", "Deduction": "306.6", "DueTo": "", "BrokerageFeeSebenarnya": "300",
		"PPHValue": "6", "PPNValue": "6.6", "NetPremiAfterPPH": "2706", "NetPremiAfterPPN": "2706.6",
		"NetPremiAfterTax": "2712.6",
	})
	// USD: Share(2) tanpa baris USD -> CARI31 kosong di layer TERAKHIR ->
	// local.currency = "" (3.2.4 dibaca apa adanya). Deduction 10.22 -> 10.
	cekBaris(t, "XOL(2)", xol[1], map[string]string{
		"Currency": "", "IDCurrency": "UJI-ID-USD", "GrossPremi": "50", "NetPremi": "45",
		"DueToValue": "45", "Deduction": "10.22", "BrokerageFeeSebenarnya": "10",
		"PPHValue": "0.2", "PPNValue": "0.22", "NetPremiAfterPPH": "45.2", "NetPremiAfterPPN": "45.22",
		"NetPremiAfterTax": "45.42",
	})
	l1 := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList(1).ValueList")
	if len(l1) != 2 {
		t.Fatalf("satu baris ValueList per layer (3.7), dapat %d", len(l1))
	}
	// 3.7.4 + 3.7.5: CARI45 = 102.2/1.022 = 100; PPH 2; PPN 2.2; AfterTax 900+(2+2.2).
	cekBaris(t, "XOL(1).Layer(1)", l1[0], map[string]string{
		"Currency": "IDR", "IDCurrency": "UJI-ID-IDR", "GrossPremi": "1000", "NetPremi": "900",
		"DueToValue": "900", "Deduction": "102.2", "LayerType": "UJI-LT", "Layer": "1",
		"LayerPartType": "UJI-PT", "LayerPart": "1", "DueTo": "", "BrokerageFeeSebenarnya": "100",
		"PPHValue": "2", "PPNValue": "2.2", "NetPremiAfterPPH": "902", "NetPremiAfterPPN": "902.2",
		"NetPremiAfterTax": "904.2",
	})
	cekBaris(t, "XOL(1).Layer(2)", l1[1], map[string]string{
		"GrossPremi": "2000", "NetPremi": "1800", "Deduction": "204.4", "Layer": "2",
		"BrokerageFeeSebenarnya": "200", "PPHValue": "4", "PPNValue": "4.4", "NetPremiAfterTax": "1808.4",
	})
	l2 := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList(2).ValueList")
	// Layer 2 tanpa baris USD: CARI kosong; PPH = toDecimal("")*0.02 = 0.
	cekBaris(t, "XOL(2).Layer(2)", l2[1], map[string]string{
		"Currency": "", "GrossPremi": "", "NetPremi": "", "Deduction": "", "BrokerageFeeSebenarnya": "",
		"PPHValue": "0", "PPNValue": "0", "NetPremiAfterPPH": "0", "NetPremiAfterTax": "0",
	})
}

func TestInsertToTreatyXOLListTanpaPPH(t *testing.T) {
	// FlagPPH bukan "true": langkah 3.6 dan 3.7.5 dilewati; TypeTax selain
	// "Inclusive" -> potongan dipakai apa adanya (tidak terlihat tanpa PPH).
	h := halamanXOL("false", "")
	if err := models.InsertToTreatyXOLList(h, idUji); err != nil {
		t.Fatal(err)
	}
	x := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")[0]
	for _, m := range []string{"BrokerageFeeSebenarnya", "PPHValue", "PPNValue", "NetPremiAfterTax"} {
		if _, ada := x[m]; ada {
			t.Errorf("%s tidak boleh terisi tanpa FlagPPH", m)
		}
	}
	sama(t, "GrossPremi", x["GrossPremi"], "3000")
}

func angka(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, err := models.AngkaTeks("uji", s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

// halamanRetro - master uji + FacultativeShareList(1) (layer 1 saja):
// GrossPremiumList IDR 100, DeductionList IDR 10.
func halamanRetro(isNew, retro string) *models.Halaman {
	m := masterUji()
	m.Nilai["FacultativeShare"] = "5"
	m.Daftar["FacultativeShareList"] = []models.Baris{{}}
	m.Daftar["FacultativeShareList(1).GrossPremiumList"] = []models.Baris{{"Currency": "IDR", "Value": "100"}}
	m.Daftar["FacultativeShareList(1).DeductionList"] = []models.Baris{{"Currency": "IDR", "Deduction": "10"}}
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.FlagPPH", "true")
	h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", isNew)
	h.Setel("PolicyTreatyIn.FlagRetroTreaty", retro)
	models.TerapkanMasterXOL(h, m)
	return h
}

func TestRetroShareHanyaUntukPolisNonPropBaru(t *testing.T) {
	h := halamanRetro("0", "")
	h.SetelDaftar("PolicyTreatyIn.TreatyXOLList", []models.Baris{{"GrossPremi": "UJI-LAMA"}})
	if err := models.InsertToTreatyXOLListRetroShare(h, idUji); err != nil {
		t.Fatal(err)
	}
	// langkah 1 selalu menghapus; langkah 2 bersyarat IsNewPolicyNonProp==1
	if n := len(h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")); n != 0 {
		t.Fatalf("IsNewPolicyNonProp 0: daftar XOL kosong, dapat %d", n)
	}
}

func TestRetroShareTanpaFlagRetro(t *testing.T) {
	h := halamanRetro("1", "false")
	if err := models.InsertToTreatyXOLListRetroShare(h, idUji); err != nil {
		t.Fatal(err)
	}
	xol := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")
	if len(xol) != 2 {
		t.Fatalf("dua angsuran master -> dua induk, dapat %d", len(xol))
	}
	// 2.3.3.4 tanpa retro: jumlah layer. Induk RetroShare TIDAK punya langkah PPH.
	cekBaris(t, "XOL(1)", xol[0], map[string]string{
		"Currency": "IDR", "IDCurrency": "UJI-ID-IDR", "GrossPremi": "3000", "NetPremi": "2700",
		"DueToValue": "2700", "Deduction": "306.6", "DueTo": "",
	})
	if _, ada := xol[0]["PPHValue"]; ada {
		t.Error("induk RetroShare tidak punya langkah FlagPPH")
	}
	// ⚠️ `Page-New InputXOL` hanya per angsuran (2.3.2): CARI layer 1 (USD) TERBAWA
	// ke layer 2 yang tak ber-USD -> dijumlah dua kali: 50+50, 45+45, 10.22+10.22.
	cekBaris(t, "XOL(2)", xol[1], map[string]string{
		"Currency": "USD", "GrossPremi": "100", "NetPremi": "90", "DueToValue": "90", "Deduction": "20.44",
	})
	l := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList(1).ValueList")
	// 2.3.7.4 (tanpa syarat): gross retro IDR 100 + 1000; di DALAM perulangannya net
	// dan potongan layer dibaca lagi dan DITAMBAHKAN: 900+900, 102.2+102.2.
	// 2.3.7.6: CARI45 tidak pernah diisi RetroShare -> BrokerageFeeSebenarnya "", PPH/PPN 0.
	cekBaris(t, "XOL(1).Layer(1)", l[0], map[string]string{
		"Currency": "IDR", "GrossPremi": "1100", "NetPremi": "1800", "DueToValue": "1800",
		"Deduction": "204.4", "Layer": "1", "DueTo": "", "BrokerageFeeSebenarnya": "", "PPHValue": "0", "PPNValue": "0",
	})
	cekBaris(t, "XOL(1).Layer(2)", l[1], map[string]string{
		"GrossPremi": "2000", "NetPremi": "1800", "DueToValue": "1800", "Deduction": "204.4", "Layer": "2",
	})
}

func TestRetroShareDenganFlagRetro(t *testing.T) {
	h := halamanRetro("1", "true")
	if err := models.InsertToTreatyXOLListRetroShare(h, idUji); err != nil {
		t.Fatal(err)
	}
	// 2.3.3.3: CARI44 = potongan retro 10 + 102.2 = 112.2 (layer 1); layer 2 = 204.4.
	// 2.3.3.4 retro: gross/net/dueto = CARI44 layer TERAKHIR (bukan jumlah), deduction 0.
	cekBaris(t, "XOL(1)", h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList")[0], map[string]string{
		"GrossPremi": "204.4", "NetPremi": "204.4", "DueToValue": "204.4", "Deduction": "0",
	})
	l := h.AmbilDaftar("PolicyTreatyIn.TreatyXOLList(1).ValueList")
	// 2.3.7.4.4: CARI19 = 10 + CARI44 (204.4) = 214.4; 2.3.7.5 retro -> CARI19.
	cekBaris(t, "XOL(1).Layer(1)", l[0], map[string]string{
		"GrossPremi": "214.4", "NetPremi": "214.4", "DueToValue": "214.4", "Deduction": "204.4",
	})
	// layer 2 tanpa baris retro: CARI19 kosong.
	cekBaris(t, "XOL(1).Layer(2)", l[1], map[string]string{
		"GrossPremi": "", "NetPremi": "", "DueToValue": "", "Deduction": "204.4",
	})
}

func TestTreatyNonPropSetSpreadingDariJenisXOL(t *testing.T) {
	h := halamanXOL("", "")
	h.Setel("PolicyTreatyIn.NetPremium", "2700")
	h.SetelDaftar("PolicyTreatyIn.SpreadingRiskList", []models.Baris{{"TreatyType": "UJI-LAMA"}})
	if err := models.TreatyNonPropSetSpreading(h); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar("PolicyTreatyIn.SpreadingRiskList")
	if len(d) != 1 {
		t.Fatalf("langkah 1 menghapus, langkah 3 mengisi baris 1: %d baris", len(d))
	}
	// 3: SpreadingTypeXOL terisi -> TreatyType = SpreadingTypeIDXOL, 100%;
	// 10 -> CountSpreading_Act: PremiumSpreaded = 2700 * 100/100.
	cekBaris(t, "Spreading(1)", d[0], map[string]string{
		"TreatyType": "UJI-SPR-ID", "Currency": "IDR", "CurrencyID": "UJI-ID-IDR",
		"SharePercentage": "100", "ClaimPercentage": "100", "PremiumSpreaded": "2700", "ClaimSpreaded": "0",
	})
	if len(h.SemuaPesan()) != 0 {
		t.Fatalf("SpreadingTypeIDXOL terisi: tanpa pesan, dapat %v", h.SemuaPesan())
	}
}

func TestTreatyNonPropSetSpreadingDariDaftarXOL(t *testing.T) {
	m := masterUji()
	m.Daftar["Share"][0]["SpreadingTypeXOL"] = ""
	m.Daftar["Share"][0]["SpreadingTypeIDXOL"] = ""
	m.Daftar["Share(1).SpreadingListXOL"] = []models.Baris{
		{"ReinsTypeID": "UJI-R1", "Pct": "6"}, {"ReinsTypeID": "UJI-R2", "Pct": "4"},
	}
	h := models.HalamanBaru()
	h.Setel("PolicyTreatyIn.Currency", "IDR")
	h.Setel("PolicyTreatyIn.NetPremium", "2700")
	models.TerapkanMasterXOL(h, m)
	if err := models.TreatyNonPropSetSpreading(h); err != nil {
		t.Fatal(err)
	}
	d := h.AmbilDaftar("PolicyTreatyIn.SpreadingRiskList")
	if len(d) != 2 {
		t.Fatalf("4.1: satu baris per SpreadingListXOL, dapat %d", len(d))
	}
	// 4.1: SharePercentage = @divide(.Pct, RNMShare 10, 20)*100 -> 60 dan 40;
	// SplitRNMSharePct = .Pct; 10 -> PremiumSpreaded = 2700*60% / 40%.
	cekBaris(t, "Spreading(1)", d[0], map[string]string{
		"TreatyType": "UJI-R1", "SplitRNMSharePct": "6", "SharePercentage": "60", "PremiumSpreaded": "1620",
	})
	cekBaris(t, "Spreading(2)", d[1], map[string]string{
		"TreatyType": "UJI-R2", "SplitRNMSharePct": "4", "SharePercentage": "40", "PremiumSpreaded": "1080",
	})
	// 7: SpreadingTypeIDXOL kosong dan FacultativeShare 0 -> pesan halaman (VERBATIM).
	if p := h.SemuaPesan(); len(p) != 1 || p[0] != "Spreading in Master data is incomplete" {
		t.Fatalf("pesan %v", p)
	}
}

func TestTreatyNonPropSetSpreadingRetro(t *testing.T) {
	h := halamanXOL("", "")
	h.Setel("PolicyTreatyIn.FlagRetroTreaty", "true")
	h.Setel("PolicyTreatyIn.NetPremium", "500")
	if err := models.TreatyNonPropSetSpreading(h); err != nil {
		t.Fatal(err)
	}
	// 6: FlagRetroTreaty -> baris 1 "ORS" 100%.
	cekBaris(t, "Spreading(1)", h.AmbilDaftar("PolicyTreatyIn.SpreadingRiskList")[0], map[string]string{
		"TreatyType": "ORS", "SharePercentage": "100", "PremiumSpreaded": "500",
	})
}
