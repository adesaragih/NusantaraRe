//go:build db

package repository_test

// Seam repository R/I Comm Life terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran -
// `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 931 ditiru (fakta WO 08-10-2026): M_RICOMM_LIFE_SUMMARY
// (ID PK + JSONDATA IS JSON) + view RICOMM_LIFE_SUMMARY, M_RICOMM_LIFE (ID PK + JSONDATA IS JSON, 0 baris), tabel flat
// RICOMM_LIFE dari berkas migrasi inti 924 YANG SAMA dengan produksi, M_SITE_DATABASE, dan dua sequence warisan; lalu
// migrasi inti 931-934 dijalankan dari berkasnya (RALAT R1). NOT NULL pada JSONDATA TIDAK ditiru (tidak terbukti di repo).
// Angka tidak bergantung NLS sesi pool (fmtAngka berargumen NLS, AngkaOracle di Go). Fixture UJI-; OPERATORID UJI-.
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
	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/repository"
	"nusantarare/modul/ricommlife/backend/services"
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

// pasang - keadaan DEV sebelum 931 (tanpa 931-934).
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
	for _, q := range []string{
		`CREATE TABLE {s}.M_RICOMM_LIFE_SUMMARY (ID VARCHAR2(10) PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_RICOMM_LIFE_SUMMARY_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.RICOMM_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM {s}.M_RICOMM_LIFE_SUMMARY a`,
		`CREATE TABLE {s}.M_RICOMM_LIFE (ID VARCHAR2(10) PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_RICOMM_LIFE_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE TABLE {s}.M_SITE_DATABASE (ID NUMBER(10), CURRENT_SITE VARCHAR2(1))`,
		`CREATE SEQUENCE {s}.M_RICOMM_LIFE_SUMMARY_SEQ START WITH 5`,
		`CREATE SEQUENCE {s}.M_RICOMM_LIFE_SEQ START WITH 44`,
		`INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (1, '1')`,
		`INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (2, '0')`,
		`INSERT INTO {s}.M_RICOMM_LIFE_SUMMARY (ID, JSONDATA) VALUES ('1000003', '{"MODIFIEDDATE":"20181205T073755.559 GMT","OPERATORID":"UJI-LAMA","pxObjClass":"ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY","USEDBY":"UJI COMM RETRO"}')`,
		`INSERT INTO {s}.M_RICOMM_LIFE_SUMMARY (ID, JSONDATA) VALUES ('1000004', '{"USEDBY":"UJI COMM B","OPERATORID":"UJI-APP"}')`,
	} {
		u.exec(t, q)
	}
	for _, q := range ddl(t, "924_ricomm_life", false, u.skema) {
		u.exec(t, q)
	}
	u.exec(t, `INSERT INTO {s}.RICOMM_LIFE (ID, IDUSEDBY, USEDBY, CONTRACT, YEAR, COMM) VALUES ('1000043', '1000004', 'UJI COMM B', 1, 2022, 0.5)`)
	t.Cleanup(func() {
		u.bongkar(t)
		_ = repo.Close()
	})
	return u
}

