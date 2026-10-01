//go:build db

package skemauji_test

// Login (M_LOGIN_GO, migrasi 902) diadu dengan Oracle.
//
// Yang hanya Oracle dapat buktikan: penguncian dan pembukaannya memakai jam
// basis data (SYSDATE), hitungan gagal mulai dari satu lagi sesudah kunci
// lewat, master organisasi dan workbasket terbaca lewat CODE, workbasket
// nonaktif bukan peran, dan ganti sandi mencabut sesi lama.
//
// ⚠️ M_ORGANIZATION, M_DIVISION, M_UNIT, M_WORKBASKET tidak dibuat migrasi
// aplikasi; di skema uji keempatnya TIRUAN berbentuk sama dengan DEV (katalog
// baca-saja 01-10-2026), dibuat dan dibuang uji ini.
//
// ⛔ Melewati bila Oracle belum dikonfigurasi. Melewati bukan lulus.

import (
	"errors"
	"testing"

	"nusantarare/inti/backend/login"
	"nusantarare/uji/skemauji"
)

func TestLoginDariOracle(t *testing.T) {
	sqlDB, skema, ctx := pasangSkemaInti(t)
	jalan := func(q string) {
		t.Helper()
		if _, err := sqlDB.ExecContext(ctx, q); err != nil {
			t.Fatalf("%s: %v", q, err)
		}
	}
	tiruan := []string{"M_UNIT", "M_DIVISION", "M_ORGANIZATION", "M_WORKBASKET"}
	t.Cleanup(func() {
		for _, n := range tiruan {
			_, _ = sqlDB.ExecContext(ctx, `DROP TABLE `+skema+`.`+n+` CASCADE CONSTRAINTS`)
		}
	})
	for _, q := range []string{
		`CREATE TABLE ` + skema + `.M_ORGANIZATION (ORGANIZATION_ID NUMBER(5) PRIMARY KEY, CODE VARCHAR2(20) UNIQUE, NAME VARCHAR2(150), IS_ACTIVE NUMBER(1))`,
		`CREATE TABLE ` + skema + `.M_DIVISION (DIVISION_ID NUMBER(10) PRIMARY KEY, ORGANIZATION_ID NUMBER(5), CODE VARCHAR2(20) UNIQUE, NAME VARCHAR2(150), SORT_ORDER NUMBER(5), IS_ACTIVE NUMBER(1))`,
		`CREATE TABLE ` + skema + `.M_UNIT (UNIT_ID NUMBER(10) PRIMARY KEY, DIVISION_ID NUMBER(10), CODE VARCHAR2(20) UNIQUE, NAME VARCHAR2(150), SORT_ORDER NUMBER(5), IS_ACTIVE NUMBER(1))`,
		`CREATE TABLE ` + skema + `.M_WORKBASKET (WORKBASKET_ID VARCHAR2(64) PRIMARY KEY, NAME VARCHAR2(150), IS_ACTIVE NUMBER(1))`,
		`INSERT INTO ` + skema + `.M_ORGANIZATION VALUES (1, 'UJI-ORG', 'Uji Organisasi', 1)`,
		`INSERT INTO ` + skema + `.M_DIVISION VALUES (11, 1, 'UJI-DIV', 'Uji Divisi', 1, 1)`,
		`INSERT INTO ` + skema + `.M_UNIT VALUES (111, 11, 'UJI-UNIT', 'Uji Unit', 1, 1)`,
		`INSERT INTO ` + skema + `.M_WORKBASKET VALUES ('ReasLifeAdmin', 'Uji Admin', 1)`,
		`INSERT INTO ` + skema + `.M_WORKBASKET VALUES ('ReasLifeSPV', 'Uji SPV', 1)`,
	} {
		jalan(q)
	}
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = repo.Close() }()
	l := login.NewLayanan(login.NewGudangOracle(repo), []byte("rahasia-uji-yang-panjangnya-32-byte!"))

	sandi, err := l.BuatPengguna(ctx, login.AkunBaru{ID: "UJI-LOGIN-1", Nama: "Uji Login",
		Organisasi: "UJI-ORG", Divisi: "UJI-DIV", Unit: "UJI-UNIT", Workbasket: []string{"ReasLifeAdmin", "ReasLifeSPV"}})
	if err != nil {
		t.Fatalf("buat pengguna: %v", err)
	}
	if _, err := l.BuatPengguna(ctx, login.AkunBaru{ID: "UJI-LOGIN-2", Nama: "Uji",
		Organisasi: "UJI-ORG", Divisi: "UJI-DIV", Unit: "UJI-TIDAK-ADA"}); !errors.Is(err, login.ErrMasterTidakAda) {
		t.Errorf("unit tidak ada: %v", err)
	}
	jalan(`UPDATE ` + skema + `.M_WORKBASKET SET IS_ACTIVE = 0 WHERE WORKBASKET_ID = 'ReasLifeSPV'`)
	p, tok, err := l.Masuk(ctx, "UJI-LOGIN-1", sandi)
	if err != nil || !p.WajibGantiSandi || len(p.Peran) != 1 || p.Peran[0] != "ReasLifeAdmin" || p.Unit != "UJI-UNIT" {
		t.Fatalf("masuk pertama: %+v %v", p, err)
	}

	// Lima sandi salah -> terkunci; sandi BENAR pun ditolak.
	for i := 0; i < login.BatasGagal; i++ {
		if _, _, err := l.Masuk(ctx, "UJI-LOGIN-1", "Sandi-Salah-000"); !errors.Is(err, login.ErrKredensial) {
			t.Fatalf("salah ke-%d: %v", i+1, err)
		}
	}
	if _, _, err := l.Masuk(ctx, "UJI-LOGIN-1", sandi); !errors.Is(err, login.ErrTerkunci) {
		t.Fatalf("sesudah 5 salah: %v", err)
	}
	// Kunci lewat (jam Oracle): satu salah lagi memulai hitungan dari 1.
	jalan(`UPDATE ` + skema + `.M_LOGIN_GO SET LOCKED_UNTIL = SYSDATE - 1/1440 WHERE LOGIN_ID = 'UJI-LOGIN-1'`)
	if _, _, err := l.Masuk(ctx, "UJI-LOGIN-1", "Sandi-Salah-000"); !errors.Is(err, login.ErrKredensial) {
		t.Fatalf("salah sesudah kunci lewat: %v", err)
	}
	var gagal int
	if err := sqlDB.QueryRowContext(ctx, `SELECT FAILED_COUNT FROM `+skema+`.M_LOGIN_GO WHERE LOGIN_ID = 'UJI-LOGIN-1'`).Scan(&gagal); err != nil || gagal != 1 {
		t.Errorf("FAILED_COUNT sesudah kunci lewat = %d (%v), mau 1", gagal, err)
	}
	_, tok, err = l.Masuk(ctx, "UJI-LOGIN-1", sandi)
	if err != nil {
		t.Fatalf("masuk sesudah kunci lewat: %v", err)
	}

	// Ganti sandi mencabut cookie lama.
	lama := login.Tandatangani([]byte("rahasia-uji-yang-panjangnya-32-byte!"), tok)
	if _, _, err := l.GantiSandi(ctx, tok, sandi, "Sandi-Baru-Uji-01"); err != nil {
		t.Fatalf("ganti sandi: %v", err)
	}
	if _, _, err := l.Sesi(ctx, lama); !errors.Is(err, login.ErrSesiTidakSah) {
		t.Errorf("cookie lama sesudah ganti sandi: %v", err)
	}
	p, _, err = l.Masuk(ctx, "UJI-LOGIN-1", "Sandi-Baru-Uji-01")
	if err != nil || p.WajibGantiSandi {
		t.Errorf("masuk dengan sandi baru: %+v %v", p, err)
	}
}
