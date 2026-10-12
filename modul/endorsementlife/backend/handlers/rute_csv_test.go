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
	g.Peserta[0].Nilai["CURRENCY"] = "IDR"
	return g
}

// Judul + baris CSV yang memenuhi kolom wajib QR (keputusan work owner 10-10-2026): kolom awal PLAN,POLICY_HOLDER,
// CERTIFICATE_NO divariasikan uji, sisanya `akhirCSV`.
const judulCSVUji = "PLAN,POLICY_HOLDER,CERTIFICATE_NO,NAME_OF_INSURED,SEX,DOB,CURRENCY,BEGIN_DATE,EXPIRED_DATE,SUM_INSURED,AGE," +
	"GROSS_PREMIUM,DEDUCTION,RI_ADMIN_FEE,BROKERAGE_FEE,TAX,PROF_COMM,CLAIM,NET_PREMIUM"

const akhirCSV = ",UJI-NAMA,M,02/01/1990,IDR,01/01/2026,31/12/2026,1000,30,1,0,0,0,0,0,0,0"

func TestRuteCSV(t *testing.T) {
	g := gudangCSV()
	h := routerBuat(g)
	unggah, csv := handlers.Prefix+"/kasus/EDML-1/unggah", handlers.Prefix+"/kasus/EDML-1/csv"
	sah := judulCSVUji + "\nUJI-PLAN,UJI-PH,UJI-N1" + akhirCSV + "\n"
	if w := mintaBerkas(t, h, unggah, "berkas", sah, false); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas %d", w.Code)
	}
	if w := mintaBerkas(t, h, unggah, "lain", sah, true); w.Code != http.StatusBadRequest {
		t.Fatalf("tanpa medan berkas %d %s", w.Code, w.Body)
	}
	if w := minta(t, h, "POST", unggah, sah, true); w.Code != http.StatusBadRequest {
		t.Fatalf("bukan multipart %d", w.Code)
	}
	// Baris 3 menolak karena CERTIFICATE_NO ganda dalam berkas (PLAN/POLICY_HOLDER beda tidak lagi ditolak, 10-10-2026).
	ganda := sah + "UJI-LAIN,UJI-PH-LAIN,UJI-N1" + akhirCSV + "\n"
	w := mintaBerkas(t, h, unggah, "berkas", ganda, true)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"total":2,"ditolak":1`) ||
		!strings.Contains(w.Body.String(), `{"baris":3,"kolom":"CERTIFICATE_NO"`) || strings.Contains(w.Body.String(), `"kolom":"PLAN"`) {
		t.Fatalf("tinjau %d %s", w.Code, w.Body)
	}
	if w := mintaBerkas(t, h, unggah, "berkas", "PLAN,CURRENCY\nUJI-PLAN,IDR\n", true); w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("judul tanpa kolom wajib %d %s", w.Code, w.Body)
	}
	if w := mintaBerkas(t, h, unggah, "berkas", judulCSVUji+",STATUS\nUJI-PLAN,UJI-PH,UJI-N3"+akhirCSV+",1\n", true); w.Code != http.StatusOK ||
		!strings.Contains(w.Body.String(), `"diabaikan":["STATUS"]`) {
		t.Fatalf("kolom asing dilaporkan %d %s", w.Code, w.Body)
	}
	w = mintaBerkas(t, h, csv, "berkas", ganda, true)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), `"galat":"1 of 2 rows have a problem; no row was saved."`) ||
		!strings.Contains(w.Body.String(), `"kolom":"CERTIFICATE_NO"`) {
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

// TestRuteCSVSesudahProses - Process Endorsement langkah final (keputusan work
// owner 08-10-2026): sesudah Save, unggah (periksa) dan tambah CSV 409 berpesan,
// peserta tidak bertambah.
func TestRuteCSVSesudahProses(t *testing.T) {
	g := gudangCSV()
	h := routerBuat(g)
	if w := minta(t, h, "POST", handlers.Prefix+"/kasus/EDML-1/simpan", `{"pilih":["UJI-D1"]}`, true); w.Code != http.StatusOK {
		t.Fatalf("simpan %d %s", w.Code, w.Body)
	}
	sah := judulCSVUji + "\nUJI-PLAN,UJI-PH,UJI-N1" + akhirCSV + "\n"
	for _, jalur := range []string{handlers.Prefix + "/kasus/EDML-1/unggah", handlers.Prefix + "/kasus/EDML-1/csv"} {
		w := mintaBerkas(t, h, jalur, "berkas", sah, true)
		if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), models.PesanSudahDiproses) {
			t.Fatalf("%s: %d %s", jalur, w.Code, w.Body)
		}
	}
	if len(g.Peserta) != 1 {
		t.Errorf("peserta bertambah: %d", len(g.Peserta))
	}
}

// TestRuteTemplatCSV - GET .../templat-csv (10-10-2026): berkas CSV unduhan; gerbang sama dengan unggahan.
func TestRuteTemplatCSV(t *testing.T) {
	g := gudangCSV()
	h := routerBuat(g)
	jalur := handlers.Prefix + "/kasus/EDML-1/templat-csv"
	if w := minta(t, h, "GET", jalur, "", false); w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa identitas %d", w.Code)
	}
	w := minta(t, h, "GET", jalur, "", true)
	if w.Code != http.StatusOK || w.Header().Get("Content-Type") != "text/csv; charset=utf-8" ||
		w.Header().Get("Content-Disposition") != `attachment; filename="Template_EDM_EDML-1.csv"` ||
		!strings.HasPrefix(w.Body.String(), "NAME_OF_INSURED,CERTIFICATE_NO,SEX,DOB,CURRENCY,") {
		t.Fatalf("templat %d %v %.120s", w.Code, w.Header(), w.Body)
	}
	tutup := gudangCSV()
	tutup.Polis["EDML-1"].Status = models.StatusKasusSelesai
	if w := minta(t, routerBuat(tutup), "GET", jalur, "", true); w.Code != http.StatusConflict {
		t.Fatalf("kasus tertutup %d %s", w.Code, w.Body)
	}
}

// TestRuteCSVGalatJudulBersih - galat judul CSV sampai ke layar TANPA awalan teknis, status tetap 422 (10-10-2026).
func TestRuteCSVGalatJudulBersih(t *testing.T) {
	h := routerBuat(gudangCSV())
	w := mintaBerkas(t, h, handlers.Prefix+"/kasus/EDML-1/unggah", "berkas", "PLAN,CURRENCY\nUJI-PLAN,IDR\n", true)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), `"galat":"Your file is missing these columns: NAME_OF_INSURED, CERTIFICATE_NO, SEX`) ||
		strings.Contains(w.Body.String(), "invalid endorsement input") || strings.Contains(w.Body.String(), "invalid CSV header") {
		t.Fatalf("galat judul %d %s", w.Code, w.Body)
	}
	w = mintaBerkas(t, h, handlers.Prefix+"/kasus/EDML-1/unggah", "berkas", judulCSVUji+"\nUJI-PLAN,\"UJI-PH,UJI-N1"+akhirCSV+"\n", true)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "The file cannot be read as a CSV file (problem in row 2)") {
		t.Fatalf("CSV rusak %d %s", w.Code, w.Body)
	}
}
