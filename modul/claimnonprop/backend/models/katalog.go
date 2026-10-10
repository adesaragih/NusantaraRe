package models

// Untuk apa berkas ini: KATALOG KOLOM - satu-satunya pemetaan properti Pega -> kolom Oracle -> golongan tipe untuk
// penyimpanan relasional Claim Non Prop (OQ-CNP-06 "pola Claim Prop + tabel XoL", keputusan work owner 09-10-2026;
// izin menyunting STRUKTUR Claim Life / Claim Prop untuk kolom tabel bersama, keputusan work owner 09-10-2026).
//
// Susunan (docs/STRUKTUR-TABEL-CLAIM-NON-PROP.md):
//   - T_GENERAL_CLAIM (kepala, tabel bersama Claim Life) - kolom yang sudah ada dipakai ulang bila maknanya sama;
//     kolom khas Non Prop ditambah migrasi 600 (ALTER ADD nullable).
//   - Tabel anak Claim Prop dipakai ulang bila kolomnya cocok, kekurangannya ditambah ALTER ADD nullable (601-605):
//     T_CLAIM_INTEREST, T_CLAIM_CLAIM_AMOUNT, T_CLAIM_SPREADING, T_CLAIM_BREAK_QS, T_CLAIM_ADJUSTMENT,
//     T_CLAIM_ADJ_SPREADING, T_CLAIM_ADJ_QUOTA_SHARE.
//   - Tabel baru khas XoL (606-608): T_CLAIM_NP_LOSS_ALLOC (Loss Allocation + To XOL, induk klaim ATAU akseptasi),
//     T_CLAIM_NP_XOL_ALLOC (XOL Allocation per layer; JENIS membedakan alokasi berjalan, alokasi lama `LossAllocation`
//     akseptasi, dan Previously Calculated `AlokasiXOLPaid`), T_CLAIM_NP_CLAIM_ACCEPT (Claim Acceptation per mata uang).
//
// ⛔ Medan TURUNAN tidak punya kolom - dihitung ulang dari basis asli (`HitungTurunan`): TotalInterestInsured,
// TotalSumInsuredIDR, TotalListClaimAmount(IDR), ListTotalEstimation (Summary XOL), SpreadingAdjustment(QS) tingkat
// klaim, AdjustmentList(n).CurencyAdjustment / .ComiteeClaim / .TotalKomite / .CommentLOD / .FlagCurrency.

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

// Desimal - golongan bertipe NUMBER(38,10).
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
	// Baru - kolom lahir di migrasi modul ini (ALTER ADD / CREATE); false = kolom tabel bersama yang dipakai ulang.
	Baru bool
}

// Tabel - satu daftar katalog di satu tabel.
type Tabel struct {
	Nama string
	// Daftar - jalur PageList sumber baris (tabel anak klaim) atau nama daftar di dalam baris
	// `ClaimData.AdjustmentList(n)` (tabel cucu akseptasi).
	Daftar string
	// KolomInduk - kolom penunjuk induk (CLAIM_ID / ADJUSTMENT_ID).
	KolomInduk string
	// Jenis - nilai kolom JENIS (tabel T_CLAIM_NP_XOL_ALLOC); "" = tabel tanpa pembeda.
	Jenis string
	Kolom []Kolom
}

func kTeks(p, k string, n int) Kolom { return Kolom{p, k, GolTeks, n, false} }
func kKode(p, k string, n int) Kolom { return Kolom{p, k, GolKode, n, false} }
func kPenanda(p, k string) Kolom     { return Kolom{p, k, GolPenanda, 16, false} }
func kUang(p, k string) Kolom        { return Kolom{p, k, GolUang, 0, false} }
func kPersen(p, k string) Kolom      { return Kolom{p, k, GolPersen, 0, false} }
func kTgl(p, k string) Kolom         { return Kolom{p, k, GolTanggal, 0, false} }
func kTglWaktu(p, k string) Kolom    { return Kolom{p, k, GolTanggalWaktu, 0, false} }

// baru - kolom yang lahir di migrasi modul ini.
func baru(k Kolom) Kolom { k.Baru = true; return k }

