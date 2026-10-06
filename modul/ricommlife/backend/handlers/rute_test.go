package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/ricommlife/backend/handlers"
	"nusantarare/modul/ricommlife/backend/services"
	"nusantarare/modul/ricommlife/backend/tiruan"
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

func csvJSON(teks string) string {
	b, _ := json.Marshal(map[string]string{"csv": teks})
	return string(b)
}

// Butir 9: prefix /api/ri-comm-life, KODE menu ricommlife.
func TestPrefixDanKode(t *testing.T) {
	if handlers.Prefix != "/api/ri-comm-life" || handlers.KodeMenu != "ricommlife" {
		t.Errorf("%s %s", handlers.Prefix, handlers.KodeMenu)
	}
}

func TestRuteRingkasan(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?usedby=comm%20b&halaman=1", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"1000004"`) || !strings.Contains(w.Body.String(), `"total":1`) ||
		!strings.Contains(w.Body.String(), `"ukuran":50`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"usedby":" UJI COMM BARU "}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"1000005"`) || !strings.Contains(w.Body.String(), `"operatorId":"UJI-ADMIN"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"usedby":"uji comm baru"}`, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "already used by ID 1000005") {
		t.Errorf("kembar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/1000005", `{"usedby":"UJI COMM BARU 2"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"usedby":"UJI COMM BARU 2"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/1000003", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"jumlahKomisi":2`) {
		t.Errorf("buka %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"usedby":"uji","asing":1}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("kolom asing %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/1000003", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("hapus View only %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/1000003/detail", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 200 {
		t.Errorf("detail View only boleh dibaca %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/1000003", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"komisiTerhapus":2`) {
		t.Errorf("hapus %d %s", w.Code, w.Body.String())
	}
	for _, jalur := range []string{"/1000003", "/1000003/detail"} {
		if w := kirim(h, "GET", handlers.Prefix+jalur, "", "UJI-ADMIN"); w.Code != 404 {
			t.Errorf("%s sesudah hapus %d", jalur, w.Code)
		}
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

// R/I COMM DETAIL (InboxRIComm): tambah dan ubah baris; USEDBY dari ringkasan; isian USEDBY ditolak (400).
func TestRuteDetail(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	w := kirim(h, "POST", handlers.Prefix+"/1000003/detail", `{"contract":"2","year":"1","comm":"0,125"}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"1000044"`) || !strings.Contains(w.Body.String(), `"usedby":"UJI COMM RETRO"`) ||
		!strings.Contains(w.Body.String(), `"idUsedBy":"1000003"`) || !strings.Contains(w.Body.String(), `"comm":"0.125"`) {
		t.Fatalf("tambah %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/1000003/detail", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"total":3`) ||
		strings.Index(w.Body.String(), `"1000040"`) > strings.Index(w.Body.String(), `"1000044"`) {
		t.Errorf("detail urut ID naik %s", w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/1000003/detail", `{"contract":"1","year":"1","comm":"1"}`, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "already exists") {
		t.Errorf("kembar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/1000003/detail", `{"contract":"9","year":"9","comm":"1","usedby":"LAIN"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("nama dari isian ditolak %d", w.Code)
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/1000003/detail/1000041", `{"contract":"1","year":"3","comm":"10.5"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"year":"3"`) {
		t.Errorf("ubah %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/1000003/detail/1000042", `{"contract":"1","year":"9","comm":"1"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("baris ringkasan lain %d %s", w.Code, w.Body.String())
	}
	for _, m := range []struct{ metode, jalur string }{{"POST", "/1000003/detail"}, {"PUT", "/1000003/detail/1000040"}} {
		if w := kirim(h, m.metode, handlers.Prefix+m.jalur, `{"contract":"9","year":"9","comm":"1"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", m.metode, w.Code)
		}
	}
}

func TestRuteUnggah(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	salah := csvJSON("USEDBY;CONTRACT;YEAR;COMM\nUJI X;1;12345;1\n")
	if w := kirim(h, "POST", handlers.Prefix+"/unggah/pratinjau", salah, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"sah":false`) || !strings.Contains(w.Body.String(), "YEAR must be") {
		t.Errorf("pratinjau %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", salah, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "Row 2: YEAR must be") {
		t.Errorf("simpan galat %d %s", w.Code, w.Body.String())
	}
	sah := csvJSON("USEDBY;CONTRACT;YEAR;COMM\nUJI X;1;1;0,5\n")
	for _, jalur := range []string{"/unggah", "/unggah/pratinjau"} {
		if w := kirim(h, "POST", handlers.Prefix+jalur, sah, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", jalur, w.Code)
		}
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", sah, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"disimpan":1`) || !strings.Contains(w.Body.String(), `"ringkasanBaru":1`) {
		t.Errorf("simpan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", `{"csv":"`+strings.Repeat("A", 9<<20)+`"}`, "UJI-ADMIN"); w.Code != 413 {
		t.Errorf("terlalu besar %d", w.Code)
	}
}
