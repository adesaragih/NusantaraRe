package main

import (
	"os"
	"strings"
	"testing"
)

// Butir 3: penolakan IS_PEGA_PROD=true terjadi SEBELUM koneksi dibuka, dan bawaan alat = uji kering.
func TestTolakProduksiSebelumKoneksi(t *testing.T) {
	isi, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	periksa, buka := strings.Index(s, "repository.PeriksaMode(*jalankan, cfg.IsPegaProd)"), strings.Index(s, "db.Open(cfg)")
	if periksa < 0 || buka < 0 || periksa > buka {
		t.Errorf("PeriksaMode (%d) harus mendahului db.Open (%d)", periksa, buka)
	}
	if !strings.Contains(s, `flag.Bool("jalankan", false,`) {
		t.Error("-jalankan harus bawaan false (uji kering)")
	}
}
