package models

// Untuk apa berkas ini: KATALOG KOLOM - satu-satunya pemetaan properti Pega ->
// kolom Oracle -> golongan tipe untuk penyimpanan relasional NB Treaty In
// (spec-penyimpanan-relasional.md ID-5..ID-31; tiket 00, 16, 17, 18, 19).
//
// Migrasi `backend/migrations/32x_*.sql`, `docs/STRUKTUR-TABEL-NB-TREATY-IN.md`,
// dan repository dicocokkan terhadap katalog ini oleh uji
// (`katalog_test.go`, `penjaga strukturkolom_test.go`). Menambah kolom berarti
// menambah satu baris di sini, satu baris DDL, dan satu baris STRUKTUR.
//
// Dasar daftar medan: `docs/SENSUS-PROPERTI-POLICYTREATYIN.md` (sensus XML
// rule terjangkau, tiket 00), dipilah ke medan yang DIBACA/DITULIS rule dan
// tampil di layar; keputusan work owner 23-09-2026 "ikuti dari data yang
// digunakan di Activity dan Section".
//
// ⛔ TIDAK disimpan, dan sebabnya:
//   - TotalPremium, TotalClaim, TotalSharePercentagePremium/Claim - TURUNAN:
//     jumlah baris SpreadingRiskList (`CountSpreading_Act` langkah 4.2-5);
//     dihitung saat dibaca (`HitungTotalSpreading`). Penjaga
//     `TestMigrasiTidakMenyimpanTotalPeserta` melarang kolom ber-awalan total.
//   - Layer/LayerType/LayerPart/LayerPartType tingkat polis - PANTULAN baris
//     pertama view (spec-penyimpanan ID-22); dibaca ulang dari view lewat
//     TREATY_IN_ID, tidak disimpan ganda.
//   - BreakDownSpreadList - KEPUTUSAN-RONDE-12 butir 3/3b (tidak dimigrasi).
//   - isApprovedtoDeptHead - P36, AC 64.
//   - Show, ViewState, pxResults - keadaan layar, bukan data.

import (
	"errors"
	"fmt"
)

// Golongan adalah kategori tipe logis sebuah kolom (spec-penyimpanan ID-14).
type Golongan string

const (
	GolTeks         Golongan = "teks"
	GolKode         Golongan = "kode"    // teks - nol di depan bermakna (ID-16)
	GolPenanda      Golongan = "penanda" // teks - "" berbeda dari "0" (ID-17)
	GolUang         Golongan = "uang"    // NUMBER(38,8)
	GolPersen       Golongan = "persen"  // NUMBER(38,8) - 12.5 = 12,5 persen
	GolTanggal      Golongan = "tanggal" // DATE, tanggal saja
	GolTanggalWaktu Golongan = "tanggal-waktu"
	GolCacah        Golongan = "cacah" // NUMBER(10)
)

// Desimal - golongan bertipe NUMBER(38,8).
func (g Golongan) Desimal() bool { return g == GolUang || g == GolPersen }

// Tanggal - golongan bertipe DATE.
func (g Golongan) Tanggal() bool { return g == GolTanggal || g == GolTanggalWaktu }

// Kolom memetakan satu properti ke satu kolom.
type Kolom struct {
	// Properti - jalur relatif pyWorkPage (tabel induk) atau nama anggota daftar
	// (tabel anak).
	Properti string
	Kolom    string
	Golongan Golongan
	// Panjang - panjang VARCHAR2; 0 untuk golongan non-teks.
	Panjang int
}

// Tabel - satu tabel katalog.
type Tabel struct {
	Nama string
	// Daftar - jalur PageList sumber baris (tabel anak); kosong untuk tabel 1:1.
	Daftar string
	Kolom  []Kolom
}

func kTeks(p, k string, n int) Kolom { return Kolom{p, k, GolTeks, n} }
func kKode(p, k string, n int) Kolom { return Kolom{p, k, GolKode, n} }
func kPenanda(p, k string) Kolom     { return Kolom{p, k, GolPenanda, 16} }
func kUang(p, k string) Kolom        { return Kolom{p, k, GolUang, 0} }
func kPersen(p, k string) Kolom      { return Kolom{p, k, GolPersen, 0} }
func kTgl(p, k string) Kolom         { return Kolom{p, k, GolTanggal, 0} }
func kCacah(p, k string) Kolom       { return Kolom{p, k, GolCacah, 0} }
func kTglWaktu(p, k string) Kolom    { return Kolom{p, k, GolTanggalWaktu, 0} }

const pt = HalamanPolis + "."