// satuTabel - 931 lalu 932 (sebagian dua kali, lalu penuh), 933 lalu 934 (sebagian dua kali, lalu penuh): pelari
// gagal di tengah lalu diulang.
func (u *ujiDB) satuTabel(t *testing.T) {
	t.Helper()
	for _, p := range []struct {
		kolom, satu string
		sebagian    int
	}{{"931_m_ricomm_life_summary_kolom", "932_m_ricomm_life_summary_satu_tabel", 2},
		{"933_m_ricomm_life_kolom", "934_m_ricomm_life_satu_tabel", 4}} {
		for _, q := range ddl(t, p.kolom, false, u.skema) {
			u.exec(t, q)
		}
		l := ddl(t, p.satu, false, u.skema)
		for _, q := range l[:p.sebagian] {
			u.exec(t, q)
		}
		for _, q := range l {
			u.exec(t, q)
		}
	}
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for _, q := range []string{
		`DROP TABLE {s}.RICOMM_LIFE PURGE`, `DROP VIEW {s}.RICOMM_LIFE`, `DROP VIEW {s}.RICOMM_LIFE_SUMMARY`,
		`DROP TABLE {s}.M_RICOMM_LIFE_SUMMARY PURGE`, `DROP TABLE {s}.M_RICOMM_LIFE PURGE`, `DROP TABLE {s}.M_SITE_DATABASE PURGE`,
		`DROP SEQUENCE {s}.M_RICOMM_LIFE_SUMMARY_SEQ`, `DROP SEQUENCE {s}.M_RICOMM_LIFE_SEQ`,
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
		t.Fatalf("%.50s: %v", q, err)
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
		t.Fatalf("%.50s: %v", q, err)
	}
	return n
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

// 931-934 maju (diulang sebagian), pelindung ORA-00904 931_down/933_down, lalu mundur 934..931 (diulang sebagian).
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	for q, mau := range map[string]int{
		`SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE_SUMMARY WHERE (ID = '1000003' AND USEDBY = 'UJI COMM RETRO' AND OPERATORID = 'UJI-LAMA'
			AND MODIFIEDDATE = '20181205T073755.559 GMT') OR (ID = '1000004' AND USEDBY = 'UJI COMM B' AND MODIFIEDDATE IS NULL)`: 2,
		`SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE WHERE ID = '1000043' AND IDUSEDBY = '1000004' AND CONTRACT = 1 AND YEAR = 2022 AND COMM = 0.5`:             1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RICOMM_LIFE_SUMMARY'`:                                     4,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RICOMM_LIFE'`:                                             6,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME LIKE 'M_RICOMM_LIFE%' AND COLUMN_NAME = 'JSONDATA'`:            0,
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME IN ('RICOMM_LIFE', 'RICOMM_LIFE_SUMMARY')`:                        0,
		`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER('{s}') AND INDEX_NAME IN ('IX_M_RICOMM_LIFE_IDUSEDBY', 'IX_M_RICOMM_LIFE_SUMMARY_NAMA')`: 2,
	} {
		if n := u.cacah(t, q); n != mau {
			t.Errorf("%d, mau %d: %s", n, mau, q)
		}
	}
	// JSONDATA sudah dibuang: 931_down / 933_down langsung TIDAK BOLEH membuang kolom.
	u.gagalORA(t, ddl(t, "933_m_ricomm_life_kolom", true, u.skema)[0], "ORA-00904")
	u.gagalORA(t, ddl(t, "931_m_ricomm_life_summary_kolom", true, u.skema)[0], "ORA-00904")
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE WHERE COMM IS NOT NULL`) +
		u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE_SUMMARY WHERE USEDBY IS NOT NULL`); n != 3 {
		t.Errorf("pelindung: kolom harus utuh (%d)", n)
	}
	// Mundur aman diulang: sebagian dulu, lalu penuh.
	for _, p := range []struct {
		satu, kolom string
		sebagian    int
	}{{"934_m_ricomm_life_satu_tabel", "933_m_ricomm_life_kolom", 4},
		{"932_m_ricomm_life_summary_satu_tabel", "931_m_ricomm_life_summary_kolom", 3}} {
		l := ddl(t, p.satu, true, u.skema)
		for _, q := range l[:p.sebagian] {
			u.exec(t, q)
		}
		for _, q := range l {
			u.exec(t, q)
		}
		for _, q := range ddl(t, p.kolom, true, u.skema) {
			u.exec(t, q)
		}
	}
	for q, mau := range map[string]int{
		`SELECT COUNT(*) FROM {s}.RICOMM_LIFE WHERE ID = '1000043' AND IDUSEDBY = '1000004' AND COMM = 0.5`:                                                                               1,
		`SELECT COUNT(*) FROM {s}.RICOMM_LIFE_SUMMARY WHERE ID = '1000003' AND USEDBY = 'UJI COMM RETRO' AND OPERATORID = 'UJI-LAMA'`:                                                     1,
		`SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE m WHERE m.ID = '1000043' AND JSON_VALUE(m.JSONDATA, '$.IDUSEDBY') = '1000004'`:                                                            1,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME IN ('ENSURE_M_RICOMM_LIFE_JSON', 'ENSURE_M_RICOMM_LIFE_SUMMARY_JSON', 'PK_RICOMM_LIFE')`: 3,
		`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER('{s}') AND INDEX_NAME = 'IX_RICOMM_LIFE_IDUSEDBY'`:                                                                      1,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME LIKE 'M_RICOMM_LIFE%'`:                                                                        4,
	} {
		if n := u.cacah(t, q); n != mau {
			t.Errorf("sesudah mundur %d, mau %d: %s", n, mau, q)
		}
	}
}

// Seam: situs + sequence warisan membentuk ID 7 karakter; ringkasan sisip/ubah kolom; rincian tulis/baca angka tanpa
// NLS; ubah rincian milik ringkasan lain ditolak; Delete beserta rincian.
func TestSeamRepositoryOracle(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	g := repository.Baru(u.repo)
	situs, err := g.Situs(u.ctx, nil)
	if err != nil || situs != "1" {
		t.Fatalf("situs %q %v", situs, err)
	}
	nomor, err := g.NomorBaru(u.ctx, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if id, err := models.BentukID(situs, nomor); err != nil || id != "1000044" {
		t.Errorf("ID rincian %q %v", id, err)
	}
	err = u.dalamTx(t, func(tx *db.Tx) error {
		if err := g.SisipRingkasan(u.ctx, tx, models.Ringkasan{ID: "1000005", UsedBy: "UJI COMM C", OperatorID: "UJI-A",
			ModifiedDate: "20261006T030405.600 GMT"}); err != nil {
			return err
		}
		if err := g.UbahRingkasan(u.ctx, tx, models.Ringkasan{ID: "1000003", UsedBy: "UJI COMM RETRO 2", OperatorID: "UJI-B",
			ModifiedDate: "20261006T030405.600 GMT"}); err != nil {
			return err
		}
		for _, k := range []models.Komisi{
			{ID: "1000044", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO 2", Contract: "1", Year: "2026", Comm: "0.5"},
			{ID: "1000045", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO 2", Contract: "99999", Year: "2021", Comm: "123456789012345678901234567890.12345678"},
		} {
			if err := g.SisipKomisi(u.ctx, tx, k); err != nil {
				return err
			}
		}
		return g.UbahKomisi(u.ctx, tx, models.Komisi{ID: "1000044", IDUsedBy: "1000003", Contract: "2", Year: "2026", Comm: "12.25"})
	})
	if err != nil {
		t.Fatal(err)
	}
	r, err := g.Ambil(u.ctx, nil, "1000003")
	if err != nil || r.UsedBy != "UJI COMM RETRO 2" || r.OperatorID != "UJI-B" {
		t.Errorf("ringkasan %+v %v", r, err)
	}
	if err := g.UbahRingkasan(u.ctx, nil, models.Ringkasan{ID: "UJI-TIADA", UsedBy: "X"}); !errors.Is(err, repository.ErrTidakAda) {
		t.Errorf("ubah ringkasan tidak ada %v", err)
	}
	d, total, err := g.DaftarKomisi(u.ctx, "1000003", 1)
	if err != nil || total != 2 || len(d) != 2 || d[0].Comm != "12.25" || d[0].Contract != "2" || d[1].Comm != "123456789012345678901234567890.12345678" {
		t.Errorf("rincian %+v %d %v", d, total, err)
	}
	if err := g.UbahKomisi(u.ctx, nil, models.Komisi{ID: "1000044", IDUsedBy: "1000005", Contract: "1", Year: "2021", Comm: "1"}); !errors.Is(err, repository.ErrTidakAda) {
		t.Errorf("ubah milik ringkasan lain %v", err)
	}
	var n int
	err = u.dalamTx(t, func(tx *db.Tx) error {
		if n, err = g.HapusKomisi(u.ctx, tx, "1000003"); err != nil {
			return err
		}
		return g.HapusRingkasan(u.ctx, tx, "1000003")
	})
	if err != nil || n != 2 {
		t.Errorf("hapus %d %v", n, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE WHERE IDUSEDBY = '1000004'`); n != 1 {
		t.Errorf("Delete menyentuh rincian ringkasan lain (%d)", n)
	}
	// Situs tidak tepat satu baris.
	u.exec(t, `INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (3, '1')`)
	if _, err := g.Situs(u.ctx, nil); !errors.Is(err, repository.ErrSitus) {
		t.Errorf("dua situs aktif %v", err)
	}
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

