package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// kasusTiruan - satu case UJI; SimpanGeneral mencatat isian dan menyimpannya.
type kasusTiruan struct {
	k   *models.Kasus
	err error
}

var errTidakAdaUji = errors.New("UJI tidak ada")

func (c kasusTiruan) BacaKasus(_ context.Context, id string) (models.Kasus, error) {
	if c.err != nil {
		return models.Kasus{}, c.err
	}
	if c.k == nil || id != c.k.CaseID {
		return models.Kasus{}, errTidakAdaUji
	}
	return *c.k, nil
}

func (c kasusTiruan) SimpanGeneral(_ context.Context, _ *db.Tx, _ string, g models.General) error {
	c.k.General = g
	return nil
}

// kasusCedingTiruan - seperti kasusTiruan, menyalin cedingIds ke cedingList (nama = kode).
type kasusCedingTiruan struct{ kasusTiruan }

func (c kasusCedingTiruan) SimpanGeneral(ctx context.Context, tx *db.Tx, id string, g models.General) error {
	for _, k := range g.CedingIDs {
		g.CedingList = append(g.CedingList, models.Ceding{ID: k, Name: k})
	}
	return c.kasusTiruan.SimpanGeneral(ctx, tx, id, g)
}

func layananKasusCeding() *services.Service {
	k := &models.Kasus{CaseID: "UJI-NB-1"}
	return services.Baru(nil).DenganKasus(kasusCedingTiruan{kasusTiruan{k: k}}).DenganTransaksi(tanpaOracle)
}

func layananKasus() *services.Service {
	k := &models.Kasus{CaseID: "UJI-NB-1", Position: "Offer", StatusWork: "Pending-Policy", InsuredName: "UJI TERTANGGUNG",
		Opportunity: models.Opportunity{BusinessProspectName: "UJI Prospek"}, General: models.General{OfferingDate: "20261001"}}
	return services.Baru(nil).DenganKasus(kasusTiruan{k: k}).DenganTransaksi(tanpaOracle).
		DenganJam(func() time.Time { return time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC) })
}

func minta(t *testing.T, svc *services.Service, metode, jalur, badan, pelaku string) (int, string) {
	t.Helper()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	if pelaku != "" {
		r.Header.Set("X-Pelaku", pelaku)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

// TestBacaKasus - GET /api/nbfacin/kasus/{caseId} (tiket 31): bentuk jawaban persis
// (kunci tingkat atas, 14 kunci opportunity, 13 kunci general), 404, 503.
func TestBacaKasus(t *testing.T) {
	kode, isi := minta(t, layananKasus(), "GET", "/api/nbfacin/kasus/UJI-NB-1", "", "")
	var j map[string]json.RawMessage
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 {
		t.Fatalf("%d %s", kode, isi)
	}
	kunci := func(raw json.RawMessage) []string {
		var m map[string]any
		_ = json.Unmarshal(raw, &m)
		var k []string
		for n := range m {
			k = append(k, n)
		}
		return k
	}
	all, _ := json.Marshal(j)
	if n := len(kunci(all)); n != 6 || len(kunci(j["opportunity"])) != 14 || len(kunci(j["general"])) != 15 {
		t.Errorf("kunci %v / opportunity %v / general %v", kunci(all), kunci(j["opportunity"]), kunci(j["general"]))
	}
	for _, harus := range []string{`"caseId":"UJI-NB-1"`, `"position":"Offer"`, `"statusWork":"Pending-Policy"`,
		`"insuredName":"UJI TERTANGGUNG"`, `"offeringDate":"01-10-2026"`, `"businessProspectName":"UJI Prospek"`, `"beginDate":""`} {
		if !strings.Contains(isi, harus) {
			t.Errorf("jawaban tanpa %s: %s", harus, isi)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1", "", ""); kode != 503 {
		t.Errorf("tanpa DB: %d", kode)
	}
	if kode, isi := minta(t, services.Baru(nil).DenganKasus(kasusTiruan{err: errors.New("ORA-UJI rincian rahasia")}).DenganTransaksi(tanpaOracle),
		"GET", "/api/nbfacin/kasus/UJI-NB-2", "", ""); kode != 500 || strings.Contains(isi, "rahasia") {
		t.Errorf("galat repository: %d %s", kode, isi)
	}
}

// TestSimpanGeneral - PUT .../general: 200 kasus baru; 400 JSON/isian; 401 tanpa
// identitas; 404 bila case tidak ada (lewat services.ErrKasusTidakAda).
func TestSimpanGeneral(t *testing.T) {
	svc := layananKasus()
	badan := `{"reffNumber":"UJI-REF","qqName":"","beginDate":"01-11-2026","offeringDate":"","endDate":"",` +
		`"policyType":"Master Policy","marketingId":"UJI-MO","day":"366","typeFacultative":"Facultative In"}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/general", badan, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"reffNumber":"UJI-REF"`) || !strings.Contains(isi, `"beginDate":"01-11-2026"`) ||
		!strings.Contains(isi, `"offeringDate":"03-10-2026"`) || !strings.Contains(isi, `"policyType":"Master Policy"`) {
		t.Fatalf("%d %s", kode, isi)
	}
	for _, u := range []struct {
		nama, jalur, badan, pelaku string
		kode                       int
	}{
		{"tanpa identitas", "/api/nbfacin/kasus/UJI-NB-1/general", badan, "", 401},
		{"JSON rusak", "/api/nbfacin/kasus/UJI-NB-1/general", `{"day":`, "UJI-USER", 400},
		{"tanggal tak sah", "/api/nbfacin/kasus/UJI-NB-1/general", `{"endDate":"31-02-2026"}`, "UJI-USER", 400},
	} {
		if kode, isi := minta(t, svc, "PUT", u.jalur, u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s: %d %s, mau %d", u.nama, kode, isi, u.kode)
		}
	}
}

// TestKasusTidakAda404 - galat repository "tidak ada" sampai ke 404 lewat services.
func TestKasusTidakAda404(t *testing.T) {
	svc := services.Baru(nil).DenganKasus(kasusTiruan{err: services.ErrKasusTidakAda}).DenganTransaksi(tanpaOracle)
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-2", "", ""); kode != 404 {
		t.Errorf("GET: %d, mau 404", kode)
	}
}

type marketingTiruan struct {
	baris []models.MarketingOfficer
	err   error
}

func (m marketingTiruan) DaftarMarketingOfficer(context.Context) ([]models.MarketingOfficer, error) {
	return m.baris, m.err
}

// TestMarketingOfficer - GET /api/nbfacin/marketing-officer: {"baris":[{"id","nama"}]}
// urutan dipertahankan, larik walau kosong, 503 tanpa DB.
func TestMarketingOfficer(t *testing.T) {
	svc := services.Baru(nil).DenganMarketingOfficer(marketingTiruan{baris: []models.MarketingOfficer{{ID: "UJI-2", Nama: "UJI B"}, {ID: "UJI-1"}}})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/marketing-officer", "", ""); kode != 200 ||
		isi != `{"baris":[{"id":"UJI-2","nama":"UJI B"},{"id":"UJI-1","nama":""}]}` {
		t.Errorf("%d %s", kode, isi)
	}
	if kode, isi := minta(t, services.Baru(nil).DenganMarketingOfficer(marketingTiruan{}), "GET", "/api/nbfacin/marketing-officer", "", ""); kode != 200 || isi != `{"baris":[]}` {
		t.Errorf("kosong: %d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/marketing-officer", "", ""); kode != 503 {
		t.Errorf("tanpa DB: %d", kode)
	}
}
