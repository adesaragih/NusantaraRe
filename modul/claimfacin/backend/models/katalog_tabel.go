package models

// Untuk apa berkas ini: SIMPUL KATALOG - tabel kepala, pohon objek -> item -> estimasi / spreading / Break QS /
// adjustment -> cucu adjustment, dan retro tingkat klaim (susunan: `katalog.go`).

// ---------------------------------------------------------------- kepala

// TabelHeaderKlaim - T_GENERAL_CLAIM: satu baris per klaim, shared PK dengan T_WORK_CLAIM. Kolom kunci (ID) dan SUMBER
// ditulis repository di luar katalog. Kolom yang sudah ada dipakai ulang bila maknanya sama dengan Claim Prop; kolom
// khas FACIN lahir di migrasi 560 (ALTER ADD nullable).
var TabelHeaderKlaim = Tabel{Nama: "T_GENERAL_CLAIM", Kolom: []Kolom{
	// nomor klaim (CLaimFaceSheet_Act 12-17)
	kKode(CD+"NoClaim", "CLAIM_NO", 64),
	// polis (CopyNB_Act; halaman polis dibaca ulang dari JSON_POLIS menurut POLICY_NO + PRODKE)
	kKode(JalurNoPolis, "POLICY_NO", 64),
	baru(kKode(JalurProdke, "PRODKE", 8)),
	kTeks(OQ+"InsuredName", "INSURED_NAME", 500),
	baru(kTeks(OQ+"QQName", "QQ_NAME", 500)),
	baru(kTeks(OQ+"CedingCoName", "CEDING_CO_NAME", 500)),
	baru(kTeks(OQ+"SobName", "SOB_NAME", 500)),
	kKode(CD+"TreatyGroupID", "TREATY_GROUP_ID", 64), // CheckLimit_Act1 4.1.3
	kKode(CD+"TreatyYear", "TREATY_YEAR", 16),
	// registrasi (InputRegisterDetail)
	kTeks(CD+"PlaNoCeding", "PLA_NO_CEDING", 255),
	kTeks(CD+"PlaNoSOB", "PLA_NO_SOB", 255),
	kTgl(CD+"DateOfLoss", "DATE_OF_LOSS"),
	kTgl(CD+"ReportDate", "REPORT_DATE"),
	kTgl(CD+"DateReceived", "DATE_RECEIVED"),
	kTeks(CD+"ReporterName", "REPORTER_NAME", 255),
	kTeks(CD+"ReporterTelp", "REPORTER_PHONE", 64),
	kKode(CD+"ReportType", "REPORT_TYPE", 16),
	kKode(CD+"InsuredRelationship", "REPORTER_STATUS", 16),
	kTeks(CD+"InsuredRelationshipOthers", "INSURED_RELATIONSHIP_OTHERS", 500),
	kTeks(CD+"ReportAddress", "REPORT_ADDRESS", 2000),
	kTeks(CD+"ReportDescription", "REPORT_DESCRIPTION", 4000),
	kKode(CD+"StsKatastrofe", "CATASTROPHE_STATUS", 32),
	kKode(CD+"NonKatastrofeType", "NON_CATASTROPHE_TYPE", 32),
	kKode(CD+"KatastrofeID", "CATASTROPHE_ID", 100),
	kTeks(CD+"KatastrofeNote", "KATASTROFE_NOTE", 1000),
	kTeks(CD+"CauseOfLoss", "CAUSE_OF_LOSS", 1000),
	kKode(CD+"CauseOfLossID", "CAUSE_OF_LOSS_ID", 100),
	kKode(CD+"ConsultantID", "CONSULTANT_ID", 100),
	kTeks(CD+"ConsultantName", "CONSULTANT_NAME", 500),
	kKode(CD+"AppointedADJID", "ADJUSTER_ID", 100),
	kTeks(CD+"AppointedADJ", "ADJUSTER_NAME", 500),
	kTeks(CD+"Location", "LOCATION_OF_LOSS", 4000),
	baru(kTeks(CD+"Country", "COUNTRY", 255)),
	baru(kKode(CD+"CountryID", "COUNTRY_ID", 64)),
	kTeks(CD+"Province", "PROVINCE", 255),
	kKode(CD+"ProvinceID", "PROVINCE_ID", 64),
	kTeks(CD+"City", "CITY", 255),
	kKode(CD+"CityID", "CITY_ID", 64),
	kTeks(CD+"District", "DISTRICT", 255),
	kKode(CD+"DistrictID", "DISTRICT_ID", 64),
	kTeks(CD+"RW", "RW", 255),
	kKode(CD+"RWID", "RW_ID", 64),
	kKode(CD+"PostalCode", "POSTAL_CODE", 16),
	kKode(CD+"Currency", "CURRENCY", 8),
	baru(kUang(CD+"DollarCurrencyVal", "CURRENCY_VALUE_IDR")),
	baru(kUang(CD+"GrossEstimate", "GROSS_ESTIMATE")),
	baru(kUang(CD+"ClaimEstimate", "CLAIM_ESTIMATE")),
	baru(kPenanda(CD+"ExGratia", "EX_GRATIA")),
	kTeks(CD+"Remark", "REMARK", 4000),
	kTeks(CD+"Remark_Close", "REMARK_CLOSE", 4000),
	// penanda alur (pyWorkPage)
	baru(kPenanda(JalurIsError, "IS_ERROR")),
	baru(kPenanda(JalurIsRegister, "IS_REGISTER")),
	baru(kPenanda(JalurIsEstimation, "IS_ESTIMATION")),
	baru(kPenanda(JalurIsAdjustment, "IS_ADJUSTMENT")),
	kPenanda(JalurIsAnyAccept, "IS_ANY_ACCEPTATION"),
	kPenanda(JalurIsCFS, "IS_CFS"),
	baru(kPenanda(JalurIsPicTransfer, "IS_PIC_TRANSFER")),
	baru(kPenanda(JalurIsTreatyOut, "IS_TREATY_OUT")),
	kPenanda(JalurAktifButton, "AKTIF_BUTTON"),
	baru(kTeks(JalurCatatan, "PY_NOTE", 255)),
	baru(kTglWaktu("StartDateRegister", "START_DATE_REGISTER")),
	baru(kTglWaktu("EndDateRegister", "END_DATE_REGISTER")),
	kTglWaktu("StartDateEstimation", "START_DATE_ESTIMATION"),
	kTglWaktu("EndDateEstimation", "END_DATE_ESTIMATION"),
	kTglWaktu("StartDateAdjustment", "START_DATE_ADJUSTMENT"),
	// pemasaran (SetMOClaim_Act)
	kKode(CD+"MarketingData.ID", "MARKETING_ID", 64),
	kKode(CD+"MarketingData.ClientID", "MARKETING_CLIENT_ID", 64),
	kTeks(CD+"MarketingData.ClientName", "MARKETING_CLIENT_NAME", 500),
	kKode(CD+"MarketingData.TeamGroup", "MARKETING_TEAM_GROUP", 64),
	kKode(CD+"MarketingData.BranchDetailID", "MARKETING_BRANCH_ID", 64),
	kTeks(CD+"MarketingData.BranchDetailName", "MARKETING_BRANCH_NAME", 500),
	// pembayaran tingkat klaim (InputAdjustment LS30 `pyWorkPage.ClaimData.Payable`)
	kKode(CD+"Payable", "PAYABLE", 16),
	kTeks(CD+"PayableTo", "PAYABLE_TO", 500),
}}

