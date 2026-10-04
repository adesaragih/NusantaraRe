package premium

import (
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
)

// fire - coverage FIRE: rate ‰, pro-rata/indemnity/loss limit %. Bawaan: TSI 1 M,
// Rate 1,5 ‰, ProRatePercent 100, IndemnityPercentage 100, LossLimit 100, basis 1.
func fire(ubah func(*Input)) Input {
	in := Input{LiniBisnis: LiniFire, MataUang: "IDR", TSI: "1000000000", Rate: "1.5", ProRatePercent: "100",
		IndemnityPercentage: "100", LossLimit: "100", CoverageBasis: "1", PctAdjustment: "0"}
	if ubah != nil {
		ubah(&in)
	}
	return in
}

// aneka - coverage ANEKA: rate %. Bawaan: TSI 1 jt, Rate 2 %, metode 3 (fix rate).
func aneka(lini LiniBisnis, ubah func(*Input)) Input {
	in := Input{LiniBisnis: lini, MataUang: "IDR", TSI: "1000000", Rate: "2", CalculateMethod: "3", ProRatePercent: "100"}
	if ubah != nil {
		ubah(&in)
	}
	return in
}

func premiTeks(t *testing.T, in Input) string {
	t.Helper()
	got, err := Calculate(in)
	if err != nil {
		t.Fatalf("%+v: %v", in, err)
	}
	return utils.FormatDecimal(got.Amount)
}

// TestPremiFire - CountPremi_ACT jalur percent. Pembagi = rate ‰ (1000) × pro-rata ×
// indemnity × loss limit (masing-masing 100) = 10⁹; adjustable ×100 → 10¹¹;
// FirstScale ×100. Presisi 20.
func TestPremiFire(t *testing.T) {
	for _, u := range []struct {
		nama string
		in   Input
		mau  string
	}{
		// 1e9 × 1,5 × 100 × 100 × 100 / 1e9 = 1.500.000 (langkah 19).
		{"basis 1", fire(nil), "1500000.00000000000000000000"},
		// × PctAdjustment 80 / 1e11 (langkah 21 menimpa 19).
		{"basis 1 adjustable", fire(func(i *Input) { i.PctAdjustment = "80" }), "1200000.00000000000000000000"},
		// NetRate 1,2 menggantikan Rate (langkah 23), lalu adjustable (25).
		{"basis 1 net rate", fire(func(i *Input) { i.NetRate = "1.2" }), "1200000.00000000000000000000"},
		{"basis 1 net rate adjustable", fire(func(i *Input) { i.NetRate = "1.2"; i.PctAdjustment = "80" }), "960000.00000000000000000000"},
		// NetRate "0" = tidak dipakai (`.NetRate!="0"&&.NetRate!=""`).
		{"net rate 0", fire(func(i *Input) { i.NetRate = "0" }), "1500000.00000000000000000000"},
		// First Loss: × FirstScale 50 / 1e11 (langkah 27); adjustable / 1e13 (29).
		{"basis 2", fire(func(i *Input) { i.CoverageBasis = "2"; i.FirstScale = "50" }), "750000.00000000000000000000"},
		{"basis 2 adjustable", fire(func(i *Input) { i.CoverageBasis = "2"; i.FirstScale = "50"; i.PctAdjustment = "80" }), "600000.00000000000000000000"},
		// EML/PML dan Sublimit: rumus premi sama dengan basis 1 (langkah 35, 43).
		{"basis 3", fire(func(i *Input) { i.CoverageBasis = "3" }), "1500000.00000000000000000000"},
		{"basis 4", fire(func(i *Input) { i.CoverageBasis = "4" }), "1500000.00000000000000000000"},
		// LossLimit kosong atau 0 = 100 (langkah 17); 50 → separuh.
		{"loss limit kosong", fire(func(i *Input) { i.LossLimit = "" }), "1500000.00000000000000000000"},
		{"loss limit teks 0", fire(func(i *Input) { i.LossLimit = "0" }), "1500000.00000000000000000000"},
		{"loss limit 50", fire(func(i *Input) { i.LossLimit = "50" }), "750000.00000000000000000000"},
		// Pro-rata dan indemnity separuh.
		{"pro-rata 50 indemnity 50", fire(func(i *Input) { i.ProRatePercent = "50"; i.IndemnityPercentage = "50" }), "375000.00000000000000000000"},
		// 1 × 5e-18 × 1e6 / 1e9 = 5e-21: tepat di tengah pada desimal ke-20 → ke atas (A37).
		{"seri ke atas", fire(func(i *Input) { i.TSI = "1"; i.Rate = "0.000000000000000005" }), "0.00000000000000000001"},
	} {
		if got := premiTeks(t, u.in); got != u.mau {
			t.Errorf("%s: %s, mau %s", u.nama, got, u.mau)
		}
	}
}

