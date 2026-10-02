package medis

// Penafsir ekspresi Pega - TANPA ambang korpus. Seluruh ekspresi di berkas ini
// SINTETIS; ambang klinis tidak disalin ke test (E20).
//
// Dibaca sesudah: ekspresi.go.

import (
	"errors"
	"testing"
)

type halamanUji map[string]string

func (h halamanUji) ambil(nama string) (string, error) { return h[nama], nil }

func nilaiUji(t *testing.T, ekspresi string, h halamanUji) string {
	t.Helper()
	got, err := evaluasiTeks(ekspresi, h)
	if err != nil {
		t.Fatalf("%s: %v", ekspresi, err)
	}
	return got
}

func TestEkspresiDasar(t *testing.T) {
	h := halamanUji{".A": "7", ".B": "teks", "param.P": "Semua", ".K": "", ".D": "1.50"}
	for _, u := range []struct{ ekspresi, mau string }{
		{`.A==7`, "true"},
		{`.A = 7`, "true"},
		{`.A>=7&&.A<8`, "true"},
		{`.A>7||.A<=7`, "true"},
		{`.B=="teks"`, "true"},
		{`.B!="teks"`, "false"},
		{`param.P=="Semua"`, "true"},
		{`.K!=""`, "false"},
		{`(.K!=""||.K!="0")`, "true"},
		{`.D==1.50`, "true"},
		{`@if(.A<=7,"X","Y")`, "X"},
		{`@if(.A>7,"X",@if(.A==7,"Z","Y"))`, "Z"},
		{`@toDecimal(@replaceAll("3,25",",","."))`, "3.25"},
		{`@divide(1,4,2)`, "0.25"},
		{`@divide(6,3,2)`, "2.00"},
	} {
		if got := nilaiUji(t, u.ekspresi, h); got != u.mau {
			t.Errorf("%s = %q, mau %q", u.ekspresi, got, u.mau)
		}
	}
}

// TestTafsirBerbedaDitolak - teks yang terbaca angka lawan angka: bila tafsir
// teks dan angka berbeda, ditolak (sikap registry NB, butir 20).
func TestTafsirBerbedaDitolak(t *testing.T) {
	h := halamanUji{".D": "1.50", ".N": "100"}
	for _, e := range []string{`.D==1.5`, `.N>=80`} {
		if _, err := evaluasiTeks(e, h); !errors.Is(err, ErrTafsirBerbeda) {
			t.Errorf("%s: galat %v, mau ErrTafsirBerbeda", e, err)
		}
	}
	// Sesudah dinormalkan @toDecimal, nilai bertipe angka: tidak ambigu.
	if got := nilaiUji(t, `@toDecimal(.N)>=80`, h); got != "true" {
		t.Errorf("%q", got)
	}
}

// TestIfMalas - cabang yang tidak dipilih tidak dievaluasi `[dugaan]`.
func TestIfMalas(t *testing.T) {
	if got := nilaiUji(t, `@if(1==1,"X",@toDecimal(.Kosong))`, halamanUji{}); got != "X" {
		t.Errorf("%q", got)
	}
}

// TestKosongTidakDitebakNol - perilaku Pega atas nilai kosong di konteks angka
// belum terverifikasi: galat, bukan nol.
func TestKosongTidakDitebakNol(t *testing.T) {
	for _, e := range []string{`@toDecimal(.K)`, `.K>=1`, `.K<1`} {
		if _, err := evaluasiTeks(e, halamanUji{".K": ""}); !errors.Is(err, ErrNilaiBukanAngka) {
			t.Errorf("%s: galat %v, mau ErrNilaiBukanAngka", e, err)
		}
	}
}

func TestBagiNolDitolak(t *testing.T) {
	if _, err := evaluasiTeks(`@divide(1,0,2)`, halamanUji{}); !errors.Is(err, ErrBagiNol) {
		t.Errorf("galat %v, mau ErrBagiNol", err)
	}
}

// TestPembulatanYangMengubahHasilDitolak - mode pembulatan @divide belum
// terverifikasi. Pembulatan yang tidak mengubah hasil akhir diterima; yang
// mengubahnya ditolak.
func TestPembulatanYangMengubahHasilDitolak(t *testing.T) {
	// 1/3 = 0,333…: semua mode memberi 0,33 atau 0,34, keduanya < 0,5.
	if got := nilaiUji(t, `@if(@divide(1,3,2)<0.5,"KECIL","BESAR")`, halamanUji{}); got != "KECIL" {
		t.Errorf("%q", got)
	}
	// 2/3 = 0,666…: half-up memberi 0,67, potong memberi 0,66 - menentukan hasil.
	_, err := evaluasiTeks(`@if(@divide(2,3,2)>=0.67,"A","B")`, halamanUji{})
	if !errors.Is(err, ErrBergantungModePembulatan) {
		t.Errorf("galat %v, mau ErrBergantungModePembulatan", err)
	}
	// Nilai hasil bagi yang DISIMPAN berbeda antar mode: ditolak.
	if _, err := evaluasiTeks(`@divide(2,3,2)`, halamanUji{}); !errors.Is(err, ErrBergantungModePembulatan) {
		t.Errorf("galat %v, mau ErrBergantungModePembulatan", err)
	}
}

func TestSintaksRusakDitolak(t *testing.T) {
	for _, e := range []string{`@if(1==1,"X"`, `.A==`, `@tidakDikenal(1)`, `"tak tertutup`} {
		if _, err := evaluasiTeks(e, halamanUji{".A": "1"}); err == nil {
			t.Errorf("%s: tanpa galat", e)
		}
	}
}
