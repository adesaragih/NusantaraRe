package handlers_test

// Rute simpan kontrak (paket 3).

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
)

const kontrakBaru = `{"reinsTypeId":"10196","bIdr":"0","idr":"1500000000.123456789","bUsd":"0","usd":""}`

func TestPostKontrakSelisihDihitungUangTeks(t *testing.T) {
	u := server(t, true)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/tahun/UJI-T1/kontrak", kontrakBaru)
	if kode != http.StatusOK {
		t.Fatalf("POST kontrak: %d %s", kode, badan)
	}
	for _, w := range []string{`"idrSelisih":"1500000000.123456789"`, `"usd":""`, `"usdSelisih":""`,
		`"reinsTypeName":"QS"`, `"treatyStartDate":"2026-01-01"`} {
		if !strings.Contains(badan, w) {
			t.Errorf("jawaban tanpa %s: %s", w, badan)
		}
	}
}

func TestPostKontrakBerIDDitolak(t *testing.T) {
	kode, _ := server(t, true).kirim(t, "POST", handlers.Prefix+"/tahun/UJI-T1/kontrak", `{"id":"UJI-X"}`)
	if kode != http.StatusBadRequest {
		t.Errorf("POST kontrak ber-id: %d", kode)
	}
}

func TestPutKontrakKosong422Verbatim(t *testing.T) {
	kode, badan := server(t, true).kirim(t, "PUT", handlers.Prefix+"/kontrak/UJI-K1", `{"reinsTypeId":"10196"}`)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, `"galat":"All value cannot be empty."`) {
		t.Errorf("PUT kontrak kosong: %d %s", kode, badan)
	}
}

func TestPutKontrakBatasTerbalik422(t *testing.T) {
	kode, badan := server(t, true).kirim(t, "PUT", handlers.Prefix+"/kontrak/UJI-K1",
		`{"reinsTypeId":"10196","bIdr":"10","idr":"5","bUsd":"0"}`)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, "MINIMUM LIMIT (IDR)") {
		t.Errorf("batas terbalik: %d %s", kode, badan)
	}
}
