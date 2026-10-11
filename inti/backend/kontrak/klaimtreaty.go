package kontrak

// Kontrak Claim Prop untuk Komite Claim Prop - keputusan work owner 08-10-2026 (pola `KlaimKomite` Claim Life).
//
// Komite Claim Prop memutus SATU baris adjustment klaim treaty inward proporsional. Kasus komite `TKMT-` dilahirkan
// Claim Prop (opsi B 07-10-2026, `AddKomiteTreatyChild_ACT`); Komite membaca klaim induknya dan, lewat
// `KomitePostAdjustment`, menulis kembali ke kasus klaim itu (`TempOpenPage`, Obj-Open-By-Handle `pxCoverInsKey`).
// Tulisan ke tabel klaim Prop (`T_GENERAL_CLAIM`, `T_CLAIM_*`, baris klaim `T_VIEW_SUGGEST`) HANYA lewat kontrak ini, di
// dalam transaksi Komite sendiri - keputusan anggota, tangga, dan tulisan balik ke klaim satu transaksi.
//
// Disediakan `claimprop` (`services.KlaimUntukKomite`), dipakai `komiteclaimprop`. Jalur halaman (`Nilai`, `Daftar`,
// kunci ubahan) memakai ejaan halaman Claim Prop - ejaan properti Pega - karena itulah yang dibaca dan ditulis XML.

import (
	"context"
	"errors"
	"time"

	"nusantarare/inti/backend/db"
)

// KlaimTreatyKomite adalah yang Komite Claim Prop butuhkan dari Claim Prop.
//
// Kasus komite Close Without Payment (TT 4, `SendCloseClaimToKomite` -> `KomitePost_Close`; perintah work owner
// 10-10-2026) tidak menunjuk baris adjustment: `adjID` KOSONG berarti "tanpa adjustment" - `BacaKlaimTreaty` menjawab
// `Adjustment = 0` dan `TulisBalikKlaimTreaty` menolak ubahan `Adjustment`.
type KlaimTreatyKomite interface {
	// BacaKlaimTreaty membaca halaman kasus klaim induk seperti dibuka Claim Prop (`pyWorkCover` / `TempOpenPage`),
	// beserta posisi baris adjustment `adjID` (kosong = kasus komite TT 4, `Adjustment = 0`). `tx` nil = baca tanpa
	// kunci (layar); terisi = di dalam transaksi Submit, sesudah `KunciKlaimTreaty`.
	//
	// Langkah XML: Section `ShowTransfer` (panel klaim induk, grid Insured Interests / Count Claim Amount / Loss
	// Allocation / Estimation List / Total Original Currency Estimation / Spreading In / Spreading Out, blok
	// deductible, penerima bayar, teks komite), `SetKomiteList_Act` S4 (Page-Copy `ClaimData`),
	// `SetDataAcceptationTreaty_Act` S4-S8, `KomitePostAdjustment` S5/S6/S7/S15/S16, `SaveAcceptation_Act` S1,
	// `SaveAcceptationTreaty_TKMT` S1-S7, `HitServiceToKasirKMT_Act` S2/S14.1, `SendEmailKlaim_KMT` S1/S5-S12.
	BacaKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string) (KlaimTreaty, error)
	// KunciKlaimTreaty = `KomitePostAdjustment` S4 (Obj-Open-By-Handle `pyWorkPage.pxCoverInsKey`, `Lock=true`,
	// `ReleaseOnCommit=true`): baris kasus klaim dikunci sampai transaksi Komite selesai. Kasus klaim yang sudah
	// ditutup -> `ErrKlaimTreatyTertutup`.
	KunciKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID string) error
	// TulisBalikKlaimTreaty menerapkan ubahan Komite atas kasus klaim induk lalu menyimpannya (`KomitePostAdjustment`
	// S19 / S36 Obj-Save `TempOpenPage`). Ubahan hanya pada jalur daftar putih `JalurHeaderKomite` dan
	// `PropAdjustmentKomite`; jalur lain -> `ErrUbahanKlaimTreatyTidakSah`.
	TulisBalikKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, adjID string, u UbahanKlaimTreaty) error
	// TutupKlaimTreaty = `KomitePost_Close` S17 (`pxForceCaseClose` TempOpenPage, `WorkStatus "Resolved-Completed"`,
	// `CloseAllSubCases true`): kasus klaim induk ditutup Resolved-Completed bersama kasus komitenya yang masih terbuka,
	// KECUALI `komiteID` (kasus komite yang sedang memutus - ditutup modul Komite sendiri di transaksi yang sama).
	// Kasus klaim yang sudah ditutup -> `ErrKlaimTreatyTertutup`.
	TutupKlaimTreaty(ctx context.Context, tx *db.Tx, klaimID, komiteID string, saat time.Time) error
}

// Jalur `KlaimTreaty.Nilai` yang bukan properti halaman klaim (kepala work object, hanya-baca): pembuat kasus klaim
// induk - `KomitePost_Close` S13 `DataIn.CARI6` / `CARI7` (`TempOpenPage.pxCreateOpName` / `pxCreateOperator`).
const (
	JalurPembuatKlaimTreaty     = "pxCreateOperator"
	JalurNamaPembuatKlaimTreaty = "pxCreateOpName"
)

