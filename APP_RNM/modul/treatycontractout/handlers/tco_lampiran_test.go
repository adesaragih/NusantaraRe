package handlers

// Uji pintu HTTP lampiran - TANPA Oracle (tiket 12). Lingkaran penuhnya
// terhadap skema uji ada di tco_db_test.go.

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/unggah"
	"nusantarare/modul/treatycontractout/services"
)

var ruteLampiranTCO = []struct{ metode, jalur string }{
	{http.MethodGet, "/api/treaty-contract-out/kategori-lampiran"},
	{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/lampiran"},
	{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/lampiran"},
	{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/lampiran/semua"},
	{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/lampiran/selaras"},
	{http.MethodGet, "/api/treaty-contract-out/tahun/1000001/lampiran/1000000001/isi"},
	{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/lampiran/1000000001/ulangi"},
	{http.MethodDelete, "/api/treaty-contract-out/tahun/1000001/lampiran/1000000001"},
}

// Kedelapan rute sampai ke handler-nya lewat Router yang sesungguhnya - bukan
// 404/405 - dan tanpa database menjawab 503. Router yang menyusun pola
// bertabrakan akan panik di sini.
func TestRuteLampiranTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range ruteLampiranTCO {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(r.metode, r.jalur, nil)
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s: kode %d, mau 503 (handler lampiran tanpa database)", r.metode, r.jalur, w.Code)
		}
	}
	isi, err := os.ReadFile("rute_treaty_contract_out.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(isi), "daftarkanRuteLampiranTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute lampiran harus didaftarkan tepat sekali dari pendaftar modul")
	}
}

func TestJawabGalatLampiranTCO(t *testing.T) {
	kasus := []struct {
		err  error
		kode int
	}{
		{services.ErrKategoriLampiranKosong, http.StatusServiceUnavailable},
		{unggah.ErrUnggahanDirBelumDisetel, http.StatusServiceUnavailable},
		{services.ErrLampiranTidakAda, http.StatusNotFound},
		{services.ErrLampiranBelumTerkirim, http.StatusConflict},
		{services.ErrLampiranSudahTerkirim, http.StatusConflict},
		{services.ErrLampiranTanpaBerkas, http.StatusConflict},
		{services.ErrBerkasSumberLampiranHilang, http.StatusConflict},
		{unggah.ErrBerkasTerlaluBesar, http.StatusRequestEntityTooLarge},
		{unggah.ErrBerkasKosong, http.StatusBadRequest},
		{services.ErrKategoriLampiranTidakDikenal, http.StatusUnprocessableEntity},
	}
	for _, k := range kasus {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, k.err)
		if w.Code != k.kode {
			t.Errorf("%v: kode %d, mau %d", k.err, w.Code, k.kode)
		}
	}
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, errors.Join(unggah.ErrBerkasKosong))
	if !strings.Contains(w.Body.String(), services.PesanTanpaBerkasTCO) {
		t.Errorf("berkas kosong tanpa pesan b376: %s", w.Body.String())
	}
}

// Unduhan: `attachment`, `nosniff`, dan nama non-ASCII tidak memecah header.
func TestKepalaUnduhanAman(t *testing.T) {
	w := httptest.NewRecorder()
	kepalaUnduhan(w, "", "kontrak 2026 é.pdf")
	d := w.Header().Get("Content-Disposition")
	if !strings.HasPrefix(d, "attachment;") || !strings.Contains(d, "filename*=utf-8''") {
		t.Errorf("Content-Disposition = %q", d)
	}
	if w.Header().Get("X-Content-Type-Options") != "nosniff" ||
		w.Header().Get("Content-Type") != "application/octet-stream" {
		t.Errorf("header: %v", w.Header())
	}
}
