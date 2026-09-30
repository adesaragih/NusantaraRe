//go:build db

package repository_test

// OQ-N2 (GILIRAN-17) terhadap skema uji Oracle: cermin OS_AKSEPTASI_KLAIM_LIFE
// mengisi NAME_OF_INSURED dan DOB DI DALAM SQL dari baris sumber peserta.
//
// ⛔ Nama dan tanggal lahir TIDAK dibaca ke Go, bahkan di uji: yang diperiksa
// cacah baris yang cocok antara cermin dan sumbernya, di dalam Oracle.
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan.

import (
	"context"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/repository"
	"nusantarare/uji/skemauji"
)

func TestCerminMengisiTertanggungDariSumber(t *testing.T) {
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	defer func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) }()
	if err := skemauji.IsiPesertaPolis(ctx, sqlDB, skema); err != nil {
		t.Fatalf("mengisi peserta polis tiruan: %v", err)
	}
	// Polis tiruan: CEDINGCO cermin datang dari `T_PREMIUM_LIST.CEDING_CO`
	// baris PROD_KE terakhir (b9226) - dua versi, yang terbaru menang.
	for _, q := range []string{
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID) VALUES ('UJI-POLIS-1710')`,
		`INSERT INTO ` + skema + `.T_WORK_POLIS (ID) VALUES ('UJI-POLIS-1711')`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, NO_POLIS, CEDING_CO, PROD_KE) VALUES ('UJI-POLIS-1710', 'UJI-POL-0001', 'UJI-CEDING-LAMA', 0)`,
		`INSERT INTO ` + skema + `.T_PREMIUM_LIST (ID, NO_POLIS, CEDING_CO, PROD_KE) VALUES ('UJI-POLIS-1711', 'UJI-POL-0001', 'UJI-CEDING-1', 1)`,
	} {
		if _, err := sqlDB.ExecContext(ctx, q); err != nil {
			t.Fatalf("polis tiruan: %v", err)
		}
	}
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	simpan := func(id, sumber string) (models.PohonKlaim, error) {
		p := models.PohonKlaim{
			Work: models.WorkClaim{ID: id, CaseID: id, Lini: inti.LiniLife, Type: "QP"},
			Klaim: models.Klaim{
				NomorKlaim: "UJI-" + id, NomorPolis: "UJI-POL-0001",
				Peserta: []models.Peserta{{
					NomorPremiList: "UJI-PL-1", NomorSertifikat: "006", SumberID: sumber, MataUang: "IDR",
					Baris: []models.BarisAdjustment{{}},
				}},
			},
		}
		tx, err := db.Mulai(ctx)
		if err != nil {
			return p, err
		}
		if err := repository.NewPohonKlaim(db).Simpan(ctx, tx, p); err != nil {
			_ = tx.Rollback()
			return p, err
		}
		return p, tx.Commit()
	}

	if _, err := simpan("CLM-UJI1710", "UJI-SRC-1"); err != nil {
		t.Fatalf("Simpan dengan baris sumber: %v", err)
	}
	var cocok int
	if err := sqlDB.QueryRowContext(ctx, `SELECT COUNT(*)
		   FROM `+skema+`.OS_AKSEPTASI_KLAIM_LIFE o, `+skema+`.M_LIFE_PREMIUM_DETAIL m
		  WHERE o.CASEID = :1 AND m.ID = :2
		    AND o.NAME_OF_INSURED = m.NAME_OF_INSURED AND o.DOB = TRUNC(m.DOB)`,
		"CLM-UJI1710", "UJI-SRC-1").Scan(&cocok); err != nil {
		t.Fatal(err)
	}
	if cocok != 1 {
		t.Errorf("baris cermin bertertanggung sama dengan sumbernya: %d, mau 1", cocok)
	}
	var ceding int
	if err := sqlDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.OS_AKSEPTASI_KLAIM_LIFE
		  WHERE CASEID = :1 AND CEDINGCO = 'UJI-CEDING-1'`, "CLM-UJI1710").Scan(&ceding); err != nil {
		t.Fatal(err)
	}
	if ceding != 1 {
		t.Errorf("CEDINGCO cermin bukan milik versi polis terbaru: %d, mau 1", ceding)
	}

	// Baris sumber yang tidak ada: GALAT, bukan cermin tanpa nama diam-diam.
	if _, err := simpan("CLM-UJI1711", "UJI-SRC-TIADA"); err == nil {
		t.Error("Simpan dengan baris sumber yang tidak ada tidak gagal")
	}
}
