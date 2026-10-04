package services

// Popup Choose Class of Construction - tiket 40: pilihan TABLEOFLIMIT untuk BIZCODE case.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrTableOfLimitTidakSiap - BIZCODE case tidak dapat diturunkan (Class of Business kosong, tidak ditemukan, atau ganda
// di BUSINESS). 409.
var ErrTableOfLimitTidakSiap = errors.New("services: kode bisnis case tidak dapat ditentukan untuk Table of Limit")

// ErrMasukanTableOfLimit - category terlalu panjang. 400.
var ErrMasukanTableOfLimit = errors.New("services: category paling banyak 50 byte")

// ErrTableOfLimitTanpaDatabase - BUSINESS / TABLEOFLIMIT / case tidak terbaca. 503.
var ErrTableOfLimitTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, Table of Limit tidak terbaca")

// lebarKategoriTOL - lebar T_TABLEOFLIMIT.CATEGORY milik case (VARCHAR2(50), migrasi 189) - tabel TUJUAN pilihan ini;
// sumbernya POOLDATA.TABLEOFLIMIT.CATEGORY VARCHAR2(4000), tetapi kategori > 50 tidak mungkin tersimpan di case.
const lebarKategoriTOL = 50

// DenganTableOfLimit memasang pembaca BUSINESS / TABLEOFLIMIT (tiket 40).
func (s *Service) DenganTableOfLimit(t repository.PembacaTableOfLimit) *Service {
	s.tableOfLimit = t
	return s
}

// TableOfLimit - pilihan Class of Construction case `id`. BIZCODE = BUSINESS.ID ber-NOTE Class of Business case di
// group business case (keputusan work owner butir 89 - form Opportunity menyimpan NAMA, bukan ID); 0 atau > 1 baris ->
// 409 (A154). Tanpa saringan tahun - sama dengan popup Pega (A161 menggantikan A153); Begin date tidak dibutuhkan.
// `kategori` kosong = semua kategori (filter RD dibuang). Tanpa identitas.
func (s *Service) TableOfLimit(ctx context.Context, id, kategori string) ([]models.BarisTableOfLimit, error) {
	kategori = strings.TrimSpace(kategori)
	if len(kategori) > lebarKategoriTOL {
		return nil, ErrMasukanTableOfLimit
	}
	if s.kasus == nil || s.tableOfLimit == nil {
		return nil, ErrTableOfLimitTanpaDatabase
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
	nama, grup := k.Opportunity.ClassOfBusiness, k.Opportunity.GroupBusinessID
	if strings.TrimSpace(nama) == "" || strings.TrimSpace(grup) == "" {
		return nil, fmt.Errorf("%w: Class of Business / Group Business case belum diisi", ErrTableOfLimitTidakSiap)
	}
	kode, err := s.tableOfLimit.KodeBisnis(ctx, nama, grup)
	if err != nil {
		return nil, err
	}
	switch len(kode) {
	case 0:
		return nil, fmt.Errorf("%w: Class of Business %q tidak ada di BUSINESS untuk group business %q", ErrTableOfLimitTidakSiap, nama, grup)
	case 1:
	default:
		return nil, fmt.Errorf("%w: Class of Business %q ganda di BUSINESS untuk group business %q", ErrTableOfLimitTidakSiap, nama, grup)
	}
	return s.tableOfLimit.DaftarTableOfLimit(ctx, kode[0], kategori)
}
