package repository

import (
	"errors"
	"strings"
	"testing"
)

func satuBaris(q string) string { return strings.Join(strings.Fields(q), " ") }

func memuat(t *testing.T, nama, q string, bagian ...string) {
	t.Helper()
	s := satuBaris(q)
	for _, b := range bagian {
		if !strings.Contains(s, b) {
			t.Errorf("%s tanpa %q:\n%s", nama, b, s)
		}
	}
}

// Grid dibaca dari kolom M_RATE_LIFE_SUMMARY (RALAT R6): kolom XML saja, saring ID / nama ber-ESCAPE, 50 per halaman.
func TestSqlDaftar(t *testing.T) {
	d := SqlDaftar("S.M_RATE_LIFE_SUMMARY", "", true)
	memuat(t, "daftar", d, "SELECT ID, USEDBY, OPERATORID, MODIFIEDDATE FROM S.M_RATE_LIFE_SUMMARY",
		`(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(USEDBY) LIKE :4 ESCAPE '\')`,
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) DESC NULLS LAST, ID DESC", "OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY")
	memuat(t, "urut nama", SqlDaftar("V", "usedby", false), "ORDER BY UPPER(USEDBY) ASC NULLS LAST, ID OFFSET")
	memuat(t, "urut tanggal", SqlDaftar("V", "modifieddate", true), "ORDER BY MODIFIEDDATE DESC NULLS LAST, ID")
	memuat(t, "urut asing", SqlDaftar("V", "ID; DROP TABLE X", false), "ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC")
	if strings.Contains(SqlDaftar("V", "ID; DROP TABLE X", false), "DROP") {
		t.Error("urut dari masukan masuk ke SQL")
	}
	memuat(t, "jumlah", SqlJumlah("V"), "SELECT COUNT(*) FROM V WHERE (:1 IS NULL")
	memuat(t, "kembar", SqlPemakaiNama("V"), "WHERE UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))")
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` {
		t.Errorf("pola %v", p)
	}
	if PolaCari("  ") != nil {
		t.Error("kosong harus nil")
	}
}

// RALAT R6: ringkasan ditulis ke kolom M_RATE_LIFE_SUMMARY - satu tabel, nol JSONDATA, TYPE tidak ditulis, nol FLAG;
// ID terpakai diperiksa di tabel itu saja.
func TestSqlTulisRingkasan(t *testing.T) {
	memuat(t, "sisip", SqlSisipRingkasan("S.M_RATE_LIFE_SUMMARY"),
		"INSERT INTO S.M_RATE_LIFE_SUMMARY (ID, USEDBY, OPERATORID, MODIFIEDDATE) VALUES (:1, :2, :3, :4)")
	memuat(t, "ubah", SqlUbahRingkasan("S.M_RATE_LIFE_SUMMARY"),
		"UPDATE S.M_RATE_LIFE_SUMMARY SET USEDBY = :1, OPERATORID = :2, MODIFIEDDATE = :3 WHERE ID = :4")
	memuat(t, "hapus", SqlHapusRingkasan("S.M_RATE_LIFE_SUMMARY"), "DELETE FROM S.M_RATE_LIFE_SUMMARY WHERE ID = :1")
	for _, q := range []string{SqlSisipRingkasan("T"), SqlUbahRingkasan("T"), SqlDaftar("T", "", false), SqlAmbil("T"),
		SqlPemakaiNama("T"), SqlJumlah("T"), kolomRingkasan} {
		if strings.Contains(q, "JSON") || strings.Contains(q, "TYPE") || strings.Contains(q, "FLAG") {
			t.Errorf("ringkasan memakai JSON/TYPE/FLAG: %s", q)
		}
	}
	memuat(t, "maks ringkasan", SqlMaksID("S.M_RATE_LIFE_SUMMARY"), "NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0)) FROM S.M_RATE_LIFE_SUMMARY")
	memuat(t, "ada ringkasan", SqlAdaID("S.M_RATE_LIFE_SUMMARY"), "SELECT COUNT(*) FROM S.M_RATE_LIFE_SUMMARY WHERE ID = :1")
	memuat(t, "hapus rate", SqlHapusRate("S.M_RATE_LIFE"), "DELETE FROM S.M_RATE_LIFE WHERE IDUSEDBY = :1")
	memuat(t, "id baru", SqlIDBaru("S.SEQ_M_RATE_LIFE_SUMMARY"), "SELECT TO_CHAR(S.SEQ_M_RATE_LIFE_SUMMARY.NEXTVAL) FROM DUAL")
	memuat(t, "maks", SqlMaksID("S.M_RATE_LIFE"), "NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0)")
}

// Rate (RALAT R7, 10-10-2026): tabel flat M_RATE_LIFE - INSERT/UPDATE kolom biasa, nol JSONDATA / JSON_OBJECT, TYPE tidak
// pernah ditulis; Edit nama = UPDATE USEDBY berkunci IDUSEDBY; kunci kembar sekali baca.
func TestSqlRate(t *testing.T) {
	memuat(t, "sisip rate", SqlSisipRate("S.M_RATE_LIFE"),
		"INSERT INTO S.M_RATE_LIFE (ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE) VALUES (:1, :2, :3, :4, :5, :6, :7)")
	memuat(t, "ubah rate", SqlUbahRate("S.M_RATE_LIFE"),
		"UPDATE S.M_RATE_LIFE SET GENDER = :1, CONTRACT = :2, AGE = :3, RATE = :4 WHERE ID = :5 AND TRIM(IDUSEDBY) = :6")
	memuat(t, "ubah nama", SqlUbahNamaRate("S.M_RATE_LIFE"), "UPDATE S.M_RATE_LIFE SET USEDBY = :1 WHERE IDUSEDBY = :2")
	for _, x := range []string{SqlSisipRate("T"), SqlUbahRate("T"), SqlUbahNamaRate("T"), SqlHapusRate("T"), SqlDaftarRate("T"),
		SqlRateDari("T", 2), SqlJumlahRate("T"), kolomRate} {
		if strings.Contains(strings.ToUpper(x), "JSON") || strings.Contains(x, "TYPE") || strings.Contains(x, "RATE_LIFE") {
			t.Errorf("SQL rate memakai JSON / TYPE / view RATE_LIFE: %s", x)
		}
	}
	memuat(t, "detail", SqlDaftarRate("S.M_RATE_LIFE"), "SELECT ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE FROM S.M_RATE_LIFE WHERE IDUSEDBY = :1",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) DESC NULLS LAST, ID DESC OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY")
	memuat(t, "jumlah", SqlJumlahRate("S.M_RATE_LIFE"), "SELECT COUNT(*) FROM S.M_RATE_LIFE WHERE IDUSEDBY = :1")
	memuat(t, "kunci", SqlRateDari("S.M_RATE_LIFE", 3), "FROM S.M_RATE_LIFE WHERE IDUSEDBY IN (:1, :2, :3)")
}

// Lapis penjaga: tulis hanya ke M_RATE_LIFE_SUMMARY dan M_RATE_LIFE (keduanya juga dibaca); nol objek baca-saja sejak RALAT R7.
func TestPeriksaTulis(t *testing.T) {
	for _, objek := range DaftarDibacaSaja {
		if err := PeriksaTulis(objek, "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
	if err := PeriksaTulis("RATE_LIFE_SUMMARY", "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
		t.Errorf("tabel flat 926 (dibuang 928) ditulis: %v", err)
	}
	for _, x := range []string{SqlSisipRate("X"), SqlUbahRate("X"), SqlUbahNamaRate("X"), SqlHapusRate("X")} {
		if err := PeriksaTulis(TabelRate, x); err != nil {
			t.Errorf("tulis rate: %v", err)
		}
	}
	if err := PeriksaTulis("RATE_LIFE", SqlUbahNamaRate("X")); !errors.Is(err, ErrBacaSaja) {
		t.Errorf("view RATE_LIFE (dibuang) ditulis: %v", err)
	}
	if err := PeriksaTulis("RATE_LIFE", SqlDaftarRate("V")); err != nil {
		t.Errorf("SELECT: %v", err)
	}
	for _, objek := range DaftarTabelDitulis {
		if err := PeriksaTulis(objek, "DELETE FROM X"); err != nil {
			t.Errorf("%s: %v", objek, err)
		}
	}
	if len(DaftarTabelDitulis) != 2 || DaftarTabelDitulis[0] != "M_RATE_LIFE_SUMMARY" || DaftarTabelDitulis[1] != "M_RATE_LIFE" {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	if len(DaftarDibacaSaja) != 0 {
		t.Errorf("dibaca saja %v (view RATE_LIFE dibuang 10-10-2026)", DaftarDibacaSaja)
	}
}
