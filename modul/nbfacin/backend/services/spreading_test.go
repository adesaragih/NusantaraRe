package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// spreadingTiruan - PenyimpanSpreading dalam memori. baca dikembalikan berurutan (yang terakhir diulang).
type spreadingTiruan struct {
	baca      []models.KasusSpreading
	kali      int
	tulis     []models.KasusSpreading
	ubahShare []bool
	treaty    []models.TreatySpreading
	kap       *models.KapasitasTreaty
	kurs      map[string]*apd.Decimal // id mata uang -> kurs
	idUang    map[string]string       // nama -> id
}

func (s *spreadingTiruan) BacaSpreading(context.Context, *db.Tx, string) (models.KasusSpreading, error) {
	i := s.kali
	if i >= len(s.baca) {
		i = len(s.baca) - 1
	}
	s.kali++
	return salinKasusSpreading(s.baca[i]), nil
}

func (s *spreadingTiruan) TulisSpreading(_ context.Context, _ *db.Tx, _ string, k models.KasusSpreading, ubah bool) error {
	s.tulis, s.ubahShare = append(s.tulis, k), append(s.ubahShare, ubah)
	return nil
}

func (s *spreadingTiruan) DaftarTreaty(context.Context, string, string) ([]models.TreatySpreading, error) {
	return s.treaty, nil
}

func (s *spreadingTiruan) KapasitasTreaty(context.Context, *apd.Decimal, string) (models.KapasitasTreaty, bool, error) {
	if s.kap == nil {
		return models.KapasitasTreaty{}, false, nil
	}
	return *s.kap, true, nil
}

func (s *spreadingTiruan) KursTerbaru(_ context.Context, id string) (*apd.Decimal, error) {
	return s.kurs[id], nil
}

func (s *spreadingTiruan) KursPada(_ context.Context, id, _ string) (*apd.Decimal, error) {
	return s.kurs[id], nil
}

func (s *spreadingTiruan) IDMataUang(_ context.Context, nama string) (string, error) {
	return s.idUang[nama], nil
}

// salinKasusSpreading - salinan dalam (tiruan tidak berbagi slice dengan pemanggil).
func salinKasusSpreading(k models.KasusSpreading) models.KasusSpreading {
	out := k
	out.Lokasi = make([]models.LokasiSpreading, len(k.Lokasi))
	for i, l := range k.Lokasi {
		out.Lokasi[i] = l
		out.Lokasi[i].Items = make([]models.ItemSpreading, len(l.Items))
		for j, it := range l.Items {
			out.Lokasi[i].Items[j] = it
			out.Lokasi[i].Items[j].Coverages = make([]models.CoverageSpreading, len(it.Coverages))
			for m, c := range it.Coverages {
				c.Spreading = append([]models.BarisSpreading{}, c.Spreading...)
				out.Lokasi[i].Items[j].Coverages[m] = c
			}
		}
	}
	return out
}

func t8(x *apd.Decimal) string { return utils.FormatDecimal(simpan8(x)) }

// kasusSatuCoverage - satu lokasi, satu item IDR, satu coverage (TSILiability 1000, TSI 1200, premi 10, diskon 2).
func kasusSatuCoverage(ps string) models.KasusSpreading {
	var share *apd.Decimal
	if ps != "" {
		share = d(ps)
	}
	return models.KasusSpreading{PercentShare: share, StartDateTime: "20250101T170000.000 GMT", Lokasi: []models.LokasiSpreading{{
		ObjectNo: "1", Items: []models.ItemSpreading{{ItemType: "BUILDING", Currency: "IDR", Coverages: []models.CoverageSpreading{{
			ID: "11", OldID: "C1", TSI: d("1200"), TSILiability: d("1000"), LimitOfLiability: d("800"), Premium: d("10"), Discount: d("2"),
		}}}},
	}}}
}

var kapasitasUji = models.KapasitasTreaty{TreatyNameQS: "QS UJI", IDTreatyQS: "Q1", TreatyNameSPL: "SPL UJI", IDTreatySPL: "S1",
	MaxLimitQSIDR: d("300"), MaxLimitSPLIDR: d("400")}

