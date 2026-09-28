package repository

// Penjaga: nol penampung berulang pada SQL berklausa pembatas baris.
//
// ⛔ Lahir dari uji asap baca-saja ke DEV (GILIRAN-12 paket 3):
// `GET /api/polis-life` menjawab 500 - `ORA-01008: not all variables bound`.
// Sebabnya `(:1 IS NULL OR w.POSITION = :1)` di kueri yang juga memakai
// `OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY`. Diuji langsung ke DEV (SELECT
// saja): penampung berulang TANPA klausa pembatas baris terikat benar; DENGAN
// klausa itu, Oracle menulis ulang kuerinya dan pengikatan per nama patah.
// Penampung unik (nilai dikirim dua kali) terikat benar di keduanya.
//
// Nol uji lain yang menangkapnya: semua uji repository berjalan tanpa Oracle.
//
// ⚠️ CAKUPANNYA DINYATAKAN: literal backtick yang memuat kata kunci SQL, bind
// angka (`:1`) maupun bernama (`:offset`). SQL yang dirakit dari potongan
// string terpisah (klausa pembatas baris di potongan lain dari WHERE-nya)
// tidak terlihat oleh penjaga ini.

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestNolPenampungBerulangDiSQLBerpembatasBaris(t *testing.T) {
	berkas, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	literal := regexp.MustCompile("(?s)`([^`]*)`")
	// Bind angka atau bernama; `HH24:MI:SS` tidak ikut (didahului huruf/angka).
	penampung := regexp.MustCompile(`(?:^|[^A-Za-z0-9_':]):(\d+|[A-Za-z_]\w*)\b`)
	sql := regexp.MustCompile(`(?i)\b(SELECT|UPDATE|INSERT|DELETE|MERGE)\b`)
	pembatas := regexp.MustCompile(`(?i)\bFETCH\s+(NEXT|FIRST)\b|\bOFFSET\b`)
	diperiksa := 0
	for _, b := range berkas {
		n := b.Name()
		if b.IsDir() || !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(n)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range literal.FindAllStringSubmatch(string(isi), -1) {
			q := m[1]
			// Literal di komentar (`OFFSET … FETCH` sebagai teks) bukan SQL.
			if !pembatas.MatchString(q) || !sql.MatchString(q) {
				continue
			}
			diperiksa++
			lihat := map[string]bool{}
			for _, p := range penampung.FindAllStringSubmatch(q, -1) {
				if lihat[p[1]] {
					t.Errorf("%s: penampung :%s berulang di SQL berpembatas baris "+
						"(ORA-01008 di Oracle); pakai penampung unik dan kirim nilainya dua kali:\n%s",
						n, p[1], strings.TrimSpace(q))
					break
				}
				lihat[p[1]] = true
			}
		}
	}
	if diperiksa < 5 {
		t.Fatalf("hanya %d SQL berpembatas baris terbaca; penelusurnya yang rusak", diperiksa)
	}
}
