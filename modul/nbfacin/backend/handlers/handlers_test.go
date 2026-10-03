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
	DaftarkanRute(mux, svc)
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
