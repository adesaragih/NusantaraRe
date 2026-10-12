package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/claimlife/backend/services"
)

// 11-10-2026: dropdown "Category" dibaca dari rute ini; identitas diperiksa SEBELUM basis data.
func TestRuteKategoriDokumen(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `mux.HandleFunc("GET /api/klaim-life/{id}/dokumen/kategori", kategoriDokumen(svc, stubPelaku))`) {
		t.Error("rute kategori dokumen tidak terdaftar")
	}
	w := httptest.NewRecorder()
	kategoriDokumen(services.New(nil), false)(w, httptest.NewRequest(http.MethodGet, "/api/klaim-life/K-1/dokumen/kategori", nil))
	if w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", w.Code)
	}
}

// 11-10-2026: "View" berbatas klaim memakai Unduh(klaimID, ...) - dokumen klaim lain = 404 (SatuDokumen berbatas klaim).
func TestRuteViewDokumenBerbatasKlaim(t *testing.T) {
	isi, err := os.ReadFile("handlers.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(isi), `mux.HandleFunc("GET /api/klaim-life/{id}/dokumen/{dokId}/isi", isiDokumenKlaim(svc, stubPelaku))`) {
		t.Error("rute View berbatas klaim tidak terdaftar")
	}
	h, err := os.ReadFile("dokumen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(h), `svc.Dokumen().BukaDokumen(r.Context(), pelaku, r.PathValue("id"), id)`) {
		t.Error("View tidak berbatas klaim")
	}
	w := httptest.NewRecorder()
	isiDokumenKlaim(services.New(nil), false)(w, httptest.NewRequest(http.MethodGet, "/api/klaim-life/K/dokumen/1/isi", nil))
	if w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", w.Code)
	}
}
