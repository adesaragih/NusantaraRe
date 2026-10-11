package services

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// cargoSubmitUji - satu cargo sah: kapal ber-ID + nama, satu coverage IDR berketerangan, spreading ORS 100 % premi = premi RNM.
func cargoSubmitUji() models.ObjekCargo {
	return models.ObjekCargo{Ship: models.KapalCargo{ID: "S1", NmShip: "UJI"}, Coverages: []models.CoverageCargo{{Coverage: "ICC A",
		CoverageNote: "UJI COV", Currency: "IDR", TSI: d("1000"), TSILiability: d("1000"), PremiNusantaraRe: d("4"),
		Spreading: []models.BarisSpreading{{TreatyType: treatyORS, SharePercentage: d("100"), PremiumSpreaded: d("4")}}}}}
}

// mataUangCargoUji - master CURRENCY tiruan: IDR sah, XXX tanpa OldID.
func mataUangCargoUji() map[string]models.MataUangCsv {
	return map[string]models.MataUangCsv{"IDR": {ID: "10026", OldID: "1", Currency: "IDR"}, "XXX": {ID: "10999", Currency: "XXX"}}
}

// layananSubmitCargo - layananSubmit (proteksi tiket 57) untuk case MARINE CARGO (GroupPanel 004) tanpa Scoring Risk; objek FIRE
// = satu lokasi kosong (atau `fire`).
func layananSubmitCargo(st *submitTiruan, cargo []models.ObjekCargo, fire ...models.ObjekFire) *Service {
	if fire == nil {
		fire = []models.ObjekFire{{ObjectNo: "1"}}
	}
	k := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{InsuredID: "UJI-ORG", ClassOfBusiness: "UJI KELAS",
		GroupBusinessID: "3"}, General: models.General{MarketingID: "UJI-MO", ReffNumber: "UJI/SLIP/1",
		StartDateTime: "20261001T170000.000 GMT", EndDateTime: "20271001T170000.000 GMT"}}
	return layananSubmit(st, fire...).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"UJI-NB-1": k}}).
		DenganTableOfLimit(tolTiruan{bisnis: satuBisnis("004", "24")}).
		DenganObjekCargo(cargoTiruan{ada: map[string][]models.ObjekCargo{"UJI-NB-1": cargo}, mu: mataUangCargoUji()})
}

// TestProteksiCoverageCargo - ProtectCoverage_Act langkah 4 MARINE: coverage / keterangan kosong, mata uang kosong / tak dikenal /
// tanpa OldID; nomor cargo 1..n; cargo tanpa coverage tidak diproteksi (K74-12); semua ber-flag ProtectMarine.CARI2.
func TestProteksiCoverageCargo(t *testing.T) {
	c1, c2 := cargoSubmitUji(), cargoSubmitUji()
	c1.Coverages = append(c1.Coverages, models.CoverageCargo{Coverage: "ICC B", Currency: "XXX"}, models.CoverageCargo{CoverageNote: "N",
		Currency: "ZZZ"})
	c2.Coverages[0].Currency = ""
	cargo := []models.ObjekCargo{c1, {}, c2}
	var p pengumpulPesan
	svc := Baru(nil).DenganObjekCargo(cargoTiruan{mu: mataUangCargoUji()})
	if err := svc.proteksiCoverageCargo(context.Background(), cargo, &p); err != nil {
		t.Fatal(err)
	}
	mau := []string{" Coverage 1 can't be empty!", "Currency on object 1 coverage 2 must be choose!",
		"Currency on object 1 coverage 3 must be choose!", "Currency on object 3 coverage 1 must be choose!"}
	if !reflect.DeepEqual(p.pesan, mau) {
		t.Errorf("%q", p.pesan)
	}
	for i, f := range p.flag {
		if !reflect.DeepEqual(f, []flagProteksi{flagMarine2}) {
			t.Errorf("flag %d: %v", i, f)
		}
	}
}