func TestHitungNR(t *testing.T) {
	k := kasusSatuCoverage("50")
	hitungNR(&k)
	c := k.Lokasi[0].Items[0].Coverages[0]
	if t8(c.TSINusantaraRe) != "500" || t8(c.PremiNusantaraRe) != "4" || t8(k.Lokasi[0].Items[0].TotalPremiumNusantaraRe) != "4" {
		t.Fatalf("NR = %s / %s / %s", t8(c.TSINusantaraRe), t8(c.PremiNusantaraRe), t8(k.Lokasi[0].Items[0].TotalPremiumNusantaraRe))
	}
}

// TestSpreadingOtomatisQSSPL - TSI RNM 500 > maxQS 300 (maxtsi 700): QS 60 % (300) + SPL 40 % (200).
func TestSpreadingOtomatisQSSPL(t *testing.T) {
	k := kasusSatuCoverage("50")
	hitungNR(&k)
	pesan, err := spreadingOtomatis(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, &k)
	if err != nil || len(pesan) != 0 {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	s := k.Lokasi[0].Items[0].Coverages[0].Spreading
	if len(s) != 2 {
		t.Fatalf("baris = %d", len(s))
	}
	mau := [][]string{{"Q1", "QS UJI", "60", "300", "2.4", "360"}, {"S1", "SPL UJI", "40", "200", "1.6", "240"}}
	for i, m := range mau {
		b := s[i]
		got := []string{b.TreatyType, b.TreatyName, t8(b.SharePercentage), t8(b.TSISpreaded), t8(b.PremiumSpreaded), t8(b.TSIGrossSpreaded)}
		if strings.Join(got, "|") != strings.Join(m, "|") {
			t.Errorf("baris %d = %v, mau %v", i, got, m)
		}
	}
}

// TestSpreadingOtomatisQSSaja - TSI RNM 250 <= maxQS: satu baris QS 100 %.
func TestSpreadingOtomatisQSSaja(t *testing.T) {
	k := kasusSatuCoverage("25")
	hitungNR(&k)
	if _, err := spreadingOtomatis(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, &k); err != nil {
		t.Fatal(err)
	}
	s := k.Lokasi[0].Items[0].Coverages[0].Spreading
	if len(s) != 1 || s[0].TreatyType != "Q1" || t8(s[0].SharePercentage) != "100" || t8(s[0].TSISpreaded) != "250" {
		t.Fatalf("spreading = %+v", s)
	}
}

func TestSpreadingOtomatisMelebihiKapasitas(t *testing.T) {
	k := kasusSatuCoverage("100") // TSI RNM 1000 > maxtsi 700
	hitungNR(&k)
	pesan, err := spreadingOtomatis(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, &k)
	if err != nil || len(pesan) != 1 || pesan[0] != "TSI in Location 1 exceeds treaty capacity" {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
}

func TestSpreadingOtomatisTanpaKapasitas(t *testing.T) {
	k := kasusSatuCoverage("50")
	pesan, err := spreadingOtomatis(context.Background(), &spreadingTiruan{}, &k)
	if err != nil || len(pesan) != 1 || !strings.Contains(pesan[0], "KAPASITAS_TREATY") {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
}

func TestSalinTemplate(t *testing.T) {
	k := kasusSatuCoverage("50")
	hitungNR(&k)
	salinTemplate(&k, []models.TemplateSpreading{{TreatyType: "T1", TreatyName: "A", SharePercentage: d("70")},
		{TreatyType: "T2", TreatyName: "B", SharePercentage: d("30")}})
	s := k.Lokasi[0].Items[0].Coverages[0].Spreading
	if len(s) != 2 || t8(s[0].TSISpreaded) != "350" || t8(s[0].PremiumSpreaded) != "2.8" || t8(s[0].TSIGrossSpreaded) != "420" ||
		t8(s[1].TSISpreaded) != "150" {
		t.Fatalf("spreading = %+v", s)
	}
}

// kasusTotal - dua coverage IDR ber-spreading T1 100 %; coverage kedua TERRORISM & SABOTAGE.
func kasusTotal() models.KasusSpreading {
	return models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20250101T170000.000 GMT", Lokasi: []models.LokasiSpreading{{
		Items: []models.ItemSpreading{{Currency: "IDR", Coverages: []models.CoverageSpreading{
			{OldID: "C1", TSI: d("1000"), TSILiability: d("1000"), LimitOfLiability: d("800"),
				Spreading: []models.BarisSpreading{{TreatyType: "T1", SharePercentage: d("100"), TSISpreaded: d("500"), PremiumSpreaded: d("4")}}},
			{OldID: "C2", CoverageNote: "TERRORISM & SABOTAGE", TSI: d("200"), TSILiability: d("200"), LimitOfLiability: d("200"),
				Spreading: []models.BarisSpreading{{TreatyType: "T1", SharePercentage: d("100"), TSISpreaded: d("100"), PremiumSpreaded: d("1")}}},
		}}},
	}}}
}

