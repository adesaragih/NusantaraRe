// Package models memuat bentuk data modul Adjuster Consultant - kelas Pega `ASM-FW-GISFW-Int-ADJUSTERCONSULTANT`
// (layar master `MstAdjusterConsultant`, folder korpus Claim Fac In dan Claim Prop; keputusan work owner 05-10-2026).
package models

// Adjuster - satu baris `POOLDATA.ADJUSTERCONSULTANT`.
type Adjuster struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	TelpNo  string `json:"telpNo"`
	// Username - akun pengubah terakhir (`USERNAME`); EditDate - waktu ubah terakhir (`EDITDATE`, `DD-MM-YYYY HH24:MI`).
	Username string `json:"username"`
	EditDate string `json:"editDate"`
	// Active - `IS_ACTIVE` (migrasi 870, keputusan work owner 05-10-2026: flag nonaktif pengganti hapus).
	Active bool `json:"active"`
}

// Isian - isian form Add / Edit. ID kosong = baris baru.
type Isian struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	TelpNo  string `json:"telpNo"`
}

// Status saringan daftar.
const (
	StatusSemua    = ""
	StatusAktif    = "active"
	StatusNonaktif = "inactive"
)

// Nilai `IS_ACTIVE`.
const (
	BenderaAktif    = "1"
	BenderaNonaktif = "0"
)

// Batas kolom (byte - semantik VARCHAR2 bawaan; katalog DEV 05-10-2026).
const (
	BatasNama     = 500
	BatasAlamat   = 1000
	BatasTelp     = 50
	BatasUsername = 50
)
