package repository

// Penjaga lintas modul (RALAT R7, keputusan work owner 10-10-2026): view `RATE_LIFE` dibuang di DEV (ORA-00942) - nol kode
// produksi Go di `modul/` dan `inti/` yang masih menyebut objek `RATE_LIFE` (kata utuh; `M_RATE_LIFE`, `M_RATE_LIFE_SUMMARY`,
// `SEQ_M_RATE_LIFE` tidak terkena). Komentar diabaikan (riwayat boleh menyebut view itu). `modul/nbfacin` dikecualikan:
// di sana `RATE_LIFE` adalah NAMA KOLOM, bukan view.

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var polaViewRateLife = regexp.MustCompile(`(^|[^A-Za-z0-9_])RATE_LIFE([^A-Za-z0-9_]|$)`)

// buangKomentarBaris - potongan sebelum "//" tiap baris (cukup untuk penjaga ini: SQL modul tidak memuat "//").
func buangKomentarBaris(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if i := strings.Index(baris, "//"); i >= 0 {
			baris = baris[:i]
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}

func TestTidakAdaViewRateLifeDiKodeProduksi(t *testing.T) {
	akar := filepath.Join("..", "..", "..", "..")
	diperiksa := 0
	for _, dir := range []string{"modul", "inti"} {
		err := filepath.WalkDir(filepath.Join(akar, dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			n := filepath.ToSlash(p)
			if d.IsDir() {
				if d.Name() == "node_modules" || strings.HasSuffix(n, "/modul/nbfacin") {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
				return nil
			}
			isi, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			diperiksa++
			for i, baris := range strings.Split(buangKomentarBaris(string(isi)), "\n") {
				if polaViewRateLife.MatchString(baris) {
					t.Errorf("%s:%d menyebut view RATE_LIFE (dibuang 10-10-2026; pakai M_RATE_LIFE): %s", n, i+1, strings.TrimSpace(baris))
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if diperiksa < 100 {
		t.Fatalf("hanya %d berkas terbaca; pemindainya yang rusak", diperiksa)
	}
	for kode, mau := range map[string]bool{
		`FROM S.RATE_LIFE WHERE`: true, `ViewRate = "RATE_LIFE"`: true, `FROM S.M_RATE_LIFE`: false,
		`M_RATE_LIFE_SUMMARY`: false, `SEQ_M_RATE_LIFE`: false, `// view RATE_LIFE dibuang`: false,
	} {
		if polaViewRateLife.MatchString(buangKomentarBaris(kode)) != mau {
			t.Errorf("pola atas %q harus %v", kode, mau)
		}
	}
}
