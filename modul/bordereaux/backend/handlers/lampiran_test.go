package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/bordereaux/backend/handlers"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/services"
	"nusantarare/modul/bordereaux/backend/tiruan"
)

func routerLampiran(g *tiruan.Gudang, p *tiruan.Penyimpanan) http.Handler {
	jam := time.Date(2026, 10, 8, 9, 30, 0, 0, time.UTC)
	l := services.BaruLayanan(g, tiruan.Transaksi).DenganJam(func() time.Time { return jam }).DenganPenyimpanan(p)
	return handlers.RouterDengan(l, true, false)
}

// unggahBerkas - POST multipart `berkas` dari sesi `akun`; lihat = menu ber-hak View only.
func unggahBerkas(h http.Handler, jalur, nama string, isi []byte, akun string, lihat []string) *httptest.ResponseRecorder {
	var badan bytes.Buffer
	m := multipart.NewWriter(&badan)
	if nama != "" {
		f, _ := m.CreateFormFile("berkas", nama)
		_, _ = f.Write(isi)
	}
	_ = m.Close()
	r := httptest.NewRequest("POST", jalur, &badan)
	r.Header.Set("Content-Type", m.FormDataContentType())
	ctx := inti.DenganAksesMenu(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: akun}), []string{handlers.KodeMenu})
	ctx = inti.DenganMenuLihat(ctx, lihat)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	return w
}

func TestRuteLampiran(t *testing.T) {
	g := tiruan.Contoh()
	g.Header["BDX-UJI.2"] = models.Header{BdxID: "BDX-UJI.2", UserInput: "UJI-MAKER", Position: "UJI-MAKER", Type: "PREMIUM",
		TypeBusiness: "FIRE"}
	p := tiruan.BaruPenyimpanan()
	h := routerLampiran(g, p)
	dasar := handlers.Prefix + "/berkas/BDX-UJI.2/lampiran"

	if w := unggahBerkas(h, dasar+"/00001", "a.pdf", []byte("ISI PDF"), "UJI-MAKER", []string{handlers.KodeMenu}); w.Code != 403 {
		t.Errorf("menu View only mengunggah: %d %s", w.Code, w.Body.String())
	}
	if w := unggahBerkas(h, dasar+"/00001", "", nil, "UJI-MAKER", nil); w.Code != 400 || !strings.Contains(w.Body.String(), "No file attached") {
		t.Errorf("tanpa berkas: %d %s", w.Code, w.Body.String())
	}
	w := unggahBerkas(h, dasar+"/00001", "Nota SOA.pdf", []byte("ISI PDF"), "UJI-MAKER", nil)
	if w.Code != 200 || strings.Contains(w.Body.String(), "UJI-IMG") {
		t.Fatalf("unggah: %d %s (IMAGEID tidak dikirim ke layar)", w.Code, w.Body.String())
	}
	var b models.Lampiran
	_ = json.Unmarshal(w.Body.Bytes(), &b)
	if b.FileName != "Nota SOA.pdf" || b.Kategori != "UJI SOA" || b.Ekstensi != "pdf" {
		t.Errorf("lampiran %+v", b)
	}

	if w := kirim(h, "GET", dasar, "", "UJI-CHK", []string{models.PeranChecker}); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `{"id":"00001","nama":"UJI SOA","cacah":1}`) {
		t.Errorf("kategori: %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", dasar+"/00001", "", "UJI-CHK", nil); w.Code != 200 || !strings.Contains(w.Body.String(), `"fileName":"Nota SOA.pdf"`) {
		t.Errorf("daftar: %d %s", w.Code, w.Body.String())
	}
	w = kirim(h, "GET", dasar+"/00001/"+b.ID+"/isi", "", "UJI-CHK", nil)
	if w.Code != 200 || w.Body.String() != "ISI PDF" || w.Header().Get("Content-Type") != "application/pdf" ||
		!strings.Contains(w.Header().Get("Content-Disposition"), "Nota SOA.pdf") {
		t.Errorf("unduh: %d %q %v", w.Code, w.Body.String(), w.Header())
	}
	if w := kirim(h, "GET", dasar+"/00001/"+b.ID+"/office", "", "UJI-CHK", nil); w.Code != 422 {
		t.Errorf("office untuk pdf: %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "GET", dasar+"/00001/TIDAK-ADA/isi", "", "UJI-CHK", nil); w.Code != 404 {
		t.Errorf("lampiran tidak ada: %d", w.Code)
	}
	if w := kirim(h, "POST", dasar+"/00001/"+b.ID+"/hapus", "", "UJI-CHK", nil); w.Code != 403 {
		t.Errorf("checker menghapus: %d %s", w.Code, w.Body.String())
	}
	if w := kirim(h, "POST", dasar+"/00001/"+b.ID+"/hapus", "", "UJI-MAKER", nil); w.Code != 200 || len(g.Lampiran) != 0 {
		t.Errorf("hapus: %d %s", w.Code, w.Body.String())
	}

	for galatnya, kode := range map[error]int{
		fmt.Errorf("%w: %q is empty or of a type that storage does not accept", penyimpanan.ErrBerkasDitolak, "a.exe"): 422,
		fmt.Errorf("%w: upload answered status 500", penyimpanan.ErrStorageGagal):                                      502,
		fmt.Errorf("%w: APPNAME (T_FOLDER_IMAGE) is not available", penyimpanan.ErrStorageBelumSiap):                   503,
	} {
		p.Gagal = galatnya
		w := unggahBerkas(h, dasar+"/00001", "a.pdf", []byte("ISI PDF"), "UJI-MAKER", nil)
		if w.Code != kode || strings.Contains(w.Body.String(), "penyimpanan:") {
			t.Errorf("%v: %d %s", galatnya, w.Code, w.Body.String())
		}
	}
}
