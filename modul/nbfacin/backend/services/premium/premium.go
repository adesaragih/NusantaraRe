// Package premium adalah Seam 3 siklus New Business Fac In: SATU pintu masuk
// perhitungan premi, `Calculate`. Percabangan lini bisnis terjadi di dalam;
// rumus per lini bukan seam (spec `docs/04-spec/03-spec-modul-terverifikasi.md`
// Modul 5, ADR-F-0001).
//
// Tiket NB-01 (tracer): rumus premi lini PA, activity `CalculatePremiPA_FacIn`.
// Tiket NB-03: satuan setiap rasio datang dari `resolver.go` (Seam 2), untuk
// tujuh lini K-018. Tiket 05: MBU. Tiket 18: FIRE, ANEKA (termasuk Bonding, A11),
// GOLF, MARINE CARGO (`lini_lain.go`); BONDING sendiri `ErrBentukBelumDiport`.
// Keputusan yang mengikat: `docs/KEPUTUSAN-30-09-2026.md` butir 4-18 dan 26-27,
// serta bab Comments tiket 01 dan 03.
package premium

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

var (
	// ErrAngkaTakTerbaca - teks masukan bukan angka desimal bertitik.
	ErrAngkaTakTerbaca = errors.New("premium: angka masukan tidak terbaca")
	// ErrRasioKosong - rate atau pro rata yang dipakai rumus tidak terisi.
	// Kosong bukan nol (ADR-U-0022, ADR-U-0027).
	ErrRasioKosong = errors.New("premium: rasio kosong")
	// ErrMelampauiPresisi - hasil antara melampaui presisi 38 digit
	// `utils.DecimalPrecision`, sehingga apd akan membulatkannya diam-diam.
	// Pengaman teknis sistem baru, BUKAN perilaku Pega (dicatat di
	// `docs/KEPUTUSAN-30-09-2026.md`).
	ErrMelampauiPresisi = errors.New("premium: hasil antara melampaui presisi desimal sistem")
	// ErrBentukBelumDiport - lini bisnis atau `CalculateMethod_FacIn` sah
	// tetapi rumusnya belum diport (tiket 04-06).
	ErrBentukBelumDiport = errors.New("premium: bentuk rumus ini belum diport")
	// ErrMetodeTakDikenal - `CalculateMethod_FacIn` di luar '1', '2', '3'.
	// [terverifikasi] Di sistem lama tidak satu pun langkah menulis
	// `.Premium` untuk nilai lain, sehingga premi lama dibiarkan apa adanya;
	// perilaku itu tidak ditiru diam-diam.
	ErrMetodeTakDikenal = errors.New("premium: CalculateMethod_FacIn tidak dikenal")
)

