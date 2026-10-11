package services

// Gerbang Class of Business (tiket 73 P0) - lini bisnis case dari BUSINESS (GROUPPANEL + OLDID) lewat BusinessType_DeT, lalu
// predikat When korpus lewat registry `rules`. Bukti `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\`:
//   - Activity\SetBusinessType_Act.xml langkah 2: Property-Map-DecisionTable `.Quotation.BusinessType` = BusinessType_DeT
//     (langkah 1 yang men-SET GroupPanel / BusinessCode di-REMARK);
//   - DecisionTable\BusinessType_DeT.xml (ASM-FW-GISFW-WORK!BUSINESSTYPE_DET): kolom `.Quotation.GroupPanel` dan
//     `.Quotation.BusinessOldId` (operator `=`), 36 baris, pyEvaluateAllRows "no" (first match, L238), pyDefaultResult "UNKNOWN".
//     Diurai sesi c3 (skrip c3_bisnisdet.py): sel utama per REPEATINGINDEX + pyOrConditions per pyRowNum (0-based: pyRowNum 7 =
//     baris 8 "Liability", L489-497); sel kosong = cocok apa pun.
//   - Isi Quotation (Activity\CopyToPolicyListPASSG.xml L12143-12229): BusinessCode 5 digit (10028), BusinessName nama Class of
//     Business berhuruf besar, GroupPanel "002" -> BusinessCode = BUSINESS.ID, BusinessName = BUSINESS.NOTE `[dugaan kuat]`.
// Sumber sistem baru: baris BUSINESS Class of Business case (NOTE + BUSINESSGROUPID, cara kode bisnis tiket 40 / OLDID tiket 62)
// yang GROUPPANEL-nya terisi - saringan `GROUPPANEL IS NOT NULL` Pega RDBList\GetOldIDBusiness_SQL.xml (keputusan work owner
// 09-10-2026); ID, OLDID, GROUPPANEL satu kueri (repository.sqlBisnisLini).
// 0 / > 1 baris ber-GROUPPANEL -> nilai kosong -> BusinessType "UNKNOWN" (default DT), predikat menilai apa adanya.

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

// barisBusinessType - satu baris BusinessType_DeT; oldID kosong = cocok apa pun.
type barisBusinessType struct {
	groupPanel string
	oldID      []string
	hasil      string
}

// businessTypeBawaan - pyDefaultResult BusinessType_DeT.
const businessTypeBawaan = "UNKNOWN"

// tabelBusinessType - BusinessType_DeT baris 1-36 urut (first match).
var tabelBusinessType = []barisBusinessType{
	{"001", nil, "Medicare"},
	{"002", []string{"03"}, "PA"},
	{"003", []string{"36", "37", "38", "39", "40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "50", "51", "52", "C2", "C3",
		"C4", "C5", "C6", "C7", "C8", "C9", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "F1", "F2", "F3", "F4", "F5"}, "Bonding"},
	{"003", []string{"95", "96", "97", "98", "99", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "B4", "B5", "B6", "B8", "B9",
		"C1"}, "BondingKBG"},
	{"003", []string{"20"}, "HE"},
	{"003", []string{"05", "55"}, "MarineHull"},
	{"003", []string{"21"}, "Glass"},
	{"003", []string{"15", "24", "89", "90", "91", "E7", "E9"}, "Liability"},
	{"003", []string{"16"}, "AllRisk"},
	{"003", []string{"06"}, "AviationHull"},
	{"003", []string{"14"}, "Burglary"},
	{"003", []string{"07"}, "Car"},
	{"003", []string{"08"}, "Ear"},
	{"003", []string{"09"}, "ElectronicEquipment"},
	{"003", []string{"19"}, "Fidelity"},
	{"003", []string{"27"}, "GolfInsurance"},
	{"003", []string{"17"}, "CIT"},
	{"003", []string{"18"}, "CIS"},
	{"003", []string{"11"}, "Boiler"},
	{"003", []string{"93"}, "Workmen"},
	{"003", []string{"10"}, "MBD"},
	{"003", []string{"12"}, "ContractorsPM"},
	{"003", []string{"79", "80", "81", "82", "83", "84", "85", "86", "87", "88"}, "CustomBond"},
	{"003", []string{"28"}, "FireStyle1"},
	{"003", []string{"B3"}, "LandRig"},
	{"003", nil, "Aneka"},
	{"004", nil, "MarineCargo"},
	{"005", nil, "Travel"},
	{"006", []string{"78"}, "OilGas"},
	{"006", []string{"22", "26", "29", "57"}, "FireStyle1"},
	{"006", []string{"01", "34", "59"}, "FireStyle2"},
	{"006", nil, "Fire"},
	{"007", []string{"02"}, "MBUCar"},
	{"007", []string{"60"}, "MBUMotorCycle"},
	{"009", []string{"L1", "L2", "L3", "L4", "L5", "L6", "L7", "L8", "L9", "L10", "L11", "L12", "L13", "L14", "L15", "L16", "L17", "L18",
		"L19", "L20", "L21"}, "Life"},
	{"008", nil, "Medicare"},
}

