package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

type rwTiruan struct {
	zip    *string
	alamat *models.AlamatBaru
	err    error
}

func (r rwTiruan) CariRW(_ context.Context, z string) ([]models.BarisRW, error) {
	*r.zip = z
	return []models.BarisRW{{ZipCode: z + "0"}}, nil
}

func (r rwTiruan) SisipRiskAddress(_ context.Context, _ *db.Tx, a models.AlamatBaru) (string, error) {
	*r.alamat = a
	return "UJI000000000001", r.err
}

// TestRW - tiket 37: saran butuh >= 3 karakter (dipangkas), simpan butuh identitas, postalCode
// dan address wajib, lebar 4000 bita, ID baru dikembalikan (W-3), galat repository diteruskan.
func TestRW(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var zip string
	var alamat models.AlamatBaru
	svc := Baru(nil).DenganRW(rwTiruan{zip: &zip, alamat: &alamat}).DenganTransaksi(tanpaTx)
	if b, err := svc.CariRW(ctx, " 123 "); err != nil || zip != "123" || len(b) != 1 {
		t.Errorf("saran: %v %v %q", b, err, zip)
	}
	for _, z := range []string{"", "  ", "12", strings.Repeat("1", 256)} {
		if _, err := svc.CariRW(ctx, z); !errors.Is(err, ErrMasukanRW) {
			t.Errorf("zip %q: %v", z, err)
		}
	}
	if _, err := Baru(nil).CariRW(ctx, "123"); !errors.Is(err, ErrRWTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	sah := models.AlamatBaru{Title: "DESA", Address: "UJI JALAN", PostalCode: "00000"}
	if id, err := svc.SimpanAlamatBaru(ctx, akun, sah); err != nil || id != "UJI000000000001" || alamat != sah {
		t.Fatalf("simpan: %q %v %+v", id, err, alamat)
	}
	if _, err := svc.SimpanAlamatBaru(ctx, inti.Pelaku{}, sah); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for nama, u := range map[string]struct {
		a     models.AlamatBaru
		pesan string
	}{
		"zip kosong":    {models.AlamatBaru{Address: "UJI"}, "postalCode wajib"},
		"alamat spasi":  {models.AlamatBaru{PostalCode: "1", Address: "  "}, "address wajib"},
		"terlalu lebar": {models.AlamatBaru{PostalCode: "1", Address: "A", CityName: strings.Repeat("U", 4001)}, "cityName paling banyak 4000"},
	} {
		if _, err := Baru(nil).SimpanAlamatBaru(ctx, akun, u.a); !errors.Is(err, ErrMasukanAlamat) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if _, err := Baru(nil).SimpanAlamatBaru(ctx, akun, sah); !errors.Is(err, ErrRWTanpaDatabase) {
		t.Errorf("simpan tanpa DB: %v", err)
	}
	galatUji := errors.New("ORA-UJI")
	if _, err := Baru(nil).DenganRW(rwTiruan{zip: &zip, alamat: &alamat, err: galatUji}).DenganTransaksi(tanpaTx).SimpanAlamatBaru(ctx, akun, sah); !errors.Is(err, galatUji) {
		t.Errorf("galat repository: %v", err)
	}
}
