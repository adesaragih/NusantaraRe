package models

// Untuk apa berkas ini: KATALOG KOLOM - satu-satunya pemetaan properti Pega -> kolom Oracle -> golongan tipe untuk
// penyimpanan relasional Claim Prop (tiket 00; STRUKTUR-TABEL-CLAIM-PROP.md bab "Bentuk mengikat 07-10-2026").
//
// Migrasi `backend/migrations/52x_*.sql`, dokumen STRUKTUR, dan repository dicocokkan terhadap katalog ini oleh uji
// (`katalog_test.go`). Menambah kolom berarti menambah satu baris di sini, satu baris DDL, dan satu baris STRUKTUR.
//
// ⛔ Medan TURUNAN tidak punya kolom - dihitung sekali jalan dari basis asli (AC 22, `HitungTurunan`):
// TotalInterestInsured, ListTotalEstimation, TotalSumInsuredIDR, TotalListClaimAmount(IDR), TotalEstimasi(IDR),
// TotalGrossEstimateTreaty, TotalGrossEstimateIDR, EstimastionReserve, TotalGrossEstimasiIDR, TotalSpread,
// ClaimData.SpreadingAdjustment(QS), AdjustmentList(n).CurencyAdjustment / .ComiteeClaim / .TotalKomite.
// ⛔ Halaman kerja yang bukan data (STRUKTUR J1): InterestListDtl, PaymentData, ObjectList, Attachment,
// AttachmentPaid, AdjustmentList(n).FacRetroList (salinan ClaimData.FacRetroList).

// Golongan adalah kategori tipe logis sebuah kolom.
type Golongan string

const (
	GolTeks         Golongan = "teks"
	GolKode         Golongan = "kode"    // teks - nol di depan bermakna
	GolPenanda      Golongan = "penanda" // teks - "" berbeda dari "0"
	GolUang         Golongan = "uang"    // NUMBER(38,10)
	GolPersen       Golongan = "persen"  // NUMBER(38,10) - 12.5 = 12,5 persen
	GolTanggal      Golongan = "tanggal" // DATE, tanggal saja
	GolTanggalWaktu Golongan = "tanggal-waktu"
)

// Desimal - golongan bertipe NUMBER(38,10) (keputusan work owner 07-10-2026).
func (g Golongan) Desimal() bool { return g == GolUang || g == GolPersen }

// Tanggal - golongan bertipe DATE.
func (g Golongan) Tanggal() bool { return g == GolTanggal || g == GolTanggalWaktu }

// Kolom memetakan satu properti ke satu kolom.
type Kolom struct {
	// Properti - jalur relatif pyWorkPage (tabel induk) atau nama anggota daftar (tabel anak).
	Properti string
	Kolom    string
	Golongan Golongan
	// Panjang - panjang VARCHAR2; 0 untuk golongan non-teks.
	Panjang int
}

// Tabel - satu tabel katalog.
type Tabel struct {
	Nama string
	// Daftar - jalur PageList sumber baris (tabel anak). Untuk tabel cucu adjustment: nama daftar di dalam baris
	// `ClaimData.AdjustmentList(n)`.
	Daftar string
	// KolomInduk - kolom penunjuk induk (CLAIM_ID / ADJUSTMENT_ID).
	KolomInduk string
	Kolom      []Kolom
}

func kTeks(p, k string, n int) Kolom { return Kolom{p, k, GolTeks, n} }
func kKode(p, k string, n int) Kolom { return Kolom{p, k, GolKode, n} }
func kPenanda(p, k string) Kolom     { return Kolom{p, k, GolPenanda, 16} }
func kUang(p, k string) Kolom        { return Kolom{p, k, GolUang, 0} }
func kPersen(p, k string) Kolom      { return Kolom{p, k, GolPersen, 0} }
func kTgl(p, k string) Kolom         { return Kolom{p, k, GolTanggal, 0} }
func kTglWaktu(p, k string) Kolom    { return Kolom{p, k, GolTanggalWaktu, 0} }

