package premium

// Rumus premi coverage FIRE, ANEKA (termasuk Bonding, A11), GOLF, MARINE CARGO -
// tiket 18 (disusun agent dari XML atas perintah work owner). Semua berkas di
// `D:\migrasi\RNM\NB FacIn\Activity\`; pemilih rumusnya CountGrossPremi_Act
// langkah 5.1-5.4 (jembatan.go).
//
// Setiap `@Math.divide` sistem lama = SATU premiKomposit: uang × rasio-rasio, dibagi
// hasil kali pembagi satuan rasio itu, dibulatkan presisi langkahnya. Pembagi
// komposit di XML (10⁴, 10⁶, 10⁸, 10⁹, 10¹¹, 10¹³) TERURAI menjadi satuan K-018:
// rate FIRE ‰ (1000), rate lain % (100), dan setiap medan persen 100.
//
// Masukan yang disediakan pemanggil (tidak dihitung di sini): pro-rata
// (`ProRatePercent`, dari aritmetika tanggal CountPremi_ACT langkah 3-9 /
// CountPremiCoverageAneka langkah 4, 14.1), `IndemnityPercentage` dan `FirstScale`
// (lookup SQL), `PctAdjustment` (flag item properti, CountPremi_ACT langkah 11-16).

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// premiKomposit - satu `@Math.divide((m × r1 × … × rn), Π pembagi, desimal)`.
func premiKomposit(m uang.Money, desimal int32, faktor ...rasio) (uang.Money, error) {
	pembagi := int64(1)
	x := m
	for _, r := range faktor {
		var err error
		if x, err = r.kali(x); err != nil {
			return uang.Money{}, err
		}
		pembagi *= r.satuan.Pembagi()
	}
	return bagiBulat(x, pembagi, desimal)
}

var (
	// ErrCoverageBasisTakDikenal - FIRE `.CoverageBasis` kosong atau di luar 1-5.
	// [terverifikasi] Setiap langkah premi CountPremi_ACT 19-50 bergerbang
	// `.CoverageBasis==n`; tanpa yang terbuka, `.Premium` lama dibiarkan apa adanya -
	// tidak ditiru diam-diam (seperti ErrMetodeTakDikenal).
	ErrCoverageBasisTakDikenal = errors.New("premium: CoverageBasis FIRE tidak dikenal")
	// ErrTeksNolAmbigu - medan yang sistem lama bandingkan sebagai TEKS dengan "0"
	// (`.NetRate!="0"`, `Local.LossLimit=="0"`) bernilai nol tetapi berbentuk lain,
	// mis. "0.0". [terverifikasi] Pega menulis desimal nol sebagai "0.0"
	// (`PctAdjustment` 49 dari 54 coverage FIRE nyata). [pertanyaan terbuka] apakah
	// pembanding Pega atas properti desimal itu teks atau angka: "0.0" akan
	// memakai NetRate 0 (premi 0) bila teks, rate biasa bila angka. Ditolak (A43).
	ErrTeksNolAmbigu = errors.New("premium: nol berbentuk teks lain dari \"0\" pada pembanding teks sistem lama")
)

// kaliLossLimitPersen - `.Premium*Local.Losslimit/100` (ANEKA/GOLF): perkalian dan
// pembagian TANPA @Math.divide, jadi tanpa pembulatan; hasil yang tidak eksak
// ditolak (`[dugaan]` operator `/` Pega atas desimal eksak).
//
// Skala hasil = skala pembilang dikurangi skala pembagi, diperpanjang hanya bila
// hasil eksaknya butuh lebih (`[dugaan]` aturan skala pilihan pembagian eksak
// Java BigDecimal; A38). [terverifikasi] Kasus nyata Kredit
// (`testdata/kasus/nb-kredit-1.json`, coverage
// `LocationList[0].Property.RiskLocation.OccupationList[0].AnekaList[0].CoverageList[0]`,
// rekonsiliasi `TestBerkasKasusP5`): pembilang berskala 4, tersimpan berskala 4,
// bukan 38 digit penuh.
func kaliLossLimitPersen(m uang.Money, lossLimit *apd.Decimal) (uang.Money, error) {
	kali, err := rasioPersen(lossLimit).kali(m)
	if err != nil {
		return uang.Money{}, err
	}
	ctx := utils.DecimalContext()
	pembagi := apd.New(Persen.Pembagi(), 0)
	hasil := new(apd.Decimal)
	if err := periksaEksak(ctx.Quo(hasil, kali.Amount, pembagi)); err != nil {
		return uang.Money{}, err
	}
	hasil.Reduce(hasil)
	if ideal := kali.Amount.Exponent - pembagi.Exponent; hasil.Exponent > ideal {
		if err := periksaEksak(ctx.Quantize(hasil, hasil, ideal)); err != nil {
			return uang.Money{}, err
		}
	}
	return uang.Money{Amount: hasil, Currency: m.Currency}, nil
}

