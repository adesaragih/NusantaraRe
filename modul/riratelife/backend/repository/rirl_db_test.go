//go:build db

package repository_test

// Seam repository R/I Rate Life (RALAT R6: ringkasan satu tabel M_RATE_LIFE_SUMMARY) terhadap Oracle NYATA (skema uji;
// POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 927 ditiru:
// M_RATE_LIFE_SUMMARY (ID PK + JSONDATA ber-constraint IS JSON) dan tabel flat 926 dari berkas migrasi YANG SAMA;
// lalu 927 + 928 dijalankan dari berkasnya (termasuk pengulangan sesudah gagal di tengah) dan jalur mundurnya.
// M_RATE_LIFE (JSON) + view RATE_LIFE, dua sequence 923 ditiru. Fixture UJI-.
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
	"nusantarare/modul/riratelife/backend/models"
	"nusantarare/modul/riratelife/backend/repository"
	"nusantarare/modul/riratelife/backend/services"
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

// pasang - keadaan DEV sebelum 927: JSON ringkasan + tabel flat 926 sebagai sumber kebenaran. 100 sama di keduanya,
// 101 hanya di flat (baru aplikasi), 102 hanya di JSON (dihapus aplikasi).
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
		`CREATE TABLE {s}.M_RATE_LIFE_SUMMARY (ID VARCHAR2(10) PRIMARY KEY, JSONDATA CLOB NOT NULL
			CONSTRAINT ENSURE_M_RATE_LIFE_SUMMARY_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE TABLE {s}.M_RATE_LIFE (ID VARCHAR2(10), JSONDATA CLOB CONSTRAINT UJI_RIRL_R CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.RATE_LIFE AS SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
			a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM {s}.M_RATE_LIFE a`,
		`CREATE SEQUENCE {s}.SEQ_M_RATE_LIFE_SUMMARY START WITH 100`,
		`CREATE SEQUENCE {s}.SEQ_M_RATE_LIFE START WITH 9000`,
		`INSERT INTO {s}.M_RATE_LIFE_SUMMARY (ID, JSONDATA) VALUES ('100', '{"USEDBY":"UJI RATE A","TYPE":"L","MODIFIEDDATE":"20240102T030405.000 GMT","OPERATORID":"UJI-LAMA","FLAG":"AP"}')`,
		`INSERT INTO {s}.M_RATE_LIFE_SUMMARY (ID, JSONDATA) VALUES ('102', '{"USEDBY":"UJI DIHAPUS","FLAG":"PM"}')`,
	} {
		u.exec(t, q)
	}
	for _, q := range ddl(t, "926_rate_life_summary_flat", false, u.skema) {
		u.exec(t, q)
	}
	u.exec(t, `INSERT INTO {s}.RATE_LIFE_SUMMARY (ID, USEDBY, TYPE, MODIFIEDDATE, OPERATORID) VALUES ('100', 'UJI RATE A 2', 'L', '20261007T010000.000 GMT', 'UJI-APP')`)
	u.exec(t, `INSERT INTO {s}.RATE_LIFE_SUMMARY (ID, USEDBY, MODIFIEDDATE, OPERATORID) VALUES ('101', 'UJI RATE BARU', '20261007T020000.000 GMT', 'UJI-APP')`)
	t.Cleanup(func() {
		u.bongkar(t)
		_ = repo.Close()
	})
	return u
}

func (u *ujiDB) bongkar(t *testing.T) {
	t.Helper()
	for _, q := range []string{
		`DROP TABLE {s}.RATE_LIFE_SUMMARY PURGE`, `DROP VIEW {s}.RATE_LIFE_SUMMARY`, `DROP VIEW {s}.RATE_LIFE`,
		`DROP TABLE {s}.M_RATE_LIFE_SUMMARY PURGE`, `DROP TABLE {s}.M_RATE_LIFE PURGE`,
		`DROP SEQUENCE {s}.SEQ_M_RATE_LIFE_SUMMARY`, `DROP SEQUENCE {s}.SEQ_M_RATE_LIFE`,
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

func (u *ujiDB) cacah(t *testing.T, q string) int {
	t.Helper()
	var n int
	if err := u.repo.QueryRowContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema)).Scan(&n); err != nil {
		t.Fatalf("%.60s: %v", q, err)
	}
	return n
}

// satuTabel - 927 lalu 928; 928 sengaja dijalankan SEBAGIAN dua kali dulu (pelari gagal di tengah lalu diulang).
func (u *ujiDB) satuTabel(t *testing.T) {
	t.Helper()
	for _, q := range ddl(t, "927_m_rate_life_summary_kolom", false, u.skema) {
		u.exec(t, q)
	}
	l := ddl(t, "928_m_rate_life_summary_satu_tabel", false, u.skema)
	for _, q := range l[:4] {
		u.exec(t, q)
	}
	for _, q := range l {
		u.exec(t, q)
	}
}

