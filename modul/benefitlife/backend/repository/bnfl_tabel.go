// Package repository adalah SATU-SATUNYA lapisan modul Benefit yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services
// (pola ririsklife: satu transaksi per simpan, prosedur PEGA_M_BENEFIT_LIFE tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// Keputusan work owner 08-10-2026 K1: SATU tabel `BENEFIT_LIFE`. Tabel Pega `M_BENEFIT_LIFE` BERGANTI NAMA menjadi
// `BENEFIT_LIFE` (nama view lama yang dibuang) dan menjadi tabel flat tanpa JSONDATA - migrasi inti 942-944. Nol JSON
// di modul ini.
package repository

// Tabel - `BENEFIT_LIFE`, kelas `ASM-FW-GISFW-Int-BENEFIT_LIFE`; kolom = KolomTabel. SATU-SATUNYA objek yang DITULIS.
const Tabel = "BENEFIT_LIFE"

// Seq - `M_BENEFIT_LIFE_SEQ` warisan (last_number 12 DEV), prosedur PEGA_M_BENEFIT_LIFE. K3: TETAP, nol sequence baru;
// nomornya dibentuk models.BentukID.
const Seq = "M_BENEFIT_LIFE_SEQ"

// KolomTabel - kolom `BENEFIT_LIFE` sesudah 944 = kolom view lama bernama sama, urutan sama (ID warisan + 943).
var KolomTabel = []string{"ID", "BENEFIT"}

// Kolom yang dibaca (nol `SELECT *`).
const kolomBaca = `ID, BENEFIT`

// DaftarTabelDitulis - penjaga modul: hanya tabel ini yang menerima INSERT/UPDATE (nol DELETE: XML tanpa Delete).
var DaftarTabelDitulis = []string{Tabel}