// Awalan jalur halaman.
const (
	CD = "ClaimData."
	TM = "TreatyInMaster."
	// OQ - `pyWorkPage.OfferFacIn.QuotationData` (Class of Business klaim, diisi SetValueClaimTNP_Act langkah 2).
	OQ = "OfferFacIn.QuotationData."
	// RC - `pyWorkPage.ClaimData.ReceiverClaim(1)`: XML Non Prop selalu memakai indeks 1 (AdjustmentDetailNP); di sini
	// halaman tunggal supaya jalurnya bukan jalur baris grid.
	RC = CD + "ReceiverClaim."
)

// Jalur daftar tingkat klaim.
const (
	DaftarInterest    = CD + "InterestList"
	DaftarClaimAmount = CD + "ListClaimAmount"
	DaftarLossAlloc   = CD + "CNPSpreadLoss"
	DaftarXOL         = CD + "SpreadingRisk"
	DaftarSpreading   = CD + "SpreadingClaim"
	DaftarBreakQS     = CD + "SpreadingBreakQS"
	DaftarAdjustment  = CD + "AdjustmentList"
	DaftarRiwayat     = CD + "SuggestList"
	DaftarTotalTSI    = CD + "TotalInterestInsured"
	DaftarSummaryXOL  = CD + "ListTotalEstimation"
	DaftarSpreadAdj   = CD + "SpreadingAdjustment"
	DaftarSpreadAdjQS = CD + "SpreadingAdjustmentQS"
	// Daftar di dalam baris AdjustmentList(n).
	AnakClaimAccept  = "ListClaimAcceptation"
	AnakLossAlloc    = "CNPSpreadLoss"
	AnakXOL          = "SpreadingRisk"
	AnakXOLLama      = "LossAllocation"
	AnakXOLDibayar   = "AlokasiXOLPaid"
	AnakSpreadIn     = "SpreadingAdjustment"
	AnakSpreadOut    = "SpreadingQuotaShare"
	AnakKomite       = "ComiteeClaim"
	AnakMataUangAdj  = "CurencyAdjustment"
	JenisXOLBerjalan = "ALOKASI"
	JenisXOLLama     = "LAMA"
	JenisXOLDibayar  = "DIBAYAR"
)

