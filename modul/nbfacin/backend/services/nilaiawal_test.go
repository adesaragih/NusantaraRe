package services

import (
	"context"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

// Tiket 68 - nilai awal coverage (keputusan work owner 08-10-2026, BUKAN dari korpus).

// TestNilaiAwalCoverage - kosong diisi (INDEMNITY "0", UNIT "0", RATE / DISCOUNT_PERCENTAGE / NET_RATE / LIMITOF_LIABILITY /
// PCT_LO_L / DISCOUNT / PREMIUM 0, LOST_LIMIT 100); isian pengguna (termasuk 0) tetap; medan khusus basis / RATE_OJK / server
// tidak disentuh.
func TestNilaiAwalCoverage(t *testing.T) {
	var c models.CoverageObjek
	nilaiAwalCoverage(&c)
	for nama, d := range map[string]*apd.Decimal{"rate": c.Rate, "discountPercentage": c.DiscountPercentage, "netRate": c.NetRate,
		"limitOfLiability": c.LimitOfLiability, "pctLol": c.PctLoL, "discount": c.Discount, "premium": c.Premium} {
		if teksD(d) != "0" {
			t.Errorf("%s = %q", nama, teksD(d))
		}
	}
	if teksD(c.LostLimit) != "100" || c.Indemnity != "0" || c.Unit != "0" {
		t.Errorf("lostLimit %q indemnity %q unit %q", teksD(c.LostLimit), c.Indemnity, c.Unit)
	}
	for nama, d := range map[string]*apd.Decimal{"firstLoss": c.FirstLoss, "firstScale": c.FirstScale, "sublimit": c.Sublimit,
		"emlPml": c.EmlPml, "rateOjk": c.RateOJK, "tsi": c.TSI, "tsiLiability": c.TSILiability, "proRatePercent": c.ProRatePercent} {
		if d != nil {
			t.Errorf("%s disentuh: %v", nama, d)
		}
	}
	isi := models.CoverageObjek{Indemnity: "UJI", Unit: "2", Rate: ds("1.5"), DiscountPercentage: ds("10"), NetRate: ds("0.3"),
		LimitOfLiability: ds("7"), PctLoL: ds("8"), Discount: ds("9"), Premium: ds("11"), LostLimit: ds("0")}
	salin := isi
	nilaiAwalCoverage(&isi)
	if isi.Indemnity != "UJI" || isi.Unit != "2" || teksD(isi.Rate) != "1.5" || teksD(isi.DiscountPercentage) != "10" ||
		teksD(isi.NetRate) != "0.3" || teksD(isi.LimitOfLiability) != "7" || teksD(isi.PctLoL) != "8" || teksD(isi.Discount) != "9" ||
		teksD(isi.Premium) != "11" || teksD(isi.LostLimit) != "0" || isi.Rate != salin.Rate {
		t.Errorf("isian pengguna berubah: %+v", isi)
	}
}

// layananNilaiAwal - case dengan periode satu tahun (prorate 100) untuk siapkanCoverage.
func layananNilaiAwal() *Service {
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1",
		General: models.General{StartDateTime: "20250101T050000.000 GMT", EndDateTime: "20260101T050000.000 GMT"}}}}
	return Baru(nil).DenganKasus(kasus)
}

func simpanSekali(t *testing.T, svc *Service, it models.ItemObjek) models.ItemObjek {
	t.Helper()
	baris := []models.ObjekFire{{Items: []models.ItemObjek{it}}}
	if err := svc.siapkanCoverage(context.Background(), "UJI-NB-1", baris); err != nil {
		t.Fatal(err)
	}
	return baris[0].Items[0]
}