// TestProteksiKapalCargo - ProtectShipData_Act (DUA spasi verbatim, flag CARI1) + "Multi Location Must Have Toprisk" hanya bila
// > 1 cargo tanpa top risk.
func TestProteksiKapalCargo(t *testing.T) {
	c1, c2 := cargoSubmitUji(), cargoSubmitUji()
	c2.Ship.NmShip = " "
	var p pengumpulPesan
	proteksiKapalCargo([]models.ObjekCargo{c1, c2}, &p)
	if !reflect.DeepEqual(p.pesan, []string{" Ship data CargoList 2 can't  be empty!", "Multi Location Must Have Toprisk"}) ||
		!reflect.DeepEqual(p.flag[0], []flagProteksi{flagMarine1}) {
		t.Errorf("%q %v", p.pesan, p.flag)
	}
	c2.IsTopRisk = true
	p = pengumpulPesan{}
	proteksiKapalCargo([]models.ObjekCargo{c1, c2}, &p)
	if len(p.pesan) != 1 {
		t.Errorf("top risk: %q", p.pesan)
	}
	c1.Ship.ID = ""
	p = pengumpulPesan{}
	proteksiKapalCargo([]models.ObjekCargo{c1}, &p)
	if !reflect.DeepEqual(p.pesan, []string{" Ship data CargoList 1 can't  be empty!"}) {
		t.Errorf("satu cargo: %q", p.pesan)
	}
}

// TestPesanSpreadingCargo - K74-13: Σshare TIDAK dibulatkan (99.999 ditahan); Σ premi vs premi RNM dibandingkan 2 desimal;
// coverage tanpa spreading -> pesan 100 %.
func TestPesanSpreadingCargo(t *testing.T) {
	c := cargoSubmitUji()
	c.Coverages[0].Spreading = []models.BarisSpreading{{SharePercentage: d("60"), PremiumSpreaded: d("2.001")},
		{SharePercentage: d("40"), PremiumSpreaded: d("2.002")}}
	if p := pesanSpreadingCargo([]models.ObjekCargo{c}); len(p) != 0 {
		t.Errorf("sah 2 desimal: %q", p)
	}
	c.Coverages[0].Spreading[1] = models.BarisSpreading{SharePercentage: d("39.999"), PremiumSpreaded: d("2.1")}
	c2 := cargoSubmitUji()
	c2.Coverages[0].Spreading = nil
	mau := []string{"% Total spreading on cargo 1 coverage 1 must be 100%!",
		"Total premium spreaded on cargo 1 coverage 1 is not equals premium RNM!",
		"% Total spreading on cargo 2 coverage 1 must be 100%!",
		"Total premium spreaded on cargo 2 coverage 1 is not equals premium RNM!"}
	if p := pesanSpreadingCargo([]models.ObjekCargo{c, c2}); !reflect.DeepEqual(p, mau) {
		t.Errorf("%q", p)
	}
}

// TestProteksiAcceptCargo - case MARINE CARGO sah lolos tanpa pesan FIRE (okupasi / Scoring / mata uang TSI); cargo rusak mendapat
// pesan cargo. Case FIRE tidak membaca pohon cargo.
func TestProteksiAcceptCargo(t *testing.T) {
	ctx := context.Background()
	pesan, err := layananSubmitCargo(submitUji(), []models.ObjekCargo{cargoSubmitUji()}).ProteksiAccept(ctx, "UJI-NB-1", nil)
	if err != nil || len(pesan) != 0 {
		t.Fatalf("cargo sah: %v %q", err, pesan)
	}
	rusak := cargoSubmitUji()
	rusak.Ship, rusak.Coverages[0].Coverage = models.KapalCargo{}, ""
	pesan, err = layananSubmitCargo(submitUji(), []models.ObjekCargo{rusak, cargoSubmitUji()}).ProteksiAccept(ctx, "UJI-NB-1", nil)
	gabung := strings.Join(pesan, "|")
	for _, m := range []string{" Coverage 1 can't be empty!", " Ship data CargoList 1 can't  be empty!", "Multi Location Must Have Toprisk"} {
		if err != nil || !strings.Contains(gabung, m) {
			t.Errorf("%s: %v %q", m, err, pesan)
		}
	}
	if strings.Contains(gabung, "Occupation") || strings.Contains(gabung, "Scoring") {
		t.Errorf("pesan FIRE ikut: %q", pesan)
	}
	// Langkah 21 ProtectCurrencyTSI (IsFire) tidak berjalan walau baris T_LOCATIONLIST ber-item dengan mata uang tak dikenal.
	fire := lokasiSah()
	fire.Items[0].Currency = "XYZ"
	pesan, err = layananSubmitCargo(submitUji(), []models.ObjekCargo{cargoSubmitUji()}, fire).ProteksiAccept(ctx, "UJI-NB-1", nil)
	if err != nil || strings.Contains(strings.Join(pesan, "|"), "Bermasalah") {
		t.Errorf("mata uang FIRE: %v %q", err, pesan)
	}
	svc := layananSubmitCargo(submitUji(), []models.ObjekCargo{rusak})
	svc.tableOfLimit = tolTiruan{bisnis: satuBisnis("006", "24")}
	if pesan, _ := svc.ProteksiAccept(ctx, "UJI-NB-1", nil); strings.Contains(strings.Join(pesan, "|"), "Ship data") {
		t.Errorf("FIRE membaca cargo: %q", pesan)
	}
}

