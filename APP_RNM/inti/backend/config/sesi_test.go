package config

import (
	"errors"
	"strings"
	"testing"
)

// Login (keputusan work owner 01-10-2026): SESI_RAHASIA pendek dan cookie
// tanpa Secure di produksi ditolak SAAT MENYALA.
func TestKonfigurasiSesiLogin(t *testing.T) {
	t.Setenv("ORACLE_DSN", "")
	t.Setenv("IS_PEGA_PROD", "false")
	t.Setenv("AUTH_STUB", "")
	t.Setenv("SESI_COOKIE_SECURE", "")

	t.Setenv("SESI_RAHASIA", "")
	cfg, err := Load()
	if err != nil || cfg.SesiRahasia != "" || !cfg.SesiCookieAman {
		t.Fatalf("bawaan: %+v %v - rahasia kosong, cookie Secure", cfg, err)
	}

	t.Setenv("SESI_RAHASIA", strings.Repeat("r", PanjangMinSesiRahasia-1))
	if _, err := Load(); !errors.Is(err, ErrKonfigurasi) {
		t.Errorf("rahasia %d karakter: %v", PanjangMinSesiRahasia-1, err)
	}
	t.Setenv("SESI_RAHASIA", strings.Repeat("r", PanjangMinSesiRahasia))
	if cfg, err := Load(); err != nil || len(cfg.SesiRahasia) != PanjangMinSesiRahasia {
		t.Errorf("rahasia %d karakter: %v", PanjangMinSesiRahasia, err)
	}

	t.Setenv("SESI_COOKIE_SECURE", "false")
	if cfg, err := Load(); err != nil || cfg.SesiCookieAman {
		t.Errorf("SESI_COOKIE_SECURE=false di non-produksi: %+v %v", cfg.SesiCookieAman, err)
	}
	t.Setenv("IS_PEGA_PROD", "true")
	if _, err := Load(); !errors.Is(err, ErrKonfigurasi) {
		t.Errorf("SESI_COOKIE_SECURE=false di produksi: %v", err)
	}
}