// TabelHeaderKlaim - T_GENERAL_CLAIM: satu baris per klaim, shared PK dengan T_WORK_CLAIM. Kolom kunci (ID) dan SUMBER
// ditulis repository di luar katalog.
var TabelHeaderKlaim = Tabel{Nama: "T_GENERAL_CLAIM", Kolom: []Kolom{
	// nomor (SaveDataToOSAksep_Act langkah 15.6.8)
	kKode(CD+"NoClaim", "CLAIM_NO", 64),
	baru(kKode(CD+"QuotationData.BusinessOldId", "BUSINESS_OLD_ID", 64)),
	// master treaty (SetValueClaimTNP_Act langkah 2, 5)
	kKode(CD+"IDMaster", "MASTER_ID", 64),
	kTeks(CD+"TreatyName", "TREATY_NAME", 500),
	kTeks(TM+"TreatyGroup", "TREATY_GROUP_NAME", 255),
	kKode(OQ+"BusinessCode", "BUSINESS_CODE", 64),
	kTeks(OQ+"BusinessName", "BUSINESS_NAME", 255),
	kTgl(CD+"StartDateTreaty", "START_DATE_TREATY"),
	kTgl(CD+"EndDateTreaty", "END_DATE_TREATY"),
	kKode(CD+"YearofAccount", "YEAR_OF_ACCOUNT", 16),
	kPersen(TM+"RNMShare", "RNM_SHARE_PCT"),
	baru(kKode(CD+"IDMasterTONP", "MASTER_ID_TO", 64)),
	// polis (CheckNoPolicy, SetEndDate_Act)
	kKode(CD+"PolicyData.PolicyNo", "POLICY_NO", 64),
	baru(kTeks(CD+"PolicyData.TreatyGroup", "POLICY_COB", 255)),
	kTgl(CD+"PolicyData.StartDateTime", "POLICY_START_DATE"),
	kTgl(CD+"PolicyData.EndDateTime", "POLICY_END_DATE"),
	kTeks(CD+"PolicyNo", "POLICY_NO_CEDING", 255),
	kTeks(CD+"InsuredName", "INSURED_NAME", 500),
	baru(kTeks(CD+"CNPReinsuranceSlip", "REINSURANCE_SLIP", 255)),
	baru(kTeks(CD+"CNPClmNoCedant", "CLAIM_NO_CEDING", 255)),
	kTeks(CD+"PlaNoCeding", "PLA_NO_CEDING", 255),
	kTeks(CD+"DLANoCeding", "DLA_NO_CEDING", 255),
	kTeks(CD+"PlaNoSOB", "PLA_NO_SOB", 255),
	// registrasi (Section OutstandingClaim / InputAcceptation)
	kTgl(CD+"DateOfLoss", "DATE_OF_LOSS"),
	kTgl(CD+"ReportDate", "REPORT_DATE"),
	kTgl(CD+"DateReceived", "DATE_RECEIVED"),
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
	baru(kTeks(CD+"CNPCircumtances", "CIRCUMSTANCES", 4000)),
	baru(kTeks(CD+"CNPSupportDoc", "SUPPORTING_DOCUMENT", 4000)),
	kTeks(CD+"Province", "PROVINCE", 255),
	kKode(CD+"ProvinceID", "PROVINCE_ID", 64),
	kKode(CD+"PostalCode", "POSTAL_CODE", 16),
	kTeks(CD+"RW", "RW", 255),
	kKode(CD+"RWID", "RW_ID", 64),
	kTeks(CD+"District", "DISTRICT", 255),
	kKode(CD+"DistrictID", "DISTRICT_ID", 64),
	kTeks(CD+"City", "CITY", 255),
	kKode(CD+"CityID", "CITY_ID", 64),
	// interest, deductible, TPL (Section OutstandingClaim tab Interests / Estimation)
	kTeks(CD+"InsuredInterest", "INSURED_INTEREST", 4000),
	kPersen(CD+"ShareCeding", "SHARE_CEDING"),
	kPenanda(CD+"DeductibleType", "DEDUCTIBLE_TYPE"),
	kKode(CD+"FormType", "DEDUCTIBLE_FORM_TYPE", 16),
	kKode(CD+"CurrencyDeductible", "DEDUCTIBLE_CURRENCY_ID", 64),
	kPersen(CD+"Amount", "DEDUCTIBLE_PCT"),
	kKode(CD+"TypeDeductible", "DEDUCTIBLE_BASIS", 16),
	kUang(CD+"TSIDeductible", "TSI_DEDUCTIBLE"),
	kUang(CD+"DeductibleValue", "DEDUCTIBLE_VALUE"),
	baru(kKode(CD+"CNPDeducMinMax", "DEDUCTIBLE_MIN_MAX", 16)),
	baru(kPenanda(CD+"IsTPL", "IS_TPL")),
	baru(kKode(CD+"TPLFormat", "TPL_FORMAT", 16)),
	baru(kKode(CD+"TPLType", "TPL_TYPE", 16)),
	baru(kPersen(CD+"PctTPL", "TPL_PCT")),
	// penanda alur (pyWorkPage)
	kPenanda("IsOutstanding", "IS_OUTSTANDING"),
	kPenanda("IsAcceptation", "IS_ACCEPTATION"),
	kPenanda("IsCFS", "IS_CFS"),
	kPenanda("IsRealisation", "IS_REALISATION"),
	kPenanda("AktifButton", "AKTIF_BUTTON"),
	kPenanda(CD+"IsPLA", "IS_PLA"),
	kKode(CD+"NoPla", "NO_PLA", 64),
	kTeks(CD+"Remark_Close", "REMARK_CLOSE", 4000),
	kPenanda("IsCloseFile", "IS_CLOSE_FILE"),
	baru(kPenanda("IsSaveToOs", "IS_SAVE_TO_OS")),
	baru(kPenanda("FlagActualPremium", "WAITING_ACTUAL_PREMIUM")),
	baru(kTeks("CNPStatusCase", "STATUS_CASE", 64)),
	baru(kPenanda("IsReject", "IS_REJECT")),
	baru(kKode("KomiteNo", "KOMITE_NO", 64)),
	// FlagPrintPla.CARI28 (halaman requestor Pega: Save to issue RNM / Save To OS bernilai berubah -> Print PLA
	// menerbitkan revisi nomor PLA, GeneratePlaCNP_Act 7) - disimpan supaya bertahan antar-permintaan.
	baru(kPenanda("FlagPrintPla.CARI28", "FLAG_PRINT_PLA")),
	// pembayaran (SetPayableTreatyNP_Act, SetAccoutNo_Act)
	kKode(CD+"Payable", "PAYABLE", 16),
	kTeks(CD+"PayableTo", "PAYABLE_TO", 500),
	kKode(CD+"MarketingData.ID", "MARKETING_ID", 64),
	kKode(CD+"MarketingData.ClientID", "MARKETING_CLIENT_ID", 64),
	kTeks(CD+"MarketingData.ClientName", "MARKETING_CLIENT_NAME", 500),
	kKode(CD+"MarketingData.TeamGroup", "MARKETING_TEAM_GROUP", 64),
	kKode(CD+"MarketingData.BranchDetailID", "MARKETING_BRANCH_ID", 64),
	kTeks(CD+"MarketingData.BranchDetailName", "MARKETING_BRANCH_NAME", 500),
	kTeks(RC+"Name", "RECEIVER_NAME", 500),
	kTeks(RC+"NameOfBank", "RECEIVER_BANK_NAME", 500),
	kTeks(RC+"BranchOfBank", "RECEIVER_BANK_BRANCH", 500),
	kTeks(RC+"NoAccount", "RECEIVER_ACCOUNT_NO", 100),
	kKode(RC+"SwiftCode", "RECEIVER_SWIFT_CODE", 64),
	kKode(RC+"IDOfBank", "RECEIVER_BANK_ID", 64),
	baru(kKode(RC+"Currency", "RECEIVER_CURRENCY", 64)),
	baru(kTeks(RC+"NameOfBank2", "RECEIVER_BANK_NAME2", 500)),
	baru(kTeks(RC+"BranchOfBank2", "RECEIVER_BANK_BRANCH2", 500)),
	baru(kTeks(RC+"NoAccount2", "RECEIVER_ACCOUNT_NO2", 100)),
	baru(kKode(RC+"SwiftCode2", "RECEIVER_SWIFT_CODE2", 64)),
	baru(kKode(RC+"IDOfBank2", "RECEIVER_BANK_ID2", 64)),
	baru(kKode(RC+"Currency2", "RECEIVER_CURRENCY2", 64)),
}}

