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
// ⛔⛔ TABEL DAN KOLOM MENGIKUTI DIAGRAM GRILLING (bab 0 butir 11-12 PROMPT
// putaran 2, keputusan work owner 03-10-2026): tepat delapan tabel
// (`Diagram-Skema-Tabel-NusantaraRe.xlsx` sheet NB Treaty In Prop/NonProp),
// kolom menurut diagram + `docs/rancangan-tabel-datar-treaty-in.md`. Kolom di
// luar keduanya HANYA bila XML membuktikan medannya DIBACA rule NB terjangkau
// (syarat, rumus, sel Section) - masing-masing bertanda `RALAT` di bawah dan
// dirinci di `docs/PERBANDINGAN-KOLOM-DIAGRAM.md` (RALAT rancangan §4sexies).
// Medan yang hanya DITULIS rule, hanya dibaca langkah berlabel `//`, atau
// hanya dibaca jalur treaty keluar (K8 butir 4) TIDAK punya kolom.
//
// Dasar sensus: `docs/SENSUS-PROPERTI-POLICYTREATYIN.md` (rule terjangkau).
//
// ⛔ TIDAK disimpan, dan sebabnya (PolicyTreatyIn 79 medan diagram = 69 kolom
// katalog + PolicyNo sebagai NOPOLIS + sembilan di bawah, lihat PERBANDINGAN):
//   - TotalPremium, TotalClaim, TotalSharePercentagePremium/Claim - TURUNAN:
//     jumlah baris SpreadingRiskList (`CountSpreading_Act` langkah 4.2-5);
//     dihitung saat dibaca (`HitungTotalSpreading`). Penjaga
//     `TestMigrasiTidakMenyimpanTotalPeserta` melarang kolom ber-awalan total.
//   - Layer/LayerType/LayerPart/LayerPartType tingkat polis - PANTULAN baris
//     pertama view (spec-penyimpanan ID-22); dibaca ulang dari view lewat
//     TREATY_IN_ID, tidak disimpan ganda.
//   - BreakDownSpreadList - KEPUTUSAN-RONDE-12 butir 3/3b (tidak dimigrasi).
//   - isApprovedtoDeptHead - P36, AC 64.
//   - IsOJKNopolis, BrokerageFee - hanya DITULIS (GeneratePolicyNoTreaty_Act
//     langkah 23 berlabel `//`; SetPPNPPH langkah 4), nol pembaca; di luar
//     diagram/rancangan.
//   - Show, ViewState, pxResults - keadaan layar, bukan data.
//
// Medan dokumen lama lain tanpa kolom: keputusan F3 per medan (WO
// 04-10-2026) - dibuang berbukti di `medan_abaikan_lama.json` bagian `pola`.

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
	GolUang         Golongan = "uang"    // NUMBER(38,10)
	GolPersen       Golongan = "persen"  // NUMBER(38,10) - 12.5 = 12,5 persen
	GolTanggal      Golongan = "tanggal" // DATE, tanggal saja
	GolTanggalWaktu Golongan = "tanggal-waktu"
	GolCacah        Golongan = "cacah" // NUMBER(10)
)

// Desimal - golongan bertipe NUMBER(38,10): diagram sheet NB Treaty In Prop F20
// "uang · persen -> angka presisi tetap, skala MINIMAL 9 desimal (P29)", dan
// J69 "NB: 100 / jumlah baris presisi 10" tersimpan utuh (perintah WO
// 04-10-2026; RALAT NUMBER(38,8) KEPUTUSAN 23-09-2026 sore).
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

