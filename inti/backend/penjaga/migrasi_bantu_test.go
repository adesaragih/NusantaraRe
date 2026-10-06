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
	"regexp"
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

// Pola baris pembuka pernyataan di teks MENTAH berkas migrasi - cara hitung
// kedua, tanpa pemisah pernyataan pelari.
var (
	polaBarisCreate      = regexp.MustCompile(`(?m)^CREATE\s`)
	polaBarisCreateTable = regexp.MustCompile(`(?m)^CREATE TABLE\s`)
)

// cacahBarisMentah menghitung baris yang cocok `pola` di setiap berkas migrasi
// MAJU, dibaca apa adanya dari disk.
func cacahBarisMentah(t *testing.T, pola *regexp.Regexp) int {
	t.Helper()
	if berkasMigrasi.muat != nil {
		t.Fatal(berkasMigrasi.muat)
	}
	n := 0
	for nama, jalur := range berkasMigrasi.asal {
		if strings.HasSuffix(nama, "_down.sql") {
			continue
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			t.Fatal(err)
		}
		n += len(pola.FindAllIndex(isi, -1))
	}
	return n
}

// ===========================================================================
// Nama tabel yang DIBUAT atau DIGANTI-NAMAI oleh migrasi Treaty In
// ===========================================================================
//
// ⛔ Lahir dari tabrakan dua keputusan, 6 Oktober 2026 — dan keduanya sah:
//
//	tco4, 29-09-2026   Treaty Contract Out tidak boleh punya tabel baru.
//	                   Penjaganya mengenali "tabel baru" dari NAMANYA:
//	                   `T_` + `TREATY…`/`MTREATY…`/`PROPORTIONAL…`.
//
//	Treaty In, 05-10   Pemilik proses menetapkan nama tabel Treaty In
//	                   mengikuti `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2
//	                   .xlsx`. Migrasi 436 menamai ulang delapan tabel
//	                   menjadi `T_TREATY_*`; 437 membuat tiga belas lagi.
//
// Sesudah 436, nama `T_TREATY_*` TIDAK LAGI menandakan Treaty Contract Out —
// dan penjaga yang menuduh tabel Treaty In sebagai pelanggaran tco4 menuduh
// hal yang benar sebagai hal yang salah. Berkas `lintasaplikasi_test.go`
// sendiri menuliskan akibatnya: penjaga semacam itu DILONGGARKAN orang,
// bukan dipatuhi.
//
// ⭐ Jadi ia DIPERSEMPIT, bukan dilonggarkan: tco4 tetap menolak setiap tabel
// baru di mana pun, KECUALI nama yang sungguh dibuat migrasi Treaty In. Yang
// memutuskan bukan letak berkasnya saja, melainkan daftar ini — supaya kode
// Treaty In yang menyebut tabel Treaty Contract Out tetap tertangkap.
func namaTabelTreatyIn(t *testing.T) map[string]bool {
	t.Helper()
	polaBuat := regexp.MustCompile(`(?i)CREATE\s+(?:TABLE|SEQUENCE|INDEX)\s+\{skema\}\.(\w+)`)
	// `ALTER TABLE … RENAME TO X` dan `RENAME A TO X` (sequence; Oracle tidak
	// punya `ALTER SEQUENCE … RENAME TO`).
	polaNamai := regexp.MustCompile(`(?i)RENAME\s+(?:\w+\s+)?TO\s+(\w+)`)
	out := map[string]bool{}
	for nama, teks := range seluruhSQL(t, false) {
		if berkasMigrasi.modul(nama) != "treatyin" {
			continue
		}
		for _, m := range polaBuat.FindAllStringSubmatch(teks, -1) {
			out[strings.ToUpper(m[1])] = true
		}
		for _, m := range polaNamai.FindAllStringSubmatch(teks, -1) {
			out[strings.ToUpper(m[1])] = true
		}
	}
	if len(out) == 0 {
		t.Fatal("nol nama tabel Treaty In terbaca; pembacanya yang rusak")
	}
	return out
}

// milikTreatyIn menjawab apakah `nama` adalah nama tabel Treaty In - persis,
// atau sebuah AWALAN dari salah satunya.
func milikTreatyIn(daftar map[string]bool, nama string) bool {
	n := strings.ToUpper(nama)
	if daftar[n] {
		return true
	}
	for ada := range daftar {
		if strings.HasPrefix(ada, n) {
			return true
		}
	}
	return false
}
