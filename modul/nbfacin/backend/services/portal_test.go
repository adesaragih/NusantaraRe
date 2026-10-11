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
	baris         []models.BarisPortal
	total         int
	err           error
	cari, pembuat *string
	// selesai - kueri selesai terakhir yang diteruskan.
	selesai        *string
	offset, ukuran *int
	// milik - bila diisi: baris per pembuat (meniru saringan CREATE_OP; total = cacah milik pembuat).
	milik map[string][]models.BarisPortal
}

func (p portalTiruan) CariPortal(_ context.Context, pembuat, cari, selesai string, offset, ukuran int) ([]models.BarisPortal, int, error) {
	if p.selesai != nil {
		*p.selesai = selesai
	}
	if p.cari != nil {
		*p.cari, *p.offset, *p.ukuran = cari, offset, ukuran
	}
	if p.pembuat != nil {
		*p.pembuat = pembuat
	}
	if p.milik != nil {
		return p.milik[pembuat], len(p.milik[pembuat]), nil
	}
	return p.baris, p.total, p.err
}

// TestCariPortal - tiket 32: ukuran 15 (A95), offset dari halaman, cari diteruskan apa
// adanya, total dari repository; 401 / 400 / 503; galat repository diteruskan.
func TestCariPortal(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var cari, pembuat string
	var offset, ukuran int
	svc := Baru(nil).DenganPortal(portalTiruan{baris: []models.BarisPortal{{CaseID: "UJI-NB-1"}}, total: 31,
		cari: &cari, pembuat: &pembuat, offset: &offset, ukuran: &ukuran})
	h, err := svc.CariPortal(ctx, akun, "uji_50%", "", 3)
	if err != nil || h.Total != 31 || h.Halaman != 3 || h.Ukuran != 15 || len(h.Baris) != 1 {
		t.Fatalf("%+v (%v)", h, err)
	}
	if cari != "uji_50%" || offset != 30 || ukuran != 15 || pembuat != "UJI-USER" {
		t.Errorf("ke repository: %q %d %d %q", cari, offset, ukuran, pembuat)
	}
	// Keputusan work owner 07-10-2026: akun A tidak melihat case akun B; total = milik sendiri.
	milik := Baru(nil).DenganPortal(portalTiruan{milik: map[string][]models.BarisPortal{
		"UJI-A": {{CaseID: "UJI-NB-1"}, {CaseID: "UJI-NB-3", NBStatus: "NB IS IN TEAM LEADER'S INBOX"}}, "UJI-B": {{CaseID: "UJI-NB-2", NBStatus: " "}}}})
	// NB status: NB_STATUS apa adanya; kosong -> "NEW NB" (awal flow InputInwardFacultativeOffer).
	if h, err := milik.CariPortal(ctx, inti.Pelaku{AkunID: "UJI-A"}, "", "", 1); err != nil || h.Total != 2 || len(h.Baris) != 2 ||
		h.Baris[0].CaseID != "UJI-NB-1" || h.Baris[1].CaseID != "UJI-NB-3" || h.Baris[0].NBStatus != "NEW NB" ||
		h.Baris[1].NBStatus != "NB IS IN TEAM LEADER'S INBOX" {
		t.Errorf("akun A: %+v %v", h, err)
	}
	if h, err := milik.CariPortal(ctx, inti.Pelaku{AkunID: "UJI-B"}, "", "", 1); err != nil || h.Total != 1 || h.Baris[0].CaseID != "UJI-NB-2" ||
		h.Baris[0].NBStatus != "NEW NB" {
		t.Errorf("akun B: %+v %v", h, err)
	}
	if _, err := svc.CariPortal(ctx, inti.Pelaku{}, "", "", 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	for _, u := range []struct {
		cari    string
		halaman int
	}{{"", 0}, {"", -1}, {strings.Repeat("é", 256), 1}, {"", halamanMaksPortal + 1}} {
		if _, err := svc.CariPortal(ctx, akun, u.cari, "", u.halaman); !errors.Is(err, ErrMasukanPortal) {
			t.Errorf("cari %d rune, halaman %d: %v", len([]rune(u.cari)), u.halaman, err)
		}
	}
	if _, err := svc.CariPortal(ctx, akun, strings.Repeat("é", 255), "", halamanMaksPortal); err != nil {
		t.Errorf("tepat di batas: %v", err)
	}
	if _, err := Baru(nil).CariPortal(ctx, akun, "", "", 1); !errors.Is(err, ErrPortalTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	galatUji := errors.New("ORA-UJI")
	if _, err := Baru(nil).DenganPortal(portalTiruan{err: galatUji}).CariPortal(ctx, akun, "", "", 1); !errors.Is(err, galatUji) {
		t.Errorf("galat repository: %v", err)
	}
}

// TestCariPortalSelesai - switch In Progress / Resolved halaman depan (permintaan work owner 11-10-2026, pola nbtreatyin): "" / "0" /
// "1" diteruskan ke repository; nilai lain 400 tanpa membaca repository.
func TestCariPortalSelesai(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	selesai := "awal"
	svc := Baru(nil).DenganPortal(portalTiruan{selesai: &selesai})
	for _, s := range []string{"", "0", "1"} {
		if _, err := svc.CariPortal(ctx, akun, "", s, 1); err != nil || selesai != s {
			t.Errorf("%q: %v diteruskan %q", s, err, selesai)
		}
	}
	for _, s := range []string{"2", "ya", " 1", "true"} {
		selesai = "tidak dibaca"
		if _, err := svc.CariPortal(ctx, akun, "", s, 1); !errors.Is(err, ErrMasukanPortal) || selesai != "tidak dibaca" {
			t.Errorf("%q: %v", s, err)
		}
	}
}
