package services_test

// Cabut peserta - OQ-M6 ditutup 29-09-2026 (GILIRAN-17). Gerbang sebelum
// basis data; jalur penuh di uji `db`.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

func TestCabutPesertaMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if err := svc.Status().CabutPeserta(ctx, services.Pelaku{}, "CLM-1", "P-1", saatUji); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v, mau ErrTanpaIdentitas", err)
	}
	for _, peran := range []string{services.PeranMedicalAdvisor, services.PeranSPV} {
		if err := svc.Status().CabutPeserta(ctx, pelakuBerperan(peran), "CLM-1", "P-1", saatUji); !errors.Is(err, services.ErrTanpaWewenang) {
			t.Errorf("%s: %v, mau ErrTanpaWewenang", peran, err)
		}
	}
	for _, k := range []struct{ klaim, peserta string }{{"", "P-1"}, {"CLM-1", " "}} {
		if err := svc.Status().CabutPeserta(ctx, pelakuAdmin(), k.klaim, k.peserta, saatUji); !errors.Is(err, services.ErrPermintaanTidakSah) {
			t.Errorf("%q/%q: %v, mau ErrPermintaanTidakSah", k.klaim, k.peserta, err)
		}
	}
	if err := svc.Status().CabutPeserta(ctx, pelakuAdmin(), "CLM-1", "P-1", saatUji); !errors.Is(err, repository.ErrTanpaOracle) {
		t.Errorf("admin: %v, mau ErrTanpaOracle", err)
	}
}
