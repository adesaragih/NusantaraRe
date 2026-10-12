package services

import (
	"context"
	"errors"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// mbTiruan - PembacaMarketingOfficer + PembacaMarketingBawaan: hasil = kandidat urut ID turun (repository); dibaca mencatat panggilan.
type mbTiruan struct {
	hasil  []string
	galat  error
	dibaca *int
}

func (m mbTiruan) DaftarMarketingOfficer(context.Context) ([]models.MarketingOfficer, error) {
	return nil, nil
}
func (m mbTiruan) MarketingBawaan(context.Context, string) ([]string, error) {
	*m.dibaca++
	return m.hasil, m.galat
}

// TestMarketingBawaan - keputusan work owner 12-10-2026: Marketing Name otomatis dari pembuat case bila MOID belum tersimpan.
// MOID terisi -> "" tanpa membaca; satu cocok -> ID itu; tidak cocok -> ""; dua cocok -> "" (koreksi work owner 12-10-2026);
// galat baca -> "" dan GET kasus TIDAK gagal; pembaca tanpa MarketingBawaan -> "".
func TestMarketingBawaan(t *testing.T) {
	ctx := context.Background()
	for nama, c := range map[string]struct {
		moid  string
		hasil []string
		galat error
		mau   string
		baca  int
	}{
		"MOID terisi": {"UJI-MO-1", []string{"UJI-MO-9"}, nil, "", 0},
		"satu cocok":  {"", []string{"UJI-MO-7"}, nil, "UJI-MO-7", 1},
		"tidak cocok": {"", nil, nil, "", 1},
		"dua cocok":   {"", []string{"UJI-MO-9", "UJI-MO-3"}, nil, "", 1},
		"galat baca":  {"", nil, errors.New("ORA-00942"), "", 1},
		"MOID spasi":  {"  ", []string{" UJI-MO-5 "}, nil, "UJI-MO-5", 1},
	} {
		ada := kasusUji()
		ada["UJI-NB-1"].General.MarketingID = c.moid
		n := 0
		svc := Baru(nil).DenganKasus(kasusTiruan{ada: ada}).DenganMarketingOfficer(mbTiruan{hasil: c.hasil, galat: c.galat, dibaca: &n})
		k, err := svc.BacaKasus(ctx, "UJI-NB-1")
		if err != nil || k.MarketingBawaan != c.mau || n != c.baca {
			t.Errorf("%s: %q dibaca %d (%v)", nama, k.MarketingBawaan, n, err)
		}
	}
	if k, err := Baru(nil).DenganKasus(kasusTiruan{ada: kasusUji()}).DenganMarketingOfficer(marketingTiruan{}).BacaKasus(ctx, "UJI-NB-1"); err != nil ||
		k.MarketingBawaan != "" {
		t.Errorf("tanpa PembacaMarketingBawaan: %q %v", k.MarketingBawaan, err)
	}
}
