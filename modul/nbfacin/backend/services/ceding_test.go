package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// cedingTiruan - penyimpan kasus tiruan untuk tiket 34: menolak kode "UJI-TOLAK" seperti
// repository (ErrCedingTidakSah), menyimpan daftar urut dengan nama sintetis.
type cedingTiruan struct{ kasusTiruan }

func (c cedingTiruan) SimpanGeneral(ctx context.Context, tx *db.Tx, id string, g models.General) error {
	daftar := []models.Ceding{}
	for i, k := range g.CedingIDs {
		if k == "UJI-TOLAK" {
			return fmt.Errorf("%w: cedingIds[%d]", repository.ErrCedingTidakSah, i)
		}
		daftar = append(daftar, models.Ceding{ID: k, Name: "UJI NAMA " + k})
	}
	g.CedingList = daftar
	return c.kasusTiruan.SimpanGeneral(ctx, tx, id, g)
}

// TestSimpanGeneralCeding - urutan pilih terjaga, kosong mengosongkan, kode tak lolos = 400
// (pesan menyebut indeksnya), kode kosong / > 50 bita = 400 sebelum basis data, ganda boleh.
func TestSimpanGeneralCeding(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	svc := Baru(nil).DenganKasus(cedingTiruan{kasusTiruan{ada: kasusUji()}}).DenganTransaksi(tanpaTx)
	k, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{CedingIDs: []string{"UJI-B", "UJI-A", "UJI-B"}})
	if err != nil || len(k.General.CedingList) != 3 || k.General.CedingList[0].ID != "UJI-B" || k.General.CedingList[1].ID != "UJI-A" ||
		k.General.CedingList[2].Name != "UJI NAMA UJI-B" {
		t.Fatalf("%+v (%v)", k.General.CedingList, err)
	}
	if k, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{}); err != nil || len(k.General.CedingList) != 0 {
		t.Errorf("kosong: %+v (%v)", k.General.CedingList, err)
	}
	_, err = svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{CedingIDs: []string{"UJI-A", "UJI-TOLAK"}})
	if !errors.Is(err, ErrMasukanGeneral) || !strings.Contains(err.Error(), "cedingIds[1]") {
		t.Errorf("kode tak lolos: %v", err)
	}
	for _, kode := range []string{"", "  ", strings.Repeat("U", 51)} {
		if _, err := Baru(nil).SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{CedingIDs: []string{"UJI-A", kode}}); !errors.Is(err, ErrMasukanGeneral) ||
			!strings.Contains(err.Error(), "cedingIds[1]") {
			t.Errorf("kode %q: %v, mau 400 sebelum DB", kode, err)
		}
	}
	if _, err := svc.SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{CedingIDs: []string{"UJI-A"}}); err != nil {
		t.Errorf("satu kode: %v", err)
	}
	galat := fmt.Errorf("%w: gabungan", repository.ErrCedingTerlaluPanjang)
	if _, err := Baru(nil).DenganKasus(galatSimpan{kasusTiruan{ada: kasusUji()}, galat}).DenganTransaksi(tanpaTx).
		SimpanGeneral(ctx, akun, "UJI-NB-1", IsianGeneral{}); !errors.Is(err, ErrMasukanGeneral) {
		t.Errorf("terlalu panjang: %v, mau 400", err)
	}
}

type galatSimpan struct {
	kasusTiruan
	err error
}

func (g galatSimpan) SimpanGeneral(context.Context, *db.Tx, string, models.General) error {
	return g.err
}