func TestTotalLokasiDanMataUang(t *testing.T) {
	k := kasusTotal()
	adaTop := totalLokasi(&k)
	if adaTop {
		t.Fatal("tidak ada top risk")
	}
	tl := k.Lokasi[0].Total
	if len(tl) != 1 || t8(tl[0].TSISpreaded) != "600" || t8(tl[0].ClaimEstimation) != "500" || t8(tl[0].PremiumSpreaded) != "5" ||
		tl[0].Currency != "IDR" || t8(tl[0].SharePercentage) != "100" {
		t.Fatalf("total lokasi = %+v", tl)
	}
	if ce := k.Lokasi[0].Items[0].Coverages[0].Spreading[0].ClaimEstimation; t8(ce) != "400" {
		t.Fatalf("ClaimEstimation baris = %s", t8(ce))
	}
	src := &spreadingTiruan{idUang: map[string]string{"IDR": idMataUangIDR}}
	pm := totalMataUang(k, adaTop)
	batas := []models.TreatySpreading{{ID: "T1", Name: "TREATY A", Limit: "550"}}
	treaty, semua, pesan, err := ringkasan(context.Background(), src, k, pm, []models.TreatySpreading{{ID: "T1", Name: "TREATY A"}}, batas)
	if err != nil {
		t.Fatal(err)
	}
	if len(treaty) != 1 || treaty[0].TreatyName != "TREATY A" {
		t.Fatalf("ringkasan treaty = %+v", treaty)
	}
	if len(semua) != 1 || t8(semua[0].TSI) != "600" || t8(semua[0].LoL) != "500" || t8(semua[0].Premium) != "5" {
		t.Fatalf("ringkasan mata uang = %+v", semua)
	}
	if len(pesan) != 1 || pesan[0] != "Total TSI Spreaded for TREATY A can't be more than 550" {
		t.Fatalf("pesan = %v", pesan)
	}
}

// TestTotalTerorismeCoveragePertama - coverage pertama bercatatan terorisme yang kuncinya sudah ada dijumlah dua kali
// (.4.8.2 + .4.8.3, verbatim); baris baru (.4.9) sekali.
func TestTotalTerorismeCoveragePertama(t *testing.T) {
	k := kasusTotal()
	k.Lokasi[0].Items = append(k.Lokasi[0].Items, models.ItemSpreading{Currency: "IDR", Coverages: []models.CoverageSpreading{{
		CoverageNote: catatanTerorisme, TSI: d("100"), TSILiability: d("100"),
		Spreading: []models.BarisSpreading{{TreatyType: "T1", SharePercentage: d("100"), TSISpreaded: d("50")}}}}})
	totalLokasi(&k)
	if got := t8(k.Lokasi[0].Total[0].TSISpreaded); got != "700" {
		t.Fatalf("TSI total = %s, mau 700 (500 + 100 + 50×2)", got)
	}
}

