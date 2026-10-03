//go:build db

// Seam HTTP security (tiket 06) terhadap skema uji Oracle NYATA - keutuhan
// rujukan sesudah PK surrogate, kaskade FK, dan hilangnya kunci berbasis nama
// hanya terbukti pada basis data sungguhan.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/uji/skemauji"
)

func TestSecurityLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"}})
	u.isiAgent([]skemauji.AgentUji{
		{ID: "UJI-R1", ClientName: "UJI REAS SATU", ClientID: "UJI-C1", StatusActive: "1"},
		{ID: "UJI-R2", ClientName: "UJI REAS DUA", ClientID: "UJI-C2", StatusActive: "1"},
		{ID: "UJI-R3", ClientName: "UJI REAS TIGA", ClientID: "UJI-C3", StatusActive: "1"},
	})
	_, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	_, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kontrak",
		map[string]string{"reinsTypeId": "10003", "treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"}, true)
	var k kontrakJSON
	_ = json.Unmarshal([]byte(badan), &k)
	reas := "/api/treaty-contract-out/tahun/" + tahun.ID + "/kontrak/" + k.ID + "/reinsurer"
	_, badan = u.minta(t, http.MethodPost, reas, map[string]string{"reinsurerId": "UJI-R1", "pctShare": "40", "ricomm": "0"}, true)
	var r struct{ Reinsurer struct{ ID string } }
	_ = json.Unmarshal([]byte(badan), &r)
	dasar := reas + "/" + r.Reinsurer.ID + "/security"

	// AC 17: nama dari master, share desimal persis. tco4: MTREATYSECURITY tanpa
	// identitas - "ID" = REAS_SECURITY terpangkas (kunci UpdateMTreatySecurity).
	kode, badan := u.minta(t, http.MethodPost, dasar, map[string]string{"reasSecurity": "UJI-R2", "pctShare": "33,33333333"}, true)
	if kode != http.StatusOK {
		t.Fatalf("POST security: %d %s", kode, badan)
	}
	var s struct{ ID, ReasID, ThnTreaty, ReasSecurity, ClientName, PctShare string }
	_ = json.Unmarshal([]byte(badan), &s)
	if s.ID != "UJI-R2" || s.ReasID != r.Reinsurer.ID || s.ThnTreaty != "2026" || s.ClientName != "UJI REAS DUA" ||
		s.PctShare != "33.33333333" {
		t.Fatalf("security baru: %+v", s)
	}
	// AC 19: tiga kolom [terbuka] NULL eksplisit.
	var kosong int
	if err := u.sqlDBMentah().QueryRowContext(u.ctx, `SELECT COUNT(*) FROM `+u.skema+`.MTREATYSECURITY
		 WHERE REAS_ID = :1 AND TRIM(REAS_SECURITY) = :2 AND TOP_ID IS NULL AND TP_TREATY IS NULL AND USER_ID IS NULL`,
		r.Reinsurer.ID, s.ID).Scan(&kosong); err != nil || kosong != 1 {
		t.Errorf("kolom [terbuka] tidak NULL: %d %v", kosong, err)
	}
	// User story 14: ganti nama -> baris yang SAMA ditimpa (UpdateMTreatySecurity
	// berkunci nama LAMA), kuncinya kini nama baru.
	kode, badan = u.minta(t, http.MethodPut, dasar+"/"+s.ID, map[string]string{"reasSecurity": "UJI-R3", "pctShare": "20"}, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"UJI-R3"`) || !strings.Contains(badan, "UJI REAS TIGA") {
		t.Fatalf("ganti nama: %d %s", kode, badan)
	}
	s.ID = "UJI-R3"
	// Dobel ditolak; security kedua diterima.
	if kode, _ := u.minta(t, http.MethodPost, dasar, map[string]string{"reasSecurity": "UJI-R3", "pctShare": "1"}, true); kode != http.StatusConflict {
		t.Errorf("dobel: %d", kode)
	}
	_, badan = u.minta(t, http.MethodPost, dasar, map[string]string{"reasSecurity": "UJI-R2", "pctShare": "5"}, true)
	var s2 struct{ ID string }
	_ = json.Unmarshal([]byte(badan), &s2)
	kode, badan = u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":2`) {
		t.Errorf("daftar: %d %s", kode, badan)
	}
	// Hapus satu baris; reinsurer induk utuh.
	if kode, _ := u.minta(t, http.MethodDelete, dasar+"/"+s.ID, nil, true); kode != http.StatusOK {
		t.Errorf("hapus: %d", kode)
	}
	if kode, badan := u.minta(t, http.MethodGet, reas, nil, true); kode != http.StatusOK || !strings.Contains(badan, `"total":1`) {
		t.Errorf("hapus security menyentuh reinsurer: %d %s", kode, badan)
	}
	// tco4: MTREATYSECURITY tanpa FK - kaskade reinsurer -> security dijalankan
	// layanan (DeleteFromTreatyReinsurer_Act), diuji tco_kaskade_db_test.go.
	_ = s2
}
