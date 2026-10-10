package models

// Untuk apa berkas ini: BARIS TABEL WARISAN yang ditulis Komite Claim Fac In - `HISTORYAKSEPTASIPEGA`
// (KomitePost_Adjustment S20-S21), `MONITORING_KLAIM_LOG` (S18-S19 "AKSEPATSI"), `CLAIMREJECTED` (KomitePost_Reject
// S14), dan bahan nomor akseptasi (S7.2.1.2.4-S7.2.1.2.7).

import "time"

// Nilai tetap `InsertHistoryAkseptasiPega_Sql` (S20).
const (
	// WorkbasketRiwayat - `InsertHistory.CARI2 := "KLAIM"`.
	WorkbasketRiwayat = "KLAIM"
	StatusRiwayatACC  = "ACCEPT"
	StatusRiwayatREJ  = "REJECT"
)

// Nilai tetap S18 / S22 (VERBATIM; "AKSEPATSI" dibaca sistem lain - ejaan tidak diperbaiki, prompt §5).
const (
	JenisLogAkseptasi  = "AKSEPATSI"
	SubProgresDiterima = "Accepted"
	SubProgresDitolak  = "Rejected"
)

// KunciInstans - padanan `pzInsKey` kasus sistem baru di tabel warisan (OS_AKSEPTASI_KLAIM.CASEID, JSON_KLAIM.IDPEGA,
// HISTORYAKSEPTASIPEGA.ID_PEGA / ID_KOMITE, MONITORING_KLAIM_LOG.IDPEGA, CLAIMREJECTED.INSKEY, SUBPROGRESSCLAIM.IDPEGA):
// ID T_WORK_CLAIM apa adanya - sama dengan Claim Fac In `models.KunciInstans` (prompt tahap 1 §6 butir 6).
func KunciInstans(id string) string { return id }

// RiwayatAkseptasi - satu baris HISTORYAKSEPTASIPEGA. `OPERATORID` tidak disebut rule-nya (tetap NULL).
type RiwayatAkseptasi struct {
	// IDPega - `TempOpenPage.pzInsKey` (klaim induk); IDKomite - `pyWorkPage.pzInsKey` (kasus komite).
	IDPega, IDKomite string
	// Status - CARI5 ACCEPT / REJECT / kode lain apa adanya; Username - `OperatorID.pyUserName`.
	Status, Username, Workbasket string
}

// LogLayanan - satu baris MONITORING_KLAIM_LOG (`InsertLogServiceClaim`).
type LogLayanan struct {
	IDPega, Parameter, JenisService, NoAkseptasi, NoDLA, StsMessage, ResponMessage string
}

// KlaimDitolak - satu baris CLAIMREJECTED (`InsertClaimRejected_Sql`, KomitePost_Reject S14.1).
type KlaimDitolak struct {
	// InsKey / ID / InsName / Label / StatusWork / Kelas - CARI1-CARI5, CARI8 kasus klaim induk.
	InsKey, ID, InsName, Label, StatusWork, Kelas string
	// PembuatNama / PembuatID - CARI6 / CARI7 (`pxCreateOpName` / `pxCreateOperator` klaim induk).
	PembuatNama, PembuatID string
	// Diperbarui - CARI9 `pyWorkPage.pxUpdateDateTime` (kasus komite); PengubahNama / PengubahID - CARI10 / CARI11.
	Diperbarui               time.Time
	PengubahNama, PengubahID string
	// Remark - CARI12 `pyWorkPage.Komite.Remarks`.
	Remark string
}

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut) dan jenisnya
// (`ParamSeq.CARI2` = kode produksi NONLIFE + "A").
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}
