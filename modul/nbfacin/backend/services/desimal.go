package services

// Uang dan persen masukan layar (tiket 39/42) - ADR-0034 (butir 94): di dalam aplikasi berbentuk desimal berskala
// tetap (*apd.Decimal), teks HANYA di batas JSON. Perubahan bentuk teks -> desimal adalah SATU fungsi (UraiDesimalIsian),
// dipakai handler; pemeriksaan ulang di services (desimalSah) memakai aturan yang sama (polaDesimal).

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
)

// polaDesimal - uang/persen NUMBER(38,8): >= 0, titik pemisah, <= 30 digit bulat, <= 8 desimal (ADR-0016).
var polaDesimal = regexp.MustCompile(`^[0-9]{1,30}(\.[0-9]{1,8})?$`)

// pesanDesimal - pesan 400 medan desimal tidak sah (sama dengan sebelum ADR-0034).
const pesanDesimal = " harus angka >= 0 dengan paling banyak 8 desimal"

// UraiDesimalIsian - teks desimal bertitik dari JSON -> *apd.Decimal. Kosong -> nil. Tidak sah -> nil dan pesan
// `<jalur> harus angka >= 0 dengan paling banyak 8 desimal` ditambahkan ke masalah. Tanpa float.
func UraiDesimalIsian(jalur, teks string, masalah *[]string) *apd.Decimal {
	if teks == "" {
		return nil
	}
	if !polaDesimal.MatchString(teks) {
		*masalah = append(*masalah, jalur+pesanDesimal)
		return nil
	}
	d, err := utils.ParseDecimal(teks)
	if err != nil {
		*masalah = append(*masalah, jalur+pesanDesimal)
		return nil
	}
	return d
}

// GalatIsianObjek - masalah urai isian objek -> galat 400 (ErrMasukanObjek); nil bila tidak ada masalah.
func GalatIsianObjek(masalah []string) error {
	if len(masalah) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrMasukanObjek, strings.Join(masalah, "; "))
}

// desimalSah - nil (kosong) atau desimal yang muat aturan polaDesimal (pemeriksaan ulang di services untuk pemanggil
// selain handler).
func desimalSah(d *apd.Decimal) bool {
	return d == nil || polaDesimal.MatchString(utils.FormatDecimal(d))
}

// nolBila - nil -> 0 (penjumlahan Loss Ratio: kosong dihitung nol, seperti Pega).
func nolBila(d *apd.Decimal) *apd.Decimal {
	if d == nil {
		return apd.New(0, 0)
	}
	return d
}
