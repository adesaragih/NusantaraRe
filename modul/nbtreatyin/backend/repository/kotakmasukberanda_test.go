package repository

// Cacah kotak masuk Beranda: SATU query ber-GROUP BY posisi; Admin = buatan akun yang masih di Admin, atasan =
// antrean yang dipegang; nilai selalu terikat, penampung urut kemunculan (keputusan work owner 06-10-2026).

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestSqlHitungKotakMasuk(t *testing.T) {
	q, args := sqlHitungKotakMasuk("S.W", "S.G", "UJI-A", true, []string{models.PosisiSecHead, models.PosisiDeptHead})
	satu := strings.Join(strings.Fields(q), " ")
	harap := "SELECT g.POSITION_NOTE, COUNT(*) FROM S.W w JOIN S.G g ON g.ID = w.ID" +
		" WHERE g.PRODKE = 0 AND w.LINI = :1 AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT IN (:2, :3))" +
		" AND ((g.POSITION_NOTE = :4 AND w.CREATE_OP = :5) OR g.POSITION_NOTE IN (:6, :7)) GROUP BY g.POSITION_NOTE"
	if satu != harap {
		t.Fatalf("query:\n %s\nharap:\n %s", satu, harap)
	}
	want := []any{models.LiniKasus, models.StatusDitolak, models.StatusSelesai, models.PosisiAdmin, "UJI-A",
		models.PosisiSecHead, models.PosisiDeptHead}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("args %#v", args)
	}
	q, args = sqlHitungKotakMasuk("S.W", "S.G", "UJI-A", false, []string{models.PosisiSecHead})
	if !strings.Contains(strings.Join(strings.Fields(q), " "), "AND (g.POSITION_NOTE IN (:4))") || len(args) != 4 {
		t.Errorf("tanpa Admin: %s %v", q, args)
	}
}
