package kontrak

// Kontrak Claim Fac In untuk Komite Claim Fac In - prompt work owner tahap 2 10-10-2026 §2 butir 1 (pola
// `KlaimTreatyNonPropKomite`: komite tanpa menu, kasus komite dibuka dari inbox klaimnya, tangga workbasket).
//
// Komite Claim Fac In memutus kasus `KMT-` yang dilahirkan Claim Fac In: TT2 satu baris adjustment (`CreateKMTNo_Act`),
// TT3 Reject Claim (`SendRejectClaimToKomite2`) dan TT4 Close Without Payment (`SendCloseClaimToKomite`) tanpa
// adjustment. Komite membaca kasus klaim induknya dan, lewat `KomitePost_Adjustment` / `KomitePost_Reject` /
// `KomitePost_CloseClaim` (korpus `Komite Claim FacIn`), menulis kembali ke kasus itu (`TempOpenPage`,
// Obj-Open-By-Handle `pxCoverInsKey`). Tulisan ke tabel klaim Fac In HANYA lewat kontrak ini, di dalam transaksi Submit
// komite.
//
// Berbeda dengan `KlaimTreaty` (kepala + satu adjustment): Fac In menulis balik di tiga tingkat pohon - objek
// (`ClaimData.ObjectList(o)`), item objek (`.ObjectItemList(i)`), dan adjustment (`.Adjustment(a)`). Disediakan
// `claimfacin` (`services.KlaimUntukKomite`), dipakai `komiteclaimfacin`. Jalur halaman memakai ejaan halaman Claim Fac
// In (ejaan properti Pega) karena itulah yang dibaca dan ditulis XML.

import (
	"context"
	"errors"
	"time"

	"nusantarare/inti/backend/db"
)

// KlaimFacInKomite adalah yang Komite Claim Fac In butuhkan dari Claim Fac In.
type KlaimFacInKomite interface {
	// BacaKlaimFacIn membaca halaman kasus klaim induk seperti dibuka Claim Fac In (`pyWorkCover` / `TempOpenPage`,
	// termasuk halaman polis `OfferFacIn`), beserta posisi baris adjustment `adjID` (ID baris T_CLAIM_ADJUSTMENT).
	// `adjID` kosong = kasus komite tanpa adjustment (TT3 / TT4). `tx` nil = baca di luar transaksi (layar).
	//
	// Langkah XML: Section `ShowTransfer` (Policy Detail(s), Object Detail, Claim Details), `SetValueKomite` S4-S12,
	// `ApprovalKomite_Act` S2, `KomitePost_Adjustment` S2/S7, `SaveAcceptation_KMT`, `SaveAccept_ACT`,
	// `SaveReject_ACT_KMT`, `HitServiceToKasirKMT_Act` S14.2, `InsertJsonClaimNonMBU_act`, `SendEmailKlaim_KMT`.
	BacaKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID, adjID string) (KlaimFacIn, error)
	// KunciKlaimFacIn = Obj-Open-By-Handle `pyWorkPage.pxCoverInsKey` `Lock=true` (`KomitePost_Adjustment` S2,
	// `KomitePost_Reject` S1, `KomitePost_CloseClaim` S1): baris kasus klaim dikunci sampai transaksi komite selesai.
	// Kasus klaim yang sudah ditutup -> `ErrKlaimFacInTertutup`.
	KunciKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID string) error
	// TulisBalikKlaimFacIn menerapkan ubahan komite atas kasus klaim induk lalu menyimpannya (Obj-Save `TempOpenPage`).
	// Ubahan hanya pada jalur daftar putih `JalurHeaderKomiteFacIn`, `PropAdjustmentKomiteFacIn`,
	// `PropObjekKomiteFacIn`, `PropItemKomiteFacIn`; jalur lain -> `ErrUbahanKlaimFacInTidakSah`. Ubahan objek / item /
	// adjustment menuntut `adjID` (pohon baris itu).
	TulisBalikKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID, adjID string, u UbahanKlaimFacIn) error
	// TutupKlaimFacIn = `pxForceCaseClose` kasus klaim induk: `KomitePost_Reject` S17 (`StatusKlaimDitolak`) dan
	// `KomitePost_CloseClaim` S14 (`StatusKlaimSelesai`). Status lain -> `ErrUbahanKlaimFacInTidakSah`.
	// `CloseAllSubCases=true`: kasus komite KMT- lain klaim itu yang masih terbuka ikut ditutup berstatus sama;
	// `komiteID` = kasus komite yang sedang memutus (dikecualikan - ditutup modul komite di akhir flow-nya).
	TutupKlaimFacIn(ctx context.Context, tx *db.Tx, klaimID, komiteID, status string, saat time.Time) error
}

