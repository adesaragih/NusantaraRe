package repository

// Test penguraian desimal yang datang dari Oracle - TANPA Oracle.
//
// Untuk apa berkas ini: uraiDesimal adalah tempat keputusan "galat dilaporkan,
// bukan ditelan" benar-benar diambil. Ronde 1 menelannya lewat
// `if d, err := ...; err == nil`, sehingga nilai yang tidak terbaca pulang
// sebagai kosong dan tidak dapat dibedakan dari nol. Perilaku itu dikunci di
// sini, bukan hanya dibicarakan di komentar.
//
// Dibaca sesudah: pohonklaim.go.

import (
	"database/sql"
	"strings"
	"testing"
)

func teks(s string) sql.NullString { return sql.NullString{String: s, Valid: true} }
func nullTeks() sql.NullString     { return sql.NullString{} }

// ⛔ Nilai yang tidak terurai menghasilkan GALAT, tidak diam-diam menjadi kosong.
func TestUraiDesimalMengembalikanGalat(t *testing.T) {
	for _, rusak := range []string{"bukan angka", "1.2.3", "12,34abc", "--5"} {
		_, err := uraiDesimal("UJI-R1", "CLAIM_AMOUNT", teks(rusak))
		if err == nil {
			t.Errorf("uraiDesimal(%q) tidak menghasilkan galat", rusak)
			continue
		}
		// Galatnya harus menyebut baris DAN kolomnya, supaya dapat ditelusuri.
		for _, wajib := range []string{"UJI-R1", "CLAIM_AMOUNT", rusak} {
			if !strings.Contains(err.Error(), wajib) {
				t.Errorf("galat %q tidak menyebut %q", err, wajib)
			}
		}
	}
}

// NULL dan teks kosong BUKAN galat - keduanya berarti "tidak ada nilai".
func TestUraiDesimalKosongBukanGalat(t *testing.T) {
	for nama, v := range map[string]sql.NullString{
		"NULL":        nullTeks(),
		"teks kosong": teks(""),
		"spasi":       teks("   "),
	} {
		d, err := uraiDesimal("UJI-R1", "CLAIM_AMOUNT", v)
		if err != nil {
			t.Errorf("%s menghasilkan galat: %v", nama, err)
		}
		if d != nil {
			t.Errorf("%s menghasilkan %v, mau nil", nama, d)
		}
	}
}

// Uang pulang tanpa berubah satu digit, termasuk yang panjang dan yang kecil.
func TestUraiUangTanpaBerubahSatuDigit(t *testing.T) {
	for _, angka := range []string{"1234567890.12345678", "0.00000001", "0", "-42.5"} {
		m, err := uraiUang("UJI-R1", "CLAIM_AMOUNT", teks(angka), "IDR")
		if err != nil {
			t.Fatalf("uraiUang(%q): %v", angka, err)
		}
		if m.Amount == nil {
			t.Fatalf("uraiUang(%q) menghasilkan nil", angka)
		}
		if got := m.Amount.Text('f'); got != angka {
			t.Errorf("uraiUang(%q) = %q - digitnya berubah", angka, got)
		}
		if m.Currency != "IDR" {
			t.Errorf("mata uang = %q, mau IDR", m.Currency)
		}
	}
}

// Galat pada rasio juga dikembalikan, dan mata uang TIDAK ikut dibawa -
// rasio bukan uang (ADR-F-0004).
func TestUraiRasioMengembalikanGalatDanBukanUang(t *testing.T) {
	if _, err := uraiRasio("UJI-R1", "PERCENT_SHARE", teks("bukan angka")); err == nil {
		t.Error("rasio rusak tidak menghasilkan galat")
	}
	r, err := uraiRasio("UJI-R1", "PERCENT_SHARE", teks("0.3"))
	if err != nil {
		t.Fatal(err)
	}
	if r.Value == nil || r.Value.Text('f') != "0.3" {
		t.Errorf("rasio = %v, mau 0.3", r.Value)
	}
	if r.Kosong() {
		t.Error("rasio terisi dilaporkan kosong")
	}
}

// Galat membawa nama kolom yang BERBEDA untuk kolom yang berbeda, supaya
// laporan menunjuk kolom yang benar.
func TestGalatMenyebutKolomYangBenar(t *testing.T) {
	_, errA := uraiUang("UJI-R1", "PREMIUM_SPREADED_NET", teks("x"), "")
	_, errB := uraiRasio("UJI-R1", "RETROCADED_SHARE", teks("x"))
	if errA == nil || errB == nil {
		t.Fatal("salah satu tidak menghasilkan galat")
	}
	if !strings.Contains(errA.Error(), "PREMIUM_SPREADED_NET") {
		t.Errorf("galat A tidak menyebut kolomnya: %v", errA)
	}
	if !strings.Contains(errB.Error(), "RETROCADED_SHARE") {
		t.Errorf("galat B tidak menyebut kolomnya: %v", errB)
	}
}
