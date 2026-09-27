package repository

// Periode nomor klaim - TANPA Oracle.
//
// Pemilik: A2 (butir o1).

import (
	"testing"
	"time"
)

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
		p := HitungPeriodeNomor(saat, k.closing)
		if p.MMYYYY != k.mauMM || p.Tahun != k.mauThn {
			t.Errorf("%s (%s, closing %d): (%q,%q), mau (%q,%q)",
				k.saat, k.apa, k.closing, p.MMYYYY, p.Tahun, k.mauMM, k.mauThn)
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
func TestPeriodeCutoverDipertahankan(t *testing.T) {
	for _, s := range []string{"2025-12-31 23:59:59", "2026-01-01 08:00:00", "2026-01-02 23:59:59"} {
		saat, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			t.Fatal(err)
		}
		p := HitungPeriodeNomor(saat, 25)
		if p.MMYYYY != "12.2025" || p.Tahun != "2025" {
			t.Errorf("%s: (%q,%q), mau (\"12.2025\",\"2025\")", s, p.MMYYYY, p.Tahun)
		}
	}
	// Sehari sesudahnya cabang itu mati.
	saat, _ := time.Parse("2006-01-02 15:04:05", "2026-01-03 08:00:00")
	if p := HitungPeriodeNomor(saat, 25); p.MMYYYY == "12.2025" {
		t.Error("cabang cutover masih menyala sesudah 02/01/2026")
	}
}

// TestTahunAdalahTeks - kunci ketiga `(CLASS, JENIS, TAHUN)`.
//
// `[data DBA]` kolomnya `VARCHAR2(5)`. Mengubahnya menjadi bilangan membuang
// bentuk aslinya (ADR-U-0022).
func TestTahunAdalahTeks(t *testing.T) {
	saat, _ := time.Parse("2006-01-02 15:04:05", "2026-06-10 10:00:00")
	p := HitungPeriodeNomor(saat, 25)
	if p.Tahun != "2026" {
		t.Errorf("tahun = %q, mau \"2026\"", p.Tahun)
	}
	if len(p.MMYYYY) != 7 || p.MMYYYY[2] != '.' {
		t.Errorf("MM.YYYY berbentuk %q, mau tujuh karakter ber-titik di posisi ketiga", p.MMYYYY)
	}
}
