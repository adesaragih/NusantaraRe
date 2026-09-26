//go:build db

package services_test

// Transisi status terhadap skema uji Oracle - tiket 04.
//
// ⛔ Seluruh test di sini MELEWATI dengan pesan bila ORACLE_DSN belum
// dikonfigurasi. Melewati bukan lulus.
//
// Dibaca sesudah: services/statusbaris.go.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/repository/skemauji"
	"nusantarare/internal/services"
)

// jejakUji menggantikan tempat jejak audit yang belum disahkan (butir am).
type jejakUji struct{ catatan []services.CatatanJejak }

func (j *jejakUji) Rekam(_ context.Context, _ *repository.Tx,
	c services.CatatanJejak) error {
	j.catatan = append(j.catatan, c)
	return nil
}

// jejakGagal meniru kegagalan DI TENGAH transaksi, sesudah dua tulisan
// pertama berhasil. ADR-U-0029 Akibat 2: test memeriksa bahwa kegagalan di
// tengah tidak meninggalkan baris separuh jadi.
type jejakGagal struct{}

var errJejakSengaja = errors.New("uji: jejak sengaja gagal")

func (jejakGagal) Rekam(context.Context, *repository.Tx, services.CatatanJejak) error {
	return errJejakSengaja
}

// repoUji membuka penanganan repository sendiri untuk test ini.
//
// Layanan tidak membuka isi dalamnya, dan memang tidak seharusnya: test yang
// memerlukan repository membukanya sendiri dan menutupnya sendiri.
func repoUji(t *testing.T) (*repository.DB, func()) {
	t.Helper()
	db, err := skemauji.BukaRepositori()
	if err != nil {
		t.Fatalf("membuka repositori uji: %v", err)
	}
	return db, func() { _ = db.Close() }
}

// pohonUjiStatus membuat satu klaim dengan satu peserta dan satu baris
// Outstanding, langsung lewat repository.
func pohonUjiStatus(t *testing.T, svc *services.Service, db *repository.DB) models.PohonKlaim {
	t.Helper()
	pohon := models.PohonKlaim{
		Work: models.WorkClaim{ID: "CLM-UJI400", Lini: models.LiniLife, Type: "QP"},
		Klaim: models.Klaim{
			NomorKlaim: "UJI-CLM-400",
			Peserta: []models.Peserta{{
				NomorSertifikat: "006", MataUang: "IDR",
				Baris: []models.BarisAdjustment{{KodeStatus: models.KodeOutstanding}},
			}},
		},
	}
	ctx := context.Background()
	err := svc.DalamTransaksi(ctx, func(tx *repository.Tx) error {
		return repository.NewPohonKlaim(db).Simpan(ctx, tx, pohon)
	})
	if err != nil {
		t.Fatalf("menyiapkan pohon: %v", err)
	}
	return pohon
}

