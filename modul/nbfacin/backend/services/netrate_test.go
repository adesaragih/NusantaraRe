package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// limaNetRate - lima coverage ber-OLDID CekNetRate_ACT dengan rate `rate[i]` (kosong = Rate nil).
func limaNetRate(rate ...string) []models.CoverageObjek {
	var hasil []models.CoverageObjek
	for i, k := range oldIDNetRate {
		c := covUji("1")
		c.OldID, c.Rate = k, nil
		if i < len(rate) && rate[i] != "" {
			c.Rate = ds(rate[i])
		}
		hasil = append(hasil, c)
	}
	return hasil
}

// TestFlagNetRate - tiket 44, CekNetRate_ACT langkah 2: kelima OLDID harfiah wajib ada; kurang satu, beda huruf, atau
// berspasi -> false; coverage tambahan tidak berpengaruh.
func TestFlagNetRate(t *testing.T) {
	if !FlagNetRate(limaNetRate()) {
		t.Error("lima lengkap mestinya true")
	}
	tambah := append(limaNetRate(), covUji("1"))
	if !FlagNetRate(tambah) {
		t.Error("coverage tambahan mestinya tidak berpengaruh")
	}
	for nama, ubah := range map[string]func([]models.CoverageObjek){
		"kurang satu": func(c []models.CoverageObjek) { c[4].OldID = "100840" },
		"huruf kecil": func(c []models.CoverageObjek) { c[0].OldID = "flexas" },
		"berspasi":    func(c []models.CoverageObjek) { c[1].OldID = "4.1A CC " },
		"ganda":       func(c []models.CoverageObjek) { c[2].OldID = "FLEXAS" },
	} {
		c := limaNetRate()
		ubah(c)
		if FlagNetRate(c) {
			t.Errorf("%s: mestinya false", nama)
		}
	}
	if FlagNetRate(nil) {
		t.Error("kosong mestinya false")
	}
}

// TestTerapkanNetRate - CalculateNetRate_ACT: NetRate = Rate / ΣRate (20 desimal) x TotalNetRate, premi dari NetRate
// presisi penuh (bukan yang dibulatkan 8), diskon tidak dihitung ulang, ΣRate 0 -> galat ber-jalur, masukan utuh.
func TestTerapkanNetRate(t *testing.T) {
	tsi, seratusD := ds("1000000000"), ds("100")
	cov := limaNetRate("1", "2", "1", "", "")
	cov[0].DiscountPercentage, cov[0].Discount = ds("10"), ds("5")
	h, err := TerapkanNetRate(cov, tsi, seratusD, ds("2"), PenyesuaianItem{})
	if err != nil {
		t.Fatal(err)
	}
	for i, mau := range [][2]string{{"0.5", "500000"}, {"1", "1000000"}, {"0.5", "500000"}, {"0", "0"}, {"0", "0"}} {
		if teksD(h[i].NetRate) != mau[0] || teksD(h[i].Premium) != mau[1] {
			t.Errorf("coverage %d: netRate %s premi %s, mau %v", i, teksD(h[i].NetRate), teksD(h[i].Premium), mau)
		}
	}
	if teksD(h[0].Discount) != "5" {
		t.Errorf("diskon dihitung ulang (DiscountStatus kosong di Pega): %s", teksD(h[0].Discount))
	}
	if cov[0].NetRate != nil {
		t.Error("larik masukan berubah")
	}
	// 1/3 pada 20 desimal: NetRate tersimpan 8 desimal, premi dari nilai 20 desimal.
	h, err = TerapkanNetRate(limaNetRate("1", "1", "1"), tsi, seratusD, ds("1"), PenyesuaianItem{})
	if err != nil || teksD(h[0].NetRate) != "0.33333333" || teksD(h[0].Premium) != "333333.33333333" {
		t.Errorf("sepertiga: %s %s %v", teksD(h[0].NetRate), teksD(h[0].Premium), err)
	}
	// TotalNetRate kosong = 0 (langkah 1 @divide("", 1, 20)): NetRate 0, premi dari Rate.
	if h, err = TerapkanNetRate(limaNetRate("1"), tsi, seratusD, nil, PenyesuaianItem{}); err != nil || teksD(h[0].NetRate) != "0" || teksD(h[0].Premium) != "1000000" {
		t.Errorf("tanpa total: %s %s %v", teksD(h[0].NetRate), teksD(h[0].Premium), err)
	}
	_, err = TerapkanNetRate(limaNetRate("0", ""), tsi, seratusD, ds("2"), PenyesuaianItem{})
	var g *galatNetRate
	if !errors.Is(err, ErrMasukanCoverage) || !errors.As(err, &g) || g.jalur != "totalNetRate" {
		t.Errorf("ΣRate 0: %v", err)
	}
}

