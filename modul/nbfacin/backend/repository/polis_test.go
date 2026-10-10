package repository

import (
	"strings"
	"testing"
)

// TestSQLPolis - tiket 71: versi dari kolom relasional JSON_POLIS (nol DATA_JSON), urut TGL_INPUT turun, sistemBaru = T_WORK_POLIS
// LINI; cadangan FACINPRODUCTION MIN(IDPEGA) NOENDORS kosong; ringkasan FACINPRODUCTION 27 kolom + ceding; bind berurutan; opsi 7
// UPPER kedua sisi (A71-2).
func TestSQLPolis(t *testing.T) {
	for nama, c := range map[string]struct {
		q    string
		mau  string
		bind int
	}{
		"versi":          {sqlVersiPolis("J", "W"), "SELECT j.PRODKE, SUBSTR(j.IDPEGA, 19), (SELECT COUNT(*) FROM W w WHERE w.ID = SUBSTR(j.IDPEGA, 19) AND w.LINI = :1) FROM J j WHERE j.NOPOLIS = :2 ORDER BY j.TGL_INPUT DESC, j.PRODKE DESC", 2},
		"versi produksi": {sqlVersiProduksi("P", "W"), "SELECT SUBSTR(x.IDPEGA, 19), (SELECT COUNT(*) FROM W w WHERE w.ID = SUBSTR(x.IDPEGA, 19) AND w.LINI = :1) FROM (SELECT MIN(p.IDPEGA) AS IDPEGA FROM P p WHERE p.NOPOLIS = :2 AND p.NOENDORS IS NULL) x WHERE x.IDPEGA IS NOT NULL", 2},
		"id versi":       {sqlIDPegaVersi("J"), "SELECT IDPEGA FROM J WHERE NOPOLIS = :1 AND PRODKE = :2 ORDER BY TGL_INPUT DESC FETCH FIRST 1 ROWS ONLY", 2},
		"ada json":       {sqlAdaJSONPolis("J"), "SELECT COUNT(*) FROM J WHERE NOPOLIS = :1", 1},
		"id produksi":    {sqlIDPegaProduksi("P"), "SELECT MIN(IDPEGA) FROM P WHERE NOPOLIS = :1 AND NOENDORS IS NULL", 1},
	} {
		if c.q != c.mau || bindBerurutan(t, nama, c.q) != c.bind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	b := sqlBarisPolis("P")
	pilih := b[len("SELECT "):strings.Index(b, " FROM P p WHERE")]
	// Sensus kolom dua cara: pecah koma tingkat atas = 27; token kolom p.X + subkueri RiSlip = 27.
	n, d := 0, 0
	for _, ch := range pilih {
		switch ch {
		case '(':
			d++
		case ')':
			d--
		case ',':
			if d == 0 {
				n++
			}
		}
	}
	if n+1 != 27 || bindBerurutan(t, "baris", b) != 3 || !strings.Contains(b, "WHERE p.IDPEGA = :2 AND p.NOPOLIS = :3") ||
		!strings.Contains(b, "KEEP (DENSE_RANK LAST ORDER BY B.SPREAD_DATE) FROM P B WHERE B.NOPOLIS = :1") ||
		!strings.Contains(b, "TO_CHAR(p.BEGINDATE, 'DD-MM-YYYY')") || strings.Contains(strings.ToUpper(b), "DATA_JSON") {
		t.Errorf("baris (%d kolom): %s", n+1, b)
	}
	if q := sqlCedingPolis("C"); bindBerurutan(t, "ceding", q) != 1 || !strings.HasPrefix(q, "SELECT CEDINGNAME, ") ||
		!strings.HasSuffix(q, "FROM C WHERE IDPEGA = :1 ORDER BY CEDINGID, CEDINGNAME") {
		t.Errorf("ceding %q", q)
	}
	if k, ada := klausaFollowing("P", CariPolisQQ); !ada || k != "UPPER(PC0.QQNAME) LIKE UPPER(:1) AND PC0.NOENDORS IS NULL" {
		t.Errorf("qq %q", k)
	}
	for _, q := range []string{sqlVersiPolis("J", "W"), sqlVersiProduksi("P", "W"), sqlIDPegaVersi("J")} {
		if strings.Contains(strings.ToUpper(q), "DATA_JSON") {
			t.Errorf("DATA_JSON: %s", q)
		}
	}
}
