package services

import (
	"errors"
	"reflect"
	"testing"
)

func pohon(medan map[string]string, larik map[string][]map[string]any) HalamanPohon {
	h := halamanKosong()
	for k, v := range medan {
		h.Medan[k] = v
	}
	for k, v := range larik {
		h.Larik[k] = v
	}
	return h
}

func nilaiLarik(xs []map[string]any, k string) []string {
	out := []string{}
	for _, x := range xs {
		out = append(out, teksSimpul(x, k))
	}
	return out
}

func hitungSelisihUji(t *testing.T, akar, lama, act HalamanPohon) HasilSelisih {
	t.Helper()
	h, err := HitungSelisih(MasukanSelisih{Aksi: AksiSelisihEDM, Akar: akar, Lama: lama, Actual: act})
	if err != nil {
		t.Fatalf("HitungSelisih: %v", err)
	}
	return h
}

// ⭐ Premium: selisih per indeks, total Amount IDR, proporsi @divide(…,20)×100,
// total per mata uang.
func TestSelisihPremi(t *testing.T) {
	akar := pohon(nil, map[string][]map[string]any{"EGNPI": {
		{"Currency": "IDR", "TreatyGroup": "FIRE", "Amount": "300", "AmountIDR": "300", "Proportion": "60"},
		{"Currency": "USD", "TreatyGroup": "MOTOR", "Amount": "20", "AmountIDR": "300"},
	}})
	lama := pohon(nil, map[string][]map[string]any{"EGNPI": {{"Amount": "100", "AmountIDR": "100", "Proportion": "50"}}})
	h := hitungSelisihUji(t, akar, lama, halamanKosong())
	vd := h.Selisih
	if got := nilaiLarik(vd.Larik["EGNPI"], "Amount"); !reflect.DeepEqual(got, []string{"200", "20"}) {
		t.Fatalf("Amount %v", got)
	}
	if vd.Medan["TotalEgnpiAmount"] != "500" {
		t.Fatalf("TotalEgnpiAmount %q", vd.Medan["TotalEgnpiAmount"])
	}
	// 200/500 = 0.4 → 40; 300/500 → 60.
	if got := nilaiLarik(vd.Larik["EGNPI"], "Proportion"); !reflect.DeepEqual(got, []string{"40", "60"}) {
		t.Fatalf("Proportion %v", got)
	}
	if vd.Medan["TotalEgnpiProportion"] != "100" {
		t.Fatalf("TotalEgnpiProportion %q", vd.Medan["TotalEgnpiProportion"])
	}
	if got := nilaiLarik(vd.Larik["TotalEgnpiAmountNP"], "Currency"); !reflect.DeepEqual(got, []string{"IDR", "USD"}) {
		t.Fatalf("TotalEgnpiAmountNP %v", vd.Larik["TotalEgnpiAmountNP"])
	}
}

// ⛔ Total selisih 0 → langkah 3 KELUAR: proporsi tidak dihitung ulang.
func TestSelisihPremiTotalNolKeluar(t *testing.T) {
	akar := pohon(nil, map[string][]map[string]any{"EGNPI": {{"Currency": "IDR", "AmountIDR": "100", "Proportion": "7"}}})
	lama := pohon(nil, map[string][]map[string]any{"EGNPI": {{"AmountIDR": "100", "Proportion": "2"}}})
	vd := hitungSelisihUji(t, akar, lama, halamanKosong()).Selisih
	if vd.Medan["TotalEgnpiAmount"] != "0" || teksSimpul(vd.Larik["EGNPI"][0], "Proportion") != "5" {
		t.Fatalf("medan %v, EGNPI %v", vd.Medan, vd.Larik["EGNPI"])
	}
	if _, ada := vd.Medan["TotalEgnpiProportion"]; ada || len(vd.Larik["TotalEgnpiAmountNP"]) != 0 {
		t.Fatalf("langkah 3 berjalan: %v", vd)
	}
}

