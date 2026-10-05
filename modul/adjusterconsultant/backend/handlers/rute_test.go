package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/adjusterconsultant/backend/handlers"
	"nusantarare/modul/adjusterconsultant/backend/services"
	"nusantarare/modul/adjusterconsultant/backend/tiruan"
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

func TestRuteAdjusterConsultant(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?status=active&q=uji", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"10001"`) || strings.Contains(w.Body.String(), `"id":"10003"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"name":"uji baru","address":"","telpNo":""}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10004"`) || !strings.Contains(w.Body.String(), `"name":"UJI BARU"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/10004", `{"name":"uji baru dua"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"name":"UJI BARU DUA"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji adjuster dua"}`, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("nama ganda %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix+"/10004/aktif", `{"active":false}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"active":false`) {
		t.Errorf("nonaktif %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/10004/aktif", `{}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("tanpa active %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji lihat"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/TIDAK-ADA", "", "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
