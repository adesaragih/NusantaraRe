//go:build db

// Seam HTTP reinsurer (tiket 05) terhadap skema uji Oracle NYATA - presisi
// desimal share dan komisi hanya terbukti benar pada basis data sungguhan.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"nusantarare/uji/skemauji"
)

func TestReinsurerKombinasiLingkaranPenuh(t *testing.T) {
	u, bersihkan := serverTCO(t)
	defer bersihkan()
	u.isiJenis([]skemauji.JenisReasuransiUji{{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1", Flag: "active"}})
	u.isiAgent([]skemauji.AgentUji{
		{ID: "UJI-R1", ClientName: "UJI REAS SATU", ClientID: "UJI-C1", StatusActive: "1"},
		{ID: "UJI-R2", ClientName: "UJI REAS DUA", ClientID: "UJI-C2", StatusActive: "1"},
		{ID: "UJI-R3", ClientName: "UJI REAS NONAKTIF", ClientID: "UJI-C3", StatusActive: "0"},
	})
	kode, badan := u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun", badanTahun("2026-01-01", "2026-12-31"), true)
	if kode != http.StatusOK {
		t.Fatalf("POST tahun: %d %s", kode, badan)
	}
	var tahun tahunJSON
	_ = json.Unmarshal([]byte(badan), &tahun)
	kode, badan = u.minta(t, http.MethodPost, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kontrak",
		map[string]string{"reinsTypeId": "10003", "treatyStartDate": "2026-01-01", "treatyEndDate": "2027-01-01"}, true)
	if kode != http.StatusOK {
		t.Fatalf("POST kontrak: %d %s", kode, badan)
	}
	var k kontrakJSON
	_ = json.Unmarshal([]byte(badan), &k)
	dasar := "/api/treaty-contract-out/tahun/" + tahun.ID + "/kontrak/" + k.ID + "/reinsurer"

	kode, badan = u.minta(t, http.MethodGet, "/api/treaty-contract-out/reinsurer-master?cari=reas", nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"total":2`) || strings.Contains(badan, "NONAKTIF") {
		t.Errorf("pemilih: %d %s", kode, badan)
	}

	kode, badan = u.minta(t, http.MethodPost, dasar,
		map[string]string{"reinsurerId": "UJI-R1", "pctShare": "33,33333333", "ricomm": "12,5", "stdRating": "A"}, true)
	if kode != http.StatusOK {
		t.Fatalf("POST reinsurer: %d %s", kode, badan)
	}
	var h struct {
		Reinsurer struct {
			ID, TreatyYear, TreatyGroupID, ReinsTypeID, Name, ClientID, PctShare, Ricomm string
		}
		TotalShare string
	}
	_ = json.Unmarshal([]byte(badan), &h)
	if len(h.Reinsurer.ID) != 7 || h.Reinsurer.TreatyYear != "2026" || h.Reinsurer.TreatyGroupID != "10001" ||
		h.Reinsurer.ReinsTypeID != "10003" || h.Reinsurer.Name != "UJI REAS SATU" || h.Reinsurer.ClientID != "UJI-C1" ||
		h.Reinsurer.PctShare != "33.33333333" || h.Reinsurer.Ricomm != "12.5" {
		t.Fatalf("reinsurer baru: %+v", h)
	}
	for _, s := range []string{"33.33333333", "33.33333334"} {
		if kode, badan := u.minta(t, http.MethodPost, dasar,
			map[string]string{"reinsurerId": "UJI-R2", "pctShare": s, "ricomm": "0"}, true); kode != http.StatusOK {
			t.Fatalf("POST %s: %d %s", s, kode, badan)
		}
	}
	// Presisi dari Oracle: tepat 100, bukan 99,99999999... atau 100,00000001.
	kode, badan = u.minta(t, http.MethodGet, dasar, nil, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"totalShare":"100.00000000"`) || !strings.Contains(badan, `"total":3`) {
		t.Errorf("daftar: %d %s", kode, badan)
	}
	// Total > 100 ditolak; reinsurer nonaktif dan kombinasi tanpa kontrak ditolak.
	if kode, badan := u.minta(t, http.MethodPost, dasar,
		map[string]string{"reinsurerId": "UJI-R1", "pctShare": "0.00000001", "ricomm": "0"}, true); kode != http.StatusUnprocessableEntity ||
		!strings.Contains(badan, "Persentase tidak boleh lebih dari 100!") {
		t.Errorf("total > 100: %d %s", kode, badan)
	}
	if kode, _ := u.minta(t, http.MethodPost, dasar,
		map[string]string{"reinsurerId": "UJI-R3", "pctShare": "0", "ricomm": "0"}, true); kode != http.StatusUnprocessableEntity {
		t.Errorf("reinsurer nonaktif: %d", kode)
	}
	if kode, _ := u.minta(t, http.MethodGet, "/api/treaty-contract-out/tahun/"+tahun.ID+"/kontrak/1999999/reinsurer", nil, true); kode != http.StatusNotFound {
		t.Errorf("kontrak tidak ada: %d", kode)
	}
	// Perbarui: turunkan baris pertama; share lamanya tidak dihitung dua kali.
	kode, badan = u.minta(t, http.MethodPut, dasar+"/"+h.Reinsurer.ID,
		map[string]string{"reinsurerId": "UJI-R1", "pctShare": "30", "ricomm": "10"}, true)
	if kode != http.StatusOK || !strings.Contains(badan, `"totalShare":"96.66666667"`) {
		t.Errorf("PUT: %d %s", kode, badan)
	}
}
