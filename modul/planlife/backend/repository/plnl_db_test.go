//go:build db

package repository_test

// Seam repository Plan terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`).
// Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 946 ditiru (fakta WO 08-10-2026): M_PRODUCT_TYPE_LIFE (ID VARCHAR2(6)
// NULLABLE TANPA PK + JSONDATA IS JSON, kunci huruf campuran) + view PRODUCT_TYPE_LIFE (teks DEV) + sequence; master
// BUSINESS (GROUPPANEL) dan BENEFIT_LIFE ditiru sebagai tabel UJI. Fixture UJI-.
//
// ⚠️ Koneksi lewat `skemauji.BukaRepositori()`.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/planlife/backend/models"
	"nusantarare/modul/planlife/backend/repository"
	"nusantarare/modul/planlife/backend/services"
	"nusantarare/uji/skemauji"
)

type ujiDB struct {
	repo  *db.DB
	skema string
	ctx   context.Context
}

func ddl(t *testing.T, kunci string, mundur bool, skema string) []string {
	t.Helper()
	sumber := fstest.MapFS{}
	for _, n := range []string{kunci + ".sql", kunci + "_down.sql"} {
		isi, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", n))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+n] = &fstest.MapFile{Data: isi}
	}
	l, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(l) != 1 {
		t.Fatalf("membaca %s: %v", kunci, err)
	}
	var out []string
	for _, q := range l[0].Pernyataan {
		out = append(out, strings.ReplaceAll(q, "{skema}", skema))
	}
	return out
}

func json(id, cover, biz, bizID, ben, benID string) string {
	return `{"ID":"` + id + `","CoverName":"` + cover + `","Business":"` + biz + `","BusinessID":"` + bizID +
		`","Benefit":"` + ben + `","BenefitID":"` + benID + `","pxObjClass":"ASM-FW-GISFW-Int-PRODUCT_TYPE_LIFE","pxCreateOperator":"UJI-OPERATOR"}`
}

func pasang(t *testing.T) *ujiDB {
	t.Helper()
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := repo.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	u := &ujiDB{repo: repo, skema: repo.Skema(), ctx: ctx}
	u.bongkar(t)
	baris := func(id, j string) string {
		return `INSERT INTO {s}.M_PRODUCT_TYPE_LIFE (ID, JSONDATA) VALUES ('` + id + `', '` + j + `')`
	}
	for _, q := range []string{
		`CREATE TABLE {s}.M_PRODUCT_TYPE_LIFE (ID VARCHAR2(6), JSONDATA CLOB CONSTRAINT ENSURE_M_PRODUCT_TYPE_LIFE CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.PRODUCT_TYPE_LIFE AS SELECT a.JSONDATA.ID, a.JSONDATA.CoverName, a.JSONDATA.Business, a.JSONDATA.BusinessID,
			a.JSONDATA.Benefit, a.JSONDATA.BenefitID FROM {s}.M_PRODUCT_TYPE_LIFE a`,
		`CREATE SEQUENCE {s}.M_PRODUCT_TYPE_LIFE_SEQ START WITH 44`,
		`CREATE TABLE {s}.BUSINESS (ID VARCHAR2(10), OLDID VARCHAR2(10), NOTE VARCHAR2(200), GROUPPANEL VARCHAR2(10))`,
		`CREATE TABLE {s}.BENEFIT_LIFE (ID VARCHAR2(10) PRIMARY KEY, BENEFIT VARCHAR2(200))`,
		`INSERT INTO {s}.BUSINESS VALUES ('9001', 'L1', 'UJI KREDIT', '009')`,
		`INSERT INTO {s}.BUSINESS VALUES ('9101', 'F1', 'UJI KEBAKARAN', '001')`,
		`INSERT INTO {s}.BENEFIT_LIFE VALUES ('100001', 'UJI RAWAT INAP')`,
		baris("100001", json("100001", "UJI Plan A", "UJI KREDIT", "9001", "UJI RAWAT INAP", "100001")),
		baris("100002", json("100002", "UJI Plan B", "UJI KREDIT", "9001", "UJI RAWAT INAP", "100001")),
	} {
		u.exec(t, q)
	}
	t.Cleanup(func() {
		u.bongkar(t)
		_ = repo.Close()
	})
	return u
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for _, q := range []string{`DROP VIEW {s}.PRODUCT_TYPE_LIFE`, `DROP TABLE {s}.PRODUCT_TYPE_LIFE PURGE`,
		`DROP TABLE {s}.M_PRODUCT_TYPE_LIFE PURGE`, `DROP SEQUENCE {s}.M_PRODUCT_TYPE_LIFE_SEQ`, `DROP TABLE {s}.BUSINESS PURGE`,
		`DROP TABLE {s}.BENEFIT_LIFE PURGE`} {
		if _, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)); err != nil &&
			!strings.Contains(err.Error(), "ORA-00942") && !strings.Contains(err.Error(), "ORA-02289") {
			t.Fatalf("membongkar %s: %v", q, err)
		}
	}
}

