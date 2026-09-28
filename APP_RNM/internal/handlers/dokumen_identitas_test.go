package handlers

// Kontrak dua sisi unduh dokumen - GILIRAN-12 paket 0.
//
// ⛔ Cacat yang ditutup: `View Office Online` dulu pranala biasa, dan browser
// yang mengikuti pranala tidak membawa `X-Pelaku`/`X-Peran`; rute isi dokumen
// menjawab 401 untuk setiap unduhan. Sisi ini mengunci bahwa rute itu MEMANG
// menolak tanpa header (jadi klien wajib mengirimnya) dan bahwa nama header
// yang dibaca di sini sama dengan yang dikirim klien. Sisi TypeScript:
// `frontend/src/services/unduhdokumen.test.ts`.

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nusantarare/internal/services"
)

func TestUnduhDokumenTanpaHeaderIdentitasDitolak401(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/dokumen/20250101120000123/isi", nil)
	// Stub MENYALA, tetapi permintaan tanpa header - seperti pranala biasa.
	p := pelakuDari(r, true)
	if p.AkunID != "" {
		t.Fatalf("pelaku tanpa header = %q, mau kosong", p.AkunID)
	}
	// Layanan menolak pelaku kosong SEBELUM menyentuh basis data.
	_, _, err := services.New(nil).Dokumen().UnduhLewatPengenal(context.Background(), p, 1)
	if !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Fatalf("galat = %v, mau ErrTanpaIdentitas", err)
	}
	w := httptest.NewRecorder()
	if !jawabGalatDokumen(w, err) || w.Code != http.StatusUnauthorized {
		t.Fatalf("kode = %d, mau 401", w.Code)
	}

	r.Header.Set("X-Pelaku", "UJI-ADMIN")
	r.Header.Set("X-Peran", "ReasLifeAdmin")
	if p := pelakuDari(r, true); p.AkunID != "UJI-ADMIN" || len(p.Peran) != 1 {
		t.Errorf("pelaku dengan header = %+v", p)
	}
}

// Nama header yang dibaca pelakuDari HARUS sama dengan yang dikirim klien.
func TestNamaHeaderIdentitasSamaDenganKlien(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "..", "frontend", "src", "store", "sesi.ts"))
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range []string{"'X-Pelaku'", "'X-Peran'"} {
		if !strings.Contains(string(isi), h) {
			t.Errorf("klien tidak mengirim header %s; pelakuDari membacanya", h)
		}
	}
}
