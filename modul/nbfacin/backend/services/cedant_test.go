package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// cedantTiruan - PenyimpanCedant dalam memori; "K404" = case tidak ada, kode "TIADA" = tidak lolos AGENT.
type cedantTiruan struct {
	ada   models.KasusCedant
	tulis []models.SimpanCedant
}

func (c *cedantTiruan) BacaCedant(_ context.Context, id string) (models.KasusCedant, error) {
	if id == "K404" {
		return models.KasusCedant{}, repository.ErrKasusTidakAda
	}
	return c.ada, nil
}

func (c *cedantTiruan) TulisCedant(_ context.Context, _ *db.Tx, id string, k models.SimpanCedant) error {
	if id == "K404" {
		return repository.ErrKasusTidakAda
	}
	for i, b := range k.Cedant {
		if b.CedingCo == "TIADA" {
			return fmt.Errorf("%w: cedant[%d]", repository.ErrCedingTidakSah, i)
		}
	}
	c.tulis = append(c.tulis, k)
	c.ada.ShareCedantType, c.ada.Cedant = k.ShareCedantType, k.Cedant
	return nil
}

func TestProteksiCedantBerlaku(t *testing.T) {
	for nama, c := range map[string]struct {
		g   gerbangCedant
		mau bool
	}{
		"kosong":           {gerbangCedant{}, false},
		"group":            {gerbangCedant{Group: true}, true},
		"on-going":         {gerbangCedant{OnGoing: true}, true},
		"edm adj":          {gerbangCedant{EdmAdjShareCedant: true}, true},
		"life menang":      {gerbangCedant{Life: true, Group: true}, false},
		"binding 4 menang": {gerbangCedant{EmailBinding4: true, OnGoing: true}, false},
	} {
		if got := proteksiCedantBerlaku(c.g); got != c.mau {
			t.Errorf("%s: %v", nama, got)
		}
	}
	if proteksiCedantBerlaku(gerbangKasus(models.KasusCedant{FlagOnGoingPolicy: "0"})) {
		t.Error("case NB baru (FLAG_ONGOING_POLICY 0): proteksi tidak berlaku")
	}
	if !proteksiCedantBerlaku(gerbangKasus(models.KasusCedant{FlagOnGoingPolicy: "2"})) {
		t.Error("FLAG_ONGOING_POLICY 2: proteksi berlaku")
	}
}