// Awalan jalur halaman.
const (
	CD = "ClaimData."
	TM = "TreatyInMaster."
	// OQ - `pyWorkPage.OfferFacIn.QuotationData` (Class of Business klaim, diisi SetValueToClaim_Act langkah 2).
	OQ = "OfferFacIn.QuotationData."
)

// Jalur daftar tingkat klaim.
const (
	DaftarEstimasi     = CD + "EstimationList"
	DaftarInterest     = CD + "InterestList"
	DaftarClaimAmount  = CD + "ListClaimAmount"
	DaftarLossAlloc    = CD + "SpreadingRisk"
	DaftarSpreading    = CD + "SpreadingClaim"
	DaftarBreakQS      = CD + "SpreadingBreakQS"
	DaftarFacRetro     = CD + "FacRetroList"
	DaftarAdjustment   = CD + "AdjustmentList"
	DaftarRiwayat      = CD + "SuggestList"
	DaftarTotalTSI     = CD + "TotalInterestInsured"
	DaftarTotalEst     = CD + "ListTotalEstimation"
	DaftarInterestDtl  = CD + "InterestListDtl"
	DaftarSpreadAdj    = CD + "SpreadingAdjustment"
	DaftarSpreadAdjQS  = CD + "SpreadingAdjustmentQS"
	AnakSpreadAdj      = "SpreadingAdjustment"
	AnakQuotaShare     = "SpreadingQuotaShare"
	AnakLossAllocation = "LossAllocation"
)

