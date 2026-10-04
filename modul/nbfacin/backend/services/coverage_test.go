package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

func teksD(d *apd.Decimal) string { return utils.FormatDecimal(d) }

// TestProrate - tiket 43, CountPremi langkah 3-8: setahun penuh = 100; 01-01 -> 01-07 = 181/365 (20 desimal) x 100;
// EDMDay menggantikan TotalDays; Begin / End kosong -> ErrPeriodeKasus.
func TestProrate(t *testing.T) {
	for nama, u := range map[string]struct {
		p   Periode
		mau string
	}{
		"setahun":     {Periode{Mulai: "20250101T050000.000 GMT", Akhir: "20260101T050000.000 GMT"}, "100"},
		"dua tahun":   {Periode{Mulai: "20240315T050000.000 GMT", Akhir: "20260315T050000.000 GMT"}, "200"},
		"setengah":    {Periode{Mulai: "20250101T050000.000 GMT", Akhir: "20250701T050000.000 GMT"}, "49.58904109589041095900"},
		"EDMDay 366":  {Periode{Mulai: "20250101T050000.000 GMT", Akhir: "20250701T050000.000 GMT", EDMDay: "366"}, "49.45355191256830601100"},
		"jam WIB 17":  {Periode{Mulai: "20241231T170000.000 GMT", Akhir: "20251231T170000.000 GMT"}, "100"},
		"satu thn+1b": {Periode{Mulai: "20250101T050000.000 GMT", Akhir: "20260201T050000.000 GMT"}, "108.49315068493150684900"},
	} {
		got, err := Prorate(u.p)
		if err != nil || got.Cmp(desimalWajib(u.mau)) != 0 {
			t.Errorf("%s: %v (%v), mau %s", nama, got, err, u.mau)
		}
	}
	for _, p := range []Periode{{Akhir: "20260101T050000.000 GMT"}, {Mulai: "20250101T050000.000 GMT"},
		{Mulai: "20250101T050000.000 GMT", Akhir: "20250701T050000.000 GMT", EDMDay: "abc"}} {
		if _, err := Prorate(p); !errors.Is(err, ErrPeriodeKasus) {
			t.Errorf("%+v: %v", p, err)
		}
	}
}

func desimalWajib(s string) *apd.Decimal {
	d, err := utils.ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return d
}

func covUji(basis string) models.CoverageObjek {
	return models.CoverageObjek{Coverage: "100815", CoverageBasis: basis, Rate: ds("1"), IndemnityPercentage: ds("100")}
}

