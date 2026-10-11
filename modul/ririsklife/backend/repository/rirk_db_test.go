//go:build db

package repository_test

// Seam repository R/I Risk terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`).
// Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 935 ditiru (fakta WO 08-10-2026): M_RIRISK_LIFE_SUMMARY (ID PK +
// JSONDATA IS JSON) + view RIRISK_LIFE_SUMMARY; M_RIRISK_LIFE (ID VARCHAR2(6) TANPA PK, JSONDATA IS JSON, kolom datar
// warisan BASI, RISK NUMBER tanpa skala, indeks warisan seperti INDEX4) + view RIRISK_LIFE; M_SITE_DATABASE; dua
// sequence warisan. Migrasi inti 935-940 dijalankan dari berkasnya. NOT NULL pada JSONDATA TIDAK ditiru (tidak terbukti).
// Fixture UJI-; OPERATORID UJI-.
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
	"nusantarare/modul/ririsklife/backend/models"
	"nusantarare/modul/ririsklife/backend/repository"
	"nusantarare/modul/ririsklife/backend/services"
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

// pasang - keadaan DEV sebelum 935.
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
	ring := func(id, json string) string {
		return `INSERT INTO {s}.M_RIRISK_LIFE_SUMMARY (ID, JSONDATA) VALUES ('` + id + `', '` + json + `')`
	}
	rinci := func(id, json, datar string) string {
		return `INSERT INTO {s}.M_RIRISK_LIFE (ID, JSONDATA, IDUSEDBY, USEDBY, YEAR, MONTH, RISK, CONTRACT) VALUES ('` + id + `', '` + json + `', ` + datar + `)`
	}
	for _, q := range []string{
		`CREATE TABLE {s}.M_RIRISK_LIFE_SUMMARY (ID VARCHAR2(10) PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_RIRISK_LIFE_SUMMARY_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.RIRISK_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM {s}.M_RIRISK_LIFE_SUMMARY a`,
		`CREATE TABLE {s}.M_RIRISK_LIFE (ID VARCHAR2(6), JSONDATA CLOB CONSTRAINT ENSURE_M_RIRISK_LIFE_JSON CHECK (JSONDATA IS JSON),
			IDUSEDBY VARCHAR2(100), USEDBY VARCHAR2(1000), YEAR VARCHAR2(10), MONTH VARCHAR2(10), RISK NUMBER, CONTRACT VARCHAR2(10))`,
		// Indeks warisan bernama lain berkolom pertama IDUSEDBY (kolom INDEX4 DEV tidak tercatat): 940 TIDAK boleh membuat
		// indeks kembar.
		`CREATE INDEX {s}.UJI_RIRISK_INDEX4 ON {s}.M_RIRISK_LIFE (IDUSEDBY)`,
		`CREATE VIEW {s}.RIRISK_LIFE AS SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.AGE, a.JSONDATA.YEAR,
			a.JSONDATA.MONTH, a.JSONDATA.RISK, a.JSONDATA.CONTRACT FROM {s}.M_RIRISK_LIFE a`,
		`CREATE TABLE {s}.M_SITE_DATABASE (ID NUMBER(10), CURRENT_SITE VARCHAR2(1))`,
		`CREATE SEQUENCE {s}.M_RIRISK_LIFE_SUMMARY_SEQ START WITH 5`,
		`CREATE SEQUENCE {s}.M_RIRISK_LIFE_SEQ START WITH 31723`,
		`INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (1, '1')`,
		`INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (2, '0')`,
		ring("1000003", `{"MODIFIEDDATE":"20191113T025753.044 GMT","OPERATORID":"UJI-LAMA","pxObjClass":"ASM-FW-GISFW-Int-RIRISK_LIFE_SUMMARY","USEDBY":"UJI RISK RETRO"}`),
		ring("1000004", `{"USEDBY":"UJI RISK B","OPERATORID":"UJI-APP"}`),
		// RISK berkoma desimal, kolom datar BASI (IDUSEDBY NULL, RISK 0) - 940 menimpanya dari JSON.
		rinci("131720", `{"CONTRACT":"1","IDUSEDBY":"1000003","RISK":"921,9","USEDBY":"UJI RISK RETRO","YEAR":"1","pxObjClass":"ASM-FW-GISFW-Int-RI_RISK_LIFE"}`,
			`NULL, 'UJI BASI', '9', NULL, 0, '9'`),
		// RISK bertitik 12 desimal, MONTH hanya di JSON.
		rinci("131721", `{"CONTRACT":"1","IDUSEDBY":"1000003","MONTH":"12","RISK":"580.894351210924","USEDBY":"UJI RISK RETRO"}`,
			`'1000003', 'UJI RISK RETRO', NULL, NULL, 580, '1'`),
		rinci("131722", `{"CONTRACT":"2","IDUSEDBY":"1000004","RISK":"5","USEDBY":"UJI RISK B","YEAR":"1"}`, `'1000004', 'UJI RISK B', '1', NULL, 5, '2'`),
	} {
		u.exec(t, q)
	}
	t.Cleanup(func() {
		u.bongkar(t)
		_ = repo.Close()
	})
	return u
}