// TabelHeaderKlaim - T_GENERAL_CLAIM: satu baris per klaim, shared PK dengan T_WORK_CLAIM. Kolom kunci (ID) dan SUMBER
// ditulis repository di luar katalog. Kolom Claim Life yang DIPAKAI ULANG untuk PROP: CLAIM_NO, POLICY_NO,
// BUSINESS_CODE, BUSINESS_NAME, KATASTROFE_NOTE. Sisanya kolom khas PROP (migrasi 520).
var TabelHeaderKlaim = Tabel{Nama: "T_GENERAL_CLAIM", Kolom: []Kolom{
	// identitas dan nomor (SaveOutstanding_Act langkah 8-20)
	kKode(CD+"NoClaim", "CLAIM_NO", 64),
	kKode(CD+"ClaimNo", "CLAIM_NO_TEMP", 64),
	// master treaty (SetValueToClaim_Act langkah 2, 5, 7-8)
	kKode(CD+"IDMaster", "MASTER_ID", 64),
	kTeks(CD+"TreatyName", "TREATY_NAME", 500),
	kKode(CD+"TreatyGroupID", "TREATY_GROUP_ID", 64),
	kTeks(CD+"TreatyGroupName", "TREATY_GROUP_NAME", 255),
	kTeks(CD+"ProportionType", "PROPORTION_TYPE", 64),
	kKode(CD+"QuotationData.BusinessCode", "BUSINESS_CODE", 64),
	kTeks(OQ+"BusinessName", "BUSINESS_NAME", 255),
	kTgl(CD+"StartDateTreaty", "START_DATE_TREATY"),
	kTgl(CD+"EndDateTreaty", "END_DATE_TREATY"),
	kKode(CD+"YearofAccount", "YEAR_OF_ACCOUNT", 16),
	kPersen(TM+"RNMShareP", "RNM_SHARE_PCT"),
	// polis (CheckNoPolicy, SetEndDate_Act)
	kKode(CD+"PolicyData.PolicyNo", "POLICY_NO", 64),
	kTgl(CD+"PolicyData.StartDateTime", "POLICY_START_DATE"),
	kTgl(CD+"PolicyData.EndDateTime", "POLICY_END_DATE"),
	kKode(CD+"Quater", "QUARTER", 16),
	kKode(CD+"YearofQuartal", "YEAR_OF_QUARTAL", 16),
	kKode(CD+"TreatyYear", "TREATY_YEAR", 16),
	kTeks(CD+"PolicyNo", "POLICY_NO_CEDING", 255),
	kTeks(CD+"InsuredName", "INSURED_NAME", 500),
	kTeks(CD+"PlaNoCeding", "PLA_NO_CEDING", 255),
	kTeks(CD+"PlaNoSOB", "PLA_NO_SOB", 255),
	kPenanda(CD+"PeriodPolicyTBA", "PERIOD_POLICY_TBA"),
	// registrasi (Section OutstandingClaim / InputAcceptation)
	kTglWaktu(CD+"DateOfLoss", "DATE_OF_LOSS"),
	kTgl(CD+"ReportDate", "REPORT_DATE"),
	kTglWaktu(CD+"DateReceived", "DATE_RECEIVED"),
	kTeks(CD+"ReporterName", "REPORTER_NAME", 255),
	kTeks(CD+"ReporterTelp", "REPORTER_PHONE", 64),
	kKode(CD+"ReportType", "REPORT_TYPE", 16),
	kKode(CD+"ReporterStatus", "REPORTER_STATUS", 16),
	kTeks(CD+"InsuredRelationshipOthers", "INSURED_RELATIONSHIP_OTHERS", 500),
	kTeks(CD+"ReportAddress", "REPORT_ADDRESS", 2000),
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
	kTeks(CD+"ReportDescription", "REPORT_DESCRIPTION", 4000),
	kTeks(CD+"Location", "LOCATION_OF_LOSS", 4000),
	kTeks(CD+"Occupation", "OCCUPATION", 4000),
	kTeks(CD+"Province", "PROVINCE", 255),
	kKode(CD+"ProvinceID", "PROVINCE_ID", 64),
	kKode(CD+"PostalCode", "POSTAL_CODE", 16),
	kTeks(CD+"RW", "RW", 255),
	kKode(CD+"RWID", "RW_ID", 64),
	kTeks(CD+"District", "DISTRICT", 255),
	kKode(CD+"DistrictID", "DISTRICT_ID", 64),
	kTeks(CD+"City", "CITY", 255),
	kKode(CD+"CityID", "CITY_ID", 64),
	// interest dan deductible (Section OutstandingClaim_Intrs / _Est)
	kTeks(CD+"InsuredInterest", "INSURED_INTEREST", 4000),
	kPersen(CD+"ShareCeding", "SHARE_CEDING"),
	kPenanda(CD+"DeductibleType", "DEDUCTIBLE_TYPE"),
	kKode(CD+"FormType", "DEDUCTIBLE_FORM_TYPE", 16),
	kKode(CD+"CurrencyDeductible", "DEDUCTIBLE_CURRENCY_ID", 64),
	kPersen(CD+"Amount", "DEDUCTIBLE_PCT"),
	kKode(CD+"TypeDeductible", "DEDUCTIBLE_BASIS", 16),
	kUang(CD+"TSIDeductible", "TSI_DEDUCTIBLE"),
	kUang(CD+"DeductibleValue", "DEDUCTIBLE_VALUE"),
	kUang(CD+"NetDeductibleValue", "NET_DEDUCTIBLE_VALUE"),
	// penanda alur (pyWorkPage)
	kPenanda("IsOutstanding", "IS_OUTSTANDING"),
	kPenanda("IsAcceptation", "IS_ACCEPTATION"),
	kPenanda("IsCFS", "IS_CFS"),
	kPenanda("ReCFS", "RE_CFS"),
	kPenanda("IsRealisation", "IS_REALISATION"),
	kPenanda("AktifButton", "AKTIF_BUTTON"),
	kPenanda(CD+"IsPLA", "IS_PLA"),
	kKode(CD+"NoPla", "NO_PLA", 64),
	kTeks(CD+"Remark", "REMARK", 4000),
	kTeks(CD+"Remark_Close", "REMARK_CLOSE", 4000),
	// pembayaran (SetPayableTreaty_Act)
	kKode(CD+"Payable", "PAYABLE", 16),
	kTeks(CD+"PayableTo", "PAYABLE_TO", 500),
	kTeks(CD+"DLANoCeding", "DLA_NO_CEDING", 255),
	kTeks(CD+"DLANoSOB", "DLA_NO_SOB", 255),
	kKode(CD+"MarketingData.ID", "MARKETING_ID", 64),
	kKode(CD+"MarketingData.ClientID", "MARKETING_CLIENT_ID", 64),
	kTeks(CD+"MarketingData.ClientName", "MARKETING_CLIENT_NAME", 500),
	kKode(CD+"MarketingData.TeamGroup", "MARKETING_TEAM_GROUP", 64),
	kKode(CD+"MarketingData.BranchDetailID", "MARKETING_BRANCH_ID", 64),
	kTeks(CD+"MarketingData.BranchDetailName", "MARKETING_BRANCH_NAME", 500),
	kTeks(CD+"ReceiverClaim.Name", "RECEIVER_NAME", 500),
	kTeks(CD+"ReceiverClaim.NameOfBank", "RECEIVER_BANK_NAME", 500),
	kTeks(CD+"ReceiverClaim.BranchOfBank", "RECEIVER_BANK_BRANCH", 500),
	kTeks(CD+"ReceiverClaim.NoAccount", "RECEIVER_ACCOUNT_NO", 100),
	kKode(CD+"ReceiverClaim.SwiftCode", "RECEIVER_SWIFT_CODE", 64),
	kKode(CD+"ReceiverClaim.IDOfBank", "RECEIVER_BANK_ID", 64),
	// waktu tahap (DT SetDateOutstanding / SetDateAcceptation)
	kTglWaktu("StartDateEstimation", "START_DATE_ESTIMATION"),
	kTglWaktu("EndDateEstimation", "END_DATE_ESTIMATION"),
	kTglWaktu("StartDateAdjustment", "START_DATE_ADJUSTMENT"),
	// ditulis modul Komite Claim Prop (KomitePostAdjustment) - dibaca layar akseptasi
	kPenanda(CD+"IsCloseFile", "IS_CLOSE_FILE"),
	kPenanda(CD+"IsReservedClaim", "IS_RESERVED_CLAIM"),
	kPenanda("IsAnyAcceptation", "IS_ANY_ACCEPTATION"),
	kPenanda(CD+"IsSubjectivity", "IS_SUBJECTIVITY"),
}}

