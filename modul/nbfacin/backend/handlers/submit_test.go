package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

type submitTiruanH struct{ tulis []models.TulisanSubmit }

func (s *submitTiruanH) KonteksSubmit(_ context.Context, id string) (models.KonteksSubmit, error) {
	if id != "UJI-NB-1" {
		return models.KonteksSubmit{}, repository.ErrKasusTidakAda
	}
	return models.KonteksSubmit{Position: "Offer", FlagOnGoingPolicy: "0", MarketingName: "UJI PEMASAR", TeamGroup: "1"}, nil
}
func (s *submitTiruanH) AkunSubmit(context.Context, string) (models.AkunSubmit, error) {
	return models.AkunSubmit{Nama: "UJI PENGGUNA"}, nil
}
func (s *submitTiruanH) RiwayatAkseptasi(context.Context, string) ([]models.RiwayatAkseptasi, error) {
	return []models.RiwayatAkseptasi{{Tanggal: "07-10-2026 22:04", PIC: "UJI PENGGUNA", Approval: "Reject", Status: "Offer", Comment: "UJI"}}, nil
}
func (s *submitTiruanH) OpsiBanding(context.Context, string, string) ([]string, error) {
	return []string{"JUW_B"}, nil
}
func (s *submitTiruanH) SebabDecline(context.Context) ([]models.BarisCauseOfDecline, error) {
	return []models.BarisCauseOfDecline{{ID: "1", TermCondition: "Rate"}, {ID: "6", TermCondition: "Period Of Insurance (POI)"},
		{ID: "15", TermCondition: "Backdated", Description: "Period Of Insurance (POI)"},
		{ID: "34", TermCondition: "Accumulation / Capacity", Description: "Accumulation / Capacity", TermCondition2: "Accumulation / Capacity"}}, nil
}
func (s *submitTiruanH) CekRW(context.Context, string, string, string, string) (models.HasilCekRW, error) {
	return models.HasilCekRW{}, nil
}
func (s *submitTiruanH) KategoriOkupasi(context.Context, []string) (map[string]string, error) {
	return map[string]string{}, nil
}
func (s *submitTiruanH) TulisSubmit(_ context.Context, _ *db.Tx, t models.TulisanSubmit) error {
	s.tulis = append(s.tulis, t)
	return nil
}

func mintaPeran(t *testing.T, svc *services.Service, metode, jalur, badan, pelaku, peran string) (int, string) {
	t.Helper()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, true)
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", pelaku)
	r.Header.Set("X-Peran", peran)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code, strings.TrimSpace(w.Body.String())
}

func kunciJSON(t *testing.T, s string) string {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("%v %s", err, s)
	}
	var k []string
	for n := range m {
		k = append(k, n)
	}
	sort.Strings(k)
	return strings.Join(k, ",")
}

// TestSubmitHandler - bentuk persis kontrak (TampilanSubmit / HasilSubmit), Reject berhasil, periksa, badan ketat 400, 401,
// 403, 422 {galat, pesan, position, nbStatus}.
func TestSubmitHandler(t *testing.T) {
	st := &submitTiruanH{}
	svc := services.Baru(nil).DenganSubmit(st).DenganTransaksi(tanpaOracle)
	jalur := "/api/nbfacin/kasus/UJI-NB-1/submit"
	kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || kunciJSON(t, isi) != "alasan,bolehSubmit,opsiBanding,pctRnmOtomatis,riwayat,sebabNonTc,sebabTc" ||
		!strings.Contains(isi, `"riwayat":[{"tanggal":"07-10-2026 22:04","pic":"UJI PENGGUNA","approval":"Reject","status":"Offer","comment":"UJI"}]`) ||
		!strings.Contains(isi, `"sebabTc":[{"termCondition":"Rate","opsi":[]},{"termCondition":"Period Of Insurance (POI)","opsi":["Backdated"]}]`) ||
		!strings.Contains(isi, `"sebabNonTc":["Accumulation / Capacity"]`) || !strings.Contains(isi, `"pctRnmOtomatis":{"rate":"","riComm":""}`) ||
		!strings.Contains(isi, `"bolehSubmit":true`) {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	kode, isi = mintaPeran(t, svc, "POST", jalur, `{"emailType":"2","comment":"UJI tolak","bandingTo":"","sebabTc":[],"sebabNonTc":""}`, "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || isi != `{"pesan":[],"position":"Offer","nbStatus":"NB IS IN UJI PEMASAR'S OPPORTUNITY"}` || len(st.tulis) != 1 {
		t.Fatalf("Reject: %d %s", kode, isi)
	}
	if kode, isi := mintaPeran(t, svc, "POST", jalur+"/periksa", `{"emailType":"4"}`, "UJI-MKT", "ReasFacInMarketing"); kode != 200 ||
		isi != `{"pesan":[],"position":"","nbStatus":""}` {
		t.Errorf("periksa: %d %s", kode, isi)
	}
	for nama, c := range map[string]struct {
		badan, pelaku, peran string
		mau                  int
	}{
		"kunci asing": {`{"emailType":"2","rasio":1}`, "UJI-MKT", "ReasFacInMarketing", 400},
		"401":         {`{"emailType":"2","comment":"x"}`, "", "", 401},
		"401 dulu":    {`bukan json`, "", "", 401},
		"403":         {`{"emailType":"2","comment":"x"}`, "UJI-UW", "ReasFacInUnderwriting", 403},
		"400 ask":     {`{"emailType":"3","comment":"x"}`, "UJI-MKT", "ReasFacInMarketing", 400},
	} {
		if kode, isi := mintaPeran(t, svc, "POST", jalur, c.badan, c.pelaku, c.peran); kode != c.mau {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	w := httptest.NewRecorder()
	tulisGalatSubmit(w, services.GalatProteksi{Pesan: []string{"Scoring Risk Cannot Empty !!", "List Payment can't be null!"}})
	if w.Code != 422 || strings.TrimSpace(w.Body.String()) != `{"galat":"Scoring Risk Cannot Empty !!; List Payment can't be null!",`+
		`"pesan":["Scoring Risk Cannot Empty !!","List Payment can't be null!"],"position":"","nbStatus":""}` {
		t.Errorf("422: %d %s", w.Code, w.Body.String())
	}
}
