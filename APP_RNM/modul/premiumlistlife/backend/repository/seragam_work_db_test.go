//go:build db

package repository_test

// Migrasi 059 terhadap skema uji Oracle: katalog = DDL, FK COVER_KEY tanpa
// ON DELETE, dan pengisian baris lama dari T_PREMIUM_LIST (brief seragam
// kolom 01-10-2026, §3 langkah 4).
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan.

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func TestMigrasi059KatalogDanPengisianBarisLama(t *testing.T) {
	// Pintu skema uji bersama - bantu_skemauji_db_test.go.
	sqlDB, skema, ctx := pasangSkemaUji(t)

	// Katalog = DDL: nama dan tipe setiap kolom T_WORK_POLIS sesudah 059.
	baris, err := sqlDB.QueryContext(ctx, `SELECT COLUMN_NAME, DATA_TYPE, CHAR_LENGTH FROM SYS.ALL_TAB_COLUMNS
		WHERE OWNER = UPPER(:1) AND TABLE_NAME = 'T_WORK_POLIS'`, skema)
	if err != nil {
		t.Fatal(err)
	}
	katalog := map[string]string{}
	for baris.Next() {
		var nama, tipe string
		var panjang int
		if err := baris.Scan(&nama, &tipe, &panjang); err != nil {
			t.Fatal(err)
		}
		if tipe == "VARCHAR2" {
			tipe = "VARCHAR2(" + strconv.Itoa(panjang) + ")"
		}
		katalog[nama] = tipe
	}
	if err := baris.Err(); err != nil {
		t.Fatal(err)
	}
	_ = baris.Close()
	ddl := bentukAkhirTabelKerja(t)["T_WORK_POLIS"]
	if len(katalog) != len(ddl) {
		t.Errorf("katalog %d kolom, DDL %d: %v lawan %v", len(katalog), len(ddl), katalog, ddl)
	}
	for k, tipe := range ddl {
		if katalog[k] != tipe {
			t.Errorf("%s: katalog %q, DDL %q", k, katalog[k], tipe)
		}
	}

	// FK COVER_KEY menunjuk dirinya sendiri, aturan hapus bawaan (NO ACTION).
	var aturan string
	if err := sqlDB.QueryRowContext(ctx, `SELECT DELETE_RULE FROM SYS.ALL_CONSTRAINTS
		WHERE OWNER = UPPER(:1) AND CONSTRAINT_NAME = 'FK_WORK_POLIS_COVER_KEY'`, skema).Scan(&aturan); err != nil {
		t.Fatalf("FK_WORK_POLIS_COVER_KEY: %v", err)
	}
	if aturan != "NO ACTION" {
		t.Errorf("DELETE_RULE %q, mau NO ACTION (tanpa ON DELETE, sama dengan klaim butir d)", aturan)
	}

	// Pengisian baris lama: satu pasangan diisi, dua pasangan tidak, nama yang
	// tidak muat VARCHAR2(128) dibiarkan kosong, baris yang sudah terisi tidak
	// ditimpa - dan pernyataan yang sama aman diulang.
	panjang := strings.Repeat("N", 129)
	for _, q := range []string{
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID) VALUES ('UJI-059-SATU')`,
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID) VALUES ('UJI-059-DUA')`,
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID) VALUES ('UJI-059-PANJANG')`,
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID, CREATE_OP_NAME) VALUES ('UJI-059-ISI', 'UJI-AKUN-BARU')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, CREATE_OP_NAME)
		   VALUES ('UJI-059-SATU', 'UJI-059-SATU', DATE '2026-01-02', 'UJI-AKUN-1')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT) VALUES ('UJI-059-DUA', 'UJI-059-DUA', DATE '2026-01-03')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT) VALUES ('UJI-059-DUA-B', 'UJI-059-DUA', DATE '2026-01-04')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, CREATE_OP_NAME)
		   VALUES ('UJI-059-PANJANG', 'UJI-059-PANJANG', DATE '2026-01-05', '` + panjang + `')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, CREATE_OP_NAME)
		   VALUES ('UJI-059-ISI', 'UJI-059-ISI', DATE '2026-01-06', 'UJI-AKUN-LAMA')`,
	} {
		if _, err := sqlDB.ExecContext(ctx, q); err != nil {
			t.Fatalf("baris tiruan: %v", err)
		}
	}
	p, err := migrasi.PernyataanLangkah(os.DirFS(".."), berkas059)
	if err != nil {
		t.Fatal(err)
	}
	isi := strings.ReplaceAll(p[len(p)-1], "{skema}", skema)
	for ulang := 0; ulang < 2; ulang++ {
		if _, err := sqlDB.ExecContext(ctx, isi); err != nil {
			t.Fatalf("pengisian baris lama (ke-%d): %v", ulang+1, err)
		}
	}
	for _, k := range []struct {
		id, tgl, nama string
	}{
		{"UJI-059-SATU", "2026-01-02", "UJI-AKUN-1"},
		{"UJI-059-DUA", "", ""},
		{"UJI-059-PANJANG", "2026-01-05", ""},
		{"UJI-059-ISI", "", "UJI-AKUN-BARU"},
	} {
		var tgl, nama, op sql.NullString
		if err := sqlDB.QueryRowContext(ctx, `SELECT TO_CHAR(TGL_CREATE, 'YYYY-MM-DD'), CREATE_OP_NAME, CREATE_OP
			FROM `+skema+`.T_WORK_POLIS WHERE ID = :1`, k.id).Scan(&tgl, &nama, &op); err != nil {
			t.Fatal(err)
		}
		if tgl.String != k.tgl || nama.String != k.nama || op.Valid {
			t.Errorf("%s: TGL_CREATE %q CREATE_OP_NAME %q CREATE_OP %v, mau %q %q kosong",
				k.id, tgl.String, nama.String, op, k.tgl, k.nama)
		}
	}
}