// TabelEstimasi - T_CLAIM_ESTIMATION <- ClaimData.EstimationList.
var TabelEstimasi = Tabel{Nama: "T_CLAIM_ESTIMATION", Daftar: DaftarEstimasi, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kTgl("EstimationDate", "ESTIMATION_DATE"),
	kKode("Type", "ESTIMATION_TYPE", 16),
	kKode("TypeLossID", "TREATY_TYPE_ID", 64),
	kTeks("TypeLoss", "TREATY_TYPE_NAME", 255),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kUang("KursValue", "KURS"),
	kUang("GrossEstimationPct", "GROSS_ESTIMATION_VALUE"),
	kUang("ConvertGrossEstimasi", "GROSS_ESTIMATION_IDR"),
	kPersen("PersenRNM", "PERSEN_RNM"),
	kUang("EstimationValue", "ESTIMATION_VALUE"),
	kUang("ConvertValue", "ESTIMATION_VALUE_IDR"),
	kKode("NoPLA", "NO_PLA", 64),
	kPenanda("PrintFaceClaim", "IS_PRINT_FACE_CLAIM"),
}}

// TabelInterest - T_CLAIM_INTEREST <- ClaimData.InterestList.
var TabelInterest = Tabel{Nama: "T_CLAIM_INTEREST", Daftar: DaftarInterest, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kTeks("ObjectName", "OBJECT_NAME", 1000),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kUang("KursObjectItem", "KURS"),
	kUang("TSIPerObject", "TSI_VALUE"),
	kUang("TSIPerObjectIDR", "TSI_VALUE_IDR"),
	kPenanda("IsAdjVal", "IS_ADJ_VALUE"),
}}

