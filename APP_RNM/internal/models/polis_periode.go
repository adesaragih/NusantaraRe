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
// owner]` "ikuti yang dari DB" - penerapannya pada `ProdDateTime` OQ-PL-13.
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
	"errors"
	"fmt"
	"time"
)

// ErrTanggalTutupBukuKosong - tabel sumbernya kosong atau tidak terbaca.
//
// ⛔ Pesannya MENYEBUT TABEL SUMBERNYA. Galat yang hanya berkata "periode
// tidak dapat ditentukan" membuat orang mencari di kode; yang menyebut
// `POOLDATA.TANGGAL_CLOSING` membuat orang membuka tabelnya.
var ErrTanggalTutupBukuKosong = errors.New(
	"models: tanggal tutup buku tidak terbaca dari POOLDATA.TANGGAL_CLOSING; " +
		"transaksi ditolak, dan TIDAK ada nilai pengganti yang dipakai")

// ErrTanggalTutupBukuTidakMasukAkal - nilainya di luar 1..31.
var ErrTanggalTutupBukuTidakMasukAkal = errors.New(
	"models: tanggal tutup buku di luar 1..31")

// JamPeriodeGMT adalah jam periode produksi - VERBATIM `T050000.000 GMT`.
//
// ⚠️ 05:00 GMT = 12:00 WIB. Jamnya ditulis dalam GMT sebab rule aslinya
// menulisnya begitu; mengubahnya menjadi "12:00 Asia/Jakarta" bermakna sama
// hari ini dan berbeda bila zona waktunya pernah bergeser.
const JamPeriodeGMT = 5

// PeriodeProduksi menghitung periode produksi sebuah transaksi.
//
// `tglTutupBuku` adalah nilai kolom `TANGGAL`; `saat` jam transaksinya.
//
// Mengembalikan tanggal **1** bulan periode, pukul `05:00 GMT`.
//
// ⛔ Nol nilai pengganti. `tglTutupBuku` <= 0 ditolak - lihat
// `ErrTanggalTutupBukuKosong` dan sebabnya.
func PeriodeProduksi(tglTutupBuku int, saat time.Time) (time.Time, error) {
	if tglTutupBuku <= 0 {
		return time.Time{}, ErrTanggalTutupBukuKosong
	}
	if tglTutupBuku > 31 {
		return time.Time{}, fmt.Errorf("%w: %d",
			ErrTanggalTutupBukuTidakMasukAkal, tglTutupBuku)
	}
	// ⛔ Zona Jakarta, sebab rule aslinya membaca
	// `@CurrentDate(..., "Asia/Jakarta")`. Membandingkan tanggal dalam UTC
	// menggeser hasilnya tujuh jam - dan tepat di sekitar tengah malam itu
	// berarti sehari.
	jakarta := zonaJakarta()
	lokal := saat.In(jakarta)

	tahun, bulan := lokal.Year(), lokal.Month()
	// b1211: `>` dan bukan `>=`.
	if lokal.Day() > tglTutupBuku {
		bulan++
		// ⛔ Pergantian tahun. `CurrentMonth+1` di Pega menghasilkan "13"
		// untuk Desember; `time.Date` menormalkannya, dan uji menahannya
		// tetap begitu.
		if bulan > time.December {
			bulan = time.January
			tahun++
		}
	}
	return time.Date(tahun, bulan, 1, JamPeriodeGMT, 0, 0, 0, time.UTC), nil
}

// DiJakarta memindahkan satu saat ke zona Asia/Jakarta.
//
// ⛔ SATU SUMBER ZONA untuk seluruh repo. Setiap aturan yang membandingkan
// TANGGAL - hari tutup buku, cutover penomoran - harus memakai hari yang sama
// dengan yang dilihat Oracle, sebab `SYSDATE` di sana adalah waktu server
// Jakarta. Dua tempat yang masing-masing memutuskan zonanya sendiri akan
// berselisih tepat di sekitar tengah malam, yaitu satu hari, yaitu satu bulan
// buku.
func DiJakarta(t time.Time) time.Time { return t.In(zonaJakarta()) }

// zonaJakarta mengembalikan zona Asia/Jakarta, dengan cadangan tetap.
//
// ⚠️ Basis data zona waktu tidak selalu ada di Windows. UTC+7 tetap zona yang
// benar untuk Jakarta; yang hilang hanya namanya.
func zonaJakarta() *time.Location {
	if loc, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return loc
	}
	return time.FixedZone("WIB", 7*60*60)
}

// PeriodeTeks menulis periode sebagai `YYYY-MM` - yang layar tampilkan.
//
// ⛔ Periode yang terpilih HARUS terlihat pemakai sebelum ia menyimpan
// (AC tiket 02), bukan tersimpan diam-diam. Transaksi yang mendarat di bulan
// yang salah karena seseorang menyimpannya lewat tengah malam adalah
// kekeliruan yang hanya dapat dicegah dengan menunjukkannya lebih dulu.
func PeriodeTeks(periode time.Time) string {
	return periode.Format("2006-01")
}
