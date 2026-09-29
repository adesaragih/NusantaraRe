package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/treaty/services"
)

// Kontrak 503 dua sisi, SISI BACKEND - lanjutan 6 Treaty Contract Out.
//
// ⛔ Sebab uji ini ada, dan itu sungguh terjadi (laporan work owner
// 29-09-2026): `Show TreatyDesc` menampilkan "Backend tidak terhubung" padahal
// backend menyala. Rute kurs menjawab 503 `{"galat": …}` - penolakan backend
// sendiri - dan klien menganggap setiap 503 datang dari proxy yang kehilangan
// upstream-nya. Kini klien membedakan keduanya dari BADANNYA: 503 yang
// membawa `galat` tak kosong adalah jawaban backend. Kontrak itu hanya
// bertahan selama setiap 503 backend memang membawa `galat` - dan itulah yang
// dikunci di sini. Pasangan uji sisi klien:
// `frontend/src/modul/treaty/keadaanGalat.test.ts`.

// badanGalat503 memeriksa satu jawaban: 503, JSON, SATU kunci `galat` tak kosong.
func badanGalat503(t *testing.T, nama string, w *httptest.ResponseRecorder) string {
	t.Helper()
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("%s: kode %d, mau 503", nama, w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("%s: Content-Type %q - klien membaca badan bukan-JSON sebagai backend mati", nama, ct)
	}
	var isi map[string]any
	if err := json.NewDecoder(w.Body).Decode(&isi); err != nil {
		t.Fatalf("%s: badan bukan JSON: %v", nama, err)
	}
	pesan, _ := isi["galat"].(string)
	if len(isi) != 1 || pesan == "" {
		t.Fatalf("%s: badan %v, mau tepat {\"galat\": \"<kalimat>\"}", nama, isi)
	}
	return pesan
}

func TestGalat503MasterKursMembawaKalimatnya(t *testing.T) {
	err := fmt.Errorf("%w: UJI sebab", services.ErrMasterKursRusak)
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, err)

	if pesan := badanGalat503(t, "master kurs rusak", w); pesan != err.Error() {
		t.Errorf("galat = %q, mau kalimat services apa adanya %q", pesan, err.Error())
	}
}

func TestGalat503RuteKursTanpaDatabaseMembawaKalimatnya(t *testing.T) {
	router := Router(services.New(nil), true)
	for _, jalur := range []string{
		"/api/treaty-contract-out/tahun/1000001/kurs",
		"/api/treaty-contract-out/tahun/1000001/kurs/konversi?dari=Rp&nilai=1&skala=8",
	} {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(http.MethodGet, jalur, nil)
		q.Header.Set("X-Pelaku", "UJI-ADMIN")
		router.ServeHTTP(w, q)
		if pesan := badanGalat503(t, jalur, w); pesan != "database belum dikonfigurasi" {
			t.Errorf("%s: galat = %q", jalur, pesan)
		}
	}
}
