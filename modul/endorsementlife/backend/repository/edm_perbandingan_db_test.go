//go:build db

package repository_test

// Perbandingan summary terhadap Oracle (08-10-2026): After = rekap kasus untuk
// setiap kolom. Fixture UJI-, nol nama orang.

import (
	"context"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

func TestPerbandinganTerhadapOracle(t *testing.T) {
	repo, skema, tutup := pasangSkema(t)
	defer tutup()
	ctx := context.Background()
	jalankan(t, repo,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST (ID, OLD_POLICY_NO, TGL_INPUT, TYPE, EDM_TYPE, PROD_KE)
		   VALUES ('EDML-902', 'UJI-PL-B', SYSDATE, 'QR', '1', 1)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM, DEDUCTION)
		   VALUES ('UJI-B1', 'EDML-902', 'UJI-PL-B', 'Old', 'IDR', 100.25, 10)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM)
		   VALUES ('UJI-B2', 'EDML-902', 'UJI-PL-B', 'Old', 'IDR', 200)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM)
		   VALUES ('UJI-B3', 'EDML-902', 'UJI-PL-B', 'New', 'IDR', 50)`,
		`INSERT INTO `+skema+`.T_PREMIUM_LIST_DETAIL (ID, PREMIUM_LIST_ID, PL_NUMBER, EDM_STATUS, CURRENCY, GROSS_PREMIUM)
		   VALUES ('UJI-B4', 'EDML-902', 'UJI-PL-B', 'Old', 'USD', 7.5)`,
	)
	g := repository.Baru(repo)
	tx, err := repo.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Tandai(ctx, tx, "EDML-902", models.StatusDelete, models.PilihanHapus{Pilih: []string{"UJI-B1"}}); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if _, err := g.HitungRekap(ctx, tx, "EDML-902", models.TypeQR); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	kel, err := g.PerbandinganKasus(ctx, nil, "EDML-902", models.TypeQR)
	if err != nil {
		t.Fatal(err)
	}
	b, err := models.SusunPerbandingan(kel, models.KolomRekapKasus[1:])
	if err != nil || len(b) != 2 {
		t.Fatalf("%+v %v", b, err)
	}
	if b[0].Before.Nilai["PREMIUM"] != "300.25" || b[0].Deleted.Nilai["PREMIUM"] != "-100.25" || b[0].Deleted.Cacah != -1 ||
		b[0].Added.Nilai["PREMIUM"] != "50" || b[0].After.Nilai["PREMIUM"] != "250" {
		t.Errorf("IDR %+v", b[0])
	}
	// ⚠️ Rekap menjumlah seluruh baris: 200 - 100.25 + 50 = 149.75 = After (250) + Deleted (-100.25).
	r, err := g.RekapKasus(ctx, nil, "EDML-902")
	if err != nil || len(r) != 2 || r[0]["PREMIUM"] != "149.75" {
		t.Fatalf("rekap %v %v", r, err)
	}
}