// TabelGeneralPolis - T_GENERAL_POLIS_TREATY: satu baris per GENERASI polis, kunci
// utama bersama T_WORK_POLIS (ID-7). Kolom kunci dan generasi ditulis
// repository di luar katalog (ID, NOPOLIS, PRODKE, NOENDORS, OLD_POLIS_ID,
// TGL_INPUT, USERNAME; IDPEGA dibuang 06-10-2026 - sama dengan ID). Generasi TERTUTUP = ada baris penerus yang
// OLD_POLIS_ID-nya menunjuk generasi ini (ID-10) - tanpa kolom penanda.
var TabelGeneralPolis = Tabel{Nama: "T_GENERAL_POLIS_TREATY", Kolom: []Kolom{
	// RALAT - halaman kerja, bukan PolicyTreatyIn; T_WORK_POLIS (milik
	// premiumlistlife) tidak punya kolomnya. PositionNote dibaca connector Flow
	// dan syarat Section; NBStatus tampil di SFAPortal_OpportunitiesList;
	// TreatyIn.ID dibaca RDBList BrowseTreatyIn ({TreatyIn.ID}) dan
	// CheckDuplicateOffer - master kontrak dimuat ulang darinya.
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
	// RALAT F3 (WO 04-10-2026) - rancangan §4.1 penentu bentuk EDM_TYPE;
	// medan dokumen lama (data guide `$.EDMType`) DIBACA syarat
	// `InputPolicyTreatyInPre_Act` langkah 10 (`.PolicyTreatyIn.EDMType=="3"`
	// -> lewati `TreatyRealizationCheckXOLList`; `PerluCekDaftarXOL`). Properti
	// lain dari `QuotationData.EdmType` (T_POLIS_QUOTATION.EDM_TYPE, J37).
	kKode(pt+"EDMType", "EDM_TYPE", 16),
	// rancangan §4.1 - ditulis dan dibaca InputPolicyTreatyInDetail_NonProp
	// langkah 10-11 (jalur XOL, K8)
	kPenanda(pt+"IsEDMInputOnNB", "IS_EDM_INPUT_ON_NB"),
	kPenanda(pt+"HasFacOut", "HAS_FAC_OUT"),
	kPenanda(pt+"FlagPPH", "FLAG_PPH"),
	kPenanda(pt+"FlagRetroTreaty", "FLAG_RETRO_TREATY"),
	kPenanda(pt+"DueTo", "DUE_TO"),
	kKode(pt+"TypeTax", "TYPE_TAX", 64),
	kKode(pt+"StatementType", "STATEMENT_TYPE", 64),
	kKode(pt+"TreatyGroupID", "TREATY_GROUP_ID", 64),
	kTeks(pt+"TreatyGroupName", "TREATY_GROUP_NAME", 255),
	kKode(pt+"TreatyGroupOldID", "TREATY_GROUP_OLD_ID", 64),
	kKode(pt+"OJKBusinessID", "OJK_BUSINESS_ID", 64),
	// rancangan §4.1 (data guide `$.IDNewBisnis`); nol rule NB - tempat medan
	// dokumen lama bagi pemuat (tiket 22)
	kKode(pt+"IDNewBisnis", "ID_NEW_BISNIS", 64),
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
	// RALAT - tampil "RNM Share" DetailPolicyTreatyIn dan DetailDeptHeadTreatyIn_UW
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
	// RALAT - tampil "Claim 100%" DetailPolicyTreatyIn; dibaca CalculatePremi_Act
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
	kUang(pt+"Deduction1", "DEDUCTION1"), // K3 (03-10-2026): uang seperti XML (pxCurrency; CountNetPremi_act 4) - menggantikan P29 di atas, tiket 07
	kUang(pt+"Deduction2", "DEDUCTION2"), // K3 (03-10-2026): uang seperti XML (pxCurrency; CountNetPremi_act 4) - menggantikan P29 di atas, tiket 07
	// RALAT - dibaca rumus SetPPNPPH langkah 4 (.PPHValue/.PPNValue =
	// .BrokerageFeeSebenarnya * ...)
	kUang(pt+"BrokerageFeeSebenarnya", "BROKERAGE_FEE_SEBENARNYA"),
	kUang(pt+"PPHValue", "PPH_VALUE"),
	kUang(pt+"PPNValue", "PPN_VALUE"),
	kUang(pt+"ShareValue", "SHARE_VALUE"),
	// persen - label layar "(%) Deduction In A/B"
	kPersen(pt+"RiCommOgp", "RI_COMM_OGP"),
	kPersen(pt+"OveriddingCommOgp", "OVERIDDING_COMM_OGP"),
	kPersen(pt+"RiCommOnp", "RI_COMM_ONP"),
	kPersen(pt+"OveriddingCommOnp", "OVERIDDING_COMM_ONP"),
	// WO 08-10-2026 (cek Copy Old: medan dokumen lama tanpa kolom): "IsSOAUpload ITU PERLU" - nol rujukan korpus,
	// disimpan apa adanya ("1" di DEV). GuaranteeFund TIDAK berkolom (WO 08-10-2026 "GUARANTEE_FUND NUMBER(38,10), buang!";
	// medan_abaikan_lama.json).
	kKode(pt+"IsSOAUpload", "IS_SOA_UPLOAD", 16),
}}

