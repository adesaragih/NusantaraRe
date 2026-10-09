package models

// Untuk apa berkas ini: KLASIFIKASI LINI PRODUK - port rule When `Claim Fac In/When/*.xml` dalam SATU tempat (spec Bab 4
// keputusan work owner 19-09-2026 "satu fungsi klasifikasi"; AC 13). Setiap When yang dipakai section / activity Claim
// Fac In ditulis sebagai fungsi atas halaman; nama fungsi = nama rule.
//
// `[terverifikasi]` Semua rule lini menguji `pyWorkPage.Quotation.BusinessType` / `.BusinessCode` (atau langsung
// `OfferFacIn.QuotationData.*`). `pyWorkPage.Quotation` = SALINAN HALAMAN UTUH `pyWorkPage.OfferFacIn.QuotationData`
// (`CopyNB_Act` langkah 8), sehingga kedua jalur dibaca dari halaman polis yang sama (`OQ` - OfferFacIn.QuotationData).
// Penyalinan itu juga berarti halaman `Quotation` TIDAK disimpan terpisah (spec Bab 17 titik 1): dibaca dari polis.
//
// ⚠️ Perbandingan Pega `=` atas teks peka huruf. `isElectronicEquipment` membandingkan dengan " ElectronicEquipment "
// (BERSPASI, verbatim XML) - ditiru apa adanya. `IsBondingAndCustomBonds` C merujuk rule `IsCustomBonds` yang tidak ada
// di ekspor (yang ada `IsCustomBond`) - Pega menilai rule tak ditemukan sebagai galat; di sini bernilai salah (PARITAS).

// OQ - halaman polis yang diuji rule lini (`pyWorkPage.OfferFacIn.QuotationData`, disalin ke `pyWorkPage.Quotation`).
const OQ = "OfferFacIn.QuotationData."

func jenisBisnis(h *Halaman) string { return h.Ambil(OQ + "BusinessType") }
func kodeBisnis(h *Halaman) string  { return h.Ambil(OQ + "BusinessCode") }

func jenisSama(h *Halaman, v string) bool { return jenisBisnis(h) == v }

// IsKPR / IsOilGas / IsFireStyle1 / IsFireStyle2 (`pyWorkPage.Quotation.BusinessType`).
func IsKPR(h *Halaman) bool        { return jenisSama(h, "KPR") }
func IsOilGas(h *Halaman) bool     { return jenisSama(h, "OilGas") }
func IsFireStyle1(h *Halaman) bool { return jenisSama(h, "FireStyle1") }
func IsFireStyle2(h *Halaman) bool { return jenisSama(h, "FireStyle2") }

// IsFire = A IsKPR OR B IsOilGas OR C IsFireStyle1 OR D IsFireStyle2 OR E BusinessType = "Fire".
func IsFire(h *Halaman) bool {
	return IsKPR(h) || IsOilGas(h) || IsFireStyle1(h) || IsFireStyle2(h) || jenisSama(h, "Fire")
}

// IsMBU = BusinessType "MBUCar" OR "MBUMotorCycle".
func IsMBU(h *Halaman) bool { return jenisSama(h, "MBUCar") || jenisSama(h, "MBUMotorCycle") }

// IsPA / IsTravel / IsMarineCargo / isGolfInsurance.
func IsPA(h *Halaman) bool            { return jenisSama(h, "PA") }
func IsTravel(h *Halaman) bool        { return jenisSama(h, "Travel") }
func IsMarineCargo(h *Halaman) bool   { return jenisSama(h, "MarineCargo") }
func IsGolfInsurance(h *Halaman) bool { return jenisSama(h, "GolfInsurance") }

