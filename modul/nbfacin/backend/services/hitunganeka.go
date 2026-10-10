package services

// Hitung coverage Aneka (tiket 73 P2) - rantai on-change sel form `NB FacIn\Section\InputCoverageAneka_FacIn.xml`, tanpa
// menyimpan. Bukti `[terverifikasi]` (riset P2, butir kunci dicek ulang sesi c3; keputusan work owner 09-10-2026 K73-12…K73-15):
//   - CountPremiCoverageAneka (CPCA, `NB FacIn\Activity\CountPremiCoverageAneka.xml`): langkah 4 (L964) ProRatePercent CASE =
//     (End - Start) / Day x 100 (skala 20) HANYA bila Status = "Coverage" (dan IsB2B != "KBRU" - dianggap bukan KBRU, OQ;
//     Edit = "EDIT" tidak pernah dikirim); langkah 6 Prorate = nilai case; 11-12 LossLimit (0 / kosong = 100); 14.1 Prorate
//     0 / kosong (metode 1) dihitung dari tahun polis + pecahan hari (TotalDays = .Day); 14-16 premi; 17 MBD (memakai nilai
//     CASE, bukan hasil 14.1).
//   - FillPremiGolf (`…\FillPremiGolf.xml`) langkah 2 L399 SELALU menimpa ProRatePercent case dari tanggal, lalu premi = CPCA.
//   - fillPremiBond (`NB FacIn\DataTransform\fillPremiBond.xml` L239): premi r4(r4(TSI·Rate/100)·PR/100) + r4(…·Loading/100), PR =
//     `.ProRatePercent` COVERAGE (kosong = 100); MinPremium: cabang 100000 bila `Policy.Currency == 1` mati (K73-14: Currency
//     tidak pernah di-SET di korpus) -> minPremium = premium.
//   - ChangeCoveragePremium (`…\ChangeCoveragePremium.xml`): day <= 0 -> 365; metode '1' -> ProRatePercent case dari tanggal;
//     `.ProRatePercent` coverage = nilai case. Premi tidak dihitung.
//   - CalculateRate (`D:\migrasi\RNM\DDL\CalculateRate.xml`, hanya varian DDL): Prorate = DiffYear x 100 + pecahan (langkah 2 bila
//     Day terisi, TotalDays = .Day) + pecahan LAGI (langkah 4 bila DiffDays != 0 atau DiffMonth != 0) - dijumlah dua kali, TIRU
//     APA ADANYA (K73-15, OQ); metode 3: langkah 3 TotalDays = TotalDayDiff = 1; ProRatePercent case = Prorate;
//     Rate = (Premium / TSI, skala 20) / ProRatePercent x 10000.
//   - ChangeDiskon (`…\DataTransform\ChangeDiskon.xml`) pre-DT sebelum CPCA: persen -> Discount = r4(Premium x % / 100); amount ->
//     % = Discount / Premium x 100. Langkah FlagDelete tidak berlaku (coverage tanpa FlagDelete di model).
//   - CountPCTLOL (`DDL\CountPCTLOL.xml`): PctLOL -> LoL = TSI x % / 100; AmountLOL -> % = LoL x 100 / TSI (skala 20).
//   - SetIndemnityRate_ACT (`NB FacIn\Activity\SetIndemnityRate_ACT.xml`): langkah 1 CARI1 = .Indemnity, CARI2 = .Unit; 2 RDB
//     `SearchIndemnityRate_SQL` (M_BI_INDEMNITY, DDL ditaruh work owner 09-10-2026; langkah selalu jalan - when `.Unit=='1'`
//     ber-pyStepsPreCondition false L451); 3 loop hasil -> LocalIndemnity = INDEMNITY_PCT (baris TERAKHIR menang); 4
//     .IndemnityPercentage = LocalIndemnity (nol baris -> kosong); 5 `.Unit == 0 || .Unit ==2` (L1079) -> 100; 7 IsMBD -> CPCA
//     Day = .Day (langkah 6 IsFire / CountPremi_ACT bukan jalur Aneka).
// Satuan `EndDateTime - StartDateTime` = hari `[dugaan kuat]` (tanggal WIB, pola tiket FIRE). Day kosong / <= 0 -> 365
// `[keputusan agent]` (fallback CountGrossPremi_Act L5092 / ChangeCoveragePremium langkah 1); FillPremiGolf dari sel % Short
// Period dipanggil TANPA Day (pembagi kosong `belum terverifikasi`) - di sini memakai .Day coverage `[keputusan agent]`.

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services/premium"
)

