package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// TestLiniSpreadingAneka - instrumen: When lini dievaluasi dari BUSINESS case (lini yang jawabannya diketahui dari TestLiniKasus /
// TestTampilanAneka).
func TestLiniSpreadingAneka(t *testing.T) {
	ctx := context.Background()
	for nama, u := range map[string]struct {
		old  string
		mau  liniSpreadingAneka
		kode []string
	}{
		"aneka":     {"30", liniSpreadingAneka{aneka: true}, nil},
		"golf":      {"27", liniSpreadingAneka{golf: true}, nil},
		"bonding":   {"36", liniSpreadingAneka{aneka: true}, nil},
		"liability": {"15", liniSpreadingAneka{aneka: true, mataUangCov24: true, mataUangCov25: true, topRisk24: true, topRisk25: true}, []string{"10021"}},
		"MBD": {"10", liniSpreadingAneka{aneka: true, engineering: true, mbd: true, mataUangCov24: true, mataUangCov25: true,
			topRisk24: true, topRisk25: true, mataUangCount: true, tsiCoverageCount: true}, nil},
	} {
		s := layananLini(satuBisnis("003", u.old, u.kode...), "UJI KELAS")
		kb, _ := s.kasusBisnisDari(ctx, "K")
		if got := liniSpreadingDari("K", kb); got != u.mau {
			t.Errorf("%s: %+v, mau %+v", nama, got, u.mau)
		}
	}
}

// lokasiSpreadUji - satu lokasi IDR: objek okupasi A (TSI 1000) dua coverage, objek Golf G (USD) satu coverage.
//   - c1: TSI 1000, premi 10, diskon 2, LoL 800, baris OR 60 % (TSIGross 77 tersimpan) + QS 40 %;
//   - c2: TSI 500, premi 5, baris OR 100 %;
//   - Golf: TSI 200, premi 4, baris OR 100 %.
func lokasiSpreadUji() []models.LokasiAneka {
	return []models.LokasiAneka{{ObjectNo: "1", RiskLocation: "JL UJI", Okupasi: []models.OkupasiAneka{{OccupationName: "OKU", Objek: []models.ObjekAneka{{
		ObjectName: "A", Currency: "IDR", TSI: d("1000"), Coverages: []models.CoverageAneka{
			{OldID: "1", TSI: d("1000"), TSILiability: d("1000"), Premium: d("10"), Discount: d("2"), LimitOfLiability: d("800"),
				Spreading: []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("60"), TSIGrossSpreaded: d("77")},
					{TreatyType: "QS", SharePercentage: d("40")}}},
			{OldID: "2", TSI: d("500"), Premium: d("5"), Spreading: []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("100")}}},
		}}}}},
		ObjekGolf: []models.ObjekAneka{{ObjectName: "G", Currency: "USD", TSI: d("200"), Coverages: []models.CoverageAneka{
			{OldID: "9", TSI: d("200"), Premium: d("4"), Spreading: []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("100")}}}}}},
	}}
}

func barisTeks(b models.BarisSpreading) string {
	return strings.Join([]string{b.TreatyType, t8(b.SharePercentage), t8(b.TSISpreaded), t8(b.PremiumSpreaded), t8(b.TSIGrossSpreaded),
		t8(b.ClaimEstimation)}, "|")
}

func totalSpreadTeks(x models.TotalSpreading) string {
	return strings.Join([]string{x.Currency, x.TreatyType, t8(x.SharePercentage), t8(x.TSISpreaded), t8(x.PremiumSpreaded),
		t8(x.ClaimEstimation), t8(x.ClaimSpreaded)}, "|")
}

// TestNRDanPersenAneka - NR dari .TSI (L3853); DT percent hanya TSISpreaded / PremiumSpreaded, TSIGross tidak disentuh; Golf
// hanya bila isGolfInsurance.
func TestNRDanPersenAneka(t *testing.T) {
	lokasi := lokasiSpreadUji()
	hitungNRAneka(lokasi, d("50"))
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	if t8(c[0].TSINusantaraRe) != "500" || t8(c[0].PremiNusantaraRe) != "4" || t8(c[1].TSINusantaraRe) != "250" ||
		t8(c[1].PremiNusantaraRe) != "2.5" || t8(lokasi[0].ObjekGolf[0].Coverages[0].PremiNusantaraRe) != "2" {
		t.Fatalf("NR: %s %s %s %s", t8(c[0].TSINusantaraRe), t8(c[0].PremiNusantaraRe), t8(c[1].TSINusantaraRe), t8(c[1].PremiNusantaraRe))
	}
	terapkanPersenAneka(lokasi, liniSpreadingAneka{aneka: true})
	for got, mau := range map[string]string{barisTeks(c[0].Spreading[0]): "OR|60|300|2.4|77|", barisTeks(c[0].Spreading[1]): "QS|40|200|1.6||",
		barisTeks(c[1].Spreading[0]): "OR|100|250|2.5||", barisTeks(lokasi[0].ObjekGolf[0].Coverages[0].Spreading[0]): "OR|100||||"} {
		if got != mau {
			t.Errorf("baris %q, mau %q", got, mau)
		}
	}
	terapkanPersenAneka(lokasi, liniSpreadingAneka{golf: true})
	if g := barisTeks(lokasi[0].ObjekGolf[0].Coverages[0].Spreading[0]); g != "OR|100|100|2||" {
		t.Errorf("Golf: %q", g)
	}
}

