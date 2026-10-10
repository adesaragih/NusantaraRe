package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/riratelife/backend/handlers"
	"nusantarare/modul/riratelife/backend/services"
	"nusantarare/modul/riratelife/backend/tiruan"
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

func TestRuteRingkasan(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"?usedby=rate%20b&halaman=1", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"102"`) || !strings.Contains(w.Body.String(), `"total":1`) ||
		!strings.Contains(w.Body.String(), `"ukuran":50`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"?urut=usedby&arah=asc", "", "UJI-ADMIN"); w.Code != 200 ||
		strings.Index(w.Body.String(), `"id":"101"`) > strings.Index(w.Body.String(), `"id":"105"`) {
		t.Errorf("urut %s", w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix, `{"usedby":" UJI RATE BARU "}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"106"`) || !strings.Contains(w.Body.String(), `"operatorId":"UJI-ADMIN"`) {
		t.Fatalf("add %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"usedby":"uji rate baru"}`, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "already used by ID 106") {
		t.Errorf("kembar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/106", `{"usedby":"UJI RATE BARU 2"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"usedby":"UJI RATE BARU 2"`) {
		t.Errorf("edit %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/101", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"jumlahRate":3`) {
		t.Errorf("buka %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/101/rate?halaman=1", "", "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"rate":"1,25"`) || !strings.Contains(w.Body.String(), `"total":3`) {
		t.Errorf("detail %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix, `{"usedby":"uji","asing":1}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("kolom asing %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/101", "", "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("hapus View only %d", w.Code)
	}
	if w := kirim(h, "DELETE", handlers.Prefix+"/101", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"rateTerhapus":3`) {
		t.Errorf("hapus %d %s", w.Code, w.Body.String())
	}
	for _, jalur := range []string{"/101", "/101/rate"} {
		if w := kirim(h, "GET", handlers.Prefix+jalur, "", "UJI-ADMIN"); w.Code != 404 {
			t.Errorf("%s sesudah hapus %d", jalur, w.Code)
		}
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-ADMIN"); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

func TestRuteUnggah(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	salah := csvJSON("USEDBY;CONTRACT;GENDER;AGE;RATE\nUJI X;1;Q;1;1\n")
	if w := kirim(h, "POST", handlers.Prefix+"/unggah/pratinjau", salah, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"sah":false`) || !strings.Contains(w.Body.String(), "GENDER must be U, M, or F") {
		t.Errorf("pratinjau %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", salah, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "Row 2: GENDER must be U, M, or F") {
		t.Errorf("simpan galat %d %s", w.Code, w.Body.String())
	}
	sah := csvJSON("USEDBY;CONTRACT;GENDER;AGE;RATE\nUJI X;1;U;1;0,5\n")
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", sah, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("View only %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah/pratinjau", sah, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
		t.Errorf("pratinjau View only %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", sah, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"disimpan":1`) || !strings.Contains(w.Body.String(), `"ringkasanBaru":1`) {
		t.Errorf("simpan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah", `{"csv":"`+strings.Repeat("A", 9<<20)+`"}`, "UJI-ADMIN"); w.Code != 413 {
		t.Errorf("terlalu besar %d", w.Code)
	}
}

// Rate Detail (View Detail.xml `InboxRIRate`): tambah dan ubah baris; R/I RATE NAME dari ringkasan yang dilihat.
func TestRuteRateDetail(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	w := kirim(h, "POST", handlers.Prefix+"/101/rate", `{"gender":"f","contract":"2","age":"","rate":"0.125"}`, "UJI-ADMIN")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"id":"9001"`) || !strings.Contains(w.Body.String(), `"usedby":"UJI RATE A"`) ||
		!strings.Contains(w.Body.String(), `"idUsedBy":"101"`) || !strings.Contains(w.Body.String(), `"rate":"0,125"`) {
		t.Fatalf("tambah %d %s", w.Code, w.Body.String())
	}
	if r := g.Ringkasan["101"]; r.OperatorID != "UJI-ADMIN" || r.UsedBy != "UJI RATE A" {
		t.Errorf("ringkasan sesudah tambah %+v", r)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/101/rate", "", "UJI-ADMIN"); w.Code != 200 || !strings.Contains(w.Body.String(), `"ukuran":20`) ||
		!strings.Contains(w.Body.String(), `"total":4`) || strings.Index(w.Body.String(), `"9001"`) > strings.Index(w.Body.String(), `"9000"`) {
		t.Errorf("detail urut ID turun %s", w.Body.String())
	}
	// Kunci (GENDER, AGE, CONTRACT) kembar di ringkasan yang sama - 9000 = U/30/1.
	if w := kirim(h, "POST", handlers.Prefix+"/101/rate", `{"gender":"U","contract":"01","age":"30","rate":"1"}`, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "already exists in R/I RATE NAME UJI RATE A (rate ID 9000)") {
		t.Errorf("kembar %d %s", w.Code, w.Body.String())
	}
	// Kunci sama di ringkasan lain boleh (9100 milik 102).
	if w := kirim(h, "POST", handlers.Prefix+"/101/rate", `{"gender":"F","contract":"2","age":"40","rate":"1"}`, "UJI-ADMIN"); w.Code != 200 {
		t.Errorf("kunci ringkasan lain %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/101/rate", `{"gender":"U","contract":"","age":"1","rate":"1"}`, "UJI-ADMIN"); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "CONTRACT is required") {
		t.Errorf("contract wajib %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/101/rate", `{"contract":"1","rate":"1","usedby":"LAIN"}`, "UJI-ADMIN"); w.Code != 400 {
		t.Errorf("nama dari isian ditolak %d", w.Code)
	}
	// Ubah: 9003 (M/30/1) -> M/31/1, rate berkoma; nama tetap dari ringkasan.
	if w := kirim(h, "PUT", handlers.Prefix+"/101/rate/9003", `{"gender":"M","contract":"1","age":"31","rate":"0,8"}`, "UJI-ADMIN"); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"age":"31"`) {
		t.Errorf("ubah %d %s", w.Code, w.Body.String())
	}
	if r := g.Rate["9003"]; r.Age != "31" || r.Rate != "0,8" || r.UsedBy != "UJI RATE A" {
		t.Errorf("rate sesudah ubah %+v", r)
	}
	// Ubah tanpa mengganti kunci sendiri bukan kembar.
	if w := kirim(h, "PUT", handlers.Prefix+"/101/rate/9003", `{"gender":"M","contract":"1","age":"31","rate":"0,9"}`, "UJI-ADMIN"); w.Code != 200 {
		t.Errorf("ubah kunci sendiri %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "PUT", handlers.Prefix+"/101/rate/9003", `{"gender":"U","contract":"1","age":"30","rate":"1"}`, "UJI-ADMIN"); w.Code != 422 {
		t.Errorf("ubah jadi kembar %d %s", w.Code, w.Body.String())
	}
	// 9100 milik ringkasan 102, bukan 101.
	if w := kirim(h, "PUT", handlers.Prefix+"/101/rate/9100", `{"contract":"1","rate":"1"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("rate ringkasan lain %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/999/rate", `{"contract":"1","rate":"1"}`, "UJI-ADMIN"); w.Code != 404 {
		t.Errorf("ringkasan tidak ada %d", w.Code)
	}
	for _, m := range []struct{ metode, jalur string }{{"POST", "/101/rate"}, {"PUT", "/101/rate/9000"}} {
		if w := kirim(h, m.metode, handlers.Prefix+m.jalur, `{"contract":"9","rate":"1"}`, "UJI-LIHAT", handlers.KodeMenu); w.Code != 403 {
			t.Errorf("%s View only %d", m.metode, w.Code)
		}
	}
}