// TestHitungCoverage - tiket 43, CountPremi langkah 19-52 (nilai dihitung tangan dari rumus korpus): TSI 1e9, Rate 1 permil,
// Prorate 100, Indemnity 100, LossLimit kosong -> 100.
func TestHitungCoverage(t *testing.T) {
	tsi, seratusD := ds("1000000000"), ds("100")
	tanpa := PenyesuaianItem{}
	adj80 := PenyesuaianItem{IsAdjustable: true, PctAdjustOther: ds("80")}
	ubah := func(basis string, f func(*models.CoverageObjek)) models.CoverageObjek {
		c := covUji(basis)
		f(&c)
		return c
	}
	for nama, u := range map[string]struct {
		c       models.CoverageObjek
		mode    string
		dm      string
		it      PenyesuaianItem
		premium string
		rate    string
		netRate string
		tsiL    string
	}{
		"b1 percent":           {covUji("1"), ModePercent, "", tanpa, "1000000", "1", "", "1000000000"},
		"b1 amount":            {ubah("1", func(c *models.CoverageObjek) { c.Rate, c.Premium = nil, ds("1000000") }), ModeAmount, "", tanpa, "1000000", "1", "", "1000000000"},
		"b1 adjust 80":         {covUji("1"), ModePercent, "", adj80, "800000", "1", "", "1000000000"},
		"b1 adjust 59 abaikan": {covUji("1"), ModePercent, "", PenyesuaianItem{true, ds("59")}, "1000000", "1", "", "1000000000"},
		"b1 net rate":          {ubah("1", func(c *models.CoverageObjek) { c.NetRate = ds("2") }), ModePercent, "", tanpa, "2000000", "1", "2", "1000000000"},
		"b1 net amount":        {ubah("1", func(c *models.CoverageObjek) { c.NetRate, c.Premium = ds("2"), ds("500000") }), ModeAmount, "", tanpa, "500000", "0.5", "0.5", "1000000000"},
		"b1 loss limit 50":     {ubah("1", func(c *models.CoverageObjek) { c.LostLimit = ds("50") }), ModePercent, "", tanpa, "500000", "1", "", "1000000000"},
		"b2 percent":           {ubah("2", func(c *models.CoverageObjek) { c.FirstScale, c.FirstLoss = ds("50"), ds("40") }), ModePercent, "", tanpa, "500000", "1", "", "400000000"},
		"b2 adjust 80":         {ubah("2", func(c *models.CoverageObjek) { c.FirstScale, c.FirstLoss = ds("50"), ds("40") }), ModePercent, "", adj80, "400000", "1", "", "400000000"},
		"b3 percent":           {ubah("3", func(c *models.CoverageObjek) { c.EmlPml = ds("25") }), ModePercent, "", tanpa, "1000000", "1", "", "250000000"},
		"b4 sublimit kosong":   {covUji("4"), ModePercent, "", tanpa, "1000000", "1", "", "1000000000"},
		"b4 sublimit 30":       {ubah("4", func(c *models.CoverageObjek) { c.Sublimit = ds("30") }), ModePercent, "", tanpa, "1000000", "1", "", "300000000"},
		// Langkah 45: PctAdjustment != 0 -> TSILiability memakai .Sublimit MENTAH (kosong = 0) - dipertahankan persis.
		"b4 adjust sublimit mentah": {covUji("4"), ModePercent, "", adj80, "800000", "1", "", "0"},
	} {
		got, err := HitungCoverage(u.c, tsi, seratusD, u.mode, u.dm, u.it)
		if err != nil {
			t.Errorf("%s: %v", nama, err)
			continue
		}
		if teksD(got.Premium) != u.premium || teksD(got.Rate) != u.rate || teksD(got.NetRate) != u.netRate || teksD(got.TSILiability) != u.tsiL ||
			teksD(got.TSI) != "1000000000" || teksD(got.ProRatePercent) != "100" {
			t.Errorf("%s: premi %s rate %s net %s tsiL %s tsi %s pro %s", nama, teksD(got.Premium), teksD(got.Rate), teksD(got.NetRate),
				teksD(got.TSILiability), teksD(got.TSI), teksD(got.ProRatePercent))
		}
	}
	// Diskon: percent 10 -> 100000; amount 250000 -> 25 %; mode diturunkan; ModeDiskonTanpa tidak menghitung.
	c := covUji("1")
	c.DiscountPercentage = ds("10")
	if got, _ := HitungCoverage(c, tsi, seratusD, ModePercent, "", tanpa); teksD(got.Discount) != "100000" {
		t.Errorf("diskon percent: %s", teksD(got.Discount))
	}
	c = covUji("1")
	c.Discount = ds("250000")
	if got, _ := HitungCoverage(c, tsi, seratusD, ModePercent, "", tanpa); teksD(got.DiscountPercentage) != "25" {
		t.Errorf("diskon amount: %s", teksD(got.DiscountPercentage))
	}
	c.DiscountPercentage = ds("10")
	if got, _ := HitungCoverage(c, tsi, seratusD, ModePercent, ModeAmount, tanpa); teksD(got.DiscountPercentage) != "25" {
		t.Errorf("modeDiskon amount eksplisit menang: %s", teksD(got.DiscountPercentage))
	}
	if got, _ := HitungCoverage(c, tsi, seratusD, ModePercent, ModeDiskonTanpa, tanpa); teksD(got.Discount) != "250000" {
		t.Errorf("tanpa diskon: %s", teksD(got.Discount))
	}
	// Pembulatan simpan 8 desimal setengah-ke-atas: premi 1.000.000 atas TSI 3e9 -> rate 1/3 permil.
	got, _ := HitungCoverage(ubah("1", func(c *models.CoverageObjek) { c.Rate, c.Premium = nil, ds("1000000") }), ds("3000000000"), seratusD, ModeAmount, "", tanpa)
	if teksD(got.Rate) != "0.33333333" {
		t.Errorf("rate 8 desimal: %s", teksD(got.Rate))
	}
	for nama, u := range map[string]struct {
		c     models.CoverageObjek
		tsi   *apd.Decimal
		mode  string
		pesan string
	}{
		"layering":       {covUji("5"), tsi, ModePercent, "Layering belum didukung"},
		"basis kosong":   {covUji(""), tsi, ModePercent, "coverageBasis harus 1-4"},
		"mode salah":     {covUji("1"), tsi, "x", "mode harus percent atau amount"},
		"amount TSI nol": {ubah("1", func(c *models.CoverageObjek) { c.Premium = ds("1") }), nil, ModeAmount, "rate tidak dapat dihitung"},
	} {
		if _, err := HitungCoverage(u.c, u.tsi, seratusD, u.mode, "", tanpa); !errors.Is(err, ErrMasukanCoverage) || !strings.Contains(err.Error(), u.pesan) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

// TestTotalPerMataUang - SumTotalTSIPremiGross_Act cabang FIRE: per mata uang urut kemunculan; TSI 0 -> rate 0.
func TestTotalPerMataUang(t *testing.T) {
	got := TotalPerMataUang([]models.ItemObjek{{Currency: "IDR", TSI: ds("1000"), TotalGrossPremi: ds("3")},
		{Currency: "USD", TotalGrossPremi: ds("5")}, {Currency: "IDR", TSI: ds("2000"), TotalGrossPremi: ds("6")}})
	if len(got) != 2 || got[0].Currency != "IDR" || teksD(got[0].TSI) != "3000" || teksD(got[0].Premium) != "9" || teksD(got[0].Rate) != "3" ||
		got[1].Currency != "USD" || teksD(got[1].Rate) != "0" || teksD(got[1].Premium) != "5" {
		t.Errorf("%+v", got)
	}
}

// TestGantiObjekCoverage - PUT: coverage dihitung ulang (rate -> percent, premi saja -> amount, keduanya kosong -> premi
// kosong), medan server dari badan diabaikan, TotalGrossPremi item = Σ premi; periode kosong -> ErrPeriodeKasus; basis 5
// -> 400 ber-indeks.
func TestGantiObjekCoverage(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var panggil int
	tersimpan := objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}
	kasus := kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": {CaseID: "UJI-NB-1",
		General: models.General{StartDateTime: "20250101T050000.000 GMT", EndDateTime: "20260101T050000.000 GMT"}},
		"UJI-NB-2": {CaseID: "UJI-NB-2"}}}
	svc := Baru(nil).DenganObjek(tersimpan).DenganTransaksi(tanpaTx).DenganPilihanItem(pilihanItemTiruan{&panggil}).DenganKasus(kasus)
	persen := covUji("1")
	persen.Premium, persen.TSI, persen.TSILiability = ds("999"), ds("1"), ds("1")
	amount := covUji("1")
	amount.Rate, amount.Premium = nil, ds("2000000")
	kosong := covUji("1")
	kosong.Rate = nil
	it := models.ItemObjek{Currency: "IDR", TSI: ds("1000000000"), Coverages: []models.CoverageObjek{persen, amount, kosong}}
	d, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", []models.ObjekFire{{ObjectType: "UJI", Items: []models.ItemObjek{it}}})
	if err != nil {
		t.Fatal(err)
	}
	cov := tersimpan.ada["UJI-NB-1"][0].Items[0]
	if teksD(cov.Coverages[0].Premium) != "1000000" || teksD(cov.Coverages[0].TSI) != "1000000000" || teksD(cov.Coverages[1].Rate) != "2" ||
		cov.Coverages[2].Premium != nil || teksD(cov.Coverages[2].TSILiability) != "1000000000" || teksD(cov.TotalGrossPremi) != "3000000" {
		t.Errorf("tersimpan: %+v", cov)
	}
	if len(d) != 1 || len(d[0].TotalPerCurrency) != 1 || teksD(d[0].TotalPerCurrency[0].Premium) != "3000000" || teksD(d[0].TotalPerCurrency[0].Rate) != "3" {
		t.Errorf("total per mata uang: %+v", d)
	}
	tersimpan.ada["UJI-NB-2"] = nil
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-2", []models.ObjekFire{{ObjectType: "UJI", Items: []models.ItemObjek{it}}}); !errors.Is(err, ErrPeriodeKasus) {
		t.Errorf("periode kosong: %v", err)
	}
	it.Coverages = []models.CoverageObjek{covUji("5")}
	if _, err := svc.GantiObjek(ctx, akun, "UJI-NB-1", []models.ObjekFire{{ObjectType: "UJI", Items: []models.ItemObjek{it}}}); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "baris[0].items[0].coverages[0].coverageBasis: Layering belum didukung") {
		t.Errorf("basis 5: %v", err)
	}
}