// layananJenjangCargo - Submit tahap uw `posisi` case MARINE CARGO: spreading FIRE tanpa item (PS 50), pohon cargo, tabel limit
// PRY (limitUji) + NONPROPANDENGG team 1 yang jenjangnya berbeda dari PRY.
func layananJenjangCargo(st *submitTiruan, posisi string, cargo []models.ObjekCargo, kurs map[string]*apd.Decimal) *Service {
	sp := models.KasusSpreading{PercentShare: d("50"), StartDateTime: "20261001T170000.000 GMT", Lokasi: []models.LokasiSpreading{{ObjectNo: "1"}}}
	svc := layananSubmitUW(st, posisi, sp).DenganTableOfLimit(tolTiruan{oldID: []string{oldIDUji}, bisnis: satuBisnis("004", "24")}).
		DenganObjekCargo(cargoTiruan{ada: map[string][]models.ObjekCargo{"UJI-NB-1": cargo}, share: d("50"), mu: mataUangCargoUji()}).
		DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR, "USD": "10001"}, kurs: kurs})
	svc.limit = limitTiruan{a: append(append([]models.BarisLimitA{}, limitUji...),
		barisLimit("M_LIMIT_NONPROPANDENGG", "1", "JUW_B", "0"), barisLimit("M_LIMIT_NONPROPANDENGG", "1", "JUW_A", "1"))}
	return svc
}

// TestCarid2Cargo - K74-8 langkah 10 Marine Cargo: Σ setiap cargo × coverage .TSI × kurs coverage × PS; top risk = .TSILiability;
// FlagDelete "1" dilewati.
func TestCarid2Cargo(t *testing.T) {
	ctx := context.Background()
	kurs := map[string]*apd.Decimal{"10001": d("2")}
	hitung := func(ubah func([]models.ObjekCargo)) string {
		c := []models.ObjekCargo{cargoSubmitUji(), cargoSubmitUji(), cargoSubmitUji()}
		c[1].Coverages[0].Currency, c[1].Coverages[0].TSI = "USD", d("500")
		c[2].FlagDelete, c[2].Coverages[0].TSI = "1", d("999999")
		ubah(c)
		d, err := layananJenjangCargo(submitUji(), posisiJuniorUWB, c, kurs).bacaDasarJenjang(ctx, "UJI-NB-1")
		if err != nil {
			t.Fatal(err)
		}
		if d.pesanKurs != "" {
			return d.pesanKurs
		}
		return utils.FormatDecimal(d.carid2())
	}
	if c := hitung(func([]models.ObjekCargo) {}); c != "1000" { // 50 % × (1000 + 2 × 500)
		t.Errorf("tanpa top risk: %s", c)
	}
	if c := hitung(func(c []models.ObjekCargo) { c[1].IsTopRisk, c[1].Coverages[0].TSILiability = true, d("300") }); c != "300" {
		t.Errorf("top risk berkurs: %s", c)
	}
}

