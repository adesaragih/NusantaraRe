package models

// Untuk apa berkas ini: BARIS TABEL WARISAN yang ditulis Komite - `HISTORYAKSEPTASIPEGA` (KomitePostAdjustment
// S22-S23) dan bahan nomor akseptasi (S14.8-S14.11). Non Prop tanpa MONITORING_KLAIM_LOG (PINDAI §7.5).

// Nilai tetap `InsertHistoryAkseptasiPega_Sql` (S22).
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

// BahanNomor - keluaran penghitung bersama (`ParamSeq.HASIL1` MM.YYYY, `ParamSeq.HASIL2` urut) dan jenisnya
// (`ParamSeq.CARI2` = kode produksi NONLIFE + "A").
type BahanNomor struct {
	Jenis  string
	MMYYYY string
	Urut   int
}
