//go:build db

package handlers_test

// DDL tiruan uji `db` = katalog DEV (lanjutan 1, code-review): salinan kolom yang ditulis tangan di
// `ddlTiruan` itulah yang dulu membuat uji `db` lulus dengan kolom yang tidak pernah ada di DEV. Tanpa Oracle.

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestDDLTiruanSamaDenganKatalogDEV(t *testing.T) {
	isi, err := os.ReadFile("../repository/testdata/katalog-dev.json")
	if err != nil {
		t.Fatal(err)
	}
	var k struct {
		Tabel map[string][]string `json:"tabel"`
		Tipe  map[string]string   `json:"tipe"`
	}
	if err := json.Unmarshal(isi, &k); err != nil {
		t.Fatal(err)
	}
	diperiksa := 0
	for _, d := range ddlTiruan {
		kolom, ada := k.Tabel[d.nama]
		if !ada {
			continue // master pendukung - di luar katalog yang diberikan
		}
		diperiksa++
		var nama []string
		for _, bagian := range strings.Split(d.kolom, ",") {
			f := strings.Fields(bagian)
			if len(f) < 2 {
				t.Fatalf("%s: kolom tiruan tak terbaca %q", d.nama, bagian)
			}
			nama = append(nama, f[0])
			if mau := k.Tipe[d.nama+"."+f[0]]; mau != "" && f[1] != mau {
				t.Errorf("%s.%s: tipe tiruan %s ≠ katalog DEV %s", d.nama, f[0], f[1], mau)
			}
		}
		if strings.Join(nama, ",") != strings.Join(kolom, ",") {
			t.Errorf("%s: kolom tiruan %v ≠ katalog DEV %v", d.nama, nama, kolom)
		}
	}
	if diperiksa != 2 {
		t.Errorf("kedua tabel produk diperiksa: %d", diperiksa)
	}
}
