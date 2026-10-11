package models

// Untuk apa berkas ini: ARITMETIKA UANG Claim Prop atas desimal - satu jalur perhitungan (spec §4, AC 18-24).
//
// Setiap rumus uang modul ini adalah ekspresi `Property-Set` activity Pega. Operatornya ditiru di sini, supaya port
// tiap langkah dapat dibaca berdampingan dengan aslinya:
//
//	+ - *                  -> Tambah, Kurang, Kali   (eksak)
//	/                      -> Bagi                   (presisi penuh 38 digit)
//	@divide(a,b,n)         -> Bagi                   (⚠️ TANPA pembulatan n desimal - penyimpangan sadar spec §4)
//	x * p / 100            -> Persen(x, p)           (kali dulu, bagi terakhir - AC 21)
//
// ⚠️ PENYIMPANGAN SADAR spec §4 butir 1-3 (keputusan work owner): `.ClaimSpreaded` hari ini ditulis di 20 titik
// dengan 4 perlakuan pembulatan (`/100`, `@divide(…,100,4)`, `,100,10`, `,100,20`). Di sini SATU perlakuan:
// presisi penuh, pembulatan hanya di batas penyimpanan NUMBER(38,10) dan di tampilan.
//
// ⛔ Nol float di seluruh berkas (ADR-0003, AC 18).

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// ErrBagiNol - pembagi bernilai nol. Pega melempar galat ekspresi di titik yang sama; di sini galatnya dikembalikan
// dan perhitungan berhenti - tidak pernah diam-diam bernilai nol.
var ErrBagiNol = errors.New("models: pembagian dengan nol")

// Kalkulator membawa galat pertama sepanjang satu rumus, supaya port rumus tetap satu baris seperti aslinya: hasil
// antara yang gagal menjadikan seluruh rumus gagal, dan galatnya dibaca SEKALI di ujung (`Galat`).
type Kalkulator struct {
	err error
}

func (k *Kalkulator) ctx() *apd.Context {
	c := utils.DecimalContext()
	c.Rounding = apd.RoundHalfUp
	return c
}

func (k *Kalkulator) catat(err error) {
	if k.err == nil && err != nil {
		k.err = err
	}
}

// Galat mengembalikan galat pertama rumus.
func (k *Kalkulator) Galat() error { return k.err }

// D mengurai literal desimal yang tertulis di rule ("100").
func (k *Kalkulator) D(lit string) *apd.Decimal {
	v, err := utils.ParseDecimal(lit)
	if err != nil {
		k.catat(fmt.Errorf("models: literal %q: %w", lit, err))
		return apd.New(0, 0)
	}
	return v
}

// Teks membaca nilai properti (kosong = 0, semantik `@toDecimal`).
func (k *Kalkulator) Teks(jalur, teks string) *apd.Decimal {
	d, err := AngkaTeks(jalur, teks)
	k.catat(err)
	if d == nil {
		return apd.New(0, 0)
	}
	return d
}

// H membaca satu jalur halaman.
func (k *Kalkulator) H(h *Halaman, jalur string) *apd.Decimal { return k.Teks(jalur, h.Ambil(jalur)) }

// B membaca satu properti baris.
func (k *Kalkulator) B(b Baris, prop string) *apd.Decimal { return k.Teks("."+prop, b[prop]) }

// Tambah = `a + b`.
func (k *Kalkulator) Tambah(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Add(r, a, b)
	k.catat(err)
	return r
}

// Kurang = `a - b`.
func (k *Kalkulator) Kurang(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Sub(r, a, b)
	k.catat(err)
	return r
}

// Kali = `a * b * ...`.
func (k *Kalkulator) Kali(a *apd.Decimal, lain ...*apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal).Set(a)
	for _, b := range lain {
		if b == nil {
			return apd.New(0, 0)
		}
		_, err := k.ctx().Mul(r, r, b)
		k.catat(err)
	}
	return r
}

// Bagi = `a / b` - presisi penuh sistem (38 digit), setengah-ke-atas di digit terakhir.
func (k *Kalkulator) Bagi(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	if b.IsZero() {
		k.catat(ErrBagiNol)
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Quo(r, a, b)
	k.catat(err)
	r.Reduce(r)
	return r
}

// Persen = `x * p1 * p2 ... / 100^n` - KALI DULU, BAGI TERAKHIR (AC 21): satu pembagian di ujung, sehingga hasilnya
// tidak bergantung pada urutan perkalian persen.
func (k *Kalkulator) Persen(x *apd.Decimal, persen ...*apd.Decimal) *apd.Decimal {
	hasil := k.Kali(x, persen...)
	pembagi := apd.New(1, 0)
	seratus := apd.New(100, 0)
	for range persen {
		pembagi = k.Kali(pembagi, seratus)
	}
	return k.Bagi(hasil, pembagi)
}

// Neg = `x * -1`.
func (k *Kalkulator) Neg(x *apd.Decimal) *apd.Decimal {
	if k.err != nil || x == nil {
		return apd.New(0, 0)
	}
	return new(apd.Decimal).Neg(x)
}

// Banding membandingkan dua desimal: -1, 0, 1.
func Banding(a, b *apd.Decimal) int {
	if a == nil {
		a = apd.New(0, 0)
	}
	if b == nil {
		b = apd.New(0, 0)
	}
	return a.Cmp(b)
}

// Lebih = `a > b`.
func Lebih(a, b *apd.Decimal) bool { return Banding(a, b) > 0 }

// Nol = `x == 0`.
func Nol(x *apd.Decimal) bool { return x == nil || x.IsZero() }

// Teks menulis desimal sebagai teks halaman (`utils.FormatDecimal`), angka nol menjadi "0".
func Teks(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	r := new(apd.Decimal)
	r.Reduce(d)
	if r.IsZero() {
		return "0"
	}
	return utils.FormatDecimal(r)
}
