package repository

// Perakit query popup Choose Business (keputusan work owner 06-10-2026: saringan kolom dicari di server):
// saringan masuk WHERE SEBELUM `FETCH FIRST 500`, nilai selalu parameter terikat, wildcard pengguna di-escape,
// kolom di luar `models.KolomSaringBisnis` tidak pernah masuk SQL.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

func TestSqlDaftarBisnisSaringanSebelumBatas(t *testing.T) {
	q, args := sqlDaftarBisnis([]string{"ID", "TREATYID"}, "UJI_SKEMA.TREATYINDETAILJOINEDM", models.SaringanBisnis{
		JenisProporsi: "Proportional",
		Kolom:         map[string]string{"TREATYID": " 1002059 ", "TREATYCONTRACTNAME": "uji_50%", "BUKAN; DROP": "x"},
	})
	harap := `SELECT ID, TREATYID FROM UJI_SKEMA.TREATYINDETAILJOINEDM WHERE PROPORTIONTYPE = :1 AND ` +
		`UPPER(TO_CHAR(TREATYID)) LIKE :2 ESCAPE '\' AND UPPER(TO_CHAR(TREATYCONTRACTNAME)) LIKE :3 ESCAPE '\' ` +
		`ORDER BY TREATYID, ID FETCH FIRST 500 ROWS ONLY`
	if q != harap {
		t.Fatalf("query:\n %s\nharap:\n %s", q, harap)
	}
	if want := []any{"Proportional", "%1002059%", `%UJI\_50\%%`}; !reflect.DeepEqual(args, want) {
		t.Errorf("args = %#v, harap %#v", args, want)
	}
	if strings.Contains(q, "BUKAN") || strings.Contains(q, "DROP") {
		t.Error("kolom di luar daftar popup masuk SQL")
	}
}

func TestSqlDaftarBisnisPemeriksaanIDDanTanpaSaringan(t *testing.T) {
	q, args := sqlDaftarBisnis([]string{"ID"}, "V", models.SaringanBisnis{ID: "UJI-D1"})
	if q != `SELECT ID FROM V WHERE ID = :1 ORDER BY TREATYID, ID FETCH FIRST 500 ROWS ONLY` || !reflect.DeepEqual(args, []any{"UJI-D1"}) {
		t.Errorf("pemeriksaan ID: %s %v", q, args)
	}
	q, args = sqlDaftarBisnis([]string{"ID"}, "V", models.SaringanBisnis{})
	if q != `SELECT ID FROM V ORDER BY TREATYID, ID FETCH FIRST 500 ROWS ONLY` || len(args) != 0 {
		t.Errorf("tanpa saringan: %s %v", q, args)
	}
}
