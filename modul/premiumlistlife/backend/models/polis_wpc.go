package models

// WPC polis - pengganti `WPCLife_Act` + `POOLDATA.GETQUARTER` /
// `GETQUARTERRETRO` (keputusan work owner 03-10-2026: "dibuat saja tanpa perlu
// menggunakan procedure itu lagi").
//
// Aturannya, dari kuartal TANGGAL PRODUKSI (periode `TANGGAL_CLOSING` saat
// polis jadi - sama dengan periode PL Number):
//
//	kuartal produksi   QR / QP                  TP / TR
//	Q1 (Jan-Mar)       30 April                 31 Mei
//	Q2 (Apr-Jun)       31 Juli                  31 Agustus
//	Q3 (Jul-Sep)       31 Oktober               30 November
//	Q4 (Okt-Des)       31 Januari tahun +1      28/29 Februari tahun +1
//
// Yaitu: AKHIR BULAN, satu (QR/QP) atau dua (TP/TR) bulan sesudah kuartalnya
// berakhir. Q4 jatuh di tahun berikutnya - padanan `WPCLife_Act` langkah
// "jika bulan 10,11,12" yang menambah tahun satu.

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ErrWPCTakDapatDihitung - Type atau periode di luar yang dikenal.
var ErrWPCTakDapatDihitung = errors.New("models: WPC tidak dapat dihitung")

// WPCPolis menghitung WPC dari Type polis dan periode produksi `MM.YYYY`.
func WPCPolis(tipe, periodeMMYYYY string) (time.Time, error) {
	var geser int
	switch strings.TrimSpace(tipe) {
	case "QR", "QP":
		geser = 1
	case "TP", "TR":
		geser = 2
	default:
		return time.Time{}, fmt.Errorf("%w: Type %q", ErrWPCTakDapatDihitung, tipe)
	}
	bagian := strings.Split(strings.TrimSpace(periodeMMYYYY), ".")
	if len(bagian) != 2 {
		return time.Time{}, fmt.Errorf("%w: periode %q bukan MM.YYYY", ErrWPCTakDapatDihitung, periodeMMYYYY)
	}
	bulan, err1 := strconv.Atoi(bagian[0])
	tahun, err2 := strconv.Atoi(bagian[1])
	if err1 != nil || err2 != nil || bulan < 1 || bulan > 12 || tahun < 1 {
		return time.Time{}, fmt.Errorf("%w: periode %q bukan MM.YYYY", ErrWPCTakDapatDihitung, periodeMMYYYY)
	}
	akhirKuartal := ((bulan-1)/3 + 1) * 3 // 3, 6, 9, 12
	sasaran := akhirKuartal + geser       // 13/14 = Januari/Februari tahun berikutnya
	// Hari ke-0 bulan SESUDAH sasaran = hari terakhir bulan sasaran; `time.Date`
	// menormalkan bulan > 12 ke tahun berikutnya (dan 29 Februari kabisat).
	return time.Date(tahun, time.Month(sasaran+1), 0, 0, 0, 0, 0, time.UTC), nil
}
