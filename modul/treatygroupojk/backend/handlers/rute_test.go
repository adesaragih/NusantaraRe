package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatygroupojk/backend/handlers"
	"nusantarare/modul/treatygroupojk/backend/services"
	"nusantarare/modul/treatygroupojk/backend/tiruan"
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

func TestRuteTreatyGroupOJK(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?q=rekayasa", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"09"`) || strings.Contains(w.Body.String(), `"id":"01"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"name":"uji cargo","nameIdn":"uji pengangkutan"}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10"`) || !strings.Contains(w.Body.String(), `"name":"UJI CARGO"`) ||
		!strings.Contains(w.Body.String(), `"orderNo":"12"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/10", `{"name":"uji cargo","nameIdn":"uji angkut"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"nameIdn":"UJI ANGKUT"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji property","nameIdn":"x"}`, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("nama kembar %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji","asing":1}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("kolom asing %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji","nameIdn":"x","orderNo":"5"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("Order No bukan isian form: %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"name":"uji lihat","nameIdn":"x"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/10", "", "UJI-ADMIN"); w.Code != 405 {
		t.Errorf("hapus harus tidak ada: %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/77", "", "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
