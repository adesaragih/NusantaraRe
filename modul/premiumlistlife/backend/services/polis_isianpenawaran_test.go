package services_test

// Form penawaran - tiket 01 bagian 3. TANPA Oracle: yang diuji urutan pagarnya.

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/premiumlistlife/backend/models"
	"nusantarare/modul/premiumlistlife/backend/services"
)

func TestFormPenawaranMenjagaPagarnyaBerurutan(t *testing.T) {
	f := services.New(nil).FormPenawaran()
	ctx := context.Background()
	w := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	p := inti.Pelaku{AkunID: "UJI-INPUTOR"}
	// Identitas lebih dulu.
	if _, err := f.Baca(ctx, inti.Pelaku{}, "NBLF-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("baca anonim: %v", err)
	}
	if err := f.Simpan(ctx, inti.Pelaku{}, "NBLF-1", models.IsianPenawaran{}, w); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("simpan anonim: %v", err)
	}
	if _, err := f.CariCeding(ctx, inti.Pelaku{}, "a"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("cari anonim: %v", err)
	}
	// Pengenal kosong dijawab "wajib diisi", bukan "ORACLE_DSN".
	if err := f.Simpan(ctx, p, "  ", models.IsianPenawaran{}, w); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: %v", err)
	}
	// Bentuk sah, tanpa Oracle: gagal terang.
	if _, err := f.Baca(ctx, p, "NBLF-1"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("baca tanpa Oracle: %v", err)
	}
	if _, err := f.CariPemegangPolis(ctx, p, "a"); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("cari tanpa Oracle: %v", err)
	}
}

// Confirm di tahap penawaran digerbangi kelengkapan isian; Decline tidak.
func TestConfirmPenawaranDigerbangiKelengkapan(t *testing.T) {
	isi, err := os.ReadFile("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	teks := string(isi)
	if !strings.Contains(teks, "keadaan.Status == models.TahapPolisPenawaran && keputusan == models.KeputusanConfirm") ||
		!strings.Contains(teks, "p.periksaPenawaranLengkap(ctx, keadaan.ID)") {
		t.Error("Confirm tahap penawaran tidak memeriksa isian Input Offer")
	}
}