// Rule anak IsAneka.
func IsMarineHull(h *Halaman) bool                { return jenisSama(h, "MarineHull") }
func IsGrowingTrees(h *Halaman) bool              { return jenisSama(h, "GrowingTrees") }
func IsExclusion(h *Halaman) bool                 { return jenisSama(h, "Exclusion") }
func IsCAR(h *Halaman) bool                       { return jenisSama(h, "Car") }
func IsMaintenance(h *Halaman) bool               { return jenisSama(h, "Maintenance") }
func IsElectronicEquipment(h *Halaman) bool       { return jenisSama(h, " ElectronicEquipment ") }
func IsAviationHull(h *Halaman) bool              { return jenisSama(h, "AviationHull") }
func IsEar(h *Halaman) bool                       { return jenisSama(h, "Ear") }
func IsAllRisk(h *Halaman) bool                   { return jenisSama(h, "AllRisk") }
func IsGlass(h *Halaman) bool                     { return jenisSama(h, "Glass") }
func IsFidelity(h *Halaman) bool                  { return jenisSama(h, "Fidelity") }
func IsBurglary(h *Halaman) bool                  { return jenisSama(h, "Burglary") }
func IsCIT(h *Halaman) bool                       { return jenisSama(h, "CIT") }
func IsCIS(h *Halaman) bool                       { return jenisSama(h, "CIS") }
func IsMBD(h *Halaman) bool                       { return jenisSama(h, "MBD") }
func IsBoiler(h *Halaman) bool                    { return jenisSama(h, "Boiler") }
func IsHE(h *Halaman) bool                        { return jenisSama(h, "HE") }
func IsLandRig(h *Halaman) bool                   { return jenisSama(h, "LandRig") }
func IsContractorsPlantMachinery(h *Halaman) bool { return jenisSama(h, "ContractorsPM") }
func IsCustomBond(h *Halaman) bool                { return jenisSama(h, "CustomBond") }
func IsBondingKBG(h *Halaman) bool                { return jenisSama(h, "BondingKBG") }

// IsYieldShortfall / IsCrime / IsEnvironmental = BusinessType "Aneka" AND BusinessCode 10167 / 10168 / 10169.
func IsYieldShortfall(h *Halaman) bool { return jenisSama(h, "Aneka") && kodeBisnis(h) == "10167" }
func IsCrime(h *Halaman) bool          { return kodeBisnis(h) == "10168" && jenisSama(h, "Aneka") }
func IsEnvironmental(h *Halaman) bool  { return kodeBisnis(h) == "10169" && jenisSama(h, "Aneka") }

// IsBillboardNeon = BusinessName "BILLBOARD/NEON SIGN" (`pyWorkPage.Quotation.BusinessName`).
func IsBillboardNeon(h *Halaman) bool { return h.Ambil(OQ+"BusinessName") == "BILLBOARD/NEON SIGN" }

// IsBillboardNeonSyariah = BusinessType "BillboardNeonSyariah" OR BusinessCode "10106".
func IsBillboardNeonSyariah(h *Halaman) bool {
	return jenisSama(h, "BillboardNeonSyariah") || kodeBisnis(h) == "10106"
}

// IsProductsLiability / IsProfessionalLiability = "Aneka" AND BusinessCode 10021 / 10023.
func IsProductsLiability(h *Halaman) bool { return jenisSama(h, "Aneka") && kodeBisnis(h) == "10021" }
func IsProfessionalLiability(h *Halaman) bool {
	return jenisSama(h, "Aneka") && kodeBisnis(h) == "10023"
}

// IsWorkmenCompensation = BusinessType "Workmen".
func IsWorkmenCompensation(h *Halaman) bool { return jenisSama(h, "Workmen") }

// IsLiability = "Liability" OR IsProductsLiability OR IsProfessionalLiability OR IsWorkmenCompensation OR BusinessCode
// "10048" OR "10184" OR IsBillboardNeon.
func IsLiability(h *Halaman) bool {
	return jenisSama(h, "Liability") || IsProductsLiability(h) || IsProfessionalLiability(h) ||
		IsWorkmenCompensation(h) || kodeBisnis(h) == "10048" || kodeBisnis(h) == "10184" || IsBillboardNeon(h)
}

