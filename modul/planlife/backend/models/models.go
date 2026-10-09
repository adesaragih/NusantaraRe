// Package models memuat bentuk data modul Plan (`planlife`) - section Pega `InboxProductType` (kelas
// `ASM-FW-GISFW-Int-PRODUCT_TYPE_LIFE`, judul "Plan" b331 / "Product Type" b644; XML
// `D:\NUSARE DEV\Menu Plan\InboxProductType.xml`, nomor = baris XML). Satu tabel `PRODUCT_TYPE_LIFE` - keputusan work
// owner 08-10-2026 K1: tabel Pega `M_PRODUCT_TYPE_LIFE` berganti nama dan menjadi flat (migrasi inti 946-948).
package models

// Plan - satu baris `PRODUCT_TYPE_LIFE` (grid `BrowseProductTypeLife_RD` b4325; urutan kolom = view lama).
type Plan struct {
	// ID - `.ID` (parameter `EditProductTypeLife_Act` b4060); '1' || LPAD(seq, 5, '0'). Tidak tampil di form / grid.
	ID string `json:"id"`
	// CoverName - "Plan Name" b812/b834 (`.CoverName` b841, pxTextInput b844); kolom grid "Plan Name" b2978 / b3543.
	CoverName string `json:"coverName"`
	// Business - "Business" b1079 (`.Business` b1105, pxAutoComplete b1107) = `BUSINESS.NOTE` (b1187/b1231).
	Business string `json:"business"`
	// BusinessID - `BUSINESS.ID` (`pyPropertyTarget` .BusinessID b1268).
	BusinessID string `json:"businessId"`
	// Benefit - "Benefit" b1432/b1455 (`.Benefit` b1461, pxAutoComplete b1464) = `BENEFIT_LIFE.BENEFIT` (b1541/b1552).
	Benefit string `json:"benefit"`
	// BenefitID - `BENEFIT_LIFE.ID` (`.Number` -> .BenefitID b1563-b1564).
	BenefitID string `json:"benefitId"`
}

// Isian - isian form (Save b6403 -> `SaveProductTypeLife_Act` b6422). Hanya TEKS: ID business / benefit diisi server
// dari pencocokan ulang ke master (K4); ID plan dari jalur (Edit) atau sequence (K3).
type Isian struct {
	CoverName string `json:"coverName"`
	Business  string `json:"business"`
	Benefit   string `json:"benefit"`
}

// PilihanBusiness - satu baris autocomplete Business (`BrowseBusinessLife_RD` b1183): kolom OLDID b1200, Note b1231,
// ID b1265.
type PilihanBusiness struct {
	ID    string `json:"id"`
	OldID string `json:"oldId"`
	Note  string `json:"note"`
}

// PilihanBenefit - satu baris autocomplete Benefit (`BrowseBenefitLife_RD` b1537): Benefit b1552, Number (= ID) b1563.
type PilihanBenefit struct {
	ID      string `json:"id"`
	Benefit string `json:"benefit"`
}

// Kolom urut grid (`pyGridSorting` true b4368): Plan Name dan Benefit dapat diurutkan (`pyColumnSorting` true b4247 /
// b4291), Business tidak (false b4269). Tanpa urutan bawaan di XML (`pySortType` NONE, `pyMaxSortOrder` 0 b2890).
const (
	UrutCoverName = "covername"
	UrutBenefit   = "benefit"
)

// KolomUrut - kolom yang boleh diurutkan; kosong = ID menaik (stabil; isi RD tidak ada di XML).
var KolomUrut = []string{UrutCoverName, UrutBenefit}

// Saringan - urutan dan halaman grid (tanpa saring: `pyGridFiltering` false b4380).
type Saringan struct {
	Urut    string
	Turun   bool
	Halaman int
}

// UkuranHalaman - `pyPageSize` 10 b4401 (`pyPageMode` Numeric b4373, `pyGridPaginator` b2709).
const UkuranHalaman = 10

// Halaman - satu halaman grid.
type Halaman struct {
	Daftar  []Plan `json:"daftar"`
	Total   int    `json:"total"`
	Halaman int    `json:"halaman"`
	Ukuran  int    `json:"ukuran"`
}

// Lebar kolom (migrasi inti 947; K1): nama 200 (pola benefitlife), ID master 10 (BUSINESS.ID / BENEFIT_LIFE.ID), ID
// plan VARCHAR2(6) warisan.
const (
	BatasNama     = 200
	BatasIDMaster = 10
	BatasID       = 6
)

// GrupBusinessLife - K4: `BUSINESS.GROUPPANEL` pilihan Business. DITURUNKAN DARI DATA, bukan XML (isi
// `BrowseBusinessLife_RD` tidak ada; hanya nama parameter `Group` b1310): ke-21 business yang dipakai plan DEV semuanya
// di grup ini, dan grup ini persis 21 baris LIFE INSURANCE (fakta WO 08-10-2026).
const GrupBusinessLife = "009"
