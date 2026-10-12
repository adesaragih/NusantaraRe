package backend

import (
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/kontrak"
)

// Modul ini menyatakan kontrak JenisBisnisFacIn (opsi 1a, keputusan work owner 12-10-2026); tanpa pernyataan ini
// inti.Ambil di Bangun menolak, dan lini polis sistem baru diam-diam "belum diperiksa".
func TestPendaftaranMembutuhkanJenisBisnis(t *testing.T) {
	p := Pendaftaran()
	ingin := inti.KontrakDari[kontrak.JenisBisnisFacIn]()
	for _, k := range p.Membutuhkan {
		if k == ingin {
			return
		}
	}
	t.Fatalf("Membutuhkan = %v, tanpa %v", p.Membutuhkan, ingin)
}
