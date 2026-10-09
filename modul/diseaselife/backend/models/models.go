// Package models memuat bentuk data modul Disease Life (`diseaselife`) - section Pega `InboxDisease` (kelas
// `ASM-FW-GISFW-Int-DISEASE_LIFE`, judul "DISEASE" b371 / `pyCaption DISEASE` b6223; XML
// `D:\NUSARE DEV\Menu Disease\InboxDisease.xml`, nomor = baris XML). Satu tabel `DISEASE_LIFE` (ID, ICD_CODE, DISEASE)
// yang SUDAH flat - keputusan work owner 08-10-2026 D1: tabel tidak di-RENAME dan kolomnya tidak diubah; migrasi modul
// hanya menambah sequence `SEQ_DISEASE_LIFE` (080) dan PK `PK_DISEASE_LIFE` (081).
package models

// Penyakit - satu baris `DISEASE_LIFE` (grid `BrowseDiseaseLife_RD` b5094).
type Penyakit struct {
	// ID - kolom grid "ID" b3758 (sel salin-tempel dari `BrowseBenefitLife_RD` b3738) bernilai `.Number` b4321; medan
	// form "Number / ID" b1021 (`DISEASE_LIFE.Number` b1050, pxTextInput b1053, disabled b1069-b1070). `.Number` = kolom
	// ID (prosedur PEGA_DISEASE_LIFE menulis ID) - tidak ada kolom NUMBER terpisah.
	ID string `json:"id"`
	// ICDCode - kolom grid "ICD Code" b3894 bernilai `.ICD_Code` b4484; medan form "ICD Code" b1200 (`.ICD_Code` b1229,
	// pxTextInput b1232).
	ICDCode string `json:"icdCode"`
	// Disease - kolom grid "Disease" b4033 bernilai `.Disease` b4629; medan form "Disease" b1470 (`.Disease` b1499,
	// pxTextArea b1502, wajib b1465 / b1515).
	Disease string `json:"disease"`
}

// Isian - isian form `InboxDisease` (Save b2102 -> `AddToList_Act` b2121). ID TIDAK diisi pengguna (disabled b1070):
// Add = dibentuk server (`BentukID`, sequence SEQ_DISEASE_LIFE), Edit = ID dari jalur (Edit b4829 -> `EditList_DT`
// b4852 dengan Number b4856, Disease b4862, ICD_Code b4868).
type Isian struct {
	ICDCode string `json:"icdCode"`
	Disease string `json:"disease"`
}

// Kolom urut grid yang dapat dipilih pengguna: ID (`pyColumnSorting` true b5014) dan ICD Code (b5036). Disease TIDAK
// dapat diurutkan (b5058), kolom Edit juga tidak (b5080).
const (
	UrutID  = "id"
	UrutICD = "icd"
)

// Saringan - saring, urut, dan halaman grid (`pyGridFiltering` true b5148, `pyGridSorting` true b5138). 97.586 baris
// DEV: saring dan halaman SELALU di server (kebutuhan teknis, PARITAS).
type Saringan struct {
	// ICDCode - saring "memuat" atas ICD Code (tanpa beda huruf).
	ICDCode string
	// Disease - saring "memuat" atas Disease (tanpa beda huruf).
	Disease string
	// Urut - UrutID (bawaan) atau UrutICD; nilai lain = UrutID.
	Urut string
	// Naik - arah menaik. Bawaan (false) = MENURUN: kolom ID `pySortType` DESC b5012, `pySortOrder` 1 b5017.
	Naik bool
	// Halaman - mulai 1.
	Halaman int
}

// UkuranHalaman - `pyPageSize` 10 b5169 (`pyPageMode` Numeric b5140, `pyGridPaginator` b3489 / b3511).
const UkuranHalaman = 10

// Halaman - satu halaman grid.
type Halaman struct {
	Daftar  []Penyakit `json:"daftar"`
	Total   int        `json:"total"`
	Halaman int        `json:"halaman"`
	Ukuran  int        `json:"ukuran"`
}

// BatasID - `DISEASE_LIFE.ID` VARCHAR2(100) warisan (fakta WO 08-10-2026).
const BatasID = 100

// BatasICDCode - `DISEASE_LIFE.ICD_CODE` VARCHAR2(100) warisan (data DEV maks 7 byte). XML tanpa batas panjang
// (pxTextInput b1232 tanpa pyMaxChars) - batas = lebar kolom.
const BatasICDCode = 100

// BatasDisease - `DISEASE_LIFE.DISEASE` VARCHAR2(1000) warisan (data DEV maks 290 byte). XML tanpa batas panjang
// (pxTextArea b1502 tanpa pyMaxChars) - batas = lebar kolom.
const BatasDisease = 1000

// BatasSaring - panjang kata saring (byte) yang diteruskan ke SQL; lebih panjang dari kolom terlebar tidak mungkin
// cocok.
const BatasSaring = BatasDisease
