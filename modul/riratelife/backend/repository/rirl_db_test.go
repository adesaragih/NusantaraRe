//go:build db

package repository_test

// Seam repository R/I Rate Life (RALAT R4, tabel flat RATE_LIFE_SUMMARY) terhadap Oracle NYATA (skema uji; POOLDATA
// tidak pernah menjadi sasaran - `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Tabel flat dari berkas migrasi inti 926
// YANG SAMA dengan produksi; objek warisan ditiru di sini (bentuk katalog DEV): M_RATE_LIFE_SUMMARY (JSON),
// M_RATE_LIFE (JSON) + view RATE_LIFE, dan dua sequence 923. Fixture UJI-.
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

// ddl926 - pernyataan migrasi inti 926 (maju atau mundur), `{skema}` diganti skema uji.
func ddl926(t *testing.T, mundur bool, skema string) []string {
	t.Helper()
	sumber := fstest.MapFS{}
	for _, n := range []string{"926_rate_life_summary_flat.sql", "926_rate_life_summary_flat_down.sql"} {
		isi, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", n))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+n] = &fstest.MapFile{Data: isi}
	}
	l, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(l) != 1 {
		t.Fatalf("membaca 926: %v", err)
	}
	var out []string
	for _, q := range l[0].Pernyataan {
		out = append(out, strings.ReplaceAll(q, "{skema}", skema))
	}
	return out
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
	for _, q := range []string{
		`CREATE TABLE {s}.M_RATE_LIFE_SUMMARY (ID VARCHAR2(10), JSONDATA CLOB CONSTRAINT UJI_RIRL_RS CHECK (JSONDATA IS JSON))`,
		`CREATE TABLE {s}.M_RATE_LIFE (ID VARCHAR2(10), JSONDATA CLOB CONSTRAINT UJI_RIRL_R CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.RATE_LIFE AS SELECT a.ID, a.JSONDATA.IDUSEDBY, a.JSONDATA.USEDBY, a.JSONDATA.TYPE, a.JSONDATA.GENDER,
			a.JSONDATA.CONTRACT, a.JSONDATA.AGE, a.JSONDATA.RATE FROM {s}.M_RATE_LIFE a`,
		`CREATE SEQUENCE {s}.SEQ_M_RATE_LIFE_SUMMARY START WITH 100`,
		`CREATE SEQUENCE {s}.SEQ_M_RATE_LIFE START WITH 9000`,
		// 100 ada di JSON warisan (belum dipindah): ID baru harus melewatinya.
		`INSERT INTO {s}.M_RATE_LIFE_SUMMARY (ID, JSONDATA) VALUES ('100', '{"USEDBY":"UJI RATE LAMA","TYPE":"L","MODIFIEDDATE":"20240102T030405.000 GMT","OPERATORID":"UJI-LAMA","FLAG":"AP"}')`,
	} {
		u.exec(t, q)
	}
	for _, q := range ddl926(t, false, u.skema) {
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
		t.Fatalf("%.50s: %v", q, err)
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

// Rumus ID: SEQ_M_RATE_LIFE_SUMMARY mulai 100, tetapi 100 terpakai di JSON warisan -> ID baru 101. Pulang-pergi tabel
// flat: Add menulis USEDBY/OPERATORID/MODIFIEDDATE, TYPE NULL; Edit tidak menyentuh TYPE. Tabel flat tanpa FLAG.
func TestDBRumusIDDanPulangPergi(t *testing.T) {
	u := pasang(t)
	l := layananOracle(u)
	r, err := l.Simpan(u.ctx, penuhDB, "", models.Isian{UsedBy: "UJI RATE DB"})
	if err != nil || r.ID != "101" || r.OperatorID != "UJI-ADMIN" || r.UsedBy != "UJI RATE DB" {
		t.Fatalf("Add %+v %v", r, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '101' AND TYPE IS NULL`); n != 1 {
		t.Errorf("TYPE baris baru harus NULL (%d)", n)
	}
	u.exec(t, `INSERT INTO {s}.RATE_LIFE_SUMMARY (ID, USEDBY, TYPE) VALUES ('200', 'UJI PINDAHAN', 'L')`)
	if _, err := l.Simpan(u.ctx, penuhDB, "200", models.Isian{UsedBy: "UJI PINDAHAN 2"}); err != nil {
		t.Fatal(err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '200' AND USEDBY = 'UJI PINDAHAN 2' AND TYPE = 'L'`); n != 1 {
		t.Error("Edit mengubah TYPE atau tidak menulis nama")
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE_SUMMARY`); n != 1 {
		t.Errorf("M_RATE_LIFE_SUMMARY berubah (%d baris)", n)
	}
}

// Transaksi: Simpan Upload sah menulis ringkasan + rate; gagal di tengah (ID ringkasan kembar, ORA-00001) membatalkan
// keduanya.
func TestDBTransaksiDanDeleteBerantai(t *testing.T) {
	u := pasang(t)
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
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '300'`) + u.cacah(t, `SELECT COUNT(*) FROM {s}.M_RATE_LIFE WHERE ID = '9300'`); n != 0 {
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
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '`+id+`'`) +
		u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE WHERE IDUSEDBY = '`+id+`'`); n != 0 {
		t.Errorf("Delete berantai menyisakan %d baris", n)
	}
}

// Alat pindah: uji kering tanpa tulis; jalankan memindah apa adanya sekali; ulangan melewati; jalur mundur 926
// memulihkan view atas M_RATE_LIFE_SUMMARY. ALTER SESSION dan transaksi di koneksi yang sama.
func TestDBPindahFlat(t *testing.T) {
	u := pasang(t)
	g := repository.Baru(u.repo)
	kon, tx, err := g.SesiPindah(u.ctx)
	if err != nil {
		t.Fatal(err)
	}
	var sidTx, sidKon string
	_ = tx.QueryRowContext(u.ctx, `SELECT SYS_CONTEXT('USERENV', 'SID') FROM DUAL`).Scan(&sidTx)
	_ = kon.QueryRowContext(u.ctx, `SELECT SYS_CONTEXT('USERENV', 'SID') FROM DUAL`).Scan(&sidKon)
	_ = tx.Rollback()
	_ = kon.Close()
	if sidTx == "" || sidTx != sidKon {
		t.Errorf("SID transaksi %q, koneksi %q", sidTx, sidKon)
	}
	lap, err := g.PindahFlat(u.ctx, false, "")
	if err != nil || lap.CacahSumber != 1 || len(lap.Baru) != 1 || lap.Ditulis || lap.PanjangMaks["MODIFIEDDATE"] != 23 {
		t.Fatalf("uji kering %+v %v", lap, err)
	}
	if lap, err = g.PindahFlat(u.ctx, true, ""); err != nil || !lap.Ditulis {
		t.Fatalf("jalankan %+v %v", lap, err)
	}
	if lap, err = g.PindahFlat(u.ctx, true, ""); err != nil || len(lap.Sama) != 1 || len(lap.Baru)+len(lap.Berubah) != 0 {
		t.Errorf("ulang %+v %v", lap, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '100' AND USEDBY = 'UJI RATE LAMA' AND TYPE = 'L'
		AND MODIFIEDDATE = '20240102T030405.000 GMT' AND OPERATORID = 'UJI-LAMA'`); n != 1 {
		t.Error("pindahan tidak apa adanya")
	}
	for _, q := range ddl926(t, true, u.skema) {
		u.exec(t, q)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '100' AND TYPE = 'L'`); n != 1 {
		t.Error("view tidak dipulihkan")
	}
}

// Delta (cutover): sesudah pemindahan penuh, Pega mengubah 100 (lebih baru) dan menulis 110 baru; aplikasi menghapus
// ringkasan 120 (ada di JSON, lebih lama dari -sejak) dan mengubah 130 sesudah Pega - putaran delta memperbarui 100,
// menyisip 110, melewati 120, dan melaporkan 130 sebagai konflik tanpa menimpanya. TYPE ikut apa adanya; FLAG di JSON
// tidak dipindah (RALAT R5).
func TestDBPindahDelta(t *testing.T) {
	u := pasang(t)
	g := repository.Baru(u.repo)
	u.exec(t, `INSERT INTO {s}.M_RATE_LIFE_SUMMARY (ID, JSONDATA) VALUES ('120', '{"USEDBY":"UJI DIHAPUS","MODIFIEDDATE":"20240102T030405.000 GMT","FLAG":"PY"}')`)
	u.exec(t, `INSERT INTO {s}.M_RATE_LIFE_SUMMARY (ID, JSONDATA) VALUES ('130', '{"USEDBY":"UJI KONFLIK","MODIFIEDDATE":"20240102T030405.000 GMT","FLAG":"PM"}')`)
	lap, err := g.PindahFlat(u.ctx, true, "")
	if err != nil || len(lap.Baru) != 3 {
		t.Fatalf("penuh %+v %v", lap, err)
	}
	batas := lap.BatasBerikut
	u.exec(t, `DELETE FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '120'`)
	u.exec(t, `UPDATE {s}.RATE_LIFE_SUMMARY SET USEDBY = 'UJI KONFLIK APLIKASI', MODIFIEDDATE = '20261006T050000.000 GMT' WHERE ID = '130'`)
	u.exec(t, `UPDATE {s}.M_RATE_LIFE_SUMMARY SET JSONDATA = '{"USEDBY":"UJI RATE LAMA 2","TYPE":"L","MODIFIEDDATE":"20261006T040628.169 GMT","OPERATORID":"UJI-PEGA","FLAG":"PM"}' WHERE ID = '100'`)
	u.exec(t, `INSERT INTO {s}.M_RATE_LIFE_SUMMARY (ID, JSONDATA) VALUES ('110', '{"USEDBY":"UJI PEGA BARU","MODIFIEDDATE":"20261006T040628.169 GMT","FLAG":"AP"}')`)
	lap, err = g.PindahFlat(u.ctx, true, batas)
	if err != nil || !lap.Ditulis || len(lap.Baru) != 1 || len(lap.Berubah) != 1 || len(lap.Dilewati) != 1 || len(lap.Konflik) != 1 ||
		lap.CacahSumber != 4 || lap.CacahFlatSesudah != 3 {
		t.Fatalf("delta %+v %v", lap, err)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE (ID = '100' AND USEDBY = 'UJI RATE LAMA 2' AND TYPE = 'L')
		OR ID = '110' OR (ID = '130' AND USEDBY = 'UJI KONFLIK APLIKASI')`); n != 3 {
		t.Errorf("hasil delta %d", n)
	}
	if n := u.cacah(t, `SELECT COUNT(*) FROM {s}.RATE_LIFE_SUMMARY WHERE ID = '120'`); n != 0 {
		t.Error("ringkasan yang dihapus aplikasi hidup lagi")
	}
}
