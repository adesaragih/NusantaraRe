package services

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// cargoTiruan - PenyimpanObjekCargo tanpa basis data.
type cargoTiruan struct {
	ada   map[string][]models.ObjekCargo
	share *apd.Decimal
	mu    map[string]models.MataUangCsv
	tulis *[]models.ObjekCargo
	// shareTulis - % Share RNM terakhir yang ditulis UbahPercentShare (nil = tidak direkam).
	shareTulis **apd.Decimal
}

func (c cargoTiruan) BacaObjekCargo(_ context.Context, id string) ([]models.ObjekCargo, error) {
	x, ok := c.ada[id]
	if !ok {
		return nil, repository.ErrKasusTidakAda
	}
	return salinCargo(x), nil
}

func (c cargoTiruan) GantiObjekCargo(_ context.Context, _ *db.Tx, id string, cargo []models.ObjekCargo) error {
	if _, ok := c.ada[id]; !ok {
		return repository.ErrKasusTidakAda
	}
	// ID baru seperti NEXTVAL: cargo tersimpan ber-ID "N<urutan>".
	simpan := salinCargo(cargo)
	for i := range simpan {
		simpan[i].ID = "N" + string(rune('1'+i))
	}
	c.ada[id] = simpan
	if c.tulis != nil {
		*c.tulis = cargo
	}
	return nil
}

func (c cargoTiruan) PercentShareKasus(context.Context, string) (*apd.Decimal, error) {
	return c.share, nil
}

func (c cargoTiruan) UbahPercentShare(_ context.Context, _ *db.Tx, _ string, ps *apd.Decimal) error {
	if c.shareTulis != nil {
		*c.shareTulis = ps
	}
	return nil
}

func (c cargoTiruan) MataUangCargo(_ context.Context, nama string) (models.MataUangCsv, bool, error) {
	m, ok := c.mu[nama]
	return m, ok, nil
}

type masterCargoTiruan map[string]string // "jenis/id" -> deskripsi; "MaxAge/id" -> warranty

func (m masterCargoTiruan) Deskripsi(_ context.Context, jenis, id string) (string, bool, error) {
	v, ok := m[jenis+"/"+id]
	return v, ok, nil
}

func (m masterCargoTiruan) WarrantyMaxAge(_ context.Context, id string) (string, bool, error) {
	v, ok := m["MaxAge/"+id]
	return v, ok, nil
}

// IDMaster - deskripsi -> ID: kunci "jenis/id" bernilai deskripsi dicari terbalik (PERSIS, peka huruf).
func (m masterCargoTiruan) IDMaster(_ context.Context, jenis, deskripsi string) (string, bool, error) {
	for k, v := range m {
		if v == deskripsi && strings.HasPrefix(k, jenis+"/") {
			return strings.TrimPrefix(k, jenis+"/"), true, nil
		}
	}
	return "", false, nil
}

func (m masterCargoTiruan) CariAgen(_ context.Context, kata string) ([]models.BarisAgenSurvei, error) {
	return []models.BarisAgenSurvei{{ID: "A1", Name: "UJI AGEN " + kata}}, nil
}

func (m masterCargoTiruan) CariMaxAge(_ context.Context, kata string) ([]models.BarisMaxAge, error) {
	return []models.BarisMaxAge{{ID: "M1", Warranty: m["MaxAge/M1"] + kata}}, nil
}

func (m masterCargoTiruan) CariKondisiMarine(_ context.Context, kata string) ([]models.BarisKondisiMarine, error) {
	return []models.BarisKondisiMarine{{ID: "K1", ConditionName: "ICC A" + kata}}, nil
}

// Daftar - pilihan jenis itu urut deskripsi.
func (m masterCargoTiruan) Daftar(_ context.Context, jenis string) ([]models.BarisMasterCargo, error) {
	var d []models.BarisMasterCargo
	for k, v := range m {
		if strings.HasPrefix(k, jenis+"/") {
			d = append(d, models.BarisMasterCargo{ID: strings.TrimPrefix(k, jenis+"/"), Description: v})
		}
	}
	slices.SortFunc(d, func(a, b models.BarisMasterCargo) int { return strings.Compare(a.Description, b.Description) })
	return d, nil
}