// TestNilaiAwalCoverageSimpanUlang - nilai awal TIDAK mengubah premi: premi simpan pertama = hitungSimpan isian asli (kosong ->
// 0), dan simpan ulang hasilnya (kini berisi nilai awal) memberi premi / rate / diskon / net rate / total yang SAMA - percent,
// amount (premi ketikan dipertahankan), kosong, diskon persen / nominal, Loss Limit kosong / 50, basis 2 / 3 / 4, item adjustable.
func TestNilaiAwalCoverageSimpanUlang(t *testing.T) {
	svc := layananNilaiAwal()
	tsi := ds("1000000000")
	basis2 := covUji("2")
	basis2.FirstLoss, basis2.FirstScale = ds("40"), ds("3")
	basis3 := covUji("3")
	basis3.EmlPml = ds("60")
	basis4 := covUji("4")
	basis4.Sublimit = ds("25")
	amount := covUji("1")
	amount.Rate, amount.Premium = nil, ds("2000000")
	kosong := covUji("1")
	kosong.Rate = nil
	dp := covUji("1")
	dp.DiscountPercentage = ds("10")
	dn := covUji("1")
	dn.Discount = ds("50000")
	ll := covUji("1")
	ll.LostLimit = ds("50")
	for nama, c := range map[string]struct {
		cov  models.CoverageObjek
		item models.ItemObjek
	}{
		"percent": {cov: covUji("1")}, "amount": {cov: amount}, "kosong": {cov: kosong}, "diskon persen": {cov: dp},
		"diskon nominal": {cov: dn}, "loss limit 50": {cov: ll}, "basis 2": {cov: basis2}, "basis 3": {cov: basis3}, "basis 4": {cov: basis4},
		"adjustable": {cov: covUji("1"), item: models.ItemObjek{IsAdjustable: true, PctAdjustOther: ds("70")}},
	} {
		it := c.item
		it.Currency, it.TSI, it.Coverages = "IDR", tsi, []models.CoverageObjek{c.cov}
		mau, err := hitungSimpan(c.cov, tsi, ds("100"), PenyesuaianItem{it.IsAdjustable, it.PctAdjustOther})
		if err != nil {
			t.Fatal(err)
		}
		r1 := simpanSekali(t, svc, it)
		p1 := r1.Coverages[0]
		mauPremi := teksD(mau.Premium)
		if mauPremi == "" {
			mauPremi = "0"
		}
		if teksD(p1.Premium) != mauPremi || p1.Unit != "0" || p1.Indemnity != "0" || teksD(p1.PctLoL) != "0" ||
			(!it.IsAdjustable && teksD(p1.TSILiability) != teksD(mau.TSILiability)) {
			t.Errorf("%s simpan 1: premi %q mau %q; %+v", nama, teksD(p1.Premium), mauPremi, p1)
		}
		r2 := simpanSekali(t, svc, r1)
		p2 := r2.Coverages[0]
		for medan, pasang := range map[string][2]*apd.Decimal{"premium": {p1.Premium, p2.Premium}, "rate": {p1.Rate, p2.Rate},
			"discount": {p1.Discount, p2.Discount}, "discountPercentage": {p1.DiscountPercentage, p2.DiscountPercentage},
			"netRate": {p1.NetRate, p2.NetRate}, "lostLimit": {p1.LostLimit, p2.LostLimit}, "tsiLiability": {p1.TSILiability, p2.TSILiability},
			"totalGrossPremi": {r1.TotalGrossPremi, r2.TotalGrossPremi}} {
			if teksD(pasang[0]) != teksD(pasang[1]) {
				t.Errorf("%s simpan ulang %s: %q -> %q", nama, medan, teksD(pasang[0]), teksD(pasang[1]))
			}
		}
	}
}

// TestNilaiAwalCoverageNetRate - jalur TerapkanNetRate (kelima coverage net rate + Total Net Rate) ikut diberi nilai awal;
// premi = TerapkanNetRate isian asli.
func TestNilaiAwalCoverageNetRate(t *testing.T) {
	svc := layananNilaiAwal()
	var cov []models.CoverageObjek
	for _, id := range oldIDNetRate {
		c := covUji("1")
		c.OldID = id
		cov = append(cov, c)
	}
	it := models.ItemObjek{Currency: "IDR", TSI: ds("1000000000"), Coverages: cov, TotalNetRate: ds("0.5")}
	mau, err := TerapkanNetRate(cov, it.TSI, ds("100"), it.TotalNetRate, PenyesuaianItem{})
	if err != nil {
		t.Fatal(err)
	}
	r := simpanSekali(t, svc, it)
	for k, c := range r.Coverages {
		if teksD(c.Premium) != teksD(mau[k].Premium) || c.Unit != "0" || teksD(c.LostLimit) != "100" || teksD(c.LimitOfLiability) != "0" {
			t.Errorf("coverage %d: %+v (mau premi %q)", k, c, teksD(mau[k].Premium))
		}
	}
}

