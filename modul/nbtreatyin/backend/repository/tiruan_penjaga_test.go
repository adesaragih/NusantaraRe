package repository_test

// Penjaga U1 (keputusan WO putaran 3, 04-10-2026) - TANPA Oracle, ikut lari di
// setiap `go test ./...`. Uji db (`*_db_test.go`, tag `db`) boleh membuat
// TIRUAN tabel warisan di skema uji, dengan syarat yang ditegakkan SATU pintu
// `buatTiruan` (tiruan_db_test.go). Penjaga ini membaca sumber uji db itu dan
// gagal bila:
//
//   - teks DDL `CREATE TABLE` / `CREATE VIEW` / `DROP TABLE` muncul di luar
//     `buatTiruan` (pintu lain tidak berpagar);
//   - `buatTiruan` tidak memanggil `config.PagarSkemaUji` (penjaga yang sama
//     dengan `uji/skemauji`) SEBELUM `ExecContext` pertamanya, atau tidak
//     mendaftarkan `t.Cleanup` yang membuang tiruannya;
//   - nama yang diberikan ke `buatTiruan` bukan konstanta berkomentar
//     "tiruan uji, bukan tabel aplikasi".

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

const tandaTiruan = "tiruan uji, bukan tabel aplikasi"

func TestTiruanUjiLewatSatuPintuBerpagar(t *testing.T) {
	berkas, err := filepath.Glob("*_db_test.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(berkas) < 6 {
		t.Fatalf("hanya %d berkas uji db - jalannya salah", len(berkas))
	}
	fset := token.NewFileSet()
	var pintu *ast.FuncDecl
	konstTiruan := map[string]bool{} // nama konstanta -> berkomentar tanda
	var panggilan []*ast.CallExpr
	for _, b := range berkas {
		f, err := parser.ParseFile(fset, b, nil, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			switch x := d.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "buatTiruan" {
					pintu = x
				}
			case *ast.GenDecl:
				if x.Tok != token.CONST {
					continue
				}
				for _, s := range x.Specs {
					vs := s.(*ast.ValueSpec)
					teks := x.Doc.Text() + vs.Doc.Text() + vs.Comment.Text()
					for _, n := range vs.Names {
						konstTiruan[n.Name] = strings.Contains(teks, tandaTiruan)
					}
				}
			}
		}
		// DDL di luar pintu.
		var fungsiKini string
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			fungsiKini = fd.Name.Name
			ast.Inspect(fd, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.BasicLit:
					if x.Kind != token.STRING {
						return true
					}
					s := strings.ToUpper(x.Value)
					for _, ddl := range []string{"CREATE TABLE", "CREATE VIEW", "DROP TABLE"} {
						if strings.Contains(s, ddl) && fungsiKini != "buatTiruan" {
							t.Errorf("%s: %s di %s - tiruan uji hanya lewat buatTiruan (U1)", fset.Position(x.Pos()), ddl, fungsiKini)
						}
					}
				case *ast.CallExpr:
					if id, ok := x.Fun.(*ast.Ident); ok && id.Name == "buatTiruan" {
						panggilan = append(panggilan, x)
					}
				}
				return true
			})
		}
	}
	if pintu == nil {
		t.Fatal("buatTiruan tidak ditemukan di uji db")
	}
	// Urutan di dalam pintu: pagar sebelum DDL pertama; Cleanup membuang.
	var posPagar, posExec, posCleanup token.Pos
	ast.Inspect(pintu.Body, func(n ast.Node) bool {
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := c.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch nama := sel.Sel.Name; {
		case nama == "PagarSkemaUji":
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "config" && posPagar == 0 {
				posPagar = c.Pos()
			}
		case nama == "ExecContext" && posExec == 0:
			posExec = c.Pos()
		case nama == "Cleanup":
			posCleanup = c.Pos()
		}
		return true
	})
	switch {
	case posPagar == 0:
		t.Error("buatTiruan tidak memanggil config.PagarSkemaUji (penjaga uji/skemauji) - POOLDATA tidak tertolak")
	case posExec == 0 || posExec < posPagar:
		t.Error("buatTiruan menjalankan DDL sebelum config.PagarSkemaUji")
	}
	if posCleanup == 0 || posCleanup < posExec {
		t.Error("buatTiruan tidak membuang tiruannya lewat t.Cleanup sesudah membuatnya")
	}
	// Setiap nama tiruan: konstanta bertanda.
	if len(panggilan) == 0 {
		t.Fatal("tidak satu pun pemanggilan buatTiruan - penjaga tidak memeriksa apa pun")
	}
	for _, c := range panggilan {
		if len(c.Args) < 5 {
			t.Errorf("%s: buatTiruan dipanggil dengan %d argumen", fset.Position(c.Pos()), len(c.Args))
			continue
		}
		id, ok := c.Args[4].(*ast.Ident)
		if !ok || !konstTiruan[id.Name] {
			t.Errorf("%s: nama tiruan harus konstanta berkomentar %q", fset.Position(c.Args[4].Pos()), tandaTiruan)
		}
	}
}
