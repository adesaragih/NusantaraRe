package services_test

// Uji sub-tab Deduction dan Reserve tab Limits Prop — `hitung_limit_rinci.go`.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// ⭐ Rincian Prop TIDAK punya Gross: pesan Activity muncul, baris TIDAK
// dirusak, total per mata uang tetap dihitung.
func TestDeduksiPropTanpaGross(t *testing.T) {
	m := services.MasukanDeduksi{Sts: "pct", Indeks: 0, DeductionList: []services.BarisDeduksi{
		{Currency: "IDR", Deduction: "10", DeductionPct: "5"},
		{Currency: "USD", Deduction: "2"},
		{Currency: "IDR", Deduction: "4"},
	}}
	h := services.HitungDeduksi(m)
	if len(h.Pesan) != 1 || h.Pesan[0] != services.PesanGrossKosong {
		t.Errorf("pesan %v", h.Pesan)
	}
	if h.DeductionList[0].Currency != "IDR" || h.DeductionList[0].Deduction != "10" {
		t.Errorf("baris dirusak: %+v", h.DeductionList[0])
	}
	tot := h.DeductionTotalList
	if len(tot) != 2 || tot[0] != nmu("IDR", "14") || tot[1] != nmu("USD", "2") {
		t.Errorf("total %v", tot)
	}
}

// Dengan Gross (bentuk Share XOL): persen → nilai, baris per mata uang, Net.
func TestDeduksiPersenDenganGross(t *testing.T) {
	m := services.MasukanDeduksi{Sts: "pct", Indeks: 0,
		DeductionList:    []services.BarisDeduksi{{Comment: "Brokerage", DeductionPct: "10"}},
		GrossPremiumList: []services.NilaiMataUang{nmu("IDR", "1000"), nmu("USD", "50")}}
	h := services.HitungDeduksi(m)
	if len(h.DeductionList) != 2 || h.DeductionList[0].Deduction != "100" || h.DeductionList[1].Deduction != "5" ||
		h.DeductionList[1].Comment != "Brokerage" {
		t.Fatalf("%+v", h.DeductionList)
	}
	if h.NetPremiumList[0] != nmu("IDR", "900") || h.NetPremiumList[1] != nmu("USD", "45") {
		t.Errorf("net %v", h.NetPremiumList)
	}
}

func TestDeduksiNilaiMembuangPersen(t *testing.T) {
	h := services.HitungDeduksi(services.MasukanDeduksi{Sts: "val", Indeks: 0,
		DeductionList: []services.BarisDeduksi{{Currency: "IDR", Deduction: "3", DeductionPct: "9"}}})
	if h.DeductionList[0].DeductionPct != "" {
		t.Errorf("%+v", h.DeductionList[0])
	}
}

func TestCadanganPremi(t *testing.T) {
	ces := []services.NilaiMataUang{nmu("IDR", "2000"), nmu("USD", "40")}
	h := services.HitungCadangan(services.MasukanCadangan{PremiumReservePct: "25", CessionList: ces})
	if len(h.ReserveList) != 2 || h.ReserveList[0] != nmu("IDR", "500") || h.ReserveList[1] != nmu("USD", "10") || len(h.Pesan) != 0 {
		t.Errorf("%v %v", h.ReserveList, h.Pesan)
	}
	h = services.HitungCadangan(services.MasukanCadangan{PremiumReservePct: "101", CessionList: ces})
	if len(h.Pesan) != 1 || h.Pesan[0] != services.PesanCadanganLebih || len(h.ReserveList) != 2 || h.ReserveList[0] != nmu("", "") {
		t.Errorf("> 100: %v %v", h.ReserveList, h.Pesan)
	}
}
