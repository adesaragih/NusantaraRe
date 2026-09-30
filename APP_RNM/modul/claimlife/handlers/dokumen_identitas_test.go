package handlers

// Kontrak dua sisi unduh dokumen - GILIRAN-12 paket 0.
//
// ⛔ Cacat yang ditutup: `View Office Online` dulu pranala biasa, dan browser
// yang mengikuti pranala tidak membawa `X-Pelaku`/`X-Peran`; rute isi dokumen
// menjawab 401 untuk setiap unduhan. Sisi ini mengunci bahwa HANDLER rute itu
// menolak tanpa header (jadi klien wajib mengirimnya) dan bahwa nama header
// yang dibaca di sini sama dengan yang dikirim klien. Sisi TypeScript:
// `frontend/src/modul/claim-life/unduhdokumen.test.ts`.

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/claimlife/services"
)

func TestUnduhDokumenTanpaHeaderIdentitasDitolak401(t *testing.T) {
	// Handler SUNGGUHAN, stub menyala, tanpa Oracle.
	h := isiDokumen(services.New(nil), true)

	r := httptest.NewRequest(http.MethodGet, "/api/dokumen/20250101120000123/isi", nil)
	r.SetPathValue("dokId", "20250101120000123")
	w := httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("tanpa header: kode = %d, mau 401", w.Code)
	}

	// Dengan header, identitasnya lolos dan permintaan sampai ke gerbang
	// berikutnya - basis data, yang di uji ini memang tidak ada.
	r = httptest.NewRequest(http.MethodGet, "/api/dokumen/20250101120000123/isi", nil)
	r.SetPathValue("dokId", "20250101120000123")
	r.Header.Set("X-Pelaku", "UJI-ADMIN")
	r.Header.Set("X-Peran", inti.PeranAdmin)
	w = httptest.NewRecorder()
	h(w, r)
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("dengan header: kode = %d, mau 503 (identitas lolos, basis data tidak ada)", w.Code)
	}
}

// Nama header yang dibaca pelakuDari HARUS sama dengan yang dikirim klien.
func TestNamaHeaderIdentitasSamaDenganKlien(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "..", "inti", "frontend", "store", "sesi.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{"'X-Pelaku'", "'X-Peran'"} {
		if !strings.Contains(string(isi), h) {
			t.Errorf("klien tidak mengirim header %s; pelakuDari membacanya", h)
		}
	}
}