// TestSalinTemplateAneka - Copy To All: TSISpreaded = share × (TSI × PS / 100) / 100, TSIGross / ClaimEstimation kosong;
// objek okupasi hanya bila IsAneka / Bonding, objek Golf hanya bila isGolfInsurance.
func TestSalinTemplateAneka(t *testing.T) {
	lokasi := lokasiSpreadUji()
	hitungNRAneka(lokasi, d("50"))
	salinTemplateAneka(lokasi, d("50"), []models.TemplateSpreading{{TreatyType: "OR", TreatyName: "UJI OR", SharePercentage: d("70")},
		{TreatyType: "QS", SharePercentage: d("30")}}, liniSpreadingAneka{aneka: true})
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	if len(c[0].Spreading) != 2 || barisTeks(c[0].Spreading[0]) != "OR|70|350|2.8||" || c[0].Spreading[0].TreatyName != "UJI OR" ||
		barisTeks(c[1].Spreading[1]) != "QS|30|75|0.75||" {
		t.Errorf("okupasi: %v", c)
	}
	if g := lokasi[0].ObjekGolf[0].Coverages[0].Spreading; len(g) != 1 || barisTeks(g[0]) != "OR|100||||" {
		t.Errorf("Golf tersentuh: %v", g)
	}
}

// TestSpreadingOtomatisGolf - satu baris ORS 10007 100 % per coverage Golf; TSILiability kosong -> TSISpreaded 0 (K73-20),
// TSIGross = TSI × PS / 100; objek okupasi tidak disentuh.
func TestSpreadingOtomatisGolf(t *testing.T) {
	lokasi := lokasiSpreadUji()
	hitungNRAneka(lokasi, d("50"))
	spreadingOtomatisGolf(lokasi, d("50"))
	g := lokasi[0].ObjekGolf[0].Coverages[0].Spreading
	if len(g) != 1 || g[0].TreatyName != "ORS" || barisTeks(g[0]) != "10007|100|0|2|100|" {
		t.Errorf("Golf: %v", g)
	}
	if len(lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading) != 2 {
		t.Error("objek okupasi tersentuh")
	}
}

// TestTotalLokasiAneka24 - cabang 2.4 (tanpa top risk): TSI dijumlah hanya coverage pertama, premi semua; LoL coverage pertama
// -> ClaimEstimation baris (L5701) dan total; coverage kedua ClaimEstimation 0; share beda -> kosong.
func TestTotalLokasiAneka24(t *testing.T) {
	lokasi := lokasiSpreadUji()
	hitungNRAneka(lokasi, d("50"))
	l := liniSpreadingAneka{aneka: true}
	terapkanPersenAneka(lokasi, l)
	total, adaTop := totalLokasiAneka(lokasi, d("50"), l)
	if adaTop || len(total) != 1 || len(total[0]) != 2 {
		t.Fatalf("total %v adaTop %v", total, adaTop)
	}
	if g := totalSpreadTeks(total[0][0]); g != "IDR|OR||300|4.9|240|" {
		t.Errorf("OR: %q", g)
	}
	if g := totalSpreadTeks(total[0][1]); g != "IDR|QS|40|200|1.6|160|" {
		t.Errorf("QS: %q", g)
	}
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	if t8(c[0].Spreading[0].ClaimEstimation) != "240" || t8(c[0].Spreading[1].ClaimEstimation) != "160" || t8(c[1].Spreading[0].ClaimEstimation) != "0" {
		t.Errorf("ClaimEstimation baris: %v", c)
	}
	if total, _ := totalLokasiAneka(lokasi, d("50"), liniSpreadingAneka{golf: true}); len(total[0]) != 0 {
		t.Errorf("bukan Aneka: %v", total)
	}
}

