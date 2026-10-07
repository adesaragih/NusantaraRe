package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
	"nusantarare/modul/nbfacin/backend/services/rules"
)

type limitTiruan struct {
	a   []models.BarisLimitA
	err error
}

func (l limitTiruan) MuatLimit(context.Context) ([]models.BarisLimitA, []models.BarisLimitB, error) {
	return l.a, nil, l.err
}

func kirim(t *testing.T, svc *services.Service, metode, jalur, isi string) (int, map[string]any) {
	t.Helper()
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, false)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(metode, jalur, strings.NewReader(isi)))
	var jawab map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &jawab); err != nil {
		t.Fatalf("%s %s: jawaban bukan JSON: %q", metode, jalur, w.Body.String())
	}
	return w.Code, jawab
}

// TestPremi - POST /api/nbfacin/premi: premi sebagai TEKS desimal (tanpa float),
// asal rumus; galat mesin 422; JSON rusak 400.
func TestPremi(t *testing.T) {
	svc := services.Baru(limitTiruan{})
	kode, j := kirim(t, svc, "POST", "/api/nbfacin/premi", `{"liniBisnis":"MARINE CARGO","mataUang":"IDR","tsi":"1000000","rate":"0.125"}`)
	if kode != 200 || j["premi"] != "1250.0000" || j["mataUang"] != "IDR" || j["asalRumus"] != "CountGPWMarinePAMbu_Act langkah 1.1.1.1.1" {
		t.Fatalf("%d %v", kode, j)
	}
	if kode, j := kirim(t, svc, "POST", "/api/nbfacin/premi", `{"liniBisnis":"UJI-LINI","tsi":"1","rate":"1"}`); kode != 422 || j["galat"] == nil {
		t.Errorf("lini tak dikenal: %d %v", kode, j)
	}
	if kode, _ := kirim(t, svc, "POST", "/api/nbfacin/premi", `{"liniBisnis":`); kode != 400 {
		t.Errorf("JSON rusak: %d", kode)
	}
	// Medan tak dikenal ditolak - salah eja tidak boleh menjadi medan kosong diam-diam.
	if kode, _ := kirim(t, svc, "POST", "/api/nbfacin/premi", `{"liniBisnis":"MARINE CARGO","tsi":"1","rait":"1"}`); kode != 400 {
		t.Errorf("medan tak dikenal: %d", kode)
	}
}

// TestAkseptasi - POST /api/nbfacin/akseptasi/langkah: jabatan dari isian (butir 58),
// jawaban menandai bahwa itu tidak aman untuk produksi.
func TestAkseptasi(t *testing.T) {
	d := func(s string) *apd.Decimal { v, _ := utils.ParseDecimal(s); return v }
	limit := []models.BarisLimitA{
		{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "SENIORUW", TeamGroup: "1", LimitBottom: d("178500000001"), LimitBottom2: d("178500000001")},
		{Tabel: "M_LIMIT_PROPERTYY", Jabatan: "KADIVTEKNIK", TeamGroup: "1", LimitBottom: d("255000000001"), LimitBottom2: d("289000000001")},
	}
	isi := `{"jabatan":"SENIORUW","anggotaGrup":false,"kasus":{
		"pyWorkPage.OfferFacIn.TotalTSINusaRe":"300000000000","pyWorkPage.OfferFacIn.TotalTSITopRisk":"0",
		"pyWorkPage.IsAdaTopRisk":"false","pyWorkPage.OfferFacIn.QuotationData.StatusBusiness":"1",
		"pyWorkPage.OfferFacIn.QuotationData.TeamGroup":"1","pyWorkPage.OfferFacIn.QuotationData.BusinessType":"Fire",
		"pyWorkPage.OfferFacIn.QuotationData.BusinessOldId":"22","pyWorkPage.OfferFacIn.IsPreferredRisk":"Preferred Risk",
		"pyWorkPage.OfferFacIn.IsBanding":"false","pyWorkPage.OfferFacIn.IsFlagReject":"false",
		"pyWorkPage.PositionNote":"ReasFacInSeniorUnderwriting"}}`
	kode, j := kirim(t, services.Baru(limitTiruan{a: limit}), "POST", "/api/nbfacin/akseptasi/langkah", isi)
	if kode != 200 || j["jabatanTujuan"] != "KADIVTEKNIK" || j["antrean"] != "ReasFacInGroupLeader" || j["selesai"] != false ||
		j["positionNoteDitulis"] != true || !strings.Contains(j["peringatan"].(string), "tidak aman untuk produksi") {
		t.Fatalf("%d %v", kode, j)
	}
	if kode, _ := kirim(t, services.Baru(limitTiruan{a: limit}), "POST", "/api/nbfacin/akseptasi/langkah", `{"kasus":{}}`); kode != 400 {
		t.Errorf("tanpa jabatan: %d", kode)
	}
	if kode, _ := kirim(t, services.Baru(limitTiruan{err: errors.New("oracle mati")}), "POST", "/api/nbfacin/akseptasi/langkah", isi); kode != 500 {
		t.Errorf("repository mati: %d", kode)
	}
	if kode, _ := kirim(t, services.Baru(nil), "POST", "/api/nbfacin/akseptasi/langkah", isi); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
	// Tabel limit kosong = belum tersambung: 503 dengan sebab, bukan 500 tanpa sebab.
	if kode, j := kirim(t, services.Baru(limitTiruan{}), "POST", "/api/nbfacin/akseptasi/langkah", isi); kode != 503 ||
		!strings.Contains(j["galat"].(string), "tabel limit") {
		t.Errorf("tabel limit kosong: %d %v", kode, j)
	}
	if kode, _ := kirim(t, services.Baru(limitTiruan{a: limit}), "POST", "/api/nbfacin/akseptasi/langkah",
		`{"jabatan":"SENIORUW","kasus":{"pyWorkPage.OfferFacIn.TotalTSINusaRe":""}}`); kode != 422 {
		t.Errorf("kasus rusak: %d", kode)
	}
}

