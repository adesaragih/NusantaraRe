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

// Ringkasan dibaca DAN ditulis di kolom RIRISK_LIFE_SUMMARY: kolom XML saja, saring ber-ESCAPE, ID menaik (b9806), 50
// lewat bind; sisip / ubah kolom bernama (sisip = ringkasan baru Simpan Upload; ubah = OPERATORID/MODIFIEDDATE).
func TestSqlRingkasan(t *testing.T) {
	memuat(t, "daftar", SqlDaftar("S.RIRISK_LIFE_SUMMARY", "", false), "SELECT ID, USEDBY, OPERATORID, MODIFIEDDATE FROM S.RIRISK_LIFE_SUMMARY",
		`(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(USEDBY) LIKE :4 ESCAPE '\')`,
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC", "OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY")
	memuat(t, "urut nama", SqlDaftar("V", "usedby", true), "ORDER BY UPPER(USEDBY) DESC NULLS LAST, ID OFFSET")
	if strings.Contains(SqlDaftar("V", "ID; DROP TABLE X", false), "DROP") {
		t.Error("urut dari masukan masuk ke SQL")
	}
	memuat(t, "kembar", SqlPemakaiNama("V"), "WHERE UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))")
	memuat(t, "sisip", SqlSisipRingkasan("S.RIRISK_LIFE_SUMMARY"),
		"INSERT INTO S.RIRISK_LIFE_SUMMARY (ID, USEDBY, OPERATORID, MODIFIEDDATE) VALUES (:1, :2, :3, :4)")
	memuat(t, "ubah", SqlUbahRingkasan("S.RIRISK_LIFE_SUMMARY"),
		"UPDATE S.RIRISK_LIFE_SUMMARY SET USEDBY = :1, OPERATORID = :2, MODIFIEDDATE = :3 WHERE ID = :4")
	memuat(t, "situs", SqlSitus("S.M_SITE_DATABASE"), "SELECT TO_CHAR(ID) FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1")
	memuat(t, "nomor", SqlNomorBaru("S.M_RIRISK_LIFE_SEQ"), "SELECT TO_CHAR(S.M_RIRISK_LIFE_SEQ.NEXTVAL) FROM DUAL")
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` || PolaCari("  ") != nil {
		t.Errorf("pola %v", p)
	}
}

// Rincian = kolom RIRISK_LIFE: CONTRACT / YEAR / MONTH teks (VARCHAR2 warisan) ditulis apa adanya; RISK NUMBER tanpa
// skala ditulis koef / 10^skala dan dibaca TM9 ber-NLS eksplisit; AGE tidak disentuh; ID menaik (b7607).
func TestSqlRincian(t *testing.T) {
	baca := "SELECT ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, MONTH, TO_CHAR(RISK, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') FROM S.RIRISK_LIFE"
	memuat(t, "daftar", SqlDaftarRincian("S.RIRISK_LIFE"), baca, "WHERE IDUSEDBY = :1",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY")
	memuat(t, "kunci", SqlRincianDari("S.RIRISK_LIFE", 3), "WHERE IDUSEDBY IN (:1, :2, :3)")
	memuat(t, "sisip", SqlSisipRincian("S.RIRISK_LIFE"),
		"INSERT INTO S.RIRISK_LIFE (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, MONTH, RISK) VALUES (:1, :2, :3, :4, :5, :6, TO_NUMBER(:7) / POWER(10, :8))")
	memuat(t, "ubah", SqlUbahRincian("S.RIRISK_LIFE"),
		"UPDATE S.RIRISK_LIFE SET CONTRACT = :1, YEAR = :2, MONTH = :3, RISK = TO_NUMBER(:4) / POWER(10, :5) WHERE ID = :6 AND IDUSEDBY = :7")
	memuat(t, "hapus", SqlHapusRincian("S.RIRISK_LIFE"), "DELETE FROM S.RIRISK_LIFE WHERE IDUSEDBY = :1")
	memuat(t, "nama", SqlUbahNamaRincian("S.RIRISK_LIFE"), "UPDATE S.RIRISK_LIFE SET USEDBY = :1 WHERE IDUSEDBY = :2")
	for _, q := range []string{SqlDaftarRincian("T"), SqlSisipRincian("T"), SqlUbahRincian("T"), SqlRincianDari("T", 1)} {
		if strings.Contains(q, "JSON") || strings.Contains(q, "AGE") {
			t.Errorf("rincian memakai JSON / AGE: %s", q)
		}
	}
	for kanonik, mau := range map[string]struct {
		koef  any
		skala int64
	}{"921.9": {"9219", 1}, "580.894351210924": {"580894351210924", 12}, "7": {"7", 0}, "0": {"0", 0}, "": {nil, 0}} {
		if koef, skala := PecahDesimal(kanonik); koef != mau.koef || skala != mau.skala {
			t.Errorf("PecahDesimal(%q) = %v,%d mau %v,%d", kanonik, koef, skala, mau.koef, mau.skala)
		}
	}
}

// Lapis penjaga: tulis hanya ke RIRISK_LIFE_SUMMARY dan RIRISK_LIFE; situs dibaca saja.
func TestPeriksaTulis(t *testing.T) {
	for _, objek := range DaftarDibacaSaja {
		if err := PeriksaTulis(objek, "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
		if err := PeriksaTulis(objek, "SELECT 1 FROM X"); err != nil {
			t.Errorf("%s SELECT: %v", objek, err)
		}
	}
	if !slices.Equal(DaftarTabelDitulis, []string{"RIRISK_LIFE_SUMMARY", "RIRISK_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	if !slices.Equal(DaftarDibacaSaja, []string{"M_SITE_DATABASE"}) {
		t.Errorf("dibaca saja %v", DaftarDibacaSaja)
	}
	if err := PeriksaTulis("M_RIRISK_LIFE_TEMP", "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
		t.Errorf("M_RIRISK_LIFE_TEMP tidak boleh disentuh (K3): %v", err)
	}
}

// K1: nol JSON dan nol nama lama M_RIRISK_LIFE* (kecuali sequence) di baris kode produksi repository.
func TestNolJSONDanNamaLamaDiRepository(t *testing.T) {
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
			if strings.Contains(kode, "JSON") || strings.Contains(kode, `"M_RIRISK_LIFE"`) || strings.Contains(kode, `"M_RIRISK_LIFE_SUMMARY"`) {
				t.Errorf("%s:%d memakai JSON / nama lama: %s", b, i+1, kode)
			}
		}
	}
}

// K2: kolom = kolom view lama (urutan sama).
func TestKolomSatuTabel(t *testing.T) {
	if !slices.Equal(KolomRincian, []string{"ID", "IDUSEDBY", "USEDBY", "AGE", "YEAR", "MONTH", "RISK", "CONTRACT"}) {
		t.Errorf("kolom rincian %v", KolomRincian)
	}
	if !slices.Equal(KolomRingkasan, []string{"ID", "USEDBY", "MODIFIEDDATE", "OPERATORID"}) {
		t.Errorf("kolom ringkasan %v", KolomRingkasan)
	}
	for _, k := range strings.Split(kolomRingkasan, ", ") {
		if !slices.Contains(KolomRingkasan, k) {
			t.Errorf("kolom ringkasan %s tidak ada di tabel", k)
		}
	}
}

// Angka dari Oracle diurai di Go tanpa bergantung NLS sesi - titik ATAU koma desimal, `TM9` tanpa nol depan.
func TestAngkaOracleTanpaNLS(t *testing.T) {
	for masuk, mau := range map[string]string{
		".9": "0.9", ",9": "0.9", "921.9": "921.9", "921,9": "921.9", "580.894351210924": "580.894351210924", "12": "12",
		"": "", "abc": "abc",
	} {
		if got := AngkaOracle(masuk); got != mau {
			t.Errorf("AngkaOracle(%q) = %q mau %q", masuk, got, mau)
		}
	}
}
