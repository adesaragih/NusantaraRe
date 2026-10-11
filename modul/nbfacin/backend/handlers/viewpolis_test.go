package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// viewTiruanH - PembacaViewPolis tiket 71 perluasan untuk uji HTTP: satu baris berisi tiap tab, total 21.
type viewTiruanH struct{}

func (viewTiruanH) CariOfferFac(context.Context, string, string, int, int) ([]models.BarisOfferFac, int, error) {
	return []models.BarisOfferFac{{CaseID: "NB-1", NoOffer: "RNM-OF-1", InsuredName: "PT UJI", SourceOfBusiness: "BRK",
		StartPeriod: "01-01-2026", EndPeriod: "01-01-2027"}}, 21, nil
}

func (viewTiruanH) CariPolisTreaty(context.Context, string, string, int, int) ([]models.BarisPolisTreaty, int, error) {
	return []models.BarisPolisTreaty{{CaseID: "NB-2", PolicyNo: "RNM-T1", InsuredName: "PT UJI", SourceOfBusiness: "BRK",
		StartPeriod: "01-01-2026", EndPeriod: "01-01-2027"}}, 21, nil
}

func (viewTiruanH) CariPolisLife(context.Context, string, string, int, int) ([]models.BarisPolisLife, int, error) {
	return []models.BarisPolisLife{{CaseID: "NBLF-3", PolicyNo: "RNM-L1", EndorsNo: "", TypePolicy: "T", SourceOfBusiness: "S",
		CedingCoName: "C", QQ: "Q", RISlipCeding: "R"}}, 21, nil
}

func (viewTiruanH) CariOfferLife(context.Context, string, string, int, int) ([]models.BarisOfferLife, int, error) {
	return []models.BarisOfferLife{{CaseID: "NBLF-4", NoOffer: "OFR-4", TypePolicy: "T", BusinessName: "B", CedingCoName: "C", QQ: "Q",
		LastStatus: "L"}}, 21, nil
}

func (viewTiruanH) BerkasPolisTreaty(_ context.Context, noPolis string) (models.BerkasPolisTreaty, bool, error) {
	if noPolis == "RNM-T2" {
		return models.BerkasPolisTreaty{Modul: models.ModulEDMTreatyIn, Kasus: "EDMT-2"}, true, nil
	}
	return models.BerkasPolisTreaty{}, false, nil
}

