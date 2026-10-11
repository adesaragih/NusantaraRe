package services

// Lookup form objek Aneka (tiket 73 P1): Vessel Name (kapal) dan Type HE (tipe kendaraan). Object HE / Merk / Draft Wording =
// teks bebas (K73-10 / K73-11), tanpa endpoint.

import (
	"context"
	"errors"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// DenganLookupAneka memasang pembaca lookup Aneka.
func (s *Service) DenganLookupAneka(l repository.PembacaLookupAneka) *Service {
	s.lookupAneka = l
	return s
}

// CariKapal - pilihan Vessel Name (selalu larik); kata = teks ketikan Nama Obyek.
func (s *Service) CariKapal(ctx context.Context, kata string) ([]models.BarisKapal, error) {
	if s.lookupAneka == nil {
		return nil, ErrObjekTanpaDatabase
	}
	return s.lookupAneka.CariKapal(ctx, kata)
}

// DaftarTipeKendaraan - pilihan Type HE (selalu larik; Pega tidak menyaring merk / objek).
func (s *Service) DaftarTipeKendaraan(ctx context.Context) ([]models.BarisTipeKendaraan, error) {
	if s.lookupAneka == nil {
		return nil, ErrObjekTanpaDatabase
	}
	return s.lookupAneka.DaftarTipeKendaraan(ctx)
}

// CariOkupasiAneka - popup Choose Occupation Aneka case `id` (selalu larik): saringan BusinessName = Class of Business case
// (`pyWorkPage.Quotation.BusinessName`; = nama Class of Business `[dugaan kuat]`, pola kasusBisnisDari), kata = kotak Search.
func (s *Service) CariOkupasiAneka(ctx context.Context, id, kata string) ([]models.PilihanOkupasiAneka, error) {
	if s.lookupAneka == nil || s.kasus == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	return s.lookupAneka.CariOkupasiAneka(ctx, strings.TrimSpace(k.Opportunity.ClassOfBusiness), kata)
}

// CariCoverageAneka - Coverage Note tab Coverage Aneka case `id` (selalu larik): BizCode = BusinessCode case (BUSINESS.ID baris
// ber-GROUPPANEL, pola lini; tidak terbaca -> tanpa saringan bisnis), kata = teks ketikan.
func (s *Service) CariCoverageAneka(ctx context.Context, id, kata string) ([]models.BarisCoverage, error) {
	if s.lookupAneka == nil || s.kasus == nil {
		return nil, ErrObjekTanpaDatabase
	}
	if !idKasusSah(id) {
		return nil, ErrKasusTidakAda
	}
	if _, err := s.kasus.BacaKasus(ctx, id); errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	} else if err != nil {
		return nil, err
	}
	kb, _ := s.kasusBisnisDari(ctx, id)
	return s.lookupAneka.CariCoverageAneka(ctx, kb.businessCode, kata)
}
