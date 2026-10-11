package models

// Untuk apa berkas ini: KLASIFIKASI LINI PRODUK - port rule When `Komite Claim FacIn/When/*.xml` yang dipakai
// `SaveAccept_ACT`, `SaveReject_ACT_KMT`, `SetDataForInformation_Act` (cabang IsFire / IsAneka / isGolfInsurance /
// IsMarineCargo / IsMBU / IsTravel / IsPA). DISALIN dari Claim Fac In `models/lini.go` (bukan impor): 47 rule When
// korpus komite identik isi dengan korpus Claim Fac In (bandingan logika + baris kondisi, 10-10-2026).
//
// Semua rule menguji `pyWorkPage.Quotation.BusinessType` / `.BusinessCode`; di kasus komite `pyWorkPage.Quotation` =
// `pyWorkCover.OfferFacIn.QuotationData` (SetValueKomite S12 `Quotation.BusinessType`, SaveReject_ACT_KMT S1 Page-Copy
// `TempOpenPage.Quotation`) - dibaca dari halaman polis klaim induk.
//
// Beda dengan Claim Fac In: `IsCustomBonds` ADA di korpus komite (`BusinessType = "CustomBond"`, kelas Work), sehingga
// cabang C `IsBondingAndCustomBonds` bernilai benar untuk CustomBond (di Claim Fac In rule itu tidak diekspor - salah).
// ⚠️ Perbandingan Pega `=` atas teks peka huruf; " ElectronicEquipment " BERSPASI ditiru apa adanya.

// OQ - halaman polis yang diuji rule lini.
const OQ = "OfferFacIn.QuotationData."

func jenisBisnis(v map[string]string) string { return v[OQ+"BusinessType"] }
func kodeBisnis(v map[string]string) string  { return v[OQ+"BusinessCode"] }

func jenis(v map[string]string, x string) bool { return jenisBisnis(v) == x }

// IsFire = IsKPR OR IsOilGas OR IsFireStyle1 OR IsFireStyle2 OR BusinessType "Fire".
func IsFire(v map[string]string) bool {
	return jenis(v, "KPR") || jenis(v, "OilGas") || jenis(v, "FireStyle1") || jenis(v, "FireStyle2") || jenis(v, "Fire")
}

// IsMBU = "MBUCar" OR "MBUMotorCycle".
func IsMBU(v map[string]string) bool { return jenis(v, "MBUCar") || jenis(v, "MBUMotorCycle") }

// IsPA / IsTravel / IsMarineCargo / isGolfInsurance.
func IsPA(v map[string]string) bool            { return jenis(v, "PA") }
func IsTravel(v map[string]string) bool        { return jenis(v, "Travel") }
func IsMarineCargo(v map[string]string) bool   { return jenis(v, "MarineCargo") }
func IsGolfInsurance(v map[string]string) bool { return jenis(v, "GolfInsurance") }

// isLiability = "Liability" OR ProductsLiability (Aneka + 10021) OR ProfessionalLiability (Aneka + 10023) OR
// WorkmenCompensation ("Workmen") OR BusinessCode 10048 / 10184 OR BillboardNeon (BusinessName "BILLBOARD/NEON SIGN").
func isLiability(v map[string]string) bool {
	return jenis(v, "Liability") || (jenis(v, "Aneka") && kodeBisnis(v) == "10021") ||
		(jenis(v, "Aneka") && kodeBisnis(v) == "10023") || jenis(v, "Workmen") || kodeBisnis(v) == "10048" ||
		kodeBisnis(v) == "10184" || v[OQ+"BusinessName"] == "BILLBOARD/NEON SIGN"
}

// isBondingAndCustomBonds = "Bonding" OR IsBondingKBG OR IsCustomBonds ("CustomBond", korpus komite).
func isBondingAndCustomBonds(v map[string]string) bool {
	return jenis(v, "Bonding") || jenis(v, "BondingKBG") || jenis(v, "CustomBond")
}

// IsAneka = A OR B OR C OR E..AA (label D tidak ada di XML).
func IsAneka(v map[string]string) bool {
	for _, x := range []string{"MarineHull", "GrowingTrees", "Exclusion", "Car", "Maintenance", " ElectronicEquipment ",
		"AviationHull", "Ear", "AllRisk", "Glass", "Fidelity", "BillboardNeonSyariah", "Burglary", "CIT", "CIS", "MBD",
		"Boiler", "HE", "Aneka", "LandRig", "ContractorsPM"} {
		if jenis(v, x) {
			return true
		}
	}
	aneka := jenis(v, "Aneka")
	return isLiability(v) || kodeBisnis(v) == "10106" || (aneka && kodeBisnis(v) == "10167") ||
		(aneka && kodeBisnis(v) == "10168") || (aneka && kodeBisnis(v) == "10169") || isBondingAndCustomBonds(v)
}

// Lini OS akseptasi (cabang SaveAccept_ACT / SaveReject_ACT_KMT).
const (
	LiniFireAneka = "fire" // IsFire / IsAneka / isGolfInsurance (SaveAccept_ACT S3; SaveReject S9-S11)
	LiniMarine    = "marine"
	LiniMBU       = "mbu"
	LiniTravel    = "travel"
	LiniPA        = "pa"
)

// LiniOS - cabang lini klaim induk. Urutan = urutan langkah SaveAccept_ACT (S3 Fire/Aneka/Golf, S4 MarineCargo, S5 MBU,
// S6 Travel, S7 PA); lini lain = "" (tanpa cabang: OS tanpa baris akseptasi).
func LiniOS(v map[string]string) string {
	switch {
	case IsFire(v) || IsAneka(v) || IsGolfInsurance(v):
		return LiniFireAneka
	case IsMarineCargo(v):
		return LiniMarine
	case IsMBU(v):
		return LiniMBU
	case IsTravel(v):
		return LiniTravel
	case IsPA(v):
		return LiniPA
	}
	return ""
}
