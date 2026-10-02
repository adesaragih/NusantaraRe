package services_test

// Premi coverage FIRE endorsement - `CalculatePremiFire` (+ cabang FIRE
// `SetLocalNonMbuProrate`). Tiket E12. Angka SINTETIS: fixture NB-15 tidak
// menyimpan medan yang dibaca activity ini.
//
// Dibaca sesudah: services/premifire.go.

import (
	"errors"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/modul/endorsmentfacin/backend/services"
)

// premiFire - polis belum berjalan, satu properti, satu item 1.000.000
// (tidak adjustable), satu coverage rate 2 dengan OldCoverage rate 1.
func premiFire(t *testing.T, ubah func(*services.MasukanPremiFire)) services.HasilPremiFire {
	t.Helper()
	h, err := services.HitungPremiFire(masukanFire(t, ubah))
	if err != nil {
		t.Fatalf("HitungPremiFire: %v", err)
	}
	return h
}

func masukanFire(t *testing.T, ubah func(*services.MasukanPremiFire)) services.MasukanPremiFire {
	t.Helper()
	m := services.MasukanPremiFire{
		IsEDM: true, IsFire: true, FlagOnGoingPolicy: "0",
		CurrencyValue: rasio(t, "2"),
		Properti: []services.PropertiFire{{
			Items: []services.ItemTSIFire{{TSIObjectItem: duit(t, "1000000"), PercentageAdjustment: rasio(t, "50")}},
			Coverages: []services.CoverageFire{{
				Rate: rasio(t, "2"), FirstLossScale: rasio(t, "100"), IndemnityPercentage: rasio(t, "100"),
				Old: services.CoverageLamaFire{
					TSI: duit(t, "1000000"), Rate: rasio(t, "1"), ProRatePercent: rasio(t, "100"),
					FirstLossScale: rasio(t, "100"), IndemnityPercentage: rasio(t, "100"), Premium: duit(t, "1000"),
				},
			}},
		}},
		MataUang: mataUangUji, MataUangRupiah: mataUangUji,
	}
	ubah(&m)
	return m
}

// berjalan - polis berjalan; dengan +12 jam pada argumen kedua, kedua selisih
// hari bulat: beg 60 hari, EDMEnd 305 hari.
func berjalan(t *testing.T, m *services.MasukanPremiFire) {
	m.FlagOnGoingPolicy = "1"
	m.StartDateTime = tanggal(t, "2026-01-01 12:00")
	m.EdmDate = tanggal(t, "2026-03-02 00:00")
	m.EndDateTime = tanggal(t, "2026-12-31 12:00")
}

func samaAngka(t *testing.T, label string, got *apd.Decimal, want string) {
	t.Helper()
	w, _, err := apd.NewFromString(want)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.Cmp(w) != 0 {
		t.Errorf("%s: %v, mau %s", label, got, want)
	}
}

func primer(h services.HasilPremiFire) services.CoverageFire { return h.Properti[0].Coverages[0] }

// TestTSIDariJumlahObjek - langkah 6: TSI dan TSISublimit SEMUA coverage
// properti = Σ TSIObjectItem item tak terhapus (TSIAdjustment bila adjustable).
func TestTSIDariJumlahObjek(t *testing.T) {
	h := premiFire(t, func(m *services.MasukanPremiFire) {
		m.Properti[0].Items = []services.ItemTSIFire{
			{TSIObjectItem: duit(t, "1000000"), PercentageAdjustment: rasio(t, "50")},
			{TSIObjectItem: duit(t, "400000"), PercentageAdjustment: rasio(t, "50"), IsAdjustableFlag: "true"},
			{TSIObjectItem: duit(t, "999"), FlagDelete: "1"},
		}
		m.Properti[0].Coverages = append(m.Properti[0].Coverages, services.CoverageFire{})
	})
	for _, c := range h.Properti[0].Coverages {
		samaUang(t, "TSI", c.TSI, duit(t, "1200000"))
		samaUang(t, "TSISublimit", c.TSISublimit, duit(t, "1200000"))
	}
	samaUang(t, "TSIAdjustment", h.Properti[0].Items[1].TSIAdjustment, duit(t, "200000"))
	// 6.2.1 tanpa syarat: item tak-adjustable juga dihitung; faktor kosong → kosong (A12).
	samaUang(t, "TSIAdjustment item 0", h.Properti[0].Items[0].TSIAdjustment, duit(t, "500000"))
	samaUang(t, "TSIAdjustment item 2", h.Properti[0].Items[2].TSIAdjustment, uang.Money{})
}