// TabelInterest - T_CLAIM_INTEREST <- ClaimData.InterestList (+ TPL, Section InputDtlInterest / SetTPLNote_Act).
var TabelInterest = Tabel{Nama: "T_CLAIM_INTEREST", Daftar: DaftarInterest, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kTeks("ObjectName", "OBJECT_NAME", 1000),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	kUang("KursObjectItem", "KURS"),
	kUang("TSIPerObject", "TSI_VALUE"),
	baru(kPenanda("IsTPL", "IS_TPL")),
	baru(kKode("TPLFormat", "TPL_FORMAT", 16)),
	baru(kKode("TPLType", "TPL_TYPE", 16)),
	baru(kPersen("TPLPct", "TPL_PCT")),
	baru(kUang("TPLAmount", "TPL_AMOUNT")),
	baru(kKode("CNPMinMax", "TPL_MIN_MAX", 16)),
	baru(kKode("TPLCurrency", "TPL_CURRENCY", 64)),
	baru(kTeks("TPLNote", "TPL_NOTE", 255)),
	baru(kKode("TPLType2", "TPL_TYPE2", 16)),
	baru(kPersen("TPLPct2", "TPL_PCT2")),
	baru(kUang("TPLAmount2", "TPL_AMOUNT2")),
	baru(kKode("CNPMinMax2", "TPL_MIN_MAX2", 16)),
	baru(kKode("CurrencyTPL2", "TPL_CURRENCY2", 64)),
}}

