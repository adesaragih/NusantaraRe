package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/benefitlife/backend/handlers"
	"nusantarare/modul/benefitlife/backend/services"
	"nusantarare/modul/benefitlife/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi{G: g}.Jalankan), adaDB, false)
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

func TestPrefixDanKode(t *testing.T) {
	if handlers.Prefix != "/api/benefit-life" || handlers.KodeMenu != "benefitlife" {
		t.Errorf("%s %s", handlers.Prefix, handlers.KodeMenu)
	}
}

// GET daftar (10, ID menurun), Save Add 201 / Edit 200, wajib 422, View only 403, tak ada 404, kolom asing 400, 503.
func TestRute(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	w := kirim(h, "GET", handlers.Prefix, "", "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"ukuran":10`) || !strings.Contains(w.Body.String(), `"total":3`) ||
		strings.Index(w.Body.String(), `"100004"`) > strings.Index(w.Body.String(), `"100001"`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"?arah=asc&benefit=uji", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 ||
		strings.Index(w.Body.String(), `"100001"`) > strings.Index(w.Body.String(), `"100004"`) {
		t.Errorf("View only boleh membaca, naik %d %s", w.Code, w.Body.String())
	}
	w = kirim(h, "POST", handlers.Prefix, `{"benefit":" uji kritis "}`, "UJI-ADMIN")
	if w.Code != 201 || !strings.Contains(w.Body.String(), `"id":"100012"`) || !strings.Contains(w.Body.String(), `"benefit":"UJI KRITIS"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/100001", `{"benefit":"uji rawat jalan"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"benefit":"UJI RAWAT JALAN"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"benefit":"  "}`, "UJI-ADMIN"); w.Code != 422 || !strings.Contains(w.Body.String(), "Benefit is required") {
		t.Errorf("wajib %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/999", `{"benefit":"UJI"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	for _, m := range []struct{ metode, jalur string }{{"POST", ""}, {"PUT", "/100001"}} {
		if w := kirim(h, m.metode, handlers.Prefix+m.jalur, `{"benefit":"UJI LIHAT"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", m.metode, w.Code)
		}
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"benefit":"UJI","id":"100099"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("ID dari isian ditolak (disabled b1013) %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/100001", "", "UJI-ADMIN"); w.Code == 200 {
		t.Errorf("DELETE tidak boleh ada (XML tanpa Delete): %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

// K3: NEXTVAL menghasilkan ID yang sudah ada -> 409 berkalimat (bukan 500 mentah).
func TestRuteIDSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 2 // -> 100002, sudah ada
	w := kirim(router(g, true), "POST", handlers.Prefix, `{"benefit":"UJI"}`, "UJI-ADMIN")
	if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "new ID 100002 from M_BENEFIT_LIFE_SEQ already exists") {
		t.Errorf("mau 409 berkalimat, dapat %d %s", w.Code, w.Body.String())
	}
}