// TestJenjangCargo - JUW B Accept case MARINE CARGO memakai M_LIMIT_NONPROPANDENGG (K74-8) -> JUW_A; FACOUT di spreading cargo =
// IsFacRetro; case FIRE (GroupPanel 006) tidak membaca cargo.
func TestJenjangCargo(t *testing.T) {
	ctx := context.Background()
	c := []models.ObjekCargo{cargoSubmitUji()}
	h, err := layananJenjangCargo(submitUji(), posisiJuniorUWB, c, nil).AjukanSubmit(ctx, pelakuPosisi(posisiJuniorUWB), "UJI-NB-1",
		models.IsiSubmit{EmailType: "1", Comment: "x"})
	if err != nil || h.Position != posisiJuniorUWA {
		t.Errorf("Marine: %v %+v", err, h)
	}
	d, err := layananJenjangCargo(submitUji(), posisiJuniorUWB, c, nil).bacaDasarJenjang(ctx, "UJI-NB-1")
	if err != nil || d.cargo == nil || d.adaFacRetro() {
		t.Errorf("tanpa FACOUT: %v %+v", err, d.cargo)
	}
	c[0].Coverages[0].Spreading[0].TreatyType = treatyFACOUT
	if d, _ := layananJenjangCargo(submitUji(), posisiJuniorUWB, c, nil).bacaDasarJenjang(ctx, "UJI-NB-1"); !d.adaFacRetro() {
		t.Error("FACOUT cargo bukan IsFacRetro")
	}
	fire := layananJenjangCargo(submitUji(), posisiJuniorUWB, c, nil)
	fire.tableOfLimit = tolTiruan{oldID: []string{oldIDUji}, bisnis: satuBisnis("006", "24")}
	if d, err := fire.bacaDasarJenjang(ctx, "UJI-NB-1"); err != nil || d.cargo != nil {
		t.Errorf("FIRE: %v %+v", err, d.cargo)
	}
}

// TestJenjangLanjutCargo - ActFlow (Senior) case MARINE CARGO: When COB tangga membaca BusinessType case (MarineCargo) ->
// M_LIMIT_NONPROPANDENGG: SENIORUW 1 -> DEPHEADUNDERWRITER 2; tabel PRY akan selesai.
func TestJenjangLanjutCargo(t *testing.T) {
	svc := layananJenjangCargo(submitUji(), posisiSeniorUW, []models.ObjekCargo{cargoSubmitUji()}, nil)
	svc.limit = limitTiruan{a: append(append([]models.BarisLimitA{}, limitUji...),
		barisLimit("M_LIMIT_NONPROPANDENGG", "1", "SENIORUW", "1"), barisLimit("M_LIMIT_NONPROPANDENGG", "1", "DEPHEADUNDERWRITER", "2"))}
	h, err := svc.AjukanSubmit(context.Background(), pelakuPosisi(posisiSeniorUW), "UJI-NB-1", models.IsiSubmit{EmailType: "1", Comment: "x"})
	if err != nil || h.Position != posisiDepHeadUW {
		t.Errorf("%v %+v", err, h)
	}
}

// TestSpreadingUWCargo - tahap uw: CheckSpreadingProtectMCargoMBU ikut (case MARINE CARGO), dengan nomor cargo / coverage.
func TestSpreadingUWCargo(t *testing.T) {
	c := cargoSubmitUji()
	c.Coverages[0].Spreading = nil
	pesan, err := layananJenjangCargo(submitUji(), posisiJuniorUWB, []models.ObjekCargo{c}, nil).proteksiSpreadingUW(context.Background(), "UJI-NB-1")
	gabung := strings.Join(pesan, "|")
	if err != nil || !strings.Contains(gabung, "% Total spreading on cargo 1 coverage 1 must be 100%!") ||
		!strings.Contains(gabung, "Total premium spreaded on cargo 1 coverage 1 is not equals premium RNM!") {
		t.Errorf("%v %q", err, pesan)
	}
}

// TestIsFlagUWMarine - tiket 74 P5: MARINE CARGO = 1 kecuali polis deklarasi (PolicyType 2) = 0.
func TestIsFlagUWMarine(t *testing.T) {
	ctx := context.Background()
	for pt, mau := range map[string]bool{"0": true, "1": true, "2": false} {
		svc := layananPolicy(submitUji(), posisiTeamLeader, models.KasusSpreading{PercentShare: d("12.5"),
			StartDateTime: "20261019T170000.000 GMT"}, 0).DenganTableOfLimit(tolTiruan{oldID: []string{oldIDUji}, bisnis: satuBisnis("004", "24")})
		svc.kasus.(kasusTiruan).ada["UJI-NB-1"].General.PolicyType = pt
		if got, err := svc.isFlagUW(ctx, "UJI-NB-1"); err != nil || got != mau {
			t.Errorf("PolicyType %s: %v %v", pt, got, err)
		}
	}
}
