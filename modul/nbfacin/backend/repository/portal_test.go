package repository

import (
	"reflect"
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

// TestSQLKotakMasuk - saringan kotak masuk Beranda: LINI FAC, STATUS_WORK bukan "Resolved-%", marketing / TL = POSITION DAN
// team case IN team akun (keputusan work owner 08-10-2026), lainnya = POSITION workbasket; bind berurutan menurut kemunculan.
func TestSQLKotakMasuk(t *testing.T) {
	tb := tabelPortal{work: "UJI.W", opp: "UJI.O", akun: "UJI.A", quo: "UJI.Q", mo: "UJI.M", general: "UJI.G"}
	// Keputusan work owner 08-10-2026: marketing / TL bersaring team case (MOID -> TEAMGROUP) IN team akun; antrean ber-team
	// kosong tidak ikut; UW tanpa saringan team.
	s := SaringanKotakMasuk{Tim: []AntreanTim{{Posisi: []string{"Offer", "ReasFacInMarketing"}, Team: []string{"1"}},
		{Posisi: []string{"ReasFacInTeamLeader"}, Team: nil}, {Posisi: []string{"ReasFacInTeamLeader"}, Team: []string{"1", "3"}}},
		Lain: []string{"ReasFacInUnderwriting"}}
	tc := "TRIM((SELECT MAX(m2.TEAMGROUP) FROM UJI.Q q2 JOIN UJI.M m2 ON m2.ID = q2.MOID WHERE q2.PARENT_ID = w.ID))"
	q, args := sqlHitungKotakMasuk(tb, s)
	if q != "SELECT w.POSITION, COUNT(*) FROM UJI.W w WHERE w.LINI = :1 AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT LIKE :2)"+
		" AND ((w.POSITION IN (:3, :4) AND "+tc+" IN (:5)) OR (w.POSITION IN (:6) AND "+tc+" IN (:7, :8)) OR w.POSITION IN (:9))"+
		" GROUP BY w.POSITION" ||
		!reflect.DeepEqual(args, []any{LiniFacIn, "Resolved-%", "Offer", "ReasFacInMarketing", "1", "ReasFacInTeamLeader", "1", "3",
			"ReasFacInUnderwriting"}) || bindBerurutan(t, "hitung", q) != 9 {
		t.Errorf("%s %v", q, args)
	}
	if !(SaringanKotakMasuk{Tim: []AntreanTim{{Posisi: []string{"Offer"}}}}).kosong() || (SaringanKotakMasuk{Lain: []string{"X"}}).kosong() {
		t.Errorf("kosong(): antrean tanpa team harus kosong")
	}
	if q := sqlTeamMarketing("UJI.M"); q != "SELECT DISTINCT TRIM(TEAMGROUP) FROM UJI.M WHERE UPPER(TRIM(AKSES_LOGIN)) = UPPER(TRIM(:1))"+
		" AND TRIM(TEAMGROUP) IS NOT NULL" {
		t.Errorf("%q", q)
	}
	if q := sqlTeamLeader("UJI.L"); q != "SELECT DISTINCT TRIM(TO_CHAR(TEAM_GROUP)) FROM UJI.L WHERE UPPER(TRIM(LOGIN)) = UPPER(TRIM(:1))"+
		" AND TEAM_GROUP IS NOT NULL" {
		t.Errorf("%q", q)
	}
	q, args = sqlDaftarMenunggu(tb, "UJI.P", SaringanKotakMasuk{Lain: []string{"ReasFacInTeamLeader"}})
	if !strings.HasSuffix(q, " AND (w.POSITION IN (:5))\nORDER BY w.TGL_UPDATE DESC NULLS LAST, w.ID DESC FETCH FIRST 500 ROWS ONLY") ||
		!reflect.DeepEqual(args, []any{AwalanIDPega, AwalanIDPega, LiniFacIn, "Resolved-%", "ReasFacInTeamLeader"}) ||
		bindBerurutan(t, "daftar", q) != 5 {
		t.Errorf("%s %v", q, args)
	}
	// Kolom inbox "New Business" (tiket kotak masuk): urutan = urutan pemindai repository.
	for _, harus := range []string{"SELECT w.ID,\n  (SELECT MAX(a.INSUREDNAME) FROM UJI.A a WHERE a.ID = o.ACCOUNT_ID),\n" +
		"  o.CLASS_OF_BUSINESS, q.SOB_NAME, q.CEDING_CO_NAME, g.START_DATE_TIME, ",
		"(SELECT MAX(m.CLIENTNAME) FROM UJI.M m WHERE m.ID = q.MOID),\n  w.FLAG_ONGOING_POLICY, TO_CHAR(w.TGL_CREATE, 'DD-MM-YYYY'), g.OFFERING_DATE, g.NB_STATUS,",
		"(SELECT MAX(h.APPROVAL) KEEP (DENSE_RANK LAST ORDER BY h.TGL_INP, h.NOURUT) FROM UJI.P h WHERE h.IDPEGA = :1 || w.ID)",
		"(SELECT MAX(h.PIC) KEEP (DENSE_RANK LAST ORDER BY h.TGL_INP, h.NOURUT) FROM UJI.P h WHERE h.IDPEGA = :2 || w.ID)",
		"TO_CHAR(w.TGL_UPDATE, 'YYYY-MM-DD HH24:MI:SS'), g.SUBMIT_TO_UW\nFROM UJI.W w", "LEFT JOIN UJI.G g ON g.ID = w.ID"} {
		if !strings.Contains(q, harus) {
			t.Errorf("inbox tanpa %q", harus)
		}
	}
	if q := sqlNamaWorkbasket("UJI.WB", 2); q != "SELECT WORKBASKET_ID, NAME FROM UJI.WB WHERE WORKBASKET_ID IN (:1, :2)" {
		t.Errorf("%q", q)
	}
}
