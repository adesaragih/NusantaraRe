//go:build db

package repository_test

// Seam repository R/I Rate Life (RALAT R6: ringkasan satu tabel M_RATE_LIFE_SUMMARY) terhadap Oracle NYATA (skema uji;
// POOLDATA tidak pernah menjadi sasaran - `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Keadaan DEV sebelum 927 ditiru:
// M_RATE_LIFE_SUMMARY (ID PK + JSONDATA ber-constraint IS JSON) dan tabel flat 926 dari berkas migrasi YANG SAMA;
// lalu 927 + 928 dijalankan dari berkasnya (termasuk pengulangan sesudah gagal di tengah) dan jalur mundurnya.
// M_RATE_LIFE (ID PK + JSONDATA ber-constraint IS JSON) + view RATE_LIFE ditiru seperti DEV, lalu 929 + 930 (RALAT R7)
// dijalankan dari berkasnya; dua sequence 923 ditiru. Fixture UJI-.
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
		// Bentuk DEV yang tercatat (fakta WO 07-10-2026): ID PK + JSONDATA CLOB ber-constraint IS JSON. NOT NULL pada
		// JSONDATA TIDAK terbukti di repo, jadi tidak ditiru (LANGKAH-WO (a) mencatat NULLABLE aslinya).
		`CREATE TABLE {s}.M_RATE_LIFE_SUMMARY (ID VARCHAR2(10) PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_RATE_LIFE_SUMMARY_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE TABLE {s}.M_RATE_LIFE (ID VARCHAR2(10) PRIMARY KEY, JSONDATA CLOB
			CONSTRAINT ENSURE_M_RATE_LIFE_JSON CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.RATE_LIFE AS SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
			a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM {s}.M_RATE_LIFE a`,
		`CREATE SEQUENCE {s}.SEQ_M_RATE_LIFE_SUMMARY START WITH 100`,
		`CREATE SEQUENCE {s}.SEQ_M_RATE_LIFE START WITH 9000`,
		// Rincian DEV: RATE berkoma / bertitik apa adanya, FLAG di JSON (tidak dipindah), 8999 yatim (ringkasan 999 tidak ada).
		`INSERT INTO {s}.M_RATE_LIFE (ID, JSONDATA) VALUES ('8997', '{"IDUSEDBY":"100","USEDBY":"UJI RATE A","GENDER":"U","CONTRACT":"1","AGE":"30","RATE":"0,5","FLAG":"AP"}')`,
		`INSERT INTO {s}.M_RATE_LIFE (ID, JSONDATA) VALUES ('8998', '{"IDUSEDBY":"100","USEDBY":"UJI RATE A","GENDER":"M","CONTRACT":"2","RATE":"1.25"}')`,
		`INSERT INTO {s}.M_RATE_LIFE (ID, JSONDATA) VALUES ('8999', '{"IDUSEDBY":"999","USEDBY":"UJI YATIM","GENDER":"F","CONTRACT":"3","AGE":"40","RATE":"2"}')`,
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

// rincianFlat - 929 lalu 930; 930 sengaja dijalankan SEBAGIAN dua kali dulu (blok isi + blok buang, lalu penuh).
func (u *ujiDB) rincianFlat(t *testing.T) {
	t.Helper()
	for _, q := range ddl(t, "929_m_rate_life_kolom", false, u.skema) {
		u.exec(t, q)
	}
	l := ddl(t, "930_m_rate_life_satu_tabel", false, u.skema)
	for _, q := range l[:2] {
		u.exec(t, q)
	}
	for _, q := range l {
		u.exec(t, q)
	}
}