// TabelGeneralPolis - T_GENERAL_POLIS: satu baris per GENERASI polis, kunci
// utama bersama T_WORK_POLIS (ID-7). Kolom kunci, generasi, dan alur ditulis
// repository di luar katalog (ID, NOPOLIS, PRODKE, OLD_POLIS_ID, IDPEGA,
// TGL_INPUT, USERNAME, TGL_TUTUP).
var TabelGeneralPolis = Tabel{Nama: "T_GENERAL_POLIS", Kolom: []Kolom{
	// alur kasus (halaman kerja)
	kTeks("PositionNote", "POSITION_NOTE", 64),
	kTeks("NBStatus", "NB_STATUS", 255),
	kKode("TreatyIn.ID", "TREATY_IN_ID", 64),
	// identitas dan rujukan
	kKode(pt+"NoOffer", "NO_OFFER", 64),
	kKode(pt+"MasterID", "MASTER_ID", 64),
	kPenanda(pt+"IsApproved", "IS_APPROVED"),
	kTeks(pt+"Suggest", "SUGGEST", 4000),
	kTglWaktu(pt+"SuggestDate", "SUGGEST_DATE"),
	kTeks(pt+"OperatorName", "OPERATOR_NAME", 128),
	kPenanda(pt+"IsNewPolicyNonProp", "IS_NEW_POLICY_NON_PROP"),
	kPenanda(pt+"HasFacOut", "HAS_FAC_OUT"),
	kPenanda(pt+"FlagPPH", "FLAG_PPH"),
	kPenanda(pt+"FlagRetroTreaty", "FLAG_RETRO_TREATY"),
	kPenanda(pt+"IsOJKNopolis", "IS_OJK_NOPOLIS"),
	kPenanda(pt+"DueTo", "DUE_TO"),
	kKode(pt+"TypeTax", "TYPE_TAX", 64),
	kKode(pt+"StatementType", "STATEMENT_TYPE", 64),
	kKode(pt+"TreatyGroupID", "TREATY_GROUP_ID", 64),
	kTeks(pt+"TreatyGroupName", "TREATY_GROUP_NAME", 255),
	kKode(pt+"TreatyGroupOldID", "TREATY_GROUP_OLD_ID", 64),
	kKode(pt+"OJKBusinessID", "OJK_BUSINESS_ID", 64),
	kKode(pt+"BizCode", "BIZ_CODE", 64),
	kTeks(pt+"BizName", "BIZ_NAME", 255),
	kKode(pt+"SOB", "SOB", 64),
	kTeks(pt+"SOBName", "SOB_NAME", 255),
	kKode(pt+"CedingCo", "CEDING_CO", 4000),
	kTeks(pt+"CedingCoName", "CEDING_CO_NAME", 4000),
	kKode(pt+"InsuredID", "INSURED_ID", 64),
	kTeks(pt+"InsuredName", "INSURED_NAME", 255),
	kTeks(pt+"MarketingOfficer", "MARKETING_OFFICER", 255),
	kKode(pt+"TreatyType", "TREATY_TYPE", 64),
	kKode(pt+"TreatyYear", "TREATY_YEAR", 16),
	kKode(pt+"Currency", "CURRENCY", 16),
	kKode(pt+"IDCurrency", "ID_CURRENCY", 64),
	kKode(pt+"ShareCurrency", "SHARE_CURRENCY", 16),
	kKode(pt+"Quartal", "QUARTAL", 16),
	kKode(pt+"YearOfQuartal", "YEAR_OF_QUARTAL", 16),
	kKode(pt+"ClaimType", "CLAIM_TYPE", 64),
	kKode(pt+"ClaimPaymentType", "CLAIM_PAYMENT_TYPE", 64),
	kKode(pt+"Installment", "INSTALLMENT", 16),
	kTeks(pt+"Remark", "REMARK", 128),
	// tanggal - ProductionDate menjadi TGL_PROD (kolom datar json_polis, ID-21)
	kTgl(pt+"StartDate", "START_DATE"),
	kTgl(pt+"EndDate", "END_DATE"),
	kTglWaktu(pt+"StatementDate", "STATEMENT_DATE"),
	kTglWaktu(pt+"ProductionDate", "TGL_PROD"),
	// uang
	kUang(pt+"GrossPremium", "GROSS_PREMIUM"),
	kUang(pt+"GrossClaim", "GROSS_CLAIM"),
	kUang(pt+"PremiOgp", "PREMI_OGP"),
	kUang(pt+"ResultOgp1", "RESULT_OGP1"),
	kUang(pt+"ResultOgp2", "RESULT_OGP2"),
	kUang(pt+"PremiOnp", "PREMI_ONP"),
	kUang(pt+"ResultOnp1", "RESULT_ONP1"),
	kUang(pt+"ResultOnp2", "RESULT_ONP2"),
	kUang(pt+"Claim", "CLAIM"),
	kUang(pt+"OutstandingClaim", "OUTSTANDING_CLAIM"),
	kUang(pt+"SalvageValue", "SALVAGE_VALUE"),
	kUang(pt+"ExcessLoss", "EXCESS_LOSS"),
	kUang(pt+"NetPremium", "NET_PREMIUM"),
	kUang(pt+"BalanceDueTo", "BALANCE_DUE_TO"),
	kUang(pt+"BalanceBeforeTax", "BALANCE_BEFORE_TAX"),
	kUang(pt+"BalanceBeforePPH", "BALANCE_BEFORE_PPH"),
	// ⛔ PERTENTANGAN WO LAWAN XML - DIIKUTI WO. `[keputusan work owner]` P29
	// (PERTANYAAN-untuk-DBA, rancangan §4.1 "Persen", spec AC 26,
	// spec-penyimpanan AC 38): DEDUCTION1/2 adalah PERSENTASE. Rule Pega
	// memperlakukan nilai halaman polis sebagai jumlah - dikurangkan dari premi
	// (CountNetPremi_act langkah 4), dibagi 1,022 (SetPPNPPH langkah 4), label
	// layar `pxCurrency`. Golongan simpan mengikuti WO; rumusnya diport apa
	// adanya (AC 79). Dicatat di HASIL-IMPLEMENTASI bab 4 dan tiket 07.
	kPersen(pt+"Deduction1", "DEDUCTION1"),
	kPersen(pt+"Deduction2", "DEDUCTION2"),
	kUang(pt+"BrokerageFee", "BROKERAGE_FEE"),
	kUang(pt+"BrokerageFeeSebenarnya", "BROKERAGE_FEE_SEBENARNYA"),
	kUang(pt+"PPHValue", "PPH_VALUE"),
	kUang(pt+"PPNValue", "PPN_VALUE"),
	kUang(pt+"ShareValue", "SHARE_VALUE"),
	// persen - label layar "(%) Deduction In A/B"
	kPersen(pt+"RiCommOgp", "RI_COMM_OGP"),
	kPersen(pt+"OveriddingCommOgp", "OVERIDDING_COMM_OGP"),
	kPersen(pt+"RiCommOnp", "RI_COMM_ONP"),
	kPersen(pt+"OveriddingCommOnp", "OVERIDDING_COMM_ONP"),
}}

