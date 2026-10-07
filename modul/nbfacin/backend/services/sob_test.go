package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

type sobTiruan struct {
	baris          []models.SOB
	total          int
	err            error
	cari           *string
	offset, ukuran *int
}

func (s sobTiruan) CariSOB(_ context.Context, cari string, offset, ukuran int) ([]models.SOB, int, error) {
	if s.cari != nil {
		*s.cari, *s.offset, *s.ukuran = cari, offset, ukuran
	}
	return s.baris, s.total, s.err
}

// TestCariSOB - tiket 33: ukuran 20 (SOB.xml pyPageSize), offset dari halaman, cari apa
// adanya ke repository; 400 / 503.
func TestCariSOB(t *testing.T) {
	var cari string
	var offset, ukuran int
	svc := Baru(nil).DenganSOB(sobTiruan{baris: []models.SOB{{ID: "UJI-G1"}}, total: 41, cari: &cari, offset: &offset, ukuran: &ukuran})
	h, err := svc.CariSOB(context.Background(), "uji", 3)
	if err != nil || h.Ukuran != 20 || h.Total != 41 || offset != 40 || ukuran != 20 || cari != "uji" {
		t.Fatalf("%+v (%v) offset %d", h, err, offset)
	}
	for _, u := range []struct {
		cari    string
		halaman int
	}{{"", 0}, {strings.Repeat("U", 256), 1}} {
		if _, err := svc.CariSOB(context.Background(), u.cari, u.halaman); !errors.Is(err, ErrMasukanSOB) {
			t.Errorf("%d/%d: %v", len(u.cari), u.halaman, err)
		}
	}
	if _, err := Baru(nil).CariSOB(context.Background(), "", 1); !errors.Is(err, ErrSOBTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
}

// generalSOBTiruan - menolak kode SOB tertentu seperti repository (ErrSOBTidakSah).
type generalSOBTiruan struct {
	kasusTiruan
	ditolak string
}

func (g generalSOBTiruan) SimpanGeneral(ctx context.Context, tx *db.Tx, id string, gen models.General) error {
	if gen.SourceOfBusinessID == g.ditolak {
		return repository.ErrSOBTidakSah
	}
	return g.kasusTiruan.SimpanGeneral(ctx, tx, id, gen)
}

// TestSimpanGeneralSOB - kode SOB diteruskan ke repository; kode yang tidak lolos syarat
// AGENT = 400 (ErrMasukanGeneral, E-4); kode > 50 bita = 400 sebelum basis data.
func TestSimpanGeneralSOB(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	svc := Baru(nil).DenganKasus(generalSOBTiruan{kasusTiruan{ada: kasusUji()}, "UJI-TOLAK"}).DenganTransaksi(tanpaTx)
	k, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{SourceOfBusinessID: "UJI-G1"})
	if err != nil || k.General.SourceOfBusinessID != "UJI-G1" {
		t.Fatalf("%+v (%v)", k.General, err)
	}
	if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{SourceOfBusinessID: "UJI-TOLAK"}); !errors.Is(err, ErrMasukanGeneral) ||
		!strings.Contains(err.Error(), "sourceOfBusinessId") {
		t.Errorf("kode tak lolos: %v, mau 400", err)
	}
	if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{SourceOfBusinessID: strings.Repeat("U", 51)}); !errors.Is(err, ErrMasukanGeneral) {
		t.Errorf("kode 51 bita: %v", err)
	}
}
