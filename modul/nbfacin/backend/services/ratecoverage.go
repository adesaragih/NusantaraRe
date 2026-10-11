package services

// Rate coverage FIRE saat Coverage / Zone diubah (tiket 52): `Activity\CountPremiRate_ACT.xml`
// (ASM-FW-GISFW-DATA-COVERAGE!COUNTPREMIRATE_ACT), dipanggil `CountPremi_ACT` bila Param.CallAct == "rate". Kesembilan
// langkah `[terverifikasi]` (keputusan work owner 06-10-2026 lanjutan: "ikutin sesuai pega saja", "iya silahkan ditambah"):
//
//	Local.PROVINSI = .Property.Province objek == "DKI JAKARTA" / "BANTEN" / "JAWA BARAT" (harfiah) ? "0.5" : "0.45"
//	Okupasi.CARI2  = .OLDID (perbandingan `==` harfiah, peka huruf)
//	FLEXAS         -> per OccupationList objek: RateFlexas_Sql (RATELOWER menurut OCCUPATIONCODE = .OccupationId);
//	                  temprate menyimpan yang TERBESAR; .Rate = temprate.pxResults(1).CARI1
//	"4.1A CC" / "OTHERS" -> .Rate = 0.1
//	"4.3"          -> .Rate = Local.PROVINSI
//	"4.2 PRGBI"    -> .Rate = RateEarthquake_Sql (ZONEOJK.RATE WHERE DESCRIPTION = .Zone) pxResults(1)
//
// Hasil kosong (tanpa okupasi / tanpa baris / NULL / Zone kosong) -> .Rate KOSONG seperti Pega, lalu premi dihitung (Rate
// kosong -> premi 0). OLDID lain -> Rate tidak diubah (tarif db-link tidak diisi, keputusan butir a).
//
// `[dugaan]` FLEXAS: perbandingan `Param.tes > .CARI1` dibaca NUMERIK (RATELOWER NUMBER(20,6); tipe properti CARI1 tidak
// ada di korpus - bila Pega membandingkan teks, hasil bisa beda); okupasi tanpa baris / RATELOWER NULL dilewati (Pega
// menambahkan nilai kosong bila okupasi PERTAMA tanpa hasil - perilaku pembandingan kosong tidak terverifikasi).

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

// CallActRate - Param.CallAct "rate": hitung-coverage menjalankan CountPremiRate_ACT lebih dulu.
const CallActRate = "rate"

// OLDID cabang CountPremiRate_ACT.
const (
	oldIDFlexas = "FLEXAS"
	oldIDGempa  = "4.2 PRGBI"
	oldID43     = "4.3"
)

// Batas masukan rate: province = lebar parameter kueri; occupationIds paling banyak 100 kode, tiap kode selebar
// T_OCCUPATIONLIST.OCCUPATION_ID (VARCHAR2(1000), migrasi 189).
const (
	lebarProvinsiRate    = 255
	batasOkupasiRate     = 100
	lebarKodeOkupasiRate = 1000
)

var (
	rateSepersepuluh = apd.New(1, -1)  // "4.1A CC" / "OTHERS"
	rateProvinsiAtas = apd.New(5, -1)  // 4.3 DKI JAKARTA / BANTEN / JAWA BARAT
	rateProvinsiLain = apd.New(45, -2) // 4.3 provinsi lain
	// oldIDRateTetap - cabang bertarif tetap 0.1.
	oldIDRateTetap = map[string]bool{"4.1A CC": true, "OTHERS": true}
	// provinsiRateAtas - Local.PROVINSI "0.5" (harfiah, peka huruf).
	provinsiRateAtas = map[string]bool{"DKI JAKARTA": true, "BANTEN": true, "JAWA BARAT": true}
)

