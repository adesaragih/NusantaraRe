package handlers_test

// Seam HTTP Copy Old (permintaan work owner 03-10-2026) di atas gudang tiruan.

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

func TestHTTPCopyOld(t *testing.T) {
	u := server(t, true)
	u.g.IsiLama("100901", `{"ID":"100901","PRODUCTNAME":"UJI LAMA","CEDING":"UJI CEDING","TREATYNUMBER":"UJI-1"}`, "")
	u.g.IsiLama("100902", `{"ID":"100902","PRODUCTNAME":"UJI DITOLAK"}`, "")
	u.g.TolakLama("100902", "TYPE terisi tetapi tidak punya kolom flat")

	kode, badan := u.minta(t, "GET", pre+"/produk-lama", "", true)
	var d struct {
		Daftar []models.ProdukLama `json:"daftar"`
		Total  int                 `json:"total"`
	}
	if kode != http.StatusOK || json.Unmarshal([]byte(badan), &d) != nil || d.Total != 2 || d.Daftar[0].ID != "100901" ||
		!d.Daftar[0].BolehDisalin || d.Daftar[1].BolehDisalin || len(d.Daftar[1].Alasan) != 1 {
		t.Fatalf("daftar produk lama: %d %s", kode, badan)
	}

	kode, badan = u.minta(t, "POST", pre+"/produk-lama/salin", `{"ids":["100901","100902"]}`, true)
	var j models.JawabanSalinLama
	if kode != http.StatusOK || json.Unmarshal([]byte(badan), &j) != nil || j.Disalin != 1 || len(j.Hasil) != 2 ||
		j.Hasil[0].Status != models.SalinDisalin || j.Hasil[1].Status != models.SalinDitolak {
		t.Fatalf("Process Copy: %d %s", kode, badan)
	}
	if _, ada := u.g.Produk["100901"]; !ada {
		t.Error("produk lama tersalin ke tabel flat")
	}
	// Sesudah disalin: tampil di grid produk, hilang dari daftar Copy Old.
	if kode, badan := u.minta(t, "GET", pre+"/produk", "", true); kode != http.StatusOK || !strings.Contains(badan, `"id":"100901"`) {
		t.Errorf("grid produk memuat produk tersalin: %d %s", kode, badan)
	}
	if kode, badan := u.minta(t, "GET", pre+"/produk-lama", "", true); kode != http.StatusOK || strings.Contains(badan, `"id":"100901"`) {
		t.Errorf("produk tersalin tidak lagi di daftar Copy Old: %d %s", kode, badan)
	}
}

func TestHTTPCopyOldDitolak(t *testing.T) {
	u := server(t, true)
	for _, k := range []struct {
		badan string
		kode  int
	}{
		{`{"ids":[]}`, http.StatusUnprocessableEntity},
		{`bukan json`, http.StatusBadRequest},
	} {
		if kode, badan := u.minta(t, "POST", pre+"/produk-lama/salin", k.badan, true); kode != k.kode {
			t.Errorf("%s: %d %s, mau %d", k.badan, kode, badan, k.kode)
		}
	}
	if kode, _ := u.minta(t, "GET", pre+"/produk-lama", "", false); kode == http.StatusOK {
		t.Error("tanpa pelaku ditolak")
	}
}
