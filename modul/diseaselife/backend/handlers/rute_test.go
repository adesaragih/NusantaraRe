package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/diseaselife/backend/handlers"
	"nusantarare/modul/diseaselife/backend/models"
	"nusantarare/modul/diseaselife/backend/services"
	"nusantarare/modul/diseaselife/backend/tiruan"
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
	if handlers.Prefix != "/api/disease-life" || handlers.KodeMenu != "diseaselife" {
		t.Errorf("%s %s", handlers.Prefix, handlers.KodeMenu)
	}
}

// GET daftar (View only boleh) dengan saring / urut / halaman sampai ke server; Save Add 201 / Edit 200; 422 (wajib,
// kembar, panjang); 403; 404; 400; nol DELETE; 503.
func TestRute(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?halaman=x", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"ukuran":10`) || !strings.Contains(w.Body.String(), `{"id":"100005","icdCode":"UJI05","disease":"UJI SHIGELOSIS"}`) ||
		!strings.Contains(w.Body.String(), `"halaman":1`) || !strings.Contains(w.Body.String(), `"total":5`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"?icd=uji0&disease=demam&urut=ICD&arah=asc&halaman=1", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"total":2`) ||
		g.Diminta != (models.Saringan{ICDCode: "uji0", Disease: "demam", Urut: models.UrutICD, Naik: true, Halaman: 1}) {
		t.Errorf("saring %d %s %+v", w.Code, w.Body.String(), g.Diminta)
	}
	w := kirim(h, "POST", handlers.Prefix, `{"icdCode":" uji99 ","disease":" uji Demam "}`, "UJI-ADMIN")
	if w.Code != 201 || !strings.Contains(w.Body.String(), `{"id":"100006","icdCode":"UJI99","disease":"UJI DEMAM"}`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/100002", `{"icdCode":"UJI02","disease":"UJI DEMAM TIFOID BERAT"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"disease":"UJI DEMAM TIFOID BERAT"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	for badan, kata := range map[string]string{
		`{"icdCode":"","disease":"UJI"}`:                                    "ICD Code is required",
		`{"icdCode":"UJI88","disease":""}`:                                  "Disease is required",
		`{"icdCode":"uji01","disease":"UJI"}`:                               "already used by ID 100001",
		`{"icdCode":"UJI88","disease":"` + strings.Repeat("A", 1001) + `"}`: "longer than 1000",
	} {
		if w := kirim(h, "POST", handlers.Prefix, badan, "UJI-ADMIN"); w.Code != 422 || !strings.Contains(w.Body.String(), kata) {
			t.Errorf("%.40s: %d %s", badan, w.Code, w.Body.String())
		}
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/999", `{"icdCode":"UJI87","disease":"UJI Q"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	for _, m := range []struct{ metode, jalur string }{{"POST", ""}, {"PUT", "/100002"}} {
		if w := kirim(h, m.metode, handlers.Prefix+m.jalur, `{"icdCode":"UJI86","disease":"UJI L"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", m.metode, w.Code)
		}
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"icdCode":"UJI85","disease":"UJI","id":"100009"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("ID dari isian ditolak (disabled b1070: server yang mengisi) %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/100002", "", "UJI-ADMIN"); w.Code == 200 {
		t.Errorf("DELETE tidak boleh ada: %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

// D1.1: NEXTVAL menghasilkan ID yang sudah ada -> 409 berkalimat.
func TestRuteIDSequenceSudahAda(t *testing.T) {
	g := tiruan.Contoh()
	g.Seq = 100001
	w := kirim(router(g, true), "POST", handlers.Prefix, `{"icdCode":"UJI84","disease":"UJI BARU"}`, "UJI-ADMIN")
	if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "new ID 100001 from SEQ_DISEASE_LIFE already exists") {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
