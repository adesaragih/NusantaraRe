package services_test

// Uji rumus tab Limits Non-Prop — `hitung_limit_np.go`.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func nmu(cur, v string) services.NilaiMataUang {
	return services.NilaiMataUang{Currency: cur, Value: v}
}

// ⭐ `TotalEgnpi`: EGNPI dicocokkan lewat Treaty Group LAYER, dijumlah per
// mata uang, urut kemunculan; grup lain tidak ikut.
//
// ⛔ BUKAN total seluruh EGNPI. Ditanyakan 8 Oktober 2026 (layer PROPERTY
// tampil 29 M, bukan total 55 M): pemilik proses memilih ikut Pega — untuk
// 55 M, tambahkan Treaty Group lain ke layer itu.
func TestTotalEgnpiPerGrupDanMataUang(t *testing.T) {
	layers := []services.LayerNP{{TreatyGroupList: []services.GrupLayerNP{{TreatyGroup: "FIRE"}, {TreatyGroup: "MARINE"}}}}
	egnpi := []services.EgnpiNP{
		{TreatyGroup: "FIRE", Currency: "IDR", Amount: "100"},
		{TreatyGroup: "ENGINEERING", Currency: "IDR", Amount: "999"},
		{TreatyGroup: "MARINE", Currency: "USD", Amount: "5"},
		{TreatyGroup: "FIRE", Currency: "IDR", Amount: "50"},
	}
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPEgnpi, Layers: layers, EGNPI: egnpi})
	got := h.Layers[0].EgnpiTotalList
	if len(got) != 2 || got[0] != nmu("IDR", "150") || got[1] != nmu("USD", "5") {
		t.Fatalf("EgnpiTotalList %v", got)
	}
	if len(layers[0].EgnpiTotalList) != 0 {
		t.Error("masukan berubah")
	}
}

// ⭐ `SetReinstatementPct`: n baris, 100/100, Amount = Limit, Additional = MDP.
func TestSetReinstatementPct(t *testing.T) {
	l := services.LayerNP{ReinstatementValue: "2", Limit: "1000", Limit2: "70", ReinstatementNote: "asamount",
		MDPList: []services.NilaiMataUang{nmu("IDR", "40"), nmu("USD", "3")}}
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPReinstatement, Layers: []services.LayerNP{l}})
	r := h.Layers[0].ReinstatementList
	if len(r) != 2 {
		t.Fatalf("%d baris", len(r))
	}
	mau := services.BarisReinstatement{ReinstatementValue: "2", ReinstatementPct: "100", ReinstatementNote: "asamount",
		AdditionalAmount1: "40", AdditionalAmount2: "3", AdditionalPct: "100",
		ReinstatementAmount1: "1000", ReinstatementAmount2: "70", ID: "1"}
	if r[1] != mau {
		t.Errorf("baris 2 %+v", r[1])
	}
}

// Ketiga DT grid Reinstatement — alamat layer = layer baris itu.
func TestHitungBarisReinstatement(t *testing.T) {
	l := services.LayerNP{Limit: "1000", Limit2: "0", MDPList: []services.NilaiMataUang{nmu("IDR", "200")},
		ReinstatementList: []services.BarisReinstatement{{AdditionalPct: "50", ReinstatementPct: "25", ReinstatementAmount1: "300"}}}
	jalan := func(aksi string) services.BarisReinstatement {
		return services.HitungLimitNP(services.MasukanLimitNP{Aksi: aksi, Layers: []services.LayerNP{l}}).Layers[0].ReinstatementList[0]
	}
	if b := jalan(services.AksiNPReinstJumlah); b.ReinstatementAmount1 != "500" || b.ReinstatementAmount2 != "0" {
		t.Errorf("CalculateReinstatement %+v", b)
	}
	if b := jalan(services.AksiNPReinstPersen); b.ReinstatementPct != "30" {
		t.Errorf("CalculateReinstatementPct %+v", b)
	}
	if b := jalan(services.AksiNPReinstTambahan); b.AdditionalAmount1 != "50" || b.AdditionalAmount2 != "0" {
		t.Errorf("ReCalculateReinstatement %+v", b)
	}
}

