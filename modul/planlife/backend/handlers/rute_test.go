package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/planlife/backend/handlers"
	"nusantarare/modul/planlife/backend/services"
	"nusantarare/modul/planlife/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan), adaDB, false)
}

func kirim(h http.Handler, metode, jalur, badan, akun string, lihat ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	ctx := inti.DenganMenuLihat(inti.DenganAksesMenu(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: akun}),
		[]string{handlers.KodeMenu}), lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func TestPrefixDanKode(t *testing.T) {
	if handlers.Prefix != "/api/plan-life" || handlers.KodeMenu != "planlife" {
		t.Errorf("%s %s", handlers.Prefix, handlers.KodeMenu)
	}
}

// GET daftar / pilihan (View only boleh), Save Add 201 / Edit 200, 422 (wajib, bukan master, kembar), 403, 404, 400,
// nol DELETE, 503.
func TestRute(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix, "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"ukuran":10`) || !strings.Contains(w.Body.String(), `"businessId":"9001"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan/business", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"oldId":"L1"`) || strings.Contains(w.Body.String(), "KEBAKARAN") {
		t.Errorf("pilihan business %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan/benefit", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), "UJI RAWAT INAP") {
		t.Errorf("pilihan benefit %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"coverName":"UJI PLAN C","business":"uji kredit","benefit":"uji meninggal"}`, "UJI-ADMIN")
	if w.Code != 201 || !strings.Contains(w.Body.String(), `"id":"100044"`) || !strings.Contains(w.Body.String(), `"business":"UJI KREDIT","businessId":"9001"`) ||
		!strings.Contains(w.Body.String(), `"benefitId":"100002"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/100001", `{"coverName":"UJI PLAN A2","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"coverName":"UJI PLAN A2"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	for badan, kata := range map[string]string{
		`{"coverName":"","business":"","benefit":""}`:                                   "Plan Name is required",
		`{"coverName":"UJI Z","business":"UJI TIDAK ADA","benefit":"UJI RAWAT INAP"}`:   "is not in the Business list",
		`{"coverName":"uji plan b","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}`: "already used by ID 100002",
	} {
		if w := kirim(h, "POST", handlers.Prefix, badan, "UJI-ADMIN"); w.Code != 422 || !strings.Contains(w.Body.String(), kata) {
			t.Errorf("%s: %d %s", badan, w.Code, w.Body.String())
		}
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/999", `{"coverName":"UJI Q","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	for _, m := range []struct{ metode, jalur string }{{"POST", ""}, {"PUT", "/100001"}} {
		if w := kirim(h, m.metode, handlers.Prefix+m.jalur, `{"coverName":"UJI L","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", m.metode, w.Code)
		}
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"coverName":"UJI","business":"UJI KREDIT","benefit":"UJI RAWAT INAP","businessId":"1"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("ID master dari isian ditolak (K4: server yang mengisi) %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/100001", "", "UJI-ADMIN"); w.Code == 200 {
		t.Errorf("DELETE tidak boleh ada: %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

// K3: NEXTVAL menghasilkan ID yang sudah ada -> 409 berkalimat.
func TestRuteIDSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 1
	w := kirim(router(g, true), "POST", handlers.Prefix, `{"coverName":"UJI BARU","business":"UJI KREDIT","benefit":"UJI RAWAT INAP"}`, "UJI-ADMIN")
	if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "new ID 100001 from M_PRODUCT_TYPE_LIFE_SEQ already exists") {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