// TabelQuotation - T_POLIS_QUOTATION, 1:1 (ID-23). Halaman `pyWorkPage.Quotation`
// dan salinannya `PolicyTreatyIn.QuotationData` (DT pra-proses langkah 14,
// GeneratePolicyNoTreaty_Act langkah 10) disimpan SATU kali.
var TabelQuotation = Tabel{Nama: "T_POLIS_QUOTATION", Kolom: []Kolom{
	kKode("ProportionalType", "PROPORTIONAL_TYPE", 32),
	kKode("MOID", "MO_ID", 64),
	kKode("BusinessCode", "BUSINESS_CODE", 64),
	kKode("BusinessOldId", "BUSINESS_OLD_ID", 16),
	kKode("GroupPanel", "GROUP_PANEL", 16),
	kTeks("BusinessName", "BUSINESS_NAME", 255),
	kKode("BusinessType", "BUSINESS_TYPE", 64),
	kKode("BusinessFac", "BUSINESS_FAC", 16),
	kKode("SourceOfBusiness", "SOURCE_OF_BUSINESS", 64),
	kTeks("SobName", "SOB_NAME", 255),
	kKode("SobLeader0", "SOB_LEADER0", 64),
	kKode("SobLeader1", "SOB_LEADER1", 64),
	kKode("CedingCo", "CEDING_CO", 4000),
	kTeks("CedingCoName", "CEDING_CO_NAME", 4000),
	kKode("InsuredID", "INSURED_ID", 64),
	kTeks("InsuredName", "INSURED_NAME", 255),
	kKode("MarketingCode", "MARKETING_CODE", 64),
	kTeks("MarketingName", "MARKETING_NAME", 255),
	kKode("TeamGroup", "TEAM_GROUP", 64),
	kKode("BranchCode", "BRANCH_CODE", 64),
	kTeks("BranchName", "BRANCH_NAME", 255),
	kTeks("NoOfferSlip", "NO_OFFER_SLIP", 4000),
	kPenanda("IsSurveyReport", "IS_SURVEY_REPORT"),
	kKode("Type", "TYPE", 64),
	kKode("EdmType", "EDM_TYPE", 16),
	kKode("OldPolicyNo", "OLD_POLICY_NO", 64),
}}

