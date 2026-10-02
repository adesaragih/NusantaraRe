package services

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/rnwfacin/backend/models"
	"nusantarare/modul/rnwfacin/backend/repository"
)

// polisDiMemori - repository.PembacaPolisLama rekaan untuk uji.
type polisDiMemori map[string]models.Kasus

func (p polisDiMemori) Baca(_ context.Context, nomor string) (models.Kasus, error) {
	k, ada := p[nomor]
	if !ada {
		return nil, repository.ErrPolisLamaTakDitemukan
	}
	return k, nil
}

// tanggaPerekam - kontrak.TanggaAkseptasiFacIn rekaan: merekam kasus yang diterimanya.
type tanggaPerekam struct{ diterima kontrak.KasusFacIn }

func (t *tanggaPerekam) Langkah(k kontrak.KasusFacIn, _ kontrak.PenggunaFacIn) (kontrak.TransisiFacIn, error) {
	t.diterima = k
	return kontrak.TransisiFacIn{Selesai: true}, nil
}

const nomorRekaan = "UJI-POLIS-1"

func polisLama() polisDiMemori {
	return polisDiMemori{nomorRekaan: models.Kasus{
		"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "1",
		"pyWorkPage.OfferFacIn.QuotationData.TeamGroup":      "1",
		"pyWorkPage.OfferFacIn.TotalTSINusaRe":               "100",
	}}
}

// TestBuatKasusRenewal - kasus baru: salinan data polis lama, StatusBusiness 2,
// OldPolicyNo = nomor polis yang dimasukkan. Data polis lama tidak berubah.
func TestBuatKasusRenewal(t *testing.T) {
	lama := polisLama()
	s := NewService(lama, &tanggaPerekam{})
	k, err := s.BuatKasusRenewal(context.Background(), Masukan{NomorPolis: nomorRekaan, TanggalRenewal: "2026-12-01", Catatan: "rekaan"})
	if err != nil {
		t.Fatal(err)
	}
	for jalur, mau := range map[string]string{
		jalurStatusBusiness: "2",
		jalurOldPolicyNo:    nomorRekaan,
		"pyWorkPage.OfferFacIn.QuotationData.TeamGroup": "1",
		"pyWorkPage.OfferFacIn.TotalTSINusaRe":          "100",
	} {
		if v, _ := k.Kasus.Nilai(jalur); v != mau {
			t.Errorf("%s = %q, mau %q", jalur, v, mau)
		}
	}
	if k.TanggalRenewal.Format("2006-01-02") != "2026-12-01" || k.Catatan != "rekaan" {
		t.Errorf("masukan tidak terbawa: %+v", k)
	}
	if v, _ := lama[nomorRekaan].Nilai(jalurStatusBusiness); v != "1" {
		t.Errorf("polis lama berubah: StatusBusiness %q", v)
	}
}

// TestBuatKasusRenewalDitolak - masukan tidak sah dan polis tak dikenal.
func TestBuatKasusRenewalDitolak(t *testing.T) {
	s := NewService(polisLama(), &tanggaPerekam{})
	for _, u := range []struct {
		nama string
		in   Masukan
		mau  error
	}{
		{"nomor polis kosong", Masukan{TanggalRenewal: "2026-12-01"}, ErrMasukan},
		{"tanggal tak terbaca", Masukan{NomorPolis: nomorRekaan, TanggalRenewal: "01/12/2026"}, ErrMasukan},
		{"polis tak dikenal", Masukan{NomorPolis: "TIDAK-ADA", TanggalRenewal: "2026-12-01"}, repository.ErrPolisLamaTakDitemukan},
	} {
		if _, err := s.BuatKasusRenewal(context.Background(), u.in); !errors.Is(err, u.mau) {
			t.Errorf("%s: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
}

// TestLangkahAkseptasiTanpaPenyesuaian - kasus renewal diserahkan ke tangga apa
// adanya, tanpa cabang khusus renewal.
func TestLangkahAkseptasiTanpaPenyesuaian(t *testing.T) {
	tg := &tanggaPerekam{}
	s := NewService(polisLama(), tg)
	k, err := s.BuatKasusRenewal(context.Background(), Masukan{NomorPolis: nomorRekaan, TanggalRenewal: "2026-12-01"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.LangkahAkseptasi(k, kontrak.PenggunaFacIn{Jabatan: "SENIORUW"}); err != nil {
		t.Fatal(err)
	}
	if v, _ := tg.diterima.Nilai(jalurStatusBusiness); v != "2" {
		t.Fatalf("tangga menerima StatusBusiness %q, mau 2", v)
	}
}