func TestPeriksaProteksiCedant(t *testing.T) {
	sah := models.BarisCedant{CedingCo: "A1", CedingCoName: "UJI A", ShareCeding: d("100")}
	for nama, c := range map[string]struct {
		baris []models.BarisCedant
		mau   []string
	}{
		"sah":          {[]models.BarisCedant{sah}, nil},
		"kosong":       {[]models.BarisCedant{}, []string{pesanCedantKosong}},
		"share nol":    {[]models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A", ShareCeding: d("0")}}, []string{pesanShareCedant}},
		"share kosong": {[]models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A"}}, []string{pesanShareCedant}},
		"share 100.5":  {[]models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A", ShareCeding: d("100.5")}}, []string{pesanShareCedant}},
		"nama kosong":  {[]models.BarisCedant{{CedingCo: "A1", ShareCeding: d("10")}}, []string{pesanCedingKosong}},
		"ceding kosong": {[]models.BarisCedant{sah, {ShareCeding: d("0")}},
			[]string{pesanShareCedant, pesanCedingKosong}},
	} {
		if got := periksaProteksiCedant(c.baris, "1", nil, nil); strings.Join(got, "|") != strings.Join(c.mau, "|") {
			t.Errorf("%s: %v", nama, got)
		}
	}
}

func TestShareOfCeding(t *testing.T) {
	isi := func(p *string) string {
		if p == nil {
			return "<nil>"
		}
		return *p
	}
	for _, c := range []struct{ tipe, ps, mau string }{
		{"0", "50", "100%"}, {"1", "12.5", "12.5%"}, {"1", "50.00000000", "50%"}, {"1", "", "%"}, {"", "50", "<nil>"},
	} {
		var ps *apd.Decimal
		if c.ps != "" {
			ps = d(c.ps)
		}
		if got := isi(shareOfCeding(c.tipe, ps)); got != c.mau {
			t.Errorf("tipe %q ps %q = %q, mau %q", c.tipe, c.ps, got, c.mau)
		}
	}
}

// kasusDuaItem - item 1 (IDR, TSI objek 1500): coverage 1000 (premi 10, diskon 2) + coverage 500 (premi 4); item 2 (mata
// uang `kedua`, TSI objek 200): coverage 200 (premi 2).
func kasusDuaItem(ps, kedua string) models.KasusSpreading {
	k := kasusSatuCoverage(ps)
	it := &k.Lokasi[0].Items[0]
	it.TSI = d("1500")
	it.Coverages = append(it.Coverages, models.CoverageSpreading{OldID: "C2", TSILiability: d("500"), Premium: d("4")})
	k.Lokasi[0].Items = append(k.Lokasi[0].Items, models.ItemSpreading{Currency: kedua, TSI: d("200"),
		Coverages: []models.CoverageSpreading{{OldID: "C3", TSILiability: d("200"), Premium: d("2")}}})
	return k
}

// TestRingkasanRNM - TSI = 50 × (1000 + 200) / 100 (coverage pertama tiap item); premi = (5 − 1) + 2 + 1; dasar per mata
// uang: IDR TSI 1500 / premi 6, USD TSI 200 / premi 1 (Share RNM: TSI × 50 / 100; Gross: premi × 100 / 50).
func TestRingkasanRNM(t *testing.T) {
	premi, jumlah, dasar := ringkasanRNM(kasusDuaItem("50", "USD"))
	if t8(premi) != "7" || len(jumlah) != 2 || len(dasar) != 2 {
		t.Fatalf("premi mentah = %s jumlah %d dasar %d", t8(premi), len(jumlah), len(dasar))
	}
	// per mata uang asli: IDR TSI 50 × 1000 / 100 = 500, premi 4 + 2; USD TSI 50 × 200 / 100 = 100, premi 1.
	if j := jumlah[0]; j.Currency != "IDR" || t8(j.TSIRNM) != "500" || t8(j.Premi) != "6" {
		t.Errorf("jumlah IDR = %+v", j)
	}
	if j := jumlah[1]; j.Currency != "USD" || t8(j.TSIRNM) != "100" || t8(j.Premi) != "1" {
		t.Errorf("jumlah USD = %+v", j)
	}
	mau := [][]string{{"IDR", "750", "6", "1500", "12"}, {"USD", "100", "1", "200", "2"}}
	for i, m := range mau {
		x := dasar[i]
		got := []string{x.Currency, t8(x.TSIRNM), t8(x.PremiumRNM), t8(x.TSIGross), t8(x.PremiumGross)}
		if strings.Join(got, "|") != strings.Join(m, "|") {
			t.Errorf("dasar %d = %v, mau %v", i, got, m)
		}
	}
	if premi, jumlah, dasar := ringkasanRNM(kasusDuaItem("", "USD")); premi != nil || len(jumlah) != 0 || len(dasar) != 0 {
		t.Fatal("tanpa % Share RNM: kosong")
	}
	if _, _, dasar := ringkasanRNM(kasusDuaItem("0", "USD")); dasar[0].PremiumGross != nil {
		t.Fatal("% Share RNM 0: premi gross kosong")
	}
}

// TestProteksiTotalPremi - ProtectShareCedant_Act 5: Σ premi CurrencyList cedant (4 desimal) wajib = TotalPremiNusaRe.
func TestProteksiTotalPremi(t *testing.T) {
	premi, _, dasar := ringkasanRNM(kasusDuaItem("50", "USD"))
	baris := []models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A", ShareCeding: d("60")},
		{CedingCo: "B2", CedingCoName: "UJI B", ShareCeding: d("40")}}
	if m := periksaProteksiCedant(baris, "1", dasar, premi); len(m) != 0 { // (60 + 40) / 100 × (6 + 1) = 7
		t.Errorf("share RNM 100 %%: %v", m)
	}
	if m := periksaProteksiCedant(baris, "0", dasar, premi); len(m) != 1 || m[0] != pesanTotalPremi { // 12 + 2 = 14
		t.Errorf("gross: %v", m)
	}
	baris[1].ShareCeding = d("39.99") // 99,99 / 100 × 7 = 6,9993 (4 desimal) ≠ 7
	if m := periksaProteksiCedant(baris, "1", dasar, premi); len(m) != 1 || m[0] != pesanTotalPremi {
		t.Errorf("share RNM 99,99 %%: %v", m)
	}
}

func layananCedant(c *cedantTiruan) *Service {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "IDR")}, idUang: map[string]string{"IDR": idMataUangIDR}}
	return Baru(nil).DenganCedant(c).DenganSpreading(src).DenganTransaksi(tanpaTx)
}