func (u *ujiDB) exec(t *testing.T, q string) {
	t.Helper()
	if _, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)); err != nil {
		t.Fatalf("%.60s: %v", q, err)
	}
}

func (u *ujiDB) gagalORA(t *testing.T, q, kode string) {
	t.Helper()
	_, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema))
	if err == nil || !strings.Contains(err.Error(), kode) {
		t.Errorf("%.60s: mau %s, dapat %v", q, kode, err)
	}
}

func (u *ujiDB) cacah(t *testing.T, q string) int {
	t.Helper()
	var n int
	if err := u.repo.QueryRowContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)).Scan(&n); err != nil {
		t.Fatalf("%.60s: %v", q, err)
	}
	return n
}

func (u *ujiDB) periksa(t *testing.T, judul string, mau map[string]int) {
	t.Helper()
	for q, n := range mau {
		if got := u.cacah(t, q); got != n {
			t.Errorf("%s: %d, mau %d: %s", judul, got, n, q)
		}
	}
}

func (u *ujiDB) lari(t *testing.T, kunci string, mundur bool, sebagian int) {
	t.Helper()
	l := ddl(t, kunci, mundur, u.skema)
	for _, q := range l[:sebagian] {
		u.exec(t, q)
	}
	for _, q := range l {
		u.exec(t, q)
	}
}

