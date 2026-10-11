package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// cargoSpreadUji - dua cargo IDR (cargo 2 top risk) + satu coverage USD; cargo 1 coverage ber-LoL.
func cargoSpreadUji() []models.ObjekCargo {
	return []models.ObjekCargo{
		{ID: "11", TradingNote: "DOMESTIK", Coverages: []models.CoverageCargo{
			{Coverage: "ICC A", Currency: "IDR", TSI: d("1000"), TSILiability: d("1000"), Premium: d("10"), Discount: d("2"),
				LimitOfLiability: d("500"), Spreading: []models.BarisSpreading{{TreatyType: "10015", SharePercentage: d("50")},
					{TreatyType: "ORS", SharePercentage: d("50")}}},
			{Coverage: "ICC B", Currency: "USD", TSI: d("100"), TSILiability: d("100"), Premium: d("1"),
				Spreading: []models.BarisSpreading{{TreatyType: "10015", SharePercentage: d("100")}}}}},
		{ID: "12", IsTopRisk: true, Coverages: []models.CoverageCargo{
			{Coverage: "ICC A", Currency: "IDR", TSI: d("2000"), TSILiability: d("2000"), Premium: d("20"), LimitOfLiability: d("900"),
				Spreading: []models.BarisSpreading{{TreatyType: "10015", SharePercentage: d("40")}}}}},
	}
}

func layananSpreadCargo(gp string, tulis *[]models.ObjekCargo, share **apd.Decimal) *Service {
	svc, c := layananCargo(gp, cargoSpreadUji(), tulis)
	c.shareTulis = share
	return svc.DenganObjekCargo(c).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{{PercentShare: d("40"),
		StartDateTime: "20260101T050000.000 GMT"}}, idUang: map[string]string{"IDR": idMataUangIDR}})
}

// TestTotalCargo - SumTSIPremiSpreadedRNM_MARINECARGO_Act: LoL hanya cargo 1, ClaimSpreaded hanya cargo top risk, share beda ->
// kosong, per mata uang; label CARI15; FlagDelete "1" tidak dijumlah.
func TestTotalCargo(t *testing.T) {
	cargo := cargoSpreadUji()
	hitungNRCargo(cargo, d("40"))
	terapkanPersenCargo(cargo)
	tot, label := totalCargo(cargo, d("40"))
	// IDR: cargo1 TSINR 400 -> 10015 200 + ORS 200; cargo2 TSINR 800 -> 10015 320 (top). LoL cargo1 = 500×40/100 = 200.
	if len(tot) != 3 || tot[0].Currency != "IDR" || tot[0].TreatyType != "10015" || t8(tot[0].TSISpreaded) != "520" ||
		t8(tot[0].ClaimSpreaded) != "320" || t8(tot[0].ClaimEstimation) != "100" || tot[0].SharePercentage != nil ||
		t8(tot[0].PremiumSpreaded) != "4.8" {
		t.Errorf("IDR 10015: %+v", tot[0])
	}
	if tot[1].TreatyType != "ORS" || t8(tot[1].ClaimEstimation) != "100" || t8(tot[1].SharePercentage) != "50" ||
		tot[2].Currency != "USD" || t8(tot[2].ClaimEstimation) != "0" {
		t.Errorf("ORS / USD: %+v / %+v", tot[1], tot[2])
	}
	if label != "Total Limit of Liability(50.0000%)" {
		t.Errorf("label: %q", label)
	}
	cargo[1].FlagDelete = "1"
	if tot, _ := totalCargo(cargo, d("40")); t8(tot[0].TSISpreaded) != "200" {
		t.Errorf("FlagDelete: %+v", tot[0])
	}
}

// TestHitungShareSpreadingCargo - NR (CountPremiAndTSIRNMCargo_ACT 7) + DT percent; tanpa menulis; 8 desimal di jawaban.
func TestHitungShareSpreadingCargo(t *testing.T) {
	var tulis []models.ObjekCargo
	svc := layananSpreadCargo("004", &tulis, nil)
	v, err := svc.HitungShareSpreadingCargo(context.Background(), "K", d("40"))
	if err != nil {
		t.Fatal(err)
	}
	c := v.Cargo[0].Coverages[0]
	if t8(c.TSINusantaraRe) != "400" || t8(c.PremiNusantaraRe) != "3.2" || t8(c.Spreading[0].TSISpreaded) != "200" ||
		t8(c.Spreading[0].PremiumSpreaded) != "1.6" || v.Cargo[0].No != 1 || v.Cargo[0].TradingNote != "DOMESTIK" || len(tulis) != 0 {
		t.Errorf("hitung-share: %+v", v.Cargo[0])
	}
	if v.LabelLoL != "Total Limit of Liability(50.0000%)" || len(v.RingkasanMataUang) != 2 {
		t.Errorf("ringkasan: %q %+v", v.LabelLoL, v.RingkasanMataUang)
	}
	if _, err := layananSpreadCargo("006", &tulis, nil).TampilanSpreadingCargo(context.Background(), "K"); !errors.Is(err, ErrBukanMarineCargo) {
		t.Errorf("FIRE: %v", err)
	}
}

