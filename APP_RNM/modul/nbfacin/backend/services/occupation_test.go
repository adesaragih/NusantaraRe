package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

type occupationTiruan struct{ kata *string }

func (o occupationTiruan) CariOccupation(_ context.Context, k string) ([]models.BarisOccupation, error) {
	*o.kata = k
	return []models.BarisOccupation{{OldID: "UJI01", Name: "UJI PABRIK"}}, nil
}

// TestCariOccupation - tiket 38: >= 2 karakter (dipangkas) dan <= 255, tanpa identitas, 503
// tanpa basis data.
func TestCariOccupation(t *testing.T) {
	ctx := context.Background()
	var kata string
	svc := Baru(nil).DenganOccupation(occupationTiruan{&kata})
	if b, err := svc.CariOccupation(ctx, " pa "); err != nil || kata != "pa" || len(b) != 1 || b[0].OldID != "UJI01" {
		t.Errorf("cari: %v %v %q", b, err, kata)
	}
	for _, k := range []string{"", "  ", "p", " p ", strings.Repeat("p", 256)} {
		if _, err := svc.CariOccupation(ctx, k); !errors.Is(err, ErrMasukanOccupation) {
			t.Errorf("cari %q: %v", k, err)
		}
	}
	if _, err := Baru(nil).CariOccupation(ctx, "pa"); !errors.Is(err, ErrOccupationTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
}