// TestHitungSimpanPremiTetap - A162 (temuan review sumbu spec): premi ketikan pengguna (mode amount) tidak bergeser saat
// Save walaupun rate dibalik tersimpan 8 desimal; pasangan basi (TSI berubah) tetap dihitung percent dari rate (A158);
// hasil di luar NUMBER(38,8) -> ErrMasukanCoverage.
func TestHitungSimpanPremiTetap(t *testing.T) {
	tsi, seratusD := ds("123456789012"), ds("100")
	amount := covUji("1")
	amount.Rate, amount.Premium = nil, ds("1000000")
	dariLayar, err := HitungCoverage(amount, tsi, seratusD, ModeAmount, "", PenyesuaianItem{})
	if err != nil || teksD(dariLayar.Rate) != "0.0081" {
		t.Fatalf("rate dibalik: %v %v", teksD(dariLayar.Rate), err)
	}
	// Bukti bahwa A158 polos memang menggeser premi (instrumen uji).
	if lama, _ := HitungCoverage(dariLayar, tsi, seratusD, ModePercent, "", PenyesuaianItem{}); teksD(lama.Premium) == "1000000" {
		t.Fatalf("percent dari rate 8 desimal mestinya bergeser, dapat %s", teksD(lama.Premium))
	}
	simpan, err := hitungSimpan(dariLayar, tsi, seratusD, PenyesuaianItem{})
	if err != nil || teksD(simpan.Premium) != "1000000" || teksD(simpan.Rate) != "0.0081" {
		t.Errorf("premi pengguna bergeser: premi %s rate %s %v", teksD(simpan.Premium), teksD(simpan.Rate), err)
	}
	// Rate diketik pengguna (percent): premi dihitung dari rate.
	persen := covUji("1")
	persen.Premium = ds("1")
	if h, err := hitungSimpan(persen, ds("1000000000"), seratusD, PenyesuaianItem{}); err != nil || teksD(h.Premium) != "1000000" || teksD(h.Rate) != "1" {
		t.Errorf("percent: %s %s %v", teksD(h.Premium), teksD(h.Rate), err)
	}
	// Pasangan basi (TSI item berubah sesudah dihitung): percent, rate pengguna dipertahankan.
	if h, err := hitungSimpan(dariLayar, ds("2000000000"), seratusD, PenyesuaianItem{}); err != nil || teksD(h.Rate) != "0.0081" || teksD(h.Premium) != "16200" {
		t.Errorf("basi: %s %s %v", teksD(h.Premium), teksD(h.Rate), err)
	}
	besar := covUji("1")
	besar.Rate = ds("999999999999999999999999999999")
	if _, err := HitungCoverage(besar, ds("999999999999999999999999999999"), seratusD, ModePercent, "", PenyesuaianItem{}); !errors.Is(err, ErrMasukanCoverage) ||
		!strings.Contains(err.Error(), "hasil premium melebihi 30 digit bulat") {
		t.Errorf("luap: %v", err)
	}
}
