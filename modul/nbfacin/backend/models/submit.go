package models

import "github.com/cockroachdb/apd/v3"

// Submit tahap marketing (tiket 57): Email Type Accept / Reject / Banding / Decline + Comment.

// Kode Email Type tahap marketing (`.EmailType`, `SaveViewSuggest` CARI9): 1 Accept, 2 Reject, 4 Banding, 7 Decline.
const (
	EmailTypeAccept  = "1"
	EmailTypeReject  = "2"
	EmailTypeBanding = "4"
	EmailTypeDecline = "7"
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

// IsiSubmit - badan Submit / periksa.
type IsiSubmit struct {
	EmailType, Comment, BandingTo string
	SebabTc                       []PilihanSebabTc
	SebabNonTc                    string
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
	// Performance - POOLDATA.FACINPERFORMANCE (InsertDataFacinPerformance_SQL); nil = tidak ditulis.
	Performance *PerformanceFacIn
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
}

// HasilCekRW - satu kode pos lokasi di tabel RW: ada sama sekali, dan ada berpasangan dengan kota / kecamatan / provinsi.
type HasilCekRW struct {
	Zip, City, District, Province bool
}
