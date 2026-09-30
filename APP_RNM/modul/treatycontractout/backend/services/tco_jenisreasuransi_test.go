package services_test

// Uji layanan jenis reasuransi Treaty Contract Out - TANPA Oracle (tiket 02).

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatycontractout/backend/repository"
	"nusantarare/modul/treatycontractout/backend/services"
)

type pembacaUji struct {
	baris []repository.JenisReasuransiTCO
	anak  []repository.JenisReasuransiTCO
	err   error
	// induk - induk yang diminta DaftarAnakTreatyLimit terakhir.
	induk *string
}

func (p pembacaUji) DaftarNonLife(context.Context) ([]repository.JenisReasuransiTCO, error) {
	return p.baris, p.err
}

func (p pembacaUji) DaftarAnakTreatyLimit(_ context.Context, induk string) ([]repository.JenisReasuransiTCO, error) {
	if p.induk != nil {
		*p.induk = induk
	}
	return p.anak, p.err
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

// Anak Treaty Limit [keputusan work owner 30-09-2026]: pilihan dari pembaca
// porsi + induk, induk wajib, kosong = gagal terang.
func TestJenisReasuransiAnakTreatyLimit(t *testing.T) {
	var diminta string
	svc := services.New(nil).JenisReasuransiTreaty().DenganPembaca(pembacaUji{
		anak:  []repository.JenisReasuransiTCO{{ID: "10004", Note: "UJI QS (R/I)", Tipe: "4"}, {ID: "10003", Note: "UJI QS", Tipe: "1"}},
		induk: &diminta,
	})
	d, err := svc.DaftarAnakTreatyLimit(context.Background(), pelakuUjiTCO, " 10003 ")
	if err != nil || len(d) != 2 || d[0].ID != "10004" || d[0].Note != "UJI QS (R/I)" || diminta != "10003" {
		t.Errorf("daftar %+v, induk %q, galat %v", d, diminta, err)
	}
	if _, err := svc.DaftarAnakTreatyLimit(context.Background(), inti.Pelaku{}, "10003"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	if _, err := svc.DaftarAnakTreatyLimit(context.Background(), pelakuUjiTCO, " "); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("tanpa induk: %v", err)
	}
	kosong := services.New(nil).JenisReasuransiTreaty().DenganPembaca(pembacaUji{})
	if _, err := kosong.DaftarAnakTreatyLimit(context.Background(), pelakuUjiTCO, "10003"); !errors.Is(err, services.ErrPilihanAnakTreatyLimitKosong) {
		t.Errorf("kosong: %v", err)
	}
	if _, err := services.New(nil).JenisReasuransiTreaty().DaftarAnakTreatyLimit(context.Background(), pelakuUjiTCO, "10003"); !errors.Is(err, services.ErrPembacaJenisReasuransiBelumDisuntik) {
		t.Errorf("bawaan harus gagal terang: %v", err)
	}
}
