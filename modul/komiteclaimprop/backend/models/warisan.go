package models

// Untuk apa berkas ini: BARIS TABEL WARISAN yang ditulis Komite - `HISTORYAKSEPTASIPEGA` (KomitePostAdjustment
// S32-S33), `MONITORING_KLAIM_LOG` (S30-S31), dan bahan nomor akseptasi (S16.5-S16.7).

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

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut) dan jenisnya
// (`ParamSeq.CARI2` = kode produksi NONLIFE + "A").
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}
