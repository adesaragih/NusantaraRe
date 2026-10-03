package services

// Tab Coverage FIRE tahap C1 - tiket 43: port `NB FacIn\Activity\CountPremi_ACT.xml` (kelas Data-Coverage) basis 1-4,
// total item (langkah 55-56) dan total per mata uang per objek (`SumTotalTSIPremiGross_Act` cabang FIRE).
//
// `[terverifikasi]` CountPremi_ACT (57 langkah, dibaca utuh lewat pengurai langkah):
//   3-4  Prorate = DiffYear*100; DiffYear = DateTimeDifference(Start, End, Y); TotalDayDiff = hari (Start+DiffYear thn)
//        -> End; TotalDays = hari 1 Jan thn(End) -> 1 Jan thn(End)+1 (penugasan kedua menimpa yang pertama).
//   7    QuotationData.EDMDay terisi -> TotalDays = EDMDay.
//   8    DiffDays (@day End - @day Start) != 0 ATAU DiffMonth != 0 -> Prorate += divide(TotalDayDiff, TotalDays, 20)*100.
//   9    jalur EDM - tidak diport (case NB).
//   10   LossLimit = .LostLimit; SubLimit = .Sublimit; .TSI = TSIObjectItem item.
//   11-16 PctAdjustment dari item (PctAdjust1 / PctAdjustOther); 17 LossLimit kosong/0 -> 100; 18 SubLimit tidak > 0 -> 100.
//   19-50 rumus per basis x mode (percent: Premium dari Rate; amount: Rate dibalik dari Premium); langkah BERURUTAN,
//        yang kemudian menimpa (PctAdjustment != 0 -> pembagi x100; NetRate terisi -> NetRate menggantikan Rate).
//   51-52 diskon (Param.DiscountStatus); 53-54 min/max rate (SetRatePolis = tahap C4) - tidak diport.
//   55-56 TotalGrossPremi item = Σ Premium coverage (basis != 5).
// Kode aksi syarat langkah Pega (2 = lanjut syarat, 3 = lewati langkah, 5 = jalankan langkah) `[dugaan]` dari konsistensi
// logika langkah 3, 8, 9, 14 - tidak tertulis di korpus.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Mode hitung = Param.PremiStatus CountPremi_ACT.
const (
	ModePercent = "percent" // Premium dari Rate
	ModeAmount  = "amount"  // Rate dibalik dari Premium
)

// ErrMasukanCoverage - isian hitung coverage tidak sah (basis, mode, pembagi nol). 400.
var ErrMasukanCoverage = errors.New("services: isian coverage tidak sah")

// ErrPeriodeKasus - Begin / End date case kosong atau tak terurai: Prorate tidak dapat dihitung. 409.
var ErrPeriodeKasus = errors.New("services: periode polis case (Begin / End date) belum diisi - premi coverage tidak dapat dihitung")

// ErrCoverageTanpaDatabase - case / COVERAGE_FACIN / COVERAGE tidak terbaca. 503.
var ErrCoverageTanpaDatabase = errors.New("services: basis data tidak dikonfigurasi, coverage tidak terbaca")

// skalaBagi - @Math.divide(..., 20); skalaTSILiability - @Math.divide(..., 100, 4).
const skalaBagi, skalaTSILiability = 20, 4

// konteksCoverage - presisi lebar, setengah-ke-atas (A155: @Math.divide Pega diduga HALF_UP `[dugaan]`). Tidak diubah
// sesudah dibuat.
var konteksCoverage = func() *apd.Context {
	c := apd.BaseContext.WithPrecision(120)
	c.Rounding = apd.RoundHalfUp
	return c
}()

var (
	seratus    = apd.New(100, 0)
	enamPuluh  = apd.New(60, 0)
	seribu     = apd.New(1000, 0)
	errBagiNol = errors.New("pembagi nol")
	pangkat    = func(n int32) *apd.Decimal { return apd.New(1, n) }
)

// kali - perkalian tepat (presisi 120).
func kali(f ...*apd.Decimal) *apd.Decimal {
	h := apd.New(1, 0)
	for _, x := range f {
		konteksCoverage.Mul(h, h, x)
	}
	return h
}

