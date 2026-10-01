package handlers_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/handlers"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

func mintaBerkas(t *testing.T, h http.Handler, jalur, medan, isi string, beridentitas bool) *httptest.ResponseRecorder {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, err := mw.CreateFormFile(medan, "UJI.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = fw.Write([]byte(isi))
	_ = mw.Close()
	r := httptest.NewRequest("POST", jalur, &b)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	if beridentitas {
		r.Header.Set("X-Pelaku", akunUji)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func gudangCSV() *tiruan.Gudang {
	g := gudangSatuKasus()
	g.Peserta[0].Nilai["PLAN"] = "UJI-PLAN"
	g.Peserta[0].Nilai["POLICY_HOLDER"] = "UJI-PH"
	return g
}

func TestRuteCSV(t *testing.T) {
	g := gudangCSV()
	h := routerBuat(g)
	unggah, csv := handlers.Prefix+"/kasus/EDMLF-1/unggah", handlers.Prefix+"/kasus/EDMLF-1/csv"
	sah := "PLAN,POLICY_HOLDER,CERTIFICATE_NO\nUJI-PLAN,UJI-PH,UJI-N1\n"
	if w := mintaBerkas(t, h, unggah, "berkas", sah, false); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas %d", w.Code)
	}
	if w := mintaBerkas(t, h, unggah, "lain", sah, true); w.Code != http.StatusBadRequest {
		t.Fatalf("tanpa medan berkas %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "POST", unggah, sah, true); w.Code != http.StatusBadRequest {
		t.Fatalf("bukan multipart %d", w.Code)
	}
	w := mintaBerkas(t, h, unggah, "berkas", sah+"UJI-LAIN,UJI-PH,UJI-N2\n", true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"total":2,"ditolak":1`) ||
		!strings.Contains(w.Body.String(), `{"baris":2,"kolom":"PLAN","pesan":"Plan di CSV tidak sesuai, mohon di cek kembali"}`) {
		t.Fatalf("tinjau %d %s", w.Code, w.Body)
	}
	if w := mintaBerkas(t, h, unggah, "berkas", "PLAN,CURRENCY\nUJI-PLAN,IDR\n", true); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("tanpa POLICY_HOLDER %d %s", w.Code, w.Body)
	}
	if w := mintaBerkas(t, h, unggah, "berkas", "PLAN,POLICY_HOLDER,STATUS\nUJI-PLAN,UJI-PH,1\n", true); w.Code != http.StatusOK ||
		!strings.Contains(w.Body.String(), `"diabaikan":["STATUS"]`) {
		t.Fatalf("kolom asing dilaporkan %d %s", w.Code, w.Body)
	}
	w = mintaBerkas(t, h, csv, "berkas", sah+"UJI-LAIN,UJI-PH,UJI-N2\n", true)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), `"galat":"1 of 2 CSV rows are rejected; no row is saved"`) ||
		!strings.Contains(w.Body.String(), `"kolom":"PLAN"`) {
		t.Fatalf("tambah ditolak %d %s", w.Code, w.Body)
	}
	w = mintaBerkas(t, h, csv, "berkas", sah, true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"disimpan":1`) {
		t.Fatalf("tambah %d %s", w.Code, w.Body)
	}
	if n := len(g.Peserta); n != 2 || g.Peserta[1].EdmStatus != models.StatusNew {
		t.Fatalf("peserta %d", n)
	}
	if w := mintaBerkas(t, h, csv, "berkas", sah, true); w.Code != http.StatusConflict {
		t.Fatalf("tambah kedua %d %s", w.Code, w.Body)
	}
}
