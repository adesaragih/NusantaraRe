package repository

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
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

// Grid: kolom = view lama (urutan sama), bawaan ID angka menaik, urut Plan Name / Benefit dari daftar putih; tulis =
// INSERT / UPDATE kolom bernama; pilihan dari master (K4).
func TestSql(t *testing.T) {
	memuat(t, "daftar", SqlDaftar("S.PRODUCT_TYPE_LIFE", "", false),
		"SELECT ID, COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID FROM S.PRODUCT_TYPE_LIFE ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY")
	memuat(t, "urut nama", SqlDaftar("V", "covername", true), "ORDER BY UPPER(COVERNAME) DESC NULLS LAST, ID OFFSET")
	memuat(t, "urut benefit", SqlDaftar("V", "benefit", false), "ORDER BY UPPER(BENEFIT) ASC NULLS LAST, ID OFFSET")
	if strings.Contains(SqlDaftar("V", "BUSINESS; DROP", false), "DROP") || strings.Contains(SqlDaftar("V", "business", false), "UPPER(BUSINESS)") {
		t.Error("urut dari masukan / Business (tidak dapat diurutkan b4269) masuk ke SQL")
	}
	memuat(t, "kembar", SqlPemakaiNama("S.PRODUCT_TYPE_LIFE"), "WHERE UPPER(TRIM(COVERNAME)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))")
	memuat(t, "sisip", SqlSisip("S.PRODUCT_TYPE_LIFE"),
		"INSERT INTO S.PRODUCT_TYPE_LIFE (ID, COVERNAME, BUSINESS, BUSINESSID, BENEFIT, BENEFITID) VALUES (:1, :2, :3, :4, :5, :6)")
	memuat(t, "ubah", SqlUbah("S.PRODUCT_TYPE_LIFE"),
		"UPDATE S.PRODUCT_TYPE_LIFE SET COVERNAME = :1, BUSINESS = :2, BUSINESSID = :3, BENEFIT = :4, BENEFITID = :5 WHERE ID = :6")
	memuat(t, "nomor", SqlNomorBaru("S.M_PRODUCT_TYPE_LIFE_SEQ"), "SELECT TO_CHAR(S.M_PRODUCT_TYPE_LIFE_SEQ.NEXTVAL) FROM DUAL")
	memuat(t, "business", SqlPilihanBusiness("S.BUSINESS"), "SELECT TO_CHAR(ID), TO_CHAR(OLDID), NOTE FROM S.BUSINESS WHERE GROUPPANEL = :1 ORDER BY ID")
	memuat(t, "benefit", SqlPilihanBenefit("S.BENEFIT_LIFE"), "SELECT ID, BENEFIT FROM S.BENEFIT_LIFE ORDER BY ID")
}

// Lapis penjaga: tulis hanya PRODUCT_TYPE_LIFE; BUSINESS dan BENEFIT_LIFE dibaca saja (jangan disentuh).
func TestPeriksaTulis(t *testing.T) {
	if !slices.Equal(DaftarTabelDitulis, []string{"PRODUCT_TYPE_LIFE"}) || !slices.Equal(DaftarDibacaSaja, []string{"BUSINESS", "BENEFIT_LIFE"}) {
		t.Errorf("%v %v", DaftarTabelDitulis, DaftarDibacaSaja)
	}
	for _, objek := range []string{"BUSINESS", "M_BUSINESS", "BENEFIT_LIFE", "PLAN_LIFE_SUMMARY"} {
		if err := PeriksaTulis(objek, "UPDATE X SET A = 1"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
}

func TestBungkus(t *testing.T) {
	if err := Bungkus(errors.New("ORA-00001: x"), "x"); !errors.Is(err, ErrKembar) {
		t.Error(err)
	}
	if err := Bungkus(errors.New("ORA-00942: x"), "x"); !errors.Is(err, ErrBelumAda) {
		t.Error(err)
	}
}

// K1: nol JSON, nol nama lama, nol DELETE di kode produksi repository; kolom = view lama.
func TestNolJSONDanNamaLama(t *testing.T) {
	if !slices.Equal(KolomTabel, []string{"ID", "COVERNAME", "BUSINESS", "BUSINESSID", "BENEFIT", "BENEFITID"}) {
		t.Errorf("kolom %v", KolomTabel)
	}
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
			kode := strings.TrimSpace(baris)
			if strings.HasPrefix(kode, "//") {
				continue
			}
			if strings.Contains(kode, "JSON") || strings.Contains(kode, `"M_PRODUCT_TYPE_LIFE"`) || strings.Contains(kode, "DELETE") {
				t.Errorf("%s:%d %s", b, i+1, kode)
			}
		}
	}
}
