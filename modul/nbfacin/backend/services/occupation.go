package services

// Saran Occupation sub-tab Surrounding Risk - tiket 38.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// minCariOccupation - permintaan sesi 0f (J-2 tiket 38, bukan keputusan work owner): saran setelah
// 2 karakter; kurang = 400. Batas Pega `[dugaan]` tidak ada (korpus tidak menyebut minimum).
// lebarCariOccupation - batas kata cari (pola A73).
const (
	minCariOccupation   = 2
	lebarCariOccupation = 255
)

// ErrMasukanOccupation - cari 1 karakter / terlalu panjang. 400. Kosong = seluruh okupasi FIRE (tiket 40).
var ErrMasukanOccupation = fmt.Errorf("services: cari kosong (semua) atau %d..%d karakter", minCariOccupation, lebarCariOccupation)

// ErrOccupationTanpaDatabase - tabel OCCUPATION tidak terbaca. 503.
var ErrOccupationTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel OCCUPATION tidak terbaca")

// DenganOccupation memasang pembaca OCCUPATION (tiket 38).
func (s *Service) DenganOccupation(o repository.PembacaOccupation) *Service {
	s.occupation = o
	return s
}

// CariOccupation - saran Occupation FIRE yang OldID atau Name-nya mengandung `cari` (A131); `cari` kosong =
// seluruh okupasi FIRE <= 500 (popup Choose Occupation, tiket 40). Tanpa identitas (pola lookup).
func (s *Service) CariOccupation(ctx context.Context, cari string) ([]models.BarisOccupation, error) {
	k := strings.TrimSpace(cari)
	if n := utf8.RuneCountInString(k); (n > 0 && n < minCariOccupation) || n > lebarCariOccupation {
		return nil, ErrMasukanOccupation
	}
	if s.occupation == nil {
		return nil, ErrOccupationTanpaDatabase
	}
	return s.occupation.CariOccupation(ctx, k)
}