// Input - satu kasus coverage dengan SEMUA angka berupa teks desimal bertitik,
// seperti data kerja Pega (berkas kasus `DDL\CONTOH`).
//
// ⚠️ Batas masukan SEMENTARA (keputusan work owner 01-10-2026). ADR-U-0022
// menaruh konversi di pemuat `repository`; pemuat NB belum ada (menunggu tabel
// flat), sehingga konversi untuk sementara terjadi SEKALI di awal fungsi tiap
// lini (`calculatePA`, `calculateMBU`, dan ketiga fungsi `lini_lain.go`), lewat
// `bacaSemua`.
// Koma desimal K-027 adalah format kolom Oracle `FACINOFFER.RATE` — tugas
// pemuat itu kelak, bukan paket ini.
type Input struct {
	LiniBisnis LiniBisnis
	// CalculateMethod - `.CalculateMethod_FacIn`, dibandingkan sebagai TEKS
	// seperti langkah 4-6 (`=='1'`, `=='2'`, `=='3'`).
	CalculateMethod string
	// MataUang ikut ke hasil apa adanya. Kosong tidak ditolak: sistem lama
	// tidak pernah menetapkan mata uang bawaan (berkas kasus PA pun kosong),
	// dan penanganannya masih [pertanyaan terbuka].
	MataUang string
	// TSI - `.TSI` coverage.
	TSI string
	// Rate - `.Rate` coverage; satuannya dari resolver lini bisnis (K-018).
	Rate string
	// ProRatePercent - `pyWorkPage.OfferFacIn.ProRatePercent`, persen. Hanya
	// dipakai metode '1'.
	ProRatePercent string
	// PctShortPeriod - `.PctShortPeriod` coverage, persen. Hanya dipakai metode
	// '2' (langkah 5 L858).
	PctShortPeriod string
	// Discount - `.Discount` coverage. Kosong = dikurangi NOL: [terverifikasi]
	// empat kasus PA nyata tanpa tag `Discount` cocok eksak dengan rumus
	// berdiskon nol (keputusan work owner 01-10-2026).
	Discount string
	// DiscountType - `param.DiscountType` activity, dibandingkan sebagai TEKS
	// seperti langkah 2-3 (`=="Percent"`, `=="Amount"`).
	DiscountType string
	// DiscountPercentage - `.DiscountPercentage` coverage, persen; dipakai
	// langkah 2.
	DiscountPercentage string
	// PremiSebelumnya - `.Premium` coverage SEBELUM activity berjalan. Langkah 2
	// membacanya sebelum langkah 4-6 menulis premi baru.
	PremiSebelumnya string

	// Loading - `.Loading` coverage MBU, satuan sama dengan rate, DIJUMLAHKAN ke
	// rate sebelum perkalian (L1144). Kosong = 0: keputusan agent A12 (menunggu
	// konfirmasi), dasarnya [terverifikasi] 93 baris coverage MBU nyata tanpa
	// Loading cocok eksak (cara hitung: tiket 05 bab Comments).
	Loading string
	// ProRatePercentCoverage - `.ProRatePercent` COVERAGE MBU (L1144). BUKAN
	// `pyWorkPage.OfferFacIn.ProRatePercent` yang dipakai PA: dua properti
	// berbeda. Kosong = 100, sesuai `@if(.ProRatePercent=="",100,.ProRatePercent)`.
	ProRatePercentCoverage string

	// Medan tiket 18 (FIRE, ANEKA/GOLF, MARINE CARGO). Teks apa adanya dari
	// coverage; artinya di lini_lain.go.
	IndemnityPercentage string // `.IndemnityPercentage` (hasil SearchIndemnityRate_SQL)
	LossLimit           string // `.LostLimit`; kosong atau 0 = 100
	PctAdjustment       string // `.PctAdjustment` FIRE; kosong = 0 (A36)
	CoverageBasis       string // `.CoverageBasis` FIRE: 1 SI, 2 First Loss, 3 EML/PML, 4 Sublimit, 5 Layering
	NetRate             string // `.NetRate` FIRE; "" dan "0" = tidak dipakai
	FirstScale          string // `.FirstScale` FIRE basis 2 (hasil SearchFirstLossScale1_SQL)
	// MBD - predikat IsMBD kasus (ANEKA/GOLF): rumus ber-indemnity bila
	// IndemnityPercentage terisi.
	MBD bool
	// MasterPolicy - MARINE CARGO `QuotationData.PolicyType == 1 && IsMOP == "MOP"`.
	MasterPolicy bool
}

// rasio - `uang.Ratio` beserta satuannya. Satuan TIDAK ditaruh di
// `uang.Ratio.Scale`: pemakai yang ada (`claimlife`) mengisinya dengan jumlah
// desimal, bukan ‰/% (keputusan work owner 30-09-2026: tipe lokal nbfacin,
// `inti` tidak diubah). Karena itu `Scale` dibiarkan nol di paket ini.
type rasio struct {
	nilai  uang.Ratio
	satuan Satuan
}

// kali - SATU-SATUNYA jembatan uang ↔ rasio: m × nilai rasio, utuh, tanpa
// pembagi satuan dan tanpa pembulatan. Pembagi satuan diterapkan pemanggil
// sekali, di tempat rumus asalnya membaginya (aturan pembagi komposit K-018).
func (r rasio) kali(m uang.Money) (uang.Money, error) {
	if m.Kosong() {
		return uang.Money{}, uang.ErrUangKosong
	}
	if r.nilai.Kosong() {
		return uang.Money{}, ErrRasioKosong
	}
	hasil := new(apd.Decimal)
	if err := periksaEksak(utils.DecimalContext().Mul(hasil, m.Amount, r.nilai.Value)); err != nil {
		return uang.Money{}, err
	}
	return uang.Money{Amount: hasil, Currency: m.Currency}, nil
}

// bacaDesimal - konversi tunggal `inti` (`utils.ParseDecimal`, ADR-U-0034).
// Kosong = kosong (nil), bukan nol (ADR-U-0022).
func bacaDesimal(medan, s string) (*apd.Decimal, error) {
	if s == "" {
		return nil, nil
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %q", ErrAngkaTakTerbaca, medan, s)
	}
	return d, nil
}

// medanAngka - satu medan teks dan TUJUAN desimalnya. Nama dan nilai dipasangkan
// lewat penunjuk tujuan, bukan lewat urutan dalam larik (CLAUDE.md §4a,
// jebakan 6: pasangkan lewat kunci, bukan urutan).
type medanAngka struct {
	nama, teks string
	ke         **apd.Decimal
}

