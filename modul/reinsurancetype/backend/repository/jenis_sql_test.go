package repository

import (
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// ID baru seperti PEGA_REINSURANCETYPE: awalan diikat + M_REINSURANCETYPE_SEQ 4 digit.
func TestSqlIDBaru(t *testing.T) {
	if q := SqlIDBaru("S.M_REINSURANCETYPE_SEQ"); q != "SELECT :1 || LPAD(TO_CHAR(S.M_REINSURANCETYPE_SEQ.NEXTVAL), 4, '0') FROM DUAL" {
		t.Errorf("ID baru %s", q)
	}
}

// Sisip / Ubah: kolom PEGA_REINSURANCETYPE + GROUPTYPE; Ubah tidak menyentuh ID.
func TestSqlSisipUbah(t *testing.T) {
	if s := satuBaris(SqlSisip("S.REINSURANCETYPE")); !strings.Contains(s, "(ID, NOTE, TYPE, SOANOTE, CODE, FLAG, USERID, NOURUT, TGLUPDATE, GROUPTYPE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)") {
		t.Errorf("sisip %s", s)
	}
	u := satuBaris(SqlUbah("S.REINSURANCETYPE"))
	if strings.Contains(strings.SplitN(u, "WHERE", 2)[0], " ID =") || !strings.HasSuffix(u, "WHERE ID = :10") {
		t.Errorf("ubah %s", u)
	}
}

// Daftar: cari ber-ESCAPE, saring Type dan Flag lewat bind; kembar tanpa beda huruf, baris sendiri dikecualikan.
func TestSqlBaca(t *testing.T) {
	d := satuBaris(SqlDaftar("S.REINSURANCETYPE"))
	for _, b := range []string{"(:6 IS NULL OR TYPE = :7)", "(:8 IS NULL OR FLAG = :9)", "UPPER(SOANOTE) LIKE :4 ESCAPE '\\'",
		"ORDER BY TYPE, UPPER(NOTE), ID"} {
		if !strings.Contains(d, b) {
			t.Errorf("daftar tanpa %q: %s", b, d)
		}
	}
	if k := satuBaris(SqlPemakaiNama("S.REINSURANCETYPE")); !strings.Contains(k, "WHERE UPPER(TRIM(NOTE)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))") {
		t.Errorf("kembar %s", k)
	}
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` {
		t.Errorf("pola %v", p)
	}
}