// tanggaPanik - tangga tiruan yang panic dengan nilai tertentu.
type tanggaPanik struct{ nilai any }

func (t tanggaPanik) Langkah(kontrak.KasusFacIn, kontrak.PenggunaFacIn) (kontrak.TransisiFacIn, error) {
	panic(t.nilai)
}

// TestAkseptasiPanik - di tingkat HTTP: panic sikap predikat (belum diport) = 422;
// panic lain (bug program) = 500 tanpa rincian ke klien.
func TestAkseptasiPanik(t *testing.T) {
	isi := `{"jabatan":"SENIORUW","kasus":{}}`
	kode, j := kirim(t, services.BaruDenganTangga(tanggaPanik{rules.PanikSikap{Predikat: "UJI", Alasan: "belum diport"}}), "POST", "/api/nbfacin/akseptasi/langkah", isi)
	if kode != 422 || !strings.Contains(j["galat"].(string), "belum diport") {
		t.Errorf("panic sikap: %d %v", kode, j)
	}
	kode, j = kirim(t, services.BaruDenganTangga(tanggaPanik{"rules: rujukan melingkar lewat UJI"}), "POST", "/api/nbfacin/akseptasi/langkah", isi)
	if kode != 500 || j["galat"] != "galat server" {
		t.Errorf("panic bug: %d %v", kode, j)
	}
}

type akunTiruan struct {
	baris []models.Akun
	total int
	err   error
}

func (a akunTiruan) CariAkun(context.Context, string, int, int) ([]models.Akun, int, error) {
	return a.baris, a.total, a.err
}

// TestCariAkun - GET /api/nbfacin/account (tiket 27): bentuk jawaban persis kontrak
// frontend sesi 0f (kolom NULL = ""), `baris` larik walau kosong, 400 halaman tak sah,
// 503 tanpa basis data, 500 galat repository tanpa rincian. Data sintetis UJI-.
func TestCariAkun(t *testing.T) {
	svc := services.Baru(nil).DenganAkun(akunTiruan{total: 21, baris: []models.Akun{
		{ID: "UJI-ID-1", GroupBusinessID: "UJI-GB", GroupBusiness: "UJI GRUP", InsuredID: "UJI-INS-1", InsuredName: "UJI NAMA"},
		{ID: "UJI-ID-2", InsuredID: "UJI-INS-2"},
	}})
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, false)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/api/nbfacin/account?cari=uji&halaman=2", nil))
	mau := `{"baris":[{"id":"UJI-ID-1","insuredId":"UJI-INS-1","insuredName":"UJI NAMA","groupBusinessId":"UJI-GB","groupBusiness":"UJI GRUP"},` +
		`{"id":"UJI-ID-2","insuredId":"UJI-INS-2","insuredName":"","groupBusinessId":"","groupBusiness":""}],"total":21,"halaman":2,"ukuran":15}`
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != mau {
		t.Fatalf("%d %s\nmau %s", w.Code, w.Body.String(), mau)
	}
	kode, j := kirim(t, services.Baru(nil).DenganAkun(akunTiruan{}), "GET", "/api/nbfacin/account", "")
	if b, ok := j["baris"].([]any); kode != 200 || !ok || len(b) != 0 || j["halaman"] != float64(1) {
		t.Errorf("kosong: %d %v", kode, j)
	}
	for _, h := range []string{"0", "-2", "satu", "1.5"} {
		if kode, _ := kirim(t, svc, "GET", "/api/nbfacin/account?halaman="+h, ""); kode != 400 {
			t.Errorf("halaman=%s: %d, mau 400", h, kode)
		}
	}
	if kode, j := kirim(t, services.Baru(nil), "GET", "/api/nbfacin/account", ""); kode != 503 || !strings.Contains(j["galat"].(string), "T_M_ACCOUNT") {
		t.Errorf("tanpa DB: %d %v", kode, j)
	}
	if kode, j := kirim(t, services.Baru(nil).DenganAkun(akunTiruan{err: errors.New("ORA-UJI rincian rahasia")}), "GET", "/api/nbfacin/account", ""); kode != 500 ||
		strings.Contains(j["galat"].(string), "rahasia") {
		t.Errorf("galat repository: %d %v", kode, j)
	}
}

