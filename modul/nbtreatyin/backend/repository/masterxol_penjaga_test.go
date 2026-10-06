package repository

// Penjaga F1 (keputusan WO putaran 3, 04-10-2026): SATU-SATUNYA pembaca JSON
// master adalah `MasterXOLDariJSON`. Tanpa Oracle.
//
//   - `TestKolomDokumenHanyaDiMasterXOLDariJSON`: kolom dokumen `JSONDATA` tidak
//     dipakai di SQL/kode berkas lain mana pun di modul ini - hanya disebut di
//     KOMENTAR (mis. `models/halaman.go`, `repository/acuan.go`,
//     `services/gudang.go`, yang menjelaskan bahwa kolom itu TIDAK dibaca).
//     Pembedaan komentar lawan kode memakai `go/scanner` (token COMMENT
//     dilewati), dan dibuktikan oleh `TestPemindaiKolomDokumenMembedakanKomentar`.
//     Migrasi `.sql` diperiksa sesudah komentar `--` dibuang.
//   - `TestMasterXOLDariJSONHanyaMengembalikanHasilUrai`: setiap `return`
//     `MasterXOLDariJSON` mengembalikan `models.MasterXOL{}` kosong (galat)
//     atau hasil `uraiMasterXOL` - penyaring daftar medan tertutup. Dengan itu
//     uji `TestUraiMasterXOLHanyaMedanDaftarTertutup` (tanpa Oracle) berlaku
//     bagi keluaran `MasterXOLDariJSON`.

import (
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// kataKolomDokumen - nama kolom dokumen JSON master, dirangkai supaya berkas
// penjaga ini sendiri tidak memuatnya sebagai literal utuh.
var kataKolomDokumen = "JSON" + "DATA"

// izinKolomDokumen - berkas (relatif ke backend/, garis miring maju) yang BOLEH
// memuat kata itu di luar komentar, dan alasannya.
var izinKolomDokumen = map[string]string{
	// satu-satunya pembaca - dibatasi lagi ke badan fungsi MasterXOLDariJSON
	"repository/masterxol.go": "MasterXOLDariJSON",
	// tiruan uji M_TREATY_IN / M_TREATY_IN_EDM bagi uji db pembaca itu sendiri
	"repository/masterxol_db_test.go": "",
}

// pemakaianKolomDokumen mengembalikan posisi token NON-komentar (literal teks,
// pengenal) yang memuat kata kolom dokumen, tanpa beda huruf besar/kecil.
func pemakaianKolomDokumen(nama string, src []byte) ([]token.Position, error) {
	fset := token.NewFileSet()
	f := fset.AddFile(nama, fset.Base(), len(src))
	var galat scanner.ErrorList
	var s scanner.Scanner
	s.Init(f, src, func(pos token.Position, msg string) { galat.Add(pos, msg) }, 0) // tanpa ScanComments: komentar tidak pernah jadi token
	var out []token.Position
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if lit != "" && strings.Contains(strings.ToUpper(lit), kataKolomDokumen) {
			out = append(out, fset.Position(pos))
		}
	}
	return out, galat.Err()
}

// rentangFungsi - posisi awal dan akhir (offset) badan satu fungsi/metode.
func rentangFungsi(t *testing.T, berkas, fungsi string) (int, int) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, berkas, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == fungsi {
			return fset.Position(fd.Pos()).Offset, fset.Position(fd.End()).Offset
		}
	}
	t.Fatalf("fungsi %s tidak ditemukan di %s", fungsi, berkas)
	return 0, 0
}

func TestPemindaiKolomDokumenMembedakanKomentar(t *testing.T) {
	kata := kataKolomDokumen
	for _, k := range []struct {
		nama   string
		src    string
		jumlah int
	}{
		{"komentar baris", "package x\n// dibaca dari view, bukan " + kata + "\n", 0},
		{"komentar blok", "package x\n/* select " + kata + " from t */\n", 0},
		{"literal SQL", "package x\nconst q = `SELECT " + kata + " FROM T`\n", 1},
		{"literal huruf kecil", "package x\nvar q = \"select a." + strings.ToLower(kata) + " from t\"\n", 1},
		{"pengenal", "package x\nvar " + kata + " = 1\n", 1},
		{"komentar dan kode", "package x\n// " + kata + "\nvar q = \"" + kata + "\"\n", 1},
	} {
		t.Run(k.nama, func(t *testing.T) {
			pos, err := pemakaianKolomDokumen("x.go", []byte(k.src))
			if err != nil {
				t.Fatal(err)
			}
			if len(pos) != k.jumlah {
				t.Fatalf("%d pemakaian, harap %d: %v", len(pos), k.jumlah, pos)
			}
		})
	}
}

