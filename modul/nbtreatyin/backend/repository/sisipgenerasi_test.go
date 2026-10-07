package repository

// Kolom IDPEGA dibuang dari T_GENERAL_POLIS_TREATY (keputusan work owner 06-10-2026): penyisipan generasi
// kasus baru tidak boleh menyebutnya lagi - ORA-00904 di skema yang kolomnya sudah di-drop.

import (
	"strings"
	"testing"
)

func TestSisipGenerasiTanpaIDPEGA(t *testing.T) {
	q := sqlSisipGenerasi("UJI_SKEMA.T_GENERAL_POLIS_TREATY")
	if strings.Contains(q, "IDPEGA") {
		t.Fatalf("INSERT generasi masih menyebut IDPEGA:\n%s", q)
	}
	if strings.Contains(q, ":4") || !strings.Contains(q, ":3") {
		t.Errorf("INSERT generasi harus tepat tiga parameter (ID, USERNAME, POSITION_NOTE):\n%s", q)
	}
}
