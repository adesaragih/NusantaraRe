package repository

// Uji SQL kesiapan Confirm (06-10-2026). TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSqlCacahKesiapan(t *testing.T) {
	q := sqlCacahKesiapan("S.T_PREMIUM_LIST_DETAIL", "S.T_PREMIUM_LIST_SPREADING", "S.T_PREMIUM_LIST_SUMMARY")
	if err := db.PeriksaSQL(q); err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"FROM S.T_PREMIUM_LIST_DETAIL d WHERE d.PREMIUM_LIST_ID = :1",
		"NVL(d.RETROCEDED_SHARE, 0) > 0", "NOT EXISTS (SELECT 1 FROM S.T_PREMIUM_LIST_SPREADING s WHERE s.DETAIL_ID = d.ID)",
		"FROM S.T_PREMIUM_LIST_SUMMARY m WHERE m.PREMIUM_LIST_ID = :3"} {
		if !strings.Contains(q, mau) {
			t.Errorf("kesiapan tanpa %q:\n%s", mau, q)
		}
	}
}