// bacaSemua - konversi semua medan sekali, di batas masukan.
func bacaSemua(medan ...medanAngka) error {
	for _, m := range medan {
		d, err := bacaDesimal(m.nama, m.teks)
		if err != nil {
			return err
		}
		*m.ke = d
	}
	return nil
}

// Calculate menghitung premi satu coverage.
func Calculate(in Input) (uang.Money, error) {
	// SatuanRate memastikan lini ada di peta K-018 (panic bila tidak), sebelum
	// masukan mana pun dibaca.
	satuanRate := SatuanRate(in.LiniBisnis)
	switch in.LiniBisnis {
	case LiniPA:
		return calculatePA(in)
	case LiniMBU:
		return calculateMBU(in)
	case LiniFire:
		return calculateFire(in)
	case LiniAneka, LiniGolf:
		return calculateAnekaGolf(in)
	case LiniMarineCargo:
		return calculateMarine(in)
	case LiniBonding:
		// [terverifikasi] tidak ada gerbang Bonding di CountGrossPremi_Act; kasus
		// Bonding membuka IsAneka dan dihitung lini ANEKA (A11).
		return uang.Money{}, fmt.Errorf("%w: kasus Bonding dihitung lewat lini ANEKA (A11)", ErrBentukBelumDiport)
	}
	return uang.Money{}, fmt.Errorf("%w: lini %s (rate %v)", ErrBentukBelumDiport, in.LiniBisnis, satuanRate)
}

// calculatePA - `.Premium` coverage PA (activity CalculatePremiPA_FacIn).
func calculatePA(in Input) (uang.Money, error) {
	var tsi, rate, proRata, diskon, periode, diskonPersen, premiLama *apd.Decimal
	if err := bacaSemua(
		medanAngka{"TSI", in.TSI, &tsi},
		medanAngka{"Rate", in.Rate, &rate},
		medanAngka{"ProRatePercent", in.ProRatePercent, &proRata},
		medanAngka{"Discount", in.Discount, &diskon},
		medanAngka{"PctShortPeriod", in.PctShortPeriod, &periode},
		medanAngka{"DiscountPercentage", in.DiscountPercentage, &diskonPersen},
		medanAngka{"PremiSebelumnya", in.PremiSebelumnya, &premiLama},
	); err != nil {
		return uang.Money{}, err
	}
	if diskon == nil {
		diskon = apd.New(0, 0) // Discount kosong = dikurangi nol (butir 17)
	}
	b := bahanPA{
		tsi:          uang.Money{Amount: tsi, Currency: in.MataUang},
		rate:         rasioRate(LiniPA, rate),
		proRata:      rasioProRata(proRata),
		periode:      rasioPeriodePendek(periode),
		diskon:       uang.Money{Amount: diskon, Currency: in.MataUang},
		diskonPersen: rasioDiskon(diskonPersen),
		premiLama:    uang.Money{Amount: premiLama, Currency: in.MataUang},
	}
	return premiPA(in.CalculateMethod, in.DiscountType, b)
}

// bahanPA - masukan rumus PA setelah dikonversi.
type bahanPA struct {
	tsi, diskon, premiLama               uang.Money
	rate, proRata, periode, diskonPersen rasio
}

