package handlers_test

// Uji jalur HTTP tombol tulis — Save dan Submit/Actions/Decline offer.

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func (_ gudangTiruan) BacaKepalaTreatyIn(_ context.Context, id string) (map[string]any, bool, error) {
	if id == "9999999" {
		return map[string]any{}, false, nil
	}
	return map[string]any{"Ceding": "ASURANSI A", "LeadingReinsSource": "DIRECT"}, true, nil
}

func (_ gudangTiruan) BacaDokumenPendaratan(_ context.Context, _ string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (_ gudangTiruan) SimpanKontrak(_ context.Context, r models.RencanaSimpan) (string, error) {
	if r.ID == "" {
		return "1001857", nil
	}
	return r.ID, nil
}

func (_ gudangTiruan) KunciBelumTerpasang(context.Context, map[string]any) ([]string, error) {
	return nil, nil
}

func (_ gudangTiruan) BacaKepalaPenyesuaian(_ context.Context, id string) (map[string]any, bool, error) {
	if id != "1001001/R01" {
		return map[string]any{}, false, nil
	}
	return map[string]any{"OLDID": "1001001", "EDMState": "1"}, true, nil
}

func (_ gudangTiruan) SimpanPenyesuaian(context.Context, models.RencanaPenyesuaian) error { return nil }

func (_ gudangTiruan) HapusPenyesuaian(context.Context, string) error { return nil }

func (_ gudangTiruan) PemegangPosisi(_ context.Context, _ string) ([]string, error) {
	return []string{"SEC1"}, nil
}

func tombol(t *testing.T, jalur, peran, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/api/treaty-in/kontrak/"+jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", "ADESAMUEL")
	r.Header.Set("X-Peran", peran)
	w := httptest.NewRecorder()
	handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true).ServeHTTP(w, r)
	return w
}

func TestRuteSaveKontrakBaru(t *testing.T) {
	w := tombol(t, "simpan", "ReasTreatyInAdmin", `{"idKontrak":"","dokumen":{"TreatyContractName":"BARU"}}`)
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d %s", w.Code, w.Body.String())
	}
	var h services.HasilSimpan
	if err := json.Unmarshal(w.Body.Bytes(), &h); err != nil || h.ID != "1001857" || h.Pesan != "Data Sudah Disimpan Dengan ID : 1001857" {
		t.Errorf("hasil %+v %v", h, err)
	}
}

func TestRuteSubmitDanGalatnya(t *testing.T) {
	if w := tombol(t, "kirim", "ReasTreatyInAdmin", `{"idKontrak":"1001001","aksi":"submit"}`); w.Code != http.StatusOK {
		t.Errorf("submit: kode %d %s", w.Code, w.Body.String())
	}
	// Pesan Activity apa adanya, 422.
	w := tombol(t, "kirim", "ReasTreatyInAdmin", `{"idKontrak":"1001001","aksi":"submit","dokumen":{"Ceding":""}}`)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "Please input Ceding") {
		t.Errorf("Ceding kosong: kode %d %s", w.Code, w.Body.String())
	}
	// Bukan pemegang posisi: 403 berpesan.
	if w := tombol(t, "kirim", "ReasTreatyInSecHead", `{"idKontrak":"1001001","aksi":"submit"}`); w.Code != http.StatusForbidden {
		t.Errorf("bukan pemegang: kode %d", w.Code)
	}
	if w := tombol(t, "kirim", "ReasTreatyInAdmin", `{"idKontrak":"9999999","aksi":"submit"}`); w.Code != http.StatusNotFound {
		t.Errorf("kontrak tak ada: kode %d", w.Code)
	}
	if w := tombol(t, "simpan", "ReasTreatyInAdmin", `bukan json`); w.Code != http.StatusBadRequest {
		t.Errorf("bukan JSON: kode %d", w.Code)
	}
}

