package services

// Premi satu coverage FIRE - tiket E12 (perintah work owner 01-10-2026:
// "port CalculatePremiFire di endorsmentfacin").
//
// Untuk apa berkas ini: `Endorsment Fac In/Activity/CalculatePremiFire.xml`
// (`ASM-FW-GISFW-Data-Coverage / CalculatePremiFire`, EDM-only - tidak ada di
// `NB FacIn`) beserta bagian yang dipanggilnya di langkah 8,
// `Activity/SetLocalNonMbuProrate.xml` langkah 1 dan 10 (cabang FIRE).
//
// Dibaca sesudah: porsiperiode.go (hariBulat), porsipembayaran.go.
//
// `[terverifikasi]` Pemanggilnya aksi klik (`runActivity`) di
// `Section/InputCoverageFire_IsUW.xml` (kelas `ASM-FW-GISFW-Data-Property`,
// 4×), `InputPerCoverageFire_IsUW.xml` (2×), `ViewCoverageFire.xml` (2×) -
// halaman primernya satu baris coverage. `[dugaan]` Baris itu salah satu
// coverage di `pyWorkPage.PropertyList(n).CoverageList` - sehingga TSI yang
// ditulis langkah 6 adalah TSI yang dipakai langkah 7/10 (A11).
//
//	2  IsEDM ∧ FlagOnGoingPolicy!=1   .ProRatePercent = .OldCoverage(1).ProRatePercent
//	3  ¬IsEDM                         .ProRatePercent = pyWorkPage.Policy.PctPremiumTariff
//	6  per PropertyList: Σ item → .TSI = .TSISublimit seluruh CoverageList-nya
//	7  FlagOnGoingPolicy!=1           Premium/PremiRp satu periode
//	8  IsEDM ∧ FlagOnGoingPolicy==1   Call SetLocalNonMbuProrate
//	10 IsEDM ∧ FlagOnGoingPolicy==1   Premium/PremiRp dua bagian (baru + lama)
//
// Langkah 1, 5, 9 ber-label `//` (di-remark) - tidak diport (butir 43).
//
// ⛔ K-046 - tiga ketidaksimetrisan rumus diport APA ADANYA:
//   - 7.1 `PremiRp` memakai `.FirstLossScale` TANPA `@if(…=="",100,…)`;
//   - 10.1 bagian baru `Premium` memakai `.ProratePercentEDMEnd` tanpa `@if`,
//     bagian lamanya `@if(.ProratePercentStartEDM=="",100,…)`;
//   - 10.1 `PremiRp` memakai porsi yang TERTUKAR (StartEDM untuk TSI baru,
//     EDMEnd untuk TSI lama) terhadap `Premium`.
//
// ⚠️ Fixture NB-15 `edm-fire-1.json` tidak menyimpan `pyWorkPage.Policy`,
// `PropertyList`, `FlagOnGoingPolicy`, maupun porsi per coverage - premi
// FIRE TIDAK dapat direkonsiliasi dengan data nyata. Seluruh test sintetis.

