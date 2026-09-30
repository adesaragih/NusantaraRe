package handlers

// Penjaga jawaban penghapusan klaim - perapian A2.
//
// Untuk apa berkas ini: memastikan `DELETE /api/klaim-life/{id}` menjawab 405
// dengan header `Allow` yang JUJUR - yaitu yang benar-benar cocok dengan rute
// yang terdaftar.
//
// ⛔ Header `Allow` yang berbohong lebih buruk daripada tidak ada: klien
// mempercayainya, lalu menembak metode yang tidak pernah terdaftar dan
// mendapat 404 yang membingungkan.

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/backend/services"
)

// metodeTerdaftarUntukKlaim membaca handlers.go dan mengumpulkan metode yang
// benar-benar didaftarkan untuk `/api/klaim-life/{id}` PERSIS - bukan
// sub-jalurnya.
func metodeTerdaftarUntukKlaim(t *testing.T) []string {
	t.Helper()
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatalf("membaca handlers.go: %v", err)
	}
	pola := regexp.MustCompile(`"([A-Z]+) (/api/klaim-life/\{id\})"`)
	var metode []string
	for _, m := range pola.FindAllStringSubmatch(string(isi), -1) {
		metode = append(metode, m[1])
	}
	if len(metode) == 0 {
		t.Fatal("nol rute terbaca untuk /api/klaim-life/{id}; pembacanya yang rusak")
	}
	sort.Strings(metode)
	return metode
}

// Header `Allow` menyebut tepat metode yang terdaftar, dikurangi DELETE.
//
// DELETE dikecualikan sebab justru DELETE-lah yang ditolak 405; menyebutnya di
// `Allow` berarti menyuruh klien mencoba lagi hal yang sama.
func TestAllowCocokDenganRuteTerdaftar(t *testing.T) {
	var mau []string
	for _, m := range metodeTerdaftarUntukKlaim(t) {
		if m == http.MethodDelete {
			continue
		}
		mau = append(mau, m)
	}
	harap := strings.Join(mau, ", ")
	if metodeKlaimDiizinkan != harap {
		t.Errorf("metodeKlaimDiizinkan = %q, sedangkan rute terdaftar "+
			"(tanpa DELETE) = %q - header Allow berbohong",
			metodeKlaimDiizinkan, harap)
	}
}

// ErrHapusFisikDilarang dijawab 405, bukan 501, dan membawa header Allow.
//
// ⛔ 501 berarti "belum bisa", yaitu janji bahwa suatu hari bisa. ADR-U-0031
// menetapkan sebaliknya: hapus fisik TIDAK PERNAH berlaku.
func TestHapusFisikDijawab405DenganAllow(t *testing.T) {
	w := httptest.NewRecorder()
	ditangani := jawabGalatHapus(w, services.ErrHapusFisikDilarang)
	if !ditangani {
		t.Fatal("galat hapus fisik tidak ditangani sama sekali")
	}
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("kode = %d, mau %d", w.Code, http.StatusMethodNotAllowed)
	}
	if got := w.Header().Get("Allow"); got != metodeKlaimDiizinkan {
		t.Errorf("header Allow = %q, mau %q", got, metodeKlaimDiizinkan)
	}
	// Kalimat kebijakannya ikut, bukan sekadar kodenya: pembaca jawaban harus
	// tahu MENGAPA, dan ADR-nya disebut.
	if badan := w.Body.String(); !strings.Contains(badan, "ADR-U-0031") {
		t.Errorf("jawaban tidak menyebut ADR-U-0031: %s", badan)
	}
}

// Galat yang dibungkus tetap terbaca - pemanggil membungkusnya dengan sebab.
func TestHapusFisikDibungkusTetap405(t *testing.T) {
	w := httptest.NewRecorder()
	bungkus := errors.Join(services.ErrHapusFisikDilarang,
		errors.New("kolom penandanya belum ada"))
	if !jawabGalatHapus(w, bungkus) {
		t.Fatal("galat terbungkus tidak ditangani")
	}
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("kode = %d, mau %d", w.Code, http.StatusMethodNotAllowed)
	}
}
