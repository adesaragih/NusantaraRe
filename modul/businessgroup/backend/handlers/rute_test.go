package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/businessgroup/backend/handlers"
	"nusantarare/modul/businessgroup/backend/services"
	"nusantarare/modul/businessgroup/backend/tiruan"
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

func TestRuteBusinessGroup(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?treatyGroup=10007", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"10013"`) || strings.Contains(w.Body.String(), "SYARIAH") ||
		strings.Contains(w.Body.String(), `"id":"10028"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"treatyGroup":[`) {
		t.Errorf("pilihan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/10020", "", "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("SYARIAH harus 404: %d", w.Code)
	}
	w := kirim(h, "POST", handlers.Prefix, `{"topId":"10009","name":"uji hull","alias":""}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10015"`) || !strings.Contains(w.Body.String(), `"treatyName":"UJI MARINE CARGO"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/10015", `{"topId":"10009","name":"uji hull","alias":"uji kapal"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"alias":"UJI KAPAL"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"topId":"10007","name":"uji motor syariah"}`, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("SYARIAH baru %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"topId":"10007","name":"uji lihat"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/10015", "", "UJI-ADMIN"); w.Code != 405 {
		t.Errorf("hapus harus tidak ada: %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
