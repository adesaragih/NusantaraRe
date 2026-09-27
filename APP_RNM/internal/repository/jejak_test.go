package repository

// Perekam jejak audit - TANPA Oracle.
//
// Pemilik: A2 (butir am).

import (
	"strings"
	"testing"
)

// TestSQLJejakBerskemaDanBerkolomLengkap - bentuk pernyataannya.
//
// ⛔ Diperiksa lewat SQL yang benar-benar dirakit, bukan lewat daftar nama:
// kolom yang lupa di-bind tidak terlihat dari daftar mana pun.
func TestSQLJejakBerskemaDanBerkolomLengkap(t *testing.T) {
	q := sqlSisipJejak("UJI.T_CLAIMLF_JEJAK")
	if err := PeriksaSQL(q); err != nil {
		t.Fatalf("pernyataan jejak tidak lolos PeriksaSQL: %v", err)
	}
	// Ketujuh kolom tabel wajib disebut - ADR-U-0007 menuntut SIAPA dan KAPAN,
	// dan keduanya hilang bila satu kolom terlewat.
	for _, kolom := range []string{
		"ID", "ADJUSTMENT_ID", "KLAIM_ID", "DARI", "KE", "AKUN_ID", "WAKTU",
	} {
		if !strings.Contains(q, kolom) {
			t.Errorf("pernyataan jejak tidak menyebut kolom %q: %s", kolom, q)
		}
	}
	// ⛔ Nol COMMIT (ADR-U-0029) dan nama tabel datang dari Qualify.
	if strings.Contains(strings.ToUpper(q), "COMMIT") {
		t.Error("pernyataan jejak memuat COMMIT")
	}
	if !strings.Contains(q, "UJI.T_CLAIMLF_JEJAK") {
		t.Errorf("nama tabel tidak berskema: %s", q)
	}
	// Tujuh kolom, tujuh penanda bind.
	for i := 1; i <= 7; i++ {
		if !strings.Contains(q, ":"+string(rune('0'+i))) {
			t.Errorf("penanda bind :%d tidak ada: %s", i, q)
		}
	}
}