type kelasBisnisTiruan struct {
	baris []models.KelasBisnis
	err   error
}

func (k kelasBisnisTiruan) KelasBisnis(context.Context, string) ([]models.KelasBisnis, error) {
	return k.baris, k.err
}

// TestKelasBisnis - GET /api/nbfacin/class-of-business (tiket 28): bentuk jawaban persis
// `{"baris":[{"id","note"}]}` (urutan repository dipertahankan), `baris` larik walau
// kosong, 400 groupBusinessId kosong/spasi, 503 tanpa basis data, 500 tanpa rincian.
// Data sintetis UJI-.
func TestKelasBisnis(t *testing.T) {
	svc := services.Baru(nil).DenganKelasBisnis(kelasBisnisTiruan{baris: []models.KelasBisnis{
		{ID: "UJI-B2", Note: "UJI A"}, {ID: "UJI-B1", Note: "UJI b & <c>"},
	}})
	mux := http.NewServeMux()
	DaftarkanRute(mux, svc, false)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/api/nbfacin/class-of-business?groupBusinessId=UJI-GB", nil))
	// Dibandingkan sesudah di-decode: encoder JSON Go meloloskan `&` `<` `>` menjadi \u00XX,
	// nilainya tetap teks apa adanya. Kunci tak dikenal = galat (bentuk persis).
	var got struct {
		Baris []struct {
			ID   string `json:"id"`
			Note string `json:"note"`
		} `json:"baris"`
	}
	dek := json.NewDecoder(strings.NewReader(w.Body.String()))
	dek.DisallowUnknownFields()
	if err := dek.Decode(&got); w.Code != 200 || err != nil || len(got.Baris) != 2 ||
		got.Baris[0].ID != "UJI-B2" || got.Baris[0].Note != "UJI A" || got.Baris[1].ID != "UJI-B1" || got.Baris[1].Note != "UJI b & <c>" {
		t.Fatalf("%d %s (%v)", w.Code, w.Body.String(), err)
	}
	if !strings.HasPrefix(w.Body.String(), `{"baris":[{"id":"UJI-B2","note":"UJI A"},{"id":"UJI-B1","note":`) {
		t.Errorf("urutan kunci/baris: %s", w.Body.String())
	}
	kode, j := kirim(t, services.Baru(nil).DenganKelasBisnis(kelasBisnisTiruan{}), "GET", "/api/nbfacin/class-of-business?groupBusinessId=UJI-GB", "")
	if b, ok := j["baris"].([]any); kode != 200 || !ok || len(b) != 0 || len(j) != 1 {
		t.Errorf("kosong: %d %v", kode, j)
	}
	for _, q := range []string{"", "?groupBusinessId=", "?groupBusinessId=%20%20", "?groupbusinessid=UJI-GB"} {
		if kode, j := kirim(t, svc, "GET", "/api/nbfacin/class-of-business"+q, ""); kode != 400 || j["galat"] == nil {
			t.Errorf("%q: %d %v, mau 400", q, kode, j)
		}
	}
	if kode, j := kirim(t, services.Baru(nil), "GET", "/api/nbfacin/class-of-business?groupBusinessId=UJI-GB", ""); kode != 503 || !strings.Contains(j["galat"].(string), "BUSINESS") {
		t.Errorf("tanpa DB: %d %v", kode, j)
	}
	if kode, j := kirim(t, services.Baru(nil).DenganKelasBisnis(kelasBisnisTiruan{err: errors.New("ORA-UJI rincian rahasia")}), "GET", "/api/nbfacin/class-of-business?groupBusinessId=UJI-GB", ""); kode != 500 ||
		strings.Contains(j["galat"].(string), "rahasia") {
		t.Errorf("galat repository: %d %v", kode, j)
	}
}