// PermintaanRate - Param.CallAct dan isian objek saat itu yang dibaca CountPremiRate_ACT (boleh belum tersimpan):
// Province = `.Property.Province`, OccupationIDs = `.OccupationId` tiap baris OccupationList.
type PermintaanRate struct {
	CallAct       string
	Province      string
	OccupationIDs []string
}

// periksa - callAct "rate" atau kosong; province / occupationIds berisi hanya bersama callAct; batas panjang / jumlah.
func (p PermintaanRate) periksa() error {
	if p.CallAct != "" && p.CallAct != CallActRate {
		return fmt.Errorf("%w: callAct harus rate atau absen", ErrMasukanCoverage)
	}
	if p.CallAct == "" && (p.Province != "" || len(p.OccupationIDs) > 0) {
		return fmt.Errorf("%w: province / occupationIds hanya dikirim bersama callAct rate", ErrMasukanCoverage)
	}
	if utf8.RuneCountInString(p.Province) > lebarProvinsiRate {
		return fmt.Errorf("%w: province paling banyak %d karakter", ErrMasukanCoverage, lebarProvinsiRate)
	}
	if len(p.OccupationIDs) > batasOkupasiRate {
		return fmt.Errorf("%w: occupationIds paling banyak %d", ErrMasukanCoverage, batasOkupasiRate)
	}
	for i, k := range p.OccupationIDs {
		if len(k) > lebarKodeOkupasiRate {
			return fmt.Errorf("%w: occupationIds[%d] paling banyak %d byte", ErrMasukanCoverage, i, lebarKodeOkupasiRate)
		}
	}
	return nil
}

// terapkanRate - CountPremiRate_ACT atas satu coverage (dipanggil hanya bila CallAct "rate").
func (s *Service) terapkanRate(ctx context.Context, c models.CoverageObjek, p PermintaanRate) (models.CoverageObjek, error) {
	var err error
	switch {
	case c.OldID == oldIDFlexas:
		c.Rate, err = s.rateFlexas(ctx, p.OccupationIDs)
	case oldIDRateTetap[c.OldID]:
		c.Rate = salinDesimal(rateSepersepuluh)
	case c.OldID == oldID43:
		c.Rate = salinDesimal(rateProvinsiLain)
		if provinsiRateAtas[p.Province] {
			c.Rate = salinDesimal(rateProvinsiAtas)
		}
	case c.OldID == oldIDGempa:
		c.Rate, err = s.rateZona(ctx, c.Zone)
	}
	return c, err
}

// rateFlexas - RATELOWER terbesar atas okupasi objek (kode kosong / ganda dibuang); tanpa nilai -> nil.
func (s *Service) rateFlexas(ctx context.Context, kode []string) (*apd.Decimal, error) {
	var unik []string
	lihat := map[string]bool{}
	for _, k := range kode {
		if k != "" && !lihat[k] {
			lihat[k] = true
			unik = append(unik, k)
		}
	}
	if len(unik) == 0 {
		return nil, nil
	}
	if s.occupation == nil {
		return nil, ErrOccupationTanpaDatabase
	}
	rate, err := s.occupation.RateLowerOkupasi(ctx, unik)
	if err != nil {
		return nil, err
	}
	var maks *apd.Decimal
	for _, k := range unik {
		if r := rate[k]; r != nil && (maks == nil || r.Cmp(maks) > 0) {
			maks = r
		}
	}
	return maks, nil
}

// rateZona - RateEarthquake_Sql; Zone kosong -> nil tanpa membaca basis data (DESCRIPTION = NULL tidak cocok baris mana pun).
func (s *Service) rateZona(ctx context.Context, zona string) (*apd.Decimal, error) {
	if zona == "" {
		return nil, nil
	}
	if s.zona == nil {
		return nil, ErrZonaTanpaDatabase
	}
	return s.zona.RateZona(ctx, zona)
}

// salinDesimal - salinan lepas (nilai paket tidak ikut terubah pemakai).
func salinDesimal(d *apd.Decimal) *apd.Decimal { return new(apd.Decimal).Set(d) }
