package services_test

// Uji `TreatyInNonSetTotal` — dibangun 8 Oktober 2026 atas permintaan
// pemilik proses *"utamakan rumus agar berjalan semua seperti di pega"*.
//
// ⭐ Ia didahulukan sebab SEPULUH seksi memanggilnya — lebih banyak daripada
// rumus mana pun yang belum dibangun.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func masukanNST() services.MasukanNonSetTotal {
	return services.MasukanNonSetTotal{
		Retention: []services.BarisNilai{
			{"Currency": "IDR", "Amount": "1000"},
			{"Currency": "USD", "Amount": "250"},
		},
		EGNPI: []services.BarisNilai{
			{"AmountIDR": "300"},
			{"AmountIDR": "700"},
		},
		Limits: []services.BarisNilai{
			{"Limit": "2000", "AdjRate": "1.5", "Deductible": "100", "PremiumEarned": "400", "MDP": "7"},
			{"Limit": "3000", "AdjRate": "2.5", "Deductible": "150", "PremiumEarned": "600", "MDP": "3"},
		},
		Share: []services.BarisNilai{
			{"NusareLimit": "90", "GrossPremium": "40", "NetPremium": "35"},
			{"NusareLimit": "10", "GrossPremium": "60", "NetPremium": "25"},
		},
		Installment: []services.BarisNilai{
			{"Amount": "500", "InstallmentPct": "50"},
			{"Amount": "500", "InstallmentPct": "50"},
		},
	}
}

func TestNonSetTotalMenjumlahKelimaLarik(t *testing.T) {
	h := services.HitungNonSetTotal(masukanNST())
	for _, u := range []struct{ nama, dapat, mau string }{
		{"TotalRetentionAmount", h.TotalRetentionAmount, "1250"},
		{"TotalEgnpiAmount", h.TotalEgnpiAmount, "1000"},
		{"TotalLimitsIOOLimit", h.TotalLimitsIOOLimit, "5000"},
		// ⛔ `4`, bukan `4.0`: `teks()` menormalkan ekor nol. Harapan
		// pertama saya keliru di sini, bukan kodenya.
		{"TotalLimitsAdjPct", h.TotalLimitsAdjPct, "4"},
		{"TotalLimitsDeductible", h.TotalLimitsDeductible, "250"},
		{"TotalLimitsPremiumEarned", h.TotalLimitsPremiumEarned, "1000"},
		{"TotalLimitsMdp", h.TotalLimitsMdp, "10"},
		{"TotalShareRnmLimit", h.TotalShareRnmLimit, "100"},
		{"TotalShareGross", h.TotalShareGross, "100"},
		{"TotalShareNet", h.TotalShareNet, "60"},
		{"TotalInstallmentAmount", h.TotalInstallmentAmount, "1000"},
		{"TotalInstallmentPct", h.TotalInstallmentPct, "100"},
	} {
		if u.dapat != u.mau {
			t.Errorf("%s = %q, mau %q", u.nama, u.dapat, u.mau)
		}
	}
}

// ⛔ LANGKAH 4 DAN 5 ADALAH DUA PERULANGAN, dan urutannya yang membuat
// pembagiannya benar.
//
// ⚠️ Bila keduanya digabung menjadi satu perulangan, baris PERTAMA dibagi
// dengan total yang baru memuat dirinya sendiri — dan proporsinya menjadi
// 100%. Uji ini memakai dua baris tak sama besar supaya kekeliruan itu
// terlihat: 30% dan 70%, bukan 100% dan 70%.
func TestNonSetTotalProporsiMemakaiTotalPENUH(t *testing.T) {
	h := services.HitungNonSetTotal(masukanNST())
	if len(h.ProporsiEgnpi) != 2 {
		t.Fatalf("proporsi %d baris, mau 2", len(h.ProporsiEgnpi))
	}
	if h.ProporsiEgnpi[0] != "30" || h.ProporsiEgnpi[1] != "70" {
		t.Errorf("proporsi %v, mau [30 70]", h.ProporsiEgnpi)
	}
	if h.TotalEgnpiProportion != "100" {
		t.Errorf("TotalEgnpiProportion = %q, mau 100", h.TotalEgnpiProportion)
	}
}

// ⭐ `TotalLimitsROL` = PremiumEarned / IOOLimit × 100 = 1000/5000 × 100.
func TestNonSetTotalROL(t *testing.T) {
	if got := services.HitungNonSetTotal(masukanNST()).TotalLimitsROL; got != "20" {
		t.Errorf("TotalLimitsROL = %q, mau 20", got)
	}
}

// ⛔ `.Currency` kontrak dari baris PERTAMA Retention — bukan dari baris
// mana pun yang kebetulan ada, dan bukan mata uang kepala kontrak.
func TestNonSetTotalCurrencyDariRetensiPertama(t *testing.T) {
	h := services.HitungNonSetTotal(masukanNST())
	if h.Currency != "IDR" {
		t.Errorf("Currency = %q, mau IDR", h.Currency)
	}
	if h.CurrencyEgnpiAmount != "IDR" {
		t.Errorf("CurrencyEgnpiAmount = %q, mau IDR", h.CurrencyEgnpiAmount)
	}
}

// ⛔ PEMBAGIAN NOL MEMBERI NOL, BUKAN GALAT.
//
// ⚠️ Kontrak tanpa EGNPI dan tanpa Limits adalah kontrak yang sah. Rumus
// yang bergalat di sini akan menghentikan layar pada kontrak yang tidak
// punya cacat apa pun.
func TestNonSetTotalLarikKosongTidakBergalat(t *testing.T) {
	h := services.HitungNonSetTotal(services.MasukanNonSetTotal{})
	for _, u := range []struct{ nama, dapat string }{
		{"TotalEgnpiProportion", h.TotalEgnpiProportion},
		{"TotalLimitsROL", h.TotalLimitsROL},
		{"TotalRetentionAmount", h.TotalRetentionAmount},
	} {
		if u.dapat != "0" {
			t.Errorf("%s = %q, mau 0", u.nama, u.dapat)
		}
	}
	// ⭐ Dan `Currency` KOSONG, bukan dikarang: Retention tidak berbaris.
	if h.Currency != "" {
		t.Errorf("Currency = %q, mau kosong", h.Currency)
	}
	if len(h.ProporsiEgnpi) != 0 {
		t.Errorf("proporsi %v, mau kosong", h.ProporsiEgnpi)
	}
}
