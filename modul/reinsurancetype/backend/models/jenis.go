// Package models memuat bentuk data modul Reinsurance Type - tabel warisan `POOLDATA.REINSURANCETYPE` (kelas Pega
// `ASM-FW-GISFW-Int-REINSURANCETYPE`, dibaca `BrowseReinsuranceType_RD`): master jenis reasuransi, bersama life dan
// non-life. Modul di luar korpus (perintah work owner 05-10-2026).
package models

// Jenis - satu baris `POOLDATA.REINSURANCETYPE`.
type Jenis struct {
	ID string `json:"id"`
	// Name - `NOTE`; dicari Pega lewat nama (`GetReinsuranceTypeBYName_SQL`).
	Name string `json:"name"`
	// Type - `1` Own Retention, `2` Treaty Out, `3` Facultative, `4` Treaty In (Prompt value Pega).
	Type      string `json:"type"`
	SoaName   string `json:"soaName"`
	Code      string `json:"code"`
	Flag      string `json:"flag"`
	NoUrut    string `json:"noUrut"`
	GroupType string `json:"groupType"`
	UserID    string `json:"userId"`
	// TglUpdate - teks format Pega `YYYYMMDDTHHMMSS.mmm GMT`; Diubah - tampilannya WIB `DD-MM-YYYY HH:MM`.
	TglUpdate string `json:"tglUpdate"`
	Diubah    string `json:"diubah"`
}

// Isian - isian form Add / Edit. ID kosong = baris baru.
type Isian struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	SoaName   string `json:"soaName"`
	Code      string `json:"code"`
	Flag      string `json:"flag"`
	NoUrut    string `json:"noUrut"`
	GroupType string `json:"groupType"`
}

// BatasTeks - seluruh kolom REINSURANCETYPE VARCHAR2(100) (katalog DEV 05-10-2026).
const BatasTeks = 100

// Nilai sah untuk baris baru / nilai yang diganti. Nilai warisan lain (Flag `1` dan kosong, Type kosong) dibiarkan
// selama tidak diubah.
var (
	Type      = []string{"1", "2", "3", "4"}
	Flag      = []string{FlagAktif, FlagNonaktif}
	GroupType = []string{"OR", "QS", "RI", "SPL"}
)

// Flag aktif / nonaktif - teks persis yang disaring pembaca Pega (`FLAG ='active'`).
const (
	FlagAktif    = "active"
	FlagNonaktif = "inactive"
	// CodeBawaan - CODE hampir seluruh baris DEV (TYPE 4 memakai kode 2 digit 10-51).
	CodeBawaan = "00"
)