// TestTotalLokasiAneka25 - cabang 2.5 (ada lokasi top risk): LoL hanya lokasi top risk, tidak ditulis ke baris; ClaimSpreaded =
// TSISpreaded coverage pertama lokasi top risk bila lini liniTopRisk25; lokasi bukan top tanpa LoL / ClaimSpreaded.
func TestTotalLokasiAneka25(t *testing.T) {
	lokasi := append(lokasiSpreadUji(), lokasiSpreadUji()...)
	lokasi[0].IsTopRisk = true
	hitungNRAneka(lokasi, d("50"))
	l := liniSpreadingAneka{aneka: true, topRisk25: true}
	terapkanPersenAneka(lokasi, l)
	total, adaTop := totalLokasiAneka(lokasi, d("50"), l)
	if !adaTop {
		t.Fatal("adaTop")
	}
	if g := totalSpreadTeks(total[0][0]); g != "IDR|OR||300|4.9|240|300" {
		t.Errorf("top OR: %q", g)
	}
	if g := totalSpreadTeks(total[1][0]); g != "IDR|OR||300|4.9|0|" {
		t.Errorf("bukan top OR: %q", g)
	}
	if lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0].ClaimEstimation != nil {
		t.Error("ClaimEstimation baris ditulis di cabang 2.5")
	}
	l.topRisk25 = false
	if total, _ := totalLokasiAneka(lokasi, d("50"), l); total[0][0].ClaimSpreaded != nil {
		t.Errorf("lini bukan top risk: %v", total[0][0])
	}
}

// TestKunciMataUangTotal - lini daftar: mata uang coverage (kosong -> objek); lini lain: objek.
func TestKunciMataUangTotal(t *testing.T) {
	lokasi := lokasiSpreadUji()
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	c[0].Currency = "USD"
	hitungNRAneka(lokasi, d("50"))
	total, _ := totalLokasiAneka(lokasi, d("50"), liniSpreadingAneka{aneka: true, mataUangCov24: true})
	var kunci []string
	for _, x := range total[0] {
		kunci = append(kunci, x.Currency+"/"+x.TreatyType)
	}
	if strings.Join(kunci, ",") != "USD/OR,USD/QS,IDR/OR" {
		t.Errorf("kunci: %v", kunci)
	}
	total, _ = totalLokasiAneka(lokasi, d("50"), liniSpreadingAneka{aneka: true})
	if len(total[0]) != 2 {
		t.Errorf("lini umum: %v", total[0])
	}
}

// TestTotalLokasiGolf - TSI dan premi SETIAP coverage dijumlah, mata uang objek, tanpa LoL / top risk.
func TestTotalLokasiGolf(t *testing.T) {
	lokasi := lokasiSpreadUji()
	g := &lokasi[0].ObjekGolf[0]
	g.Coverages = append(g.Coverages, models.CoverageAneka{OldID: "8", TSI: d("100"), Premium: d("2"),
		Spreading: []models.BarisSpreading{{TreatyType: "OR", SharePercentage: d("50")}}})
	hitungNRAneka(lokasi, d("50"))
	l := liniSpreadingAneka{golf: true}
	terapkanPersenAneka(lokasi, l)
	total := totalLokasiGolf(lokasi, l)
	if len(total[0]) != 1 || totalSpreadTeks(total[0][0]) != "USD|OR||125|2.5||" {
		t.Errorf("Golf: %v", total)
	}
	if total := totalLokasiGolf(lokasi, liniSpreadingAneka{aneka: true}); len(total[0]) != 0 {
		t.Errorf("bukan Golf: %v", total)
	}
}

// TestSpreadingOtomatisEngineering - beda Engineering vs FIRE: coverage kedua (TSI RNM 500 > maxQS 300, TSI RNM 100 % 200) dipotong
// dengan share = MaxLimitQSIDR / TSISpreaded sebelum dipotong, lalu SPL ditambahkan walau TSI RNM 500 <= maxtsi 700 (FIRE tidak).
func TestSpreadingOtomatisEngineering(t *testing.T) {
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "IDR", Coverages: []models.CoverageAneka{
		{OldID: "1", TSI: d("400"), TSILiability: d("400"), Premium: d("8")},
		{OldID: "2", TSI: d("2000"), TSILiability: d("1000"), Premium: d("20")},
	}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	pesan, err := spreadingOtomatisEngineering(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, k, lokasi)
	if err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	if len(c[0].Spreading) != 1 || barisTeks(c[0].Spreading[0]) != "Q1|100|200|4|200|" {
		t.Errorf("coverage 1: %v", c[0].Spreading)
	}
	if len(c[1].Spreading) != 2 || barisTeks(c[1].Spreading[0]) != "Q1|60|300|6|600|" || barisTeks(c[1].Spreading[1]) != "S1|40|200|4|400|" {
		t.Errorf("coverage 2: %v", c[1].Spreading)
	}
}

