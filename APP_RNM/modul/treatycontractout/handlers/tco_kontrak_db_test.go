//go:build db

// Seam HTTP kontrak treaty (tiket 04) terhadap skema uji Oracle NYATA.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/uji/skemauji"
)

type kontrakJSON struct {
	ID, IDTreatyYear, ReinsTypeID, ReinsTypeName, TreatyStartDate, TreatyEndDate, UserID, TglUpdate string
}

// AC 5-10, 41 dan tanggal akhir bawaan: tahun -> kontrak baru ('1'+6 digit,
// nama jenis dari master), daftar, perbarui tanpa baris baru, 409 dobel,
// 422 tahun mulai beda / periode terbalik / jenis di luar daftar, 404 tahun
// tidak ada, jejak.
func TestKontrakTreatyLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{
		{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"},
		{ID: "10005", Note: "UJI SURPLUS", Tipe: "2", Flag: "active"},
		{ID: "10004", Note: "UJI BLACKLIST", Tipe: "1", Flag: "active"},
	})
	kode, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	if kode != http.StatusOK {
		t.Fatalf("POST tahun: %d %s", kode, badan)
	}
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	dasar := "/api/treaty-contract-out/tahun/" + tahun.ID + "/kontrak"

	kode, badan = u.minta(t, http.MethodGet, dasar+"/akhir-bawaan?mulai=2026-01-01", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"treatyEndDate":"2027-01-01"`) {
		t.Errorf("akhir bawaan: %d %s", kode, badan)
	}

	masuk := map[string]string{"reinsTypeId": "10003", "reinsTypeName": "KARANGAN KLIEN",
		"treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"}
	kode, badan = u.minta(t, http.MethodPost, dasar, masuk, true)
	if kode != http.StatusOK {
		t.Fatalf("POST kontrak: %d %s", kode, badan)
	}
	var k kontrakJSON
	_ = json.Unmarshal([]byte(badan), &k)
	if len(k.ID) != 7 || !strings.HasPrefix(k.ID, "1") || k.IDTreatyYear != tahun.ID ||
		k.ReinsTypeName != "UJI QUOTA SHARE" || k.UserID != "UJI-ADMIN" || k.TreatyEndDate != "2027-01-01" {
		t.Fatalf("kontrak baru: %+v", k)
	}

	// Dobel: jenis yang sama di tahun yang sama.
	kode, badan = u.minta(t, http.MethodPost, dasar, masuk, true)
	if kode != http.StatusConflict || !strings.Contains(badan, "Data sudah pernah di Input") || !strings.Contains(badan, k.ID) {
		t.Errorf("dobel: %d %s", kode, badan)
	}
	// Gerbang 422.
	for nama, ubah := range map[string]map[string]string{
		"tahun mulai beda": {"reinsTypeId": "10005", "treatyStartDate": "2027-01-01", "treatyEndDate": "2028-01-01"},
		"periode terbalik": {"reinsTypeId": "10005", "treatyStartDate": "2026-06-01", "treatyEndDate": "2026-01-01"},
		"jenis blacklist":  {"reinsTypeId": "10004", "treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"},
	} {
		if kode, badan := u.minta(t, http.MethodPost, dasar, ubah, true); kode != http.StatusUnprocessableEntity {
			t.Errorf("%s: %d %s", nama, kode, badan)
		}
	}
	if kode, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun/1999999/kontrak", masuk, true); kode != http.StatusNotFound {
		t.Errorf("tahun tidak ada: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, http.MethodPost, dasar, map[string]string{"id": "1000099", "reinsTypeId": "10005"}, true); kode != http.StatusBadRequest {
		t.Errorf("POST ber-id: %d", kode)
	}

	// Perbarui: jenis lain, tanpa baris baru.
	ubah := map[string]string{"reinsTypeId": "10005", "treatyStartDate": "2026-02-01", "treatyEndDate": "2027-01-31"}
	kode, badan = u.minta(t, http.MethodPut, dasar+"/"+k.ID, ubah, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"reinsTypeName":"UJI SURPLUS"`) {
		t.Fatalf("PUT: %d %s", kode, badan)
	}
	kode, badan = u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":1`) || !strings.Contains(badan, `"treatyStartDate":"2026-02-01"`) {
		t.Errorf("daftar sesudah PUT: %d %s", kode, badan)
	}
}
