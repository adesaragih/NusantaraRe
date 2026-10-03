package repository

import (
	"strings"
	"testing"
)

// TestSQLOccupation - tiket 38: TYPE persis, OLDID ATAU NAME mengandung (UPPER, ESCAPE), urut RD
// NAME lalu OLDID, batas terikat; empat bind berurutan.
func TestSQLOccupation(t *testing.T) {
	q := sqlCariOccupation("UJI.OCCUPATION")
	mau := `SELECT OLDID, NAME FROM UJI.OCCUPATION WHERE TYPE = :1 AND (UPPER(OLDID) LIKE :2 ESCAPE '\' OR UPPER(NAME) LIKE :3 ESCAPE '\') ORDER BY NAME, OLDID FETCH FIRST :4 ROWS ONLY`
	if q != mau {
		t.Errorf("SQL\n%s\nmau\n%s", q, mau)
	}
	for i := 1; i <= 4; i++ {
		if strings.Count(q, ":"+string(rune('0'+i))) != 1 {
			t.Errorf("bind :%d tidak tepat satu", i)
		}
	}
	if TipeOccupationFire != "FIRE" || BatasSaranOccupation != 500 {
		t.Error("parameter RD: TYPE FIRE, pyMaxRecords 500")
	}
	if got := PolaCari(" pab_rik% "); got != `%PAB\_RIK\%%` {
		t.Errorf("pola %q", got)
	}
}
