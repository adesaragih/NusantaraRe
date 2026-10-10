package models

import "github.com/cockroachdb/apd/v3"

// Submit tahap marketing (tiket 57): Email Type Accept / Reject / Banding / Decline + Comment.

// Kode Email Type (`.EmailType` marketing / `.EmailTypeTL` Team Leader, `SaveViewSuggest` CARI9): 1 Accept, 2 Reject,
// 3 Ask (TL, `DDL\EmailTypeTL.xml`), 4 Banding (marketing), 7 Decline.
const (
	EmailTypeAccept  = "1"
	EmailTypeReject  = "2"
	EmailTypeAsk     = "3"
	EmailTypeBanding = "4"
	EmailTypeDecline = "7"
	// EmailTypeRevise - tahap Binding (`DDL\EmailTypeBinding.xml` 9 Revise, tiket 64).
	EmailTypeRevise = "9"
)

// RiwayatAkseptasi - satu baris POOLDATA.HISTORYAKSEPTASIPRODUCTION (grid `OfferFacIn.ViewSuggest`). Tanggal = TGL_INP
// `DD-MM-YYYY HH:MM` (jam dinding WIB yang ditulis Pega / layanan ini, tanpa konversi zona).
type RiwayatAkseptasi struct {
	Tanggal, PIC, Approval, Status, Comment string
}

// BarisSebabTc - satu T/C popup Cause of Decline; Opsi = pilihan Option (sumbernya tidak di korpus -> selalu kosong).
type BarisSebabTc struct {
	TermCondition string
	Opsi          []string
}

// BarisCauseOfDecline - satu baris POOLDATA.CAUSEOFDECLINE apa adanya (data `D:\migrasi\RNM\DDL\CAUSEOFDECLINE.csv`, 41 baris;
// DDL-nya tidak ada di korpus).
type BarisCauseOfDecline struct {
	ID, TermCondition, Description, TermCondition2 string
}

// PilihanSebabTc - pilihan pengguna atas satu T/C (`.CheckBox1` / `.CheckBox2` = Prioritas "I" / "II").
type PilihanSebabTc struct {
	TermCondition, Prioritas, Description1, Description2, PctRNM, PctSOB string
}

// IsiSubmit - badan Submit / periksa. DeductibleAcceptance / CommentDeductible - isian Team Leader (tiket 58;
// `.OfferFacIn.IsDeductibleAcceptance` / `.CommentDeductible`); marketing mengirim false / "".
type IsiSubmit struct {
	EmailType, Comment, BandingTo string
	SebabTc                       []PilihanSebabTc
	SebabNonTc                    string
	DeductibleAcceptance          bool
	CommentDeductible             string
	// DaysValidity - tahap Underwriting (tiket 61): Validity Of Acceptance (hari, teks bulat 0-9999); "" = mode bawaan
	// GetDateValidity_ACT.
	DaysValidity string
	// ReceivedRiSlip - tahap Binding (tiket 64): checkbox Received RI Slip (`.OfferFacIn.ReceivedRiSlip`); hanya Accept.
	ReceivedRiSlip bool
}

// KonteksSubmit - keadaan kasus yang dibaca Submit: posisi / FlagOnGoingPolicy (T_WORK_POLIS), marketing (MARKETINGOFFICER
// lewat T_QUOTATIONDATA.MOID), dan medan riwayat / performance.
type KonteksSubmit struct {
	Position, FlagOnGoingPolicy string
	MOID                        string
	// MarketingName = MARKETINGOFFICER.CLIENTNAME; TeamGroup = MARKETINGOFFICER.TEAMGROUP; MOLeader =
	// MARKETINGOFFICER.MOLEADER (saringan GetDataMarketing_SQL). Kosong bila MOID kosong / tidak ada.
	MarketingName, TeamGroup, MOLeader    string
	InsuredName, SobName, ClassOfBusiness string
	// PercentShare = T_GENERAL_POLIS.PERCENT_SHARE (.OfferFacIn.PercentShare), RICommision = PAY_RI_COMMISION; nil = kosong.
	PercentShare, RICommision *apd.Decimal
	// DeductibleAcceptance (teks Pega "true" / "false" / "") / CommentDeductible - T_GENERAL_POLIS (migrasi 210, tiket 58).
	DeductibleAcceptance, CommentDeductible string
	// Tiket 61 (tahap Underwriting): CreateOp = T_WORK_POLIS.CREATE_OP (pxCreateOperator); StartDateTime / DateToUW =
	// T_GENERAL_POLIS teks Pega DateTime; DaysValidity / DateValidity / WaitingBindDate = isian tersimpan (migrasi 212).
	CreateOp, StartDateTime, DateToUW           string
	DaysValidity, DateValidity, WaitingBindDate string
	// ConfirmBinding ("1" / "0" / "") / ReceivedRiSlip ("true" / "false" / "") - T_GENERAL_POLIS (migrasi 213, tiket 64).
	ConfirmBinding, ReceivedRiSlip string
	// MarketingCode - MARKETINGOFFICER.CLIENTID lewat MOID (`.QuotationData.MarketingCode`: GetMOID_Sql `clientid =
	// MarketingCode`); penanda group E (tiket 63).
	MarketingCode string
	// StatusWork - T_WORK_POLIS.STATUS_WORK ("Resolved-…" = case selesai, tiket 69a-2).
	StatusWork string
}

