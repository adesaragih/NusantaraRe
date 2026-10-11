package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// 11-10-2026 - Recommendation (tombol Save Medical Check): peserta aktif, dibekukan bila diputus (b7323).
func TestSQLRekomendasi(t *testing.T) {
	q := strings.Join(strings.Fields(sqlSetelRekomendasi("S.P")), " ")
	for _, mau := range []string{"UPDATE S.P SET RECOMMENDATION = :1", "WHERE ID = :2 AND CLAIM_ID = :3",
		"STS_HAPUS IS NULL", "(STS_REJECT IS NULL OR STS_REJECT NOT IN ('1', '2'))"} {
		if !strings.Contains(q, mau) {
			t.Errorf("Recommendation tanpa %q:\n%s", mau, q)
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Errorf("PeriksaSQL: %v", err)
	}
}

func TestKolomPesertaMembacaRekomendasi(t *testing.T) {
	for _, k := range kolomPeserta {
		if k.Nama == "RECOMMENDATION" {
			return
		}
	}
	t.Error("RECOMMENDATION tidak dibaca dari T_CLAIMLF_PREMIUMLIST_DETAIL")
}
