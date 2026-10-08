package services_test

// Uji ALUR tab Share Non-Prop — rantai Activity per tombol (Update Summary,
// spreading bernama/manual, % RNM Share baris, Update Total) dan pemuatan
// dari tabel pendaratan. Uji per-rumus ada di `hitung_share_np_test.go`.

import (
	"encoding/json"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

// sumberTiruan — kedua RD spreading.
type sumberTiruan struct {
	induk []models.SusunanSpreading
	anak  map[string][]models.SusunanSpreading
}

func (s sumberTiruan) Induk(_, _ string) []models.SusunanSpreading { return s.induk }
func (s sumberTiruan) Anak(induk, _ string) []models.SusunanSpreading {
	return s.anak[induk]
}

// Bentuk susunan nyata `PROPORTIONALARRG` (induk 10227, tahun 1000610).
var sumberQS = sumberTiruan{
	induk: []models.SusunanSpreading{
		{ReinsTypeID: "10227", ReinsTypeName: "2022 QS 145M TRT", TreatyYearID: "1000610"},
		{ReinsTypeID: "10007", ReinsTypeName: "ORS", TreatyYearID: "1000610"},
	},
	anak: map[string][]models.SusunanSpreading{
		"10227": {
			{ReinsTypeID: "10028", ReinsTypeName: "QS (OR)", ParentReinsTypeID: "10227", Pct: "40"},
			{ReinsTypeID: "10004", ReinsTypeName: "QS (R/I)", ParentReinsTypeID: "10227", Pct: "60"},
		},
		"10007": {{ReinsTypeID: "10007", ReinsTypeName: "ORS", ParentReinsTypeID: "10007", Pct: "100"}},
	},
}

func layerShare() []services.LayerNP {
	return []services.LayerNP{
		{LayerType: "Layer", Layer: "1", LayerPartType: "Layer", LayerPart: "1", Cover: "risk",
			Currency: "IDR", Limit: "1000000000", Currency2: "USD", Limit2: "0",
			TreatyGroupList: []services.GrupLayerNP{{TreatyGroup: "PROPERTY", TreatyGroupID: "10002"}},
			MDPList:         []services.NilaiMataUang{nmu("IDR", "50000000")}},
		{LayerType: "Layer", Layer: "2", LayerPartType: "Layer", LayerPart: "1",
			Currency: "IDR", Limit: "2000000000", Currency2: "USD", Limit2: "100000",
			TreatyGroupList: []services.GrupLayerNP{{TreatyGroup: "PROPERTY", TreatyGroupID: "10002"}},
			MDPList:         []services.NilaiMataUang{nmu("IDR", "30000000"), nmu("USD", "1000")}},
	}
}

func shareAwal(rnm string) models.ShareNP {
	return models.ShareNP{RNMShare: rnm, BrokeragePercent: "10", RNMShareAcrossTheBoard: "true", Total: map[string][]models.NilaiMataUang{}}
}

func nilaiDi(t *testing.T, daftar []models.NilaiMataUang, cur string) string {
	t.Helper()
	for _, v := range daftar {
		if v.Currency == cur {
			return v.Value
		}
	}
	t.Fatalf("mata uang %s tidak ada di %+v", cur, daftar)
	return ""
}

func TestNonAddItemShareSatuBarisPerLayer(t *testing.T) {
	s := shareAwal("40")
	pesan := services.NonAddItemShare(&s, layerShare())
	if len(pesan) != 0 || len(s.Share) != 2 {
		t.Fatalf("pesan %v, baris %d", pesan, len(s.Share))
	}
	b := s.Share[0]
	// 100% Limit RNM = Limit × @divide(40,100,4); Limit2 < 1 tidak ditambah.
	if len(b.RnmLimitList) != 1 || nilaiDi(t, b.RnmLimitList, "IDR") != "400000000" {
		t.Errorf("RnmLimitList %+v", b.RnmLimitList)
	}
	if nilaiDi(t, b.GrossPremiumList, "IDR") != "20000000" {
		t.Errorf("Gross %+v", b.GrossPremiumList)
	}
	// [13.5] Net := Gross DIKOMENTARI — Net kosong sampai SetBrokerage.
	if len(b.NetPremiumList) != 0 {
		t.Errorf("Net mestinya kosong: %+v", b.NetPremiumList)
	}
	if b.RNMShare != "40" || b.Cover != "risk" || b.TreatyGroupList[0].TreatyGroupID != "10002" {
		t.Errorf("medan layer %+v", b)
	}
	if nilaiDi(t, s.Share[1].RnmLimitList, "USD") != "40000" {
		t.Errorf("Limit2 layer 2 %+v", s.Share[1].RnmLimitList)
	}
}

func TestUpdateSummarySpreadingBernamaDanBrokerage(t *testing.T) {
	s := shareAwal("40")
	// S1: Spreading Type baris lama dibawa ke baris baru berindeks sama.
	s.Share = []models.BarisShareNP{{SpreadingTypeXOL: "2022 QS 145M TRT"}, {SpreadingTypeXOL: "ORS"}}
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSummary, Share: s, Layers: layerShare()}, sumberQS)
	b := h.Share.Share[0]
	if b.SpreadingTypeIDXOL != "10227" || len(b.SpreadingListXOL) != 2 || b.SpreadingTotalPctXOL != "100" {
		t.Fatalf("spreading %+v", b)
	}
	// OR = @divide(40,100,9) × 400.000.000; R/I 60%.
	if nilaiDi(t, b.RNMSpreadedListXOL, "IDR") != "160000000" || nilaiDi(t, b.RNMSpreadedListRIXOL, "IDR") != "240000000" {
		t.Errorf("OR/RI %+v %+v", b.RNMSpreadedListXOL, b.RNMSpreadedListRIXOL)
	}
	// Brokerage fee 10% dari Gross(1); Net = Gross − Deduction.
	if len(b.DeductionList) != 1 || b.DeductionList[0].Comment != "Brokerage fee" || b.DeductionList[0].Deduction != "2000000" {
		t.Errorf("deduksi %+v", b.DeductionList)
	}
	if nilaiDi(t, b.NetPremiumList, "IDR") != "18000000" {
		t.Errorf("Net %+v", b.NetPremiumList)
	}
	// Net OR = 40% (QS (OR)).
	if nilaiDi(t, b.RNMSpreadedListNetXOL, "IDR") != "7200000" {
		t.Errorf("Net OR %+v", b.RNMSpreadedListNetXOL)
	}
	// ORS: OR Limit 100% (`FetchQSfromMasterXOL` [19.2]).
	o := h.Share.Share[1]
	if nilaiDi(t, o.RNMSpreadedListXOL, "IDR") != "800000000" || nilaiDi(t, o.RNMSpreadedListRIXOL, "IDR") != "0" {
		t.Errorf("ORS %+v / %+v", o.RNMSpreadedListXOL, o.RNMSpreadedListRIXOL)
	}
	if nilaiDi(t, h.Share.Total["TotalSpreadedRnmProp"], "IDR") != "960000000" {
		t.Errorf("Total OR %+v", h.Share.Total["TotalSpreadedRnmProp"])
	}
	// Net OR baris ORS: `TreatyInSetBrokerage` hanya mengenal `QS (OR)`, dan
	// persennya TERBAWA dari baris sebelumnya — 40% × (12 jt − 1,2 jt).
	if nilaiDi(t, o.RNMSpreadedListNetXOL, "IDR") != "4320000" {
		t.Errorf("Net OR baris ORS %+v", o.RNMSpreadedListNetXOL)
	}
	if nilaiDi(t, h.Share.Total["TotalSpreadedNetPremi"], "IDR") != "11520000" {
		t.Errorf("Total OR Net %+v", h.Share.Total["TotalSpreadedNetPremi"])
	}
}