// ⭐ adj → PE → ROL; mdp → MDP/MDPMin → ROL.
func TestDetailCalculationAdjMDPDanROL(t *testing.T) {
	kurs := []services.KursNP{{Currency: "IDR", Conversion: "1"}, {Currency: "USD", Conversion: "15000"}}
	l := services.LayerNP{Limit: "1000000", Limit2: "100", CurrencyRelation: "OR", AdjRate: "10", MDPPct: "80", MDPMinPct: "50",
		EgnpiTotalList: []services.NilaiMataUang{nmu("IDR", "500000")}}
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPAdj, Layers: []services.LayerNP{l}, Kurs: kurs})
	got := h.Layers[0]
	if len(got.PremiumEarnedList) != 1 || got.PremiumEarnedList[0] != nmu("IDR", "50000") {
		t.Fatalf("PE %v", got.PremiumEarnedList)
	}
	// ROL = 50.000 / 1.000.000 × 100 = 5
	if got.ROLPct != "5" {
		t.Errorf("ROL %q", got.ROLPct)
	}
	h = services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPMDP, Layers: h.Layers, Kurs: kurs})
	got = h.Layers[0]
	if got.MDPList[0] != nmu("IDR", "40000") || got.MDPMinList[0] != nmu("IDR", "25000") {
		t.Errorf("MDP %v Min %v", got.MDPList, got.MDPMinList)
	}
	if len(h.Pesan) != 0 {
		t.Errorf("pesan %v", h.Pesan)
	}
}

// ⚠️ Ditiru apa adanya: relasi AND menjumlahkan limit dua kali; limit nol
// memberi pesan dan ROL tidak berubah.
func TestROLRelasiANDDanLimitKosong(t *testing.T) {
	kurs := []services.KursNP{{Currency: "IDR", Conversion: "1"}}
	l := services.LayerNP{Limit: "1000", CurrencyRelation: "AND", ROLPct: "7",
		EgnpiTotalList:    []services.NilaiMataUang{nmu("IDR", "1")},
		PremiumEarnedList: []services.NilaiMataUang{nmu("IDR", "100")}}
	if pesan := services.DetailCalculationROL(&l, kurs); pesan != nil || l.ROLPct != "5" {
		t.Errorf("AND: ROL %q pesan %v — mau 100/2000×100 = 5", l.ROLPct, pesan)
	}
	k := services.LayerNP{ROLPct: "7"}
	if pesan := services.DetailCalculationROL(&k, kurs); len(pesan) != 1 || pesan[0] != services.PesanLimitKosong || k.ROLPct != "7" {
		t.Errorf("limit kosong: %v %q", pesan, k.ROLPct)
	}
}

func TestMDPPersenKurangDariSatuMemberiPesan(t *testing.T) {
	l := services.LayerNP{Limit: "10", MDPPct: "0.5", PremiumEarnedList: []services.NilaiMataUang{nmu("IDR", "100")}}
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPMDP, Layers: []services.LayerNP{l}})
	if len(h.Pesan) == 0 || h.Pesan[0] != services.PesanMDPKosong {
		t.Errorf("pesan %v", h.Pesan)
	}
	if h.Layers[0].MDPList[0].Value != "0.5" {
		t.Errorf("MDP tetap dihitung: %v", h.Layers[0].MDPList)
	}
}

