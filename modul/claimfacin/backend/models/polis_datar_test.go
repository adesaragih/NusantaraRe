package models_test

// Uji dokumen polis dari tabel datar NB Fac In (perintah work owner 10-10-2026: "pertama cek json_polis, jika json nya
// kosong ambil dari flat nya"): pohon PARENT_ID / SRC_PATH menjadi bentuk DATA_JSON Pega, lalu TerapkanPolis menulis
// halaman yang sama dengan polis Pega. Fixture `UJI-`.

import (
	"encoding/json"
	"testing"

	"nusantarare/modul/claimfacin/backend/models"
)

func TestDokumenPolisBerisi(t *testing.T) {
	for dok, mau := range map[string]bool{
		`{"QuotationData":{"BusinessType":"Fire"}}`: true,
		`{"PolicyData":{"PolicyNo":"UJI-1"}}`:       true,
		`{"LocationList":[]}`:                       true,
		`{"CaseIDNewNB":"NB-1"}`:                    false, // penunjuk kasus sistem baru (DEV: {"xxxxxxxxxxx":"xx-999999"})
		`{}`:                                        false,
		``:                                          false,
		`bukan json`:                                true, // galat urai tetap terlihat (TerapkanPolis)
	} {
		if got := models.DokumenPolisBerisi([]byte(dok)); got != mau {
			t.Errorf("DokumenPolisBerisi(%s) = %v, mau %v", dok, got, mau)
		}
	}
}

func TestKasusPenunjukPolis(t *testing.T) {
	if k := models.KasusPenunjukPolis([]byte(`{"CaseIDNewNB":"NB-184351"}`)); k != "NB-184351" {
		t.Fatalf("penunjuk: %q", k)
	}
	if k := models.KasusPenunjukPolis([]byte(`{"QuotationData":{}}`)); k != "" {
		t.Fatalf("dokumen polis bukan penunjuk: %q", k)
	}
}

func TestJenisBisnisNBDecisionTable(t *testing.T) {
	// BusinessType_DeT (SetBusinessType_Act 2): GroupPanel x BusinessOldId, baris pertama yang cocok; bawaan "UNKNOWN".
	for _, c := range [][3]string{
		{"006", "78", "OilGas"}, {"006", "22", "FireStyle1"}, {"006", "34", "FireStyle2"}, {"006", "UJ", "Fire"},
		{"003", "16", "AllRisk"}, {"003", "36", "Bonding"}, {"003", "A1", "BondingKBG"}, {"003", "28", "FireStyle1"},
		{"003", "UJ", "Aneka"}, {"007", "02", "MBUCar"}, {"007", "60", "MBUMotorCycle"}, {"007", "L5", "UNKNOWN"},
		{"002", "03", "PA"}, {"002", "UJ", "UNKNOWN"}, {"004", "", "MarineCargo"}, {"001", "", "Medicare"},
		{"009", "L3", "Life"}, {"009", "", "UNKNOWN"}, {"", "", "UNKNOWN"},
	} {
		if got := models.JenisBisnisNB(c[0], c[1]); got != c[2] {
			t.Errorf("JenisBisnisNB(%q, %q) = %q, mau %q", c[0], c[1], got, c[2])
		}
	}
}

// bacaUji - pembaca tabel datar tiruan: baris per (tabel, induk), disaring PARENT_ID.
func bacaUji(baris map[[2]string][]models.BarisPolisDatar) func(models.SisiPolisDatar, []string) ([]models.BarisPolisDatar, error) {
	return func(s models.SisiPolisDatar, induk []string) ([]models.BarisPolisDatar, error) {
		ada := map[string]bool{}
		for _, id := range induk {
			ada[id] = true
		}
		var out []models.BarisPolisDatar
		for _, b := range baris[[2]string{s.Tabel, s.Induk}] {
			if ada[b.Induk] {
				out = append(out, b)
			}
		}
		return out, nil
	}
}

