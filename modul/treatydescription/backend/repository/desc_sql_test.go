package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

// ID baru PERSIS seperti PEGA_TREATYDESC: '1' || LPAD(TREATY_DESCRIPTION_SEQ.NEXTVAL, 4, '0'); awalan diikat.
func TestSqlIDBaru(t *testing.T) {
	if q, mau := SqlIDBaru("S.TREATY_DESCRIPTION_SEQ"), "SELECT :1 || LPAD(TO_CHAR(S.TREATY_DESCRIPTION_SEQ.NEXTVAL), 4, '0') FROM DUAL"; q != mau {
		t.Errorf("ID baru\ndapat %s\nmau   %s", q, mau)
	}
	if AwalanID != "1" || Sequence != "TREATY_DESCRIPTION_SEQ" {
		t.Errorf("awalan %q sequence %q", AwalanID, Sequence)
	}
}

// Sisip dan Ubah menulis keempat kolom PEGA_TREATYDESC; Ubah tidak pernah menyentuh ID.
func TestSqlSisipUbah(t *testing.T) {
	if s := satuBaris(SqlSisip("S.TREATYDESC")); s != "INSERT INTO S.TREATYDESC (ID, DESCNAME, ISXOL, STATUSAKTIF) VALUES (:1, :2, :3, :4)" {
		t.Errorf("sisip %s", s)
	}
	u := satuBaris(SqlUbah("S.TREATYDESC"))
	if u != "UPDATE S.TREATYDESC SET DESCNAME = :1, ISXOL = :2, STATUSAKTIF = :3 WHERE ID = :4" {
		t.Errorf("ubah %s", u)
	}
}

// Daftar dan Ambil HANYA membaca TREATYDESC (perintah work owner 05-10-2026: "hanya baca dari tabel yang saya kasih");
// saringan jenis (NULL = Non XOL) dan status (NULL = aktif), urut ID.
func TestSqlBaca(t *testing.T) {
	d := satuBaris(SqlDaftar("S.TREATYDESC"))
	for _, b := range []string{
		"SELECT ID, DESCNAME, NVL(ISXOL, '0'), STATUSAKTIF FROM S.TREATYDESC WHERE",
		"(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\\' OR UPPER(DESCNAME) LIKE :3 ESCAPE '\\')",
		"(:4 IS NULL OR NVL(ISXOL, '0') = :5)",
		"(:6 IS NULL OR DECODE(STATUSAKTIF, '0', '0', '1') = :7)",
		"ORDER BY ID",
	} {
		if !strings.Contains(d, b) {
			t.Errorf("daftar tanpa %q:\n%s", b, d)
		}
	}
	a := satuBaris(SqlAmbil("S.TREATYDESC"))
	if a != "SELECT ID, DESCNAME, NVL(ISXOL, '0'), STATUSAKTIF FROM S.TREATYDESC WHERE ID = :1" {
		t.Errorf("ambil %s", a)
	}
	for _, q := range []string{d, a, SqlPemakaiNama("S.TREATYDESC")} {
		if strings.Contains(q, "JOIN") || strings.Count(q, "FROM ") != 1 {
			t.Errorf("membaca tabel lain: %s", q)
		}
	}
	if n := satuBaris(SqlPemakaiNama("S.TREATYDESC")); !strings.Contains(n, "UPPER(TRIM(DESCNAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))") {
		t.Errorf("pemakai nama %s", n)
	}
	if len(DaftarTabelDitulis) != 1 || DaftarTabelDitulis[0] != Tabel || len(DaftarTabelDibacaSaja) != 0 {
		t.Errorf("modul ini hanya memakai %s: tulis %v, baca %v", Tabel, DaftarTabelDitulis, DaftarTabelDibacaSaja)
	}
}

// Setiap SQL lolos PeriksaSQL (nol COMMIT/ROLLBACK).
func TestSeluruhSQLLolosPeriksa(t *testing.T) {
	for _, q := range []string{SqlDaftar("S.T"), SqlAmbil("S.T"), SqlPemakaiNama("S.T"), SqlIDBaru("S.Q"),
		SqlSisip("S.T"), SqlUbah("S.T")} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
	}
}