// 928 dari sumber kebenaran flat: 100 diperbarui, 101 disisip, 102 (dihapus aplikasi) dibuang; JSONDATA, constraint
// IS JSON, dan tabel flat hilang; indeks nama baru ada. Jalur mundur 928 + 927 membangun ulang flat dan JSONDATA.
func TestDBMigrasiSatuTabelDanMundur(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY WHERE (ID = '100' AND USEDBY = 'UJI RATE A 2' AND TYPE = 'L'
		AND OPERATORID = 'UJI-APP') OR (ID = '101' AND USEDBY = 'UJI RATE BARU' AND TYPE IS NULL)`); n != 2 {
		t.Errorf("isi sesudah 928: %d", n)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY`); n != 2 {
		t.Errorf("102 harus dibuang (%d baris)", n)
	}
	for q, mau := range map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RATE_LIFE_SUMMARY'`:                              5,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RATE_LIFE_SUMMARY' AND COLUMN_NAME = 'JSONDATA'`: 0,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'ENSURE_M_RATE_LIFE_SUMMARY_JSON'`:             0,
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME = 'RATE_LIFE_SUMMARY'`:                                   0,
		`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER('{s}') AND INDEX_NAME = 'IX_M_RATE_LIFE_SUMMARY_NAMA'`:                          1,
	} {
		if n := u.cacah(t, q); n != mau {
			t.Errorf("%d, mau %d: %s", n, mau, q)
		}
	}
	for _, kunci := range []string{"928_m_rate_life_summary_satu_tabel", "927_m_rate_life_summary_kolom"} {
		for _, q := range ddl(t, kunci, true, u.skema) {
			u.exec(t, q)
		}
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID IN ('100', '101')`); n != 2 {
		t.Errorf("flat dipulihkan %d", n)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY m WHERE m.ID = '100'
		AND JSON_VALUE(m.JSONDATA, '$.USEDBY') = 'UJI RATE A 2' AND JSON_VALUE(m.JSONDATA, '$.TYPE') = 'L'`); n != 1 {
		t.Error("JSONDATA dibangun ulang dari kolom")
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

// Rumus ID: SEQ_M_RATE_LIFE_SUMMARY mulai 100 - 100 dan 101 terpakai di M_RATE_LIFE_SUMMARY -> ID baru 102. Pulang-
// pergi kolom: Add menulis USEDBY/OPERATORID/MODIFIEDDATE, TYPE NULL; Edit tidak menyentuh TYPE.
func TestDBRumusIDDanPulangPergi(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	l := layananOracle(u)
	r, err := l.Simpan(u.ctx, penuhDB, "", models.Isian{UsedBy: "UJI RATE DB"})
	if err != nil || r.ID != "102" || r.OperatorID != "UJI-ADMIN" || r.UsedBy != "UJI RATE DB" {
		t.Fatalf("Add %+v %v", r, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY WHERE ID = '102' AND TYPE IS NULL`); n != 1 {
		t.Errorf("TYPE baris baru harus NULL (%d)", n)
	}
	if _, err := l.Simpan(u.ctx, penuhDB, "100", models.Isian{UsedBy: "UJI RATE A 3"}); err != nil {
		t.Fatal(err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY WHERE ID = '100' AND USEDBY = 'UJI RATE A 3' AND TYPE = 'L'`); n != 1 {
		t.Error("Edit mengubah TYPE atau tidak menulis nama")
	}
	if _, err := l.Simpan(u.ctx, penuhDB, "", models.Isian{UsedBy: " uji rate a 3 "}); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("nama kembar lewat indeks nama: %v", err)
	}
}

// Transaksi: gagal di tengah (ID ringkasan kembar, ORA-00001) membatalkan ringkasan + rate; Simpan Upload sah lalu
// Delete berantai menghapus ringkasan beserta rate-nya.
func TestDBTransaksiDanDeleteBerantai(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	g := repository.Baru(u.repo)
	tx, err := u.repo.Mulai(u.ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = g.SisipRingkasan(u.ctx, tx, models.Ringkasan{ID: "300", UsedBy: "UJI GAGAL"})
	if err == nil {
		err = g.SisipRate(u.ctx, tx, models.Rate{ID: "9300", IDUsedBy: "300", UsedBy: "UJI GAGAL", Gender: "U", Contract: "1", Age: "30", Rate: "0,5"})
	}
	if err == nil {
		err = g.SisipRingkasan(u.ctx, tx, models.Ringkasan{ID: "300", UsedBy: "UJI GAGAL"})
	}
	_ = tx.Rollback()
	if !errors.Is(err, repository.ErrKembar) {
		t.Fatalf("mau ErrKembar, dapat %v", err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY WHERE ID = '300'`) + u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE ID = '9300'`); n != 0 {
		t.Errorf("rollback meninggalkan %d baris", n)
	}
	l := layananOracle(u)
	h, err := l.SimpanUnggah(u.ctx, penuhDB, services.PermintaanUnggah{CSV: "USEDBY;CONTRACT;GENDER;AGE;RATE\nUJI DEL;1;U;30;0,5\nUJI DEL;2;U;30;0,6\n"})
	if err != nil || h.Disimpan != 2 || h.RingkasanBaru != 1 {
		t.Fatalf("unggah %+v %v", h, err)
	}
	id := h.Ringkasan[0].ID
	hh, err := l.Hapus(u.ctx, penuhDB, id)
	if err != nil || hh.RateTerhapus != 2 {
		t.Fatalf("hapus %+v %v", hh, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY WHERE ID = '`+id+`'`) +
		u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE WHERE IDUSEDBY = '`+id+`'`); n != 0 {
		t.Errorf("Delete berantai menyisakan %d baris", n)
	}
}
