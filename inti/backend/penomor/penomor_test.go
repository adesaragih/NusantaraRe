package penomor

// Periode nomor klaim - TANPA Oracle.
//
// Pemilik: A2 (butir o1).

import (
	"testing"
	"time"
)

// bantuHitungPeriode menjalankan HitungPeriodeNomor dan menggagalkan uji bila
// aturannya menolak masukannya.
//
// ⚠️ Fungsi bantu ini lahir 28-09-2026, saat `HitungPeriodeNomor` mulai
// mengembalikan galat - sebab pergeserannya kini dikerjakan
// `models.PeriodeProduksi`, yang MENOLAK hari tutup buku yang kosong atau di
// luar 1..31 alih-alih menebaknya.
func bantuHitungPeriode(t *testing.T, saat time.Time, hariClosing int) PeriodeNomor {
	t.Helper()
	p, err := HitungPeriodeNomor(saat, hariClosing)
	if err != nil {
		t.Fatalf("HitungPeriodeNomor(%s, %d): %v", saat.Format(time.RFC3339), hariClosing, err)
	}
	return p
}

// TestPeriodeNomorMengikutiHariTutupBuku - `[data DBA]` `ADD_MONTHS(+1)` bila
// hari melewati `TANGGAL_CLOSING.TANGGAL`.
func TestPeriodeNomorMengikutiHariTutupBuku(t *testing.T) {
	for _, k := range []struct {
		saat          string
		closing       int
		mauMM, mauThn string
		apa           string
	}{
		{"2026-03-20 10:00:00", 25, "03.2026", "2026", "sebelum tutup buku"},
		{"2026-03-25 10:00:00", 25, "03.2026", "2026", "tepat di hari tutup buku"},
		{"2026-03-26 10:00:00", 25, "04.2026", "2026", "lewat tutup buku"},
		// ⛔ ADD_MONTHS menggeser TAHUN juga - berbeda dari SaveAdjustment_Act,
		// yang tahunnya tidak ikut karena kedua cabang @if-nya identik.
		{"2026-12-26 10:00:00", 25, "01.2027", "2027", "Desember lewat tutup buku: tahun ikut"},
		{"2026-12-20 10:00:00", 25, "12.2026", "2026", "Desember sebelum tutup buku"},
		// Hari tutup buku yang berbeda menggeser ambangnya.
		{"2026-03-16 10:00:00", 15, "04.2026", "2026", "tutup buku tanggal 15"},
	} {
		saat, err := time.Parse("2006-01-02 15:04:05", k.saat)
		if err != nil {
			t.Fatal(err)
		}
		p := bantuHitungPeriode(t, saat, k.closing)
		if p.MMYYYY != k.mauMM || p.Tahun != k.mauThn {
			t.Errorf("%s (%s, closing %d): (%q,%q), mau (%q,%q)",
				k.saat, k.apa, k.closing, p.MMYYYY, p.Tahun, k.mauMM, k.mauThn)
		}
	}
}

// TestPeriodeNomorTidakMelompatiBulanPendek mengunci cacat yang sudah pernah
// ada di sini.
//
// ⛔ SEJARAHNYA. Ronde sebelumnya menggeser bulan dengan `saat.AddDate(0,1,0)`,
// dan Go melimpahkan tanggal yang tidak ada: 31 Januari + 1 bulan = 3 MARET,
// sehingga periodenya `03.2026` dan FEBRUARI terlewat. Oracle `ADD_MONTHS`
// menjepit ke 29 Februari, yaitu `02.2026`.
//
// Ia menyala pada tanggal 29-31 bulan yang penggantinya lebih pendek, dan
// hasilnya nomor yang terbukukan ke bulan yang salah tanpa satu pun galat.
// Itulah sebabnya kasusnya diuji satu per satu, bukan dipercayakan pada satu
// contoh.
func TestPeriodeNomorTidakMelompatiBulanPendek(t *testing.T) {
	for _, k := range []struct {
		saat  string
		mauMM string
		apa   string
	}{
		{"2026-01-29 10:00:00", "02.2026", "29 Januari: Februari 2026 hanya 28 hari"},
		{"2026-01-30 10:00:00", "02.2026", "30 Januari"},
		{"2026-01-31 10:00:00", "02.2026", "31 Januari - kasus yang dahulu jadi Maret"},
		{"2026-03-31 10:00:00", "04.2026", "31 Maret: April hanya 30 hari"},
		{"2026-05-31 10:00:00", "06.2026", "31 Mei: Juni hanya 30 hari"},
		{"2026-08-31 10:00:00", "09.2026", "31 Agustus: September hanya 30 hari"},
		{"2026-10-31 10:00:00", "11.2026", "31 Oktober: November hanya 30 hari"},
	} {
		saat, err := time.Parse("2006-01-02 15:04:05", k.saat)
		if err != nil {
			t.Fatal(err)
		}
		p := bantuHitungPeriode(t, saat, 25)
		if p.MMYYYY != k.mauMM {
			t.Errorf("%s (%s): %q, mau %q", k.saat, k.apa, p.MMYYYY, k.mauMM)
		}
	}
}

