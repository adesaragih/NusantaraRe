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

// Grid BrowseCauseofLossLife_RD: kolom XML saja (ID, CAUSEOFLOSS), tanpa saring (b4036), ID ANGKA menaik tetap
// (b3923 / b3929), 10 lewat bind; nama kembar tanpa beda huruf; tulis = INSERT / UPDATE kolom bernama; nol DELETE.
func TestSql(t *testing.T) {
	memuat(t, "daftar", SqlDaftar("S.CAUSEOFLOSS_LIFE"), "SELECT ID, CAUSEOFLOSS FROM S.CAUSEOFLOSS_LIFE ORDER BY",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY")
	if strings.Contains(SqlDaftar("S.CAUSEOFLOSS_LIFE"), "WHERE") || strings.Contains(SqlJumlah("T"), "WHERE") {
		t.Error("grid tanpa saring (pyGridFiltering false b4036)")
	}
	memuat(t, "jumlah", SqlJumlah("S.CAUSEOFLOSS_LIFE"), "SELECT COUNT(*) FROM S.CAUSEOFLOSS_LIFE")
	memuat(t, "ambil", SqlAmbil("S.CAUSEOFLOSS_LIFE"), "SELECT ID, CAUSEOFLOSS FROM S.CAUSEOFLOSS_LIFE WHERE ID = :1")
	memuat(t, "pemakai", SqlPemakaiNama("S.CAUSEOFLOSS_LIFE"),
		"WHERE UPPER(TRIM(CAUSEOFLOSS)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID")
	memuat(t, "nomor", SqlNomorBaru("S.M_CAUSEOFLOSS_LIFE_SEQ"), "SELECT TO_CHAR(S.M_CAUSEOFLOSS_LIFE_SEQ.NEXTVAL) FROM DUAL")
	memuat(t, "ada", SqlAdaID("S.CAUSEOFLOSS_LIFE"), "SELECT COUNT(*) FROM S.CAUSEOFLOSS_LIFE WHERE ID = :1")
	memuat(t, "sisip", SqlSisip("S.CAUSEOFLOSS_LIFE"), "INSERT INTO S.CAUSEOFLOSS_LIFE (ID, CAUSEOFLOSS) VALUES (:1, :2)")
	memuat(t, "ubah", SqlUbah("S.CAUSEOFLOSS_LIFE"), "UPDATE S.CAUSEOFLOSS_LIFE SET CAUSEOFLOSS = :1 WHERE ID = :2")
}

// Lapis penjaga: tulis hanya ke CAUSEOFLOSS_LIFE; objek milik aplikasi lain (prompt bab 1) ditolak.
func TestPeriksaTulis(t *testing.T) {
	if !slices.Equal(DaftarTabelDitulis, []string{"CAUSEOFLOSS_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	for _, objek := range []string{"M_CAUSE_OF_LOSS", "D_CAUSE_OF_LOSS", "V_M_CAUSE_OF_LOSS", "T_LISTCAUSEOFLOSS", "M_CAUSEOFLOSS_LIFE"} {
		if err := PeriksaTulis(objek, "UPDATE X SET A = 1"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
	if err := PeriksaTulis("DUAL", "SELECT 1 FROM DUAL"); err != nil {
		t.Errorf("SELECT: %v", err)
	}
}

// ORA-00001 (PK SYS_C008825) = ErrKembar; tabel / sequence tidak ada = ErrBelumAda.
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

// K1: nol JSON, nol nama lama M_CAUSEOFLOSS_LIFE (kecuali sequence), nol DELETE di baris kode produksi repository;
// kolom = kolom view lama.
func TestNolJSONDanNamaLama(t *testing.T) {
	if !slices.Equal(KolomTabel, []string{"ID", "CAUSEOFLOSS"}) {
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
			if strings.Contains(kode, "JSON") || strings.Contains(kode, `"M_CAUSEOFLOSS_LIFE"`) || strings.Contains(kode, "DELETE") {
				t.Errorf("%s:%d memakai JSON / nama lama / DELETE: %s", b, i+1, kode)
			}
		}
	}
}
