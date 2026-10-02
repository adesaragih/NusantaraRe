package repository

// Tabel FLAT produk - keputusan work owner 02-10-2026 (K2, K5; tiket 01 bab bertanggal 02-10-2026): satu produk = satu
// baris induk `M_PRODUCTNAME_LIFE` (sisi umum + sisi inward digabung, grilling Q1b) + tujuh tabel anak berawalan
// `M_PRODUCTNAME_LIFE_` (grilling D2). DDL: `migrations/140`–`147`; kolom dan tipe: `docs/STRUKTUR-TABEL-…md`.

// Nama tabel flat.
const (
	TabelFlatInduk    = "M_PRODUCTNAME_LIFE"
	TabelFlatLien     = "M_PRODUCTNAME_LIFE_LIEN"
	TabelFlatDokumen  = "M_PRODUCTNAME_LIFE_DOCCLAIM"
	TabelFlatPlan     = "M_PRODUCTNAME_LIFE_PLAN"
	TabelFlatFinUW    = "M_PRODUCTNAME_LIFE_FINUW"
	TabelFlatUWLimit  = "M_PRODUCTNAME_LIFE_UWLIMIT"
	TabelFlatOutward  = "M_PRODUCTNAME_LIFE_OUTWARD"
	TabelFlatKomentar = "M_PRODUCTNAME_LIFE_COMMENT"
)

// DaftarTabelFlat - induk lalu ketujuh anak, urutan migrasi 140–147.
var DaftarTabelFlat = []string{TabelFlatInduk, TabelFlatLien, TabelFlatDokumen, TabelFlatPlan, TabelFlatFinUW,
	TabelFlatUWLimit, TabelFlatOutward, TabelFlatKomentar}
