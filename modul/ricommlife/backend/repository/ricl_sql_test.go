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

// Ringkasan dibaca DAN ditulis di kolom M_RICOMM_LIFE_SUMMARY (RALAT R1): kolom XML saja, saring ber-ESCAPE, ID menaik
// (b9857), 50 lewat bind; sisip / ubah kolom bernama.
func TestSqlRingkasan(t *testing.T) {
	memuat(t, "daftar", SqlDaftar("S.M_RICOMM_LIFE_SUMMARY", "", false), "SELECT ID, USEDBY, OPERATORID, MODIFIEDDATE FROM S.M_RICOMM_LIFE_SUMMARY",
		`(:1 IS NULL OR UPPER(ID) LIKE :2 ESCAPE '\') AND (:3 IS NULL OR UPPER(USEDBY) LIKE :4 ESCAPE '\')`,
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC", "OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY")
	memuat(t, "urut nama", SqlDaftar("V", "usedby", true), "ORDER BY UPPER(USEDBY) DESC NULLS LAST, ID OFFSET")
	if strings.Contains(SqlDaftar("V", "ID; DROP TABLE X", false), "DROP") {
		t.Error("urut dari masukan masuk ke SQL")
	}
	memuat(t, "kembar", SqlPemakaiNama("V"), "WHERE UPPER(TRIM(USEDBY)) = UPPER(TRIM(:1)) AND ID <> NVL(:2, CHR(0))")
	memuat(t, "sisip", SqlSisipRingkasan("S.M_RICOMM_LIFE_SUMMARY"),
		"INSERT INTO S.M_RICOMM_LIFE_SUMMARY (ID, USEDBY, OPERATORID, MODIFIEDDATE) VALUES (:1, :2, :3, :4)")
	memuat(t, "ubah", SqlUbahRingkasan("S.M_RICOMM_LIFE_SUMMARY"),
		"UPDATE S.M_RICOMM_LIFE_SUMMARY SET USEDBY = :1, OPERATORID = :2, MODIFIEDDATE = :3 WHERE ID = :4")
	memuat(t, "situs", SqlSitus("S.M_SITE_DATABASE"), "SELECT TO_CHAR(ID) FROM S.M_SITE_DATABASE WHERE CURRENT_SITE = :1")
	memuat(t, "nomor", SqlNomorBaru("S.M_RICOMM_LIFE_SEQ"), "SELECT TO_CHAR(S.M_RICOMM_LIFE_SEQ.NEXTVAL) FROM DUAL")
	if p := PolaCari(" a_b% "); p != `%A\_B\%%` || PolaCari("  ") != nil {
		t.Errorf("pola %v", p)
	}
}

// Rincian = kolom M_RICOMM_LIFE (RALAT R1): kolom bernama, angka tanpa NLS, ID menaik (b9515).
func TestSqlKomisi(t *testing.T) {
	baca := "SELECT ID, IDUSEDBY, USEDBY, TO_CHAR(CONTRACT, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), " +
		"TO_CHAR(YEAR, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), TO_CHAR(COMM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') FROM S.M_RICOMM_LIFE"
	memuat(t, "daftar", SqlDaftarKomisi("S.M_RICOMM_LIFE"), baca, "WHERE IDUSEDBY = :1",
		"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '^[0-9]+$')) ASC NULLS LAST, ID ASC OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY")
	memuat(t, "kunci", SqlKomisiDari("S.M_RICOMM_LIFE", 3), "WHERE IDUSEDBY IN (:1, :2, :3)")
	memuat(t, "sisip", SqlSisipKomisi("S.M_RICOMM_LIFE"),
		"INSERT INTO S.M_RICOMM_LIFE (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM) VALUES (:1, :2, :3, TO_NUMBER(:4), TO_NUMBER(:5), TO_NUMBER(:6) / POWER(10, :7))")
	memuat(t, "ubah", SqlUbahKomisi("S.M_RICOMM_LIFE"),
		"UPDATE S.M_RICOMM_LIFE SET CONTRACT = TO_NUMBER(:1), YEAR = TO_NUMBER(:2), COMM = TO_NUMBER(:3) / POWER(10, :4) WHERE ID = :5 AND IDUSEDBY = :6")
	memuat(t, "nama", SqlUbahNamaKomisi("S.M_RICOMM_LIFE"), "UPDATE S.M_RICOMM_LIFE SET USEDBY = :1 WHERE IDUSEDBY = :2")
	memuat(t, "hapus", SqlHapusKomisi("S.M_RICOMM_LIFE"), "DELETE FROM S.M_RICOMM_LIFE WHERE IDUSEDBY = :1")
	for kanonik, mau := range map[string]struct {
		koef  any
		skala int64
	}{"12.05": {"1205", 2}, "0.5": {"5", 1}, "7": {"7", 0}, "0": {"0", 0}, "": {nil, 0}, "100": {"100", 0}} {
		if koef, skala := PecahDesimal(kanonik); koef != mau.koef || skala != mau.skala {
			t.Errorf("PecahDesimal(%q) = %v,%d mau %v,%d", kanonik, koef, skala, mau.koef, mau.skala)
		}
	}
}

// Lapis penjaga: tulis hanya ke M_RICOMM_LIFE_SUMMARY dan M_RICOMM_LIFE; situs dibaca saja.
func TestPeriksaTulis(t *testing.T) {
	for _, objek := range DaftarDibacaSaja {
		if err := PeriksaTulis(objek, "DELETE FROM X"); !errors.Is(err, ErrBacaSaja) {
			t.Errorf("%s ditulis: %v", objek, err)
		}
		if err := PeriksaTulis(objek, "SELECT 1 FROM X"); err != nil {
			t.Errorf("%s SELECT: %v", objek, err)
		}
	}
	if !slices.Equal(DaftarTabelDitulis, []string{"M_RICOMM_LIFE_SUMMARY", "M_RICOMM_LIFE"}) {
		t.Errorf("tabel ditulis %v", DaftarTabelDitulis)
	}
	if !slices.Equal(DaftarDibacaSaja, []string{"M_SITE_DATABASE"}) {
		t.Errorf("dibaca saja %v", DaftarDibacaSaja)
	}
	if err := PeriksaTulis("SESI", "ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '.,'"); !errors.Is(err, ErrBacaSaja) {
		t.Errorf("pernyataan sesi harus ditolak (alat pindah dihapus): %v", err)
	}
}

// RALAT R1: nol JSON dan nol objek lama (view RICOMM_LIFE_SUMMARY, tabel flat RICOMM_LIFE) di baris kode produksi
// repository - komentar boleh menyebutnya sebagai riwayat.
func TestNolJSONDanObjekLamaDiRepository(t *testing.T) {
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
			if strings.Contains(kode, "JSON") || strings.Contains(kode, `"RICOMM_LIFE"`) || strings.Contains(kode, `"RICOMM_LIFE_SUMMARY"`) {
				t.Errorf("%s:%d memakai JSON / objek lama: %s", b, i+1, kode)
			}
		}
	}
}

// Kolom = kolom view ringkasan lama dan tabel flat 924 (keputusan work owner 08-10-2026).
func TestKolomSatuTabel(t *testing.T) {
	if !slices.Equal(KolomKomisi, []string{"ID", "IDUSEDBY", "USEDBY", "CONTRACT", "YEAR", "COMM"}) {
		t.Errorf("kolom rincian %v", KolomKomisi)
	}
	if !slices.Equal(KolomRingkasan, []string{"ID", "USEDBY", "MODIFIEDDATE", "OPERATORID"}) {
		t.Errorf("kolom ringkasan %v", KolomRingkasan)
	}
	for _, k := range strings.Split(kolomRingkasan, ", ") {
		if !slices.Contains(KolomRingkasan, k) {
			t.Errorf("kolom ringkasan %s tidak ada di tabel", k)
		}
	}
	for _, k := range []string{"IDUSEDBY", "USEDBY", "CONTRACT", "YEAR", "COMM"} {
		if !strings.Contains(kolomKomisi(), k) {
			t.Errorf("kolom rincian %s tidak dibaca", k)
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
}
