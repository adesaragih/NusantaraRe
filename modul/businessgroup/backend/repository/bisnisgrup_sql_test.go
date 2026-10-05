package repository

import (
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// ID baru seperti UPSERT_BUSINESSGROUP: situs aktif + BUSINESSGROUP_SEQ 4 digit; situs diikat, bukan literal.
func TestSqlIDBaru(t *testing.T) {
	q := SqlIDBaru("S.M_SITE_DATABASE", "S.BUSINESSGROUP_SEQ")
	if mau := "SELECT ID || LPAD(TO_CHAR(S.BUSINESSGROUP_SEQ.NEXTVAL), 4, '0') FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1"; q != mau {
		t.Errorf("ID baru\ndapat %s\nmau   %s", q, mau)
	}
}

// Daftar: tanpa SYARIAH (saringan Pega NotEndsWith), saring Treaty Group lewat bind, urut nama.
func TestSqlDaftar(t *testing.T) {
	q := satuBaris(SqlDaftar("S.BUSINESSGROUP"))
	for _, b := range []string{"WHERE NVL(UPPER(TRIM(NOTE)), ' ') NOT LIKE '%SYARIAH'", "(:6 IS NULL OR TOPID = :7)",
		"UPPER(ALIASNAME) LIKE :4 ESCAPE '\\'", "ORDER BY UPPER(NOTE), ID"} {
		if !strings.Contains(q, b) {
			t.Errorf("daftar tanpa %q: %s", b, q)
		}
	}
	if q := SqlDaftarTreaty("S.TREATYGROUP"); q != "SELECT ID, TREATYGROUPNAME FROM S.TREATYGROUP ORDER BY UPPER(TREATYGROUPNAME), ID" {
		t.Errorf("treaty %s", q)
	}
	for _, tb := range DaftarTabelDitulis {
		if tb != Tabel {
			t.Errorf("modul ini hanya menulis %s, bukan %s", Tabel, tb)
		}
	}
}
