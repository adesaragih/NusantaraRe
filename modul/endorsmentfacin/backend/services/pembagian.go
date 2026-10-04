package services

// Aritmetika desimal eksak bersama porsi periode (porsiperiode.go,
// porsipembayaran.go) dan premi FIRE (premifire.go).
//
// ⚠️ Fungsi sejenis ada di modul lain (nbfacin, treatycontractout); kandidat
// `inti/backend/utils`, tetapi inti di luar jatah modul ini - dicatat di
// laporan, tidak dipindah.

import (
	"errors"
	"math/big"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

var (
	// ErrPresisiPremi - hasil kali/jumlah/bagi melampaui presisi desimal
	// aplikasi (38 digit) sehingga tidak eksak.
	ErrPresisiPremi = errors.New("endorsement: hasil aritmetika tidak eksak pada presisi 38 digit")
	// ErrBagiNol - penyebut nol.
	ErrBagiNol = errors.New("endorsement: pembagian dengan nol")
)

// kaliEksak, tambahEksak - aritmetika presisi aplikasi; hasil yang tidak eksak
// ditolak, bukan dibulatkan diam-diam.
func kaliEksak(a, b *apd.Decimal) (*apd.Decimal, error) {
	d := new(apd.Decimal)
	cond, err := utils.DecimalContext().Mul(d, a, b)
	if err != nil {
		return nil, err
	}
	if cond.Inexact() {
		return nil, ErrPresisiPremi
	}
	return d, nil
}

func tambahEksak(a, b *apd.Decimal) (*apd.Decimal, error) {
	d := new(apd.Decimal)
	cond, err := utils.DecimalContext().Add(d, a, b)
	if err != nil {
		return nil, err
	}
	if cond.Inexact() {
		return nil, ErrPresisiPremi
	}
	return d, nil
}

// bagiHalfUp - `@Math.divide(a, b, n)`: HALF_UP pada n desimal (A37,
// dikonfirmasi work owner), dihitung eksak lewat bilangan bulat lalu dijaga
// tidak melampaui presisi aplikasi (utils.DecimalPrecision, 38 digit).
//
// ⚠️ `[dugaan]` Untuk hasil NEGATIF, setengah dibulatkan menjauhi nol seperti
// `RoundingMode.HALF_UP` Java; yang terbukti di kasus nyata hanya nilai
// positif.
func bagiHalfUp(a, b *apd.Decimal, n int32) (*apd.Decimal, error) {
	if b.IsZero() {
		return nil, ErrBagiNol
	}
	num := a.Coeff.MathBigInt()
	den := b.Coeff.MathBigInt()
	k := int64(a.Exponent) - int64(b.Exponent) + int64(n)
	pangkat := func(e int64) *big.Int { return new(big.Int).Exp(big.NewInt(10), big.NewInt(e), nil) }
	if k >= 0 {
		num.Mul(num, pangkat(k))
	} else {
		den.Mul(den, pangkat(-k))
	}
	hasil, sisa := new(big.Int).QuoRem(num, den, new(big.Int))
	if new(big.Int).Lsh(sisa, 1).Cmp(den) >= 0 {
		hasil.Add(hasil, big.NewInt(1))
	}
	if a.Negative != b.Negative && hasil.Sign() != 0 {
		hasil.Neg(hasil)
	}
	d := apd.NewWithBigInt(new(apd.BigInt).SetMathBigInt(hasil), -n)
	if d.NumDigits() > utils.DecimalPrecision {
		return nil, ErrPresisiPremi
	}
	return d, nil
}
