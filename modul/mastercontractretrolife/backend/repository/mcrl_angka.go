package repository

// Menulis desimal ke kolom NUMBER tanpa bergantung NLS sesi.
//
// ⛔ Jalur baca sudah kebal NLS (`db.FmtDesimal`). Jalur tulis di modul lain
// membind teks "apa adanya" dan bergantung pada `NLS_NUMERIC_CHARACTERS` sesi
// - sesi ber-NLS Indonesia membaca "1.5" sebagai galat. Di sini desimal
// dipecah menjadi koefisien BULAT (teks) dan skala, lalu dirakit Oracle:
// `TO_NUMBER(:koef) / POWER(10, :skala)`. Bilangan bulat tidak punya pemisah
// desimal, dan pembagian dengan pangkat sepuluh pada NUMBER eksak - nol float
// di jalur mana pun (ADR-U-0003).

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"
)

// angka menyusun ekspresi SQL untuk dua penampung: koefisien dan skala.
func angka(koef, skala string) string {
	return fmt.Sprintf(`(TO_NUMBER(%s) / POWER(10, %s))`, koef, skala)
}

// PecahAngka memecah desimal menjadi koefisien bulat (teks) dan skala.
// nil = NULL (skala 0). Contoh: 1500000000.10 -> ("150000000010", 2).
func PecahAngka(d *apd.Decimal) (any, int64) {
	if d == nil {
		return nil, 0
	}
	s := d.Text('f')
	tanda := ""
	if strings.HasPrefix(s, "-") {
		tanda, s = "-", s[1:]
	}
	bulat, pecahan, _ := strings.Cut(s, ".")
	koef := strings.TrimLeft(bulat+pecahan, "0")
	if koef == "" {
		return "0", 0
	}
	return tanda + koef, int64(len(pecahan))
}

// argAngka menambahkan dua argumen bind untuk satu desimal.
func argAngka(args []any, d *apd.Decimal) []any {
	koef, skala := PecahAngka(d)
	return append(args, koef, skala)
}
