package services

// `Activity/TreatyInNonSetTotal.xml` — SATU lintasan yang menjumlahkan
// kelima page list kontrak Non-Prop sekaligus.
//
// ---------------------------------------------------------------------
// ⭐ MENGAPA RUMUS INI YANG DIDAHULUKAN
// ---------------------------------------------------------------------
// Ia dipanggil SEPULUH seksi — `TreatyInActualLimits`, `TreatyInActualShare`,
// `TreatyInActualSumary`, `TreatyInFacultativeShareCalculation`,
// `TreatyInActualFacultativeShareCalculation`,
// `TreatyInTabsNPValueDifferenceProRate`, dan kawan-kawannya — lebih banyak
// daripada rumus mana pun yang belum dibangun.
//
// ---------------------------------------------------------------------
// LANGKAHNYA, APA ADANYA
// ---------------------------------------------------------------------
//
//	1  Property-Remove kesebelas total
//	2  .Currency = Retention(1).Currency ; .CurrencyEgnpiAmount = "IDR"
//	3  per Retention:    TotalRetentionAmount += .Amount
//	4  per EGNPI:        TotalEgnpiAmount += .AmountIDR
//	5  per EGNPI:        .Proportion = .AmountIDR / TotalEgnpiAmount * 100
//	                     TotalEgnpiProportion += .Proportion
//	6  per Limits:       TotalLimitsIOOLimit    += .Limit
//	                     TotalLimitsAdjPct      += .AdjRate
//	                     TotalLimitsDeductible  += .Deductible
//	                     TotalLimitsPremiumEarned += .PremiumEarned
//	                     TotalLimitsMdp         += .MDP
//	                     TotalLimitsROL = TotalLimitsPremiumEarned /
//	                                      TotalLimitsIOOLimit * 100
//	7  per Share:        TotalShareRnmLimit += .NusareLimit
//	                     TotalShareGross    += .GrossPremium
//	                     TotalShareNet      += .NetPremium
//	8  per Installment:  TotalInstallmentAmount += .Amount
//	                     TotalInstallmentPct    += .InstallmentPct
//
// ---------------------------------------------------------------------
// ⚠️ DISALIN APA ADANYA WALAU JANGGAL — rumus, bukan alamat
// ---------------------------------------------------------------------
//   - Langkah 4 dan 5 adalah DUA perulangan atas larik yang SAMA, dan
//     urutannya yang membuat pembagian di langkah 5 memakai total PENUH.
//     Menggabungkannya menjadi satu perulangan akan membagi dengan total
//     yang masih tumbuh, dan proporsi baris pertama menjadi 100%.
//   - `TotalLimitsROL` dihitung ULANG di dalam perulangan langkah 6, bukan
//     sesudahnya. Hasil akhirnya sama (ia memakai kedua total berjalan pada
//     putaran terakhir) — dan bentuknya dipertahankan supaya yang membaca
//     ekspor menemukan hal yang sama di sini.
//   - `.Currency` kontrak diambil dari baris PERTAMA Retention. Bila
//     Retention kosong, ia menjadi kosong — bukan mata uang kontrak.
//
// ⛔ MURNI: nol baca, nol tulis basis data.

import "github.com/cockroachdb/apd/v3"

// AksiNonSetTotal - nama aksi rute.
const AksiNonSetTotal = "non-set-total"

// BarisNilai - satu baris page list yang hanya dibaca nilainya.
//
// ⭐ Satu bentuk untuk kelima larik: yang berbeda hanya NAMA medan yang
// dibaca, dan nama itu dipilih di tempat penjumlahannya.
type BarisNilai = map[string]string

// MasukanNonSetTotal - kelima page list yang dijumlahkan.
type MasukanNonSetTotal struct {
	Retention   []BarisNilai `json:"Retention"`
	EGNPI       []BarisNilai `json:"EGNPI"`
	Limits      []BarisNilai `json:"Limits"`
	Share       []BarisNilai `json:"Share"`
	Installment []BarisNilai `json:"Installment"`
}