func TestDaftarTreaty(t *testing.T) {
	master := []models.TreatySpreading{{ID: "1", Name: "A", Limit: "10"}, {ID: "2", Name: "SF-HRE X", Limit: "5"},
		{ID: "3", Name: "A", Limit: "20"}, {ID: "4", Name: "B TRT", Limit: "1"}}
	dropdown, batas := daftarTreaty(master, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), true)
	ids := func(t []models.TreatySpreading) string {
		var s []string
		for _, x := range t {
			s = append(s, x.ID)
		}
		return strings.Join(s, ",")
	}
	if ids(dropdown) != "3,10007,10015" || ids(batas) != "3,4,10007" {
		t.Fatalf("dropdown %s batas %s", ids(dropdown), ids(batas))
	}
	_, batas = daftarTreaty(master, time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), true)
	if ids(batas) != "2,3,4,10007" {
		t.Fatalf("batas dalam jendela SF-HRE = %s", ids(batas))
	}
}

func TestPeriksaJenisTreaty(t *testing.T) {
	batas := []models.TreatySpreading{{ID: "T1"}}
	if p := periksaJenisTreaty([]models.TotalSpreading{{TreatyType: "T1"}, {TreatyType: "10015"}}, batas); len(p) != 0 {
		t.Fatalf("pesan %v", p)
	}
	if p := periksaJenisTreaty([]models.TotalSpreading{{TreatyType: "X"}}, batas); len(p) != 1 || p[0] != pesanJenisTreaty {
		t.Fatalf("pesan %v", p)
	}
}

func layananSpreading(src *spreadingTiruan) *Service {
	return Baru(nil).DenganSpreading(src).DenganTransaksi(tanpaTx)
}

var pelakuUji = inti.Pelaku{AkunID: "uji"}

func TestSimpanSpreading(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	badan := []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{{OldID: "C1",
		Spreading: []models.BarisSpreading{{TreatyType: " 10015 ", TreatyName: "FACOUT", SharePercentage: d("40")}}}}}}}}
	if _, err := layananSpreading(src).SimpanSpreading(context.Background(), pelakuUji, "NB-1", d("50"), badan); err != nil {
		t.Fatal(err)
	}
	if len(src.tulis) != 1 || !src.ubahShare[0] {
		t.Fatalf("tulis %d ubahShare %v", len(src.tulis), src.ubahShare)
	}
	c := src.tulis[0].Lokasi[0].Items[0].Coverages[0]
	b := c.Spreading[0]
	if t8(c.TSINusantaraRe) != "500" || b.TreatyType != "10015" || t8(b.TSISpreaded) != "200" || t8(b.PremiumSpreaded) != "1.6" ||
		t8(b.TSIGrossSpreaded) != "240" || t8(b.ClaimEstimation) != "160" {
		t.Fatalf("ditulis %+v / %+v", c, b)
	}
}

func TestSimpanSpreadingBentukBerubah(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}}
	badan := []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{{OldID: "LAIN"}}}}}}
	if _, err := layananSpreading(src).SimpanSpreading(context.Background(), pelakuUji, "NB-1", d("50"), badan); !errors.Is(err, ErrSpreadingBerubah) {
		t.Fatalf("err = %v", err)
	}
	if _, err := layananSpreading(src).SimpanSpreading(context.Background(), pelakuUji, "NB-1", d("50"), nil); !errors.Is(err, ErrSpreadingBerubah) {
		t.Fatalf("err = %v", err)
	}
}

func TestSimpanSpreadingShareLebih100(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}}
	badan := []models.LokasiSpreading{{Items: []models.ItemSpreading{{Coverages: []models.CoverageSpreading{{OldID: "C1",
		Spreading: []models.BarisSpreading{{TreatyType: "A", SharePercentage: d("60")}, {TreatyType: "B", SharePercentage: d("50")}}}}}}}}
	_, err := layananSpreading(src).SimpanSpreading(context.Background(), pelakuUji, "NB-1", d("50"), badan)
	if !errors.Is(err, ErrMasukanSpreading) || !strings.Contains(err.Error(), pesanShareLebih100) || len(src.tulis) != 0 {
		t.Fatalf("err = %v tulis %d", err, len(src.tulis))
	}
	if _, err := layananSpreading(src).SimpanSpreading(context.Background(), pelakuUji, "NB-1", d("101"), badan); !errors.Is(err, ErrMasukanSpreading) {
		t.Fatalf("percentShare 101: err = %v", err)
	}
}

