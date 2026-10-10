package handlers

import (
	"context"

	"github.com/cockroachdb/apd/v3"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

type lookupAnekaTiruan struct{ kata *string }

func (l lookupAnekaTiruan) IndemnityPct(_ context.Context, indemnity, unit string) ([]*apd.Decimal, error) {
	if l.kata != nil {
		*l.kata = indemnity + "|" + unit
	}
	return []*apd.Decimal{apd.New(80, 0)}, nil
}

func (l lookupAnekaTiruan) CariCoverageAneka(_ context.Context, kode, kata string) ([]models.BarisCoverage, error) {
	*l.kata = kode + "|" + kata
	return []models.BarisCoverage{{ID: "100901", OldID: "A01", Nama: "UJI LIABILITY COVER"}}, nil
}

func (l lookupAnekaTiruan) CariOkupasiAneka(_ context.Context, bisnis, kata string) ([]models.PilihanOkupasiAneka, error) {
	*l.kata = bisnis + "|" + kata
	return []models.PilihanOkupasiAneka{{ID: "10003865", Name: "UJI OKUPASI"}}, nil
}

func (l lookupAnekaTiruan) CariKapal(_ context.Context, kata string) ([]models.BarisKapal, error) {
	*l.kata = kata
	return []models.BarisKapal{{ID: "99", Nama: "UJI KAPAL", Flag: "UJI", Tahun: "2001", NRT: "1", GRT: "2", DWT: "3",
		Klasifikasi: "K", Konstruksi: "BAJA"}}, nil
}

func (l lookupAnekaTiruan) DaftarTipeKendaraan(context.Context) ([]models.BarisTipeKendaraan, error) {
	return nil, nil
}

// TestLookupAneka - GET /kapal?cari= (kata diteruskan apa adanya; bentuk kabel persis) dan /tipe-kendaraan (larik walau kosong);
// 503 tanpa DB.
func TestLookupAneka(t *testing.T) {
	var kata string
	svc := services.Baru(nil).DenganLookupAneka(lookupAnekaTiruan{&kata})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kapal?cari=uji%20kapal", "", "")
	if kode != 200 || kata != "uji kapal" || isi != `{"baris":[{"id":"99","nama":"UJI KAPAL","flag":"UJI","tahun":"2001","nrt":"1","grt":"2","dwt":"3","klasifikasi":"K","konstruksi":"BAJA"}]}` {
		t.Errorf("kapal: %d %q %s", kode, kata, isi)
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/tipe-kendaraan", "", ""); kode != 200 || isi != `{"baris":[]}` {
		t.Errorf("tipe: %d %s", kode, isi)
	}
	for _, j := range []string{"/api/nbfacin/kapal", "/api/nbfacin/tipe-kendaraan"} {
		if kode, _ := minta(t, services.Baru(nil), "GET", j, "", ""); kode != 503 {
			t.Errorf("%s tanpa DB: %d", j, kode)
		}
	}
}

// TestOkupasiAneka - GET …/kasus/{caseId}/occupation-aneka?q=: BusinessName = Class of Business case, q diteruskan; kabel
// {baris:[{id,name}]} persis; 404 case asing; 503 tanpa DB.
func TestOkupasiAneka(t *testing.T) {
	var kata string
	k := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{ClassOfBusiness: " UJI LIABILITY "}}
	svc := services.Baru(nil).DenganLookupAneka(lookupAnekaTiruan{&kata}).DenganKasus(kasusTiruan{k: k})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/occupation-aneka?q=pabrik", "", "")
	if kode != 200 || kata != "UJI LIABILITY|pabrik" || isi != `{"baris":[{"id":"10003865","name":"UJI OKUPASI"}]}` {
		t.Errorf("%d %q %s", kode, kata, isi)
	}
	tanpaKasus := services.Baru(nil).DenganLookupAneka(lookupAnekaTiruan{&kata}).DenganKasus(kasusTiruan{err: repository.ErrKasusTidakAda})
	if kode, _ := minta(t, tanpaKasus, "GET", "/api/nbfacin/kasus/UJI-NB-2/occupation-aneka", "", ""); kode != 404 {
		t.Errorf("404: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1/occupation-aneka", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestCoverageAneka - GET …/kasus/{caseId}/coverage-aneka?q=: BizCode = BUSINESS.ID ber-GROUPPANEL case, q diteruskan; bentuk
// BarisCoverage FIRE; 404; 503.
func TestCoverageAneka(t *testing.T) {
	var kata string
	k := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{ClassOfBusiness: "UJI LIABILITY", GroupBusinessID: "3"}}
	svc := services.Baru(nil).DenganLookupAneka(lookupAnekaTiruan{&kata}).DenganKasus(kasusTiruan{k: k}).
		DenganTableOfLimit(tolTiruan{bisnis: []models.BarisBisnisLini{{ID: "10021", OldID: "15", GroupPanel: "003"}}})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/coverage-aneka?q=liab", "", "")
	if kode != 200 || kata != "10021|liab" || isi != `{"baris":[{"id":"100901","oldId":"A01","nama":"UJI LIABILITY COVER"}]}` {
		t.Errorf("%d %q %s", kode, kata, isi)
	}
	tanpaKasus := services.Baru(nil).DenganLookupAneka(lookupAnekaTiruan{&kata}).DenganKasus(kasusTiruan{err: repository.ErrKasusTidakAda})
	if kode, _ := minta(t, tanpaKasus, "GET", "/api/nbfacin/kasus/UJI-NB-2/coverage-aneka", "", ""); kode != 404 {
		t.Errorf("404: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1/coverage-aneka", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}