// TestNilaiAwalBasis - tiket 68 tambahan: SUBLIMIT 100, FIRST_LOSS 0, EML_PML 0 (kosong saja; FIRST_SCALE tidak) diisi SEBELUM
// hitung: basis 2 / 3 / 4 biasa -> TSILiability dan premi sama dengan kosong; TEPI basis 4 adjustable (Sublimit mentah,
// langkah 45-50) -> TSILiability kosong 0 menjadi TSI, premi tetap.
func TestNilaiAwalBasis(t *testing.T) {
	var c models.CoverageObjek
	nilaiAwalBasis(&c)
	if teksD(c.Sublimit) != "100" || teksD(c.FirstLoss) != "0" || teksD(c.EmlPml) != "0" || c.FirstScale != nil {
		t.Errorf("kosong: %+v", c)
	}
	isi := models.CoverageObjek{Sublimit: ds("0"), FirstLoss: ds("40"), EmlPml: ds("60")}
	nilaiAwalBasis(&isi)
	if teksD(isi.Sublimit) != "0" || teksD(isi.FirstLoss) != "40" || teksD(isi.EmlPml) != "60" {
		t.Errorf("isian berubah: %+v", isi)
	}
	svc := layananNilaiAwal()
	tsi := ds("1000000000")
	for _, basis := range []string{"2", "3", "4"} {
		cov := covUji(basis)
		mau, _ := hitungSimpan(cov, tsi, ds("100"), PenyesuaianItem{})
		r := simpanSekali(t, svc, models.ItemObjek{Currency: "IDR", TSI: tsi, Coverages: []models.CoverageObjek{cov}})
		p := r.Coverages[0]
		if teksD(p.TSILiability) != teksD(mau.TSILiability) || teksD(p.Premium) != teksD(mau.Premium) || teksD(p.Sublimit) != "100" ||
			teksD(p.FirstLoss) != "0" || teksD(p.EmlPml) != "0" || p.FirstScale != nil {
			t.Errorf("basis %s: %+v (mau TSILiability %q premi %q)", basis, p, teksD(mau.TSILiability), teksD(mau.Premium))
		}
	}
	adj := PenyesuaianItem{true, ds("70")}
	cov := covUji("4")
	mau, _ := hitungSimpan(cov, tsi, ds("100"), adj)
	r := simpanSekali(t, svc, models.ItemObjek{Currency: "IDR", TSI: tsi, Coverages: []models.CoverageObjek{cov}, IsAdjustable: true,
		PctAdjustOther: ds("70")})
	if p := r.Coverages[0]; teksD(mau.TSILiability) != "0" || teksD(p.TSILiability) != "1000000000" || teksD(p.Premium) != teksD(mau.Premium) {
		t.Errorf("tepi basis 4 adjustable: TSILiability %q (kosong %q), premi %q mau %q", teksD(p.TSILiability), teksD(mau.TSILiability),
			teksD(p.Premium), teksD(mau.Premium))
	}
}

// TestNilaiAwalTidakMengubahPemanggil - larik coverage pemanggil tetap kosong (nilai awal dipasang pada salinan).
func TestNilaiAwalTidakMengubahPemanggil(t *testing.T) {
	cov := []models.CoverageObjek{covUji("4")}
	r := simpanSekali(t, layananNilaiAwal(), models.ItemObjek{Currency: "IDR", TSI: ds("1000000000"), Coverages: cov})
	if cov[0].Sublimit != nil || cov[0].FirstLoss != nil || cov[0].Premium != nil || teksD(r.Coverages[0].Sublimit) != "100" {
		t.Errorf("pemanggil %+v / hasil %+v", cov[0], r.Coverages[0])
	}
}
