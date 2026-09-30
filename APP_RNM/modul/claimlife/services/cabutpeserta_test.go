package services_test

// Cabut peserta - OQ-M6 ditutup 29-09-2026 (GILIRAN-17). Gerbang sebelum
// basis data; jalur penuh di uji `db`.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/services"
)

func TestCabutPesertaMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if err := svc.Status().CabutPeserta(ctx, inti.Pelaku{}, "CLM-1", "P-1", saatUji); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v, mau ErrTanpaIdentitas", err)
	}
	for _, peran := range []string{inti.PeranMedicalAdvisor, inti.PeranSPV} {
		if err := svc.Status().CabutPeserta(ctx, pelakuBerperan(peran), "CLM-1", "P-1", saatUji); !errors.Is(err, inti.ErrTanpaWewenang) {
			t.Errorf("%s: %v, mau ErrTanpaWewenang", peran, err)
		}
	}
	for _, k := range []struct{ klaim, peserta string }{{"", "P-1"}, {"CLM-1", " "}} {
		if err := svc.Status().CabutPeserta(ctx, pelakuAdmin(), k.klaim, k.peserta, saatUji); !errors.Is(err, galat.ErrPermintaanTidakSah) {
			t.Errorf("%q/%q: %v, mau ErrPermintaanTidakSah", k.klaim, k.peserta, err)
		}
	}
	if err := svc.Status().CabutPeserta(ctx, pelakuAdmin(), "CLM-1", "P-1", saatUji); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("admin: %v, mau ErrTanpaOracle", err)
	}
}