// Pemicu sel on-change form coverage Aneka.
const (
	PemicuMetode       = "metode"       // Calculate Method L2364 -> CPCA
	PemicuProRate      = "proRate"      // % Pro Rate L2649 -> fillPremiBond
	PemicuShortPeriod  = "shortPeriod"  // % Short Period L2941 -> FillPremiGolf
	PemicuDay          = "day"          // Day L3258 -> ChangeCoveragePremium
	PemicuCurrency     = "currency"     // Currency L3557 -> CPCA Status=Coverage
	PemicuTSI          = "tsi"          // TSI L4069 -> CPCA Status=Coverage
	PemicuLoading      = "loading"      // % Loading L4501 -> CPCA
	PemicuScale        = "scale"        // % Scale L4793 -> FillPremiGolf
	PemicuRate         = "rate"         // % Rate L5168 -> CPCA Status=Coverage
	PemicuIndemnity    = "indemnity"    // Indemnity / Unit L5548 / L5843 -> CPCA (pencarian tabel = OQ)
	PemicuPremium      = "premium"      // Premium L6339 -> CalculateRate
	PemicuDiskonPersen = "diskonPersen" // % Discount L6609 -> ChangeDiskon persen + CPCA
	PemicuDiskonJumlah = "diskonJumlah" // Discount L6981 -> ChangeDiskon amount + CPCA
	PemicuLoL          = "lol"          // Limit of Liability L7649 -> CountPCTLOL AmountLOL
	PemicuPctLoL       = "pctLol"       // % LoL L7922 -> CountPCTLOL PctLOL
)

var pemicuAneka = map[string]bool{PemicuMetode: true, PemicuProRate: true, PemicuShortPeriod: true, PemicuDay: true,
	PemicuCurrency: true, PemicuTSI: true, PemicuLoading: true, PemicuScale: true, PemicuRate: true, PemicuIndemnity: true,
	PemicuPremium: true, PemicuDiskonPersen: true, PemicuDiskonJumlah: true, PemicuLoL: true, PemicuPctLoL: true}

// ErrMasukanHitungAneka - pemicu / isian hitung coverage Aneka tidak sah. 400.
var ErrMasukanHitungAneka = errors.New("services: isian hitung coverage Aneka tidak sah")

// KonteksHitungAneka - data case yang dibaca rantai: periode polis (teks Pega), lini (MBD), ProRatePercent CASE (masuk / keluar);
// CariIndemnity = pencarian M_BI_INDEMNITY (nil = tanpa basis data -> pemicu indemnity 503).
type KonteksHitungAneka struct {
	Mulai, Akhir   string
	MBD            bool
	ProRatePercent *apd.Decimal
	CariIndemnity  func(indemnity, unit string) ([]*apd.Decimal, error)
}

// polaIndemnity - INDEMNITY NUMBER(4,0): selain bilangan bulat <= 4 digit tidak dikueri (tidak ada baris yang cocok; mencegah
// konversi Oracle ORA-01722) `[keputusan agent]`.
var polaIndemnity = regexp.MustCompile(`^[0-9]{1,4}$`)

