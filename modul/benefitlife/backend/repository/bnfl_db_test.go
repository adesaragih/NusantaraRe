//go:build db

package repository_test

// Seam repository Benefit terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`).
// Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 942 ditiru (fakta WO 08-10-2026): M_BENEFIT_LIFE (ID VARCHAR2(10) PK +
// JSONDATA IS JSON, kunci `Benefit` huruf campuran, `Number` di sebagian baris) + view BENEFIT_LIFE
// (`a.JSONDATA.Benefit`) + M_BENEFIT_LIFE_SEQ. Migrasi inti 942-944 dijalankan dari berkasnya. Fixture UJI-.
//
// ⚠️ Koneksi lewat `skemauji.BukaRepositori()`, BUKAN `skemauji.Buka()` (penjaga Claim Life
// `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci cacah pemanggilnya).

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
	"nusantarare/modul/benefitlife/backend/models"
	"nusantarare/modul/benefitlife/backend/repository"
	"nusantarare/modul/benefitlife/backend/services"
	"nusantarare/uji/skemauji"
)

type ujiDB struct {
	repo  *db.DB
	skema string
	ctx   context.Context
}

// ddl - pernyataan migrasi inti `kunci` (maju atau mundur), `{skema}` diganti skema uji.
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

// pasang - keadaan DEV sebelum 942.
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
	baris := func(id, json string) string {
		return `INSERT INTO {s}.M_BENEFIT_LIFE (ID, JSONDATA) VALUES ('` + id + `', '` + json + `')`
	}
	for _, q := range []string{
		`CREATE TABLE {s}.M_BENEFIT_LIFE (ID VARCHAR2(10) NOT NULL PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_BENEFIT_LIFE_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.BENEFIT_LIFE AS SELECT a.ID, a.JSONDATA.Benefit FROM {s}.M_BENEFIT_LIFE a`,
		`CREATE SEQUENCE {s}.M_BENEFIT_LIFE_SEQ START WITH 12`,
		baris("100001", `{"Benefit":"UJI RAWAT INAP","Number":"100001","pxObjClass":"ASM-FW-GISFW-Int-BENEFIT_LIFE","pxCreateOperator":"UJI-OPERATOR"}`),
		baris("100002", `{"Benefit":"UJI Meninggal Dunia","pxObjClass":"ASM-FW-GISFW-Int-BENEFIT_LIFE"}`),
		baris("100004", `{"Benefit":"UJI CACAT TETAP","Number":"100004"}`),
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
	for _, q := range []string{`DROP VIEW {s}.BENEFIT_LIFE`, `DROP TABLE {s}.BENEFIT_LIFE PURGE`, `DROP TABLE {s}.M_BENEFIT_LIFE PURGE`,
		`DROP SEQUENCE {s}.M_BENEFIT_LIFE_SEQ`} {
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

// gagalORA - pernyataan `q` HARUS gagal dengan kode Oracle `kode`.
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

// lari - langkah `kunci`: `sebagian` pernyataan pertama dulu (pelari gagal di tengah), lalu penuh (diulang).
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

// K2 sebelum migrasi: ekspresi 944 (`m.JSONDATA.Benefit`) = view lama, baris demi baris.
func TestDBEkspresiSamaDenganView(t *testing.T) {
	u := pasang(t)
	u.periksa(t, "ekspresi = view", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_BENEFIT_LIFE m JOIN {s}.BENEFIT_LIFE v ON v.ID = m.ID WHERE m.JSONDATA.Benefit = v.BENEFIT`: 3,
		`SELECT COUNT(*) FROM {s}.M_BENEFIT_LIFE m WHERE m.JSONDATA.Benefit IS NULL`:                                            0,
		`SELECT COUNT(*) FROM {s}.M_BENEFIT_LIFE m WHERE JSON_VALUE(m.JSONDATA, '$.Benefit') = m.JSONDATA.Benefit`:              3,
	})
}

// 942-944 maju (diulang sebagian): nama lama dan view hilang, kolom = ID, BENEFIT; isi = JSON apa adanya (huruf tidak
// diubah migrasi); PK ikut RENAME; pelindung ORA-00904; mundur 944..942 memulihkan nama, view, JSONDATA.
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.lari(t, "942_benefit_life_ganti_nama", false, 1)
	u.lari(t, "943_benefit_life_kolom", false, 0)
	u.lari(t, "944_benefit_life_satu_tabel", false, 2)
	u.periksa(t, "sesudah 944", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME = 'M_BENEFIT_LIFE'`:                            0,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME = 'BENEFIT_LIFE'`:                                  0,
		`SELECT COUNT(*) FROM SYS.ALL_TABLES WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'BENEFIT_LIFE'`:                                1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'BENEFIT_LIFE'`:                           2,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'BENEFIT_LIFE' AND CONSTRAINT_TYPE = 'P'`: 1,
		`SELECT COUNT(*) FROM {s}.BENEFIT_LIFE WHERE (ID = '100001' AND BENEFIT = 'UJI RAWAT INAP')
			OR (ID = '100002' AND BENEFIT = 'UJI Meninggal Dunia') OR (ID = '100004' AND BENEFIT = 'UJI CACAT TETAP')`: 3,
	})
	for _, k := range []string{"943_benefit_life_kolom", "942_benefit_life_ganti_nama"} {
		u.gagalORA(t, ddl(t, k, true, u.skema)[0], "ORA-00904")
	}
	u.lari(t, "944_benefit_life_satu_tabel", true, 1)
	u.lari(t, "943_benefit_life_kolom", true, 0)
	l := ddl(t, "942_benefit_life_ganti_nama", true, u.skema)
	u.exec(t, l[0])
	u.exec(t, l[1])
	for i, q := range l {
		// Pelindung ORA-00942 (tabel sudah berganti nama balik saat diulang) = DITOLERANSI Bongkar.
		if _, err := u.repo.ExecContext(u.ctx, q); err != nil && !(i == 0 && strings.Contains(err.Error(), "ORA-00942")) {
			t.Fatalf("942 mundur %d: %v", i, err)
		}
	}
	u.periksa(t, "sesudah mundur", map[string]int{
		`SELECT COUNT(*) FROM {s}.BENEFIT_LIFE WHERE ID = '100002' AND BENEFIT = 'UJI Meninggal Dunia'`:                          1,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME = 'BENEFIT_LIFE'`:                           1,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'ENSURE_M_BENEFIT_LIFE_JSON'`: 1,
	})
}

// K2: BENEFIT NULL padahal JSON punya Benefit -> pemeriksaan 944 GAGAL KERAS (ORA-01407) SEBELUM JSONDATA dibuang;
// tidak ada yang berubah.
func TestDBPemeriksaanK2Berhenti(t *testing.T) {
	u := pasang(t)
	u.lari(t, "942_benefit_life_ganti_nama", false, 0)
	u.lari(t, "943_benefit_life_kolom", false, 0)
	l := ddl(t, "944_benefit_life_satu_tabel", false, u.skema)
	u.exec(t, l[0])
	u.exec(t, `UPDATE {s}.BENEFIT_LIFE SET BENEFIT = NULL WHERE ID = '100002'`)
	u.gagalORA(t, l[1], "ORA-01407")
	u.exec(t, `UPDATE {s}.BENEFIT_LIFE SET BENEFIT = 'UJI LAIN' WHERE ID = '100002'`)
	u.gagalORA(t, l[1], "ORA-01407")
	u.periksa(t, "JSONDATA utuh", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'BENEFIT_LIFE' AND COLUMN_NAME = 'JSONDATA'`: 1,
		`SELECT COUNT(*) FROM {s}.BENEFIT_LIFE WHERE ID IS NOT NULL`:                                                                       3,
	})
	u.exec(t, l[0]) // isi ulang (aman diulang) -> pemeriksaan lolos -> buang
	u.exec(t, l[1])
	u.exec(t, l[2])
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

// Seam: Add = '1' || LPAD(12, 5) = 100012; Edit; grid ID menurun; K3 ID dari sequence yang sudah ada = ErrIDTerpakai.
func TestDBLayanan(t *testing.T) {
	u := pasang(t)
	u.lari(t, "942_benefit_life_ganti_nama", false, 0)
	u.lari(t, "943_benefit_life_kolom", false, 0)
	u.lari(t, "944_benefit_life_satu_tabel", false, 0)
	l := layananOracle(u)
	penuh := services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	b, err := l.Simpan(u.ctx, penuh, "", models.Isian{Benefit: " uji kritis "})
	if err != nil || b != (models.Benefit{ID: "100012", Benefit: "UJI KRITIS"}) {
		t.Fatalf("add %+v %v", b, err)
	}
	if b, err := l.Simpan(u.ctx, penuh, "100002", models.Isian{Benefit: "uji meninggal"}); err != nil || b.Benefit != "UJI MENINGGAL" {
		t.Errorf("edit %+v %v", b, err)
	}
	h, err := l.Daftar(u.ctx, models.Saringan{})
	if err != nil || h.Total != 4 || h.Daftar[0].ID != "100012" || h.Daftar[3].ID != "100001" {
		t.Errorf("daftar %+v %v", h, err)
	}
	u.exec(t, `INSERT INTO {s}.BENEFIT_LIFE (ID, BENEFIT) VALUES ('100013', 'UJI DITULIS PEGA')`)
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{Benefit: "UJI"}); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("K3 mau ErrIDTerpakai, dapat %v", err)
	}
}