// Agen - "Agen/<id>" bernilai nama (nilai "BELUM" = master agen belum terpasang); alamat / kota / negara tetap "UJI".
func (m masterCargoTiruan) Agen(_ context.Context, id string) (models.AgenSurvei, bool, error) {
	v, ok := m["Agen/"+id]
	if m["Agen/*"] == "BELUM" {
		return models.AgenSurvei{}, false, repository.ErrMasterBelumAda
	}
	if !ok {
		return models.AgenSurvei{}, false, nil
	}
	return models.AgenSurvei{AgentID: id, FullName: v, Address: "UJI ALAMAT", City: "UJI KOTA", Country: "UJI NEGARA"}, true, nil
}

// layananCargo - case "K" MARINE CARGO (GroupPanel 004), periode 01-01-2026 .. 02-07-2026 (182 hari), PS 25.
func layananCargo(gp string, lama []models.ObjekCargo, tulis *[]models.ObjekCargo) (*Service, cargoTiruan) {
	k := &models.Kasus{CaseID: "K", Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS", GroupBusinessID: "3"},
		General: models.General{StartDateTime: "20260101T050000.000 GMT", EndDateTime: "20260702T050000.000 GMT"}}
	c := cargoTiruan{ada: map[string][]models.ObjekCargo{"K": lama}, share: d("25"), tulis: tulis,
		mu: map[string]models.MataUangCsv{"IDR": {ID: "10026", OldID: "1", Currency: "IDR"}}}
	svc := Baru(nil).DenganObjekCargo(c).DenganTransaksi(tanpaTx).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"K": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: satuBisnis(gp, "24")})
	return svc, c
}

func cargoUji() models.ObjekCargo {
	return models.ObjekCargo{TradingID: "T1", GoodID: "G1", ConveyanceID: "C1", PackingID: "P1", FromRute: "A", ToRute: "B",
		SailDate: "15-03-2026", InvoiceDate: "01-03-2026", LC: models.LCCargo{LC: true, EndDate: "31-12-2026"},
		Ship: models.KapalCargo{ID: "S1", NmShip: "UJI", YMake1: "2001", GRT: "1,5", NRT: "2,25", DWT: "3"}}
}

func masterUji() masterCargoTiruan {
	return masterCargoTiruan{"Trading/T1": "DOMESTIK", "Goods/G1": "BARANG", "Conveyance/C1": "KAPAL LAUT", "Packing/P1": "PETI",
		"MaxAge/M1": "WARRANTY 25 TAHUN"}
}

// TestBacaObjekCargo - GET: tanggal Pega -> DD-MM-YYYY, cargo tanpa coverage -> satu coverage disintesis (ProRate
// AppendCoverageMarineCargo_ACT, tidak disimpan), Currency ID / OldID, total gross per mata uang (×100).
func TestBacaObjekCargo(t *testing.T) {
	lama := []models.ObjekCargo{
		{ID: "11", SailDate: "20260315T050000.000 GMT", Coverages: []models.CoverageCargo{}},
		{ID: "12", Coverages: []models.CoverageCargo{{Currency: "IDR", TSI: d("1000"), Premium: d("15")},
			{Currency: "IDR", TSI: d("3000"), Premium: d("5")}, {Currency: "USD", TSI: d("0"), Premium: d("1")}}},
	}
	svc, c := layananCargo("004", lama, nil)
	v, err := svc.BacaObjekCargo(context.Background(), "K")
	if err != nil {
		t.Fatal(err)
	}
	if v[0].Objek.SailDate != "15-03-2026" || len(v[0].Objek.Coverages) != 1 || t8(v[0].Objek.Coverages[0].ProRatePercent) != "49.863014" {
		t.Errorf("cargo 1: %+v", v[0].Objek)
	}
	if len(c.ada["K"][0].Coverages) != 0 {
		t.Error("coverage sintesis tersimpan")
	}
	tot := v[1].TotalGross
	if len(tot) != 2 || tot[0].Currency != "IDR" || t8(tot[0].TSI) != "4000" || t8(tot[0].Premium) != "20" || t8(tot[0].Rate) != "0.5" ||
		tot[1].Currency != "USD" || t8(tot[1].Rate) != "0" {
		t.Errorf("total: %+v", tot)
	}
	if v[1].Objek.Coverages[0].CurrencyID != "10026" || v[1].Objek.Coverages[0].CurrencyOldID != "1" || v[1].Objek.Coverages[2].CurrencyID != "" {
		t.Errorf("mata uang: %+v", v[1].Objek.Coverages)
	}
	if _, err := svc.BacaObjekCargo(context.Background(), "X"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("tanpa case: %v", err)
	}
}

