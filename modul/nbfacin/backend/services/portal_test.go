package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

type portalTiruan struct {
	baris          []models.BarisPortal
	total          int
	err            error
	cari           *string
	offset, ukuran *int
}

func (p portalTiruan) CariPortal(_ context.Context, cari string, offset, ukuran int) ([]models.BarisPortal, int, error) {
	if p.cari != nil {
		*p.cari, *p.offset, *p.ukuran = cari, offset, ukuran
	}
	return p.baris, p.total, p.err
}

// TestCariPortal - tiket 32: ukuran 15 (A95), offset dari halaman, cari diteruskan apa
// adanya, total dari repository; 401 / 400 / 503; galat repository diteruskan.
func TestCariPortal(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var cari string
	var offset, ukuran int
	svc := Baru(nil).DenganPortal(portalTiruan{baris: []models.BarisPortal{{CaseID: "UJI-NB-1"}}, total: 31,
		cari: &cari, offset: &offset, ukuran: &ukuran})
	h, err := svc.CariPortal(ctx, akun, "uji_50%", 3)
	if err != nil || h.Total != 31 || h.Halaman != 3 || h.Ukuran != 15 || len(h.Baris) != 1 {
		t.Fatalf("%+v (%v)", h, err)
	}
	if cari != "uji_50%" || offset != 30 || ukuran != 15 {
		t.Errorf("ke repository: %q %d %d", cari, offset, ukuran)
	}
	if _, err := svc.CariPortal(ctx, inti.Pelaku{}, "", 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for _, u := range []struct {
		cari    string
		halaman int
	}{{"", 0}, {"", -1}, {strings.Repeat("é", 256), 1}, {"", halamanMaksPortal + 1}} {
		if _, err := svc.CariPortal(ctx, akun, u.cari, u.halaman); !errors.Is(err, ErrMasukanPortal) {
			t.Errorf("cari %d rune, halaman %d: %v", len([]rune(u.cari)), u.halaman, err)
		}
	}
	if _, err := svc.CariPortal(ctx, akun, strings.Repeat("é", 255), halamanMaksPortal); err != nil {
		t.Errorf("tepat di batas: %v", err)
	}
	if _, err := Baru(nil).CariPortal(ctx, akun, "", 1); !errors.Is(err, ErrPortalTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	galatUji := errors.New("ORA-UJI")
	if _, err := Baru(nil).DenganPortal(portalTiruan{err: galatUji}).CariPortal(ctx, akun, "", 1); !errors.Is(err, galatUji) {
		t.Errorf("galat repository: %v", err)
	}
}
