package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatygroup/backend/handlers"
	"nusantarare/modul/treatygroup/backend/services"
	"nusantarare/modul/treatygroup/backend/tiruan"
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

func TestRuteTreatyGroup(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?ojk=09", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"10009"`) || strings.Contains(w.Body.String(), `"id":"10007"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"ojk":[`) {
		t.Errorf("pilihan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/10007", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"anak":[`) || strings.Contains(w.Body.String(), "SYARIAH") {
		t.Errorf("buka %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"ojkId":"01","name":"uji livestock","soaName":""}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10010"`) || !strings.Contains(w.Body.String(), `"coaId":"10013"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/10010", `{"ojkId":"01","name":"uji livestock dua"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"name":"UJI LIVESTOCK DUA"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"ojkId":"01","name":"uji agriculture"}`, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("nama kembar %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"ojkId":"01","name":"uji","coaId":"1"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("COAID bukan isian form: %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"ojkId":"01","name":"uji lihat"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/10010", "", "UJI-ADMIN"); w.Code != 405 {
		t.Errorf("hapus harus tidak ada: %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/19999", "", "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