// businessTypeDari - BusinessType_DeT first match atas (GroupPanel, BusinessOldId); tanpa baris cocok -> "UNKNOWN".
func businessTypeDari(groupPanel, oldID string) string {
	for _, b := range tabelBusinessType {
		if b.groupPanel != groupPanel {
			continue
		}
		if len(b.oldID) == 0 || dalam(b.oldID, oldID) {
			return b.hasil
		}
	}
	return businessTypeBawaan
}

// LiniKasus - `businessType` + `lini` GET kasus (tiket 73 P0) dan `tampilanAneka` (tiket 73 P1): hasil predikat When korpus.
type LiniKasus struct {
	BusinessType                                           string
	Fire, Aneka, AnekaNotBonding, Liability, Bonding, Golf bool
	// MarineCargo - tiket 74 P0: When IsMarineCargo (`NB FacIn\When\IsMarineCargo.xml` L44 `G OR H`, BusinessType "MarineCargo"
	// = GroupPanel 004).
	MarineCargo   bool
	TampilanAneka TampilanAneka
	// JenisCsv - tiket 73 P6: template CSV objek case (models.JenisCsvFire / JenisCsvCargo / TemplateCsv*; "" = tanpa template),
	// lihat templateCsvKasus.
	JenisCsv string
}

// templateCsvBerurutan - SetDownloadMasterCSV langkah 2-6 (urutan dipertahankan: yang terakhir benar menang).
var templateCsvBerurutan = []struct{ predikat, jenis string }{
	{"IsAneka", models.TemplateCsvAneka}, {"IsMarineHull", models.TemplateCsvMarineHull},
	{"IsContractorsPlantMachinery", models.TemplateCsvAnekaCPM}, {"IsFire", models.JenisCsvFire},
	{"IsMarineCargo", models.JenisCsvCargo},
}

// templateCsvKasus - lihat LiniKasus.JenisCsv.
func templateCsvKasus(id string, kb kasusBisnis) string {
	jenis := ""
	for _, p := range templateCsvBerurutan {
		if salahSatu(id, kb, p.predikat) {
			jenis = p.jenis
		}
	}
	return jenis
}

// TampilanAneka - sakelar form objek Aneka (tiket 73 P1, S1) - semua dari rules.Eval:
//   - Kontainer: "C1".."C9" yang visible-when-nya benar, `[terverifikasi]` NB FacIn\Section\InputDtlObjectAneka_FacIn.xml
//     L413 / L9678 / L17776 / L19883 / L23090 / L24543 / L27780 / L32104 / L34670;
//   - KolomTahunJumlah IsObjectWithQuantityYear; HE IsHE; Golf isGolfInsurance; TanamanHasil isGrowingTrees || IsYieldShortfall;
//     Engineering IsEngineering (MaintenancePeriod); Bonding IsBonding (Draft Wording); Briguna IsCreditBriguna (L4166 / L4465);
//   - TopRisk: sel `.Property.IsTopRisk` NB FacIn\Section\InputDtlObjectLocation_FacIn.xml (cell 17, pyCondition L6301) `IsMBD ||
//     IsHE || IsContractorsPlantMachinery||IsElectronicEquipment||IsCMI || IsMarineHull||IsLiability`;
//   - SectionLol IsObjectSectionAneka (K73-4: perbandingan BusinessCode 2 karakter ditiru apa adanya);
//   - JumlahSatuan `pyWorkPage.Quotation.BusinessCode!=20` (L23869, OTHER).
//
// Gerbang: grid tab Object Aneka NB FacIn\Section\InputDtlObject_FacIn.xml L36463 `IsAnekaNotBonding || IsHE || IsGolfInsurance ||
// IsGrowingTrees` ATAU L22724 `IsBonding || IsCustomBonds`; selain itu kontainer [] dan semua false.
type TampilanAneka struct {
	Kontainer                                                                        []string
	KolomTahunJumlah, HE, Golf, TanamanHasil, Engineering, Bonding, Briguna, TopRisk bool
	SectionLol, JumlahSatuan                                                         bool
	// Tiket 73 P2 - form coverage `NB FacIn\Section\InputCoverageAneka_FacIn.xml`: CurrencyCoverage = visible-when sel Currency
	// (L4052, verbatim di currencyCoverageAneka); MBD = IsMBD (sel Indemnity / Indemnity Unit / % Indemnity L5826 / L6132 / L6322).
	CurrencyCoverage, MBD bool
}

