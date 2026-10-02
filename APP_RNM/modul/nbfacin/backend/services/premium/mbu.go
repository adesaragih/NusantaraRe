// Berkas ini memuat rumus premi lini MBU (tiket NB-05).

package premium

import (
	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// calculateMBU - `.Premium` coverage MBU.
//
// Asal: NB FacIn\Activity\FillPremiMBU_FacIn.xml
// (ASM-FW-GISFW-DATA-COVERAGE / FILLPREMIMBU_FACIN) langkah 2.2.1.3, L1144:
//
//	.Premium = @Math.divide((@Math.divide((.TSI*(.Rate+.Loading)),100,4)*@if(.ProRatePercent=="",100,.ProRatePercent)),100,4)
//
// [terverifikasi] Langkah 2.2.1.3 satu-satunya langkah activity ini yang menulis
// `.Premium`. Ia ber-`pyStepsPreCondition=false`; MBU bukan PA, jadi P-11 berlaku:
// tetap jalan (butir 18). Seluruh 93 baris coverage dari lima kasus MBU nyata
// cocok eksak (`mbu_test.go`).
//
// ⚠️ "L1626" yang dikutip K-018 dan tiket 05 sebagai bentuk kedua BUKAN rumus
// premi: baris itu ada di `pyStepsPreCondParamsWhen` langkah 2.2.1.5, syarat
// `.MinPremium = 100000` (premi < 100.000 dan `pyWorkPage.Policy.Currency == 1`).
// `.MinPremium` dan langkah diskon 2.2.1.7-2.2.1.8 tidak menulis `.Premium` dan
// belum ada pemakainya, jadi tidak diport. L1144 TIDAK mengurangkan diskon.
func calculateMBU(in Input) (uang.Money, error) {
	var tsiNilai, rateNilai, loading, proRataNilai *apd.Decimal
	if err := bacaSemua(
		medanAngka{"TSI", in.TSI, &tsiNilai},
		medanAngka{"Rate", in.Rate, &rateNilai},
		medanAngka{"Loading", in.Loading, &loading},
		medanAngka{"ProRatePercentCoverage", in.ProRatePercentCoverage, &proRataNilai},
	); err != nil {
		return uang.Money{}, err
	}
	tsi := uang.Money{Amount: tsiNilai, Currency: in.MataUang}
	if rateNilai == nil {
		return uang.Money{}, ErrRasioKosong
	}
	if loading == nil {
		loading = apd.New(0, 0) // keputusan agent A12
	}
	// `(.Rate+.Loading)` dijumlahkan dulu - satu rasio ber-satuan rate MBU.
	rateLoading := new(apd.Decimal)
	if err := periksaEksak(utils.DecimalContext().Add(rateLoading, rateNilai, loading)); err != nil {
		return uang.Money{}, err
	}
	rate := rasioRate(LiniMBU, rateLoading)
	if proRataNilai == nil {
		proRataNilai = apd.New(100, 0) // @if(.ProRatePercent=="",100,…)
	}
	proRata := rasioProRata(proRataNilai)

	// Bentuk BERSARANG: pembagi dalam (100) menempel pada rate MBU (%), pembagi
	// luar (100) pada pro-rata (%). Tidak diratakan menjadi ÷10.000: pembulatan
	// dalam terjadi lebih dulu.
	kali, err := rate.kali(tsi)
	if err != nil {
		return uang.Money{}, err
	}
	// Asal: FillPremiMBU_FacIn langkah 2.2.1.3 L1144, @Math.divide dalam, presisi 4.
	dalam, err := bagiBulat(kali, rate.satuan.Pembagi(), 4)
	if err != nil {
		return uang.Money{}, err
	}
	kaliProRata, err := proRata.kali(dalam)
	if err != nil {
		return uang.Money{}, err
	}
	// Asal: FillPremiMBU_FacIn langkah 2.2.1.3 L1144, @Math.divide luar, presisi 4.
	return bagiBulat(kaliProRata, proRata.satuan.Pembagi(), 4)
}