// TabelClaimAmount - T_CLAIM_CLAIM_AMOUNT <- ClaimData.ListClaimAmount.
var TabelClaimAmount = Tabel{Nama: "T_CLAIM_CLAIM_AMOUNT", Daftar: DaftarClaimAmount, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kUang("IDR", "KURS"),
	kUang("ClaimAmount", "CLAIM_AMOUNT"),
	kUang("NetDeductibleValue", "NET_DEDUCTIBLE_VALUE"),
	kUang("Value", "VALUE"),
	kUang("USD", "VALUE_IDR"),
	kPenanda("Note", "NOTE"),
}}

// TabelLossAlloc - T_CLAIM_LOSS_ALLOCATION <- ClaimData.SpreadingRisk (RALAT STRUKTUR J1 butir 1, migrasi 524).
var TabelLossAlloc = Tabel{Nama: "T_CLAIM_LOSS_ALLOCATION", Daftar: DaftarLossAlloc, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kKode("TreatyType", "TREATY_TYPE_ID", 64),
	kTeks("TreatyName", "TREATY_NAME", 255),
	kPersen("SharePercentage", "SHARE_PERCENTAGE"),
	kUang("ClaimSpreaded", "CLAIM_SPREADED"),
	kUang("ClaimEstimation", "CLAIM_ESTIMATION"),
	kUang("PremiumSpreaded", "KURS"),
	kPenanda("IsOldData", "IS_OLD_DATA"),
}}

// TabelSpreading - T_CLAIM_SPREADING <- ClaimData.SpreadingClaim.
var TabelSpreading = Tabel{Nama: "T_CLAIM_SPREADING", Daftar: DaftarSpreading, KolomInduk: "CLAIM_ID", Kolom: kolomSpreadingKlaim()}

// TabelBreakQS - T_CLAIM_BREAK_QS <- ClaimData.SpreadingBreakQS.
var TabelBreakQS = Tabel{Nama: "T_CLAIM_BREAK_QS", Daftar: DaftarBreakQS, KolomInduk: "CLAIM_ID", Kolom: kolomSpreadingKlaim()}

func kolomSpreadingKlaim() []Kolom {
	return []Kolom{
		kKode("TreatyType", "TREATY_ID", 64),
		kTeks("TreatyName", "TREATY_NAME", 255),
		kPersen("SharePercentage", "SHARE_PERCENTAGE"),
		kUang("ClaimSpreaded", "CLAIM_SPREADED"),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
		kPenanda("IsOldData", "IS_OLD_DATA"),
	}
}

// TabelFacRetro - T_CLAIM_FAC_RETRO <- ClaimData.FacRetroList.
var TabelFacRetro = Tabel{Nama: "T_CLAIM_FAC_RETRO", Daftar: DaftarFacRetro, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kKode("ReinsurerID", "REINSURER_ID", 64),
	kTeks("ReinsurerName", "REINSURER_NAME", 500),
	kPersen("PctShareAllObj", "SHARE_PCT"),
	kPersen("RiCommAllObj", "RI_COMMISSION_PCT"),
	kTeks("AdditionalInfo", "ADDITIONAL_INFO", 1000),
	kUang("TotalEstimasiReas", "TOTAL_ESTIMATION_REINS"),
}}

