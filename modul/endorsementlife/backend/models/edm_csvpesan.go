package models

// Kalimat penolakan CSV untuk user awam (keputusan work owner 10-10-2026: "tampilan dan kalimat nya disederhanakan
// agar lebih mudah di mengerti oleh user") - SATU sumber untuk Upload CSV dan Add CSV Data. Aturan validasi tidak
// berubah; hanya kalimatnya. Nama kolom disebut apa adanya, dengan contoh format bila relevan.

import "fmt"

// Kalimat per baris; %s pertama = nama kolom.
const (
	PesanCSVTanggal     = "%s must be a date in DD/MM/YYYY format, e.g. 22/11/1990."
	PesanCSVRibuan      = "%s must not use thousands separators, e.g. 1500000.50 or 1500000,50."
	PesanCSVBukanAngka  = "%s must be a number, e.g. 1500000.50."
	PesanCSVDesimal     = "%s has more than 8 decimal places."
	PesanCSVBulat       = "%s must be a whole number (max 5 digits)."
	PesanCSVFactor      = "%s must be a number with a dot as decimal separator, e.g. 0.75."
	PesanCSVTeksPanjang = "%s is too long (max 255 characters)."
	PesanCSVJumlahMedan = "This row has %d values, but the header has %d columns."
	PesanCSVJudulKurang = "Your file is missing these columns: %s."
	PesanCSVJudulGanda  = "Column %s appears more than once in the header."
	PesanCSVJudulKosong = "The header has an empty column name (column %d)."
)

// PesanGalatCSV - galat CSV berkalimat layar bersih: `Error()` tetap membawa galat asalnya (log, `errors.Is`),
// `PesanLayar()` hanya kalimat untuk user (tanpa awalan teknis seperti "models:" / "invalid CSV header:").
type PesanGalatCSV struct {
	Asal  error
	Pesan string
}

func (g PesanGalatCSV) Error() string      { return g.Asal.Error() + ": " + g.Pesan }
func (g PesanGalatCSV) Unwrap() error      { return g.Asal }
func (g PesanGalatCSV) PesanLayar() string { return g.Pesan }

func galatCSV(asal error, format string, a ...any) error {
	return PesanGalatCSV{Asal: asal, Pesan: fmt.Sprintf(format, a...)}
}
