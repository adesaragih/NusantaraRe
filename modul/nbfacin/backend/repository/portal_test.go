package repository

import (
	"strings"
	"testing"
)

// TestSQLPortal - teks SQL tiket 32: HANYA LINI Fac In (bind :1, case Life tidak tampil) dan HANYA buatan pelaku
// (CREATE_OP = :2, keputusan work owner 07-10-2026), cari tidak peka huruf atas ID dan BUSINESS_PROSPECT_NAME dengan
// ESCAPE, terbaru dulu, paging bernomor sesudah saringan; cacah memakai saringan yang sama (total ikut tersaring).
func TestSQLPortal(t *testing.T) {
	tb := tabelPortal{work: "UJI.W", opp: "UJI.O", akun: "UJI.A", quo: "UJI.Q", mo: "UJI.M", general: "UJI.G"}
	saring := ` WHERE w.LINI = :1 AND w.CREATE_OP = :2 AND (UPPER(w.ID) LIKE :3 ESCAPE '\' OR UPPER(o.BUSINESS_PROSPECT_NAME) LIKE :4 ESCAPE '\')`
	for _, u := range []struct {
		cari              bool
		akhir, tidakBoleh string
	}{
		{false, " WHERE w.LINI = :1 AND w.CREATE_OP = :2\nORDER BY w.TGL_CREATE DESC, w.ID DESC OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY", "LIKE"},
		{true, saring + "\nORDER BY w.TGL_CREATE DESC, w.ID DESC OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY", ":7"},
	} {
		q := sqlCariPortal(tb, u.cari)
		if !strings.HasSuffix(q, u.akhir) || strings.Contains(q, u.tidakBoleh) || strings.Contains(q, "SELECT *") {
			t.Errorf("cari=%v:\n%s", u.cari, q)
		}
		for _, harus := range []string{"SELECT w.ID, o.BUSINESS_PROSPECT_NAME, o.GROUP_BUSINESS,",
			"(SELECT MAX(a.INSUREDNAME) FROM UJI.A a WHERE a.ID = o.ACCOUNT_ID)",
			"(SELECT MAX(m.CLIENTNAME) FROM UJI.M m WHERE m.ID = q.MOID)", "w.STATUS_WORK, g.NB_STATUS\nFROM UJI.W w",
			"LEFT JOIN UJI.O o ON o.ID = w.ID", "LEFT JOIN UJI.Q q ON q.PARENT_ID = w.ID", "LEFT JOIN UJI.G g ON g.ID = w.ID"} {
			if !strings.Contains(q, harus) {
				t.Errorf("cari=%v tanpa %q", u.cari, harus)
			}
		}
	}
	if got := sqlCacahPortal(tb, true); got != "SELECT COUNT(*) FROM UJI.W w LEFT JOIN UJI.O o ON o.ID = w.ID"+saring {
		t.Errorf("cacah %q", got)
	}
	if got := sqlCacahPortal(tb, false); got != "SELECT COUNT(*) FROM UJI.W w LEFT JOIN UJI.O o ON o.ID = w.ID WHERE w.LINI = :1 AND w.CREATE_OP = :2" {
		t.Errorf("cacah %q", got)
	}
}
