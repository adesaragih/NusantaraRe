package repository

import (
	"strings"
	"testing"
)

// ID baru persis `GetIDConsultanAdj_SQL` Pega: ID situs aktif + sequence 4 digit; situs aktif diikat, bukan literal.
func TestSqlIDBaru(t *testing.T) {
	q := SqlIDBaru("S.M_SITE_DATABASE", "S.ADJUSTERCONSULTANT_SEQ")
	mau := "SELECT ID || LPAD(TO_CHAR(S.ADJUSTERCONSULTANT_SEQ.NEXTVAL), 4, '0') FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1"
	if q != mau {
		t.Errorf("SQL ID baru\ndapat %s\nmau   %s", q, mau)
	}
}

// Daftar: cari "memuat" ber-ESCAPE di empat kolom, saring status lewat bind, urut nama, berbatas.
func TestSqlDaftar(t *testing.T) {
	q := strings.Join(strings.Fields(SqlDaftar("S.ADJUSTERCONSULTANT")), " ")
	for _, bagian := range []string{"FROM S.ADJUSTERCONSULTANT", "UPPER(NAME) LIKE :3 ESCAPE '\\'", "(:6 IS NULL OR IS_ACTIVE = :7)",
		"ORDER BY UPPER(NAME), ID FETCH FIRST :8 ROWS ONLY"} {
		if !strings.Contains(q, bagian) {
			t.Errorf("SQL daftar tanpa %q:\n%s", bagian, q)
		}
	}
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` {
		t.Errorf("pola %v", p)
	}
	if PolaCari("  ") != nil {
		t.Error("kata kosong = nil")
	}
}