// TabelQuotation - T_POLIS_QUOTATION, 1:1 (ID-23; diagram J35-J37 "10 medan
// skalar"). Halaman `pyWorkPage.Quotation` dan salinannya
// `PolicyTreatyIn.QuotationData` (preACT langkah 14.9, DT pra-proses langkah 14,
// GeneratePolicyNoTreaty_Act langkah 10) disimpan SATU kali.
//
// ⛔ Dibuang (di luar diagram, nol pembaca NB terjangkau): BusinessType (hanya
// ditulis DT BusinessType_DeT 14.8; turunan GroupPanel + BusinessOldId),
// SobName/SobLeader0/SobLeader1 (pembacanya CheckDataMkt langkah 7 `//` dan
// jalur treaty keluar), CedingCo/CedingCoName (tingkat polis di T_GENERAL_POLIS_TREATY,
// R47).
//
// RALAT 06-10-2026 (keputusan work owner, "CheckDataMkt ikuti aja itu semua"): SEMUA
// medan yang ditulis CheckDataMkt langkah 4 disimpan - MarketingCode, TeamGroup,
// BranchCode, BranchName (sebelumnya dibuang, nol pembaca). TeamGroup dibaca aturan
// tim Sec Head. Panjang kolom = kolom sumber MARKETINGOFFICER (VARCHAR2 100).
var TabelQuotation = Tabel{Nama: "T_POLIS_QUOTATION", Kolom: []Kolom{
	// diagram - 10 medan
	kKode("ProportionalType", "PROPORTIONAL_TYPE", 32),
	kKode("MOID", "MO_ID", 64),
	kKode("BusinessCode", "BUSINESS_CODE", 64),
	kKode("BusinessOldId", "BUSINESS_OLD_ID", 16),
	kKode("GroupPanel", "GROUP_PANEL", 16),
	kKode("SourceOfBusiness", "SOURCE_OF_BUSINESS", 64),
	kKode("Type", "TYPE", 64),
	kKode("EdmType", "EDM_TYPE", 16),
	kKode("OldPolicyNo", "OLD_POLICY_NO", 64),
	kTeks("MarketingName", "MARKETING_NAME", 255),
	// RALAT - dibaca rule terjangkau / tampil di Section NB
	//   BusinessName   InputPolicyTreatyInPre_Act langkah 2 (syarat + CARI2);
	//                  tampil SFAPortal_OpportunitiesList (A.Quotation.BusinessName)
	//   BusinessFac    SaveViewSuggest CARI7 (TYPE riwayat produksi);
	//                  GetListOpportunity filter E (A.Quotation.BusinessFac = T)
	//   InsuredID      InputPolicyTreatyInDetail_preACT langkah 3 (rumus) dan
	//                  14.3 (syarat Quotation.InsuredID == "")
	//   InsuredName    preACT langkah 3 dan 14.1; tampil SFAPortal_OpportunitiesList
	//   NoOfferSlip    tampil DetailPolicyTreatyIn (diisi) dan DetailDeptHeadTreatyIn_UW
	//   IsSurveyReport tampil + wajib DetailPolicyTreatyIn; syarat tombol Survey Report
	kTeks("BusinessName", "BUSINESS_NAME", 255),
	kKode("BusinessFac", "BUSINESS_FAC", 16),
	kKode("InsuredID", "INSURED_ID", 64),
	kTeks("InsuredName", "INSURED_NAME", 255),
	kTeks("NoOfferSlip", "NO_OFFER_SLIP", 4000),
	kPenanda("IsSurveyReport", "IS_SURVEY_REPORT"),
	// RALAT 06-10-2026 - CheckDataMkt langkah 4 (ClientID, TeamGroup, BranchDetailID, BranchDetailName)
	kKode("MarketingCode", "MARKETING_CODE", 100),
	kKode("TeamGroup", "TEAM_GROUP", 100),
	kKode("BranchCode", "BRANCH_CODE", 100),
	kTeks("BranchName", "BRANCH_NAME", 255),
}}

