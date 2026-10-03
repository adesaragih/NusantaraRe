package premium

// Total premi per item properti dan per mata uang lokasi FIRE - tiket 07 (akumulasi,
// lintas lini), `SumTotalTSIPremiGross_Act` (`D:\migrasi\RNM\NB FacIn\Activity\`,
// ASM-FW-GISFW-DATA-COVERAGE!SUMTOTALTSIPREMIGROSS_ACT) langkah 1 (`IsFire`).
// Langkah 2-6 (GOLF, ANEKA, MARINE, PA, MBU) belum diport: tanpa nilai pembanding di
// fixture.

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// ErrSeriDivide - `@divide(a,b,20)` jatuh tepat di tengah. `@divide` bukan
// `@Math.divide`: A37 (setengah-ke-atas) hanya terbukti untuk yang kedua, jadi mode
// fungsi ini `belum terverifikasi` dan kasus seri ditolak, bukan ditebak.
var ErrSeriDivide = errors.New("premium: @divide jatuh tepat di tengah, mode pembulatannya belum terverifikasi")

// ErrMataUangKosong - item properti tanpa `.Currency`: entri daftar per mata uang
// dicocokkan lewat nama mata uang (1.3.4.1 `Local.Currency==.Name`).
var ErrMataUangKosong = errors.New("premium: item properti tanpa mata uang")

// ItemProperti - satu `.Property.PropertyItemList(n)`: `.Currency`, `.TSIObjectItem`,
// dan `.Premium` tiap `.CoverageList(m)` (teks desimal bertitik).
type ItemProperti struct {
	MataUang      string
	TSIObjectItem string
	PremiCoverage []string
}

// TotalMataUang - satu `.Property.TotalTSIPremiGrossList(n)`.
//
// LewatDouble = entri ini menerima penjumlahan cabang 1.3.4 (≥2 item bermata uang
// sama). [terverifikasi] Di sistem lama item ke-2 dst. DIBULATKAN KE `double` lalu
// diakumulasi `Decimal`: `Local.TSI` dan `Local.Premi` dideklarasikan `double`
// (SumTotalTSIPremiGross_Act `pyLocalParameters` L293-301; `Local.PremiCov` Decimal,
// L305-307). [terverifikasi] Model itu cocok 2/2 lokasi (`edm-fire-1` akar + OldData;
// eksak 0/2, double murni 0/2) - batas bukti: satu kasus. Port ini TETAP EKSAK -
// keputusan work owner butir 63 (A48): perubahan perilaku yang disengaja; paritas
// double hanya lewat ADR bila kelak diwajibkan (CLAUDE.md §7). Nilai entri LewatDouble
// karena itu belum tentu sama dengan sistem lama, dan rekonsiliasi melaporkannya belum
// tercakup.
type TotalMataUang struct {
	TSI, Premium uang.Money
	Rate         *apd.Decimal
	LewatDouble  bool
}

// TotalLokasi - hasil langkah 1 untuk satu lokasi: `.TotalGrossPremi` tiap item
// (urut item) dan daftar per mata uang (urut kemunculan pertama).
type TotalLokasi struct {
	PremiItem   []uang.Money
	PerMataUang []TotalMataUang
}

