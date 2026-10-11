package handlers

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// mintaMultipartCsv - POST multipart dengan bagian bernama `medan` (isi CSV) dan satu bagian lain di depannya.
func mintaMultipartCsv(t *testing.T, svc *services.Service, medan, isi string) (int, string) {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	_ = mw.WriteField("lain", "x")
	f, err := mw.CreateFormFile(medan, "uji.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte(isi))
	_ = mw.Close()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest("POST", "/api/nbfacin/kasus/UJI-NB-1/objek/csv/berkas", &b)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r.Header.Set("X-Pelaku", "UJI-USER")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

// TestUnggahCsvKabel - tiket 75 lewat HTTP: bukan multipart / tanpa medan berkas 400; bagian `berkas` sampai ke layanan (tanpa
// basis data 503); rute status / aktif / baris / simpan / ulangi terdaftar; halaman bukan angka 400.
func TestUnggahCsvKabel(t *testing.T) {
	svc := services.Baru(nil)
	if kode, _ := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/objek/csv/berkas", "a;b", "UJI-USER"); kode != 400 {
		t.Errorf("bukan multipart: %d", kode)
	}
	if kode, isi := mintaMultipartCsv(t, svc, "lain2", "x"); kode != 400 || !strings.Contains(isi, "medan berkas") {
		t.Errorf("tanpa berkas: %d %s", kode, isi)
	}
	if kode, _ := mintaMultipartCsv(t, svc, "berkas", strings.Join(services.KolomCsvObjek, ";")+"\n"); kode != 503 {
		t.Errorf("berkas tanpa basis data: %d", kode)
	}
	for _, k := range []struct{ metode, jalur string }{{"GET", "/objek/csv/aktif"}, {"GET", "/objek/csv/7"},
		{"POST", "/objek/csv/7/simpan"}, {"POST", "/objek/csv/7/ulangi"}} {
		if kode, isi := minta(t, svc, k.metode, "/api/nbfacin/kasus/UJI-NB-1"+k.jalur, "", "UJI-USER"); kode != 503 {
			t.Errorf("%s %s: %d %s", k.metode, k.jalur, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek/csv/7/baris?halaman=x&ukuran=10", "", ""); kode != 400 {
		t.Errorf("halaman: %d", kode)
	}
}

// TestStatusUnggahKabel - galat / pesan berawalan nomor baris (kecuali yang sudah berawalan "baris N"), dibuat DD-MM-YYYY HH:mm.
func TestStatusUnggahKabel(t *testing.T) {
	k := keStatusUnggahKabel(services.StatusUnggahCsv{Unggah: models.UnggahCsv{ID: "7", Status: "valid",
		Dibuat: time.Date(2026, 10, 10, 13, 5, 0, 0, time.UTC)}, Galat: []models.PesanUnggahCsv{{NoBaris: 2, Teks: "baris 2 kolom LOL: x"},
		{NoBaris: 3, Teks: "mata uang"}, {Teks: "berkas"}}, AdaObjek: true})
	if !reflect.DeepEqual(k.Galat, []string{"baris 2 kolom LOL: x", "baris 3: mata uang", "berkas"}) || k.Dibuat != "10-10-2026 13:05" ||
		!k.AdaObjek || k.Pesan == nil {
		t.Errorf("%+v", k)
	}
}

// TestGalatCsvBukanFire - bug DEV 11-10-2026: CSV template FIRE untuk case bukan FIRE -> 409 dengan pesan apa adanya.
func TestGalatCsvBukanFire(t *testing.T) {
	w := httptest.NewRecorder()
	tulisGalat(w, services.ErrCsvBukanFire)
	if w.Code != 409 || !strings.Contains(w.Body.String(), "Upload CSV template FIRE hanya untuk Class of Business FIRE (Pega "+
		"ChooseActivityMappingCSV_Act). CSV untuk Class of Business ini belum tersedia di sistem baru.") {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}

// TestGalatCsvTemplateAneka - tiket 73 P6 (K73-60): case Aneka bertemplate Aneka umum / CPM -> 409 dengan pesan apa adanya.
func TestGalatCsvTemplateAneka(t *testing.T) {
	w := httptest.NewRecorder()
	tulisGalat(w, services.ErrCsvTemplateAnekaBelumAda)
	if w.Code != 409 || !strings.Contains(w.Body.String(), "Template CSV untuk Class of Business ini (Pega TemplateCSVAneka / "+
		"TemplateCSVAneka_CPM) belum tersedia di sistem baru - unggah CSV baru tersedia untuk Marine Hull.") {
		t.Errorf("%d %s", w.Code, w.Body.String())
	}
}
