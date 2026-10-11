package handlers

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// cargoTiruanH - PenyimpanObjekCargo dalam memori (satu case "UJI-NB-1").
type cargoTiruanH struct{ d *[]models.ObjekCargo }

func (c cargoTiruanH) BacaObjekCargo(_ context.Context, id string) ([]models.ObjekCargo, error) {
	if id != "UJI-NB-1" {
		return nil, repository.ErrKasusTidakAda
	}
	return append([]models.ObjekCargo{}, *c.d...), nil
}

func (c cargoTiruanH) GantiObjekCargo(_ context.Context, _ *db.Tx, _ string, cargo []models.ObjekCargo) error {
	*c.d = append([]models.ObjekCargo{}, cargo...)
	for i := range *c.d {
		(*c.d)[i].ID = "7" + string(rune('0'+i))
	}
	return nil
}

func (cargoTiruanH) PercentShareKasus(context.Context, string) (*apd.Decimal, error) {
	return apd.New(25, 0), nil
}

func (cargoTiruanH) UbahPercentShare(context.Context, *db.Tx, string, *apd.Decimal) error { return nil }

func (cargoTiruanH) MataUangCargo(_ context.Context, nama string) (models.MataUangCsv, bool, error) {
	if nama == "IDR" {
		return models.MataUangCsv{ID: "10026", OldID: "1", Currency: "IDR"}, true, nil
	}
	return models.MataUangCsv{}, false, nil
}

type masterCargoH struct{}

func (masterCargoH) Deskripsi(_ context.Context, jenis, id string) (string, bool, error) {
	return strings.ToUpper(jenis) + " " + id, true, nil
}
func (masterCargoH) WarrantyMaxAge(context.Context, string) (string, bool, error) {
	return "W", true, nil
}
func (masterCargoH) IDMaster(_ context.Context, _, deskripsi string) (string, bool, error) {
	return "ID-" + deskripsi, deskripsi != "", nil
}
func (masterCargoH) CariAgen(_ context.Context, kata string) ([]models.BarisAgenSurvei, error) {
	return []models.BarisAgenSurvei{{ID: "A1", Name: kata, Contact: "C"}}, nil
}
func (masterCargoH) CariMaxAge(_ context.Context, kata string) ([]models.BarisMaxAge, error) {
	return []models.BarisMaxAge{{ID: "M1", Warranty: kata}}, nil
}
func (masterCargoH) CariKondisiMarine(_ context.Context, kata string) ([]models.BarisKondisiMarine, error) {
	return nil, nil
}
func (masterCargoH) Daftar(_ context.Context, jenis string) ([]models.BarisMasterCargo, error) {
	return []models.BarisMasterCargo{{ID: "1", Description: strings.ToUpper(jenis) + " A"}}, nil
}
func (masterCargoH) Agen(_ context.Context, id string) (models.AgenSurvei, bool, error) {
	return models.AgenSurvei{AgentID: id}, true, nil
}

// TestMasterCargoKabel - GET /api/nbfacin/marine/{jenis}: master {baris: [{id, description}]}; agent / max-age / coverage bentuk
// api.ts; tanpa master terpasang semua 503.
func TestMasterCargoKabel(t *testing.T) {
	svc := layananCargoH("004", nil, true)
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/marine/trading", "", ""); kode != 200 || isi != `{"baris":[{"id":"1","description":"TRADING A"}]}` {
		t.Errorf("%d %s", kode, isi)
	}
	for jalur, mau := range map[string]string{
		"agent?q=Uji":    `{"baris":[{"id":"A1","name":"Uji","address":"","city":"","country":"","contact":"C"}]}`,
		"max-age?q=25":   `{"baris":[{"id":"M1","warranty":"25"}]}`,
		"coverage?q=ICC": `{"baris":[]}`,
	} {
		if kode, isi := minta(t, svc, "GET", "/api/nbfacin/marine/"+jalur, "", ""); kode != 200 || isi != mau {
			t.Errorf("%s: %d %s", jalur, kode, isi)
		}
	}
	tanpa := layananCargoH("004", nil, false)
	for _, j := range []string{"agent", "max-age", "coverage", "trading", "tipe-kapal"} {
		if kode, _ := minta(t, tanpa, "GET", "/api/nbfacin/marine/"+j, "", ""); kode != 503 {
			t.Errorf("tanpa master %s: %d", j, kode)
		}
	}
}

