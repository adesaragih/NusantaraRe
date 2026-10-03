//go:build db

package repository_test

// Uji seam repository terhadap Oracle SUNGGUHAN (skema uji, `make test-db`):
// pulang-pergi halaman lewat katalog, generasi tertutup ditolak, nomor polis
// sekali dan unik, dan pembatalan transaksi yang tidak menyisakan baris.
//
// Tanpa instance Oracle (ORACLE_DSN / skema uji), seluruh test di sini
// MELEWATI dengan pesan. ⛔ Fixture berawalan UJI-.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
	"nusantarare/uji/skemauji"
)

func pasang(t *testing.T) (*sql.DB, string, context.Context, *intidb.DB) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	t.Cleanup(func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) })
	repo, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = repo.Close() })
	return sqlDB, skema, ctx, repo
}

func dalamTx(t *testing.T, ctx context.Context, d *intidb.DB, f func(tx *intidb.Tx) error) error {
	t.Helper()
	tx, err := d.Mulai(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := f(tx); err != nil {
		return err
	}
	return tx.Commit()
}

func TestPulangPergiHalamanLewatKatalog(t *testing.T) {
	_, _, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-1"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.Setel("PolicyTreatyIn.PremiOgp", "830.82191780804")
	h.Setel("PolicyTreatyIn.RiCommOgp", "12.5")
	h.Setel("PolicyTreatyIn.StartDate", "2026-10-01")
	h.Setel("PolicyTreatyIn.StatementDate", "2026-10-03 09:15:00")
	h.Setel("PolicyTreatyIn.BizCode", "006")
	h.Setel("Quotation.BusinessOldId", "01")
	h.Setel("Quotation.BusinessFac", "T")
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-10015", "SharePercentage": "33.3333"}, {"TreatyType": "UJI-10218"}})
	h.SetelDaftar(models.DaftarAngsuran, []models.Baris{{"InstallmentNo": "1", "DueDate": "2026-11-01", "Premium": "1.5"}})
	h.SetelDaftar(models.JalurAnak(models.DaftarAngsuran, 1, "InstallmentList"), []models.Baris{{"InstallmentNo": "1", "PremiumAfterTax": "0.75"}})
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		return g.SimpanHalaman(ctx, tx, id, h)
	}); err != nil {
		t.Fatal(err)
	}
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	for j, harap := range map[string]string{
		"PolicyTreatyIn.PremiOgp":                  "830.82191780804", // presisi penuh (AC 23)
		"PolicyTreatyIn.RiCommOgp":                 "12.5",
		"PolicyTreatyIn.StartDate":                 "2026-10-01",
		"PolicyTreatyIn.StatementDate":             "2026-10-03 09:15:00",
		"PolicyTreatyIn.BizCode":                   "006", // nol di depan bertahan
		"Quotation.BusinessOldId":                  "01",
		"PolicyTreatyIn.QuotationData.BusinessFac": "T",
	} {
		if got := b.Ambil(j); got != harap {
			t.Errorf("%s = %q, harap %q", j, got, harap)
		}
	}
	sp := b.AmbilDaftar(models.DaftarSpreading)
	if len(sp) != 2 || sp[0]["SharePercentage"] != "33.3333" || sp[1]["TreatyType"] != "UJI-10218" {
		t.Errorf("spreading %+v", sp)
	}
	rinci := b.AmbilDaftar(models.JalurAnak(models.DaftarAngsuran, 1, "InstallmentList"))
	if len(rinci) != 1 || rinci[0]["PremiumAfterTax"] != "0.75" {
		t.Errorf("rincian angsuran %+v", rinci)
	}
	// NB: baris dihapus, NOURUT dinomori ulang (ID-12)
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-10218"}})
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
		t.Fatal(err)
	}
	b, _ = g.BacaHalaman(ctx, nil, id)
	if sp := b.AmbilDaftar(models.DaftarSpreading); len(sp) != 1 || sp[0]["TreatyType"] != "UJI-10218" {
		t.Errorf("sesudah hapus baris: %+v", sp)
	}
}

func TestNomorPolisSekaliDanUnik(t *testing.T) { // AC 31, 74
	_, _, ctx, d := pasang(t)
	g := repository.Baru(d)
	for i := 1; i <= 2; i++ {
		id := fmt.Sprintf("UJI-NB-NO-%d", i)
		if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI") }); err != nil {
			t.Fatal(err)
		}
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-1", "UJI-QR.T1.10.2026.00001") }); err != nil {
		t.Fatal(err)
	}
	err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-1", "UJI-QR.T1.10.2026.00002") })
	if !errors.Is(err, repository.ErrNomorPolisSudahAda) {
		t.Fatalf("nomor kedua untuk berkas yang sama: %v", err)
	}
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SetelNomorPolis(ctx, tx, "UJI-NB-NO-2", "UJI-QR.T1.10.2026.00001") }); err == nil {
		t.Fatal("dua berkas bernomor sama (PRODKE sama) harus ditolak indeks unik")
	}
}

func TestGenerasiTertutupDitolakDanPembatalanUtuh(t *testing.T) { // ID-10, AC 29
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-TUTUP"
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI") }); err != nil {
		t.Fatal(err)
	}
	// pembatalan: kasus kedua disisipkan lalu transaksi digagalkan
	err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, "UJI-NB-BATAL", "UJI-AKUN", "UJI"); err != nil {
			return err
		}
		return errors.New("UJI-gagal di tengah")
	})
	if err == nil {
		t.Fatal("harap galat")
	}
	if _, err := g.Keadaan(ctx, nil, "UJI-NB-BATAL"); !errors.Is(err, repository.ErrKasusTidakAda) {
		t.Fatalf("baris tersisa sesudah pembatalan: %v", err)
	}
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`UPDATE %s.T_GENERAL_POLIS SET TGL_TUTUP = SYSDATE WHERE ID = :1`, skema), id); err != nil {
		t.Fatal(err)
	}
	err = dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, models.HalamanBaru()) })
	if !errors.Is(err, repository.ErrGenerasiTertutup) {
		t.Fatalf("generasi tertutup: %v", err)
	}
}