// TestGantiObjekCargo - PUT: medan server (catatan master, MaxAgeDetail, AGE, GRT koma -> titik, tanggal Pega), coverage server
// (Currency ID, Discount, TSILiability, RNM PS 25, 8 desimal, ProRate), wajib isi, gerbang Marine Cargo, master belum terpasang.
func TestGantiObjekCargo(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	var tulis []models.ObjekCargo
	svc, _ := layananCargo("004", []models.ObjekCargo{}, &tulis)
	svc.DenganMasterCargo(masterUji())
	c := cargoUji()
	c.MaxAge = "M1"
	c.TradingNote = "KIRIMAN KLIEN" // server menimpa
	c.Coverages = []models.CoverageCargo{{Coverage: "ICC A", Currency: "IDR", TSI: d("500000000000000"), Rate: d("0.1"),
		Premium: d("500000000000"), DiscountPercentage: d("10"), LimitOfLiability: d("1.123456789")}}
	masuk := []models.ObjekCargo{c}
	v, err := svc.GantiObjekCargo(ctx, akun, "K", masuk)
	if err != nil {
		t.Fatal(err)
	}
	w := tulis[0]
	if w.TradingNote != "DOMESTIK" || w.GoodNote != "BARANG" || w.ConveyanceName != "KAPAL LAUT" || w.ConveyanceNote != "KAPAL LAUT" ||
		w.PackingNote != "PETI" || w.MaxAgeDetail != "WARRANTY 25 TAHUN" {
		t.Errorf("catatan: %+v", w)
	}
	if w.SailDate != "20260315T050000.000 GMT" || w.LC.EndDate != "20261231T050000.000 GMT" || w.Ship.Age != "25" ||
		w.Ship.GRT != "1.5" || w.Ship.NRT != "2.25" {
		t.Errorf("objek: sail %q lc %q age %q grt %q nrt %q", w.SailDate, w.LC.EndDate, w.Ship.Age, w.Ship.GRT, w.Ship.NRT)
	}
	cv := w.Coverages[0]
	if cv.CurrencyID != "10026" || t8(cv.Discount) != "50000000000" || t8(cv.TSILiability) != "500000000000000" ||
		t8(cv.TSINusantaraRe) != "125000000000000" || t8(cv.PremiNusantaraRe) != "112500000000" ||
		utils.FormatDecimal(cv.LimitOfLiability) != "1.12345679" || t8(cv.ProRatePercent) != "49.863014" {
		t.Errorf("coverage: %+v", cv)
	}
	if masuk[0].SailDate != "15-03-2026" || masuk[0].TradingNote != "KIRIMAN KLIEN" {
		t.Error("masukan pemanggil diubah")
	}
	if v[0].Objek.SailDate != "15-03-2026" || len(v[0].TotalGross) != 1 {
		t.Errorf("baca ulang: %+v", v[0])
	}
	// Wajib isi + tanggal.
	salah := cargoUji()
	salah.TradingID, salah.SailDate = "", "31-02-2026"
	if _, err := svc.GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{salah}); !errors.Is(err, ErrMasukanObjek) ||
		!strings.Contains(err.Error(), "tradingId wajib diisi") || !strings.Contains(err.Error(), "sailDate bukan tanggal") {
		t.Errorf("wajib / tanggal: %v", err)
	}
	// ID master tidak ada -> 400.
	salah = cargoUji()
	salah.GoodID = "GX"
	if _, err := svc.GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{salah}); !errors.Is(err, ErrMasukanObjek) {
		t.Errorf("master tidak ada: %v", err)
	}
	// Bukan Marine Cargo -> 409.
	fire, _ := layananCargo("006", []models.ObjekCargo{}, nil)
	if _, err := fire.DenganMasterCargo(masterUji()).GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{cargoUji()}); !errors.Is(err, ErrBukanMarineCargo) {
		t.Errorf("FIRE: %v", err)
	}
	if _, err := svc.GantiObjekCargo(ctx, inti.Pelaku{}, "K", nil); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

