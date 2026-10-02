//go:build gagalkompilasi

// Berkas ini WAJIB GAGAL dikompilasi, dan hanya ikut dibangun dengan tag
// `gagalkompilasi` (TestUangTambahRasioGagalKompilasi). Ia membuktikan jaminan
// tipe ADR-F-0004 untuk tipe LOKAL `rasio`, yang tidak terjangkau dari
// `testdata/`: uang dan rasio tidak dapat dijumlahkan, dan rasio tidak dapat
// dipakai sebagai uang.

package premium

import "nusantarare/inti/backend/uang"

var _ = uang.Money{} + rasio{}

func init() {
	_, _ = uang.Money{}.Add(rasio{})
}
