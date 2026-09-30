package models

// Periode produksi dan tanggal tutup buku - tiket 02 PremiumList Life.
//
// Untuk apa berkas ini: menentukan periode produksi sebuah transaksi dari
// tanggal tutup buku yang berlaku. Seluruhnya MURNI, dan jamnya diserahkan
// pemanggil - aturan periode harus dapat diuji tanpa menunggu tanggal nyata.
//
// Sumbernya `Activity/SubmitPremiumList_Act.xml`, dibaca sebagai pohon:
//
// ⛔ RALAT 28-09-2026 (sensus remark GILIRAN-12): langkah 2-4 di bawah memang
// hidup, tetapi hasilnya hanya mengalir ke `TempGenerate.CARI1/2` (tidak
// dibaca rule mana pun; SQL penomoran memakai `ParamSeq`) dan ke langkah 15
// yang TER-REMARK (`//` b3398). Aturan periode yang HIDUP ada dua: nomor PL
// digulir `PROC_GENERATE_SEQUENCE_NUMBER` (membaca `TANGGAL_CLOSING` sendiri,
// SUMBER-PENOMORAN-DBA.md), dan `ProdDateTime` digeser `InsertJsonPolisLife_Act`
// langkah 4 (b1065, nilai b1092, gerbang b1170 `>25` TERTANAM). Bentuk
// hitungnya ditiru dari b1092; ambangnya dari tabel sesuai `[keputusan work
// owner]` "ikuti yang dari DB" - untuk periode yang ditampilkan dan penomoran
// PL. ⛔ OQ-PL-13 DITUTUP 29-09-2026 (GILIRAN-17): `ProdDateTime` sendiri IKUT
// XML - 25 tertanam (`AmbangProdDateTimePega`).
//
//	langkah 2 b715   `RDB-List` -> `GETTanggalClosing_SQL`
//	                  (`SELECT * FROM POOLDATA.TANGGAL_CLOSING`)
//	langkah 3 b918   `Local.TglProd = TglProd.pxResults(1).TANGGAL`
//	langkah 3 b966   `Local.TglProd = @if(Local.TglProd=="",25,Local.TglProd)`
//	langkah 4 b1211  `Local.NextMonth = @if(@toDecimal(Local.currentdate) >
//	                   Local.TglProd, CurrentMonth+1, CurrentMonth)`
//	                 `@CurrentDate("yyyy") + NextMonth + "01T050000.000 GMT"`
//
// ⛔ PERBEDAAN SENGAJA DARI PEGA, dan satu-satunya di berkas ini: baris b966
// diam-diam memakai **25** ketika tabelnya kosong. Kami **menolak**
// transaksinya. Fallback diam membukukan transaksi ke periode yang SALAH
// tanpa meninggalkan jejak, dan kekeliruannya baru terlihat saat tutup buku
// bulan berikutnya - ketika angkanya sudah terlanjur salah. Itu kegagalan
// tersembunyi, bukan ketahanan. `[keputusan work owner]`
//
// ⛔ PERBANDINGANNYA `>`, BUKAN `>=`. Baris b1211 (Submit langkah 4) dan b1170
// (InsertJsonPolisLife langkah 4) keduanya memakai `>` (b3491 milik langkah 15
// yang ter-remark), jadi transaksi pada tanggal yang SAMA dengan tanggal tutup buku tetap
// di periode berjalan. Satu tanda sama dengan menggeser sehari penuh
// transaksi ke bulan yang salah.
//
// Dibaca sesudah: polis_penawaran.go.

import (
	"time"
)

// AmbangProdDateTimePega - OQ-PL-13 DITUTUP 29-09-2026 (GILIRAN-17)
// `[keputusan work owner]`: `ProdDateTime` ikut XML, ambang 25 TERTANAM -
// `InsertJsonPolisLife_Act` langkah 4 (b1065, hidup), gerbang b1170
// `@toDecimal(Local.currentdate)>25`, nilai b1092 `@CurrentDate("yyyy",
// "Asia/Jakarta")+Local.NextMonth+"01T050000.000 GMT"`.
//
// ⚠️ Nol pemakai hari ini: `ProdDateTime` hanya hidup di JSON halaman (langkah
// 5 b1237), dan `JSON_POLIS` tidak lagi ditulis (pl1); ia juga bukan muatan
// `convertJsonNusareToProduction` (`tglInput` = `.pxCreateDateTime`). Penjaga
// "nol ambang tertanam" (ambangperiode_test.go) tetap berlaku untuk periode
// dan penomoran. Uji: `TestAmbangProdDateTimeTertanamB1170`.
const AmbangProdDateTimePega = 25

// PeriodeTeks menulis periode sebagai `YYYY-MM` - yang layar tampilkan.
//
// ⛔ Periode yang terpilih HARUS terlihat pemakai sebelum ia menyimpan
// (AC tiket 02), bukan tersimpan diam-diam. Transaksi yang mendarat di bulan
// yang salah karena seseorang menyimpannya lewat tengah malam adalah
// kekeliruan yang hanya dapat dicegah dengan menunjukkannya lebih dulu.
func PeriodeTeks(periode time.Time) string {
	return periode.Format("2006-01")
}
