package handlers_test

// Seam HTTP simpan produk (paket 3) di atas gudang tiruan.

import (
	"net/http"
	"strings"
	"testing"
)

const badanUji = `{"umum":{"productName":"UJI PRODUK","ceding":"UJI CEDING SATU","cedingId":"L0UJI1","sobName":"",
  "riComm":"0","cause":"ANY CAUSE","causeId":"100004"},
  "inward":{"policyHolder":"UJI-ORG-1","policyHolderName":"UJI PEMEGANG","begin":"2026-03-01"}}`

func TestHTTPSimpanBaruLaluUbah(t *testing.T) {
	u := server(t, true)
	isiMaster(u.g)
	kode, badan := u.minta(t, "POST", pre+"/produk", badanUji, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"100044"`) || !strings.Contains(badan, `"createOp":"UJI-PELAKU"`) {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, "PUT", pre+"/produk/100044", strings.Replace(badanUji, "UJI PRODUK", "UJI UBAH", 1), true)
	if kode != http.StatusOK || !strings.Contains(badan, `"productName":"UJI UBAH"`) || len(u.g.Umum) != 3 {
		t.Errorf("PUT: %d %s (%d produk)", kode, badan, len(u.g.Umum))
	}
}

func TestHTTPSimpanGalatBerkalimat(t *testing.T) {
	u := server(t, true)
	isiMaster(u.g)
	if kode, badan := u.minta(t, "POST", pre+"/produk", `{"id":"100001"}`, true); kode != http.StatusBadRequest ||
		!strings.Contains(badan, "server assigns it") {
		t.Errorf("ID dari klien: %d %s", kode, badan)
	}
	if kode, badan := u.minta(t, "PUT", pre+"/produk/100044", `{"id":"100045"}`, true); kode != http.StatusBadRequest ||
		!strings.Contains(badan, "differs") {
		t.Errorf("id badan ≠ jalur: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, "POST", pre+"/produk", `{`, true); kode != http.StatusBadRequest {
		t.Errorf("JSON rusak: %d", kode)
	}
	kode, badan := u.minta(t, "POST", pre+"/produk", strings.Replace(badanUji, `"riComm":"0"`, `"riComm":"abc"`, 1), true)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(badan, "Deduction (%)") {
		t.Errorf("angka tak sah 422 berlabel: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, "PUT", pre+"/produk/100999", badanUji, true); kode != http.StatusNotFound {
		t.Errorf("ubah produk tak ada: %d", kode)
	}
	if kode, _ := u.minta(t, "POST", pre+"/produk", badanUji, false); kode != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: %d", kode)
	}
}