// kolomJumlahKlaim - baris "Claim Amount" per mata uang (ListClaimAmount; ListClaimAcceptation akseptasi berkelas sama
// `ASM-FW-GISFW-Data-TreatyInTotal`, AddListClaimNP_Act). `tambahan` = kolom yang lahir di migrasi modul ini pada tabel
// bersama (T_CLAIM_CLAIM_AMOUNT); pada tabel baru semua kolom baru.
func kolomJumlahKlaim(tabelBaru bool) []Kolom {
	lama := func(k Kolom) Kolom {
		if tabelBaru {
			return baru(k)
		}
		return k
	}
	return []Kolom{
		lama(kKode("CurrencyID", "CURRENCY_ID", 64)),
		lama(kTeks("Currency", "CURRENCY_NAME", 64)),
		lama(kUang("AltValue", "KURS")),
		lama(kUang("Value", "VALUE")),
		lama(kUang("USD", "VALUE_IDR")),
		lama(kUang("CNPDeductible", "NET_DEDUCTIBLE_VALUE")),
		lama(kPenanda("Note", "NOTE")),
		baru(kUang("TPL", "TPL")),
		baru(kUang("AdjusterFee", "ADJUSTER_FEE")),
		baru(kUang("Salvage", "SALVAGE")),
		baru(kUang("CNPOthersFee", "OTHERS_FEE")),
		baru(kPersen("PctProrateClaim", "PROPORTION_PCT")),
		baru(kUang("ClaimAmountCedant", "CLAIM_AMOUNT_CEDANT")),
		baru(kUang("ClaimAmountAdjust", "CLAIM_AMOUNT_ADJUST")),
		baru(kUang("CNPTSI", "TSI_VALUE")),
		baru(kPenanda("CNPFlagOuts", "IS_LOCKED")),
	}
}

// TabelClaimAmount - T_CLAIM_CLAIM_AMOUNT <- ClaimData.ListClaimAmount.
var TabelClaimAmount = Tabel{Nama: "T_CLAIM_CLAIM_AMOUNT", Daftar: DaftarClaimAmount, KolomInduk: "CLAIM_ID",
	Kolom: kolomJumlahKlaim(false)}

// kolomLossAlloc - baris "Loss Allocation" (`CNPSpreadLoss`, kelas ASM-FW-GISFW-Data-SpreadingRisk).
func kolomLossAlloc() []Kolom {
	return []Kolom{
		baru(kKode("CurrencyID", "CURRENCY_ID", 64)),
		baru(kTeks("Currency", "CURRENCY_NAME", 64)),
		baru(kTeks("TreatyName", "TREATY_NAME", 255)),
		baru(kPersen("ClaimPercentage", "CLAIM_PERCENTAGE")),
		baru(kUang("ClaimAmountAdjust", "CLAIM_AMOUNT_ADJUST")),
		baru(kUang("AdjusterFee", "ADJUSTER_FEE")),
		baru(kUang("Salvage", "SALVAGE")),
		baru(kUang("CNPOthersFee", "OTHERS_FEE")),
		baru(kPenanda("CNPFlagXOL", "TO_XOL")),
		baru(kPenanda("CNPFlagOuts", "IS_LOCKED")),
	}
}

// TabelLossAlloc - T_CLAIM_NP_LOSS_ALLOC <- ClaimData.CNPSpreadLoss.
var TabelLossAlloc = Tabel{Nama: "T_CLAIM_NP_LOSS_ALLOC", Daftar: DaftarLossAlloc, KolomInduk: "CLAIM_ID",
	Kolom: kolomLossAlloc()}

// kolomXOL - baris "XOL Allocation" (`SpreadingRisk`, CountLossAllocation_act langkah 17.2.5 / 17.2.14; akseptasi
// AdjClaimCNP_Act langkah 2-11).
func kolomXOL() []Kolom {
	return []Kolom{
		baru(kKode("TreatyType", "TREATY_TYPE_ID", 64)),
		baru(kTeks("TreatyName", "TREATY_NAME", 255)),
		baru(kKode("CurrencyID", "CURRENCY_ID", 64)),
		baru(kTeks("Currency", "CURRENCY_NAME", 64)),
		baru(kUang("Kurs", "KURS")),
		baru(kUang("KursIDR", "KURS_IDR")),
		baru(kUang("ClaimEstimation", "CLAIM_ESTIMATION")),
		baru(kUang("ClaimAmountAdjust", "CLAIM_AMOUNT_ADJUST")),
		baru(kUang("AdjClaimValue", "ADJ_CLAIM_VALUE")),
		baru(kUang("TotalClaim", "TOTAL_CLAIM")),
		baru(kPersen("ClaimPercentage", "CLAIM_PERCENTAGE")),
		baru(kUang("ClaimSpreaded", "CLAIM_SPREADED")),
		baru(kUang("AdjusterFee", "ADJUSTER_FEE")),
		baru(kUang("Salvage", "SALVAGE")),
		baru(kUang("CNPOthersFee", "OTHERS_FEE")),
		baru(kUang("TotalSpread", "TOTAL_SPREAD")),
		baru(kUang("ClaimAmountIDR", "CLAIM_AMOUNT_IDR")),
		baru(kPersen("CNPProrateClaim", "PRORATE_PCT")),
		baru(kUang("CNPLimit", "LIMIT_VALUE")),
		baru(kUang("CNPLimitFull", "LIMIT_FULL")),
		baru(kKode("CurrLayerOri", "LAYER_CURRENCY", 64)),
		baru(kUang("CNPMDP", "MDP_VALUE")),
		baru(kPersen("CNPPctReinstate", "REINSTATE_PCT")),
		baru(kKode("Layer", "LAYER", 64)),
		baru(kKode("LayerType", "LAYER_TYPE", 64)),
		baru(kKode("LayerPart", "LAYER_PART", 64)),
		baru(kKode("LayerPartType", "LAYER_PART_TYPE", 64)),
		baru(kPenanda("IsEditClaim", "IS_EDIT_CLAIM")),
		baru(kPenanda("CNPFlagOuts", "IS_LOCKED")),
		baru(kUang("TotalClaimRNM", "TOTAL_CLAIM_RNM")),
		baru(kUang("CNPReinstatement", "REINSTATEMENT")),
		baru(kUang("CNPReinstatementRNM", "REINSTATEMENT_RNM")),
	}
}