// K2: kelima ekspresi 948 = view lama baris demi baris; ID tabel = JSONDATA.ID.
func TestDBEkspresiSamaDenganView(t *testing.T) {
	u := pasang(t)
	u.periksa(t, "ekspresi = view", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_PRODUCT_TYPE_LIFE m JOIN {s}.PRODUCT_TYPE_LIFE v ON v.ID = m.JSONDATA.ID
			WHERE m.ID = m.JSONDATA.ID AND m.JSONDATA.CoverName = v.COVERNAME AND m.JSONDATA.Business = v.BUSINESS
			AND m.JSONDATA.BusinessID = v.BUSINESSID AND m.JSONDATA.Benefit = v.BENEFIT AND m.JSONDATA.BenefitID = v.BENEFITID`: 2,
	})
}

// 946-948 maju (diulang sebagian): tabel, 6 kolom urutan view, PK (ID NOT NULL), isi huruf asli; pelindung ORA-00904;
// mundur 948..946 memulihkan nama, view, JSONDATA, ID NULLABLE.
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.lari(t, "946_product_type_life_ganti_nama", false, 1)
	u.lari(t, "947_product_type_life_kolom", false, 0)
	u.lari(t, "948_product_type_life_satu_tabel", false, 4)
	u.periksa(t, "sesudah 948", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME = 'M_PRODUCT_TYPE_LIFE'`:                                            0,
		`SELECT COUNT(*) FROM SYS.ALL_TABLES WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE'`:                                                1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE'`:                                           6,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'ID' AND NULLABLE = 'N'`: 1,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'PK_PRODUCT_TYPE_LIFE'`:                                   1,
		`SELECT COUNT(*) FROM {s}.PRODUCT_TYPE_LIFE WHERE ID = '100001' AND COVERNAME = 'UJI Plan A' AND BUSINESSID = '9001' AND BENEFITID = '100001'`:       1,
	})
	for _, k := range []string{"947_product_type_life_kolom", "946_product_type_life_ganti_nama"} {
		u.gagalORA(t, ddl(t, k, true, u.skema)[0], "ORA-00904")
	}
	u.lari(t, "948_product_type_life_satu_tabel", true, 1)
	u.lari(t, "947_product_type_life_kolom", true, 0)
	l := ddl(t, "946_product_type_life_ganti_nama", true, u.skema)
	u.exec(t, l[0])
	u.exec(t, l[1])
	for i, q := range l {
		if _, err := u.repo.ExecContext(u.ctx, q); err != nil && !(i == 0 && strings.Contains(err.Error(), "ORA-00942")) {
			t.Fatalf("946 mundur %d: %v", i, err)
		}
	}
	u.periksa(t, "sesudah mundur", map[string]int{
		`SELECT COUNT(*) FROM {s}.PRODUCT_TYPE_LIFE WHERE ID = '100002' AND COVERNAME = 'UJI Plan B' AND BENEFITID = '100001'`:                                 1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'ID' AND NULLABLE = 'Y'`: 1,
	})
}

// K1.4 / K2: ID NULL, ID ganda, ID <> JSONDATA.ID, isi berbeda -> pemeriksaan 948 GAGAL KERAS ORA-12899 sebelum PK /
// DROP; JSONDATA utuh.
func TestDBPemeriksaanBerhenti(t *testing.T) {
	u := pasang(t)
	u.lari(t, "946_product_type_life_ganti_nama", false, 0)
	u.lari(t, "947_product_type_life_kolom", false, 0)
	l := ddl(t, "948_product_type_life_satu_tabel", false, u.skema)
	u.exec(t, l[0])
	u.exec(t, `UPDATE {s}.PRODUCT_TYPE_LIFE SET ID = '100009' WHERE ID = '100002'`) // ID tabel <> JSONDATA.ID
	u.gagalORA(t, l[1], "ORA-12899")
	u.exec(t, `UPDATE {s}.PRODUCT_TYPE_LIFE SET ID = NULL WHERE ID = '100009'`)
	u.gagalORA(t, l[1], "ORA-12899")
	u.exec(t, `UPDATE {s}.PRODUCT_TYPE_LIFE SET ID = '100001' WHERE ID IS NULL`) // ganda
	u.gagalORA(t, l[1], "ORA-12899")
	u.exec(t, `UPDATE {s}.PRODUCT_TYPE_LIFE SET ID = '100002' WHERE ROWID = (SELECT MAX(ROWID) FROM {s}.PRODUCT_TYPE_LIFE)`)
	u.exec(t, `UPDATE {s}.PRODUCT_TYPE_LIFE SET BENEFITID = NULL WHERE ID = '100002'`)
	u.gagalORA(t, l[2], "ORA-12899")
	u.periksa(t, "JSONDATA utuh", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'PRODUCT_TYPE_LIFE' AND COLUMN_NAME = 'JSONDATA'`: 1,
	})
}

func layananOracle(u *ujiDB) *services.Layanan {
	return services.BaruLayanan(repository.Baru(u.repo), func(ctx context.Context, fn func(tx *db.Tx) error) error {
		tx, err := u.repo.Mulai(ctx)
		if err != nil {
			return err
		}
		if err := fn(tx); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

// Seam: pilihan grup 009; Add 100044 dengan nama / ID master; K4 422; K5 kembar; K3 409.
func TestDBLayanan(t *testing.T) {
	u := pasang(t)
	for _, k := range []string{"946_product_type_life_ganti_nama", "947_product_type_life_kolom", "948_product_type_life_satu_tabel"} {
		u.lari(t, k, false, 0)
	}
	l := layananOracle(u)
	if b, err := l.PilihanBusiness(u.ctx); err != nil || len(b) != 1 || b[0].OldID != "L1" {
		t.Errorf("pilihan %+v %v", b, err)
	}
	penuh := services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	p, err := l.Simpan(u.ctx, penuh, "", models.Isian{CoverName: "UJI Plan C", Business: "uji kredit", Benefit: "uji rawat inap"})
	if err != nil || p.ID != "100044" || p.Business != "UJI KREDIT" || p.BusinessID != "9001" || p.BenefitID != "100001" {
		t.Fatalf("add %+v %v", p, err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{CoverName: "UJI D", Business: "UJI KEBAKARAN", Benefit: "UJI RAWAT INAP"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("K4 %v", err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{CoverName: "uji plan a", Business: "UJI KREDIT", Benefit: "UJI RAWAT INAP"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("K5 %v", err)
	}
	u.exec(t, `INSERT INTO {s}.PRODUCT_TYPE_LIFE (ID, COVERNAME) VALUES ('100045', 'UJI DITULIS PEGA')`)
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{CoverName: "UJI E", Business: "UJI KREDIT", Benefit: "UJI RAWAT INAP"}); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("K3 %v", err)
	}
}
