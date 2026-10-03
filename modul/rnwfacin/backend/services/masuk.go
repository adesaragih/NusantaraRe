// Package services memuat layanan modul RNW Fac In (tiket R01, lingkup layanan
// saja - keputusan work owner 01-10-2026, butir 55): kasus renewal dibuat dari
// polis lama, lalu diserahkan ke tangga akseptasi NB lewat kontrak - tanpa
// salinan mesin dan tanpa cabang khusus renewal.
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/rnwfacin/backend/models"
	"nusantarare/modul/rnwfacin/backend/repository"
)

// Jalur properti yang ditulis pintu masuk renewal. `[terverifikasi]` keduanya
// diikat section `PeriodeRenewal` (`D:\migrasi\RNM\RNW Fac In\Section\PeriodeRenewal.xml`,
// `<pyValue>.QuotationData.OldPolicyNo` / `.QuotationData.StatusBusiness`, kelas
// ASM-FW-GISFW-Data-OfferFacIn); fixture kasus RNW nyata ber-StatusBusiness 2.
const (
	jalurStatusBusiness = "pyWorkPage.OfferFacIn.QuotationData.StatusBusiness"
	jalurOldPolicyNo    = "pyWorkPage.OfferFacIn.QuotationData.OldPolicyNo"
)

// ErrMasukan - masukan layar masuk renewal tidak sah.
var ErrMasukan = errors.New("services: masukan renewal tidak sah")

// Masukan - layar masuk renewal: No. Polis + Renewal Date + Note.
type Masukan struct {
	NomorPolis     string
	TanggalRenewal string // utils.TanggalSaja
	Catatan        string
}

// Service - layanan pintu masuk renewal.
type Service struct {
	polis  repository.PembacaPolisLama
	tangga kontrak.TanggaAkseptasiFacIn
}

// NewService merakit layanan; tangga disediakan nbfacin lewat kontrak.
func NewService(polis repository.PembacaPolisLama, tangga kontrak.TanggaAkseptasiFacIn) *Service {
	return &Service{polis: polis, tangga: tangga}
}

// BuatKasusRenewal - tombol OK. `[keputusan rancangan]` gerbang masuk ditetapkan,
// bukan diport (tiket R01: tidak ada perilaku terekam): kasus baru = salinan data
// polis lama sebagai nilai awal, StatusBusiness 2, OldPolicyNo = nomor polis.
// Penyetelan nilai periode baru milik R02.
func (s *Service) BuatKasusRenewal(ctx context.Context, in Masukan) (models.KasusRenewal, error) {
	nomor := strings.TrimSpace(in.NomorPolis)
	if nomor == "" {
		return models.KasusRenewal{}, fmt.Errorf("%w: nomor polis kosong", ErrMasukan)
	}
	tanggal, err := utils.ParseTanggal(in.TanggalRenewal)
	if err != nil {
		return models.KasusRenewal{}, fmt.Errorf("%w: tanggal renewal: %v", ErrMasukan, err)
	}
	lama, err := s.polis.Baca(ctx, nomor)
	if err != nil {
		return models.KasusRenewal{}, err
	}
	baru := lama.Salin()
	baru[jalurStatusBusiness] = models.StatusBusinessRenewal
	baru[jalurOldPolicyNo] = nomor
	return models.KasusRenewal{Kasus: baru, TanggalRenewal: tanggal, Catatan: in.Catatan}, nil
}

// LangkahAkseptasi - satu langkah tangga NB atas kasus renewal, apa adanya.
func (s *Service) LangkahAkseptasi(k models.KasusRenewal, p kontrak.PenggunaFacIn) (kontrak.TransisiFacIn, error) {
	return s.tangga.Langkah(k.Kasus, p)
}
