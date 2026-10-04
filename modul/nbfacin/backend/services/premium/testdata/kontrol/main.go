// Kontrol untuk TestUangTambahRasioGagalKompilasi: berkas ini WAJIB
// terkompilasi. Bila ia gagal, kegagalan berkas pasangannya tidak membuktikan
// apa-apa - pembangunnya yang rusak.
package main

import "nusantarare/inti/backend/uang"

func main() {
	m := uang.Money{}
	_, _ = m.Add(uang.Money{})
}