// ---------------------------------------------------------------- objek dan item

// TabelObjek - T_CLAIM_OBJECT (baru) <- ClaimData.ObjectList. Medan identitas objek ditulis sebagai salinan baca;
// `LengkapiObjek` menimpanya dari polis saat dimuat (KUNCI_POLIS = letak objek di halaman polis).
var TabelObjek = Tabel{Nama: "T_CLAIM_OBJECT", Daftar: DaftarObjek, KolomInduk: "CLAIM_ID", Stabil: true,
	Kolom: []Kolom{
		baru(kKode(PropKunciPolis, "KUNCI_POLIS", 64)),
		baru(kKode("ObjectID", "OBJECT_REF", 200)),
		baru(kTeks("ObjectName", "OBJECT_NAME", 2000)),
		baru(kPenanda("CFS", "CFS")),
		baru(kPenanda("PrintFaceClaim", "PRINT_FACE_CLAIM")),
		baru(kPenanda("PlaStatus", "PLA_STATUS")),
		baru(kPenanda("IsFacretro", "IS_FAC_RETRO")),
		baru(kPenanda("IsMoreThanTreatyLimit", "IS_MORE_THAN_TREATY_LIMIT")),
		baru(kPenanda("IsKomite", "IS_KOMITE")),
		baru(kPenanda("IsPrintAccept", "IS_PRINT_ACCEPT")),
		baru(kTeks("RemarksPLA", "REMARKS_PLA", 4000)),
		baru(kTeks("RemarksDLA", "REMARKS_DLA", 4000)),
		baru(kKode("ShareRetro", "SHARE_RETRO", 16)),
		baru(kKode("NoDLA", "NO_DLA", 64)),
		baru(kPenanda("DLAStatus", "DLA_STATUS")),
		baru(kPenanda("SaveSpreading", "SAVE_SPREADING")),
	}}

