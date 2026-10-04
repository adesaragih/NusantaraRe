package handlers_test

// Rute HTTP Company Detail di atas gudang tiruan: bentuk jawaban dan kode status.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/companydetail/backend/handlers"
	"nusantarare/modul/companydetail/backend/models"
	"nusantarare/modul/companydetail/backend/services"
	"nusantarare/modul/companydetail/backend/tiruan"
)

func router(g *tiruan.Gudang, adaDB bool) http.Handler {
	return handlers.RouterDengan(services.BaruLayanan(g, tiruan.Transaksi), adaDB, true)
}

func kirim(t *testing.T, h http.Handler, metode, jalur, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
	r.Header.Set("X-Pelaku", "UJI-ADMIN")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

var idAnak = url.PathEscape("ASM-SFAGIS-WORK-ORG ORG-115")

func TestRuteDaftarPilihanIndukDanDetail(t *testing.T) {
	h := router(tiruan.Contoh(), true)
	w := kirim(t, h, "GET", handlers.Prefix+"?q=uji&halaman=1&ukuran=10", "")
	var d struct {
		Daftar []map[string]any `json:"daftar"`
		Total  int              `json:"total"`
	}
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &d) != nil || d.Total != 2 {
		t.Fatalf("daftar %d %s", w.Code, w.Body.String())
	}
	for _, kunci := range []string{"id", "idView", "nama", "title", "npwp", "countryName", "businessField", "parentName"} {
		if _, ada := d.Daftar[0][kunci]; !ada {
			t.Errorf("baris daftar tanpa %q", kunci)
		}
	}
	if w := kirim(t, h, "GET", handlers.Prefix+"/pilihan", ""); w.Code != 200 || !strings.Contains(w.Body.String(), `"negara":[{"id":"100901"`) {
		t.Errorf("pilihan %d %s", w.Code, w.Body.String())
	}
	if w := kirim(t, h, "GET", handlers.Prefix+"/induk?q=induk&kecuali="+url.QueryEscape("ASM-SFAGIS-WORK-ORG ORG-115"), ""); w.Code != 200 ||
		!strings.Contains(w.Body.String(), `"idView":"ORG-100"`) || strings.Contains(w.Body.String(), "ORG-115") {
		t.Errorf("induk %d %s", w.Code, w.Body.String())
	}
	w = kirim(t, h, "GET", handlers.Prefix+"/"+idAnak, "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"idView":"ORG-115"`) || !strings.Contains(w.Body.String(), `"telfax":[{"type":"2"`) {
		t.Errorf("detail %d %s", w.Code, w.Body.String())
	}
}

