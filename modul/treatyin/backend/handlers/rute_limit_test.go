package handlers_test

// Uji rute rumus tab Limits — Non-Prop, Deduction, Reserve, dan
// autocomplete Class of Business.

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

func TestRuteHitungLimitTabLimits(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	for _, c := range []struct {
		jalur, isi, mau string
	}{
		{"/api/treaty-in/hitung/limit-np",
			`{"aksi":"adj","indeks":0,"kurs":[{"Currency":"IDR","Conversion":"1"}],"layers":[{"Limit":"1000","AdjRate":"10","EgnpiTotalList":[{"Currency":"IDR","Value":"500"}]}]}`,
			`"ROLPct":"5"`},
		{"/api/treaty-in/hitung/limit-deduksi",
			`{"sts":"","indeks":0,"DeductionList":[{"Currency":"IDR","Deduction":"3"},{"Currency":"IDR","Deduction":"4"}]}`,
			`"DeductionTotalList":[{"Currency":"IDR","CurrencyID":"","Value":"7"}]`},
		{"/api/treaty-in/hitung/limit-cadangan",
			`{"PremiumReservePct":"10","CessionList":[{"Currency":"IDR","Value":"500"}]}`,
			`"Value":"50"`},
		// Tab EGNPI - `TreatyInEGNPIListValue` lalu `TreatyInNPSetTotal`:
		// 2 x 100 = 200 IDR, dan satu-satunya baris memegang 100 % proporsi.
		// Tab Maximum Retention - `TreatyInNPSetTotal(retention)`:
		// dua baris IDR digabung menjadi satu, 3,5 M + 0,5 M.
		{"/api/treaty-in/hitung/retensi",
			`{"aksi":"total","retensi":[{"Currency":"IDR","Amount":"3500000000"},{"Currency":"IDR","Amount":"500000000"}]}`,
			`"TotalRetentionAmountNP":[{"Currency":"IDR","CurrencyID":"","Value":"4000000000"}]`},
		{"/api/treaty-in/hitung/egnpi",
			`{"aksi":"nilai","kurs":[{"Currency":"USD","Conversion":"100"}],"egnpi":[{"Currency":"USD","Amount":"2"}]}`,
			`"TotalEgnpiAmount":"200"`},
	} {
		if w := kirim(t, h, http.MethodPost, c.jalur, "", badan(c.isi)); w.Code != http.StatusUnauthorized {
			t.Errorf("%s tanpa identitas: %d", c.jalur, w.Code)
		}
		if w := kirim(t, h, http.MethodPost, c.jalur, "AKUN-UJI", badan(`{`)); w.Code != http.StatusBadRequest {
			t.Errorf("%s badan rusak: %d", c.jalur, w.Code)
		}
		w := kirim(t, h, http.MethodPost, c.jalur, "AKUN-UJI", badan(c.isi))
		if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), c.mau) {
			t.Errorf("%s: %d %s", c.jalur, w.Code, w.Body.String())
		}
	}
}

func TestRuteKelasBisnisTanpaGrupKosong(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	w := kirim(t, h, http.MethodGet, "/api/treaty-in/warisan/kelas-bisnis", "AKUN-UJI", nil)
	if w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "[]" {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
