package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
)

// cargoProduksiUji - satu cargo berkolom Marine lengkap, satu coverage dua baris spreading.
func cargoProduksiUji() models.ObjekCargo {
	c := cargoSubmitUji()
	c.PackingID, c.PackingNote, c.GoodNote, c.TradingNote, c.ConveyanceNote = "P1", "UJI PACK", "UJI GOOD", "UJI TRADE", "UJI CONV"
	c.FromRute, c.ToRute, c.SailDate = "UJI DARI", "UJI KE", "20261015T060000.000 GMT"
	c.Coverages[0].Rate = d("0.125")
	c.Coverages[0].Spreading = []models.BarisSpreading{
		{TreatyType: treatyORS, SharePercentage: d("60"), TSISpreaded: d("600"), PremiumSpreaded: d("2.4"), ClaimEstimation: d("7")},
		{TreatyType: "QS", SharePercentage: d("40"), TSISpreaded: d("400"), PremiumSpreaded: d("1.6")}}
	return c
}

// TestTanggalSailFacIn - K74-9b: "dd/MM/yyyy hh:mm:ss" WIB 12-jam apa adanya (13:00 WIB -> 01:00:00); kosong / rusak = NULL.
func TestTanggalSailFacIn(t *testing.T) {
	for masuk, mau := range map[string]string{
		"20261015T060000.000 GMT": "15/10/2026 01:00:00", "20261015T050000.000 GMT": "15/10/2026 12:00:00",
		"20261014T170000.000 GMT": "15/10/2026 12:00:00", "": "", "15/10/2026": "",
	} {
		if got := tanggalSailFacIn(masuk); got != mau {
			t.Errorf("%q: %q, mau %q", masuk, got, mau)
		}
	}
}

// TestBarisFacInOfferCargo - langkah 16: kolom Marine per cargo × coverage × spreading; alamat / okupasi / TOP_RISK kosong.
func TestBarisFacInOfferCargo(t *testing.T) {
	dasar := models.BarisFacInOffer{IDPega: "UJI-NB-1", BusinessCode: "B1", GroupType: "2"}
	c2 := cargoProduksiUji()
	c2.PackingID, c2.Ship.ID, c2.SailDate = "P2", "", ""
	b, err := barisFacInOfferCargo(dasar, []models.ObjekCargo{cargoProduksiUji(), c2})
	if err != nil || len(b) != 4 {
		t.Fatalf("%v %d", err, len(b))
	}
	x := b[1]
	if x.IDPega != "UJI-NB-1" || x.BusinessCode != "B1" || x.ObjectNo != "P1" || x.ObjectName != "UJI PACK" || x.ObjectItem != "UJI GOOD" ||
		x.TradingName != "UJI TRADE" || x.ConveyanceName != "UJI CONV" || x.IDShip != "S1" || x.FromRute != "UJI DARI" ||
		x.ToRute != "UJI KE" || x.SailDate != "15/10/2026 01:00:00" || x.Coverage != "ICC A" || x.CoverageNote != "UJI COV" ||
		x.CurrID != "IDR" || x.Rate != "0.125" || x.JnReas != "QS" || teksPolis(x.PercentShareSpreaded) != "40" ||
		teksPolis(x.TSI) != "400" || teksPolis(x.Premi) != "1.6" {
		t.Errorf("%+v", x)
	}
	if x.Province != "" || x.RoadName != "" || x.RiskLocation != "" || x.OccupationCode != "" || x.AccumulationCode != "" ||
		x.TopRisk != "" || x.NoSerial != "" {
		t.Errorf("kolom non-Marine terisi: %+v", x)
	}
	if b[2].ObjectNo != "P2" || b[2].IDShip != "" || b[2].SailDate != "" || b[0].JnReas != treatyORS {
		t.Errorf("cargo 2: %+v", b[2])
	}
}

// TestIDKapalFacIn - IDSHIP VARCHAR2(20): 20 karakter lolos, 21 -> ErrIDKapalFacIn bernomor cargo (tidak dipotong).
func TestIDKapalFacIn(t *testing.T) {
	c := cargoProduksiUji()
	c.Ship.ID = strings.Repeat("K", 20)
	if b, err := barisFacInOfferCargo(models.BarisFacInOffer{}, []models.ObjekCargo{c}); err != nil || b[0].IDShip != c.Ship.ID {
		t.Errorf("20: %v", err)
	}
	c.Ship.ID = strings.Repeat("é", 11) // 11 rune, 22 bita
	panjang := c
	panjang.Ship.ID = strings.Repeat("K", 21)
	for _, x := range []models.ObjekCargo{c, panjang} {
		_, err := barisFacInOfferCargo(models.BarisFacInOffer{}, []models.ObjekCargo{cargoProduksiUji(), x})
		if !errors.Is(err, ErrIDKapalFacIn) || !strings.Contains(err.Error(), "(cargo 2)") {
			t.Errorf("%q: %v", x.Ship.ID, err)
		}
		if _, err := barisFacInProductionCargo(dasarProduksi{}, []models.ObjekCargo{x}); !errors.Is(err, ErrIDKapalFacIn) {
			t.Errorf("produksi %q: %v", x.Ship.ID, err)
		}
	}
}

