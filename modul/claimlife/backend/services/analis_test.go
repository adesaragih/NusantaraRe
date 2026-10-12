package services_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/claimlife/backend/services"
)

// 12-10-2026 - tombol Save layar Claim Analis (CLAIM ANALYST / NOTES per peserta).

func TestPeriksaIsianAnalis(t *testing.T) {
	aktif := map[string]bool{"P1": true}
	for _, u := range []struct {
		nama string
		is   services.IsianAnalis
		mau  error
	}{
		{"kosong", services.IsianAnalis{}, galat.ErrPermintaanTidakSah},
		{"peserta asing", services.IsianAnalis{Catatan: map[string]string{"P9": "x"}}, galat.ErrPermintaanTidakSah},
		{"kepanjangan", services.IsianAnalis{Catatan: map[string]string{"P1": strings.Repeat("a", services.BatasCatatanAnalis+1)}},
			services.ErrCatatanAnalisKepanjangan},
		// Byte, bukan karakter: 334 x "é" (2 byte) = 668 byte lolos; 501 x = 1002 byte ditolak.
		{"multibyte lolos", services.IsianAnalis{Catatan: map[string]string{"P1": strings.Repeat("é", 334)}}, nil},
		{"multibyte kepanjangan", services.IsianAnalis{Catatan: map[string]string{"P1": strings.Repeat("é", 501)}},
			services.ErrCatatanAnalisKepanjangan},
		{"tepat batas", services.IsianAnalis{Catatan: map[string]string{"P1": strings.Repeat("a", services.BatasCatatanAnalis)}}, nil},
		{"dikosongkan", services.IsianAnalis{Catatan: map[string]string{"P1": ""}}, nil},
	} {
		err := services.PeriksaIsianAnalis(u.is, aktif)
		if (u.mau == nil && err != nil) || (u.mau != nil && !errors.Is(err, u.mau)) {
			t.Errorf("%s: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
}

func TestSimpanAnalisMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	is := services.IsianAnalis{Catatan: map[string]string{"P1": "UJI"}}
	if err := svc.SimpanAnalis().Simpan(ctx, inti.Pelaku{}, "CLM-1", is); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v, mau ErrTanpaIdentitas", err)
	}
	spv := inti.Pelaku{AkunID: "UJI-AKUN", Peran: []string{inti.PeranSPV}}
	if err := svc.SimpanAnalis().Simpan(ctx, spv, " ", is); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("tanpa pengenal: %v, mau ErrPermintaanTidakSah", err)
	}
	if err := svc.SimpanAnalis().Simpan(ctx, spv, "CLM-1", is); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa database: %v, mau ErrTanpaOracle", err)
	}
}

// Gerbang tahap + pemegang + kasus terbuka tetap DIPANGGIL (statik; kebenarannya di models dan uji db).
func TestSimpanAnalisMemanggilGerbang(t *testing.T) {
	isi, err := os.ReadFile("analis.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"WajibIdentitas(pelaku)", "PastikanKasusTerbuka(ctx, klaimID)",
		"tahap != models.TahapClaimAnalis", "WajibPemegangTahap(pelaku, tahap)", "PeriksaIsianAnalis(is, aktif)",
		"SetelCatatanAnalis(ctx, tx, klaimID, id, v)"} {
		if !strings.Contains(string(isi), mau) {
			t.Errorf("analis.go tidak lagi memanggil %s", mau)
		}
	}
}
