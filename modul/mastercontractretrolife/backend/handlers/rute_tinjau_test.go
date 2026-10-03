package handlers_test

// Perbaikan /code-review (01-10-2026): badan dibatasi dan dibaca SESUDAH identitas, medan wajib yang
// kosong tercatat di log server.

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestBadanTerlaluBesarDitolak413(t *testing.T) {
	u := server(t, true)
	besar := `{"treatyYear":"` + strings.Repeat("9", 70<<10) + `"}`
	kode, badan := u.kirim(t, http.MethodPost, "/api/master-contract-retro-life/tahun", besar)
	if kode != http.StatusRequestEntityTooLarge || !strings.Contains(badan, `"galat"`) {
		t.Errorf("badan 70 KiB: %d %s", kode, badan)
	}
	if len(u.g.Tahun) != 1 {
		t.Errorf("nol baris boleh lahir: %d", len(u.g.Tahun))
	}
}

func TestTanpaIdentitasDitolakSebelumBadanDibaca(t *testing.T) {
	u := server(t, true)
	req, err := http.NewRequest(http.MethodPost, u.srv.URL+"/api/master-contract-retro-life/tahun", bytes.NewBufferString("bukan json"))
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("tanpa X-Pelaku, badan rusak: %d, mau 401", res.StatusCode)
	}
}

func TestWajibIsiMencatatMedanYangKosong(t *testing.T) {
	log := tangkapLog(t)
	u := server(t, true)
	kode, _ := u.kirim(t, http.MethodPost, "/api/master-contract-retro-life/tahun", `{"treatyYear":"2026"}`)
	if kode != http.StatusUnprocessableEntity {
		t.Fatalf("wajib isi: %d", kode)
	}
	if s := log.String(); !strings.Contains(s, "UNDERWRITING YEAR") || !strings.Contains(s, "START DATE") {
		t.Errorf("log tanpa medan kosong: %q", s)
	}
}
