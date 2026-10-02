package acceptance

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// folderLimit - ekspor tabel limit `POOLDATA.M_LIMIT_*` dari work owner (01-10-2026,
// `docs/KEPUTUSAN-30-09-2026.md` butir 42), disalin byte demi byte dari `DDL\`.
const folderLimit = "testdata/limit"

// bacaCSVLimit - satu CSV ekspor tabel limit (pemisah `;`). ReadAll menolak baris
// yang lebarnya beda dari header.
func bacaCSVLimit(t *testing.T, jalur string) [][]string {
	t.Helper()
	f, err := os.Open(jalur)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ';'
	baris, err := r.ReadAll()
	if err != nil {
		t.Fatalf("%s: %v", jalur, err)
	}
	return baris
}

// TestFixtureLimitTanpaNamaLogin - penjaga: keenam tabel ada, kolom `NAMA` dan `LOGIN`
// tidak pernah ikut (K-025, CLAUDE.md §4 butir 10), dan setiap baris selebar header.
func TestFixtureLimitTanpaNamaLogin(t *testing.T) {
	berkas, err := filepath.Glob(filepath.Join(folderLimit, "*.csv"))
	if err != nil {
		t.Fatal(err)
	}
	var nama []string
	for _, b := range berkas {
		nama = append(nama, filepath.Base(b))
	}
	sort.Strings(nama)
	mau := "M_LIMIT_ENGINEERINGG.csv,M_LIMIT_FINANCIALINS.csv,M_LIMIT_NONPROPANDENGG.csv," +
		"M_LIMIT_PROPERTYY.csv,M_LIMIT_PROPERTY_NON_PREFERREDD.csv,M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL.csv"
	if strings.Join(nama, ",") != mau {
		t.Fatalf("berkas %v", nama)
	}
	for _, b := range berkas {
		for _, k := range bacaCSVLimit(t, b)[0] {
			if k == "NAMA" || k == "LOGIN" {
				t.Errorf("%s: kolom %s tidak boleh ada", filepath.Base(b), k)
			}
		}
	}
}
