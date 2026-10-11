package services

import (
	"context"
	"errors"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

func samaD(t *testing.T, nama string, got *apd.Decimal, mau string) {
	t.Helper()
	if got == nil || got.Cmp(desimalWajib(mau)) != 0 {
		t.Errorf("%s: %v, mau %s", nama, got, mau)
	}
}

// periode uji: 1 Jan 2026 - 1 Jan 2027 (teks Pega 05:00 GMT = 12:00 WIB).
func konteksAnekaUji(mulai, akhir string, pr string) *KonteksHitungAneka {
	k := &KonteksHitungAneka{Mulai: mulai + "T050000.000 GMT", Akhir: akhir + "T050000.000 GMT"}
	if pr != "" {
		k.ProRatePercent = desimalWajib(pr)
	}
	return k
}

// TestHitungAnekaSampelKorpus - instrumen: premi coverage nyata `D:\migrasi\RNM\DDL\CONTOH` (riset P2 bagian 6, 127/127 cocok
// dengan rumus CPCA): NB-143055 L134 (metode 1, prorate CASE 50,684…), NB-147721 L119 (metode 3, LostLimit 58,08… tanpa skala),
// NB-183524 L101 (Golf metode 3), NB-183603 L120 (LostLimit 85). Pemicu tanpa Status -> ProRate case dipakai apa adanya.
func TestHitungAnekaSampelKorpus(t *testing.T) {
	for _, c := range []struct {
		nama, pr string
		cov      models.CoverageAneka
		mau      string
	}{
		{"NB-143055", "50.68493150684931506800", models.CoverageAneka{CalculateMethodFacIn: "1", TSI: ds("2500000.0"),
			Rate: ds("0.939549"), ProRatePercent: ds("49.86301369863013698600"), Loading: ds("0"), LostLimit: ds("100"), Day: "365"},
			"11905.2442"},
		{"NB-147721", "", models.CoverageAneka{CalculateMethodFacIn: "3", TSI: ds("79994175360"), Rate: ds("0.8"), Loading: ds("0"),
			LostLimit: ds("58.08219178082190000"), Day: "365"}, "371698962.76865742"},
		{"NB-183524 Golf", "", models.CoverageAneka{CalculateMethodFacIn: "3", TSI: ds("1356444375"), Rate: ds("2"), Loading: ds("0"),
			LostLimit: ds("100"), Day: "365"}, "27128887.5"},
		{"NB-183603", "100", models.CoverageAneka{CalculateMethodFacIn: "1", TSI: ds("10000000000"), Rate: ds("0.01"),
			Loading: ds("0"), LostLimit: ds("85"), Day: "365"}, "850000"},
	} {
		k := konteksAnekaUji("20260101", "20270101", c.pr)
		got, err := HitungCoverageAneka(k, c.cov, PemicuLoading)
		if err != nil {
			t.Fatalf("%s: %v", c.nama, err)
		}
		samaD(t, c.nama, got.Premium, c.mau)
	}
}

// TestHitungAnekaPemicu - rantai per pemicu (K73-12…K73-15).
func TestHitungAnekaPemicu(t *testing.T) {
	dasar := models.CoverageAneka{CalculateMethodFacIn: "1", TSI: ds("1000000"), Rate: ds("0.5"), Loading: ds("0"), Day: "365"}
	// Status="Coverage" (tsi / rate / currency): ProRate case dihitung ulang dari tanggal (CPCA langkah 4), nilai masuk diabaikan.
	for _, p := range []string{PemicuTSI, PemicuRate, PemicuCurrency} {
		k := konteksAnekaUji("20260101", "20260702", "10")
		got, err := HitungCoverageAneka(k, dasar, p)
		if err != nil {
			t.Fatal(err)
		}
		samaD(t, p+" prorate case", k.ProRatePercent, "49.86301370") // 182 / 365 x 100, simpan8
		samaD(t, p+" premi", got.Premium, "2493.1507")               // r4(1e6 x 0,5 x 49,8630137 / 1e4)
	}
	// Tanpa Status (metode / loading / diskon / indemnity): ProRate case dari badan dipakai apa adanya.
	k := konteksAnekaUji("20260101", "20260702", "10")
	got, _ := HitungCoverageAneka(k, dasar, PemicuMetode)
	samaD(t, "metode premi", got.Premium, "500")
	samaD(t, "metode prorate tetap", k.ProRatePercent, "10")
	// 14.1: metode 1 dengan ProRate case 0 / kosong -> tahun polis + pecahan (TotalDays = .Day); nilai case TIDAK diubah.
	k = konteksAnekaUji("20260101", "20270702", "")
	got, _ = HitungCoverageAneka(k, dasar, PemicuLoading)
	samaD(t, "14.1 premi", got.Premium, "7493.1507") // (100 + 182/365 x 100) = 149,86301…
	if k.ProRatePercent != nil {
		t.Errorf("14.1 menulis prorate case: %v", k.ProRatePercent)
	}
	// FillPremiGolf (shortPeriod / scale): prorate case SELALU dari tanggal (L399), Golf pro-rata manual tidak bertahan.
	k = konteksAnekaUji("20260101", "20270101", "0")
	got, _ = HitungCoverageAneka(k, dasar, PemicuScale)
	samaD(t, "scale prorate", k.ProRatePercent, "100")
	samaD(t, "scale premi", got.Premium, "5000")
	// fillPremiBond (proRate): PR = .ProRatePercent COVERAGE (kosong = 100); minPremium = premium.
	bond := models.CoverageAneka{TSI: ds("1000"), Rate: ds("2"), ProRatePercent: ds("50"), Loading: ds("10")}
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", "77"), bond, PemicuProRate)
	samaD(t, "bond premi", got.Premium, "11") // r4(20 x 50 / 100) = 10 + r4(10 x 10 / 100) = 1
	samaD(t, "bond minPremium", got.MinPremium, "11")
	bond.ProRatePercent = nil
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", ""), bond, PemicuProRate)
	samaD(t, "bond PR kosong = 100", got.Premium, "22")
	// Day (ChangeCoveragePremium): metode 1 -> prorate case dari tanggal, coverage.ProRatePercent = case; premi tetap.
	k = konteksAnekaUji("20260101", "20260702", "10")
	c := dasar
	c.Premium, c.Day = ds("123"), "0"
	got, _ = HitungCoverageAneka(k, c, PemicuDay)
	samaD(t, "day prorate (Day 0 -> 365)", k.ProRatePercent, "49.86301370")
	samaD(t, "day coverage prorate", got.ProRatePercent, "49.86301370")
	samaD(t, "day premi tetap", got.Premium, "123")
	c.CalculateMethodFacIn = "3"
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20260702", "10"), c, PemicuDay)
	samaD(t, "day metode 3: salin case", got.ProRatePercent, "10")
	// CalculateRate (premium): pecahan prorate DIJUMLAH DUA KALI (K73-15), rate = premi / TSI / prorate x 10000.
	c = dasar
	c.Premium = ds("4958.9041")
	k = konteksAnekaUji("20260101", "20260702", "10")
	got, _ = HitungCoverageAneka(k, c, PemicuPremium)
	samaD(t, "calculateRate prorate", k.ProRatePercent, "99.72602740") // 2 x 182/365 x 100
	samaD(t, "calculateRate rate", got.Rate, "0.49725275")             // (4958,9041 / 1e6) / 99,7260274 x 1e4
	// Diskon: persen -> amount dari premi SEBELUM CPCA; amount -> persen.
	c = dasar
	c.Premium, c.DiscountPercentage = ds("1000"), ds("10")
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", "100"), c, PemicuDiskonPersen)
	samaD(t, "diskon dari premi lama", got.Discount, "100")
	samaD(t, "premi baru", got.Premium, "5000")
	c.Discount, c.DiscountPercentage = ds("250"), nil
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", "100"), c, PemicuDiskonJumlah)
	samaD(t, "persen dari jumlah", got.DiscountPercentage, "25")
	// LoL.
	c = dasar
	c.PctLoL = ds("20")
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", ""), c, PemicuPctLoL)
	samaD(t, "LoL", got.LimitOfLiability, "200000")
	c.LimitOfLiability, c.PctLoL = ds("250000"), nil
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", ""), c, PemicuLoL)
	samaD(t, "%LoL", got.PctLoL, "25")
	// MBD: premi ber-indemnity memakai prorate CASE.
	k = &KonteksHitungAneka{Mulai: "20260101T050000.000 GMT", Akhir: "20270101T050000.000 GMT", MBD: true, ProRatePercent: ds("100")}
	c = dasar
	c.IndemnityPercentage = ds("50")
	got, _ = HitungCoverageAneka(k, c, PemicuLoading)
	samaD(t, "MBD", got.Premium, "2500")
	// MBD (17.1) memakai prorate CASE apa adanya - case 0 -> premi 0 (bukan prorate 14.1 dari tahun polis).
	k = &KonteksHitungAneka{Mulai: "20260101T050000.000 GMT", Akhir: "20270101T050000.000 GMT", MBD: true, ProRatePercent: ds("0")}
	got, _ = HitungCoverageAneka(k, c, PemicuLoading)
	samaD(t, "MBD case 0", got.Premium, "0")
	k.MBD, k.ProRatePercent = false, ds("0")
	got, _ = HitungCoverageAneka(k, c, PemicuLoading)
	samaD(t, "non-MBD case 0 -> 14.1", got.Premium, "5000")
	// Metode tak dikenal: premi tetap; pemicu tak dikenal: 400; tanggal case kosong pada pemicu ber-tanggal: ErrPeriodeKasus.
	c = dasar
	c.CalculateMethodFacIn, c.Premium = "", ds("7")
	got, _ = HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", "100"), c, PemicuLoading)
	samaD(t, "metode kosong", got.Premium, "7")
	if _, err := HitungCoverageAneka(konteksAnekaUji("20260101", "20270101", ""), dasar, "x"); !errors.Is(err, ErrMasukanHitungAneka) {
		t.Errorf("pemicu: %v", err)
	}
	if _, err := HitungCoverageAneka(&KonteksHitungAneka{}, dasar, PemicuTSI); !errors.Is(err, ErrPeriodeKasus) {
		t.Errorf("periode: %v", err)
	}
}

// TestMedanServerAneka - K73-13 / K73-14 dengan sampel: NB-143055 Share 99,5 -> PremiNR 11845,717979, TSINR 2487500, TSILiability = TSI;
// NB-183524 Golf Share 20 -> PremiNR 5425777,5, TSINR 271288875, TSILiability kosong; minPremium = premium; share kosong -> NR kosong;
// coverage tidak dikirim (nil) tidak disentuh.
func TestMedanServerAneka(t *testing.T) {
	lok := []models.LokasiAneka{{
		Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Coverages: []models.CoverageAneka{{TSI: ds("2500000"),
			Premium: ds("11905.2442")}}}, {}}}},
		ObjekGolf: []models.ObjekAneka{{Coverages: []models.CoverageAneka{{TSI: ds("1356444375"), Premium: ds("27128887.5"),
			Discount: ds("0"), TSILiability: ds("9")}}}},
	}}
	medanServerAneka(lok, ds("99.5"))
	o := lok[0].Okupasi[0].Objek[0].Coverages[0]
	samaD(t, "PremiNR", o.PremiNusantaraRe, "11845.717979")
	samaD(t, "TSINR", o.TSINusantaraRe, "2487500")
	samaD(t, "TSILiability", o.TSILiability, "2500000")
	samaD(t, "minPremium", o.MinPremium, "11905.2442")
	if lok[0].Okupasi[0].Objek[1].Coverages != nil {
		t.Error("coverage tidak dikirim disentuh")
	}
	medanServerAneka(lok, ds("20"))
	g := lok[0].ObjekGolf[0].Coverages[0]
	samaD(t, "Golf PremiNR", g.PremiNusantaraRe, "5425777.5")
	samaD(t, "Golf TSINR", g.TSINusantaraRe, "271288875")
	if g.TSILiability != nil {
		t.Errorf("Golf TSILiability: %v", g.TSILiability)
	}
	medanServerAneka(lok, nil)
	if lok[0].ObjekGolf[0].Coverages[0].PremiNusantaraRe != nil {
		t.Error("share kosong tetap mengisi NR")
	}
}

