package models

import (
	"context"
	"testing"
	"time"
)

// acuanUji - Acuan tiruan: hanya bacaan yang dipakai uji; selainnya panik (antarmuka tersemat nil).
type acuanUji struct {
	Acuan
	jenis map[string]BarisReinsType
}

func (a acuanUji) JenisReasXOL(_ context.Context, nama string) ([]BarisReinsType, error) {
	if j, ada := a.jenis[nama]; ada {
		return []BarisReinsType{j}, nil
	}
	return nil, nil
}

// REINSURANCETYPE type 4 DEV (baca-saja 09-10-2026): XL 1ST LAYER = 10046, XL 2ND LAYER = 10036.
var jenisDEV = map[string]BarisReinsType{
	"XL 1ST LAYER": {ID: "10046", Note: "XL 1ST LAYER"},
	"XL 2ND LAYER": {ID: "10036", Note: "XL 2ND LAYER"},
	"XL 3RD LAYER": {ID: "10037", Note: "XL 3RD LAYER"},
}

func konteksUji() *Konteks {
	return &Konteks{Acuan: acuanUji{jenis: jenisDEV}, Pelaku: "UJI-PELAKU",
		Sekarang: time.Date(2026, 10, 9, 10, 0, 0, 0, Jakarta)}
}

// masterDEV1001789 - layer, kurs, MDP, dan spreading master 1001789 DEV (M_TREATY_IN.JSONDATA, baca-saja 09-10-2026):
// layer 1 limit 3,1 M retensi 2,4 M; layer 2 limit 6,5 M retensi 5,5 M; layer 3 limit 25 M retensi 12 M; RNM 30%;
// Share(1) 100% (spreading TRT, nama diganti UJI-) dipecah QS (R/I) 60% / QS (OR) 40%. ID master diberi awalan UJI-.
func masterDEV1001789() MasterTreaty {
	grup := []string{"PROPERTY", "ENGINEERING", "GENERAL ACCIDENT"}
	layer := func(n, limit, ded, mdp string) LimitXOL {
		return LimitXOL{Layer: n, LayerType: "layer", LayerPart: n, LayerPartType: "layer", Currency: "IDR",
			Currency2: "USD", Limit: limit, Limit2: "0", Deductible: ded, Deductible2: "0", ReinstatementPct: "100",
			IsCombineMDP: "false", NoRIPCalculation: "false", TreatyGroups: grup,
			MDPList: []NilaiMataUang{{Currency: "IDR", Value: mdp}}}
	}
	return MasterTreaty{ID: "UJI-1001789", ProportionType: ProporsiMaster, AccountingModeNonProp: ModeLoss, RNMShare: "30",
		Limits: []LimitXOL{
			layer("1", "3100000000", "2400000000", "1311300000.0000036000000000"),
			layer("2", "6500000000", "5500000000", "526500000.000000"),
			layer("3", "25000000000", "12000000000", "562499999.99999964000000000"),
		},
		CurrencyList: []KursMaster{{Currency: "IDR", Conversion: "1"}, {Currency: "USD", Conversion: "15500"}},
		Share: []ShareXOL{{SpreadingTypeIDXOL: "10241", SpreadingTypeXOL: "UJI QS TRT", SpreadingTotalPctXOL: "100",
			SpreadingListXOL: []SpreadingMaster{{ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", Pct: "60"},
				{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", Pct: "40"}}}},
	}
}

// halamanDEV3998 - masukan klaim CLMNP-3998 DEV (JSON_KLAIM, baca-saja 09-10-2026): klaim 10 M IDR, share ceding 100%,
// tanpa deductible / TPL / fee, Loss Allocation OR 100% ke XOL.
func halamanDEV3998() *Halaman {
	h := HalamanBaru()
	h.Setel(CD+"ShareCeding", "100")
	h.Setel(CD+"DeductibleType", "false")
	h.Setel(TM+"RNMShare", "30")
	h.Setel(TM+"TreatyGroup", "PROPERTY")
	h.SetelDaftar(DaftarClaimAmount, []Baris{{"Currency": "IDR", "CurrencyID": "10026", "AltValue": "1.0",
		"Value": "10000000000", "TPL": "0", "AdjusterFee": "0", "Salvage": "0", "CNPOthersFee": "0",
		"CNPTSI": "10000000000"}})
	h.SetelDaftar(DaftarLossAlloc, []Baris{{"Currency": "IDR", "CurrencyID": "10026", "TreatyName": "OR",
		"ClaimPercentage": "100", "CNPFlagXOL": "true"}})
	return h
}

func cekAngka(t *testing.T, apa, got, want string) {
	t.Helper()
	if !SamaAngka(got, want) {
		t.Errorf("%s = %q, ingin %q", apa, got, want)
	}
}

// Kasus DEV CLMNP-3998: hasil Pega tersimpan (SpreadingRisk, ListTotalEstimation, SpreadingClaim, SpreadingBreakQS di
// JSON_KLAIM) = hasil port CountClaimTNP_Act + CountLossAllocation_act.
func TestHitungKlaimSamaDenganKasusDEV3998(t *testing.T) {
	h := halamanDEV3998()
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	klaim := h.AmbilDaftar(DaftarClaimAmount)[0]
	cekAngka(t, "USD", klaim["USD"], "10000000000.0")
	cekAngka(t, "ClaimAmountCedant", klaim["ClaimAmountCedant"], "10000000000")
	cekAngka(t, "ClaimAmountAdjust", klaim["ClaimAmountAdjust"], "10000000000")
	cekAngka(t, "PctProrateClaim", klaim["PctProrateClaim"], "100.00000")
	cekAngka(t, "Loss Allocation ClaimAmountAdjust", h.AmbilDaftar(DaftarLossAlloc)[0]["ClaimAmountAdjust"], "10000000000")

	xol := h.AmbilDaftar(DaftarXOL)
	if len(xol) != 3 {
		t.Fatalf("SpreadingRisk %d baris, ingin 3 (UR + 2 layer; layer 3 tidak tercapai): %v", len(xol), xol)
	}
	ingin := []struct{ nama, tipe, total, spread, limit, mdp, idr string }{
		{"UR", "UR", "2400000000", "0", "", "", ""},
		{"XL 1ST LAYER", "10046", "3100000000", "930000000", "3100000000", "1311300000.0000036", "3100000000"},
		{"XL 2ND LAYER", "10036", "4500000000", "1350000000", "6500000000", "526500000", "4500000000"},
	}
	for i, w := range ingin {
		x := xol[i]
		if x["TreatyName"] != w.nama || x["TreatyType"] != w.tipe {
			t.Errorf("baris %d = %s/%s, ingin %s/%s", i+1, x["TreatyName"], x["TreatyType"], w.nama, w.tipe)
		}
		cekAngka(t, w.nama+" TotalClaim", x["TotalClaim"], w.total)
		cekAngka(t, w.nama+" ClaimEstimation", x["ClaimEstimation"], w.total)
		cekAngka(t, w.nama+" ClaimSpreaded", x["ClaimSpreaded"], w.spread)
		if w.limit != "" {
			cekAngka(t, w.nama+" CNPLimit", x["CNPLimit"], w.limit)
			cekAngka(t, w.nama+" CNPMDP", x["CNPMDP"], w.mdp)
			cekAngka(t, w.nama+" ClaimAmountIDR", x["ClaimAmountIDR"], w.idr)
			cekAngka(t, w.nama+" ClaimPercentage", x["ClaimPercentage"], "30")
		}
	}
	if err := SusunSummaryXOL(h); err != nil {
		t.Fatal(err)
	}
	s := h.AmbilDaftar(DaftarSummaryXOL)
	if len(s) != 1 {
		t.Fatalf("Summary XOL %d baris, ingin 1", len(s))
	}
	cekAngka(t, "Summary IDR", s[0]["IDR"], "10000000000")
	cekAngka(t, "Summary Value", s[0]["Value"], "2280000000")
	cekAngka(t, "Summary ClaimAmount", s[0]["ClaimAmount"], "7600000000")
	cekAngka(t, "Summary CNPTotalClaim", s[0]["CNPTotalClaim"], "10000000000")

	spr := h.AmbilDaftar(DaftarSpreading)
	if len(spr) != 1 || spr[0]["TreatyType"] != "10241" {
		t.Fatalf("SpreadingClaim = %v", spr)
	}
	cekAngka(t, "Spreading ClaimSpreaded", spr[0]["ClaimSpreaded"], "2280000000")
	cekAngka(t, "Spreading ClaimAmountIDR", spr[0]["ClaimAmountIDR"], "7600000000")
	qs := h.AmbilDaftar(DaftarBreakQS)
	if len(qs) != 2 {
		t.Fatalf("SpreadingBreakQS = %v", qs)
	}
	for i, w := range []struct{ tipe, spread, idr string }{{"10004", "1368000000", "4560000000"}, {"10028", "912000000", "3040000000"}} {
		if qs[i]["TreatyType"] != w.tipe {
			t.Errorf("Break QS %d TreatyType = %s, ingin %s", i+1, qs[i]["TreatyType"], w.tipe)
		}
		cekAngka(t, "Break QS ClaimSpreaded", qs[i]["ClaimSpreaded"], w.spread)
		cekAngka(t, "Break QS ClaimAmountIDR", qs[i]["ClaimAmountIDR"], w.idr)
	}
}

// Rumus XML (CountLossAllocation_act langkah 17.2.14): bila klaim habis teralokasi di layer 1, layer berikut tidak
// dibuat.
func TestHitungXOLBerhentiSaatKlaimHabis(t *testing.T) {
	h := halamanDEV3998()
	h.AmbilDaftar(DaftarClaimAmount)[0]["Value"] = "4000000000"
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	xol := h.AmbilDaftar(DaftarXOL)
	if len(xol) != 2 {
		t.Fatalf("SpreadingRisk %d baris, ingin 2 (UR + layer 1): %v", len(xol), xol)
	}
	cekAngka(t, "UR", xol[0]["TotalClaim"], "2400000000")
	cekAngka(t, "layer 1", xol[1]["TotalClaim"], "1600000000")
	cekAngka(t, "layer 1 RNM", xol[1]["ClaimSpreaded"], "480000000")
}

// Loss Allocation tanpa To XOL tidak masuk waterfall (langkah 13.1).
func TestHitungXOLTanpaToXOLKosong(t *testing.T) {
	h := halamanDEV3998()
	h.AmbilDaftar(DaftarLossAlloc)[0]["CNPFlagXOL"] = "false"
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	if n := len(h.AmbilDaftar(DaftarXOL)); n != 0 {
		t.Fatalf("SpreadingRisk %d baris, ingin 0", n)
	}
}

// Layer di luar TreatyGroup kasus tidak dipakai (langkah 11-12).
func TestHitungXOLLayerHanyaTreatyGroupKasus(t *testing.T) {
	h := halamanDEV3998()
	h.Setel(TM+"TreatyGroup", "MARINE")
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	if n := len(h.AmbilDaftar(DaftarXOL)); n != 0 {
		t.Fatalf("SpreadingRisk %d baris, ingin 0 (tidak ada layer bergrup MARINE)", n)
	}
}

func TestNamaLayerXOL(t *testing.T) {
	for _, c := range []struct{ layer, tipe, ingin string }{
		{"1", "layer", "XL 1ST LAYER"},
		{"2", "layer", "XL 2ND LAYER"},
		{"3", "layer", "XL 3RD LAYER"},
		{"4", "layer", "XL 4TH LAYER"},
		{"12", "sublayer", "XL 1ST SUB LAYER"},
	} {
		if got := NamaLayerXOL(c.layer, c.tipe); got != c.ingin {
			t.Errorf("NamaLayerXOL(%q,%q) = %q, ingin %q", c.layer, c.tipe, got, c.ingin)
		}
	}
}

// AdjClaimAmount_Act: Claim Amount layer 1 diubah -> RNM ikut, Summary / Spreading disusun ulang.
func TestAdjClaimAmountLayerPertama(t *testing.T) {
	h := halamanDEV3998()
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	h.AmbilDaftar(DaftarXOL)[1]["TotalClaim"] = "3000000000"
	if err := AdjClaimAmount(h, 2); err != nil {
		t.Fatal(err)
	}
	xol := h.AmbilDaftar(DaftarXOL)
	cekAngka(t, "layer 1 ClaimSpreaded", xol[1]["ClaimSpreaded"], "900000000")
	// layer terakhir menyerap selisih 100 jt terhadap Loss Allocation 10 M (langkah 6-7)
	cekAngka(t, "layer 2 AdjClaimValue", xol[2]["AdjClaimValue"], "100000000")
	cekAngka(t, "layer 2 TotalClaim", xol[2]["TotalClaim"], "4600000000")
	cekAngka(t, "layer 2 ClaimSpreaded", xol[2]["ClaimSpreaded"], "1380000000")
}

// SetActualPremium_ACT: Waiting For Actual Premium menimpa Claim RNM layer terakhir dengan konstanta 50.000.
func TestSetActualPremiumMenimpaLayerTerakhir(t *testing.T) {
	h := halamanDEV3998()
	if err := HitungKlaim(konteksUji(), h, masterDEV1001789()); err != nil {
		t.Fatal(err)
	}
	h.Setel("FlagActualPremium", "true")
	if err := SetActualPremium(h); err != nil {
		t.Fatal(err)
	}
	xol := h.AmbilDaftar(DaftarXOL)
	cekAngka(t, "layer terakhir ClaimSpreaded", xol[len(xol)-1]["ClaimSpreaded"], PremiAktual)
}
