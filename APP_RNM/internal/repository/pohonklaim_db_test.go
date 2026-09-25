//go:build db

// Seam `repository` terhadap skema uji Oracle NYATA - tiket 14.
//
// Untuk apa berkas ini: membuktikan bahwa migrasi benar-benar membangun tujuh
// tabel dengan relasinya, bahwa satu pohon klaim dapat ditulis dan dibaca
// kembali sampai tingkat terdalam, dan bahwa menghapus klaim mengkaskade sampai
// cicit.
//
// Jalankan: make test-db   (perlu ORACLE_DSN dan ORACLE_SCHEMA)
//
// Tanpa instance Oracle, seluruh test di sini MELEWATI dengan pesan - bukan
// lulus diam-diam.
package repository_test

import (
	"context"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/pkg/utils"
)

func siapkanPohon(t *testing.T) (*repository.DB, *repository.PohonKlaim, func()) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
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
	uang := func(s string) models.Money {
		m, err := models.NewMoney(s, "IDR")
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	return models.PohonKlaim{
		Work: models.WorkClaim{
			ID: "CLM-UJI900", Lini: models.LiniLife, Type: "UJI-TYPE", CaseID: "UJI-CASE-900",
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
						Retro: []models.SpreadingRetro{
							{ID: "UJI-RT-1", ReinsurerName: "UJI-REINSURER", Amount: uang("250000.25")},
							{ID: "UJI-RT-2", ReinsurerName: "UJI-REINSURER-2", Amount: uang("0.00000001")},
						},
					}},
				}},
			}},
		},
	}
}

// Migrasi idempoten: menjalankannya dua kali tidak menambah apa pun.
func TestMigrasiIdempoten(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()

	lap, err := db.JalankanMigrasi(context.Background())
	if err != nil {
		t.Fatalf("migrasi kedua gagal: %v", err)
	}
	if len(lap.Dijalankan) != 0 {
		t.Errorf("migrasi kedua menjalankan ulang %v - tidak idempoten", lap.Dijalankan)
	}
	if len(lap.Dilewati) == 0 {
		t.Error("tidak satu pun langkah dilaporkan dilewati")
	}
}

// Jalur mundur benar-benar membongkar, dan migrasi dapat dijalankan lagi
// sesudahnya.
func TestJalurMundurDiuji(t *testing.T) {
	db, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	if _, err := db.BongkarMigrasi(ctx); err != nil {
		t.Fatalf("jalur mundur gagal: %v", err)
	}
	lap, err := db.JalankanMigrasi(ctx)
	if err != nil {
		t.Fatalf("migrasi ulang sesudah mundur gagal: %v", err)
	}
	if len(lap.Dijalankan) == 0 {
		t.Error("sesudah dibongkar, migrasi ulang tidak menjalankan apa pun")
	}
}

// Satu pohon utuh ditulis dalam SATU transaksi, ke tabel relasional DAN ke
// baris datar warisan.
func TestSimpanPohonMenulisDuaTempatDalamSatuTransaksi(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)

	tx, err := db.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatalf("menyimpan pohon: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	// Penulisan kedua benar-benar terjadi: satu baris datar per baris adjustment.
	n, err := repo.CacahBarisLama(ctx, p.Work.CaseID)
	if err != nil {
		t.Fatal(err)
	}
	if n != p.Klaim.CacahBaris() {
		t.Errorf("baris datar %d, mau %d", n, p.Klaim.CacahBaris())
	}
}

// Pembacaan turun sampai CICIT - tingkat keenam pohon.
func TestBacaSampaiCicit(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)

	tx, _ := db.Mulai(ctx)
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	_ = tx.Commit()

	spr, err := repo.AmbilSpreading(ctx, p.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	daftar := spr["UJI-A-1"]
	if len(daftar) != 1 {
		t.Fatalf("spreading %d, mau 1", len(daftar))
	}
	if got := utils.FormatDecimal(daftar[0].IDR.Amount); got != "500000.5" {
		t.Errorf("IDR = %q, mau 500000.5", got)
	}
	if len(daftar[0].Retro) != 2 {
		t.Fatalf("spreading retro %d, mau 2 - cicit tidak terbaca", len(daftar[0].Retro))
	}
	// Uang di tingkat terdalam tetap utuh, digit demi digit.
	if got := utils.FormatDecimal(daftar[0].Retro[1].Amount.Amount); got != "0.00000001" {
		t.Errorf("jumlah retro = %q, mau 0.00000001", got)
	}
}

// Menghapus klaim mengkaskade sampai tingkat terdalam. Test memeriksa CICIT
// ikut hilang, bukan hanya anak langsungnya.
func TestHapusMengkaskadeSampaiCicit(t *testing.T) {
	db, repo, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()
	p := contohPohon(t)

	tx, _ := db.Mulai(ctx)
	if err := repo.Simpan(ctx, tx, p); err != nil {
		_ = tx.Rollback()
		t.Fatal(err)
	}
	_ = tx.Commit()

	tx2, _ := db.Mulai(ctx)
	if err := repo.Hapus(ctx, tx2, p.Work.ID, p.Work.CaseID); err != nil {
		_ = tx2.Rollback()
		t.Fatalf("menghapus pohon: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatal(err)
	}

	spr, err := repo.AmbilSpreading(ctx, p.Work.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(spr) != 0 {
		t.Errorf("sesudah hapus masih ada %d kelompok spreading - kaskade tidak sampai cicit", len(spr))
	}
	n, err := repo.CacahBarisLama(ctx, p.Work.CaseID)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("baris datar warisan tersisa %d", n)
	}
}

// Pembongkaran data lama berjalan di atas tabel tiruan warisan, dan jumlah
// baris sesudah migrasi sama dengan sebelumnya.
func TestBongkarDataLamaDariTabelTiruan(t *testing.T) {
	_, _, bersihkan := siapkanPohon(t)
	defer bersihkan()
	ctx := context.Background()

	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		t.Skipf("lewati: %v", err)
	}
	defer func() { _ = sqlDB.Close() }()

	masuk := []repository.BarisLama{
		{ID: "L1", CASEID: "UJI-CASE-800", NO_CLAIM: "UJI-CLM-8", CERTIFICATE_NO: "006",
			CURRENCY: "IDR", CLAIM_AMOUNT: "100.00000001", STS_REJECT: "0"},
		{ID: "L2", CASEID: "UJI-CASE-800", NO_CLAIM: "UJI-CLM-8", CERTIFICATE_NO: "006",
			CURRENCY: "IDR", CLAIM_AMOUNT: "200", STS_REJECT: "1"},
	}
	if err := skemauji.IsiBarisLama(ctx, sqlDB, skema, masuk); err != nil {
		t.Fatalf("mengisi tabel tiruan: %v", err)
	}

	pohon, lap := repository.BongkarBarisLama(masuk)
	if len(pohon) != 1 {
		t.Fatalf("klaim terbentuk %d, mau 1", len(pohon))
	}
	if lap.AdjustmentTerbentuk != len(masuk) {
		t.Errorf("adjustment %d, mau %d", lap.AdjustmentTerbentuk, len(masuk))
	}
	if lap.BarisHardcode != 1 {
		t.Errorf("baris hardcode %d, mau 1", lap.BarisHardcode)
	}
}