func TestSpreadingManualLaluUpdateTotal(t *testing.T) {
	s := shareAwal("40")
	s.Share = []models.BarisShareNP{{SpreadingTypeXOL: "2022 QS 145M TRT"}, {}}
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSummary, Share: s, Layers: layerShare()}, sumberQS)
	// Baris 2 tanpa Spreading Type → spreading manual: 25 + 15 = 40.
	s2 := h.Share
	s2.Share[1].SpreadingListXOL = []models.BarisSpreadingNP{{ReinsTypeID: "10227", Pct: "25"}, {ReinsTypeID: "10007", Pct: "15"}}
	h2 := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingPct, Share: s2, Indeks: 1}, sumberQS)
	if len(h2.Pesan) != 0 || len(h2.PesanBaris) != 0 {
		t.Errorf("total 40 = RNM 40, nol pesan: %v %v", h2.Pesan, h2.PesanBaris)
	}
	if h2.Share.Share[1].SpreadingListXOL[0].ReinsTypeName != "2022 QS 145M TRT" {
		t.Errorf("nama induk dari RD: %+v", h2.Share.Share[1].SpreadingListXOL[0])
	}
	h3 := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareTotal, Share: h2.Share}, sumberQS)
	tot := h3.Share.Total
	// Gross baris 2 = 30 jt × 0,4 = 12 jt, dibagi 25/40 + 15/40 → 12 jt.
	// Baris 1 (bernama) TIDAK dijumlah — hanya blok [24].
	if nilaiDi(t, tot["TotalShareGrossNP"], "IDR") != "12000000" {
		t.Errorf("Total Gross %+v", tot["TotalShareGrossNP"])
	}
	if len(tot["TotalShareRnmNP"]) != 0 || len(tot["TotalShareGrossMinNP"]) != 0 {
		t.Errorf("RNM %+v, Gross Min %+v", tot["TotalShareRnmNP"], tot["TotalShareGrossMinNP"])
	}
	if len(h3.Share.LimitShareSummaryList) != 2 || h3.Share.LimitShareSummaryList[0].Note != "Layer1 of Layer1" {
		t.Errorf("summary %+v", h3.Share.LimitShareSummaryList)
	}
}