// ⛔ Kontrak 1001855 (8 Oktober 2026): angka tersimpan berbentuk ketikan
// Indonesia. Dulu terbaca NOL — Total ROL tetap 0 sesudah Update Total, dan
// ROL % berhenti di "limit kosong".
func TestAngkaKetikanIndonesiaTerbaca(t *testing.T) {
	layers := []services.LayerNP{{ROLPct: "14,379"}, {ROLPct: "6,756"}, {ROLPct: "1,659"}}
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPTotal, Layers: layers})
	if h.TotalLimitsROL != "22.794" {
		t.Errorf("Total ROL %q, mau 14,379 + 6,756 + 1,659 = 22.794", h.TotalLimitsROL)
	}
	kurs := []services.KursNP{}
	l := services.LayerNP{Limit: "16.000.000.000,00", Limit2: "0", CurrencyRelation: "OR", AdjRate: "4,183",
		EgnpiTotalList: []services.NilaiMataUang{nmu("IDR", "55000000000")}}
	h = services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPAdj, Layers: []services.LayerNP{l}, Kurs: kurs})
	// PE = 55 M × 4,183% ; ROL = 2.300.650.000 / 16 M → 0,14379063 × 100
	if pe := h.Layers[0].PremiumEarnedList; len(pe) != 1 || pe[0] != nmu("IDR", "2300650000") {
		t.Errorf("PE %v", pe)
	}
	if h.Layers[0].ROLPct != "14.379063" || len(h.Pesan) != 0 {
		t.Errorf("ROL %q pesan %v", h.Layers[0].ROLPct, h.Pesan)
	}
	// Bentuk kabel TIDAK ditafsir ulang: "1.659" tetap 1,659, bukan 1.659.
	k := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPTotal, Layers: []services.LayerNP{{ROLPct: "1.659"}}})
	if k.TotalLimitsROL != "1.659" {
		t.Errorf("kabel ditafsir ulang: %q", k.TotalLimitsROL)
	}
}

// ⭐ Update Total + Summary of Limit.
func TestTotalDanRingkasanLimit(t *testing.T) {
	layers := []services.LayerNP{
		{LayerType: "layer", Layer: "1", LayerPartType: "layer", LayerPart: "1", Currency: "IDR", Currency2: "USD",
			Limit: "100", Limit2: "10", Deductible: "5", ROLPct: "2",
			PremiumEarnedList: []services.NilaiMataUang{nmu("IDR", "7")}, MDPList: []services.NilaiMataUang{nmu("IDR", "6")}},
		{LayerType: "layer", Layer: "1", LayerPartType: "layer", LayerPart: "1", Currency: "IDR",
			Limit: "50", Deductible: "1", ROLPct: "3",
			PremiumEarnedList: []services.NilaiMataUang{nmu("IDR", "3")}, MDPList: []services.NilaiMataUang{nmu("IDR", "4"), nmu("USD", "1")}},
	}
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: services.AksiNPTotal, Layers: layers})
	if io := h.Total[services.TotalLimitIOO]; len(io) != 2 || io[0] != nmu("IDR", "150") || io[1] != nmu("USD", "10") {
		t.Errorf("Total 100%% Limit %v", io)
	}
	if pe := h.Total[services.TotalLimitPE]; len(pe) != 2 || pe[0] != nmu("IDR", "10") || pe[1] != nmu("USD", "0") {
		t.Errorf("Total PE %v", pe)
	}
	if h.TotalLimitsROL != "5" {
		t.Errorf("Total ROL %q", h.TotalLimitsROL)
	}
	if len(h.LimitSummaryList) != 1 {
		t.Fatalf("ringkasan %v", h.LimitSummaryList)
	}
	r := h.LimitSummaryList[0]
	if r.Note != "layer1 of layer1" || r.Limit != "150" || r.MDP != "10" || r.MDP2 != "1" || r.Deductible != "6" {
		t.Errorf("ringkasan %+v", r)
	}
}

// Kontrak Non-Prop yang dibuka: Reinstatement dibangun pada pohon layar.
func TestSiapkanReinstatementPohon(t *testing.T) {
	pohon := []map[string]any{{"ReinstatementValue": "1", "Limit": "9", "Reinstatement_List": []map[string]any{}}}
	services.SiapkanReinstatementPohon(pohon)
	r, ok := pohon[0]["Reinstatement_List"].([]map[string]any)
	if !ok || len(r) != 1 || r[0]["ReinstatementAmount1"] != "9" || r[0]["ID"] != "1" {
		t.Fatalf("%#v", pohon[0]["Reinstatement_List"])
	}
}
