package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/handlers"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// akunUji - header identitas stub (`inti.PelakuDari` bila stub aktif).
const akunUji = "UJI-AKUN-1"

func routerUji(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi, nil), adaDB, true)
}

func minta(t *testing.T, h http.Handler, metode, jalur, badan string, beridentitas bool) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	if beridentitas {
		r.Header.Set("X-Pelaku", akunUji)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func gudangSatuKasus() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["EDMLF-1"] = &tiruan.Polis{ID: "EDMLF-1", OldPolicyNo: "UJI-PL-1", EdmType: "1", TglInput: "2026-10-01 09:00:00",
		Kepala: map[string]string{"TYPE": "QR"}}
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-D1", PolisID: "EDMLF-1", PLNumber: "UJI-PL-1",
		EdmStatus: models.StatusOld, Nilai: map[string]string{"CERTIFICATE_NO": "UJI-C1"}})
	return g
}

func TestRuteBacaMenjawab(t *testing.T) {
	h := routerUji(gudangSatuKasus(), true)
	w := minta(t, h, "GET", handlers.Prefix+"/inbox", "", true)
	if w.Code != http.StatusOK {
		t.Fatalf("inbox %d %s", w.Code, w.Body)
	}
	var inbox struct {
		Baris []map[string]any `json:"baris"`
		Total int              `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &inbox); err != nil || inbox.Total != 1 || inbox.Baris[0]["caseId"] != "EDMLF-1" {
		t.Fatalf("inbox = %s (%v)", w.Body, err)
	}
	w = minta(t, h, "GET", handlers.Prefix+"/kasus/EDMLF-1", "", true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"policyNo":"UJI-PL-1"`) {
		t.Fatalf("kasus %d %s", w.Code, w.Body)
	}
	w = minta(t, h, "GET", handlers.Prefix+"/kasus/EDMLF-1/peserta?halaman=1", "", true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"CERTIFICATE_NO":"UJI-C1"`) {
		t.Fatalf("peserta %d %s", w.Code, w.Body)
	}
}

func TestRuteBacaMenolak(t *testing.T) {
	h := routerUji(gudangSatuKasus(), true)
	for _, k := range []struct {
		jalur string
		ident bool
		kode  int
	}{
		{"/inbox", false, http.StatusUnauthorized},
		{"/kasus/EDMLF-1", false, http.StatusUnauthorized},
		{"/kasus/NBLF-1", true, http.StatusNotFound},
		{"/kasus/EDMLF-9/peserta", true, http.StatusNotFound},
	} {
		if w := minta(t, h, "GET", handlers.Prefix+k.jalur, "", k.ident); w.Code != k.kode || !strings.Contains(w.Body.String(), `"galat"`) {
			t.Errorf("%s ident=%v: %d %s, mau %d", k.jalur, k.ident, w.Code, w.Body, k.kode)
		}
	}
	if w := minta(t, routerUji(gudangSatuKasus(), false), "GET", handlers.Prefix+"/inbox", "", true); w.Code != http.StatusServiceUnavailable {
		t.Errorf("tanpa Oracle: %d", w.Code)
	}
}
