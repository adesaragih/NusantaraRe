package repository

import (
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// ID baru seperti PEGA_M_TREATYGROUP: situs aktif + TREATYGROUP_SEQ 4 digit; situs diikat, bukan literal.
func TestSqlIDBaru(t *testing.T) {
	q := SqlIDBaru("S.M_SITE_DATABASE", "S.TREATYGROUP_SEQ")
	if mau := "SELECT ID || LPAD(TO_CHAR(S.TREATYGROUP_SEQ.NEXTVAL), 4, '0') FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1"; q != mau {
		t.Errorf("ID baru\ndapat %s\nmau   %s", q, mau)
	}
}

// Sisip menulis kolom PEGA_TREATYGROUP + COAID; Ubah juga menulis COAID (ikut OJK), tidak pernah ID dan OLDID.
func TestSqlSisipUbah(t *testing.T) {
	s := satuBaris(SqlSisip("S.TREATYGROUP"))
	if !strings.Contains(s, "(ID, OJKBUSINESSID, OJKBUSINESSNAME, OJKBUSINESSNAMEIDN, TREATYGROUPNAME, TREATYGROUPSOANAME, TGLUPDATE, USERID, ORDERNO, COAID) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)") ||
		strings.Contains(s, "OLDID") {
		t.Errorf("sisip %s", s)
	}
	u := satuBaris(SqlUbah("S.TREATYGROUP"))
	set := strings.SplitN(u, "WHERE", 2)[0]
	if !strings.Contains(set, "COAID = :9") {
		t.Errorf("ubah tanpa COAID: %s", u)
	}
	for _, k := range []string{"OLDID", " ID ="} {
		if strings.Contains(set, k) {
			t.Errorf("ubah menyentuh %s: %s", k, u)
		}
	}
	if !strings.HasSuffix(u, "WHERE ID = :10") {
		t.Errorf("ubah %s", u)
	}
}

// Daftar bersaring OJK dan bergabung nama COA; anak tanpa SYARIAH; COA se-OJK paling sering.
func TestSqlBaca(t *testing.T) {
	d := satuBaris(SqlDaftar("S.TREATYGROUP", "S.BUSINESSGROUP"))
	for _, b := range []string{"LEFT JOIN S.BUSINESSGROUP b ON b.ID = t.COAID", "(:6 IS NULL OR t.OJKBUSINESSID = :7)",
		"ORDER BY LPAD(t.ORDERNO, 10, '0'), UPPER(t.TREATYGROUPNAME), t.ID"} {
		if !strings.Contains(d, b) {
			t.Errorf("daftar tanpa %q: %s", b, d)
		}
	}
	if a := SqlAnak("S.BUSINESSGROUP"); !strings.Contains(a, "TOPID = :1 AND NVL(UPPER(TRIM(NOTE)), ' ') NOT LIKE '%SYARIAH'") {
		t.Errorf("anak %s", a)
	}
	if c := satuBaris(SqlCoaSeOjk("S.TREATYGROUP")); !strings.Contains(c, "GROUP BY COAID ORDER BY COUNT(*) DESC, COAID FETCH FIRST 1 ROWS ONLY") {
		t.Errorf("coa %s", c)
	}
	c := satuBaris(SqlCoaPerOjk("S.TREATYGROUP", "S.BUSINESSGROUP"))
	for _, b := range []string{"ROW_NUMBER() OVER (PARTITION BY t.OJKBUSINESSID ORDER BY COUNT(*) DESC, t.COAID)",
		"LEFT JOIN S.BUSINESSGROUP b ON b.ID = t.COAID", "WHERE t.COAID IS NOT NULL", "WHERE rn = 1"} {
		if !strings.Contains(c, b) {
			t.Errorf("coa per OJK tanpa %q: %s", b, c)
		}
	}
	for _, tb := range DaftarTabelDitulis {
		if tb != Tabel {
			t.Errorf("modul ini hanya menulis %s, bukan %s", Tabel, tb)
		}
	}
}
