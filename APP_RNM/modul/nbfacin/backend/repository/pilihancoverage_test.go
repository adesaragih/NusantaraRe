package repository

import "testing"

// TestSQLPilihanCoverage - tiket 43: popup Choose Coverage dari tabel COVERAGE (butir 96: TYPE FIRE peka huruf, aktif =
// '1' atau kosong, NAME Contains tidak peka huruf, DISTINCT enam kolom laporan RD, urut NAME lalu OLDID, maks 500) dan
// coverage otomatis (COVERAGE, lima kode urut korpus).
func TestSQLPilihanCoverage(t *testing.T) {
	if q := sqlCariCoverage("UJI.C", true); q != `SELECT ID, OLDID, NAME FROM (SELECT DISTINCT ID, BUSINESSCODE, NAME, TYPE, OLDID, ACTIVESTATUS FROM UJI.C WHERE TYPE = :1 AND (ACTIVESTATUS = :2 OR ACTIVESTATUS IS NULL) AND UPPER(NAME) LIKE :3 ESCAPE '\') ORDER BY NAME, OLDID, ID FETCH FIRST :4 ROWS ONLY` {
		t.Errorf("dengan kata: %q", q)
	}
	if q := sqlCariCoverage("UJI.C", false); q != "SELECT ID, OLDID, NAME FROM (SELECT DISTINCT ID, BUSINESSCODE, NAME, TYPE, OLDID, ACTIVESTATUS FROM UJI.C WHERE TYPE = :1 AND (ACTIVESTATUS = :2 OR ACTIVESTATUS IS NULL)) ORDER BY NAME, OLDID, ID FETCH FIRST :3 ROWS ONLY" {
		t.Errorf("tanpa kata: %q", q)
	}
	if q := sqlCoverageOtomatis("UJI.C"); q != "SELECT ID, OLDID, NAME FROM UJI.C WHERE ID IN (:1, :2, :3, :4, :5) ORDER BY ID, NAME, OLDID" {
		t.Errorf("otomatis: %q", q)
	}
	if TipeCoverageFire != "FIRE" || BatasCoverage != 500 || StatusAktifCoverage != "1" || TabelCoverage != "COVERAGE" {
		t.Error("parameter RD / nama tabel DDL berubah")
	}
	mau := []string{"100815", "100828", "100829", "100825", "100840"}
	for i, k := range mau {
		if len(kodeCoverageOtomatis) != len(mau) || kodeCoverageOtomatis[i] != k {
			t.Fatalf("kode otomatis %v, mau urutan korpus %v", kodeCoverageOtomatis, mau)
		}
	}
}