func TestSalinSpreadingTidakMenyimpan(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	tpl := []models.TemplateSpreading{{TreatyType: "Q1", TreatyName: "QS", SharePercentage: d("50")}, {TreatyType: "X9", TreatyName: "LAIN", SharePercentage: d("50")}}
	got, err := layananSpreading(src).SalinSpreading(context.Background(), "NB-1", d("50"), tpl)
	if err != nil {
		t.Fatal(err)
	}
	if len(src.tulis) != 0 || len(got.Template) != 2 || len(got.Pesan) != 1 || got.Pesan[0] != pesanJenisTreaty {
		t.Fatalf("tulis %d template %d pesan %v", len(src.tulis), len(got.Template), got.Pesan)
	}
	if _, err := layananSpreading(src).SalinSpreading(context.Background(), "NB-1", d("50"), nil); !errors.Is(err, ErrMasukanSpreading) {
		t.Fatalf("template kosong: err = %v", err)
	}
}

func TestHitungShareSpreading(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("")}, kap: &kapasitasUji,
		idUang: map[string]string{"IDR": idMataUangIDR}}
	got, err := layananSpreading(src).HitungShareSpreading(context.Background(), "NB-1", d("50"))
	if err != nil {
		t.Fatal(err)
	}
	c := got.Lokasi[0].Items[0].Coverages[0]
	if len(src.tulis) != 0 || t8(got.PercentShare) != "50" || len(c.Spreading) != 2 || t8(c.TSINusantaraRe) != "500" {
		t.Fatalf("tulis %d tampilan %+v", len(src.tulis), c)
	}
	if len(got.RingkasanMataUang) != 1 || t8(got.RingkasanMataUang[0].TSI) != "500" {
		t.Fatalf("ringkasan %+v", got.RingkasanMataUang)
	}
}

func TestSpreadingTanpaDatabase(t *testing.T) {
	if _, err := Baru(nil).TampilanSpreading(context.Background(), "NB-1"); !errors.Is(err, ErrSpreadingTanpaDatabase) {
		t.Fatalf("err = %v", err)
	}
}

// TestPertahankanSpreading - Save Object: coverage ber-oldId sama di posisi sama mendapat spreading lamanya (dihitung ulang
// dari TSI baru); coverage baru kosong; % Share RNM tidak ditulis.
func TestPertahankanSpreading(t *testing.T) {
	lama := kasusSatuCoverage("50")
	lama.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "T1", SharePercentage: d("100")}}
	baru := kasusSatuCoverage("50")
	baru.Lokasi[0].Items[0].Coverages[0].TSILiability = d("2000")
	baru.Lokasi[0].Items[0].Coverages = append(baru.Lokasi[0].Items[0].Coverages, models.CoverageSpreading{OldID: "C9"})
	src := &spreadingTiruan{baca: []models.KasusSpreading{baru}}
	if err := layananSpreading(src).pertahankanSpreading(context.Background(), nil, "NB-1", lama); err != nil {
		t.Fatal(err)
	}
	if len(src.tulis) != 1 || src.ubahShare[0] {
		t.Fatalf("tulis %d ubahShare %v", len(src.tulis), src.ubahShare)
	}
	cv := src.tulis[0].Lokasi[0].Items[0].Coverages
	if len(cv[0].Spreading) != 1 || t8(cv[0].Spreading[0].TSISpreaded) != "1000" || len(cv[1].Spreading) != 0 {
		t.Fatalf("coverage %+v", cv)
	}
	// Tanpa isi spreading lama: tidak menulis.
	src2 := &spreadingTiruan{baca: []models.KasusSpreading{baru}}
	kosong := kasusSatuCoverage("50")
	kosong.PercentShare = nil
	if err := layananSpreading(src2).pertahankanSpreading(context.Background(), nil, "NB-1", kosong); err != nil || len(src2.tulis) != 0 {
		t.Fatalf("err %v tulis %d", err, len(src2.tulis))
	}
}

