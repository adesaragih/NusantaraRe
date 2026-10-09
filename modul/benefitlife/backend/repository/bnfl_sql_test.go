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

// Grid BrowseBenefitLife_RD: kolom XML saja (ID, BENEFIT), saring ber-ESCAPE, ID ANGKA menurun bawaan (b4391), 10 lewat
// bind; tulis = INSERT / UPDATE kolom bernama; nol DELETE.
func TestSql(t *testing.T) {
	memuat(t, "daftar", SqlDaftar("S.BENEFIT_LIFE", false), "SELECT ID, BENEFIT FROM S.BENEFIT_LIFE",
		`(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(BENEFIT) LIKE :4 ESCAPE '\')`,
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) DESC NULLS LAST, ID DESC", "OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY")
	memuat(t, "naik", SqlDaftar("V", true), "ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET")
	memuat(t, "jumlah", SqlJumlah("S.BENEFIT_LIFE"), "SELECT COUNT(*) FROM S.BENEFIT_LIFE WHERE (:1 IS NULL")
	memuat(t, "ambil", SqlAmbil("S.BENEFIT_LIFE"), "SELECT ID, BENEFIT FROM S.BENEFIT_LIFE WHERE ID = :1")
	memuat(t, "nomor", SqlNomorBaru("S.M_BENEFIT_LIFE_SEQ"), "SELECT TO_CHAR(S.M_BENEFIT_LIFE_SEQ.NEXTVAL) FROM DUAL")
	memuat(t, "ada", SqlAdaID("S.BENEFIT_LIFE"), "SELECT COUNT(*) FROM S.BENEFIT_LIFE WHERE ID = :1")
	memuat(t, "sisip", SqlSisip("S.BENEFIT_LIFE"), "INSERT INTO S.BENEFIT_LIFE (ID, BENEFIT) VALUES (:1, :2)")
	memuat(t, "ubah", SqlUbah("S.BENEFIT_LIFE"), "UPDATE S.BENEFIT_LIFE SET BENEFIT = :1 WHERE ID = :2")
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` || PolaCari("  ") != nil {
		t.Errorf("pola %v", p)
	}
}

// Lapis penjaga: tulis hanya ke BENEFIT_LIFE; objek lain (mis. M_BENEFIT, milik aplikasi lain) ditolak.
func TestPeriksaTulis(t *testing.T) {
	if !slices.Equal(DaftarTabelDitulis, []string{"BENEFIT_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	for _, objek := range []string{"M_BENEFIT", "M_BENEFIT_LIFE", "T_BENEFIT", "M_PLAN_BENEFIT"} {
		if err := PeriksaTulis(objek, "UPDATE X SET A = 1"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
	if err := PeriksaTulis("DUAL", "SELECT 1 FROM DUAL"); err != nil {
		t.Errorf("SELECT: %v", err)
	}
}

// ORA-00001 (PK SYS_C009031) = ErrKembar; tabel / sequence tidak ada = ErrBelumAda.
func TestBungkus(t *testing.T) {
	if err := Bungkus(errors.New("ORA-00001: unique constraint violated"), "x"); !errors.Is(err, ErrKembar) {
		t.Errorf("ORA-00001 %v", err)
	}
	for _, kode := range []string{"ORA-00942", "ORA-00904", "ORA-02289"} {
		if err := Bungkus(errors.New(kode+": x"), "x"); !errors.Is(err, ErrBelumAda) {
			t.Errorf("%s %v", kode, err)
		}
	}
}

// K1: nol JSON dan nol nama lama M_BENEFIT_LIFE (kecuali sequence) di baris kode produksi repository; kolom = kolom
// view lama.
func TestNolJSONDanNamaLama(t *testing.T) {
	if !slices.Equal(KolomTabel, []string{"ID", "BENEFIT"}) {
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
			if strings.Contains(kode, "JSON") || strings.Contains(kode, `"M_BENEFIT_LIFE"`) || strings.Contains(kode, "DELETE") {
				t.Errorf("%s:%d memakai JSON / nama lama / DELETE: %s", b, i+1, kode)
			}
		}
	}
}
