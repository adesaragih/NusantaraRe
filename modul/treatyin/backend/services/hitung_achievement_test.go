package services_test

// Uji `GetAchievement` — `hitung_achievement.go`.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func barisAch(grup, kuartal, tahun, cur, idCur, prem, net, paid, os string) models.BarisAchievement {
	return models.BarisAchievement{TreatyType: "QUOTA SHARE", TreatyGroupName: grup, Quarter: kuartal, QuarterYear: tahun,
		Currency: cur, IDCurrency: idCur, Premium: prem, NetPremium: net, PaidClaim: paid, OutstandingClaim: os, NoPolis: "P1"}
}

func TestAchievementTempelHitungDanTotal(t *testing.T) {
	m := services.MasukanAchievement{
		RNMShareP: "50",
		Limits: []services.LimitAchievement{{TreatyType: "QUOTA SHARE", Detail: []services.DetailAchievement{
			{TreatyGroup: "PROPERTY", EPIList: []services.NilaiMataUang{{Currency: "IDR", Value: "1000"}}},
			{TreatyGroup: "PROPERTY"},
			{TreatyGroup: "MARINE"},
		}}},
	}
	baris := []models.BarisAchievement{
		barisAch("PROPERTY", "1", "2024", "IDR", "10026", "100", "80", "10", "6"),
		barisAch("PROPERTY", "2", "2024", "USD", "10001", "2", "1", "0", "0"),
	}
	h := services.HitungAchievement(m, baris, map[string]string{"10001": "15000"})
	d := h.Limits[0]
	// Detail 1: dua baris + baris total.
	if len(d[0].AchievementLists) != 3 {
		t.Fatalf("detail 1: %d baris", len(d[0].AchievementLists))
	}
	r := d[0].AchievementLists[0]
	if r["Quarter"] != "Q 1" || r["IncuredClaim"] != "16" || r["Total"] != "64" || r["LossRatio"] != "20" {
		t.Errorf("baris IDR %v", r)
	}
	if d[0].AchievementLists[1]["PREMIUMtoIDR"] != "30000" {
		t.Errorf("konversi USD %v", d[0].AchievementLists[1])
	}
	// Detail 2 (grup sama dengan sebelumnya): baris NOL.
	if z := d[1].AchievementLists[0]; z["PREMIUM"] != "0" || z["FlagTreaty"] != "1" {
		t.Errorf("baris nol %v", z)
	}
	// Total & Achievement % = (30100 / 0,5) / 1000 × 100.
	if d[0].TotalAchPremium != "30100" || d[0].AchievementPct != "6020" {
		t.Errorf("total %s ach %s", d[0].TotalAchPremium, d[0].AchievementPct)
	}
	if tot := d[0].AchievementLists[2]; tot["Currency"] != " Total In IDR" || tot["PREMIUM"] != "30100" {
		t.Errorf("baris total %v", tot)
	}
	if len(d[0].CurrencyList) != 2 || d[0].CurrencyList[0]["Parameter"] != " Based on Gross " {
		t.Errorf("CurrencyList %v", d[0].CurrencyList)
	}
	if !h.FlagExcel || len(h.Kuartal) != 2 || len(h.TahunKuartal) != 1 {
		t.Errorf("flag %v kuartal %v tahun %v", h.FlagExcel, h.Kuartal, h.TahunKuartal)
	}
}

// Saringan As At / Quarter Year; RNMShareP kosong → Achievement % kosong.
func TestAchievementSaringanDanTanpaShare(t *testing.T) {
	m := services.MasukanAchievement{Cari: true, AsAt: "1", Tahun: "2024",
		Limits: []services.LimitAchievement{{TreatyType: "QUOTA SHARE", Detail: []services.DetailAchievement{
			{TreatyGroup: "PROPERTY", EPIList: []services.NilaiMataUang{{Currency: "IDR", Value: "1000"}}}}}}}
	baris := []models.BarisAchievement{
		barisAch("PROPERTY", "1", "2024", "IDR", "10026", "100", "80", "0", "0"),
		barisAch("PROPERTY", "2", "2024", "IDR", "10026", "100", "80", "0", "0"),
		barisAch("PROPERTY", "1", "2025", "IDR", "10026", "100", "80", "0", "0"),
	}
	h := services.HitungAchievement(m, baris, nil)
	if n := len(h.Limits[0][0].AchievementLists); n != 2 { // satu baris lolos + total
		t.Errorf("%d baris", n)
	}
	if h.Limits[0][0].AchievementPct != "" {
		t.Errorf("Achievement %% tanpa RNMShareP: %q", h.Limits[0][0].AchievementPct)
	}
}