// TestPremiFireMasukanTidakDiubah - hasil salinan.
func TestPremiFireMasukanTidakDiubah(t *testing.T) {
	m := masukanFire(t, func(*services.MasukanPremiFire) {})
	if _, err := services.HitungPremiFire(m); err != nil {
		t.Fatal(err)
	}
	if !m.Properti[0].Coverages[0].TSI.Kosong() || !m.Properti[0].Coverages[0].Premium.Kosong() {
		t.Error("masukan termutasi")
	}
}

// TestPolisBelumBerjalan - langkah 2 + 7.1: ProRatePercent dari OldCoverage,
// Premium = TSI × Rate × ProRata × FLS × IP / 1e9, PremiRp = idem × kurs.
func TestPolisBelumBerjalan(t *testing.T) {
	c := primer(premiFire(t, func(*services.MasukanPremiFire) {}))
	samaAngka(t, "ProRatePercent", c.ProRatePercent.Value, "100")
	samaUang(t, "Premium", c.Premium, duit(t, "2000")) // 1e6×2×100×100×100/1e9
	samaUang(t, "PremiRp", c.PremiRp, duit(t, "4000"))
}

// TestBukanEDMPakaiTarifPolis - langkah 3.
func TestBukanEDMPakaiTarifPolis(t *testing.T) {
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) { m.IsEDM, m.PctPremiumTariff = false, rasio(t, "50") }))
	samaAngka(t, "ProRatePercent", c.ProRatePercent.Value, "50")
	samaUang(t, "Premium", c.Premium, duit(t, "1000"))
}

// TestBukanEDMPolisBerjalanTanpaPremi - langkah 7 butuh ongoing!=1, langkah
// 8/10 butuh IsEDM: keduanya tidak jalan.
func TestBukanEDMPolisBerjalanTanpaPremi(t *testing.T) {
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) { berjalan(t, m); m.IsEDM = false }))
	if !c.Premium.Kosong() || !c.PremiRp.Kosong() {
		t.Errorf("Premium %v PremiRp %v", c.Premium, c.PremiRp)
	}
}

// TestCoverageDihapusPremiNol - 7.2.
func TestCoverageDihapusPremiNol(t *testing.T) {
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) { m.Properti[0].Coverages[0].FlagDelete = "1" }))
	samaUang(t, "Premium", c.Premium, duit(t, "0"))
	samaUang(t, "PremiRp", c.PremiRp, duit(t, "0"))
}

// TestPembulatanEmpatDesimalSetengahKeAtas - @Math.divide(…,1e9,4) HALF_UP
// (A37). 0.00125 membedakan HALF_UP (0.0013) dari HALF_EVEN/DOWN (0.0012).
func TestPembulatanEmpatDesimalSetengahKeAtas(t *testing.T) {
	for rate, mau := range map[string]string{"1.23456": "0.0012", "1.25": "0.0013", "1.24999": "0.0012"} {
		c := primer(premiFire(t, func(m *services.MasukanPremiFire) {
			m.Properti[0].Items = []services.ItemTSIFire{{TSIObjectItem: duit(t, "1")}}
			m.Properti[0].Coverages[0].Rate = rasio(t, rate)
		}))
		samaUang(t, "Premium rate "+rate, c.Premium, duit(t, mau))
	}
}

// TestBawaanSeratusHanyaDiPremium - K-046: 7.1 PremiRp memakai `.FirstLossScale`
// TANPA `@if`; FLS kosong → galat. IP kosong → 100 di keduanya.
func TestBawaanSeratusHanyaDiPremium(t *testing.T) {
	_, err := services.HitungPremiFire(masukanFire(t, func(m *services.MasukanPremiFire) {
		m.Properti[0].Coverages[0].FirstLossScale = uang.Ratio{}
	}))
	if !errors.Is(err, services.ErrNilaiPremiKosong) {
		t.Errorf("galat %v, mau ErrNilaiPremiKosong", err)
	}
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) {
		m.Properti[0].Coverages[0].IndemnityPercentage = uang.Ratio{}
		m.CurrencyValue = uang.Ratio{} // @if(…=="",1,…)
	}))
	samaUang(t, "Premium", c.Premium, duit(t, "2000"))
	samaUang(t, "PremiRp", c.PremiRp, duit(t, "2000"))
}