// bagi - @Math.divide(a, b, skala): hasil bagi dibulatkan setengah-ke-atas pada `skala` desimal; b nol -> errBagiNol.
func bagi(a, b *apd.Decimal, skala int32) (*apd.Decimal, error) {
	if b.IsZero() {
		return nil, errBagiNol
	}
	h := new(apd.Decimal)
	if _, err := konteksCoverage.Quo(h, a, b); err != nil {
		return nil, err
	}
	if _, err := konteksCoverage.Quantize(h, h, -skala); err != nil {
		return nil, err
	}
	return h, nil
}

// simpan8 - nilai tersimpan NUMBER(38,8): setengah-ke-atas 8 desimal (A148/A155), nol ekor dibuang; nil tetap nil.
// Quantize gagal (hasil > 112 digit bulat) -> nilai dikembalikan TANPA dibulatkan, sehingga pemeriksaan `muatNumber38`
// sesudahnya menolaknya (bukan galat yang dibuang diam-diam).
func simpan8(d *apd.Decimal) *apd.Decimal {
	if d == nil {
		return nil
	}
	h := new(apd.Decimal)
	if _, err := konteksCoverage.Quantize(h, d, -8); err != nil {
		return d
	}
	h.Reduce(h)
	return h
}

// samaDesimal - nil sama dengan nil; selain itu dibandingkan nilainya.
func samaDesimal(a, b *apd.Decimal) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return a.Cmp(b) == 0
}

// periksaBasis - Coverage Basis yang dapat dihitung tahap C1; "" = sah, selain itu pesan penolakan (diawali nama medan).
func periksaBasis(basis string) string {
	switch strings.TrimSpace(basis) {
	case "1", "2", "3", "4":
		return ""
	case "5":
		return "coverageBasis: Layering belum didukung (coverageBasis 5, tahap C2-C4)"
	}
	return "coverageBasis harus 1-4"
}

// Periode - periode polis case untuk Prorate (langkah 3-8).
type Periode struct {
	Mulai, Akhir string // teks Pega START_DATE_TIME / END_DATE_TIME
	EDMDay       string // QuotationData.EDMDay (T_QUOTATIONDATA.EDM_DAY)
}

// tanggalJakarta - langkah 3: teks Pega -> tanggal kalender Asia/Jakarta (pukul dinormalkan).
func tanggalJakarta(s string) (time.Time, bool) {
	t, err := time.Parse(BentukWaktuPega, s)
	if err != nil {
		return time.Time{}, false
	}
	w := t.In(WIB)
	return time.Date(w.Year(), w.Month(), w.Day(), 0, 0, 0, 0, time.UTC), true
}

func hari(a, b time.Time) int64 { return int64(b.Sub(a) / (24 * time.Hour)) }

