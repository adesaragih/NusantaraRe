package repository

import (
	"database/sql"
	"strings"
	"testing"
)

// TestSQLLimitKolomPersis - teks SQL membaca PERSIS kolom yang dipakai rule Pega
// (tiket 11/12) dari tabel yang sudah dikualifikasi, tanpa NAMA dan LOGIN
// (CLAUDE.md §4 butir 10) dan tanpa `SELECT *`.
func TestSQLLimitKolomPersis(t *testing.T) {
	for _, u := range []struct{ sql, mau string }{
		{sqlBentukA("UJI.M_LIMIT_PROPERTYY"), "SELECT JABATAN, TEAM_GROUP, LIMIT_BOTTOM, LIMIT_BOTTOM2 FROM UJI.M_LIMIT_PROPERTYY"},
		{sqlBentukB("UJI.M_LIMIT_FINANCIALINS"), "SELECT JABATAN, LIMITBOND_BOTTOM, LIMITCREDITCL_BOTTOM, LIMITCREDITNCL_BOTTOM FROM UJI.M_LIMIT_FINANCIALINS"},
	} {
		if u.sql != u.mau {
			t.Errorf("SQL %q, mau %q", u.sql, u.mau)
		}
		for _, terlarang := range []string{"NAMA", "LOGIN", "*"} {
			if strings.Contains(u.sql, terlarang) {
				t.Errorf("SQL %q memuat %q", u.sql, terlarang)
			}
		}
	}
}

// TestTabelBentukA - lima tabel bentuk A, ejaan DDL apa adanya (huruf ganda di akhir
// adalah nama tabel sungguhan, `D:\migrasi\RNM\DDL\M_LIMIT_*.txt`).
func TestTabelBentukA(t *testing.T) {
	mau := "M_LIMIT_PROPERTYY,M_LIMIT_PROPERTY_NON_PREFERREDD,M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL,M_LIMIT_ENGINEERINGG,M_LIMIT_NONPROPANDENGG"
	if got := strings.Join(TabelBentukA, ","); got != mau {
		t.Errorf("tabel bentuk A %s", got)
	}
	if TabelBentukB != "M_LIMIT_FINANCIALINS" {
		t.Errorf("tabel bentuk B %s", TabelBentukB)
	}
}

// TestUraiBarisLimit - NUMBER dibaca teks lalu diurai desimal: NULL = nil (bukan
// nol), angka rusak = galat yang menyebut tabel dan kolom.
func TestUraiBarisLimit(t *testing.T) {
	b, err := uraiBarisA("M_LIMIT_PROPERTYY", teks("SENIORUW"), teks("1"), teks("178500000001"), nul())
	if err != nil || b.Jabatan != "SENIORUW" || b.TeamGroup != "1" || b.LimitBottom.String() != "178500000001" || b.LimitBottom2 != nil {
		t.Fatalf("baris A %+v (%v)", b, err)
	}
	if _, err := uraiBarisA("M_LIMIT_PROPERTYY", teks("X"), teks("1"), teks("1,5"), nul()); err == nil ||
		!strings.Contains(err.Error(), "M_LIMIT_PROPERTYY") || !strings.Contains(err.Error(), "LIMIT_BOTTOM") {
		t.Errorf("angka rusak: %v", err)
	}
	f, err := uraiBarisB(teks("DIREKTUR TEKNIK"), teks("5000000001"), nul(), teks("0"))
	if err != nil || f.Jabatan != "DIREKTUR TEKNIK" || f.LimitBond.String() != "5000000001" || f.LimitCreditCL != nil || !f.LimitCreditNCL.IsZero() {
		t.Fatalf("baris B %+v (%v)", f, err)
	}
}

func teks(s string) sql.NullString { return sql.NullString{String: s, Valid: true} }
func nul() sql.NullString          { return sql.NullString{} }