// setIndemnityRate - SetIndemnityRate_ACT langkah 1-5 dan 7.
func setIndemnityRate(k *KonteksHitungAneka, c *models.CoverageAneka) error {
	if k.CariIndemnity == nil {
		return ErrObjekTanpaDatabase
	}
	var lokal *apd.Decimal
	if ind := strings.TrimSpace(c.Indemnity); polaIndemnity.MatchString(ind) {
		baris, err := k.CariIndemnity(ind, c.Unit)
		if err != nil {
			return err
		}
		for _, b := range baris {
			lokal = b
		}
	}
	c.IndemnityPercentage = lokal
	if u := strings.TrimSpace(c.Unit); u != "" {
		if n, err := strconv.Atoi(u); err == nil && (n == 0 || n == 2) {
			c.IndemnityPercentage = apd.New(100, 0)
		}
	}
	if k.MBD {
		return cpca(k, c, false)
	}
	return nil
}

var tigaRatusLima = apd.New(365, 0)

// hariAneka - `.Day` sebagai pembagi: kosong / bukan bilangan / <= 0 -> 365.
func hariAneka(day string) *apd.Decimal {
	d, err := utils.ParseDecimal(strings.TrimSpace(day))
	if err != nil || d.Sign() <= 0 {
		return tigaRatusLima
	}
	return d
}

// prorateTanggal - `@Math.divide((EndDateTime - StartDateTime), day, 20) * 100`; tanggal tidak terbaca -> galat.
func prorateTanggal(k KonteksHitungAneka, day string) (*apd.Decimal, error) {
	mulai, ok1 := tanggalJakarta(k.Mulai)
	akhir, ok2 := tanggalJakarta(k.Akhir)
	if !ok1 || !ok2 {
		return nil, ErrPeriodeKasus
	}
	h, err := bagi(apd.New(hari(mulai, akhir), 0), hariAneka(day), skalaBagi)
	if err != nil {
		return nil, err
	}
	return kali(h, seratus), nil
}

// selisihTahun - DiffYear / DiffDays / DiffMonth / TotalDayDiff dan panjang tahun ulang-tahun polis ke-DiffYear (CPCA 14.1.1,
// CalculateRate langkah 1). DiffYear = tahun penuh (DateTimeDifference "Y" dibulatkan ke bawah `[dugaan]`, pola Prorate FIRE).
func selisihTahun(k KonteksHitungAneka) (diffYear int, diffDays, diffMonth int, totalDayDiff, panjangUlang, panjangKalender int64, err error) {
	mulai, ok1 := tanggalJakarta(k.Mulai)
	akhir, ok2 := tanggalJakarta(k.Akhir)
	if !ok1 || !ok2 {
		return 0, 0, 0, 0, 0, 0, ErrPeriodeKasus
	}
	diffYear = akhir.Year() - mulai.Year()
	if mulai.AddDate(diffYear, 0, 0).After(akhir) {
		diffYear--
	}
	diffDays = akhir.Day() - mulai.Day()
	diffMonth = int(akhir.Month()) - int(mulai.Month())
	ulang := mulai.AddDate(diffYear, 0, 0)
	totalDayDiff = hari(ulang, akhir)
	panjangUlang = hari(ulang, mulai.AddDate(diffYear+1, 0, 0))
	awalThn := akhir.AddDate(0, 0, 0)
	awalThn = awalThn.AddDate(0, 0, -awalThn.YearDay()+1)
	panjangKalender = hari(awalThn, awalThn.AddDate(1, 0, 0))
	return diffYear, diffDays, diffMonth, totalDayDiff, panjangUlang, panjangKalender, nil
}

