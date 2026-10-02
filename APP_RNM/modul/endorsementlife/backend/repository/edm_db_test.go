//go:build db

package repository_test

// Uji terhadap skema uji Oracle (`uji/skemauji`) - migrasi 480/481 dan
// pembacaan modul Endorsement Life. Tanpa ORACLE_DSN seluruh test di sini
// MELEWATI dengan pesan.
//
// ⛔ Tabel warisan `JSON_POLIS` tidak dibuat skema uji bersama; test ini
// membuat tiruannya sendiri (lima kolom yang dibaca) dan membuangnya lagi.
// Fixture berawalan UJI-, nol nama orang.

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
	"nusantarare/uji/skemauji"
)

// pasangSkema membuka skema uji dan menjalankan SELURUH migrasi terdaftar
// (termasuk 051-056 PremiumList dan 480/481/976 modul ini), atau melewati test.
//
// ⚠️ Satu koneksi saja - `skemauji.BukaRepositori()`: pemanggil
// fungsi `Buka` skemauji dicacah penjaga Claim Life
// (`TestSetiapPemanggilBukaMemeriksaBolehDilewati`), dan pelari migrasi
// `inti/backend/migrasi` cukup dengan koneksi repository. `BolehDilewati`
// tetap diperiksa.
func pasangSkema(t *testing.T) (*intidb.DB, string, func()) {
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
		_ = repo.Close()
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	skema := repo.Skema()
	bongkar := func() {
		_, _ = repo.ExecContext(ctx, `DROP TABLE `+skema+`.JSON_POLIS PURGE`)
		_, _ = migrasi.Bongkar(ctx, repo, skemauji.SumberMigrasi()...)
		_, _ = repo.ExecContext(ctx, `DROP TABLE `+skema+`.T_MIGRASI CASCADE CONSTRAINTS`)
	}
	bongkar()
	if _, err := migrasi.Jalankan(ctx, repo, skemauji.SumberMigrasi()...); err != nil {
		t.Fatalf("menjalankan migrasi: %v", err)
	}
	tiruanJSON := `CREATE TABLE ` + skema + `.JSON_POLIS (IDPEGA VARCHAR2(64), NOPOLIS VARCHAR2(255),
	    PRODKE NUMBER, TGL_INPUT DATE, DATA_JSON CLOB)`
	if _, err := repo.ExecContext(ctx, tiruanJSON); err != nil && !strings.Contains(err.Error(), "ORA-00955") {
		t.Fatalf("tiruan JSON_POLIS: %v", err)
	}
	return repo, skema, func() {
		bongkar()
		_ = repo.Close()
	}
}

func jalankan(t *testing.T, repo *intidb.DB, q ...string) {
	t.Helper()
	for _, s := range q {
		if _, err := repo.ExecContext(context.Background(), s); err != nil {
			t.Fatalf("%.80s: %v", s, err)
		}
	}
}

// TestMigrasi480ProdKeBawaanSatuDanIndexVersi - E1: baris new business yang
// `INSERT`-nya tidak menyebut `PROD_KE` (bentuk PremiumList) menjadi versi 1;
// ketiga index pencari versi ada; sequence kasus menerbitkan angka.
func TestMigrasi480ProdKeBawaanSatuDanIndexVersi(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	jalankan(t, repo, `INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT) VALUES ('UJI-NB-1', 'UJI-NB-1', SYSDATE)`)
	var prodKe sql.NullInt64
	if err := repo.QueryRowContext(ctx, `SELECT PROD_KE FROM `+skema+`.T_PREMIUM_LIST WHERE ID = 'UJI-NB-1'`).Scan(&prodKe); err != nil {
		t.Fatal(err)
	}
	if !prodKe.Valid || prodKe.Int64 != 1 {
		t.Fatalf("PROD_KE baris NB = %v, mau 1 (DEFAULT 1, ralat E1)", prodKe)
	}
	var n int
	if err := repo.QueryRowContext(ctx, `SELECT COUNT(*) FROM SYS.ALL_INDEXES WHERE OWNER = :1
	    AND INDEX_NAME IN ('IDX_PL_NOPOLIS_PRODKE', 'IDX_PL_OLD_POLICY_NO', 'IDX_PLD_PL_NUMBER')`, strings.ToUpper(skema)).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Fatalf("%d dari 3 index pencari versi", n)
	}
	tx, err := repo.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	urut, err := repo.NomorBerikut(ctx, tx, "SEQ_WORK_EDM_LIFE")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := models.RakitPengenalKasus(urut); err != nil {
		t.Fatal(err)
	}
}

