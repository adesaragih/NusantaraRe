package repository

// Pembantu penjaga statik milik Treaty Contract Out - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): penjaga Treaty dulu meminjam pembantu uji
// Claim Life (`batasanpemakaian_test.go`, `migrasibatas_test.go`,
// `migrasi_test.go`). Modul Treaty kini membawa salinannya sendiri, supaya
// uji satu modul tidak pernah bergantung pada berkas uji modul lain.
//
// Paket 8: penjaga "nol tabel baru Treaty" (tco4), yang membaca migrasi SEMUA
// modul, pindah ke `inti/backend/penjaga/lintasaplikasi_test.go` - pembaca migrasinya
// (`seluruhSQL`, `sumberMigrasi`) ikut pergi dari berkas ini.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// akarModul menunjuk folder APP_RNM dari folder paket ini.
const akarModul = "../../.."

// berkasGoSelainTest mengumpulkan seluruh berkas .go yang BUKAN test.
func berkasGoSelainTest(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	err := filepath.Walk(filepath.FromSlash(akarModul), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// frontend dan hasil bangun tidak memuat kode Go yang relevan.
			switch info.Name() {
			case "frontend", "node_modules", "bin", ".git":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		isi, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		hasil[filepath.ToSlash(p)] = string(isi)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) == 0 {
		t.Fatal("nol berkas .go terbaca; pembacanya yang rusak, bukan kodenya")
	}
	return hasil
}

// berkasSumberProduksi - berkas .go bukan-test DAN berkas .sql.
func berkasSumberProduksi(t *testing.T) map[string]string {
	t.Helper()
	hasil := berkasGoSelainTest(t)
	err := filepath.Walk(filepath.FromSlash(akarModul),
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				switch info.Name() {
				case "frontend", "node_modules", "bin", ".git":
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.EqualFold(filepath.Ext(p), ".sql") {
				return nil
			}
			isi, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			hasil[filepath.ToSlash(p)] = string(isi)
			return nil
		})
	if err != nil {
		t.Fatal(err)
	}
	return hasil
}

// cariBerkas mencari isi berkas yang jalurnya berakhiran `akhiran`.
func cariBerkas(peta map[string]string, akhiran string) (string, bool) {
	for nama, isi := range peta {
		if strings.HasSuffix(nama, akhiran) {
			return isi, true
		}
	}
	return "", false
}

// buangKomentarGo menghapus baris komentar `//`.
func buangKomentarGo(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "//") {
			continue
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}

// buangKomentarSumber menghapus baris komentar - `//` untuk Go, `--` untuk SQL.
func buangKomentarSumber(nama, isi string) string {
	awalan := "//"
	if strings.EqualFold(filepath.Ext(nama), ".sql") {
		awalan = "--"
	}
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), awalan) {
			continue
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}
