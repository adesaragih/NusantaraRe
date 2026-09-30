package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/inti/backend/galat"
	"nusantarare/modul/treatycontractout/backend/services"
)

// Keputusan work owner 30-09-2026: di form Add tahun treaty, Start Date
// mengisi End Date (+1 tahun, aturan kontrak) - dihitung server, tanpa basis data.
func TestAkhirBawaanTahunBaru(t *testing.T) {
	router := Router(services.New(nil), true)
	minta := func(kueri string, pelaku bool) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		q := httptest.NewRequest(http.MethodGet, "/api/treaty-contract-out/tahun/akhir-bawaan"+kueri, nil)
		if pelaku {
			q.Header.Set("X-Pelaku", "UJI-ADMIN")
		}
		router.ServeHTTP(w, q)
		return w
	}
	for mulai, mau := range map[string]string{
		"2026-06-01": "2027-06-01",
		"2024-02-29": "2025-02-28", // ADD_MONTHS: 29 Feb -> 28 Feb
		"2025-02-28": "2026-02-28",
	} {
		w := minta("?mulai="+mulai, true)
		var isi map[string]string
		if w.Code != http.StatusOK || json.NewDecoder(w.Body).Decode(&isi) != nil || isi["endDate"] != mau {
			t.Errorf("mulai %s: kode %d, endDate %q, mau %q", mulai, w.Code, isi["endDate"], mau)
		}
	}
	if w := minta("", true); w.Code != http.StatusBadRequest {
		t.Errorf("tanpa mulai: kode %d, mau 400", w.Code)
	}
	if w := minta("?mulai=2026-06-01", false); w.Code != http.StatusUnauthorized {
		t.Errorf("tanpa identitas: kode %d, mau 401", w.Code)
	}
}

// Keputusan work owner 30-09-2026: 400 berbahasa Inggris - sentinel `inti/`
// diganti (`services.TeksInggrisTCO`, diuji di services).
func TestPesanTCOBerbahasaInggris(t *testing.T) {
	w := httptest.NewRecorder()
	jawabGalatTreatyContractOut(w, fmt.Errorf("%w: start date (mulai) is required", galat.ErrPermintaanTidakSah))
	if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "services: invalid request: start date") {
		t.Errorf("400: kode %d %s", w.Code, w.Body.String())
	}
}