// Status penutupan kasus klaim induk oleh komite (`pxForceCaseClose` WorkStatus).
const (
	StatusKlaimDitolak = "Resolved-Rejected"  // KomitePost_Reject S17
	StatusKlaimSelesai = "Resolved-Completed" // KomitePost_CloseClaim S14
)

// Jalur Nilai pembuat kasus klaim induk (`TempOpenPage.pxCreateOperator` / `.pxCreateOpName`, KomitePost_Reject S14.1).
const (
	JalurPembuatKlaim     = "pxCreateOperator"
	JalurNamaPembuatKlaim = "pxCreateOpName"
)

// KlaimFacIn adalah salinan baca halaman kasus klaim induk.
type KlaimFacIn struct {
	// Nilai - jalur halaman -> nilai (`pyID`, `ClaimData.NoClaim`, `ClaimData.CauseOfLoss`,
	// `OfferFacIn.QuotationData.BusinessOldId`, `OfferFacIn.PolicyData.PolicyNo`, ..., `JalurPembuatKlaim`).
	Nilai map[string]string
	// Daftar - jalur daftar -> baris (`ClaimData.ObjectList`, `ClaimData.ObjectList(1).ObjectItemList`,
	// `ClaimData.ObjectList(1).ObjectItemList(1).Adjustment`, `...EstimationList`, `...SpreadingList`, ...). Baris =
	// properti -> nilai.
	Daftar map[string][]map[string]string
	// Objek / Item / Adjustment - posisi (1..n) baris `adjID` di pohon objek -> item -> adjustment; 0 bila `adjID`
	// kosong.
	Objek, Item, Adjustment int
	// Tertutup - kasus klaim sudah ditutup (`T_WORK_CLAIM.STATUS_WORK` terisi).
	Tertutup bool
}

// UbahanKlaimFacIn adalah ubahan komite atas kasus klaim induk.
type UbahanKlaimFacIn struct {
	// Header - jalur halaman -> nilai baru; hanya kunci `JalurHeaderKomiteFacIn`.
	Header map[string]string
	// Adjustment / Objek / Item - properti baris adjustment `adjID`, objeknya, dan item objeknya -> nilai baru; hanya
	// kunci `PropAdjustmentKomiteFacIn` / `PropObjekKomiteFacIn` / `PropItemKomiteFacIn`.
	Adjustment map[string]string
	Objek      map[string]string
	Item       map[string]string
	// Riwayat - baris `ClaimData.Chronology` baru (`ChronologyInsertion_DT`, `Data.CARI12`).
	Riwayat []RiwayatKlaimFacIn
}

// RiwayatKlaimFacIn adalah satu baris kronologi kasus klaim (`ChronologyInsertion_DT`, T_VIEW_SUGGEST).
type RiwayatKlaimFacIn struct {
	// Teks - `Data.CARI12` ("Accepted by <jabatan> - KMT-..." / "Rejected by <jabatan> - KMT-...").
	Teks string
	// Pelaku - akun pemutus (`OperatorID.pyUserName`).
	Pelaku string
	// Tingkat - jabatan tangga pemutus (`.ASMUserID`; `ASMNoteType` "Committee").
	Tingkat string
	Saat    time.Time
}

