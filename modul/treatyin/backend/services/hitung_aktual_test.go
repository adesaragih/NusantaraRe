package services

import (
	"errors"
	"reflect"
	"testing"
)

func hitungAktualUji(t *testing.T, m MasukanAktual) HasilAktual {
	t.Helper()
	h, err := HitungAktual(m)
	if err != nil {
		t.Fatalf("HitungAktual: %v", err)
	}
	return h
}

// Satu layer FIRE (Limit 1000, AdjRate 10) di Treaty In; Actual GNPI FIRE IDR 500.
func akarAktualUji() HalamanPohon {
	return pohon(map[string]string{"TotalEgnpiAmount": "400", "RNMShare": "40", "BrokeragePercent": "10", "FacultativeShare": "0"},
		map[string][]map[string]any{
			"Limits": {{
				"LayerType": "Layer ", "Layer": "1", "Currency": "IDR", "Limit": "1000", "AdjRate": "10", "ROLPct": "2",
				"TreatyGroupList":   []any{map[string]any{"TreatyGroup": "FIRE"}},
				"EgnpiTotalList":    []any{map[string]any{"Currency": "IDR", "Value": "400"}},
				"PremiumEarnedList": []any{map[string]any{"Currency": "IDR", "Value": "40"}},
				"MDPList":           []any{map[string]any{"Currency": "IDR", "Value": "40"}},
			}},
			"Share": {{"Layer": "1", "SpreadingTypeXOL": "QS 2024", "SpreadingListXOL": []any{map[string]any{"ReinsTypeName": "QS (OR)", "Pct": "40"}}}},
			"EGNPI": {{"TreatyGroup": "FIRE", "Currency": "IDR", "Amount": "400", "AmountIDR": "400"}},
		})
}

func actualUji() HalamanPohon {
	return pohon(nil, map[string][]map[string]any{
		"EGNPI":  {{"TreatyGroup": "FIRE", "Currency": "IDR", "CurrencyID": "10026", "Amount": "500", "AmountIDR": "500"}},
		"Limits": {{"Layer": "LAMA"}},
	})
}

// ⭐ Actual GNPI · Update Total: Limits Actual = salinan Limits Treaty In;
// EGNPI layer dari Actual GNPI ber-Treaty Group sama; PE = EGNPI × AdjRate%;
// MDP = PE (MDP% dipaksa 100); selisih premi dinolkan bila Actual < estimasi.
func TestAktualNilai(t *testing.T) {
	h := hitungAktualUji(t, MasukanAktual{Aksi: AksiAktualNilai, Akar: akarAktualUji(), Actual: actualUji()})
	l := h.Actual.Larik["Limits"][0]
	if teksSimpul(l, "Layer") != "1" || teksSimpul(l, "MDPPct") != "100" || teksSimpul(l, "AddendumStatus") != "1" {
		t.Fatalf("layer Actual %v", l)
	}
	if got := nilaiLarik(larikSimpul(l, "EgnpiTotalList"), "Value"); !reflect.DeepEqual(got, []string{"500"}) {
		t.Fatalf("EgnpiTotalList %v", got)
	}
	if got := nilaiLarik(larikSimpul(l, "PremiumEarnedList"), "Value"); !reflect.DeepEqual(got, []string{"50"}) {
		t.Fatalf("PremiumEarnedList %v", got)
	}
	if got := nilaiLarik(larikSimpul(l, "MDPList"), "Value"); !reflect.DeepEqual(got, []string{"50"}) {
		t.Fatalf("MDPList %v", got)
	}
	if h.Actual.Medan["TotalEgnpiAmount"] != "500" || h.Actual.Medan["TotalEgnpiProportion"] != "100" {
		t.Fatalf("total Actual GNPI %v", h.Actual.Medan)
	}
	if got := nilaiLarik(h.Actual.Larik["TotalLimitPremiEarnNP"], "Value"); !reflect.DeepEqual(got, []string{"50"}) {
		t.Fatalf("TotalLimitPremiEarnNP %v", got)
	}
	if h.Actual.Medan["TotalLimitsROL"] != "2" {
		t.Fatalf("TotalLimitsROL %q", h.Actual.Medan["TotalLimitsROL"])
	}
	if h.Selisih == nil {
		t.Fatal("ValueDifference tidak disusun")
	}
	vl := h.Selisih.Larik["Limits"][0]
	if got := nilaiLarik(larikSimpul(vl, "PremiumEarnedList"), "Value"); !reflect.DeepEqual(got, []string{"10"}) {
		t.Fatalf("selisih PE %v", got)
	}
	// Share Actual: RNM 40% atas Limit 1000; QS (OR) 40% dari Share Treaty In.
	s := h.Actual.Larik["Share"][0]
	if got := nilaiLarik(larikSimpul(s, "RnmLimitList"), "Value"); !reflect.DeepEqual(got, []string{"400"}) {
		t.Fatalf("RnmLimitList %v", got)
	}
	if got := nilaiLarik(larikSimpul(s, "RNMSpreadedListXOL"), "Value"); !reflect.DeepEqual(got, []string{"160"}) {
		t.Fatalf("RNMSpreadedListXOL %v", got)
	}
}

