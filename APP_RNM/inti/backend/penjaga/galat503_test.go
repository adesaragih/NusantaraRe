package penjaga

// Kontrak 503 dua sisi, SISI BACKEND - penjaga statik lintas modul.
//
// Refactor bentuk B (30-09-2026): dipindah dari `internal/handlers/
// galat503_test.go`. Dua uji perilakunya (rute kurs Treaty) ikut modul
// Treaty; penjaga statiknya tinggal di sini dan kini membaca folder handlers
// SEMUA modul - dulu `filepath.Glob("*.go")` di satu folder, yang akan
// menyempit diam-diam begitu modul pertama pindah.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// kode503 - `http.StatusServiceUnavailable` atau literal `503`.
func kode503(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.SelectorExpr:
		x, ok := v.X.(*ast.Ident)
		return ok && x.Name == "http" && v.Sel.Name == "StatusServiceUnavailable"
	case *ast.BasicLit:
		return v.Kind == token.INT && v.Value == "503"
	}
	return false
}

// panggilGalatTulis mengenali `galat.Tulis(...)` - juga bila paketnya diimpor
// dengan alias `intigalat` karena berkasnya memakai nama `galat` untuk hal lain.
func panggilGalatTulis(fun ast.Expr) bool {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Tulis" {
		return false
	}
	x, ok := sel.X.(*ast.Ident)
	return ok && (x.Name == "galat" || x.Name == "intigalat")
}

// Penjaga statik: SETIAP 503 di handler ditulis lewat `galat()` dengan pesan
// tak kosong. 503 yang ditulis `http.Error` (teks biasa) atau `WriteHeader`
// polos sampai di klien tanpa `galat` - dan klien benar menyebutnya "backend
// tidak terhubung".
//
// Dibaca sebagai pohon sintaks (seperti `fungsiPaket`, `tertutup409_test.go`),
// bukan baris: argumen yang terpecah ke baris berikutnya, komentar, dan nama
// fungsi yang kebetulan berakhiran `galat` tidak mengecohnya.
// ⚠️ Batasnya: kode status di dalam VARIABEL (`galat(w, kode, …)`) tidak
// terlihat - penjaga ini membaca konstanta dan literal saja.
func TestSetiap503HandlerLewatGalat(t *testing.T) {
	berkas := berkasHandler(t)
	fset := token.NewFileSet()
	jumlah := 0
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, b, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		sah := map[ast.Expr]bool{}
		ast.Inspect(f, func(n ast.Node) bool {
			c, ok := n.(*ast.CallExpr)
			if !ok || len(c.Args) != 3 || !kode503(c.Args[1]) {
				return true
			}
			// Refactor bentuk B (30-09-2026): `galat()` kini `galat.Tulis()`
			// dari paket bersama `inti/backend/galat`.
			if !panggilGalatTulis(c.Fun) {
				return true
			}
			sah[c.Args[1]] = true
			jumlah++
			if lit, ok := c.Args[2].(*ast.BasicLit); ok && lit.Kind == token.STRING && (lit.Value == `""` || lit.Value == "``") {
				t.Errorf("%s: 503 dengan galat kosong - klien membacanya sebagai backend mati", fset.Position(c.Pos()))
			}
			return true
		})
		ast.Inspect(f, func(n ast.Node) bool {
			if e, ok := n.(ast.Expr); ok && kode503(e) && !sah[e] {
				t.Errorf("%s: 503 tidak lewat galat()", fset.Position(e.Pos()))
			}
			return true
		})
	}
	if jumlah < 10 {
		t.Fatalf("hanya %d jawaban 503 lewat galat() terbaca; pembacanya yang rusak", jumlah)
	}
	t.Logf("%d jawaban 503 lewat galat()", jumlah)
}