func TestSpreadingManualTidakSamaRNMBerpesan(t *testing.T) {
	s := shareAwal("40")
	s.Share = []models.BarisShareNP{{RNMShare: "40", SpreadingListXOL: []models.BarisSpreadingNP{{Pct: "10"}}}}
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingPct, Share: s}, sumberQS)
	// Pesan medan `.SpreadingTotalPctXOL` — menempel pada barisnya, bukan di kepala tab.
	if len(h.Pesan) != 0 || len(h.PesanBaris) != 1 || h.PesanBaris[0] != (services.PesanBarisShare{Indeks: 0, Pesan: services.PesanShareSpreading}) {
		t.Errorf("pesan %v / baris %v", h.Pesan, h.PesanBaris)
	}
}

func TestSpreadingTambahHapusMenghitungTotalPct(t *testing.T) {
	s := shareAwal("40")
	s.Share = []models.BarisShareNP{{SpreadingListXOL: []models.BarisSpreadingNP{{Pct: "30"}}}}
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingTambah, Share: s}, sumberQS)
	if n := len(h.Share.Share[0].SpreadingListXOL); n != 2 || h.Share.Share[0].SpreadingListXOL[1].Pct != "0" {
		t.Fatalf("tambah: %d %+v", n, h.Share.Share[0].SpreadingListXOL)
	}
	h2 := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingHapus, Share: h.Share, Baris: 0}, sumberQS)
	if n := len(h2.Share.Share[0].SpreadingListXOL); n != 1 || h2.Share.Share[0].SpreadingTotalPctXOL != "0" {
		t.Errorf("hapus: %d total %q", n, h2.Share.Share[0].SpreadingTotalPctXOL)
	}
}

func TestSpreadingTypeDropdownMengisiAnakDanTotal(t *testing.T) {
	s := shareAwal("40")
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSummary, Share: s, Layers: layerShare()[:1]}, sumberQS)
	s2 := h.Share
	s2.Share[0].SpreadingTypeXOL = "2022 QS 145M TRT"
	h2 := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingType, Share: s2, Indeks: 0}, sumberQS)
	b := h2.Share.Share[0]
	if len(b.SpreadingListXOL) != 2 || nilaiDi(t, b.RNMSpreadedListXOL, "IDR") != "160000000" {
		t.Errorf("spreading %+v / %+v", b.SpreadingListXOL, b.RNMSpreadedListXOL)
	}
	// `IsUpdate = 0` → Net ikut dibagi; GetNilaiTotal mengisi keempat total.
	if nilaiDi(t, b.RNMSpreadedListNetXOL, "IDR") != "7200000" || nilaiDi(t, h2.Share.Total["TotalSpreadedNetPremiRI"], "IDR") != "10800000" {
		t.Errorf("Net %+v total %+v", b.RNMSpreadedListNetXOL, h2.Share.Total["TotalSpreadedNetPremiRI"])
	}
}

func TestSummaryShareIDRdanUSDBarisPenuh(t *testing.T) {
	rows := []models.BarisShareNP{
		{LayerType: "Layer", Layer: "1", LayerPartType: "Layer", LayerPart: "1",
			RnmLimitList:     []models.NilaiMataUang{nmu("IDR", "100"), nmu("USD", "7"), nmu("SGD", "9")},
			GrossPremiumList: []models.NilaiMataUang{nmu("IDR", "10")},
			NetPremiumList:   []models.NilaiMataUang{nmu("IDR", "9")},
			DeductionList:    []models.BarisDeduksiShare{{Currency: "IDR", Deduction: "1"}}},
		{LayerType: "Layer", Layer: "2"},
	}
	r := services.SummaryLimitShare(rows)
	if len(r) != 1 {
		t.Fatalf("layer tanpa 100%% Limit tidak masuk: %+v", r)
	}
	if r[0].Limit != "100" || r[0].Limit2 != "7" || r[0].MDP != "10" || r[0].NetPremi != "9" || r[0].Deductible != "1" || r[0].MDP2 != "" {
		t.Errorf("%+v", r[0])
	}
}