// Rumus ID lewat layanan: ringkasan = site 1 || LPAD(M_RICOMM_LIFE_SUMMARY_SEQ 5) = 1000005; rincian = site 1 ||
// LPAD(M_RICOMM_LIFE_SEQ 44) = 1000044 (1000043 terpakai dari 924, dilewati bila perlu). Pulang-pergi kolom.
func TestDBRumusIDDanPulangPergi(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	l := layananOracle(u)
	r, err := l.Simpan(u.ctx, penuhDB, "", models.Isian{UsedBy: "UJI COMM DB"})
	if err != nil || r.ID != "1000005" || r.OperatorID != "UJI-ADMIN" {
		t.Fatalf("ringkasan %+v %v", r, err)
	}
	k, err := l.SimpanKomisi(u.ctx, penuhDB, r.ID, "", models.IsianKomisi{Contract: "99999", Year: "2026", Comm: "0,00000001"})
	if err != nil || k.ID != "1000044" {
		t.Fatalf("rincian %+v %v", k, err)
	}
	d, err := l.DaftarKomisi(u.ctx, r.ID, 1)
	if err != nil || d.Total != 1 || d.Daftar[0] != (models.Komisi{ID: "1000044", IDUsedBy: "1000005", UsedBy: "UJI COMM DB",
		Contract: "99999", Year: "2026", Comm: "0.00000001"}) {
		t.Errorf("pulang-pergi %+v %v", d, err)
	}
	// Edit nama ringkasan menyalin nama ke rinciannya (satu transaksi).
	if _, err := l.Simpan(u.ctx, penuhDB, r.ID, models.Isian{UsedBy: "UJI COMM DB 2"}); err != nil {
		t.Fatal(err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE WHERE IDUSEDBY = '1000005' AND USEDBY = 'UJI COMM DB 2'`); n != 1 {
		t.Errorf("salinan nama rincian %d", n)
	}
}

// Satu transaksi ringkasan + rincian: kegagalan di tengah (ID rincian kembar, ORA-00001) membatalkan KEDUANYA.
func TestDBTransaksiGagalRollbackKeduanya(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	g := repository.Baru(u.repo)
	ring0, kom0 := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE_SUMMARY`), u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE`)
	err := u.dalamTx(t, func(tx *db.Tx) error {
		if err := g.SisipRingkasan(u.ctx, tx, models.Ringkasan{ID: "1000009", UsedBy: "UJI GAGAL", OperatorID: "UJI-A",
			ModifiedDate: "20261006T030405.600 GMT"}); err != nil {
			return err
		}
		k := models.Komisi{ID: "1000090", IDUsedBy: "1000009", UsedBy: "UJI GAGAL", Contract: "1", Year: "2026", Comm: "1"}
		if err := g.SisipKomisi(u.ctx, tx, k); err != nil {
			return err
		}
		return g.SisipKomisi(u.ctx, tx, k)
	})
	if !errors.Is(err, repository.ErrKembar) {
		t.Fatalf("mau ErrKembar, dapat %v", err)
	}
	if r, k := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE_SUMMARY`), u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE`); r != ring0 || k != kom0 {
		t.Errorf("rollback: ringkasan %d->%d, rincian %d->%d", ring0, r, kom0, k)
	}
}

// Kembar (CONTRACT, YEAR) di satu ringkasan ditolak; rincian milik ringkasan LAIN tidak bisa diubah lewat ringkasan
// ini (ErrKomisiTidakAda); Delete ringkasan menghapus rinciannya saja (berantai, satu transaksi).
func TestDBKembarKepemilikanDanDeleteBerantai(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	l := layananOracle(u)
	if _, err := l.SimpanKomisi(u.ctx, penuhDB, "1000003", "", models.IsianKomisi{Contract: "1", Year: "2026", Comm: "1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := l.SimpanKomisi(u.ctx, penuhDB, "1000003", "", models.IsianKomisi{Contract: "01", Year: "2026", Comm: "2"}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("kembar %v", err)
	}
	if _, err := l.SimpanKomisi(u.ctx, penuhDB, "1000003", "1000043", models.IsianKomisi{Contract: "9", Year: "2030", Comm: "9"}); !errors.Is(err, services.ErrKomisiTidakAda) {
		t.Errorf("ubah rincian milik 1000004 lewat 1000003: %v", err)
	}
	if _, err := l.SimpanKomisi(u.ctx, penuhDB, "1000003", "", models.IsianKomisi{Contract: "2", Year: "2026", Comm: "2"}); err != nil {
		t.Fatal(err)
	}
	h, err := l.Hapus(u.ctx, penuhDB, "1000003")
	if err != nil || h.KomisiTerhapus != 2 {
		t.Fatalf("hapus %+v %v", h, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE WHERE IDUSEDBY = '1000003'`) +
		u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE_SUMMARY WHERE ID = '1000003'`); n != 0 {
		t.Errorf("sisa sesudah Delete berantai %d", n)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RICOMM_LIFE WHERE ID = '1000043' AND CONTRACT = 1 AND YEAR = 2022`); n != 1 {
		t.Errorf("rincian ringkasan lain berubah / terhapus (%d)", n)
	}
}
