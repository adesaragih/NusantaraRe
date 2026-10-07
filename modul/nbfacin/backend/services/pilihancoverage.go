package services

// Pilihan coverage tab Coverage FIRE - tiket 43: popup Choose Coverage (tabel COVERAGE, butir 96) dan lima coverage otomatis
// (COVERAGE, AddCoverageAutoFire).

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// lebarCariCoverage - batas kata cari (pola A73).
const lebarCariCoverage = 255

// DenganCoverage memasang pembaca COVERAGE (tiket 43).
func (s *Service) DenganCoverage(c repository.PembacaCoverage) *Service {
	s.coverage = c
	return s
}

// CariCoverage - popup Choose Coverage: NamaCoverage mengandung `cari` (tidak peka huruf); kosong = semua coverage FIRE.
// Tanpa identitas (pola lookup).
func (s *Service) CariCoverage(ctx context.Context, cari string) ([]models.BarisCoverage, error) {
	k := strings.TrimSpace(cari)
	if utf8.RuneCountInString(k) > lebarCariCoverage {
		return nil, fmt.Errorf("%w: cari paling banyak %d karakter", ErrMasukanCoverage, lebarCariCoverage)
	}
	if s.coverage == nil {
		return nil, ErrCoverageTanpaDatabase
	}
	return s.coverage.CariCoverage(ctx, k)
}

// CoverageOtomatis - lima coverage AddCoverageAutoFire (urut korpus). Tanpa identitas.
func (s *Service) CoverageOtomatis(ctx context.Context) ([]models.BarisCoverage, error) {
	if s.coverage == nil {
		return nil, ErrCoverageTanpaDatabase
	}
	return s.coverage.CoverageOtomatis(ctx)
}
