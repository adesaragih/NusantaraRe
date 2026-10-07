package repository

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// desimalUji - teks desimal uji -> *apd.Decimal (panik bila salah tulis di uji).
func desimalUji(s string) *apd.Decimal {
	d, err := utils.ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

// TestDesimalPulangPergi - ADR-0034 (butir 94): ikatDesimal -> (Oracle TO_NUMBER/TO_CHAR TM9) -> bacaDesimal mengembalikan
// nilai SAMA PERSIS. TM9 membuang nol depan (".5") dan nol ekor; disimulasikan. Kosong = NULL = nil.
func TestDesimalPulangPergi(t *testing.T) {
	for _, u := range []struct{ masuk, tm9 string }{
		{"0.5", ".5"}, {"0.12345678", ".12345678"}, {"123456789012345678901234567890.12345678", "123456789012345678901234567890.12345678"},
		{"123456789012345678901234567890", "123456789012345678901234567890"}, {"1500000000.10", "1500000000.1"}, {"100", "100"},
	} {
		d := desimalUji(u.masuk)
		ikat := ikatDesimal(d)
		if ikat != u.masuk {
			t.Errorf("ikat %q = %v", u.masuk, ikat)
		}
		balik, err := bacaDesimal("UJI", "UJI_KOLOM", &sql.NullString{String: u.tm9, Valid: true})
		if err != nil || balik.Cmp(d) != 0 {
			t.Errorf("%q: baca %q = %v (%v), mau sama nilai", u.masuk, u.tm9, balik, err)
		}
	}
	if ikatDesimal(nil) != nil {
		t.Error("nil harus diikat NULL")
	}
	if d, err := bacaDesimal("UJI", "UJI_KOLOM", &sql.NullString{}); d != nil || err != nil {
		t.Errorf("NULL -> %v %v", d, err)
	}
	for _, u := range []struct{ ekspresi, harus string }{
		{angkaMasuk(":5"), "TO_NUMBER(:5, 'FM999999999999999999999999999999D99999999', 'NLS_NUMERIC_CHARACTERS="},
		{angkaKeluar("i.TSI_OBJECT_ITEM"), "TO_CHAR(i.TSI_OBJECT_ITEM, 'TM9', 'NLS_NUMERIC_CHARACTERS="},
	} {
		if !strings.HasPrefix(u.ekspresi, u.harus) {
			t.Errorf("konversi NLS satu tempat berubah: %s", u.ekspresi)
		}
	}
}