// TestSpreadingOtomatisEngineeringTopRisk - K73-23 (A): TSI RNM top risk diisi seperti cabang IsFire (L2494, coverage pertama
// lokasi top risk) -> 500 > maxQS 300 dipecah QS 60 % / 300 + SPL 40 % / 200 (sebelumnya pembagi 0: QS 0 % / SPL 100 %).
func TestSpreadingOtomatisEngineeringTopRisk(t *testing.T) {
	lokasi := []models.LokasiAneka{{IsTopRisk: true, Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "IDR",
		Coverages: []models.CoverageAneka{{OldID: "1", TSI: d("1000"), TSILiability: d("1000"), Premium: d("10")}}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	if pesan, err := spreadingOtomatisEngineering(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, k, lokasi); err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	s := lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading
	if len(s) != 2 || barisTeks(s[0]) != "Q1|60|300|3|300|" || barisTeks(s[1]) != "S1|40|200|2|200|" {
		t.Errorf("top risk: %v", s)
	}
	if pesan, _ := spreadingOtomatisEngineering(context.Background(), &spreadingTiruan{}, k, lokasi); len(pesan) != 1 ||
		!strings.Contains(pesan[0], "KAPASITAS_TREATY") {
		t.Errorf("tanpa kapasitas: %v", pesan)
	}
}

// TestKasusSpreadingAneka - adaptor hilir: lini umum satu item per objek (TSI objek); lini Count khusus satu item per coverage
// (mata uang coverage, kosong -> objek; TSI coverage hanya liniTSICoverageCount); ringkasanRNM memberi premi / TSI per mata uang.
func TestKasusSpreadingAneka(t *testing.T) {
	lokasi := lokasiSpreadUji()
	k := models.KasusSpreading{PercentShare: d("50")}
	// Objek okupasi hanya bila IsAneka (Count langkah 7), objek Golf hanya bila isGolfInsurance (langkah 8) - tiket 73 P5.
	_, jumlah, dasar := ringkasanRNM(kasusSpreadingAneka(k, lokasi, liniSpreadingAneka{aneka: true, golf: true}))
	var got []string
	for _, j := range jumlah {
		got = append(got, j.Currency+":"+t8(j.Premi)+":"+t8(j.TSIRNM))
	}
	if strings.Join(got, ",") != "IDR:6.5:500,USD:2:100" || len(dasar) != 2 || t8(dasar[0].TSIGross) != "1000" {
		t.Errorf("umum: %v / %v", got, dasar)
	}
	for l, mau := range map[liniSpreadingAneka]string{{aneka: true}: "IDR:6.5:500", {golf: true}: "USD:2:100"} {
		_, jumlah, _ := ringkasanRNM(kasusSpreadingAneka(k, lokasi, l))
		if len(jumlah) != 1 || jumlah[0].Currency+":"+t8(jumlah[0].Premi)+":"+t8(jumlah[0].TSIRNM) != mau {
			t.Errorf("gerbang lini %+v: %v", l, jumlah)
		}
	}
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Currency = "USD"
	for nama, u := range map[string]struct {
		l   liniSpreadingAneka
		mau string
	}{
		"TSI coverage": {liniSpreadingAneka{aneka: true, mataUangCount: true, tsiCoverageCount: true}, "USD:4:500,IDR:2.5:250"},
		"BuilderRisk":  {liniSpreadingAneka{aneka: true, mataUangCount: true}, "USD:4:0,IDR:2.5:0"},
	} {
		_, jumlah, _ := ringkasanRNM(kasusSpreadingAneka(k, lokasi, u.l))
		got = nil
		for _, j := range jumlah {
			got = append(got, j.Currency+":"+t8(j.Premi)+":"+t8(j.TSIRNM))
		}
		if strings.Join(got, ",") != u.mau {
			t.Errorf("%s: %v, mau %s", nama, got, u.mau)
		}
	}
}

// TestSamaBentukAneka - jumlah dan oldId per posisi.
func TestSamaBentukAneka(t *testing.T) {
	a := lokasiSpreadUji()
	if !samaBentukAneka(a, lokasiSpreadUji()) {
		t.Fatal("sama")
	}
	for nama, ubah := range map[string]func(l []models.LokasiAneka) []models.LokasiAneka{
		"lokasi":  func(l []models.LokasiAneka) []models.LokasiAneka { return append(l, l[0]) },
		"okupasi": func(l []models.LokasiAneka) []models.LokasiAneka { l[0].Okupasi = nil; return l },
		"golf":    func(l []models.LokasiAneka) []models.LokasiAneka { l[0].ObjekGolf = nil; return l },
		"coverage": func(l []models.LokasiAneka) []models.LokasiAneka {
			l[0].Okupasi[0].Objek[0].Coverages = l[0].Okupasi[0].Objek[0].Coverages[:1]
			return l
		},
		"oldId": func(l []models.LokasiAneka) []models.LokasiAneka {
			l[0].ObjekGolf[0].Coverages[0].OldID = "X"
			return l
		},
	} {
		if samaBentukAneka(a, ubah(lokasiSpreadUji())) {
			t.Errorf("%s: dianggap sama", nama)
		}
	}
}

// TestSpreadingOtomatisEngineeringPotongPecahan - beda 4: QS hasil pecah (75 %) dipotong pada coverage kedua -> share =
// MaxLimitQSIDR 300 / TSISpreaded sebelum dipotong 375 × 100 = 80 (FIRE: 300 / TSI RNM 500 = 60); SPL ditambahkan untuk setiap
// baris sesudah pemotongan (verbatim For Each, pola FIRE).
func TestSpreadingOtomatisEngineeringPotongPecahan(t *testing.T) {
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "IDR", Coverages: []models.CoverageAneka{
		{OldID: "1", TSI: d("800"), TSILiability: d("800"), Premium: d("8")},
		{OldID: "2", TSI: d("1000"), TSILiability: d("1000"), Premium: d("20")},
	}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	if pesan, err := spreadingOtomatisEngineering(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, k, lokasi); err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	var got []string
	for _, b := range c[1].Spreading {
		got = append(got, b.TreatyType+":"+t8(b.SharePercentage)+":"+t8(b.TSISpreaded))
	}
	if strings.Join(got, ",") != "Q1:80:300,S1:25:125,S1:20:100,S1:20:100" {
		t.Errorf("coverage 2: %v", got)
	}
	if len(c[0].Spreading) != 2 || barisTeks(c[0].Spreading[0]) != "Q1|75|300|3|300|" {
		t.Errorf("coverage 1: %v", c[0].Spreading)
	}
}

// TestSpreadingOtomatisEngineeringKurs - objek USD (kurs 2): maxQS / maxSPL dikurs (150 / 200); coverage pertama TSI RNM 100
// (QS 100 %), coverage kedua TSI RNM 200 > 150 dipotong dengan pembilang Local.maxtsiQS (K73-25 P1): share 150 / 200 × 100 = 75
// (verbatim Pega: MaxLimitQSIDR mentah 300 / 200 = 150), SPL 25 % / 50. Kurs tidak ada -> pesan, spreading tersimpan utuh.
func TestSpreadingOtomatisEngineeringKurs(t *testing.T) {
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "USD", Coverages: []models.CoverageAneka{
		{OldID: "1", TSI: d("200"), TSILiability: d("200"), Premium: d("10")},
		{OldID: "2", TSI: d("400"), TSILiability: d("400"), Premium: d("10")}}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	src := &spreadingTiruan{kap: &kapasitasUji, idUang: map[string]string{"USD": "U1"}, kurs: map[string]*apd.Decimal{"U1": d("2")}}
	if pesan, err := spreadingOtomatisEngineering(context.Background(), src, k, lokasi); err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	c := lokasi[0].Okupasi[0].Objek[0].Coverages
	if len(c[0].Spreading) != 1 || t8(c[0].Spreading[0].SharePercentage) != "100" || t8(c[0].Spreading[0].TSISpreaded) != "100" {
		t.Errorf("coverage 1: %v", c[0].Spreading)
	}
	if s := c[1].Spreading; len(s) != 2 || barisTeks(s[0]) != "Q1|75|150|3.75|150|" || barisTeks(s[1]) != "S1|25|50|1.25|50|" {
		t.Errorf("coverage 2: %v", s)
	}
	lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "LAMA"}}
	if pesan, _ := spreadingOtomatisEngineering(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, k, lokasi); len(pesan) != 1 ||
		!strings.HasPrefix(pesan[0], "Kurs USD ke IDR tidak ada") || lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0].TreatyType != "LAMA" {
		t.Errorf("tanpa kurs: %v", pesan)
	}
}

