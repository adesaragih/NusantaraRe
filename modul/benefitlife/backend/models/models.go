// Package models memuat bentuk data modul Benefit (`benefitlife`) - section Pega `InboxBenefit` (kelas
// `ASM-FW-GISFW-Int-BENEFIT_LIFE`, judul "INSURANCE BENEFIT" b313; XML `D:\NUSARE DEV\Menu Benefit\InboxBenefit.xml`,
// nomor = baris XML). Satu tabel `BENEFIT_LIFE` (ID, BENEFIT) - keputusan work owner 08-10-2026 K1: tabel Pega
// `M_BENEFIT_LIFE` berganti nama dan menjadi flat (migrasi inti 942-944), JSONDATA dibuang.
package models

// Benefit - satu baris `BENEFIT_LIFE` (grid `BrowseBenefitLife_RD` b3411/b4451).
type Benefit struct {
	// ID - kolom grid "ID" b3431 bernilai `.Number` b3858; medan form "Number / ID" b964 (`.Number` b993, pxTextInput
	// disabled b1012-b1013). `.Number` = ID (fakta WO 08-10-2026: kunci JSON Number selalu sama dengan ID) - tidak ada
	// kolom NUMBER terpisah.
	ID string `json:"id"`
	// Benefit - kolom grid "Benefit" b3570 bernilai `.Benefit` b4020; medan form "Benefit" b1143 (`.Benefit` b1172,
	// pxTextArea b1175, wajib b1137/b1188).
	Benefit string `json:"benefit"`
}

// Isian - isian form `InboxBenefit` (Save b1723/b1772 -> `AddToList_Act` b1791). ID TIDAK diisi pengguna (disabled
// b1013): Add = dibentuk server (`BentukID`), Edit = ID dari jalur (Edit b4220 -> `EditList_DT` b4243 dengan Number,
// Benefit b4247-b4255).
type Isian struct {
	Benefit string `json:"benefit"`
}

// Saringan - filter dan arah urut grid (`pyGridFiltering` b4505, `pyGridSorting` b4493).
type Saringan struct {
	// ID - saring "memuat" atas ID.
	ID string
	// Benefit - saring "memuat" atas Benefit (tanpa beda huruf).
	Benefit string
	// Naik - ID menaik; bawaan (false) = ID MENURUN: `pySortType` DESC b4391, `pySortOrder` 1 b4397.
	Naik bool
	// Halaman - mulai 1.
	Halaman int
}

// UkuranHalaman - `pyPageSize` 10 b4526 (`pyPageMode` Numeric b4498, `pyGridPaginator` b3162).
const UkuranHalaman = 10

// Halaman - satu halaman grid.
type Halaman struct {
	Daftar  []Benefit `json:"daftar"`
	Total   int       `json:"total"`
	Halaman int       `json:"halaman"`
	Ukuran  int       `json:"ukuran"`
}

// BatasBenefit - panjang Benefit (byte) = lebar `BENEFIT_LIFE.BENEFIT` VARCHAR2(200) (migrasi inti 943; K1: pola
// ririsklife `USEDBY` 200 karena XML tidak menetapkan batas - pxTextArea tanpa pyMaxChars b1175). Data DEV maks 48
// byte (fakta WO 08-10-2026).
const BatasBenefit = 200

// BatasID - `BENEFIT_LIFE.ID` VARCHAR2(10) warisan (PK `SYS_C009031`).
const BatasID = 10
