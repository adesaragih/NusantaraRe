package handlers_test

// Regresi temuan /code-review 01-10-2026 di seam HTTP.

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
	"nusantarare/inti/backend/unggah"
)

// #8 - berkas melampaui batas dijawab 413 berkalimat, bukan "Tidak ada file yg diattach".
func TestHTTPLampiranTerlaluBesar413(t *testing.T) {
	srv, _ := serverLampiran(t)
	var buf bytes.Buffer
	m := multipart.NewWriter(&buf)
	w, _ := m.CreateFormFile("berkas", "UJI besar.pdf")
	_, _ = w.Write(bytes.Repeat([]byte("x"), int(unggah.BatasUkuranUnggahan+galat.BatasFormulir)+1024))
	_ = m.Close()
	req, _ := http.NewRequest("POST", srv.URL+pre+"/produk/100007/lampiran", &buf)
	req.Header.Set("Content-Type", m.FormDataContentType())
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if res.StatusCode != http.StatusRequestEntityTooLarge || !strings.Contains(string(b), "exceeds") {
		t.Errorf("berkas terlalu besar: %d %s", res.StatusCode, b)
	}
}

// #15 - badan JSON dibatasi.
func TestHTTPBadanJSONTerlaluBesar413(t *testing.T) {
	u := server(t, true)
	badan := `{"umum":{"productName":"` + strings.Repeat("x", 5<<20) + `"}}`
	if kode, b := u.minta(t, "POST", pre+"/produk", badan, true); kode != http.StatusRequestEntityTooLarge {
		t.Errorf("badan > 4 MiB: %d %.200s", kode, b)
	}
}

// #10 - `batas` autocomplete.
func TestHTTPMasterBatas(t *testing.T) {
	u := server(t, true)
	isiMaster(u.g)
	if kode, b := u.minta(t, "GET", pre+"/master/ceding?cari=uji&batas=1", "", true); kode != http.StatusOK || !strings.Contains(b, `"total":1`) {
		t.Errorf("batas 1: %d %s", kode, b)
	}
	if kode, b := u.minta(t, "GET", pre+"/master/ceding?batas=x", "", true); kode != http.StatusBadRequest {
		t.Errorf("batas bukan angka: %d %s", kode, b)
	}
}