// TabelCeding - T_POLIS_CEDING ← QuotationData.CedingCoList (ID-24), anak
// T_POLIS_QUOTATION lewat QUOTATION_ID (diagram O39; kolom R43
// "CEDING_CO_ID ← .CedingCo · CEDING_CO_NAME ← .CedingCoName").
var TabelCeding = Tabel{Nama: "T_POLIS_CEDING", Daftar: HalamanPolis + ".QuotationData.CedingCoList", Kolom: []Kolom{
	kKode("CedingCo", "CEDING_CO_ID", 64),
	kTeks("CedingCoName", "CEDING_CO_NAME", 255),
}}

// TabelAngsuran - T_POLIS_INSTALMENT ← PolicyTreatyIn.ListInstallment (ID-26;
// diagram J53 "10 medan" = rancangan §4.3 tanpa PAYMENT_DATE, yang hanya ada di
// anak - diagram R61).
var TabelAngsuran = Tabel{Nama: "T_POLIS_INSTALMENT", Daftar: DaftarAngsuran, Kolom: []Kolom{
	kCacah("InstallmentNo", "INSTALLMENT_NO"), // cacah - ID-14
	kTgl("DueDate", "DUE_DATE"),
	kPersen("InstallmentPercentage", "INSTALLMENT_PERCENTAGE"),
	kUang("Premium", "PREMIUM"),
	kUang("PaymentTotal", "PAYMENT_TOTAL"),
	kUang("PremiumAfterPPH", "PREMIUM_AFTER_PPH"),
	kUang("PremiumAfterPPN", "PREMIUM_AFTER_PPN"),
	kUang("PremiumAfterTax", "PREMIUM_AFTER_TAX"),
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
	// RALAT - ditulis lalu DIBACA rumus InputPolicyTreatyInDetail_preACT
	// langkah 18.3.4.1 (.PaymentTotalAfterPPN = .PaymentTotal+.PPN;
	// Local.PPNins = .PPN; Local.PPHins = .PPh; Local.PaymentTotalAfterPPN/Tax =
	// .PaymentTotalAfterPPN/Tax), ada di dokumen lama (data guide)
	kUang("PPN", "PPN"),
	kUang("PPh", "PPH"),
	kUang("PaymentTotalAfterPPN", "PAYMENT_TOTAL_AFTER_PPN"),
	kUang("PaymentTotalAfterTax", "PAYMENT_TOTAL_AFTER_TAX"),
}}

