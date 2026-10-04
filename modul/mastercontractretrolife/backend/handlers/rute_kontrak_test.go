package handlers_test

// Rute simpan kontrak (paket 3).

import (
	"net/http"
	"strings"
	"testing"

	"nusantarare/modul/mastercontractretrolife/backend/handlers"
)

// REINS TYPE OR (10200): QS (10196) sudah dipakai UJI-K1 di UJI-T1, dan satu jenis
// hanya boleh satu kontrak per tahun (04-10-2026).
const kontrakBaru = `{"reinsTypeId":"10200","bIdr":"0","idr":"1500000000.123456789","bUsd":"0","usd":"250000.5"}`

func TestPostKontrakSelisihDihitungUangTeks(t *testing.T) {
	u := server(t, true)
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/tahun/UJI-T1/kontrak", kontrakBaru)
	if kode != http.StatusOK {
		t.Fatalf("POST kontrak: %d %s", kode, badan)
	}
	for _, w := range []string{`"idrSelisih":"1500000000.123456789"`, `"usd":"250000.5"`, `"usdSelisih":"250000.5"`,
		`"reinsTypeName":"OR"`, `"treatyStartDate":"2026-01-01"`} {
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
		`{"reinsTypeId":"10196","bIdr":"10","idr":"5","bUsd":"0","usd":"10"}`)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, "MINIMUM LIMIT (IDR)") {
		t.Errorf("batas terbalik: %d %s", kode, badan)
	}
}

// Satu REINS TYPE satu kontrak per tahun treaty (04-10-2026).
func TestPostKontrakJenisGandaDitolak(t *testing.T) {
	u := server(t, true)
	ganda := `{"reinsTypeId":"10196","bIdr":"0","idr":"100","bUsd":"0","usd":"10"}`
	kode, badan := u.kirim(t, "POST", handlers.Prefix+"/tahun/UJI-T1/kontrak", ganda)
	if kode == http.StatusOK || !strings.Contains(badan, "REINS TYPE QS is already used in this treaty year.") {
		t.Errorf("kontrak QS kedua di UJI-T1: %d %s", kode, badan)
	}
}
