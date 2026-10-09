package handlers_test

// Uji jalur HTTP layar Adjustment - daftar dan satu penyesuaian.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"nusantarare/modul/treatyinadjustment/backend/handlers"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/repository"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

func (g gudangTiruan) DaftarPenyesuaianWarisan(context.Context) ([]models.BarisPenyesuaian, error) {
	return g.daftarPenyesuaian, g.galatPenyesuaian
}

func (g gudangTiruan) BacaPenyesuaianPendaratan(_ context.Context, id string) (models.Penyesuaian, error) {
	if g.galatPenyesuaian != nil {
		return models.Penyesuaian{}, g.galatPenyesuaian
	}
	p, ada := g.penyesuaian[id]
	if !ada {
		return models.Penyesuaian{}, fmt.Errorf("%w: %s", repository.ErrPenyesuaianTidakAda, id)
	}
	return p, nil
}

func routerPenyesuaian() http.Handler {
	p := models.Penyesuaian{
		ID: "1000080/R02", IDAsal: "1000080/R01",
		Baru: models.SisiPenyesuaian{Medan: map[string]string{"EDMState": "2"}, Larik: map[string][]map[string]string{}},
		Lama: models.SisiPenyesuaian{Medan: map[string]string{}, Larik: map[string][]map[string]string{}},
	}
	return handlers.RouterDengan(services.LayananDengan(gudangTiruan{
		daftarPenyesuaian: []models.BarisPenyesuaian{{ID: p.ID, IDAsal: p.IDAsal, JenisPenyesuaian: "2"}},
		penyesuaian:       map[string]models.Penyesuaian{p.ID: p},
	}), true, true)
}

func jalurSatu(id string) string {
	return "/api/treaty-in-adjustment/penyesuaian-warisan/satu?id=" + url.QueryEscape(id)
}

func TestDaftarPenyesuaianMenjawab200(t *testing.T) {
	w := minta(t, routerPenyesuaian(), "/api/treaty-in-adjustment/penyesuaian-warisan", "UJI")
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200 (badan %s)", w.Code, w.Body.String())
	}
	var d []models.BarisPenyesuaian
	if err := json.Unmarshal(w.Body.Bytes(), &d); err != nil || len(d) != 1 || d[0].IDAsal != "1000080/R01" {
		t.Fatalf("badan tidak sesuai: %v %s", err, w.Body.String())
	}
}

// ⛔ Pengenal BERGARIS MIRING sampai utuh - itulah alasan ia dibawa lewat
// parameter kueri, bukan segmen jalur.
func TestSatuPenyesuaianBergarisMiringMenjawab200(t *testing.T) {
	w := minta(t, routerPenyesuaian(), jalurSatu("1000080/R02"), "UJI")
	if w.Code != http.StatusOK {
		t.Fatalf("kode %d, mau 200 (badan %s)", w.Code, w.Body.String())
	}
	var p models.Penyesuaian
	if err := json.Unmarshal(w.Body.Bytes(), &p); err != nil {
		t.Fatal(err)
	}
	if p.ID != "1000080/R02" || p.Baru.Medan["EDMState"] != "2" {
		t.Errorf("badan tidak sesuai: %+v", p)
	}
}

func TestSatuPenyesuaianTanpaIdentitasMenjawab401(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), true, false)
	for _, j := range []string{"/api/treaty-in-adjustment/penyesuaian-warisan", jalurSatu("1000080/R02")} {
		if w := minta(t, h, j, ""); w.Code != http.StatusUnauthorized {
			t.Errorf("%s: kode %d, mau 401", j, w.Code)
		}
	}
}

func TestSatuPenyesuaianTanpaPengenalMenjawab400(t *testing.T) {
	w := minta(t, routerPenyesuaian(), "/api/treaty-in-adjustment/penyesuaian-warisan/satu", "UJI")
	if w.Code != http.StatusBadRequest {
		t.Errorf("kode %d, mau 400", w.Code)
	}
}

func TestSatuPenyesuaianTidakAdaMenjawab404(t *testing.T) {
	w := minta(t, routerPenyesuaian(), jalurSatu("9999999/R01"), "UJI")
	if w.Code != http.StatusNotFound {
		t.Errorf("kode %d, mau 404", w.Code)
	}
}

// Galat tak terduga -> 500 berkalimat umum; sebab aslinya ke log, bukan
// ke layar.
func TestSatuPenyesuaianGalatGudangMenjawab500(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{
		galatPenyesuaian: errors.New("ORA-00942: table or view does not exist"),
	}), true, true)
	w := minta(t, h, jalurSatu("1000080/R02"), "UJI")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("kode %d, mau 500", w.Code)
	}
	if strings.Contains(w.Body.String(), "ORA-00942") {
		t.Errorf("sebab Oracle bocor ke layar: %s", w.Body.String())
	}
}

func TestPenyesuaianTanpaBasisDataMenjawab503(t *testing.T) {
	h := handlers.RouterDengan(services.LayananDengan(gudangTiruan{}), false, true)
	if w := minta(t, h, jalurSatu("1000080/R02"), "UJI"); w.Code != http.StatusServiceUnavailable {
		t.Errorf("kode %d, mau 503", w.Code)
	}
}
