package handlers

// Uji pintu HTTP kaskade hapus - TANPA Oracle (tiket 10).

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/internal/services"
)

func TestRuteKaskadeTerdaftarDanTanpaDatabase503(t *testing.T) {
	router := Router(services.New(nil), true)
	dasar := "/api/treaty-contract-out/tahun/1000001/kontrak/1000003"
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodGet, dasar + "/dampak-hapus"},
		{http.MethodDelete, dasar + "?reinsurer=0&security=0&business=0&bersama=0"},
		{http.MethodGet, dasar + "/reinsurer/1000007/dampak-hapus"},
		{http.MethodDelete, dasar + "/reinsurer/1000007?security=0"},
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(r.metode, r.jalur, nil)
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s %s: kode %d, mau 503", r.metode, r.jalur, w.Code)
		}
	}
	isi, _ := os.ReadFile("rute_treaty_contract_out.go")
	if strings.Count(string(isi), "daftarkanRuteKaskadeTCO(mux, svc, stubPelaku)") != 1 {
		t.Error("rute kaskade harus didaftarkan tepat sekali")
	}
	// Klausul tidak punya jalur hapus di konteks ini (AC 44).
	sumber, _ := os.ReadFile("tco_kaskade.go")
	if strings.Contains(string(sumber), "klausul") && strings.Contains(string(sumber), `"DELETE `+"/api/treaty-contract-out/tahun/{id}/klausul") {
		t.Error("jalur hapus klausul")
	}
}

func TestKonfirmasiDariKueri(t *testing.T) {
	q := httptest.NewRequest(http.MethodDelete, "/x?reinsurer=2&security=3&business=1", nil)
	k, ok := konfirmasiDariKueri(q, "reinsurer", "security", "business")
	if !ok || k.Reinsurer != 2 || k.Security != 3 || k.Business != 1 {
		t.Errorf("%+v %v", k, ok)
	}
	for _, jalur := range []string{"/x?reinsurer=2&security=3", "/x?reinsurer=-1&security=0&business=0", "/x?reinsurer=a&security=0&business=0"} {
		if _, ok := konfirmasiDariKueri(httptest.NewRequest(http.MethodDelete, jalur, nil), "reinsurer", "security", "business"); ok {
			t.Errorf("%s diterima", jalur)
		}
	}
}

func TestJawabGalatKaskadeTCO(t *testing.T) {
	for _, err := range []error{services.ErrDampakBerubah, services.ErrKaskadeTidakUtuh, services.ErrKontrakBeranak, services.ErrTahunBeranak} {
		w := httptest.NewRecorder()
		jawabGalatTreatyContractOut(w, err)
		if w.Code != http.StatusConflict {
			t.Errorf("%v: %d", err, w.Code)
		}
	}
}

// OQ-TCO-19 (keputusan work owner 29-09-2026, "tidak perlu"): rute simpan
// kontrak utuh dibuang - penyimpanan per panel seperti Pega.
func TestRuteSimpanUtuhDibuang(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, r := range []struct{ metode, jalur string }{
		{http.MethodPost, "/api/treaty-contract-out/tahun/1000001/kontrak-utuh"},
		{http.MethodPut, "/api/treaty-contract-out/tahun/1000001/kontrak/1000003/utuh"},
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(r.metode, r.jalur, strings.NewReader("{}"))
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if w.Code == http.StatusServiceUnavailable || w.Code == http.StatusOK {
			t.Errorf("%s %s masih dilayani (%d)", r.metode, r.jalur, w.Code)
		}
	}
	isi, _ := os.ReadFile("rute_treaty_contract_out.go")
	if strings.Contains(string(isi), "SimpanUtuh") || strings.Contains(string(isi), "kontrak-utuh") {
		t.Error("rute simpan utuh masih terdaftar")
	}
}
