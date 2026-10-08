package models

// Untuk apa berkas ini: DAFTAR KERJA penyetuju - worklist assignment "KomiteRouter" (`Flow/KomiteTreaty_Flow.xml`,
// router `KomiteRouter` S6: `AssignTo` = operator baris tangga pertama ber-keputusan 0). Halaman awal modul; nol tombol
// "buat baru" (`pyCanCreateWorkObject=false`, kasus lahir dari Claim Prop).
//
// Kolom mengikuti konvensi daftar kerja Komite Claim Life (prompt §8): worklist `KomiteRouter` tidak punya section
// inbox sendiri di korpus.

import "time"

// BarisKerja - satu kasus komite di daftar kerja penyetuju.
type BarisKerja struct {
	KasusID string `json:"kasusId"`
	// KlaimID - kasus klaim induk (CLMP-, `pxCoverInsKey`); NoKlaim - `ClaimData.NoClaim` induknya.
	KlaimID string `json:"klaimId"`
	NoKlaim string `json:"noKlaim"`
	// Tingkat - `KOMITE_URUT` baris tangga yang menunggu; Count / Loop - `KomiteCount` / `KomiteLoop`.
	Tingkat int    `json:"tingkat"`
	Count   int    `json:"komiteCount"`
	Loop    int    `json:"komiteLoop"`
	Jabatan string `json:"jabatan"`
	// Nilai / MataUang - "Adjustment RNM" baris adjustment yang diputus (uang TEKS, nol float).
	Nilai    string `json:"nilai"`
	MataUang string `json:"mataUang"`
	// StatusBaris - `.AcceptanceStatus` baris adjustment (kode apa adanya).
	StatusBaris string    `json:"statusBaris"`
	StatusWork  string    `json:"statusWork"`
	TglUpdate   time.Time `json:"tglUpdate"`
}
