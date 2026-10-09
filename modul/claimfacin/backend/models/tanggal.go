package models

// Untuk apa berkas ini: TANGGAL HALAMAN. Halaman menyimpan tanggal sebagai teks "2006-01-02" (Date) atau
// "2006-01-02 15:04:05" (DateTime, zona Jakarta) - satu bentuk pertukaran dengan repository (`utils.TanggalSaja` /
// `utils.TanggalWaktu`). Dokumen Pega lama memakai "yyyyMMdd" dan "yyyyMMddTHHmmss.SSS GMT" (katalog DEV
// JSON_KLAIM 07-10-2026, AC 12); `TanggalPega` membaca keduanya.

import (
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/utils"
)

// UraiTanggal membaca teks tanggal halaman (Date atau DateTime) atau tanggal Pega. ok=false bila kosong/tak terurai.
func UraiTanggal(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	if t, err := utils.ParseTanggal(s); err == nil {
		return t, true
	}
	if t, ok := TanggalPega(s); ok {
		return t, true
	}
	return time.Time{}, false
}

// TanggalPega membaca "yyyyMMdd" atau "yyyyMMddTHHmmss.SSS GMT" (GMT diubah ke Jakarta).
func TanggalPega(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if t, err := time.Parse("20060102", s); err == nil {
		return t, true
	}
	if strings.HasSuffix(s, " GMT") {
		if t, err := time.Parse("20060102T150405.000 MST", s); err == nil {
			l := t.In(Jakarta)
			return time.Date(l.Year(), l.Month(), l.Day(), l.Hour(), l.Minute(), l.Second(), 0, time.UTC), true
		}
	}
	return time.Time{}, false
}

// hariSaja memotong waktu - pembandingan `@CompareDates` dilakukan pada tingkat HARI.
//
// `[dugaan]` Isi fungsi `@CompareDates` Pega tidak ada di ekspor; activity membandingkan DateTime (DateOfLoss,
// DateReceived) dengan `@CurrentDate("dd/MM/yyyy")` (tanpa jam). Membandingkan pada tingkat hari mencegah DOL hari ini
// pukul 10.00 dianggap "lebih dari hari ini".
func hariSaja(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// SesudahTanggal meniru `@CompareDates(a, b) == true`: a lebih kemudian dari b (tingkat hari). Kosong/tak terurai =
// false (pembandingan Pega atas properti kosong tidak pernah benar di activity yang dibaca).
func SesudahTanggal(a, b string) bool {
	ta, oka := UraiTanggal(a)
	tb, okb := UraiTanggal(b)
	if !oka || !okb {
		return false
	}
	return hariSaja(ta).After(hariSaja(tb))
}

// SamaTanggal - kedua tanggal pada hari yang sama.
func SamaTanggal(a, b string) bool {
	ta, oka := UraiTanggal(a)
	tb, okb := UraiTanggal(b)
	return oka && okb && hariSaja(ta).Equal(hariSaja(tb))
}

// TeksTanggal menulis tanggal halaman ("2006-01-02").
func TeksTanggal(t time.Time) string { return t.Format("2006-01-02") }

// TambahTahun meniru `@DateTime.addCalendar(x, n, 0, …)` - bentuk teks masukan dipertahankan (Date atau DateTime).
func TambahTahun(s string, n int) string {
	t, ok := UraiTanggal(s)
	if !ok {
		return ""
	}
	h := t.AddDate(n, 0, 0)
	if len(strings.TrimSpace(s)) > 10 {
		return h.Format("2006-01-02 15:04:05")
	}
	return TeksTanggal(h)
}

// YMD menulis tanggal sebagai "yyyyMMdd" - bentuk `@substring(@FormatDateTime(x,"yyyyMMdd …"),0,8)` (SaveOutstanding_Act
// langkah 29, parameter TreatyYearTreatyin_SQL).
func YMD(s string) string {
	t, ok := UraiTanggal(s)
	if !ok {
		return ""
	}
	return t.Format("20060102")
}

// TambahBulan meniru `@addCalendar(x, 0, bulan, 0, ...)`: bulan kosong / tak terurai = 0. Bentuk teks masukan
// dipertahankan.
func TambahBulan(s, bulan string) string {
	t, ok := UraiTanggal(s)
	if !ok {
		return ""
	}
	n, err := strconv.Atoi(strings.TrimSpace(bulan))
	if err != nil {
		n = 0
	}
	h := t.AddDate(0, n, 0)
	if len(strings.TrimSpace(s)) > 10 {
		return h.Format("2006-01-02 15:04:05")
	}
	return TeksTanggal(h)
}