// TabelItem - T_CLAIM_OBJECT_ITEM (baru) <- ObjectList(o).ObjectItemList.
var TabelItem = Tabel{Nama: "T_CLAIM_OBJECT_ITEM", Daftar: AnakItem, KolomInduk: "OBJECT_ID", DenganKlaim: true,
	Stabil: true, Kolom: []Kolom{
		baru(kTeks("ObjectItemName", "OBJECT_ITEM_NAME", 2000)),
		baru(kKode("ObjectItemID", "OBJECT_ITEM_REF", 100)),
		baru(kKode("IndexPropertyItem", "INDEX_PROPERTY_ITEM", 16)),
		baru(kKode("IndexCoverage", "INDEX_COVERAGE", 16)),
		baru(kKode("IndexAneka", "INDEX_ANEKA", 16)),
		baru(kKode("IndexOccupation", "INDEX_OCCUPATION", 16)),
		baru(kKode("CoverageID", "COVERAGE_ID", 100)),
		baru(kKode("CoverageOLDID", "COVERAGE_OLD_ID", 100)),
		baru(kTeks("CoverageNote", "COVERAGE_NOTE", 4000)),
		baru(kKode("OccupationId", "OCCUPATION_ID", 100)),
		baru(kTeks("OccupationName", "OCCUPATION_NAME", 2000)),
		baru(kKode("CurrencyID", "CURRENCY_ID", 64)),
		baru(kTeks("Currency", "CURRENCY_NAME", 64)),
		baru(kUang("KursObjectItem", "KURS_OBJECT_ITEM")),
		baru(kUang("TSIPerObject", "TSI_PER_OBJECT")),
		baru(kUang("TSINusare", "TSI_NUSARE")),
		baru(kUang("ValueTSINusareIDR", "VALUE_TSI_NUSARE_IDR")),
		baru(kUang("LimitofLiability", "LIMIT_OF_LIABILITY")),
		baru(kUang("PremiNusare", "PREMI_NUSARE")),
		baru(kUang("TotalGrossPremi", "TOTAL_GROSS_PREMI")),
		baru(kUang("TotalPremiumNusantaraRe", "TOTAL_PREMIUM_NUSARE")),
		// deductible (section Estimasi "Deductible Info")
		baru(kPenanda("DeductibleType", "DEDUCTIBLE_TYPE")),
		baru(kKode("FormType", "DEDUCTIBLE_FORM_TYPE", 16)),
		baru(kKode("CurrencyDeductible", "DEDUCTIBLE_CURRENCY_ID", 64)),
		baru(kPersen("Amount", "DEDUCTIBLE_PCT")),
		baru(kKode("TypeDeductible", "DEDUCTIBLE_BASIS", 16)),
		baru(kUang("TSIDeductible", "TSI_DEDUCTIBLE")),
		baru(kUang("ClaimDeductible", "CLAIM_DEDUCTIBLE")),
		baru(kUang("DeductibleValue", "DEDUCTIBLE_VALUE")),
		baru(kUang("NetDeductibleValue", "NET_DEDUCTIBLE_VALUE")),
		// total estimasi / spreading (CheckEstimateValue, CountSpreadingClaim_ACT, CopySpreading_Act)
		baru(kUang("TotalEstimasi", "TOTAL_ESTIMASI")),
		baru(kUang("TotalGrossEstimasi", "TOTAL_GROSS_ESTIMASI")),
		baru(kUang("TotalGrossEstimasiIDR", "TOTAL_GROSS_ESTIMASI_IDR")),
		baru(kUang("TotalClaimSpreaded", "TOTAL_CLAIM_SPREADED")),
		baru(kUang("TotalEstimationValueinIDR", "TOTAL_ESTIMATION_VALUE_IDR")),
		baru(kUang("EstimastionReserve", "ESTIMATION_RESERVE")),
		baru(kPersen("TotalSharePercentage", "TOTAL_SHARE_PCT")),
		baru(kUang("TotalTSISpreaded", "TOTAL_TSI_SPREADED")),
		baru(kUang("TotalPremiumSpreaded", "TOTAL_PREMIUM_SPREADED")),
		// penanda
		baru(kPenanda("PrintFaceClaim", "PRINT_FACE_CLAIM")),
		baru(kPenanda("IsFacretro", "IS_FAC_RETRO")),
		baru(kPenanda("IsMoreThanTreatyLimit", "IS_MORE_THAN_TREATY_LIMIT")),
		baru(kPenanda("IsKomite", "IS_KOMITE")),
		baru(kKode("PLARev_Treaty", "PLA_REV_TREATY", 16)),
		// adjustment tingkat item
		baru(kPenanda("AdjustmentVal", "ADJUSTMENT_VAL")),
		baru(kPenanda("IsAdjVal", "IS_ADJ_VAL")),
		baru(kKode("PaymentType", "PAYMENT_TYPE", 16)),
		baru(kKode(PropIndeksAdj, "INDEX_ADJUSTMENT", 16)),
		baru(kUang("TotalHasilClaim", "TOTAL_HASIL_CLAIM")),
		baru(kUang("TotalEstimasiReas", "TOTAL_ESTIMASI_REAS")),
		baru(kPersen("PercentReas", "PERCENT_REAS")),
		baru(kPersen("PercentHandlingFee", "PERCENT_HANDLING_FEE")),
		baru(kPersen("PercentFacoutClaim", "PERCENT_FACOUT_CLAIM")),
	}}

