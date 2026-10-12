package services_test

// Rincian spreading MANUAL (`SetSpreadingXOL` [5.3.8]–[5.9]) — permintaan
// pemakai 11 Oktober 2026: Reins Type selain ORS harus memperlihatkan anak
// susunannya seperti Pega. Pembanding: dokumen Pega 1002304 (RNM Share 19,
// baris "2026 QS 212M TRT" Pct 14, RNM Limit layer 7.837.500.000) — angka
// harapan disalin PERSIS dari `M_TREATY_IN.JSONDATA` kontrak itu.

import (
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

var sumber1002304 = sumberTiruan{
	induk: []models.SusunanSpreading{
		{ReinsTypeID: "10260", ReinsTypeName: "2026 QS 212M TRT", TreatyYearID: "1000700"},
		{ReinsTypeID: "10007", ReinsTypeName: "ORS", TreatyYearID: "1000700"},
	},
	anak: map[string][]models.SusunanSpreading{
		"10260": {
			{ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", ParentReinsTypeID: "10260", Pct: "60"},
			{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", ParentReinsTypeID: "10260", Pct: "40"},
		},
	},
}

// samaNilai - teks desimal sama secara angka.
func samaNilai(t *testing.T, label, got, harap string) {
	t.Helper()
	a, _, err1 := apd.NewFromString(got)
	b, _, err2 := apd.NewFromString(harap)
	if err1 != nil || err2 != nil || a.Cmp(b) != 0 {
		t.Errorf("%s = %q, harap %s", label, got, harap)
	}
}

func shareSpreadingManual(baris ...models.BarisSpreadingNP) models.ShareNP {
	s := shareAwal("19")
	s.Share = []models.BarisShareNP{{
		RNMShare:         "19",
		RnmLimitList:     []models.NilaiMataUang{{Currency: "IDR", Value: "7837500000"}},
		GrossPremiumList: []models.NilaiMataUang{{Currency: "IDR", Value: "358769020"}},
		SpreadingListXOL: baris,
	}}
	return s
}

func TestSpreadingManualMemecahAnakSusunanSepertiPega(t *testing.T) {
	s := shareSpreadingManual(models.BarisSpreadingNP{ReinsTypeID: "10260", Pct: "14"})
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingPct, Share: s, Indeks: 0}, sumber1002304)
	sp := h.Share.Share[0].SpreadingListXOL[0]
	if sp.ReinsTypeName != "2026 QS 212M TRT" || sp.Currency != "IDR" {
		t.Fatalf("baris spreading %q / %q", sp.ReinsTypeName, sp.Currency)
	}
	samaNilai(t, "Value", sp.Value, "5775000000.000000000024750000000000")
	if len(sp.BreakDownSprdListXOL) != 2 {
		t.Fatalf("rincian %+v", sp.BreakDownSprdListXOL)
	}
	ri, or := sp.BreakDownSprdListXOL[0], sp.BreakDownSprdListXOL[1]
	if ri.ReinsTypeName != "QS (R/I)" || ri.SharePct != "60" || or.ReinsTypeName != "QS (OR)" || or.SharePct != "40" {
		t.Errorf("rincian %+v", sp.BreakDownSprdListXOL)
	}
	samaNilai(t, "Amount QS (R/I)", ri.Amount, "3465000000.00000000001485000000")
	samaNilai(t, "Amount QS (OR)", or.Amount, "2310000000.00000000000990000000")
	samaNilai(t, "Spreading OR Limit", nilaiDi(t, sp.RNMSpreadedListXOL, "IDR"), "2310000000.00000000000990000000000000000000000000000000")
	samaNilai(t, "Spreading R/I Limit", nilaiDi(t, sp.RNMSpreadedListRIXOL, "IDR"), "3465000000.00000000001485000000000000000000000000000000")
	samaNilai(t, "Spreading OR MDP", nilaiDi(t, sp.RNMSpreadedListGrossXOL, "IDR"), "105742448.000000000000453181920000000000000000000000000000000")
	samaNilai(t, "Spreading R/I MDP", nilaiDi(t, sp.RNMSpreadedListGrossRIXOL, "IDR"), "158613672.000000000000679772880000000000000000000000000000000")
	samaNilai(t, "Total OR Limit", nilaiDi(t, h.Share.Total["TotalSpreadedRnmProp"], "IDR"), "2310000000.0000000000099")
	samaNilai(t, "Total R/I Limit", nilaiDi(t, h.Share.Total["TotalSpreadedRnmRIProp"], "IDR"), "3465000000.00000000001485")
}

// ORS tanpa anak di master tetap 100 % OR (keputusan 10 Oktober 2026).
func TestSpreadingManualORSTanpaAnakSeratusPersenOR(t *testing.T) {
	s := shareSpreadingManual(models.BarisSpreadingNP{ReinsTypeID: "10007", Pct: "19"})
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingPct, Share: s, Indeks: 0}, sumber1002304)
	sp := h.Share.Share[0].SpreadingListXOL[0]
	if len(sp.BreakDownSprdListXOL) != 1 || sp.BreakDownSprdListXOL[0].ReinsTypeName != "ORS" || sp.BreakDownSprdListXOL[0].SharePct != "100" {
		t.Fatalf("rincian ORS %+v", sp.BreakDownSprdListXOL)
	}
	samaNilai(t, "Amount ORS", sp.BreakDownSprdListXOL[0].Amount, "7837500000")
	samaNilai(t, "Spreading OR Limit", nilaiDi(t, sp.RNMSpreadedListXOL, "IDR"), "7837500000")
	if len(sp.RNMSpreadedListRIXOL) != 0 {
		t.Errorf("R/I ORS mestinya kosong: %+v", sp.RNMSpreadedListRIXOL)
	}
	samaNilai(t, "Total OR Limit", nilaiDi(t, h.Share.Total["TotalSpreadedRnmProp"], "IDR"), "7837500000")
}

// Saat kontrak dibuka (`SiapkanShareNP`) rinciannya ikut dihitung.
func TestSiapkanShareNPMemecahSpreadingManual(t *testing.T) {
	s := shareSpreadingManual(models.BarisSpreadingNP{ReinsTypeID: "10260", ReinsTypeName: "2026 QS 212M TRT", Pct: "14"})
	s.Share[0].RnmLimitList = nil
	layer := []services.LayerNP{{LayerType: "Layer", Layer: "1", LayerPartType: "Layer", LayerPart: "1",
		Currency: "IDR", Limit: "41250000000", MDPList: []services.NilaiMataUang{nmu("IDR", "1888258000")}}}
	services.SiapkanShareNP(&s, layer, "1009999", sumber1002304, "20260101")
	sp := s.Share[0].SpreadingListXOL[0]
	samaNilai(t, "Value", sp.Value, "5775000000.000000000024750000000000")
	if len(sp.BreakDownSprdListXOL) != 2 {
		t.Fatalf("rincian %+v", sp.BreakDownSprdListXOL)
	}
	samaNilai(t, "Amount QS (OR)", sp.BreakDownSprdListXOL[1].Amount, "2310000000.00000000000990000000")
}
