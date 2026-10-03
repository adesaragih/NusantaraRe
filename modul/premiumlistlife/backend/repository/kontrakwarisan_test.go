package repository

import (
	"testing"
)

// TestPenulisWarisanMengisiKunciBacaKlaim - kunci baca = kunci tulis.
//
// Claim Life membaca berkunci `PL_NUMBER` (+ `CERTIFICATE_NO`); penulis
// mengisi `PL_NUMBER` dengan nomor yang BARU terbit di transaksi yang sama,
// bukan dari kolom sumber yang mungkin masih kosong.
func TestPenulisWarisanMengisiKunciBacaKlaim(t *testing.T) {
	for _, k := range kolomPesertaWarisan {
		switch k.Kolom {
		case "PL_NUMBER":
			if k.Sumber != "" {
				t.Errorf("PL_NUMBER bersumber %q; ia harus nomor yang dirakit nilaiSalinWarisan", k.Sumber)
			}
		case "CERTIFICATE_NO":
			if k.Sumber != "d.CERTIFICATE_NO" {
				t.Errorf("CERTIFICATE_NO bersumber %q", k.Sumber)
			}
		}
	}
	arg := nilaiSalinWarisan("UJI-PL-9", "UJI-W", BarisWarisan{})
	for i, k := range kolomPesertaWarisan {
		if k.Kolom == "PL_NUMBER" && arg[i] != "UJI-PL-9" {
			t.Errorf("PL_NUMBER dikirim %v", arg[i])
		}
	}
}
