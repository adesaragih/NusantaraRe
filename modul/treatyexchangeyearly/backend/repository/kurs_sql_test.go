package repository

import (
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// ID baru seperti PEGA_TREATYEXCHANGE: situs aktif + TREATYEXCHANGE_SEQ 4 digit; situs diikat, bukan literal.
func TestSqlIDBaru(t *testing.T) {
	q := SqlIDBaru("S.M_SITE_DATABASE", "S.TREATYEXCHANGE_SEQ")
	if mau := "SELECT ID || LPAD(TO_CHAR(S.TREATYEXCHANGE_SEQ.NEXTVAL), 4, '0') FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1"; q != mau {
		t.Errorf("ID baru\ndapat %s\nmau   %s", q, mau)
	}
}

// Sisip menulis kolom PEGA_TREATYEXCHANGE tanpa QURRENCYID; Ubah lewat ROWID dan tidak menyentuh ID, QURRENCYID, DATEIN.
func TestSqlSisipUbah(t *testing.T) {
	s := satuBaris(SqlSisip("S.TREATYEXCHANGEYEARLY"))
	if !strings.Contains(s, "(ID, TREATYYEAR, STARTDATE, ENDDATE, TOIDR, TOUSD, USERID, DATEIU, IDCURRENCY, CURRENCY, QUARTER, DATEIN) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12)") ||
		strings.Contains(s, "QURRENCYID") {
		t.Errorf("sisip %s", s)
	}
	u := satuBaris(SqlUbah("S.TREATYEXCHANGEYEARLY"))
	set := strings.SplitN(u, "WHERE", 2)[0]
	for _, k := range []string{"QURRENCYID", "DATEIN", " ID ="} {
		if strings.Contains(set, k) {
			t.Errorf("ubah menyentuh %s: %s", k, u)
		}
	}
	if !strings.HasSuffix(u, "WHERE ROWID = CHARTOROWID(:11)") {
		t.Errorf("ubah harus lewat ROWID: %s", u)
	}
}

// Kembar: tahun + mata uang + quarter, baris sendiri (ROWID) dikecualikan; daftar bergabung nama mata uang.
func TestSqlBaca(t *testing.T) {
	k := satuBaris(SqlKembar("S.TREATYEXCHANGEYEARLY"))
	if !strings.Contains(k, "TREATYYEAR = :1 AND IDCURRENCY = :2 AND QUARTER = :3 AND (:4 IS NULL OR ROWID <> CHARTOROWID(:5))") {
		t.Errorf("kembar %s", k)
	}
	d := satuBaris(SqlDaftar("S.TREATYEXCHANGEYEARLY t LEFT JOIN S.CURRENCY c ON c.ID = t.IDCURRENCY"))
	for _, b := range []string{"SELECT ROWIDTOCHAR(t.ROWID), t.ID,", "(:5 IS NULL OR t.TREATYYEAR = :6)",
		"ORDER BY t.TREATYYEAR DESC, t.CURRENCY, t.QUARTER, t.ID"} {
		if !strings.Contains(d, b) {
			t.Errorf("daftar tanpa %q: %s", b, d)
		}
	}
	for _, tb := range DaftarTabelDitulis {
		if tb != Tabel {
			t.Errorf("modul ini hanya menulis %s, bukan %s", Tabel, tb)
		}
	}
}
