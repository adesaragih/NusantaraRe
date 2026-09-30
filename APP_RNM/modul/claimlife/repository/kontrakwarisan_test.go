package repository

// Kontrak hilir PremiumList → Claim Life atas `M_LIFE_PREMIUM_DETAIL` - tiket 08.
// TANPA Oracle.
//
// ⛔ Kenapa uji ini ada: sejak pl2 tabel itu punya DUA pemilik di repo ini -
// PENULISnya (polis_warisan.go, modul PremiumList Life) dan PEMBACAnya
// (pesertapolis.go, modul Claim Life). Pelajaran envelope `galat` dan
// `IsCheck`: masing-masing sisi benar menurut dirinya sendiri, dan hanya
// PERTEMUANNYA yang dapat salah. Maka yang diuji di sini pertemuannya:
// setiap kolom yang dibaca Claim Life harus diisi penulis PremiumList.

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// letakPenulisPremiumList - berkas penulis tabel warisan milik PremiumList.
//
// Refactor bentuk B (30-09-2026): uji kontrak ini dulu membaca variabel
// `kolomPesertaWarisan` langsung - keduanya satu paket. Kini penulisnya
// tinggal di modul PremiumList, dan modul Claim Life tidak boleh
// mengimpornya; maka deklarasinya dibaca dari SUMBERnya sebagai pohon
// sintaks. Yang diuji tetap PERTEMUAN kedua sisi.
const letakPenulisPremiumList = "../../../modul/premiumlistlife/repository/polis_warisan.go"

// kolomPenulis adalah satu baris `kolomPesertaWarisan`: kolom dan sumbernya.
type kolomPenulis struct{ Kolom, Sumber string }

// kolomPenulisPremiumList mengurai `var kolomPesertaWarisan = []struct{...}{...}`
// dari sumber penulis PremiumList.
func kolomPenulisPremiumList(t *testing.T) []kolomPenulis {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, letakPenulisPremiumList, nil, 0)
	if err != nil {
		t.Fatalf("penulis PremiumList tidak terbaca: %v", err)
	}
	var hasil []kolomPenulis
	ast.Inspect(f, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "kolomPesertaWarisan" || len(vs.Values) != 1 {
			return true
		}
		lit, ok := vs.Values[0].(*ast.CompositeLit)
		if !ok {
			t.Fatal("kolomPesertaWarisan bukan literal komposit; pembacanya yang rusak")
		}
		for _, e := range lit.Elts {
			baris, ok := e.(*ast.CompositeLit)
			if !ok || len(baris.Elts) < 2 {
				t.Fatalf("baris kolomPesertaWarisan tak terbaca di %s", fset.Position(e.Pos()))
			}
			var teks [2]string
			for i := 0; i < 2; i++ {
				bl, ok := baris.Elts[i].(*ast.BasicLit)
				if !ok || bl.Kind != token.STRING {
					t.Fatalf("medan ke-%d baris kolomPesertaWarisan bukan teks di %s", i, fset.Position(baris.Pos()))
				}
				teks[i], _ = strconv.Unquote(bl.Value)
			}
			hasil = append(hasil, kolomPenulis{Kolom: teks[0], Sumber: teks[1]})
		}
		return false
	})
	if len(hasil) < 10 {
		t.Fatalf("hanya %d kolom penulis PremiumList terbaca; pembacanya yang rusak", len(hasil))
	}
	return hasil
}

// kolomDitulisWarisan adalah seluruh kolom yang diisi `PesertaWarisan.Ganti`.
func kolomDitulisWarisan(t *testing.T) map[string]bool {
	ada := map[string]bool{"ID": true}
	for _, k := range kolomPenulisPremiumList(t) {
		ada[k.Kolom] = true
	}
	return ada
}

// TestKolomBacaClaimLifeDiisiPenulisPremiumList - kontrak dua sisi pl2.
//
// Dua pembaca Claim Life: `kolomSalin` (salinan ke klaim) dan daftar pilih
// `sqlCariPeserta` (layar Find Insured, padanan `GetPesertaClaim_sql1`).
func TestKolomBacaClaimLifeDiisiPenulisPremiumList(t *testing.T) {
	ditulis := kolomDitulisWarisan(t)

	q, _ := sqlCariPeserta("S.M", "UJI-PL", "", "", 10)
	m := regexp.MustCompile(`(?s)SELECT (.*?)\s+FROM`).FindStringSubmatch(q)
	if m == nil {
		t.Fatal("daftar pilih sqlCariPeserta tidak terbaca; pembacanya yang rusak")
	}
	for nama, kolom := range map[string][]string{
		"kolomSalin":     NamaKolomSalinPeserta(),
		"sqlCariPeserta": namaKolomDaftarPilih(m[1]),
	} {
		if len(kolom) < 5 {
			t.Fatalf("%s: hanya %d kolom terbaca; pembacanya yang rusak", nama, len(kolom))
		}
		// ⚠️ Instrumen diuji atas jawaban yang sudah diketahui: ujung pertama
		// dan terakhir tiap daftar harus terbaca, termasuk yang ber-TO_CHAR.
		mau := map[string][2]string{
			"kolomSalin":     {"ID", "CLAIM_AMOUNT"},
			"sqlCariPeserta": {"PL_NUMBER", "EDMSTATUS"},
		}[nama]
		if kolom[0] != mau[0] || kolom[len(kolom)-1] != mau[1] {
			t.Fatalf("%s terbaca %s…%s, mau %s…%s; pembacanya yang rusak",
				nama, kolom[0], kolom[len(kolom)-1], mau[0], mau[1])
		}
		var hilang []string
		for _, k := range kolom {
			if !ditulis[k] {
				hilang = append(hilang, k)
			}
		}
		sort.Strings(hilang)
		if len(hilang) > 0 {
			t.Errorf("%s (Claim Life) membaca %v dari M_LIFE_PREMIUM_DETAIL, "+
				"tetapi penulis PremiumList tidak mengisinya", nama, hilang)
		}
	}
}

// TestPesertaNBTetapHidupDiJalurBacaKlaim - AC 26 spec Claim Life, tiket 08.
//
// ⛔ Jalur NB tidak mengisi `EDMSTATUS` (`d.EDM_STATUS` kosong di unggahan
// NB), jadi baris warisannya ber-`EDMSTATUS` NULL. Penyaring naif
// `NOT IN ('Delete','Batal')` membuang NULL di Oracle - dan seluruh peserta
// NB hilang dari layar klaim.
func TestPesertaNBTetapHidupDiJalurBacaKlaim(t *testing.T) {
	var sumber string
	for _, k := range kolomPenulisPremiumList(t) {
		if k.Kolom == "EDMSTATUS" {
			sumber = k.Sumber
		}
	}
	if sumber != "d.EDM_STATUS" {
		t.Fatalf("EDMSTATUS bersumber %q, mau d.EDM_STATUS (TempValue.EDMStatus)", sumber)
	}
	// NULL tiba di pembaca sebagai teks kosong.
	if !PesertaHidup("") {
		t.Error("peserta ber-EDMSTATUS kosong (NB) disaring keluar dari jalur baca klaim")
	}
	if !strings.Contains(penyaringHidup, "IS NULL") {
		t.Errorf("penyaring hidup %q tidak meloloskan NULL", penyaringHidup)
	}
}
