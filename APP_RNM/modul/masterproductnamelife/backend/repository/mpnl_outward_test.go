package repository

// Teks SQL `BrowseReinstypeOR_SQL` b84 (paket 9) - tanpa Oracle.

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSQLReinstypeORMengikutiRDB(t *testing.T) {
	q := sqlReinstypeOR("S.TREATYCONTRACT_LIFE", "S.TREATYYEAR_LIFE")
	for _, w := range []string{
		"FROM S.TREATYCONTRACT_LIFE tc JOIN S.TREATYYEAR_LIFE ty ON ty.ID = tc.IDTREATYYEAR",
		"WHERE TO_DATE(:1, 'DD/MM/YYYY') >= tc.TREATYSTARTDATE",
		"AND TO_DATE(:2, 'DD/MM/YYYY') <= tc.TREATYENDDATE",
		"AND tc.REINSTYPEID = '10200'",
		"SELECT tc.REINSTYPEID, tc.REINSTYPENAME, ty.TREATYYEAR, ty.UNDERWRITINGYEAR",
	} {
		if !strings.Contains(rata(q), w) {
			t.Errorf("SQL tanpa %q:\n%s", w, rata(q))
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	if strings.Contains(q, "FOR UPDATE") || strings.Contains(q, "TREATYCONTRACTID") {
		t.Error("pembaca master tidak mengunci; kolom TREATYCONTRACTID tidak ada di kedua tabel (OQ-MPNL-15)")
	}
}

func TestReinstypeORTanggalKosongNolBarisTanpaOracle(t *testing.T) {
	g := &Gudang{}
	for _, d := range [][2]string{{"", "01/01/2027"}, {"01/01/2026", ""}} {
		b, err := g.DaftarReinstypeOR(context.Background(), nil, d[0], d[1])
		if err != nil || b == nil || len(b) != 0 {
			t.Errorf("%v: %v %v", d, b, err)
		}
	}
}
