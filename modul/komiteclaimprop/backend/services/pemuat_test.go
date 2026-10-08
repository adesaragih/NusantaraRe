package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/services"
)

type sumberUji struct{ pegaGagal error }

func (sumberUji) BacaOSLama(context.Context) ([]models.BarisOSLama, error) {
	t := time.Date(2026, 1, 2, 0, 0, 0, 0, models.Jakarta)
	k := models.AwalanKunciKlaimLama
	return []models.BarisOSLama{
		{CaseID: k + "UJI1", StsReject: "0", Tanggal: t}, {CaseID: k + "UJI1", StsReject: "1", Tanggal: t.Add(time.Hour)},
		{CaseID: k + "UJI2", StsReject: "1", Tanggal: t},
		{CaseID: k + "UJI3", StsReject: "1", Tanggal: t},
		{CaseID: k + "UJI4", StsReject: "1", Tanggal: t}, {CaseID: k + "UJI4", StsReject: "0", Tanggal: t.Add(time.Hour)},
	}, nil
}

func (sumberUji) BacaRiwayatLama(context.Context) ([]models.RiwayatLama, error) {
	k := models.AwalanKunciKlaimLama
	return []models.RiwayatLama{{IDPega: k + "UJI1", IDKomite: "X TKMT-UJI1", Status: "ACCEPT"},
		{IDPega: k + "UJI2", IDKomite: "X TKMT-UJI2", Status: "ACCEPT"},
		{IDPega: k + "UJI2", IDKomite: "X TKMT-UJI2", Status: "ACCEPT"}}, nil
}

func (s sumberUji) BacaKomitePega(context.Context) ([]models.KomitePega, error) {
	if s.pegaGagal != nil {
		return nil, s.pegaGagal
	}
	return []models.KomitePega{{Kunci: "X TKMT-UJI1", Cover: "ASM-FW-GCNMFW-WORK-CLAIMTREATY CLMP-UJI1", StatusWork: "New"}}, nil
}

func TestSensusLamaUjiKering(t *testing.T) {
	s, pg, err := services.SensusLama(context.Background(), sumberUji{})
	if err != nil || pg != nil {
		t.Fatal(err, pg)
	}
	// UJI4 berbaris berlaku 0 (Claim Prop memuatnya); UJI1-UJI3 ditunda.
	if s.KlaimLama != 4 || s.Ditunda != 3 || s.DenganBlob != 1 || s.TanpaTangga != 1 || s.TanpaKomite != 1 {
		t.Fatalf("sensus %+v", s)
	}
	if s.Baris[1].Klaim != "CLMP-UJI2" || s.Baris[1].Riwayat != 2 || s.Baris[1].KomiteRiwayat != 1 ||
		s.Baris[1].Sebab != models.SebabTanpaTangga {
		t.Fatalf("baris UJI2 %+v", s.Baris[1])
	}
	_, pg, err = services.SensusLama(context.Background(), sumberUji{pegaGagal: errors.New("ORA-00942")})
	if err != nil || pg == nil {
		t.Fatal("DATAPEGA tak terbaca dilaporkan, sensus tetap jalan")
	}
}

func TestJalankanLamaDitolak(t *testing.T) {
	if !errors.Is(services.JalankanLama(true), services.ErrPemuatProduksi) {
		t.Fatal("IS_PEGA_PROD=true ditolak")
	}
	if !errors.Is(services.JalankanLama(false), services.ErrPemuatTanpaSumberTangga) {
		t.Fatal("tanpa sumber tangga ditolak (jangan mengarang tangga)")
	}
	if models.UsulLama("true") != "1" || models.UsulLama("false") != "0" || models.UsulLama("") != "0" {
		t.Fatal("usul benar/salah/kosong -> 1/0/0")
	}
}
