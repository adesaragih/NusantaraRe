package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type caseNBTiruan struct {
	err     error
	pembuat *string
}

func (c caseNBTiruan) PengenalBerikut(context.Context, *db.Tx) (string, error) {
	return "NB-900001", nil
}
func (c caseNBTiruan) SisipCase(_ context.Context, _ *db.Tx, _, pembuat string) error {
	if c.pembuat != nil {
		*c.pembuat = pembuat
	}
	return nil
}
func (c caseNBTiruan) SisipOpportunity(context.Context, *db.Tx, string, models.Opportunity) error {
	return c.err
}

func tanpaOracle(_ context.Context, fn func(*db.Tx) error) error { return fn(nil) }

const badanSah = `{"estimatedClosingDate":"31-10-2026","businessProspectName":"UJI Prospek","accountId":"UJI-AKUN",` +
	`"insuredId":"UJI-INS","groupBusinessId":"UJI-GB","groupBusiness":"UJI GRUP","classOfBusiness":"UJI KELAS",` +
	`"typeOfInward":"Facultative","typeOfFacultative":"Facultative In","phase":"Proposal","stage":"Opportunity",` +
	`"opportunitySource":"","businessStatus":"New Business","description":""}`

// kirimOpportunity - POST /api/nbfacin/opportunity; pelaku "" = tanpa identitas.
func kirimOpportunity(t *testing.T, svc *services.Service, stub bool, pelaku, sesi, badan string) (int, string) {
	t.Helper()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, stub)
	r := httptest.NewRequest("POST", "/api/nbfacin/opportunity", strings.NewReader(badan))
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	if sesi != "" {
		r = r.WithContext(inti.DenganPelakuSesi(r.Context(), inti.Pelaku{AkunID: sesi}))
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

// TestBuatOpportunity - POST /api/nbfacin/opportunity (tiket 29): 201 persis
// {"caseId":"NB-<n>"}; pembuat dari sesi login (didahulukan) atau stub X-Pelaku; 401
// tanpa identitas; 400 JSON rusak / isian tak sah; 503 tanpa DB; 500 tanpa rincian.
// Data sintetis UJI-.
func TestBuatOpportunity(t *testing.T) {
	var pembuat string
	svc := services.Baru(nil).DenganCaseNB(caseNBTiruan{pembuat: &pembuat}).DenganTransaksi(tanpaOracle)
	if kode, isi := kirimOpportunity(t, svc, false, "", "UJI-SESI", badanSah); kode != 201 || isi != `{"caseId":"NB-900001"}` || pembuat != "UJI-SESI" {
		t.Fatalf("sesi: %d %s, pembuat %q", kode, isi, pembuat)
	}
	if kode, _ := kirimOpportunity(t, svc, true, "UJI-STUB", "UJI-SESI", badanSah); kode != 201 || pembuat != "UJI-SESI" {
		t.Errorf("sesi harus mengalahkan stub: %d, pembuat %q", kode, pembuat)
	}
	if kode, _ := kirimOpportunity(t, svc, true, "UJI-STUB", "", badanSah); kode != 201 || pembuat != "UJI-STUB" {
		t.Errorf("stub: %d, pembuat %q", kode, pembuat)
	}
	for _, u := range []struct {
		nama         string
		svc          *services.Service
		stub         bool
		pelaku       string
		badan        string
		kode         int
		pesan, bukan string
	}{
		{"stub mati", svc, false, "UJI-STUB", badanSah, 401, "identitas", ""},
		{"tanpa identitas", svc, true, "", badanSah, 401, "identitas", ""},
		{"JSON rusak", svc, true, "UJI-STUB", `{"phase":`, 400, "JSON", ""},
		{"bukan teks", svc, true, "UJI-STUB", `{"phase":1}`, 400, "JSON", ""},
		{"isian tak sah", svc, true, "UJI-STUB", strings.Replace(badanSah, `"31-10-2026"`, `"31-02-2026"`, 1), 400, "estimatedClosingDate", ""},
		{"tanpa DB", services.Baru(nil), true, "UJI-STUB", badanSah, 503, "basis data", ""},
		{"galat Oracle", services.Baru(nil).DenganCaseNB(caseNBTiruan{err: errors.New("ORA-UJI rincian rahasia")}).DenganTransaksi(tanpaOracle),
			true, "UJI-STUB", badanSah, 500, "galat server", "rahasia"},
	} {
		kode, isi := kirimOpportunity(t, u.svc, u.stub, u.pelaku, "", u.badan)
		var j map[string]string
		if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != u.kode || !strings.Contains(j["galat"], u.pesan) ||
			(u.bukan != "" && strings.Contains(isi, u.bukan)) {
			t.Errorf("%s: %d %s, mau %d berisi %q", u.nama, kode, isi, u.kode, u.pesan)
		}
	}
}
