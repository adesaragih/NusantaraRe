package handlers

// Penjaga: SETIAP rute pengubah Claim Life menjawab kasus tertutup dengan 409.
//
// ⛔ Lahir 28-09-2026 (GILIRAN-11, saat menyusun panduan uji layar): enam
// handler - akseptasi, hapus, komite, putaran, tahap, tolak - tidak mengenal
// `ErrKasusSudahTertutup`, sehingga layanan yang MENOLAK kasus tertutup dengan
// benar (penjaga butir bb) sampai ke layar sebagai 500 "gagal ...". Rute DOL
// pernah punya cacat yang sama. Penjaga layanan tidak melihatnya: yang salah
// terjemahannya, bukan penolakannya.
//
// ⛔ PER FUNGSI, bukan per berkas (temuan /code-review GILIRAN-11). Ronde
// pertama mencari nama galatnya di BERKAS pemilik handler - jadi handler
// kedua di berkas yang sama lolos berkat terjemahan milik handler pertama.
// Kini tubuh handler itu sendiri wajib menyebutnya, atau memanggil fungsi
// paket ini yang (secara transitif) menyebutnya - pola `jawabGalat...`.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// fungsiPaket memetakan fungsi tingkat-paket handlers ke tubuhnya.
func fungsiPaket(t *testing.T) map[string]*ast.FuncDecl {
	t.Helper()
	berkas, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	hasil := map[string]*ast.FuncDecl{}
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, b, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Recv == nil && fn.Body != nil {
				hasil[fn.Name.Name] = fn
			}
		}
	}
	return hasil
}

// menerjemahkanTertutup - tubuh fungsi menyebut `services.ErrKasusSudahTertutup`,
// atau memanggil fungsi paket yang menyebutnya.
func menerjemahkanTertutup(nama string, fungsi map[string]*ast.FuncDecl, dilihat map[string]bool) bool {
	fn, ada := fungsi[nama]
	if !ada || dilihat[nama] {
		return false
	}
	dilihat[nama] = true
	ketemu := false
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if ketemu {
			return false
		}
		switch x := n.(type) {
		case *ast.SelectorExpr:
			// Refactor bentuk B (30-09-2026): galatnya kini tinggal di
			// inti/kontrak (dibagi dengan Komite); `services.` tetap dikenali.
			if id, ok := x.X.(*ast.Ident); ok && (id.Name == "services" || id.Name == "kontrak") &&
				x.Sel.Name == "ErrKasusSudahTertutup" {
				ketemu = true
			}
		case *ast.CallExpr:
			if id, ok := x.Fun.(*ast.Ident); ok && menerjemahkanTertutup(id.Name, fungsi, dilihat) {
				ketemu = true
			}
		}
		return !ketemu
	})
	return ketemu
}

// dikecualikanTertutup - handler rute pengubah yang memang tidak menjumpai
// kasus tertutup, beserta alasannya.
var dikecualikanTertutup = map[string]string{
	"daftarKlaim": "pendaftaran MELAHIRKAN kasus; belum ada kasus untuk ditutup",
}

func TestSetiapHandlerPengubahMenjawabKasusTertutup409(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	pola := regexp.MustCompile(`mux\.HandleFunc\(\s*"(POST|PUT|DELETE) /api/klaim-life[^"]*",\s*([a-zA-Z]+)\(`)
	cocok := pola.FindAllStringSubmatch(string(isi), -1)
	if len(cocok) < 10 {
		t.Fatalf("hanya %d rute pengubah terbaca; pembacanya yang rusak", len(cocok))
	}
	fungsi := fungsiPaket(t)
	for _, m := range cocok {
		nama := m[2]
		if _, ada := dikecualikanTertutup[nama]; ada {
			continue
		}
		if _, ada := fungsi[nama]; !ada {
			t.Errorf("%s: fungsi handler tidak ditemukan", nama)
			continue
		}
		if !menerjemahkanTertutup(nama, fungsi, map[string]bool{}) {
			t.Errorf("%s tidak menerjemahkan ErrKasusSudahTertutup (tubuhnya maupun "+
				"fungsi yang dipanggilnya); kasus tertutup akan sampai ke layar sebagai 500", nama)
		}
	}
}
