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

func TestSusunPolisDatarAnekaTampilanNB(t *testing.T) {
	// Tab Object Aneka tampilan polis NB (work owner 11-10-2026 "tampilan view polis yang belum tergenerate tolong
	// diperbaiki. cek tampilan nya di nb facin!"): medan form objek NB (migrasi nbfacin 214) - kolom tambahan T_ANEKALIST,
	// halaman .MarineHull / .VehicleHE / .AviationHull / .MemorandumAccept - sampai ke baris AnekaList halaman klaim.
	const okup = "LocationList/Property/RiskLocation/OccupationList"
	akar := models.BarisPolisDatar{ID: "UJI-NB-3"}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_LOCATIONLIST", "T_GENERAL_POLIS"}: {{ID: "31", Induk: "UJI-NB-3"}},
		{"T_PROPERTY", "T_LOCATIONLIST"}:      {{ID: "41", Induk: "31", Nilai: map[string]string{"ObjectNo": "1"}}},
		{"T_RISKLOCATION", "T_PROPERTY"}:      {{ID: "51", Induk: "41", Nilai: map[string]string{"ASMAddress": "UJI JALAN"}}},
		{"T_OCCUPATIONLIST", "T_RISKLOCATION"}: {{ID: "61", Induk: "51", Src: okup,
			Nilai: map[string]string{"OccupationId": "UJI-OC1"}}},
		{"T_ANEKALIST", "T_OCCUPATIONLIST"}: {{ID: "71", Induk: "61", Src: okup + "/AnekaList",
			Nilai: map[string]string{"ObjectName": "UJI KAPAL", "NO_KONTRAK": "UJI-K1", "NO_KTP": "UJI-KTP", "Others": "UJI CATATAN",
				"LoanAccountNumber": "UJI-L1", "LimitOfLiability": "100", "NoOfTree": "7"}}},
		{"T_MARINEHULL", "T_ANEKALIST"}:       {{ID: "72", Induk: "71", Nilai: map[string]string{"NRTRT": "10", "LENGTHM": "50"}}},
		{"T_VEHICLEHE", "T_ANEKALIST"}:        {{ID: "73", Induk: "71", Nilai: map[string]string{"ChassisNumber": "UJI-CH1"}}},
		{"T_AVIATIONHULL", "T_ANEKALIST"}:     {{ID: "74", Induk: "71", Nilai: map[string]string{"Pilots": "2"}}},
		{"T_MEMORANDUMACCEPT", "T_ANEKALIST"}: {{ID: "75", Induk: "71", Nilai: map[string]string{"Collateral": "UJI JAMINAN"}}},
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
	an := h.AmbilDaftar(models.DaftarLokasiPolis + "(1).Property.RiskLocation.OccupationList(1).AnekaList")
	if len(an) != 1 {
		t.Fatalf("AnekaList: %+v", an)
	}
	for k, mau := range map[string]string{
		"NO_KONTRAK": "UJI-K1", "NO_KTP": "UJI-KTP", "Others": "UJI CATATAN", "LoanAccountNumber": "UJI-L1",
		"LimitOfLiability": "100", "NoOfTree": "7", "MarineHull.NRTRT": "10", "MarineHull.LENGTHM": "50",
		"VehicleHE.ChassisNumber": "UJI-CH1", "AviationHull.Pilots": "2", "MemorandumAccept.Collateral": "UJI JAMINAN",
	} {
		if an[0][k] != mau {
			t.Errorf("AnekaList %s = %q, mau %q", k, an[0][k], mau)
		}
	}
	for _, tb := range []string{"T_MARINEHULL", "T_AVIATIONHULL", "T_MEMORANDUMACCEPT"} {
		if !models.TabelTampilanNB(tb) {
			t.Errorf("%s harus tabel tampilan NB (dimuat hanya bila berindeks induk)", tb)
		}
	}
}

