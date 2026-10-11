package models

// TGL_PROD `JSON_POLIS` endorsement - keputusan work owner 08-10-2026.
//
// Tanggal produksi saat Confirm (saat NO_ENDORS terbit), dari periode produksi
// TANGGAL_CLOSING (`penomor.HitungPeriodeNomor` inti, dipanggil services):
//
//	periode = bulan Confirm (Jakarta)  → TGL_PROD = tanggal Confirm (Jakarta)
//	periode bergeser ke bulan berikut  → TGL_PROD = tanggal 1 bulan periode
//
// Aturan SAMA dengan PremiumList Life (`premiumlistlife/backend/models/polis_tglprod.go`)
// tetapi ditulis sendiri - nol impor lintas modul. Tanpa jam.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ErrTanggalProduksi - periode bukan MM.YYYY.
var ErrTanggalProduksi = errors.New("models: production date cannot be determined")

// TanggalProduksi menghitung TGL_PROD dari saat Confirm dan periode produksi `MM.YYYY`.
func TanggalProduksi(saat time.Time, periodeMMYYYY string) (time.Time, error) {
	bagian := strings.Split(strings.TrimSpace(periodeMMYYYY), ".")
	if len(bagian) != 2 {
		return time.Time{}, fmt.Errorf("%w: period %q is not MM.YYYY", ErrTanggalProduksi, periodeMMYYYY)
	}
	bulan, err1 := strconv.Atoi(bagian[0])
	tahun, err2 := strconv.Atoi(bagian[1])
	if err1 != nil || err2 != nil || bulan < 1 || bulan > 12 || tahun < 1 {
		return time.Time{}, fmt.Errorf("%w: period %q is not MM.YYYY", ErrTanggalProduksi, periodeMMYYYY)
	}
	lokal := saat.In(zonaJakartaProduksi())
	if lokal.Year() == tahun && int(lokal.Month()) == bulan {
		return time.Date(lokal.Year(), lokal.Month(), lokal.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC), nil
}

// zonaJakartaProduksi - Asia/Jakarta, cadangan UTC+7 (pola `penomor.zonaJakarta`).
func zonaJakartaProduksi() *time.Location {
	if loc, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return loc
	}
	return time.FixedZone("WIB", 7*60*60)
}
