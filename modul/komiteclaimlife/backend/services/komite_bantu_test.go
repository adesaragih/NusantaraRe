package services

// Pembantu uji milik Komite Claim Life - TANPA Oracle.
//
// Refactor bentuk B (30-09-2026): uji Komite dulu meminjam `badanFungsi` dari
// uji PremiumList (`polis_summary_test.go`). PremiumList kini modul sendiri;
// salinannya ikut Komite. Isinya sama.

import (
	"strings"
	"testing"
)

// badanFungsi memotong badan satu fungsi dari teks sumber, mulai `kepala`.
func badanFungsi(t *testing.T, teks, kepala string) string {
	t.Helper()
	i := strings.Index(teks, kepala)
	if i < 0 {
		t.Fatalf("fungsi %q tidak ditemukan", kepala)
	}
	j := strings.Index(teks[i:], "\n}\n")
	if j < 0 {
		t.Fatalf("ujung fungsi %q tidak ditemukan", kepala)
	}
	return teks[i : i+j]
}
