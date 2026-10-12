package models

// TGL_PROD `JSON_POLIS` - keputusan work owner 06-10-2026.
//
// Untuk apa berkas ini: padanan `ProdDateTime` Pega (`InsertJsonPolisLife_Act`
// langkah 2/4 → CARI36 → `P_TGL_PROD`) dengan aturan periode yang SAMA dengan
// PL Number (`penomor.HitungPeriodeNomor`, TANGGAL_CLOSING) - bukan ambang 25
// tertanam (`AmbangProdDateTimePega`, polis_periode.go).
//
// Pega: tanggal hari ini; bila hari > 25, tanggal 1 bulan berikutnya. Di sini:
//
//	periode nomor = bulan transaksi (Jakarta)  → TGL_PROD = tanggal transaksi (Jakarta)
//	periode nomor ≠ bulan transaksi             → TGL_PROD = tanggal 1 bulan periode
//
// Tanpa jam (DATE dari `To_date(..,'DD/MM/YYYY')`).

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ErrTanggalProduksi - periode nomor bukan MM.YYYY.
var ErrTanggalProduksi = errors.New("models: tanggal produksi tidak dapat dihitung")

// TanggalProduksi menghitung TGL_PROD dari saat Confirm dan periode PL Number (`MM.YYYY`).
func TanggalProduksi(saat time.Time, periodeMMYYYY string) (time.Time, error) {
	bagian := strings.Split(strings.TrimSpace(periodeMMYYYY), ".")
	if len(bagian) != 2 {
		return time.Time{}, fmt.Errorf("%w: periode %q bukan MM.YYYY", ErrTanggalProduksi, periodeMMYYYY)
	}
	bulan, err1 := strconv.Atoi(bagian[0])
	tahun, err2 := strconv.Atoi(bagian[1])
	if err1 != nil || err2 != nil || bulan < 1 || bulan > 12 || tahun < 1 {
		return time.Time{}, fmt.Errorf("%w: periode %q bukan MM.YYYY", ErrTanggalProduksi, periodeMMYYYY)
	}
	lokal := saat.In(zonaJakartaTglProd())
	if lokal.Year() == tahun && int(lokal.Month()) == bulan {
		return time.Date(lokal.Year(), lokal.Month(), lokal.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC), nil
}

// zonaJakartaTglProd - Asia/Jakarta, cadangan UTC+7 (pola penomor.zonaJakarta).
func zonaJakartaTglProd() *time.Location {
	if loc, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		return loc
	}
	return time.FixedZone("WIB", 7*60*60)
}
