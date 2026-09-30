package config

// Uji PELAKSANA_STORAGE dan TCO_PEKERJA_LAMPIRAN_INTERVAL (OQ-TCO-08/09).
// Garam di sini palsu (awalan UJI-).

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func kosongkanEnvDasar(t *testing.T) {
	t.Helper()
	for _, k := range []string{"ORACLE_DSN", "ORACLE_SCHEMA", "AUTH_STUB", "IS_PEGA_PROD",
		"PELAKSANA_STORAGE", "STORAGE_TOKEN_SALT", EnvIntervalPekerjaLampiranTCO} {
		t.Setenv(k, "")
	}
}

func TestPelaksanaStorageBawaanStub(t *testing.T) {
	kosongkanEnvDasar(t)
	for _, v := range []string{"", "stub", " STUB "} {
		t.Setenv("PELAKSANA_STORAGE", v)
		c, err := Load()
		if err != nil || c.PelaksanaStorage != PelaksanaStorageStub {
			t.Errorf("%q: %q %v", v, c.PelaksanaStorage, err)
		}
	}
	t.Setenv("STORAGE_TOKEN_SALT", "UJI-GARAM-PALSU")
	t.Setenv("PELAKSANA_STORAGE", "Nyata")
	if c, err := Load(); err != nil || c.PelaksanaStorage != PelaksanaStorageNyata {
		t.Errorf("nyata: %q %v", c.PelaksanaStorage, err)
	}
}

// ⛔ Salah ketik tidak jatuh diam-diam ke stub, dan `nyata` tanpa garam ditolak
// saat menyala - pesannya menyebut NAMA kunci, tidak pernah nilainya.
func TestPelaksanaStorageDitolak(t *testing.T) {
	kosongkanEnvDasar(t)
	t.Setenv("STORAGE_TOKEN_SALT", "UJI-GARAM-PALSU")
	t.Setenv("PELAKSANA_STORAGE", "nyta")
	if _, err := Load(); !errors.Is(err, ErrKonfigurasi) || strings.Contains(err.Error(), "UJI-GARAM-PALSU") {
		t.Errorf("salah ketik: %v", err)
	}
	for _, garam := range []string{"", "   "} {
		t.Setenv("STORAGE_TOKEN_SALT", garam)
		t.Setenv("PELAKSANA_STORAGE", "nyata")
		_, err := Load()
		if !errors.Is(err, ErrKonfigurasi) || !strings.Contains(err.Error(), "STORAGE_TOKEN_SALT") {
			t.Errorf("nyata tanpa garam %q: %v", garam, err)
		}
	}
}

func TestIntervalPekerjaLampiranTCO(t *testing.T) {
	kosongkanEnvDasar(t)
	for raw, mau := range map[string]time.Duration{"": 0, "0": 0, "1s": time.Second, "30s": 30 * time.Second, " 1m ": time.Minute} {
		t.Setenv(EnvIntervalPekerjaLampiranTCO, raw)
		c, err := Load()
		if err != nil || c.IntervalPekerjaLampiranTCO != mau {
			t.Errorf("%q: %v %v", raw, c.IntervalPekerjaLampiranTCO, err)
		}
	}
	// Batas bawah 1s: `30ms` hampir pasti salah ketik untuk `30s`.
	for _, raw := range []string{"30", "sebentar", "-1s", "30ms", "1ns"} {
		t.Setenv(EnvIntervalPekerjaLampiranTCO, raw)
		if _, err := Load(); !errors.Is(err, ErrKonfigurasi) || !strings.Contains(err.Error(), EnvIntervalPekerjaLampiranTCO) {
			t.Errorf("%q diterima: %v", raw, err)
		}
	}
}

// TestModulAktifDiurai - MODUL_AKTIF (refactor bentuk B paket 6): dipisah koma,
// spasi dan huruf besar dibuang, ganda dan kosong dilewati; kosong = nil (semua).
func TestModulAktifDiurai(t *testing.T) {
	kosongkanEnvDasar(t)
	for raw, mau := range map[string][]string{
		"":          nil,
		" , ":       nil,
		"claimlife": {"claimlife"},
		"KomiteClaimLife, claimlife,,KOMITECLAIMLIFE ": {"komiteclaimlife", "claimlife"},
	} {
		t.Setenv("MODUL_AKTIF", raw)
		c, err := Load()
		if err != nil {
			t.Fatalf("%q: %v", raw, err)
		}
		if !reflect.DeepEqual(c.ModulAktif, mau) {
			t.Errorf("MODUL_AKTIF=%q -> %v, mau %v", raw, c.ModulAktif, mau)
		}
	}
}
