package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/reinsurancetype/backend/handlers"
	"nusantarare/modul/reinsurancetype/backend/services"
	"nusantarare/modul/reinsurancetype/backend/tiruan"
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

func TestRuteReinsuranceType(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?type=2&flag=active", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"10262"`) || strings.Contains(w.Body.String(), `"id":"10196"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	baru := `{"name":"uji xl 6th layer","type":"4","soaName":"","code":"","flag":"active","noUrut":"","groupType":""}`
	w := kirim(h, "POST", handlers.Prefix, baru, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10263"`) || !strings.Contains(w.Body.String(), `"code":"00"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, baru, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("nama kembar %d", w.Code)
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/10263", `{"name":"uji xl 6th layer","type":"4","code":"46","flag":"inactive"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"flag":"inactive"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/10263", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"code":"46"`) {
		t.Errorf("buka %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji","asing":1}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("kolom asing %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, baru, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/10263", "", "UJI-ADMIN"); w.Code != 405 {
		t.Errorf("hapus harus tidak ada: %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/19999", "", "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