// Prorate - langkah 3-8 CountPremi_ACT. DiffYear = tahun penuh Start -> End `[dugaan]` (DateTimeDifference "Y" dibulatkan
// ke bawah). Begin / End kosong -> ErrPeriodeKasus; EDMDay bukan bilangan / 0 -> ErrPeriodeKasus.
func Prorate(p Periode) (*apd.Decimal, error) {
	mulai, ok1 := tanggalJakarta(p.Mulai)
	akhir, ok2 := tanggalJakarta(p.Akhir)
	if !ok1 || !ok2 {
		return nil, ErrPeriodeKasus
	}
	diffYear := akhir.Year() - mulai.Year()
	if mulai.AddDate(diffYear, 0, 0).After(akhir) {
		diffYear--
	}
	diffDays := akhir.Day() - mulai.Day()
	diffMonth := int(akhir.Month()) - int(mulai.Month())
	totalDayDiff := hari(mulai.AddDate(diffYear, 0, 0), akhir)
	awalThn := time.Date(akhir.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	totalDays := apd.New(hari(awalThn, awalThn.AddDate(1, 0, 0)), 0)
	if s := strings.TrimSpace(p.EDMDay); s != "" {
		d, err := utils.ParseDecimal(s)
		if err != nil || d.IsZero() {
			return nil, fmt.Errorf("%w: Day case %q bukan bilangan hari", ErrPeriodeKasus, s)
		}
		totalDays = d
	}
	prorate := kali(apd.New(int64(diffYear), 0), seratus)
	if diffDays != 0 || diffMonth != 0 {
		porsi, err := bagi(apd.New(totalDayDiff, 0), totalDays, skalaBagi)
		if err != nil {
			return nil, fmt.Errorf("%w: TotalDays nol", ErrPeriodeKasus)
		}
		konteksCoverage.Add(prorate, prorate, kali(porsi, seratus))
	}
	return prorate, nil
}

// PenyesuaianItem - data item yang dibaca langkah 11-16.
type PenyesuaianItem struct {
	IsAdjustable   bool
	PctAdjustOther *apd.Decimal
}

// pctAdjustment - langkah 11-15. item.PctAdjust1 tidak ada di model (tiket 39 tidak memport; A156) = kosong/0, sehingga:
// PctAdjustOther >= 60 -> PctAdjustment = PctAdjustOther (langkah 14, dan 15 bila <= 100), selain itu 0. Langkah 16
// (menimpa item.PctAdjustOther dengan PctAdjust1) tidak berlaku karena PctAdjust1 kosong.
func pctAdjustment(it PenyesuaianItem) *apd.Decimal {
	if it.PctAdjustOther != nil && it.PctAdjustOther.Cmp(enamPuluh) >= 0 {
		return it.PctAdjustOther
	}
	return apd.New(0, 0)
}

// Mode diskon = Param.DiscountStatus (langkah 51-52). ModeDiskonTanpa - diskon tidak dihitung.
const ModeDiskonTanpa = "-"

// HitungCoverage - CountPremi_ACT langkah 10-52 atas satu coverage. `tsi` = TSIObjectItem item, `prorate` hasil Prorate,
// `modeDiskon` "percent" / "amount" (frontend, = Param.DiscountStatus), "" = diturunkan (A157), ModeDiskonTanpa = tidak
// dihitung. Basis 5 -> ErrMasukanCoverage "Layering belum didukung"; amount dengan pembagi nol -> ErrMasukanCoverage.
// Hasil: salinan coverage dengan TSI, TSILiability, ProRatePercent, PctAdjustment, Premium (percent) / Rate & NetRate
// (amount), Discount / DiscountPercentage - dibulatkan simpan8.
func HitungCoverage(c models.CoverageObjek, tsi, prorate *apd.Decimal, mode, modeDiskon string, it PenyesuaianItem) (models.CoverageObjek, error) {
	if alasan := periksaBasis(c.CoverageBasis); alasan != "" {
		return c, fmt.Errorf("%w: %s", ErrMasukanCoverage, alasan)
	}
	basis := strings.TrimSpace(c.CoverageBasis)
	if mode != ModePercent && mode != ModeAmount {
		return c, fmt.Errorf("%w: mode harus percent atau amount", ErrMasukanCoverage)
	}
	// Langkah 10, 17, 18.
	c.TSI = nolBila(tsi)
	tsiCov := c.TSI
	lossLimit := c.LostLimit
	if lossLimit == nil || lossLimit.IsZero() {
		lossLimit = seratus
	}
	subLimit := c.Sublimit
	if subLimit == nil || subLimit.Sign() <= 0 {
		subLimit = seratus
	}
	pro, indem := prorate, nolBila(c.IndemnityPercentage)
	adj := pctAdjustment(it)
	c.PctAdjustment = adj
	adaAdj := !adj.IsZero()
	adaNet := c.NetRate != nil && !c.NetRate.IsZero()

	// Langkah 19-50: faktor First Scale (basis 2) dan pembagi dasar 1e9 (basis 1/3/4) atau 1e11 (basis 2); x100 bila
	// PctAdjustment != 0. TSILiability menurut basis (basis 4: Local.SubLimit, tetapi .Sublimit MENTAH bila PctAdjustment
	// atau NetRate terisi - langkah 45-50, dipertahankan persis).
	faktor := []*apd.Decimal{tsiCov}
	dasar := int32(9)
	var err error
	switch basis {
	case "1":
		c.TSILiability = tsiCov
	case "2":
		if c.TSILiability, err = bagi(kali(tsiCov, nolBila(c.FirstLoss)), seratus, skalaTSILiability); err != nil {
			return c, err
		}
		faktor = append(faktor, nolBila(c.FirstScale))
		dasar = 11
	case "3":
		if c.TSILiability, err = bagi(kali(tsiCov, nolBila(c.EmlPml)), seratus, skalaTSILiability); err != nil {
			return c, err
		}
	case "4":
		sl := subLimit
		if adaAdj || adaNet {
			sl = nolBila(c.Sublimit)
		}
		if c.TSILiability, err = bagi(kali(tsiCov, sl), seratus, skalaTSILiability); err != nil {
			return c, err
		}
	}
	faktor = append(faktor, pro, indem)
	if adaAdj {
		faktor = append(faktor, adj)
		dasar += 2
	}
	faktor = append(faktor, lossLimit)
	if mode == ModePercent {
		r := nolBila(c.Rate)
		if adaNet {
			r = c.NetRate
		}
		premi, err := bagi(kali(append([]*apd.Decimal{r}, faktor...)...), pangkat(dasar), skalaBagi)
		if err != nil {
			return c, err
		}
		c.Premium = premi
	} else {
		penyebut := kali(faktor...)
		rate, err := bagi(kali(nolBila(c.Premium), pangkat(dasar)), penyebut, skalaBagi)
		if errors.Is(err, errBagiNol) {
			return c, fmt.Errorf("%w: rate tidak dapat dihitung dari premi - TSI / Prorate / Indemnity / Loss Limit bernilai 0", ErrMasukanCoverage)
		}
		if err != nil {
			return c, err
		}
		c.Rate = rate
		if adaNet {
			c.NetRate = rate
		}
	}
	// Langkah 51-52. Mode kosong diturunkan (A157): DiscountPercentage terisi -> percent; kalau tidak, Discount terisi ->
	// amount. Premi kosong / 0 pada amount -> DiscountPercentage tidak dihitung (tetap) `[dugaan]` (Pega membagi nol).
	if modeDiskon == "" {
		switch {
		case c.DiscountPercentage != nil:
			modeDiskon = ModePercent
		case c.Discount != nil:
			modeDiskon = ModeAmount
		}
	}
	switch {
	case c.Premium == nil:
	case modeDiskon == ModePercent:
		if c.Discount, err = bagi(kali(c.Premium, nolBila(c.DiscountPercentage)), seratus, skalaBagi); err != nil {
			return c, err
		}
	case modeDiskon == ModeAmount && !c.Premium.IsZero():
		dp, err := bagi(nolBila(c.Discount), c.Premium, skalaBagi)
		if err != nil {
			return c, err
		}
		c.DiscountPercentage = kali(dp, seratus)
	}
	c.ProRatePercent = pro
	for _, h := range []struct {
		medan string
		d     **apd.Decimal
	}{{"tsi", &c.TSI}, {"tsiLiability", &c.TSILiability}, {"proRatePercent", &c.ProRatePercent}, {"premium", &c.Premium},
		{"rate", &c.Rate}, {"netRate", &c.NetRate}, {"discount", &c.Discount}, {"discountPercentage", &c.DiscountPercentage},
		{"pctAdjustment", &c.PctAdjustment}} {
		*h.d = simpan8(*h.d)
		// Hasil wajib muat NUMBER(38,8) - pola desimal isian (30 digit bulat, 8 desimal): 400, bukan ORA-01438.
		if !desimalSah(*h.d) {
			return c, fmt.Errorf("%w: hasil %s melebihi 30 digit bulat", ErrMasukanCoverage, h.medan)
		}
	}
	return c, nil
}

// hitungSimpan - satu coverage saat PUT (Param.PremiStatus tidak tersimpan, jadi mode diturunkan dari isinya):
//   - rate dan premi kosong -> mode percent dengan rate 0 dan tanpa diskon; premi dibiarkan kosong (A158);
//   - rate kosong, premi terisi -> amount (A158);
//   - keduanya terisi -> percent (A158), KECUALI percent tidak menghasilkan kembali premi yang dikirim sedangkan amount
//     menghasilkan kembali rate yang dikirim: pasangan itu lahir dari mode amount (rate = rate dibalik yang dibulatkan
//     8 desimal), jadi premi pengguna dipertahankan (A162). Tanpa ini premi ketikan pengguna bergeser saat Save
//     (TSI 123456789012, premi 1000000 -> 999999.9909972).
func hitungSimpan(c models.CoverageObjek, tsi, prorate *apd.Decimal, it PenyesuaianItem) (models.CoverageObjek, error) {
	switch {
	case c.Rate == nil && c.Premium == nil:
		h, err := HitungCoverage(c, tsi, prorate, ModePercent, ModeDiskonTanpa, it)
		h.Premium = nil
		return h, err
	case c.Rate == nil:
		return HitungCoverage(c, tsi, prorate, ModeAmount, "", it)
	}
	h, err := HitungCoverage(c, tsi, prorate, ModePercent, "", it)
	if err != nil || c.Premium == nil || samaDesimal(h.Premium, simpan8(c.Premium)) {
		return h, err
	}
	if a, errA := HitungCoverage(c, tsi, prorate, ModeAmount, "", it); errA == nil && samaDesimal(a.Rate, simpan8(c.Rate)) {
		return a, nil
	}
	return h, nil
}

// periksaCoverage - isian coverage item ke-m baris ke-n (tanpa basis data): jumlah, lebar teks, desimal sah, basis 1-4.
// Rate tidak dipaksa wajib (P-3).
func periksaCoverage(n, m int, cov []models.CoverageObjek) []string {
	if len(cov) > batasItem {
		return []string{fmt.Sprintf("baris[%d].items[%d].coverages paling banyak %d", n, m, batasItem)}
	}
	var masalah []string
	for k, c := range cov {
		awal := fmt.Sprintf("baris[%d].items[%d].coverages[%d].", n, m, k)
		if alasan := periksaBasis(c.CoverageBasis); alasan != "" {
			masalah = append(masalah, awal+alasan)
		}
		for _, l := range lebarCoverage {
			if len(l.nilai(c)) > l.n {
				masalah = append(masalah, fmt.Sprintf("%s%s paling banyak %d byte", awal, l.nama, l.n))
			}
		}
		for _, u := range desimalCoverage {
			if !desimalSah(u.nilai(c)) {
				masalah = append(masalah, awal+u.nama+pesanDesimal)
			}
		}
	}
	return masalah
}

// lebarMataUangCoverage - T_COVERAGELIST.CURRENCY_CODE VARCHAR2(10) (migrasi 193).
const lebarMataUangCoverage = 10

// lebarCoverage - lebar kolom (BYTE) migrasi 193 medan teks coverage.
var lebarCoverage = []struct {
	nama  string
	nilai func(models.CoverageObjek) string
	n     int
}{
	{"coverage", func(c models.CoverageObjek) string { return c.Coverage }, 50},
	{"oldId", func(c models.CoverageObjek) string { return c.OldID }, 50},
	{"coverageNote", func(c models.CoverageObjek) string { return c.CoverageNote }, 500},
	{"coverageBasis", func(c models.CoverageObjek) string { return c.CoverageBasis }, 50},
	{"day", func(c models.CoverageObjek) string { return c.Day }, 50},
	{"indemnity", func(c models.CoverageObjek) string { return c.Indemnity }, 50},
	{"conditions", func(c models.CoverageObjek) string { return c.Conditions }, 500},
}

// desimalCoverage - medan desimal masukan coverage (medan server tsi / tsiLiability / proRatePercent diabaikan).
var desimalCoverage = []struct {
	nama  string
	nilai func(models.CoverageObjek) *apd.Decimal
}{
	{"rate", func(c models.CoverageObjek) *apd.Decimal { return c.Rate }},
	{"rateOjk", func(c models.CoverageObjek) *apd.Decimal { return c.RateOJK }},
	{"firstLoss", func(c models.CoverageObjek) *apd.Decimal { return c.FirstLoss }},
	{"discountPercentage", func(c models.CoverageObjek) *apd.Decimal { return c.DiscountPercentage }},
	{"netRate", func(c models.CoverageObjek) *apd.Decimal { return c.NetRate }},
	{"limitOfLiability", func(c models.CoverageObjek) *apd.Decimal { return c.LimitOfLiability }},
	{"pctLol", func(c models.CoverageObjek) *apd.Decimal { return c.PctLoL }},
	{"indemnityPercentage", func(c models.CoverageObjek) *apd.Decimal { return c.IndemnityPercentage }},
	{"firstScale", func(c models.CoverageObjek) *apd.Decimal { return c.FirstScale }},
	{"sublimit", func(c models.CoverageObjek) *apd.Decimal { return c.Sublimit }},
	{"lostLimit", func(c models.CoverageObjek) *apd.Decimal { return c.LostLimit }},
	{"emlPml", func(c models.CoverageObjek) *apd.Decimal { return c.EmlPml }},
	{"discount", func(c models.CoverageObjek) *apd.Decimal { return c.Discount }},
	{"premium", func(c models.CoverageObjek) *apd.Decimal { return c.Premium }},
}

// periodeKasus - Begin / End / Day case (ErrKasusTidakAda / ErrCoverageTanpaDatabase).
func (s *Service) periodeKasus(ctx context.Context, id string) (Periode, error) {
	if s.kasus == nil {
		return Periode{}, ErrCoverageTanpaDatabase
	}
	k, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return Periode{}, ErrKasusTidakAda
	}
	if err != nil {
		return Periode{}, err
	}
	return Periode{Mulai: k.General.StartDateTime, Akhir: k.General.EndDateTime, EDMDay: k.General.Day}, nil
}

