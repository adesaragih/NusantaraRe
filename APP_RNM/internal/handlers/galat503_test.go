package handlers

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/internal/services"
)

// Kontrak 503 dua sisi, SISI BACKEND - lanjutan 6 Treaty Contract Out.
//
// ⛔ Sebab uji ini ada, dan itu sungguh terjadi (laporan work owner
// 29-09-2026): `Show TreatyDesc` menampilkan "Backend tidak terhubung" padahal
// backend menyala. Rute kurs menjawab 503 `{"galat": …}` - penolakan backend
// sendiri - dan klien menganggap setiap 503 datang dari proxy yang kehilangan
// upstream-nya. Kini klien membedakan keduanya dari BADANNYA: 503 yang
// membawa `galat` tak kosong adalah jawaban backend. Kontrak itu hanya
// bertahan selama setiap 503 backend memang membawa `galat` - dan itulah yang
// dikunci di sini. Pasangan uji sisi klien:
// `frontend/src/lib/keadaanGalat.test.ts`.

// badanGalat503 memeriksa satu jawaban: 503, JSON, SATU kunci `galat` tak kosong.
func badanGalat503(t *testing.T, nama string, w *httptest.ResponseRecorder) string {
	t.Helper()
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("%s: kode %d, mau 503", nama, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("%s: Content-Type %q - klien membaca badan bukan-JSON sebagai backend mati", nama, ct)
	}
	var isi map[string]any
	if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
		t.Fatalf("%s: badan bukan JSON: %v", nama, err)
	}
	pesan, _ := isi["galat"].(string)
	if len(isi) != 1 || pesan == "" {
		t.Fatalf("%s: badan %v, mau tepat {\"galat\": \"<kalimat>\"}", nama, isi)
	}
	return pesan
}

func TestGalat503MasterKursMembawaKalimatnya(t *testing.T) {
	err := fmt.Errorf("%w: UJI sebab", services.ErrMasterKursRusak)
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, err)

	if pesan := badanGalat503(t, "master kurs rusak", w); pesan != err.Error() {
		t.Errorf("galat = %q, mau kalimat services apa adanya %q", pesan, err.Error())
	}
}

func TestGalat503RuteKursTanpaDatabaseMembawaKalimatnya(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, jalur := range []string{
		"/api/treaty-contract-out/tahun/1000001/kurs",
		"/api/treaty-contract-out/tahun/1000001/kurs/konversi?dari=Rp&nilai=1&skala=8",
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(http.MethodGet, jalur, nil)
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if pesan := badanGalat503(t, jalur, w); pesan != "database belum dikonfigurasi" {
			t.Errorf("%s: galat = %q", jalur, pesan)
		}
	}
}

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
	berkas, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
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
			// dari paket bersama `inti/galat`.
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
