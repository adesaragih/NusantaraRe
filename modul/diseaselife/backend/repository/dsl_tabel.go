// Package repository adalah SATU-SATUNYA lapisan modul Disease Life yang berbicara ke Oracle. Setiap nilai diikat
// (`:n`, go-ora mengikat menurut URUTAN kemunculan, setiap penampung unik); nama objek lewat `Qualify`; nol COMMIT -
// transaksi milik services (pola causeoflosslife: satu transaksi per simpan, prosedur PEGA_DISEASE_LIFE tidak
// dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// Keputusan work owner 08-10-2026 D1: tabel `DISEASE_LIFE` SUDAH flat - TIDAK di-RENAME, kolomnya TIDAK diubah (dibaca
// juga claimlife dan ditulis prosedur Pega PEGA_DISEASE_LIFE yang masih VALID). Migrasi modul hanya menambah sequence
// `SEQ_DISEASE_LIFE` (080) dan PK `PK_DISEASE_LIFE` (081).
package repository

// Tabel - `DISEASE_LIFE`, kelas `ASM-FW-GISFW-Int-DISEASE_LIFE`; kolom = KolomTabel. SATU-SATUNYA objek yang DITULIS.
// Dibaca juga modul claimlife (pencarian diagnosa, kueri tidak berubah).
const Tabel = "DISEASE_LIFE"

// Seq - `SEQ_DISEASE_LIFE` (migrasi 080, D1.1). BUKAN M_DISEASE_LIFE_SEQ warisan (rumus Pega bertabrakan dengan data
// impor - models/idbaru.go).
const Seq = "SEQ_DISEASE_LIFE"

// KolomTabel - kolom `DISEASE_LIFE` warisan, urutan katalog DEV (fakta WO 08-10-2026), tidak diubah migrasi mana pun.
var KolomTabel = []string{"ID", "ICD_CODE", "DISEASE"}

// Kolom yang dibaca (nol `SELECT *`).
const kolomBaca = `ID, ICD_CODE, DISEASE`

// DaftarTabelDitulis - penjaga modul: hanya tabel ini yang menerima INSERT/UPDATE (nol DELETE: XML tanpa Delete).
var DaftarTabelDitulis = []string{Tabel}
