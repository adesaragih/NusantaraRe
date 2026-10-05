package handlers_test

// Rute HTTP Accounts di atas gudang tiruan: status, badan, dan identitas pelaku.

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/accounts/backend/handlers"
	"nusantarare/modul/accounts/backend/models"
	"nusantarare/modul/accounts/backend/services"
	"nusantarare/modul/accounts/backend/tiruan"
)

func server(t *testing.T, adaDB bool) (*httptest.Server, *tiruan.Gudang) {
	t.Helper()
	g := tiruan.Contoh()
	srv := httptest.NewServer(handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi), adaDB, true))
	t.Cleanup(srv.Close)
	return srv, g
}

func minta(t *testing.T, metode, url, badan string, pelaku bool) (int, string) {
	t.Helper()
	var isi io.Reader
	if badan != "" {
		isi = bytes.NewBufferString(badan)
	}
	req, _ := http.NewRequest(metode, url, isi)
	if pelaku {
		req.Header.Set("X-Pelaku", "UJI-ADMIN")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

func TestRuteBacaDanTulis(t *testing.T) {
	srv, g := server(t, true)
	dasar := srv.URL + handlers.Prefix
	kode, b := minta(t, "GET", dasar+"?cari=uji+pt&halaman=1&ukuran=1", "", false)
	var h services.HalamanDaftar
	if err := json.Unmarshal([]byte(b), &h); kode != http.StatusOK || err != nil || h.Total != 2 || len(h.Daftar) != 1 {
		t.Errorf("daftar: %d %s", kode, b)
	}
	if kode, b := minta(t, "GET", dasar+"/pilihan", "", false); kode != http.StatusOK || !strings.Contains(b, `"note":"UJI GROUP A"`) {
		t.Errorf("pilihan: %d %s", kode, b)
	}
	if kode, b := minta(t, "GET", dasar+"/organisasi?cari=dua", "", false); kode != http.StatusOK || !strings.Contains(b, `"idView":"ORG-102"`) {
		t.Errorf("organisasi: %d %s", kode, b)
	}
	badan := `{"insuredId":"` + models.AwalanOrg + `ORG-102","groupBusinessId":"10001","description":"UJI"}`
	kode, b = minta(t, "POST", dasar, badan, true)
	if kode != http.StatusOK || !strings.Contains(b, `"idView":"ACC-4916562"`) || !strings.Contains(b, `"createOp":"UJI-ADMIN"`) {
		t.Errorf("tambah: %d %s", kode, b)
	}
	if kode, _ := minta(t, "POST", dasar, badan, true); kode != http.StatusConflict {
		t.Errorf("tambah ganda = 409: %d", kode)
	}
	if kode, b := minta(t, "POST", dasar, `{"insuredId":"","groupBusinessId":""}`, true); kode != http.StatusUnprocessableEntity ||
		!strings.Contains(b, "Insured Name: Value cannot be blank") {
		t.Errorf("wajib isi = 422: %d %s", kode, b)
	}
	if kode, _ := minta(t, "POST", dasar, `{"insuredId":"x","id":"ACC-9"}`, true); kode != http.StatusBadRequest {
		t.Errorf("medan asing (id dari klien) = 400: %d", kode)
	}
	if kode, _ := minta(t, "POST", dasar, badan, false); kode != http.StatusUnauthorized {
		t.Errorf("tanpa pelaku = 401: %d", kode)
	}
	// Akun hanya diinput sekali: nol rute ubah, baca-satu, atau hapus.
	ubah := `{"insuredId":"` + models.AwalanOrg + `ORG-101","groupBusinessId":"10002","description":"UJI ubah"}`
	for _, k := range []struct{ metode, jalur string }{{"PUT", "/ACC-4916560"}, {"GET", "/ACC-4916560"}, {"DELETE", "/ACC-4916560"}} {
		if kode, _ := minta(t, k.metode, dasar+k.jalur, ubah, true); kode != http.StatusNotFound {
			t.Errorf("%s %s = 404: %d", k.metode, k.jalur, kode)
		}
	}
	for _, metode := range []string{"PUT", "PATCH", "DELETE"} {
		if kode, _ := minta(t, metode, dasar, ubah, true); kode != http.StatusMethodNotAllowed {
			t.Errorf("%s %s = 405: %d", metode, handlers.Prefix, kode)
		}
	}
	if len(g.Disisip) != 1 {
		t.Errorf("tulisan: %d sisip", len(g.Disisip))
	}
}

func TestRuteTanpaDatabase503(t *testing.T) {
	srv, _ := server(t, false)
	if kode, _ := minta(t, "GET", srv.URL+handlers.Prefix, "", false); kode != http.StatusServiceUnavailable {
		t.Errorf("tanpa database: %d", kode)
	}
}