// TabelAdjustment - T_CLAIM_ADJUSTMENT <- ClaimData.AdjustmentList (+ DataCommitteeTreaty). ID baris STABIL: properti
// `ID` baris membawa kunci baris tersimpan; KOMITE_ID tidak ditulis dari layar (penautan komite = OQ-CP-16).
var TabelAdjustment = Tabel{Nama: "T_CLAIM_ADJUSTMENT", Daftar: DaftarAdjustment, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kKode("Type", "ADJUSTMENT_TYPE", 16),
	kTeks("FormType", "FORM_TYPE", 255),
	kKode("PaymentType", "PAYMENT_TYPE", 16),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kUang("KursIDR", "KURS"),
	kPersen("PersenRNM", "PERSEN_RNM"),
	kTeks("TreatyName", "LOSS_ALLOCATION_NAME", 255),
	kPersen("ShareLossAllocation", "LOSS_ALLOCATION_SHARE"),
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
	kKode("AcceptedNo", "ACCEPTED_NO", 64),
	kTglWaktu("AcceptedDate", "ACCEPTED_DATE"),
	kKode("AcceptanceStatus", "ACCEPTANCE_STATUS", 16),
	kPenanda("IsApproved", "IS_APPROVED"),
	kPenanda("IsKomite", "IS_KOMITE"),
	kPenanda("IsSubjectivity", "IS_SUBJECTIVITY"),
	kTeks("SubjectivityNote", "SUBJECTIVITY_NOTE", 1000),
	kTeks("Notes", "NOTES", 4000),
	kKode("DLA_No", "DLA_NO", 64),
	kTeks("RemarksDLA", "REMARKS_DLA", 4000),
	kPenanda("IsFacRetro", "IS_FAC_RETRO"),
	kPenanda("IsPrintAccept", "IS_PRINT_ACCEPT"),
	kTeks("DataCommitteeTreaty.CircumCauseOfLoss", "KOMITE_CIRCUM_CAUSE_OF_LOSS", 4000),
	kTeks("DataCommitteeTreaty.AdjusterFee", "KOMITE_ADJUSTER_FEE", 4000),
	kTeks("DataCommitteeTreaty.Remarks", "KOMITE_REMARKS", 4000),
	kTeks("DataCommitteeTreaty.Salvage", "KOMITE_SALVAGE", 4000),
	kTeks("DataCommitteeTreaty.LegalLiability", "KOMITE_LEGAL_LIABILITY", 4000),
	kTeks("DataCommitteeTreaty.ExtentOfLoss", "KOMITE_EXTENT_OF_LOSS", 4000),
	kTeks("DataCommitteeTreaty.Occupation", "KOMITE_OCCUPATION", 4000),
	kKode("pxCreateOperator", "CREATED_BY", 64),
	kTeks("pxCreateOpName", "CREATED_BY_NAME", 128),
	kTglWaktu("pxCreateDateTime", "CREATED_AT"),
}}

// TabelAdjSpreading - T_CLAIM_ADJ_SPREADING <- AdjustmentList(n).SpreadingAdjustment.
var TabelAdjSpreading = Tabel{Nama: "T_CLAIM_ADJ_SPREADING", Daftar: AnakSpreadAdj, KolomInduk: "ADJUSTMENT_ID", Kolom: []Kolom{
	kKode("TreatyType", "TREATY_ID", 64),
	kTeks("TreatyName", "TREATY_NAME", 255),
	kPersen("SharePercentage", "SHARE_PERCENTAGE"),
	kUang("ClaimSpreaded", "CLAIM_SPREADED"),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kUang("PremiumSpreaded", "PREMIUM_SPREADED"),
	kTeks("NoAccount", "BANK_ACCOUNT_NO", 100),
	kKode("IDOfBank", "BANK_ID", 64),
}}

// TabelAdjQuotaShare - T_CLAIM_ADJ_QUOTA_SHARE <- AdjustmentList(n).SpreadingQuotaShare.
var TabelAdjQuotaShare = Tabel{Nama: "T_CLAIM_ADJ_QUOTA_SHARE", Daftar: AnakQuotaShare, KolomInduk: "ADJUSTMENT_ID", Kolom: []Kolom{
	kKode("TreatyType", "TREATY_ID", 64),
	kTeks("TreatyName", "TREATY_NAME", 255),
	kPersen("SharePercentage", "SHARE_PERCENTAGE"),
	kUang("ClaimSpreaded", "CLAIM_SPREADED"),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
}}