// premiPA - `.Premium` akhir activity CalculatePremiPA_FacIn.
//
// Asal: NB FacIn\Activity\CalculatePremiPA_FacIn.xml
// (ASM-FW-GISFW-DATA-COVERAGE / CALCULATEPREMIPA_FACIN). Tiga langkah menulis
// `.Premium`, masing-masing bergerbang `.CalculateMethod_FacIn` sebagai teks:
//
//	langkah 4 L713  '1' prorata       (@Math.divide((.TSI * .Rate),100000,4)*pyWorkPage.OfferFacIn.ProRatePercent) - .Discount
//	langkah 5 L858  '2' short period  (@Math.divide((.TSI*.Rate),1000,4)*@Math.divide((.PctShortPeriod),100,4))  - .Discount
//	langkah 6 L1003 '3' fixrate       (@Math.divide((.TSI*.Rate),1000,4)*1)-.Discount
//
// ⛔ Langkah 7 (L1146, `pyStepsPreCondition=false`, syarat tersimpan
// `.CalculateMethod_FacIn==1`) TIDAK dijalankan: labelnya `//` = di-remark,
// tidak dipakai lagi (keputusan work owner 01-10-2026, butir 43). Data
// menguatkannya: kasus PA metode 1 nyata cocok eksak dengan langkah 4 dan
// meleset 0,0048 dari langkah 7; tiga kasus metode 3 cocok eksak dengan
// langkah 6 (butir 13). Konsisten dengan butir 18 (pengecualian P-11 hanya PA):
// di antara activity premi yang diport, hanya langkah ini yang berlabel `//`.
//
// Langkah 2-3 berjalan SEBELUM langkah 4-6 dan bergerbang `param.DiscountType`:
//
//	langkah 2 "Percent"  .Discount = @Math.divide((.DiscountPercentage *.Premium),100,20)
//	langkah 3 "Amount"   .DiscountPercentage = @Math.divide(.Discount,.Premium,20)*100
//
// Langkah 2 membaca `.Premium` LAMA (premi baru belum ditulis) dan mengubah
// `.Discount` yang lalu dikurangkan langkah 4-6. Langkah 3 hanya menulis
// `.DiscountPercentage` - tidak memengaruhi premi, dan belum ada pemakainya, jadi
// tidak diport.
func premiPA(metode, jenisDiskon string, b bahanPA) (uang.Money, error) {
	diskon := b.diskon
	if jenisDiskon == "Percent" {
		// Asal: CalculatePremiPA_FacIn langkah 2, presisi 20.
		kali, err := b.diskonPersen.kali(b.premiLama)
		if err != nil {
			return uang.Money{}, err
		}
		if diskon, err = bagiBulat(kali, b.diskonPersen.satuan.Pembagi(), 20); err != nil {
			return uang.Money{}, err
		}
	}
	rumus, ada := rumusPA[metode]
	if !ada {
		return uang.Money{}, fmt.Errorf("%w: %q", ErrMetodeTakDikenal, metode)
	}
	premi, err := rumus.hitung(b)
	if err != nil {
		return uang.Money{}, err
	}
	// `- .Discount` di luar `@Math.divide`, tanpa pembulatan lagi.
	return premi.Sub(diskon)
}

// rumusPA - `.CalculateMethod_FacIn` → rumus CalculatePremiPA_FacIn yang dipakai,
// beserta asalnya untuk laporan rekonsiliasi (AsalRumus).
var rumusPA = map[string]struct {
	asal   string
	hitung func(bahanPA) (uang.Money, error)
}{
	"1": {"CalculatePremiPA_FacIn langkah 4 L713", func(b bahanPA) (uang.Money, error) { return premiPAProrata(b.tsi, b.rate, b.proRata) }},
	"2": {"CalculatePremiPA_FacIn langkah 5 L858", func(b bahanPA) (uang.Money, error) { return premiPAPeriodePendek(b.tsi, b.rate, b.periode) }},
	"3": {"CalculatePremiPA_FacIn langkah 6 L1003", func(b bahanPA) (uang.Money, error) { return premiPAFixrate(b.tsi, b.rate) }},
}

// asalMBU - rumus premi MBU yang diport (A4).
const asalMBU = "FillPremiMBU_FacIn L1144"

// AsalRumus - rule dan langkah sistem lama yang menghasilkan premi `in`, untuk
// laporan rekonsiliasi (tiket 16); kosong bila lini/metodenya belum diport.
func AsalRumus(in Input) string {
	switch in.LiniBisnis {
	case LiniPA:
		return rumusPA[in.CalculateMethod].asal
	case LiniMBU:
		return asalMBU
	case LiniFire:
		switch in.CoverageBasis {
		case "1", "2", "3", "4":
			return "CountPremi_ACT langkah 19-50 basis " + in.CoverageBasis
		}
	case LiniAneka:
		return "CountPremiCoverageAneka langkah 14-17"
	case LiniGolf:
		return "FillPremiGolf langkah 9-12"
	case LiniMarineCargo:
		return "CountGPWMarinePAMbu_Act langkah 1.1.1.1.1"
	}
	return ""
}

// premiPAProrata - langkah 4 L713 tanpa diskon:
// round4((TSI × Rate) / 100000) × ProRatePercent.
func premiPAProrata(tsi uang.Money, rate, proRata rasio) (uang.Money, error) {
	kali, err := rate.kali(tsi)
	if err != nil {
		return uang.Money{}, err
	}
	// Pembagi komposit 100000 = 1.000 (Rate ber-‰) × 100 (ProRatePercent
	// ber-%), walaupun ProRatePercent dikalikan DI LUAR pembagian — aturan
	// penguraian K-018.
	bagi, err := bagiBulat(kali, rate.satuan.Pembagi()*proRata.satuan.Pembagi(),
		// Asal: CalculatePremiPA_FacIn langkah 4 L713, presisi 4.
		4)
	if err != nil {
		return uang.Money{}, err
	}
	// Urutan L713: round(a,4) × b, bukan round(a × b, 4). ProRatePercent
	// dikalikan utuh, tanpa pembulatan.
	return proRata.kali(bagi)
}