// TestUbahStatusMencerminkanTigaTingkat - baris, peserta, dan header berubah
// bersama-sama dalam satu transaksi.
func TestUbahStatusMencerminkanTigaTingkat(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()

	db, tutupDB := repoUji(t)
	defer tutupDB()
	pohon := pohonUjiStatus(t, svc, db)
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)
	peserta, err := baca.AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil || len(peserta) != 1 {
		t.Fatalf("membaca peserta: %v (%d peserta)", err, len(peserta))
	}
	perBaris, err := baca.AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca baris: %v", err)
	}
	adj := perBaris[peserta[0].ID]
	if len(adj) != 1 {
		t.Fatalf("baris = %d, mau 1", len(adj))
	}

	jejak := &jejakUji{}
	saat := time.Date(2026, 9, 26, 22, 0, 0, 0, time.UTC)
	err = svc.Status().DenganJejak(jejak).Ubah(ctx,
		services.Pelaku{AkunID: "UJI-AKUN"},
		pohon.Work.ID, peserta[0].ID, adj[0].ID, models.StatusAksep, saat)
	if err != nil {
		t.Fatalf("Ubah: %v", err)
	}

	// Baris berubah, dan tanggal akseptasinya terstempel.
	perBaris, err = baca.AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca ulang baris: %v", err)
	}
	if got := perBaris[peserta[0].ID][0].KodeStatus; got != models.KodeAksep {
		t.Errorf("kode baris = %q, mau %q", got, models.KodeAksep)
	}
	if perBaris[peserta[0].ID][0].TanggalAkseptasi.IsZero() {
		t.Error("ACCEPTATION_DATE tidak terstempel saat baris diaksep")
	}
	// Peserta ikut tercermin.
	peserta, err = baca.AmbilPeserta(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca ulang peserta: %v", err)
	}
	// Header ikut tercermin.
	hdr, err := baca.AmbilHeader(ctx, pohon.Work.ID)
	if err != nil || hdr == nil {
		t.Fatalf("membaca header: %v", err)
	}
	if hdr.KodeStatus != models.KodeAksep {
		t.Errorf("kode header = %q, mau %q", hdr.KodeStatus, models.KodeAksep)
	}
	if len(jejak.catatan) != 1 || jejak.catatan[0].AkunID != "UJI-AKUN" {
		t.Errorf("jejak = %+v, mau satu catatan ber-AkunID UJI-AKUN", jejak.catatan)
	}

	// ⛔ Transisi KEDUA atas baris yang sama ditolak - kefinalan berlaku juga
	// terhadap basis data, bukan hanya di dalam proses ini.
	err = svc.Status().DenganJejak(jejak).Ubah(ctx,
		services.Pelaku{AkunID: "UJI-AKUN"},
		pohon.Work.ID, peserta[0].ID, adj[0].ID, models.StatusDitolak, saat)
	if !errors.Is(err, services.ErrBarisSudahFinal) {
		t.Errorf("transisi kedua: galat = %v, mau ErrBarisSudahFinal", err)
	}
}

// TestKegagalanDiTengahTidakMeninggalkanSeparuhJadi - ADR-U-0029 Akibat 2.
func TestKegagalanDiTengahTidakMeninggalkanSeparuhJadi(t *testing.T) {
	svc, tutup := siapkanPendaftaran(t)
	defer tutup()

	db, tutupDB := repoUji(t)
	defer tutupDB()
	pohon := pohonUjiStatus(t, svc, db)
	ctx := context.Background()
	baca := repository.NewKlaimLife(db)
	peserta, _ := baca.AmbilPeserta(ctx, pohon.Work.ID)
	perBaris, _ := baca.AmbilBaris(ctx, pohon.Work.ID)
	adj := perBaris[peserta[0].ID]

	saat := time.Date(2026, 9, 26, 22, 0, 0, 0, time.UTC)
	err := svc.Status().DenganJejak(jejakGagal{}).Ubah(ctx,
		services.Pelaku{AkunID: "UJI-AKUN"},
		pohon.Work.ID, peserta[0].ID, adj[0].ID, models.StatusAksep, saat)
	if !errors.Is(err, errJejakSengaja) {
		t.Fatalf("galat = %v, mau errJejakSengaja", err)
	}

	// Jejak gagal SESUDAH ketiga tulisan; seluruhnya wajib ikut batal.
	perBaris, err = baca.AmbilBaris(ctx, pohon.Work.ID)
	if err != nil {
		t.Fatalf("membaca ulang baris: %v", err)
	}
	if got := perBaris[peserta[0].ID][0].KodeStatus; got != models.KodeOutstanding {
		t.Errorf("baris = %q sesudah transaksi batal, mau tetap %q",
			got, models.KodeOutstanding)
	}
	hdr, err := baca.AmbilHeader(ctx, pohon.Work.ID)
	if err != nil || hdr == nil {
		t.Fatalf("membaca header: %v", err)
	}
	if hdr.KodeStatus == models.KodeAksep {
		t.Error("header tertinggal Aksep padahal transaksinya batal - " +
			"pencerminan separuh jadi")
	}
}