// ---------------------------------------------------------------- di bawah item

// TabelEstimasi - T_CLAIM_ESTIMATION <- ObjectItemList(i).EstimationList (kolom Claim Prop + ADD 562).
var TabelEstimasi = Tabel{Nama: "T_CLAIM_ESTIMATION", Daftar: AnakEstimasi, KolomInduk: "OBJECT_ITEM_ID",
	DenganKlaim: true, Kolom: []Kolom{
		kTgl("EstimationDate", "ESTIMATION_DATE"),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
		kUang("KursValue", "KURS"),
		kUang("GrossEstimationPct", "GROSS_ESTIMATION_VALUE"),
		kUang("ConvertGrossEstimasi", "GROSS_ESTIMATION_IDR"),
		kPersen("PersenRNM", "PERSEN_RNM"),
		kUang("EstimationValue", "ESTIMATION_VALUE"),
		kUang("ConvertValue", "ESTIMATION_VALUE_IDR"),
		kPenanda("PrintFaceClaim", "IS_PRINT_FACE_CLAIM"),
		baru(kUang("GrossEstimationPctMBU", "GROSS_ESTIMATION_MBU")),
		baru(kUang("Deductible", "DEDUCTIBLE")),
		baru(kUang("NetEstimationValue", "NET_ESTIMATION_VALUE")),
		baru(kPenanda("EstimasiMoreThanTSI", "IS_MORE_THAN_TSI")),
		baru(kKode("pxCreateOperator", "CREATED_BY", 64)),
	}}

