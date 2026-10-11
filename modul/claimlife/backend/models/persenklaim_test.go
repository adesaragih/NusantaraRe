package models

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/uang"
)

// Keputusan work owner 11-10-2026: "share nusantara re nya 100 dan claim rnm nya adalah 50. berarti percent claim nya
// adalah 50% dan claim paid nya 50. sehingga claim rnm tidak boleh melebihi share nusantara re".
func TestPersenKlaimDiturunkanDariClaimRNM(t *testing.T) {
	m := func(s string) uang.Money {
		v, err := uang.NewMoney(s, "IDR")
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	for _, u := range []struct{ klaim, share, mau string }{
		{"50", "100", "50"},
		{"100", "100", "100"},
		{"1", "3", "33.333"},
		{"0", "100", "0"},
		// Kasus kesamaan dengan frontend desimal.ts persenKlaimTeks (11-10-2026) - pecahan half-up 5 desimal.
		{"2", "3", "66.667"},
		{"665", "5850038.75", "0.011"},
		{"1", "200000", "0.001"},
		{"1", "400000", "0"},
		{"94695.705", "100000", "94.696"},
	} {
		got, err := HitungPersenKlaim(m(u.klaim), m(u.share))
		if err != nil || got != u.mau {
			t.Errorf("%s/%s = %q (%v), mau %q", u.klaim, u.share, got, err, u.mau)
		}
	}
	if got, _ := HitungPersenKlaim(m("50"), uang.Money{}); got != "" {
		t.Errorf("share kosong: %q, mau kosong", got)
	}
	if got, _ := HitungPersenKlaim(m("50"), m("0")); got != "" {
		t.Errorf("share nol: %q, mau kosong", got)
	}
	// Batas: Claim RNM > Share RNM ditolak (tanpa pengecualian kode bisnis bila kode kosong).
	if err := PeriksaClaimGrossTerhadapShare(m("101"), m("100"), ""); !errors.Is(err, ErrClaimGrossMelebihiShare) {
		t.Errorf("melebihi share lolos: %v", err)
	}
	if err := PeriksaClaimGrossTerhadapShare(m("100"), m("100"), ""); err != nil {
		t.Errorf("sama dengan share ditolak: %v", err)
	}
	// JSON baris membawa turunan: persenKlaim 50, klaimDibayar = Claim RNM.
	b, err := json.Marshal(BarisAdjustment{ID: "UJI-A1", JumlahKlaim: m("50"), ShareNusantaraRe: m("100")})
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{`"persenKlaim":"50"`, `"klaimDibayar":{"amount":"50"`, `"shareNusantaraRe":{"amount":"100"`} {
		if !strings.Contains(string(b), k) {
			t.Errorf("JSON baris tanpa %s: %s", k, b)
		}
	}
}

// 11-10-2026: JSON baris membawa OutStanding Date (jejak), data bank, dan spreading per lapisan + rekap per reinsurer
// (dijumlahkan backend dengan Money.Add).
func TestJSONBarisBankDanSpreading(t *testing.T) {
	m := func(s string) uang.Money {
		v, err := uang.NewMoney(s, "IDR")
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	retro := func(nama, gross, net string) SpreadingRetro {
		return SpreadingRetro{ReinsurerName: nama, PremiumSpreadedGross: m(gross), PremiumSpreadedNet: m(net), Amount: m(gross)}
	}
	b, err := json.Marshal(BarisAdjustment{ID: "UJI-A1", TanggalOutstanding: "2026-10-11 09:00:00", PayableTo: "UJI PAY",
		NamaBank: "UJI BANK", CabangBank: "UJI CABANG", NomorRekening: "000123"})
	if err != nil {
		t.Fatal(err)
	}
	// 11-10-2026: spreading tidak lagi di JSON baris (diambil dari PremiumList per peserta).
	if strings.Contains(string(b), "\"spreading\"") {
		t.Errorf("JSON baris masih membawa spreading: %s", b)
	}
	sp, err := json.Marshal(SusunSpreading([]Spreading{
		{TreatyTypeName: "QS", RetrocadedShare: m("1000.5"), Retro: []SpreadingRetro{retro("UJI RE A", "100", "90"), retro("UJI RE B", "50", "45")}},
		{TreatyTypeName: "SURPLUS", RetrocadedShare: m("200"), Retro: []SpreadingRetro{retro("UJI RE A", "20.5", "18")}},
	}))
	if err != nil {
		t.Fatal(err)
	}
	b = append(b, sp...)
	for _, k := range []string{`"tanggalOutstanding":"2026-10-11 09:00:00"`, `"payableTo":"UJI PAY"`, `"namaBank":"UJI BANK"`,
		`"cabangBank":"UJI CABANG"`, `"nomorRekening":"000123"`, `"treatyTypeName":"QS"`,
		`"retrocededShare":{"amount":"1000.5"`, `"nama":"UJI RE B"`,
		// Rekap: UJI RE A = 100 + 20.5 gross, 90 + 18 net.
		`{"reinsurer":"UJI RE A","currency":"IDR","gross":{"amount":"120.5"`, `"net":{"amount":"108"`,
		// TOTAL per mata uang: 120.5 + 50 gross, 108 + 45 net.
		`"total":[{"currency":"IDR","gross":{"amount":"170.5","currency":"IDR"},"net":{"amount":"153","currency":"IDR"}}]`} {
		if !strings.Contains(string(b), k) {
			t.Errorf("JSON baris tanpa %s: %s", k, b)
		}
	}
	kosong, _ := json.Marshal(SusunSpreading(nil))
	if string(kosong) != `{"lapisan":[],"rekap":[],"total":[]}` {
		t.Errorf("spreading kosong harus daftar kosong, bukan null: %s", kosong)
	}
}
