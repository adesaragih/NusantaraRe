//go:build db

package repository_test

// Seam repository Cover Life terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`).
// Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 085 ditiru (fakta WO 08-10-2026): M_COVER_LIFE (ID VARCHAR2(10) NOT
// NULL PK + JSONDATA IS JSON, kunci `Cover` huruf campuran, kunci `Note` TIDAK ada) + view COVER_LIFE
// (`SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note`) + M_COVER_LIFE_SEQ (last_number 5). Migrasi modul 085-086
// dijalankan dari berkasnya. Fixture UJI-.
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
	"nusantarare/modul/coverlife/backend/models"
	"nusantarare/modul/coverlife/backend/repository"
	"nusantarare/modul/coverlife/backend/services"
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

// pasang - keadaan DEV sebelum 085 (bentuk DEV, nilai UJI-).
func pasang(t *testing.T) *ujiDB {
	t.Helper()
	u := buka(t)
	u.bongkar(t)
	baris := func(id, json string) string {
		return `INSERT INTO {s}.M_COVER_LIFE (ID, JSONDATA) VALUES ('` + id + `', '` + json + `')`
	}
	for _, q := range []string{
		`CREATE TABLE {s}.M_COVER_LIFE (ID VARCHAR2(10) NOT NULL PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_COVER_LIFE_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.COVER_LIFE AS SELECT ID, a.JSONDATA.Cover, a.JSONDATA.Note FROM {s}.M_COVER_LIFE a`,
		`CREATE SEQUENCE {s}.M_COVER_LIFE_SEQ START WITH 5`,
		baris("100001", `{"Cover":"UJI KECELAKAAN DIRI","pxObjClass":"ASM-FW-GISFW-Int-COVER_LIFE","pyRuleHarness":"UJI"}`),
		baris("100002", `{"Cover":"UJI JIWA","pxObjClass":"ASM-FW-GISFW-Int-COVER_LIFE"}`),
		baris("100003", `{"Cover":"UJI Kesehatan"}`),
		baris("100004", `{"Cover":"UJI RIDER","pyRuleHarness":"UJI"}`),
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
	for _, q := range []string{`DROP VIEW {s}.COVER_LIFE`, `DROP TABLE {s}.M_COVER_LIFE PURGE`, `DROP SEQUENCE {s}.M_COVER_LIFE_SEQ`} {
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

// Bukti sebelum migrasi: ekspresi 086 (`m.JSONDATA.Cover`, `m.JSONDATA.Note`) = view lama, baris demi baris; Note NULL
// di keduanya (DECODE menganggap NULL = NULL).
func TestDBEkspresiSamaDenganView(t *testing.T) {
	u := pasang(t)
	u.periksa(t, "ekspresi = view", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_COVER_LIFE m JOIN {s}.COVER_LIFE v ON v.ID = m.ID
			WHERE DECODE(m.JSONDATA.Cover, v.COVER, 1, 0) = 1 AND DECODE(m.JSONDATA.Note, v.NOTE, 1, 0) = 1`: 4,
		`SELECT COUNT(*) FROM {s}.COVER_LIFE WHERE NOTE IS NULL`: 4,
	})
}

// 085-086 maju (diulang sebagian): nama TETAP M_COVER_LIFE (TABLE), kolom = ID, COVER, NOTE (urutan view), nol
// JSONDATA, PK tetap, nol objek COVER_LIFE; isi = JSON apa adanya (huruf tidak diubah); NOTE NULL TIDAK menghentikan;
// pelindung ORA-00904; mundur 086..085 memulihkan JSONDATA, constraint, view PERSIS.
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.lari(t, "085_cover_life_kolom", false, 0)
	u.lari(t, "086_cover_life_satu_tabel", false, 2)
	u.periksa(t, "sesudah 086", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME = 'COVER_LIFE'`:                                   0,
		`SELECT COUNT(*) FROM SYS.ALL_TABLES WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE'`:                                   1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE'`:                              3,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA'`: 0,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE'
			AND ((COLUMN_NAME = 'ID' AND COLUMN_ID = 1) OR (COLUMN_NAME = 'COVER' AND COLUMN_ID = 2) OR (COLUMN_NAME = 'NOTE' AND COLUMN_ID = 3))`: 3,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE' AND CONSTRAINT_TYPE = 'P'`: 1,
		`SELECT COUNT(*) FROM {s}.M_COVER_LIFE WHERE NOTE IS NULL AND ((ID = '100001' AND COVER = 'UJI KECELAKAAN DIRI')
			OR (ID = '100002' AND COVER = 'UJI JIWA') OR (ID = '100003' AND COVER = 'UJI Kesehatan') OR (ID = '100004' AND COVER = 'UJI RIDER'))`: 4,
	})
	u.gagalORA(t, ddl(t, "085_cover_life_kolom", true, u.skema)[0], "ORA-00904")
	u.lari(t, "086_cover_life_satu_tabel", true, 0)
	u.lari(t, "085_cover_life_kolom", true, 0)
	u.periksa(t, "sesudah mundur", map[string]int{
		`SELECT COUNT(*) FROM {s}.COVER_LIFE WHERE ID = '100003' AND COVER = 'UJI Kesehatan' AND NOTE IS NULL`:                 1,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME = 'COVER_LIFE'`:                           1,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'ENSURE_M_COVER_LIFE_JSON'`: 1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE'`:                  2,
	})
}

// C1.3: COVER / NOTE BERBEDA dari view -> pemeriksaan 086 GAGAL KERAS (ORA-01407) SEBELUM JSONDATA dan view dibuang;
// NULL = NULL (NOTE 4/4) lolos.
func TestDBPemeriksaanC13(t *testing.T) {
	u := pasang(t)
	u.lari(t, "085_cover_life_kolom", false, 0)
	l := ddl(t, "086_cover_life_satu_tabel", false, u.skema)
	for _, rusak := range []string{
		`UPDATE {s}.M_COVER_LIFE SET COVER = NULL WHERE ID = '100002'`,
		`UPDATE {s}.M_COVER_LIFE SET NOTE = 'UJI ISI' WHERE ID = '100001'`,
		`UPDATE {s}.M_COVER_LIFE SET COVER = 'UJI Jiwa' WHERE ID = '100002'`,
	} {
		u.exec(t, l[0])
		u.exec(t, rusak)
		u.gagalORA(t, l[1], "ORA-01407")
	}
	u.periksa(t, "JSONDATA dan view utuh", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_COVER_LIFE' AND COLUMN_NAME = 'JSONDATA'`: 1,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME = 'COVER_LIFE'`:                                       1,
		`SELECT COUNT(*) FROM {s}.M_COVER_LIFE WHERE ID IS NOT NULL`:                                                                       4,
	})
	for _, q := range l { // isi ulang (aman diulang) -> pemeriksaan lolos walau NOTE NULL -> buang -> view
		u.exec(t, q)
	}
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

// Seam: Add = '1' || LPAD(5, 5) = 100005; Edit (Note); kembar ditolak; grid ID menaik 50; ID dari sequence yang sudah
// ada = ErrIDTerpakai.
func TestDBLayanan(t *testing.T) {
	u := pasang(t)
	u.lari(t, "085_cover_life_kolom", false, 0)
	u.lari(t, "086_cover_life_satu_tabel", false, 0)
	l := layananOracle(u)
	penuh := services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	c, err := l.Simpan(u.ctx, penuh, "", models.Isian{Cover: " UJI Penyakit Kritis ", Note: "uji catatan"})
	if err != nil || c != (models.Cover{ID: "100005", Cover: "UJI Penyakit Kritis", Note: "uji catatan"}) {
		t.Fatalf("add %+v %v", c, err)
	}
	if c, err := l.Simpan(u.ctx, penuh, "100002", models.Isian{Cover: "UJI JIWA", Note: "uji catatan jiwa"}); err != nil || c.Note != "uji catatan jiwa" {
		t.Errorf("edit %+v %v", c, err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{Cover: "uji kesehatan"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kembar %v", err)
	}
	h, err := l.Daftar(u.ctx, models.Saringan{})
	if err != nil || h.Total != 5 || h.Ukuran != 50 || h.Daftar[0].ID != "100001" || h.Daftar[4].ID != "100005" {
		t.Errorf("daftar %+v %v", h, err)
	}
	u.exec(t, `INSERT INTO {s}.M_COVER_LIFE (ID, COVER) VALUES ('100007', 'UJI DITULIS PEGA')`)
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{Cover: "UJI X"}); err != nil { // NEXTVAL 6 -> 100006
		t.Errorf("100006 %v", err)
	}
	if _, err := l.Simpan(u.ctx, penuh, "", models.Isian{Cover: "UJI Y"}); !errors.Is(err, services.ErrIDTerpakai) {
		t.Errorf("C2 mau ErrIDTerpakai, dapat %v", err)
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

// K0 di skema uji: pelari SUNGGUHAN atas migrasi modul ini + 900 / 901 / 909 inti menjalankan 085, 086 SEBELUM 900 dan
// 957 SESUDAHNYA; baris menu lahir MASTER TREATY URUTAN 16 DIMIGRASI '1'; diulang = nol langkah.
func TestDBUrutanPelari085Lalu900Lalu957(t *testing.T) {
	u := pasang(t)
	sumber := sumberModul(t)
	cacahT := func(nama string) int {
		n := 0
		if err := u.repo.QueryRowContext(u.ctx, strings.ReplaceAll(`SELECT COUNT(*) FROM {s}.T_MIGRASI WHERE NAMA = '`+nama+`'`,
			"{s}", u.skema)).Scan(&n); err != nil && !strings.Contains(err.Error(), "ORA-00942") {
			t.Fatalf("membaca T_MIGRASI: %v", err)
		}
		return n
	}
	if cacahT("085_cover_life_kolom") > 0 {
		t.Skip("skema uji sudah memuat migrasi modul ini (skemauji.Pasang); uji urutan butuh skema tanpa 085")
	}
	sudah := cacahT("900_m_nav_menu")
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
	i085, i957 := slices.Index(lap.Dijalankan, "085_cover_life_kolom"), slices.Index(lap.Dijalankan, "957_menu_coverlife")
	if i085 < 0 || i957 < i085 {
		t.Errorf("dijalankan %v", lap.Dijalankan)
	}
	if sudah == 0 {
		if i900 := slices.Index(lap.Dijalankan, "900_m_nav_menu"); !(i085 < i900 && i900 < i957) {
			t.Errorf("urutan 085 -> 900 -> 957 dilanggar: %v", lap.Dijalankan)
		}
	}
	for _, x := range lap.ObjekSudahAda {
		if strings.Contains(x, "COVER") {
			t.Errorf("dilewati, objeknya sudah ada: %s", x)
		}
	}
	u.periksa(t, "menu", map[string]int{
		`SELECT COUNT(*) FROM {s}.M_NAV_MENU WHERE KODE = 'coverlife' AND LABEL = 'Cover Life'
			AND GROUPMENU = 'MASTER TREATY' AND MODUL = 'coverlife' AND URUTAN = 16 AND DIMIGRASI = '1'`: 1,
		`SELECT COUNT(*) FROM {s}.T_MIGRASI WHERE NAMA IN ('085_cover_life_kolom', '086_cover_life_satu_tabel', '957_menu_coverlife')`: 3,
		`SELECT COUNT(*) FROM {s}.M_COVER_LIFE`: 4,
	})
	if ulang, err := migrasi.Jalankan(u.ctx, u.repo, intiSumber, sumber); err != nil || len(ulang.Dijalankan) != 0 {
		t.Errorf("diulang: %+v %v", ulang, err)
	}
}
