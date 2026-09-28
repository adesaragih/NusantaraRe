package utils

// Nol NaN dan nol Infinity yang lolos menjadi desimal - 28-09-2026.
//
// ⛔ SEBAB PENJAGA INI ADA. `apd.NewFromString` MENERIMA "NaN", "Infinity",
// "Inf", dan "-Infinity": ia mengikuti spesifikasi desimal, dan di sana
// keduanya nilai yang sah. Di sistem ini tidak - uang tidak pernah tak-hingga,
// dan pangsa tidak pernah bukan-bilangan.
//
// Ia ketahuan lewat uji unggahan CSV tiket 04, bukan lewat tinjauan: satu uji
// yang sengaja mencoba nilai aneh menemukan bahwa keempatnya lolos sampai ke
// nilai uang. Fungsi ini satu-satunya jalan masuk teks-ke-desimal
// (ADR-U-0034), dan salah satu pemanggilnya membaca uang dari JSON - yaitu
// batas yang dilewati permintaan dari luar.

import (
	"errors"
	"testing"
)

func TestParseDecimalMenolakYangTidakBerhingga(t *testing.T) {
	for _, s := range []string{
		"NaN", "nan", "-NaN", "sNaN",
		"Infinity", "infinity", "-Infinity", "Inf", "-Inf", "+Inf",
	} {
		d, err := ParseDecimal(s)
		if err == nil {
			t.Errorf("%q diterima menjadi %v; mau ditolak", s, FormatDecimal(d))
			continue
		}
		if !errors.Is(err, ErrBukanDesimal) {
			t.Errorf("%q: galat %v, mau ErrBukanDesimal", s, err)
		}
	}
}

func TestParseDecimalTetapMenerimaBilanganBiasa(t *testing.T) {
	// ⚠️ Penjaga yang menolak terlalu banyak sama berbahayanya. Bentuk di
	// bawah ini SAH dan harus tetap lolos - termasuk notasi eksponen, yang
	// dipakai sebagian pengekspor.
	for _, s := range []string{
		"0", "-0", "1", "-1.5", "0.00000001", "1e5", "-2E-3",
		"12345678901234567890.12345678",
	} {
		if _, err := ParseDecimal(s); err != nil {
			t.Errorf("%q ditolak: %v", s, err)
		}
	}
}