// TabelXOL - T_CLAIM_NP_XOL_ALLOC (JENIS ALOKASI) <- ClaimData.SpreadingRisk.
var TabelXOL = Tabel{Nama: "T_CLAIM_NP_XOL_ALLOC", Daftar: DaftarXOL, KolomInduk: "CLAIM_ID", Jenis: JenisXOLBerjalan,
	Kolom: kolomXOL()}

// kolomSpreading - baris Spreading List / Break QS (`SpreadingClaim`, `SpreadingBreakQS`; T_CLAIM_SPREADING /
// T_CLAIM_BREAK_QS Claim Prop + kolom fee Non Prop).
func kolomSpreading() []Kolom {
	return []Kolom{
		kKode("TreatyType", "TREATY_ID", 64),
		kTeks("TreatyName", "TREATY_NAME", 255),
		kPersen("SharePercentage", "SHARE_PERCENTAGE"),
		kUang("ClaimSpreaded", "CLAIM_SPREADED"),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
		baru(kUang("AdjusterFee", "ADJUSTER_FEE")),
		baru(kUang("Salvage", "SALVAGE")),
		baru(kUang("CNPOthersFee", "OTHERS_FEE")),
		baru(kUang("ClaimAmountIDR", "CLAIM_AMOUNT_IDR")),
		baru(kPenanda("CNPFlagOuts", "IS_LOCKED")),
	}
}

// TabelSpreading - T_CLAIM_SPREADING <- ClaimData.SpreadingClaim.
var TabelSpreading = Tabel{Nama: "T_CLAIM_SPREADING", Daftar: DaftarSpreading, KolomInduk: "CLAIM_ID",
	Kolom: kolomSpreading()}

// TabelBreakQS - T_CLAIM_BREAK_QS <- ClaimData.SpreadingBreakQS.
var TabelBreakQS = Tabel{Nama: "T_CLAIM_BREAK_QS", Daftar: DaftarBreakQS, KolomInduk: "CLAIM_ID", Kolom: kolomSpreading()}