// ⚠️ Limits: Currency Premium Earned ke EgnpiTotalList (disalin apa adanya);
// ROL selisih dari ActualValue.
func TestSelisihLimitsJanggalDisalin(t *testing.T) {
	akar := pohon(nil, map[string][]map[string]any{"Limits": {{
		"Layer": "1", "Limit": "1000", "AdjRate": "5",
		"EgnpiTotalList":    []any{map[string]any{"Currency": "IDR", "Value": "50"}},
		"PremiumEarnedList": []any{map[string]any{"Currency": "USD", "Value": "9"}},
		"MDPList":           []any{map[string]any{"Currency": "IDR", "Value": "8"}},
	}}})
	lama := pohon(map[string]string{"TotalLimitsROL": "2"}, map[string][]map[string]any{"Limits": {{
		"Limit": "400", "AdjRate": "5",
		"EgnpiTotalList":    []any{map[string]any{"Value": "20"}},
		"PremiumEarnedList": []any{map[string]any{"Value": "4"}},
	}}})
	act := pohon(map[string]string{"TotalLimitsROL": "7"}, nil)
	vd := hitungSelisihUji(t, akar, lama, act).Selisih
	l := vd.Larik["Limits"][0]
	if teksSimpul(l, "Limit") != "600" || teksSimpul(l, "AdjRate") != "0" || teksSimpul(l, "Layer") != "1" {
		t.Fatalf("layer %v", l)
	}
	egnpi := larikSimpul(l, "EgnpiTotalList")
	if teksSimpul(egnpi[0], "Currency") != "USD" || teksSimpul(egnpi[0], "Value") != "30" {
		t.Fatalf("EgnpiTotalList %v", egnpi)
	}
	pe := larikSimpul(l, "PremiumEarnedList")
	if teksSimpul(pe[0], "Value") != "5" || teksSimpul(pe[0], "Currency") != "" {
		t.Fatalf("PremiumEarnedList %v", pe)
	}
	if vd.Medan["TotalLimitsROL"] != "5" {
		t.Fatalf("TotalLimitsROL %q", vd.Medan["TotalLimitsROL"])
	}
}

// ⭐ Share: EDMState 3 → skalar dari ActualValue; Deduction membuang Brokerage
// fee bernilai kosong, Net disusun dari Gross selisih; total tanpa mata uang
// kosong.
func TestSelisihShareDanDeduksi(t *testing.T) {
	akar := pohon(map[string]string{"EDMState": "3", "RNMShare": "99"}, map[string][]map[string]any{"Share": {{
		"Layer":            "1",
		"RnmLimitList":     []any{map[string]any{"Currency": "IDR", "Value": "500"}},
		"GrossPremiumList": []any{map[string]any{"Currency": "IDR", "Value": "100"}},
		"NetPremiumList":   []any{map[string]any{"Currency": "IDR", "Value": "90"}},
		"DeductionList":    []any{map[string]any{"Currency": "IDR", "Deduction": "10", "DeductionPct": "10"}},
	}}})
	lama := pohon(map[string]string{"RNMShare": "40", "BrokeragePercent": "10"}, map[string][]map[string]any{"Share": {{
		"GrossPremiumList": []any{map[string]any{"Currency": "IDR", "Value": "60"}},
		"NetPremiumList":   []any{map[string]any{"Currency": "IDR", "Value": "54"}},
	}}})
	act := pohon(map[string]string{"RNMShare": "50"}, nil)
	h := hitungSelisihUji(t, akar, lama, act)
	vd := h.Selisih
	if vd.Medan["RNMShare"] != "10" {
		t.Fatalf("RNMShare %q (mau dari ActualValue 50 − 40)", vd.Medan["RNMShare"])
	}
	s := vd.Larik["Share"][0]
	if got := larikSimpul(s, "DeductionList"); len(got) != 0 {
		t.Fatalf("Brokerage fee tidak dibuang: %v", got)
	}
	// Net = Gross selisih (40) − Σ deduksi bermata uang IDR (nol).
	if net := larikSimpul(s, "NetPremiumList"); len(net) != 1 || teksSimpul(net[0], "Value") != "40" {
		t.Fatalf("NetPremiumList %v", net)
	}
	tampil := larikSimpul(s, "RnmLimitListDisplay")
	if len(tampil) != 2 || teksSimpul(tampil[0], "Value") != "500" || teksSimpul(tampil[1], "Currency") != "" {
		t.Fatalf("RnmLimitListDisplay %v", tampil)
	}
	// DeductionTotalList ber-mata uang kosong TIDAK masuk TotalShareDeductionNP.
	if got := vd.Larik["TotalShareDeductionNP"]; len(got) != 0 {
		t.Fatalf("TotalShareDeductionNP %v", got)
	}
	if got := nilaiLarik(vd.Larik["TotalShareGrossNP"], "Value"); !reflect.DeepEqual(got, []string{"40"}) {
		t.Fatalf("TotalShareGrossNP %v", got)
	}
	if len(vd.Larik["LimitShareSummaryList"]) != 1 {
		t.Fatalf("LimitShareSummaryList %v", vd.Larik["LimitShareSummaryList"])
	}
}

