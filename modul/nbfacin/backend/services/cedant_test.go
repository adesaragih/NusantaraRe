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
		if got := periksaProteksiCedant(c.baris); strings.Join(got, "|") != strings.Join(c.mau, "|") {
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

// kasusDuaItem - item 1: coverage 1000 (premi 10, diskon 2) + coverage 500 (premi 4); item 2: coverage 200 (premi 2).
func kasusDuaItem(ps string) models.KasusSpreading {
	k := kasusSatuCoverage(ps)
	it := &k.Lokasi[0].Items[0]
	it.Coverages = append(it.Coverages, models.CoverageSpreading{OldID: "C2", TSILiability: d("500"), Premium: d("4")})
	k.Lokasi[0].Items = append(k.Lokasi[0].Items, models.ItemSpreading{Currency: "IDR",
		Coverages: []models.CoverageSpreading{{OldID: "C3", TSILiability: d("200"), Premium: d("2")}}})
	return k
}

// TestTotalRNM - TSI = 50 × (1000 + 200) / 100 (coverage pertama tiap item); premi = (5 − 1) + 2 + 1.
func TestTotalRNM(t *testing.T) {
	tsi, premi := totalRNM(kasusDuaItem("50"))
	if t8(tsi) != "600" || t8(premi) != "7" {
		t.Fatalf("total = %s / %s", t8(tsi), t8(premi))
	}
	if tsi, premi := totalRNM(kasusDuaItem("")); tsi != nil || premi != nil {
		t.Fatal("tanpa % Share RNM: kosong")
	}
}

func layananCedant(c *cedantTiruan) *Service {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50")}}
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
	if !errors.Is(err, ErrMasukanCedant) || !strings.HasSuffix(err.Error(), ": "+pesanTipeCedantWajib+"; "+pesanCedantKosong) ||
		len(c.tulis) != 0 {
		t.Fatalf("err = %v tulis %d", err, len(c.tulis))
	}
	sah := []models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A", ShareCeding: d("100")}}
	if _, err := svc.SimpanCedant(context.Background(), pelakuUji, "NB-1", "0", sah); err != nil || len(c.tulis) != 1 {
		t.Fatalf("sah: err %v tulis %d", err, len(c.tulis))
	}
}
