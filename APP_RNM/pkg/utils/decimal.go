// Package utils memuat helper publik yang tidak bergantung pada internal/.
// Arah ketergantungan hanya satu arah: internal/ boleh mengimpor pkg/,
// pkg/ tidak pernah mengimpor internal/.
package utils

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"
)

// DecimalPrecision adalah presisi desimal seluruh sistem.
//
// Keputusan DECIDED-TEKNIS 18-09: cockroachdb/apd dengan presisi 38;
// shopspring/decimal ditolak. Angka 38 sejajar dengan kolom uang
// NUMBER(38,8) (ADR-U-0003 - ADR-U-0016).
//
// Dinyatakan DI SATU TEMPAT, yaitu di sini. Tidak ada paket lain yang boleh
// membuat apd.Context sendiri.
const DecimalPrecision = 38

var decimalContext = apd.BaseContext.WithPrecision(DecimalPrecision)

// DecimalContext mengembalikan konteks aritmetika desimal.
//
// Yang dikembalikan adalah SALINAN. Konteks apd punya medan yang dapat
// ditulis - Precision di antaranya - sehingga membagikan pointer ke satu
// nilai bersama membuat presisi 38 dapat diubah satu pemanggil untuk seluruh
// proses. Presisi dinyatakan di satu tempat, dan tetap begitu.
func DecimalContext() *apd.Context {
	salinan := *decimalContext
	return &salinan
}

// ErrBukanDesimal dikembalikan bila teks tidak dapat dibaca sebagai desimal.
var ErrBukanDesimal = errors.New("teks bukan bilangan desimal")

// ParseDecimal mengubah teks menjadi desimal. Ini adalah SATU-SATUNYA fungsi
// konversi teks-ke-desimal untuk seluruh batas procedure (ADR-U-0034), dan
// konversi tipe dilakukan sekali saat masuk (ADR-U-0022).
//
// ⛔ JANGAN memakainya untuk kode, penanda, atau enumerasi. Nilai seperti
// "006" harus tetap teks: mengubahnya menjadi bilangan mengembalikannya
// sebagai "6" dan memecahkan penggolong (ADR-U-0022).
// ⛔ `NaN` DAN `Infinity` DITOLAK, ditambahkan 28-09-2026. `apd.NewFromString`
// MENERIMA keduanya - ia mengikuti spesifikasi desimal, dan di sana keduanya
// nilai yang sah. Di sini tidak: uang tidak pernah tak-hingga, dan pangsa
// tidak pernah bukan-bilangan.
//
// Lubangnya nyata dan bukan teoretis. Fungsi ini satu-satunya jalan masuk
// teks-ke-desimal (ADR-U-0034), dan salah satu pemanggilnya adalah pembaca
// uang dari JSON - yaitu batas yang dilewati permintaan dari luar. Badan
// permintaan berisi `{"amount":"NaN"}` akan lolos seluruh validasi, lalu
// gagal di Oracle dengan galat yang tidak menyebut sebabnya - atau, lebih
// buruk, tersimpan di jalur yang tidak memeriksanya.
//
// ⚠️ Menolaknya di SINI menutup kesepuluh pemanggilnya sekaligus. Menolaknya
// di tiap pemanggil berarti sepuluh tempat untuk lupa.
func ParseDecimal(s string) (*apd.Decimal, error) {
	if s == "" {
		return nil, fmt.Errorf("%w: teks kosong", ErrBukanDesimal)
	}
	d, _, err := apd.NewFromString(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrBukanDesimal, s)
	}
	if d.Form != apd.Finite {
		return nil, fmt.Errorf("%w: %q bukan bilangan berhingga", ErrBukanDesimal, s)
	}
	return d, nil
}

// FormatDecimal mengubah desimal menjadi teks tanpa notasi eksponen.
// Ini pasangan ParseDecimal, dan satu-satunya jalan keluar ke teks.
//
// Nilai nil menjadi teks kosong, bukan "0" — kolom kosong dan kolom bernilai
// nol adalah dua hal berbeda (ADR-U-0027: seluruh kolom nullable).
func FormatDecimal(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	return d.Text('f')
}