// TestMedanServerAnekaDiskonDanPUT - PremiNR memotong diskon (Share x Premium / 100 - Share x Discount / 100); PUT objek-aneka
// mengisi medan server dari PERCENT_SHARE case untuk coverage yang dikirim.
func TestMedanServerAnekaDiskonDanPUT(t *testing.T) {
	lok := []models.LokasiAneka{{ObjekGolf: []models.ObjekAneka{{Coverages: []models.CoverageAneka{{TSI: ds("2000"),
		Premium: ds("1000"), Discount: ds("100")}}}}}}
	medanServerAneka(lok, ds("50"))
	samaD(t, "PremiNR berdiskon", lok[0].ObjekGolf[0].Coverages[0].PremiNusantaraRe, "450")
	var tulis []models.LokasiAneka
	svc := layananAneka("10076", &tulis)
	svc.objekAneka = anekaTiruan{ada: map[string][]models.LokasiAneka{"K": {}}, tulis: &tulis, share: ds("50")}
	l := lokasiAnekaUji(models.ObjekAneka{Currency: "IDR", Coverages: []models.CoverageAneka{{TSI: ds("2000"), Premium: ds("1000"),
		Discount: ds("100"), MinPremium: ds("1"), TSILiability: ds("1")}}})
	if _, err := svc.GantiObjekAneka(context.Background(), inti.Pelaku{AkunID: "UJI"}, "K", l, nil); err != nil {
		t.Fatal(err)
	}
	c := tulis[0].Okupasi[0].Objek[0].Coverages[0]
	samaD(t, "PUT PremiNR", c.PremiNusantaraRe, "450")
	samaD(t, "PUT TSINR", c.TSINusantaraRe, "1000")
	samaD(t, "PUT TSILiability", c.TSILiability, "2000")
	samaD(t, "PUT minPremium", c.MinPremium, "1000")
}