// TestBacaKasusDanVersiTerhadapOracle - kotak masuk, kepala, peserta, dan
// versi berjalan dari kedua sumber.
func TestBacaKasusDanVersiTerhadapOracle(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	jalankan(t, repo,
		// Versi NB sistem baru: PROD_KE kosong di INSERT → 1.
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, TYPE) VALUES ('UJI-NB-2', 'UJI-NB-2', SYSDATE, 'QR')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, CERTIFICATE_NO) VALUES ('UJI-D-NB', 'UJI-NB-2', 'UJI-PL-2', 'UJI-C1')`,
		// Kasus EDM terbuka atas polis itu.
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, OLD_POLICY_NO, EDM_TYPE, PROD_KE, TYPE)
		   VALUES ('EDMLF-901', 'EDMLF-901', SYSDATE, 'UJI-PL-2', '1', 2, 'QR')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PARENT_ID, PL_NUMBER, CERTIFICATE_NO, EDM_STATUS, SUM_INSURED)
		   VALUES ('UJI-D-EDM', 'EDMLF-901', 'UJI-D-NB', 'UJI-PL-2', 'UJI-C1', 'Old', 0.5)`,
		// Polis warisan: dua versi, yang kedua Batal.
		`INSERT INTO `+skema+`.JSON_POLIS (IDPEGA, NOPOLIS, PRODKE, TGL_INPUT, DATA_JSON)
		   VALUES ('UJI-IDPEGA-1', 'UJI-PL-3', NULL, SYSDATE - 2, '{"EdmType":""}')`,
		`INSERT INTO `+skema+`.JSON_POLIS (IDPEGA, NOPOLIS, PRODKE, TGL_INPUT, DATA_JSON)
		   VALUES ('UJI-IDPEGA-2', 'UJI-PL-3', 2, SYSDATE - 1, '{"EdmType":"3"}')`,
	)
	g := repository.Baru(repo)
	baris, total, err := g.Inbox(ctx, 1, 20)
	if err != nil || total != 1 || len(baris) != 1 || baris[0].CaseID != "EDMLF-901" || baris[0].EndorsementNo != "UJI-PL-2" {
		t.Fatalf("inbox = %+v, %d, %v", baris, total, err)
	}
	k, err := g.AmbilKasus(ctx, nil, "EDMLF-901", false)
	if err != nil || k.ProdKe != 2 || k.Kepala["TYPE"] != "QR" {
		t.Fatalf("kasus = %+v, %v", k, err)
	}
	if _, err := g.AmbilKasus(ctx, nil, "UJI-NB-2", false); err != repository.ErrTidakAda {
		t.Errorf("baris NB terbaca sebagai kasus: %v", err)
	}
	p, err := g.DaftarPeserta(ctx, "EDMLF-901", 1, 20)
	if err != nil || len(p) != 1 || p[0].ParentID != "UJI-D-NB" || p[0].Nilai["SUM_INSURED"] != "0.5" {
		t.Fatalf("peserta = %+v, %v", p, err)
	}
	v, ada, err := g.VersiBerjalan(ctx, nil, "UJI-PL-2", 0)
	if err != nil || !ada || v.ID != "UJI-NB-2" || v.Jenis != models.SumberAplikasi || v.ProdKe != 1 {
		t.Fatalf("versi NB = %+v, %v, %v (kasus terbuka bukan versi)", v, ada, err)
	}
	v, ada, err = g.VersiBerjalan(ctx, nil, "UJI-PL-3", 0)
	if err != nil || !ada || v.ID != "UJI-IDPEGA-2" || v.ProdKe != 2 || !v.SudahBatal() {
		t.Fatalf("versi warisan = %+v, %v, %v", v, ada, err)
	}
	v, _, err = g.VersiBerjalan(ctx, nil, "UJI-PL-3", 2)
	if err != nil || v.ID != "UJI-IDPEGA-1" || v.ProdKe != 1 {
		t.Fatalf("versi sebelum 2 = %+v, %v (PRODKE kosong dibaca 1)", v, err)
	}
	if _, ada, _ = g.VersiBerjalan(ctx, nil, "UJI-PL-TIDAK-ADA", 0); ada {
		t.Error("polis tak dikenal ditemukan")
	}
}
