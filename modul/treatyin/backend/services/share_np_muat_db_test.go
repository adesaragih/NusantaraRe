//go:build db

package services_test

import (
	"context"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

// ⭐ Kontrak dibuka utuh lewat services: kolom 445 belum ada (toleran),
// RNM Share dari salinan baris (1001270) atau TREATYINDETAIL (1000402).
func TestMuatKontrakShareNPUrutanSumber(t *testing.T) {
	cfg, _ := config.Load()
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	l := services.LayananDengan(repository.Baru(d))
	for id, mau := range map[string]string{"1001270": "40", "1000402": "10"} {
		k, err := l.BacaKontrakWarisan(context.Background(), inti.Pelaku{AkunID: "AKUN-UJI-DB"}, id)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		if k.ShareNP.RNMShare != mau {
			t.Errorf("%s: RNMShare %q, mau %q", id, k.ShareNP.RNMShare, mau)
		}
	}
}
