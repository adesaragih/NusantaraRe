package repository

// Perekam jejak audit - TANPA Oracle.
//
// Pemilik: A2 (butir am).

import (
	"os"
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
	// Kedelapan kolom tabel wajib disebut - ADR-U-0007 menuntut SIAPA dan KAPAN,
	// dan keduanya hilang bila satu kolom terlewat. KOMENTAR sejak migrasi 021
	// (OQ-M5, GILIRAN-17): alasan penolakan Admin.
	for _, kolom := range []string{
		"ID", "ADJUSTMENT_ID", "KLAIM_ID", "DARI", "KE", "AKUN_ID", "WAKTU", "KOMENTAR",
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
	// Delapan kolom, delapan penanda bind.
	for i := 1; i <= 8; i++ {
		if !strings.Contains(q, ":"+string(rune('0'+i))) {
			t.Errorf("penanda bind :%d tidak ada: %s", i, q)
		}
	}
}

// OQ-M5 DITUTUP (GILIRAN-17): alasan penolakan disimpan di kolom komentar
// baru jejak klaim - migrasi 021, lebar dari preseden KOMITE_COMMENT (013),
// sebab korpus tidak mengekspor rule Property `Remarks`/`KomiteComment`.
func TestMigrasi021KolomKomentarJejak(t *testing.T) {
	naik, err := os.ReadFile("migrations/021_kolom_komentar_jejak.sql")
	if err != nil {
		t.Fatal(err)
	}
	turun, err := os.ReadFile("migrations/021_kolom_komentar_jejak_down.sql")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(naik), "ALTER TABLE {skema}.T_CLAIMLF_JEJAK ADD (") ||
		!strings.Contains(string(naik), "KOMENTAR VARCHAR2(4000)") {
		t.Errorf("021 tidak menambah KOMENTAR VARCHAR2(4000):\n%s", naik)
	}
	if !strings.Contains(string(turun), "ALTER TABLE {skema}.T_CLAIMLF_JEJAK DROP (") ||
		!strings.Contains(string(turun), "KOMENTAR") {
		t.Errorf("021_down tidak membuang KOMENTAR:\n%s", turun)
	}
}
