package services_test

// Kasus polis baru - GILIRAN-13 paket 1. TANPA Oracle.

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

func TestBuatKasusPolisMenjagaPagarnyaBerurutan(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	w := time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)
	// Identitas lebih dulu - bahkan sebelum benderanya dinilai.
	if _, err := svc.KasusPolis().Buat(ctx, services.Pelaku{}, "9", w); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v", err)
	}
	p := services.Pelaku{AkunID: "UJI-INPUTOR"}
	// Bendera di luar "0"/"1": permintaan tidak sah, dan sebabnya dikenali.
	_, err := svc.KasusPolis().Buat(ctx, p, "2", w)
	if !errors.Is(err, services.ErrPermintaanTidakSah) || !errors.Is(err, models.ErrFlagPolisTidakSah) {
		t.Errorf("bendera asing: %v", err)
	}
	// Bendera sah, tanpa Oracle: gagal terang, bukan kasus hantu.
	if _, err := svc.KasusPolis().Buat(ctx, p, models.FlagPolisPremium, w); !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("tanpa Oracle: %v", err)
	}
}
