package models_test

// Uji periode produksi - tiket 02 PremiumList Life.
//
// ⛔ Jamnya diserahkan, jadi seluruh uji di sini berjalan tanpa menunggu
// tanggal nyata - AC tiket 02.

import (
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/inti/penomor"
)

// wib menyusun waktu dalam zona Jakarta.
func wib(tahun int, bulan time.Month, hari, jam int) time.Time {
	return time.Date(tahun, bulan, hari, jam, 0, 0, 0,
		time.FixedZone("WIB", 7*60*60))
}

func TestPerbandinganLebihBesarBukanLebihBesarSama(t *testing.T) {
	// ⛔ b1211 dan b1170 keduanya memakai `>` (b3491 ter-remark, sensus
	// 28-09-2026). Satu tanda sama dengan menggeser SEHARI PENUH transaksi ke
	// bulan yang salah.
	const tutupBuku = 25

	// Tepat pada tanggal tutup buku -> TETAP di periode berjalan.
	sama, err := penomor.PeriodeProduksi(tutupBuku, wib(2026, time.September, 25, 23))
	if err != nil {
		t.Fatal(err)
	}
	if sama.Month() != time.September || sama.Year() != 2026 {
		t.Errorf("tanggal 25 -> %s, mau September 2026", models.PeriodeTeks(sama))
	}

	// Sehari sesudahnya -> bulan berikutnya.
	sesudah, err := penomor.PeriodeProduksi(tutupBuku, wib(2026, time.September, 26, 0))
	if err != nil {
		t.Fatal(err)
	}
	if sesudah.Month() != time.October {
		t.Errorf("tanggal 26 -> %s, mau Oktober", models.PeriodeTeks(sesudah))
	}
}

func TestPeriodeSelaluTanggalSatuJamLimaGMT(t *testing.T) {
	// VERBATIM `+"01T050000.000 GMT"`.
	p, err := penomor.PeriodeProduksi(25, wib(2026, time.September, 26, 10))
	if err != nil {
		t.Fatal(err)
	}
	if p.Day() != 1 {
		t.Errorf("hari = %d, mau 1", p.Day())
	}
	if p.UTC().Hour() != penomor.JamPeriodeGMT {
		t.Errorf("jam UTC = %d, mau %d", p.UTC().Hour(), penomor.JamPeriodeGMT)
	}
	if p.Minute() != 0 || p.Second() != 0 {
		t.Errorf("periode bukan pada jam bulat: %s", p)
	}
}

func TestPergantianTahunBukanBulanTigaBelas(t *testing.T) {
	// ⛔ `CurrentMonth+1` di Pega menghasilkan "13" untuk Desember.
	p, err := penomor.PeriodeProduksi(25, wib(2026, time.December, 31, 9))
	if err != nil {
		t.Fatal(err)
	}
	if p.Year() != 2027 || p.Month() != time.January {
		t.Errorf("Desember lewat tutup buku -> %s, mau 2027-01", models.PeriodeTeks(p))
	}
}

func TestTabelKosongDITOLAK_BukanDiamDiamPakai25(t *testing.T) {
	// ⛔ INTI TIKET 02, dan satu-satunya perbedaan sengaja dari Pega. Baris
	// b966 `@if(Local.TglProd=="",25,Local.TglProd)` diam-diam memakai 25;
	// kami menolak. Fallback diam membukukan transaksi ke periode yang SALAH
	// tanpa meninggalkan jejak.
	for _, kosong := range []int{0, -1} {
		_, err := penomor.PeriodeProduksi(kosong, wib(2026, time.September, 26, 10))
		if !errors.Is(err, penomor.ErrTanggalTutupBukuKosong) {
			t.Errorf("tutup buku %d: %v, mau ErrTanggalTutupBukuKosong", kosong, err)
		}
	}
	// ⛔ Dan pesannya MENYEBUT TABEL SUMBERNYA. Galat yang hanya berkata
	// "periode tidak dapat ditentukan" membuat orang mencari di kode.
	pesan := penomor.ErrTanggalTutupBukuKosong.Error()
	if !strings.Contains(pesan, "POOLDATA.TANGGAL_CLOSING") {
		t.Errorf("pesan tidak menyebut tabel sumbernya: %s", pesan)
	}
	// Dan ia menyatakan bahwa nol pengganti dipakai.
	if !strings.Contains(strings.ToUpper(pesan), "TIDAK ADA NILAI PENGGANTI") {
		t.Errorf("pesan tidak menyatakan nol pengganti: %s", pesan)
	}
}

func TestTanggalTutupBukuTidakMasukAkalDitolak(t *testing.T) {
	_, err := penomor.PeriodeProduksi(32, wib(2026, time.September, 10, 10))
	if !errors.Is(err, penomor.ErrTanggalTutupBukuTidakMasukAkal) {
		t.Errorf("tutup buku 32: %v", err)
	}
}

func TestZonaJakartaMenentukanHarinya(t *testing.T) {
	// ⚠️ Rule aslinya membaca `@CurrentDate(..., "Asia/Jakarta")`.
	// Membandingkan tanggal dalam UTC menggeser hasilnya tujuh jam - dan
	// tepat di sekitar tengah malam itu berarti SEHARI.
	//
	// 2026-09-25 20:00 UTC = 2026-09-26 03:00 WIB -> sudah lewat tutup buku.
	utcMalam := time.Date(2026, time.September, 25, 20, 0, 0, 0, time.UTC)
	p, err := penomor.PeriodeProduksi(25, utcMalam)
	if err != nil {
		t.Fatal(err)
	}
	if p.Month() != time.October {
		t.Errorf("20:00 UTC tanggal 25 = 03:00 WIB tanggal 26 -> %s, mau Oktober",
			models.PeriodeTeks(p))
	}
}

func TestPeriodeTeksUntukLayar(t *testing.T) {
	// Periode HARUS terlihat pemakai sebelum ia menyimpan.
	p, err := penomor.PeriodeProduksi(25, wib(2026, time.September, 10, 10))
	if err != nil {
		t.Fatal(err)
	}
	if got := models.PeriodeTeks(p); got != "2026-09" {
		t.Errorf("PeriodeTeks = %q, mau 2026-09", got)
	}
}