// TestPeriodeNomorMenolakHariTutupBukuTakMasukAkal - gagal terang.
//
// ⛔ Nol nilai pengganti, sejalan dengan tiket 02: hari tutup buku yang kosong
// menggeser periode SELURUH nomor yang terbit hari itu, dan menebaknya
// membuat kekeliruannya baru terlihat saat tutup buku bulan berikutnya.
func TestPeriodeNomorMenolakHariTutupBukuTakMasukAkal(t *testing.T) {
	saat, _ := time.Parse("2006-01-02 15:04:05", "2026-06-10 10:00:00")
	for _, hari := range []int{0, -1, 32, 99} {
		if _, err := HitungPeriodeNomor(saat, hari); err == nil {
			t.Errorf("hari tutup buku %d diterima; mau ditolak", hari)
		}
	}
}

// TestPeriodeCutoverDipertahankan - cabang yang sudah lewat, sengaja tetap.
//
// `[data DBA]` `IF TRUNC(v_now) <= TO_DATE('02/01/2026','DD/MM/YYYY')`.
//
// ⛔ Ia tidak akan menyala lagi hari ini, dan itu BUKAN alasan menghapusnya:
// migrasi data (tiket 13) menguraikan nomor LAMA, dan nomor yang lahir sebelum
// tanggal itu memakai periode `12.2025` meski dibuat Januari 2026.
//
// ⚠️ JAMNYA DIBACA DI JAKARTA sejak 28-09-2026, sama dengan pergeserannya -
// `TRUNC(v_now)` di procedure membaca `SYSDATE`, jam server Jakarta. Karena
// itu kasus-kasus di bawah memakai jam siang, yang harinya sama di UTC maupun
// di Jakarta; kasus batasnya diuji tersendiri di bawah.
func TestPeriodeCutoverDipertahankan(t *testing.T) {
	for _, s := range []string{"2025-12-31 12:00:00", "2026-01-01 08:00:00", "2026-01-02 12:00:00"} {
		saat, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			t.Fatal(err)
		}
		p := bantuHitungPeriode(t, saat, 25)
		if p.MMYYYY != "12.2025" || p.Tahun != "2025" {
			t.Errorf("%s: (%q,%q), mau 12.2025 dan 2025", s, p.MMYYYY, p.Tahun)
		}
	}
	// Sehari sesudahnya cabang itu mati.
	saat, _ := time.Parse("2006-01-02 15:04:05", "2026-01-03 08:00:00")
	if p := bantuHitungPeriode(t, saat, 25); p.MMYYYY == "12.2025" {
		t.Error("cabang cutover masih menyala sesudah 02/01/2026")
	}
}

// TestBatasCutoverDibacaDiJakarta mengunci zona pembacaannya.
//
// ⛔ 02-01-2026 pukul 23:59:59 UTC SUDAH tanggal 3 di Jakarta, jadi cutover
// TIDAK berlaku untuknya. Sebelum 28-09-2026 fungsi ini membaca harinya di
// zona `saat` sendiri dan menjawab sebaliknya. Yang dikunci di sini bukan
// selera: dua tempat yang masing-masing memutuskan zonanya sendiri berselisih
// sehari tepat di sekitar tengah malam, dan sehari di sini berarti satu bulan
// buku.
func TestBatasCutoverDibacaDiJakarta(t *testing.T) {
	saat, err := time.Parse(time.RFC3339, "2026-01-02T23:59:59Z")
	if err != nil {
		t.Fatal(err)
	}
	p := bantuHitungPeriode(t, saat, 25)
	if p.MMYYYY == "12.2025" {
		t.Error("cutover menyala untuk saat yang di Jakarta sudah 3 Januari")
	}
	// Dan yang di Jakarta masih tanggal 2 tetap kena cutover.
	saat2, _ := time.Parse(time.RFC3339, "2026-01-02T12:00:00Z")
	if p := bantuHitungPeriode(t, saat2, 25); p.MMYYYY != "12.2025" {
		t.Errorf("cutover tidak menyala untuk 2 Januari Jakarta: %q", p.MMYYYY)
	}
}

// TestTahunAdalahTeks - kunci ketiga `(CLASS, JENIS, TAHUN)`.
//
// `[data DBA]` kolomnya `VARCHAR2(5)`. Mengubahnya menjadi bilangan membuang
// bentuk aslinya (ADR-U-0022).
func TestTahunAdalahTeks(t *testing.T) {
	saat, _ := time.Parse("2006-01-02 15:04:05", "2026-06-10 10:00:00")
	p := bantuHitungPeriode(t, saat, 25)
	if p.Tahun != "2026" {
		t.Errorf("tahun = %q, mau 2026", p.Tahun)
	}
	if len(p.MMYYYY) != 7 || p.MMYYYY[2] != '.' {
		t.Errorf("MM.YYYY berbentuk %q, mau tujuh karakter ber-titik di posisi ketiga", p.MMYYYY)
	}
}
