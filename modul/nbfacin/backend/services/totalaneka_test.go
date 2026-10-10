package services

import (
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

func totalTeks(d []models.TotalCoverage) [][4]string {
	var h [][4]string
	for _, t := range d {
		h = append(h, [4]string{t.Currency, teksD(t.TSI), teksD(t.Premium), teksD(t.Rate)})
	}
	return h
}

func samaTotal(t *testing.T, nama string, got []models.TotalCoverage, mau [][4]string) {
	t.Helper()
	g := totalTeks(got)
	if len(g) != len(mau) {
		t.Errorf("%s: %v, mau %v", nama, g, mau)
		return
	}
	for i := range mau {
		beda := g[i][0] != mau[i][0]
		for j := 1; j < 4; j++ {
			beda = beda || desimalWajib(g[i][j]).Cmp(desimalWajib(mau[i][j])) != 0
		}
		if beda {
			t.Errorf("%s baris %d: %v, mau %v", nama, i, g[i], mau[i])
		}
	}
}

func covPremi(mataUang, tsi, premi string) models.CoverageAneka {
	return models.CoverageAneka{Currency: mataUang, TSI: ds(tsi), Premium: ds(premi)}
}

// TestTotalAnekaLokasi - SumTotalTSIPremiGross_Act langkah 2 (Golf) / 3a / 3b: per lokasi x mata uang (urutan kemunculan), premi
// coverage mentah, rate = premi/TSI x 100, TSI 0 -> rate 0; tanpa gerbang -> kosong.
func TestTotalAnekaLokasi(t *testing.T) {
	objekUSD := models.ObjekAneka{Currency: "USD", TSI: ds("1000"), Coverages: []models.CoverageAneka{
		covPremi("IDR", "400", "10"), covPremi("USD", "600", "5")}}
	objekIDR := models.ObjekAneka{Currency: "IDR", TSI: ds("2000"), Coverages: []models.CoverageAneka{covPremi("IDR", "2000", "20")}}
	objekNol := models.ObjekAneka{Currency: "EUR", Coverages: []models.CoverageAneka{covPremi("EUR", "", "7")}}
	golf := models.LokasiAneka{ObjekGolf: []models.ObjekAneka{objekUSD, objekIDR, objekNol}}
	samaTotal(t, "Golf", totalAnekaLokasi("K", kasusBisnis{businessType: "GolfInsurance", businessCode: "10999"}, golf),
		[][4]string{{"USD", "1000", "15", "1.5"}, {"IDR", "2000", "20", "1"}, {"EUR", "0", "7", "0"}})
	// Aneka biasa (Liability): okupasi dilebur; premi ke mata uang OBJEK (3a). Objek Golf diabaikan.
	lok := models.LokasiAneka{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{objekUSD}}, {Objek: []models.ObjekAneka{objekIDR}}},
		ObjekGolf: []models.ObjekAneka{objekIDR}}
	samaTotal(t, "Liability", totalAnekaLokasi("K", kasusBisnis{businessType: "Liability", businessCode: "10999"}, lok),
		[][4]string{{"USD", "1000", "15", "1.5"}, {"IDR", "2000", "20", "1"}})
	// Tipe khusus (HE, 3b): TSI dan premi per COVERAGE, mata uang coverage.
	samaTotal(t, "HE", totalAnekaLokasi("K", kasusBisnis{businessType: "HE", businessCode: "10999"}, lok),
		[][4]string{{"IDR", "2400", "30", "1.25"}, {"USD", "600", "5", "0.83333333"}})
	if got := totalAnekaLokasi("K", kasusBisnis{businessType: "Fire", businessCode: "10999"}, lok); len(got) != 0 {
		t.Errorf("Fire: %v", totalTeks(got))
	}
}