// 930: kolom berisi nilai JSON APA ADANYA (RATE berkoma tetap berkoma), baris yatim ikut, FLAG tidak dipindah,
// JSONDATA + constraint + view hilang, indeks ada; jalur mundur 930 + 929 membangun ulang JSONDATA dan view.
func TestDBMigrasiRincianFlatDanMundur(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	u.rincianFlat(t)
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE (ID = '8997' AND IDUSEDBY = '100' AND RATE = '0,5' AND AGE = '30')
		OR (ID = '8998' AND RATE = '1.25' AND AGE IS NULL) OR (ID = '8999' AND IDUSEDBY = '999' AND GENDER = 'F')`); n != 3 {
		t.Errorf("isi rincian apa adanya %d", n)
	}
	for q, mau := range map[string]int{
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RATE_LIFE'`:                              8,
		`SELECT COUNT(*) FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = UPPER('{s}') AND TABLE_NAME = 'M_RATE_LIFE' AND COLUMN_NAME = 'JSONDATA'`: 0,
		`SELECT COUNT(*) FROM SYS.ALL_CONSTRAINTS WHERE OWNER = UPPER('{s}') AND CONSTRAINT_NAME = 'ENSURE_M_RATE_LIFE_JSON'`:             0,
		`SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE OWNER = UPPER('{s}') AND OBJECT_NAME = 'RATE_LIFE'`:                                   0,
		`SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = UPPER('{s}') AND INDEX_NAME = 'IX_M_RATE_LIFE_IDUSEDBY'`:                      1,
	} {
		if n := u.cacah(t, q); n != mau {
			t.Errorf("%d, mau %d: %s", n, mau, q)
		}
	}
	// 930 dibuang JSONDATA-nya tetapi (seolah) tidak tercatat: 929_down langsung TIDAK BOLEH membuang kolom.
	pelindung := ddl(t, "929_m_rate_life_kolom", true, u.skema)[0]
	u.gagalORA(t, pelindung, "ORA-00904")
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE RATE IS NOT NULL`); n != 3 {
		t.Errorf("pelindung 929_down: kolom rincian harus utuh (%d)", n)
	}
	// 930_down aman diulang: dua blok pertama dua kali, lalu penuh.
	turun := ddl(t, "930_m_rate_life_satu_tabel", true, u.skema)
	for _, q := range turun[:2] {
		u.exec(t, q)
	}
	for _, kunci := range []string{"930_m_rate_life_satu_tabel", "929_m_rate_life_kolom"} {
		for _, q := range ddl(t, kunci, true, u.skema) {
			u.exec(t, q)
		}
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE WHERE ID = '8997' AND IDUSEDBY = '100' AND RATE = '0,5'`); n != 1 {
		t.Error("view RATE_LIFE dan JSONDATA dibangun ulang dari kolom")
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
	u.gagalORA(t, ddl(t, "927_m_rate_life_summary_kolom", true, u.skema)[0], "ORA-00904")
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY WHERE USEDBY IS NOT NULL`); n != 2 {
		t.Errorf("pelindung 927_down: kolom ringkasan harus utuh (%d)", n)
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
	u.rincianFlat(t)
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
	u.rincianFlat(t)
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
		u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE IDUSEDBY = '`+id+`'`); n != 0 {
		t.Errorf("Delete berantai menyisakan %d baris", n)
	}
}

// Rate milik ringkasan LAIN: Edit lewat ringkasan 101 atas rate 8997 (milik 100) ditolak ErrRateTidakAda (layanan) dan
// ErrTidakAda (repository, WHERE ID AND IDUSEDBY); Delete ringkasan 101 tidak menyentuh rate milik 100. Isi 8997 utuh.
func TestDBRateRingkasanLainDitolak(t *testing.T) {
	u := pasang(t)
	u.satuTabel(t)
	u.rincianFlat(t)
	utuh := `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE ID = '8997' AND IDUSEDBY = '100' AND GENDER = 'U'
		AND CONTRACT = '1' AND AGE = '30' AND RATE = '0,5'`
	l := layananOracle(u)
	if _, err := l.SimpanRate(u.ctx, penuhDB, "101", "8997", models.IsianRate{Gender: "M", Contract: "9", Age: "50", Rate: "7"}); !errors.Is(err, services.ErrRateTidakAda) {
		t.Errorf("SimpanRate rate ringkasan lain: mau ErrRateTidakAda, dapat %v", err)
	}
	g := repository.Baru(u.repo)
	tx, err := u.repo.Mulai(u.ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = g.UbahRate(u.ctx, tx, models.Rate{ID: "8997", IDUsedBy: "101", Gender: "M", Contract: "9", Age: "50", Rate: "7"})
	_ = tx.Rollback()
	if !errors.Is(err, repository.ErrTidakAda) {
		t.Errorf("UbahRate rate ringkasan lain: mau ErrTidakAda, dapat %v", err)
	}
	h, err := l.Hapus(u.ctx, penuhDB, "101")
	if err != nil || h.RateTerhapus != 0 {
		t.Errorf("Delete ringkasan 101 %+v %v", h, err)
	}
	if n := u.cacah(t, utuh) + u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE IDUSEDBY = '100'`); n != 3 {
		t.Errorf("rate milik 100 berubah / terhapus (%d)", n)
	}
}