// HasilNonSetTotal - kesebelas total, berikut kedua medan mata uang dan
// `Proportion` tiap baris EGNPI yang langkah 5 tulis balik.
type HasilNonSetTotal struct {
	Currency            string `json:"Currency"`
	CurrencyEgnpiAmount string `json:"CurrencyEgnpiAmount"`

	TotalRetentionAmount string `json:"TotalRetentionAmount"`

	TotalEgnpiAmount     string `json:"TotalEgnpiAmount"`
	TotalEgnpiProportion string `json:"TotalEgnpiProportion"`
	// ProporsiEgnpi - `.Proportion` tiap baris EGNPI, SEJAJAR urutan masukan.
	ProporsiEgnpi []string `json:"proporsiEgnpi"`

	TotalLimitsIOOLimit      string `json:"TotalLimitsIOOLimit"`
	TotalLimitsAdjPct        string `json:"TotalLimitsAdjPct"`
	TotalLimitsDeductible    string `json:"TotalLimitsDeductible"`
	TotalLimitsPremiumEarned string `json:"TotalLimitsPremiumEarned"`
	TotalLimitsMdp           string `json:"TotalLimitsMdp"`
	TotalLimitsROL           string `json:"TotalLimitsROL"`

	TotalShareRnmLimit string `json:"TotalShareRnmLimit"`
	TotalShareGross    string `json:"TotalShareGross"`
	TotalShareNet      string `json:"TotalShareNet"`

	TotalInstallmentAmount string `json:"TotalInstallmentAmount"`
	TotalInstallmentPct    string `json:"TotalInstallmentPct"`
}

// jumlahMedan - SIGMA satu medan atas satu larik.
func jumlahMedan(baris []BarisNilai, medan string) *apd.Decimal {
	t := apd.New(0, 0)
	for _, b := range baris {
		t = tambah(t, angka(b[medan]))
	}
	return t
}

// persenDari - `a / b * 100`, dan NOL bila pembaginya nol.
//
// ⚠️ Pega menghasilkan nol untuk pembagian ini, bukan galat: `@divide` atas
// properti kosong memberi kosong, dan kosong dibaca nol oleh perkalian
// berikutnya. Mengembalikan galat di sini akan menghentikan layar pada
// kontrak yang sah-sah saja — kontrak tanpa EGNPI, misalnya.
func persenDari(a, b *apd.Decimal) *apd.Decimal {
	if b.IsZero() {
		return apd.New(0, 0)
	}
	return kali(bagiBulatDes(a, b, 20), apd.New(100, 0))
}

// HitungNonSetTotal - `TreatyInNonSetTotal`, langkah demi langkah.
func HitungNonSetTotal(m MasukanNonSetTotal) HasilNonSetTotal {
	var h HasilNonSetTotal

	// [2] — `.Currency` dari baris PERTAMA Retention; `CurrencyEgnpiAmount`
	// tetapan `"IDR"`.
	if len(m.Retention) > 0 {
		h.Currency = m.Retention[0]["Currency"]
	}
	h.CurrencyEgnpiAmount = "IDR"

	// [3]
	h.TotalRetentionAmount = teks(jumlahMedan(m.Retention, "Amount"))

	// [4] — SELURUH larik dulu; langkah [5] membaginya dengan total PENUH.
	totalEgnpi := jumlahMedan(m.EGNPI, "AmountIDR")
	h.TotalEgnpiAmount = teks(totalEgnpi)

	// [5]
	proporsi := apd.New(0, 0)
	h.ProporsiEgnpi = make([]string, 0, len(m.EGNPI))
	for _, b := range m.EGNPI {
		p := persenDari(angka(b["AmountIDR"]), totalEgnpi)
		h.ProporsiEgnpi = append(h.ProporsiEgnpi, teks(p))
		proporsi = tambah(proporsi, p)
	}
	h.TotalEgnpiProportion = teks(proporsi)

	// [6]
	ioo := jumlahMedan(m.Limits, "Limit")
	pe := jumlahMedan(m.Limits, "PremiumEarned")
	h.TotalLimitsIOOLimit = teks(ioo)
	h.TotalLimitsAdjPct = teks(jumlahMedan(m.Limits, "AdjRate"))
	h.TotalLimitsDeductible = teks(jumlahMedan(m.Limits, "Deductible"))
	h.TotalLimitsPremiumEarned = teks(pe)
	h.TotalLimitsMdp = teks(jumlahMedan(m.Limits, "MDP"))
	h.TotalLimitsROL = teks(persenDari(pe, ioo))

	// [7]
	h.TotalShareRnmLimit = teks(jumlahMedan(m.Share, "NusareLimit"))
	h.TotalShareGross = teks(jumlahMedan(m.Share, "GrossPremium"))
	h.TotalShareNet = teks(jumlahMedan(m.Share, "NetPremium"))

	// [8]
	h.TotalInstallmentAmount = teks(jumlahMedan(m.Installment, "Amount"))
	h.TotalInstallmentPct = teks(jumlahMedan(m.Installment, "InstallmentPct"))

	return h
}
