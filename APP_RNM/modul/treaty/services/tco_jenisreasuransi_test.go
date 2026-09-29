package services_test

// Uji layanan jenis reasuransi Treaty Contract Out - TANPA Oracle (tiket 02).

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti"
	"nusantarare/modul/treaty/repository"
	"nusantarare/modul/treaty/services"
)

type pembacaUji struct {
	baris []repository.JenisReasuransiTCO
	err   error
}

func (p pembacaUji) DaftarNonLife(context.Context) ([]repository.JenisReasuransiTCO, error) {
	return p.baris, p.err
}

var pelakuUjiTCO = inti.Pelaku{AkunID: "UJI-ADMIN"}

func TestJenisReasuransiTanpaIdentitasDitolak(t *testing.T) {
	svc := services.New(nil).JenisReasuransiTreaty().DenganPembaca(pembacaUji{
		baris: []repository.JenisReasuransiTCO{{ID: "10003", Note: "UJI QS", Tipe: "1"}},
	})
	_, err := svc.Daftar(context.Background(), inti.Pelaku{})
	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

func TestJenisReasuransiBawaanGagalTerang(t *testing.T) {
	// ⛔ Tanpa suntikan, layanan TIDAK diam-diam mengembalikan daftar kosong.
	_, err := services.New(nil).JenisReasuransiTreaty().Daftar(context.Background(), pelakuUjiTCO)
	if !errors.Is(err, services.ErrPembacaJenisReasuransiBelumDisuntik) {
		t.Errorf("bawaan harus gagal terang, dapat %v", err)
	}
}

func TestJenisReasuransiMasterKosongAdalahKegagalan(t *testing.T) {
	svc := services.New(nil).JenisReasuransiTreaty().DenganPembaca(pembacaUji{})
	_, err := svc.Daftar(context.Background(), pelakuUjiTCO)
	if !errors.Is(err, services.ErrMasterJenisReasuransiKosong) {
		t.Errorf("master kosong harus ErrMasterJenisReasuransiKosong (ADR-0015), dapat %v", err)
	}
}

func TestJenisReasuransiMeneruskanGalatPembaca(t *testing.T) {
	sebab := errors.New("UJI: oracle putus")
	svc := services.New(nil).JenisReasuransiTreaty().DenganPembaca(pembacaUji{err: sebab})
	if _, err := svc.Daftar(context.Background(), pelakuUjiTCO); !errors.Is(err, sebab) {
		t.Errorf("galat pembaca ditelan: %v", err)
	}
}

func TestJenisReasuransiMembawaTigaMedanApaAdanya(t *testing.T) {
	svc := services.New(nil).JenisReasuransiTreaty().DenganPembaca(pembacaUji{
		baris: []repository.JenisReasuransiTCO{
			{ID: "10003", Note: "UJI QS", Tipe: "1"},
			{ID: "00007", Note: "UJI SURPLUS", Tipe: "2"},
		},
	})
	daftar, err := svc.Daftar(context.Background(), pelakuUjiTCO)
	if err != nil {
		t.Fatal(err)
	}
	if len(daftar) != 2 || daftar[1].ID != "00007" || daftar[0].Note != "UJI QS" || daftar[1].Tipe != "2" {
		t.Errorf("daftar berubah: %+v", daftar)
	}
	// Layanan TIDAK menyaring ulang: saringannya di SQL. Urutan dipertahankan.
	if daftar[0].ID != "10003" {
		t.Error("urutan dari pembaca harus dipertahankan")
	}
}
