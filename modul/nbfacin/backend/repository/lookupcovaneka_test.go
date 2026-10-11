package repository

import "testing"

// TestSQLCoverageAneka - BrowseCoverageFacIn_RD Type ANEKA: TYPE + aktif selalu, BUSINESSCODE (peka huruf) dan NAME (tak peka
// huruf) hanya bila diisi; bentuk DISTINCT / urutan = FIRE; bind 1..n pada keempat kombinasi.
func TestSQLCoverageAneka(t *testing.T) {
	awal := "SELECT ID, OLDID, NAME FROM (SELECT DISTINCT ID, BUSINESSCODE, NAME, TYPE, OLDID, ACTIVESTATUS FROM S.C WHERE TYPE = :1 AND (ACTIVESTATUS = :2 OR ACTIVESTATUS IS NULL)"
	for _, c := range []struct {
		bisnis, kata bool
		mau          string
	}{
		{false, false, awal + ") ORDER BY NAME, OLDID, ID FETCH FIRST :3 ROWS ONLY"},
		{true, false, awal + " AND BUSINESSCODE = :3) ORDER BY NAME, OLDID, ID FETCH FIRST :4 ROWS ONLY"},
		{false, true, awal + ` AND UPPER(NAME) LIKE :3 ESCAPE '\') ORDER BY NAME, OLDID, ID FETCH FIRST :4 ROWS ONLY`},
		{true, true, awal + ` AND BUSINESSCODE = :3 AND UPPER(NAME) LIKE :4 ESCAPE '\') ORDER BY NAME, OLDID, ID FETCH FIRST :5 ROWS ONLY`},
	} {
		if q := sqlCariCoverageAneka("S.C", c.bisnis, c.kata); q != c.mau || bindTerbesar(t, q) < 3 {
			t.Errorf("%v/%v: %q", c.bisnis, c.kata, q)
		}
	}
	if TipeCoverageAneka != "ANEKA" {
		t.Error("literal")
	}
}
