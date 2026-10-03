package handlers_test

// Rute HTTP Marketing Officer di atas gudang tiruan: bentuk jawaban dan kode status.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/marketingofficer/backend/handlers"
	"nusantarare/modul/marketingofficer/backend/services"
	"nusantarare/modul/marketingofficer/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi), adaDB, true)
}

func kirim(t *testing.T, h http.Handler, metode, jalur, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", "UJI-ADMIN")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestRuteDaftarDanPilihan(t *testing.T) {
	h := router(tiruan.Contoh(), true)
	w := kirim(t, h, "GET", handlers.Prefix, "")
	var d struct {
		Daftar []map[string]any `json:"daftar"`
		Total  int              `json:"total"`
	}
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &d) != nil || d.Total != 3 {
		t.Fatalf("daftar %d %s", w.Code, w.Body.String())
	}
	for _, kunci := range []string{"id", "clientId", "clientName", "clientId2", "moLeader", "moStatus", "aksesLogin", "statusAkun"} {
		if _, ada := d.Daftar[0][kunci]; !ada {
			t.Errorf("baris daftar tanpa %q: %v", kunci, d.Daftar[0])
		}
	}
	w = kirim(t, h, "GET", handlers.Prefix+"/pilihan", "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"leader":[{"id":"10000101"`) {
		t.Errorf("pilihan %d %s", w.Code, w.Body.String())
	}
}

func TestRuteTambahUbahDanGalat(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	badan := `{"aksesLogin":"UJI-MKT01","leader":false,"leaderId":"10000101","branchParent":"UJI-P00","branchDetailId":"UJI-B01","aktif":true}`
	w := kirim(t, h, "POST", handlers.Prefix, badan)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10000201"`) || !strings.Contains(w.Body.String(), `"clientId":"CON-1002"`) {
		t.Fatalf("tambah %d %s", w.Code, w.Body.String())
	}
	for _, k := range []struct {
		metode, jalur, badan string
		kode                 int
		pesan                string
	}{
		{"POST", handlers.Prefix, badan, http.StatusConflict, "already has an active marketing officer row"},
		{"POST", handlers.Prefix, `{"aksesLogin":""}`, http.StatusUnprocessableEntity, "Login Account is required"},
		{"POST", handlers.Prefix, `{"tidakDikenal":1}`, http.StatusBadRequest, "not valid JSON"},
		{"PUT", handlers.Prefix + "/19999999", badan, http.StatusNotFound, "Marketing officer not found"},
	} {
		w := kirim(t, h, k.metode, k.jalur, k.badan)
		if w.Code != k.kode || !strings.Contains(w.Body.String(), k.pesan) {
			t.Errorf("%s %s: %d %s, mau %d %q", k.metode, k.jalur, w.Code, w.Body.String(), k.kode, k.pesan)
		}
		if strings.Contains(w.Body.String(), "invalid input:") || strings.Contains(w.Body.String(), "already active:") {
			t.Errorf("pesan layar memuat nama galat penanda: %s", w.Body.String())
		}
	}
	w = kirim(t, h, "PUT", handlers.Prefix+"/10000201", strings.Replace(badan, `"aktif":true`, `"aktif":false`, 1))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"moStatus":"2"`) {
		t.Errorf("nonaktifkan %d %s", w.Code, w.Body.String())
	}
}

func TestRuteTanpaDatabase503(t *testing.T) {
	w := kirim(t, router(tiruan.Contoh(), false), "GET", handlers.Prefix, "")
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("tanpa Oracle: %d", w.Code)
	}
}

func TestRuteLog(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	badan := `{"aksesLogin":"UJI-MKT02","leader":false,"leaderId":"10000101","branchParent":"UJI-P00","branchDetailId":"UJI-B02","aktif":true}`
	if w := kirim(t, h, "PUT", handlers.Prefix+"/10000103", badan); w.Code != 200 {
		t.Fatalf("ubah %d %s", w.Code, w.Body.String())
	}
	w := kirim(t, h, "GET", handlers.Prefix+"/10000103/log", "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"kolom":"AKSES_LOGIN","sebelum":"UJIOPERATORLAMA","sesudah":"UJI-MKT02"`) {
		t.Errorf("log %d %s", w.Code, w.Body.String())
	}
	if w := kirim(t, h, "GET", handlers.Prefix+"/19999999/log", ""); w.Code != http.StatusNotFound {
		t.Errorf("log MO tidak ada: %d", w.Code)
	}
}
