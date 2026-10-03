package services

// Tab Coverage FIRE tahap C2 - tiket 44: ‰ Total Net Rate item dibagi ke coverage.
//
// `[terverifikasi]` `NB FacIn\Activity\CekNetRate_ACT.xml` (kelas Data-Coverage; dipanggil CountPremi_ACT langkah 5 bila
// Param.CallAct == "rate"): langkah 2 menelusuri `.Property.PropertyItemList(Param.IdxPropertyItem).CoverageList` dan
// menandai lima .OLDID (perbandingan `==` harfiah: "FLEXAS", "4.1A CC", "4.3", "4.2 PRGBI", "OTHERS"); kelimanya ada ->
// FlagNetRate "true". Langkah 3 menyalin flag ke item; langkah 4: flag "false" -> item .TotalNetRate = 0 (NetRate
// coverage TIDAK disentuh).
// `[terverifikasi]` `NB FacIn\Activity\CalculateNetRate_ACT.xml` (kelas Data-PropertyItem): langkah 1 TotalNetrate =
// @divide(.TotalNetRate, 1, 20); langkah 2 (flag "true") TotalRate = Σ .Rate CoverageList; langkah 3 (flag "true") per
// coverage .NetRate = @divide(.Rate, TotalRate, 20) * TotalNetrate, lalu Call CountPremi_ACT PremiStatus "percent"
// (DiscountStatus tidak diisi -> langkah diskon 51-52 tidak berjalan).

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

// oldIDNetRate - lima .OLDID CekNetRate_ACT langkah 2 (harfiah, peka huruf, tanpa pemangkasan).
var oldIDNetRate = []string{"FLEXAS", "4.1A CC", "4.3", "4.2 PRGBI", "OTHERS"}

// FlagNetRate - CekNetRate_ACT: true hanya bila kelima oldIDNetRate ada di antara coverage item.
func FlagNetRate(cov []models.CoverageObjek) bool {
	ada := map[string]bool{}
	for _, c := range cov {
		ada[c.OldID] = true
	}
	for _, k := range oldIDNetRate {
		if !ada[k] {
			return false
		}
	}
	return true
}

// TerapkanNetRate - CalculateNetRate_ACT langkah 1-3 untuk item ber-flag: NetRate tiap coverage = Rate / ΣRate (20
// desimal) x TotalNetRate (kosong = 0), lalu CountPremi_ACT mode percent tanpa diskon (ModeDiskonTanpa = DiscountStatus
// kosong). ΣRate = 0 -> ErrMasukanCoverage (A165; Pega membagi nol). Galat = *galatNetRate berjalur relatif item.
// Larik masukan tidak diubah.
func TerapkanNetRate(cov []models.CoverageObjek, tsi, prorate, totalNetRate *apd.Decimal, it PenyesuaianItem) ([]models.CoverageObjek, error) {
	jumlahRate := apd.New(0, 0)
	for _, c := range cov {
		konteksCoverage.Add(jumlahRate, jumlahRate, nolBila(c.Rate))
	}
	if jumlahRate.IsZero() {
		return nil, &galatNetRate{"totalNetRate", "Σ Rate coverage = 0, Total Net Rate tidak dapat dibagi"}
	}
	tnr, err := bagi(nolBila(totalNetRate), apd.New(1, 0), skalaBagi)
	if err != nil {
		return nil, &galatNetRate{"totalNetRate", err.Error()}
	}
	hasil := make([]models.CoverageObjek, len(cov))
	for k, c := range cov {
		porsi, err := bagi(nolBila(c.Rate), jumlahRate, skalaBagi)
		if err != nil {
			return nil, &galatNetRate{fmt.Sprintf("coverages[%d]", k), err.Error()}
		}
		c.NetRate = kali(porsi, tnr)
		if hasil[k], err = HitungCoverage(c, tsi, prorate, ModePercent, ModeDiskonTanpa, it); err != nil {
			return nil, &galatNetRate{fmt.Sprintf("coverages[%d]", k), err.Error()}
		}
	}
	return hasil, nil
}

// galatNetRate - galat TerapkanNetRate: jalur relatif item ("totalNetRate", "coverages[k]") + alasan. Tipe tersendiri
// (bukan fmt.Errorf seperti galat coverage lain) karena satu galat dipakai dua rute dengan jalur berbeda: POST
// hitung-net-rate memakainya apa adanya (dibuka sebagai ErrMasukanCoverage -> 400), PUT objek mengambil `jalur` lewat
// errors.As lalu memberi awalan "baris[n].items[m]." (ErrMasukanObjek -> 400).
type galatNetRate struct{ jalur, alasan string }

func (g *galatNetRate) Error() string {
	return ErrMasukanCoverage.Error() + ": " + g.jalur + ": " + g.alasan
}
func (g *galatNetRate) Unwrap() error { return ErrMasukanCoverage }

// HasilNetRate - jawaban POST hitung-net-rate.
type HasilNetRate struct {
	Coverages    []models.CoverageObjek
	TotalNetRate *apd.Decimal
	Flag         bool
}

// HitungNetRateKasus - POST hitung-net-rate: tanpa menyimpan. Case selalu dibaca (404 / 503 sama dengan hitung-coverage).
// Flag false -> TotalNetRate 0, coverage dikembalikan apa adanya (CekNetRate langkah 4), periode tidak diperiksa; flag
// true -> TerapkanNetRate dengan Prorate periode case.
func (s *Service) HitungNetRateKasus(ctx context.Context, id string, cov []models.CoverageObjek, tsi, totalNetRate *apd.Decimal,
	it PenyesuaianItem) (HasilNetRate, error) {
	if masalah := periksaCoverage("", cov); len(masalah) > 0 {
		return HasilNetRate{}, fmt.Errorf("%w: %s", ErrMasukanCoverage, strings.Join(masalah, "; "))
	}
	if !idKasusSah(id) {
		return HasilNetRate{}, ErrKasusTidakAda
	}
	periode, err := s.periodeKasus(ctx, id)
	if err != nil {
		return HasilNetRate{}, err
	}
	if cov == nil {
		cov = []models.CoverageObjek{}
	}
	if !FlagNetRate(cov) {
		return HasilNetRate{Coverages: cov, TotalNetRate: apd.New(0, 0)}, nil
	}
	prorate, err := Prorate(periode)
	if err != nil {
		return HasilNetRate{}, err
	}
	hasil, err := TerapkanNetRate(cov, tsi, prorate, totalNetRate, it)
	if err != nil {
		return HasilNetRate{}, err
	}
	return HasilNetRate{Coverages: hasil, TotalNetRate: simpan8(nolBila(totalNetRate)), Flag: true}, nil
}
