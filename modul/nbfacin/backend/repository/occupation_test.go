package repository

import (
	"strings"
	"testing"
)

// TestSQLOccupation - tiket 38: TYPE persis, OLDID ATAU NAME mengandung (UPPER, ESCAPE), urut RD
// NAME lalu OLDID, batas terikat; empat bind berurutan.
func TestSQLOccupation(t *testing.T) {
	q := sqlCariOccupation("UJI.OCCUPATION")
	mau := `SELECT OLDID, NAME, KDRISKEXPOSURE FROM UJI.OCCUPATION WHERE TYPE = :1 AND (UPPER(OLDID) LIKE :2 ESCAPE '\' OR UPPER(NAME) LIKE :3 ESCAPE '\') ORDER BY NAME, OLDID FETCH FIRST :4 ROWS ONLY`
	if q != mau {
		t.Errorf("SQL\n%s\nmau\n%s", q, mau)
	}
	for i := 1; i <= 4; i++ {
		if strings.Count(q, ":"+string(rune('0'+i))) != 1 {
			t.Errorf("bind :%d tidak tepat satu", i)
		}
	}
	// tiket 40: kata kosong = seluruh okupasi FIRE.
	if q := sqlSemuaOccupation("UJI.OCCUPATION"); q != "SELECT OLDID, NAME, KDRISKEXPOSURE FROM UJI.OCCUPATION WHERE TYPE = :1 ORDER BY NAME, OLDID FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("semua: %q", q)
	}
	if TipeOccupationFire != "FIRE" || BatasSaranOccupation != 500 {
		t.Error("parameter RD: TYPE FIRE, pyMaxRecords 500")
	}
	if got := PolaCari(" pab_rik% "); got != `%PAB\_RIK\%%` {
		t.Errorf("pola %q", got)
	}
}

// TestSQLRateLowerOkupasi - RateFlexas_Sql untuk banyak kode: OCCUPATIONCODE IN (bind berurutan), RATELOWER lewat
// angkaKeluar, urut kode lalu ID (baris pertama per kode), tanpa saringan lain.
func TestSQLRateLowerOkupasi(t *testing.T) {
	if q := sqlRateLowerOkupasi("UJI.OCCUPATION", 3); q != "SELECT OCCUPATIONCODE, "+angkaKeluar("RATELOWER")+
		" FROM UJI.OCCUPATION WHERE OCCUPATIONCODE IN (:1, :2, :3) ORDER BY OCCUPATIONCODE, ID" {
		t.Errorf("%s", q)
	}
}
