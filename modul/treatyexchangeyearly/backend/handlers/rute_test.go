package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyexchangeyearly/backend/handlers"
	"nusantarare/modul/treatyexchangeyearly/backend/services"
	"nusantarare/modul/treatyexchangeyearly/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi), adaDB, false)
}

// kirim - permintaan dari sesi login `akun`; `lihat` = menu ber-hak View only.
func kirim(h http.Handler, metode, jalur, badan, akun string, lihat ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	ctx := inti.DenganMenuLihat(inti.DenganAksesMenu(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: akun}),
		[]string{handlers.KodeMenu}), lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func TestRuteTreatyExchangeYearly(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?tahun=2024", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"kunci":"AAA4"`) || strings.Contains(w.Body.String(), `"treatyYear":"2025"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"mataUang":[`) || !strings.Contains(w.Body.String(), `"tahun":["2025","2024","2019"]`) {
		t.Errorf("pilihan %d %s", w.Code, w.Body.String())
	}
	baru := `{"treatyYear":"2026","idCurrency":"10025","startDate":"2026-07-01","endDate":"2027-06-30","toIdr":"19000","toUsd":"","quarter":"0"}`
	w := kirim(h, "POST", handlers.Prefix, baru, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10115"`) || !strings.Contains(w.Body.String(), `"startDate":"20260701T000000.000 GMT"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, baru, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("kembar %d", w.Code)
	}
	ubah := `{"kunci":"AAA3","treatyYear":"2025","idCurrency":"10025","startDate":"2025-07-01","endDate":"2026-06-30","toIdr":"18950","toUsd":"","quarter":"0"}`
	if w := kirim(h, "PUT", handlers.Prefix, ubah, "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"toIdr":"18950"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix, baru, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("edit tanpa kunci %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, ubah, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("add berkunci %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, baru, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 405 {
		t.Errorf("hapus harus tidak ada: %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
