package deidentifikasi

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// folderFixture - fixture rekonsiliasi hasil alat ini (tiket 15, keputusan work
// owner 01-10-2026 butir 38). Nama berkasnya netral: siklus-lini-urut.
const folderFixture = "../../services/premium/testdata/kasus"

// TestFixtureKasusSudahDibersihkan - penjaga: setiap fixture sudah melewati alat ini
// (semua daun kunci/jalur BUANG kosong) dan nama berkasnya tidak memuat nomor kasus.
// Berkas mentah yang disalin ke sini membuat uji ini gagal.
func TestFixtureKasusSudahDibersihkan(t *testing.T) {
	berkas, err := filepath.Glob(filepath.Join(folderFixture, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(berkas) != 5 {
		t.Fatalf("%d fixture di %s, mau 5", len(berkas), folderFixture)
	}
	namaSah := regexp.MustCompile(`^(nb|rnw|edm)-[a-z]+-[0-9]\.json$`)
	for _, b := range berkas {
		if !namaSah.MatchString(filepath.Base(b)) {
			t.Errorf("%s: nama bukan siklus-lini-urut", filepath.Base(b))
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		if err := SudahBersih(isi); err != nil {
			t.Errorf("%s: %v", filepath.Base(b), err)
		}
	}
}
