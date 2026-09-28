package models_test

// Uji rumus `IMAGEID` - butir be, ralat 28-09-2026.

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
)

func TestStempelImageIDMeniruFF9(t *testing.T) {
	// `TO_CHAR(SYSTIMESTAMP,'YYYYMMDDHH24MISSFF9')` - 14 angka waktu + 9
	// angka pecahan detik, TANPA titik. Go menulis nanodetik dengan titik;
	// titik itu harus hilang, dan uji ini yang menjaganya.
	saat := time.Date(2026, 9, 28, 13, 45, 7, 123456789, time.UTC)
	got := models.StempelImageID(saat)
	if got != "20260928134507123456789" {
		t.Errorf("stempel = %q, mau %q", got, "20260928134507123456789")
	}
	if len(got) != 23 {
		t.Errorf("panjang = %d, mau 23 (14 + FF9)", len(got))
	}
	if strings.Contains(got, ".") {
		t.Errorf("stempel memuat titik: %q; Oracle tidak menuliskannya untuk FF9", got)
	}
	// Nanodetik NOL tetap sembilan angka - `FF9` selalu sembilan.
	nol := models.StempelImageID(time.Date(2026, 9, 28, 13, 45, 7, 0, time.UTC))
	if !strings.HasSuffix(nol, "000000000") {
		t.Errorf("pecahan nol tidak ditulis sembilan angka: %q", nol)
	}
}

func TestMasukanImageIDMemakaiAwalanLimaHuruf(t *testing.T) {
	// ⛔ `'ASMPP'` - LIMA huruf, VERBATIM `GenerateImageID_SQL.xml` b86.
	// Bukan `ASMAPP`, yang awalan TOKEN penyimpanan (butir an) dan hal yang
	// berbeda. Keduanya berdampingan di modul ini dan berbeda satu huruf.
	if models.AwalanImageID != "ASMPP" {
		t.Fatalf("awalan = %q, mau %q", models.AwalanImageID, "ASMPP")
	}
	m := models.MasukanImageID(
		time.Date(2026, 9, 28, 13, 45, 7, 123456789, time.UTC),
		"0123456789abcdef0123456789abcdef")
	// ⛔ GUID naik ke HURUF BESAR: Oracle mengubah RAW menjadi heksa huruf
	// besar saat menyambungnya ke teks, dan masukan yang berbeda huruf
	// menghasilkan hash yang berbeda sama sekali.
	if m != "ASMPP202609281345071234567890123456789ABCDEF0123456789ABCDEF" {
		t.Errorf("masukan = %q", m)
	}
}

// TestImageIDLiteralTetap mengunci rumusnya dengan satu nilai terhitung.
//
// ⛔ `[murni]` Nilainya MD5 dari teks yang tertulis di uji ini, dihitung
// sekali lalu dibekukan. Ia menjaga hal yang paling mudah bergeser tanpa
// terlihat: huruf besar/kecil keluarannya, dan pilihan fungsi hash-nya.
// Pasangannya uji `db` yang membandingkannya dengan `STANDARD_HASH` Oracle
// atas teks yang SAMA - lihat `repository/imageid_db_test.go`.
func TestImageIDLiteralTetap(t *testing.T) {
	const masukan = "ASMPP202609281345071234567890123456789ABCDEF0123456789ABCDEF"
	got := models.ImageIDDari(masukan)
	// MD5 heksa HURUF BESAR, 32 karakter.
	if len(got) != models.PanjangImageID {
		t.Errorf("panjang = %d, mau %d", len(got), models.PanjangImageID)
	}
	if got != strings.ToUpper(got) {
		t.Errorf("keluaran berhuruf kecil: %q; Oracle menuliskannya huruf besar", got)
	}
	if !regexp.MustCompile(`^[0-9A-F]{32}$`).MatchString(got) {
		t.Errorf("keluaran bukan heksa huruf besar 32 karakter: %q", got)
	}
	// ⛔ Nilai tetap. Bila baris ini berubah, rumusnya yang berubah - dan
	// setiap `IMAGEID` yang sudah tersimpan menjadi tak terjangkau.
	const tetap = "759CFD629759D8AA5DA0AEB2DC3A2CFB"
	if got != tetap {
		t.Errorf("IMAGEID = %s, mau %s (rumus berubah?)", got, tetap)
	}
}

func TestImageIDBaruTidakDapatDitebak(t *testing.T) {
	// ⛔ Dua panggilan pada JAM YANG SAMA tetap berbeda - itulah guna
	// `SYS_GUID()`. Ronde pertama memakai cap waktu dokumen sebagai IMAGEID;
	// siapa pun yang tahu kapan sebuah berkas diunggah dapat menyusun kunci
	// penyimpanannya.
	saat := time.Date(2026, 9, 28, 13, 45, 7, 0, time.UTC)
	a, err := models.ImageIDBaru(saat)
	if err != nil {
		t.Fatal(err)
	}
	b, err := models.ImageIDBaru(saat)
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Error("dua IMAGEID pada jam yang sama identik; SYS_GUID tidak tertiru")
	}
	for _, v := range []string{a, b} {
		if len(v) != models.PanjangImageID {
			t.Errorf("panjang %q = %d", v, len(v))
		}
	}
	// Dan ia TIDAK memuat cap waktunya apa adanya.
	if strings.Contains(a, "20260928") {
		t.Errorf("IMAGEID memuat cap waktu terbaca: %q", a)
	}
}
