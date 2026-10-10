//go:build db

package repository_test

// Seam repository Cause Of Loss Life terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran -
// `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 090 ditiru (fakta WO 08-10-2026): M_CAUSEOFLOSS_LIFE
// (ID VARCHAR2(10) NOT NULL PK + JSONDATA IS JSON, kunci `CauseofLoss` huruf campuran; baris 100001
// `"CauseofLoss":""`) + view CAUSEOFLOSS_LIFE (`a.JSONDATA.CauseofLoss`) + M_CAUSEOFLOSS_LIFE_SEQ (last_number 5).
// Migrasi modul 090-092 dijalankan dari berkasnya. Fixture UJI-.
//
// ⚠️ Koneksi lewat `skemauji.BukaRepositori()`, BUKAN `skemauji.Buka()` (penjaga Claim Life
// `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci cacah pemanggilnya).
// ⚠️ DITULIS dan DIKOMPILASI (`go vet -tags=db`), BELUM PERNAH DIJALANKAN - executor tidak terhubung ke Oracle.

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/causeoflosslife/backend/models"
	"nusantarare/modul/causeoflosslife/backend/repository"
	"nusantarare/modul/causeoflosslife/backend/services"
	"nusantarare/uji/skemauji"
)

type ujiDB struct {
	repo  *db.DB
	skema string
	ctx   context.Context
}