func TestSusunPolisDatarFire(t *testing.T) {
	const cov = "LocationList/Property/PropertyItemList/CoverageList"
	akar := models.BarisPolisDatar{ID: "UJI-NB-1", Nilai: map[string]string{"PolicyData.StartDateTime": "20260101T050000.000 GMT",
		"PolicyData.EndDateTime": "20261231T050000.000 GMT", "PercentShare": "25"}}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_QUOTATIONDATA", "T_GENERAL_POLIS"}: {{ID: "11", Induk: "UJI-NB-1", Nilai: map[string]string{"NoOfferSlip": "UJI-OFR",
			"SobName": "UJI SOB", "CedingCoName": "UJI CEDING", "MOID": "UJI-MO"}}},
		{"T_CEDINGCOLIST", "T_QUOTATIONDATA"}: {{ID: "21", Induk: "11", Nilai: map[string]string{"CedingCoName": "UJI CEDING"}}},
		{"T_LOCATIONLIST", "T_GENERAL_POLIS"}: {{ID: "31", Induk: "UJI-NB-1"}, {ID: "32", Induk: "UJI-NB-1"}},
		{"T_PROPERTY", "T_LOCATIONLIST"}: {
			{ID: "41", Induk: "31", Nilai: map[string]string{"ObjectNo": "1", "ObjectName": "UJI GEDUNG"}},
			{ID: "42", Induk: "32", Nilai: map[string]string{"ObjectNo": "2", "ObjectName": "UJI GUDANG"}}},
		{"T_RISKLOCATION", "T_PROPERTY"}: {{ID: "51", Induk: "41", Nilai: map[string]string{"ASMAddress": "UJI JALAN"}}},
		{"T_PROPERTYITEMLIST", "T_PROPERTY"}: {{ID: "61", Induk: "41", Nilai: map[string]string{"ItemType": "UJI BANGUNAN",
			"Currency": "IDR", "TSIObjectItem": "1000000000", "PropertyItemNo": "1"}}},
		{"T_COVERAGELIST", "T_PROPERTYITEMLIST"}: {
			{ID: "71", Induk: "61", Src: cov, Nilai: map[string]string{"OLDID": "UJI-C1", "Coverage": "UJI-COV1"}},
			{ID: "72", Induk: "61", Src: cov, Nilai: map[string]string{"OLDID": "UJI-C2", "TSINusantaraRe": "250000000"}}},
		{"T_SPREADINGLIST", "T_COVERAGELIST"}: {
			{ID: "81", Induk: "72", Src: cov + "/SpreadingList", Nilai: map[string]string{"TreatyType": "10003", "SharePercentage": "60"}},
			{ID: "82", Induk: "72", Src: cov + "/SpreadingList", Nilai: map[string]string{"TreatyType": "10015", "SharePercentage": "40"}}},
	}
	doc, err := models.SusunPolisDatar(akar, bacaUji(baris))
	if err != nil {
		t.Fatal(err)
	}
	models.LengkapiPolisDatar(doc, "UJI-RNM-F.009", models.TambahanPolisDatar{BusinessCode: "UJI-B1", BusinessName: "UJI FIRE",
		BusinessOldId: "UJ", GroupPanel: "006", InsuredName: "UJI TERTANGGUNG", StatusBusiness: "UJI BARU",
		MarketingName: "UJI MARKETING", TeamGroup: "UJI-TIM"}, map[string]string{"IDR": "10026"})
	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	h := models.HalamanBaru()
	if err := models.PilihPolis(h, "UJI-RNM-F.009", "0", b); err != nil {
		t.Fatal(err)
	}
	if !models.IsFire(h) {
		t.Fatalf("lini: BusinessType %q", h.Ambil(models.OQ+"BusinessType"))
	}
	for jalur, mau := range map[string]string{
		models.JalurNoPolis: "UJI-RNM-F.009", models.JalurMulaiPolis: "2026-01-01", models.JalurAkhirPolis: "2026-12-31",
		models.AwalanPolis + ".PercentShare": "25", models.OQ + "NoOfferSlip": "UJI-OFR", models.OQ + "BusinessCode": "UJI-B1",
		models.OQ + "BusinessName": "UJI FIRE", models.OQ + "InsuredName": "UJI TERTANGGUNG", models.OQ + "SobName": "UJI SOB",
		models.OQ + "StatusBusiness": "UJI BARU", models.OQ + "MarketingName": "UJI MARKETING", models.OQ + "TeamGroup": "UJI-TIM",
		models.OQ + "MOID": "UJI-MO",
	} {
		if got := h.Ambil(jalur); got != mau {
			t.Errorf("%s = %q, mau %q", jalur, got, mau)
		}
	}
	lok := h.AmbilDaftar(models.DaftarLokasiPolis)
	if len(lok) != 2 || lok[0]["Property.ObjectName"] != "UJI GEDUNG" || lok[1]["Property.ObjectNo"] != "2" ||
		lok[0]["Property.RiskLocation.ASMAddress"] != "UJI JALAN" {
		t.Fatalf("LocationList: %+v", lok)
	}
	if c := h.AmbilDaftar(models.DaftarCedingCo); len(c) != 1 || c[0]["CedingCoName"] != "UJI CEDING" {
		t.Fatalf("CedingCoList: %+v", c)
	}
	item := h.AmbilDaftar(models.DaftarLokasiPolis + "(1).Property.PropertyItemList")
	if len(item) != 1 || item[0]["CurrencyID"] != "10026" || item[0]["TSIObjectItem"] != "1000000000" {
		t.Fatalf("PropertyItemList: %+v", item)
	}
	covs := h.AmbilDaftar(models.DaftarLokasiPolis + "(1).Property.PropertyItemList(1).CoverageList")
	if len(covs) != 2 || covs[0]["IndexCoverage"] != "1" || covs[1]["IndexCoverage"] != "2" || covs[1]["TSINusantaraRe"] != "250000000" {
		t.Fatalf("CoverageList: %+v", covs)
	}
	sp := h.AmbilDaftar(models.DaftarLokasiPolis + "(1).Property.PropertyItemList(1).CoverageList(2).SpreadingList")
	if len(sp) != 2 || sp[0]["TreatyType"] != "10003" || sp[1]["SharePercentage"] != "40" {
		t.Fatalf("SpreadingList: %+v", sp)
	}
	if n := len(h.AmbilDaftar(models.DaftarLokasiPolis + "(2).Property.PropertyItemList")); n != 0 {
		t.Fatalf("lokasi 2 tanpa item, dapat %d", n)
	}
}