// lari - langkah maju `kunci`: `sebagian` pernyataan pertama dulu (pelari gagal di tengah), lalu penuh (diulang).
func (u *ujiDB) lari(t *testing.T, eksekusi func(t *testing.T, q string), kunci string, mundur bool, sebagian int) {
	t.Helper()
	l := ddl(t, kunci, mundur, u.skema)
	for _, q := range l[:sebagian] {
		eksekusi(t, q)
	}
	for _, q := range l {
		eksekusi(t, q)
	}
}

// satuTabel - 935-940 (tanpa 941: M_NAV_MENU tidak ditiru), masing-masing diulang sebagian.
func (u *ujiDB) satuTabel(t *testing.T) {
	t.Helper()
	for _, p := range []struct {
		kunci    string
		sebagian int
	}{{"935_ririsk_life_summary_ganti_nama", 1}, {"936_ririsk_life_summary_kolom", 0}, {"937_ririsk_life_summary_satu_tabel", 2},
		{"938_ririsk_life_ganti_nama", 1}, {"939_ririsk_life_kolom", 0}, {"940_ririsk_life_satu_tabel", 4}} {
		u.lari(t, u.exec, p.kunci, false, p.sebagian)
	}
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for _, q := range []string{
		`DROP VIEW {s}.RIRISK_LIFE`, `DROP VIEW {s}.RIRISK_LIFE_SUMMARY`, `DROP TABLE {s}.RIRISK_LIFE PURGE`,
		`DROP TABLE {s}.RIRISK_LIFE_SUMMARY PURGE`, `DROP TABLE {s}.M_RIRISK_LIFE PURGE`, `DROP TABLE {s}.M_RIRISK_LIFE_SUMMARY PURGE`,
		`DROP TABLE {s}.M_SITE_DATABASE PURGE`, `DROP SEQUENCE {s}.M_RIRISK_LIFE_SUMMARY_SEQ`, `DROP SEQUENCE {s}.M_RIRISK_LIFE_SEQ`,
	} {
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

// gagalORA - pernyataan `q` HARUS gagal dengan kode Oracle `kode` (pelindung gagal-keras jalur mundur).
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

// 935-940 maju (diulang sebagian): nama lama dan view hilang, kolom = kolom view, kolom basi DITIMPA dari JSON, RISK
// "921,9" -> 921.9 dan "580.894351210924" utuh, PK ada, INDEX4 berkolom IDUSEDBY dipakai (nol indeks kembar);
// pelindung ORA-00904; mundur 940..935 (diulang sebagian) memulihkan nama, view, JSONDATA.
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	u.periksa(t, "sesudah 940", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME IN ('M_RIRISK_LIFE', 'M_RIRISK_LIFE_SUMMARY')`:       0,
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME LIKE 'RIRISK%'`:                                          0,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'RIRISK_LIFE_SUMMARY'`:                          4,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'RIRISK_LIFE'`:                                  8,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME LIKE 'RIRISK_LIFE%' AND COLUMN_NAME = 'JSONDATA'`: 0,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE_SUMMARY WHERE ID = '1000003' AND USEDBY = 'UJI RISK RETRO' AND OPERATORID = 'UJI-LAMA'
			AND MODIFIEDDATE = '20191113T025753.044 GMT'`: 1,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE ID = '131720' AND IDUSEDBY = '1000003' AND USEDBY = 'UJI RISK RETRO' AND YEAR = '1'
			AND MONTH IS NULL AND CONTRACT = '1' AND RISK = 921.9`: 1,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE ID = '131721' AND MONTH = '12' AND YEAR IS NULL AND RISK = 580.894351210924`: 1,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE AGE IS NOT NULL`:                                                             0,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'PK_RIRISK_LIFE'`:             1,
		`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER('{s}') AND INDEX_NAME = 'IX_RIRISK_LIFE_IDUSEDBY'`:             0,
		`SELECT COUNT(*) FROM SYS.ALL_IND_COLUMNS WHERE TABLE_OWNER = UPPER('{s}') AND TABLE_NAME = 'RIRISK_LIFE'
			AND COLUMN_NAME = 'IDUSEDBY' AND COLUMN_POSITION = 1`: 1,
		`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER('{s}') AND INDEX_NAME = 'IX_RIRISK_LIFE_SUMMARY_NAMA'`: 1,
	})
	// JSONDATA sudah dibuang: setiap mundur yang membuang kolom / memulihkan nama GAGAL KERAS.
	for _, k := range []string{"939_ririsk_life_kolom", "938_ririsk_life_ganti_nama", "936_ririsk_life_summary_kolom", "935_ririsk_life_summary_ganti_nama"} {
		u.gagalORA(t, ddl(t, k, true, u.skema)[0], "ORA-00904")
	}
	for _, p := range []struct {
		kunci    string
		sebagian int
	}{{"940_ririsk_life_satu_tabel", 3}, {"939_ririsk_life_kolom", 0}, {"938_ririsk_life_ganti_nama", 2},
		{"937_ririsk_life_summary_satu_tabel", 3}, {"936_ririsk_life_summary_kolom", 0}, {"935_ririsk_life_summary_ganti_nama", 2}} {
		l := ddl(t, p.kunci, true, u.skema)
		for _, q := range l[:p.sebagian] {
			u.exec(t, q)
		}
		for i, q := range l {
			// Pelindung ORA-00942 (tabel sudah berganti nama balik saat diulang) = DITOLERANSI Bongkar.
			if _, err := u.repo.ExecContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)); err != nil &&
				!(i == 0 && strings.Contains(err.Error(), "ORA-00942")) {
				t.Fatalf("%s mundur %d: %v", p.kunci, i, err)
			}
		}
	}
	u.periksa(t, "sesudah mundur", map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_VIEWS WHERE OWNER = UPPER('{s}') AND VIEW_NAME IN ('RIRISK_LIFE', 'RIRISK_LIFE_SUMMARY')`:                                         2,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE ID = '131720' AND IDUSEDBY = '1000003' AND RISK = '921.9'`:                                                          1,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE_SUMMARY WHERE ID = '1000003' AND USEDBY = 'UJI RISK RETRO' AND OPERATORID = 'UJI-LAMA'`:                                   1,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME IN ('ENSURE_M_RIRISK_LIFE_JSON', 'ENSURE_M_RIRISK_LIFE_SUMMARY_JSON')`: 2,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'PK_RIRISK_LIFE'`:                                                    0,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RIRISK_LIFE' AND COLUMN_NAME = 'AGE'`:                                  0,
	})
}

