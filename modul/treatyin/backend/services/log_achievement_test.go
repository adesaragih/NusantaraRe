package services_test

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

// Gudang tiruan bersama: Submit Achievement tidak dipakai uji lain.
func (g *gudangTiruan) CatatLogAchievement(_ context.Context, _, _ string, _ []repository.BarisLogSiap) error {
	return nil
}

// gudangLog - menangkap baris log yang tombol Submit kirim.
type gudangLog struct {
	*gudangTiruan
	master, operator string
	baris            []repository.BarisLogSiap
}

func (g *gudangLog) CatatLogAchievement(_ context.Context, master, operator string, b []repository.BarisLogSiap) error {
	g.master, g.operator, g.baris = master, operator, b
	return nil
}

func TestSubmitAchievementMenyisipkanTiapBarisKecualiQuarterKosong(t *testing.T) {
	g := &gudangLog{gudangTiruan: &gudangTiruan{}}
	h, err := services.LayananDengan(g).CatatLogAchievement(context.Background(), admin, services.MasukanLogAchievement{
		IDKontrak: "1000080/R01",
		Baris: []models.BarisLogAchievement{
			{Quarter: "1", QUARTERYEAR: "2024", CurrencyID: "10000", Currency: "IDR", PREMIUM: "1000.50", LossRatio: "12.3456"},
			{Quarter: "", PREMIUM: "999"},
			{Quarter: "2", QUARTERYEAR: "2024", Currency: "USD", PREMIUM: "", Total: "-5"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if h.Disisipkan != 2 || h.Dilewati != 1 {
		t.Fatalf("hasil = %+v, ingin 2 disisipkan dan 1 dilewati", h)
	}
	if g.master != "1000080/R01" || g.operator != admin.AkunID {
		t.Errorf("master/operator = %q/%q", g.master, g.operator)
	}
	if g.baris[0].Angka[0].String() != "1000.50" || g.baris[0].Angka[9].String() != "12.3456" {
		t.Errorf("angka baris 1 = %v / %v", g.baris[0].Angka[0], g.baris[0].Angka[9])
	}
	if g.baris[1].Angka[0] != nil {
		t.Error("Premium kosong wajib NULL, bukan 0")
	}
	if g.baris[1].Angka[8].String() != "-5" {
		t.Errorf("Total baris 2 = %v", g.baris[1].Angka[8])
	}
}

func TestSubmitAchievementMenolakAngkaTakTerbaca(t *testing.T) {
	g := &gudangLog{gudangTiruan: &gudangTiruan{}}
	_, err := services.LayananDengan(g).CatatLogAchievement(context.Background(), admin, services.MasukanLogAchievement{
		IDKontrak: "1000080", Baris: []models.BarisLogAchievement{{Quarter: "1", PREMIUM: "1.000,00"}},
	})
	if !errors.Is(err, services.ErrTombolDitolak) {
		t.Fatalf("galat = %v, ingin ErrTombolDitolak", err)
	}
	if g.baris != nil {
		t.Error("baris tertulis padahal ditolak")
	}
}

func TestSubmitAchievementTanpaPengenalDitolak(t *testing.T) {
	_, err := services.LayananDengan(&gudangTiruan{}).CatatLogAchievement(context.Background(), admin, services.MasukanLogAchievement{})
	if !errors.Is(err, services.ErrTombolDitolak) {
		t.Fatalf("galat = %v", err)
	}
}
