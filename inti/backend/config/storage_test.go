package config

// Uji MODUL_AKTIF. PELAKSANA_STORAGE dan TCO_PEKERJA_LAMPIRAN_INTERVAL dibuang
// bersama fitur lampiran Treaty Contract Out (keputusan work owner 08-10-2026).

import (
	"reflect"
	"testing"
)

func kosongkanEnvDasar(t *testing.T) {
	t.Helper()
	for _, k := range []string{"ORACLE_DSN", "ORACLE_SCHEMA", "AUTH_STUB", "IS_PEGA_PROD",
		"STORAGE_TOKEN_SALT"} {
		t.Setenv(k, "")
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
