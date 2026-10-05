package repository

import (
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// Daftar: cari "memuat" ber-ESCAPE di ID, NAME, NAMEIDN; urut Order No sebagai angka.
func TestSqlDaftar(t *testing.T) {
	q := satuBaris(SqlDaftar("S.TREATYGROUPOJK"))
	for _, bagian := range []string{"SELECT ID, NAME, NAMEIDN, ORDERNO FROM S.TREATYGROUPOJK", "UPPER(NAMEIDN) LIKE :4 ESCAPE '\\'",
		"ORDER BY LPAD(ORDERNO, 10, '0'), ID"} {
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

// ID dan Order No baru: seluruh baris dikunci dulu, lalu nomor tertinggi dari nilai yang seluruhnya angka.
func TestSqlIDBaru(t *testing.T) {
	if q := SqlKunci("S.TREATYGROUPOJK"); q != "SELECT ID FROM S.TREATYGROUPOJK FOR UPDATE" {
		t.Errorf("kunci %s", q)
	}
	if q := SqlNomorTertinggi("S.TREATYGROUPOJK"); q != "SELECT NVL(MAX(TO_NUMBER(ID)), 0) FROM S.TREATYGROUPOJK WHERE REGEXP_LIKE(ID, '^[0-9]+$')" {
		t.Errorf("nomor tertinggi %s", q)
	}
	if q := SqlOrderTertinggi("S.TREATYGROUPOJK"); q != "SELECT NVL(MAX(TO_NUMBER(ORDERNO)), 0) FROM S.TREATYGROUPOJK WHERE REGEXP_LIKE(ORDERNO, '^[0-9]+$')" {
		t.Errorf("order tertinggi %s", q)
	}
}

// Edit tidak menyentuh ID dan ORDERNO.
func TestSqlUbah(t *testing.T) {
	if q := SqlUbah("S.TREATYGROUPOJK"); q != "UPDATE S.TREATYGROUPOJK SET NAME = :1, NAMEIDN = :2 WHERE ID = :3" {
		t.Errorf("ubah %s", q)
	}
}

// Kembar: tanpa beda huruf dan spasi tepi, baris sendiri dikecualikan.
func TestSqlPemakai(t *testing.T) {
	q := satuBaris(SqlPemakai("S.TREATYGROUPOJK", "NAME"))
	if !strings.Contains(q, "WHERE UPPER(TRIM(NAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))") {
		t.Errorf("pemakai %s", q)
	}
}