func layananCargoH(gp string, d *[]models.ObjekCargo, master bool) *services.Service {
	k := &models.Kasus{CaseID: "UJI-NB-1", Opportunity: models.Opportunity{ClassOfBusiness: "UJI KELAS", GroupBusinessID: "3"},
		General: models.General{StartDateTime: "20260101T050000.000 GMT", EndDateTime: "20270101T050000.000 GMT"}}
	svc := services.Baru(nil).DenganObjekCargo(cargoTiruanH{d}).DenganTransaksi(tanpaOracle).DenganKasus(kasusTiruan{k: k}).
		DenganTableOfLimit(tolTiruan{bisnis: []models.BarisBisnisLini{{ID: "10999", OldID: "24", GroupPanel: gp}}})
	if master {
		svc.DenganMasterCargo(masterCargoH{})
	}
	return svc
}

func kunciPeta(m map[string]any) []string {
	k := make([]string, 0, len(m))
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

// TestObjekCargo - GET/PUT …/objek-cargo (tiket 74): kunci kabel PERSIS kontrak tiap tingkat, larik walau kosong, medan server
// dari server (badan diabaikan), tanggal DD-MM-YYYY, coverage sintesis; 400 / 401 / 404 / 409 / 503.
func TestObjekCargo(t *testing.T) {
	var d []models.ObjekCargo
	svc := layananCargoH("004", &d, true)
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek-cargo", "", ""); kode != 200 || isi != `{"cargo":[]}` {
		t.Fatalf("kosong: %d %s", kode, isi)
	}
	badan := `{"cargo":[{"id":"","isTopRisk":true,"tradingId":"T1","tradingNote":"KLIEN","goodId":"G1","conveyanceId":"C1",
		"packingId":"P1","fromRute":"A","toRute":"B","sailDate":"15-03-2026","invoiceDate":"","maxAge":"M1","maxAgeDetail":"KLIEN",
		"surveyAgent":{"agentId":"A1","fullName":"UJI AGEN","address":"","city":"","country":"ID"},
		"lc":{"lc":true,"bank":"B","number":"","condition":"","remark":"","endDate":"31-12-2026"},
		"ship":{"id":"S1","nmShip":"UJI","shipTypeId":"","shipTypeName":"","grt":"1,5","nrt":"","dwt":"","age":"7","flagName":"",
		"construction":"","yMake1":"2001","yMake2":"","remark":"","register":"","imo":"","formerShipName":"","classification":"","rebuilt":""},
		"additionalShip":[{"id":"","register":"R1","nmShip":"","voyageNo":"V1","shipTypeId":"","flagName":"","nrt":"","dwt":"","grt":"","yMake1":"","construction":""}],
		"flagDelete":"1","totalGross":[{"currency":"X","tsi":"1","premium":"1","rate":"1"}],
		"coverages":[{"coverage":"ICC A","coverageNote":"N","coverageInitial":"I","currency":"IDR","currencyId":"999","currencyOldId":"9",
		"rate":"0.1","limitOfLiability":"","discountPercentage":"10","tsi":"1000","premium":"1","minPremium":"","discount":"5",
		"tsiLiability":"5","tsiNusantaraRe":"5","premiNusantaraRe":"5","proRatePercent":"5",
		"warrantyExtension":[{"coverage":"W1","coverageNote":"WN"}],
		"deductibles":[{"typeDeductible":"1","pctDeductible":"10","minMax":"","currency":"","typeDeductible2":"","pctDeductible2":"",
		"condition":"","amount":"","inputCondition":"","timeExcess":"","descriptions":"DESK"}]}]}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-cargo", badan, "UJI-USER")
	if kode != 200 {
		t.Fatalf("PUT: %d %s", kode, isi)
	}
	var out struct {
		Cargo []map[string]any `json:"cargo"`
	}
	if err := json.Unmarshal([]byte(isi), &out); err != nil || len(out.Cargo) != 1 {
		t.Fatalf("%v %s", err, isi)
	}
	c := out.Cargo[0]
	mau := []string{"additionalShip", "blNumber", "conveyanceId", "conveyanceName", "conveyanceNote", "coverages", "flagDelete", "fromRute",
		"goodId", "goodNote", "id", "invoiceDate", "invoiceNumber", "isTopRisk", "lc", "maxAge", "maxAgeDetail", "packingId", "packingNote",
		"sailDate", "ship", "surveyAgent", "toRute", "totalGross", "tradingId", "tradingNote"}
	if !reflect.DeepEqual(kunciPeta(c), mau) {
		t.Errorf("kunci cargo: %v", kunciPeta(c))
	}
	ship := c["ship"].(map[string]any)
	if !reflect.DeepEqual(kunciPeta(ship), []string{"age", "classification", "construction", "dwt", "flagName", "formerShipName", "grt",
		"id", "imo", "nmShip", "nrt", "rebuilt", "register", "remark", "shipTypeId", "shipTypeName", "yMake1", "yMake2"}) {
		t.Errorf("kunci ship: %v", kunciPeta(ship))
	}
	cov := c["coverages"].([]any)[0].(map[string]any)
	if !reflect.DeepEqual(kunciPeta(cov), []string{"coverage", "coverageInitial", "coverageNote", "currency", "currencyId", "currencyOldId",
		"deductibles", "discount", "discountPercentage", "limitOfLiability", "minPremium", "premiNusantaraRe", "premium", "proRatePercent",
		"rate", "tsi", "tsiLiability", "tsiNusantaraRe", "warrantyExtension"}) {
		t.Errorf("kunci coverage: %v", kunciPeta(cov))
	}
	if c["tradingNote"] != "TRADING T1" || c["conveyanceName"] != "CONVEYANCE C1" || c["maxAgeDetail"] != "W" || c["flagDelete"] != "" ||
		c["sailDate"] != "15-03-2026" || ship["age"] != "25" || ship["grt"] != "1.5" ||
		c["lc"].(map[string]any)["endDate"] != "31-12-2026" || len(c["additionalShip"].([]any)) != 1 {
		t.Errorf("objek: %s", isi)
	}
	// Premium dikirim 1 (isian) dibiarkan; medan server dihitung ulang: Discount = 1×10/100, RNM PS 25.
	if cov["currencyId"] != "10026" || cov["discount"] != "0.1" || cov["tsiLiability"] != "1000" || cov["tsiNusantaraRe"] != "250" ||
		cov["premiNusantaraRe"] != "0.225" || cov["proRatePercent"] != "100" ||
		cov["deductibles"].([]any)[0].(map[string]any)["descriptions"] != "DESK" {
		t.Errorf("coverage: %v", cov)
	}
	tg := c["totalGross"].([]any)
	if len(tg) != 1 || tg[0].(map[string]any)["currency"] != "IDR" || tg[0].(map[string]any)["rate"] != "0.1" {
		t.Errorf("totalGross: %v", tg)
	}
	// coverages absen -> coverage lama cargo ber-id sama dipertahankan.
	badan2 := `{"cargo":[{"id":"70","tradingId":"T1","goodId":"G1","conveyanceId":"C1","packingId":"P1","fromRute":"A","toRute":"B"}]}`
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-cargo", badan2, "UJI-USER"); kode != 200 ||
		!strings.Contains(isi, `"coverage":"ICC A"`) {
		t.Errorf("pertahankan: %d %s", kode, isi)
	}
	for _, k := range []struct {
		svc    *services.Service
		badan  string
		pelaku string
		kode   int
		jalur  string
	}{
		{svc, `{"cargo":[{"tradingId":"T1","goodId":"G1","conveyanceId":"C1","packingId":"P1","fromRute":"A","toRute":"B","coverages":[{"tsi":"1,5"}]}]}`,
			"UJI-USER", 400, "cargo[0].coverages[0].tsi"},
		{svc, `{"cargo":`, "UJI-USER", 400, ""},
		{svc, `{"cargo":[]}`, "", 401, ""},
		{layananCargoH("006", &d, true), `{"cargo":[]}`, "UJI-USER", 409, ""},
		{layananCargoH("004", &d, false), `{"cargo":[{"tradingId":"T9","goodId":"G1","conveyanceId":"C1","packingId":"P1","fromRute":"A","toRute":"B"}]}`,
			"UJI-USER", 503, ""},
	} {
		kode, isi := minta(t, k.svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek-cargo", k.badan, k.pelaku)
		if kode != k.kode || !strings.Contains(isi, k.jalur) {
			t.Errorf("%s: %d %s", k.badan, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/LAIN/objek-cargo", "", ""); kode != 404 {
		t.Errorf("404: %d", kode)
	}
}

// TestHitungCoverageCargoKabel - POST …/hitung-coverage-cargo: urutan sama, medan server terisi; pemicu tak dikenal 400.
func TestHitungCoverageCargoKabel(t *testing.T) {
	var d []models.ObjekCargo
	svc := layananCargoH("004", &d, false)
	kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage-cargo",
		`{"baris":[{"coverage":{"currency":"IDR","tsi":"1000","rate":"0.5"},"pemicu":"currency"},{"coverage":{"tsi":"3000","premium":"1000"},"pemicu":"premium"}]}`, "")
	if kode != 200 || !strings.Contains(isi, `"premium":"5"`) || !strings.Contains(isi, `"currencyId":"10026"`) ||
		!strings.Contains(isi, `"rate":"33.3333"`) || !strings.Contains(isi, `"tsiNusantaraRe":"250"`) {
		t.Errorf("hitung: %d %s", kode, isi)
	}
	if kode, _ := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage-cargo", `{"baris":[{"coverage":{},"pemicu":"x"}]}`, ""); kode != 400 {
		t.Errorf("pemicu: %d", kode)
	}
}

// TestSpreadingCargo - GET / PUT / hitung-share …/spreading-cargo (tiket 74 P3): kunci kabel persis; PUT membawa baris (share +
// tsiGrossSpreaded / claimEstimation isian) lalu menulis pohon cargo; bentuk beda / bukan Marine 409.
func TestSpreadingCargo(t *testing.T) {
	d := []models.ObjekCargo{{ID: "70", TradingNote: "DOMESTIK", Coverages: []models.CoverageCargo{{Coverage: "ICC A", Currency: "IDR",
		TSI: apd.New(1000, 0), TSILiability: apd.New(1000, 0), Premium: apd.New(10, 0), LimitOfLiability: apd.New(500, 0),
		Spreading: []models.BarisSpreading{{TreatyType: "10015", SharePercentage: apd.New(100, 0)}}}}}}
	sp := &spreadingTiruanH{tersimpan: models.KasusSpreading{PercentShare: apd.New(40, 0), StartDateTime: "20260101T050000.000 GMT"}}
	svc := layananCargoH("004", &d, true).DenganSpreading(sp)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/spreading-cargo", "", "")
	if kode != 200 {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(isi), &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(kunciPeta(out), []string{"cargo", "labelLimitOfLiability", "percentShare", "pesan", "ringkasanMataUang",
		"ringkasanTreaty", "template", "treaty"}) || out["labelLimitOfLiability"] != "Total Limit of Liability(50.0000%)" {
		t.Errorf("kunci: %v %v", kunciPeta(out), out["labelLimitOfLiability"])
	}
	c := out["cargo"].([]any)[0].(map[string]any)
	cov := c["coverages"].([]any)[0].(map[string]any)
	if !reflect.DeepEqual(kunciPeta(c), []string{"conveyanceNote", "coverages", "goodNote", "no", "packingNote", "tradingNote"}) ||
		!reflect.DeepEqual(kunciPeta(cov), []string{"coverage", "coverageNote", "currency", "discount", "limitOfLiability", "premiNusantaraRe",
			"premium", "rate", "spreading", "tsi", "tsiLiability", "tsiNusantaraRe"}) {
		t.Errorf("kunci cargo / coverage: %v / %v", kunciPeta(c), kunciPeta(cov))
	}
	badan := `{"percentShare":"40","cargo":[{"coverages":[{"coverage":"ICC A","spreading":[{"treatyType":"10015","treatyName":"FACOUT",
		"sharePercentage":"100","tsiGrossSpreaded":"7","tsiSpreaded":"9","claimEstimation":"123","premiumSpreaded":"9"}]}]}]}`
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/spreading-cargo", badan, "UJI-USER"); kode != 200 ||
		!strings.Contains(isi, `"claimEstimation":"123"`) || !strings.Contains(isi, `"tsiGrossSpreaded":"7"`) ||
		!strings.Contains(isi, `"tsiSpreaded":"400"`) {
		t.Errorf("PUT: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/spreading-cargo/hitung-share", `{"percentShare":"50"}`, ""); kode != 200 ||
		!strings.Contains(isi, `"tsiNusantaraRe":"500"`) {
		t.Errorf("hitung-share: %d %s", kode, isi)
	}
	beda := `{"percentShare":"40","cargo":[{"coverages":[{"coverage":"LAIN","spreading":[]}]}]}`
	if kode, _ := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/spreading-cargo", beda, "UJI-USER"); kode != 409 {
		t.Errorf("bentuk: %d", kode)
	}
	if kode, _ := minta(t, layananCargoH("006", &d, true).DenganSpreading(sp), "GET", "/api/nbfacin/kasus/UJI-NB-1/spreading-cargo", "", ""); kode != 409 {
		t.Errorf("FIRE: %d", kode)
	}
	if kode, _ := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/spreading-cargo", `{"percentShare":"40"}`, "UJI-USER"); kode != 400 {
		t.Errorf("tanpa cargo: %d", kode)
	}
}
