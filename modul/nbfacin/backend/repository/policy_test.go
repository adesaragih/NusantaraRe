package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// TestSQLHitungAkumulasi - RDBList\GetCountAccumulation.xml L86-89: FACINPRODUCTION dikecualikan :3 (pzInsKey), FACINOFFER :6
// (pyID); ENDDATE teks dd/mm/rrrr; tanpa COMMIT.
func TestSQLHitungAkumulasi(t *testing.T) {
	q := sqlHitungAkumulasi("P.FACINPRODUCTION", "P.FACINOFFER")
	for _, a := range []string{"select SUM(CARID1)AS CARID1 FROM (",
		"from P.FACINPRODUCTION where ACCUMULATIONCODE = :1 AND TO_DATE (ENDDATE, 'dd/mm/rrrr') >= TO_DATE (:2, 'dd/mm/rrrr') and IDPEGA != :3",
		"UNION ALL", "from P.FACINOFFER where ACCUMULATIONCODE = :4 AND TO_DATE (ENDDATE, 'dd/mm/rrrr') >= TO_DATE (:5, 'dd/mm/rrrr')and IDPEGA != :6)"} {
		if !strings.Contains(q, a) {
			t.Errorf("tanpa %q\n%s", a, q)
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
}
