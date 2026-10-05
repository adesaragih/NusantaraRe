// Package models memuat bentuk data modul Treaty Group - tabel warisan `POOLDATA.TREATYGROUP` (kelas Pega
// `ASM-FW-GISFW-Int-TREATYGROUP`, dibaca `BrowseTreatyGroup_RD`). Modul di luar korpus (perintah work owner 05-10-2026).
package models

// Grup - satu baris `POOLDATA.TREATYGROUP`, ditambah nama grup bisnis COA (`BUSINESSGROUP.NOTE`, dibaca saja).
type Grup struct {
	ID    string `json:"id"`
	OldID string `json:"oldId"`
	// OJK - salinan baris `TREATYGROUPOJK` saat disimpan (OJKBUSINESSID, OJKBUSINESSNAME, OJKBUSINESSNAMEIDN, ORDERNO).
	OjkID      string `json:"ojkId"`
	OjkName    string `json:"ojkName"`
	OjkNameIDN string `json:"ojkNameIdn"`
	OrderNo    string `json:"orderNo"`
	Name       string `json:"name"`
	SoaName    string `json:"soaName"`
	// TglUpdate - teks warisan format Pega `YYYYMMDDTHHMMSS.mmm GMT`; Diubah - tampilannya WIB `DD-MM-YYYY HH:MM`.
	TglUpdate string `json:"tglUpdate"`
	Diubah    string `json:"diubah"`
	UserID    string `json:"userId"`
	// COA - `BUSINESSGROUP.ID` grup bisnis utama OJK-nya (keputusan work owner 05-10-2026: opsi A, tidak diketik; ikut
	// OJK yang dipilih).
	CoaID   string `json:"coaId"`
	CoaName string `json:"coaName"`
}

// BisnisGrup - satu baris `BUSINESSGROUP` yang `TOPID`-nya grup ini (View; dibaca saja).
type BisnisGrup struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// Detail - satu grup dan grup bisnis anaknya (tanpa yang berakhiran SYARIAH).
type Detail struct {
	Grup
	Anak []BisnisGrup `json:"anak"`
}

// Ojk - satu baris `TREATYGROUPOJK` (pilihan OJK Business; dibaca saja) dan COA OJK itu: COAID yang paling sering di
// grup se-OJK, lalu terkecil (kosong = belum ada). Form menampilkannya begitu OJK dipilih (perintah work owner
// 05-10-2026: "ubah pas pilih OJK Business").
type Ojk struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	NameIDN string `json:"nameIdn"`
	OrderNo string `json:"orderNo"`
	CoaID   string `json:"coaId"`
	CoaName string `json:"coaName"`
}

// Isian - isian form Add / Edit. ID kosong = baris baru.
type Isian struct {
	ID      string `json:"id"`
	OjkID   string `json:"ojkId"`
	Name    string `json:"name"`
	SoaName string `json:"soaName"`
}

// Batas kolom (byte - katalog DEV 05-10-2026).
const (
	BatasTeks    = 1000 // OJKBUSINESS*, TREATYGROUPNAME, TREATYGROUPSOANAME
	BatasUserID  = 100
	BatasOrderNo = 10
)