// KlaimTreaty adalah salinan baca halaman kasus klaim induk.
type KlaimTreaty struct {
	// Nilai - jalur halaman -> nilai (`ClaimData.NoClaim`, `TreatyInMaster.Ceding`,
	// `OfferFacIn.QuotationData.BusinessOldId`, ...).
	Nilai map[string]string
	// Daftar - jalur daftar -> baris (`ClaimData.InterestList`, `ClaimData.AdjustmentList`,
	// `ClaimData.AdjustmentList(n).SpreadingAdjustment`, ...). Baris = properti -> nilai.
	Daftar map[string][]map[string]string
	// Adjustment - posisi (1..n) baris `adjID` di `ClaimData.AdjustmentList`; 0 = kasus komite tanpa adjustment (TT 4).
	Adjustment int
	// Tertutup - kasus klaim sudah ditutup (`T_WORK_CLAIM.STATUS_WORK` terisi).
	Tertutup bool
}

// UbahanKlaimTreaty adalah ubahan Komite atas kasus klaim induk.
type UbahanKlaimTreaty struct {
	// Header - jalur halaman -> nilai baru; hanya kunci `JalurHeaderKomite`.
	Header map[string]string
	// Adjustment - properti baris adjustment `adjID` -> nilai baru; hanya kunci `PropAdjustmentKomite`.
	Adjustment map[string]string
	// Riwayat - baris `ClaimData.SuggestList` baru (`InsertChronology_DT`, `KomitePostAdjustment` S8-S10).
	Riwayat []RiwayatKlaimTreaty
	// FacRetro - `SaveAcceptationTreaty_TKMT` S6.3: baris `ClaimData.FacRetroList` yang ditambahkan HANYA bila daftar
	// itu masih kosong (S5 `SizeRetro`, S6.3 `Local.SizeRetro<1`). Kunci baris = properti `FacRetroList`.
	FacRetro []map[string]string
}

// RiwayatKlaimTreaty adalah satu baris riwayat kasus klaim (`InsertChronology_DT`).
type RiwayatKlaimTreaty struct {
	// Teks - `DataChronology.CARI1` ("Accepted by <jabatan>" / "Rejected by <jabatan>").
	Teks string
	// Pelaku - akun penyetuju (`OperatorID.pyUserName`).
	Pelaku string
	// Tingkat - label peran (`IsCedingConfirm`); keputusan 17-09: jabatan tangga, bukan nama yang di-hardcode.
	Tingkat string
	Saat    time.Time
}

// JalurHeaderKomite - jalur header kasus klaim yang boleh ditulis Komite, beserta langkah XML penulisnya.
var JalurHeaderKomite = map[string]string{
	"ClaimData.IsCloseFile":     "KomitePostAdjustment S11 (<- Adjustment.IsProposeClose)",
	"ClaimData.IsReservedClaim": "KomitePostAdjustment S11 (<- Adjustment.IsPropReserved)",
	"IsAnyAcceptation":          "KomitePostAdjustment S14 (tingkat akhir disetujui)",
	"AktifButton":               "KomitePostAdjustment S22 (tingkat akhir) dan S25 (ditolak)",
	"ClaimData.IsSubjectivity":  "KomitePostAdjustment S24 (<- pyWorkPage.IsSubjectivity)",
	"IsOutstanding":             "SaveAcceptation_Act S6 (KomitePostAdjustment S17)",
	"IsCFS":                     "SaveAcceptation_Act S6 (KomitePostAdjustment S17)",
}

// PropAdjustmentKomite - properti baris adjustment yang boleh ditulis Komite, beserta langkah XML penulisnya.
var PropAdjustmentKomite = map[string]string{
	"IsApproved":       "KomitePostAdjustment S15",
	"AcceptedNo":       "KomitePostAdjustment S16.9",
	"AcceptanceStatus": "KomitePostAdjustment S16.9 (1) dan S25 (2)",
	"AcceptedDate":     "KomitePostAdjustment S16.9",
	"IsKomite":         "KomitePostAdjustment S23 (subjectivity)",
	"IsSubjectivity":   "KomitePostAdjustment S24",
	"SubjectivityNote": "KomitePostAdjustment S24",
	"Notes":            "KomitePostAdjustment S27",
	"IsFacRetro":       "SaveAcceptationTreaty_TKMT S7 (KomitePostAdjustment S21)",
	"IsPrintAccept":    "SaveAcceptationTreaty_TKMT S9 (KomitePostAdjustment S21)",
	"IDOfBank":         "HitServiceToKasirKMT_Act S13 (KomitePostAdjustment S34)",
	"StatusKasir":      "HitServiceToKasirKMT_Act S14.5 (KomitePostAdjustment S34)",
}

var (
	// ErrKlaimTreatyTidakAda - kasus klaim induk tidak ada.
	ErrKlaimTreatyTidakAda = errors.New("kontrak: kasus klaim treaty induk tidak ada")
	// ErrKlaimTreatyTertutup - kasus klaim induk sudah ditutup.
	ErrKlaimTreatyTertutup = errors.New("kontrak: kasus klaim treaty induk sudah ditutup")
	// ErrAdjustmentTreatyTidakAda - baris adjustment `adjID` tidak ada di kasus klaim induk.
	ErrAdjustmentTreatyTidakAda = errors.New("kontrak: baris adjustment tidak ada di kasus klaim induk")
	// ErrUbahanKlaimTreatyTidakSah - ubahan di luar daftar putih.
	ErrUbahanKlaimTreatyTidakSah = errors.New("kontrak: ubahan klaim treaty di luar daftar putih Komite")
)
