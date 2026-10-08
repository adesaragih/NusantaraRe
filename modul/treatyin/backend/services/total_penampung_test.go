package services_test

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// gudangTotal - gudang tiruan yang MEMBAWA larik total tersimpan.
type gudangTotal struct {
	*gudangTiruan
	total map[string][]map[string]any
	galat error
}

func (g gudangTotal) BacaTotalPenampung(_ context.Context, _ string) (map[string][]map[string]any, error) {
	return g.total, g.galat
}

// ⭐ Laporan pemakai 8 Oktober 2026: total Share Prop "No items" sesudah
// Save sampai Refresh ditekan. Total yang TERSIMPAN wajib sampai ke jawaban
// kontrak, berkunci ejaan Pega, supaya form menyemainya ke penampung.
func TestTotalTersimpanIkutDibukaKontrak(t *testing.T) {
	tot := map[string][]map[string]any{
		"TotalShareRnmProp":    {{"Currency": "IDR", "CurrencyID": "10000", "Value": "15000000000"}},
		"TotalSpreadedRnmProp": {{"Currency": "IDR", "CurrencyID": "10000", "Value": "6000000000"}},
	}
	g := gudangTotal{gudangTiruan: &gudangTiruan{kontrakWarisan: kontrakWarisanUji()}, total: tot}
	k, err := services.LayananDengan(g).BacaKontrakWarisan(context.Background(), pelakuAda, "1000797")
	if err != nil {
		t.Fatal(err)
	}
	if got := k.PenampungLarik["TotalShareRnmProp"]; len(got) != 1 || got[0]["Value"] != "15000000000" {
		t.Errorf("TotalShareRnmProp = %v", got)
	}
	if got := k.PenampungLarik["TotalSpreadedRnmProp"]; len(got) != 1 || got[0]["Value"] != "6000000000" {
		t.Errorf("TotalSpreadedRnmProp = %v", got)
	}
	// Larik yang tidak tersimpan TIDAK dikarang kosong — tab menyemai sendiri.
	if _, ada := k.PenampungLarik["TotalSpreadedRnmRIProp"]; ada {
		t.Error("larik tanpa baris tidak boleh dikirim")
	}
}

// ⛔ Galat membaca total diteruskan — total kosong dan total yang gagal
// dibaca terlihat sama di layar.
func TestGalatBacaTotalDiteruskan(t *testing.T) {
	rusak := errors.New("oracle tidak menjawab")
	g := gudangTotal{gudangTiruan: &gudangTiruan{kontrakWarisan: kontrakWarisanUji()}, galat: rusak}
	if _, err := services.LayananDengan(g).BacaKontrakWarisan(context.Background(), pelakuAda, "1000797"); !errors.Is(err, rusak) {
		t.Fatalf("galat = %v, ingin %v", err, rusak)
	}
}
