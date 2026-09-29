//go:build db

package services_test

// GILIRAN-17 paket 2 terhadap skema uji Oracle - OQ-M1, OQ-M5, OQ-M6.
//
// ⛔ Seluruh test di sini MELEWATI dengan pesan bila ORACLE_DSN belum
// dikonfigurasi. Melewati bukan lulus.

import (
	"context"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

// pohonUjiG17 - satu klaim, satu peserta, satu baris berkode `kode`
// ("" = baru lahir, belum Save to RNM).
func pohonUjiG17(t *testing.T, svc *services.Service, db *repository.DB, id, kode string) models.PohonKlaim {
	t.Helper()
	pohon := models.PohonKlaim{
		Work: models.WorkClaim{ID: id, Lini: models.LiniLife, Type: "QP"},
		Klaim: models.Klaim{
			NomorKlaim: "UJI-" + id,
			Peserta: []models.Peserta{{
				NomorSertifikat: "017", MataUang: "IDR",
				Baris: []models.BarisAdjustment{{KodeStatus: kode}},
			}},
		},
	}
	ctx := context.Background()
	if err := svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		return repository.NewPohonKlaim(db).Simpan(ctx, tx, pohon)
	}); err != nil {
		t.Fatalf("menyiapkan pohon: %v", err)
	}
	return pohon
}

// OQ-M1: penanda "sudah Save to RNM" diturunkan dari status baris.
func TestSudahSaveRNMDariStatusBaris(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)

	baru := pohonUjiG17(t, svc, db, "CLM-UJI1701", "")
	if sudah, err := baca.SudahSaveRNM(ctx, baru.Work.ID); err != nil || sudah {
		t.Errorf("baris baru lahir: sudah = %v, %v; mau belum", sudah, err)
	}
	disimpan := pohonUjiG17(t, svc, db, "CLM-UJI1702", models.KodeOutstanding)
	if sudah, err := baca.SudahSaveRNM(ctx, disimpan.Work.ID); err != nil || !sudah {
		t.Errorf("baris Outstanding: sudah = %v, %v; mau sudah", sudah, err)
	}
}

// OQ-M5: Remarks sampai ke jejak transisi penolakan, dan kolom KOMENTAR
// (migrasi 021) menerima tulisan delapan kolom perekam Oracle.
func TestTolakMenyimpanRemarksDiJejak(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)
	pelaku := services.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranRejectOutstanding}}
	saat := time.Date(2026, 9, 29, 10, 0, 0, 0, time.UTC)

	pohon := pohonUjiG17(t, svc, db, "CLM-UJI1703", models.KodeOutstanding)
	peserta, _ := baca.AmbilPeserta(ctx, pohon.Work.ID)
	perBaris, _ := baca.AmbilBaris(ctx, pohon.Work.ID)
	jejak := &jejakUji{}
	if err := svc.Status().DenganJejak(jejak).Tolak(ctx, pelaku, pohon.Work.ID,
		perBaris[peserta[0].ID][0].ID, "UJI alasan tolak", saat); err != nil {
		t.Fatalf("Tolak: %v", err)
	}
	if len(jejak.catatan) != 1 || jejak.catatan[0].Komentar != "UJI alasan tolak" {
		t.Errorf("jejak: %+v, mau satu catatan berkomentar", jejak.catatan)
	}

	kedua := pohonUjiG17(t, svc, db, "CLM-UJI1704", models.KodeOutstanding)
	peserta, _ = baca.AmbilPeserta(ctx, kedua.Work.ID)
	perBaris, _ = baca.AmbilBaris(ctx, kedua.Work.ID)
	if err := svc.Status().DenganJejak(services.PerekamJejakOracle(svc)).Tolak(ctx, pelaku,
		kedua.Work.ID, perBaris[peserta[0].ID][0].ID, "UJI alasan Oracle", saat); err != nil {
		t.Fatalf("Tolak lewat perekam Oracle (KOMENTAR, migrasi 021): %v", err)
	}
}

// OQ-M6: cabut = penanda; pembaca menyembunyikannya; tidak dapat diulang.
func TestCabutPesertaMenandaiDanMenyembunyikan(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()
	db, tutupDB := repoUji(t)
	defer tutupDB()
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)

	pohon := pohonUjiG17(t, svc, db, "CLM-UJI1705", "")
	peserta, err := baca.AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("peserta: %v (%d)", err, len(peserta))
	}
	if err := svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		return baca.CabutPeserta(ctx, tx, pohon.Work.ID, peserta[0].ID)
	}); err != nil {
		t.Fatalf("CabutPeserta: %v", err)
	}
	if sisa, err := baca.AmbilPeserta(ctx, pohon.Work.ID); err != nil || len(sisa) != 0 {
		t.Errorf("sesudah cabut: %d peserta, %v; mau 0", len(sisa), err)
	}
	if baris, err := baca.AmbilBaris(ctx, pohon.Work.ID); err != nil || len(baris) != 0 {
		t.Errorf("sesudah cabut: baris peserta tercabut masih terbaca (%d), %v", len(baris), err)
	}
	if err := svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		return baca.CabutPeserta(ctx, tx, pohon.Work.ID, peserta[0].ID)
	}); err == nil {
		t.Error("mencabut dua kali tidak gagal")
	}
}