// currencyCoverageAneka - pyCondition sel Currency form coverage Aneka (L4052): `IsPortRisk || IsBuilderRisk||IsObjectWithQuantityYear
// || isAllRisk || IsMarineHull || IsMarineHullOffshore || IsMBD || IsHE || IsAviationHull || isBurglary`.
var currencyCoverageAneka = []string{"IsPortRisk", "IsBuilderRisk", "IsObjectWithQuantityYear", "isAllRisk", "IsMarineHull",
	"IsMarineHullOffshore", "IsMBD", "IsHE", "IsAviationHull", "isBurglary"}

// kontainerAneka - C1..C9: daftar When yang di-OR-kan (nama persis korpus).
var kontainerAneka = []struct {
	nama     string
	predikat []string
}{
	{"C1", []string{"IsAsuransiKredit", "IsEnvironmental", "IsCrime", "isGolfInsurance", "isCitCis", "IsLiability", "IsWhetherIndex",
		"isFidelity", "isEAR", "isGlass", "IsBonding", "isCAR", "IsCivilEngineeringCompletedRisks", "isBillboardNeonSyariah",
		"IsCustomBonds", "IsBoiler", "IsComprehensiveMachineriesInsurance"}},
	{"C2", []string{"IsPortRisk", "IsBuilderRisk", "IsMarineHull", "IsMarineHullOffshore"}},
	{"C3", []string{"isAllrisk", "isBurglary"}},
	{"C4", []string{"IsObjectWithQuantityYear"}},
	{"C5", []string{"IsObjectWithQuantityYear", "IsHE"}},
	{"C6", []string{"IsHE"}},
	{"C7", []string{"isAviationHull"}},
	{"C8", []string{"isGrowingTrees", "IsYieldShortfall"}},
	{"C9", []string{"isElectronicEquipment"}},
}

// gerbangAneka - grid tab Object Aneka (dua kontainer InputDtlObject_FacIn).
var gerbangAneka = []string{"IsAnekaNotBonding", "IsHE", "IsGolfInsurance", "IsGrowingTrees", "IsBonding", "IsCustomBonds"}

// topRiskAneka - kondisi sel IsTopRisk lokasi Aneka.
var topRiskAneka = []string{"IsMBD", "IsHE", "IsContractorsPlantMachinery", "IsElectronicEquipment", "IsCMI", "IsMarineHull", "IsLiability"}

// predikatLini - nama When registry per medan LiniKasus.
var predikatLini = []struct {
	nama string
	ke   func(*LiniKasus) *bool
}{
	{"IsFire", func(l *LiniKasus) *bool { return &l.Fire }},
	{"IsAneka", func(l *LiniKasus) *bool { return &l.Aneka }},
	{"IsAnekaNotBonding", func(l *LiniKasus) *bool { return &l.AnekaNotBonding }},
	{"IsLiability", func(l *LiniKasus) *bool { return &l.Liability }},
	{"IsBondingAndCustomBonds", func(l *LiniKasus) *bool { return &l.Bonding }},
	{"isGolfInsurance", func(l *LiniKasus) *bool { return &l.Golf }},
	{"IsMarineCargo", func(l *LiniKasus) *bool { return &l.MarineCargo }},
}

// kasusBisnis - rules.Kasus jalur Quotation yang dibaca predikat lini (dicacah sesi c3 dari registry_gen.go, rekursif):
// BusinessType dalam empat bentuk jalur (pyWorkPage.Quotation, pyWorkPage.OfferFacIn.QuotationData, pyWorkPage.Policy.Quotation,
// `.Quotation` relatif) - satu nilai; BusinessCode = BUSINESS.ID; BusinessName = BUSINESS.NOTE (Class of Business).
type kasusBisnis struct{ businessType, businessCode, businessName string }

func (k kasusBisnis) Nilai(jalur string) (string, bool) {
	switch jalur {
	case "pyWorkPage.Quotation.BusinessType", "pyWorkPage.OfferFacIn.QuotationData.BusinessType",
		"pyWorkPage.Policy.Quotation.BusinessType", ".Quotation.BusinessType",
		".OfferFacIn.QuotationData.BusinessType": // IsMarineCargo G - When kelas WORK, jalur relatif = pyWorkPage
		return k.businessType, true
	case "pyWorkPage.Quotation.BusinessCode", "pyWorkPage.OfferFacIn.QuotationData.BusinessCode":
		return k.businessCode, true
	case "pyWorkPage.Quotation.BusinessName", "pyWorkPage.OfferFacIn.QuotationData.BusinessName",
		".OfferFacIn.QuotationData.BusinessName":
		return k.businessName, true
	}
	return "", false
}