// TabelAdjLossAlloc - T_CLAIM_ADJ_LOSS_ALLOCATION <- AdjustmentList(n).LossAllocation.
var TabelAdjLossAlloc = Tabel{Nama: "T_CLAIM_ADJ_LOSS_ALLOCATION", Daftar: AnakLossAllocation, KolomInduk: "ADJUSTMENT_ID", Kolom: []Kolom{
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kKode("TreatyType", "TREATY_TYPE_ID", 64),
	kTeks("TreatyName", "TREATY_NAME", 255),
	kPersen("SharePercentage", "SHARE_PERCENTAGE"),
	kUang("ClaimSpreaded", "CLAIM_SPREADED"),
	kUang("ClaimEstimation", "CLAIM_ESTIMATION"),
	kUang("PremiumSpreaded", "KURS"),
}}

// TabelAnakKlaim - anak langsung T_GENERAL_CLAIM, urutan tulis.
var TabelAnakKlaim = []Tabel{TabelEstimasi, TabelInterest, TabelClaimAmount, TabelLossAlloc, TabelSpreading,
	TabelBreakQS, TabelFacRetro}

// TabelCucuAdjustment - anak T_CLAIM_ADJUSTMENT.
var TabelCucuAdjustment = []Tabel{TabelAdjSpreading, TabelAdjQuotaShare, TabelAdjLossAlloc}

// SemuaTabel - seluruh tabel katalog (header, anak, adjustment, cucu) - untuk uji katalog.
func SemuaTabel() []Tabel {
	out := []Tabel{TabelHeaderKlaim}
	out = append(out, TabelAnakKlaim...)
	out = append(out, TabelAdjustment)
	out = append(out, TabelCucuAdjustment...)
	return out
}

// PropID - properti tersembunyi baris adjustment yang membawa ID baris tersimpan (T_CLAIM_ADJUSTMENT.ID). Tidak
// pernah ditulis layar; dipakai repository untuk menyunting di tempat (AC 63).
const PropID = "ID"

// PropKomiteID - properti baris adjustment yang membawa KOMITE_ID tersimpan (baca saja dari layar).
const PropKomiteID = "KomiteID"

// ProyeksiKatalog - halaman sebagaimana ia KELUAR dari penyimpanan: hanya medan berkolom. Dipakai gudang tiruan supaya
// uji seam HTTP melihat apa yang Oracle simpan.
func ProyeksiKatalog(h *Halaman) *Halaman {
	s := HalamanBaru()
	for _, k := range TabelHeaderKlaim.Kolom {
		if v := h.Ambil(k.Properti); v != "" {
			s.Setel(k.Properti, v)
		}
	}
	saring := func(t Tabel, b []Baris, simpanID bool) []Baris {
		var out []Baris
		for _, x := range b {
			nb := Baris{}
			for _, k := range t.Kolom {
				if v := x[k.Properti]; v != "" {
					nb[k.Properti] = v
				}
			}
			if simpanID {
				for _, p := range []string{PropID, PropKomiteID} {
					if v := x[p]; v != "" {
						nb[p] = v
					}
				}
			}
			out = append(out, nb)
		}
		return out
	}
	for _, t := range TabelAnakKlaim {
		s.SetelDaftar(t.Daftar, saring(t, h.AmbilDaftar(t.Daftar), false))
	}
	adj := h.AmbilDaftar(DaftarAdjustment)
	s.SetelDaftar(DaftarAdjustment, saring(TabelAdjustment, adj, true))
	for i := range adj {
		for _, c := range TabelCucuAdjustment {
			j := JalurAnak(DaftarAdjustment, i+1, c.Daftar)
			s.SetelDaftar(j, saring(c, h.AmbilDaftar(j), false))
		}
	}
	return s
}
