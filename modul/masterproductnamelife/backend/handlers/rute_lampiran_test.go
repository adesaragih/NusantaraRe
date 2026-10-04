package handlers_test

// Seam HTTP lampiran (paket 8) di atas gudang tiruan + penyimpanan stub di folder sementara.

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/handlers"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

func serverLampiran(t *testing.T) (*httptest.Server, *tiruan.Gudang) {
	t.Helper()
	g := tiruan.Baru()
	g.IsiJSON("100007", `{"ID":"100007"}`, "")
	g.AppName = "UJI-APP"
	l := services.BaruLayanan(g, g.Transaksi, nil).DenganPenyimpanan(services.PenyimpananLokal(t.TempDir()))
	srv := httptest.NewServer(handlers.RouterDengan(l, true, true))
	t.Cleanup(srv.Close)
	return srv, g
}

func kirimBerkasUji(t *testing.T, alamat, nama, isi string) (int, string) {
	t.Helper()
	var buf bytes.Buffer
	m := multipart.NewWriter(&buf)
	if nama != "" {
		w, _ := m.CreateFormFile("berkas", nama)
		_, _ = w.Write([]byte(isi))
	}
	_ = m.Close()
	req, _ := http.NewRequest("POST", alamat, &buf)
	req.Header.Set("Content-Type", m.FormDataContentType())
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

func ambil(t *testing.T, metode, alamat string) (int, []byte, http.Header) {
	t.Helper()
	req, _ := http.NewRequest(metode, alamat, nil)
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, b, res.Header
}

func TestHTTPLampiranUnggahUnduhHapus(t *testing.T) {
	srv, _ := serverLampiran(t)
	dasar := srv.URL + pre + "/produk/100007/lampiran"
	kode, badan := kirimBerkasUji(t, dasar, "UJI nota.pdf", "ISI PDF")
	if kode != http.StatusOK || !strings.Contains(badan, `"status":"terunggah"`) || !strings.Contains(badan, `"fileMimeType":"pdf"`) {
		t.Fatalf("unggah: %d %s", kode, badan)
	}
	var a struct{ ID string }
	_ = json.Unmarshal([]byte(badan), &a)
	kode, isi, kepala := ambil(t, "GET", dasar+"/"+a.ID+"/unduh")
	if kode != http.StatusOK || string(isi) != "ISI PDF" || kepala.Get("Content-Type") != "application/pdf" ||
		!strings.Contains(kepala.Get("Content-Disposition"), "UJI%20nota.pdf") {
		t.Errorf("unduh: %d %q %v", kode, isi, kepala)
	}
	kode, isi, _ = ambil(t, "GET", dasar+"/unduh-semua")
	if z, err := zip.NewReader(bytes.NewReader(isi), int64(len(isi))); kode != http.StatusOK || err != nil || len(z.File) != 1 {
		t.Errorf("unduh semua: %d %v", kode, err)
	}
	if kode, b, _ := ambil(t, "GET", dasar+"/"+a.ID+"/office"); kode != http.StatusUnprocessableEntity {
		t.Errorf("office untuk pdf: %d %s", kode, b)
	}
	if kode, b, _ := ambil(t, "DELETE", dasar+"/"+a.ID); kode != http.StatusOK {
		t.Errorf("hapus: %d %s", kode, b)
	}
	if kode, b, _ := ambil(t, "GET", dasar); kode != http.StatusOK || !strings.Contains(string(b), `"total":0`) {
		t.Errorf("daftar kosong sesudah hapus: %d %s", kode, b)
	}
}

func TestHTTPLampiranGalat(t *testing.T) {
	srv, g := serverLampiran(t)
	dasar := srv.URL + pre + "/produk/100007/lampiran"
	if kode, badan := kirimBerkasUji(t, dasar, "", ""); kode != http.StatusUnprocessableEntity ||
		!strings.Contains(badan, "Tidak ada file yg diattach") {
		t.Errorf("tanpa berkas: %d %s", kode, badan)
	}
	if kode, badan := kirimBerkasUji(t, dasar, "UJI.xyzzy", "x"); kode != http.StatusUnprocessableEntity ||
		!strings.Contains(badan, "not supported") {
		t.Errorf("jenis tak didukung: %d %s", kode, badan)
	}
	_, _ = kirimBerkasUji(t, dasar, "UJI.xlsx", "x")
	if kode, badan := kirimBerkasUji(t, dasar, "UJI.XLSX", "y"); kode != http.StatusConflict {
		t.Errorf("nama sama: %d %s", kode, badan)
	}
	var id string
	for k := range g.Lampiran {
		id = k
	}
	if kode, b, _ := ambil(t, "GET", dasar+"/"+id+"/office"); kode != http.StatusServiceUnavailable ||
		!strings.Contains(string(b), "View Office Online needs the file in storage") {
		t.Errorf("View Office Online atas berkas stub lokal: %d %s", kode, b)
	}
	if kode, _, _ := ambil(t, "GET", srv.URL+pre+"/produk/100999/lampiran"); kode != http.StatusNotFound {
		t.Errorf("produk tidak ada: %d", kode)
	}
	if kode, _, _ := ambil(t, "GET", dasar+"/TIDAKADA/unduh"); kode != http.StatusNotFound {
		t.Errorf("lampiran tidak ada: %d", kode)
	}
}
