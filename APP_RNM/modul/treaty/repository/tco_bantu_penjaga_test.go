package repository

// Pembantu penjaga statik milik Treaty Contract Out - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): penjaga Treaty dulu meminjam pembantu uji
// Claim Life (`batasanpemakaian_test.go`, `migrasibatas_test.go`,
// `migrasi_test.go`). Modul Treaty kini membawa salinannya sendiri, supaya
// uji satu modul tidak pernah bergantung pada berkas uji modul lain. Isinya
// sama; satu beda yang disengaja: `seluruhSQL` membaca SETIAP folder
// `migrations/` di bawah akar aplikasi - bukan hanya milik satu modul - sebab
// penjaga "nol tabel baru Treaty" menjaga migrasi semua modul.

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/inti/migrasi"
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

// sumberMigrasi mengumpulkan SETIAP folder `migrations/` di bawah akar
// aplikasi sebagai sumber pelari migrasi.
func sumberMigrasi(t *testing.T) []fs.FS {
	t.Helper()
	var sumber []fs.FS
	err := filepath.Walk(filepath.FromSlash(akarModul), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		switch info.Name() {
		case "frontend", "node_modules", "bin", ".git":
			return filepath.SkipDir
		case "migrations":
			// Pelari membaca `migrations/*.sql` dari akar sumbernya.
			sumber = append(sumber, os.DirFS(filepath.Dir(p)))
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sumber) == 0 {
		t.Fatal("nol folder migrations terbaca; pembacanya yang rusak")
	}
	return sumber
}

// seluruhSQL - isi setiap langkah migrasi semua modul, per nama berkas.
func seluruhSQL(t *testing.T, mundur bool) map[string]string {
	t.Helper()
	langkah, err := migrasi.Daftar(mundur, sumberMigrasi(t)...)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, m := range langkah {
		out[m.Nama] = strings.Join(m.Pernyataan, "\n")
	}
	return out
}