// K2: konversi RISK 940 TIDAK bergantung NLS sesi - di sesi berpemisah desimal KOMA (`,.`) hasilnya sama dengan sesi
// bertitik: "921,9" -> 921.9, "580.894351210924" -> 580.894351210924. Satu koneksi terkunci (db.Koneksi) supaya
// ALTER SESSION berlaku untuk pernyataan 940.
func TestDBKonversiRiskTanpaNLS(t *testing.T) {
	u := pasang(t)
	for _, k := range []string{"935_ririsk_life_summary_ganti_nama", "936_ririsk_life_summary_kolom",
		"937_ririsk_life_summary_satu_tabel", "938_ririsk_life_ganti_nama", "939_ririsk_life_kolom"} {
		u.lari(t, u.exec, k, false, 0)
	}
	kon, err := u.repo.Koneksi(u.ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = kon.Close() }()
	if _, err := kon.ExecContext(u.ctx, `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = ',.'`); err != nil {
		t.Fatal(err)
	}
	for _, q := range ddl(t, "940_ririsk_life_satu_tabel", false, u.skema) {
		if _, err := kon.ExecContext(u.ctx, q); err != nil {
			t.Fatalf("940 di NLS koma: %v", err)
		}
	}
	u.periksa(t, "NLS koma", map[string]int{
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE (ID = '131720' AND RISK = 921.9) OR (ID = '131721' AND RISK = 580.894351210924)
			OR (ID = '131722' AND RISK = 5)`: 3,
	})
}

func (u *ujiDB) dalamTx(t *testing.T, fn func(tx *db.Tx) error) error {
	t.Helper()
	tx, err := u.repo.Mulai(u.ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// layananOracle - aturan modul di atas Oracle uji; transaksi = db.Tx sungguhan (Commit / Rollback).
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

var penuhDB = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}

// Seam: rincian ID '1' || LPAD(31723, 5) = 131723; RISK pulang-pergi tanpa NLS (bertitik, 12 desimal); CONTRACT /
// YEAR / MONTH teks; kepemilikan IDUSEDBY; kembar (CONTRACT, YEAR, MONTH); Delete berantai satu transaksi.
func TestDBLayananRincian(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	l := layananOracle(u)
	k, err := l.SimpanRincian(u.ctx, penuhDB, "1000003", "", models.IsianRincian{Contract: "7", Month: "6", Risk: "921,9"})
	if err != nil || k.ID != "131723" {
		t.Fatalf("tambah %+v %v", k, err)
	}
	d, err := l.DaftarRincian(u.ctx, "1000003", 1)
	if err != nil || d.Total != 3 || d.Ukuran != 200 || d.Daftar[1] != (models.Rincian{ID: "131721", IDUsedBy: "1000003",
		UsedBy: "UJI RISK RETRO", Contract: "1", Month: "12", Risk: "580.894351210924"}) || d.Daftar[2].Risk != "921.9" {
		t.Errorf("daftar %+v %v", d, err)
	}
	if _, err := l.SimpanRincian(u.ctx, penuhDB, "1000003", "", models.IsianRincian{Contract: "07", Month: "06", Risk: "1"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kembar %v", err)
	}
	if _, err := l.SimpanRincian(u.ctx, penuhDB, "1000003", "131722", models.IsianRincian{Contract: "9", Year: "9", Risk: "9"}); !errors.Is(err, services.ErrRincianTidakAda) {
		t.Errorf("rincian milik 1000004 lewat 1000003: %v", err)
	}
	g := repository.Baru(u.repo)
	if err := g.UbahRincian(u.ctx, nil, models.Rincian{ID: "131722", IDUsedBy: "1000003", Contract: "9", Risk: "9"}); !errors.Is(err, repository.ErrTidakAda) {
		t.Errorf("UbahRincian ringkasan lain %v", err)
	}
	h, err := l.Hapus(u.ctx, penuhDB, "1000003")
	if err != nil || h.RincianTerhapus != 3 {
		t.Fatalf("hapus %+v %v", h, err)
	}
	u.periksa(t, "Delete berantai", map[string]int{
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE IDUSEDBY = '1000003'`:                                0,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE_SUMMARY WHERE ID = '1000003'`:                              0,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE ID = '131722' AND IDUSEDBY = '1000004' AND RISK = 5`: 1,
	})
}