func TestSusunPolisDatarPetaSisi(t *testing.T) {
	// setiap sisi menunjuk tabel yang ada di peta; induknya T_GENERAL_POLIS atau tabel peta lain.
	for _, s := range models.SisiPolisDatarSemua() {
		if _, ada := models.TabelPolisDatarDari(s.Tabel); !ada {
			t.Errorf("sisi %s <- %s: tabel tidak ada di peta", s.Tabel, s.Induk)
		}
		if _, ada := models.TabelPolisDatarDari(s.Induk); !ada {
			t.Errorf("sisi %s <- %s: induk tidak ada di peta", s.Tabel, s.Induk)
		}
		if len(s.Jalur) == 0 {
			t.Errorf("sisi %s <- %s tanpa jalur", s.Tabel, s.Induk)
		}
	}
}

func TestSusunPolisDatarAneka(t *testing.T) {
	// bentuk polis NB sistem baru di DEV 10-10-2026: LocationList > Property > RiskLocation > OccupationList >
	// AnekaList > CoverageList > SpreadingList (SRC_PATH), tanpa PropertyItemList.
	const okup = "LocationList/Property/RiskLocation/OccupationList"
	akar := models.BarisPolisDatar{ID: "UJI-NB-2"}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_LOCATIONLIST", "T_GENERAL_POLIS"}: {{ID: "31", Induk: "UJI-NB-2"}},
		{"T_PROPERTY", "T_LOCATIONLIST"}:      {{ID: "41", Induk: "31", Nilai: map[string]string{"ObjectNo": "1"}}},
		{"T_RISKLOCATION", "T_PROPERTY"}:      {{ID: "51", Induk: "41", Nilai: map[string]string{"ASMAddress": "UJI JALAN"}}},
		{"T_OCCUPATIONLIST", "T_RISKLOCATION"}: {{ID: "61", Induk: "51", Src: okup,
			Nilai: map[string]string{"OccupationId": "UJI-OC1", "OccupationName": "UJI OKUPASI"}}},
		{"T_ANEKALIST", "T_OCCUPATIONLIST"}: {{ID: "71", Induk: "61", Src: okup + "/AnekaList",
			Nilai: map[string]string{"ObjectName": "UJI MESIN", "TSI": "5000000", "Currency.Name": "IDR"}}},
		{"T_COVERAGELIST", "T_ANEKALIST"}: {{ID: "81", Induk: "71", Src: okup + "/AnekaList/CoverageList",
			Nilai: map[string]string{"OLDID": "UJI-C9", "TSINusantaraRe": "5000000"}}},
		{"T_SPREADINGLIST", "T_COVERAGELIST"}: {{ID: "91", Induk: "81", Src: okup + "/AnekaList/CoverageList/SpreadingList",
			Nilai: map[string]string{"TreatyType": "10003", "SharePercentage": "100"}}},
	}
	doc, err := models.SusunPolisDatar(akar, bacaUji(baris))
	if err != nil {
		t.Fatal(err)
	}
	if k := models.KodeMataUangDatar(doc); len(k) != 1 || k[0] != "IDR" {
		t.Fatalf("kode mata uang: %v", k)
	}
	models.LengkapiPolisDatar(doc, "UJI-RNM-F.010", models.TambahanPolisDatar{GroupPanel: "003", BusinessOldId: "UJ",
		StatusBusiness: "New Business"}, map[string]string{"IDR": "10026"})
	b, _ := json.Marshal(doc)
	h := models.HalamanBaru()
	if err := models.PilihPolis(h, "UJI-RNM-F.010", "0", b); err != nil {
		t.Fatal(err)
	}
	if h.Ambil(models.OQ+"StatusBusiness") != "1" { // Section Periode: StatusBusiness 1 = "New Business"
		t.Fatalf("StatusBusiness %q", h.Ambil(models.OQ+"StatusBusiness"))
	}
	if h.Ambil(models.OQ+"BusinessType") != "Aneka" || !models.IsAneka(h) {
		t.Fatalf("lini Aneka: %q", h.Ambil(models.OQ+"BusinessType"))
	}
	calon := models.CalonObjek(h)
	if len(calon) != 1 || calon[0].Medan["ObjectName"] != "UJI OKUPASI" || calon[0].Medan["ObjectLocation"] != "UJI JALAN" {
		t.Fatalf("calon objek Aneka: %+v", calon)
	}
	cov := h.AmbilDaftar(calon[0].Daftar["CoverageList"])
	if len(cov) != 1 || cov[0]["OLDID"] != "UJI-C9" || cov[0]["IndexAneka"] != "1" || cov[0]["IndexOccupation"] != "1" {
		t.Fatalf("CoverageList calon (%s): %+v", calon[0].Daftar["CoverageList"], cov)
	}
	if sp := h.AmbilDaftar(calon[0].Daftar["CoverageList"] + "(1).SpreadingList"); len(sp) != 1 || sp[0]["SharePercentage"] != "100" {
		t.Fatalf("SpreadingList Aneka: %+v", sp)
	}
	// SetIndexObject_DT: AnekaList ber-IdxOccupation -> GettsiAneka memilih aneka ber-IndexAneka (bukan aneka terakhir)
	an := h.AmbilDaftar(models.DaftarLokasiPolis + "(1).Property.RiskLocation.OccupationList(1).AnekaList")
	if len(an) != 1 || an[0]["IdxOccupation"] != "1" || an[0]["IdxAneka"] != "1" || an[0]["IdxLocation"] != "1" ||
		an[0]["Currency.Name"] != "IDR" || an[0]["Currency.ID"] != "10026" {
		t.Fatalf("AnekaList: %+v", an)
	}
}