// TestViewPolisHandler - tiket 71 perluasan: bentuk kabel (kunci + urutan dikunci untuk frontend 9d), hanya login, filter kosong ->
// baris [] total 0, opsi tak sah / halaman bukan angka 400, 401.
func TestViewPolisHandler(t *testing.T) {
	svc := layananKasus().DenganViewPolis(viewTiruanH{})
	for jalur, mau := range map[string]string{
		"/api/nbfacin/polis/offer?opsi=5&filter=PT&halaman=2": `{"baris":[{"caseId":"NB-1","noOffer":"RNM-OF-1","insuredName":"PT UJI",` +
			`"sourceOfBusiness":"BRK","qqName":"","startPeriod":"01-01-2026","endPeriod":"01-01-2027"}],"total":21,"halaman":2,"ukuran":20}`,
		"/api/nbfacin/polis/treaty?opsi=4&filter=PT": `{"baris":[{"caseId":"NB-2","policyNo":"RNM-T1","insuredName":"PT UJI",` +
			`"sourceOfBusiness":"BRK","startPeriod":"01-01-2026","endPeriod":"01-01-2027"}],"total":21,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/life?filter=PT": `{"baris":[{"caseId":"NBLF-3","policyNo":"RNM-L1","endorsNo":"","typePolicy":"T",` +
			`"sourceOfBusiness":"S","cedingCoName":"C","qq":"Q","riSlipCeding":"R"}],"total":21,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/life-offer?filter=PT": `{"baris":[{"caseId":"NBLF-4","noOffer":"OFR-4","typePolicy":"T","businessName":"B",` +
			`"cedingCoName":"C","qq":"Q","lastStatus":"L"}],"total":21,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/offer?opsi=2&filter=":  `{"baris":[],"total":0,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/treaty?opsi=2&filter=": `{"baris":[],"total":0,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/life?filter=%22%27":    `{"baris":[],"total":0,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/life-offer?filter=+":   `{"baris":[],"total":0,"halaman":1,"ukuran":20}`,
		"/api/nbfacin/polis/life-offer?halaman=3":  `{"baris":[],"total":0,"halaman":3,"ukuran":20}`,
	} {
		if kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-LAIN", ""); kode != 200 || strings.TrimSpace(isi) != mau {
			t.Errorf("%s: %d\n%s\n%s", jalur, kode, isi, mau)
		}
	}
	for jalur, mau := range map[string]int{
		"/api/nbfacin/polis/offer?opsi=7&filter=X":           400,
		"/api/nbfacin/polis/offer?opsi=3&filter=X":           400,
		"/api/nbfacin/polis/offer?filter=X":                  400,
		"/api/nbfacin/polis/treaty?opsi=6&filter=X":          400,
		"/api/nbfacin/polis/treaty?opsi=&filter=":            400,
		"/api/nbfacin/polis/life?filter=X&halaman=satu":      400,
		"/api/nbfacin/polis/life-offer?filter=X&halaman=0":   400,
		"/api/nbfacin/polis/treaty?opsi=2&filter=X&halaman=": 200,
	} {
		if kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-LAIN", ""); kode != mau {
			t.Errorf("%s: %d %s", jalur, kode, isi)
		}
	}
	for _, jalur := range []string{"/api/nbfacin/polis/offer?opsi=2&filter=X", "/api/nbfacin/polis/treaty?opsi=2&filter=X",
		"/api/nbfacin/polis/life?filter=X", "/api/nbfacin/polis/life-offer?filter=X"} {
		if kode, _ := mintaPeran(t, svc, "GET", jalur, "", "", ""); kode != 401 {
			t.Errorf("401 %s: %d", jalur, kode)
		}
		if kode, _ := mintaPeran(t, layananKasus(), "GET", jalur, "", "UJI-LAIN", ""); kode != 503 {
			t.Errorf("503 %s: %d", jalur, kode)
		}
	}
}

// TestBerkasPolisTreatyHandler - GET /api/nbfacin/polis/treaty/berkas: {modul, kasus}; tanpa berkas 404; nopolis kosong 400; 401; 503.
func TestBerkasPolisTreatyHandler(t *testing.T) {
	svc := layananKasus().DenganViewPolis(viewTiruanH{})
	if kode, isi := mintaPeran(t, svc, "GET", "/api/nbfacin/polis/treaty/berkas?nopolis=RNM-T2", "", "UJI-LAIN", ""); kode != 200 ||
		strings.TrimSpace(isi) != `{"modul":"edmtreatyin","kasus":"EDMT-2"}` {
		t.Errorf("200: %d %s", kode, isi)
	}
	for jalur, mau := range map[string]int{"/api/nbfacin/polis/treaty/berkas?nopolis=RNM-LAMA": 404,
		"/api/nbfacin/polis/treaty/berkas?nopolis=": 400, "/api/nbfacin/polis/treaty/berkas": 400} {
		if kode, isi := mintaPeran(t, svc, "GET", jalur, "", "UJI-LAIN", ""); kode != mau {
			t.Errorf("%s: %d %s", jalur, kode, isi)
		}
	}
	if kode, _ := mintaPeran(t, svc, "GET", "/api/nbfacin/polis/treaty/berkas?nopolis=RNM-T2", "", "", ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	if kode, _ := mintaPeran(t, layananKasus(), "GET", "/api/nbfacin/polis/treaty/berkas?nopolis=RNM-T2", "", "UJI-LAIN", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}
