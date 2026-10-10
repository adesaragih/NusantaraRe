package handlers

import (
	"context"

	"encoding/json"
	"github.com/cockroachdb/apd/v3"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// anekaTiruanH - PenyimpanObjekAneka dalam memori (satu case "UJI-NB-1").
type anekaTiruanH struct{ d *[]models.LokasiAneka }

func (anekaTiruanH) PercentShareKasus(context.Context, string) (*apd.Decimal, error) {
	return apd.New(40, 0), nil
}

func (a anekaTiruanH) BacaObjekAneka(_ context.Context, id string) ([]models.LokasiAneka, error) {
	if id != "UJI-NB-1" {
		return nil, repository.ErrKasusTidakAda
	}
	return *a.d, nil
}

func (anekaTiruanH) ProRateKasus(context.Context, string) (*apd.Decimal, error) { return nil, nil }

func (anekaTiruanH) UbahPercentShare(context.Context, *db.Tx, string, *apd.Decimal) error { return nil }

func (a anekaTiruanH) GantiObjekAneka(_ context.Context, _ *db.Tx, id string, lokasi []models.LokasiAneka, _ *apd.Decimal) error {
	if id != "UJI-NB-1" {
		return repository.ErrKasusTidakAda
	}
	*a.d = lokasi
	return nil
}

type mataUangTiruanH struct{}

func (mataUangTiruanH) DaftarJenisItem(context.Context) ([]models.JenisItem, error) { return nil, nil }
func (mataUangTiruanH) DaftarMataUang(context.Context) ([]string, error) {
	return []string{"IDR", "USD"}, nil
}

// kunciAneka - kunci objek JSON terurut.
func kunciAneka(t *testing.T, m map[string]any) []string {
	t.Helper()
	var k []string
	for n := range m {
		k = append(k, n)
	}
	sort.Strings(k)
	return k
}

// TestObjekAneka - GET/PUT …/objek-aneka (tiket 73 P1): bentuk kabel (kunci persis tiap tingkat), larik walau kosong, desimal
// teks, tanggal DD-MM-YYYY bolak-balik, suretyType dari server; 400 desimal ber-jalur / JSON rusak; 401; 404.
func TestObjekAneka(t *testing.T) {
	var d []models.LokasiAneka
	svc := services.Baru(nil).DenganObjekAneka(anekaTiruanH{&d}).DenganTransaksi(tanpaOracle).
		DenganPilihanItem(mataUangTiruanH{})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek-aneka", "", ""); kode != 200 || isi != `{"lokasi":[],"proRatePercent":""}` {
		t.Fatalf("kosong: %d %s", kode, isi)
	}
	badan := `{"lokasi":[{"objectNo":"1","isTopRisk":true,"riskLocation":"UJI ALAMAT","okupasi":[{"occupationId":"UJI01",` +
		`"occupationName":"UJI OKUPASI","objek":[{"objectName":"UJI KAPAL","currency":"IDR","tsi":"1500000.50","startDate":"01-10-2026",` +
		`"briguna":{"salary":"7500000","dateOfBirth":"31-12-1990"},"marineHull":{"typeOfVessel":"TUG","grtrt":"120"},` +
		`"memorandum":{"suretyType":"KLIEN","keyPerson":"UJI PIC"}}]},{"occupationId":"UJI02","objek":[]}],` +
		`"objekGolf":[{"objectName":"UJI GOLF","currency":"USD","limitOfLiability":"250000"}]}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-aneka", badan, "UJI-USER")
	if kode != 200 {
		t.Fatalf("%d %s", kode, isi)
	}
	var j struct {
		Lokasi []map[string]any `json:"lokasi"`
	}
	if err := json.Unmarshal([]byte(isi), &j); err != nil || len(j.Lokasi) != 1 {
		t.Fatalf("%v %s", err, isi)
	}
	lok := j.Lokasi[0]
	mauLokasi := "buildingNo city country district isTopRisk objectNo objekGolf okupasi province riskAddressId riskLocation roadName roadType territory totalPerCurrency zipCode"
	if k := strings.Join(kunciAneka(t, lok), " "); k != mauLokasi {
		t.Errorf("kunci lokasi: %s", k)
	}
	ok := lok["okupasi"].([]any)
	if len(ok) != 2 || strings.Join(kunciAneka(t, ok[0].(map[string]any)), " ") != "objek occupationId occupationName" {
		t.Fatalf("okupasi: %v", ok)
	}
	if o2 := ok[1].(map[string]any)["objek"].([]any); o2 == nil || len(o2) != 0 {
		t.Errorf("objek kosong bukan larik: %v", ok[1])
	}
	obj := ok[0].(map[string]any)["objek"].([]any)[0].(map[string]any)
	mauObjek := "almRiskId areaHectar aviationHull briguna builtIn coverages currency draftWordingId endDate flag geographicalLimits invoiceNo " +
		"limitIdemnityDay limitIdemnityMonth limitOfLiability maintenancePeriod marineHull memorandum noKontrak noKtp noNpwp noOfTree " +
		"objectName others quantity section serialNo startDate tsi tsiOld unitQuantity vehicleHe year"
	if k := strings.Join(kunciAneka(t, obj), " "); k != mauObjek {
		t.Errorf("kunci objek: %s", k)
	}
	for anak, mau := range map[string]string{
		"briguna":      "dateOfBirth depositAccountNo idCardNo institution job loanAccountNumber loanName loanType salary",
		"marineHull":   "breadthm classification construction depthm dwtrt grtrt lengthm nrtrt typeOfVessel",
		"aviationHull": "maxPassengers pilots registrationMarks specialRentalUses specialUses standardUses",
		"vehicleHe":    "brand brandName chassisNumber engineNumber objectIdHe objectNameHe type typeName",
		"memorandum": "accumulation applicationNumber buildingStatus collateral conclusionsProposals insuredSiup keyPerson " +
			"originalRate principalsBackground principalsData principalsGuarantee principalsOutGo summaryFinancial " +
			"supportingDocument suretyType typeOfWording",
	} {
		if k := strings.Join(kunciAneka(t, obj[anak].(map[string]any)), " "); k != mau {
			t.Errorf("kunci %s: %s", anak, k)
		}
	}
	for _, harus := range []string{`"tsi":"1500000.50"`, `"startDate":"01-10-2026"`, `"dateOfBirth":"31-12-1990"`,
		`"salary":"7500000"`, `"grtrt":"120"`, `"suretyType":""`, `"keyPerson":"UJI PIC"`, `"limitOfLiability":"250000"`,
		`"tsiOld":""`, `"isTopRisk":true`} {
		if !strings.Contains(isi, harus) {
			t.Errorf("tanpa %s: %s", harus, isi)
		}
	}
	if len(d) != 1 || d[0].Okupasi[0].Objek[0].StartDate != "20261001T050000.000 GMT" || d[0].ObjekGolf[0].Currency != "USD" {
		t.Errorf("model tersimpan: %+v", d)
	}
	// P2: kunci `coverages` absen -> nil (dipertahankan); `[]` -> larik kosong (ganti utuh).
	if d[0].Okupasi[0].Objek[0].Coverages != nil {
		t.Errorf("coverages absen jadi %v", d[0].Okupasi[0].Objek[0].Coverages)
	}
	if kode, _ := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-aneka", `{"lokasi":[{"objekGolf":[{"coverages":[]}]}]}`,
		"UJI-USER"); kode != 200 || d[0].ObjekGolf[0].Coverages == nil {
		t.Errorf("coverages [] jadi nil: %d", kode)
	}
	for nama, u := range map[string]struct {
		badan, pelaku string
		kode          int
		pesan         string
	}{
		"desimal koma": {`{"lokasi":[{"objekGolf":[{"tsi":"1,5"}]}]}`, "UJI-USER", 400, "lokasi[0].objekGolf[0].tsi"},
		"gaji teks":    {`{"lokasi":[{"okupasi":[{"objek":[{"briguna":{"salary":"x"}}]}]}]}`, "UJI-USER", 400, "lokasi[0].okupasi[0].objek[0].briguna.salary"},
		"JSON rusak":   {`{"lokasi":`, "UJI-USER", 400, "daftar lokasi Aneka"},
		"tanpa pelaku": {badan, "", 401, ""},
	} {
		if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-aneka", u.badan, u.pelaku); kode != u.kode ||
			!strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-2/objek-aneka", "", ""); kode != 404 {
		t.Errorf("404: %d", kode)
	}
}

// TestObjekAnekaProteksi422 - kabel 422 `{galat, pelanggaran:[{lokasi, okupasi, objek, medan, pesan}]}`; okupasi null untuk Golf;
// pesan tanpa akhiran spasi + `\n` literal; tidak ada yang ditulis.
func TestObjekAnekaProteksi422(t *testing.T) {
	var d []models.LokasiAneka
	k := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{ClassOfBusiness: "UJI HE", GroupBusinessID: "3"}}
	svc := services.Baru(nil).DenganObjekAneka(anekaTiruanH{&d}).DenganTransaksi(tanpaOracle).DenganPilihanItem(mataUangTiruanH{}).
		DenganKasus(kasusTiruan{k: k}).DenganTableOfLimit(tolTiruan{bisnis: []models.BarisBisnisLini{{ID: "10999", OldID: "20", GroupPanel: "003"}}})
	badan := `{"lokasi":[{"okupasi":[{"objek":[{"year":"2020","vehicleHe":{"chassisNumber":"C","engineNumber":"E","type":"1",` +
		`"objectNameHe":"TBA","brand":"3"}}]}],"objekGolf":[{"year":"2021","vehicleHe":{"chassisNumber":"C","engineNumber":"E",` +
		`"type":"1","objectIdHe":"2"}}]}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-aneka", badan, "UJI-USER")
	mau := `{"galat":"Nama Obyek tidak boleh TBA; Merk Obyek harus diisi","pelanggaran":[` +
		`{"lokasi":0,"okupasi":0,"objek":0,"medan":"vehicleHe.objectNameHe","pesan":"Nama Obyek tidak boleh TBA"},` +
		`{"lokasi":0,"okupasi":null,"objek":0,"medan":"vehicleHe.brandName","pesan":"Merk Obyek harus diisi"}]}`
	if kode != 422 || isi != mau || d != nil {
		t.Errorf("%d %s (tertulis %v)", kode, isi, d)
	}
}

// TestHitungCoverageAneka - POST …/hitung-coverage-aneka: `{proRatePercent, baris:[{coverage, pemicu}]}` -> `{proRatePercent, baris}`
// urutan sama, kunci coverage = kabel CoverageAneka; ProRate case berantai antar baris (tsi = Status Coverage -> dari tanggal);
// 400 pemicu / desimal; 404 case asing.
func TestHitungCoverageAneka(t *testing.T) {
	k := &models.Kasus{CaseID: "UJI-NB-1", General: models.General{StartDateTime: "20260101T050000.000 GMT",
		EndDateTime: "20270101T050000.000 GMT"}}
	svc := services.Baru(nil).DenganKasus(kasusTiruan{k: k})
	badan := `{"proRatePercent":"50","baris":[` +
		`{"pemicu":"loading","coverage":{"calculateMethodFacIn":"1","tsi":"1000000","rate":"0.5","loading":"0","day":"365"}},` +
		`{"pemicu":"tsi","coverage":{"calculateMethodFacIn":"1","tsi":"2000000","rate":"0.5","day":"365"}},` +
		`{"pemicu":"metode","coverage":{"calculateMethodFacIn":"1","tsi":"1000000","rate":"0.5","day":"365"}}]}`
	kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage-aneka", badan, "")
	var j struct {
		ProRatePercent string           `json:"proRatePercent"`
		Baris          []map[string]any `json:"baris"`
	}
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || len(j.Baris) != 3 {
		t.Fatalf("%d %s", kode, isi)
	}
	if j.ProRatePercent != "100" || j.Baris[0]["premium"] != "2500" || j.Baris[1]["premium"] != "10000" ||
		j.Baris[2]["premium"] != "5000" {
		t.Errorf("rantai: %s", isi)
	}
	mau := "calculateMethod calculateMethodFacIn coverage coverageNote currency day deductibles discount discountPercentage indemnity " +
		"indemnityPercentage limitOfLiability loading lostLimit minPremium oldId pctLol pctShortPeriod premiNusantaraRe premium " +
		"proRatePercent rate subLimitNote tsi tsiLiability tsiNusantaraRe unit"
	if got := strings.Join(kunciAneka(t, j.Baris[0]), " "); got != mau {
		t.Errorf("kunci coverage: %s", got)
	}
	for nama, u := range map[string]struct {
		badan string
		kode  int
		pesan string
	}{
		"pemicu":     {`{"baris":[{"pemicu":"x","coverage":{}}]}`, 400, "baris[0].pemicu"},
		"desimal":    {`{"baris":[{"pemicu":"tsi","coverage":{"tsi":"1,5"}}]}`, 400, "baris[0].coverage.tsi"},
		"prorate":    {`{"proRatePercent":"x","baris":[]}`, 400, "proRatePercent"},
		"JSON rusak": {`{"baris":`, 400, "hitung coverage Aneka"},
	} {
		if kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage-aneka", u.badan, ""); kode != u.kode ||
			!strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	tanpa := services.Baru(nil).DenganKasus(kasusTiruan{err: repository.ErrKasusTidakAda})
	if kode, _ := minta(t, tanpa, "POST", "/api/nbfacin/kasus/UJI-NB-2/hitung-coverage-aneka", `{"baris":[]}`, ""); kode != 404 {
		t.Errorf("404: %d", kode)
	}
	// PUT objek-aneka: proRatePercent desimal tidak sah -> 400.
	var d []models.LokasiAneka
	put := services.Baru(nil).DenganObjekAneka(anekaTiruanH{&d}).DenganTransaksi(tanpaOracle)
	if kode, isi := minta(t, put, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-aneka", `{"lokasi":[],"proRatePercent":"1,5"}`, "UJI-USER"); kode != 400 ||
		!strings.Contains(isi, "proRatePercent") {
		t.Errorf("PUT prorate: %d %s", kode, isi)
	}
}

// TestHitungIndemnityAneka - pemicu indemnity lewat HTTP: kunci Indemnity + Unit ke M_BI_INDEMNITY; tanpa pembaca -> 503.
func TestHitungIndemnityAneka(t *testing.T) {
	k := &models.Kasus{CaseID: "UJI-NB-1", General: models.General{StartDateTime: "20260101T050000.000 GMT",
		EndDateTime: "20270101T050000.000 GMT"}}
	var kata string
	svc := services.Baru(nil).DenganKasus(kasusTiruan{k: k}).DenganLookupAneka(lookupAnekaTiruan{&kata})
	badan := `{"baris":[{"pemicu":"indemnity","coverage":{"indemnity":"30","unit":"1"}}]}`
	if kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage-aneka", badan, ""); kode != 200 ||
		kata != "30|1" || !strings.Contains(isi, `"indemnityPercentage":"80"`) {
		t.Errorf("%d %q %s", kode, kata, isi)
	}
	if kode, _ := minta(t, services.Baru(nil).DenganKasus(kasusTiruan{k: k}), "POST", "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage-aneka",
		badan, ""); kode != 503 {
		t.Errorf("tanpa pembaca: %d", kode)
	}
}