// TestHitungNetRateKasus - POST: flag false -> TotalNetRate 0, coverage apa adanya, periode tidak diperiksa (case tetap
// dibaca: tidak ada -> 404); flag true -> dibagi; isian tak sah -> ErrMasukanCoverage ber-jalur coverages[k].
func TestHitungNetRateKasus(t *testing.T) {
	ctx := context.Background()
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1",
		General: models.General{StartDateTime: "20250101T050000.000 GMT", EndDateTime: "20260101T050000.000 GMT"}}}}
	svc := Baru(nil).DenganKasus(kasus)
	tanpaFlag := []models.CoverageObjek{covUji("1")}
	tanpaFlag[0].Premium = ds("7")
	tanpaPeriode := Baru(nil).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1"}}})
	h, err := tanpaPeriode.HitungNetRateKasus(ctx, "UJI-NB-1", tanpaFlag, ds("1000000000"), ds("2"), PenyesuaianItem{})
	if err != nil || h.Flag || teksD(h.TotalNetRate) != "0" || teksD(h.Coverages[0].Premium) != "7" {
		t.Errorf("flag false: %+v %v", h, err)
	}
	if _, err := svc.HitungNetRateKasus(ctx, "UJI-NB-9", tanpaFlag, nil, nil, PenyesuaianItem{}); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("case tidak ada: %v", err)
	}
	if _, err := Baru(nil).HitungNetRateKasus(ctx, "UJI-NB-1", tanpaFlag, nil, nil, PenyesuaianItem{}); !errors.Is(err, ErrCoverageTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
	h, err = svc.HitungNetRateKasus(ctx, "UJI-NB-1", limaNetRate("1", "1", "2"), ds("1000000000"), ds("2"), PenyesuaianItem{})
	if err != nil || !h.Flag || teksD(h.TotalNetRate) != "2" || teksD(h.Coverages[2].NetRate) != "1" {
		t.Errorf("flag true: %+v %v", h, err)
	}
	salah := limaNetRate("1")
	salah[3].CoverageBasis = "5"
	if _, err := svc.HitungNetRateKasus(ctx, "UJI-NB-1", salah, ds("1"), ds("2"), PenyesuaianItem{}); !errors.Is(err, ErrMasukanCoverage) ||
		!strings.Contains(err.Error(), "coverages[3].coverageBasis: Layering") {
		t.Errorf("basis 5: %v", err)
	}
}

// TestGantiObjekNetRate - PUT: item ber-flag dengan Total Net Rate -> dibagi; ber-flag tanpa Total Net Rate -> hitung
// biasa (A166), Total Net Rate tetap kosong; tanpa flag -> 0; ΣRate 0 -> 400 ber-jalur item.
func TestGantiObjekNetRate(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var panggil int
	tersimpan := objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1",
		General: models.General{StartDateTime: "20250101T050000.000 GMT", EndDateTime: "20260101T050000.000 GMT"}}}}
	svc := Baru(nil).DenganObjek(tersimpan).DenganTransaksi(tanpaTx).DenganPilihanItem(pilihanItemTiruan{&panggil}).DenganKasus(kasus)
	item := func(tnr string, cov []models.CoverageObjek) models.ItemObjek {
		it := models.ItemObjek{Currency: "IDR", TSI: ds("1000000000"), Coverages: cov}
		if tnr != "" {
			it.TotalNetRate = ds(tnr)
		}
		return it
	}
	manual := limaNetRate("1")
	manual[0].NetRate = ds("3")
	baris := []models.ObjekFire{{ObjectType: "UJI", Items: []models.ItemObjek{
		item("2", limaNetRate("1", "1", "2")), item("", manual), item("9", []models.CoverageObjek{covUji("1")})}}}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", baris); err != nil {
		t.Fatal(err)
	}
	it := tersimpan.ada["UJI-NB-1"][0].Items
	if teksD(it[0].Coverages[2].NetRate) != "1" || teksD(it[0].Coverages[2].Premium) != "1000000" || teksD(it[0].TotalNetRate) != "2" ||
		teksD(it[0].TotalGrossPremi) != "2000000" {
		t.Errorf("dibagi: %+v", it[0])
	}
	if teksD(it[1].Coverages[0].NetRate) != "3" || teksD(it[1].Coverages[0].Premium) != "3000000" || it[1].TotalNetRate != nil {
		t.Errorf("flag tanpa total: %+v", it[1])
	}
	if teksD(it[2].TotalNetRate) != "0" {
		t.Errorf("tanpa flag: %s", teksD(it[2].TotalNetRate))
	}
	baris = []models.ObjekFire{{ObjectType: "UJI", Items: []models.ItemObjek{item("2", limaNetRate())}}}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", baris); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "baris[0].items[0].totalNetRate: Σ Rate coverage = 0") {
		t.Errorf("ΣRate 0: %v", err)
	}
}
