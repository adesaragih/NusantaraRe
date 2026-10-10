package models

// Untuk apa berkas ini: BARIS TABEL WARISAN yang ditulis Komite - `HISTORYAKSEPTASIPEGA` (KomitePostAdjustment
// S32-S33), `MONITORING_KLAIM_LOG` (S30-S31), `CLAIMREJECTED` (KomitePost_Close S13), dan bahan nomor akseptasi
// (S16.5-S16.7).

import "time"

// Nilai tetap `InsertHistoryAkseptasiPega_Sql` (S32).
const (
	// WorkbasketRiwayat - `InsertHistory.CARI2 := "KLAIM"`.
	WorkbasketRiwayat = "KLAIM"
	StatusRiwayatACC  = "ACCEPT"
	StatusRiwayatREJ  = "REJECT"
)

// RiwayatAkseptasi - satu baris HISTORYAKSEPTASIPEGA. `OPERATORID` tidak disebut rule-nya (tetap NULL).
type RiwayatAkseptasi struct {
	// IDPega - `TempOpenPage.pzInsKey` (klaim induk); IDKomite - `pyWorkPage.pzInsKey` (kasus komite).
	IDPega, IDKomite string
	// Status - CARI5 ACCEPT / REJECT / kode lain apa adanya; Username - `OperatorID.pyUserName`.
	Status, Username, Workbasket string
}

// JenisServiceAkseptasi - S30 `Param.JenisService := "AKSEPTASI"`.
const JenisServiceAkseptasi = "AKSEPTASI"

// LogLayanan - satu baris MONITORING_KLAIM_LOG (`InsertLogServiceClaim`).
type LogLayanan struct {
	IDPega, Parameter, JenisService, NoAkseptasi, NoDLA, StsMessage, ResponMessage string
}

// SusunLogAkseptasi = S30: `pzInsKey + " / " + PolicyNo + " / " + stsReject`; NoAkseptasi = `OutputData.START_DATE`
// (ditulis S16.8 = nomor akseptasi; kosong bila S16 tidak berjalan); NoDLA ""; StsMessage / ResponMessage = jawaban
// Connect-REST
// `KonversiKlaimNonLife` S29 - di sistem baru panggilan itu efek outbox asinkron, jadi kosong di sini (pola Claim
// Prop). `stsReject` = `TempOpenPage.stsReject` (SaveAcceptation_Act S1 "1"; tidak ditulis bila subjectivity).
func SusunLogAkseptasi(klaimID, noPolis, stsReject, nomor string) LogLayanan {
	return LogLayanan{IDPega: KunciInstans(klaimID), Parameter: KunciInstans(klaimID) + " / " + noPolis + " / " + stsReject,
		JenisService: JenisServiceAkseptasi, NoAkseptasi: nomor}
}

// Nilai tetap baris CLAIMREJECTED kasus klaim treaty (KomitePost_Close S13.1, `DataIn.CARI4` / `CARI5` / `CARI8` =
// `TempOpenPage.pyLabel` / `pyStatusWork` / `pxObjClass`): bentuk baris CLMP warisan di DEV (237 baris, 10-10-2026) -
// LABEL "ClaimTreaty", STATUSWORK klaim saat diputus ("New" untuk jalur close: klaim masih terbuka).
const (
	LabelKlaimTreaty   = "ClaimTreaty"
	StatusKlaimTerbuka = "New"
)

// KlaimDitolak - satu baris CLAIMREJECTED (`InsertClaimRejected_Sql`, KomitePost_Close S13 bila disetujui).
type KlaimDitolak struct {
	// InsKey / ID / InsName / Label / StatusWork / Kelas - CARI1-CARI5, CARI8 kasus klaim induk.
	InsKey, ID, InsName, Label, StatusWork, Kelas string
	// PembuatNama / PembuatID - CARI6 / CARI7 (`pxCreateOpName` / `pxCreateOperator` klaim induk).
	PembuatNama, PembuatID string
	// Diperbarui - CARI9 `pyWorkPage.pxUpdateDateTime` (kasus komite); PengubahNama / PengubahID - CARI10 / CARI11
	// (`TempOpenPage.pxUpdateOpName` / `pxUpdateOperator`: kasus klaim baru saja disimpan Komite - penyetuju).
	Diperbarui               time.Time
	PengubahNama, PengubahID string
	// Remark - CARI12 `pyWorkPage.Komite.Remarks` (= Remarks pop-up close, `ClaimData.Remark_Close` klaim induk).
	Remark string
}

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut) dan jenisnya
// (`ParamSeq.CARI2` = kode produksi NONLIFE + "A").
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}