// TabelAdjustment - T_CLAIM_ADJUSTMENT <- ClaimData.AdjustmentList (AddAkseptasiCNP_Act, AdjustmentDetailNP). ID baris
// STABIL: properti `ID` membawa kunci baris tersimpan; KOMITE_ID ditulis penyerahan komite, tidak dari layar.
var TabelAdjustment = Tabel{Nama: "T_CLAIM_ADJUSTMENT", Daftar: DaftarAdjustment, KolomInduk: "CLAIM_ID", Kolom: []Kolom{
	kKode("Type", "ADJUSTMENT_TYPE", 16),
	kKode("PaymentType", "PAYMENT_TYPE", 16),
	baru(kKode("CNPIndexInterim", "INTERIM_INDEX", 16)),
	kPersen("PersenRNM", "PERSEN_RNM"),
	kUang("TotalEstimasiValue", "TOTAL_ESTIMATION_VALUE"),
	kUang("ValueAdjustment", "VALUE_ADJUSTMENT_IDR"),
	kKode("IndividualRiskType", "INDIVIDUAL_RISK_TYPE", 16),
	kPersen("IndividualRiskPercentage", "INDIVIDUAL_RISK_PCT"),
	kUang("IndividualRiskValue", "INDIVIDUAL_RISK_VALUE"),
	kKode("Payable", "PAYABLE", 16),
	kTeks("PayableTo", "PAYABLE_TO", 500),
	kTeks("NameOfBank", "BANK_NAME", 500),
	kTeks("BranchOfBank", "BANK_BRANCH", 500),
	kTeks("NoAccount", "BANK_ACCOUNT_NO", 100),
	kKode("SwiftCode", "SWIFT_CODE", 64),
	kKode("IDOfBank", "BANK_ID", 64),
	kTeks("Currency", "CURRENCY_NAME", 64),
	baru(kTeks("NameOfBank2", "BANK_NAME2", 500)),
	baru(kTeks("BranchOfBank2", "BANK_BRANCH2", 500)),
	baru(kTeks("NoAccount2", "BANK_ACCOUNT_NO2", 100)),
	baru(kKode("SwiftCode2", "SWIFT_CODE2", 64)),
	baru(kKode("IDOfBank2", "BANK_ID2", 64)),
	baru(kTeks("Currency2", "BANK_CURRENCY2", 64)),
	kTeks("DLANoCeding", "DLA_NO_CEDING", 255),
	kPenanda("DirectToKasir", "IS_DIRECT_TO_KASIR"),
	kTeks("StatusKasir", "STATUS_KASIR", 1000),
	baru(kPenanda("FlagErrorKasir", "FLAG_ERROR_KASIR")),
	kKode("AcceptedNo", "ACCEPTED_NO", 64),
	kTglWaktu("AcceptedDate", "ACCEPTED_DATE"),
	kKode("AcceptanceStatus", "ACCEPTANCE_STATUS", 16),
	kPenanda("IsKomite", "IS_KOMITE"),
	kPenanda("IsSubjectivity", "IS_SUBJECTIVITY"),
	kTeks("SubjectivityNote", "SUBJECTIVITY_NOTE", 1000),
	kKode("DLA_No", "DLA_NO", 64),
	kTeks("RemarksDLA", "REMARKS_DLA", 4000),
	kPenanda("IsPrintAccept", "IS_PRINT_ACCEPT"),
	kTeks("DataCommitteeTreaty.CircumCauseOfLoss", "KOMITE_CIRCUM_CAUSE_OF_LOSS", 4000),
	kTeks("DataCommitteeTreaty.Remarks", "KOMITE_REMARKS", 4000),
	kTeks("DataCommitteeTreaty.Occupation", "KOMITE_OCCUPATION", 4000),
	kKode("pxCreateOperator", "CREATED_BY", 64),
	kTeks("pxCreateOpName", "CREATED_BY_NAME", 128),
	kTglWaktu("pxCreateDateTime", "CREATED_AT"),
}}

// kolomSpreadAkseptasi - Spreading In / Out akseptasi (`SpreadingAdjustment` / `SpreadingQuotaShare`, AdjClaimCNP_Act
// langkah 15): kolom Claim Prop + fee, Total Claim, RNM Net Claim. `spreadIn` - T_CLAIM_ADJ_SPREADING sudah punya
// PREMIUM_SPREADED dan rekening (BANK_ACCOUNT_NO / BANK_ID, SetAccoutNo_Act); T_CLAIM_ADJ_QUOTA_SHARE belum punya
// PREMIUM_SPREADED.
func kolomSpreadAkseptasi(spreadIn bool) []Kolom {
	premi := kUang("PremiumSpreaded", "PREMIUM_SPREADED")
	if !spreadIn {
		premi = baru(premi)
	}
	var rekening []Kolom
	if spreadIn {
		rekening = []Kolom{kTeks("NoAccount", "BANK_ACCOUNT_NO", 100), kKode("IDOfBank", "BANK_ID", 64)}
	}
	return append([]Kolom{
		kKode("TreatyType", "TREATY_ID", 64),
		kTeks("TreatyName", "TREATY_NAME", 255),
		kPersen("SharePercentage", "SHARE_PERCENTAGE"),
		kUang("ClaimSpreaded", "CLAIM_SPREADED"),
		kKode("CurrencyID", "CURRENCY_ID", 64),
		kTeks("Currency", "CURRENCY_NAME", 64),
		premi,
		baru(kUang("AdjusterFee", "ADJUSTER_FEE")),
		baru(kUang("Salvage", "SALVAGE")),
		baru(kUang("CNPOthersFee", "OTHERS_FEE")),
		baru(kUang("TotalClaim", "TOTAL_CLAIM")),
		baru(kUang("NetClaim", "NET_CLAIM")),
	}, rekening...)
}