// LiniKasus - lini bisnis case `id`. LUNAK untuk GET kasus: galat baca / predikat ditolak (ErrTafsirBerbeda) -> medan itu false,
// dicatat log. Tanpa pembaca kasus / BUSINESS -> nilai nol.
func (s *Service) LiniKasus(ctx context.Context, id string) LiniKasus {
	var l LiniKasus
	kb, ada := s.kasusBisnisDari(ctx, id)
	if !ada {
		return l
	}
	l.BusinessType = kb.businessType
	for _, p := range predikatLini {
		v, err := rules.Eval(p.nama, kb)
		if err != nil {
			log.Printf("nbfacin: lini kasus %s: %s: %v", id, p.nama, err)
			continue
		}
		*p.ke(&l) = v
	}
	l.TampilanAneka = tampilanAneka(id, kb)
	l.JenisCsv = templateCsvKasus(id, kb)
	return l
}

// kasusBisnisDari - BusinessType / BusinessCode / BusinessName case `id` untuk predikat When (lini, tampilan dan proteksi
// Aneka). LUNAK: tanpa pembaca / case tidak terbaca -> ada=false; BUSINESS bukan tepat satu baris ber-GROUPPANEL -> nilai kosong
// (BusinessType "UNKNOWN"), dicatat log.
func (s *Service) kasusBisnisDari(ctx context.Context, id string) (kasusBisnis, bool) {
	if s.kasus == nil || s.tableOfLimit == nil || !idKasusSah(id) {
		return kasusBisnis{}, false
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if err != nil {
		if !errors.Is(err, repository.ErrKasusTidakAda) {
			log.Printf("nbfacin: lini kasus %s: %v", id, err)
		}
		return kasusBisnis{}, false
	}
	nama, grup := strings.TrimSpace(k.Opportunity.ClassOfBusiness), strings.TrimSpace(k.Opportunity.GroupBusinessID)
	// Baris BUSINESS ber-GROUPPANEL (saringan Pega, keputusan work owner 09-10-2026); tepat satu baris yang dipakai.
	var b models.BarisBisnisLini
	if nama != "" && grup != "" {
		v, err := s.tableOfLimit.BisnisLini(ctx, nama, grup)
		if err != nil {
			log.Printf("nbfacin: lini kasus %s: BUSINESS: %v", id, err)
		}
		if len(v) == 1 {
			b = v[0]
		} else {
			log.Printf("nbfacin: lini kasus %s: BUSINESS %q group %q = %d baris ber-GROUPPANEL - BusinessType UNKNOWN", id, nama, grup, len(v))
		}
	}
	return kasusBisnis{businessType: businessTypeDari(strings.TrimSpace(b.GroupPanel), strings.TrimSpace(b.OldID)),
		businessCode: strings.TrimSpace(b.ID), businessName: nama}, true
}

// salahSatu - OR predikat lewat rules.Eval; galat (ErrTafsirBerbeda) = predikat itu salah, dicatat.
func salahSatu(id string, kb kasusBisnis, nama ...string) bool {
	hasil := false
	for _, n := range nama {
		v, err := rules.Eval(n, kb)
		if err != nil {
			log.Printf("nbfacin: tampilan Aneka kasus %s: %s: %v", id, n, err)
			continue
		}
		hasil = hasil || v
	}
	return hasil
}

// tampilanAneka - lihat TampilanAneka.
func tampilanAneka(id string, kb kasusBisnis) TampilanAneka {
	t := TampilanAneka{Kontainer: []string{}}
	if !salahSatu(id, kb, gerbangAneka...) {
		return t
	}
	for _, k := range kontainerAneka {
		if salahSatu(id, kb, k.predikat...) {
			t.Kontainer = append(t.Kontainer, k.nama)
		}
	}
	t.KolomTahunJumlah = salahSatu(id, kb, "IsObjectWithQuantityYear")
	t.HE, t.Golf = salahSatu(id, kb, "IsHE"), salahSatu(id, kb, "isGolfInsurance")
	t.TanamanHasil = salahSatu(id, kb, "isGrowingTrees", "IsYieldShortfall")
	t.Engineering, t.Bonding = salahSatu(id, kb, "IsEngineering"), salahSatu(id, kb, "IsBonding")
	t.Briguna, t.TopRisk = salahSatu(id, kb, "IsCreditBriguna"), salahSatu(id, kb, topRiskAneka...)
	t.SectionLol = salahSatu(id, kb, "IsObjectSectionAneka")
	t.JumlahSatuan = !samaDuaPuluh(kb.businessCode)
	t.CurrencyCoverage, t.MBD = salahSatu(id, kb, currencyCoverageAneka...), salahSatu(id, kb, "IsMBD")
	return t
}

// samaDuaPuluh - `BusinessCode = 20` (pembanding angka Pega: teks "20" atau nilai angka 20).
func samaDuaPuluh(kode string) bool {
	k := strings.TrimSpace(kode)
	if k == "20" {
		return true
	}
	if n, err := strconv.ParseFloat(k, 64); err == nil {
		return n == 20
	}
	return false
}
