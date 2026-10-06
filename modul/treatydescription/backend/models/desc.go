// Package models memuat bentuk data modul Treaty Description - tabel warisan `POOLDATA.TREATYDESC` (kelas Pega
// `ASM-FW-GISFW-Int-TREATYDESC`, dibaca `BrowseTreatyDesc_RD`, ditulis prosedur `PEGA_TREATYDESC`). Modul di luar
// korpus (perintah work owner 05-10-2026).
package models

// Nilai kolom ISXOL dan STATUSAKTIF (semuanya VARCHAR2 di TREATYDESC).
const (
	NonXOL   = "0"
	XOL      = "1"
	Aktif    = "1"
	Nonaktif = "0"
)

// Desc - satu baris `TREATYDESC`.
type Desc struct {
	ID       string `json:"id"`
	DescName string `json:"descName"`
	// IsXOL - "0" Non XOL, "1" XOL; NULL dibaca "0" (seluruh baris DEV "0").
	IsXOL string `json:"isXol"`
	// StatusAktif - teks apa adanya: "1", "0", atau "" (NULL - baris lama Pega, 6 dari 13 di DEV). Aktif = bukan "0".
	StatusAktif string `json:"statusAktif"`
	Aktif       bool   `json:"aktif"`
}

// Isian - isian form Add / Edit. ID kosong = baris baru (diisi dari jalur, bukan badan).
type Isian struct {
	ID          string `json:"id"`
	DescName    string `json:"descName"`
	IsXOL       string `json:"isXol"`
	StatusAktif string `json:"statusAktif"`
}

// BatasNama - `DESCNAME VARCHAR2(100)` (katalog DEV 05-10-2026), dalam byte.
const BatasNama = 100