func TestSusunPolisDatarCargoTampilanNB(t *testing.T) {
	// Tab Object / Coverage Marine Cargo tampilan polis NB (work owner 11-10-2026 "cek tampilan nya di nb facin!"): Ship,
	// LC, Survey Agent, Additional Ship (anak Ship), Warranty Extension (AdditionalCoverage) sampai ke baris CargoList.
	akar := models.BarisPolisDatar{ID: "UJI-NB-4"}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_CARGOLIST", "T_GENERAL_POLIS"}: {{ID: "31", Induk: "UJI-NB-4",
			Nilai: map[string]string{"TradingNote": "UJI DAGANG", "InvoiceDate": "20260105", "MaxAgeDetail": "UJI MAKS"}}},
		{"T_SHIP", "T_CARGOLIST"}:        {{ID: "41", Induk: "31", Nilai: map[string]string{"NM_SHIP": "UJI KAPAL", "ID": "UJI-S1"}}},
		{"T_ADDITIONALSHIP", "T_SHIP"}:   {{ID: "42", Induk: "41", Nilai: map[string]string{"NM_SHIP": "UJI KAPAL 2", "VOYAGE_NO": "UJI-V1"}}},
		{"T_LC", "T_CARGOLIST"}:          {{ID: "43", Induk: "31", Nilai: map[string]string{"LC": "true", "LCBANK": "UJI BANK"}}},
		{"T_SURVEYAGENT", "T_CARGOLIST"}: {{ID: "44", Induk: "31", Nilai: map[string]string{"pyFullName": "UJI SURVEYOR"}}},
		{"T_COVERAGELIST", "T_CARGOLIST"}: {{ID: "51", Induk: "31", Src: "CargoList/CoverageList",
			Nilai: map[string]string{"CoverageNote": "UJI ICC A"}}},
		{"T_ADDITIONALCOVERAGE", "T_COVERAGELIST"}: {{ID: "52", Induk: "51", Nilai: map[string]string{"Coverage": "UJI-W1"}}},
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
	cg := h.AmbilDaftar(models.AwalanPolis + ".CargoList")
	if len(cg) != 1 {
		t.Fatalf("CargoList: %+v", cg)
	}
	for k, mau := range map[string]string{"TradingNote": "UJI DAGANG", "InvoiceDate": "20260105", "MaxAgeDetail": "UJI MAKS",
		"Ship.NM_SHIP": "UJI KAPAL", "Ship.ID": "UJI-S1", "LC.LCBANK": "UJI BANK", "SurveyAgent.pyFullName": "UJI SURVEYOR"} {
		if cg[0][k] != mau {
			t.Errorf("CargoList %s = %q, mau %q", k, cg[0][k], mau)
		}
	}
	if as := h.AmbilDaftar(models.AwalanPolis + ".CargoList(1).Ship.AdditionalShip"); len(as) != 1 || as[0]["VOYAGE_NO"] != "UJI-V1" {
		t.Errorf("AdditionalShip: %+v", as)
	}
	if w := h.AmbilDaftar(models.AwalanPolis + ".CargoList(1).CoverageList(1).AdditionalCoverage"); len(w) != 1 || w[0]["Coverage"] != "UJI-W1" {
		t.Errorf("AdditionalCoverage: %+v", w)
	}
	for _, tb := range []string{"T_SHIP", "T_LC", "T_SURVEYAGENT", "T_ADDITIONALSHIP", "T_ADDITIONALCOVERAGE"} {
		if !models.TabelTampilanNB(tb) {
			t.Errorf("%s harus tabel tampilan NB (dimuat hanya bila berindeks induk)", tb)
		}
	}
}

func TestSusunPolisDatarPaymentTampilanNB(t *testing.T) {
	// Tab Premium Deduction / Payment tampilan polis NB (keputusan work owner 11-10-2026: nilai TERSIMPAN): PAY_* akar ->
	// PolicyData.Payment, CurrencyList PAY_* -> Policy.Payment (jalur Pega `PaymentCurrencyList`), ListInstallment.
	akar := models.BarisPolisDatar{ID: "UJI-NB-5", Nilai: map[string]string{"PolicyData.Payment.RICommision": "12.5",
		"PolicyData.Payment.Installment": "2"}}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_CURRENCYLIST", "T_GENERAL_POLIS"}: {{ID: "31", Induk: "UJI-NB-5",
			Nilai: map[string]string{"Name": "IDR", "Policy.Payment.Commision": "125", "SumTotalPayment": "800"}}},
		{"T_LISTINSTALLMENT", "T_CURRENCYLIST"}: {{ID: "41", Induk: "31",
			Nilai: map[string]string{"InstallmentNo": "1", "PaymentTotal": "400"}}},
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
	if v := h.Ambil(models.AwalanPolis + ".PolicyData.Payment.RICommision"); v != "12.5" {
		t.Errorf("RICommision %q", v)
	}
	cu := h.AmbilDaftar(models.AwalanPolis + ".CurrencyList")
	if len(cu) != 1 || cu[0]["Policy.Payment.Commision"] != "125" || cu[0]["SumTotalPayment"] != "800" {
		t.Fatalf("CurrencyList: %+v", cu)
	}
	if ci := h.AmbilDaftar(models.AwalanPolis + ".CurrencyList(1).ListInstallment"); len(ci) != 1 || ci[0]["PaymentTotal"] != "400" {
		t.Errorf("ListInstallment: %+v", ci)
	}
	for _, tb := range []string{"T_CURRENCYLIST", "T_LISTINSTALLMENT"} {
		if !models.TabelTampilanNB(tb) {
			t.Errorf("%s harus tabel tampilan NB", tb)
		}
	}
	akarSpec, _ := models.TabelPolisDatarDari(models.TabelAkarPolisDatar)
	ada := map[string]string{}
	for _, k := range akarSpec.Kolom {
		ada[k.Nama] = k.Medan
	}
	for kol, mau := range map[string]string{"PAY_RI_COMMISION": "PolicyData.Payment.RICommision",
		"PAY_PCT_BROKERAGE_FEE": "PolicyData.Payment.PctBrokerageFee", "PAY_INSTALLMENT": "PolicyData.Payment.Installment"} {
		if ada[kol] != mau {
			t.Errorf("peta T_GENERAL_POLIS %s -> %q, mau %q", kol, ada[kol], mau)
		}
	}
}