func TestKolomDokumenHanyaDiMasterXOLDariJSON(t *testing.T) {
	akar := filepath.Join("..") // backend/
	awalPembaca, akhirPembaca := rentangFungsi(t, filepath.Join(akar, "repository", "masterxol.go"), "MasterXOLDariJSON")
	diperiksa := 0
	err := filepath.WalkDir(akar, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(akar, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		isi, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		switch filepath.Ext(p) {
		case ".go":
			diperiksa++
			pos, err := pemakaianKolomDokumen(rel, isi)
			if err != nil {
				return err
			}
			fungsi, diizinkan := izinKolomDokumen[rel]
			for _, ps := range pos {
				switch {
				case !diizinkan:
					t.Errorf("%s:%d memakai %s di luar komentar - satu-satunya pembaca JSON master adalah repository.MasterXOLDariJSON (F1)",
						rel, ps.Line, kataKolomDokumen)
				case fungsi != "" && (ps.Offset < awalPembaca || ps.Offset > akhirPembaca):
					t.Errorf("%s:%d memakai %s di luar badan %s (F1)", rel, ps.Line, kataKolomDokumen, fungsi)
				}
			}
		case ".sql":
			diperiksa++
			for i, b := range strings.Split(string(isi), "\n") {
				if j := strings.Index(b, "--"); j >= 0 {
					b = b[:j]
				}
				if strings.Contains(strings.ToUpper(b), kataKolomDokumen) {
					t.Errorf("%s:%d memakai %s di SQL migrasi (F1)", rel, i+1, kataKolomDokumen)
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 50 { // backend/ memuat >100 berkas .go dan 18 .sql
		t.Fatalf("hanya %d berkas diperiksa - jalannya salah, bukan modulnya bersih", diperiksa)
	}
}

// F1 (tinjauan spec P3 (c)2): layar modul tidak pernah menyentuh kolom dokumen
// JSON master - nol kemunculan di `modul/nbtreatyin/frontend/` (kode maupun
// komentar, huruf besar maupun kecil). Data master sampai ke layar hanya lewat
// halaman kerja hasil `MasterXOLDariJSON` (daftar medan tertutup).
func TestFrontendTanpaKolomDokumen(t *testing.T) {
	akar := filepath.Join("..", "..", "frontend")
	diperiksa := 0
	err := filepath.WalkDir(akar, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		switch filepath.Ext(p) {
		case ".ts", ".tsx", ".css", ".json":
		default:
			return nil
		}
		diperiksa++
		isi, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for i, b := range strings.Split(string(isi), "\n") {
			if strings.Contains(strings.ToUpper(b), kataKolomDokumen) {
				t.Errorf("%s:%d memakai %s - layar tidak membaca kolom dokumen master (F1)", filepath.ToSlash(p), i+1, kataKolomDokumen)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 20 { // frontend/ memuat >30 berkas .ts/.tsx
		t.Fatalf("hanya %d berkas diperiksa - jalannya salah, bukan layarnya bersih", diperiksa)
	}
}

func TestMasterXOLDariJSONHanyaMengembalikanHasilUrai(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "masterxol.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var fd *ast.FuncDecl
	for _, d := range f.Decls {
		if x, ok := d.(*ast.FuncDecl); ok && x.Name.Name == "MasterXOLDariJSON" {
			fd = x
		}
	}
	if fd == nil {
		t.Fatal("MasterXOLDariJSON tidak ditemukan")
	}
	jumlahUrai := 0
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		if _, ok := n.(*ast.FuncLit); ok {
			return false // return di dalam closure bukan return fungsi ini
		}
		r, ok := n.(*ast.ReturnStmt)
		if !ok || len(r.Results) == 0 {
			return true
		}
		switch x := r.Results[0].(type) {
		case *ast.CompositeLit:
			if len(x.Elts) == 0 {
				return true // models.MasterXOL{} bersama galat
			}
		case *ast.CallExpr:
			if id, ok := x.Fun.(*ast.Ident); ok && id.Name == "uraiMasterXOL" {
				jumlahUrai++
				return true
			}
		}
		t.Errorf("masterxol.go:%d: return MasterXOLDariJSON bukan models.MasterXOL{} dan bukan uraiMasterXOL(...) - "+
			"keluaran yang tidak lewat penyaring daftar medan tertutup (F1)", fset.Position(r.Pos()).Line)
		return true
	})
	if jumlahUrai != 1 {
		t.Errorf("MasterXOLDariJSON memanggil uraiMasterXOL di %d return, harap tepat 1", jumlahUrai)
	}
}