// CariKasusGanda - ParamInProdCase.CARI10-15 ProtectCoverage_Act langkah 2 (tiket 73 P4): InsuredName sudah tanpa "." dan
// spasi; Mulai / Akhir dd/MM/yyyy WIB; IDPegaSendiri = pzInsKey case (K73-31).
type CariKasusGanda struct {
	InsuredName, SobName, BusinessCode, NoOfferSlip, Mulai, Akhir, IDPegaSendiri string
}

// AkunSubmit - akun pelaku dari M_LOGIN_GO: Nama (`OperatorID.pyUserName`), Divisi (`OperatorID.pyOrgDivision` =
// DIVISION_CODE). Kosong bila akun tidak ada.
type AkunSubmit struct {
	Nama, Divisi string
}

// TulisanSubmit - seluruh tulisan satu Submit (satu transaksi).
type TulisanSubmit struct {
	CaseID, IDPega         string
	PosisiLama, PosisiBaru string
	// StatusWorkBaru - T_WORK_POLIS.STATUS_WORK sesudah Submit; "" = tidak diubah.
	StatusWorkBaru string
	NBStatus       string
	// Riwayat - POOLDATA.HISTORYAKSEPTASIPRODUCTION (InsertViewSuggest_SQL); TglInp = jam dinding WIB.
	Riwayat RiwayatProduksi
	// Pega - POOLDATA.HISTORYAKSEPTASIPEGA (InsertHistoryAkseptasiPega_Sql).
	Pega RiwayatPega
	// Performance - POOLDATA.FACINPERFORMANCE (InsertDataFacinPerformance_SQL); nil = tidak ditulis. PerformanceLanjut -
	// baris kedua InsertFacInPerformance_Act langkah 13 (Accept tahap Underwriting: OUTSTANDING tahap berikutnya); nil = tidak.
	Performance, PerformanceLanjut *PerformanceFacIn
	// Validitas - Submit tahap Underwriting (tiket 61, migrasi 212); nil = tidak diubah.
	Validitas *ValiditasUW
	// Deductible - isian Team Leader disimpan (T_GENERAL_POLIS, migrasi 210); nil = tidak diubah.
	Deductible *DeductibleTL
	// KeUW - teks Pega DateTime Accept Team Leader (IsJuniorUW_B: DATE_TO_UW = ini, SUBMIT_TO_UW = ini bila kosong); "" = tidak.
	KeUW string
	// FlagBaru - T_WORK_POLIS.FLAG_ONGOING_POLICY sesudah Submit (tiket 64: tahap Binding -> Policy "1" / Offer "0"); "" = tetap.
	FlagBaru string
	// StatusBinding - CONFIRM_BINDING / RECEIVED_RI_SLIP (migrasi 213, tiket 64); nil = tidak diubah.
	StatusBinding *StatusBinding
	// Binding - Accept tahap Underwriting ke jalur akhir (tiket 63: SAVE JSON_OFFER + FACINOFFER + Transition38); nil = tidak.
	Binding *TulisanBinding
	// Policy - selesai Policy (tiket 69a-2, SaveJsonPolicyFacIn_Act); nil = tidak. NBStatusTetap - NB_STATUS tidak ditulis
	// (jalur selesai Policy tidak men-SET NBStatus).
	Policy        *TulisanPolicy
	NBStatusTetap bool
}

// StatusBinding - teks Pega: ConfirmBinding "1" / "0", ReceivedRiSlip "true" / "false".
type StatusBinding struct {
	ConfirmBinding, ReceivedRiSlip string
}

