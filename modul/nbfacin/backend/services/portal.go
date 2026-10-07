package services

// Daftar case NB di portal Opportunity - tiket 32 (permintaan work owner 03-10-2026 lewat
// sesi 0f: "nb yang sudah di create, muncul disini ... harus ada case id nya, group
// business, insured name, status").

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// UkuranHalamanPortal - A95: 15 baris, pola keputusan work owner tiket 27 (butir 73.2).
const UkuranHalamanPortal = 15

// batasCariPortal - A96: kolom tercari terlebar BUSINESS_PROSPECT_NAME VARCHAR2(255).
const batasCariPortal = 255

// halamanMaksPortal - offset (halaman-1)*ukuran tetap dalam int32 (bind Oracle).
const halamanMaksPortal = (1<<31-1)/UkuranHalamanPortal + 1

// ErrMasukanPortal - halaman bukan bilangan bulat >= 1 atau cari terlalu panjang. 400.
var ErrMasukanPortal = fmt.Errorf("services: halaman harus bilangan bulat >= 1 dan cari paling banyak %d karakter", batasCariPortal)

// ErrPortalTanpaDatabase - layanan tanpa basis data. 503.
var ErrPortalTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, daftar case NB tidak terbaca")

// HasilPortal - satu halaman daftar case NB.
type HasilPortal struct {
	Baris          []models.BarisPortal
	Total, Halaman int
	Ukuran         int
}

// DenganPortal memasang pembaca daftar portal (tiket 32).
func (s *Service) DenganPortal(p repository.PembacaPortal) *Service {
	s.portal = p
	return s
}

// CariPortal - halaman ke-`halaman` (mulai 1) case NB yang ID atau Business Prospect
// Name-nya "mengandung" `cari` (tidak peka huruf); cari kosong = semua. Urutan
// pemeriksaan: identitas (401) -> masukan (400) -> basis data (503).
func (s *Service) CariPortal(ctx context.Context, pelaku inti.Pelaku, cari string, halaman int) (HasilPortal, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return HasilPortal{}, err
	}
	if halaman < 1 || halaman > halamanMaksPortal || utf8.RuneCountInString(cari) > batasCariPortal {
		return HasilPortal{}, ErrMasukanPortal
	}
	if s.portal == nil {
		return HasilPortal{}, ErrPortalTanpaDatabase
	}
	baris, total, err := s.portal.CariPortal(ctx, cari, (halaman-1)*UkuranHalamanPortal, UkuranHalamanPortal)
	if err != nil {
		return HasilPortal{}, err
	}
	return HasilPortal{Baris: baris, Total: total, Halaman: halaman, Ukuran: UkuranHalamanPortal}, nil
}
