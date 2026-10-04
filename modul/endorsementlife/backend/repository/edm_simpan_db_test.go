//go:build db

package repository_test

// `Save` terhadap Oracle: jurnal balik idempoten (R06) dan rekap mata uang
// dihitung ulang (R08). Fixture UJI-, nol nama orang.

import (
	"context"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

func TestSimpanTerhadapOracle(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	jalankan(t, repo,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, OLD_POLICY_NO, TGL_INPUT, TYPE, EDM_TYPE, PROD_KE)
		   VALUES ('EDMLF-901', 'UJI-PL-S', SYSDATE, 'QR', '1', 2)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM, DEDUCTION, COMM)
		   VALUES ('UJI-S1', 'EDMLF-901', 'UJI-PL-S', 'Old', 'IDR', 100.25, 10, 1)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM)
		   VALUES ('UJI-S2', 'EDMLF-901', 'UJI-PL-S', 'Old', 'IDR', 200)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM, DEDUCTION)
		   VALUES ('UJI-S3', 'EDMLF-901', 'UJI-PL-S', 'Old', 'USD', 7.5, 0.5)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM)
		   VALUES ('UJI-S4', 'EDMLF-901', 'UJI-PL-S', 'New', 'IDR', 50)`,
	)
	g := repository.Baru(repo)
	dalamTx := func(f func(tx *intidb.Tx) error) {
		t.Helper()
		tx, err := repo.Mulai(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := f(tx); err != nil {
			_ = tx.Rollback()
			t.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
	}
	pilih := models.PilihanHapus{Pilih: []string{"UJI-S1", "UJI-S4"}}
	dalamTx(func(tx *intidb.Tx) error {
		n, err := g.Tandai(ctx, tx, "EDMLF-901", models.StatusDelete, pilih)
		if err != nil || n != 1 {
			t.Fatalf("tandai pertama %d %v (baris New tidak tersentuh)", n, err)
		}
		// Pembalikan kedua mustahil: baris itu tidak lagi Old.
		if n, err = g.Tandai(ctx, tx, "EDMLF-901", models.StatusDelete, pilih); err != nil || n != 0 {
			t.Fatalf("tandai kedua %d %v", n, err)
		}
		if n, err = g.HitungRekap(ctx, tx, "EDMLF-901", models.TypeQR); err != nil || n != 2 {
			t.Fatalf("rekap %d %v", n, err)
		}
		// Hitung ulang: baris lama terhapus, tidak berganda.
		n, err = g.HitungRekap(ctx, tx, "EDMLF-901", models.TypeQR)
		if err == nil && n != 2 {
			t.Fatalf("rekap ulang %d", n)
		}
		return err
	})
	var gp, ded, comm, status string
	if err := repo.QueryRowContext(ctx, `SELECT TO_CHAR(GROSS_PREMIUM, 'TM9'), TO_CHAR(DEDUCTION, 'TM9'), TO_CHAR(COMM, 'TM9'), EDM_STATUS
	    FROM `+skema+`.T_PREMIUM_LIST_DETAIL WHERE ID = 'UJI-S1'`).Scan(&gp, &ded, &comm, &status); err != nil {
		t.Fatal(err)
	}
	if gp != "-100.25" || ded != "-10" || comm != "1" || status != models.StatusDelete {
		t.Fatalf("UJI-S1 = %s %s %s %s (COMM di luar 32 kolom)", gp, ded, comm, status)
	}
	r, err := g.RekapKasus(ctx, nil, "EDMLF-901")
	if err != nil || len(r) != 2 {
		t.Fatalf("rekap %v %v", r, err)
	}
	if r[0]["CURRENCY"] != "IDR" || r[0]["PREMIUM"] != "149.75" || r[0]["BALANCE"] != "154.75" || r[0]["COMMISSION"] != "1" ||
		r[1]["CURRENCY"] != "USD" || r[1]["BALANCE"] != "7" {
		t.Fatalf("rekap = %v", r)
	}
	ada, err := g.AdaRekap(ctx, nil, "EDMLF-901")
	if err != nil || !ada {
		t.Fatalf("AdaRekap %v %v", ada, err)
	}
	dalamTx(func(tx *intidb.Tx) error {
		n, err := g.Tandai(ctx, tx, "EDMLF-901", models.StatusBatal, models.PilihanHapus{Semua: true, Kecuali: []string{"UJI-S3"}})
		if err == nil && n != 1 {
			t.Fatalf("DELETE ALL berpengecualian menandai %d, mau 1 (UJI-S2)", n)
		}
		return err
	})
}