// HitungCoverageKasus - POST hitung-coverage: Prorate dari periode case, tanpa menyimpan. Tanpa identitas tidak perlu
// (tidak menulis) - tetap ikut kebijakan rute (pola lookup).
func (s *Service) HitungCoverageKasus(ctx context.Context, id string, c models.CoverageObjek, tsi *apd.Decimal, mode,
	modeDiskon string, it PenyesuaianItem) (models.CoverageObjek, error) {
	if modeDiskon != "" && modeDiskon != ModePercent && modeDiskon != ModeAmount {
		return c, fmt.Errorf("%w: modeDiskon harus percent atau amount", ErrMasukanCoverage)
	}
	if !idKasusSah(id) {
		return c, ErrKasusTidakAda
	}
	periode, err := s.periodeKasus(ctx, id)
	if err != nil {
		return c, err
	}
	prorate, err := Prorate(periode)
	if err != nil {
		return c, err
	}
	return HitungCoverage(c, tsi, prorate, mode, modeDiskon, it)
}

// siapkanCoverage - PUT objek: setiap coverage dihitung ulang (modeSimpan), TotalGrossPremi item = Σ Premium tersimpan
// (A159; item tanpa coverage -> nil). Periode case dibaca hanya bila ada coverage.
func (s *Service) siapkanCoverage(ctx context.Context, id string, baris []models.ObjekFire) error {
	var prorate *apd.Decimal
	for n := range baris {
		// Salinan item: larik item pemanggil tidak diubah (larik kosong tetap bukan nil).
		if baris[n].Items != nil {
			baris[n].Items = append(make([]models.ItemObjek, 0, len(baris[n].Items)), baris[n].Items...)
		}
		for m := range baris[n].Items {
			it := &baris[n].Items[m]
			if len(it.Coverages) == 0 {
				it.TotalGrossPremi = nil
				continue
			}
			if prorate == nil {
				periode, err := s.periodeKasus(ctx, id)
				if err != nil {
					return err
				}
				if prorate, err = Prorate(periode); err != nil {
					return err
				}
			}
			total := apd.New(0, 0)
			hasil := make([]models.CoverageObjek, len(it.Coverages))
			for k, c := range it.Coverages {
				c2, err := hitungSimpan(c, it.TSI, prorate, PenyesuaianItem{it.IsAdjustable, it.PctAdjustOther})
				if err != nil {
					return fmt.Errorf("%w: baris[%d].items[%d].coverages[%d]: %v", ErrMasukanObjek, n, m, k, err)
				}
				hasil[k] = c2
				if c2.Premium != nil {
					konteksCoverage.Add(total, total, c2.Premium)
				}
			}
			it.Coverages = hasil
			it.TotalGrossPremi = simpan8(total)
		}
	}
	return nil
}