// layananSpreadAneka - case "K" ber-BUSINESS (003, old); objek Aneka tersimpan lokasiSpreadUji; % Share RNM tersimpan 50.
func layananSpreadAneka(old string) (*Service, anekaTiruan, *[]models.LokasiAneka, **apd.Decimal) {
	var tulis []models.LokasiAneka
	var share *apd.Decimal
	pr := d("12.5")
	a := anekaTiruan{ada: map[string][]models.LokasiAneka{"K": lokasiSpreadUji()}, tulis: &tulis, share: d("50"), shareTulis: &share,
		proRate: &pr}
	k := &models.Kasus{CaseID: "K", Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS", GroupBusinessID: "3"}}
	svc := Baru(nil).DenganObjekAneka(a).DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: satuBisnis("003", old)}).DenganPilihanItem(pilihanItemTiruan{panggil: new(int)}).
		DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT",
			Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}}}}, idUang: map[string]string{"IDR": "10026"}})
	return svc, a, &tulis, &share
}

// badanSpreadUji - badan PUT berbentuk sama dengan lokasiSpreadUji; baris c1 diganti OR 100 %.
func badanSpreadUji() []models.LokasiAneka {
	b := lokasiSpreadUji()
	b[0].Okupasi[0].Objek[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: " OR ", TreatyName: "UJI", SharePercentage: d("100")}}
	return b
}

