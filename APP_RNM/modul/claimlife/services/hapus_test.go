package services_test

// Penghapusan klaim - TANPA Oracle.
//
// Pemilik: tiket 15. Dibaca sesudah: hapus.go.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/services"
)

func pelakuHapus() inti.Pelaku {
	return inti.Pelaku{AkunID: "UJI-AKUN", Peran: []string{services.PeranHapusKlaim}}
}

// TestDampakMenghitungTiapJenisTerpisah - satu angka total menyembunyikan
// tingkat mana yang ternyata lebih besar dari dugaan, dan justru itu yang
// membuat pengguna menekan Batal.
func TestDampakMenghitungTiapJenisTerpisah(t *testing.T) {
	d := models.DampakHapus{
		Header: 1, Peserta: 2, Adjustment: 4, Spreading: 8, SpreadingRetro: 16,
		Dokumen: 3, WorkClaim: 1, BarisDatarWarisan: 5,
	}
	// ⛔ Total TIDAK memuat baris datar warisan: nasibnya belum diputuskan
	// (AC 11, `[terbuka - work owner]`), dan menjumlahkannya ke dalam "yang
	// akan ikut terhapus" berarti menjawab pertanyaan itu lewat sebuah angka.
	if got := d.Total(); got != 35 {
		t.Errorf("Total = %d, mau 35 (tanpa baris datar warisan)", got)
	}
	if got := d.TotalTermasukWarisan(); got != 40 {
		t.Errorf("TotalTermasukWarisan = %d, mau 40", got)
	}
	var nol models.DampakHapus
	if !nol.Kosong() {
		t.Error("dampak nol tidak dinyatakan kosong")
	}
	if d.Kosong() {
		t.Error("dampak berisi dinyatakan kosong")
	}
	// Rinciannya wajib terbaca satu per satu di jejak audit dan di layar.
	teks := d.String()
	for _, potong := range []string{"header 1", "peserta 2", "adjustment 4", "spreading 8",
		"spreading retro 16", "dokumen 3", "baris work 1", "baris datar warisan 5"} {
		if !strings.Contains(teks, potong) {
			t.Errorf("rincian %q tidak ada di %q", potong, teks)
		}
	}
}

// TestHapusMenuntutPeranDanIdentitas - dua pertanyaan berbeda, dua galat.
func TestHapusMenuntutPeranDanIdentitas(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if _, err := svc.Penghapusan().Hapus(ctx, inti.Pelaku{
		Peran: []string{services.PeranHapusKlaim}}, "CLM-1"); !errors.Is(
		err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: galat = %v, mau ErrTanpaIdentitas", err)
	}
	if _, err := svc.Penghapusan().Hapus(ctx,
		inti.Pelaku{AkunID: "UJI-AKUN"}, "CLM-1"); !errors.Is(
		err, inti.ErrTanpaWewenang) {
		t.Errorf("tanpa peran: galat = %v, mau ErrTanpaWewenang", err)
	}
	if _, err := svc.Penghapusan().Hapus(ctx, pelakuHapus(), "CLM-1"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: galat = %v, mau ErrTanpaOracle", err)
	}
}

// TestDampakTidakMenULIS - jalur pratinjau tidak punya satu pun tulisan untuk
// dibatalkan; itulah yang membuat "Batal" benar-benar membatalkan.
func TestDampakTidakMenulis(t *testing.T) {
	svc := services.New(nil)
	_, err := svc.Penghapusan().Dampak(context.Background(), pelakuHapus(), "CLM-1")
	if !errors.Is(err, db.ErrTanpaOracle) {
		t.Fatalf("galat = %v, mau ErrTanpaOracle", err)
	}
	// Pengenal kosong ditolak sebelum menyentuh apa pun.
	if _, err := svc.Penghapusan().Dampak(
		context.Background(), pelakuHapus(), "  "); !errors.Is(
		err, galat.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: galat = %v, mau ErrPermintaanTidakSah", err)
	}
}
