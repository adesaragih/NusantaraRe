//go:build db

package repository_test

// Seam repository R/I Comm Life terhadap Oracle NYATA (skema uji; POOLDATA tidak pernah menjadi sasaran -
// `uji/skemauji`). Tanpa ORACLE_DSN MELEWATI. Tabel flat RICOMM_LIFE dari berkas migrasi inti 924 YANG SAMA dengan
// produksi; objek warisan ditiru di sini (bentuk dicek work owner 06-10-2026): M_RICOMM_LIFE_SUMMARY (JSON) + view
// RICOMM_LIFE_SUMMARY, M_RICOMM_LIFE (JSON), M_SITE_DATABASE, dan dua sequence warisan. Sesi sengaja ber-NLS koma:
// angka tidak boleh bergantung NLS. Fixture UJI-.
//
// ⚠️ Koneksi lewat `skemauji.BukaRepositori()`, BUKAN `skemauji.Buka()` (penjaga Claim Life
// `TestSetiapPemanggilBukaMemeriksaBolehDilewati` mengunci cacah pemanggilnya).

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/repository"
	"nusantarare/uji/skemauji"
)

type ujiDB struct {
	repo  *db.DB
	skema string
	ctx   context.Context
}

// ddl924 - pernyataan migrasi inti 924 (maju atau mundur), `{skema}` diganti skema uji.
func ddl924(t *testing.T, mundur bool, skema string) []string {
	t.Helper()
	sumber := fstest.MapFS{}
	for _, n := range []string{"924_ricomm_life.sql", "924_ricomm_life_down.sql"} {
		isi, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "inti", "backend", "migrations", n))
		if err != nil {
			t.Fatal(err)
		}
		sumber["migrations/"+n] = &fstest.MapFile{Data: isi}
	}
	l, err := migrasi.Daftar(mundur, sumber)
	if err != nil || len(l) != 1 {
		t.Fatalf("membaca 924: %v", err)
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
		`CREATE TABLE {s}.M_RICOMM_LIFE_SUMMARY (ID VARCHAR2(10), JSONDATA CLOB CONSTRAINT UJI_RICL_RS CHECK (JSONDATA IS JSON))`,
		`CREATE VIEW {s}.RICOMM_LIFE_SUMMARY AS SELECT a.ID, a.JSONDATA.USEDBY, a.JSONDATA.MODIFIEDDATE, a.JSONDATA.OPERATORID FROM {s}.M_RICOMM_LIFE_SUMMARY a`,
		`CREATE TABLE {s}.M_RICOMM_LIFE (ID VARCHAR2(10), JSONDATA CLOB CONSTRAINT UJI_RICL_R CHECK (JSONDATA IS JSON))`,
		`CREATE TABLE {s}.M_SITE_DATABASE (ID NUMBER(10), CURRENT_SITE VARCHAR2(1))`,
		`CREATE SEQUENCE {s}.M_RICOMM_LIFE_SUMMARY_SEQ START WITH 5`,
		`CREATE SEQUENCE {s}.M_RICOMM_LIFE_SEQ START WITH 44`,
		`INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (1, '1')`,
		`INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (2, '0')`,
		`INSERT INTO {s}.M_RICOMM_LIFE_SUMMARY (ID, JSONDATA) VALUES ('1000003', '{"MODIFIEDDATE":"20181205T073755.559 GMT","OPERATORID":"UJI-LAMA","pxObjClass":"ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY","USEDBY":"UJI COMM RETRO"}')`,
		`INSERT INTO {s}.M_RICOMM_LIFE (ID, JSONDATA) VALUES ('1000001', '{"IDUSEDBY":"1000003","USEDBY":"UJI COMM RETRO","CONTRACT":"1","YEAR":"2","COMM":"0.5"}')`,
	} {
		u.exec(t, q)
	}
	for _, q := range ddl924(t, false, u.skema) {
		u.exec(t, q)
	}
	u.exec(t, `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = ',.'`)
	t.Cleanup(func() {
		u.bongkar(t)
		_ = repo.Close()
	})
	return u
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