// kolomSpreading - baris spreading item (SpreadingList polis / SpreadingClaim klaim): kolom Claim Prop + ADD 563.
func kolomSpreading() []Kolom {
	return []Kolom{
		kKode("TreatyType", "TREATY_ID", 64),
		kTeks("TreatyName", "TREATY_NAME", 255),
		kPersen("SharePercentage", "SHARE_PERCENTAGE"),
		kUang("ClaimSpreaded", "CLAIM_SPREADED"),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
		baru(kUang("TSISpreaded", "TSI_SPREADED")),
		baru(kUang("PremiumSpreaded", "PREMIUM_SPREADED")),
	}
}

// Nilai kolom JENIS T_CLAIM_SPREADING (ADD 563; baris Claim Prop / Non Prop tetap NULL).
const (
	JenisSpreadPolis = "POLIS"
	JenisSpreadKlaim = "KLAIM"
)

// Tabel spreading dan Break QS tingkat item.
var (
	TabelSpreadPolis = Tabel{Nama: "T_CLAIM_SPREADING", Daftar: AnakSpreadPolis, KolomInduk: "OBJECT_ITEM_ID",
		DenganKlaim: true, Jenis: JenisSpreadPolis, Kolom: kolomSpreading()}
	TabelSpreadKlaim = Tabel{Nama: "T_CLAIM_SPREADING", Daftar: AnakSpreadKlaim, KolomInduk: "OBJECT_ITEM_ID",
		DenganKlaim: true, Jenis: JenisSpreadKlaim, Kolom: kolomSpreading()}
	TabelBreakQS = Tabel{Nama: "T_CLAIM_BREAK_QS", Daftar: AnakBreakQS, KolomInduk: "OBJECT_ITEM_ID",
		DenganKlaim: true, Kolom: []Kolom{
			kKode("TreatyType", "TREATY_ID", 64),
			kTeks("TreatyName", "TREATY_NAME", 255),
			kPersen("SharePercentage", "SHARE_PERCENTAGE"),
			kKode("CurrencyID", "CURRENCY_ID", 64),
			kTeks("Currency", "CURRENCY_NAME", 64),
		}}
)

// kolomSpreadAdj - Spreading In / Quota Share adjustment (`SpreadingAdjustment` / `SpreadingQuotaShare`).
func kolomSpreadAdj() []Kolom {
	return []Kolom{
		kKode("TreatyType", "TREATY_ID", 64),
		kTeks("TreatyName", "TREATY_NAME", 255),
		kPersen("SharePercentage", "SHARE_PERCENTAGE"),
		kUang("ClaimSpreaded", "CLAIM_SPREADED"),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
	}
}

// kolomRetro - baris retro fakultatif (`FacRetroList`).
func kolomRetro() []Kolom {
	return []Kolom{
		kKode("ReinsurerID", "REINSURER_ID", 64),
		kTeks("ReinsurerName", "REINSURER_NAME", 500),
		kPersen("PctShareAllObj", "SHARE_PCT"),
		kPersen("RiCommAllObj", "RI_COMMISSION_PCT"),
		kTeks("AdditionalInfo", "ADDITIONAL_INFO", 1000),
		kUang("TotalEstimasiReas", "TOTAL_ESTIMATION_REINS"),
	}
}

// Tabel cucu adjustment.
var (
	TabelAdjSpreading = Tabel{Nama: "T_CLAIM_ADJ_SPREADING", Daftar: AnakAdjSpread, KolomInduk: "ADJUSTMENT_ID",
		Kolom: kolomSpreadAdj()}
	TabelAdjQuotaShare = Tabel{Nama: "T_CLAIM_ADJ_QUOTA_SHARE", Daftar: AnakAdjQS, KolomInduk: "ADJUSTMENT_ID",
		Kolom: kolomSpreadAdj()}
	TabelAdjRetro = Tabel{Nama: "T_CLAIM_FAC_RETRO", Daftar: AnakFacRetro, KolomInduk: "ADJUSTMENT_ID",
		DenganKlaim: true, Kolom: kolomRetro()}
)

