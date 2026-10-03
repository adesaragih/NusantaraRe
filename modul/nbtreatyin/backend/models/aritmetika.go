package models

// Untuk apa berkas ini: ARITMETIKA EKSPRESI PEGA atas desimal.
//
// Setiap rumus uang modul ini adalah ekspresi `Property-Set` di aktivitas Pega
// (`docs/INVENTARIS-XML.md` bab 6). Operatornya ditiru di sini satu per satu,
// supaya port tiap langkah dapat dibaca berdampingan dengan aslinya:
//
//	+ - *                  -> Tambah, Kurang, Kali   (eksak)
//	/                      -> Bagi                   (presisi penuh 38 digit)
//	@divide(a,b,n)         -> BagiBulat(a, b, n)     (bagi, lalu bulatkan n desimal)
//	@Math.divide(a,b,n)    -> BagiBulat(a, b, n)     (fungsi yang sama)
//
// ⛔ Nol float di seluruh berkas (ADR-U-0003, spec §5.6 AC 25).
//
// ⚠️ PEMBULATAN `@divide` = SETENGAH-KE-ATAS, mengikuti konvensi yang sudah
// ditetapkan di repo ini untuk fungsi Pega yang sama (modul nbfacin,
// `services/premium/premium.go` `bulatkan`, keputusan A37). Mode pembulatan
// internal Pega tidak ada di ekspor; konvensi repo dipakai supaya dua modul
// tidak membulatkan `@divide` dengan dua cara.
//
// ⚠️ PEMBAGIAN `/` TANPA `@divide`: Pega membaginya dengan BigDecimal; konteks
// presisinya tidak ada di ekspor. Di sini presisi penuh sistem
// (`utils.DecimalPrecision` = 38), setengah-ke-atas. Hasilnya disimpan ke
// kolom NUMBER(38,8) (spec-penyimpanan ID-14), jadi selisih di digit ke-38
// tidak pernah sampai ke penyimpanan.

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// ErrBagiNol - pembagi bernilai nol. Pega melempar galat ekspresi di titik
// yang sama dan langkahnya gagal; di sini galatnya dikembalikan dan
// perhitungan berhenti - tidak pernah diam-diam bernilai nol.
var ErrBagiNol = errors.New("models: pembagian dengan nol")

// kalkulator membawa galat pertama sepanjang satu rumus, supaya port rumus
// tetap satu baris seperti aslinya: hasil antara yang gagal menjadikan seluruh
// rumus gagal, dan galatnya dibaca SEKALI di ujung.
type kalkulator struct {
	err error
}

func (k *kalkulator) ctx() *apd.Context {
	c := utils.DecimalContext()
	c.Rounding = apd.RoundHalfUp
	return c
}

func (k *kalkulator) catat(err error) {
	if k.err == nil && err != nil {
		k.err = err
	}
}

// d mengurai literal desimal yang tertulis di rule ("100", "102.2").
func (k *kalkulator) d(lit string) *apd.Decimal {
	v, err := utils.ParseDecimal(lit)
	if err != nil {
		k.catat(fmt.Errorf("models: literal %q: %w", lit, err))
		return apd.New(0, 0)
	}
	return v
}

func (k *kalkulator) nol(x *apd.Decimal) bool { return x == nil || x.IsZero() }

// Tambah = `a + b`.
func (k *kalkulator) Tambah(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Add(r, a, b)
	k.catat(err)
	return r
}

// Kurang = `a - b`.
func (k *kalkulator) Kurang(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Sub(r, a, b)
	k.catat(err)
	return r
}

// Kali = `a * b`.
func (k *kalkulator) Kali(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Mul(r, a, b)
	k.catat(err)
	return r
}

// Bagi = `a / b` tanpa `@divide` - presisi penuh.
func (k *kalkulator) Bagi(a, b *apd.Decimal) *apd.Decimal {
	if k.err != nil || a == nil || b == nil {
		return apd.New(0, 0)
	}
	if k.nol(b) {
		k.catat(ErrBagiNol)
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Quo(r, a, b)
	k.catat(err)
	return r
}

// BagiBulat = `@divide(a, b, desimal)` / `@Math.divide(a, b, desimal)`.
func (k *kalkulator) BagiBulat(a, b *apd.Decimal, desimal int32) *apd.Decimal {
	q := k.Bagi(a, b)
	if k.err != nil {
		return apd.New(0, 0)
	}
	return k.Bulatkan(q, desimal)
}

// Bulatkan membulatkan ke `desimal` angka di belakang koma, setengah-ke-atas.
func (k *kalkulator) Bulatkan(x *apd.Decimal, desimal int32) *apd.Decimal {
	if k.err != nil || x == nil {
		return apd.New(0, 0)
	}
	r := new(apd.Decimal)
	_, err := k.ctx().Quantize(r, x, -desimal)
	k.catat(err)
	return r
}
