package repository

import (
	"strings"
	"time"
)

// tanggalRD - tanggal parameter RD ke bentuk TERSIMPAN `YYYYMMDD`, bentuk
// `TREATYYEAR.STARTDATE/ENDDATE` yang dibandingkan sebagai TEKS.
//
// ⭐ 9 Oktober 2026 - laporan pemakai: dropdown Reins Type spreading layar
// Adjustment KOSONG untuk Treaty Group PROPERTY. Commencement panel New yang
// sudah disunting tiba sebagai bentuk kabel (`DD-MM-YYYY`) atau bentuk kotak
// tanggal (`YYYY-MM-DD`); `'20260701' <= '08-10-2026'` sebagai teks SALAH,
// sehingga setiap susunan tersaring habis.
//
// Bentuk yang tidak dikenali lewat apa adanya (perilaku lama).
func tanggalRD(s string) string {
	t := strings.TrimSpace(s)
	for _, pola := range []string{"20060102", "2006-01-02", "02-01-2006", "02/01/2006"} {
		if d, err := time.Parse(pola, t); err == nil {
			return d.Format("20060102")
		}
	}
	// Stempel Pega `20260701T000000.000 GMT` - delapan aksara pertama.
	if len(t) > 8 {
		if d, err := time.Parse("20060102", t[:8]); err == nil {
			return d.Format("20060102")
		}
	}
	return t
}