// ⛔ PremiumIncome keluar bila Total EGNPI TREATY IN nol — proporsi Actual tidak dihitung.
func TestAktualPremiKeluarBilaTotalTreatyInNol(t *testing.T) {
	akar := akarAktualUji()
	akar.Medan["TotalEgnpiAmount"] = ""
	h := hitungAktualUji(t, MasukanAktual{Aksi: AksiAktualNilai, Akar: akar, Actual: actualUji()})
	if _, ada := h.Actual.Medan["TotalEgnpiProportion"]; ada || len(h.Actual.Larik["TotalEgnpiAmountNP"]) != 0 {
		t.Fatalf("langkah 4 berjalan: %v", h.Actual.Medan)
	}
	if h.Actual.Medan["TotalEgnpiAmount"] != "500" {
		t.Fatalf("langkah 2 tetap berjalan: %v", h.Actual.Medan)
	}
}

// ⚠️ Fac Share > 0: RnmShareDeducted AKAR ditulis; baris Facultative lahir.
func TestAktualShareFacultative(t *testing.T) {
	akar := akarAktualUji()
	akar.Medan["FacultativeShare"] = "10"
	act := actualUji()
	act.Larik["Limits"] = akar.Larik["Limits"]
	h := hitungAktualUji(t, MasukanAktual{Aksi: AksiAktualShare, Akar: akar, Actual: act})
	if h.Akar.Medan["RnmShareDeducted"] != "30" {
		t.Fatalf("RnmShareDeducted %v", h.Akar.Medan)
	}
	if len(h.Actual.Larik["FacultativeShareList"]) != 1 {
		t.Fatalf("FacultativeShareList %v", h.Actual.Larik["FacultativeShareList"])
	}
	// Share = 30% (40 − 10) atas Limit 1000.
	if got := nilaiLarik(larikSimpul(h.Actual.Larik["Share"][0], "RnmLimitList"), "Value"); !reflect.DeepEqual(got, []string{"300"}) {
		t.Fatalf("RnmLimitList %v", got)
	}
	// ⚠️ ringkasan facultative AKAR dibuang (SummaryLimitActualFacShare).
	if got, ada := h.Akar.Larik["LimitFacShareSummaryList"]; !ada || len(got) != 0 {
		t.Fatalf("LimitFacShareSummaryList akar %v", got)
	}
	if h.Selisih == nil || h.Selisih.Medan["FacultativeShare"] != "0" {
		t.Fatalf("selisih FacultativeShare %v", h.Selisih)
	}
}

// ⚠️ Update Total Actual Limits: baris mata uang ditambahkan ke total PE/MDP
// AKAR; total facultative AKAR dibuang.
func TestAktualLimitsJanggalDisalin(t *testing.T) {
	akar := akarAktualUji()
	act := pohon(nil, map[string][]map[string]any{"Limits": {{
		"Currency": "IDR", "Limit": "100", "PremiumEarnedList": []any{map[string]any{"Currency": "USD", "Value": "3"}},
	}}})
	h := hitungAktualUji(t, MasukanAktual{Aksi: AksiAktualLimits, Akar: akar, Actual: act})
	if got := h.Akar.Larik["TotalLimitPremiEarnNP"]; len(got) != 1 || teksSimpul(got[0], "Currency") != "USD" {
		t.Fatalf("TotalLimitPremiEarnNP akar %v", got)
	}
	if got, ada := h.Akar.Larik["TotalFacShareNetNP"]; !ada || len(got) != 0 {
		t.Fatalf("TotalFacShareNetNP akar %v", got)
	}
	if got := nilaiLarik(h.Actual.Larik["TotalLimitIOONP"], "Value"); !reflect.DeepEqual(got, []string{"100", "0"}) {
		t.Fatalf("TotalLimitIOONP %v", got)
	}
	if len(h.Actual.Larik["LimitSummaryList"]) != 1 || h.Selisih != nil {
		t.Fatalf("ringkasan %v / selisih %v", h.Actual.Larik["LimitSummaryList"], h.Selisih)
	}
}

// ⚠️ % RNM Share rincian (EDMState 3): baris sebaran DITAMBAHKAN ke total AKAR.
func TestAktualBarisMenambahTotalAkar(t *testing.T) {
	akar := akarAktualUji()
	akar.Larik["TotalSpreadedNetPremi"] = []map[string]any{{"Currency": "IDR", "Value": "1"}}
	act := pohon(nil, map[string][]map[string]any{
		"Limits": {{"Currency": "IDR", "Limit": "1000", "MDPList": []any{map[string]any{"Currency": "IDR", "Value": "100"}}}},
		"Share":  {{"RNMShare": "50"}},
	})
	h := hitungAktualUji(t, MasukanAktual{Aksi: AksiAktualBaris, Akar: akar, Actual: act, Indeks: 0})
	s := h.Actual.Larik["Share"][0]
	if got := nilaiLarik(larikSimpul(s, "GrossPremiumList"), "Value"); !reflect.DeepEqual(got, []string{"50"}) {
		t.Fatalf("Gross %v", got)
	}
	if got := h.Akar.Larik["TotalSpreadedNetPremi"]; len(got) != 2 {
		t.Fatalf("TotalSpreadedNetPremi akar %v", got)
	}
}

func TestAktualAksiTakDikenal(t *testing.T) {
	if _, err := HitungAktual(MasukanAktual{Aksi: "x"}); !errors.Is(err, ErrMasukanTidakSah) {
		t.Fatalf("err %v", err)
	}
}
