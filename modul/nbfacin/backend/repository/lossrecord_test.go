package repository

import (
	"regexp"
	"strings"
	"testing"
)

// TestSQLLossRecord - kolom yang ditulis / dibaca ada di migrasi 206; jumlah bind = jumlah kolom; uang lewat angkaMasuk;
// hapus anak -> induk; kepala dibaca lewat case FAC (LINI); jabatan dari LOGIN_ID.
func TestSQLLossRecord(t *testing.T) {
	kolom := map[string]map[string]bool{}
	for _, m := range regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.(\w+) \((.*?)\n\)`).FindAllStringSubmatch(bacaMigrasi(t, "206_loss_record_kasus.sql"), -1) {
		kolom[m[1]] = map[string]bool{}
		for _, b := range strings.Split(m[2], "\n") {
			if f := strings.Fields(b); len(f) > 0 && f[0] != "CONSTRAINT" {
				kolom[m[1]][f[0]] = true
			}
		}
	}
	if len(kolom) != 2 || len(kolom[TabelLossRecord]) != 12 || len(kolom[TabelLossRecordList]) != 13 {
		t.Fatalf("migrasi 206 terbaca %v", kolom)
	}
	for _, q := range []string{sqlSisipKepalaLossRecord(TabelLossRecord), sqlSisipBarisLossRecord(TabelLossRecordList)} {
		m := regexp.MustCompile(`INSERT INTO (\w+) \(([^)]*)\) VALUES \((.*)\)$`).FindStringSubmatch(q)
		cols := strings.Split(m[2], ", ")
		nilai := regexp.MustCompile(`TO_NUMBER\(:\d+[^)]*\)|:\d+`).FindAllString(m[3], -1)
		if len(cols) != len(nilai) || len(cols) != len(kolom[m[1]]) {
			t.Errorf("%s: %d kolom / %d nilai / %d di migrasi", m[1], len(cols), len(nilai), len(kolom[m[1]]))
		}
		for _, c := range cols {
			if !kolom[m[1]][c] {
				t.Errorf("%s.%s ditulis, tidak ada di 206", m[1], c)
			}
		}
	}
	tb := tabelLossRecord{work: "UJI.W", general: "UJI.G", kepala: TabelLossRecord, baris: TabelLossRecordList}
	for _, m := range regexp.MustCompile(`\bk\.([A-Z_0-9]+)`).FindAllStringSubmatch(sqlBacaKepalaLossRecord(tb), -1) {
		if !kolom[TabelLossRecord][m[1]] {
			t.Errorf("T_LOSSRECORD.%s dibaca, tidak ada", m[1])
		}
	}
	for _, m := range regexp.MustCompile(`\bb\.([A-Z_0-9]+)`).FindAllStringSubmatch(sqlBacaBarisLossRecord(tb), -1) {
		if !kolom[TabelLossRecordList][m[1]] {
			t.Errorf("T_LOSSRECORDLIST.%s dibaca, tidak ada", m[1])
		}
	}
	if !strings.Contains(sqlBacaKepalaLossRecord(tb), "WHERE w.ID = :1 AND w.LINI = :2") ||
		!strings.HasPrefix(sqlHapusLossRecord(tb)[0], "DELETE FROM T_LOSSRECORDLIST ") ||
		sqlJabatanAkun("UJI.L") != "SELECT JOB_POSITION FROM UJI.L WHERE LOGIN_ID = :1" {
		t.Error("kepala / hapus / jabatan")
	}
	if q := sqlNoteKelasBisnis("UJI.B", 2); q != "SELECT NOTE FROM UJI.B WHERE UPPER(TRIM(NOTE)) IN (:1, :2) ORDER BY ID" {
		t.Errorf("note: %s", q)
	}
	// Saran COB: dengan kata, NOTE sama persis (UPPER(TRIM) = :2) lebih dulu, lalu ID DESC, batas :3; tanpa kata seperti RD.
	if q := sqlSaranKelasBisnis("UJI.B", true); q != `SELECT ID, NOTE FROM UJI.B WHERE NOTE IS NOT NULL AND UPPER(NOTE) LIKE :1 ESCAPE '\'`+
		` ORDER BY CASE WHEN UPPER(TRIM(NOTE)) = :2 THEN 0 ELSE 1 END, ID DESC FETCH FIRST :3 ROWS ONLY` {
		t.Errorf("saran dengan kata: %s", q)
	}
	if q := sqlSaranKelasBisnis("UJI.B", false); q != "SELECT ID, NOTE FROM UJI.B WHERE NOTE IS NOT NULL ORDER BY ID DESC FETCH FIRST :1 ROWS ONLY" {
		t.Errorf("saran tanpa kata: %s", q)
	}
}
