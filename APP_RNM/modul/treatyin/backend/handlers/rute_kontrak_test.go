package handlers_test

// Uji jalur HTTP tiket 14: 201 / 400 / 401 / 404 / 422 / 503.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nusantarare/modul/treatyin/backend/handlers"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func badanSah() []byte {
	b, _ := json.Marshal(services.MasukanKontrak{
		IDCedant: 1, IDAsalBisnis: 2, SifatProporsi: models.SifatNonProporsional,
		TanggalMulai: "2026-01-01", TanggalBerakhir: "2026-12-31",
		NamaKontrak: "Kontrak Uji", KodeMataUangKontrak: 3, PersenBagianNure: "12.5",
		BagianNureSeragam: "0", MemakaiBordereaux: "0", CaraPembukuan: "X",
		MemakaiProrata: "0", RetroBerganda: "0", KeadaanSiklusHidup: "DRAFT",
	})
	return b
}

func kirim(t *testing.T, h http.Handler, metode, jalur, pelaku string, badan []byte) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if badan == nil {
		r = httptest.NewRequest(metode, jalur, nil)
	} else {
		r = httptest.NewRequest(metode, jalur, bytes.NewReader(badan))
	}
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestBuatKontrakMenjawab201(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "AKUN-UJI", badanSah())

	if w.Code != http.StatusCreated {
		t.Fatalf("mau 201, dapat %d (badan %s)", w.Code, w.Body.String())
	}
	var hasil services.HasilBuatKontrak
	if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}
	if hasil.IDKontrak == 0 || hasil.IDVersi == 0 {
		t.Errorf("pengenal tidak dikembalikan: %+v", hasil)
	}
}

func TestBuatKontrakTanpaIdentitasMenjawab401(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	if w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "", badanSah()); w.Code != http.StatusUnauthorized {
		t.Fatalf("mau 401, dapat %d", w.Code)
	}
}

// JSON rusak adalah BENTUK permintaan yang salah (400); JSON sah yang isinya
// ditolak gerbang adalah 422. Dua pertanyaan berbeda.
func TestBadanRusakMenjawab400DanGerbangMenjawab422(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	if w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "AKUN-UJI", []byte("{bukan json")); w.Code != http.StatusBadRequest {
		t.Errorf("JSON rusak: mau 400, dapat %d", w.Code)
	}

	var m services.MasukanKontrak
	_ = json.Unmarshal(badanSah(), &m)
	m.SifatProporsi = "entah"
	rusak, _ := json.Marshal(m)
	w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "AKUN-UJI", rusak)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("gerbang: mau 422, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

// INV-04 ditolak basis data -> 409, dan pesannya sampai ke layar.
func TestNomorUrutVersiGandaMenjawab409(t *testing.T) {
	g := gudangTiruan{galat: services.ErrNomorUrutVersiGanda}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "AKUN-UJI", badanSah())

	if w.Code != http.StatusConflict {
		t.Fatalf("mau 409, dapat %d (badan %s)", w.Code, w.Body.String())
	}
}

func TestBacaKontrakMenjawab200(t *testing.T) {
	satu := int64(1)
	g := gudangTiruan{kontrak: models.KontrakDenganVersi{
		Kontrak: models.Kontrak{ID: 777, SifatProporsi: models.SifatNonProporsional},
		Versi:   []models.VersiKontrak{{ID: 888, IDKontrak: 777, NomorUrutVersi: &satu}},
	}}
	h := handlers.RouterDengan(services.LayananDengan(g), true, true)

	w := kirim(t, h, http.MethodGet, handlers.Prefix+"/kontrak/777", "AKUN-UJI", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("mau 200, dapat %d (badan %s)", w.Code, w.Body.String())
	}
	var k models.KontrakDenganVersi
	if err := json.Unmarshal(w.Body.Bytes(), &k); err != nil {
		t.Fatalf("badan bukan JSON: %v", err)
	}
	if k.Kontrak.ID != 777 || len(k.Versi) != 1 {
		t.Errorf("bentuk tidak sesuai: %+v", k)
	}
}

func TestBacaKontrakPengenalBukanAngkaMenjawab400(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, true)

	if w := kirim(t, h, http.MethodGet, handlers.Prefix+"/kontrak/tujuh", "AKUN-UJI", nil); w.Code != http.StatusBadRequest {
		t.Fatalf("mau 400, dapat %d", w.Code)
	}
}

func TestBacaKontrakTidakAdaMenjawab404(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{galat: services.ErrKontrakTidakAda}), true, true)

	if w := kirim(t, h, http.MethodGet, handlers.Prefix+"/kontrak/999", "AKUN-UJI", nil); w.Code != http.StatusNotFound {
		t.Fatalf("mau 404, dapat %d", w.Code)
	}
}

func TestKontrakTanpaBasisDataMenjawab503(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), false, true)

	if w := kirim(t, h, http.MethodPost, handlers.Prefix+"/kontrak", "AKUN-UJI", badanSah()); w.Code != http.StatusServiceUnavailable {
		t.Fatalf("mau 503, dapat %d", w.Code)
	}
}