func TestRNMBarisMenghapusSpreadingType(t *testing.T) {
	s := shareAwal("40")
	s.Share = []models.BarisShareNP{{SpreadingTypeXOL: "2022 QS 145M TRT"}}
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSummary, Share: s, Layers: layerShare()[:1]}, sumberQS)
	s2 := h.Share
	s2.Share[0].RNMShare = "20"
	h2 := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareRNMBaris, Share: s2, Indeks: 0, Layers: layerShare()[:1]}, sumberQS)
	b := h2.Share.Share[0]
	// `FetchQSfromMasterXOL` dipanggil TANPA induk — spreading baris terhapus.
	if b.SpreadingTypeXOL != "" || len(b.SpreadingListXOL) != 0 {
		t.Errorf("spreading mestinya kosong: %+v", b)
	}
	if nilaiDi(t, b.RnmLimitList, "IDR") != "200000000" || nilaiDi(t, b.RNMSpreadedListRIXOL, "IDR") != "200000000" {
		t.Errorf("RNM baris 20%%: %+v / %+v", b.RnmLimitList, b.RNMSpreadedListRIXOL)
	}
}

func TestDeduksiPanelMenghormatiAutoCalculate(t *testing.T) {
	s := shareAwal("40")
	s.Share = []models.BarisShareNP{{
		GrossPremiumList: []models.NilaiMataUang{nmu("IDR", "1000"), nmu("USD", "50")},
		DeductionList:    []models.BarisDeduksiShare{{Comment: "X", DeductionPct: "10", DeductionPctCalculate: "false"}},
	}}
	h := services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareDeduksi, Share: s, Sts: "pct"}, sumberQS)
	b := h.Share.Share[0]
	// Bendera bukan `true` → langkah 5 (mata uang berikutnya) dilompati.
	if len(b.DeductionList) != 1 || b.DeductionList[0].Deduction != "100" {
		t.Errorf("deduksi %+v", b.DeductionList)
	}
	if nilaiDi(t, b.NetPremiumList, "IDR") != "900" || nilaiDi(t, b.NetPremiumList, "USD") != "50" {
		t.Errorf("Net %+v", b.NetPremiumList)
	}
}

func TestShareDariPendaratan(t *testing.T) {
	sp := models.SharePendaratan{
		Share: []map[string]any{
			{"Layer": "1", "RNMShare": "", "DeductionList": []map[string]any{{"Comment": "Brokerage fee", "DeductionPct": "2.5"}}},
			{"Layer": "2", "RNMShare": "35"},
		},
		ShareReins: []map[string]any{{"ReinsName": "SINGAPORE RE", "SharePct": "5"}},
	}
	s := services.ShareDariPendaratan(sp, map[string]string{"FacultativeShare": "10"})
	if s.RNMShare != "35" || s.BrokeragePercent != "2.5" || s.RNMShareAcrossTheBoard != "true" || s.RnmShareDeducted != "25" {
		t.Errorf("%+v", s)
	}
	if len(s.ShareReins) != 1 || s.ShareReins[0].ReinsName != "SINGAPORE RE" {
		t.Errorf("reins %+v", s.ShareReins)
	}
	b, _ := json.Marshal(s)
	if strings.Contains(string(b), "null") {
		t.Errorf("larik nil bocor: %s", b)
	}
}

func TestSiapkanShareNPMenurunkanNilaiDariLayer(t *testing.T) {
	sp := models.SharePendaratan{Share: []map[string]any{
		{"Layer": "1", "RNMShare": "40", "SpreadingTypeXOL": "2022 QS 145M TRT",
			"SpreadingListXOL": []map[string]any{{"ReinsTypeName": "QS (OR)", "Pct": "40"}, {"ReinsTypeName": "QS (R/I)", "Pct": "60"}},
			"GrossPremiumList": []map[string]any{{"Currency": "IDR", "Value": "20000000"}},
			"NetPremiumList":   []map[string]any{{"Currency": "IDR", "Value": "18000000"}},
			"DeductionList":    []map[string]any{{"Comment": "Brokerage fee", "Currency": "IDR", "Deduction": "2000000", "DeductionPct": "10"}}},
	}}
	s := services.ShareDariPendaratan(sp, map[string]string{})
	services.SiapkanShareNP(&s, layerShare()[:1], "1001270")
	b := s.Share[0]
	if nilaiDi(t, b.RnmLimitList, "IDR") != "400000000" || b.TreatyGroupList[0].TreatyGroup != "PROPERTY" {
		t.Errorf("turunan layer %+v", b)
	}
	if nilaiDi(t, s.Total["TotalSpreadedRnmProp"], "IDR") != "160000000" || nilaiDi(t, s.Total["TotalSpreadedNetPremi"], "IDR") != "7200000" {
		t.Errorf("total OR %+v net %+v", s.Total["TotalSpreadedRnmProp"], s.Total["TotalSpreadedNetPremi"])
	}
	if nilaiDi(t, b.DeductionTotalList, "IDR") != "2000000" || len(s.LimitShareSummaryList) != 1 {
		t.Errorf("deduksi total %+v summary %+v", b.DeductionTotalList, s.LimitShareSummaryList)
	}
}