// TestPremiFireDitolak - basis 5 (Layering, butir 30), basis kosong/tak dikenal,
// medan wajib kosong, dan nol berbentuk teks lain dari "0" pada pembanding teks.
func TestPremiFireDitolak(t *testing.T) {
	for _, u := range []struct {
		nama string
		in   Input
		mau  error
	}{
		{"basis 5 Layering", fire(func(i *Input) { i.CoverageBasis = "5" }), ErrBentukBelumDiport},
		{"basis tak dikenal", fire(func(i *Input) { i.CoverageBasis = "9" }), ErrCoverageBasisTakDikenal},
		{"basis kosong", fire(func(i *Input) { i.CoverageBasis = "" }), ErrCoverageBasisTakDikenal},
		// `.NetRate!="0"` dan `Local.LossLimit=="0"` dibandingkan sebagai TEKS (A43).
		{"net rate 0.0", fire(func(i *Input) { i.NetRate = "0.0" }), ErrTeksNolAmbigu},
		{"loss limit 0.0", fire(func(i *Input) { i.LossLimit = "0.0" }), ErrTeksNolAmbigu},
		{"indemnity kosong", fire(func(i *Input) { i.IndemnityPercentage = "" }), ErrRasioKosong},
		{"first scale kosong", fire(func(i *Input) { i.CoverageBasis = "2" }), ErrRasioKosong},
	} {
		if _, err := Calculate(u.in); !errors.Is(err, u.mau) {
			t.Errorf("%s: galat %v, mau %v", u.nama, err, u.mau)
		}
	}
}

// TestPremiAnekaGolf - CountPremiCoverageAneka / FillPremiGolf (rumus identik):
// round4(TSI×Rate×P/1e4) + round4(TSI×Rate×P×Loading/1e6), lalu × LossLimit / 100.
func TestPremiAnekaGolf(t *testing.T) {
	for _, lini := range []LiniBisnis{LiniAneka, LiniGolf} {
		for _, u := range []struct {
			nama string
			in   Input
			mau  string
		}{
			// Metode 3: P = 100 → 1e6 × 2 × 100 / 1e4 = 20.000.
			{"fix rate", aneka(lini, nil), "20000"},
			// Metode 1: P = ProRatePercent 50; loading 10 → 10.000 + 1.000.
			{"pro rata + loading", aneka(lini, func(i *Input) { i.CalculateMethod = "1"; i.ProRatePercent = "50"; i.Loading = "10" }), "11000"},
			// Metode 2: P = PctShortPeriod 25.
			{"short period", aneka(lini, func(i *Input) { i.CalculateMethod = "2"; i.PctShortPeriod = "25" }), "5000"},
			// LossLimit 75: 11.000 × 75 / 100.
			{"loss limit", aneka(lini, func(i *Input) {
				i.CalculateMethod = "1"
				i.ProRatePercent = "50"
				i.Loading = "10"
				i.LossLimit = "75"
			}), "8250"},
			// MBD + indemnity 50: round4(1e6×2×100×50/1e6) + round4(…×10/1e8) = 10.000 + 1.000.
			{"MBD indemnity", aneka(lini, func(i *Input) { i.MBD = true; i.IndemnityPercentage = "50"; i.Loading = "10" }), "11000"},
			// MBD tanpa indemnity → rumus biasa (langkah 17.2 / 12.2 menimpa).
			{"MBD tanpa indemnity", aneka(lini, func(i *Input) { i.MBD = true }), "20000"},
			// Bukan MBD: indemnity terisi pun tidak dipakai.
			{"indemnity tanpa MBD", aneka(lini, func(i *Input) { i.IndemnityPercentage = "50" }), "20000"},
			// `Local.Losslimit==0` NUMERIK di ANEKA/GOLF: "0.0" = bawaan 100.
			{"loss limit 0.0", aneka(lini, func(i *Input) { i.LossLimit = "0.0" }), "20000"},
			// round4 di DALAM: 1 × 0,00123 × 100 / 1e4 = 0,0000123 → 0.
			{"pembulatan 4", aneka(lini, func(i *Input) { i.TSI = "1"; i.Rate = "0.00123" }), "0"},
		} {
			got := premiTeks(t, u.in)
			if !samaNilai(t, got, u.mau) {
				t.Errorf("%s %s: %s, mau %s", lini, u.nama, got, u.mau)
			}
		}
	}
	// 1 × 0,005 × 100 / 1e4 = 0,00005: seri di desimal ke-4 → ke atas (A37).
	if got := premiTeks(t, aneka(LiniAneka, func(i *Input) { i.TSI = "1"; i.Rate = "0.005" })); got != "0.0001" {
		t.Errorf("seri: %s, mau 0.0001", got)
	}
	// ANEKA langkah 14.1: pro-rata 0 atau kosong dihitung ulang dari tanggal (belum
	// diport). GOLF tidak punya langkah itu: pro-rata 0 = premi 0.
	for _, pr := range []string{"0", ""} {
		if _, err := Calculate(aneka(LiniAneka, func(i *Input) { i.CalculateMethod = "1"; i.ProRatePercent = pr })); !errors.Is(err, ErrBentukBelumDiport) {
			t.Errorf("ANEKA pro-rata %q: galat %v", pr, err)
		}
	}
	if got := premiTeks(t, aneka(LiniGolf, func(i *Input) { i.CalculateMethod = "1"; i.ProRatePercent = "0" })); !samaNilai(t, got, "0") {
		t.Errorf("GOLF pro-rata 0: %s", got)
	}
	if _, err := Calculate(aneka(LiniAneka, func(i *Input) { i.CalculateMethod = "4" })); !errors.Is(err, ErrMetodeTakDikenal) {
		t.Errorf("metode 4: galat %v", err)
	}
}

