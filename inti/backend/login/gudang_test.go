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
		"sisip":      sqlSisipAkun("S.M_LOGIN_GO", "S.M_LOGIN_GO_CONTACT_SEQ"),
		"sisip wb":   sqlSisipWorkbasket("S.M_LOGIN_GO_WORKBASKET"),
		"menu":       sqlMenu("S.M_LOGIN_GO_MENU"),
		"wb semua":   sqlWorkbasketSemua("S.M_LOGIN_GO_WORKBASKET"),
		"sisip menu": sqlSisipMenu("S.M_LOGIN_GO_MENU"),
		// Identitas (migrasi 905, 03-10-2026).
		"pemakai username": sqlPemakaiUsername("S.M_LOGIN_GO"),
		"pemakai email":    sqlPemakaiEmail("S.M_LOGIN_GO"),
		// Kelola User (01-10-2026).
		"ringkas":      sqlDaftarAkun("S.M_LOGIN_GO", true),
		"kunci admin":  sqlKunciAdmin("S.M_LOGIN_GO", "S.M_LOGIN_GO_MENU"),
		"hitung admin": sqlHitungAdmin("S.M_LOGIN_GO", "S.M_LOGIN_GO_MENU"),
		"ubah profil":  sqlUbahProfil("S.M_LOGIN_GO"),
		"setel aktif":  sqlSetelAktif("S.M_LOGIN_GO"),
		"buka kunci":   sqlBukaKunci("S.M_LOGIN_GO"),
		"atur sandi":   sqlAturSandi("S.M_LOGIN_GO"),
		"wajib ganti":  sqlSetelWajibGanti("S.M_LOGIN_GO"),
		"hapus milik":  sqlHapusMilik("S.M_LOGIN_GO_MENU"),
		"m org":        sqlMasterOrganisasi("S.M_ORGANIZATION"),
		"m divisi":     sqlMasterDivisi("S.M_DIVISION", "S.M_ORGANIZATION"),
		"m unit":       sqlMasterUnit("S.M_UNIT", "S.M_DIVISION"),
		"m wb":         sqlMasterWorkbasket("S.M_WORKBASKET"),
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

// Daftar akun tanpa bind (seluruh baris) - tetap bersih dari literal dan COMMIT,
// dan nol kolom rahasia.
func TestSQLDaftarAkunTanpaRahasia(t *testing.T) {
	q := sqlDaftarAkun("S.M_LOGIN_GO", false)
	if err := db.PeriksaSQL(q); err != nil {
		t.Error(err)
	}
	for _, w := range []string{"PASSWORD_HASH", "SESSION_VERSION", "'"} {
		if strings.Contains(q, w) {
			t.Errorf("daftar akun memuat %q:\n%s", w, q)
		}
	}
}

// Penjaga admin terakhir: kunci baris admin aktif lebih dulu (FOR UPDATE),
// hitung sesudah perubahan - keduanya atas bendera aktif dan KODE Kelola User.
func TestSQLAdminTerakhirMengunciLaluMenghitung(t *testing.T) {
	k := sqlKunciAdmin("S.M_LOGIN_GO", "S.M_LOGIN_GO_MENU")
	h := sqlHitungAdmin("S.M_LOGIN_GO", "S.M_LOGIN_GO_MENU")
	syarat := "WHERE l.IS_ACTIVE = :1 AND EXISTS (SELECT 1 FROM S.M_LOGIN_GO_MENU m WHERE m.LOGIN_ID = l.LOGIN_ID AND m.MENU_KODE = :2)"
	if !strings.Contains(k, syarat) || !strings.Contains(k, "ORDER BY l.LOGIN_ID") || !strings.HasSuffix(strings.TrimSpace(k), "FOR UPDATE") {
		t.Errorf("kunci admin:\n%s", k)
	}
	if !strings.Contains(h, syarat) || !strings.HasPrefix(h, "SELECT COUNT(*) FROM S.M_LOGIN_GO l") {
		t.Errorf("hitung admin:\n%s", h)
	}
}