// Tabel cucu akseptasi.
var (
	TabelAdjSpreading  = Tabel{Nama: "T_CLAIM_ADJ_SPREADING", Daftar: AnakSpreadIn, KolomInduk: "ADJUSTMENT_ID", Kolom: kolomSpreadAkseptasi(true)}
	TabelAdjQuotaShare = Tabel{Nama: "T_CLAIM_ADJ_QUOTA_SHARE", Daftar: AnakSpreadOut, KolomInduk: "ADJUSTMENT_ID", Kolom: kolomSpreadAkseptasi(false)}
	TabelAdjClaimAcc   = Tabel{Nama: "T_CLAIM_NP_CLAIM_ACCEPT", Daftar: AnakClaimAccept, KolomInduk: "ADJUSTMENT_ID", Kolom: kolomJumlahKlaim(true)}
	TabelAdjLossAlloc  = Tabel{Nama: "T_CLAIM_NP_LOSS_ALLOC", Daftar: AnakLossAlloc, KolomInduk: "ADJUSTMENT_ID", Kolom: kolomLossAlloc()}
	TabelAdjXOL        = Tabel{Nama: "T_CLAIM_NP_XOL_ALLOC", Daftar: AnakXOL, KolomInduk: "ADJUSTMENT_ID", Jenis: JenisXOLBerjalan, Kolom: kolomXOL()}
	TabelAdjXOLLama    = Tabel{Nama: "T_CLAIM_NP_XOL_ALLOC", Daftar: AnakXOLLama, KolomInduk: "ADJUSTMENT_ID", Jenis: JenisXOLLama, Kolom: kolomXOL()}
	TabelAdjXOLDibayar = Tabel{Nama: "T_CLAIM_NP_XOL_ALLOC", Daftar: AnakXOLDibayar, KolomInduk: "ADJUSTMENT_ID", Jenis: JenisXOLDibayar, Kolom: kolomXOL()}
)

// TabelAnakKlaim - anak langsung T_GENERAL_CLAIM, urutan tulis.
var TabelAnakKlaim = []Tabel{TabelInterest, TabelClaimAmount, TabelLossAlloc, TabelXOL, TabelSpreading, TabelBreakQS}

// TabelCucuAdjustment - anak T_CLAIM_ADJUSTMENT.
var TabelCucuAdjustment = []Tabel{TabelAdjClaimAcc, TabelAdjLossAlloc, TabelAdjXOL, TabelAdjXOLLama, TabelAdjXOLDibayar,
	TabelAdjSpreading, TabelAdjQuotaShare}

// SemuaTabel - seluruh daftar katalog (header, anak, adjustment, cucu).
func SemuaTabel() []Tabel {
	out := []Tabel{TabelHeaderKlaim}
	out = append(out, TabelAnakKlaim...)
	out = append(out, TabelAdjustment)
	out = append(out, TabelCucuAdjustment...)
	return out
}

// PropID - properti tersembunyi baris akseptasi yang membawa ID baris tersimpan (T_CLAIM_ADJUSTMENT.ID).
const PropID = "ID"

// PropKomiteID - properti baris akseptasi yang membawa KOMITE_ID tersimpan (baca saja dari layar).
const PropKomiteID = "KomiteID"

// JalurAdj - jalur daftar di dalam baris AdjustmentList(n).
func JalurAdj(n int, anak string) string { return JalurAnak(DaftarAdjustment, n, anak) }

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
			j := JalurAdj(i+1, c.Daftar)
			s.SetelDaftar(j, saring(c, h.AmbilDaftar(j), false))
		}
	}
	riw := h.AmbilDaftar(DaftarRiwayat)
	s.SetelDaftar(DaftarRiwayat, SalinDaftar(riw))
	return s
}