// IsBondingAndCustomBonds = "Bonding" OR IsBondingKBG OR IsCustomBonds (rule tidak diekspor -> salah).
func IsBondingAndCustomBonds(h *Halaman) bool { return jenisSama(h, "Bonding") || IsBondingKBG(h) }

// IsBonding = `.Quotation.BusinessType` "Bonding" OR IsBondingKBG.
func IsBonding(h *Halaman) bool { return jenisSama(h, "Bonding") || IsBondingKBG(h) }

// IsAneka = logika A OR B OR C OR E..AA (label D tidak ada di XML).
func IsAneka(h *Halaman) bool {
	return IsLiability(h) || IsMarineHull(h) || IsGrowingTrees(h) || IsExclusion(h) || IsCAR(h) ||
		IsMaintenance(h) || IsElectronicEquipment(h) || IsAviationHull(h) || IsEar(h) || IsAllRisk(h) ||
		IsGlass(h) || IsFidelity(h) || IsBillboardNeonSyariah(h) || IsBurglary(h) || IsCIT(h) || IsCIS(h) ||
		IsMBD(h) || IsBoiler(h) || IsHE(h) || jenisSama(h, "Aneka") || IsLandRig(h) ||
		IsContractorsPlantMachinery(h) || IsYieldShortfall(h) || IsCrime(h) || IsEnvironmental(h) ||
		IsBondingAndCustomBonds(h)
}

// IsSPK = `pyWorkPage.OfferFacIn.IsB2B = "SPK"` (Decision4 "B2B" Register_Flow).
func IsSPK(h *Halaman) bool { return h.Ambil("OfferFacIn.IsB2B") == "SPK" }

// IsBackStage = `.pyNote = "Back"` (Decision5 / Decision8 Register_Flow).
func IsBackStage(h *Halaman) bool { return h.Ambil("pyNote") == "Back" }

// IsNotTravelPA = `.Policy.Quotation.BusinessType != "Travel" AND != "PA"`.
func IsNotTravelPA(h *Halaman) bool { return !IsTravel(h) && !IsPA(h) }

// isPA_PNC = `pyWorkPage.Quotation.GroupPanel = "002"`.
func IsPAPNC(h *Halaman) bool { return h.Ambil(OQ+"GroupPanel") == "002" }

// IsEDM = `pyWorkPage.Quotation.StatusBusiness = 3`.
func IsEDM(h *Halaman) bool { return h.Ambil(OQ+"StatusBusiness") == "3" }

// IsClaim = `pyWorkPage.pyWorkIDPrefix = "CLM-"` - setiap kasus modul ini (kunci tetap, OQ-CFI-02).
func IsClaim(*Halaman) bool { return true }

// Lini - golongan grid objek / item per lini (section memilih varian dengan When yang sama).
type Lini string

// Golongan grid (InputRegisterDetail / InputEstimasiDetail / ShowObjectAdj).
const (
	LiniFire   Lini = "fire"
	LiniAneka  Lini = "aneka"
	LiniGolf   Lini = "golf"
	LiniMarine Lini = "marine"
	LiniMBU    Lini = "mbu"
	LiniPA     Lini = "pa"
	LiniTravel Lini = "travel"
	LiniLain   Lini = ""
)

// GolonganLini - varian grid yang tampil untuk kasus ini. Urutan uji mengikuti urutan cabang `InsertObjects_dt` langkah
// 8.1-8.7 (Fire, MBU, PA, Travel, Golf, Aneka, MarineCargo). Lini yang tidak dikenali -> LiniLain (nol grid objek; AC 17).
func GolonganLini(h *Halaman) Lini {
	switch {
	case IsFire(h):
		return LiniFire
	case IsMBU(h):
		return LiniMBU
	case IsPA(h):
		return LiniPA
	case IsTravel(h):
		return LiniTravel
	case IsGolfInsurance(h):
		return LiniGolf
	case IsAneka(h):
		return LiniAneka
	case IsMarineCargo(h):
		return LiniMarine
	}
	return LiniLain
}
