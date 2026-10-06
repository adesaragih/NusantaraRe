// Package models memuat bentuk data modul Treaty Group OJK - tabel warisan `POOLDATA.TREATYGROUPOJK`, lini bisnis OJK
// induk setiap Treaty Group (`TREATYGROUP.OJKBUSINESSID`). Modul di luar korpus (perintah work owner 05-10-2026).
package models

// Ojk - satu baris `POOLDATA.TREATYGROUPOJK`.
type Ojk struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	NameIDN string `json:"nameIdn"`
	OrderNo string `json:"orderNo"`
}

// Isian - isian form Add / Edit. ID kosong = baris baru. Order No bukan isian (perintah work owner 05-10-2026:
// "ORDERNO hide aja, isi sesuai max dari order no") - Add mengisinya, Edit membiarkannya.
type Isian struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	NameIDN string `json:"nameIdn"`
}

// Batas kolom (byte - semantik VARCHAR2 bawaan; katalog DEV 05-10-2026: keempat kolom VARCHAR2(50)).
const (
	BatasID    = 50
	BatasNama  = 50
	BatasOrder = 50
)
