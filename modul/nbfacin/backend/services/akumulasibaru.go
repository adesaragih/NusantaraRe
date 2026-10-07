package services

// Form Add New akumulasi popup Choose Accumulation (tiket 46; work owner 04-10-2026 "tombol add accumulationnya mana?").
//
// `[terverifikasi]` `NB FacIn\Activity\SaveAccumulation_Act.xml`: Note dan CZone wajib berisi
// (`@PropertyHasValue`), PostalCode dan negara (`InputData.CARI34` = NationInitial) tidak kosong; selain itu langkah 10
// Property-Set-Messages dengan pesan verbatim "Postal code, Nation, CZone and Accumulation Description can't be
// null!". Lolos -> RDB-List UpdateMasterAccumulation = prosedur RDBMASTERACCUMULATION (DDL): cek ganda upper(Note) + zip
// ("Error master item Akumulasi Sudah Ada"), ID negara-zip-lpad(seq,6), Note huruf besar - port-nya sudah ada di mesin
// master inti (master "accumulation", ADR-0043 tanpa CALL), jadi simpan DIDELEGASIKAN ke sana: satu port, satu jejak ubah
// (CREATE_OP = pelaku, 882).
//
// Wajib section `Section\InputAccumulationCov.xml` (Accumulation Type*, Province*, Key Word*, Scope Area*) diperiksa
// di sini juga (keputusan A183: layar menandainya wajib; tanpa pesan verbatim di korpus -> "<medan> wajib diisi").
// Tombol "Add" tipe akumulasi (ModalAccumulation_FacIn) TIDAK dibangun - diganti menu Master Accumulated Type.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	inti "nusantarare/inti/backend"
	masterservices "nusantarare/inti/backend/master/services"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

var (
	// ErrWajibAkumulasi - pesan verbatim SaveAccumulation_Act langkah 10. 400.
	ErrWajibAkumulasi = errors.New("Postal code, Nation, CZone and Accumulation Description can't be null!")
	// ErrAkumulasiSudahAda - pesan verbatim prosedur RDBMASTERACCUMULATION (tanpa markup HTML-nya). 409.
	ErrAkumulasiSudahAda = errors.New("Error master item Akumulasi Sudah Ada")
)

const (
	// masterAkumulasi - kunci master ACCUMULATION di mesin inti.
	masterAkumulasi = "accumulation"
	// UkuranHalamanZip - baris per halaman Choose Zip Code (work owner 04-10-2026 "15 list perpage").
	UkuranHalamanZip = 15
)

// OpsiScopeArea - daftar pilihan `.ScopeArea` `[terverifikasi]` `DDL\ScopeArea.xml`
// (ASM-FW-GISFW-INT-ACCUMULATION!SCOPEAREA, pyStandardValue; dikirim work owner 04-10-2026). Keputusan A186: nilai di
// luar daftar ditolak 400 (layar Pega berupa daftar pilihan).
var OpsiScopeArea = []string{"AREA", "DISTRICT", "CITY", "PROVINCE", "COUNTRY"}

// PenambahMaster - penambah baris master (mesin `inti/backend/master/services.Service`).
type PenambahMaster interface {
	Tambah(ctx context.Context, pelaku inti.Pelaku, kunci string, masukan map[string]string) (string, error)
}

// IsianAkumulasi - badan POST /api/nbfacin/akumulasi.
type IsianAkumulasi struct {
	Accumulation, AccumulationType, Note, Keyword, ScopeArea, CZone, CZoneID, ProvinceID, ZipCode, Negara string
}

// DenganAddAkumulasi memasang pembaca form Add New dan penambah master-nya.
func (s *Service) DenganAddAkumulasi(p repository.PembacaAddAkumulasi, m PenambahMaster) *Service {
	s.addAkumulasi, s.penambahMaster = p, m
	return s
}

