package handlers_test

// Uji jalur HTTP tombol Add Revision / Add Adjustment Premium dan `Choose`.

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nusantarare/modul/treatyinadjustment/backend/handlers"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

func (g gudangTiruan) DaftarMasterPilihan(_ context.Context, hanyaNonProp bool) ([]models.BarisMasterPilihan, error) {
	if !hanyaNonProp {
		return g.master, nil
	}
	out := []models.BarisMasterPilihan{}
	for _, b := range g.master {
		if b.SifatProporsi == "NonProportional" {
			out = append(out, b)
		}
	}
	return out, nil
}

func (g gudangTiruan) BacaDokumenMaster(_ context.Context, id string) (models.SisiPenyesuaian, bool, error) {
	d, ada := g.dokumenMaster[id]
	return d, ada, nil
}

func (g gudangTiruan) AdaRevisi(context.Context, string) (bool, error) { return false, nil }

func routerDraf() http.Handler {
	return handlers.RouterDengan(services.LayananDengan(gudangTiruan{
		master: []models.BarisMasterPilihan{
			{ID: "1000506", SifatProporsi: "NonProportional"},
			{ID: "1000507", SifatProporsi: "Proportional"},
		},
		dokumenMaster: map[string]models.SisiPenyesuaian{"1000506": {
			Medan: map[string]string{"ProportionType": "NonProportional"},
			Larik: map[string][]map[string]string{},
		}},
	}), true, true)
}

func kirimDraf(t *testing.T, badan string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodPost, "/api/treaty-in-adjustment/penyesuaian-warisan/draf", strings.NewReader(badan))
	r.Header.Set("X-Pelaku", "UJI")
	w := httptest.NewRecorder()
	routerDraf().ServeHTTP(w, r)
	return w
}

func TestDaftarMasterPickerPerJenis(t *testing.T) {
	for jenis, cacah := range map[string]int{"revisi": 2, "premi": 1} {
		w := minta(t, routerDraf(), "/api/treaty-in-adjustment/penyesuaian-warisan/master?jenis="+jenis, "UJI")
		var d []models.BarisMasterPilihan
		if w.Code != http.StatusOK || json.Unmarshal(w.Body.Bytes(), &d) != nil || len(d) != cacah {
			t.Errorf("%s: kode %d, %d baris, mau %d (badan %s)", jenis, w.Code, len(d), cacah, w.Body.String())
		}
	}
	if w := minta(t, routerDraf(), "/api/treaty-in-adjustment/penyesuaian-warisan/master?jenis=lain", "UJI"); w.Code != http.StatusBadRequest {
		t.Errorf("jenis lain: kode %d, mau 400", w.Code)
	}
}

func TestChooseMengembalikanDrafTanpaMenyimpan(t *testing.T) {
	w := kirimDraf(t, `{"id":"1000506","internalType":"3","materialType":"1"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d (badan %s)", w.Code, w.Body.String())
	}
	var p models.Penyesuaian
	if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	if p.ID != "1000506/R01" || p.IDAsal != "1000506" || p.Baru.Medan["EDMState"] != "3" || p.Lama.Medan["ID"] != "1000506" {
		t.Errorf("draf %+v", p)
	}
}

func TestChooseMenolakMasukanSalah(t *testing.T) {
	for badan, kode := range map[string]int{
		`bukan json`:                          http.StatusBadRequest,
		`{"id":"1000506","internalType":"9"}`: http.StatusBadRequest,
		`{"id":"9999999","internalType":"1"}`: http.StatusNotFound,
	} {
		if w := kirimDraf(t, badan); w.Code != kode {
			t.Errorf("%s: kode %d, mau %d", badan, w.Code, kode)
		}
	}
}
