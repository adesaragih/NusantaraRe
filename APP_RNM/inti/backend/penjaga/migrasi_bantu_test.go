package penjaga

// Pembantu penjaga migrasi lintas modul - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): penjaga higiene migrasi (BOM, CR, jalur
// mundur, skema eksplisit, kaskade terdaftar, ...) dulu membaca folder
// `migrations/` Claim Life, yang ketika itu memuat migrasi SEMUA modul. Sejak
// tiap modul membawa foldernya sendiri, penjaga di paket ini membaca
// GABUNGAN seluruh folder `migrations/` di bawah akar aplikasi - dilihat
// sebagai satu folder `migrations/`, isinya sama dengan folder tunggal lama.
//
// ⛔ Penghitung per modul TIDAK tinggal di sini: penjaga di paket ini hanya
// memuat aturan yang berlaku untuk setiap berkas migrasi, sehingga modul yang
// menambah migrasi tidak pernah dipaksa menyunting berkas uji ini.

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

// gabunganMigrasi memetakan `migrations/<nama>` ke berkasnya di disk.
type gabunganMigrasi struct {
	asal map[string]string
	muat error
}

// berkasMigrasi - sumber yang dibaca penjaga-penjaga migrasi di paket ini.
var berkasMigrasi = muatGabunganMigrasi()

func muatGabunganMigrasi() *gabunganMigrasi {
	g := &gabunganMigrasi{asal: map[string]string{}}
	g.muat = filepath.Walk(akarAplikasi, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		switch info.Name() {
		case "frontend", "node_modules", "bin", ".git", "dist", "unggahan":
			return filepath.SkipDir
		case "migrations":
			entri, err := os.ReadDir(p)
			if err != nil {
				return err
			}
			for _, e := range entri {
				if e.IsDir() {
					continue
				}
				if lama, ganda := g.asal[e.Name()]; ganda {
					return fmt.Errorf("berkas migrasi %s ada di dua folder: %s dan %s", e.Name(), lama, p)
				}
				g.asal[e.Name()] = filepath.Join(p, e.Name())
			}
			return filepath.SkipDir
		}
		return nil
	})
	if g.muat == nil && len(g.asal) == 0 {
		g.muat = fmt.Errorf("nol berkas migrasi di bawah %s; pembacanya yang rusak", akarAplikasi)
	}
	return g
}

func (g *gabunganMigrasi) jalur(nama string) (string, error) {
	if g.muat != nil {
		return "", g.muat
	}
	dir, dasar := path.Split(nama)
	if j, ada := g.asal[dasar]; ada && dir == "migrations/" {
		return j, nil
	}
	return "", &fs.PathError{Op: "open", Path: nama, Err: fs.ErrNotExist}
}

// Open memenuhi fs.FS.
func (g *gabunganMigrasi) Open(nama string) (fs.File, error) {
	j, err := g.jalur(nama)
	if err != nil {
		return nil, err
	}
	return os.Open(j)
}

// ReadFile memenuhi fs.ReadFileFS.
func (g *gabunganMigrasi) ReadFile(nama string) ([]byte, error) {
	j, err := g.jalur(nama)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(j)
}

// ReadDir memenuhi fs.ReadDirFS - hanya untuk folder `migrations`.
func (g *gabunganMigrasi) ReadDir(nama string) ([]fs.DirEntry, error) {
	if g.muat != nil {
		return nil, g.muat
	}
	if nama != "migrations" {
		return nil, &fs.PathError{Op: "readdir", Path: nama, Err: fs.ErrNotExist}
	}
	var out []fs.DirEntry
	for _, j := range g.asal {
		info, err := os.Stat(j)
		if err != nil {
			return nil, err
		}
		out = append(out, fs.FileInfoToDirEntry(info))
	}
	sort.Slice(out, func(i, k int) bool { return out[i].Name() < out[k].Name() })
	return out, nil
}

// seluruhSQL - isi setiap langkah migrasi semua modul, per nama berkas.
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

// gabungSemua - seluruh pernyataan maju semua modul, disambung.
func gabungSemua(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, isi := range seluruhSQL(t, false) {
		b.WriteString(isi)
		b.WriteString("\n")
	}
	return b.String()
}

// modul menyebut modul pemilik sebuah berkas migrasi, dibaca dari letaknya:
// `modul/<nama>/.../migrations/` milik `<nama>`, `inti/.../migrations/` milik
// `inti`.
//
// Refactor bentuk B paket 8: menggantikan `milikClaimLife` yang membagi
// menurut NOMOR (`nama < "050_"`). Pemilik kini dibaca dari tempat berkasnya,
// sehingga penjaga per modul tidak bergantung pada rentang nomor modul lain.
//
// Struktur tim satu folder per modul (30-09-2026): folder migrasi kini
// `modul/<nama>/backend/migrations/` dan `inti/backend/migrations/`, jadi
// pemilik diambil dari segmen jalur sesudah akar aplikasi - bukan dari nama
// folder di atas `migrations/` (yang kini selalu `backend`).
func (g *gabunganMigrasi) modul(nama string) string {
	j, ada := g.asal[nama]
	if !ada {
		return ""
	}
	return pemilikJalur(j)
}

// pemilikJalur menyebut pemilik sebuah berkas di bawah akar aplikasi: nama
// modul untuk `modul/<nama>/...`, `inti` untuk `inti/...`, kosong selainnya.
func pemilikJalur(j string) string {
	rel, err := filepath.Rel(akarAplikasi, j)
	if err != nil {
		return ""
	}
	bagian := strings.Split(filepath.ToSlash(rel), "/")
	switch {
	case len(bagian) >= 3 && bagian[0] == "modul":
		return bagian[1]
	case len(bagian) >= 2 && bagian[0] == "inti":
		return "inti"
	}
	return ""
}
