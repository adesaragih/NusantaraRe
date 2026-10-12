package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// 12-10-2026 - catatan Claim Analyst: NOTES peserta aktif klaim.
func TestSQLCatatanAnalis(t *testing.T) {
	q := strings.Join(strings.Fields(sqlSetelCatatanAnalis("S.P")), " ")
	if q != "UPDATE S.P SET NOTES = :1 WHERE ID = :2 AND CLAIM_ID = :3 AND STS_HAPUS IS NULL" {
		t.Errorf("SQL catatan Claim Analyst:\n%s", q)
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Errorf("PeriksaSQL: %v", err)
	}
}
