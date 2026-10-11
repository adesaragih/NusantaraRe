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

// 11-10-2026 - tombol Save layar Medical Check (RECOMMENDATION per peserta).

func TestPeriksaIsianMedis(t *testing.T) {
	aktif := map[string]bool{"P1": true}
	for _, u := range []struct {
		nama string
		is   services.IsianMedis
		mau  error
	}{
		{"kosong", services.IsianMedis{}, galat.ErrPermintaanTidakSah},
		{"peserta asing", services.IsianMedis{Rekomendasi: map[string]string{"P9": "x"}}, galat.ErrPermintaanTidakSah},
		{"kepanjangan", services.IsianMedis{Rekomendasi: map[string]string{"P1": strings.Repeat("a", services.BatasRekomendasi+1)}},
			services.ErrRekomendasiKepanjangan},
		// Byte, bukan karakter: 334 x "é" (2 byte) = 668 byte lolos; 501 x = 1002 byte ditolak.
		{"multibyte lolos", services.IsianMedis{Rekomendasi: map[string]string{"P1": strings.Repeat("é", 334)}}, nil},
		{"multibyte kepanjangan", services.IsianMedis{Rekomendasi: map[string]string{"P1": strings.Repeat("é", 501)}},
			services.ErrRekomendasiKepanjangan},
		{"tepat batas", services.IsianMedis{Rekomendasi: map[string]string{"P1": strings.Repeat("a", services.BatasRekomendasi)}}, nil},
		{"dikosongkan", services.IsianMedis{Rekomendasi: map[string]string{"P1": ""}}, nil},
	} {
		err := services.PeriksaIsianMedis(u.is, aktif)
		if (u.mau == nil && err != nil) || (u.mau != nil && !errors.Is(err, u.mau)) {
			t.Errorf("%s: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
}

func TestSimpanMedisMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	is := services.IsianMedis{Rekomendasi: map[string]string{"P1": "UJI"}}
	if err := svc.SimpanMedis().Simpan(ctx, inti.Pelaku{}, "CLM-1", is); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("anonim: %v, mau ErrTanpaIdentitas", err)
	}
	medis := inti.Pelaku{AkunID: "UJI-AKUN", Peran: []string{inti.PeranMedicalAdvisor}}
	if err := svc.SimpanMedis().Simpan(ctx, medis, " ", is); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("tanpa pengenal: %v, mau ErrPermintaanTidakSah", err)
	}
	if err := svc.SimpanMedis().Simpan(ctx, medis, "CLM-1", is); !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("tanpa database: %v, mau ErrTanpaOracle", err)
	}
}

// Gerbang tahap + pemegang + kasus terbuka tetap DIPANGGIL (statik; kebenarannya di models dan uji db).
func TestSimpanMedisMemanggilGerbang(t *testing.T) {
	isi, err := os.ReadFile("medis.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{"WajibIdentitas(pelaku)", "PastikanKasusTerbuka(ctx, klaimID)",
		"tahap != models.TahapMedicalCheck", "WajibPemegangTahap(pelaku, tahap)", "PeriksaIsianMedis(is, aktif)",
		"SetelRekomendasi(ctx, tx, klaimID, id, v)"} {
		if !strings.Contains(string(isi), mau) {
			t.Errorf("medis.go tidak lagi memanggil %s", mau)
		}
	}
}
