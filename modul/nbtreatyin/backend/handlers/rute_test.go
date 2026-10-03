package handlers_test

// Uji seam 1 - HTTP di atas layanan dengan gudang tiruan. Yang diuji: rute
// terpasang di prefix modul, pelaku dibaca dari header stub, dan galat
// layanan sampai sebagai kode HTTP yang benar.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/nbtreatyin/backend/handlers"
	"nusantarare/modul/nbtreatyin/backend/services"
	"nusantarare/modul/nbtreatyin/backend/tiruan"
)

func jam() time.Time { return time.Date(2026, 10, 3, 9, 0, 0, 0, time.UTC) }

func server(g *tiruan.Gudang) http.Handler {
	var l *services.Layanan
	if g == nil {
		l = services.Baru(nil, jam)
	} else {
		l = services.Baru(g, jam)
	}
	return handlers.Router(l, true)
}

func minta(t *testing.T, h http.Handler, metode, jalur, pelaku, peran string, badan any) *httptest.ResponseRecorder {
	t.Helper()
	var b bytes.Buffer
	if badan != nil {
		if err := json.NewEncoder(&b).Encode(badan); err != nil {
			t.Fatal(err)
		}
	}
	r := httptest.NewRequest(metode, jalur, &b)
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	if peran != "" {
		r.Header.Set("X-Peran", peran)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestTanpaOracle503(t *testing.T) {
	w := minta(t, server(nil), "GET", "/api/nb-treaty-in/kasus", "UJI-A", "", nil)
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(w.Body.String(), `"galat"`) {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}

func TestTanpaIdentitas401(t *testing.T) {
	w := minta(t, server(tiruan.Baru()), "POST", "/api/nb-treaty-in/kasus", "", "", nil)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
}

func TestAlurHTTP(t *testing.T) {
	g := tiruan.Baru()
	s := server(g)
	w := minta(t, s, "POST", "/api/nb-treaty-in/kasus", "UJI-A", "ReasTreatyInAdmin", nil)
	if w.Code != http.StatusCreated {
		t.Fatalf("buat: %d %s", w.Code, w.Body)
	}
	var k struct{ ID string }
	if err := json.Unmarshal(w.Body.Bytes(), &k); err != nil || k.ID == "" {
		t.Fatalf("buat: %v %s", err, w.Body)
	}
	jalur := "/api/nb-treaty-in/kasus/" + k.ID

	if w := minta(t, s, "GET", "/api/nb-treaty-in/kasus/UJI-404", "UJI-A", "", nil); w.Code != http.StatusNotFound {
		t.Fatalf("kasus tak ada: %d", w.Code)
	}
	w = minta(t, s, "GET", jalur, "UJI-A", "ReasTreatyInAdmin", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"bolehKerja":true`) {
		t.Fatalf("buka: %d %s", w.Code, w.Body)
	}
	// bukan anggota antrean -> 403
	if w := minta(t, s, "POST", jalur+"/kirim", "UJI-B", "UJI-LAIN", map[string]any{}); w.Code != http.StatusForbidden {
		t.Fatalf("bukan anggota: %d %s", w.Code, w.Body)
	}
	// medan wajib kosong -> 422 beserta pesannya
	badan := map[string]any{"halaman": map[string]any{"nilai": map[string]string{"PolicyTreatyIn.IsApproved": "1"}}}
	w = minta(t, s, "POST", jalur+"/kirim", "UJI-A", "ReasTreatyInAdmin", badan)
	if w.Code != http.StatusUnprocessableEntity || !strings.Contains(w.Body.String(), "Marketing Officer") {
		t.Fatalf("validasi: %d %s", w.Code, w.Body)
	}
	// aksi hitung karangan -> 400
	w = minta(t, s, "POST", jalur+"/hitung", "UJI-A", "ReasTreatyInAdmin", map[string]any{"urutan": []map[string]string{{"aksi": "UJI-Karangan"}}})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("aksi karangan: %d %s", w.Code, w.Body)
	}
	// SATU bentuk permintaan hitung: urutan kosong (termasuk bentuk lama aksi/param) -> 400
	w = minta(t, s, "POST", jalur+"/hitung", "UJI-A", "ReasTreatyInAdmin", map[string]any{"aksi": "CountNetPremi"})
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "urutan hitung kosong") {
		t.Fatalf("tanpa urutan: %d %s", w.Code, w.Body)
	}
	// pilih bisnis atas kontrak yang tidak ada -> 422 (AC 37: galat ditampilkan)
	w = minta(t, s, "POST", jalur+"/pilih-bisnis", "UJI-A", "ReasTreatyInAdmin", map[string]any{"idDetail": "UJI-X"})
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("kontrak tak ada: %d %s", w.Code, w.Body)
	}
	// nomor polis di posisi admin -> 409
	w = minta(t, s, "POST", jalur+"/nomor-polis", "UJI-A", "ReasTreatyInAdmin", map[string]any{})
	if w.Code != http.StatusConflict {
		t.Fatalf("nomor polis di admin: %d %s", w.Code, w.Body)
	}
	// daftar portal: grid hanya di wadah ReasTreatyInAdmin (portal_test.go)
	if w := minta(t, s, "GET", "/api/nb-treaty-in/kasus", "UJI-A", "ReasTreatyInAdmin", nil); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), k.ID) {
		t.Fatalf("daftar: %d %s", w.Code, w.Body)
	}
	if w := minta(t, s, "GET", "/api/nb-treaty-in/kasus", "UJI-A", "", nil); w.Code != http.StatusForbidden {
		t.Fatalf("daftar tanpa antrean: %d %s", w.Code, w.Body)
	}
	if w := minta(t, s, "GET", "/api/nb-treaty-in/acuan", "UJI-A", "", nil); w.Code != http.StatusOK {
		t.Fatalf("acuan: %d %s", w.Code, w.Body)
	}
	if w := minta(t, s, "GET", jalur+"/riwayat", "UJI-A", "", nil); w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatalf("riwayat: %d %s", w.Code, w.Body)
	}
	if w := minta(t, s, "PUT", jalur, "UJI-A", "ReasTreatyInAdmin", "bukan-objek"); w.Code != http.StatusBadRequest {
		t.Fatalf("badan rusak: %d %s", w.Code, w.Body)
	}
}