// TestGantiObjekMempertahankanSpreading - Save Object: spreading dibaca di dalam transaksi sebelum objek ditulis, lalu
// dipasang ulang (baca kedua = pohon baru) dan ditulis tanpa mengubah % Share RNM.
func TestGantiObjekMempertahankanSpreading(t *testing.T) {
	lama := kasusSatuCoverage("50")
	lama.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "T1", SharePercentage: d("100")}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{lama, kasusSatuCoverage("50")}}
	svc := Baru(nil).DenganObjek(objekTiruan{ada: map[string][]models.ObjekFire{"UJI-NB-1": {}}}).DenganTransaksi(tanpaTx).
		DenganSpreading(src)
	if _, err := svc.GantiObjek(context.Background(), pelakuUji, "UJI-NB-1", []models.ObjekFire{}); err != nil {
		t.Fatal(err)
	}
	if src.kali != 2 || len(src.tulis) != 1 || src.ubahShare[0] {
		t.Fatalf("baca %d tulis %d ubahShare %v", src.kali, len(src.tulis), src.ubahShare)
	}
	if s := src.tulis[0].Lokasi[0].Items[0].Coverages[0].Spreading; len(s) != 1 || t8(s[0].TSISpreaded) != "500" {
		t.Fatalf("spreading = %+v", s)
	}
}

// TestSpreadingOtomatisTanpaKursTidakMengubah - satu item non-IDR tanpa kurs: spreading tersimpan seluruh lokasi utuh.
func TestSpreadingOtomatisTanpaKursTidakMengubah(t *testing.T) {
	k := kasusSatuCoverage("50")
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "LAMA", SharePercentage: d("100")}}
	k.Lokasi = append(k.Lokasi, models.LokasiSpreading{Items: []models.ItemSpreading{{Currency: "USD",
		Coverages: []models.CoverageSpreading{{OldID: "C2", TSILiability: d("1")}}}}})
	pesan, err := spreadingOtomatis(context.Background(), &spreadingTiruan{kap: &kapasitasUji}, &k)
	if err != nil || len(pesan) != 1 || !strings.Contains(pesan[0], "Kurs USD") {
		t.Fatalf("pesan %v err %v", pesan, err)
	}
	if s := k.Lokasi[0].Items[0].Coverages[0].Spreading; len(s) != 1 || s[0].TreatyType != "LAMA" {
		t.Fatalf("spreading lokasi 1 berubah: %+v", s)
	}
}

// TestRingkasanNamaTreatyTerbawa - jenis treaty tanpa nama di dropdown memakai nama baris sebelumnya (19.5.5, verbatim);
// batas dijumlah dengan pembulatan 4 desimal tiap baris (20.5).
func TestRingkasanNamaTreatyTerbawa(t *testing.T) {
	src := &spreadingTiruan{idUang: map[string]string{"IDR": idMataUangIDR}}
	pm := []models.TotalSpreading{{Currency: "IDR", TreatyType: "T1", TSISpreaded: d("0.00004")},
		{Currency: "IDR", TreatyType: "X9", TSISpreaded: d("1")}, {Currency: "IDR", TreatyType: "T1", TSISpreaded: d("0.00004")}}
	treaty, _, pesan, err := ringkasan(context.Background(), src, models.KasusSpreading{}, pm,
		[]models.TreatySpreading{{ID: "T1", Name: "A"}}, []models.TreatySpreading{{ID: "T1", Name: "A", Limit: "0.00005"}})
	if err != nil {
		t.Fatal(err)
	}
	if treaty[1].TreatyName != "A" {
		t.Fatalf("nama terbawa = %q", treaty[1].TreatyName)
	}
	if len(pesan) != 0 { // 0.00004 -> 0 tiap baris; tanpa pembulatan per baris 0.00008 > 0.00005
		t.Fatalf("pesan = %v", pesan)
	}
}

