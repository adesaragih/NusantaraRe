package repository

import "testing"

// TestSQLUbahAkumulasiCoverage - tiket 60: hanya dua kolom akumulasi, satu baris per ID, bind berurutan.
func TestSQLUbahAkumulasiCoverage(t *testing.T) {
	q := sqlUbahAkumulasiCoverage("UJI.C")
	if q != "UPDATE UJI.C SET ACCUMULATION_CODE = :1, ACCUMULATION_DESCRIPTION = :2 WHERE ID = :3" || bindBerurutan(t, "akumulasi", q) != 3 {
		t.Errorf("%q", q)
	}
}
