package repository

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"nusantarare/modul/ricommlife/backend/models"
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

// Grid ringkasan dibaca dari VIEW RICOMM_LIFE_SUMMARY: kolom XML saja, saring ber-ESCAPE, ID menaik (b9857), 50 lewat bind.
func TestSqlRingkasan(t *testing.T) {
	memuat(t, "daftar", SqlDaftar("S.RICOMM_LIFE_SUMMARY", "", false), "SELECT ID, USEDBY, OPERATORID, MODIFIEDDATE FROM S.RICOMM_LIFE_SUMMARY",
		`(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(USEDBY) LIKE :4 ESCAPE '\')`,
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC", "OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY")
	memuat(t, "urut nama", SqlDaftar("V", "usedby", true), "ORDER BY UPPER(USEDBY) DESC NULLS LAST, ID OFFSET")
	if strings.Contains(SqlDaftar("V", "ID; DROP TABLE X", false), "DROP") {
		t.Error("urut dari masukan masuk ke SQL")
	}
	memuat(t, "kembar", SqlPemakaiNama("V"), "WHERE UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))")
	memuat(t, "sisip", SqlSisipRingkasan("S.M_RICOMM_LIFE_SUMMARY"),
		"INSERT INTO S.M_RICOMM_LIFE_SUMMARY (ID, JSONDATA) VALUES (:1, JSON_OBJECT('USEDBY' VALUE :2, 'OPERATORID' VALUE :3, 'MODIFIEDDATE' VALUE :4, 'pxObjClass' VALUE :5 ABSENT ON NULL))")
	memuat(t, "baca json", SqlBacaJSON("S.M_RICOMM_LIFE_SUMMARY"), "SELECT JSONDATA FROM S.M_RICOMM_LIFE_SUMMARY WHERE ID = :1 FOR UPDATE")
	memuat(t, "tulis json", SqlTulisJSON("S.M_RICOMM_LIFE_SUMMARY"), "UPDATE S.M_RICOMM_LIFE_SUMMARY SET JSONDATA = :1 WHERE ID = :2")
	memuat(t, "situs", SqlSitus("S.M_SITE_DATABASE"), "SELECT TO_CHAR(ID) FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1")
	memuat(t, "nomor", SqlNomorBaru("S.M_RICOMM_LIFE_SEQ"), "SELECT TO_CHAR(S.M_RICOMM_LIFE_SEQ.NEXTVAL) FROM DUAL")
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` || PolaCari("  ") != nil {
		t.Errorf("pola %v", p)
	}
}

// Rincian = tabel FLAT RICOMM_LIFE: kolom bernama (bukan JSONDATA), angka tanpa NLS, ID menaik (b9515).
func TestSqlKomisiFlat(t *testing.T) {
	baca := "SELECT ID, IDUSEDBY, USEDBY, TO_CHAR(CONTRACT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), " +
		"TO_CHAR(YEAR, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), TO_CHAR(COMM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') FROM S.RICOMM_LIFE"
	memuat(t, "daftar", SqlDaftarKomisi("S.RICOMM_LIFE"), baca, "WHERE IDUSEDBY = :1",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY")
	memuat(t, "kunci", SqlKomisiDari("S.RICOMM_LIFE", 3), "WHERE IDUSEDBY IN (:1, :2, :3)")
	memuat(t, "sisip", SqlSisipKomisi("S.RICOMM_LIFE"),
		"INSERT INTO S.RICOMM_LIFE (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM) VALUES (:1, :2, :3, TO_NUMBER(:4), TO_NUMBER(:5), TO_NUMBER(:6) / POWER(10, :7))")
	memuat(t, "ubah", SqlUbahKomisi("S.RICOMM_LIFE"),
		"UPDATE S.RICOMM_LIFE SET CONTRACT = TO_NUMBER(:1), YEAR = TO_NUMBER(:2), COMM = TO_NUMBER(:3) / POWER(10, :4) WHERE ID = :5 AND IDUSEDBY = :6")
	memuat(t, "nama", SqlUbahNamaKomisi("S.RICOMM_LIFE"), "UPDATE S.RICOMM_LIFE SET USEDBY = :1 WHERE IDUSEDBY = :2")
	memuat(t, "hapus", SqlHapusKomisi("S.RICOMM_LIFE"), "DELETE FROM S.RICOMM_LIFE WHERE IDUSEDBY = :1")
	for _, q := range []string{SqlDaftarKomisi("T"), SqlSisipKomisi("T"), SqlUbahKomisi("T"), SqlKomisiDari("T", 1)} {
		if strings.Contains(q, "JSON") {
			t.Errorf("rincian flat memakai JSON: %s", q)
		}
	}
	for kanonik, mau := range map[string]struct {
		koef  any
		skala int64
	}{"12.05": {"1205", 2}, "0.5": {"5", 1}, "7": {"7", 0}, "0": {"0", 0}, "": {nil, 0}, "100": {"100", 0}} {
		if koef, skala := PecahDesimal(kanonik); koef != mau.koef || skala != mau.skala {
			t.Errorf("PecahDesimal(%q) = %v,%d mau %v,%d", kanonik, koef, skala, mau.koef, mau.skala)
		}
	}
	for masuk, mau := range map[string]string{".5": "0.5", "12.50": "12.5", "7": "7", "": ""} {
		if got := AngkaOracle(masuk); got != mau {
			t.Errorf("AngkaOracle(%q) = %q mau %q", masuk, got, mau)
		}
	}
}

// Lapis penjaga: tulis hanya ke M_RICOMM_LIFE_SUMMARY dan RICOMM_LIFE; M_RICOMM_LIFE, view, situs dibaca saja.
func TestPeriksaTulis(t *testing.T) {
	for _, objek := range DaftarDibacaSaja {
		if err := PeriksaTulis(objek, "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
		if err := PeriksaTulis(objek, "SELECT 1 FROM X"); err != nil {
			t.Errorf("%s SELECT: %v", objek, err)
		}
	}
	if !slices.Equal(DaftarTabelDitulis, []string{"M_RICOMM_LIFE_SUMMARY", "RICOMM_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	if err := PeriksaTulis(TabelKomisi, SqlKunciKomisi("S.RICOMM_LIFE")); err != nil {
		t.Errorf("kunci tabel flat: %v", err)
	}
	if !slices.Contains(DaftarDibacaSaja, "M_RICOMM_LIFE") {
		t.Error("M_RICOMM_LIFE harus dibaca saja (butir 3)")
	}
}

// Oracle DEV menolak JSON_MERGEPATCH (ORA-00907, RALAT R3 riratelife): nol pemakaian di SQL modul ini.
func TestNolJSONMergepatchDiSQL(t *testing.T) {
	berkas, _ := filepath.Glob("*.go")
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		for i, baris := range strings.Split(string(isi), "\n") {
			if strings.Contains(baris, "JSON_MERGEPATCH") && !strings.HasPrefix(strings.TrimSpace(baris), "//") {
				t.Errorf("%s:%d memakai JSON_MERGEPATCH", b, i+1)
			}
		}
	}
}

// Butir 6: kunci lama dan pxObjClass dipertahankan; kosong = dibuang.
func TestTerapkanKunci(t *testing.T) {
	lama := `{"MODIFIEDDATE":"20181205T073755.559 GMT","OPERATORID":"UJI-LAMA","pxObjClass":"ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY","USEDBY":"UJI <A>&B","N":1.50}`
	baru, err := TerapkanKunci(lama, map[string]string{"USEDBY": "UJI C", "OPERATORID": "UJI-BARU", "X": ""})
	if err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{`"USEDBY":"UJI C"`, `"OPERATORID":"UJI-BARU"`, `"pxObjClass":"` + models.KelasRingkasan + `"`,
		`"MODIFIEDDATE":"20181205T073755.559 GMT"`, `"N":1.50`} {
		if !strings.Contains(baru, mau) {
			t.Errorf("tanpa %s: %s", mau, baru)
		}
	}
	for _, rusak := range []string{"[1]", "bukan json", "null"} {
		if _, err := TerapkanKunci(rusak, map[string]string{"A": "1"}); !errors.Is(err, ErrJSONRusak) {
			t.Errorf("%q: %v", rusak, err)
		}
	}
	if v, ada := TeksKunci(`{"COMM":12.5,"YEAR":" 1 ","X":null}`, "COMM"); v != "12.5" || !ada {
		t.Errorf("TeksKunci angka %q %v", v, ada)
	}
	if v, _ := TeksKunci(`{"YEAR":" 1 "}`, "YEAR"); v != "1" {
		t.Errorf("TeksKunci teks %q", v)
	}
	if _, ada := TeksKunci(`{"X":null}`, "X"); ada {
		t.Error("null harus tidak ada")
	}
}

// Kolom tabel flat = kolom VIEW warisan RICOMM_LIFE (butir 1), kolom view ringkasan = definisinya.
func TestKolomCocokDenganView(t *testing.T) {
	if !slices.Equal(KolomViewKomisiLama, []string{"ID", "IDUSEDBY", "USEDBY", "CONTRACT", "YEAR", "COMM"}) {
		t.Errorf("kolom view lama %v", KolomViewKomisiLama)
	}
	if !slices.Equal(KolomViewRingkasan, []string{"ID", "USEDBY", "MODIFIEDDATE", "OPERATORID"}) {
		t.Errorf("kolom view ringkasan %v", KolomViewRingkasan)
	}
	for _, k := range strings.Split(kolomRingkasan, ", ") {
		if !slices.Contains(KolomViewRingkasan, k) {
			t.Errorf("kolom ringkasan %s tidak ada di view", k)
		}
	}
}

// Butir 4: angka dari Oracle diurai di Go tanpa bergantung NLS sesi - titik ATAU koma desimal, `TM9` tanpa nol depan.
func TestAngkaOracleTanpaNLS(t *testing.T) {
	for masuk, mau := range map[string]string{
		".5": "0.5", ",5": "0.5", "12.50": "12.5", "12,50": "12.5", "7": "7", "2026": "2026", " 0 ": "0",
		"123456789012345678901234567890.12345678": "123456789012345678901234567890.12345678",
		"123456789012345678901234567890,12345678": "123456789012345678901234567890.12345678",
		"-1.5": "-1.5", "-,5": "-0.5", "": "", "1.234,5": "1.234,5", "abc": "abc",
	} {
		if got := AngkaOracle(masuk); got != mau {
			t.Errorf("AngkaOracle(%q) = %q mau %q", masuk, got, mau)
		}
	}
	if err := PeriksaTulis("SESI", SqlSesiNLS); err != nil {
		t.Errorf("setelan sesi pindah ditolak penjaga: %v", err)
	}
	if err := PeriksaTulis("SESI", "ALTER SESSION SET NLS_DATE_FORMAT = 'YYYY'"); !errors.Is(err, ErrBacaSaja) {
		t.Errorf("pernyataan sesi lain harus ditolak: %v", err)
	}
}
