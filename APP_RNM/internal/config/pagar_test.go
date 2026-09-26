package config

// Pagar operasi merusak - TANPA Oracle.
//
// Untuk apa berkas ini: dua jalur dapat MENGHAPUS tabel di skema yang ditunjuk
// ORACLE_SCHEMA - test bertag `db`, dan flag `-migrate-down` di cmd/api.
// Keduanya memakai pagar yang sama, dan pagar itu tinggal di sini justru supaya
// hanya ada satu: cmd/api tidak boleh mengimpor paket penunjang test, sehingga
// menaruhnya di paket skema uji akan memaksa jalur kedua menyalin logikanya -
// dan salinan yang berselisih adalah cara paling mudah kehilangan tabel.
//
// `-migrate` TIDAK dipagari: ia hanya membuat objek, tidak pernah menghapus.
//
// Dibaca sesudah: config.go.

import (
	"errors"
	"strings"
	"testing"
)

// Pagar menuntut DUA syarat, dan menolak bila salah satunya tidak terpenuhi.
func TestPagarSkemaUjiMenuntutDuaSyarat(t *testing.T) {
	kasus := []struct {
		nama   string
		skema  string
		diakui string
		mauSah bool
	}{
		{"tanpa env sama sekali", "SKEMA_UJI_DBA", "", false},
		{"env salah eja", "SKEMA_UJI_DBA", "ya", false},
		{"env false", "SKEMA_UJI_DBA", "false", false},
		{"skema warisan meski env true", NamaSkemaWarisan, "true", false},
		{"skema warisan huruf kecil", "pooldata", "true", false},
		{"turunan skema warisan", "POOLDATA_DEV", "true", false},
		{"skema warisan tanpa env", NamaSkemaWarisan, "", false},
		{"keduanya benar", "SKEMA_UJI_DBA", "true", true},
		{"keduanya benar, env huruf besar", "SKEMA_UJI_DBA", "TRUE", true},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			err := PagarSkemaUji(k.skema, k.diakui)
			if k.mauSah {
				if err != nil {
					t.Errorf("pagar menolak yang seharusnya sah: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("pagar meloloskan skema yang seharusnya ditolak")
			}
			if !errors.Is(err, ErrBukanSkemaUji) {
				t.Errorf("galatnya bukan ErrBukanSkemaUji: %v", err)
			}
			if !strings.Contains(err.Error(), EnvSkemaUji) {
				t.Errorf("pesan tidak menyebut %s: %v", EnvSkemaUji, err)
			}
		})
	}
}

// Config membaca pengakuan itu dari env, dan meneruskannya ke pagar yang sama.
func TestPastikanSkemaUjiMembacaEnv(t *testing.T) {
	t.Setenv("ORACLE_DSN", "oracle://uji")
	t.Setenv("ORACLE_SCHEMA", "SKEMA_UJI_DBA")

	t.Setenv(EnvSkemaUji, "")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SkemaUjiDiakui {
		t.Error("env kosong seharusnya tidak diakui")
	}
	if err := cfg.PastikanSkemaUji(); err == nil {
		t.Error("tanpa env, operasi merusak seharusnya ditolak")
	}

	t.Setenv(EnvSkemaUji, "true")
	cfg, err = Load()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.SkemaUjiDiakui {
		t.Fatal("env true tidak terbaca; pembacanya yang rusak")
	}
	if err := cfg.PastikanSkemaUji(); err != nil {
		t.Errorf("skema uji yang sah seharusnya lolos: %v", err)
	}

	// Skema warisan tetap ditolak walau env-nya diakui.
	t.Setenv("ORACLE_SCHEMA", NamaSkemaWarisan)
	cfg, err = Load()
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.PastikanSkemaUji(); err == nil {
		t.Errorf("%s seharusnya ditolak meski %s=true", NamaSkemaWarisan, EnvSkemaUji)
	}
}