func TestRuteCreateUbahDanGalat(t *testing.T) {
	g := tiruan.Contoh()
	h := router(g, true)
	badan := `{"nama":"UJI Baru","title":"PT.","npwp":"","country":"001","businessField":"01","parentId":"","note":"",` +
		`"pic":[{"userIdentifier":"","nama":"UJI PIC","position":"Staff","gender":"2","email":"","dateOfBirth":"","phone":""}],` +
		`"alamat":[{"asal":"","type":"2","address":"UJI Jalan","telfax":[{"type":"3","code":"","no":"UJI-HP"}]}]}`
	w := kirim(t, h, "POST", handlers.Prefix, badan)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"idView":"ORG-121"`) || !strings.Contains(w.Body.String(), `"userIdentifier":"PIC-1"`) {
		t.Fatalf("create %d %s", w.Code, w.Body.String())
	}
	for _, k := range []struct {
		metode, jalur, badan string
		kode                 int
		isi                  string
	}{
		{"POST", handlers.Prefix, `{"nama":""}`, 422, "Organization Name is required"},
		{"POST", handlers.Prefix, `{"nama":"x","asing":1}`, 400, "not valid JSON"},
		{"PUT", handlers.Prefix + "/UJI-TIDAK-ADA", `{"nama":"x","country":"001","businessField":"01"}`, 404, "Organization not found"},
		{"GET", handlers.Prefix + "/UJI-TIDAK-ADA", "", 404, "Organization not found"},
	} {
		if w := kirim(t, h, k.metode, k.jalur, k.badan); w.Code != k.kode || !strings.Contains(w.Body.String(), k.isi) {
			t.Errorf("%s %s: %d %s, mau %d %q", k.metode, k.jalur, w.Code, w.Body.String(), k.kode, k.isi)
		}
	}
	ubah := `{"nama":"UJI Anak Usaha Baru","title":"PT.","npwp":"","country":"001","businessField":"23",` +
		`"parentId":"ASM-SFAGIS-WORK-ORG ORG-100","note":"","pic":[],"alamat":[]}`
	if w := kirim(t, h, "PUT", handlers.Prefix+"/"+idAnak, ubah); w.Code != 200 || !strings.Contains(w.Body.String(), `"nama":"UJI ANAK USAHA BARU"`) {
		t.Errorf("ubah %d %s", w.Code, w.Body.String())
	}
	if len(g.PIC["ASM-SFAGIS-WORK-ORG ORG-115"]) != 0 || len(g.Alamat["ASM-SFAGIS-WORK-ORG ORG-115"]) != 0 {
		t.Error("PIC dan alamat yang dibuang dari grid tidak terhapus")
	}
}

func TestRuteTanpaDatabase503(t *testing.T) {
	if w := kirim(t, router(tiruan.Contoh(), false), "GET", handlers.Prefix, ""); w.Code != 503 {
		t.Errorf("tanpa database %d", w.Code)
	}
}

func TestRutePeriksaNama(t *testing.T) {
	h := router(tiruan.Contoh(), true)
	w := kirim(t, h, "GET", handlers.Prefix+"/periksa-nama?nama="+url.QueryEscape("PT UJI Anak Usaha"), "")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"titleDalamNama":"PT."`) || !strings.Contains(w.Body.String(), `"jenis":"sama"`) {
		t.Errorf("periksa nama %d %s", w.Code, w.Body.String())
	}
}

func TestRuteCopyOld(t *testing.T) {
	g := tiruan.Contoh()
	g.Lama = []models.OrgLama{{ID: "ASM-SFAGIS-WORK-ORG ORG-118", IDView: "ORG-118", Nama: "UJI LAMA", Baru: true, Isi: []string{}}}
	h := router(g, true)
	kirimMenu := func(metode, jalur, badan string, menu ...string) *httptest.ResponseRecorder {
		r := httptest.NewRequest(metode, jalur, strings.NewReader(badan))
		r.Header.Set("X-Pelaku", "UJI-ADMIN")
		if menu != nil {
			r = r.WithContext(inti.DenganAksesMenu(r.Context(), menu))
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	if w := kirimMenu("GET", handlers.Prefix+"/hak", "", "companydetail"); w.Code != 200 || !strings.Contains(w.Body.String(), `"copyOld":false`) {
		t.Errorf("hak bukan superadmin %d %s", w.Code, w.Body.String())
	}
	if w := kirimMenu("GET", handlers.Prefix+"/lama", "", "companydetail"); w.Code != 403 {
		t.Errorf("lama bukan superadmin %d %s", w.Code, w.Body.String())
	}
	if w := kirimMenu("GET", handlers.Prefix+"/hak", "", "kelolauser"); !strings.Contains(w.Body.String(), `"copyOld":true`) {
		t.Errorf("hak superadmin %s", w.Body.String())
	}
	if w := kirimMenu("GET", handlers.Prefix+"/lama", "", "kelolauser"); w.Code != 200 || !strings.Contains(w.Body.String(), `"idView":"ORG-118"`) {
		t.Errorf("lama superadmin %d %s", w.Code, w.Body.String())
	}
	w := kirimMenu("POST", handlers.Prefix+"/lama/salin", `{"ids":["ASM-SFAGIS-WORK-ORG ORG-118"]}`, "kelolauser")
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"status":"disalin"`) || !strings.Contains(w.Body.String(), `"disalin":1`) {
		t.Errorf("salin %d %s", w.Code, w.Body.String())
	}
}