// TestPremiMarine - CountGPWMarinePAMbu_Act 1.1.1.1.1-2: round4(TSI × Rate / 100);
// master policy = 0.
func TestPremiMarine(t *testing.T) {
	m := Input{LiniBisnis: LiniMarineCargo, MataUang: "IDR", TSI: "1000000", Rate: "0.125"}
	if got := premiTeks(t, m); got != "1250.0000" {
		t.Errorf("marine: %s", got)
	}
	m.MasterPolicy = true
	if got := premiTeks(t, m); !samaNilai(t, got, "0") {
		t.Errorf("master policy: %s", got)
	}
	if got := premiTeks(t, Input{LiniBisnis: LiniMarineCargo, MataUang: "IDR", TSI: "1", Rate: "0.005"}); got != "0.0001" {
		t.Errorf("seri: %s, mau 0.0001 (A37)", got)
	}
}

// TestPremiBondingLewatAneka - kasus Bonding dihitung lewat lini ANEKA (A11);
// LiniBonding sendiri tidak punya rumus.
func TestPremiBondingLewatAneka(t *testing.T) {
	if _, err := Calculate(aneka(LiniBonding, nil)); !errors.Is(err, ErrBentukBelumDiport) {
		t.Errorf("BONDING: galat %v", err)
	}
}

// TestLossLimitSkalaPembilang - `.Premium*Local.Losslimit/100` mempertahankan skala
// pembilang (kasus nyata Kredit: `834892000.0000`), diperpanjang hanya bila perlu.
func TestLossLimitSkalaPembilang(t *testing.T) {
	for _, u := range []struct{ premi, loss, mau string }{
		{"834892000.0000", "100", "834892000.0000"},
		{"11000.0000", "75", "8250.0000"},
		{"0.0001", "50", "0.00005"},
		{"0.0001", "100", "0.0001"},
	} {
		m := uang.Money{Amount: desimalUji(t, u.premi), Currency: "IDR"}
		got, err := kaliLossLimitPersen(m, desimalUji(t, u.loss))
		if err != nil || utils.FormatDecimal(got.Amount) != u.mau {
			t.Errorf("%s × %s: dapat %v (%v), mau %s", u.premi, u.loss, got.Amount, err, u.mau)
		}
	}
}

func desimalUji(t *testing.T, s string) *apd.Decimal {
	t.Helper()
	d, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
