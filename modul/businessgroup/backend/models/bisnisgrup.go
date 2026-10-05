// Package models memuat bentuk data modul Business Group - tabel warisan `POOLDATA.BUSINESSGROUP` (kelas Pega
// `ASM-FW-GISFW-Int-BUSINESSGROUP`, dibaca `BrowseBusinessGroup_RD`). Modul di luar korpus (perintah work owner
// 05-10-2026).
package models

// BisnisGrup - satu baris `POOLDATA.BUSINESSGROUP`.
type BisnisGrup struct {
	ID string `json:"id"`
	// Name - `NOTE` (nama grup bisnis).
	Name  string `json:"name"`
	Alias string `json:"alias"`
	// TopID - `TOPID` = `TREATYGROUP.ID` induknya; TreatyName - salinan `TREATYGROUPNAME` induk saat disimpan.
	TopID      string `json:"topId"`
	TreatyName string `json:"treatyName"`
}

// TreatyGroup - satu baris `TREATYGROUP` (pilihan Treaty Group; dibaca saja).
type TreatyGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Isian - isian form Add / Edit. ID kosong = baris baru.
type Isian struct {
	ID    string `json:"id"`
	TopID string `json:"topId"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// BatasTeks - kelima kolom BUSINESSGROUP VARCHAR2(4000) (katalog DEV 05-10-2026).
const BatasTeks = 4000

// AkhiranSyariah - grup bisnis berakhiran ini tidak dikelola modul ini (saringan Pega `BrowseBusinessGroup_RD`
// `.Note NotEndsWith "SYARIAH"`; keputusan work owner 05-10-2026: "tidak ada syariah").
const AkhiranSyariah = "SYARIAH"