// premiPAPeriodePendek - langkah 5 L858 tanpa diskon:
// round4((TSI × Rate) / 1000) × round4(PctShortPeriod / 100).
func premiPAPeriodePendek(tsi uang.Money, rate, periode rasio) (uang.Money, error) {
	kali, err := rate.kali(tsi)
	if err != nil {
		return uang.Money{}, err
	}
	// Asal: CalculatePremiPA_FacIn langkah 5 L858, kedua @Math.divide presisi 4,
	// masing-masing dibulatkan SEBELUM dikalikan.
	dasar, err := bagiBulat(kali, rate.satuan.Pembagi(), 4)
	if err != nil {
		return uang.Money{}, err
	}
	f, err := periode.faktor(4)
	if err != nil {
		return uang.Money{}, err
	}
	return f.kali(dasar)
}

// premiPAFixrate - langkah 6 L1003 tanpa diskon: round4((TSI × Rate) / 1000) × 1.
func premiPAFixrate(tsi uang.Money, rate rasio) (uang.Money, error) {
	kali, err := rate.kali(tsi)
	if err != nil {
		return uang.Money{}, err
	}
	// Asal: CalculatePremiPA_FacIn langkah 6 L1003, presisi 4. `*1` sesudahnya
	// tidak mengubah nilai maupun skala, sehingga tidak ditulis.
	return bagiBulat(kali, rate.satuan.Pembagi(), 4)
}

// bagiBulat - `@Math.divide(m, pembagi, desimal)`: bagi eksak, lalu bulatkan ke
// `desimal` angka di belakang koma.
//
// Mode: setengah-ke-atas (A37, mengganti butir 14 - dikonfirmasi work owner
// 01-10-2026, butir 56). [terverifikasi] Pemotongan tersingkir (kasus PA nyata .61275755 →
// .6128). [terverifikasi] Setengah-genap tersingkir: dari tiga coverage FIRE nyata
// yang hasil baginya jatuh TEPAT di tengah pada desimal ke-21, dua (kasus EDM) ber-
// digit sebelumnya genap dan tersimpan dibulatkan ke atas (…660|5 → …661,
// …214|5 → …215) - setengah-genap akan memberi …660 dan …214; yang ketiga ber-digit
// ganjil, tidak membedakan mode (TestBerkasKasusP5; register A37 bab Ralat).
// Hanya nilai positif yang terbukti; apd RoundHalfUp = menjauhi nol.
func bagiBulat(m uang.Money, pembagi int64, desimal int32) (uang.Money, error) {
	d, err := bagiBulatDesimal(m.Amount, pembagi, desimal)
	if err != nil {
		return uang.Money{}, err
	}
	return uang.Money{Amount: d, Currency: m.Currency}, nil
}

// bagiBulatDesimal - inti bagiBulat untuk desimal telanjang (dipakai juga rasio).
func bagiBulatDesimal(x *apd.Decimal, pembagi int64, desimal int32) (*apd.Decimal, error) {
	ctx := utils.DecimalContext()
	hasilBagi := new(apd.Decimal)
	if err := periksaEksak(ctx.Quo(hasilBagi, x, apd.New(pembagi, 0))); err != nil {
		return nil, err
	}
	return bulatkan(hasilBagi, desimal)
}

// bulatkan - d ke `desimal` angka di belakang koma, setengah-ke-atas (bagiBulat).
// Konteksnya salinan `utils.DecimalContext()` (presisi tetap 38); hanya mode
// pembulatannya yang ditetapkan di sini.
func bulatkan(d *apd.Decimal, desimal int32) (*apd.Decimal, error) {
	ctx := utils.DecimalContext()
	ctx.Rounding = apd.RoundHalfUp
	hasil := new(apd.Decimal)
	if _, err := ctx.Quantize(hasil, d, -desimal); err != nil {
		return nil, err
	}
	return hasil, nil
}

// periksaEksak membungkus hasil satu operasi apd yang wajib eksak (kali,
// bagi): galat apd diteruskan, dan hasil yang dibulatkan presisi 38 digit
// menjadi ErrMelampauiPresisi.
func periksaEksak(kondisi apd.Condition, err error) error {
	if err != nil {
		return err
	}
	if kondisi.Inexact() {
		return ErrMelampauiPresisi
	}
	return nil
}