func TestRuteRevision(t *testing.T) {
	// Gudang tiruan rute: kontrak belum tuntas — ditolak 422, bukan 404.
	if w := tombol(t, "revisi", "ReasTreatyInAdmin", `{"idKontrak":"1001001"}`); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("belum tuntas: kode %d %s", w.Code, w.Body.String())
	}
	if w := tombol(t, "revisi", "ReasTreatyInSecHead", `{"idKontrak":"1001001"}`); w.Code != http.StatusForbidden {
		t.Errorf("bukan Admin: kode %d", w.Code)
	}
	if w := tombol(t, "revisi", "ReasTreatyInAdmin", `bukan json`); w.Code != http.StatusBadRequest {
		t.Errorf("bukan JSON: kode %d", w.Code)
	}
}

func tombolPenyesuaian(t *testing.T, jalur, peran, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/api/treaty-in/penyesuaian/"+jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", "ADESAMUEL")
	r.Header.Set("X-Peran", peran)
	w := httptest.NewRecorder()
	handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true).ServeHTTP(w, r)
	return w
}

func TestRutePenyesuaian(t *testing.T) {
	if w := tombolPenyesuaian(t, "simpan", "ReasTreatyInAdmin", `{"id":"1001001/R01","baru":{"medan":{"Comment":"x"}}}`); w.Code != http.StatusOK {
		t.Errorf("simpan: kode %d %s", w.Code, w.Body.String())
	}
	if w := tombolPenyesuaian(t, "kirim", "ReasTreatyInAdmin", `{"id":"1001001/R01","aksi":"submit"}`); w.Code != http.StatusOK {
		t.Errorf("kirim: kode %d %s", w.Code, w.Body.String())
	}
	if w := tombolPenyesuaian(t, "hapus", "ReasTreatyInAdmin", `{"id":"1001001/R01"}`); w.Code != http.StatusOK {
		t.Errorf("hapus: kode %d %s", w.Code, w.Body.String())
	}
	if w := tombolPenyesuaian(t, "simpan", "ReasTreatyInAdmin", `{"id":"9999999/R01","baru":{}}`); w.Code != http.StatusNotFound {
		t.Errorf("tak ada: kode %d", w.Code)
	}
}

func (_ gudangTiruan) NamaAplikasiSimpanan(context.Context) (string, error) { return "rnmtest", nil }

func (_ gudangTiruan) CatatLampiran(context.Context, models.LampiranBaru) (string, error) {
	return "20261008120000123", nil
}

func TestRuteUnggahLampiranMultipart(t *testing.T) {
	var badan bytes.Buffer
	mw := multipart.NewWriter(&badan)
	_ = mw.WriteField("kategori", "00002")
	fw, _ := mw.CreateFormFile("berkas", "a.pdf")
	_, _ = fw.Write([]byte("%PDF-1.4 isi berkas uji"))
	_ = mw.Close()
	r := httptest.NewRequest(http.MethodPost, "/api/treaty-in/kontrak/1001001/lampiran", &badan)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	r.Header.Set("X-Pelaku", "ADESAMUEL")
	r.Header.Set("X-Peran", "ReasTreatyInAdmin")
	w := httptest.NewRecorder()
	handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true).ServeHTTP(w, r)
	// Tiruan rute tanpa penyimpanan tersambung → 503 berpesan, bukan 500.
	if w.Code != http.StatusUnprocessableEntity && w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode %d %s", w.Code, w.Body.String())
	}
	// Bukan multipart → 400.
	r2 := httptest.NewRequest(http.MethodPost, "/api/treaty-in/kontrak/1001001/lampiran", strings.NewReader("{}"))
	r2.Header.Set("Content-Type", "application/json")
	r2.Header.Set("X-Pelaku", "ADESAMUEL")
	w2 := httptest.NewRecorder()
	handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true).ServeHTTP(w2, r2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("bukan multipart: kode %d", w2.Code)
	}
}

func (_ gudangTiruan) BacaObjekSimpanan(context.Context, string) (models.ObjekSimpanan, bool, error) {
	return models.ObjekSimpanan{}, false, nil
}

func (_ gudangTiruan) PerbaruiObjekSimpanan(context.Context, models.ObjekSimpanan, string) error {
	return nil
}

func (_ gudangTiruan) HapusLampiran(context.Context, string, string, string) error { return nil }

func (_ gudangTiruan) UbahKategoriLampiran(context.Context, string, []models.PerubahanKategori) error {
	return nil
}