// TabelCeding - T_POLIS_CEDING ← QuotationData.CedingCoList (ID-24).
var TabelCeding = Tabel{Nama: "T_POLIS_CEDING", Daftar: HalamanPolis + ".QuotationData.CedingCoList", Kolom: []Kolom{
	kKode("CedingCo", "CEDING_CO", 64),
	kTeks("CedingCoName", "CEDING_CO_NAME", 255),
}}

// TabelAngsuran - T_POLIS_INSTALMENT ← PolicyTreatyIn.ListInstallment (ID-26).
var TabelAngsuran = Tabel{Nama: "T_POLIS_INSTALMENT", Daftar: DaftarAngsuran, Kolom: []Kolom{
	kCacah("InstallmentNo", "INSTALLMENT_NO"), // cacah - ID-14
	kTgl("DueDate", "DUE_DATE"),
	kTgl("PaymentDate", "PAYMENT_DATE"),
	kPersen("InstallmentPercentage", "INSTALLMENT_PERCENTAGE"),
	kUang("Premium", "PREMIUM"),
	kUang("PaymentTotal", "PAYMENT_TOTAL"),
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
	kUang("PPN", "PPN"),
	kUang("PPh", "PPH"),
	kUang("PaymentTotalAfterPPN", "PAYMENT_TOTAL_AFTER_PPN"),
	kUang("PaymentTotalAfterTax", "PAYMENT_TOTAL_AFTER_TAX"),
}}

// TabelAngsuranRinci - T_POLIS_INSTALMENT_DETAIL ← ListInstallment().InstallmentList
// (hanya non-proporsional, rancangan §3.2).
var TabelAngsuranRinci = Tabel{Nama: "T_POLIS_INSTALMENT_DETAIL", Daftar: "InstallmentList", Kolom: []Kolom{
	kCacah("InstallmentNo", "INSTALLMENT_NO"), // cacah - ID-14
	kTgl("DueDate", "DUE_DATE"),
	kTgl("PaymentDate", "PAYMENT_DATE"),
	kPersen("InstallmentPercentage", "INSTALLMENT_PERCENTAGE"),
	kUang("Premium", "PREMIUM"),
	kUang("PaymentTotal", "PAYMENT_TOTAL"),
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
	kUang("PPN", "PPN"),
	kUang("PPh", "PPH"),
	kUang("PremiumAfterPPN", "PREMIUM_AFTER_PPN"),
	kUang("PremiumAfterTax", "PREMIUM_AFTER_TAX"),
}}

// TabelSpreading - T_POLIS_SPREADING ← PolicyTreatyIn.SpreadingRiskList (ID-28).
var TabelSpreading = Tabel{Nama: "T_POLIS_SPREADING", Daftar: DaftarSpreading, Kolom: []Kolom{
	kKode("TreatyType", "TREATY_TYPE", 64),
	kTeks("TreatyName", "TREATY_NAME", 255),
	kKode("Currency", "CURRENCY", 16),
	kKode("CurrencyID", "CURRENCY_ID", 64),
	kPersen("SharePercentage", "SHARE_PERCENTAGE"),
	kPersen("SplitRNMSharePct", "SPLIT_RNM_SHARE_PCT"),
	kPersen("ClaimPercentage", "CLAIM_PERCENTAGE"),
	kUang("PremiumSpreaded", "PREMIUM_SPREADED"),
	kUang("ClaimSpreaded", "CLAIM_SPREADED"),
}}

// TabelXOL - T_POLIS_XOL ← PolicyTreatyIn.TreatyXOLList (ID-29, ID-30).
var TabelXOL = Tabel{Nama: "T_POLIS_XOL", Daftar: HalamanPolis + ".TreatyXOLList", Kolom: []Kolom{
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
	kUang("GrossPremi", "GROSS_PREMI"),
	kUang("NetPremi", "NET_PREMI"),
	kUang("Deduction", "DEDUCTION"),
	kPenanda("DueTo", "DUE_TO"),
	kUang("DueToValue", "DUE_TO_VALUE"),
	kUang("BrokerageFeeSebenarnya", "BROKERAGE_FEE_SEBENARNYA"),
	kUang("PPHValue", "PPH_VALUE"),
	kUang("PPNValue", "PPN_VALUE"),
	kUang("NetPremiAfterPPH", "NET_PREMI_AFTER_PPH"),
	kUang("NetPremiAfterPPN", "NET_PREMI_AFTER_PPN"),
	kUang("NetPremiAfterTax", "NET_PREMI_AFTER_TAX"),
}}

