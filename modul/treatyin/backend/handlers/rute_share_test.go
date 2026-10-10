package handlers_test

// Uji rute tab Share Non-Prop — `POST /hitung/share-np` (alur Update
// Summary dari NOL: layer diketik di tab Limits, tanpa data tersimpan) dan
// kedua daftar pilihan.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

func TestRuteHitungShareNPDariNol(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	const jalur = "/api/treaty-in/hitung/share-np"
	// Update Summary: RNM 40, Brokerage 10, satu layer 1 M dengan MDP 50 jt.
	isi := `{"aksi":"summary","share":{"RNMShare":"40","BrokeragePercent":"10","RNMShareAcrossTheBoard":"true","Total":{}},
		"layers":[{"LayerType":"Layer","Layer":"1","LayerPartType":"Layer","LayerPart":"1","Currency":"IDR","Limit":"1000000000",
		"MDPList":[{"Currency":"IDR","Value":"50000000"}],"TreatyGroupList":[{"TreatyGroup":"PROPERTY","TreatyGroupID":"10002"}]}]}`
	if w := kirim(t, h, http.MethodPost, jalur, "", badan(isi)); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", w.Code)
	}
	if w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{`)); w.Code != http.StatusBadRequest {
		t.Errorf("badan rusak: %d", w.Code)
	}
	w := kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(isi))
	b := w.Body.String()
	if w.Code != http.StatusOK {
		t.Fatalf("%d %s", w.Code, b)
	}
	for _, mau := range []string{
		`"RnmLimitList":[{"Currency":"IDR","CurrencyID":"","Value":"400000000"}]`,    // Limit × 40%
		`"GrossPremiumList":[{"Currency":"IDR","CurrencyID":"","Value":"20000000"}]`, // MDP × 40%
		`"Comment":"Brokerage fee"`,
		`"NetPremiumList":[{"Currency":"IDR","CurrencyID":"","Value":"18000000"}]`, // Gross − 10%
		`"TotalSpreadedRnmProp":[]`, // nol Spreading Type → nol bagian OR
		// Baris tanpa Spreading Type → `SetSpreadingXOL` [9]: total spreading
		// 0 ≠ RNM Share 40 — pesan Activity, persis Pega.
		services.PesanShareSpreading,
	} {
		if !strings.Contains(b, mau) {
			t.Errorf("jawaban tanpa %s\n%s", mau, b)
		}
	}
	// RNM Share kosong → pesan Activity, Share kosong.
	w = kirim(t, h, http.MethodPost, jalur, "AKUN-UJI", badan(`{"aksi":"summary","share":{"Total":{}},"layers":[{"Limit":"1"}]}`))
	if !strings.Contains(w.Body.String(), services.PesanShareNilaiKosong) {
		t.Errorf("RNM kosong: %s", w.Body.String())
	}
}

func TestRutePilihanShareNP(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	for _, jalur := range []string{"/api/treaty-in/warisan/spreading-induk?treatyGroupId=10002&mulai=20230401", "/api/treaty-in/warisan/reasuradur-share"} {
		if w := kirim(t, h, http.MethodGet, jalur, "", nil); w.Code != http.StatusUnauthorized {
			t.Errorf("%s tanpa identitas: %d", jalur, w.Code)
		}
		if w := kirim(t, h, http.MethodGet, jalur, "AKUN-UJI", nil); w.Code != http.StatusOK {
			t.Errorf("%s: %d %s", jalur, w.Code, w.Body.String())
		}
	}
}
