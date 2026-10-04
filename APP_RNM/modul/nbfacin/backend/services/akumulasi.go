package services

// Popup Choose Accumulation Code - tiket 46 butir 2-3: pencarian akumulasi dan saran autocomplete. Tanpa identitas
// (pola lookup). Choose (pemeriksaan zip) dan Copy Accumulation di frontend.
//
// Pemilihan jalur `[terverifikasi]` `NB FacIn\Activity\GetDataAccumulation_act.xml`:
//   langkah 5  (City != "" || District != "" || CenterTransStatus != "") && PostalCode == "" -> MapHeight 1 (jalur SQL);
//   langkah 6  (MapHeight 1) berurutan, semuanya ke pyReportContentPage: 6.2 CenterTransStatus -> GetSummaryRiskAccumPolis_Sql,
//              6.3 City -> GetAccumulationProvince_SQL, 6.4 District -> GetAccumulationDistrict_SQL, 6.5 SyariahStatus ->
//              GetAccumulation_SQL (kelas Int-ACCUMULATION - aturan itu TIDAK ADA di korpus; hanya Int-ACCUMULATION_LIFE);
//              6.6 (CenterTransStatus) salin unik menurut CARI10; 6.7 (tanpa CenterTransStatus) salin pyReportContentPage.
//   selain itu -> SearchRiskAccumulation_RD (data page layar).
// Keputusan A172: kecamatan > kota (RDB-List kemudian menimpa halaman yang sama `[dugaan]`); nomor polis bersama kota /
// kecamatan -> kosong (langkah 6.6 di Pega memindai baris kota / kecamatan yang tanpa CARI10 sehingga hanya menghasilkan
// satu baris kosong); syariahStatus diabaikan (aturan 6.5 tidak ada; RD tidak menyaringnya).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// ErrMasukanAkumulasi - parameter popup akumulasi tidak sah (terlalu panjang / jenis saran tak dikenal). 400.
var ErrMasukanAkumulasi = errors.New("services: parameter akumulasi tidak sah")

// ErrSaranBelumTersedia - jenis saran yang sumbernya belum terverifikasi (nation / province / accumtype / czone: DDL
// tidak ada, A178). 501.
var ErrSaranBelumTersedia = errors.New("services: sumber saran ini belum terverifikasi (DDL belum ada)")

// ErrAkumulasiTanpaDatabase - ACCUMULATION / RW / CITY / DISTRICT / JSON_POLIS tidak terbaca. 503.
var ErrAkumulasiTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, akumulasi tidak terbaca")

// lebarParamAkumulasi - batas tiap parameter kueri (pola A73).
const lebarParamAkumulasi = 255

// jenisSaranMenunggu - jenis saran yang diminta kontrak tetapi sumbernya belum ada di DDL (A178).
var jenisSaranMenunggu = map[string]bool{"nation": true, "province": true, "accumtype": true, "czone": true}

// DenganAkumulasi memasang pembaca akumulasi (tiket 46).
func (s *Service) DenganAkumulasi(a repository.PembacaAkumulasi) *Service {
	s.akumulasi = a
	return s
}

// paramAkumulasi - satu parameter kueri: nama (untuk pesan 400) dan nilai yang dipangkas di tempat.
type paramAkumulasi struct {
	nama  string
	nilai *string
}

// periksaParamAkumulasi - pangkas setiap nilai; > 255 karakter -> 400 (pesan urut parameter).
func periksaParamAkumulasi(param ...paramAkumulasi) error {
	var masalah []string
	for _, p := range param {
		*p.nilai = strings.TrimSpace(*p.nilai)
		if utf8.RuneCountInString(*p.nilai) > lebarParamAkumulasi {
			masalah = append(masalah, fmt.Sprintf("%s paling banyak %d karakter", p.nama, lebarParamAkumulasi))
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanAkumulasi, strings.Join(masalah, "; "))
	}
	return nil
}

// CariAkumulasi - GET /api/nbfacin/akumulasi.
func (s *Service) CariAkumulasi(ctx context.Context, f models.SaringAkumulasi) ([]models.BarisAkumulasi, error) {
	if err := periksaParamAkumulasi(paramAkumulasi{"id", &f.ID}, paramAkumulasi{"policyNo", &f.PolicyNo},
		paramAkumulasi{"note", &f.Note}, paramAkumulasi{"postalCode", &f.PostalCode}, paramAkumulasi{"syariahStatus", &f.SyariahStatus},
		paramAkumulasi{"provinceId", &f.ProvinceID}, paramAkumulasi{"cityId", &f.CityID}, paramAkumulasi{"districtId", &f.DistrictID},
		paramAkumulasi{"czone", &f.CZone}, paramAkumulasi{"keyword", &f.Keyword}); err != nil {
		return nil, err
	}
	if s.akumulasi == nil {
		return nil, ErrAkumulasiTanpaDatabase
	}
	jalurSQL := (f.CityID != "" || f.DistrictID != "" || f.PolicyNo != "") && f.PostalCode == ""
	switch {
	case !jalurSQL:
		return s.akumulasi.CariAkumulasiRD(ctx, f)
	case f.PolicyNo != "" && (f.CityID != "" || f.DistrictID != ""):
		return []models.BarisAkumulasi{}, nil
	case f.PolicyNo != "":
		return s.akumulasi.AkumulasiPolis(ctx, f.PolicyNo)
	case f.DistrictID != "":
		return s.akumulasi.AkumulasiWilayah(ctx, true, f.DistrictID)
	default:
		return s.akumulasi.AkumulasiWilayah(ctx, false, f.CityID)
	}
}

// SaranAkumulasi - GET /api/nbfacin/akumulasi/saran/{jenis}?q=&induk=.
func (s *Service) SaranAkumulasi(ctx context.Context, jenis, kata, induk string) ([]models.SaranAkumulasi, error) {
	if err := periksaParamAkumulasi(paramAkumulasi{"q", &kata}, paramAkumulasi{"induk", &induk}); err != nil {
		return nil, err
	}
	if jenisSaranMenunggu[jenis] {
		return nil, fmt.Errorf("%w: %s", ErrSaranBelumTersedia, jenis)
	}
	if !repository.JenisSaranDidukung(jenis) {
		return nil, fmt.Errorf("%w: jenis saran %q tidak dikenal", ErrMasukanAkumulasi, jenis)
	}
	if s.akumulasi == nil {
		return nil, ErrAkumulasiTanpaDatabase
	}
	return s.akumulasi.Saran(ctx, jenis, kata, induk)
}