// TabelLayerXOL - T_POLIS_XOL_LAYER ← TreatyXOLList().ValueList (ID-29).
var TabelLayerXOL = Tabel{Nama: "T_POLIS_XOL_LAYER", Daftar: "ValueList", Kolom: []Kolom{
	kKode("Layer", "LAYER", 64),
	kKode("LayerType", "LAYER_TYPE", 64),
	kKode("LayerPart", "LAYER_PART", 64),
	kKode("LayerPartType", "LAYER_PART_TYPE", 64),
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
	kUang("GrossPremi", "GROSS_PREMI"),
	kUang("NetPremi", "NET_PREMI"),
	kUang("Deduction", "DEDUCTION"),
	kPenanda("DueTo", "DUE_TO"),
	kUang("DueToValue", "DUE_TO_VALUE"),
	kUang("BrokerageFeeSebenarnya", "BROKERAGE_FEE_SEBENARNYA"),
	kUang("PPHValue", "PPH_VALUE"),
	kUang("PPNValue", "PPN_VALUE"),
	kUang("NetPremiAfterPPH", "NET_PREMI_AFTER_PPH"),
	kUang("NetPremiAfterPPN", "NET_PREMI_AFTER_PPN"),
	kUang("NetPremiAfterTax", "NET_PREMI_AFTER_TAX"),
}}

// TabelUsulan - T_POLIS_SUGGEST ← PolicyTreatyIn.SuggestList.
//
// ⭐ RALAT rancangan §4bis.1: tabel usulan DIBUTUHKAN. `SaveViewSuggest`
// menulis `HISTORYAKSEPTASIPRODUCTION` hanya bila `Quotation.BusinessFac == "F"`
// - termasuk kalang berketerangan "UNTUK TREATY" - sedangkan treaty bernilai
// "T" (`SetCategoryAttach`). Untuk kasus treaty, daftar usulan selama ini
// hanya tersimpan di dokumen JSON; tanpa tabel ini catatan pengguna hilang
// (AC 71, 72). OPERATOR_ID = identitas akses login (P4, AC 39).
var TabelUsulan = Tabel{Nama: "T_POLIS_SUGGEST", Daftar: DaftarUsulan, Kolom: []Kolom{
	kTeks("Suggest", "SUGGEST", 4000),
	kPenanda("IsApproved", "IS_APPROVED"),
	kTglWaktu("Date", "SUGGEST_DATE"),
	kTeks("OperatorName", "OPERATOR_NAME", 128),
	kKode("OperatorID", "OPERATOR_ID", 64),
	kPenanda("IsSave", "IS_SAVE"),
}}

// SemuaTabel - urutan tulis (induk lebih dulu).
var SemuaTabel = []Tabel{
	TabelGeneralPolis, TabelQuotation, TabelCeding, TabelAngsuran, TabelAngsuranRinci,
	TabelSpreading, TabelXOL, TabelLayerXOL, TabelUsulan,
}

// ErrBentukTidakSah - halaman membawa baris yang tidak boleh dimiliki jenis
// polisnya (spec-penyimpanan ID-26, ID-35; AC 31, 33).
var ErrBentukTidakSah = errors.New("models: bentuk halaman tidak sah untuk jenis proporsinya")

// PeriksaBentukSimpan menolak polis PROPORSIONAL yang membawa baris XOL
// (`TreatyXOLList`, AC 33) atau rincian angsuran bersarang
// (`ListInstallment().InstallmentList`, AC 31) - keduanya milik bentuk
// non-proporsional. Jenisnya dibaca dari QuotationData, lalu Quotation.
func PeriksaBentukSimpan(h *Halaman) error {
	jenis := h.Ambil(HalamanPolis + ".QuotationData.ProportionalType")
	if jenis == "" {
		jenis = h.Ambil(HalamanQuotation + ".ProportionalType")
	}
	if jenis != JenisProporsional {
		return nil
	}
	if len(h.AmbilDaftar(TabelXOL.Daftar)) > 0 {
		return fmt.Errorf("%w: polis Proportional membawa baris TreatyXOLList", ErrBentukTidakSah)
	}
	for i := range h.AmbilDaftar(DaftarAngsuran) {
		if len(h.AmbilDaftar(JalurAnak(DaftarAngsuran, i+1, TabelAngsuranRinci.Daftar))) > 0 {
			return fmt.Errorf("%w: polis Proportional membawa rincian angsuran bersarang", ErrBentukTidakSah)
		}
	}
	return nil
}
