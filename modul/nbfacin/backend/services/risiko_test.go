package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

type riskTiruan struct {
	saring         *models.SaringRisk
	offset, ukuran *int
}

func (r riskTiruan) CariRisk(_ context.Context, s models.SaringRisk, offset, ukuran int) ([]models.RiskAddress, int, error) {
	*r.saring, *r.offset, *r.ukuran = s, offset, ukuran
	return []models.RiskAddress{{ID: "UJI-R1"}}, 37, nil
}

// TestCariRisk - tiket 36: tanpa saringan = 400 (SearchRiskAddressAct), ukuran 10, saringan
// diteruskan apa adanya, 503 tanpa DB, halaman/panjang tak sah = 400.
func TestCariRisk(t *testing.T) {
	var s models.SaringRisk
	var offset, ukuran int
	svc := Baru(nil).DenganRisk(riskTiruan{&s, &offset, &ukuran})
	h, err := svc.CariRisk(context.Background(), models.SaringRisk{City: "UJI KOTA"}, 3)
	if err != nil || h.Ukuran != 10 || h.Total != 37 || offset != 20 || ukuran != 10 || s.City != "UJI KOTA" {
		t.Fatalf("%+v (%v) offset %d", h, err, offset)
	}
	for nama, u := range map[string]struct {
		s       models.SaringRisk
		halaman int
		pesan   string
	}{
		"tanpa saringan":  {models.SaringRisk{}, 1, "minimal satu saringan"},
		"spasi saja":      {models.SaringRisk{Address: "  ", Territory: " "}, 1, "minimal satu saringan"},
		"halaman 0":       {models.SaringRisk{City: "UJI"}, 0, "halaman"},
		"terlalu panjang": {models.SaringRisk{Address: strings.Repeat("U", 256)}, 1, "address paling banyak 255"},
	} {
		if _, err := svc.CariRisk(context.Background(), u.s, u.halaman); !errors.Is(err, ErrMasukanRisk) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if _, err := Baru(nil).CariRisk(context.Background(), models.SaringRisk{City: "UJI"}, 1); !errors.Is(err, ErrRiskTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	if _, err := Baru(nil).CariRisk(context.Background(), models.SaringRisk{}, 1); !errors.Is(err, ErrMasukanRisk) {
		t.Errorf("tanpa saringan, tanpa DB: %v - mau 400 dulu", err)
	}
}