// TestSetIndemnityRate - SetIndemnityRate_ACT: kunci Indemnity + Unit apa adanya; baris TERAKHIR menang; nol baris -> kosong;
// Unit 0 / 2 -> 100 (L1079); Indemnity bukan bilangan bulat 4 digit -> tidak dikueri; IsMBD -> CPCA (premi ber-indemnity),
// selain itu premi tetap; tanpa pencari -> 503.
func TestSetIndemnityRate(t *testing.T) {
	var kunci []string
	cari := func(baris ...*apd.Decimal) func(string, string) ([]*apd.Decimal, error) {
		return func(i, u string) ([]*apd.Decimal, error) {
			kunci = append(kunci, i+"|"+u)
			return baris, nil
		}
	}
	c := models.CoverageAneka{CalculateMethodFacIn: "1", TSI: ds("1000000"), Rate: ds("0.5"), Loading: ds("0"), Day: "365",
		Indemnity: "30", Unit: "1", Premium: ds("9")}
	k := &KonteksHitungAneka{Mulai: "20260101T050000.000 GMT", Akhir: "20270101T050000.000 GMT", MBD: true, ProRatePercent: ds("100"),
		CariIndemnity: cari(ds("10"), ds("40"))}
	got, err := HitungCoverageAneka(k, c, PemicuIndemnity)
	if err != nil || len(kunci) != 1 || kunci[0] != "30|1" {
		t.Fatalf("%v %v", err, kunci)
	}
	samaD(t, "baris terakhir", got.IndemnityPercentage, "40")
	samaD(t, "MBD premi", got.Premium, "2000") // 1e6 x 0,5 x 100 x 40 / 1e6
	for _, u := range []string{"0", "2"} {
		c2 := c
		c2.Unit = u
		got, _ = HitungCoverageAneka(&KonteksHitungAneka{CariIndemnity: cari(ds("10")), ProRatePercent: ds("100")}, c2, PemicuIndemnity)
		samaD(t, "unit "+u, got.IndemnityPercentage, "100")
		samaD(t, "non-MBD premi tetap", got.Premium, "9")
	}
	c.IndemnityPercentage = ds("55")
	got, _ = HitungCoverageAneka(&KonteksHitungAneka{CariIndemnity: cari()}, c, PemicuIndemnity)
	if got.IndemnityPercentage != nil {
		t.Errorf("nol baris: %v", got.IndemnityPercentage)
	}
	kunci = nil
	c.Indemnity = "1,5"
	if got, _ = HitungCoverageAneka(&KonteksHitungAneka{CariIndemnity: cari(ds("10"))}, c, PemicuIndemnity); len(kunci) != 0 ||
		got.IndemnityPercentage != nil {
		t.Errorf("indemnity tidak sah dikueri: %v %v", kunci, got.IndemnityPercentage)
	}
	if _, err := HitungCoverageAneka(&KonteksHitungAneka{}, c, PemicuIndemnity); !errors.Is(err, ErrObjekTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
}