func TestTampilanDanSimpanCedant(t *testing.T) {
	c := &cedantTiruan{ada: models.KasusCedant{SobName: "UJI SOB", CedingUmum: []models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A"}}}}
	svc := layananCedant(c)
	got, err := svc.TampilanCedant(context.Background(), "NB-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.SobName != "UJI SOB" || got.Wajib || t8(got.PercentShare) != "50" || t8(got.TotalTSIRNM) != "600" || len(got.CedingUmum) != 1 {
		t.Fatalf("tampilan %+v", got)
	}
	baris := []models.BarisCedant{{CedingCo: "A1", ShareCeding: d("60")}, {ShareCeding: nil}}
	if _, err := svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "1", baris); err != nil {
		t.Fatal(err)
	}
	if len(c.tulis) != 1 || c.tulis[0].ShareOfCeding == nil || *c.tulis[0].ShareOfCeding != "50%" || len(c.tulis[0].Cedant) != 2 {
		t.Fatalf("tulis %+v", c.tulis)
	}
	if _, err := svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "", nil); err != nil || len(c.tulis) != 2 ||
		c.tulis[1].ShareOfCeding != nil {
		t.Fatalf("tipe kosong: err %v tulis %+v", err, c.tulis)
	}
	for nama, x := range map[string]struct {
		id, tipe string
		baris    []models.BarisCedant
		mau      error
		teks     string
	}{
		"tipe":      {"NB-1", "2", nil, ErrMasukanCedant, "shareCedantType"},
		"share 101": {"NB-1", "0", []models.BarisCedant{{CedingCo: "A1", ShareCeding: d("101")}}, ErrMasukanCedant, pesanShareCeding},
		"kode":      {"NB-1", "0", []models.BarisCedant{{CedingCo: strings.Repeat("x", 51)}}, ErrMasukanCedant, "cedingCo"},
		"agent":     {"NB-1", "0", []models.BarisCedant{{CedingCo: "TIADA"}}, ErrMasukanCedant, "cedant[0]"},
		"404":       {"K404", "0", nil, ErrKasusTidakAda, ""},
		"total 110": {"NB-1", "0", []models.BarisCedant{{CedingCo: "A1", ShareCeding: d("50")}, {CedingCo: "B2", ShareCeding: d("60")}},
			ErrMasukanCedant, pesanTotalShareCedant},
	} {
		_, err := svc.SimpanCedant(context.Background(), pelakuUji, x.id, x.tipe, x.baris)
		if !errors.Is(err, x.mau) || !strings.Contains(err.Error(), x.teks) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if _, err := Baru(nil).TampilanCedant(context.Background(), "NB-1"); !errors.Is(err, ErrCedantTanpaDatabase) {
		t.Errorf("tanpa basis data: %v", err)
	}
}

// TestSimpanCedantOnGoing - FLAG_ONGOING_POLICY 2: Share Cedant Type wajib dan proteksi ProtectShareCedant_Act berlaku.
func TestSimpanCedantOnGoing(t *testing.T) {
	c := &cedantTiruan{ada: models.KasusCedant{FlagOnGoingPolicy: "2"}}
	svc := layananCedant(c)
	got, err := svc.TampilanCedant(context.Background(), "NB-1")
	if err != nil || !got.Wajib {
		t.Fatalf("wajib %v err %v", got.Wajib, err)
	}
	_, err = svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "", []models.BarisCedant{})
	if !errors.Is(err, ErrMasukanCedant) ||
		!strings.HasSuffix(err.Error(), ": "+pesanTipeCedantWajib+"; "+pesanTotalPremi+"; "+pesanCedantKosong) ||
		len(c.tulis) != 0 {
		t.Fatalf("err = %v tulis %d", err, len(c.tulis))
	}
	sah := []models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A", ShareCeding: d("100")}}
	if _, err := svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "1", sah); err != nil || len(c.tulis) != 1 {
		t.Fatalf("sah: err %v tulis %d", err, len(c.tulis))
	}
}