// TotalFireLokasi - langkah 1.2-1.3.5 untuk satu lokasi.
//
// [terverifikasi] 1.2 `Property-Remove .Property.TotalTSIPremiGrossList` → daftar
// dibangun ulang. 1.3.1-1.3.3 `Local.PremiCov = 0`, `+= .Premium` tiap coverage,
// `.TotalGrossPremi = Local.PremiCov` (Decimal - penjumlahan eksak). 1.3.4
// entri ber-`.Name` sama: `.TSI += Local.TSI`, `.Premium += Local.Premi`, `.Rate`
// dihitung ulang; 1.3.5 tanpa entri: entri baru dengan nilai item. `.Rate =
// @if(.TSI=0,0,@divide(.Premium,.TSI,20)*1000)`. Diuji terhadap tiga lokasi FIRE nyata
// (rekonsiliasi `TestBerkasKasusP5`).
func TotalFireLokasi(items []ItemProperti) (TotalLokasi, error) {
	var hasil TotalLokasi
	for i, it := range items {
		if it.MataUang == "" {
			return TotalLokasi{}, fmt.Errorf("%w: item %d", ErrMataUangKosong, i+1)
		}
		premi := uang.Money{Amount: apd.New(0, 0), Currency: it.MataUang}
		for j, p := range it.PremiCoverage {
			d, err := angkaWajib(fmt.Sprintf("item %d coverage %d Premium", i+1, j+1), p)
			if err != nil {
				return TotalLokasi{}, err
			}
			if premi, err = premi.Add(uang.Money{Amount: d, Currency: it.MataUang}); err != nil {
				return TotalLokasi{}, err
			}
		}
		tsiItem, err := angkaWajib(fmt.Sprintf("item %d TSIObjectItem", i+1), it.TSIObjectItem)
		if err != nil {
			return TotalLokasi{}, err
		}
		tsi := uang.Money{Amount: tsiItem, Currency: it.MataUang}
		hasil.PremiItem = append(hasil.PremiItem, premi)
		k := -1
		for n, e := range hasil.PerMataUang {
			if e.TSI.Currency == it.MataUang {
				k = n
				break
			}
		}
		if k < 0 {
			hasil.PerMataUang = append(hasil.PerMataUang, TotalMataUang{TSI: tsi, Premium: premi})
			k = len(hasil.PerMataUang) - 1
		} else {
			e := &hasil.PerMataUang[k]
			e.LewatDouble = true
			if e.TSI, err = tsi.Add(e.TSI); err != nil {
				return TotalLokasi{}, err
			}
			if e.Premium, err = premi.Add(e.Premium); err != nil {
				return TotalLokasi{}, err
			}
		}
		e := &hasil.PerMataUang[k]
		if e.Rate, err = rateTotal(e.Premium.Amount, e.TSI.Amount); err != nil {
			return TotalLokasi{}, err
		}
	}
	return hasil, nil
}

// rateTotal - `@if(.TSI=0,0,@divide(.Premium,.TSI,20)*1000)`.
func rateTotal(premi, tsi *apd.Decimal) (*apd.Decimal, error) {
	if tsi.IsZero() {
		return apd.New(0, 0), nil
	}
	ctx := utils.DecimalContext()
	hasilBagi := new(apd.Decimal)
	kondisi, err := ctx.Quo(hasilBagi, premi, tsi)
	if err != nil {
		return nil, err
	}
	keAtas, err := bulatkan(hasilBagi, 20)
	if err != nil {
		return nil, err
	}
	// Seri hanya mungkin bila hasil bagi eksak (berhingga dalam 38 digit).
	if kondisi&apd.Inexact == 0 {
		genap := new(apd.Decimal)
		c := utils.DecimalContext()
		c.Rounding = apd.RoundHalfEven
		if _, err := c.Quantize(genap, hasilBagi, -20); err != nil {
			return nil, err
		}
		if genap.Cmp(keAtas) != 0 {
			return nil, fmt.Errorf("%w: %s", ErrSeriDivide, utils.FormatDecimal(hasilBagi))
		}
	}
	rate := new(apd.Decimal)
	if err := periksaEksak(ctx.Mul(rate, keAtas, apd.New(1000, 0))); err != nil {
		return nil, err
	}
	return rate, nil
}

// angkaWajib - teks desimal bertitik; kosong ditolak (kosong bukan nol).
func angkaWajib(medan, teks string) (*apd.Decimal, error) {
	if teks == "" {
		return nil, fmt.Errorf("%w: %s kosong", ErrAngkaTakTerbaca, medan)
	}
	d, err := utils.ParseDecimal(teks)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrAngkaTakTerbaca, medan, err)
	}
	return d, nil
}
