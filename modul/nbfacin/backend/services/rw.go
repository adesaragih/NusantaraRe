package services

// Popup Add alamat risiko - tiket 37: saran Zip Code dari RW, simpan alamat baru ke RISKADDRESS.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// minZipRW - I-2: saran setelah 3 karakter; kurang = 400. lebarRW - batas kata cari (pola A73).
const (
	minZipRW    = 3
	lebarCariRW = 255
	// lebarRiskAddress - kolom RISKADDRESS VARCHAR2(4000) (DDL), dihitung bita.
	lebarRiskAddress = 4000
)

// ErrMasukanRW - zipCode kosong / < 3 karakter / terlalu panjang. 400.
var ErrMasukanRW = fmt.Errorf("services: zipCode wajib %d..%d karakter", minZipRW, lebarCariRW)

// ErrMasukanAlamat - isian alamat baru tidak sah (postalCode / address wajib, lebar kolom). 400.
var ErrMasukanAlamat = errors.New("services: isian alamat risiko tidak sah")

// ErrRWTanpaDatabase - tabel RW / RISKADDRESS tidak terbaca. 503.
var ErrRWTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, tabel RW / RISKADDRESS tidak terbaca")

// DenganRW memasang pembaca RW dan penulis RISKADDRESS (tiket 37).
func (s *Service) DenganRW(r repository.PenyimpanRW) *Service {
	s.rw = r
	return s
}

// CariRW - saran Zip Code yang DIAWALI `zip` (A123). Tanpa identitas (pola lookup).
func (s *Service) CariRW(ctx context.Context, zip string) ([]models.BarisRW, error) {
	z := strings.TrimSpace(zip)
	if n := utf8.RuneCountInString(z); n < minZipRW || n > lebarCariRW {
		return nil, ErrMasukanRW
	}
	if s.rw == nil {
		return nil, ErrRWTanpaDatabase
	}
	return s.rw.CariRW(ctx, z)
}

// periksaAlamat - I-4: postalCode dan address wajib (bukan hanya spasi); tiap medan <= 4000
// bita. Title tidak dicocokkan ke daftar.
func periksaAlamat(a models.AlamatBaru) error {
	var masalah []string
	if strings.TrimSpace(a.PostalCode) == "" {
		masalah = append(masalah, "postalCode wajib diisi")
	}
	if strings.TrimSpace(a.Address) == "" {
		masalah = append(masalah, "address wajib diisi")
	}
	for _, m := range []struct{ nama, nilai string }{{"nationName", a.NationName}, {"provinceName", a.ProvinceName},
		{"districtName", a.DistrictName}, {"cityName", a.CityName}, {"territoryName", a.TerritoryName}, {"title", a.Title},
		{"address", a.Address}, {"postalCode", a.PostalCode}} {
		if len(m.nilai) > lebarRiskAddress {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d byte", m.nama, lebarRiskAddress))
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanAlamat, strings.Join(masalah, "; "))
	}
	return nil
}

// SimpanAlamatBaru - simpan alamat baru ke master RISKADDRESS (butir 81: port Go prosedur
// InsertUpdateRISKADDRESS); mengembalikan ID baru (W-3). Identitas wajib (data master) ->
// isian -> basis data. Akumulasi (SaveAccumulationByRiskAddress_Act) ditunda (W-2).
func (s *Service) SimpanAlamatBaru(ctx context.Context, pelaku inti.Pelaku, a models.AlamatBaru) (string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	if err := periksaAlamat(a); err != nil {
		return "", err
	}
	if s.rw == nil || s.transaksi == nil {
		return "", ErrRWTanpaDatabase
	}
	var id string
	err := s.transaksi(ctx, func(tx *db.Tx) error {
		var err error
		id, err = s.rw.SisipRiskAddress(ctx, tx, a)
		return err
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
