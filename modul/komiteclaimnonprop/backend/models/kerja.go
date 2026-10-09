package models

// Untuk apa berkas ini: DAFTAR KERJA penyetuju - worklist assignment "KomiteRouter" (`Flow/KomiteTreaty_Flow.xml`,
// router `KomiteRouter` S6: `AssignTo` = operator baris tangga pertama ber-keputusan 0). Halaman awal modul; nol tombol
// "buat baru" (kasus lahir dari Claim Non Prop).
//
// Kolom mengikuti daftar kerja Komite Claim Prop (worklist `KomiteRouter` tidak punya section inbox sendiri di korpus),
// TANPA nilai / mata uang: baris akseptasi Non Prop tidak punya satu nilai (nilai per layer; `ValueAdjustment` tanpa
// penulis, OQ-CNP-37).

import "time"

// BarisKerja - satu kasus komite di daftar kerja penyetuju.
type BarisKerja struct {
	KasusID string `json:"kasusId"`
	// KlaimID - kasus klaim induk (CLMNP-, `pxCoverInsKey`); NoKlaim - `ClaimData.NoClaim` induknya.
	KlaimID string `json:"klaimId"`
	NoKlaim string `json:"noKlaim"`
	// Tingkat - `KOMITE_URUT` baris tangga yang menunggu; Count / Loop - `KomiteCount` / `KomiteLoop`.
	Tingkat int    `json:"tingkat"`
	Count   int    `json:"komiteCount"`
	Loop    int    `json:"komiteLoop"`
	Jabatan string `json:"jabatan"`
	// StatusBaris - `.AcceptanceStatus` baris adjustment; layanan memasang labelnya (`LabelStatusBaris`).
	StatusBaris string    `json:"statusBaris"`
	StatusWork  string    `json:"statusWork"`
	TglUpdate   time.Time `json:"tglUpdate"`
}
