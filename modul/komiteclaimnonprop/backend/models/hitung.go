package models

// Untuk apa berkas ini: ARITMETIKA UANG langkah tingkat akhir (CNPLayerList, Kasir) atas desimal - bagian yang dipakai
// dari `modul/claimnonprop/backend/models/hitung.go` (disalin, bukan diimpor); aturannya sama dengan Claim Non Prop.
//
// Setiap rumus uang modul ini adalah ekspresi `Property-Set` activity Pega. Operatornya ditiru di sini, supaya port
// tiap langkah dapat dibaca berdampingan dengan aslinya:
//
//	+ - *                  -> Tambah, Kurang, Kali   (eksak)
//	/                      -> Bagi                   (presisi penuh 38 digit)
//	@divide(a,b,n)         -> BagiPega               (presisi penuh; pembagi nol = 0, lihat di bawah)
//
// ⚠️ PENYIMPANGAN SADAR yang diwarisi pola Claim Prop: rumus Pega memakai beberapa skala `@divide` (10, 20) dan `/`
// tanpa skala; di sini SATU perlakuan - presisi penuh, pembulatan hanya di batas penyimpanan NUMBER(38,10) dan di
// tampilan (PARITAS).
//
// `[inferensi]` Pembagi nol pada `@divide` / `/` bernilai 0, bukan galat: `@if` Pega mengevaluasi KEDUA cabangnya
// (fungsi biasa), dan rumus XoL menulis `@if(.TreatyType=="UR", ..., @divide(x*100, .ClaimPercentage, 20))` untuk baris
// UR yang ClaimPercentage-nya 0 (AdjClaimAmount_Act langkah 2, CountLossAllocation_act langkah 17.2.5). Rumus itu
// berjalan di produksi tanpa galat - pembagian nolnya tidak melempar (`BagiPega`).
//
// ⛔ Nol float di seluruh berkas (ADR-0003).

import (
	"errors"

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

// Teks membaca nilai properti (kosong = 0, semantik `@toDecimal`).
func (k *Kalkulator) Teks(jalur, teks string) *apd.Decimal {
	d, err := desimal(jalur, teks)
	k.catat(err)
	if d == nil {
		return apd.New(0, 0)
	}
	return d
}

// B membaca satu properti baris.
func (k *Kalkulator) B(b map[string]string, prop string) *apd.Decimal {
	return k.Teks("."+prop, b[prop])
}

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

// BagiPega = `@divide(a, b, n)` / `a / b` rumus XoL: seperti Bagi, tetapi pembagi nol bernilai 0 (lihat kepala
// berkas, `[inferensi]`).
func (k *Kalkulator) BagiPega(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil || b.IsZero() {
		return apd.New(0, 0)
	}
	return k.Bagi(a, b)
}

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