// lossLimitFire - CountPremi_ACT langkah 17 L3572 `Local.LossLimit=="" ||
// Local.LossLimit=="0"` → 100: pembanding TEKS, jadi nol berbentuk lain ditolak (A43).
func lossLimitFire(teks string, d *apd.Decimal) (*apd.Decimal, error) {
	if teks == "" || teks == "0" {
		return apd.New(100, 0), nil
	}
	if d.IsZero() {
		return nil, fmt.Errorf("%w: LossLimit %q", ErrTeksNolAmbigu, teks)
	}
	return d, nil
}

// lossLimitBawaan - ANEKA/GOLF `Local.Losslimit==0||Local.Losslimit==""` → 100:
// pembanding ANGKA (CountPremiCoverageAneka L2480, FillPremiGolf L1580).
func lossLimitBawaan(d *apd.Decimal) *apd.Decimal {
	if d == nil || d.IsZero() {
		return apd.New(100, 0)
	}
	return d
}

var errLayering = errors.New("CoverageBasis 5 (Layering) tidak dipakai di sistem baru (butir 30)")

// calculateFire - CountPremi_ACT (ASM-FW-GISFW-DATA-COVERAGE!COUNTPREMI_ACT), jalur
// `Param.PremiStatus == "percent"`. Langkah yang terbuka saling menimpa menurut
// urutannya; hasil akhirnya:
//
//	basis 1/3/4  langkah 19/35/43  TSI×Rate×Prorate×Indemnity×LossLimit / 1e9   (20)
//	basis 2      langkah 27        … × FirstScale                    / 1e11  (20)
//	adjustable   langkah 21/29/37/45  … × PctAdjustment, pembagi × 100
//	net rate     langkah 23-25/31-33/39-41/47-49  NetRate menggantikan Rate
//
// Basis 3 dan 4 hanya mengubah TSILiability, bukan rumus premi. Basis 5 = Layering.
func calculateFire(in Input) (uang.Money, error) {
	var tsi, rate, netRate, proRata, indemnity, lossLimit, adj, firstScale *apd.Decimal
	if err := bacaSemua(
		medanAngka{"TSI", in.TSI, &tsi},
		medanAngka{"Rate", in.Rate, &rate},
		medanAngka{"NetRate", in.NetRate, &netRate},
		medanAngka{"ProRatePercent", in.ProRatePercent, &proRata},
		medanAngka{"IndemnityPercentage", in.IndemnityPercentage, &indemnity},
		medanAngka{"LossLimit", in.LossLimit, &lossLimit},
		medanAngka{"PctAdjustment", in.PctAdjustment, &adj},
		medanAngka{"FirstScale", in.FirstScale, &firstScale},
	); err != nil {
		return uang.Money{}, err
	}
	// `.NetRate!="0"&&.NetRate!=""` - dibandingkan sebagai TEKS (L4681).
	if in.NetRate != "0" && in.NetRate != "" {
		if netRate.IsZero() {
			return uang.Money{}, fmt.Errorf("%w: NetRate %q", ErrTeksNolAmbigu, in.NetRate)
		}
		rate = netRate
	}
	lossLimit, err := lossLimitFire(in.LossLimit, lossLimit)
	if err != nil {
		return uang.Money{}, err
	}
	faktor := []rasio{rasioRate(LiniFire, rate), rasioProRata(proRata), rasioPersen(indemnity), rasioPersen(lossLimit)}
	switch in.CoverageBasis {
	case "1", "3", "4":
	case "2":
		faktor = append(faktor, rasioPersen(firstScale))
	case "5":
		return uang.Money{}, fmt.Errorf("%w: %v", ErrBentukBelumDiport, errLayering)
	default:
		return uang.Money{}, fmt.Errorf("%w: %q", ErrCoverageBasisTakDikenal, in.CoverageBasis)
	}
	if adj != nil && !adj.IsZero() { // `.PctAdjustment!=0`; kosong = 0 (A36)
		faktor = append(faktor, rasioPersen(adj))
	}
	return premiKomposit(uang.Money{Amount: tsi, Currency: in.MataUang}, 20, faktor...)
}