// prorateTahunCPCA - CPCA 14.1: DiffYear x 100 (+ TotalDayDiff / TotalDays x 100 bila DiffDays != 0 atau DiffMonth != 0);
// TotalDays = panjang tahun ulang-tahun polis, diganti `.Day` bila Day terisi (14.1.2 - Section selalu mengirim Day).
func prorateTahunCPCA(k KonteksHitungAneka, day string) (*apd.Decimal, error) {
	dy, dd, dm, diff, panjang, _, err := selisihTahun(k)
	if err != nil {
		return nil, err
	}
	total := apd.New(panjang, 0)
	if strings.TrimSpace(day) != "" {
		total = hariAneka(day)
	}
	p := kali(apd.New(int64(dy), 0), seratus)
	if dd != 0 || dm != 0 {
		porsi, err := bagi(apd.New(diff, 0), total, skalaBagi)
		if err != nil {
			return nil, err
		}
		konteksCoverage.Add(p, p, kali(porsi, seratus))
	}
	return p, nil
}

// prorateCalculateRate - CalculateRate langkah 1-4 (pecahan dijumlah dua kali, K73-15).
func prorateCalculateRate(k KonteksHitungAneka, c models.CoverageAneka) (*apd.Decimal, error) {
	dy, dd, dm, diff, _, kalender, err := selisihTahun(k)
	if err != nil {
		return nil, err
	}
	// Langkah 1: TotalDays L458 ditimpa L563 = panjang tahun kalender tahun Akhir.
	total, selisih := apd.New(kalender, 0), apd.New(diff, 0)
	p := kali(apd.New(int64(dy), 0), seratus)
	tambah := func() error {
		porsi, err := bagi(selisih, total, skalaBagi)
		if err != nil {
			return err
		}
		konteksCoverage.Add(p, p, kali(porsi, seratus))
		return nil
	}
	if strings.TrimSpace(c.Day) != "" { // langkah 2
		total = hariAneka(c.Day)
		if err := tambah(); err != nil {
			return nil, err
		}
	}
	if c.CalculateMethodFacIn == "3" { // langkah 3
		total, selisih = apd.New(1, 0), apd.New(1, 0)
	}
	if dd != 0 || dm != 0 { // langkah 4
		if err := tambah(); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func teksAtauNol(d *apd.Decimal) string {
	if d == nil {
		return "0"
	}
	return utils.FormatDecimal(d)
}

func teksAtauKosong(d *apd.Decimal) string {
	if d == nil {
		return ""
	}
	return utils.FormatDecimal(d)
}

// premiCPCA - langkah 11-17 dengan Prorate `p` (metode 1; MBD memakai nilai case `pKasus`). Metode tak dikenal / kosong ->
// langkah 14-16 dilewati, premi tetap.
func premiCPCA(c models.CoverageAneka, p, pKasus *apd.Decimal, mbd bool) (*apd.Decimal, error) {
	metode := c.CalculateMethodFacIn
	if metode != "1" && metode != "2" && metode != "3" {
		return c.Premium, nil
	}
	pakai := p
	if mbd && c.IndemnityPercentage != nil {
		pakai = pKasus
	}
	if metode == "1" && (pakai == nil || pakai.IsZero()) {
		return apd.New(0, 0), nil // TSI x Rate x 0 = 0
	}
	in := premium.Input{LiniBisnis: premium.LiniAneka, CalculateMethod: metode, MataUang: c.Currency, TSI: teksAtauNol(c.TSI),
		Rate: teksAtauNol(c.Rate), ProRatePercent: teksAtauKosong(pakai), PctShortPeriod: teksAtauNol(c.PctShortPeriod),
		Loading: teksAtauKosong(c.Loading), IndemnityPercentage: teksAtauKosong(c.IndemnityPercentage),
		LossLimit: teksAtauKosong(c.LostLimit), MBD: mbd}
	m, err := premium.Calculate(in)
	if err != nil {
		return nil, err
	}
	return m.Amount, nil
}

// cpca - CountPremiCoverageAneka; `status` = "Coverage" menjalankan langkah 4.
func cpca(k *KonteksHitungAneka, c *models.CoverageAneka, status bool) error {
	if status {
		pr, err := prorateTanggal(*k, c.Day)
		if err != nil {
			return err
		}
		k.ProRatePercent = simpan8(pr)
	}
	p := k.ProRatePercent
	if c.CalculateMethodFacIn == "1" && (p == nil || p.IsZero()) {
		var err error
		if p, err = prorateTahunCPCA(*k, c.Day); err != nil {
			return err
		}
	}
	premi, err := premiCPCA(*c, p, k.ProRatePercent, k.MBD)
	if err != nil {
		return err
	}
	c.Premium = simpan8(premi)
	return nil
}

// fillPremiGolf - langkah 2 (ProRatePercent case dari tanggal, selalu) lalu premi.
func fillPremiGolf(k *KonteksHitungAneka, c *models.CoverageAneka) error {
	pr, err := prorateTanggal(*k, c.Day)
	if err != nil {
		return err
	}
	k.ProRatePercent = simpan8(pr)
	premi, err := premiCPCA(*c, k.ProRatePercent, k.ProRatePercent, k.MBD)
	if err != nil {
		return err
	}
	c.Premium = simpan8(premi)
	return nil
}

// fillPremiBond - premi Bond dari `.ProRatePercent` coverage (kosong = 100); minPremium = premium (sakelar 100000 mati, K73-14).
func fillPremiBond(c *models.CoverageAneka) error {
	pr := c.ProRatePercent
	if pr == nil {
		pr = seratus
	}
	dasar, err := bagi(kali(nolBila(c.TSI), nolBila(c.Rate)), seratus, 4)
	if err != nil {
		return err
	}
	pokok, err := bagi(kali(dasar, pr), seratus, 4)
	if err != nil {
		return err
	}
	tambahan, err := bagi(kali(pokok, nolBila(c.Loading)), seratus, 4)
	if err != nil {
		return err
	}
	jumlah := new(apd.Decimal)
	konteksCoverage.Add(jumlah, pokok, tambahan)
	c.Premium = simpan8(jumlah)
	c.MinPremium = c.Premium
	return nil
}

// changeCoveragePremium - day <= 0 -> 365; metode 1 -> ProRatePercent case dari tanggal; coverage = nilai case.
func changeCoveragePremium(k *KonteksHitungAneka, c *models.CoverageAneka) error {
	if c.CalculateMethodFacIn == "1" {
		pr, err := prorateTanggal(*k, c.Day)
		if err != nil {
			return err
		}
		k.ProRatePercent = simpan8(pr)
	}
	c.ProRatePercent = k.ProRatePercent
	return nil
}

// calculateRate - ProRatePercent case dari tanggal (dua kali, K73-15), Rate = (Premium / TSI) / ProRate x 10000. TSI / Prorate
// nol -> rate tidak diubah (pembagian nol di Pega `belum terverifikasi`).
func calculateRate(k *KonteksHitungAneka, c *models.CoverageAneka) error {
	pr, err := prorateCalculateRate(*k, *c)
	if err != nil {
		return err
	}
	k.ProRatePercent = simpan8(pr)
	if c.TSI == nil || c.TSI.IsZero() || pr.IsZero() {
		return nil
	}
	r, err := bagi(nolBila(c.Premium), c.TSI, skalaBagi)
	if err != nil {
		return err
	}
	h := new(apd.Decimal)
	if _, err := konteksCoverage.Quo(h, r, pr); err != nil {
		return err
	}
	c.Rate = simpan8(kali(h, sepuluhRibu))
	return nil
}

// changeDiskon - pre-DT ChangeDiskon (premi SEBELUM CPCA).
func changeDiskon(c *models.CoverageAneka, persen bool) error {
	switch {
	case persen && c.DiscountPercentage != nil:
		d, err := bagi(kali(nolBila(c.Premium), c.DiscountPercentage), seratus, 4)
		if err != nil {
			return err
		}
		c.Discount = simpan8(d)
	case !persen && c.Discount != nil && c.Premium != nil && !c.Premium.IsZero():
		d, err := bagi(c.Discount, c.Premium, skalaBagi)
		if err != nil {
			return err
		}
		c.DiscountPercentage = simpan8(kali(d, seratus))
	}
	return nil
}

// countPctLoL - CountPCTLOL: `amount` (sel LoL) -> %LoL; selain itu (sel %LoL) -> LoL. TSI nol / kosong -> % tidak diubah.
func countPctLoL(c *models.CoverageAneka, amount bool) error {
	if !amount {
		c.LimitOfLiability = simpan8(kali(nolBila(c.TSI), nolBila(c.PctLoL), apd.New(1, -2)))
		return nil
	}
	if c.TSI == nil || c.TSI.IsZero() {
		return nil
	}
	p, err := bagi(kali(nolBila(c.LimitOfLiability), seratus), c.TSI, skalaBagi)
	if err != nil {
		return err
	}
	c.PctLoL = simpan8(p)
	return nil
}

// HitungCoverageAneka - satu baris: rantai pemicu atas salinan coverage; konteks ProRatePercent case diperbarui.
func HitungCoverageAneka(k *KonteksHitungAneka, c models.CoverageAneka, pemicu string) (models.CoverageAneka, error) {
	var err error
	switch pemicu {
	case PemicuMetode, PemicuLoading:
		err = cpca(k, &c, false)
	case PemicuIndemnity:
		err = setIndemnityRate(k, &c)
	case PemicuCurrency, PemicuTSI, PemicuRate:
		err = cpca(k, &c, true)
	case PemicuProRate:
		err = fillPremiBond(&c)
	case PemicuShortPeriod, PemicuScale:
		err = fillPremiGolf(k, &c)
	case PemicuDay:
		err = changeCoveragePremium(k, &c)
	case PemicuPremium:
		err = calculateRate(k, &c)
	case PemicuDiskonPersen, PemicuDiskonJumlah:
		if err = changeDiskon(&c, pemicu == PemicuDiskonPersen); err == nil {
			err = cpca(k, &c, false)
		}
	case PemicuLoL, PemicuPctLoL:
		err = countPctLoL(&c, pemicu == PemicuLoL)
	default:
		return c, fmt.Errorf("%w: pemicu %q tidak dikenal", ErrMasukanHitungAneka, pemicu)
	}
	if errors.Is(err, ErrPeriodeKasus) || errors.Is(err, ErrObjekTanpaDatabase) {
		return c, err
	}
	if err != nil {
		return c, fmt.Errorf("%w: %v", ErrMasukanHitungAneka, err)
	}
	return c, nil
}

// BarisHitungAneka - satu permintaan hitung (badan POST).
type BarisHitungAneka struct {
	Coverage models.CoverageAneka
	Pemicu   string
}

// HitungCoverageAnekaKasus - POST …/hitung-coverage-aneka: periode case + lini (IsMBD) dibaca; baris dihitung berurutan, ProRate
// case berantai antar baris (seperti satu halaman kerja Pega). Tanpa menyimpan.
func (s *Service) HitungCoverageAnekaKasus(ctx context.Context, id string, proRate *apd.Decimal, baris []BarisHitungAneka) ([]models.CoverageAneka,
	*apd.Decimal, error) {
	var masalah []string
	for i, b := range baris {
		if !pemicuAneka[b.Pemicu] {
			masalah = append(masalah, "baris["+strconv.Itoa(i)+"].pemicu tidak dikenal")
		}
	}
	if len(masalah) > 0 {
		return nil, nil, fmt.Errorf("%w: %s", ErrMasukanHitungAneka, strings.Join(masalah, "; "))
	}
	if !idKasusSah(id) {
		return nil, nil, ErrKasusTidakAda
	}
	periode, err := s.periodeKasus(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	kb, _ := s.kasusBisnisDari(ctx, id)
	k := &KonteksHitungAneka{Mulai: periode.Mulai, Akhir: periode.Akhir, MBD: salahSatu(id, kb, "IsMBD"), ProRatePercent: proRate}
	if s.lookupAneka != nil {
		k.CariIndemnity = func(indemnity, unit string) ([]*apd.Decimal, error) {
			return s.lookupAneka.IndemnityPct(ctx, indemnity, unit)
		}
	}
	hasil := make([]models.CoverageAneka, 0, len(baris))
	for i, b := range baris {
		c, err := HitungCoverageAneka(k, b.Coverage, b.Pemicu)
		if err != nil {
			return nil, nil, fmt.Errorf("baris[%d]: %w", i, err)
		}
		hasil = append(hasil, c)
	}
	return hasil, k.ProRatePercent, nil
}

// ProRateKasusAneka - % Pro Rate case untuk GET objek-aneka: T_GENERAL_POLIS.PRO_RATE_PERCENT (migrasi 255, K73-12); NULL /
// tidak terbaca -> dari tanggal (CPCA langkah 4, Day 365); tanggal tidak terbaca -> nil.
func (s *Service) ProRateKasusAneka(ctx context.Context, id string) *apd.Decimal {
	if s.objekAneka != nil {
		if pr, err := s.objekAneka.ProRateKasus(ctx, id); err == nil && pr != nil {
			return pr
		}
	}
	periode, err := s.periodeKasus(ctx, id)
	if err != nil {
		return nil
	}
	pr, err := prorateTanggal(KonteksHitungAneka{Mulai: periode.Mulai, Akhir: periode.Akhir}, "")
	if err != nil {
		return nil
	}
	return simpan8(pr)
}

// medanServerAneka - medan server coverage yang DIKIRIM (Coverages != nil) saat PUT, keluarga CountPremiAndTSIRNM* (K73-13;
// `NB FacIn\Activity\CountPremiAndTSIRNMFireAnekaGolf_ACT.xml` langkah 7 Aneka L3853 / L3905, langkah 8 Golf L6499 / L6551):
// TSINusantaraRe = Share x TSI / 100 (skala 20); PremiNusantaraRe = Share x Premium / 100 - Share x Discount / 100 (skala 20 tiap suku);
// TSILiability = TSI untuk objek okupasi (CountGrossPremi_Act L5352), kosong untuk Golf; MinPremium = Premium (K73-14).
// Share = T_GENERAL_POLIS.PERCENT_SHARE; kosong -> medan RNM kosong.
func medanServerAneka(lokasi []models.LokasiAneka, share *apd.Decimal) {
	isi := func(c *models.CoverageAneka, golf bool) {
		c.MinPremium = c.Premium
		c.TSILiability = nil
		if !golf {
			c.TSILiability = c.TSI
		}
		c.TSINusantaraRe, c.PremiNusantaraRe = nil, nil
		if share == nil {
			return
		}
		if t, err := bagi(kali(share, nolBila(c.TSI)), seratus, skalaBagi); err == nil {
			c.TSINusantaraRe = simpan8(t)
		}
		p, err1 := bagi(kali(share, nolBila(c.Premium)), seratus, skalaBagi)
		d, err2 := bagi(kali(share, nolBila(c.Discount)), seratus, skalaBagi)
		if err1 == nil && err2 == nil {
			h := new(apd.Decimal)
			konteksCoverage.Sub(h, p, d)
			c.PremiNusantaraRe = simpan8(h)
		}
	}
	for i := range lokasi {
		for j := range lokasi[i].Okupasi {
			for n := range lokasi[i].Okupasi[j].Objek {
				for k := range lokasi[i].Okupasi[j].Objek[n].Coverages {
					isi(&lokasi[i].Okupasi[j].Objek[n].Coverages[k], false)
				}
			}
		}
		for n := range lokasi[i].ObjekGolf {
			for k := range lokasi[i].ObjekGolf[n].Coverages {
				isi(&lokasi[i].ObjekGolf[n].Coverages[k], true)
			}
		}
	}
}
