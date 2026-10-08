package handlers_test

// Uji rute rumus tab Prop yang membaca penampung halaman — `POST /hitung/akumulasi`
// (Accumulation ← Reporting Period) dan `POST /hitung/share-prop` (Share ← Limits).

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

func TestRuteHitungTabProp(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)
	for _, c := range []struct{ jalur, isi, mau string }{
		{"/api/treaty-in/hitung/akumulasi",
			`{"aksi":"periode","AccumulationPeriod":"half","ReportingStart":"20250101","ReportingEnd":"20251231","AccumulationList":[]}`,
			`{"Period":"H 2","ReportDate":"20250701","SubDays":"","SubDueDate":""}`},
		{"/api/treaty-in/hitung/share-prop",
			`{"aksi":"share","RNMShareP":"10","OptionLimit":"1","Limits":[{"TreatyType":"SURPLUS","Detail":[{"TreatyGroup":"FIRE","SpreadingTypeID":"X","CessionList":[{"Currency":"IDR","Value":"1000"}]}]}]}`,
			`"TotalShareRnmProp":[{"Currency":"IDR","CurrencyID":"","Value":"100"}]`},
	} {
		if w := kirim(t, h, http.MethodPost, c.jalur, "", badan(c.isi)); w.Code != http.StatusUnauthorized {
			t.Errorf("%s tanpa identitas: %d", c.jalur, w.Code)
		}
		if w := kirim(t, h, http.MethodPost, c.jalur, "AKUN-UJI", badan(`{`)); w.Code != http.StatusBadRequest {
			t.Errorf("%s badan rusak: %d", c.jalur, w.Code)
		}
		if w := kirim(t, h, http.MethodPost, c.jalur, "AKUN-UJI", badan(`{"aksi":"hapus"}`)); w.Code != http.StatusUnprocessableEntity {
			t.Errorf("%s aksi tak dikenal: %d", c.jalur, w.Code)
		}
		w := kirim(t, h, http.MethodPost, c.jalur, "AKUN-UJI", badan(c.isi))
		if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), c.mau) {
			t.Errorf("%s: %d %s — mau %s", c.jalur, w.Code, w.Body.String(), c.mau)
		}
	}
}