// sumberCatat mencatat Treaty Group yang diminta RD induk.
type sumberCatat struct {
	sumberTiruan
	grup *[]string
}

func (s sumberCatat) Induk(grup, mulai string) []models.SusunanSpreading {
	*s.grup = append(*s.grup, grup)
	return s.sumberTiruan.Induk(grup, mulai)
}

// ⭐ `SetSpreadingXOL` [5.3.2] mengirim `TempSprd.TreatyGroupID` — tak pernah
// diisi untuk baris Non-Prop, jadi filter B dilewati (grup KOSONG). Spreading
// Type bernama memakai `.TreatyGroupList(1).TreatyGroupID`.
func TestGrupRDIndukMenurutJalurnya(t *testing.T) {
	var diminta []string
	src := sumberCatat{sumberTiruan: sumberQS, grup: &diminta}
	s := shareAwal("40")
	s.Share = []models.BarisShareNP{{RNMShare: "40",
		TreatyGroupList:  []models.GrupShareNP{{TreatyGroup: "PROPERTY", TreatyGroupID: "10002"}},
		SpreadingListXOL: []models.BarisSpreadingNP{{ReinsTypeID: "10227", Pct: "40"}}}}
	services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingPct, Share: s}, src)
	if len(diminta) == 0 || diminta[0] != "" {
		t.Errorf("spreading manual meminta grup %q, mau kosong", diminta)
	}
	diminta = nil
	s.Share[0].SpreadingTypeXOL = "2022 QS 145M TRT"
	services.HitungShareNP(services.MasukanShareNP{Aksi: services.AksiShareSpreadingType, Share: s}, src)
	if len(diminta) == 0 || diminta[0] != "10002" {
		t.Errorf("Spreading Type meminta grup %q, mau 10002", diminta)
	}
}

// ⭐ Urutan sumber skalar akar (keputusan pemakai 7 Oktober 2026): kolom
// `T_TREATY_REVISION` → salinan Pega di baris Share → `TREATYINDETAIL`.
func TestTerapkanAkarShareMenurutUrutanSumber(t *testing.T) {
	salinan := models.ShareNP{RNMShare: "35", BrokeragePercent: "2.5", RNMShareAcrossTheBoard: "true", FacultativeShare: "5"}
	detail := models.ShareDetailWarisan{RNMShare: "40", BrokeragePercent: "3"}

	// 1. Kolom menang atas segalanya.
	s := salinan
	services.TerapkanAkarShare(&s, map[string]string{"RNMShare": "45", "BrokeragePercent": "1", "RNMShareAcrossTheBoard": "false"}, detail)
	if s.RNMShare != "45" || s.BrokeragePercent != "1" || s.RNMShareAcrossTheBoard != "false" || s.RnmShareDeducted != "40" {
		t.Errorf("kolom: %+v", s)
	}
	// 2. Tanpa kolom → salinan baris.
	s = salinan
	services.TerapkanAkarShare(&s, map[string]string{}, detail)
	if s.RNMShare != "35" || s.BrokeragePercent != "2.5" || s.RnmShareDeducted != "30" {
		t.Errorf("salinan: %+v", s)
	}
	// 3. Tanpa kolom dan salinan → TREATYINDETAIL.
	s = models.ShareNP{RNMShareAcrossTheBoard: "true"}
	services.TerapkanAkarShare(&s, map[string]string{}, detail)
	if s.RNMShare != "40" || s.BrokeragePercent != "3" {
		t.Errorf("detail: %+v", s)
	}
	// Nol sumber → kosong, TIDAK ditebak.
	s = models.ShareNP{}
	services.TerapkanAkarShare(&s, map[string]string{}, models.ShareDetailWarisan{})
	if s.RNMShare != "" || s.BrokeragePercent != "" {
		t.Errorf("tebakan: %+v", s)
	}
}
