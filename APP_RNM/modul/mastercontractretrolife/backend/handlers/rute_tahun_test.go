package handlers_test

// Rute simpan tahun treaty (paket 2): POST baru tanpa id, PUT /{id}; tahun
// treaty ABADI - nol rute hapus.

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
)

func (u *uji) kirim(t *testing.T, metode, jalur, badan string) (int, string) {
	t.Helper()
	req, err := http.NewRequest(metode, u.srv.URL+jalur, bytes.NewBufferString(badan))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Pelaku", "UJI-PELAKU")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

const tahunBaru = `{"treatyYear":"2027","underwritingYear":"2027","startDate":"2027-01-01","endDate":"2027-12-31"}`

func TestPostTahunBaruIDDariServer(t *testing.T) {
	u := server(t, true)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/tahun", tahunBaru)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"1000044"`) || !strings.Contains(badan, `"userId":"UJI-PELAKU"`) {
		t.Fatalf("POST tahun: %d %s", kode, badan)
	}
}

func TestPostTahunBerIDDitolak400(t *testing.T) {
	u := server(t, true)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/tahun", `{"id":"UJI-X","treatyYear":"2027"}`)
	if kode != http.StatusBadRequest || !strings.Contains(badan, "server assigns") {
		t.Errorf("POST ber-id: %d %s", kode, badan)
	}
}

func TestPostTahunKosong422PesanVerbatim(t *testing.T) {
	u := server(t, true)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/tahun", `{"treatyYear":"2027"}`)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, `"galat":"Value cannot be empty."`) {
		t.Errorf("POST kosong: %d %s", kode, badan)
	}
}

func TestPutTahunIDJalurDanBadanBeda400(t *testing.T) {
	u := server(t, true)
	kode, _ := u.kirim(t, "PUT", handlers.Prefix+"/tahun/UJI-T1", `{"id":"UJI-LAIN","treatyYear":"2026"}`)
	if kode != http.StatusBadRequest {
		t.Errorf("PUT id beda: %d", kode)
	}
	kode, badan := u.kirim(t, "PUT", handlers.Prefix+"/tahun/UJI-T1",
		`{"treatyYear":"2026","underwritingYear":"2025","startDate":"2026-01-01","endDate":"2026-12-31"}`)
	if kode != http.StatusOK || !strings.Contains(badan, `"underwritingYear":"2025"`) {
		t.Errorf("PUT tahun: %d %s", kode, badan)
	}
}

func TestBadanBukanJSON400(t *testing.T) {
	kode, _ := server(t, true).kirim(t, "POST", handlers.Prefix+"/tahun", `{bukan json`)
	if kode != http.StatusBadRequest {
		t.Errorf("badan rusak: %d", kode)
	}
}

// ⛔ Tahun treaty ABADI (spec AC 2, `[fakta bisnis — work owner]`): nol
// rute DELETE atas tahun - di mux maupun di teks sumber handlers.
func TestTahunTreatyTidakDapatDihapus(t *testing.T) {
	kode, _ := server(t, true).kirim(t, "DELETE", handlers.Prefix+"/tahun/UJI-T1", "")
	if kode != http.StatusMethodNotAllowed && kode != http.StatusNotFound {
		t.Errorf("DELETE tahun: %d, mau 405/404", kode)
	}
	pola := regexp.MustCompile(`"DELETE "\s*\+\s*Prefix\s*\+\s*"/tahun`)
	berkas, _ := os.ReadDir(".")
	for _, b := range berkas {
		if strings.HasSuffix(b.Name(), "_test.go") || !strings.HasSuffix(b.Name(), ".go") {
			continue
		}
		isi, _ := os.ReadFile(b.Name())
		if pola.Match(isi) {
			t.Errorf("%s mendaftarkan DELETE tahun treaty", b.Name())
		}
	}
}
