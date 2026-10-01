package login

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

var _ Gudang = (*GudangOracle)(nil)

func TestSQLGudangBersihDanBerbind(t *testing.T) {
	for nama, q := range map[string]string{
		"ambil":      sqlAmbilAkun("S.M_LOGIN_GO"),
		"workbasket": sqlWorkbasket("S.M_LOGIN_GO_WORKBASKET", "S.M_WORKBASKET"),
		"gagal":      sqlCatatGagal("S.M_LOGIN_GO"),
		"berhasil":   sqlCatatBerhasil("S.M_LOGIN_GO"),
		"ganti":      sqlGantiSandi("S.M_LOGIN_GO"),
		"versi":      sqlNaikkanVersi("S.M_LOGIN_GO"),
		"unit":       sqlInfoUnit("S.M_UNIT", "S.M_DIVISION"),
		"divisi":     sqlInfoDivisi("S.M_DIVISION", "S.M_ORGANIZATION"),
		"organisasi": sqlInfoOrganisasi("S.M_ORGANIZATION"),
		"wb aktif":   sqlWorkbasketAktif("S.M_WORKBASKET"),
		"sisip":      sqlSisipAkun("S.M_LOGIN_GO"),
		"sisip wb":   sqlSisipWorkbasket("S.M_LOGIN_GO_WORKBASKET"),
	} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", nama, err)
		}
		if !strings.Contains(q, ":1") {
			t.Errorf("%s tanpa bind:\n%s", nama, q)
		}
		// Nol literal bertanda kutip: setiap nilai lewat bind.
		if regexp.MustCompile(`'[^']*'`).MatchString(q) {
			t.Errorf("%s memuat literal:\n%s", nama, q)
		}
	}
}

// Kunci dihitung dan diperiksa dengan jam Oracle (SYSDATE) - satu jam saja.
func TestKunciMemakaiJamOracle(t *testing.T) {
	if q := sqlAmbilAkun("S.M_LOGIN_GO"); !strings.Contains(q, "CASE WHEN LOCKED_UNTIL > SYSDATE THEN 1 ELSE 0 END") {
		t.Errorf("status kunci tidak dihitung Oracle:\n%s", q)
	}
	q := sqlCatatGagal("S.M_LOGIN_GO")
	for _, w := range []string{
		"FAILED_COUNT = CASE WHEN LOCKED_UNTIL <= SYSDATE THEN 1 ELSE FAILED_COUNT + 1 END",
		">= :1", "THEN SYSDATE + :2 / 1440", "WHERE LOGIN_ID = :3",
	} {
		if !strings.Contains(q, w) {
			t.Errorf("catat gagal tanpa %q:\n%s", w, q)
		}
	}
}

// Workbasket nonaktif di master tidak menjadi peran.
func TestPeranHanyaWorkbasketAktif(t *testing.T) {
	q := sqlWorkbasket("S.M_LOGIN_GO_WORKBASKET", "S.M_WORKBASKET")
	if !strings.Contains(q, "JOIN S.M_WORKBASKET w ON w.WORKBASKET_ID = l.WORKBASKET_ID") ||
		!strings.Contains(q, "w.IS_ACTIVE = :2") {
		t.Errorf("peran tidak disaring master:\n%s", q)
	}
}

// Ganti sandi hanya dari versi sesi yang dibaca.
func TestGantiSandiDikunciVersi(t *testing.T) {
	q := sqlGantiSandi("S.M_LOGIN_GO")
	if !strings.Contains(q, "SESSION_VERSION = SESSION_VERSION + 1") || !strings.Contains(q, "AND SESSION_VERSION = :4") {
		t.Errorf("ganti sandi:\n%s", q)
	}
}