// TotalPerMataUang - SumTotalTSIPremiGross_Act cabang FIRE atas satu objek: per mata uang item (urut kemunculan),
// TSI = Σ TSIObjectItem, Premium = Σ TotalGrossPremi, Rate = Premium / TSI (20 desimal) x 1000, TSI 0 -> 0; disimpan8.
func TotalPerMataUang(item []models.ItemObjek) []models.TotalCoverage {
	hasil := []models.TotalCoverage{}
	posisi := map[string]int{}
	for _, it := range item {
		i, ada := posisi[it.Currency]
		if !ada {
			i = len(hasil)
			posisi[it.Currency] = i
			hasil = append(hasil, models.TotalCoverage{Currency: it.Currency, TSI: apd.New(0, 0), Premium: apd.New(0, 0)})
		}
		konteksCoverage.Add(hasil[i].TSI, hasil[i].TSI, nolBila(it.TSI))
		konteksCoverage.Add(hasil[i].Premium, hasil[i].Premium, nolBila(it.TotalGrossPremi))
	}
	for i := range hasil {
		t := &hasil[i]
		if t.TSI.IsZero() {
			t.Rate = apd.New(0, 0)
		} else if r, err := bagi(t.Premium, t.TSI, skalaBagi); err == nil {
			t.Rate = simpan8(kali(r, seribu))
		} // bagi hanya gagal bila hasil melampaui presisi 120 digit - Rate dibiarkan kosong (baca-saja, tidak disimpan)
		t.TSI, t.Premium = simpan8(t.TSI), simpan8(t.Premium)
	}
	return hasil
}
