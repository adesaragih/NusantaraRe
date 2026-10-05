package rute_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/menu"
	"nusantarare/inti/backend/templat"
	"nusantarare/inti/backend/templat/rute"
)

func mux(t *testing.T, stub bool) *http.ServeMux {
	t.Helper()
	k, err := templat.NewKatalog([]templat.Slot{{
		Kode: "uji.premi.fire", Menu: "uji", Grup: "UJI Menu", Nama: "PREMIUM · FIRE", Ekstensi: ".csv", Pemisah: ';',
		JumlahKolom: 2, NamaUnduhan: "UJI - PREMIUM · FIRE.csv", Bawaan: []byte("A;B\n"),
	}})
	if err != nil {
		t.Fatal(err)
	}
	m := http.NewServeMux()
	rute.Pasang(m, templat.NewLayanan(k, templat.NewGudangMemori()), stub)
	return m
}

// sesi - permintaan dari akun login yang memegang `menuAkun`.
func sesi(r *http.Request, akun string, menuAkun ...string) *http.Request {
	ctx := inti.DenganAksesMenu(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: akun}), menuAkun)
	return r.WithContext(ctx)
}

func kirim(m http.Handler, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	return w
}

func unggahan(t *testing.T, kode, isi, catatan string) *http.Request {
	t.Helper()
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	f, _ := mw.CreateFormFile("berkas", "baru.csv")
	_, _ = f.Write([]byte(isi))
	_ = mw.WriteField("catatan", catatan)
	_ = mw.Close()
	r := httptest.NewRequest("POST", rute.Prefix+"/"+kode+"/unggah", &b)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	return r
}

func TestAksesKelolaHanyaPemegangTemplateManager(t *testing.T) {
	m := mux(t, false)
	if w := kirim(m, httptest.NewRequest("GET", rute.Prefix, nil)); w.Code != 401 {
		t.Errorf("tanpa sesi %d", w.Code)
	}
	if w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix, nil), "UJI-1", "uji")); w.Code != 403 {
		t.Errorf("pemegang menu pemilik saja tidak boleh mengelola: %d", w.Code)
	}
	w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix, nil), "UJI-IT", menu.KodeTemplateManager))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"kode":"uji.premi.fire"`) || !strings.Contains(w.Body.String(), `"aktif":null`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
}

func TestUnduhVersiAktifUntukPemegangMenuPemilik(t *testing.T) {
	m := mux(t, false)
	w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix+"/uji.premi.fire/unduh", nil), "UJI-1", "uji"))
	if w.Code != 200 || w.Body.String() != "A;B\n" || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/csv") {
		t.Fatalf("unduh %d %q", w.Code, w.Body.String())
	}
	cd := w.Header().Get("Content-Disposition")
	if !strings.Contains(cd, `filename="UJI - PREMIUM _ FIRE.csv"`) || !strings.Contains(cd, "filename*=UTF-8''UJI%20-%20PREMIUM%20%C2%B7%20FIRE.csv") {
		t.Errorf("content-disposition %q", cd)
	}
	if w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix+"/uji.premi.fire/unduh?versi=0", nil), "UJI-1", "uji")); w.Code != 403 {
		t.Errorf("versi tertentu hanya pengelola: %d", w.Code)
	}
	if w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix+"/uji.premi.fire/unduh", nil), "UJI-2", "lain")); w.Code != 403 {
		t.Errorf("menu lain ditolak: %d", w.Code)
	}
	if w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix+"/uji.tidak/unduh", nil), "UJI-IT", menu.KodeTemplateManager)); w.Code != 404 {
		t.Errorf("slot tak ada: %d", w.Code)
	}
}

func TestUnggahPeriksaAktifkan(t *testing.T) {
	m := mux(t, false)
	it := func(r *http.Request) *http.Request { return sesi(r, "UJI-IT", menu.KodeTemplateManager) }
	if w := kirim(m, it(unggahan(t, "uji.premi.fire", "A;B;C\n", ""))); w.Code != 422 || !strings.Contains(w.Body.String(), "needs 2 columns") {
		t.Errorf("kolom salah %d %s", w.Code, w.Body.String())
	}
	if w := kirim(m, it(unggahan(t, "uji.premi.fire", "A;X\nUJI;1\n", "judul B diganti"))); w.Code != 200 || !strings.Contains(w.Body.String(), `"versi":1`) {
		t.Fatalf("unggah %d %s", w.Code, w.Body.String())
	}
	w := kirim(m, it(httptest.NewRequest("GET", rute.Prefix+"/uji.premi.fire/riwayat", nil)))
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"diunggahOleh":"UJI-IT"`) || !strings.Contains(w.Body.String(), `"catatan":"judul B diganti"`) {
		t.Errorf("riwayat %d %s", w.Code, w.Body.String())
	}
	if w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix+"/uji.premi.fire/unduh", nil), "UJI-1", "uji")); w.Body.String() != "A;X\nUJI;1\n" {
		t.Errorf("unduh versi aktif %q", w.Body.String())
	}
	r := it(httptest.NewRequest("POST", rute.Prefix+"/uji.premi.fire/aktifkan", strings.NewReader(`{"versi":0}`)))
	if w := kirim(m, r); w.Code != 200 {
		t.Errorf("aktifkan bawaan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(m, sesi(httptest.NewRequest("GET", rute.Prefix+"/uji.premi.fire/unduh", nil), "UJI-1", "uji")); w.Body.String() != "A;B\n" {
		t.Errorf("kembali ke bawaan %q", w.Body.String())
	}
	r = it(httptest.NewRequest("POST", rute.Prefix+"/uji.premi.fire/aktifkan", strings.NewReader(`{"versi":7}`)))
	if w := kirim(m, r); w.Code != 404 {
		t.Errorf("versi tak ada %d", w.Code)
	}
	r = it(httptest.NewRequest("POST", rute.Prefix+"/uji.premi.fire/aktifkan", strings.NewReader(`{}`)))
	if w := kirim(m, r); w.Code != 400 {
		t.Errorf("badan tanpa versi %d", w.Code)
	}
}

func TestStubTanpaIdentitasTidakBolehUnggah(t *testing.T) {
	m := mux(t, true)
	if w := kirim(m, unggahan(t, "uji.premi.fire", "A;B\n", "")); w.Code != 401 {
		t.Errorf("stub tanpa X-Pelaku %d", w.Code)
	}
	r := unggahan(t, "uji.premi.fire", "A;B\n", "")
	r.Header.Set("X-Pelaku", "UJI-DEV")
	if w := kirim(m, r); w.Code != 200 {
		t.Errorf("stub dengan X-Pelaku %d %s", w.Code, w.Body.String())
	}
}