// TekanBinding - tombol Binding (tiket 64, `Activity\UpdateTglBindJsonOffer_Act.xml`): satu transaksi.
type TekanBinding struct {
	CaseID, IDPega string
	// Position / Flag - penjaga serentak: kasus masih di posisi + flag ini.
	Position, Flag string
	// JSONData - JSON_OFFER.JSONDATA minimal (tiket 63 D1); TGL_BIND = SYSDATE.
	JSONData string
	// Riwayat - baris "Offer has been bind." (HISTORYAKSEPTASIPRODUCTION).
	Riwayat RiwayatProduksi
}

// TulisanBinding - Utility16 SaveJsonOfferFacIn_Act (JSON_OFFER, port PEGA_M_JSON_OFFER tanpa CALL) + SaveOfferProduction_Act
// (FACINOFFER) + Transition38 (FLAG_ONGOING_POLICY).
type TulisanBinding struct {
	// FlagBaru - T_WORK_POLIS.FLAG_ONGOING_POLICY (Transition38 `.FlagOnGoingPolicy = 2`).
	FlagBaru string
	// JSONData - JSON_OFFER.JSONDATA minimal (keputusan work owner D1); BulanNoOffer - "MM.YYYY" nomor offer baru.
	JSONData, BulanNoOffer string
	// Offer - baris FACINOFFER (satu per baris spreading); NOOFFER / SPREAD_DATE diisi repository dari JSON_OFFER.
	Offer []BarisFacInOffer
}

// BarisFacInOffer - satu baris POOLDATA.FACINOFFER cabang FIRE (InsertOfferProduction_Sql, 50 kolom; NOOFFER dan SPREAD_DATE
// dari JSON_OFFER di repository). Tanggal = teks "DD/MM/YYYY HH24:MI:SS" jam dinding WIB ("" = NULL); angka *apd.Decimal (nil
// = NULL); teks "" = NULL.
type BarisFacInOffer struct {
	IDPega, BusinessCode, CedingCo, SOB, BeginDate, EndDate                            string
	Province, Kabupaten, Kecamatan, Kelurahan, RoadName, BuildingNo, ZipCode, ObjectNo string
	ObjectName, RiskLocation, ObjectItem, Coverage, CoverageNote, Rate, JnReas, CurrID string
	AccumulationCode, OccupationCode, MarketingOfficerCode, InsuredName, RISlipCeding  string
	GroupType, GroupName, TopRisk                                                      string
	// NoSerial - NOSERIAL (tiket 73 P4: `.SerialNo` objek Aneka cabang 13); "" = NULL (Pega "NULL", cabang lain).
	NoSerial                                                    string
	PercentRNM, PercentShareSpreaded, TSI, Premi, PercentRIComm *apd.Decimal
}

// ValiditasUW - DAYS_VALIDITY / DATE_VALIDITY / WAITING_BIND_DATE / COMMENT_TEXT T_GENERAL_POLIS (tanggal = teks Pega DateTime).
type ValiditasUW struct {
	Days, Date, WaitingBind, Comment string
}

// DeductibleTL - IS_DEDUCTIBLE_ACCEPTANCE (teks "true" / "false") dan COMMENT_DEDUCTIBLE.
type DeductibleTL struct {
	Acceptance, Comment string
}

// RiwayatProduksi - kolom InsertViewSuggest_SQL.
type RiwayatProduksi struct {
	TypePolis, Posisi, PIC, TglInp, Div, Type, Putaran, Approval, Keterangan, AksesLogin string
	PercentRNM                                                                           *apd.Decimal
}

// RiwayatPega - kolom InsertHistoryAkseptasiPega_Sql (ID_KOMITE tidak diisi jalur ini).
type RiwayatPega struct {
	Status, Username, Workbasket string
}

// PerformanceFacIn - kolom InsertDataFacinPerformance_SQL (RISLIPRECEIVED tidak ada di sistem baru -> NULL).
type PerformanceFacIn struct {
	PositionDoc, StatusDoc, Username, InsuredName, SOB, MOLeader, COB string
	// RISlipReceived - RISLIPRECEIVED VARCHAR2(10) (`.OfferFacIn.ReceivedRiSlip` "true" / "false"; tiket 64); "" = NULL.
	RISlipReceived string
}

// HasilCekRW - satu kode pos lokasi di tabel RW: ada sama sekali, dan ada berpasangan dengan kota / kecamatan / provinsi.
type HasilCekRW struct {
	Zip, City, District, Province bool
}
