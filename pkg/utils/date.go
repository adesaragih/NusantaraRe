package utils

import (
	"errors"
	"fmt"
	"time"
)

// Dua bentuk tanggal yang dipakai sistem lama (ADR-U-0022). Keduanya
// dinyatakan di sini supaya tidak ada layout tanggal yang ditulis ulang
// tersebar di handler atau repository.
const (
	// TanggalSaja untuk kolom bertipe tanggal.
	TanggalSaja = "2006-01-02"
	// TanggalWaktu untuk kolom bertipe stempel waktu.
	TanggalWaktu = "2006-01-02 15:04:05"
)

// ErrBukanTanggal dikembalikan bila teks tidak cocok dengan kedua bentuk.
var ErrBukanTanggal = errors.New("teks bukan tanggal yang dikenal")

// ParseTanggal membaca teks dalam salah satu dari dua bentuk di atas.
// Konversi dilakukan sekali saat masuk (ADR-U-0022).
func ParseTanggal(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("%w: teks kosong", ErrBukanTanggal)
	}
	for _, layout := range []string{TanggalWaktu, TanggalSaja} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("%w: %q", ErrBukanTanggal, s)
}

// FormatTanggal menulis tanggal tanpa komponen waktu.
// Waktu nol menjadi teks kosong, bukan "0001-01-01" (ADR-U-0027).
func FormatTanggal(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TanggalSaja)
}

// FormatTanggalWaktu menulis tanggal beserta waktunya.
func FormatTanggalWaktu(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TanggalWaktu)
}