// TestGantiObjekCargoPertahankan - cargo ber-ID lama: coverage nil -> coverage lama (dengan spreading) dipertahankan; coverage
// dikirim -> spreading + ProRate lama per posisi; master belum terpasang: ID sama -> catatan lama, ID baru -> 503.
func TestGantiObjekCargoPertahankan(t *testing.T) {
	ctx, akun := context.Background(), inti.Pelaku{AkunID: "UJI-USER"}
	lamaCov := models.CoverageCargo{Coverage: "ICC A", TSI: d("100"), ProRatePercent: d("77"),
		Spreading: []models.BarisSpreading{{TreatyType: "10001", SharePercentage: d("100")}}}
	lama := cargoUji()
	lama.ID, lama.TradingNote, lama.GoodNote, lama.ConveyanceName, lama.ConveyanceNote, lama.PackingNote = "11", "D", "B", "KL", "KN", "P"
	lama.SailDate, lama.FlagDelete, lama.Ship.Age = "20260315T050000.000 GMT", "1", "99"
	lama.Coverages = []models.CoverageCargo{lamaCov}
	var tulis []models.ObjekCargo
	svc, _ := layananCargo("004", []models.ObjekCargo{lama}, &tulis)
	tanpaCov := cargoUji()
	tanpaCov.ID = "11"
	if _, err := svc.GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{tanpaCov}); err != nil {
		t.Fatal(err)
	}
	w := tulis[0]
	if w.TradingNote != "D" || w.ConveyanceName != "KL" || w.ConveyanceNote != "KN" || len(w.Coverages) != 1 || w.FlagDelete != "1" ||
		w.Ship.Age != "25" ||
		w.Coverages[0].Coverage != "ICC A" || len(w.Coverages[0].Spreading) != 1 {
		t.Errorf("pertahankan: %+v", w)
	}
	svc2, _ := layananCargo("004", []models.ObjekCargo{lama}, &tulis)
	kirim := cargoUji()
	kirim.ID = "11"
	kirim.Coverages = []models.CoverageCargo{{Coverage: "ICC B", TSI: d("200")}, {Coverage: "BARU"}}
	if _, err := svc2.GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{kirim}); err != nil {
		t.Fatal(err)
	}
	c0, c1 := tulis[0].Coverages[0], tulis[0].Coverages[1]
	if c0.Coverage != "ICC B" || len(c0.Spreading) != 1 || t8(c0.ProRatePercent) != "77" || len(c1.Spreading) != 0 ||
		t8(c1.ProRatePercent) != "49.863014" {
		t.Errorf("coverage dikirim: %+v / %+v", c0, c1)
	}
	tanpaTahun := cargoUji()
	tanpaTahun.ID, tanpaTahun.Ship.YMake1, tanpaTahun.Ship.Age = "11", "", "1"
	svc3, _ := layananCargo("004", []models.ObjekCargo{lama}, &tulis)
	if _, err := svc3.GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{tanpaTahun}); err != nil || tulis[0].Ship.Age != "99" {
		t.Errorf("Y_MAKE1 kosong: age %q %v", tulis[0].Ship.Age, err)
	}
	baru := cargoUji()
	baru.ID, baru.TradingID = "11", "T2"
	if _, err := svc2.GantiObjekCargo(ctx, akun, "K", []models.ObjekCargo{baru}); !errors.Is(err, ErrMasterCargoBelumAda) {
		t.Errorf("master belum ada: %v", err)
	}
}

