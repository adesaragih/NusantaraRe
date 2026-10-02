package handlers_test

// Seam HTTP aksi tambahan dari XML (paket 9): `Generate`, `Copy`, `On Retention`.

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHTTPGenerateSeeDetailCSV(t *testing.T) {
	u := server(t, true)
	res, err := http.Post(u.srv.URL+pre+"/produk/generate", "application/json", strings.NewReader(badanUji))
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", res.StatusCode)
	}
	req, _ := http.NewRequest("POST", u.srv.URL+pre+"/produk/generate", strings.NewReader(badanUji))
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "text/csv; charset=utf-8" ||
		!strings.Contains(res.Header.Get("Content-Disposition"), "SeeDetail.csv") ||
		!strings.HasPrefix(string(b), "TYPE,TYPE_CEDING,CEDING,GRUP,PRODUCTNAME,") ||
		!strings.Contains(string(b), "\r\nRider,SUPRLUS,UJI CEDING SATU,HEALTH,UJI PRODUK,") {
		t.Errorf("Generate: %d %v\n%s", res.StatusCode, res.Header, b)
	}
	if len(u.g.Umum) != 2 || u.g.Komit != 0 {
		t.Error("Generate tidak menyimpan apa pun")
	}
	if kode, badan := u.minta(t, "POST", pre+"/produk/generate", "{", true); kode != http.StatusBadRequest {
		t.Errorf("badan rusak: %d %s", kode, badan)
	}
}

func TestHTTPCopyDanOnRetention(t *testing.T) {
	u := server(t, true)
	isiMaster(u.g)
	salin := strings.Replace(badanUji, `{"umum":`, `{"salinanDari":"UJI-001","umum":`, 1)
	kode, badan := u.minta(t, "POST", pre+"/produk", salin, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"100044"`) || strings.Contains(badan, "salinanDari") {
		t.Fatalf("Copy: %d %s", kode, badan)
	}
	// Daftar yang dapat disunting datang dari badan form (klien menyalin halaman);
	// produk asal tidak tersentuh.
	if !strings.Contains(u.g.Umum["100044"], `"CREATEOP":"UJI-PELAKU"`) || !strings.Contains(u.g.Umum["UJI-001"], `"Plan":"UJI PLAN"`) {
		t.Errorf("salinan: %s", u.g.Umum["100044"])
	}
	if kode, badan := u.minta(t, "POST", pre+"/produk", strings.Replace(salin, "UJI-001", "UJI-999", 1), true); kode != http.StatusNotFound ||
		!strings.Contains(badan, "UJI-999") {
		t.Errorf("produk asal tidak ada: %d %s", kode, badan)
	}
	if kode, badan := u.minta(t, "PUT", pre+"/produk/100044", salin, true); kode != http.StatusBadRequest ||
		!strings.Contains(badan, "salinanDari") {
		t.Errorf("salinanDari pada PUT: %d %s", kode, badan)
	}
	hitung := strings.Replace(badanUji, `{"umum":`, `{"hitungOutward":true,"umum":`, 1)
	if kode, badan := u.minta(t, "PUT", pre+"/produk/100044", hitung, true); kode != http.StatusOK ||
		!strings.Contains(badan, `"outwardList":[]`) || u.g.MintaOR[0] != "01/03/2026" {
		t.Errorf("On Retention diubah: %d %s %v", kode, badan, u.g.MintaOR)
	}
}