// TestTotalShareCedantTepat100 - C-7: Σ tepat 100 (eksak, termasuk desimal) lolos; 100,00000001 ditolak.
func TestTotalShareCedantTepat100(t *testing.T) {
	c := &cedantTiruan{}
	svc := layananCedant(c)
	sah := []models.BarisCedant{{CedingCo: "A1", ShareCeding: d("33.33333333")}, {CedingCo: "B2", ShareCeding: d("66.66666667")}}
	if _, err := svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "1", sah); err != nil || len(c.tulis) != 1 {
		t.Fatalf("tepat 100: err %v tulis %d", err, len(c.tulis))
	}
	sah[1].ShareCeding = d("66.66666668")
	if _, err := svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "1", sah); !errors.Is(err, ErrMasukanCedant) ||
		!strings.HasSuffix(err.Error(), ": "+pesanTotalShareCedant) || len(c.tulis) != 1 {
		t.Fatalf("100,00000001: err %v tulis %d", err, len(c.tulis))
	}
}

// TestTotalKeIDR - C-8: Σ per mata uang × kurs (IDR = ID 10026 -> 1; selain itu kurs pada Begin date); satu mata uang tanpa
// kurs -> kedua total kosong dan mata uangnya dilaporkan; item tanpa mata uang = "(kosong)".
func TestTotalKeIDR(t *testing.T) {
	_, jumlah, _ := ringkasanRNM(kasusDuaItem("50", "USD"))
	src := &spreadingTiruan{idUang: map[string]string{"IDR": idMataUangIDR, "USD": "UJI-USD"},
		kurs: map[string]*apd.Decimal{"UJI-USD": d("15000")}}
	tsi, premi, tidak, err := totalKeIDR(context.Background(), src, "20250101T170000.000 GMT", jumlah)
	if err != nil || t8(tsi) != "1500500" || t8(premi) != "15006" || len(tidak) != 0 { // 500 + 100 × 15000; 6 + 1 × 15000
		t.Fatalf("IDR = %s / %s tidak %v err %v", t8(tsi), t8(premi), tidak, err)
	}
	delete(src.kurs, "UJI-USD")
	if tsi, premi, tidak, err := totalKeIDR(context.Background(), src, "", jumlah); err != nil || tsi != nil || premi != nil ||
		len(tidak) != 1 || tidak[0] != "USD" {
		t.Fatalf("tanpa kurs USD: %v / %v tidak %v err %v", tsi, premi, tidak, err)
	}
	_, jumlah, _ = ringkasanRNM(kasusDuaItem("50", ""))
	if _, _, tidak, _ := totalKeIDR(context.Background(), src, "", jumlah); len(tidak) != 1 || tidak[0] != mataUangKosong {
		t.Fatalf("item tanpa mata uang: %v", tidak)
	}
	if tsi, _, tidak, _ := totalKeIDR(context.Background(), src, "", nil); tsi != nil || tidak == nil || len(tidak) != 0 {
		t.Fatalf("tanpa jumlah: %v %v", tsi, tidak)
	}
}
