package penjaga

// Penjaga IMPOR lintas modul - refactor bentuk B paket 8 (30-09-2026),
// diperbarui struktur tim satu folder per modul (30-09-2026).
//
// Aturannya, ditegakkan atas SETIAP berkas .go (termasuk berkas uji dan berkas
// bertag `db` - pengurai tidak menimbang tag build):
//
//  1. `modul/X/...` hanya mengimpor `inti/...` dan `modul/X/...`. Ketergantungan
//     lintas modul lewat antarmuka di `inti/backend/kontrak`, dinyatakan di
//     `Pendaftaran()` `backend/modul.go` dan disambung perakit `inti.Rakit`.
//     Modul TIDAK mengimpor daftar modul `inti/backend/daftar` (impor
//     melingkar, dan modul tidak mengenal modul lain). Berkas UJI sebuah modul
//     boleh memakai `uji/skemauji` - pengecualian BERNAMA: paket itu sengaja
//     mengenal semua modul, sebab ia membangun skema uji UTUH (migrasi semua
//     modul, fixture Claim Life dan PremiumList). Ia satu-satunya jalur lintas
//     modul yang disahkan, hanya untuk berkas uji; `uji/lintasmodul` tidak.
//  2. `inti/...` hanya mengimpor `inti/...` - inti tidak mengenal modul. SATU
//     pengecualian: daftar bangkitan `inti/backend/daftar` (paket itu saja,
//     bukan subfoldernya) mengimpor PAKET AKAR setiap modul,
//     `modul/<nama>/backend` - bukan paket di dalamnya.
//  3. `cmd/...` hanya mengimpor `inti/...` - ia memasang modul dari daftar
//     `inti/backend/daftar`, tidak pernah satu modul langsung.
//
// Dan satu aturan LETAK: kode Go sebuah modul hanya di `modul/<nama>/backend/`.
// `uji/...` sengaja bebas: ia memang tempat yang mengenal semua modul.
//
// ⛔ Nol nama modul sungguhan di berkas ini: kasus uji memakai `alfa`/`beta`.

import (
	"go/parser"
	"go/token"
	"os"
	"path"
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
		if path.Dir(rel) == folderDaftar && akarPaketModul(dalam) {
			return "" // daftar bangkitan: satu impor paket akar per modul
		}
		if !keInti {
			return "inti tidak mengenal modul: hanya boleh mengimpor inti/... (kecuali daftar bangkitan " +
				folderDaftar + ", dan hanya paket akar modul/<nama>/backend)"
		}
	case strings.HasPrefix(rel, "cmd/"):
		if !keInti {
			return "cmd hanya mengimpor inti/...; modul dipasang lewat daftar " + folderDaftar
		}
	case strings.HasPrefix(rel, "modul/"):
		milik := "modul/" + strings.SplitN(strings.TrimPrefix(rel, "modul/"), "/", 2)[0]
		switch {
		case dalam == folderDaftar:
			return "modul tidak mengimpor daftar modul " + folderDaftar + " (impor melingkar; modul tidak mengenal modul lain)"
		case keInti, dalam == milik, strings.HasPrefix(dalam, milik+"/"):
		case uji && dalam == "uji/skemauji":
		default:
			return "modul hanya mengimpor inti/... dan dirinya sendiri; ketergantungan lintas modul lewat " +
				"inti/backend/kontrak, dinyatakan di Pendaftaran() backend/modul.go"
		}
	}
	return ""
}

// folderDaftar - paket daftar modul bangkitan.
const folderDaftar = "inti/backend/daftar"

// akarPaketModul menjawab apakah `dalam` (jalur impor tanpa `nusantarare/`)
// adalah PAKET AKAR sebuah modul, `modul/<nama>/backend`.
func akarPaketModul(dalam string) bool {
	b := strings.Split(dalam, "/")
	return len(b) == 3 && b[0] == "modul" && b[1] != "" && b[2] == "backend"
}

