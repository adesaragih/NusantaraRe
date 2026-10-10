// Package repository adalah SATU-SATUNYA lapisan modul Cover Life yang berbicara ke Oracle. Setiap nilai diikat (`:n`,
// go-ora mengikat menurut URUTAN kemunculan); nama objek lewat `Qualify`; nol COMMIT - transaksi milik services (pola
// causeoflosslife: satu transaksi per simpan, prosedur PEGA_M_COVER_LIFE tidak dipanggil).
//
// Untuk apa berkas ini: nama objek Oracle dan kolom modul ini, ditulis SEKALI.
//
// Keputusan work owner 08-10-2026 C1: SATU tabel `M_COVER_LIFE` - nama TETAP (TANPA RENAME), menjadi tabel flat tanpa
// JSONDATA (migrasi modul 085-086). View lama yang dulu membaca tabel ini DIBUANG 086 dan TIDAK disebut kode mana pun
// (uji TestKodeTidakMenyebutViewLama). Nol JSON di modul ini.
package repository

// Tabel - `M_COVER_LIFE`, kelas `ASM-FW-GISFW-Int-COVER_LIFE`; kolom = KolomTabel. SATU-SATUNYA objek yang DITULIS dan
// dibaca. Tidak ada pembaca lain di repo (fakta WO).
const Tabel = "M_COVER_LIFE"

// Seq - `M_COVER_LIFE_SEQ` warisan (last_number 5 DEV), prosedur PEGA_M_COVER_LIFE. C2: TETAP, nol sequence baru;
// nomornya dibentuk models.BentukID.
const Seq = "M_COVER_LIFE_SEQ"

// KolomTabel - kolom `M_COVER_LIFE` sesudah 086 = kolom view lama bernama sama, urutan sama (ID warisan + 085).
var KolomTabel = []string{"ID", "COVER", "NOTE"}

// Kolom yang dibaca (nol `SELECT *`).
const kolomBaca = `ID, COVER, NOTE`

// DaftarTabelDitulis - penjaga modul: hanya tabel ini yang menerima INSERT/UPDATE (nol DELETE: XML tanpa Delete).
var DaftarTabelDitulis = []string{Tabel}
