package services_test

// Kasus polis baru - GILIRAN-13 paket 1. TANPA Oracle.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/inti/galat"
	"nusantarare/modul/premiumlistlife/models"
	"nusantarare/modul/premiumlistlife/services"
)

func TestBuatKasusPolisMenjagaPagarnyaBerurutan(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	w := time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)
	// Identitas lebih dulu - bahkan sebelum benderanya dinilai.
	if _, err := svc.KasusPolis().Buat(ctx, inti.Pelaku{}, "9", w); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v", err)
	}
	p := inti.Pelaku{AkunID: "UJI-INPUTOR"}
	// Bendera di luar "0"/"1": permintaan tidak sah, dan sebabnya dikenali.
	_, err := svc.KasusPolis().Buat(ctx, p, "2", w)
	if !errors.Is(err, galat.ErrPermintaanTidakSah) || !errors.Is(err, models.ErrFlagPolisTidakSah) {
		t.Errorf("bendera asing: %v", err)
	}
	// Bendera sah, tanpa Oracle: gagal terang, bukan kasus hantu.
	if _, err := svc.KasusPolis().Buat(ctx, p, models.FlagPolisPremium, w); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: %v", err)
	}
}
