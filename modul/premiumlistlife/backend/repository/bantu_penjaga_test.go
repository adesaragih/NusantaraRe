package repository

// Pembantu uji milik PremiumList Life - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): uji PremiumList dulu meminjam pembantu uji
// Claim Life (`batasanpemakaian_test.go`, `migrasi_test.go`,
// `gandawarisan_test.go`) dan folder migrasi Claim Life yang ditanam ke biner.
// Modul ini kini membawa salinannya sendiri; isinya sama. Satu beda yang
// disengaja: `berkasMigrasi` dibaca dari DISK (`../migrations`), bukan dari
// `modul.go` - uji dalam paket repository tidak boleh mengimpor paket perakit
// modulnya sendiri (siklus impor), dan isi berkasnya identik.

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

// akarModul menunjuk folder APP_RNM dari folder paket ini.
const akarModul = "../../../.." // APP_RNM, dari modul/<nama>/backend/repository

// berkasMigrasiDisk membuka folder `migrations/` modul ini sebagai fs.FS.
type berkasMigrasiDisk struct{ fs.FS }

// ReadFile membaca satu berkas, jalur `migrations/<nama>`.
func (b berkasMigrasiDisk) ReadFile(nama string) ([]byte, error) { return fs.ReadFile(b.FS, nama) }

// ReadDir membaca isi folder, jalur `migrations`.
func (b berkasMigrasiDisk) ReadDir(nama string) ([]fs.DirEntry, error) { return fs.ReadDir(b.FS, nama) }

// berkasMigrasi - folder `migrations/` modul ini (050-099, MODUL.md).
var berkasMigrasi = berkasMigrasiDisk{os.DirFS("..")}

// seluruhSQL - isi setiap langkah migrasi modul ini, per nama berkas.
func seluruhSQL(t *testing.T, mundur bool) map[string]string {
	t.Helper()
	langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, m := range langkah {
		out[m.Nama] = strings.Join(m.Pernyataan, "\n")
	}
	return out
}

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

// urutanPenampung - urutan nomor penampung `:n` sebagaimana muncul di teks.
func urutanPenampung(q string) string {
	var got []string
	for _, m := range regexp.MustCompile(`:(\d)`).FindAllStringSubmatch(q, -1) {
		got = append(got, m[1])
	}
	return strings.Join(got, "")
}