// TestHitungCoverageCargo - CountCoverageMarine (premium.Calculate, 4 desimal), CountTSIandRateMarineCargo_Act, Currency,
// RNM PS 25.
func TestHitungCoverageCargo(t *testing.T) {
	svc, _ := layananCargo("004", []models.ObjekCargo{}, nil)
	h, err := svc.HitungCoverageCargo(context.Background(), "K", []BarisHitungCargo{
		{models.CoverageCargo{TSI: d("500000000000000"), Rate: d("0.1"), DiscountPercentage: d("10")}, PemicuCargoTSI},
		{models.CoverageCargo{TSI: d("3000"), Premium: d("1000")}, PemicuCargoPremium},
		{models.CoverageCargo{Currency: "IDR", TSI: d("1000"), Rate: d("0.33333")}, PemicuCargoCurrency},
		{models.CoverageCargo{Premium: d("10")}, PemicuCargoPremium},
	})
	if err != nil {
		t.Fatal(err)
	}
	if t8(h[0].Premium) != "500000000000" || t8(h[0].Discount) != "50000000000" || t8(h[0].TSILiability) != "500000000000000" ||
		t8(h[0].TSINusantaraRe) != "125000000000000" || t8(h[0].PremiNusantaraRe) != "112500000000" {
		t.Errorf("tsi: %+v", h[0])
	}
	if t8(h[1].Rate) != "33.3333" || t8(h[1].Discount) != "0" {
		t.Errorf("premium: %+v", h[1])
	}
	if h[2].CurrencyID != "10026" || h[2].CurrencyOldID != "1" || t8(h[2].Premium) != "3.3333" {
		t.Errorf("currency: %+v", h[2])
	}
	if h[3].Rate != nil {
		t.Errorf("TSI kosong: rate %v", h[3].Rate)
	}
	if _, err := svc.HitungCoverageCargo(context.Background(), "K", []BarisHitungCargo{{Pemicu: "x"}}); !errors.Is(err, ErrMasukanObjek) {
		t.Errorf("pemicu: %v", err)
	}
	if _, err := svc.HitungCoverageCargo(context.Background(), "K", []BarisHitungCargo{{models.CoverageCargo{Currency: "XYZ"}, PemicuCargoCurrency}}); !errors.Is(err, ErrMasukanObjek) {
		t.Errorf("currency tak dikenal: %v", err)
	}
}

// TestUmurKapalDanProRate - CountAgeShip (Y_MAKE1 kosong / tak terurai -> tidak diubah) dan batas 100 ProRate.
func TestUmurKapalDanProRate(t *testing.T) {
	for _, k := range []struct {
		sail, ymake, mau string
		ok               bool
	}{{"20260315T050000.000 GMT", "2001", "25", true}, {"20260315T050000.000 GMT", "", "", false},
		{"20260315T050000.000 GMT", "TUA", "", false}, {"", "2001", "", false}} {
		if v, ok := umurKapal(k.sail, k.ymake); v != k.mau || ok != k.ok {
			t.Errorf("%q %q: %q %v", k.sail, k.ymake, v, ok)
		}
	}
	svc, _ := layananCargo("004", nil, nil)
	svc.kasus = kasusTiruan{ada: map[string]*models.Kasus{"K": {CaseID: "K", General: models.General{StartDateTime: "20250101T050000.000 GMT",
		EndDateTime: "20260102T050000.000 GMT"}}}}
	if pr := svc.proRateCargo(context.Background(), "K"); t8(pr) != "100" {
		t.Errorf("ProRate 366 hari: %s", t8(pr))
	}
}
