package services

// Gerbang masukan bersama jalur simpan produk.
//
// ⛔ Angka diperiksa sebagai DESIMAL presisi arbitrer dan disimpan sebagai
// teksnya (ADR-0003); koma diterima sebagai titik desimal (Pega membersihkan
// koma dengan `@replaceAll(…, ",", ".")`, `SetProductNameInward` 3.1 b1046),
// koma DAN titik sekaligus ditolak. Lebih dari 38 digit bermakna DITOLAK,
// tidak dibulatkan Oracle diam-diam.
// ⛔ Tanggal `YYYY-MM-DD` (API) atau `dd/MM/yyyy` (bentuk Pega).
// ⛔ Setiap penolakan menyebut LABEL medan VERBATIM korpus (PARITAS §3), dan
// SEMUA penolakan dilaporkan sekaligus (tiket 05), dipisah `; `.

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

var (
	// ErrIDDariKlien - produk baru membawa ID (400).
	ErrIDDariKlien = errors.New("services: a new product must not carry an id; the server assigns it")
	// ErrMasukanTidakSah - nilai tidak dapat diterima (422).
	ErrMasukanTidakSah = errors.New("services: value is not valid")
)

// GalatValidasi - satu atau lebih penolakan; pesan layar = semuanya.
type GalatValidasi struct {
	Pesan []string
}

func (g GalatValidasi) Error() string { return "services: " + strings.Join(g.Pesan, "; ") }

// PesanLayar - seluruh penolakan, urutan form.
func (g GalatValidasi) PesanLayar() string { return strings.Join(g.Pesan, "; ") }

// Is membuat errors.Is(err, ErrMasukanTidakSah) benar.
func (GalatValidasi) Is(target error) bool { return target == ErrMasukanTidakSah }

// periksa mengumpulkan penolakan.
type periksa struct{ pesan []string }

func (p *periksa) tolak(format string, a ...any) {
	p.pesan = append(p.pesan, fmt.Sprintf(format, a...))
}

func (p *periksa) galat() error {
	if len(p.pesan) == 0 {
		return nil
	}
	return GalatValidasi{Pesan: p.pesan}
}

// desimal memeriksa satu angka di tempat; kosong tetap kosong.
func (p *periksa) desimal(label string, v *string) *apd.Decimal {
	s := strings.TrimSpace(*v)
	if s == "" {
		*v = ""
		return nil
	}
	if strings.Contains(s, ",") {
		if strings.Contains(s, ".") {
			p.tolak("%s %q contains both a comma and a dot", label, s)
			return nil
		}
		s = strings.ReplaceAll(s, ",", ".")
	}
	// Notasi eksponen bukan isian angka Pega (`pxNumber`).
	d, err := utils.ParseDecimal(s)
	if err != nil || strings.ContainsAny(s, "eE") {
		p.tolak("%s %q is not a number", label, s)
		return nil
	}
	ringkas := new(apd.Decimal).Set(d)
	ringkas.Reduce(ringkas)
	if ringkas.NumDigits() > utils.DecimalPrecision {
		p.tolak("%s %q has more than %d significant digits", label, s, utils.DecimalPrecision)
		return nil
	}
	*v = s
	return d
}

// tanggal memeriksa satu tanggal di tempat dan menyeragamkannya ke `YYYY-MM-DD`.
func (p *periksa) tanggal(label string, v *string) {
	s := strings.TrimSpace(*v)
	if s == "" {
		*v = ""
		return
	}
	for _, bentuk := range []string{utils.TanggalSaja, "02/01/2006"} {
		if t, err := time.Parse(bentuk, s); err == nil {
			*v = t.Format(utils.TanggalSaja)
			return
		}
	}
	p.tolak("%s %q is not a date (YYYY-MM-DD)", label, s)
}

// panjang - kolom datar bertipe VARCHAR2(n) `[data DBA]`.
func (p *periksa) panjang(label, v string, maks int) {
	if n := len(v); n > maks {
		p.tolak("%s is %d bytes long; the column holds at most %d", label, n, maks)
	}
}

// tidakLebihBesar - batas bawah ≤ batas atas (AC 32, keputusan tertulis tiket 03/06, R17).
func (p *periksa) tidakLebihBesar(labelBawah string, bawah *apd.Decimal, labelAtas string, atas *apd.Decimal) {
	if bawah != nil && atas != nil && bawah.Cmp(atas) > 0 {
		p.tolak("%s must not be greater than %s", labelBawah, labelAtas)
	}
}
