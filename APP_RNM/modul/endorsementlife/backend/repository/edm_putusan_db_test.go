//go:build db

package repository_test

// Keputusan terhadap Oracle: riwayat, peresmian versi, anti-dobel, rekap
// warisan `M_LIFE_PREMIUM_SUMMARY`, penolakan. Fixture UJI-, nol nama orang.
//
// ⚠️ Tiruan `M_LIFE_PREMIUM_SUMMARY` dibuat bila belum ada (bentuk kolom dari
// `repository.KolomRekapWarisan`: uang NUMBER, teks VARCHAR2) dan HANYA yang
// dibuat uji ini yang dibuang lagi - tiruan bersama `skemauji` tidak disentuh.

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

func TestPutuskanTerhadapOracle(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	var kolom []string
	for _, k := range repository.KolomRekapWarisan {
		tipe := "NUMBER"
		switch k {
		case "COB", "CURRENCY", "PL_NUMBER", "PL_NUMBER_EDM", "IDPEGA":
			tipe = "VARCHAR2(255)"
		}
		kolom = append(kolom, k+" "+tipe)
	}
	dibuat := false
	if _, err := repo.ExecContext(ctx, `CREATE TABLE `+skema+`.M_LIFE_PREMIUM_SUMMARY (ID VARCHAR2(50), `+strings.Join(kolom, ", ")+`)`); err == nil {
		dibuat = true
		jalankan(t, repo, `CREATE SEQUENCE `+skema+`.M_LIFE_PREMIUM_SUMMARY_SEQ`)
	} else if !strings.Contains(err.Error(), "ORA-00955") {
		t.Fatal(err)
	}
	defer func() {
		if dibuat {
			_, _ = repo.ExecContext(ctx, `DROP TABLE `+skema+`.M_LIFE_PREMIUM_SUMMARY PURGE`)
			_, _ = repo.ExecContext(ctx, `DROP SEQUENCE `+skema+`.M_LIFE_PREMIUM_SUMMARY_SEQ`)
		} else {
			_, _ = repo.ExecContext(ctx, `DELETE FROM `+skema+`.M_LIFE_PREMIUM_SUMMARY WHERE IDPEGA = 'EDMLF-903'`)
		}
	}()
	jalankan(t, repo,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, ID_PEGA, TGL_INPUT, TYPE) VALUES ('UJI-NB-P', 'UJI-NB-P', SYSDATE, 'TR')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER) VALUES ('UJI-NBD-P', 'UJI-NB-P', 'UJI-PL-P')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, OLD_POLICY_NO, TGL_INPUT, TYPE, EDM_TYPE, PROD_KE, BUSINESS_NAME)
		   VALUES ('EDMLF-903', 'UJI-PL-P', SYSDATE, 'TR', '1', 2, 'UJI-COB')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM_REFUND_RETRO)
		   VALUES ('UJI-P1', 'EDMLF-903', 'UJI-PL-P', 'Old', 'IDR', 10.5)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM_REFUND_RETRO)
		   VALUES ('UJI-P2', 'EDMLF-903', 'UJI-PL-P', 'New', 'IDR', 2)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, OLD_POLICY_NO, TGL_INPUT, TYPE, EDM_TYPE, PROD_KE) VALUES ('EDMLF-904', 'UJI-PL-Q', SYSDATE, 'QR', '1', 2)`,
	)
	g := repository.Baru(repo)
	tx, err := repo.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	gagal := func(err error) {
		t.Helper()
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if _, err := g.HitungRekap(ctx, tx, "EDMLF-903", models.TypeTR); err != nil {
		gagal(err)
	}
	for _, s := range []string{"Decline", "Accept"} {
		if err := g.SisipRiwayat(ctx, tx, models.RiwayatTulis{KasusID: "EDMLF-903", Waktu: "2026-10-01 10:00:00", PIC: "UJI-AKUN-1", Status: s}); err != nil {
			gagal(err)
		}
	}
	v, ada, err := g.VersiBerjalan(ctx, tx, "UJI-PL-P", 0)
	if err != nil || !ada || v.ID != "UJI-NB-P" || v.ProdKe != 1 {
		gagal(err)
	}
	nomor, prodKe, _ := models.NomorEndorsement("UJI-PL-P", v)
	if dobel, err := g.AdaVersiResmi(ctx, tx, "UJI-PL-P", prodKe); err != nil || dobel {
		gagal(err)
	}
	n, err := g.Resmikan(ctx, tx, models.ResmiKasus{ID: "EDMLF-903", NomorPolis: "UJI-PL-P", ProdKe: prodKe, Nomor: nomor, StatusJenis: models.StatusJenis(models.TypeTR)})
	if err != nil || n != 2 {
		gagal(err)
	}
	if c, err := g.TulisRekapWarisan(ctx, tx, repository.RekapWarisanTulis{KasusID: "EDMLF-903", NomorPolis: "UJI-PL-P", Nomor: nomor, COB: "UJI-COB"}); err != nil || c != 1 {
		gagal(err)
	}
	if err := g.Tolak(ctx, tx, "EDMLF-904"); err != nil {
		gagal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	var noPolis, noEndors, statuss string
	var prod int
	if err := repo.QueryRowContext(ctx, `SELECT NO_POLIS, NO_ENDORS, STATUSS, PROD_KE FROM `+skema+`.T_PREMIUM_LIST WHERE ID = 'EDMLF-903'`).
		Scan(&noPolis, &noEndors, &statuss, &prod); err != nil {
		t.Fatal(err)
	}
	if noPolis != "UJI-PL-P" || noEndors != "UJI-PL-P/01" || statuss != models.StatusKasusSelesai || prod != 2 {
		t.Fatalf("kepala %s %s %s %d", noPolis, noEndors, statuss, prod)
	}
	var lama, jenis, edm string
	if err := repo.QueryRowContext(ctx, `SELECT STATUS_OLD, STATUS, PL_NUMBER_EDM FROM `+skema+`.T_PREMIUM_LIST_DETAIL WHERE ID = 'UJI-P1'`).
		Scan(&lama, &jenis, &edm); err != nil || lama != "1" || jenis != "1" || edm != "UJI-PL-P/01" {
		t.Fatalf("peserta Old %s %s %s %v", lama, jenis, edm, err)
	}
	if err := repo.QueryRowContext(ctx, `SELECT STATUS_OLD FROM `+skema+`.T_PREMIUM_LIST_DETAIL WHERE ID = 'UJI-P2'`).Scan(&lama); err != nil || lama != "0" {
		t.Fatalf("peserta New STATUS_OLD %s %v", lama, err)
	}
	if dobel, _ := g.AdaVersiResmi(ctx, nil, "UJI-PL-P", 2); !dobel {
		t.Error("anti-dobel tidak melihat versi yang baru resmi")
	}
	if v, _, _ := g.VersiBerjalan(ctx, nil, "UJI-PL-P", 0); v.ID != "EDMLF-903" || v.ProdKe != 2 {
		t.Errorf("versi berjalan sesudah Confirm %+v", v)
	}
	var premi, cob, idpega string
	if err := repo.QueryRowContext(ctx, `SELECT TO_CHAR(PREMIUM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), COB, IDPEGA FROM `+skema+`.M_LIFE_PREMIUM_SUMMARY
	    WHERE PL_NUMBER_EDM = 'UJI-PL-P/02'`).Scan(&premi, &cob, &idpega); err != nil || premi != "12.5" || cob != "UJI-COB" || idpega != "EDMLF-903" {
		t.Fatalf("rekap warisan %s %s %s %v", premi, cob, idpega, err)
	}
	r, err := g.Riwayat(ctx, nil, "EDMLF-903")
	if err != nil || len(r) != 2 || r[0].No != 2 || r[0].Status != "Accept" || r[1].Status != "Decline" || r[0].Tanggal != "2026-10-01 10:00:00" {
		t.Fatalf("riwayat %+v %v", r, err)
	}
	if err := repo.QueryRowContext(ctx, `SELECT STATUSS FROM `+skema+`.T_PREMIUM_LIST WHERE ID = 'EDMLF-904'`).Scan(&statuss); err != nil || statuss != models.StatusKasusDitolak {
		t.Fatalf("tolak %s %v", statuss, err)
	}
	tx2, _ := repo.Mulai(ctx)
	defer func() { _ = tx2.Rollback() }()
	if err := g.Tolak(ctx, tx2, "EDMLF-904"); err == nil {
		t.Error("kasus tertutup dapat ditolak lagi")
	}
}