func TestSusunPolisDatarTampilanNB(t *testing.T) {
	// Tampilan polis NB hanya-baca (work owner 10-10-2026 "IKUTI MODUL NB FACIN"): sub-tab Object Address (Building
	// Construction), Surrounding Risk, Occupation (TableOfLimit), FEA, Loss Record (CoinsData), dan Deductible dibaca
	// dari tabel datar ke jalur halaman yang dibaca layar.
	const cov = "LocationList/Property/PropertyItemList/CoverageList"
	akar := models.BarisPolisDatar{ID: "UJI-NB-2"}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_LOCATIONLIST", "T_GENERAL_POLIS"}: {{ID: "31", Induk: "UJI-NB-2", Nilai: map[string]string{"LossRatio1YearAmount": "12.5"}}},
		{"T_PROPERTY", "T_LOCATIONLIST"}:      {{ID: "41", Induk: "31", Nilai: map[string]string{"ObjectNo": "1"}}},
		{"T_BUILDINGCONSTRUCTION", "T_PROPERTY"}: {{ID: "42", Induk: "41", Nilai: map[string]string{"RoofType": "14",
			"NumberOfFloor": "2"}}},
		{"T_SURROUNDINGRISK", "T_PROPERTY"}: {{ID: "43", Induk: "41", Nilai: map[string]string{"FrontOccupation": "UJI TOKO",
			"FloodAreaStatus": "0"}}},
		{"T_OCCUPATIONLIST", "T_PROPERTY"}: {{ID: "44", Induk: "41", Src: "LocationList/Property/OccupationList",
			Nilai: map[string]string{"OccupationId": "UJI-OC1"}}},
		{"T_TABLEOFLIMIT", "T_OCCUPATIONLIST"}:   {{ID: "45", Induk: "44", Nilai: map[string]string{"Description": "UJI KELAS I"}}},
		{"T_LISTCAUSEOFLOSS", "T_PROPERTY"}:      {{ID: "46", Induk: "41", Nilai: map[string]string{"Amount": "100", "DateOfLoss": "20250101"}}},
		{"T_COINSDATA", "T_LISTCAUSEOFLOSS"}:     {{ID: "47", Induk: "46", Nilai: map[string]string{"CoinsName": "UJI KOAS"}}},
		{"T_FEALIST", "T_LOCATIONLIST"}:          {{ID: "48", Induk: "31", Nilai: map[string]string{"APAR": "3", "PrivateTruckBrigade": "1"}}},
		{"T_PROPERTYITEMLIST", "T_PROPERTY"}:     {{ID: "61", Induk: "41", Nilai: map[string]string{"Currency": "IDR"}}},
		{"T_COVERAGELIST", "T_PROPERTYITEMLIST"}: {{ID: "71", Induk: "61", Src: cov, Nilai: map[string]string{"OLDID": "UJI-C1"}}},
		{"T_DEDUCTIBLELIST", "T_COVERAGELIST"}: {{ID: "91", Induk: "71", Src: cov + "/DeductibleList",
			Nilai: map[string]string{"TypeDeductible": "7", "Amount": "5000"}}},
	}
	doc, err := models.SusunPolisDatar(akar, bacaUji(baris))
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	h := models.HalamanBaru()
	if err := models.TerapkanPolis(h, b); err != nil {
		t.Fatal(err)
	}
	l := models.DaftarLokasiPolis
	lok := h.AmbilDaftar(l)
	if len(lok) != 1 {
		t.Fatalf("LocationList: %+v", lok)
	}
	for k, mau := range map[string]string{
		"Property.BuildingConstruction.RoofType":      "14",
		"Property.BuildingConstruction.NumberOfFloor": "2",
		"Property.SurroundingRisk.FrontOccupation":    "UJI TOKO",
		"Property.SurroundingRisk.FloodAreaStatus":    "0",
		"LossRatio1YearAmount":                        "12.5",
	} {
		if lok[0][k] != mau {
			t.Errorf("LocationList(1).%s = %q, mau %q", k, lok[0][k], mau)
		}
	}
	periksa := func(daftar, k, mau string) {
		t.Helper()
		rows := h.AmbilDaftar(daftar)
		if len(rows) != 1 || rows[0][k] != mau {
			t.Errorf("%s: %+v, mau %s = %q", daftar, rows, k, mau)
		}
	}
	periksa(l+"(1).Property.OccupationList", "TableOfLimit.Description", "UJI KELAS I")
	periksa(l+"(1).Property.ListCauseOfLoss", "CoinsData.CoinsName", "UJI KOAS")
	periksa(l+"(1).FEAList", "PrivateTruckBrigade", "1")
	periksa(l+"(1).Property.PropertyItemList(1).CoverageList(1).DeductibleList", "Amount", "5000")
}