// TabelAdjustment - T_CLAIM_ADJUSTMENT <- ObjectItemList(i).Adjustment (ID stabil: T_GENERAL_KOMITE menunjuknya).
// Kolom Claim Prop + ADD 564 (khas FAC: InputAdjustment).
var TabelAdjustment = Tabel{Nama: "T_CLAIM_ADJUSTMENT", Daftar: AnakAdj, KolomInduk: "OBJECT_ITEM_ID",
	DenganKlaim: true, Stabil: true, Kolom: []Kolom{
		kKode("PaymentType", "PAYMENT_TYPE", 16),
		kTeks("FormType", "FORM_TYPE", 255),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
		kUang("KursValue", "KURS"),
		kPersen("PersenRNM", "PERSEN_RNM"),
		kUang("GrossAdjustment", "GROSS_ADJUSTMENT"),
		kUang("GrossAdjustmentIDR", "GROSS_ADJUSTMENT_IDR"),
		kUang("GrossValue", "GROSS_VALUE"),
		kKode("IndividualRiskType", "INDIVIDUAL_RISK_TYPE", 16),
		kPersen("IndividualRiskPercentage", "INDIVIDUAL_RISK_PCT"),
		kUang("IndividualRiskValue", "INDIVIDUAL_RISK_VALUE"),
		kUang("IndividualRiskRNM", "INDIVIDUAL_RISK_RNM"),
		kUang("ProposeAdjustmentValue", "PROPOSE_ADJUSTMENT_VALUE"),
		kUang("AdjustmentValue", "ADJUSTMENT_VALUE"),
		kUang("ValueAdjustment", "VALUE_ADJUSTMENT_IDR"),
		kUang("TotalEstimasiValue", "TOTAL_ESTIMATION_VALUE"),
		kUang("AdjusterFeeValue", "ADJUSTER_FEE_VALUE"),
		kUang("SalvageValue", "SALVAGE_VALUE"),
		kKode("Payable", "PAYABLE", 16),
		kTeks("PayableTo", "PAYABLE_TO", 500),
		kTeks("NameOfBank", "BANK_NAME", 500),
		kTeks("BranchOfBank", "BANK_BRANCH", 500),
		kTeks("NoAccount", "BANK_ACCOUNT_NO", 100),
		kKode("SwiftCode", "SWIFT_CODE", 64),
		kKode("IDOfBank", "BANK_ID", 64),
		kTeks("DLANoCeding", "DLA_NO_CEDING", 255),
		kTeks("DLANoSOB", "DLA_NO_SOB", 255),
		kPenanda("DirectToKasir", "IS_DIRECT_TO_KASIR"),
		kTeks("StatusKasir", "STATUS_KASIR", 1000),
		milikKomite(kKode("AcceptedNo", "ACCEPTED_NO", 64)),
		milikKomite(kTglWaktu("AcceptedDate", "ACCEPTED_DATE")),
		milikKomite(kKode("AcceptanceStatus", "ACCEPTANCE_STATUS", 16)),
		milikKomite(kPenanda("IsApproved", "IS_APPROVED")),
		kPenanda("IsKomite", "IS_KOMITE"),
		kKode("DLA_No", "DLA_NO", 64),
		kTeks("RemarksDLA", "REMARKS_DLA", 4000),
		kPenanda("IsFacRetro", "IS_FAC_RETRO"),
		kPenanda("IsPrintAccept", "IS_PRINT_ACCEPT"),
		kTeks("DataCommitteFacin.CircumCauseOfLoss", "KOMITE_CIRCUM_CAUSE_OF_LOSS", 4000),
		kTeks("DataCommitteFacin.AdjusterFee", "KOMITE_ADJUSTER_FEE", 4000),
		kTeks("DataCommitteFacin.Remarks", "KOMITE_REMARKS", 4000),
		kTeks("DataCommitteFacin.Salvage", "KOMITE_SALVAGE", 4000),
		kTeks("DataCommitteFacin.LegalLiability", "KOMITE_LEGAL_LIABILITY", 4000),
		kTeks("DataCommitteFacin.ExtentOfLoss", "KOMITE_EXTENT_OF_LOSS", 4000),
		kKode("pxCreateOperator", "CREATED_BY", 64),
		kTeks("pxCreateOpName", "CREATED_BY_NAME", 128),
		kTglWaktu("pxCreateDateTime", "CREATED_AT"),
		baru(kPenanda("NetForCollection", "NET_FOR_COLLECTION")),
		baru(kTeks("CurrencyEstimasi", "CURRENCY_ESTIMASI", 64)),
		baru(kUang("EstimationValue", "ESTIMATION_VALUE")),
		baru(kKode("UploadLOD", "CURRENCY_CHOICE_ID", 64)),
		baru(kUang("SurveyExpenses", "SURVEY_EXPENSES")),
		baru(kPersen("VAT", "VAT_PCT")),
		baru(kUang("VATValue", "VAT_VALUE")),
		baru(kUang("ProfessionalFee", "PROFESSIONAL_FEE")),
		baru(kPenanda("ExGratia", "EX_GRATIA")),
		baru(kPersen("TotalSharePersen", "TOTAL_SHARE_PCT")),
		baru(kUang("TotalSpreadAdjustment", "TOTAL_SPREAD_ADJUSTMENT")),
		baru(kUang("TotalSpreadBreakQs", "TOTAL_SPREAD_BREAK_QS")),
		baru(kKode(PropPilihCedant, "CEDANT_CHOICE", 64)),
	},
	Anak: []*Tabel{&TabelAdjSpreading, &TabelAdjQuotaShare, &TabelAdjRetro}}

