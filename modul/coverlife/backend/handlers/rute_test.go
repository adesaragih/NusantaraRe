package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/coverlife/backend/handlers"
	"nusantarare/modul/coverlife/backend/services"
	"nusantarare/modul/coverlife/backend/tiruan"
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
	if handlers.Prefix != "/api/cover-life" || handlers.KodeMenu != "coverlife" {
		t.Errorf("%s %s", handlers.Prefix, handlers.KodeMenu)
	}
}

// GET daftar (View only boleh), Save Add 201 / Edit 200, 422 (wajib, kembar, panjang), 403, 404, 400, nol DELETE, 503.
func TestRute(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?halaman=x", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"ukuran":50`) || !strings.Contains(w.Body.String(), `{"id":"100001","cover":"UJI KECELAKAAN DIRI","note":""}`) ||
		!strings.Contains(w.Body.String(), `"halaman":1`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"cover":" uji Penyakit Kritis ","note":"uji catatan"}`, "UJI-ADMIN")
	if w.Code != 201 || !strings.Contains(w.Body.String(), `{"id":"100005","cover":"uji Penyakit Kritis","note":"uji catatan"}`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/100002", `{"cover":"UJI JIWA BERJANGKA","note":""}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"cover":"UJI JIWA BERJANGKA"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	for badan, kata := range map[string]string{
		`{"cover":""}`:                                                 "Cover is required",
		`{"cover":"uji kesehatan"}`:                                    "already used by ID 100003",
		`{"cover":"` + strings.Repeat("A", 201) + `"}`:                 "longer than 200",
		`{"cover":"UJI N","note":"` + strings.Repeat("A", 1001) + `"}`: "Note is longer than 1000",
	} {
		if w := kirim(h, "POST", handlers.Prefix, badan, "UJI-ADMIN"); w.Code != 422 || !strings.Contains(w.Body.String(), kata) {
			t.Errorf("%.40s: %d %s", badan, w.Code, w.Body.String())
		}
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/999", `{"cover":"UJI Q"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	for _, m := range []struct{ metode, jalur string }{{"POST", ""}, {"PUT", "/100002"}} {
		if w := kirim(h, m.metode, handlers.Prefix+m.jalur, `{"cover":"UJI L"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", m.metode, w.Code)
		}
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"cover":"UJI","id":"100009"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("ID dari isian ditolak (form tanpa medan ID: server yang mengisi) %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/100002", "", "UJI-ADMIN"); w.Code == 200 {
		t.Errorf("DELETE tidak boleh ada: %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

// C2: NEXTVAL menghasilkan ID yang sudah ada -> 409 berkalimat.
func TestRuteIDSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 1
	w := kirim(router(g, true), "POST", handlers.Prefix, `{"cover":"UJI BARU"}`, "UJI-ADMIN")
	if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "new ID 100001 from M_COVER_LIFE_SEQ already exists") {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
