// Package repository adalah SATU-SATUNYA lapisan modul Cause Of Loss Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services
// (pola benefitlife: satu transaksi per simpan, prosedur PEGA_M_CAUSEOFLOSS_LIFE tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// Keputusan work owner 08-10-2026 K1: SATU tabel `CAUSEOFLOSS_LIFE`. Tabel Pega `M_CAUSEOFLOSS_LIFE` BERGANTI NAMA
// menjadi `CAUSEOFLOSS_LIFE` (nama view lama yang dibuang) dan menjadi tabel flat tanpa JSONDATA - migrasi modul
// 090-092. Nol JSON di modul ini.
package repository

// Tabel - `CAUSEOFLOSS_LIFE`, kelas `ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE`; kolom = KolomTabel. SATU-SATUNYA objek yang
// DITULIS. Dibaca juga modul masterproductnamelife (pemilih Cause Of Loss) - nama dan kolomnya tidak berubah.
const Tabel = "CAUSEOFLOSS_LIFE"

// Seq - `M_CAUSEOFLOSS_LIFE_SEQ` warisan (last_number 5 DEV), prosedur PEGA_M_CAUSEOFLOSS_LIFE. K3: TETAP, nol sequence
// baru; nomornya dibentuk models.BentukID.
const Seq = "M_CAUSEOFLOSS_LIFE_SEQ"

// KolomTabel - kolom `CAUSEOFLOSS_LIFE` sesudah 092 = kolom view lama bernama sama, urutan sama (ID warisan + 091).
var KolomTabel = []string{"ID", "CAUSEOFLOSS"}

// Kolom yang dibaca (nol `SELECT *`).
const kolomBaca = `ID, CAUSEOFLOSS`

// DaftarTabelDitulis - penjaga modul: hanya tabel ini yang menerima INSERT/UPDATE (nol DELETE: XML tanpa Delete).
var DaftarTabelDitulis = []string{Tabel}