func TestSusunPolisDatarClausesTampilanNB(t *testing.T) {
	// Tab Clauses tampilan polis NB (FIRE): T_CLAUSELIST (nbfacin 197, satu baris per argumen) -> OfferFacIn.ClauseList.
	akar := models.BarisPolisDatar{ID: "UJI-NB-6"}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_CLAUSELIST", "T_GENERAL_POLIS"}: {
			{ID: "31", Induk: "UJI-NB-6", Nilai: map[string]string{"SeqNo": "1", "ClauseCode": "UJI-CL1", "ClauseTitle": "UJI JUDUL", "ArgumentCount": "2"}},
			{ID: "32", Induk: "UJI-NB-6", Nilai: map[string]string{"SeqNo": "1", "ClauseCode": "UJI-CL1", "ClauseTitle": "UJI JUDUL", "ArgumentCount": "2"}},
		},
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
	if cl := h.AmbilDaftar(models.AwalanPolis + ".ClauseList"); len(cl) != 2 || cl[0]["ClauseCode"] != "UJI-CL1" || cl[1]["SeqNo"] != "1" {
		t.Fatalf("ClauseList: %+v", cl)
	}
	if !models.TabelTampilanNB("T_CLAUSELIST") {
		t.Error("T_CLAUSELIST harus tabel tampilan NB")
	}
}

func TestSusunPolisDatarScoringTampilanNB(t *testing.T) {
	// Tab Scoring Risk tampilan polis NB (FIRE, nilai tersimpan - keputusan work owner 11-10-2026): lima tabel rancangan
	// nbfacin 205 -> OfferFacIn.ScoringRisk; kolom penentu faktor (FACTOR_GROUP / FACTOR_NAME, aturan loader V-30) dan opsi
	// (CHECK_BOX / OPTION_NO) ikut dimuat.
	akar := models.BarisPolisDatar{ID: "UJI-NB-7"}
	baris := map[[2]string][]models.BarisPolisDatar{
		{"T_SCORINGRISK", "T_GENERAL_POLIS"}: {
			{ID: "41", Induk: "UJI-NB-7", Nilai: map[string]string{"FinalScore": "80.5", "NoteFinalScore": "Above Average"}},
		},
		{"T_SCORINGRESULT", "T_SCORINGRISK"}: {
			{ID: "42", Induk: "41", Nilai: map[string]string{"AboveAverage": "true", "AcceptableNotApproval": "true"}},
		},
		{"T_DATASCORINGRISKLIST", "T_SCORINGRISK"}: {
			{ID: "43", Induk: "41", Nilai: map[string]string{"OccupationNote": "UJI OKUPASI", "Score": "80.5"}},
			{ID: "44", Induk: "41", Nilai: map[string]string{"IndexLocation": "2"}},
		},
		{"T_SCORING_FACTOR", "T_DATASCORINGRISKLIST"}: {
			{ID: "45", Induk: "43", Nilai: map[string]string{"FactorGroup": "LossRatio", "FactorName": "LossRatio",
				"PercentLossRatio": "12.5", "ScorePerFactor": "15"}},
		},
		{"T_SCORING_OPTION", "T_SCORING_FACTOR"}: {
			{ID: "46", Induk: "45", Nilai: map[string]string{"CheckBox": "true", "OptionNo": "5", "Score": "15"}},
		},
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
	sr := models.AwalanPolis + ".ScoringRisk"
	if h.Ambil(sr+".NoteFinalScore") != "Above Average" || h.Ambil(sr+".ScoringResult.AcceptableNotApproval") != "true" {
		t.Errorf("ScoringRisk: %q / %q", h.Ambil(sr+".NoteFinalScore"), h.Ambil(sr+".ScoringResult.AcceptableNotApproval"))
	}
	if d := h.AmbilDaftar(sr + ".DataScoringRiskList"); len(d) != 2 || d[1]["IndexLocation"] != "2" {
		t.Fatalf("DataScoringRiskList: %+v", d)
	}
	f := h.AmbilDaftar(sr + ".DataScoringRiskList(1).ScoringFactor")
	if len(f) != 1 || f[0]["FactorGroup"] != "LossRatio" || f[0]["FactorName"] != "LossRatio" || f[0]["PercentLossRatio"] != "12.5" {
		t.Fatalf("ScoringFactor: %+v", f)
	}
	o := h.AmbilDaftar(sr + ".DataScoringRiskList(1).ScoringFactor(1).ScoringOption")
	if len(o) != 1 || o[0]["CheckBox"] != "true" || o[0]["OptionNo"] != "5" {
		t.Fatalf("ScoringOption: %+v", o)
	}
	for _, tb := range []string{"T_SCORINGRISK", "T_SCORINGRESULT", "T_DATASCORINGRISKLIST", "T_SCORING_FACTOR", "T_SCORING_OPTION"} {
		if !models.TabelTampilanNB(tb) {
			t.Errorf("%s harus tabel tampilan NB", tb)
		}
	}
}
