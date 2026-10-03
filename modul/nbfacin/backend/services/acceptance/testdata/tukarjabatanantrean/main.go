// Berkas ini WAJIB GAGAL dikompilasi (tiket NB-11): antrean dan kode jabatan
// adalah dua tipe berbeda; menukarnya ditangkap kompilator. Dibangun oleh
// TestJabatanAntreanGagalKompilasi; `go build ./...` melewati folder testdata.
package main

import "nusantarare/modul/nbfacin/backend/services/acceptance"

func main() {
	var t acceptance.Transisi
	t.JabatanTujuan = t.Antrean
	t.Antrean = t.JabatanTujuan
}