// calculateAnekaGolf - CountPremiCoverageAneka (…!COUNTPREMICOVERAGEANEKA) langkah
// 14-17 dan FillPremiGolf (…!FILLPREMIGOLF) langkah 9-12 - rumusnya identik:
//
//	round4(TSI×Rate×P / 1e4) + round4(TSI×Rate×P×Loading / 1e6), lalu × LossLimit / 100
//	MBD (IsMBD dan IndemnityPercentage terisi): × Indemnity, pembagi 1e6 / 1e8
//
// P: metode 1 ProRatePercent, 2 PctShortPeriod, 3 = 100. Langkah 1-3 (EDM) tidak diport.
func calculateAnekaGolf(in Input) (uang.Money, error) {
	var tsiNilai, rate, proRata, periode, loading, indemnity, lossLimit *apd.Decimal
	if err := bacaSemua(
		medanAngka{"TSI", in.TSI, &tsiNilai},
		medanAngka{"Rate", in.Rate, &rate},
		medanAngka{"ProRatePercent", in.ProRatePercent, &proRata},
		medanAngka{"PctShortPeriod", in.PctShortPeriod, &periode},
		medanAngka{"Loading", in.Loading, &loading},
		medanAngka{"IndemnityPercentage", in.IndemnityPercentage, &indemnity},
		medanAngka{"LossLimit", in.LossLimit, &lossLimit},
	); err != nil {
		return uang.Money{}, err
	}
	var p rasio
	switch in.CalculateMethod {
	case "1":
		// Langkah 14.1 (CountPremiCoverageAneka L3311) `Local.Prorate == 0 ||
		// Local.Prorate = ""`: pro-rata dihitung ulang dari tanggal polis (14.1.1-14.1.3,
		// @DateTime) - belum diport. FillPremiGolf tidak punya langkah itu.
		if in.LiniBisnis != LiniGolf && (proRata == nil || proRata.IsZero()) {
			return uang.Money{}, fmt.Errorf("%w: pro-rata ANEKA dari tanggal (langkah 14.1)", ErrBentukBelumDiport)
		}
		p = rasioProRata(proRata)
	case "2":
		p = rasioPeriodePendek(periode)
	case "3":
		p = rasioProRata(apd.New(100, 0))
	default:
		return uang.Money{}, fmt.Errorf("%w: %q", ErrMetodeTakDikenal, in.CalculateMethod)
	}
	if loading == nil {
		// Kosong = 0: keputusan agent A42 (menunggu konfirmasi), meniru A12 MBU.
		// [terverifikasi] satu-satunya coverage ANEKA nyata ber-Loading "0", jadi
		// tafsir kosong belum teruji data.
		loading = apd.New(0, 0)
	}
	dasar := []rasio{rasioRate(in.LiniBisnis, rate), p}
	// Langkah 17.x / 12.x menimpa 14-16 hanya untuk IsMBD; `.IndemnityPercentage==""`
	// kembali ke rumus tanpa indemnity.
	if in.MBD && in.IndemnityPercentage != "" {
		dasar = append(dasar, rasioPersen(indemnity))
	}
	tsi := uang.Money{Amount: tsiNilai, Currency: in.MataUang}
	pokok, err := premiKomposit(tsi, 4, dasar...)
	if err != nil {
		return uang.Money{}, err
	}
	tambahan, err := premiKomposit(tsi, 4, append(dasar, rasioPersen(loading))...)
	if err != nil {
		return uang.Money{}, err
	}
	jumlah, err := pokok.Add(tambahan)
	if err != nil {
		return uang.Money{}, err
	}
	return kaliLossLimitPersen(jumlah, lossLimitBawaan(lossLimit))
}

// calculateMarine - CountGPWMarinePAMbu_Act (ASM-FW-GISFW-WORK!COUNTGPWMARINEPAMBU_ACT)
// langkah 1.1.1.1.1 `@Math.divide((.TSI*.Rate),100,4)`; 1.1.1.1.2 master policy → 0.
func calculateMarine(in Input) (uang.Money, error) {
	var tsi, rate *apd.Decimal
	if err := bacaSemua(medanAngka{"TSI", in.TSI, &tsi}, medanAngka{"Rate", in.Rate, &rate}); err != nil {
		return uang.Money{}, err
	}
	if in.MasterPolicy {
		return uang.Money{Amount: apd.New(0, 0), Currency: in.MataUang}, nil
	}
	return premiKomposit(uang.Money{Amount: tsi, Currency: in.MataUang}, 4, rasioRate(LiniMarineCargo, rate))
}
