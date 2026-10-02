//go:build db

package skemauji_test

// M_NAV_MENU diadu dengan Oracle - migrasi 900 + 901 (menu datar, keputusan
// work owner 30-09-2026).
//
// Yang hanya Oracle dapat buktikan: blok berpelindung katalog 901 benar-benar
// membuang lima butir, kunci tamu, indeks, dan kolom PARENT_ID; mengulangnya
// (pelari yang gagal di tengah) tidak galat; jalur mundurnya mengembalikan
// semuanya dan juga aman diulang; dan CHECK GROUPMENU tetap menolak golongan
// lain.
//
// ⛔ Melewati bila Oracle belum dikonfigurasi. Melewati bukan lulus.

import (
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/uji/skemauji"
)

// Menggantikan `TestMenuIsiAwalDariOracleIdempotenDanBerCheck` (20 kelompok +
// 5 butir, INSERT 900 diulang) - brief menu datar 30-09-2026.
func TestMenuDatarDariOracleIdempotenDanBerCheck(t *testing.T) {
	// Pintu skema uji inti bersama - bantu_inti_db_test.go.
	sqlDB, skema, ctx := pasangSkemaInti(t)
	var err error

	satu := func(q string, arg ...any) int {
		t.Helper()
		var n int
		if err := sqlDB.QueryRowContext(ctx, q, arg...).Scan(&n); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
		return n
	}
	bentuk := func() (baris, parent, fk, ix int) {
		t.Helper()
		baris = satu(`SELECT COUNT(*) FROM ` + skema + `.M_NAV_MENU`)
		parent = satu(`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER(:1)
			AND TABLE_NAME = 'M_NAV_MENU' AND COLUMN_NAME = 'PARENT_ID'`, skema)
		fk = satu(`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER(:1)
			AND TABLE_NAME = 'M_NAV_MENU' AND CONSTRAINT_NAME = 'FK_M_NAV_MENU_INDUK'`, skema)
		ix = satu(`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER(:1)
			AND TABLE_NAME = 'M_NAV_MENU' AND INDEX_NAME = 'IX_M_NAV_MENU_PARENT'`, skema)
		return
	}
	mau := func(tahap string, baris, sisa int) {
		t.Helper()
		b, p, f, i := bentuk()
		if b != baris || p != sisa || f != sisa || i != sisa {
			t.Errorf("%s: %d baris, PARENT_ID %d, FK %d, indeks %d - mau %d baris dan ketiganya %d",
				tahap, b, p, f, i, baris, sisa)
		}
	}
	jalankan := func(berkas string) {
		t.Helper()
		pernyataan, err := migrasi.PernyataanLangkah(inti.SumberMigrasi(), berkas)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range pernyataan {
			if _, err := sqlDB.ExecContext(ctx, strings.ReplaceAll(p, "{skema}", skema)); err != nil {
				t.Fatalf("%s: %v", berkas, err)
			}
		}
	}

	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = repo.Close() }()
	// Pembaca backend (paket 2): 20 baris modul, empat golongan, satu tingkat -
	// SEBELUM dan SESUDAH 901 (lihat jalur mundur di bawah).
	pembaca := func(tahap string) {
		t.Helper()
		baris, err := menu.NewPembaca(repo).Baca(ctx)
		if err != nil {
			t.Fatalf("%s: pembaca menu: %v", tahap, err)
		}
		if len(baris) != 20 {
			t.Errorf("%s: pembaca membaca %d baris, mau 20 baris modul", tahap, len(baris))
		}
		m := menu.Susun(baris, []string{"claimlife", "premiumlistlife", "komiteclaimlife", "treatycontractout"})
		n := 0
		for _, g := range m.Golongan {
			n += len(g.Modul)
		}
		if len(m.Golongan) != 4 || n != 20 {
			t.Errorf("%s: %d golongan, %d modul - mau 4 dan 20", tahap, len(m.Golongan), n)
		}
	}

	// Pasang menjalankan 900 lalu 901: 20 baris, satu per modul, datar.
	mau("sesudah Pasang (900 + 901)", 20, 0)
	pembaca("sesudah 901")
	if n := satu(`SELECT COUNT(*) FROM ` + skema + `.M_NAV_MENU WHERE KODE <> MODUL`); n != 0 {
		t.Errorf("%d baris ber-KODE bukan nama modulnya", n)
	}

	// 901 diulang - keadaan pelari yang gagal di tengah lalu mengulang.
	jalankan("901_m_nav_menu_datar.sql")
	mau("sesudah 901 diulang", 20, 0)

	// Jalur mundur: kolom, FK, indeks, dan lima butir kembali; diulang aman.
	jalankan("901_m_nav_menu_datar_down.sql")
	mau("sesudah 901 mundur", 25, 1)
	// Keadaan DEV sebelum work owner menjalankan 901: PARENT_ID dan lima butir
	// masih ada - pembaca baru tetap benar (KODE = MODUL).
	pembaca("sebelum 901 (jalur mundur)")
	if n := satu(`SELECT COUNT(*) FROM ` + skema + `.M_NAV_MENU WHERE PARENT_ID IS NOT NULL
		AND KODE IN ('inbox', 'register', 'premiumlist', 'komite', 'tco-tahun')`); n != 5 {
		t.Errorf("jalur mundur mengembalikan %d butir, mau 5", n)
	}
	jalankan("901_m_nav_menu_datar_down.sql")
	mau("sesudah 901 mundur diulang", 25, 1)
	jalankan("901_m_nav_menu_datar.sql")
	mau("sesudah 901 lagi", 20, 0)

	// CHECK GROUPMENU menolak golongan di luar keempatnya.
	_, err = sqlDB.ExecContext(ctx, `INSERT INTO `+skema+`.M_NAV_MENU (ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN) `+
		`VALUES (`+skema+`.SEQ_M_NAV_MENU.NEXTVAL, 'uji-lain', 'Uji', 'LAIN', 'uji', 1)`)
	if err == nil || !strings.Contains(err.Error(), "ORA-02290") {
		t.Errorf("GROUPMENU 'LAIN' diterima atau galatnya bukan ORA-02290: %v", err)
	}
}
