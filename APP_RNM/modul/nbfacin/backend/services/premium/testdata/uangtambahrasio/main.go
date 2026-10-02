// Berkas ini WAJIB GAGAL dikompilasi (tiket NB-01, ADR-F-0004): uang dan
// rasio adalah dua tipe berbeda yang tidak dapat dijumlahkan. Dibangun oleh
// TestUangTambahRasioGagalKompilasi; `go build ./...` melewati folder testdata.
package main

import "nusantarare/inti/backend/uang"

func main() {
	m, r := uang.Money{}, uang.Ratio{}
	_ = m + r
	_, _ = m.Add(r)
}