// sumberModul - folder `migrations/` modul ini.
func sumberModul(t *testing.T) fstest.MapFS {
	t.Helper()
	entri, err := os.ReadDir(filepath.Join("..", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	sumber := fstest.MapFS{}
	for _, e := range entri {
		isi, err := os.ReadFile(filepath.Join("..", "migrations", e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+e.Name()] = &fstest.MapFile{Data: isi}
	}
	return sumber
}

// ddl - pernyataan migrasi modul `kunci` (maju atau mundur), `{skema}` diganti skema uji.
func ddl(t *testing.T, kunci string, mundur bool, skema string) []string {
	t.Helper()
	l, err := migrasi.Daftar(mundur, sumberModul(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range l {
		if migrasi.KunciLangkah(x.Nama) == kunci {
			var out []string
			for _, q := range x.Pernyataan {
				out = append(out, strings.ReplaceAll(q, "{skema}", skema))
			}
			return out
		}
	}
	t.Fatalf("langkah %s tidak ada", kunci)
	return nil
}

// buka - koneksi skema uji (lewati tanpa ORACLE_DSN).
func buka(t *testing.T) *ujiDB {
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
	return &ujiDB{repo: repo, skema: repo.Skema(), ctx: ctx}
}

// pasang - keadaan DEV sebelum 090 (bentuk DEV, nilai UJI-).
func pasang(t *testing.T) *ujiDB {
	t.Helper()
	u := buka(t)
	u.bongkar(t)
	baris := func(id, json string) string {
		return `INSERT INTO {s}.M_CAUSEOFLOSS_LIFE (ID, JSONDATA) VALUES ('` + id + `', '` + json + `')`
	}
	for _, q := range []string{
		`CREATE TABLE {s}.M_CAUSEOFLOSS_LIFE (ID VARCHAR2(10) NOT NULL PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_CAUSEOFLOSS_LIFE_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.CAUSEOFLOSS_LIFE AS SELECT a.ID, a.JSONDATA.CauseofLoss FROM {s}.M_CAUSEOFLOSS_LIFE a`,
		`CREATE SEQUENCE {s}.M_CAUSEOFLOSS_LIFE_SEQ START WITH 5`,
		baris("100001", `{"CauseofLoss":"","pxObjClass":"ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE","pyRuleHarness":"UJI"}`),
		baris("100002", `{"CauseofLoss":"UJI SAKIT","pxObjClass":"ASM-FW-GISFW-Int-CAUSEOFLOSS_LIFE"}`),
		baris("100003", `{"CauseofLoss":"UJI Kecelakaan"}`),
		baris("100004", `{"CauseofLoss":"UJI SEMUA SEBAB","pyRuleHarness":"UJI"}`),
	} {
		u.exec(t, q)
	}
	t.Cleanup(func() {
		u.bongkar(t)
		_ = u.repo.Close()
	})
	return u
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for _, q := range []string{`DROP VIEW {s}.CAUSEOFLOSS_LIFE`, `DROP TABLE {s}.CAUSEOFLOSS_LIFE PURGE`, `DROP TABLE {s}.M_CAUSEOFLOSS_LIFE PURGE`,
		`DROP SEQUENCE {s}.M_CAUSEOFLOSS_LIFE_SEQ`} {
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

// K2 sebelum migrasi: ekspresi 092 (`m.JSONDATA.CauseofLoss`) = view lama, baris demi baris - termasuk 100001
// (`""` -> NULL di keduanya; DECODE menganggap NULL = NULL).
func TestDBEkspresiSamaDenganView(t *testing.T) {
	u := pasang(t)
	u.periksa(t, "ekspresi = view", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_CAUSEOFLOSS_LIFE m JOIN {s}.CAUSEOFLOSS_LIFE v ON v.ID = m.ID
			WHERE DECODE(m.JSONDATA.CauseofLoss, v.CAUSEOFLOSS, 1, 0) = 1`: 4,
		`SELECT COUNT(*) FROM {s}.M_CAUSEOFLOSS_LIFE m WHERE m.JSONDATA.CauseofLoss IS NULL`:    1,
		`SELECT COUNT(*) FROM {s}.CAUSEOFLOSS_LIFE WHERE CAUSEOFLOSS IS NULL AND ID = '100001'`: 1,
	})
}

// 090-092 maju (diulang sebagian): nama lama dan view hilang, kolom = ID, CAUSEOFLOSS; isi = JSON apa adanya (huruf
// tidak diubah); 100001 NULL dan TIDAK menghentikan pemeriksaan; PK ikut RENAME; pelindung ORA-00904; mundur 092..090
// memulihkan nama, view, JSONDATA.
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.lari(t, "090_causeofloss_life_ganti_nama", false, 1)
	u.lari(t, "091_causeofloss_life_kolom", false, 0)
	u.lari(t, "092_causeofloss_life_satu_tabel", false, 2)
	u.periksa(t, "sesudah 092", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME = 'M_CAUSEOFLOSS_LIFE'`:                            0,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME = 'CAUSEOFLOSS_LIFE'`:                                  0,
		`SELECT COUNT(*) FROM SYS.ALL_TABLES WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE'`:                                1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE'`:                           2,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND CONSTRAINT_TYPE = 'P'`: 1,
		`SELECT COUNT(*) FROM {s}.CAUSEOFLOSS_LIFE WHERE (ID = '100001' AND CAUSEOFLOSS IS NULL)
			OR (ID = '100002' AND CAUSEOFLOSS = 'UJI SAKIT') OR (ID = '100003' AND CAUSEOFLOSS = 'UJI Kecelakaan')
			OR (ID = '100004' AND CAUSEOFLOSS = 'UJI SEMUA SEBAB')`: 4,
	})
	for _, k := range []string{"091_causeofloss_life_kolom", "090_causeofloss_life_ganti_nama"} {
		u.gagalORA(t, ddl(t, k, true, u.skema)[0], "ORA-00904")
	}
	u.lari(t, "092_causeofloss_life_satu_tabel", true, 1)
	u.lari(t, "091_causeofloss_life_kolom", true, 0)
	l := ddl(t, "090_causeofloss_life_ganti_nama", true, u.skema)
	u.exec(t, l[0])
	u.exec(t, l[1])
	for i, q := range l {
		// Pelindung ORA-00942 (tabel sudah berganti nama balik saat diulang) = DITOLERANSI Bongkar.
		if _, err := u.repo.ExecContext(u.ctx, q); err != nil && !(i == 0 && strings.Contains(err.Error(), "ORA-00942")) {
			t.Fatalf("090 mundur %d: %v", i, err)
		}
	}
	u.periksa(t, "sesudah mundur", map[string]int{
		`SELECT COUNT(*) FROM {s}.CAUSEOFLOSS_LIFE WHERE ID = '100003' AND CAUSEOFLOSS = 'UJI Kecelakaan'`:                           1,
		`SELECT COUNT(*) FROM {s}.CAUSEOFLOSS_LIFE WHERE ID = '100001' AND CAUSEOFLOSS IS NULL`:                                      1,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME = 'CAUSEOFLOSS_LIFE'`:                           1,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'ENSURE_M_CAUSEOFLOSS_LIFE_JSON'`: 1,
	})
}

// K1.4: kolom BERBEDA dari view -> pemeriksaan 092 GAGAL KERAS (ORA-01407) SEBELUM JSONDATA dibuang (tiga cabang:
// kolom NULL padahal view terisi, kolom terisi padahal view NULL, isi berbeda); NULL = NULL (100001) lolos.
func TestDBPemeriksaanK14(t *testing.T) {
	u := pasang(t)
	u.lari(t, "090_causeofloss_life_ganti_nama", false, 0)
	u.lari(t, "091_causeofloss_life_kolom", false, 0)
	l := ddl(t, "092_causeofloss_life_satu_tabel", false, u.skema)
	for _, rusak := range []string{
		`UPDATE {s}.CAUSEOFLOSS_LIFE SET CAUSEOFLOSS = NULL WHERE ID = '100002'`,
		`UPDATE {s}.CAUSEOFLOSS_LIFE SET CAUSEOFLOSS = 'UJI ISI' WHERE ID = '100001'`,
		`UPDATE {s}.CAUSEOFLOSS_LIFE SET CAUSEOFLOSS = 'UJI Sakit' WHERE ID = '100002'`,
	} {
		u.exec(t, l[0])
		u.exec(t, rusak)
		u.gagalORA(t, l[1], "ORA-01407")
	}
	u.periksa(t, "JSONDATA utuh", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'CAUSEOFLOSS_LIFE' AND COLUMN_NAME = 'JSONDATA'`: 1,
		`SELECT COUNT(*) FROM {s}.CAUSEOFLOSS_LIFE WHERE ID IS NOT NULL`:                                                                       4,
	})
	u.exec(t, l[0]) // isi ulang (aman diulang) -> pemeriksaan lolos walau 100001 NULL -> buang
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

// Seam: Add = '1' || LPAD(5, 5) = 100005; Edit; kembar ditolak; grid ID menaik; K3 ID dari sequence yang sudah ada =
// ErrIDTerpakai.
func TestDBLayanan(t *testing.T) {
	u := pasang(t)
	for _, k := range []string{"090_causeofloss_life_ganti_nama", "091_causeofloss_life_kolom", "092_causeofloss_life_satu_tabel"} {
		u.lari(t, k, false, 0)
	}
	l := layananOracle(u)
	penuh := services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	c, err := l.Simpan(u.ctx, penuh, "", models.Isian{CauseOfLoss: " UJI Bencana "})
	if err != nil || c != (models.CauseOfLoss{ID: "100005", CauseOfLoss: "UJI Bencana"}) {
		t.Fatalf("add %+v %v", c, err)
	}
	if c, err := l.Simpan(u.ctx, penuh, "100002", models.Isian{CauseOfLoss: "UJI SAKIT BERAT"}); err != nil || c.CauseOfLoss != "UJI SAKIT BERAT" {
		t.Errorf("edit %+v %v", c, err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{CauseOfLoss: "uji kecelakaan"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kembar %v", err)
	}
	h, err := l.Daftar(u.ctx, models.Saringan{})
	if err != nil || h.Total != 5 || h.Daftar[0] != (models.CauseOfLoss{ID: "100001"}) || h.Daftar[4].ID != "100005" {
		t.Errorf("daftar %+v %v", h, err)
	}
	u.exec(t, `INSERT INTO {s}.CAUSEOFLOSS_LIFE (ID, CAUSEOFLOSS) VALUES ('100007', 'UJI DITULIS PEGA')`)
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{CauseOfLoss: "UJI X"}); err != nil { // NEXTVAL 6 -> 100006
		t.Errorf("100006 %v", err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{CauseOfLoss: "UJI Y"}); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("K3 mau ErrIDTerpakai, dapat %v", err)
	}
}

// salinInti - langkah inti `nama` (maju + `_down`) dari migrasi inti SUNGGUHAN.
func salinInti(t *testing.T, ke fstest.MapFS, nama ...string) {
	t.Helper()
	for _, n := range nama {
		for _, b := range []string{n + ".sql", n + "_down.sql"} {
			isi, err := fs.ReadFile(inti.SumberMigrasi(), "migrations/"+b)
			if err != nil {
				t.Fatal(err)
			}
			ke["migrations/"+b] = &fstest.MapFile{Data: isi}
		}
	}
}

// K0 di skema uji: pelari SUNGGUHAN (`migrasi.Jalankan`) atas migrasi modul ini + 900 / 901 / 909 inti menjalankan 090,
// 091, 092 SEBELUM 900 dan 955 SESUDAHNYA; baris menu lahir MASTER TREATY URUTAN 14 DIMIGRASI '1'; diulang = nol langkah.
// Langkah inti yang SUDAH tercatat di skema uji (skema yang sudah dipasang) tidak dijalankan ulang dan tidak dibongkar.
func TestDBUrutanPelari090Lalu900Lalu955(t *testing.T) {
	u := pasang(t)
	sumber := sumberModul(t)
	// T_MIGRASI belum ada (ORA-00942) = skema kosong: 900 dijalankan DAN dibongkar uji ini.
	sudah := 0
	if err := u.repo.QueryRowContext(u.ctx, strings.ReplaceAll(`SELECT COUNT(*) FROM {s}.T_MIGRASI WHERE NAMA = '900_m_nav_menu'`,
		"{s}", u.skema)).Scan(&sudah); err != nil && !strings.Contains(err.Error(), "ORA-00942") {
		t.Fatalf("membaca T_MIGRASI: %v", err)
	}
	intiSumber := fstest.MapFS{}
	salinInti(t, intiSumber, "900_m_nav_menu", "901_m_nav_menu_datar", "909_m_nav_menu_master_treaty")
	lap, err := migrasi.Jalankan(u.ctx, u.repo, intiSumber, sumber)
	if err != nil {
		t.Fatalf("pelari: %v (laporan %+v)", err, lap)
	}
	t.Cleanup(func() {
		if _, err := migrasi.Bongkar(u.ctx, u.repo, sumber); err != nil {
			t.Errorf("bongkar modul: %v", err)
		}
		if sudah == 0 {
			if _, err := migrasi.Bongkar(u.ctx, u.repo, intiSumber); err != nil {
				t.Errorf("bongkar inti: %v", err)
			}
		}
	})
	i090, i955 := slices.Index(lap.Dijalankan, "090_causeofloss_life_ganti_nama"), slices.Index(lap.Dijalankan, "955_menu_causeoflosslife")
	if i090 < 0 || i955 < i090 {
		t.Errorf("dijalankan %v", lap.Dijalankan)
	}
	if sudah == 0 {
		if i900 := slices.Index(lap.Dijalankan, "900_m_nav_menu"); !(i090 < i900 && i900 < i955) {
			t.Errorf("urutan 090 -> 900 -> 955 dilanggar: %v", lap.Dijalankan)
		}
	}
	for _, x := range lap.ObjekSudahAda {
		if strings.Contains(x, "CAUSEOFLOSS") {
			t.Errorf("dilewati, objeknya sudah ada: %s", x)
		}
	}
	u.periksa(t, "menu", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_NAV_MENU WHERE KODE = 'causeoflosslife' AND LABEL = 'Cause Of Loss Life'
			AND GROUPMENU = 'MASTER TREATY' AND MODUL = 'causeoflosslife' AND URUTAN = 14 AND DIMIGRASI = '1'`: 1,
		`SELECT COUNT(*) FROM {s}.T_MIGRASI WHERE NAMA IN ('090_causeofloss_life_ganti_nama', '091_causeofloss_life_kolom',
			'092_causeofloss_life_satu_tabel', '955_menu_causeoflosslife')`: 4,
		`SELECT COUNT(*) FROM {s}.CAUSEOFLOSS_LIFE`: 4,
	})
	if ulang, err := migrasi.Jalankan(u.ctx, u.repo, intiSumber, sumber); err != nil || len(ulang.Dijalankan) != 0 {
		t.Errorf("diulang: %+v %v", ulang, err)
	}
}
