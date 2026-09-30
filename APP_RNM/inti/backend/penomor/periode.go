package penomor

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
