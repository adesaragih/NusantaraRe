package models

// Nomor akseptasi Komite - tiket 04a. TANPA Oracle.

import (
	"errors"
	"testing"

	"nusantarare/inti/penomor"
)

// TestNomorAkseptasiKomiteDariLiteral - langkah 4.11 / 4.12.
func TestNomorAkseptasiKomiteDariLiteral(t *testing.T) {
	for _, u := range []struct{ tipe, mau string }{
		{"QR", "RNML-AL01.09.26.00007"},
		{"QP", "RNML-AL01.09.26.00007"},
		{"TP", "RNML-ARL01.09.26.00007"},
		{"TR", "RNML-ARL01.09.26.00007"},
	} {
		got, err := NomorAkseptasiKomite("RNML-", u.tipe, "L01", "09.2026", 7)
		if err != nil {
			t.Fatalf("%s: %v", u.tipe, err)
		}
		if got != u.mau {
			t.Errorf("%s: %q, mau %q", u.tipe, got, u.mau)
		}
	}
}

// TestSatuPenghitungUntukKeduaCabang - JENIS disusun di 4.8, sebelum cabang.
func TestSatuPenghitungUntukKeduaCabang(t *testing.T) {
	if JenisPenghitungKomite("RNML-") != "RNML-A" {
		t.Errorf("jenis %q, mau RNML-A (SUMBER-PENOMORAN-DBA)", JenisPenghitungKomite("RNML-"))
	}
	if ClassPenghitungKomiteLife != "ASM-FW-GCNMFW-Work-KomiteLife" {
		t.Errorf("class %q", ClassPenghitungKomiteLife)
	}
}

// TestBahanKosongDitolak - nomor tanpa satu bagiannya tetap tampak sah.
func TestBahanKosongDitolak(t *testing.T) {
	if _, err := NomorAkseptasiKomite("", "QR", "L01", "09.2026", 1); !errors.Is(err, penomor.ErrAwalanProduksiKosong) {
		t.Errorf("awalan kosong: %v", err)
	}
	if _, err := NomorAkseptasiKomite("RNML-", "FAC", "L01", "09.2026", 1); !errors.Is(err, penomor.ErrTipePLTanpaCabang) {
		t.Errorf("tipe asing: %v", err)
	}
	if _, err := NomorAkseptasiKomite("RNML-", "QR", " ", "09.2026", 1); !errors.Is(err, penomor.ErrKodeBisnisKosong) {
		t.Errorf("kode bisnis kosong: %v", err)
	}
	if _, err := NomorAkseptasiKomite("RNML-", "QR", "L01", "9.26", 1); err == nil {
		t.Error("periode tak berbentuk diterima")
	}
}