// TestPolisBerjalanDuaBagian - 8 + 10.1: porsi hari ÷ 365 (6 desimal) × 100,
// Premium = bagian baru (porsi sisa) + bagian lama (porsi berjalan). PremiRp
// memakai porsi TERTUKAR (K-046).
func TestPolisBerjalanDuaBagian(t *testing.T) {
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) { berjalan(t, m) }))
	samaAngka(t, "StartEDM", c.ProratePercentStartEDM.Value, "16.4384") // 60/365 = 0.164384
	samaAngka(t, "EDMEnd", c.ProratePercentEDMEnd.Value, "83.5616")     // 305/365 = 0.835616
	samaAngka(t, "ProRatePercent", c.ProRatePercent.Value, "100")
	if c.CalculateMethod != "1" {
		t.Errorf("CalculateMethod %q", c.CalculateMethod)
	}
	// baru 1e6×2×83.5616×1e4/1e9 = 1671.232; lama 1e6×1×16.4384×1e4/1e9 = 164.384
	samaUang(t, "Premium", c.Premium, duit(t, "1835.616"))
	// tertukar: 1e6×2×16.4384… = 328.768 + 1e6×1×83.5616… = 835.616 → 1164.384 × 2
	samaUang(t, "PremiRp", c.PremiRp, duit(t, "2328.768"))
}

// TestPolisBerjalanBukanFire - SetLocalNonMbuProrate langkah 1 saja:
// ProRatePercent = 0; porsi EDM tetap nilai lama coverage.
func TestPolisBerjalanBukanFire(t *testing.T) {
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) {
		berjalan(t, m)
		m.IsFire = false
		m.Properti[0].Coverages[0].ProratePercentStartEDM = rasio(t, "50")
		m.Properti[0].Coverages[0].ProratePercentEDMEnd = rasio(t, "50")
	}))
	samaAngka(t, "ProRatePercent", c.ProRatePercent.Value, "0")
	// baru 1e6×2×50×1e4/1e9 = 1000; lama 1e6×1×50×1e4/1e9 = 500
	samaUang(t, "Premium", c.Premium, duit(t, "1500"))
}

// TestPolisBerjalanCoverageDihapus - 10.2: Old.Premium − (Old.Premium × porsi
// sisa, 4 desimal) ÷ Old.ProRatePercent.
func TestPolisBerjalanCoverageDihapus(t *testing.T) {
	c := primer(premiFire(t, func(m *services.MasukanPremiFire) {
		berjalan(t, m)
		m.Properti[0].Coverages[0].FlagDelete = "1"
	}))
	samaUang(t, "Premium", c.Premium, duit(t, "164.384")) // 1000 − 83561.6/100
	samaUang(t, "PremiRp", c.PremiRp, duit(t, "328.768"))
}

// TestPolisBerjalanHariTakBulatDitolak - jam sama di kedua tanggal + 12 jam →
// n,5 hari; `@DateTimeDifference(…,D)` atas pecahan belum terverifikasi.
func TestPolisBerjalanHariTakBulatDitolak(t *testing.T) {
	_, err := services.HitungPremiFire(masukanFire(t, func(m *services.MasukanPremiFire) {
		berjalan(t, m)
		m.StartDateTime, m.EdmDate, m.EndDateTime = tanggal(t, "2026-01-01 05:00"), tanggal(t, "2026-03-02 05:00"), tanggal(t, "2026-12-31 05:00")
	}))
	if !errors.Is(err, services.ErrSatuanSelisihWaktuBelumTerverifikasi) {
		t.Errorf("galat %v", err)
	}
}

// TestPolisBerjalanTanggalKosong - tanggal tak terisi → galat.
func TestPolisBerjalanTanggalKosong(t *testing.T) {
	_, err := services.HitungPremiFire(masukanFire(t, func(m *services.MasukanPremiFire) {
		berjalan(t, m)
		m.EdmDate = time.Time{}
	}))
	if !errors.Is(err, services.ErrTanggalPorsiPeriodeKosong) {
		t.Errorf("galat %v", err)
	}
}

// TestCoveragePrimerDiLuarDaftar - indeks salah → galat.
func TestCoveragePrimerDiLuarDaftar(t *testing.T) {
	_, err := services.HitungPremiFire(masukanFire(t, func(m *services.MasukanPremiFire) { m.IndeksCoverage = 3 }))
	if !errors.Is(err, services.ErrCoveragePrimerTidakAda) {
		t.Errorf("galat %v", err)
	}
}
