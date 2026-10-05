package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/menu"
	"nusantarare/modul/bordereaux/backend/handlers"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
	"nusantarare/modul/bordereaux/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	jam := time.Date(2026, 10, 4, 9, 30, 0, 0, time.UTC)
	l := services.BaruLayanan(g, tiruan.Transaksi).DenganJam(func() time.Time { return jam })
	return handlers.RouterDengan(l, adaDB, false)
}

// kirim - permintaan dari sesi login `akun` berperan `peran`; menu = menu akun (Kelola User = superadmin).
func kirim(h http.Handler, metode, jalur, badan, akun string, peran []string, menuAkun ...string) *httptest.ResponseRecorder {
	return kirimHak(h, metode, jalur, badan, akun, peran, nil, menuAkun...)
}

// kirimHak - seperti kirim, dengan menu ber-hak LIHAT (`M_LOGIN_GO_MENU.HAK`, migrasi 914).
func kirimHak(h http.Handler, metode, jalur, badan, akun string, peran, lihat []string, menuAkun ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	ctx := inti.DenganAksesMenu(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: akun, Peran: peran}), menuAkun)
	ctx = inti.DenganMenuLihat(ctx, lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func TestRuteBordereaux(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	var admin []string // Input Data = hak menu PENUH, bukan workbasket (keputusan work owner 04-10-2026)
	lihat := []string{handlers.KodeMenu}
	if w := kirim(h, "GET", handlers.Prefix+"?ceding=satu", "", "UJI-CHK", []string{models.PeranChecker}); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"bdxId":"BDX-UJI.1"`) || !strings.Contains(w.Body.String(), `"putuskan":true`) {
		t.Errorf("daftar %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-MAKER", admin); w.Code != 200 || !strings.Contains(w.Body.String(), `"bolehBuat":true`) {
		t.Errorf("pilihan %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-SUPER", nil, menu.KodeKelolaUser); !strings.Contains(w.Body.String(), `"bolehBuat":true`) {
		t.Error("superadmin = pemegang Kelola User boleh Input Data")
	}
	badan := `{"type":"PREMIUM","business":"FIRE","masterId":"UJI-TI-1","reportStart":"01-07-2026","reportEnd":"30-09-2026","reffNoSoa":"","reffNoBdx":"","baris":[{"COB":"UJI"}]}`
	if w := kirimHak(h, "POST", handlers.Prefix+"/simpan", badan, "UJI-TAMU", nil, lihat, handlers.KodeMenu); w.Code != 403 {
		t.Errorf("menu View only %d %s", w.Code, w.Body.String())
	}
	if w := kirimHak(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-TAMU", nil, lihat, handlers.KodeMenu); !strings.Contains(w.Body.String(), `"bolehBuat":false`) {
		t.Errorf("View only tanpa Input Data: %s", w.Body.String())
	}
	w := kirim(h, "POST", handlers.Prefix+"/simpan", badan, "UJI-MAKER", admin)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"bdxId":"BDX-2026.10.04.00000"`) {
		t.Fatalf("simpan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/berkas/BDX-2026.10.04.00000", "", "UJI-MAKER", admin); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"ubah":true`) || !strings.Contains(w.Body.String(), `"COB":"UJI"`) {
		t.Errorf("buka %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/berkas/BDX-UJI.1/submit", `{"setuju":false,"komentar":""}`, "UJI-CHK", []string{models.PeranChecker}); w.Code != 422 {
		t.Errorf("tolak tanpa komentar %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix+"/berkas/BDX-UJI.1/hapus", ``, "UJI-CHK", []string{models.PeranChecker}); w.Code != 403 {
		t.Errorf("hapus di checker %d", w.Code)
	}
	csv := "A;B\r\nx;y\r\n"
	if w := kirim(h, "POST", handlers.Prefix+"/unggah-csv", `{"type":"PREMIUM","business":"FIRE","csv":"`+csv+`"}`, "UJI-MAKER", admin); w.Code != 400 {
		// \r\n mentah di string JSON = JSON tidak sah
		t.Errorf("csv mentah %d", w.Code)
	}
	if w := kirim(h, "POST", handlers.Prefix+"/unggah-csv", `{"type":"PREMIUM","business":"FIRE","csv":"A;B\nx;y\n"}`, "UJI-MAKER", admin); w.Code != 422 ||
		!strings.Contains(w.Body.String(), "needs 51 columns") {
		t.Errorf("csv salah templat %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/master-treaty?ceding=UJI-AG-1", "", "UJI-MAKER", admin); !strings.Contains(w.Body.String(), `"id":"UJI-TI-2"`) {
		t.Errorf("treaty %s", w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/chart", "", "UJI-MAKER", admin); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `{"business":"FIRE","type":"PREMIUM","cedingId":"UJI-AG-1","cedingName":"UJI CEDANT SATU","jumlah":1}`) {
		t.Errorf("chart %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", handlers.Prefix+"/berkas/TIDAK-ADA", "", "UJI-MAKER", admin); w.Code != 404 {
		t.Errorf("tak ada %d", w.Code)
	}
	if w := kirim(router(g, false), "GET", handlers.Prefix, "", "UJI-MAKER", admin); w.Code != 503 {
		t.Errorf("tanpa db %d", w.Code)
	}
}

// Copy Old Data: hanya superadmin (pemegang Kelola User); badan Process Copy {"ids": [...]}.
func TestRuteCopyOld(t *testing.T) {
	g := tiruan.Contoh()
	g.JSON["BDX-UJI.1"] = `{"BDX_ID":"BDX-UJI.1","TYPE":"PREMIUM","TYPE_BUSINESS":"FIRE","PremiumFireList":[{"COB":"UJI-LAMA"}]}`
	h := router(g, true)
	if w := kirim(h, "GET", handlers.Prefix+"/lama", "", "UJI-MAKER", nil); w.Code != 403 {
		t.Errorf("bukan superadmin %d", w.Code)
	}
	if w := kirim(h, "GET", handlers.Prefix+"/pilihan", "", "UJI-SUPER", nil, menu.KodeKelolaUser); !strings.Contains(w.Body.String(), `"copyOld":true`) {
		t.Error("tombol Copy Old Data untuk superadmin")
	}
	if w := kirim(h, "GET", handlers.Prefix+"/lama", "", "UJI-SUPER", nil, menu.KodeKelolaUser); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"id":"BDX-UJI.1"`) || !strings.Contains(w.Body.String(), `"barisJson":1`) {
		t.Errorf("daftar lama %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/lama/salin", `{"ids":["BDX-UJI.1"]}`, "UJI-SUPER", nil, menu.KodeKelolaUser); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"status":"disalin"`) || !strings.Contains(w.Body.String(), `"disalin":1`) {
		t.Errorf("salin %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", handlers.Prefix+"/lama/salin", `{"ids":[]}`, "UJI-SUPER", nil, menu.KodeKelolaUser); w.Code != 422 {
		t.Errorf("tanpa ID %d", w.Code)
	}
}
