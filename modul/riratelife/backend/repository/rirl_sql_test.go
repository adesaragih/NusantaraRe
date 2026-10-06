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

// Grid dibaca dari VIEW RATE_LIFE_SUMMARY: kolom XML saja, saring ID / nama ber-ESCAPE, 50 per halaman lewat bind.
func TestSqlDaftar(t *testing.T) {
	d := SqlDaftar("S.RATE_LIFE_SUMMARY", "", true)
	memuat(t, "daftar", d, "SELECT ID, USEDBY, OPERATORID, MODIFIEDDATE FROM S.RATE_LIFE_SUMMARY",
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

// RALAT R4: ringkasan ditulis ke TABEL FLAT RATE_LIFE_SUMMARY - kolom bernama, nol JSONDATA; TYPE/FLAG tidak ditulis.
// ID terpakai diperiksa di tabel flat DAN JSON warisan M_RATE_LIFE_SUMMARY.
func TestSqlTulisRingkasan(t *testing.T) {
	memuat(t, "sisip", SqlSisipRingkasan("S.RATE_LIFE_SUMMARY"),
		"INSERT INTO S.RATE_LIFE_SUMMARY (ID, USEDBY, OPERATORID, MODIFIEDDATE) VALUES (:1, :2, :3, :4)")
	memuat(t, "ubah", SqlUbahRingkasan("S.RATE_LIFE_SUMMARY"),
		"UPDATE S.RATE_LIFE_SUMMARY SET USEDBY = :1, OPERATORID = :2, MODIFIEDDATE = :3 WHERE ID = :4")
	memuat(t, "hapus", SqlHapusRingkasan("S.RATE_LIFE_SUMMARY"), "DELETE FROM S.RATE_LIFE_SUMMARY WHERE ID = :1")
	for _, q := range []string{SqlSisipRingkasan("T"), SqlUbahRingkasan("T"), SqlDaftar("T", "", false), SqlAmbil("T")} {
		if strings.Contains(q, "JSON") || strings.Contains(q, "TYPE") || strings.Contains(q, "FLAG") {
			t.Errorf("ringkasan flat memakai JSON/TYPE/FLAG: %s", q)
		}
	}
	memuat(t, "maks ringkasan", SqlMaksIDRingkasan("S.RATE_LIFE_SUMMARY", "S.M_RATE_LIFE_SUMMARY"),
		"NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0)) FROM (SELECT ID FROM S.RATE_LIFE_SUMMARY UNION ALL SELECT ID FROM S.M_RATE_LIFE_SUMMARY)")
	memuat(t, "ada ringkasan", SqlAdaIDRingkasan("S.RATE_LIFE_SUMMARY", "S.M_RATE_LIFE_SUMMARY"),
		"SELECT COUNT(*) FROM (SELECT ID FROM S.RATE_LIFE_SUMMARY WHERE ID = :1 UNION ALL SELECT ID FROM S.M_RATE_LIFE_SUMMARY WHERE ID = :2)")
	memuat(t, "baca json rate", SqlBacaJSON("S.M_RATE_LIFE"), "SELECT JSONDATA FROM S.M_RATE_LIFE WHERE ID = :1 FOR UPDATE")
	memuat(t, "tulis json", SqlTulisJSON("S.M_RATE_LIFE"), "UPDATE S.M_RATE_LIFE SET JSONDATA = :1 WHERE ID = :2")
	memuat(t, "rate milik", SqlIDRateMilik("S.RATE_LIFE"), "SELECT ID FROM S.RATE_LIFE WHERE IDUSEDBY = :1")
	memuat(t, "hapus rate", SqlHapusRate("S.M_RATE_LIFE", "S.RATE_LIFE"),
		"DELETE FROM S.M_RATE_LIFE WHERE ID IN (SELECT ID FROM S.RATE_LIFE WHERE IDUSEDBY = :1)")
	memuat(t, "id baru", SqlIDBaru("S.SEQ_M_RATE_LIFE_SUMMARY"), "SELECT TO_CHAR(S.SEQ_M_RATE_LIFE_SUMMARY.NEXTVAL) FROM DUAL")
	memuat(t, "maks", SqlMaksID("S.M_RATE_LIFE"), "NVL(MAX(TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$'))), 0)")
}

// Rate: TYPE tidak pernah ditulis; CONTRACT kosong = kunci tidak ditulis (ABSENT ON NULL); kunci kembar sekali baca.
func TestSqlRate(t *testing.T) {
	s := satuBaris(SqlSisipRate("S.M_RATE_LIFE"))
	memuat(t, "sisip rate", s, "JSON_OBJECT('IDUSEDBY' VALUE :2, 'USEDBY' VALUE :3, 'GENDER' VALUE :4, 'CONTRACT' VALUE :5, 'AGE' VALUE :6, 'RATE' VALUE :7 ABSENT ON NULL)")
	if strings.Contains(s, "'TYPE'") {
		t.Error("TYPE tidak boleh ditulis")
	}
	memuat(t, "detail", SqlDaftarRate("S.RATE_LIFE"), "SELECT ID, IDUSEDBY, USEDBY, GENDER, CONTRACT, AGE, RATE FROM S.RATE_LIFE WHERE IDUSEDBY = :1",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) DESC NULLS LAST, ID DESC OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY")
	memuat(t, "kunci", SqlRateDari("S.RATE_LIFE", 3), "WHERE IDUSEDBY IN (:1, :2, :3)")
}

// Lapis penjaga: tulis hanya ke tabel flat RATE_LIFE_SUMMARY dan M_RATE_LIFE; JSON ringkasan warisan dan view rate
// dibaca saja (K-F1).
func TestPeriksaTulis(t *testing.T) {
	for _, objek := range DaftarDibacaSaja {
		if err := PeriksaTulis(objek, "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
	if err := PeriksaTulis(TabelRingkasanJSON, SqlTulisJSON("X.M_RATE_LIFE_SUMMARY")); !errors.Is(err, ErrBacaSaja) {
		t.Errorf("M_RATE_LIFE_SUMMARY ditulis: %v", err)
	}
	if err := PeriksaTulis(TabelRate, SqlTulisJSON("X.M_RATE_LIFE")); err != nil {
		t.Errorf("ubah rate: %v", err)
	}
	if err := PeriksaTulis(ViewRate, SqlDaftarRate("V")); err != nil {
		t.Errorf("SELECT view: %v", err)
	}
	for _, objek := range DaftarTabelDitulis {
		if err := PeriksaTulis(objek, "DELETE FROM X"); err != nil {
			t.Errorf("%s: %v", objek, err)
		}
	}
	if len(DaftarTabelDitulis) != 2 || DaftarTabelDitulis[0] != "RATE_LIFE_SUMMARY" || DaftarTabelDitulis[1] != "M_RATE_LIFE" {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	if len(DaftarDibacaSaja) != 2 || DaftarDibacaSaja[0] != "RATE_LIFE" || DaftarDibacaSaja[1] != "M_RATE_LIFE_SUMMARY" {
		t.Errorf("dibaca saja %v", DaftarDibacaSaja)
	}
}
