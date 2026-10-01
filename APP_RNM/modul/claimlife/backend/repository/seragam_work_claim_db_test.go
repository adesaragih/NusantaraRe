//go:build db

package repository_test

// Migrasi 023 terhadap skema uji Oracle: katalog = DDL, dan pengaman CASE_ID
// menolak baris yang CASE_ID-nya berbeda dari ID (keputusan work owner
// 01-10-2026).
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan.

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func TestMigrasi023KatalogDanPengamanCaseID(t *testing.T) {
	// Pintu skema uji yang sama dengan uji pohon klaim - pohonklaim_db_test.go.
	d, _, bersih := siapkanPohon(t)
	defer bersih()
	ctx := context.Background()
	skema := d.Skema()

	katalog := func(tabel string) map[string]string {
		t.Helper()
		baris, err := d.QueryContext(ctx, `SELECT COLUMN_NAME, DATA_TYPE, CHAR_LENGTH FROM SYS.ALL_TAB_COLUMNS
			WHERE OWNER = UPPER(:1) AND TABLE_NAME = :2`, skema, tabel)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = baris.Close() }()
		out := map[string]string{}
		for baris.Next() {
			var nama, tipe string
			var panjang int
			if err := baris.Scan(&nama, &tipe, &panjang); err != nil {
				t.Fatal(err)
			}
			if tipe == "VARCHAR2" {
				tipe = "VARCHAR2(" + strconv.Itoa(panjang) + ")"
			}
			out[nama] = tipe
		}
		if err := baris.Err(); err != nil {
			t.Fatal(err)
		}
		return out
	}
	ddl := bentukAkhirKlaim(t, "T_WORK_CLAIM", "T_GENERAL_CLAIM")
	work := katalog("T_WORK_CLAIM")
	if len(work) != len(ddl["T_WORK_CLAIM"]) {
		t.Errorf("T_WORK_CLAIM katalog %v, DDL %v", work, ddl["T_WORK_CLAIM"])
	}
	for k, tipe := range ddl["T_WORK_CLAIM"] {
		if work[k] != tipe {
			t.Errorf("T_WORK_CLAIM.%s: katalog %q, DDL %q", k, work[k], tipe)
		}
	}
	if got := katalog("T_GENERAL_CLAIM")["TYPE"]; got != "VARCHAR2(32)" {
		t.Errorf("T_GENERAL_CLAIM.TYPE di katalog: %q", got)
	}

	// Pengaman: CASE_ID dikembalikan sementara dengan satu baris berselisih.
	p, err := migrasi.PernyataanLangkah(os.DirFS(".."), berkas023)
	if err != nil {
		t.Fatal(err)
	}
	jalan := func(q string) error {
		_, err := d.ExecContext(ctx, strings.ReplaceAll(q, "{skema}", skema))
		return err
	}
	for _, q := range []string{
		`ALTER TABLE {skema}.T_WORK_CLAIM ADD (CASE_ID VARCHAR2(64))`,
		`INSERT INTO {skema}.T_WORK_CLAIM (ID, LINI, CASE_ID) VALUES ('CLMLF-UJI-023', 'LIFE', 'UJI-CASE-LAIN')`,
	} {
		if err := jalan(q); err != nil {
			t.Fatalf("persiapan: %v", err)
		}
	}
	if err := jalan(p[0]); err == nil || !strings.Contains(err.Error(), "ORA-02293") {
		t.Fatalf("pengaman CASE_ID tidak menolak baris berselisih: %v", err)
	}
	// Disamakan: pengaman lolos, lalu kolomnya terbuang beserta CHECK-nya.
	if err := jalan(`UPDATE {skema}.T_WORK_CLAIM SET CASE_ID = ID WHERE ID = 'CLMLF-UJI-023'`); err != nil {
		t.Fatal(err)
	}
	for _, i := range []int{0, 0, 7} { // pengaman dua kali: aman diulang
		if err := jalan(p[i]); err != nil {
			t.Fatalf("pernyataan %d: %v", i+1, err)
		}
	}
	if _, ada := katalog("T_WORK_CLAIM")["CASE_ID"]; ada {
		t.Error("CASE_ID masih ada sesudah pernyataan 8")
	}
}
