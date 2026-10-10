// Package models memuat bentuk data modul Cause Of Loss Life (`causeoflosslife`) - section Pega `InboxCauseofLossLife`
// (kelas `ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE`, judul "CAUSE OF LOSS" b343; XML
// `D:\NUSARE DEV\Menu Cause Of Loss\InboxCauseofLossLife.xml`, nomor = baris XML). Satu tabel `CAUSEOFLOSS_LIFE` (ID,
// CAUSEOFLOSS) - keputusan work owner 08-10-2026 K1: tabel Pega `M_CAUSEOFLOSS_LIFE` berganti nama dan menjadi flat
// (migrasi modul 090-092), JSONDATA dibuang.
package models

// CauseOfLoss - satu baris `CAUSEOFLOSS_LIFE` (grid `BrowseCauseofLossLife_RD` b2968 / b3981).
type CauseOfLoss struct {
	// ID - kolom grid "ID" b2989 bernilai `.ID` b3415; medan form "ID" b823 (`.ID` b852, pxTextInput b855, disabled
	// b871-b872).
	ID string `json:"id"`
	// CauseOfLoss - kolom grid "Cause of Loss" b3127 bernilai `.CauseofLoss` b3562; medan form "Cause of Loss" b1001
	// (`.CauseofLoss` b1027, pxTextInput b1029, wajib b996 / b1044). Kosong = NULL (baris 100001 DEV).
	CauseOfLoss string `json:"causeOfLoss"`
}

// Isian - isian form `InboxCauseofLossLife` (Save b5366 -> `AddToList_Act` b5385). ID TIDAK diisi pengguna (disabled
// b872): Add = dibentuk server (`BentukID`), Edit = ID dari jalur (Edit b3747 -> `EditList_DT` b3770 dengan ID b3774
// dan CauseofLoss b3780).
type Isian struct {
	CauseOfLoss string `json:"causeOfLoss"`
}

// Saringan - halaman grid. XML TANPA saring (`pyGridFiltering` false b4036) dan tanpa urut pilihan pengguna (ketiga
// kolom `pyColumnSorting` false b3925 / b3947 / b3969): satu-satunya masukan adalah halaman.
type Saringan struct {
	// Halaman - mulai 1.
	Halaman int
}

// UkuranHalaman - `pyPageSize` 10 b4057 (`pyPageMode` Numeric b4028, `pyGridPaginator` b2724).
const UkuranHalaman = 10

// Halaman - satu halaman grid.
type Halaman struct {
	Daftar  []CauseOfLoss `json:"daftar"`
	Total   int           `json:"total"`
	Halaman int           `json:"halaman"`
	Ukuran  int           `json:"ukuran"`
}

// BatasCauseOfLoss - panjang Cause of Loss (byte) = lebar `CAUSEOFLOSS_LIFE.CAUSEOFLOSS` VARCHAR2(200) (migrasi 091;
// K1: pola Benefit 943 karena XML tidak menetapkan batas - pxTextInput tanpa pyMaxChars b1029). Data DEV maks 9 byte
// (fakta WO 08-10-2026).
const BatasCauseOfLoss = 200

// BatasID - `CAUSEOFLOSS_LIFE.ID` VARCHAR2(10) warisan (PK `SYS_C008825`).
const BatasID = 10
