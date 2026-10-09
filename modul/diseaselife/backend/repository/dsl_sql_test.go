package repository

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"nusantarare/modul/diseaselife/backend/models"
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

// Grid BrowseDiseaseLife_RD: kolom XML saja (ID, ICD_CODE, DISEASE), saring ICD Code / Disease "memuat" tanpa beda
// huruf (b5148), bawaan ID ANGKA menurun (b5012), urut pilihan ICD Code (b5036), halaman lewat bind; ICD kembar tanpa
// beda huruf; tulis = INSERT / UPDATE kolom bernama.
func TestSql(t *testing.T) {
	const t1 = "S.DISEASE_LIFE"
	memuat(t, "daftar", SqlDaftar(t1, models.UrutID, false), "SELECT ID, ICD_CODE, DISEASE FROM S.DISEASE_LIFE WHERE",
		"(:1 IS NULL OR UPPER(ICD_CODE) LIKE :2 ESCAPE '\\') AND (:3 IS NULL OR UPPER(DISEASE) LIKE :4 ESCAPE '\\')",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) DESC NULLS LAST, ID DESC OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY")
	memuat(t, "daftar naik", SqlDaftar(t1, "", true), "ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET")
	memuat(t, "daftar icd", SqlDaftar(t1, models.UrutICD, true),
		"ORDER BY ICD_CODE ASC NULLS LAST, TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET")
	if strings.Contains(SqlDaftar(t1, "disease", false), "ORDER BY DISEASE") {
		t.Error("kolom Disease tidak dapat diurutkan (pyColumnSorting false b5058)")
	}
	memuat(t, "jumlah", SqlJumlah(t1), "SELECT COUNT(*) FROM S.DISEASE_LIFE WHERE (:1 IS NULL OR UPPER(ICD_CODE) LIKE :2")
	memuat(t, "ambil", SqlAmbil(t1), "SELECT ID, ICD_CODE, DISEASE FROM S.DISEASE_LIFE WHERE ID = :1")
	memuat(t, "pemakai", SqlPemakaiICD(t1), "WHERE UPPER(TRIM(ICD_CODE)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0)) ORDER BY ID FETCH FIRST 5 ROWS ONLY")
	memuat(t, "nomor", SqlNomorBaru("S.SEQ_DISEASE_LIFE"), "SELECT TO_CHAR(S.SEQ_DISEASE_LIFE.NEXTVAL) FROM DUAL")
	memuat(t, "ada", SqlAdaID(t1), "SELECT COUNT(*) FROM S.DISEASE_LIFE WHERE ID = :1")
	memuat(t, "sisip", SqlSisip(t1), "INSERT INTO S.DISEASE_LIFE (ID, ICD_CODE, DISEASE) VALUES (:1, :2, :3)")
	memuat(t, "ubah", SqlUbah(t1), "UPDATE S.DISEASE_LIFE SET ICD_CODE = :1, DISEASE = :2 WHERE ID = :3")
}

// ⛔ 97.586 baris: SETIAP kombinasi saring / urut menghasilkan kueri BERHALAMAN (OFFSET / FETCH NEXT terikat) - tidak
// ada jalan yang memuat seluruh DISEASE_LIFE (D3, kebutuhan teknis).
func TestSqlDaftarSelaluBerhalaman(t *testing.T) {
	for _, urut := range []string{"", models.UrutID, models.UrutICD, "lain"} {
		for _, naik := range []bool{false, true} {
			q := SqlDaftar("S.DISEASE_LIFE", urut, naik)
			if !strings.HasSuffix(q, "OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY") {
				t.Errorf("urut %q naik %v tanpa halaman: %s", urut, naik, q)
			}
		}
	}
	arg := ArgSaring(models.Saringan{ICDCode: " a0_% ", Disease: ""})
	if len(arg) != 4 || arg[0] != `%A0\_\%%` || arg[1] != arg[0] || arg[2] != nil || arg[3] != nil {
		t.Errorf("arg saring %#v", arg)
	}
}

// Lapis penjaga: tulis hanya ke DISEASE_LIFE; objek lain (M_DISEASE_LIFE JSON lama, tabel klaim) ditolak.
func TestPeriksaTulis(t *testing.T) {
	if !slices.Equal(DaftarTabelDitulis, []string{"DISEASE_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	for _, objek := range []string{"M_DISEASE_LIFE", "T_CLAIMLF_DIAGNOSE", "T_CLAIMLF_PREMIUMLIST_DETAIL", "M_DISEASE_LIFE_SEQ"} {
		if err := PeriksaTulis(objek, "UPDATE X SET A = 1"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
	}
	if err := PeriksaTulis("DUAL", "SELECT 1 FROM DUAL"); err != nil {
		t.Errorf("SELECT: %v", err)
	}
}

// ORA-00001 (PK_DISEASE_LIFE) = ErrKembar; tabel / sequence tidak ada = ErrBelumAda.
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

// D1: nol JSON, nol objek lama (M_DISEASE_LIFE / M_DISEASE_LIFE_SEQ / PEGA_*), nol DELETE di baris kode produksi
// repository; kolom = kolom warisan.
func TestNolJSONDanObjekLama(t *testing.T) {
	if !slices.Equal(KolomTabel, []string{"ID", "ICD_CODE", "DISEASE"}) {
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
			for _, larang := range []string{"JSON", "M_DISEASE_LIFE", "PEGA_", "DELETE"} {
				if strings.Contains(kode, larang) {
					t.Errorf("%s:%d memuat %s: %s", b, i+1, larang, kode)
				}
			}
		}
	}
}
