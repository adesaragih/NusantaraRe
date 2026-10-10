package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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
func (s *submitTiruanH) TulisTekanBinding(context.Context, *db.Tx, models.TekanBinding) error {
	return nil
}

func (s *submitTiruanH) HitungAkumulasiBerjalan(context.Context, string, string, string, string) (int, error) {
	return 0, nil
}

func (s *submitTiruanH) KasusGanda(context.Context, models.CariKasusGanda) ([]string, error) {
	return nil, nil
}

func (s *submitTiruanH) PenerimaBinding(context.Context, string, string) (string, []string, error) {
	return "", nil, nil
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
	if kode != 200 || kunciJSON(t, isi) != "alasan,bolehSubmit,commentDeductible,confirmBinding,dateValidity,daysValidity,deductibleAcceptance,deductibleBedaMataUang,opsiBanding,pctRnmOtomatis,receivedRiSlip,riwayat,sebabNonTc,sebabTc,tahap,waitingBindDate" ||
		!strings.Contains(isi, `"riwayat":[{"tanggal":"07-10-2026 22:04","pic":"UJI PENGGUNA","approval":"Reject","status":"Offer","comment":"UJI"}]`) ||
		!strings.Contains(isi, `"sebabTc":[{"termCondition":"Rate","opsi":[]},{"termCondition":"Period Of Insurance (POI)","opsi":["Backdated"]}]`) ||
		!strings.Contains(isi, `"sebabNonTc":["Accumulation / Capacity"]`) || !strings.Contains(isi, `"pctRnmOtomatis":{"rate":"","riComm":""}`) ||
		!strings.Contains(isi, `"bolehSubmit":true`) {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	kode, isi = mintaPeran(t, svc, "POST", jalur, `{"emailType":"2","comment":"UJI tolak","bandingTo":"","sebabTc":[],"sebabNonTc":""}`, "UJI-MKT", "ReasFacInMarketing")
	if kode != 200 || isi != `{"pesan":[],"position":"Offer","nbStatus":"NB IS IN UJI PEMASAR'S OPPORTUNITY","pesanBinding":[],"nomorPolis":""}` || len(st.tulis) != 1 {
		t.Fatalf("Reject: %d %s", kode, isi)
	}
	// Decline tanpa proteksi (Banding kini diproteksi, tiket 57b - layanan uji ini tanpa kasus).
	if kode, isi := mintaPeran(t, svc, "POST", jalur+"/periksa", `{"emailType":"7"}`, "UJI-MKT", "ReasFacInMarketing"); kode != 200 ||
		isi != `{"pesan":[],"position":"","nbStatus":"","pesanBinding":[],"nomorPolis":""}` {
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
		`"pesan":["Scoring Risk Cannot Empty !!","List Payment can't be null!"],"position":"","nbStatus":"","pesanBinding":[],"nomorPolis":""}` {
		t.Errorf("422: %d %s", w.Code, w.Body.String())
	}
}

// submitTiruanUW - kasus UJI-NB-1 di posisi Junior UW B (tahap uw, tiket 61).
type submitTiruanUW struct{ submitTiruanH }

func (s *submitTiruanUW) KonteksSubmit(context.Context, string) (models.KonteksSubmit, error) {
	return models.KonteksSubmit{Position: "ReasFacInJuniorUnderwritingB", FlagOnGoingPolicy: "0", MarketingName: "UJI PEMASAR",
		TeamGroup: "1", StartDateTime: "20991231T170000.000 GMT"}, nil
}

// TestSubmitUWHandler - tiket 61 lewat HTTP: GET tahap uw + validitas; POST …/submit menerima daysValidity (Reject ->
// marketing, validitas tersimpan); POST …/submit/validitas 200 / 400 / 403; kunci JSON TampilanSubmit / ValiditasSubmit dua
// arah dengan api.ts.
func TestSubmitUWHandler(t *testing.T) {
	st := &submitTiruanUW{}
	svc := services.Baru(nil).DenganSubmit(st).DenganTransaksi(tanpaOracle)
	jalur, peran := "/api/nbfacin/kasus/UJI-NB-1/submit", "ReasFacInJuniorUnderwritingB"
	kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-UW", peran)
	if kode != 200 || !strings.Contains(isi, `"tahap":"uw"`) || !strings.Contains(isi, `"daysValidity":"","dateValidity":"","waitingBindDate":""`) ||
		!strings.Contains(isi, `"bolehSubmit":true`) {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	kode, isi = mintaPeran(t, svc, "POST", jalur, `{"emailType":"2","comment":"UJI tolak","bandingTo":"","sebabTc":[],"sebabNonTc":"",`+
		`"deductibleAcceptance":false,"commentDeductible":"","daysValidity":"3"}`, "UJI-UW", peran)
	if kode != 200 || isi != `{"pesan":[],"position":"ReasFacInMarketing","nbStatus":"NB IS IN UJI PEMASAR'S INBOX","pesanBinding":[],"nomorPolis":""}` || len(st.tulis) != 1 ||
		st.tulis[0].Validitas == nil || st.tulis[0].Validitas.Days != "3" || st.tulis[0].Validitas.Comment != "UJI tolak" {
		t.Fatalf("POST: %d %s %+v", kode, isi, st.tulis)
	}
	if kode, isi := mintaPeran(t, svc, "POST", jalur+"/periksa", `{"emailType":"3","deductibleAcceptance":false,"commentDeductible":"","daysValidity":""}`,
		"UJI-UW", peran); kode != 200 || isi != `{"pesan":[],"position":"","nbStatus":"","pesanBinding":[],"nomorPolis":""}` {
		t.Errorf("periksa: %d %s", kode, isi)
	}
	kode, isi = mintaPeran(t, svc, "POST", jalur+"/validitas", `{"daysValidity":"10"}`, "UJI-UW", peran)
	if kode != 200 || kunciJSON(t, isi) != "dateValidity,daysValidity,waitingBindDate" || !strings.Contains(isi, `"daysValidity":"10"`) {
		t.Errorf("validitas edit: %d %s", kode, isi)
	}
	kode, isi = mintaPeran(t, svc, "POST", jalur+"/validitas", `{"daysValidity":""}`, "UJI-UW", peran)
	if kode != 200 || !strings.Contains(isi, `"daysValidity":"30"`) { // Begin date 2099 -> batas 30
		t.Errorf("validitas bawaan: %d %s", kode, isi)
	}
	for nama, c := range map[string]struct {
		badan, peran string
		kode         int
	}{
		"400 huruf": {`{"daysValidity":"x"}`, peran, 400}, "400 kunci asing": {`{"hari":"1"}`, peran, 400},
		"403": {`{"daysValidity":"1"}`, "ReasFacInMarketing", 403},
	} {
		if kode, isi := mintaPeran(t, svc, "POST", jalur+"/validitas", c.badan, "UJI-UW", c.peran); kode != c.kode {
			t.Errorf("validitas %s: %d %s", nama, kode, isi)
		}
	}
	for antarmuka, gob := range map[string][]string{"TampilanSubmit": kunciTagJSON(tampilanSubmitKabel{}), "ValiditasSubmit": kunciTagJSON(validitasSubmitKabel{})} {
		if ts := kunciInterfaceTS(t, antarmuka); !reflect.DeepEqual(ts, gob) {
			t.Errorf("%s: api.ts %v, handler %v", antarmuka, ts, gob)
		}
	}
	if n := len(kunciInterfaceTS(t, "TampilanSubmit")); n != 16 { // tiket 64: + confirmBinding, receivedRiSlip
		t.Errorf("instrumen: TampilanSubmit %d kunci, mau 16", n)
	}
}

// TestTekanBindingHandler - POST …/submit/binding (tiket 64): 401 tanpa sesi, 400 badan bukan {}, 409 kasus bukan tahap Binding;
// badan POST …/submit menerima SEMUA kunci api.ts IsiSubmit (termasuk receivedRiSlip, dua arah dengan pesan 400 handler).
func TestTekanBindingHandler(t *testing.T) {
	st := &submitTiruanH{}
	svc := services.Baru(nil).DenganSubmit(st).DenganTransaksi(tanpaOracle)
	jalur := "/api/nbfacin/kasus/UJI-NB-1/submit/binding"
	for nama, c := range map[string]struct {
		badan, pelaku string
		kode          int
	}{
		"401": {`{}`, "", 401}, "400": {`{"x":1}`, "UJI-MKT", 400}, "409 bukan binding": {`{}`, "UJI-MKT", 409},
	} {
		if kode, isi := mintaPeran(t, svc, "POST", jalur, c.badan, c.pelaku, "ReasFacInMarketing"); kode != c.kode {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	nilai := map[string]string{"sebabTc": "[]", "deductibleAcceptance": "false", "receivedRiSlip": "false"}
	var bagian []string
	kunci := kunciInterfaceTS(t, "IsiSubmit")
	for _, k := range kunci {
		v, ada := nilai[k]
		if !ada {
			v = `""`
		}
		bagian = append(bagian, `"`+k+`":`+v)
	}
	if len(kunci) != 9 || kunci[8] != "receivedRiSlip" {
		t.Fatalf("instrumen: IsiSubmit %v", kunci)
	}
	kode, isi := mintaPeran(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/submit", "{"+strings.Join(bagian, ",")+"}", "UJI-MKT", "ReasFacInMarketing")
	if kode == 400 && strings.Contains(isi, "badan permintaan harus") {
		t.Errorf("IsiSubmit api.ts ditolak handler: %s", isi)
	}
	for _, k := range kunci {
		if kode, isi := mintaPeran(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/submit", `{"x":1}`, "UJI-MKT", "ReasFacInMarketing"); kode != 400 ||
			!strings.Contains(isi, `\"`+k+`\"`) {
			t.Errorf("pesan 400 tanpa kunci %s: %s", k, isi)
		}
	}
}

// TestHasilSubmitNomorPolis - tiket 69a-2: nomorPolis diteruskan dari services (selesai Policy) ke kabel.
func TestHasilSubmitNomorPolis(t *testing.T) {
	if k := keHasilSubmitKabel(services.HasilSubmit{Position: "P", NomorPolis: "UJI-F1.10.2026.00001"}); k.NomorPolis != "UJI-F1.10.2026.00001" ||
		k.Pesan == nil || k.PesanBinding == nil {
		t.Errorf("%+v", k)
	}
}
