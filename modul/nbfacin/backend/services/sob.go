package services

// Pilihan SOB untuk popup Change SOB - tiket 33.

import (
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// UkuranHalamanSOB - `[terverifikasi]` `NB FacIn\Harness\SOB.xml` pyPageSize 20 (L3885).
const UkuranHalamanSOB = 20

// batasCariSOB - A103: pola A73 (kolom AGENT VARCHAR2(1000), cari lebih panjang dari 255 ditolak).
const batasCariSOB = 255

const halamanMaksSOB = (1<<31-1)/UkuranHalamanSOB + 1

// ErrMasukanSOB - halaman bukan bilangan bulat >= 1 atau cari terlalu panjang. 400.
var ErrMasukanSOB = fmt.Errorf("services: halaman harus bilangan bulat >= 1 dan cari paling banyak %d karakter", batasCariSOB)

// ErrSOBTanpaDatabase - tabel AGENT tidak terbaca. 503.
var ErrSOBTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel AGENT tidak terbaca")

// HasilSOB - satu halaman pilihan SOB.
type HasilSOB struct {
	Baris          []models.SOB
	Total, Halaman int
	Ukuran         int
}

// DenganSOB memasang pembaca pilihan SOB (tiket 33).
func (s *Service) DenganSOB(p repository.PembacaSOB) *Service {
	s.sob = p
	return s
}

// CariSOB - halaman ke-`halaman` SOB yang lolos syarat work owner dan (bila cari) ID /
// ClientID / nama-nya "mengandung" `cari`, tidak peka huruf. Tanpa identitas (A104, pola
// lookup tiket 27/28).
func (s *Service) CariSOB(ctx context.Context, cari string, halaman int) (HasilSOB, error) {
	if halaman < 1 || halaman > halamanMaksSOB || utf8.RuneCountInString(cari) > batasCariSOB {
		return HasilSOB{}, ErrMasukanSOB
	}
	if s.sob == nil {
		return HasilSOB{}, ErrSOBTanpaDatabase
	}
	baris, total, err := s.sob.CariSOB(ctx, cari, (halaman-1)*UkuranHalamanSOB, UkuranHalamanSOB)
	if err != nil {
		return HasilSOB{}, err
	}
	return HasilSOB{Baris: baris, Total: total, Halaman: halaman, Ukuran: UkuranHalamanSOB}, nil
}