import (
	"errors"
	"fmt"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

var (
	// ErrNilaiPremiKosong - nilai kosong ikut dalam perkalian/pembagian rumus
	// premi di tempat yang korpus TIDAK membungkus dengan `@if(…=="",…)`.
	// Perilaku Pega atas aritmetika dengan "" belum terverifikasi.
	ErrNilaiPremiKosong = errors.New("endorsement: nilai kosong dalam rumus premi FIRE; perilaku Pega atas aritmetika nilai kosong belum terverifikasi")
	// ErrCoveragePrimerTidakAda - indeks halaman primer di luar PropertyList.
	ErrCoveragePrimerTidakAda = errors.New("endorsement: coverage primer tidak ada di PropertyList")
)

// ItemTSIFire - satu baris `PropertyList(n).PropertyItemList`
// (`ASM-FW-GISFW-Data-PropertyItem`).
type ItemTSIFire struct {
	TSIObjectItem        uang.Money
	PercentageAdjustment uang.Ratio
	IsAdjustableFlag     string // dibandingkan sebagai teks dengan "true"
	FlagDelete           string // dibandingkan sebagai ANGKA dengan 1 (6.2.2/6.2.3)
	// TSIAdjustment - keluaran 6.2.1. Kosong bila salah satu faktornya kosong
	// (A12).
	TSIAdjustment uang.Money
}

// CoverageLamaFire - `.OldCoverage(1)` coverage primer.
type CoverageLamaFire struct {
	TSI                                                       uang.Money
	Rate, ProRatePercent, FirstLossScale, IndemnityPercentage uang.Ratio
	Premium                                                   uang.Money
}

// CoverageFire - satu baris `CoverageList` (`ASM-FW-GISFW-Data-Coverage`).
// Medan yang ditulis activity ditandai "keluaran".
type CoverageFire struct {
	TSI, TSISublimit                          uang.Money // keluaran langkah 6
	Rate, FirstLossScale, IndemnityPercentage uang.Ratio
	FlagDelete                                string // dibandingkan sebagai teks dengan "1" (7.x/10.x)
	ProRatePercent                            uang.Ratio
	// ProratePercentStartEDM, ProratePercentEDMEnd, ProRatePeriodePercent,
	// CalculateMethod - keluaran SetLocalNonMbuProrate (langkah 8).
	ProratePercentStartEDM, ProratePercentEDMEnd, ProRatePeriodePercent uang.Ratio
	CalculateMethod                                                     string
	Premium, PremiRp                                                    uang.Money // keluaran 7/10
	Old                                                                 CoverageLamaFire
}

// PropertiFire - satu baris `pyWorkPage.PropertyList`.
type PropertiFire struct {
	Items     []ItemTSIFire
	Coverages []CoverageFire
}

// MasukanPremiFire - halaman kerja yang dibaca activity.
type MasukanPremiFire struct {
	// IsEDM, IsFire - rule `When`, dinilai pemanggil (registry EDM belum
	// tersedia di modul ini).
	IsEDM, IsFire bool
	// FlagOnGoingPolicy - `pyWorkPage.FlagOnGoingPolicy`, dibandingkan
	// sebagai ANGKA dengan 1.
	FlagOnGoingPolicy string
	// PctPremiumTariff, CurrencyValue - `pyWorkPage.Policy.*`.
	PctPremiumTariff, CurrencyValue uang.Ratio
	// StartDateTime, EndDateTime - `pyWorkPage.Policy.*`; EdmDate -
	// `pyWorkPage.Quotation.EdmDate`. ⚠️ Halaman `pyWorkPage.Policy` /
	// `.Quotation`, BUKAN `OfferFacIn.PolicyData` / `.QuotationData` - korpus
	// memakai keduanya; hubungannya belum terverifikasi.
	StartDateTime, EndDateTime, EdmDate time.Time
	Properti                            []PropertiFire
	// IndeksProperti, IndeksCoverage - letak halaman primer (A11).
	IndeksProperti, IndeksCoverage int
	// MataUang - mata uang `Local.tsiobject` (awal 0, langkah 6.1).
	MataUang string
	// MataUangRupiah - mata uang `PremiRp`. Kodenya tidak ada di korpus;
	// konfigurasi pemanggil.
	MataUangRupiah string
}

// HasilPremiFire - PropertyList sesudah activity (salinan; masukan tidak
// diubah).
type HasilPremiFire struct {
	Properti []PropertiFire
}

// HitungPremiFire - `CalculatePremiFire` langkah 2-10 berurutan.
func HitungPremiFire(m MasukanPremiFire) (HasilPremiFire, error) {
	prop := salinProperti(m.Properti)
	if m.IndeksProperti < 0 || m.IndeksProperti >= len(prop) ||
		m.IndeksCoverage < 0 || m.IndeksCoverage >= len(prop[m.IndeksProperti].Coverages) {
		return HasilPremiFire{}, ErrCoveragePrimerTidakAda
	}
	berjalan, err := samaKode(m.FlagOnGoingPolicy, models.FlagPolisBerjalan)
	if err != nil {
		return HasilPremiFire{}, err
	}
	c := &prop[m.IndeksProperti].Coverages[m.IndeksCoverage]

	// 2 (IsEDM T=2, ongoing!=1 T='') dan 3 (IsEDM T=3 F='').
	if m.IsEDM && !berjalan {
		c.ProRatePercent = c.Old.ProRatePercent
	}
	if !m.IsEDM {
		c.ProRatePercent = m.PctPremiumTariff
	}

	// 6 - TSI per properti.
	for i := range prop {
		if err := tsiProperti(&prop[i], m.MataUang); err != nil {
			return HasilPremiFire{}, err
		}
	}

	// 7 - polis belum berjalan (tanpa gerbang IsEDM).
	if !berjalan {
		if err := premiSatuPeriode(c, m); err != nil {
			return HasilPremiFire{}, err
		}
	}
	// 8 + 10 - polis berjalan, endorsement.
	if m.IsEDM && berjalan {
		if err := porsiNonMbuFire(c, m); err != nil {
			return HasilPremiFire{}, err
		}
		if err := premiDuaBagian(c, m); err != nil {
			return HasilPremiFire{}, err
		}
	}
	return HasilPremiFire{Properti: prop}, nil
}

// tsiProperti - langkah 6.1-6.3 untuk satu properti.
//
// ⚠️ A12 - 6.2.3 ditulis `.IsAdjustableFlag="true"` (`=` tunggal). Ditafsir
// sebagai perbandingan: `[terverifikasi]` `=` tunggal dipakai di prakondisi
// 62 berkas korpus (`grep -rlE '<pyStepsPreCondParamsWhen>[^<]*[^=!<>]="'`),
// termasuk deret `X="a"||X="b"` yang hanya bermakna sebagai perbandingan.
// Semantik Pega-nya sendiri `[belum terverifikasi]`.
func tsiProperti(p *PropertiFire, mataUang string) error {
	tsi := uang.Money{Amount: apd.New(0, 0), Currency: mataUang} // 6.1
	for j := range p.Items {
		it := &p.Items[j]
		// 6.2.1 - tanpa syarat, semua item.
		it.TSIAdjustment = uang.Money{}
		if !it.TSIObjectItem.Kosong() && !it.PercentageAdjustment.Kosong() {
			kali, err := kaliEksak(it.TSIObjectItem.Amount, it.PercentageAdjustment.Value)
			if err != nil {
				return err
			}
			d, err := bagiHalfUp(kali, apd.New(100, 0), 4)
			if err != nil {
				return err
			}
			it.TSIAdjustment = uang.Money{Amount: d, Currency: it.TSIObjectItem.Currency}
		}
		dihapus, err := samaKode(it.FlagDelete, models.FlagDihapus)
		if err != nil {
			return err
		}
		if dihapus {
			continue
		}
		tambah := it.TSIObjectItem // 6.2.2
		if it.IsAdjustableFlag == models.NilaiBenar {
			tambah = it.TSIAdjustment // 6.2.3
		}
		if tambah.Kosong() {
			return fmt.Errorf("%w (TSI item properti)", ErrNilaiPremiKosong)
		}
		if tsi, err = tsi.Add(tambah); err != nil {
			return err
		}
	}
	for k := range p.Coverages { // 6.3
		p.Coverages[k].TSI, p.Coverages[k].TSISublimit = tsi, tsi
	}
	return nil
}

// premiSatuPeriode - langkah 7.
//
//	7.1 Premium = @Math.divide(TSI*Rate*ProRatePercent*@if(FLS=="",100,FLS)*@if(IP=="",100,IP), 1e9, 4)
//	    PremiRp = @Math.divide(TSI*Rate*ProRatePercent*FLS*@if(IP=="",100,IP), 1e9, 4) * @if(CurrencyValue=="",1,CurrencyValue)
//	7.2 FlagDelete=="1": Premium = PremiRp = 0
func premiSatuPeriode(c *CoverageFire, m MasukanPremiFire) error {
	if c.FlagDelete == models.FlagDihapus {
		c.Premium = uang.Money{Amount: apd.New(0, 0), Currency: m.MataUang}
		c.PremiRp = uang.Money{Amount: apd.New(0, 0), Currency: m.MataUangRupiah}
		return nil
	}
	premi, err := sukuPremi(c.TSI, c.Rate, c.ProRatePercent, bawaan(c.FirstLossScale, 100), bawaan(c.IndemnityPercentage, 100))
	if err != nil {
		return err
	}
	rp, err := sukuPremi(c.TSI, c.Rate, c.ProRatePercent, c.FirstLossScale, bawaan(c.IndemnityPercentage, 100))
	if err != nil {
		return err
	}
	if rp, err = kaliEksak(rp, bawaan(m.CurrencyValue, 1).Value); err != nil {
		return err
	}
	c.Premium = uang.Money{Amount: premi, Currency: m.MataUang}
	c.PremiRp = uang.Money{Amount: rp, Currency: m.MataUangRupiah}
	return nil
}

// porsiNonMbuFire - `SetLocalNonMbuProrate` langkah 1 dan 10 (cabang FIRE).
// Langkah 2-8 hanya menulis lokal/param yang tidak dibaca cabang FIRE; langkah
// 9 cabang PA. Gerbang IsEDM langkah 10 dijamin gerbang pemanggil (langkah 8
// `CalculatePremiFire`: IsEDM ∧ ongoing==1).
//
//	1    CalculateMethod = 1, ProRatePercent = 0, ProRatePeriodePercent = 0
//	10   IsFire ∧ IsEDM (10.1 `pyStepsPreCondition=false` → tetap jalan, P-11):
//	     begEDM = @DateTimeDifference(Policy.StartDateTime, Quotation.EdmDate + 12 jam, D)
//	     EDMEnd = @DateTimeDifference(Quotation.EdmDate, Policy.EndDateTime + 12 jam, D)
//	     ProratePercentStartEDM = @Math.divide(begEDM, 365, 6) * 100
//	     ProratePercentEDMEnd   = @Math.divide(EDMEnd, 365, 6) * 100
//	     ProRatePercent         = jumlah keduanya
//
// ⛔ Lokal `dateDifferentBegEDM`/`EDMEnd` bertipe Decimal, dan +12 jam hanya
// pada argumen kedua: dengan jam yang sama di kedua tanggal, selisihnya SELALU
// n,5 hari. Hasil `@DateTimeDifference(…,D)` atas pecahan hari belum
// terverifikasi → `ErrSatuanSelisihWaktuBelumTerverifikasi` (A09), bukan
// tebakan - sejalan dengan 150,5 hari kasus edm-fire-1.
func porsiNonMbuFire(c *CoverageFire, m MasukanPremiFire) error {
	nol := uang.Ratio{Value: apd.New(0, 0), Scale: 0}
	c.CalculateMethod, c.ProRatePercent, c.ProRatePeriodePercent = models.CalculateMethodProRata, nol, nol
	if !m.IsFire {
		return nil
	}
	if m.StartDateTime.IsZero() || m.EndDateTime.IsZero() || m.EdmDate.IsZero() {
		return ErrTanggalPorsiPeriodeKosong
	}
	setengahHari := 12 * time.Hour
	beg, err := hariBulat(m.EdmDate.Add(setengahHari).Sub(m.StartDateTime))
	if err != nil {
		return err
	}
	akhir, err := hariBulat(m.EndDateTime.Add(setengahHari).Sub(m.EdmDate))
	if err != nil {
		return err
	}
	persen := func(hari int64) (uang.Ratio, error) {
		d, err := bagiHalfUp(apd.New(hari, 0), apd.New(365, 0), 6)
		if err != nil {
			return uang.Ratio{}, err
		}
		if d, err = kaliEksak(d, apd.New(100, 0)); err != nil {
			return uang.Ratio{}, err
		}
		return uang.Ratio{Value: d, Scale: 6}, nil
	}
	if c.ProratePercentStartEDM, err = persen(beg); err != nil {
		return err
	}
	if c.ProratePercentEDMEnd, err = persen(akhir); err != nil {
		return err
	}
	jumlah, err := tambahEksak(c.ProratePercentStartEDM.Value, c.ProratePercentEDMEnd.Value)
	if err != nil {
		return err
	}
	c.ProRatePercent = uang.Ratio{Value: jumlah, Scale: 6}
	return nil
}

// premiDuaBagian - langkah 10.
//
//	10.1 Premium = div4(TSI*Rate*EDMEnd*@if(FLS)*@if(IP)) + div4(Old.TSI*Old.Rate*@if(StartEDM)*@if(Old.FLS)*@if(Old.IP))
//	     PremiRp = (div4(TSI*Rate*StartEDM*@if(FLS)*@if(IP)) + div4(Old.TSI*Old.Rate*@if(EDMEnd)*@if(Old.FLS)*@if(Old.IP))) * @if(CurrencyValue)
//	10.2 FlagDelete=="1":
//	     Premium = Old.Premium − @Math.divide(@Math.divide(Old.Premium*EDMEnd, 1, 4), Old.ProRatePercent, 4)
//	     PremiRp = Premium * @if(CurrencyValue)
//
// (div4(x) = @Math.divide(x, 1e9, 4); @if(v) = @if(v=="",100,v).)
func premiDuaBagian(c *CoverageFire, m MasukanPremiFire) error {
	kurs := bawaan(m.CurrencyValue, 1).Value
	lama := c.Old
	if c.FlagDelete == models.FlagDihapus {
		if lama.Premium.Kosong() || c.ProratePercentEDMEnd.Kosong() || lama.ProRatePercent.Kosong() {
			return fmt.Errorf("%w (10.2)", ErrNilaiPremiKosong)
		}
		kali, err := kaliEksak(lama.Premium.Amount, c.ProratePercentEDMEnd.Value)
		if err != nil {
			return err
		}
		if kali, err = bagiHalfUp(kali, apd.New(1, 0), 4); err != nil {
			return err
		}
		if kali, err = bagiHalfUp(kali, lama.ProRatePercent.Value, 4); err != nil {
			return err
		}
		premi, err := lama.Premium.Sub(uang.Money{Amount: kali, Currency: lama.Premium.Currency})
		if err != nil {
			return err
		}
		rp, err := kaliEksak(premi.Amount, kurs)
		if err != nil {
			return err
		}
		c.Premium, c.PremiRp = premi, uang.Money{Amount: rp, Currency: m.MataUangRupiah}
		return nil
	}
	// premiBaruLama - div4(TSI baru × porsiBaru …) + div4(TSI lama × porsiLama …).
	premiBaruLama := func(porsiBaru, porsiLama uang.Ratio) (*apd.Decimal, error) {
		baru, err := sukuPremi(c.TSI, c.Rate, porsiBaru, bawaan(c.FirstLossScale, 100), bawaan(c.IndemnityPercentage, 100))
		if err != nil {
			return nil, err
		}
		sebelum, err := sukuPremi(lama.TSI, lama.Rate, porsiLama, bawaan(lama.FirstLossScale, 100), bawaan(lama.IndemnityPercentage, 100))
		if err != nil {
			return nil, err
		}
		return tambahEksak(baru, sebelum)
	}
	premi, err := premiBaruLama(c.ProratePercentEDMEnd, bawaan(c.ProratePercentStartEDM, 100))
	if err != nil {
		return err
	}
	rp, err := premiBaruLama(c.ProratePercentStartEDM, bawaan(c.ProratePercentEDMEnd, 100)) // K-046 tertukar
	if err != nil {
		return err
	}
	if rp, err = kaliEksak(rp, kurs); err != nil {
		return err
	}
	c.Premium = uang.Money{Amount: premi, Currency: m.MataUang}
	c.PremiRp = uang.Money{Amount: rp, Currency: m.MataUangRupiah}
	return nil
}

// sukuPremi - @Math.divide(tsi × r1 × … , 1e9, 4). Faktor kosong →
// ErrNilaiPremiKosong (pembungkus `@if` sudah diterapkan pemanggil lewat
// `bawaan`).
func sukuPremi(tsi uang.Money, faktor ...uang.Ratio) (*apd.Decimal, error) {
	if tsi.Kosong() {
		return nil, fmt.Errorf("%w (TSI)", ErrNilaiPremiKosong)
	}
	hasil := tsi.Amount
	for _, f := range faktor {
		if f.Kosong() {
			return nil, ErrNilaiPremiKosong
		}
		var err error
		if hasil, err = kaliEksak(hasil, f.Value); err != nil {
			return nil, err
		}
	}
	return bagiHalfUp(hasil, apd.New(1, 9), 4)
}

// bawaan - `@if(v=="", n, v)`.
func bawaan(v uang.Ratio, n int64) uang.Ratio {
	if v.Kosong() {
		return uang.Ratio{Value: apd.New(n, 0), Scale: 0}
	}
	return v
}

// salinProperti - salinan dalam; medan uang/rasio dibagi pakai (tidak pernah
// dimutasi di tempat - setiap penulisan mengganti pointer).
func salinProperti(src []PropertiFire) []PropertiFire {
	out := make([]PropertiFire, len(src))
	for i, p := range src {
		out[i] = PropertiFire{
			Items:     append([]ItemTSIFire(nil), p.Items...),
			Coverages: append([]CoverageFire(nil), p.Coverages...),
		}
	}
	return out
}