// TestTampilanSpreadingAneka - GET: % Share RNM, pohon objek dengan baris tersimpan, total per lokasi lini Aneka (objek Golf tidak
// dijumlah untuk case bukan Golf), ringkasan; grid total per lokasi tersembunyi (bukan MBD).
func TestTampilanSpreadingAneka(t *testing.T) {
	svc, _, _, _ := layananSpreadAneka("30")
	v, err := svc.TampilanSpreadingAneka(context.Background(), "K")
	if err != nil {
		t.Fatal(err)
	}
	if t8(v.PercentShare) != "50" || len(v.Lokasi) != 1 || v.Lokasi[0].Location != "JL UJI" || v.TampilTotalLokasi {
		t.Fatalf("tampilan: %+v", v)
	}
	if len(v.Lokasi[0].Total) != 2 || v.Lokasi[0].Total[0].TreatyType != "OR" || len(v.RingkasanMataUang) != 1 ||
		v.RingkasanMataUang[0].Currency != "IDR" {
		t.Errorf("total %v ringkasan %v", v.Lokasi[0].Total, v.RingkasanMataUang)
	}
	if len(v.Lokasi[0].ObjekGolf) != 1 || len(v.Lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading) != 2 || v.Pesan == nil {
		t.Errorf("pohon: %+v", v.Lokasi[0])
	}
	svc, _, _, _ = layananSpreadAneka("10")
	if v, _ := svc.TampilanSpreadingAneka(context.Background(), "K"); !v.TampilTotalLokasi {
		t.Error("MBD: grid total per lokasi tidak tampil")
	}
	if _, err := Baru(nil).TampilanSpreadingAneka(context.Background(), "K"); !errors.Is(err, ErrSpreadingTanpaDatabase) {
		t.Errorf("tanpa DB: %v", err)
	}
}

// TestSimpanSpreadingAneka - PUT (tanpa pembaca tahap = baris badan dipakai): % Share RNM + pohon ditulis dalam transaksi,
// % Pro Rate tersimpan tidak diubah, NR / DT / ClaimEstimation dihitung ulang 8 desimal; badan tidak diubah.
func TestSimpanSpreadingAneka(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	svc, a, tulis, share := layananSpreadAneka("30")
	badan := badanSpreadUji()
	if _, err := svc.SimpanSpreadingAneka(ctx, akun, "K", d("40"), badan); err != nil {
		t.Fatal(err)
	}
	if t8(*share) != "40" || t8(*a.proRate) != "12.5" {
		t.Errorf("share %s proRate %s", t8(*share), t8(*a.proRate))
	}
	c := (*tulis)[0].Okupasi[0].Objek[0].Coverages
	if len(c[0].Spreading) != 1 || barisTeks(c[0].Spreading[0]) != "OR|100|400|3.2||320" || c[0].Spreading[0].TreatyName != "UJI" ||
		t8(c[0].TSINusantaraRe) != "400" || barisTeks(c[1].Spreading[0]) != "OR|100|200|2||0" {
		t.Errorf("tulis: %v", c)
	}
	if badan[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0].TreatyType != " OR " {
		t.Error("badan diubah")
	}
	beda := badanSpreadUji()
	beda[0].ObjekGolf[0].Coverages[0].OldID = "X"
	if _, err := svc.SimpanSpreadingAneka(ctx, akun, "K", d("40"), beda); !errors.Is(err, ErrSpreadingBerubah) {
		t.Errorf("bentuk: %v", err)
	}
	lebih := badanSpreadUji()
	lebih[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0].SharePercentage = d("101")
	if _, err := svc.SimpanSpreadingAneka(ctx, akun, "K", d("40"), lebih); !errors.Is(err, ErrMasukanSpreading) {
		t.Errorf("share 101: %v", err)
	}
	if _, err := svc.SimpanSpreadingAneka(ctx, akun, "K", d("101"), badan); !errors.Is(err, ErrMasukanSpreading) {
		t.Errorf("percentShare 101: %v", err)
	}
	if _, err := svc.SimpanSpreadingAneka(ctx, inti.Pelaku{}, "K", d("40"), badan); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
	a.ada["K"] = []models.LokasiAneka{{ObjectNo: "1"}}
	if _, err := svc.SimpanSpreadingAneka(ctx, akun, "K", d("40"), []models.LokasiAneka{{}}); !errors.Is(err, ErrMasukanSpreading) ||
		!strings.Contains(err.Error(), pesanObjekKosong) {
		t.Errorf("objek kosong: %v", err)
	}
}

