package repository

import "testing"

// TestSQLPilihanCoverage - tiket 43: popup Choose Coverage (COVERAGE_FACIN, TYPE FIRE, DISTINCT kolom laporan RD tanpa
// ACTIVESTATUS - A160, urut NamaCoverage lalu OLDID, maks 500) dan coverage otomatis (COVERAGE, lima kode urut korpus).
func TestSQLPilihanCoverage(t *testing.T) {
	if q := sqlCariCoverage("UJI.CF", true); q != `SELECT ID, OLDID, NAMACOVERAGE FROM (SELECT DISTINCT TO_CHAR(ID) AS ID, BIZCODE, NAMACOVERAGE, TYPE, OLDID FROM UJI.CF WHERE TYPE = :1 AND UPPER(NAMACOVERAGE) LIKE :2 ESCAPE '\') ORDER BY NAMACOVERAGE, OLDID, ID FETCH FIRST :3 ROWS ONLY` {
		t.Errorf("dengan kata: %q", q)
	}
	if q := sqlCariCoverage("UJI.CF", false); q != "SELECT ID, OLDID, NAMACOVERAGE FROM (SELECT DISTINCT TO_CHAR(ID) AS ID, BIZCODE, NAMACOVERAGE, TYPE, OLDID FROM UJI.CF WHERE TYPE = :1) ORDER BY NAMACOVERAGE, OLDID, ID FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("tanpa kata: %q", q)
	}
	if q := sqlCoverageOtomatis("UJI.C"); q != "SELECT ID, OLDID, NAME FROM UJI.C WHERE ID IN (:1, :2, :3, :4, :5) ORDER BY ID, NAME, OLDID" {
		t.Errorf("otomatis: %q", q)
	}
	if TipeCoverageFire != "FIRE" || BatasCoverage != 500 || TabelCoverageFacIn != "COVERAGE_FACIN" || TabelCoverage != "COVERAGE" {
		t.Error("parameter RD / nama tabel DDL berubah")
	}
	mau := []string{"100815", "100828", "100829", "100825", "100840"}
	for i, k := range mau {
		if len(kodeCoverageOtomatis) != len(mau) || kodeCoverageOtomatis[i] != k {
			t.Fatalf("kode otomatis %v, mau urutan korpus %v", kodeCoverageOtomatis, mau)
		}
	}
}
