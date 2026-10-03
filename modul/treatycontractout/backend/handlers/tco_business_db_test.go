//go:build db

// Seam HTTP business (tiket 07) terhadap skema uji Oracle NYATA.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/uji/skemauji"
)

// AC 21-23, 63/64: kode + nama dari master, nonaktif tetap terbaca, pembaruan
// menulis seluruh medan, hapus satu tabel, jejak.
func TestBusinessKombinasiLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"}})
	u.isiBusiness([]skemauji.BusinessUji{
		{ID: "UJI-B1", Note: "UJI BISNIS SATU", BusinessGroupID: "UJI-G1"},
		{ID: "UJI-B2", Note: "UJI BISNIS DUA", BusinessGroupID: "UJI-G1"},
		{ID: "UJI-B3", Note: "UJI TANPA GRUP"},
	})
	_, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	_, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kontrak",
		map[string]string{"reinsTypeId": "10003", "treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"}, true)
	var k kontrakJSON
	_ = json.Unmarshal([]byte(badan), &k)
	if k.ID == "" {
		t.Fatalf("kontrak: %s", badan)
	}
	dasar := "/api/treaty-contract-out/tahun/" + tahun.ID + "/kontrak/" + k.ID + "/business"

	kode, badan := u.minta(t, http.MethodGet, "/api/treaty-contract-out/business-master", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":2`) || strings.Contains(badan, "TANPA GRUP") {
		t.Errorf("master: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodPost, dasar, map[string]string{"bizCode": "UJI-B1", "isActive": "1"}, true)
	if kode != http.StatusOK {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	var b struct{ ID, BizName, TreatyYearID, IsActive string }
	_ = json.Unmarshal([]byte(badan), &b)
	if len(b.ID) != 7 || b.BizName != "UJI BISNIS SATU" || b.TreatyYearID != tahun.ID {
		t.Fatalf("baru: %+v", b)
	}
	if kode, badan := u.minta(t, http.MethodPost, dasar, map[string]string{"bizCode": "UJI-B1", "isActive": "1"}, true); kode != http.StatusConflict ||
		!strings.Contains(badan, "Data has already been entered") {
		t.Errorf("dobel: %d %s", kode, badan)
	}
	// Nonaktifkan + ganti kode: SELURUH medan tersimpan; baris tetap terbaca.
	kode, badan = u.minta(t, http.MethodPut, dasar+"/"+b.ID, map[string]string{"bizCode": "UJI-B2", "isActive": "0"}, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"bizName":"UJI BISNIS DUA"`) || !strings.Contains(badan, `"aktif":false`) {
		t.Fatalf("PUT: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":1`) || !strings.Contains(badan, `"isActive":"0"`) {
		t.Errorf("daftar nonaktif: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodDelete, dasar+"/"+b.ID, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, "Data with ID "+b.ID+" successfully deleted") {
		t.Errorf("hapus: %d %s", kode, badan)
	}
}