// TestHitungShareSpreadingAneka - tanpa pembaca tahap (= uw): Golf -> spreading otomatis ORS; Aneka -> baris tersimpan dihitung
// ulang dengan % Share baru; tidak menulis.
func TestHitungShareSpreadingAneka(t *testing.T) {
	ctx := context.Background()
	svc, _, tulis, _ := layananSpreadAneka("27")
	v, err := svc.HitungShareSpreadingAneka(ctx, "K", d("40"))
	if err != nil {
		t.Fatal(err)
	}
	if g := v.Lokasi[0].ObjekGolf[0].Coverages[0].Spreading; len(g) != 1 || g[0].TreatyType != "10007" || len(*tulis) != 0 {
		t.Errorf("Golf: %v / tulis %d", g, len(*tulis))
	}
	if len(v.Lokasi[0].Total) != 1 || v.Lokasi[0].Total[0].Currency != "USD" {
		t.Errorf("total Golf: %v", v.Lokasi[0].Total)
	}
	svc, _, _, _ = layananSpreadAneka("30")
	v, _ = svc.HitungShareSpreadingAneka(ctx, "K", d("40"))
	if b := v.Lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0]; t8(b.TSISpreaded) != "240" || t8(v.PercentShare) != "40" {
		t.Errorf("Aneka: %v", b)
	}
}

// TestBacaSpreadingRingkas - pohon FIRE berisi item -> apa adanya (objek Aneka tidak dibaca); tanpa item -> adaptor Aneka.
func TestBacaSpreadingRingkas(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _ := layananSpreadAneka("30")
	k, err := svc.bacaSpreadingRingkas(ctx, nil, "K")
	if err != nil || len(k.Lokasi) != 1 || len(k.Lokasi[0].Items) != 1 || k.Lokasi[0].Items[0].Currency != "IDR" || t8(k.PercentShare) != "50" {
		t.Fatalf("Aneka: %+v %v", k, err)
	}
	fire := kasusSatuCoverage("50")
	svc.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{fire}})
	if k, _ := svc.bacaSpreadingRingkas(ctx, nil, "K"); len(k.Lokasi[0].Items) != 1 || k.Lokasi[0].Items[0].ItemType != "BUILDING" {
		t.Errorf("FIRE: %+v", k)
	}
}

// TestGantiObjekAnekaHitungUlangSpreading - Save tab Object Aneka dengan % Share RNM tersimpan: spreading coverage yang dikirim
// ulang dipertahankan dan dihitung ulang (DT percent + ClaimEstimation 2.4), 8 desimal.
func TestGantiObjekAnekaHitungUlangSpreading(t *testing.T) {
	svc, _, tulis, _ := layananSpreadAneka("30")
	masuk := lokasiSpreadUji()
	masuk[0].ObjekGolf = nil
	for i := range masuk[0].Okupasi[0].Objek[0].Coverages {
		masuk[0].Okupasi[0].Objek[0].Coverages[i].Spreading = nil
	}
	if _, err := svc.GantiObjekAneka(context.Background(), inti.Pelaku{AkunID: "UJI"}, "K", masuk, nil); err != nil {
		t.Fatal(err)
	}
	c := (*tulis)[0].Okupasi[0].Objek[0].Coverages
	if barisTeks(c[0].Spreading[0]) != "OR|60|300|2.4|77|240" || barisTeks(c[1].Spreading[0]) != "OR|100|250|2.5||0" {
		t.Errorf("spreading: %v", c)
	}
}

// TestHapusPesanMBD - langkah 24 (K73-22): limit GetTreatyName (ORS 150 M / 181,5 M menurut Begin date) <= TSISpreaded total
// lokasi berjenis sama -> true; limit tak terurai tidak dibandingkan.
func TestHapusPesanMBD(t *testing.T) {
	k := models.KasusSpreading{StartDateTime: "20250101T170000.000 GMT"}
	lok := func(jenis, tsi string) []models.LokasiSpreading {
		return []models.LokasiSpreading{{Total: []models.TotalSpreading{{TreatyType: jenis, TSISpreaded: d(tsi)}}}}
	}
	batas := []models.TreatySpreading{{ID: "Q1", Limit: "300"}, {ID: "X", Limit: "abc"}, {ID: treatyORS, Limit: limitORS}}
	for nama, u := range map[string]struct {
		k   models.KasusSpreading
		l   []models.LokasiSpreading
		mau bool
	}{
		"sama dengan limit": {k, lok("Q1", "300"), true},
		"di bawah limit":    {k, lok("Q1", "299.99"), false},
		"jenis lain":        {k, lok("Q2", "999"), false},
		"limit bukan angka": {k, lok("X", "999"), false},
		"ORS 150 M":         {k, lok(treatyORS, "150000000000"), true},
		"ORS sesudah Juli":  {models.KasusSpreading{StartDateTime: "20260630T170000.000 GMT"}, lok(treatyORS, "150000000000"), false},
		"ORS 181,5 M":       {models.KasusSpreading{StartDateTime: "20260630T170000.000 GMT"}, lok(treatyORS, "181500000000"), true},
	} {
		if got := hapusPesanMBD(u.k, u.l, batas); got != u.mau {
			t.Errorf("%s: %v", nama, got)
		}
	}
}

