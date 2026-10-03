package login

// Migrasi 902 - M_LOGIN_GO dan M_LOGIN_GO_WORKBASKET. TANPA Oracle.

import (
	"os"
	"reflect"
	"strconv"
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

// Migrasi 903 - M_LOGIN_GO_MENU (Kelola User, keputusan work owner
// 01-10-2026): menu per akun, lalu SETIAP akun yang sudah ada mendapat SEMUA
// menu - termasuk Kelola User - supaya tak seorang pun kehilangan layarnya
// dan selalu ada admin. Isi daftarnya dijaga penjaga menu (`penjaga`).
func TestMigrasi903MenuPerAkun(t *testing.T) {
	p := baca902(t, "903_m_login_go_menu.sql")
	if len(p) != 2 {
		t.Fatalf("903 terbaca %d pernyataan, mau 2 (tabel, isi awal)", len(p))
	}
	tabel, kolom := migrasi.KolomCreateTable(p[0])
	if tabel != "M_LOGIN_GO_MENU" || !reflect.DeepEqual(kolom, []string{"LOGIN_ID", "MENU_KODE", "TGL_CREATE"}) {
		t.Errorf("M_LOGIN_GO_MENU: %s %v", tabel, kolom)
	}
	for _, w := range []string{
		"LOGIN_ID   VARCHAR2(64) NOT NULL",
		"MENU_KODE  VARCHAR2(64) NOT NULL",
		"CONSTRAINT PK_M_LOGIN_GO_MENU PRIMARY KEY (LOGIN_ID, MENU_KODE)",
		"CONSTRAINT FK_M_LOGIN_GO_MENU_LOGIN FOREIGN KEY (LOGIN_ID) REFERENCES {skema}.M_LOGIN_GO (LOGIN_ID)",
	} {
		if !strings.Contains(p[0], w) {
			t.Errorf("M_LOGIN_GO_MENU tanpa %q", w)
		}
	}
	if !strings.HasPrefix(p[1], "INSERT INTO {skema}.M_LOGIN_GO_MENU (LOGIN_ID, MENU_KODE)\nSELECT l.LOGIN_ID, k.KODE\n  FROM {skema}.M_LOGIN_GO l\n CROSS JOIN (") {
		t.Errorf("isi awal 903: %q", p[1])
	}
	semua := strings.Join(p, "\n")
	for _, w := range []string{"ON DELETE", "COMMIT"} {
		if strings.Contains(semua, w) {
			t.Errorf("903 memuat %q", w)
		}
	}
}

func TestMigrasi903Mundur(t *testing.T) {
	p := baca902(t, "903_m_login_go_menu_down.sql")
	if mau := []string{"DROP TABLE {skema}.M_LOGIN_GO_MENU CASCADE CONSTRAINTS"}; !reflect.DeepEqual(p, mau) {
		t.Errorf("903_down:\n dapat %q\n mau   %q", p, mau)
	}
}

// Migrasi 904 - kolom kontak M_LOGIN_GO (Kelola User 03-10-2026): empat kolom berbahasa Inggris, semuanya
// NULLABLE (akun lama tidak punya isinya), lebar = batas `PeriksaKontak`.
func TestMigrasi904KolomKontak(t *testing.T) {
	p := baca902(t, "904_m_login_go_kontak.sql")
	if len(p) != 1 {
		t.Fatalf("904 terbaca %d pernyataan, mau 1", len(p))
	}
	satu := strings.Join(strings.Fields(p[0]), " ")
	mau := "ALTER TABLE {skema}.M_LOGIN_GO ADD ( EMAIL VARCHAR2(254), PHONE_NUMBER VARCHAR2(30), EMPLOYEE_ID VARCHAR2(30), JOB_POSITION VARCHAR2(150) )"
	if satu != mau {
		t.Errorf("904:\n dapat %q\n mau   %q", satu, mau)
	}
	if strings.Contains(strings.ToUpper(satu), "NOT NULL") {
		t.Error("kolom kontak harus NULLABLE - akun yang sudah ada tidak punya isinya")
	}
	for _, k := range []struct {
		kolom string
		lebar int
	}{{"EMAIL", PanjangMaksEmail}, {"PHONE_NUMBER", PanjangMaksTelepon}, {"EMPLOYEE_ID", PanjangMaksNIK}, {"JOB_POSITION", PanjangMaksJabatan}} {
		if !strings.Contains(satu, k.kolom+" VARCHAR2("+strconv.Itoa(k.lebar)+")") {
			t.Errorf("lebar kolom %s harus %d (batas PeriksaKontak)", k.kolom, k.lebar)
		}
	}
}

func TestMigrasi904Mundur(t *testing.T) {
	p := baca902(t, "904_m_login_go_kontak_down.sql")
	if mau := []string{"ALTER TABLE {skema}.M_LOGIN_GO DROP (EMAIL, PHONE_NUMBER, EMPLOYEE_ID, JOB_POSITION)"}; !reflect.DeepEqual(p, mau) {
		t.Errorf("904_down:\n dapat %q\n mau   %q", p, mau)
	}
}