// TabelAngsuranRinci - T_POLIS_INSTALMENT_DETAIL ← ListInstallment().InstallmentList
// (hanya non-proporsional, rancangan §3.2; diagram R58 "11 medan", R61 "punya
// PAYMENT_DATE"). PPN/PPh baris anak (preACT 18.3.4.2.1) hanya DITULIS, nol
// pembaca - tidak punya kolom.
var TabelAngsuranRinci = Tabel{Nama: "T_POLIS_INSTALMENT_DETAIL", Daftar: "InstallmentList", Kolom: []Kolom{
	kCacah("InstallmentNo", "INSTALLMENT_NO"), // cacah - ID-14
	kTgl("DueDate", "DUE_DATE"),
	kTgl("PaymentDate", "PAYMENT_DATE"),
	kPersen("InstallmentPercentage", "INSTALLMENT_PERCENTAGE"),
	kUang("Premium", "PREMIUM"),
	kUang("PaymentTotal", "PAYMENT_TOTAL"),
	kUang("PremiumAfterPPH", "PREMIUM_AFTER_PPH"),
	kUang("PremiumAfterPPN", "PREMIUM_AFTER_PPN"),
	kUang("PremiumAfterTax", "PREMIUM_AFTER_TAX"),
	kKode("Currency", "CURRENCY", 16),
	kKode("IDCurrency", "ID_CURRENCY", 64),
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

// TabelSurvei - T_POLIS_SURVEY ← PolicyTreatyIn.QuotationData.SurveyReportList: popup Historical Survey
// Report (`Harness/HistoricalSurveyReport` ← tombol Survey Report). Keputusan work owner 06-10-2026,
// membatalkan K7 (03-10-2026): tabel kesembilan, di luar diagram grilling. Kolom = sel grid
// `Section/HistoricalSurveyReportDtl`: DateofSurvey pxDateTime, SurveyedBy teks, LossPrevention pxNumber
// 4 desimal, Remarks pxDropdown (nilai standar).
var TabelSurvei = Tabel{Nama: "T_POLIS_SURVEY", Daftar: DaftarSurvei, Kolom: []Kolom{
	kTgl("DateofSurvey", "DATE_OF_SURVEY"),
	kTeks("SurveyedBy", "SURVEYED_BY", 255),
	kPersen("LossPrevention", "LOSS_PREVENTION"),
	kKode("Remarks", "REMARKS", 64),
}}

// SemuaTabel - urutan tulis (induk lebih dulu).
var SemuaTabel = []Tabel{
	TabelGeneralPolis, TabelQuotation, TabelCeding, TabelAngsuran, TabelAngsuranRinci,
	TabelSpreading, TabelXOL, TabelLayerXOL, TabelSurvei,
}

// NilaiQuotation - nilai satu medan T_POLIS_QUOTATION: QuotationData bila
// terisi, selain itu Quotation (ID-23: keduanya disimpan SATU kali; layar
// menyunting QuotationData, aktivitas menulis Quotation, `SalinQuotation`
// menyelaraskannya di pra-proses). Saat dimuat keduanya diisi nilai yang sama.
func NilaiQuotation(h *Halaman, prop string) string {
	if v := h.Ambil(HalamanPolis + ".QuotationData." + prop); v != "" {
		return v
	}
	return h.Ambil(HalamanQuotation + "." + prop)
}

// keturunanKatalog - daftar anak tabel polis beserta cucunya (urutan tulis).
var keturunanKatalog = []struct {
	anak Tabel
	cucu *Tabel
}{
	{TabelCeding, nil},
	{TabelAngsuran, &TabelAngsuranRinci},
	{TabelSpreading, nil},
	{TabelXOL, &TabelLayerXOL},
	{TabelSurvei, nil},
}

// ProyeksiKatalog - halaman sebagaimana ia KELUAR dari penyimpanan: hanya
// medan yang punya kolom di delapan tabel diagram (katalog ini), nilai
// Quotation/QuotationData disatukan (`NilaiQuotation`), baris anak hanya
// kolom katalognya. Kunci di luar katalog (PolicyTreatyIn.PolicyNo = NOPOLIS)
// dan SuggestList (HISTORYAKSEPTASIPRODUCTION) TIDAK ikut - milik
// penulisnya sendiri. Normalisasi angka/tanggal Oracle tidak ditiru.
//
// Dipakai gudang tiruan supaya uji seam HTTP melihat apa yang Oracle simpan:
// medan tanpa kolom hilang sesudah disimpan, persis repository.
func ProyeksiKatalog(h *Halaman) *Halaman {
	s := HalamanBaru()
	for _, k := range TabelGeneralPolis.Kolom {
		if v := h.Ambil(k.Properti); v != "" {
			s.Setel(k.Properti, v)
		}
	}
	for _, k := range TabelQuotation.Kolom {
		if v := NilaiQuotation(h, k.Properti); v != "" {
			s.Setel(HalamanQuotation+"."+k.Properti, v)
			s.Setel(HalamanPolis+".QuotationData."+k.Properti, v)
		}
	}
	saring := func(t Tabel, b []Baris) []Baris {
		var out []Baris
		for _, x := range b {
			nb := Baris{}
			for _, k := range t.Kolom {
				if v := x[k.Properti]; v != "" {
					nb[k.Properti] = v
				}
			}
			out = append(out, nb)
		}
		return out
	}
	for _, kt := range keturunanKatalog {
		baris := h.AmbilDaftar(kt.anak.Daftar)
		if len(baris) == 0 {
			continue
		}
		s.SetelDaftar(kt.anak.Daftar, saring(kt.anak, baris))
		if kt.cucu == nil {
			continue
		}
		for i := range baris {
			j := JalurAnak(kt.anak.Daftar, i+1, kt.cucu.Daftar)
			if cucu := h.AmbilDaftar(j); len(cucu) > 0 {
				s.SetelDaftar(j, saring(*kt.cucu, cucu))
			}
		}
	}
	return s
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
