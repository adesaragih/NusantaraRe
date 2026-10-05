package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatydescription/backend/handlers"
	"nusantarare/modul/treatydescription/backend/services"
	"nusantarare/modul/treatydescription/backend/tiruan"
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

func TestRuteTreatyDescription(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?xol=1", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"10015"`) || strings.Contains(w.Body.String(), `"id":"10001"`) {
		t.Errorf("daftar XOL %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"?q=limit&status=1", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"10001"`) || !strings.Contains(w.Body.String(), `"id":"10004"`) {
		t.Errorf("daftar cari %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/10001", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"aktif":true`) || strings.Contains(w.Body.String(), "dipakai") {
		t.Errorf("buka %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"descName":"uji portfolio","isXol":"0","statusAktif":"1"}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"10018"`) || !strings.Contains(w.Body.String(), `"descName":"UJI PORTFOLIO"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/10018", `{"descName":"uji portfolio dua","isXol":"1","statusAktif":"0"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"isXol":"1"`) || !strings.Contains(w.Body.String(), `"aktif":false`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"descName":"uji treaty limit","isXol":"0","statusAktif":"1"}`, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("nama kembar %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"descName":"uji","isXol":"0","statusAktif":"1","aktif":true}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("kolom tampilan bukan isian form: %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"descName":"uji lihat","isXol":"0","statusAktif":"1"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/10018", "", "UJI-ADMIN"); w.Code != 405 {
		t.Errorf("hapus harus tidak ada: %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/19999", "", "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}