// TestSalinSpreadingCargo - CopyToAllSpreadingCargoMBU_ACT: template menggantikan baris setiap coverage; TSISpreaded = share ×
// (TSI × PS / 100) / 100, PremiumSpreaded = share × PremiNR / 100.
func TestSalinSpreadingCargo(t *testing.T) {
	svc := layananSpreadCargo("004", nil, nil)
	v, err := svc.SalinSpreadingCargo(context.Background(), inti.Pelaku{AkunID: "UJI"}, "K", d("40"),
		[]models.TemplateSpreading{{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: d("100")}})
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range v.Cargo {
		for m, cv := range c.Coverages {
			if len(cv.Spreading) != 1 || cv.Spreading[0].TreatyType != "10015" || t8(cv.Spreading[0].PremiumSpreaded) != t8(cv.PremiNusantaraRe) ||
				t8(cv.Spreading[0].TSISpreaded) != t8(persenDari(cv.TSI, d("40"))) {
				t.Errorf("cargo %d coverage %d: %+v", i, m, cv.Spreading)
			}
		}
	}
	if _, err := svc.SalinSpreadingCargo(context.Background(), inti.Pelaku{AkunID: "UJI"}, "K", d("40"), nil); !errors.Is(err, ErrMasukanSpreading) {
		t.Errorf("template kosong: %v", err)
	}
}

// TestSimpanSpreadingCargo - PUT: baris badan (tanpa pembaca tahap), hitungan Save, tulis pohon cargo + % Share RNM; bentuk beda
// 409; share baris tidak sah 400.
func TestSimpanSpreadingCargo(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var tulis []models.ObjekCargo
	var share *apd.Decimal
	svc := layananSpreadCargo("004", &tulis, &share)
	badan := cargoSpreadUji()
	badan[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: " 10015 ", TreatyName: "FACOUT", SharePercentage: d("100")}}
	if _, err := svc.SimpanSpreadingCargo(ctx, akun, "K", d("40"), badan); err != nil {
		t.Fatal(err)
	}
	c := tulis[0].Coverages[0]
	if t8(share) != "40" || len(c.Spreading) != 1 || c.Spreading[0].TreatyType != "10015" || t8(c.Spreading[0].TSISpreaded) != "400" ||
		t8(c.Spreading[0].PremiumSpreaded) != "3.2" || t8(c.PremiNusantaraRe) != "3.2" {
		t.Errorf("tulis: %+v share %v", c, share)
	}
	beda := cargoSpreadUji()
	beda[1].Coverages[0].Coverage = "LAIN"
	if _, err := svc.SimpanSpreadingCargo(ctx, akun, "K", d("40"), beda); !errors.Is(err, ErrSpreadingBerubah) {
		t.Errorf("bentuk: %v", err)
	}
	lebih := cargoSpreadUji()
	lebih[0].Coverages[0].Spreading[0].SharePercentage = d("101")
	if _, err := svc.SimpanSpreadingCargo(ctx, akun, "K", d("40"), lebih); !errors.Is(err, ErrMasukanSpreading) {
		t.Errorf("share 101: %v", err)
	}
}

// TestKasusSpreadingCargo - sumber ringkasanRNM (P4): per mata uang premi = Σ PremiNR, TSIRNM = PS × Σ TSILiability / 100.
func TestKasusSpreadingCargo(t *testing.T) {
	k := kasusSpreadingCargo(models.KasusSpreading{PercentShare: d("40")}, cargoSpreadUji())
	_, jumlah, _ := ringkasanRNM(k)
	if len(jumlah) != 2 || jumlah[0].Currency != "IDR" || t8(jumlah[0].TSIRNM) != "1200" || t8(jumlah[0].Premi) != "11.2" ||
		jumlah[1].Currency != "USD" || t8(jumlah[1].TSIRNM) != "40" {
		t.Errorf("jumlah: %+v", jumlah)
	}
}

// TestSpreadingRingkasCargo - P4: case MARINE CARGO tanpa item FIRE -> sumber Payment / Deduction / Cedant = pohon cargo; case
// bukan Marine -> pohon FIRE apa adanya.
func TestSpreadingRingkasCargo(t *testing.T) {
	svc := layananSpreadCargo("004", nil, nil)
	k, err := svc.bacaSpreadingRingkas(context.Background(), nil, "K")
	if err != nil {
		t.Fatal(err)
	}
	if _, jumlah, _ := ringkasanRNM(k); len(k.Lokasi) != 2 || len(jumlah) != 2 || t8(jumlah[0].Premi) != "11.2" {
		t.Errorf("Marine: %d lokasi %+v", len(k.Lokasi), jumlah)
	}
	fire := layananSpreadCargo("006", nil, nil)
	if k, err := fire.bacaSpreadingRingkas(context.Background(), nil, "K"); err != nil || len(k.Lokasi) != 0 {
		t.Errorf("FIRE: %v %+v", err, k)
	}
}
