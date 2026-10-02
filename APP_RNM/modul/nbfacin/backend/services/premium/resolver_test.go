package premium

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// harusPanic menjalankan f dan gagal bila f TIDAK panic.
func harusPanic(t *testing.T, nama string, f func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Errorf("%s: tidak panic", nama)
		}
	}()
	f()
}

// TestSatuanRatePetaK018 - Seam 2 (NB-03): peta lini bisnis → satuan rate.
// Harapan disalin dari tabel K-018 `docs/00-KEPUTUSAN-WORK-OWNER.md` ("Satuan rate
// per lini bisnis DIKUNCI"), bukan dari kode.
func TestSatuanRatePetaK018(t *testing.T) {
	for _, u := range []struct {
		lini LiniBisnis
		mau  Satuan
	}{
		{LiniPA, PerMil},
		{LiniFire, PerMil},
		{LiniMBU, Persen},
		{LiniAneka, Persen},
		{LiniBonding, Persen},
		{LiniGolf, Persen},
		{LiniMarineCargo, Persen},
	} {
		if got := SatuanRate(u.lini); got != u.mau {
			t.Errorf("%s: satuan %v, mau %v", u.lini, got, u.mau)
		}
	}
}

// TestPembagiSatuan - sumbangan tiap satuan ke pembagi komposit (K-018: ‰ = 1.000,
// % = 100). Bahwa 100000 pada PA L713 terbaca sebagai per mille dibuktikan lewat
// Calculate: `TestPremiPAContohHitung` kasus "prorata pembagi komposit" (premi
// 1.250, bukan 125.000) dan rekonsiliasi kasus PA metode 1.
func TestPembagiSatuan(t *testing.T) {
	if p := PerMil.Pembagi(); p != 1000 {
		t.Errorf("‰ pembagi %d, mau 1000", p)
	}
	if p := Persen.Pembagi(); p != 100 {
		t.Errorf("%% pembagi %d, mau 100", p)
	}
	harusPanic(t, "satuan nol", func() { Satuan(0).Pembagi() })
}

// TestSatuanProRata - ProRatePercent selalu persen, terlepas dari lini bisnis
// (K-018 tabel sumbangan faktor, butir 27).
func TestSatuanProRata(t *testing.T) {
	if got := SatuanProRata(); got != Persen {
		t.Errorf("ProRatePercent %v, mau %%", got)
	}
}

// TestSatuanRateLiniTakDikenalPanic - lini di luar peta → gagal keras, bukan
// satuan bawaan. Ejaan dibandingkan persis: kode lini di data produksi
// dipetakan tiket 10, bukan ditebak di sini.
func TestSatuanRateLiniTakDikenalPanic(t *testing.T) {
	// "Layering": tidak dipakai di sistem baru (keputusan work owner 01-10-2026, butir 30).
	for _, lini := range []LiniBisnis{"", "TRAVEL", "pa", "Fire", "MARINECARGO", "Layering"} {
		harusPanic(t, string(lini), func() { SatuanRate(lini) })
	}
}

// TestCalculateLiniDikenalTanpaRumus - lini yang ADA di peta tetapi tanpa rumus
// sendiri mengembalikan galat, bukan panic dan bukan angka. Sejak tiket 18 hanya
// BONDING (dihitung lewat ANEKA, A11); FIRE/ANEKA/GOLF/MARINE CARGO: lini_lain_test.go.
func TestCalculateLiniDikenalTanpaRumus(t *testing.T) {
	for _, lini := range []LiniBisnis{LiniBonding} {
		in := kasusDasar()
		in.LiniBisnis = lini
		if _, err := Calculate(in); !errors.Is(err, ErrBentukBelumDiport) {
			t.Errorf("%s: galat %v, mau ErrBentukBelumDiport", lini, err)
		}
	}
}

// bolehMembangunRasio - berkas yang boleh memuat literal `rasio{…}`, dengan
// alasannya. Daftar eksplisit, bukan pola.
var bolehMembangunRasio = map[string]string{
	"resolver.go": "fungsi pembangun rasio (rasioRate, rasioProRata, rasioPeriodePendek, rasioDiskon, rasioPersen, rasio.faktor)",
	"kompilasi_gagal_rasio.go": "berkas wajib-gagal-kompilasi bertag gagalkompilasi, " +
		"tidak pernah ikut build (TestUangTambahRasioGagalKompilasi)",
}

// TestRasioHanyaDibangunDiResolver - penjaga "resolver satu-satunya pengisi
// skala" (NB-03): literal `rasio{…}` hanya boleh ada di resolver.go. Kode lain
// memakai `rasioRate` / `rasioProRata`, sehingga tidak ada yang menyetel satuan
// sendiri.
func TestRasioHanyaDibangunDiResolver(t *testing.T) {
	berkas, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	diperiksa := 0
	for _, b := range berkas {
		n := b.Name()
		if !strings.HasSuffix(n, ".go") || strings.HasSuffix(n, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, n, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		diperiksa++
		ast.Inspect(f, func(node ast.Node) bool {
			if lit, ok := node.(*ast.CompositeLit); ok {
				if id, ok := lit.Type.(*ast.Ident); ok && id.Name == "rasio" && bolehMembangunRasio[n] == "" {
					t.Errorf("%s: literal rasio{…} di luar resolver.go - pakai fungsi pembangun di resolver.go", fset.Position(lit.Pos()))
				}
			}
			return true
		})
	}
	if diperiksa < 2 {
		t.Fatalf("hanya %d berkas diperiksa; pembacanya yang rusak", diperiksa)
	}
}

// TestSatuanRatePanikBertipe - lini di luar peta skala panic dengan premium.PanikLini
// (dipilah pemanggil lewat errors.As, bukan teks), pesannya tetap menyebut lininya.
func TestSatuanRatePanikBertipe(t *testing.T) {
	defer func() {
		r := recover()
		p, ok := r.(PanikLini)
		if !ok || p.Lini != "UJI-LINI" || !strings.Contains(p.Error(), "UJI-LINI") {
			t.Fatalf("panic %#v", r)
		}
	}()
	SatuanRate("UJI-LINI")
}