// TambahAkumulasi - POST /api/nbfacin/akumulasi -> ID baru dan Note yang tersimpan (huruf besar, seperti prosedur).
func (s *Service) TambahAkumulasi(ctx context.Context, pelaku inti.Pelaku, a IsianAkumulasi) (string, string, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return "", "", err
	}
	for _, v := range []*string{&a.Accumulation, &a.AccumulationType, &a.Note, &a.Keyword, &a.ScopeArea, &a.CZone,
		&a.CZoneID, &a.ProvinceID, &a.ZipCode, &a.Negara} {
		*v = strings.TrimSpace(*v)
	}
	if a.ZipCode == "" || a.Negara == "" || a.CZone == "" || a.Note == "" {
		return "", "", ErrWajibAkumulasi
	}
	var kosong []string
	for _, w := range []struct{ nama, nilai string }{{"accumulation", a.Accumulation}, {"provinceId", a.ProvinceID},
		{"keyword", a.Keyword}, {"scopeArea", a.ScopeArea}} {
		if w.nilai == "" {
			kosong = append(kosong, w.nama+" wajib diisi")
		}
	}
	if len(kosong) > 0 {
		return "", "", fmt.Errorf("%w: %s", ErrMasukanAkumulasi, strings.Join(kosong, "; "))
	}
	if !slices.Contains(OpsiScopeArea, a.ScopeArea) {
		return "", "", fmt.Errorf("%w: scopeArea harus salah satu dari %s", ErrMasukanAkumulasi, strings.Join(OpsiScopeArea, ", "))
	}
	if s.penambahMaster == nil {
		return "", "", ErrAkumulasiTanpaDatabase
	}
	id, err := s.penambahMaster.Tambah(ctx, pelaku, masterAkumulasi, map[string]string{
		"accumulation": a.Accumulation, "accumulationType": a.AccumulationType, "note": a.Note, "keyword": a.Keyword,
		"scopeArea": a.ScopeArea, "cZone": a.CZone, "cZoneId": a.CZoneID, "provinceId": a.ProvinceID, "zipCode": a.ZipCode,
		masterservices.KunciNegara: a.Negara,
	})
	switch {
	case errors.Is(err, masterservices.ErrSudahAda):
		return "", "", ErrAkumulasiSudahAda
	case errors.Is(err, masterservices.ErrMasukanMaster):
		return "", "", fmt.Errorf("%w: %s", ErrMasukanAkumulasi, strings.TrimPrefix(err.Error(), masterservices.ErrMasukanMaster.Error()+": "))
	case errors.Is(err, masterservices.ErrMasterTanpaDatabase):
		return "", "", ErrAkumulasiTanpaDatabase
	case err != nil:
		return "", "", err
	}
	return id, strings.ToUpper(a.Note), nil
}

// CZoneZip - GET /api/nbfacin/akumulasi/czone?zip= : CZone pertama zip itu; zip kosong / tanpa CZone = kosong.
func (s *Service) CZoneZip(ctx context.Context, zip string) (models.CZoneZip, error) {
	if err := periksaParamAkumulasi(paramAkumulasi{"zip", &zip}); err != nil {
		return models.CZoneZip{}, err
	}
	if zip == "" {
		return models.CZoneZip{}, nil
	}
	if s.addAkumulasi == nil {
		return models.CZoneZip{}, ErrAkumulasiTanpaDatabase
	}
	c, _, err := s.addAkumulasi.CZoneZip(ctx, zip)
	return c, err
}

// HalamanZip - satu halaman Choose Zip Code.
type HalamanZip struct {
	Baris                []models.ZipAkumulasi
	Total, Nomor, Ukuran int
}

// CariZipAkumulasi - GET /api/nbfacin/akumulasi/zipcode?provinceName=&q=&halaman= (halaman >= 1; di luar jangkauan =
// baris kosong).
func (s *Service) CariZipAkumulasi(ctx context.Context, provinsi, kata string, nomor int) (HalamanZip, error) {
	if err := periksaParamAkumulasi(paramAkumulasi{"provinceName", &provinsi}, paramAkumulasi{"q", &kata}); err != nil {
		return HalamanZip{}, err
	}
	if nomor < 1 {
		return HalamanZip{}, fmt.Errorf("%w: halaman mulai dari 1", ErrMasukanAkumulasi)
	}
	if s.addAkumulasi == nil {
		return HalamanZip{}, ErrAkumulasiTanpaDatabase
	}
	baris, total, err := s.addAkumulasi.CariZip(ctx, provinsi, kata, nomor, UkuranHalamanZip)
	if err != nil {
		return HalamanZip{}, err
	}
	return HalamanZip{Baris: baris, Total: total, Nomor: nomor, Ukuran: UkuranHalamanZip}, nil
}