// TestPesanMBDDihapus - GET: TSI ORS 400 M melampaui batas ORS -> pesan batas treaty; case MBD -> langkah 24 menghapusnya.
func TestPesanMBDDihapus(t *testing.T) {
	for old, kosong := range map[string]bool{"30": false, "10": true} {
		svc, a, _, _ := layananSpreadAneka(old)
		a.ada["K"][0].Okupasi[0].Objek[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: treatyORS, SharePercentage: d("100"),
			TSISpreaded: d("400000000000")}}
		v, err := svc.TampilanSpreadingAneka(context.Background(), "K")
		if err != nil {
			t.Fatal(err)
		}
		if (len(v.Pesan) == 0) != kosong || (!kosong && !strings.Contains(v.Pesan[0], "Total TSI Spreaded for ORS")) {
			t.Errorf("BUSINESS %s: pesan %v", old, v.Pesan)
		}
	}
}

// TestSpreadingOtomatisEngineeringSisaPresisi - K73-26 (c) + K73-27 (ii): TSI RNM 700, maxQS 300 / maxSPL 400 -> pecah QS 300 / 700
// = 42,857142857142857143 %; TSISpreaded 300,000…001 tidak memicu pemotongan (dibandingkan 8 desimal), dan baris SPL 57,142857 %
// (TSI 400 > maxQS 300) TIDAK dipotong karena pemicu .19.3 hanya untuk baris QS: QS 300 + SPL 400 = 700, tanpa baris tambahan.
// Verbatim Pega (L19342 tanpa saringan jenis): SPL 75 % / 300 + SPL "beda 5" 25 % / 175.
func TestSpreadingOtomatisEngineeringSisaPresisi(t *testing.T) {
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "IDR", Coverages: []models.CoverageAneka{
		{OldID: "1", TSI: d("1400"), TSILiability: d("1400"), Premium: d("10")}}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	if pesan, err := spreadingOtomatisEngineering(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, k, lokasi); err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	var got []string
	for _, b := range lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading {
		got = append(got, barisTeks(b))
	}
	if strings.Join(got, ",") != "Q1|42.85714286|300|2.14285714|300|,S1|57.14285714|400|2.85714286|400|" {
		t.Errorf("baris: %v", got)
	}
}

// TestSpreadingOtomatisEngineeringTopRiskKurs - K73-23 (A) objek USD (kurs 2) di lokasi top risk: tsitoprisk = 500 × 2 = 1000
// (IDR, L2441) dipakai mencari KAPASITAS_TREATY; TSIRNMTopRisk1 = 1000 × 50 % = 500 (IDR, L2494) -> TSIRNMTopRisk = 500 / 2 = 250
// (L16245) > maxQS 150 -> pecah QS 60 % / 150 + SPL 40 % / 100.
func TestSpreadingOtomatisEngineeringTopRiskKurs(t *testing.T) {
	lokasi := []models.LokasiAneka{{IsTopRisk: true, Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "USD",
		Coverages: []models.CoverageAneka{{OldID: "1", TSI: d("500"), TSILiability: d("500"), Premium: d("10")}}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	src := &spreadingTiruan{kap: &kapasitasUji, idUang: map[string]string{"USD": "U1"}, kurs: map[string]*apd.Decimal{"U1": d("2")}}
	if pesan, err := spreadingOtomatisEngineering(context.Background(), src, k, lokasi); err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	s := lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading
	if len(s) != 2 || barisTeks(s[0]) != "Q1|60|150|3|150|" || barisTeks(s[1]) != "S1|40|100|2|100|" {
		t.Errorf("baris: %v", s)
	}
	if strings.Join(src.tsiKap, ",") != "1000" {
		t.Errorf("TSI kapasitas: %v", src.tsiKap)
	}
}

// TestSpreadingOtomatisEngineeringBatasDibulatkan - K73-26: maxtsiQS juga dibulatkan 8 desimal sebelum dibandingkan. Objek USD
// kurs 7: maxQS = 300 / 7 = 42,857142857142857143 (skala 20) -> 42,85714286; TSI RNM 42,857142858 sama sesudah pembulatan, jadi
// baris QS tidak dipotong dan tidak ada SPL "beda 5" tambahan (dua baris hasil pecah).
func TestSpreadingOtomatisEngineeringBatasDibulatkan(t *testing.T) {
	lokasi := []models.LokasiAneka{{Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{Currency: "USD", Coverages: []models.CoverageAneka{
		{OldID: "1", TSI: d("85.714285716"), TSILiability: d("85.714285716"), Premium: d("10")}}}}}}}}
	hitungNRAneka(lokasi, d("50"))
	k := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT"}
	src := &spreadingTiruan{kap: &kapasitasUji, idUang: map[string]string{"USD": "U1"}, kurs: map[string]*apd.Decimal{"U1": d("7")}}
	if pesan, err := spreadingOtomatisEngineering(context.Background(), src, k, lokasi); err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	if s := lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading; len(s) != 2 || s[0].TreatyType != "Q1" || s[1].TreatyType != "S1" {
		t.Errorf("baris: %v", s)
	}
}
