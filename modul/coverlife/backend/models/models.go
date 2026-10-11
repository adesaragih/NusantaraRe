// Package models memuat bentuk data modul Cover Life (`coverlife`) - section Pega `InboxCoverLife` (kelas
// `ASM-FW-GISFW-Int-COVER_LIFE`, judul "COVER" b368; XML `D:\NUSARE DEV\Menu Cover\InboxCoverLife.xml`, nomor = baris
// XML). Satu tabel `M_COVER_LIFE` (ID, COVER, NOTE) - keputusan work owner 08-10-2026 C1: tabel Pega `M_COVER_LIFE`
// TETAP bernama sama (TANPA RENAME) dan menjadi flat (migrasi modul 085-086), JSONDATA dan view lamanya dibuang.
package models

// Cover - satu baris `M_COVER_LIFE` (grid `BrowseCoverLife_RD` b4026).
type Cover struct {
	// ID - kolom grid "ID" b3023 (sel salin-tempel dari `BrowseCauseofLossLife_RD` b3002) bernilai `.ID` b3449. Form
	// `InboxCoverLife` TIDAK punya medan ID (hanya Cover b848 dan Note b1026) - ID dibentuk server.
	ID string `json:"id"`
	// Cover - kolom grid "Cover" b3161 bernilai `.Cover` b3596; medan form "Cover" b848 (`.Cover` b877, pxTextInput
	// b880, wajib b843 / b895).
	Cover string `json:"cover"`
	// Note - medan form "Note" b1026 (`.Note` b1055, pxTextArea b1058, tidak wajib b1020 / b1070). TIDAK tampil di grid
	// (grid hanya ID, Cover, Edit), tetapi diisi `EditList_DT` (b3823) sehingga ikut dikirim API untuk Edit. Kosong =
	// NULL (4/4 baris DEV).
	Note string `json:"note"`
}

// Isian - isian form `InboxCoverLife` (Save b5407 -> `AddToList_Act` b5426). Tanpa ID: Add = dibentuk server
// (`BentukID`), Edit = ID dari jalur (Edit b3784 -> `EditList_DT` b3807 dengan ID b3811, Cover b3817, Note b3823).
type Isian struct {
	Cover string `json:"cover"`
	Note  string `json:"note"`
}

// Saringan - halaman grid. XML TANPA saring (`pyGridFiltering` false b4080) dan tanpa urut pilihan pengguna (ketiga
// kolom `pyColumnSorting` false b3969 / b3991 / b4013): satu-satunya masukan adalah halaman.
type Saringan struct {
	// Halaman - mulai 1.
	Halaman int
}

// UkuranHalaman - `pyPageSize` 50 b4101 (`pyPageMode` Numeric b4072, `pyGridPaginator` b2758 / b2780).
const UkuranHalaman = 50

// Halaman - satu halaman grid.
type Halaman struct {
	Daftar  []Cover `json:"daftar"`
	Total   int     `json:"total"`
	Halaman int     `json:"halaman"`
	Ukuran  int     `json:"ukuran"`
}

// BatasCover - panjang Cover (byte) = lebar `M_COVER_LIFE.COVER` VARCHAR2(200) (migrasi 085; C1.1). XML tidak
// menetapkan batas (pxTextInput tanpa pyMaxChars b880). Data DEV maks 17 byte (fakta WO 08-10-2026).
const BatasCover = 200

// BatasNote - panjang Note (byte) = lebar `M_COVER_LIFE.NOTE` VARCHAR2(1000) (migrasi 085; C1.1 "lebar dari XML atau
// 1000" - pxTextArea tanpa pyMaxChars b1058).
const BatasNote = 1000

// BatasID - `M_COVER_LIFE.ID` VARCHAR2(10) warisan (PK `SYS_C009203`).
const BatasID = 10