func TestSimpanSpreadingTanpaObjek(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{{}}}
	_, err := layananSpreading(src).SimpanSpreading(context.Background(), pelakuUji, "NB-1", d("50"), []models.LokasiSpreading{})
	if !errors.Is(err, ErrMasukanSpreading) || !strings.Contains(err.Error(), pesanObjekKosong) || len(src.tulis) != 0 {
		t.Fatalf("err = %v tulis %d", err, len(src.tulis))
	}
}

// TestTampilanSpreadingTotalPremiNR - GET: totalPremiumRnm = Σ premi NR tersimpan; jenis treaty tak dikenal tanpa pesan.
func TestTampilanSpreadingTotalPremiNR(t *testing.T) {
	k := kasusSatuCoverage("50")
	k.Lokasi[0].Items[0].Coverages[0].PremiNusantaraRe = d("4")
	k.Lokasi[0].Items[0].Coverages[0].Spreading = []models.BarisSpreading{{TreatyType: "X9", SharePercentage: d("100")}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{k}, idUang: map[string]string{"IDR": idMataUangIDR}}
	got, err := layananSpreading(src).TampilanSpreading(context.Background(), "NB-1")
	if err != nil {
		t.Fatal(err)
	}
	if t8(got.Lokasi[0].Items[0].TotalPremiumNusantaraRe) != "4" || len(got.Pesan) != 0 {
		t.Fatalf("total %s pesan %v", t8(got.Lokasi[0].Items[0].TotalPremiumNusantaraRe), got.Pesan)
	}
}

// CatatanJenisTreaty - "Q1" / "QS1" ber-NOTE QS, selain itu tanpa.
func (s *spreadingTiruan) CatatanJenisTreaty(_ context.Context, id string) (string, error) {
	if id == "Q1" || id == "QS1" {
		return "UJI QS FIRE", nil
	}
	return "UJI SURPLUS", nil
}

// TestSalinSpreadingProteksi - CalcultePersentageSpeading_Act: baris pertama wajib QS (kecuali ORS / FACOUT, NOTE dicari
// per ID, peka huruf) dan jenis treaty tidak boleh kembar; pesan verbatim, 400, tidak ada yang dihitung.
func TestSalinSpreadingProteksi(t *testing.T) {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusSatuCoverage("50")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	tpl := func(jenis ...string) []models.TemplateSpreading {
		var out []models.TemplateSpreading
		for _, j := range jenis {
			out = append(out, models.TemplateSpreading{TreatyType: j, TreatyName: "UJI", SharePercentage: d("10")})
		}
		return out
	}
	for nama, c := range map[string]struct {
		jenis []string
		mau   string
	}{
		"QS dulu":           {[]string{"Q1", "S1"}, ""},
		"ORS pertama":       {[]string{"10007", "S1"}, ""},
		"FACOUT pertama":    {[]string{" 10015 "}, ""},
		"tanpa QS":          {[]string{"S1", "Q1"}, pesanTanpaQS},
		"kembar":            {[]string{"Q1", "Q1"}, pesanJenisKembar},
		"kembar & tanpa QS": {[]string{"S1", "S1"}, pesanJenisKembar + "; " + pesanTanpaQS},
	} {
		_, err := layananSpreading(src).SalinSpreading(context.Background(), "NB-1", d("50"), tpl(c.jenis...))
		switch {
		case c.mau == "" && err != nil:
			t.Errorf("%s: err = %v", nama, err)
		case c.mau != "" && (!errors.Is(err, ErrMasukanSpreading) || !strings.HasSuffix(err.Error(), ": "+c.mau)):
			t.Errorf("%s: err = %v, mau %q", nama, err, c.mau)
		}
	}
}