// pelanggaranLetak menjawab mengapa berkas Go `rel` salah letak; kosong = benar.
func pelanggaranLetak(rel string) string {
	if !strings.HasPrefix(rel, "modul/") {
		return ""
	}
	b := strings.Split(rel, "/")
	if len(b) < 4 || b[2] != "backend" {
		return "kode Go modul tinggal di modul/<nama>/backend/ - bukan langsung di modul/ atau di sebelah frontend/ dan docs/"
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
		if alasan := pelanggaranLetak(rel); alasan != "" {
			t.Errorf("%s - %s", rel, alasan)
		}
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
		{"modul/alfa/backend/services/x.go", "nusantarare/modul/beta/backend/models", false},
		{"modul/alfa/backend/services/x_test.go", "nusantarare/modul/beta/backend/models", false},
		{"modul/alfa/backend/services/x.go", "nusantarare/modul/beta/backend", false},
		{"modul/alfa/backend/services/x.go", "nusantarare/inti/backend/daftar", false},
		{"modul/alfa/backend/modul.go", "nusantarare/inti/backend/daftar", false},
		{"modul/alfa/backend/services/x.go", "nusantarare/uji/skemauji", false},
		{"modul/alfa/backend/services/x_test.go", "nusantarare/uji/lintasmodul", false},
		{"modul/alfa/backend/services/x_test.go", "nusantarare/uji/skemauji", true},
		{"modul/alfa/backend/services/x.go", "nusantarare/modul/alfa/backend/models", true},
		{"modul/alfa/backend/modul.go", "nusantarare/modul/alfa/backend/handlers", true},
		// Tabrakan awalan: `alfab` diawali nama modul ini, tetapi modul LAIN.
		{"modul/alfa/backend/services/x.go", "nusantarare/modul/alfab/backend/models", false},
		{"modul/alfa/backend/services/x.go", "nusantarare/inti/backend/kontrak", true},
		{"modul/alfa/backend/services/x.go", "nusantarare/inti/backend", true},
		// Daftar bangkitan: paket akar modul saja, dan hanya dari paket daftar.
		{"inti/backend/daftar/modul_alfa_gen.go", "nusantarare/modul/alfa/backend", true},
		{"inti/backend/daftar/modul_alfa_gen.go", "nusantarare/modul/alfa/backend/services", false},
		{"inti/backend/daftar/bangkit/main.go", "nusantarare/modul/alfa/backend", false},
		{"inti/backend/kontrak/x.go", "nusantarare/modul/alfa/backend", false},
		{"inti/backend/kontrak/x.go", "nusantarare/modul/alfa/backend/models", false},
		{"inti/backend/penjaga/x_test.go", "nusantarare/uji/skemauji", false},
		{"inti/backend/db/x.go", "nusantarare/inti/backend/galat", true},
		{"cmd/api/main.go", "nusantarare/inti/backend/daftar", true},
		{"cmd/api/main.go", "nusantarare/modul/alfa/backend", false},
		{"cmd/api/main.go", "nusantarare/modul/alfa/backend/handlers", false},
		{"uji/lintasmodul/x_test.go", "nusantarare/modul/alfa/backend/services", true},
		{"modul/alfa/backend/services/x.go", "github.com/sijms/go-ora/v2", true},
	} {
		if dapat := pelanggaranImpor(k.rel, k.imp) == ""; dapat != k.boleh {
			t.Errorf("%s -> %s: boleh=%v, mau %v", k.rel, k.imp, dapat, k.boleh)
		}
	}
	// Aturan letak, dua arah.
	for rel, benar := range map[string]bool{
		"modul/alfa/backend/modul.go":           true,
		"modul/alfa/backend/services/x.go":      true,
		"modul/daftar.go":                       false,
		"modul/alfa/modul.go":                   false,
		"modul/alfa/frontend/x.go":              false,
		"inti/backend/daftar/modul_alfa_gen.go": true,
	} {
		if dapat := pelanggaranLetak(rel) == ""; dapat != benar {
			t.Errorf("letak %s: benar=%v, mau %v", rel, dapat, benar)
		}
	}
}
