package handlers_test

// Uji jalur HTTP tombol `Copy` — `GET /kontrak-warisan/{id}/salin` dan `POST /kontrak/salin`.
//
// Kepala tiruan paket ini tidak `Resolve Complete`, jadi yang diuji di sini
// JALURNYA dan pemetaan galatnya; aturan `TreatyInCopy` diuji di
// `services/salin_kontrak_test.go`.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/services"
)

func TestRuteSalinTerpasangDanMenegakkanSyaratTampil(t *testing.T) {
	// Sumber belum Resolve Complete — 422 berpesan, bukan 404 rute.
	if w := tombol(t, "salin", "ReasTreatyInAdmin", `{"idSumber":"1001001"}`); w.Code != http.StatusUnprocessableEntity {
		t.Errorf("POST: kode %d %s", w.Code, w.Body.String())
	}
	// Bukan workbasket Admin — 403.
	if w := tombol(t, "salin", "ReasTreatyInSecHead", `{"idSumber":"1001001"}`); w.Code != http.StatusForbidden {
		t.Errorf("POST SecHead: kode %d %s", w.Code, w.Body.String())
	}
	if w := tombol(t, "salin", "ReasTreatyInAdmin", `{bukan json`); w.Code != http.StatusBadRequest {
		t.Errorf("JSON rusak: kode %d", w.Code)
	}
	r := httptest.NewRequest(http.MethodGet, "/api/treaty-in/kontrak-warisan/1001001/salin", nil)
	r.Header.Set("X-Pelaku", "ADESAMUEL")
	r.Header.Set("X-Peran", "ReasTreatyInAdmin")
	w := httptest.NewRecorder()
	handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true).ServeHTTP(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("GET draf: kode %d %s", w.Code, w.Body.String())
	}
}