// TestBarisFacInProductionCargo - langkah 22: kolom Marine + RICOMM / BROKERAGE / LOL (ClaimEstimation kosong 0); MOID case;
// OBJECT_NAME_ID "NULL"; LAYERNO "NULL"; IDX_OBJECT_ITEM / TOP_RISK / TSI100 / DEDUCTION2 tidak diisi.
func TestBarisFacInProductionCargo(t *testing.T) {
	k := models.KonteksSubmit{RICommision: d("10"), PercentShare: d("12.5"), MOID: " MO-UJI "}
	dp := dasarProduksi{id: "UJI-NB-1", kodeBisnis: "B1", groupType: "2", k: k, pctBrokerage: d("5")}
	b, err := barisFacInProductionCargo(dp, []models.ObjekCargo{cargoProduksiUji()})
	if err != nil || len(b) != 2 {
		t.Fatalf("%v %d", err, len(b))
	}
	x, y := b[0], b[1]
	if x.ObjectNo != "P1" || x.ObjectName != "UJI PACK" || x.ObjectItem != "UJI GOOD" || x.TradingName != "UJI TRADE" ||
		x.ConveyanceName != "UJI CONV" || x.IDShip != "S1" || x.FromRute != "UJI DARI" || x.ToRute != "UJI KE" ||
		x.SailDate != "15/10/2026 01:00:00" || x.Coverage != "ICC A" || x.CurrID != "IDR" || teksPolis(x.Rate) != "0.125" ||
		x.JnReas != treatyORS || teksPolis(x.TSI) != "600" || teksPolis(x.Premi) != "2.4" || teksPolis(x.PercentShareSpreaded) != "60" {
		t.Errorf("%+v", x)
	}
	if teksPolis(x.LoL) != "7" || teksPolis(y.LoL) != "0" || teksPolis(x.RIComm) != "0.24" || teksPolis(y.BrokerageFee) != "0.08" ||
		teksPolis(x.PercentRIComm) != "10" || teksPolis(x.PctBrokerageFee) != "5" {
		t.Errorf("hitungan: LoL %v/%v RIComm %v Brokerage %v", teksPolis(x.LoL), teksPolis(y.LoL), teksPolis(x.RIComm),
			teksPolis(y.BrokerageFee))
	}
	if x.MOIDKasus != "MO-UJI" || x.ObjectNameID != "NULL" || x.LayerNo != layerNoNull || x.IdxObjectItem != "" || x.TopRisk != "" ||
		x.TSI100 != nil || x.Province != "" || x.OccupationCode != "" || x.IDPega != idPega("UJI-NB-1") || x.BusinessCode != "B1" {
		t.Errorf("kolom tetap: %+v", x)
	}
}

// TestProduksiCargoJalur - gerbang 409 tiket 74 dicabut: Binding dan selesai Policy case MARINE CARGO menulis baris cargo,
// M_RISK_LOSS_PROFILE = Aneka (K74-11); case FIRE tetap cabang FIRE; ID kapal kepanjangan -> ErrIDKapalFacIn.
func TestProduksiCargoJalur(t *testing.T) {
	ctx := context.Background()
	cargo := []models.ObjekCargo{cargoProduksiUji()}
	svc := layananSubmitCargo(submitUji(), cargo)
	k := submitUji().konteks["UJI-NB-1"]
	tb, err := svc.tulisanBinding(ctx, "UJI-NB-1", k)
	if err != nil || len(tb.Offer) != 2 || tb.Offer[0].IDShip != "S1" || tb.Offer[0].ObjectNo != "P1" {
		t.Fatalf("Binding: %v %+v", err, tb)
	}
	// FIRE (006): pohon cargo tidak dibaca - fixture tanpa spreading FIRE berhenti di cabang FIRE, bukan cabang cargo.
	svc.tableOfLimit = tolTiruan{bisnis: satuBisnis("006", "24")}
	if _, err := svc.tulisanBinding(ctx, "UJI-NB-1", k); err == nil || !strings.Contains(err.Error(), "FACINOFFER tidak dapat disusun") {
		t.Errorf("FIRE Binding: %v", err)
	}
	pol := layananPolicy(submitUji(), posisiTeamLeader, models.KasusSpreading{PercentShare: d("12.5"),
		StartDateTime: "20261019T170000.000 GMT"}, 0).
		DenganObjekCargo(cargoTiruan{ada: map[string][]models.ObjekCargo{"UJI-NB-1": cargo}, mu: mataUangCargoUji()})
	pol.tableOfLimit = tolTiruan{oldID: []string{oldIDUji}, bisnis: satuBisnis("004", "24")}
	p, err := pol.tulisanPolicy(ctx, "UJI-NB-1", k, inti.Pelaku{AkunID: "UJI-USER"})
	if err != nil || len(p.Produksi) != 2 || p.Produksi[1].JnReas != "QS" || p.Produksi[0].SailDate != "15/10/2026 01:00:00" {
		t.Fatalf("Policy: %v %+v", err, p)
	}
	if r := p.RisikoLoss; r.TSINusantaraRe != nil || teksPolis(r.PremiNusantaraRe) != "0" || teksPolis(r.PremiumOR) != "0" ||
		teksPolis(r.PremiumQS) != "0" || teksPolis(r.PremiumFacout) != "0" {
		t.Errorf("risiko %+v", r)
	}
	cargo[0].Ship.ID = strings.Repeat("K", 21)
	if _, err := pol.tulisanPolicy(ctx, "UJI-NB-1", k, inti.Pelaku{AkunID: "UJI-USER"}); !errors.Is(err, ErrIDKapalFacIn) {
		t.Errorf("Policy IDSHIP: %v", err)
	}
	if _, err := svc.tulisanBinding(ctx, "UJI-NB-1", k); errors.Is(err, ErrIDKapalFacIn) {
		t.Errorf("FIRE tak terkena IDSHIP: %v", err)
	}
	svc.tableOfLimit = tolTiruan{bisnis: satuBisnis("004", "24")}
	if _, err := svc.tulisanBinding(ctx, "UJI-NB-1", k); !errors.Is(err, ErrIDKapalFacIn) {
		t.Errorf("Binding IDSHIP: %v", err)
	}
}
