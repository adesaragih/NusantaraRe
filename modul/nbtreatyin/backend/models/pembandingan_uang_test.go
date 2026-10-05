package models

// Spec-penyimpanan AC 25 / ID-20 (RALAT P11, 04-10-2026): port rumus TIDAK
// PERNAH membandingkan dua nilai uang satu sama lain - sebab 176 rule terjangkau
// pun tidak (sisir XML di tiket 18 bab P11). Yang ada di XML, dan ditiru di
// sini apa adanya:
//
//	uang lawan nol       `.BalanceDueTo>=0` / `<0` (SetDueTo_act 1-2),
//	                     `.Deduction1!="" || .Deduction1!=0` (CountNetPremi_act 6),
//	                     `.Claim!=0&&.Claim!=""` (CountOGPONP_Act 8),
//	                     `@if(.Limit>0,...)` (InputPolicyTreatyInDetail_preACT 17),
//	                     `.PremiOgp == "" ||.PremiOgp == "0"` (CountOGPONP_Act 1-2)
//	persen lawan 100     `.RiCommOgp>100` dst. (CountResult*), `local.pcttotal>100`
//	                     (SetValidateInstallment_Act 4), angsuran 3.3
//	penanda lawan 1      `ListAgent.pxResults(1).STS_PKP == 1` (SetPPNPPH 4)
//
// Pembandingan rasio dua nilai uang bertoleransi `<=0.01` hanya ada di langkah
// berlabel `//` (CountOverridingCommOgp_Act 1/2/4, CountOverridingCommOnp_Act
// 1/2, CountRiCommOgp_act 4, CountRiCommOnp_act 1) - dinonaktifkan, tidak
// diport; `Local.NETPREMI>200000000.00` hanya di CekLimitTreatyAcc_Act (K2,
// tidak dibangun).

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// tetapanPembanding - satu-satunya ruas kanan `.Cmp(...)` yang sah: literal
// rule (`100`, `1`), bukan nilai halaman.
func tetapanPembanding(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name == "seratus"
	case *ast.CallExpr: // apd.New(<literal>, <literal>)
		sel, ok := v.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "New" || len(v.Args) != 2 {
			return false
		}
		for _, a := range v.Args {
			if _, lit := a.(*ast.BasicLit); !lit {
				return false
			}
		}
		return true
	}
	return false
}

// medanUang - nama properti (ruas terakhir) setiap kolom bergolongan uang di
// katalog delapan tabel.
func medanUang() map[string]bool {
	m := map[string]bool{}
	for _, tb := range SemuaTabel {
		for _, k := range tb.Kolom {
			if k.Golongan == GolUang {
				p := k.Properti
				if i := strings.LastIndex(p, "."); i >= 0 {
					p = p[i+1:]
				}
				m[p] = true
			}
		}
	}
	return m
}

// medanDibaca - nama properti yang dibaca sebuah panggilan `x.teks("…")` /
// `h.Ambil(pt+"…")`: literal teks terakhir argumen pertamanya.
func medanDibaca(e ast.Expr) string {
	c, ok := e.(*ast.CallExpr)
	if !ok || len(c.Args) == 0 {
		return ""
	}
	a := c.Args[0]
	for {
		b, ok := a.(*ast.BinaryExpr)
		if !ok {
			break
		}
		a = b.Y
	}
	lit, ok := a.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	s := strings.Trim(lit.Value, "\"`")
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	return s
}

func TestPortTidakMembandingkanDuaNilaiUang(t *testing.T) { // spec-penyimpanan AC 25, ID-20
	fset := token.NewFileSet()
	isi, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var berkas []*ast.File
	for _, e := range isi {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, e.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		berkas = append(berkas, f)
	}
	uang := medanUang()
	if len(uang) == 0 {
		t.Fatal("nol medan uang di katalog; pembacanya yang rusak")
	}
	cmp := 0
	for _, f := range berkas {
		ast.Inspect(f, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.CallExpr:
				sel, ok := v.Fun.(*ast.SelectorExpr)
				if !ok || sel.Sel.Name != "Cmp" || len(v.Args) != 1 {
					return true
				}
				cmp++
				if !tetapanPembanding(v.Args[0]) {
					t.Errorf("%s: `.Cmp` lawan nilai, bukan tetapan rule - pembandingan dua nilai "+
						"uang wajib bertoleransi/terbulatkan (AC 25), dan XML terjangkau tidak memilikinya",
						fset.Position(v.Pos()))
				}
			case *ast.BinaryExpr:
				if v.Op != token.EQL && v.Op != token.NEQ {
					return true
				}
				kiri, kanan := medanDibaca(v.X), medanDibaca(v.Y)
				if kiri != "" && kanan != "" && (uang[kiri] || uang[kanan]) {
					t.Errorf("%s: teks %s dibandingkan persis dengan teks %s (AC 25)",
						fset.Position(v.Pos()), kiri, kanan)
				}
			}
			return true
		})
	}
	if cmp == 0 {
		t.Fatal("nol `.Cmp` terbaca; pembacanya yang rusak")
	}
	t.Logf("%d pembandingan `.Cmp` diperiksa", cmp)
}

// Pembandingan uang LAWAN NOL tetap eksak, sama dengan XML - toleransi di sini
// mengubah perilaku: `SetDueTo_act` langkah 1 `.BalanceDueTo>=0`, langkah 2
// `.BalanceDueTo<0`. Hitung tangan: -0,00000001 < 0 -> DueTo "0";
// 0,00000001 >= 0 -> "1"; 0 >= 0 -> "1".
func TestTandaUangLawanNolEksakSepertiXML(t *testing.T) { // spec-penyimpanan AC 25, ID-20
	for balance, harap := range map[string]string{
		"-0.00000001": "0",
		"0.00000001":  "1",
		"0":           "1",
	} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.BalanceDueTo", balance)
		if err := SetDueTo(h); err != nil {
			t.Fatal(err)
		}
		if got := h.Ambil("PolicyTreatyIn.DueTo"); got != harap {
			t.Errorf("BalanceDueTo %s: DueTo %q, harap %q", balance, got, harap)
		}
	}
}