// Save ringkasan (keputusan work owner 08-10-2026): Add = site || LPAD(5, 6) = 1000005; nama kembar ditolak lewat
// UPPER(TRIM(USEDBY)); Edit nama mengganti USEDBY rinciannya dalam satu transaksi.
func TestDBSimpanRingkasan(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	l := layananOracle(u)
	r, err := l.Simpan(u.ctx, penuhDB, "", models.Isian{UsedBy: " UJI RISK DB "})
	if err != nil || r.ID != "1000005" || r.OperatorID != "UJI-ADMIN" || r.UsedBy != "UJI RISK DB" {
		t.Fatalf("add %+v %v", r, err)
	}
	if _, err := l.Simpan(u.ctx, penuhDB, "", models.Isian{UsedBy: "uji risk db"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kembar %v", err)
	}
	if _, err := l.Simpan(u.ctx, penuhDB, "1000003", models.Isian{UsedBy: "UJI RISK RETRO 2"}); err != nil {
		t.Fatal(err)
	}
	u.periksa(t, "edit nama", map[string]int{
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE_SUMMARY WHERE ID = '1000003' AND USEDBY = 'UJI RISK RETRO 2' AND OPERATORID = 'UJI-ADMIN'`: 1,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE IDUSEDBY = '1000003' AND USEDBY = 'UJI RISK RETRO 2'`:                                2,
		`SELECT COUNT(*) FROM {s}.RIRISK_LIFE WHERE IDUSEDBY = '1000004' AND USEDBY = 'UJI RISK B'`:                                      1,
	})
}

// Simpan Upload: ringkasan baru = site || LPAD(5, 6) = 1000005, rincian '1' || LPAD; satu transaksi ringkasan +
// rincian - ID rincian kembar (ORA-00001, PK_RIRISK_LIFE) membatalkan KEDUANYA.
func TestDBUnggahDanRollback(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	l := layananOracle(u)
	h, err := l.SimpanUnggah(u.ctx, penuhDB, services.PermintaanUnggah{CSV: "USEDBY;CONTRACT;YEAR;MONTH;RISK\nUJI BARU;1;1;;0,5\nuji risk b;3;;24;1.25\n"})
	if err != nil || h.Disimpan != 2 || h.RingkasanBaru != 1 || h.Ringkasan[0].ID != "1000005" {
		t.Fatalf("unggah %+v %v", h, err)
	}
	g := repository.Baru(u.repo)
	ring0, rinc0 := u.cacah(t, `SELECT COUNT(*) FROM {s}.RIRISK_LIFE_SUMMARY`), u.cacah(t, `SELECT COUNT(*) FROM {s}.RIRISK_LIFE`)
	err = u.dalamTx(t, func(tx *db.Tx) error {
		if err := g.SisipRingkasan(u.ctx, tx, models.Ringkasan{ID: "1000009", UsedBy: "UJI GAGAL", OperatorID: "UJI-A"}); err != nil {
			return err
		}
		k := models.Rincian{ID: "199990", IDUsedBy: "1000009", UsedBy: "UJI GAGAL", Contract: "1", Risk: "1"}
		if err := g.SisipRincian(u.ctx, tx, k); err != nil {
			return err
		}
		return g.SisipRincian(u.ctx, tx, k)
	})
	if !errors.Is(err, repository.ErrKembar) {
		t.Fatalf("mau ErrKembar, dapat %v", err)
	}
	if r, k := u.cacah(t, `SELECT COUNT(*) FROM {s}.RIRISK_LIFE_SUMMARY`), u.cacah(t, `SELECT COUNT(*) FROM {s}.RIRISK_LIFE`); r != ring0 || k != rinc0 {
		t.Errorf("rollback: ringkasan %d->%d, rincian %d->%d", ring0, r, rinc0, k)
	}
}
