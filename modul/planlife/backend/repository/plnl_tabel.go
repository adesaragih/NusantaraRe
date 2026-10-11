// Package repository adalah SATU-SATUNYA lapisan modul Plan yang berbicara ke Oracle. Setiap nilai diikat (`:n`,
// go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services (pola
// benefitlife; prosedur PEGA_M_PRODUCT_TYPE_LIFE tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// Keputusan work owner 08-10-2026 K1: SATU tabel `PRODUCT_TYPE_LIFE`. Tabel Pega `M_PRODUCT_TYPE_LIFE` BERGANTI NAMA
// menjadi `PRODUCT_TYPE_LIFE` (nama view lama yang dibuang) dan menjadi tabel flat tanpa JSONDATA - migrasi inti
// 946-948. Nol JSON di modul ini.
package repository

// Tabel - `PRODUCT_TYPE_LIFE`, kelas `ASM-FW-GISFW-Int-PRODUCT_TYPE_LIFE`; kolom = KolomTabel. SATU-SATUNYA objek
// yang DITULIS.
const Tabel = "PRODUCT_TYPE_LIFE"

// Seq - `M_PRODUCT_TYPE_LIFE_SEQ` warisan (last_number 44 DEV). K3: TETAP, nol sequence baru.
const Seq = "M_PRODUCT_TYPE_LIFE_SEQ"

// Master yang DIBACA saja (K4) - TIDAK pernah ditulis.
const (
	// MasterBusiness - `BUSINESS` (kelas `ASM-FW-GISFW-Int-BUSINESS`, `BrowseBusinessLife_RD` b1183-b1189).
	MasterBusiness = "BUSINESS"
	// MasterBenefit - `BENEFIT_LIFE` (modul benefitlife; `BrowseBenefitLife_RD` b1537-b1543).
	MasterBenefit = "BENEFIT_LIFE"
)

// KolomTabel - kolom `PRODUCT_TYPE_LIFE` sesudah 948 = kolom view lama, urutan sama (ID warisan + 947).
var KolomTabel = []string{"ID", "COVERNAME", "BUSINESS", "BUSINESSID", "BENEFIT", "BENEFITID"}

// Kolom yang dibaca (nol `SELECT *`).
const kolomBaca = `ID, COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID`

// DaftarTabelDitulis - penjaga modul: hanya tabel ini yang menerima INSERT/UPDATE (nol DELETE: XML tanpa Delete).
var DaftarTabelDitulis = []string{Tabel}

// DaftarDibacaSaja - hanya SELECT.
var DaftarDibacaSaja = []string{MasterBusiness, MasterBenefit}
