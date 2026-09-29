package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/internal/services"
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
// `frontend/src/lib/keadaanGalat.test.ts`.

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
	if len(isi) != 1 || strings.TrimSpace(pesan) == "" {
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
		"/api/treaty-contract-out/tahun/1000682/kurs",
		"/api/treaty-contract-out/tahun/1000682/kurs/konversi?dari=Rp&nilai=1&skala=8",
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

// Penjaga statik: SETIAP 503 di handler ditulis lewat `galat()`. 503 yang
// ditulis `http.Error` (teks biasa) atau `WriteHeader` polos sampai di klien
// tanpa `galat` - dan klien benar menyebutnya "backend tidak terhubung".
func TestSetiap503HandlerLewatGalat(t *testing.T) {
	berkas, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	// Pesannya boleh di baris berikutnya (`dokumen.go`); kosong dilarang.
	lewatGalat := regexp.MustCompile(`galat\(w, http\.StatusServiceUnavailable,`)
	pesanKosong := regexp.MustCompile(`galat\(w, http\.StatusServiceUnavailable,\s*""\)`)
	jumlah := 0
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		for i, baris := range strings.Split(string(isi), "\n") {
			if !strings.Contains(baris, "StatusServiceUnavailable") && !strings.Contains(baris, "503") {
				continue
			}
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			switch {
			case pesanKosong.MatchString(baris):
				t.Errorf("%s:%d: 503 dengan galat kosong - klien membacanya sebagai backend mati", b, i+1)
			case lewatGalat.MatchString(baris):
				jumlah++
			default:
				t.Errorf("%s:%d: 503 tidak lewat galat(): %s", b, i+1, strings.TrimSpace(baris))
			}
		}
	}
	if jumlah == 0 {
		t.Error("nol 503 ditemukan - penjaga ini tidak lagi membaca berkas yang benar")
	}
	t.Logf("%d jawaban 503 lewat galat()", jumlah)
}