// JalurHeaderKomiteFacIn - jalur kepala kasus klaim yang boleh ditulis Komite Claim Fac In, beserta langkah XML
// penulisnya. Tidak dimuat (PARITAS komite): `ClaimData.FlagOnGoingCommitte` (S7.2.1.5 / S15), `ClaimData.IDObjectUpdate`
// (S13), `ClaimData.ClaimComitee(<LAST>)` (KomitePost_Reject S6), `stsReject` (SaveAccept_ACT S9),
// `ClaimData.osAkseptasi` (SaveAccept_ACT S3-S7) - tanpa kolom di katalog Claim Fac In dan tanpa pembaca di layar.
var JalurHeaderKomiteFacIn = map[string]string{
	"AktifButton":               "KomitePost_Adjustment S7.2.1.7 (tingkat akhir: 0)",
	"ClaimData.IsCloseFile":     "KomitePost_Adjustment S13 (<- Adjustment.IsProposeClose)",
	"ClaimData.IsReservedClaim": "KomitePost_Adjustment S13 (<- Adjustment.IsPropReserved)",
}

// PropAdjustmentKomiteFacIn - properti baris adjustment yang boleh ditulis Komite Claim Fac In, beserta langkah XML
// penulisnya.
var PropAdjustmentKomiteFacIn = map[string]string{
	"AcceptanceStatus": "KomitePost_Adjustment S7.2.1.5 (1, tingkat akhir) dan S7.2.1.6 (2, tolak)",
	"AcceptedNo":       "KomitePost_Adjustment S7.2.1.5",
	"AcceptedDate":     "KomitePost_Adjustment S7.2.1.5",
	"Notes":            "KomitePost_Adjustment S7.2.1.5 (<- Comment)",
	"IsApproved":       "KomitePost_Adjustment S7.2.1.7 (<- AcceptStatus, tingkat akhir)",
	"IsFacRetro":       "SaveAcceptation_KMT S3 / S4 (KomitePost_Adjustment S7.2.1.16)",
	"IsPrintAccept":    "KomitePost_Adjustment S12.1.1",
	"StatusKasir":      "HitServiceToKasirKMT_Act S14.5 (KomitePost_Adjustment S25.2.1.1)",
	"IDOfBank":         "HitServiceToKasirKMT_Act S12-S13 (KomitePost_Adjustment S25.2.1.1)",
}

// PropObjekKomiteFacIn - properti objek (`ClaimData.ObjectList(o)`) baris adjustment yang boleh ditulis Komite Claim Fac
// In. `IsKomiteApprove` (S7.2.1.5 / S7.2.1.6) dan `KomiteList` (S7.2.1.8) tidak dimuat: tanpa kolom di katalog Claim Fac
// In (tangga dibaca dari tabel komite).
var PropObjekKomiteFacIn = map[string]string{
	"IsPrintAccept": "KomitePost_Adjustment S12 (@if(AcceptStatus==\"2\",\"1\",\"\"), tingkat akhir)",
	"DLAStatus":     "SaveAcceptation_KMT S5 (0)",
	"IsFacretro":    "SaveAcceptation_KMT S3 / S4",
}

// PropItemKomiteFacIn - properti item objek (`.ObjectItemList(i)`) baris adjustment yang boleh ditulis Komite Claim Fac
// In.
var PropItemKomiteFacIn = map[string]string{
	"IsFacretro": "SaveAcceptation_KMT S3 / S4",
}

var (
	// ErrKlaimFacInTidakAda - kasus klaim induk tidak ada.
	ErrKlaimFacInTidakAda = errors.New("kontrak: kasus klaim Fac In induk tidak ada")
	// ErrKlaimFacInTertutup - kasus klaim induk sudah ditutup.
	ErrKlaimFacInTertutup = errors.New("kontrak: kasus klaim Fac In induk sudah ditutup")
	// ErrAdjustmentFacInTidakAda - baris adjustment `adjID` tidak ada di kasus klaim induk.
	ErrAdjustmentFacInTidakAda = errors.New("kontrak: baris adjustment tidak ada di kasus klaim Fac In induk")
	// ErrUbahanKlaimFacInTidakSah - ubahan di luar daftar putih.
	ErrUbahanKlaimFacInTidakSah = errors.New("kontrak: ubahan klaim Fac In di luar daftar putih Komite")
)