// TabelRetroTreaty - T_CLAIM_FAC_RETRO tingkat klaim (ADJUSTMENT_ID NULL) <- ClaimData.FacRetroTreaty (pengganti
// FacRetroList polis yang ditulis CheckLimit_Act1 9.1.1.4).
var TabelRetroTreaty = Tabel{Nama: "T_CLAIM_FAC_RETRO", Daftar: DaftarFacRetroTreaty, KolomInduk: "CLAIM_ID",
	Kolom: kolomRetro()}

func init() {
	TabelItem.Anak = []*Tabel{&TabelEstimasi, &TabelSpreadPolis, &TabelSpreadKlaim, &TabelBreakQS, &TabelAdjustment}
	TabelObjek.Anak = []*Tabel{&TabelItem}
}

// SemuaTabel - seluruh simpul katalog (kepala lebih dulu, lalu pohon objek pra-urut, lalu retro tingkat klaim).
func SemuaTabel() []*Tabel {
	out := []*Tabel{&TabelHeaderKlaim}
	var jalan func(t *Tabel)
	jalan = func(t *Tabel) {
		out = append(out, t)
		for _, a := range t.Anak {
			jalan(a)
		}
	}
	jalan(&TabelObjek)
	return append(out, &TabelRetroTreaty)
}

// ProyeksiKatalog - halaman sebagaimana ia KELUAR dari penyimpanan: hanya medan berkolom (gudang tiruan memakainya
// supaya uji seam HTTP melihat apa yang Oracle simpan). Kronologi disalin utuh.
func ProyeksiKatalog(h *Halaman) *Halaman {
	s := HalamanBaru()
	for _, k := range TabelHeaderKlaim.Kolom {
		if v := h.Ambil(k.Properti); v != "" {
			s.Setel(k.Properti, v)
		}
	}
	var salin func(t *Tabel, daftar string)
	salin = func(t *Tabel, daftar string) {
		var out []Baris
		for _, x := range h.AmbilDaftar(daftar) {
			nb := Baris{}
			for _, k := range t.Kolom {
				if v := x[k.Properti]; v != "" {
					nb[k.Properti] = v
				}
			}
			if t.Stabil {
				for _, p := range []string{PropID, PropKomiteID} {
					if v := x[p]; v != "" {
						nb[p] = v
					}
				}
			}
			out = append(out, nb)
		}
		s.SetelDaftar(daftar, out)
		for n := range out {
			for _, a := range t.Anak {
				salin(a, JalurAnak(daftar, n+1, a.Daftar))
			}
		}
	}
	salin(&TabelObjek, DaftarObjek)
	salin(&TabelRetroTreaty, DaftarFacRetroTreaty)
	s.SetelDaftar(DaftarKronologi, SalinDaftar(h.AmbilDaftar(DaftarKronologi)))
	return s
}
