package penjaga

// Penjaga IMPOR lintas modul - refactor bentuk B paket 8 (30-09-2026).
//
// Aturan bentuk B, ditegakkan atas SETIAP berkas .go (termasuk berkas uji dan
// berkas bertag `db` - pengurai tidak menimbang tag build):
//
//  1. `modul/X/...` hanya mengimpor `inti/...` dan `modul/X/...`. Ketergantungan
//     lintas modul lewat antarmuka di `inti/kontrak`, disambung di daftar
//     modul (`modul/daftar.go`). Berkas UJI sebuah modul boleh memakai
//     penunjang uji netral `uji/skemauji` (skema Oracle tiruan).
//  2. `inti/...` hanya mengimpor `inti/...` - inti tidak mengenal modul.
//  3. `cmd/...` hanya mengimpor `inti/...` dan daftar modul `nusantarare/modul`
//     - ia memasang modul dari daftar, tidak pernah satu modul langsung.
//
// `modul/daftar.go` (paket `nusantarare/modul`) dan `uji/...` sengaja bebas:
// keduanya memang tempat yang mengenal semua modul.

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const jalurModulGo = "nusantarare"

// pelanggaranImpor menjawab mengapa berkas `rel` (jalur relatif akar aplikasi,
// bergaris-miring) tidak boleh mengimpor `imp`; kosong = boleh.
func pelanggaranImpor(rel, imp string) string {
	if imp != jalurModulGo && !strings.HasPrefix(imp, jalurModulGo+"/") {
		return "" // pustaka standar dan pihak ketiga di luar aturan ini
	}
	dalam := strings.TrimPrefix(strings.TrimPrefix(imp, jalurModulGo), "/")
	keInti := dalam == "inti" || strings.HasPrefix(dalam, "inti/")
	uji := strings.HasSuffix(rel, "_test.go")
	switch {
	case strings.HasPrefix(rel, "inti/"):
		if !keInti {
			return "inti tidak mengenal modul: hanya boleh mengimpor inti/..."
		}
	case strings.HasPrefix(rel, "cmd/"):
		if !keInti && dalam != "modul" {
			return "cmd memasang modul dari daftar (nusantarare/modul), bukan satu modul langsung"
		}
	case strings.HasPrefix(rel, "modul/"):
		bagian := strings.SplitN(strings.TrimPrefix(rel, "modul/"), "/", 2)
		if len(bagian) < 2 {
			return "" // modul/daftar.go - daftar modul, sengaja mengenal semuanya
		}
		milik := "modul/" + bagian[0]
		switch {
		case keInti, dalam == milik, strings.HasPrefix(dalam, milik+"/"):
		case uji && dalam == "uji/skemauji":
		default:
			return "modul hanya mengimpor inti/... dan dirinya sendiri; ketergantungan lintas modul lewat inti/kontrak, disambung di modul/daftar.go"
		}
	}
	return ""
}

// imporSetiapBerkasGo - setiap berkas .go di bawah inti/, modul/, cmd/, uji/,
// beserta impornya, berkunci jalur relatif akar aplikasi.
func imporSetiapBerkasGo(t *testing.T) map[string][]string {
	t.Helper()
	hasil := map[string][]string{}
	fset := token.NewFileSet()
	for _, akar := range []string{"inti", "modul", "cmd", "uji"} {
		err := filepath.Walk(filepath.Join(akarAplikasi, akar), func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(p, ".go") {
				return nil
			}
			f, err := parser.ParseFile(fset, p, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(akarAplikasi, p)
			if err != nil {
				return err
			}
			var imp []string
			for _, s := range f.Imports {
				jalur, err := strconv.Unquote(s.Path.Value)
				if err != nil {
					return err
				}
				imp = append(imp, jalur)
			}
			hasil[filepath.ToSlash(rel)] = imp
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	return hasil
}

func TestModulTidakMengimporModulLain(t *testing.T) {
	berkas := imporSetiapBerkasGo(t)
	var urut []string
	for rel := range berkas {
		urut = append(urut, rel)
	}
	sort.Strings(urut)
	diperiksa := map[string]int{}
	for _, rel := range urut {
		lapis := strings.SplitN(rel, "/", 2)[0]
		diperiksa[lapis]++
		for _, imp := range berkas[rel] {
			if alasan := pelanggaranImpor(rel, imp); alasan != "" {
				t.Errorf("%s mengimpor %s - %s", rel, imp, alasan)
			}
		}
	}
	t.Logf("berkas .go diperiksa per lapis: %v", diperiksa)
	// ⛔ Penjaga yang membaca nol berkas di satu lapis lulus atas apa pun.
	for _, lapis := range []string{"inti", "modul", "cmd", "uji"} {
		if diperiksa[lapis] == 0 {
			t.Fatalf("nol berkas .go di %s/; pembacanya yang rusak", lapis)
		}
	}
	if diperiksa["modul"] < 100 {
		t.Fatalf("hanya %d berkas modul terbaca; pembacanya yang rusak", diperiksa["modul"])
	}
}

// TestAturanImporMenggigit - aturannya sendiri diuji dua arah, supaya penjaga
// di atas tidak lulus karena aturannya longgar.
func TestAturanImporMenggigit(t *testing.T) {
	for _, k := range []struct {
		rel, imp string
		boleh    bool
	}{
		{"modul/claimlife/services/x.go", "nusantarare/modul/komite/models", false},
		{"modul/claimlife/services/x_test.go", "nusantarare/modul/premiumlist/models", false},
		{"modul/komite/services/x.go", "nusantarare/modul", false},
		{"modul/komite/services/x.go", "nusantarare/uji/skemauji", false},
		{"modul/komite/services/x_test.go", "nusantarare/uji/lintasmodul", false},
		{"modul/komite/services/x_test.go", "nusantarare/uji/skemauji", true},
		{"modul/komite/services/x.go", "nusantarare/modul/komite/models", true},
		{"modul/komite/modul.go", "nusantarare/modul/komite/handlers", true},
		{"modul/komite/services/x.go", "nusantarare/modul/komiteb/models", false},
		{"modul/komite/services/x.go", "nusantarare/inti/kontrak", true},
		{"modul/daftar.go", "nusantarare/modul/treaty/services", true},
		{"inti/kontrak/x.go", "nusantarare/modul/claimlife/models", false},
		{"inti/penjaga/x_test.go", "nusantarare/uji/skemauji", false},
		{"inti/db/x.go", "nusantarare/inti/galat", true},
		{"cmd/api/main.go", "nusantarare/modul", true},
		{"cmd/api/main.go", "nusantarare/modul/treaty/handlers", false},
		{"uji/lintasmodul/x_test.go", "nusantarare/modul/komite/services", true},
		{"modul/claimlife/services/x.go", "github.com/sijms/go-ora/v2", true},
	} {
		if dapat := pelanggaranImpor(k.rel, k.imp) == ""; dapat != k.boleh {
			t.Errorf("%s -> %s: boleh=%v, mau %v", k.rel, k.imp, dapat, k.boleh)
		}
	}
}
