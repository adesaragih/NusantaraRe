//go:build db

package lintasmodul_test

// Pembantu uji lintas modul ber-tag `db` - salinan APA ADANYA.
//
// Refactor bentuk B (30-09-2026): `siapkanPohon` dan `contohPohon` milik uji
// repository Claim Life (`internal/repository/pohonklaim_db_test.go`). Uji
// integrasi Claim Life + Komite di paket ini memerlukan keduanya; paket uji
// luar tidak dapat dipinjam lintas folder, maka disalin.

import (
	"context"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/inti"
	"nusantarare/inti/db"
	intiuang "nusantarare/inti/uang"
)

func siapkanPohon(t *testing.T) (*db.DB, *repository.PohonKlaim, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		// Hanya "Oracle belum dikonfigurasi" yang dilewati. Salah
		// konfigurasi, menunjuk produksi, atau menunjuk skema yang bukan
		// skema uji harus MENGGAGALKAN - test ini menghapus tabel.
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	return db, repository.NewPohonKlaim(db), func() {
		_ = skemauji.Bongkar(ctx, sqlDB, skema)
		_ = db.Close()
		_ = sqlDB.Close()
	}
}

// contohPohon membuat satu klaim lengkap sampai tingkat terdalam.
// ⛔ Nol nama orang, nol nomor polis nyata.
func contohPohon(t *testing.T) models.PohonKlaim {
	t.Helper()
	uang := func(s string) intiuang.Money {
		m, err := intiuang.NewMoney(s, "IDR")
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	rasio := func(s string) intiuang.Ratio {
		r, err := intiuang.NewRatio(s, 8)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	return models.PohonKlaim{
		Work: models.WorkClaim{
			ID: "CLM-UJI900", Lini: inti.LiniLife, Type: "UJI-TYPE", CaseID: "UJI-CASE-900",
		},
		Klaim: models.Klaim{
			ID: "CLM-UJI900", NomorKlaim: "UJI-CLM-9", NomorPolis: "UJI-POL-9",
			NamaBisnis: "UJI BISNIS", KodeStatus: "0",
			Peserta: []models.Peserta{{
				ID: "UJI-P-1", NomorSertifikat: "006", MataUang: "IDR",
				Baris: []models.BarisAdjustment{{
					ID: "UJI-A-1", KodeStatus: "1", JumlahKlaim: uang("1234567890.12345678"),
					Spreading: []models.Spreading{{
						ID: "UJI-S-1", TreatyTypeName: "UJI-TREATY", TreatyYearLife: "2026",
						IDR: uang("500000.5"), Currency: "IDR",
						// RetrocadedShare UANG (diralat 26-09-2026 menurut
						// SpreadingClaimLife_Act langkah 8.2.1.4-7); Rate rasio.
						// Keduanya sengaja bertipe berbeda.
						RetrocadedShare: uang("0.12345678"),
						Rate:            rasio("0.075"),
						Retro: []models.SpreadingRetro{
							{ID: "UJI-RT-1", ReinsurerName: "UJI-REINSURER",
								Amount: uang("250000.25"), PercentShare: rasio("0.6"),
								Rate: rasio("0.0125"), PremiumSpreadedGross: uang("99999.99999999"),
								PremiumSpreadedNet: uang("88888.88888888"),
								// COMMISION dan OVR_COMM adalah PERSEN, bukan uang.
								Commision: rasio("1234.5"), OvrComm: rasio("0.00000001")},
							{ID: "UJI-RT-2", ReinsurerName: "UJI-REINSURER-2",
								Amount: uang("0.00000001"), PercentShare: rasio("0.4")},
						},
					}},
				}},
			}},
		},
	}
}
