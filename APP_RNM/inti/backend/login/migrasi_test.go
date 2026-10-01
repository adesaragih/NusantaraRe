package login

// Migrasi 902 - M_LOGIN_GO dan M_LOGIN_GO_WORKBASKET. TANPA Oracle.

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func baca902(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(os.DirFS(".."), nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestMigrasi902BentukTabelLogin(t *testing.T) {
	p := baca902(t, "902_m_login_go.sql")
	if len(p) != 3 {
		t.Fatalf("902 terbaca %d pernyataan, mau 3 (dua tabel, satu index)", len(p))
	}
	tabel, kolom := migrasi.KolomCreateTable(p[0])
	if tabel != "M_LOGIN_GO" || !reflect.DeepEqual(kolom, []string{
		"LOGIN_ID", "NAME", "PASSWORD_HASH", "ORGANIZATION_CODE", "DIVISION_CODE", "UNIT_CODE",
		"IS_ACTIVE", "FAILED_COUNT", "LOCKED_UNTIL", "MUST_CHANGE_PASSWORD", "SESSION_VERSION",
		"LAST_LOGIN", "TGL_CREATE", "TGL_UPDATE",
	}) {
		t.Errorf("M_LOGIN_GO: %s %v", tabel, kolom)
	}
	for _, w := range []string{
		"LOGIN_ID             VARCHAR2(64) NOT NULL",
		"PASSWORD_HASH        VARCHAR2(255) NOT NULL",
		// CODE, bukan ID, dari master organisasi (keputusan work owner).
		"ORGANIZATION_CODE    VARCHAR2(20),",
		"DIVISION_CODE        VARCHAR2(20),",
		"UNIT_CODE            VARCHAR2(20),",
		"IS_ACTIVE            VARCHAR2(1) DEFAULT '1' NOT NULL",
		"MUST_CHANGE_PASSWORD VARCHAR2(1) DEFAULT '1' NOT NULL",
		"SESSION_VERSION      NUMBER(10) DEFAULT 1 NOT NULL",
		"CONSTRAINT PK_M_LOGIN_GO PRIMARY KEY (LOGIN_ID)",
	} {
		if !strings.Contains(p[0], w) {
			t.Errorf("M_LOGIN_GO tanpa %q", w)
		}
	}
	tabel, kolom = migrasi.KolomCreateTable(p[1])
	if tabel != "M_LOGIN_GO_WORKBASKET" || !reflect.DeepEqual(kolom, []string{"LOGIN_ID", "WORKBASKET_ID", "TGL_CREATE"}) {
		t.Errorf("M_LOGIN_GO_WORKBASKET: %s %v", tabel, kolom)
	}
	for _, w := range []string{
		"CONSTRAINT PK_M_LOGIN_GO_WORKBASKET PRIMARY KEY (LOGIN_ID, WORKBASKET_ID)",
		"FOREIGN KEY (LOGIN_ID) REFERENCES {skema}.M_LOGIN_GO (LOGIN_ID)",
	} {
		if !strings.Contains(p[1], w) {
			t.Errorf("M_LOGIN_GO_WORKBASKET tanpa %q", w)
		}
	}
	if p[2] != "CREATE INDEX {skema}.IX_M_LOGIN_GO_WB_WORKBASKET ON {skema}.M_LOGIN_GO_WORKBASKET (WORKBASKET_ID)" {
		t.Errorf("index workbasket: %q", p[2])
	}
	semua := strings.Join(p, "\n")
	// ⛔ Tanpa FK ke master yang tidak dibuat migrasi ini (keputusan work
	// owner): keberadaan CODE dan WORKBASKET_ID diperiksa aplikasi.
	for _, w := range []string{"M_UNIT", "M_DIVISION", "M_ORGANIZATION", "{skema}.M_WORKBASKET", "ON DELETE", "COMMIT"} {
		if strings.Contains(semua, w) {
			t.Errorf("902 memuat %q", w)
		}
	}
}

func TestMigrasi902Mundur(t *testing.T) {
	p := baca902(t, "902_m_login_go_down.sql")
	mau := []string{
		"DROP TABLE {skema}.M_LOGIN_GO_WORKBASKET CASCADE CONSTRAINTS",
		"DROP TABLE {skema}.M_LOGIN_GO CASCADE CONSTRAINTS",
	}
	if !reflect.DeepEqual(p, mau) {
		t.Errorf("902_down:\n dapat %q\n mau   %q", p, mau)
	}
}