// ⭐ Installment: jadwal disalin per mata uang, Amount = @divide(Net × %, 100, 4).
func TestSelisihAngsuran(t *testing.T) {
	akar := pohon(nil, map[string][]map[string]any{
		// Net selisih disusun ulang dari GROSS selisih (Deduction [3]).
		"Share": {{"GrossPremiumList": []any{map[string]any{"Currency": "IDR", "Value": "100.005"}}}},
		"Installment": {{"Currency": "IDR", "InstallmentList": []any{
			map[string]any{"Installment": "1", "InstallmentPct": "50", "DueDate": "20240101"},
			map[string]any{"Installment": "2", "InstallmentPct": "50"},
		}}},
	})
	vd := hitungSelisihUji(t, akar, halamanKosong(), halamanKosong()).Selisih
	inst := vd.Larik["Installment"]
	if len(inst) != 1 {
		t.Fatalf("Installment %v", inst)
	}
	if got := nilaiLarik(larikSimpul(inst[0], "InstallmentList"), "Amount"); !reflect.DeepEqual(got, []string{"50.0025", "50.0025"}) {
		t.Fatalf("Amount %v", got)
	}
	if teksSimpul(inst[0], "AmountTotal") != "100.005" || teksSimpul(inst[0], "PctTotal") != "100" {
		t.Fatalf("total %v", inst[0])
	}
	if got := nilaiLarik(vd.Larik["TotalInstallmentNP"], "Value"); !reflect.DeepEqual(got, []string{"100.005"}) {
		t.Fatalf("TotalInstallmentNP %v", got)
	}
}

// ⭐ Pro Rate: ValueBeforeProrate = sebelum; Share × Pro Rate %, deduksi TIDAK
// (langkahnya berblok `//`).
func TestSelisihProrata(t *testing.T) {
	akar := pohon(map[string]string{"IsProRate": "true", "ProRatePercent": "50"}, map[string][]map[string]any{"Share": {{
		"GrossPremiumList":   []any{map[string]any{"Currency": "IDR", "Value": "100"}},
		"DeductionTotalList": []any{map[string]any{"Currency": "IDR", "Value": "10"}},
	}}})
	h := hitungSelisihUji(t, akar, halamanKosong(), halamanKosong())
	if h.SebelumProrata == nil {
		t.Fatal("ValueBeforeProrate tidak ada")
	}
	if got := nilaiLarik(larikSimpul(h.SebelumProrata.Larik["Share"][0], "GrossPremiumList"), "Value"); !reflect.DeepEqual(got, []string{"100"}) {
		t.Fatalf("sebelum %v", got)
	}
	s := h.Selisih.Larik["Share"][0]
	if got := nilaiLarik(larikSimpul(s, "GrossPremiumList"), "Value"); !reflect.DeepEqual(got, []string{"50"}) {
		t.Fatalf("Gross %v", got)
	}
	if got := nilaiLarik(h.Selisih.Larik["TotalShareGrossNP"], "Value"); !reflect.DeepEqual(got, []string{"50"}) {
		t.Fatalf("TotalShareGrossNP %v", got)
	}
	tanpa := hitungSelisihUji(t, pohon(map[string]string{"IsProRate": "false"}, nil), halamanKosong(), halamanKosong())
	if tanpa.SebelumProrata != nil {
		t.Fatal("ValueBeforeProrate lahir tanpa Pro Rate")
	}
}

func TestSelisihAksiTakDikenal(t *testing.T) {
	if _, err := HitungSelisih(MasukanSelisih{Aksi: "x"}); !errors.Is(err, ErrMasukanTidakSah) {
		t.Fatalf("err %v", err)
	}
}
