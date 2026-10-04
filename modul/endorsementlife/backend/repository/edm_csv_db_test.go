//go:build db

package repository_test

// `Add CSV Data` terhadap Oracle: penyisipan bertipe (uang teks, tanggal),
// pembuangan baris `New`, dan pembatalan transaksi. Fixture UJI-.

import (
	"context"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

func TestCSVTerhadapOracle(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	jalankan(t, repo,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, OLD_POLICY_NO, TGL_INPUT, TYPE, EDM_TYPE, PROD_KE)
		   VALUES ('EDMLF-902', 'UJI-PL-C', SYSDATE, 'QR', '1', 2)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CERTIFICATE_NO, PLAN, POLICY_HOLDER)
		   VALUES ('UJI-C2', 'EDMLF-902', 'UJI-PL-C', 'Old', 'UJI-B', 'UJI-PLAN-B', 'UJI-PH')`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CERTIFICATE_NO, PLAN, POLICY_HOLDER)
		   VALUES ('UJI-C1', 'EDMLF-902', 'UJI-PL-C', 'Old', 'UJI-A', 'UJI-PLAN', 'UJI-PH')`,
	)
	g := repository.Baru(repo)
	a, ada, err := g.AcuanCSV(ctx, nil, "EDMLF-902")
	if err != nil || !ada || a.Plan != "UJI-PLAN" {
		t.Fatalf("acuan %+v %v %v (peserta pertama urutan grid)", a, ada, err)
	}
	nilai, pesan := models.RapikanBarisCSV(1, map[string]string{"PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH",
		"GROSS_PREMIUM": "-0,00000001", "DOB": "02/01/1990", "STNC": "03/04/2026", "AGE": "36"}, a)
	if len(pesan) != 0 {
		t.Fatal(pesan)
	}
	tx, err := repo.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n, err := g.SisipPesertaCSV(ctx, tx, "EDMLF-902", "UJI-PL-C", []models.BarisCSV{{Nomor: 1, Nilai: nilai}}); err != nil || n != 1 {
		_ = tx.Rollback()
		t.Fatalf("sisip %d %v", n, err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	var gp, dob, stnc, status, pl string
	var umur int
	if err := repo.QueryRowContext(ctx, `SELECT TO_CHAR(GROSS_PREMIUM, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,'''), TO_CHAR(DOB, 'YYYY-MM-DD'),
	    STNC, EDM_STATUS, PL_NUMBER, AGE FROM `+skema+`.T_PREMIUM_LIST_DETAIL WHERE PREMIUM_LIST_ID = 'EDMLF-902' AND EDM_STATUS = 'New'`).
		Scan(&gp, &dob, &stnc, &status, &pl, &umur); err != nil {
		t.Fatal(err)
	}
	if gp != "-.00000001" && gp != "-0.00000001" || dob != "1990-01-02" || stnc != "03/04/2026" || pl != "UJI-PL-C" || umur != 36 {
		t.Fatalf("baris CSV = %s %s %s %s %s %d", gp, dob, stnc, status, pl, umur)
	}
	// Pembatalan: baris yang sudah disisip ikut hilang bersama transaksinya.
	tx, _ = repo.Mulai(ctx)
	if n, err := g.HapusPesertaBaru(ctx, tx, "EDMLF-902"); err != nil || n != 1 {
		t.Fatalf("hapus New %d %v", n, err)
	}
	_, _ = g.SisipPesertaCSV(ctx, tx, "EDMLF-902", "UJI-PL-C", []models.BarisCSV{{Nomor: 2, Nilai: nilai}})
	_ = tx.Rollback()
	var cacah int
	_ = repo.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.T_PREMIUM_LIST_DETAIL WHERE PREMIUM_LIST_ID = 'EDMLF-902'`).Scan(&cacah)
	if cacah != 3 {
		t.Fatalf("sesudah rollback %d baris, mau 3 (2 Old + 1 New semula)", cacah)
	}
}
