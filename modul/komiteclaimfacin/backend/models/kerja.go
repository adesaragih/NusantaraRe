package models

// Untuk apa berkas ini: DAFTAR KERJA penyetuju - worklist assignment "KomiteRouter" (`Flow/Komite_Flow.xml`, assignment
// WorkList Custom): kasus `KMT-` yang menunggu akun pelaku atau workbasket aktifnya di tingkat berjalan. Tampil sebagai
// tabel komite di bawah inbox Claim Fac In (perintah work owner 09-10-2026, commit dcd1522e: modul tanpa menu). Nol
// tombol "buat baru" (kasus lahir dari Claim Fac In).
//
// Kolom mengikuti daftar kerja Komite Claim Prop / Non Prop (worklist `KomiteRouter` tidak punya section inbox sendiri di
// korpus) ditambah jenis penyerahan (Section ShowTransfer LS1: ADJUSTMENT / REJECT / CLOSE).

import "time"

// BarisKerja - satu kasus komite di daftar kerja penyetuju.
type BarisKerja struct {
	KasusID string `json:"kasusId"`
	// KlaimID - kasus klaim induk (CLM-, `pxCoverInsKey`); NoKlaim - `ClaimData.NoClaim` induknya.
	KlaimID string `json:"klaimId"`
	NoKlaim string `json:"noKlaim"`
	// Jenis - label TransferType (ADJUSTMENT / REJECT / CLOSE).
	Jenis string `json:"jenis"`
	// Tingkat - `KOMITE_URUT` baris tangga yang menunggu; Count / Loop - `KomiteCount` / `KomiteLoop`.
	Tingkat int    `json:"tingkat"`
	Count   int    `json:"komiteCount"`
	Loop    int    `json:"komiteLoop"`
	Jabatan string `json:"jabatan"`
	// StatusBaris - `.AcceptanceStatus` baris adjustment TT2 (kosong TT3 / TT4); layanan memasang labelnya.
	StatusBaris string    `json:"statusBaris"`
	StatusWork  string    `json:"statusWork"`
	TglUpdate   time.Time `json:"tglUpdate"`
}