// Seam: situs + sequence warisan membentuk ID 7 karakter; ringkasan JSON sisip/ubah (pxObjClass dipertahankan);
// rincian flat tulis/baca angka tanpa NLS; Delete beserta rincian.
func TestSeamRepositoryOracle(t *testing.T) {
	u := pasang(t)
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
		if err := g.SisipRingkasan(u.ctx, tx, models.Ringkasan{ID: "1000005", UsedBy: "UJI COMM B", OperatorID: "UJI-A",
			ModifiedDate: "20261006T030405.600 GMT"}); err != nil {
			return err
		}
		if err := g.UbahRingkasan(u.ctx, tx, models.Ringkasan{ID: "1000003", UsedBy: "UJI COMM RETRO 2", OperatorID: "UJI-B",
			ModifiedDate: "20261006T030405.600 GMT"}); err != nil {
			return err
		}
		for _, k := range []models.Komisi{
			{ID: "1000044", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO 2", Contract: "1", Year: "2026", Comm: "0.5"},
			{ID: "1000045", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO 2", Contract: "99999", Year: "1", Comm: "123456789012345678901234567890.12345678"},
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
	var kelas string
	if err := u.repo.QueryRowContext(u.ctx, fmt.Sprintf(`SELECT JSON_VALUE(JSONDATA, '$.pxObjClass') FROM %s.M_RICOMM_LIFE_SUMMARY WHERE ID = '1000005'`, u.skema)).Scan(&kelas); err != nil ||
		kelas != models.KelasRingkasan {
		t.Errorf("pxObjClass %q %v", kelas, err)
	}
	d, total, err := g.DaftarKomisi(u.ctx, "1000003", 1)
	if err != nil || total != 2 || len(d) != 2 || d[0].Comm != "12.25" || d[0].Contract != "2" || d[1].Comm != "123456789012345678901234567890.12345678" {
		t.Errorf("rincian %+v %d %v", d, total, err)
	}
	if err := g.UbahKomisi(u.ctx, nil, models.Komisi{ID: "1000044", IDUsedBy: "1000005", Contract: "1", Year: "1", Comm: "1"}); !errors.Is(err, repository.ErrTidakAda) {
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
	// Situs tidak tepat satu baris.
	u.exec(t, `INSERT INTO {s}.M_SITE_DATABASE (ID, CURRENT_SITE) VALUES (3, '1')`)
	if _, err := g.Situs(u.ctx, nil); !errors.Is(err, repository.ErrSitus) {
		t.Errorf("dua situs aktif %v", err)
	}
}

// Alat pindah: uji kering membaca M_RICOMM_LIFE tanpa menulis; jalankan menyisipkan sekali, ulangan melewati.
func TestSeamPindahFlatOracle(t *testing.T) {
	u := pasang(t)
	g := repository.Baru(u.repo)
	lap, err := g.PindahFlat(u.ctx, false, false)
	if err != nil || lap.Sumber != 1 || lap.AkanDitulis != 1 || lap.Ditulis {
		t.Fatalf("uji kering %+v %v", lap, err)
	}
	if lap, err = g.PindahFlat(u.ctx, true, false); err != nil || !lap.Ditulis {
		t.Fatalf("jalankan %+v %v", lap, err)
	}
	if lap, err = g.PindahFlat(u.ctx, true, false); err != nil || lap.SudahSama != 1 || lap.AkanDitulis != 0 {
		t.Errorf("ulang %+v %v", lap, err)
	}
	d, _, err := g.DaftarKomisi(u.ctx, "1000003", 1)
	if err != nil || len(d) != 1 || d[0] != (models.Komisi{ID: "1000001", IDUsedBy: "1000003", UsedBy: "UJI COMM RETRO", Contract: "1", Year: "2", Comm: "0.5"}) {
		t.Errorf("hasil pindah %+v %v", d, err)
	}
	// Jalur mundur 924 memulihkan view atas M_RICOMM_LIFE.
	for _, q := range ddl924(t, true, u.skema) {
		u.exec(t, q)
	}
	var id string
	if err := u.repo.QueryRowContext(u.ctx, fmt.Sprintf(`SELECT ID FROM %s.RICOMM_LIFE WHERE IDUSEDBY = '1000003'`, u.skema)).Scan(&id); err != nil || id != "1000001" {
		t.Errorf("view dipulihkan %q %v", id, err)
	}
}
