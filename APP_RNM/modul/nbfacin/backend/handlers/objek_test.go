package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type objekTiruan struct{ d *[]models.ObjekFire }

func (o objekTiruan) BacaObjek(context.Context, string) ([]models.ObjekFire, error) {
	return append([]models.ObjekFire{}, *o.d...), nil
}

func (o objekTiruan) GantiObjek(_ context.Context, _ *db.Tx, _ string, baris []models.ObjekFire) error {
	*o.d = baris
	return nil
}

// TestObjek - GET/PUT /api/nbfacin/kasus/{caseId}/objek (tiket 35, 38, 39, 40, 41, 42): 34 kunci persis kontrak
// ObjekFire, boolean JSON, `baris` larik walau kosong; PUT 200 baca ulang, 400/401/503.
func TestObjek(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle)
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/objek", "", ""); kode != 200 || isi != `{"baris":[]}` {
		t.Fatalf("kosong: %d %s", kode, isi)
	}
	badan := `{"baris":[{"objectNo":"1","objectType":"UJI TIPE","isTopRisk":true,"isMaterialDamage":false,"numberOfFloor":"3","otherType":"UJI LAIN",` +
		`"ownership":"2","isHotWorkProcess":true,"surroundingRisk":{"left":{"occupation":"UJI01","construction":"UJI KONSTRUKSI",` +
		`"distance":"12.50","note":"UJI PABRIK"},"floodAreaStatus":"2","housekeepingRemark":"UJI CATATAN"}}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", badan, "UJI-USER")
	var j struct {
		Baris []map[string]any `json:"baris"`
	}
	kunciMau := []string{"objectNo", "objectType", "objectName", "isMaterialDamage", "isTopRisk", "roadType", "roadName",
		"buildingNo", "zipCode", "country", "riskLocation", "territory", "city", "district", "province", "riskAddressId",
		"numberOfFloor", "roofType", "wallType", "floorType", "partitionType", "supportWallType", "otherType",
		"ownership", "isProductionProcess", "isHotWorkProcess", "isFlammableItem", "surroundingRisk", "items", "occupations", "fea", "lossRecords", "lossRatio", "internalLossRecords"}
	if err := json.Unmarshal([]byte(isi), &j); err == nil && len(j.Baris) == 1 {
		for _, k := range kunciMau {
			if _, ada := j.Baris[0][k]; !ada {
				t.Errorf("kunci %q tidak ada (kontrak ObjekFire)", k)
			}
		}
	}
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || len(j.Baris) != 1 || len(j.Baris[0]) != len(kunciMau) ||
		j.Baris[0]["isTopRisk"] != true || j.Baris[0]["otherType"] != "UJI LAIN" || j.Baris[0]["numberOfFloor"] != "3" {
		t.Fatalf("%d %s", kode, isi)
	}
	for _, u := range []struct {
		nama, badan, pelaku string
		svc                 *services.Service
		kode                int
	}{
		{"tanpa identitas", badan, "", svc, 401},
		{"JSON rusak", `{"baris":[`, "UJI-USER", svc, 400},
		{"tipe kosong", `{"baris":[{"objectType":""}]}`, "UJI-USER", svc, 400},
		{"jarak tidak sah", `{"baris":[{"objectType":"UJI","surroundingRisk":{"back":{"distance":"1,5"}}}]}`, "UJI-USER", svc, 400},
		{"tanpa DB", badan, "UJI-USER", services.Baru(nil), 503},
	} {
		if kode, isi := minta(t, u.svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s: %d %s, mau %d", u.nama, kode, isi, u.kode)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1/objek", "", ""); kode != 503 {
		t.Errorf("GET tanpa DB: %d", kode)
	}
}

// TestObjekSurroundingRisk - tiket 38: medan bersarang pulang-pergi utuh lewat pemetaan
// eksplisit, bentuk persis kontrak SurroundingRisk / SisiRisiko; pesan 400 menyebut sisi.
func TestObjekSurroundingRisk(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle)
	sisi := `{"occupation":"UJI01","construction":"UJI K","distance":"0","note":"UJI N"}`
	sekitar := `{"front":` + sisi + `,"left":` + sisi + `,"back":` + sisi + `,"right":` + sisi +
		`,"housekeepingStatus":"0","floodAreaStatus":"2","floodArea":"UJI AREA","housekeepingRemark":"UJI R"}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek",
		`{"baris":[{"objectType":"UJI","ownership":"1","isProductionProcess":true,"isFlammableItem":true,"surroundingRisk":`+sekitar+`}]}`, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"surroundingRisk":`+sekitar+`,"items":[],"occupations":[],"fea":[],"lossRecords":[],"lossRatio":{"oneYearAmount":"0","oneYearPercent":"0","threeFiveYearAmount":"0","threeFiveYearPercent":"0"},"internalLossRecords":[]}`) ||
		!strings.Contains(isi, `"ownership":"1","isProductionProcess":true,"isHotWorkProcess":false,"isFlammableItem":true`) {
		t.Fatalf("%d %s", kode, isi)
	}
	if len(d) != 1 || d[0].SurroundingRisk.Right.Note != "UJI N" || d[0].SurroundingRisk.FloodArea != "UJI AREA" {
		t.Errorf("model: %+v", d)
	}
	kode, isi = minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek",
		`{"baris":[{"objectType":"UJI","surroundingRisk":{"front":{"distance":"100000.01"}}}]}`, "UJI-USER")
	if kode != 400 || !strings.Contains(isi, "baris[0].surroundingRisk.front.distance") {
		t.Errorf("jarak: %d %s", kode, isi)
	}
}

type occupationTiruan struct{}

func (occupationTiruan) CariOccupation(_ context.Context, k string) ([]models.BarisOccupation, error) {
	return []models.BarisOccupation{{OldID: "UJI01", Name: "UJI " + k, KdRiskExposure: "03"}}, nil
}

// TestOccupation - GET /api/nbfacin/occupation?cari= (tiket 38): bentuk persis kontrak
// BarisOccupation; tanpa identitas; 400 < 2 karakter; 503.
func TestOccupation(t *testing.T) {
	svc := services.Baru(nil).DenganOccupation(occupationTiruan{})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/occupation?cari=pa", "", ""); kode != 200 || isi != `{"baris":[{"oldId":"UJI01","name":"UJI pa","kdRiskExposure":"03"}]}` {
		t.Errorf("%d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/occupation", "", ""); kode != 200 || isi != `{"baris":[{"oldId":"UJI01","name":"UJI ","kdRiskExposure":"03"}]}` {
		t.Errorf("tanpa cari (tiket 40): %d %s", kode, isi)
	}
	for jalur, mau := range map[string]int{"/api/nbfacin/occupation?cari=p": 400} {
		if kode, isi := minta(t, svc, "GET", jalur, "", ""); kode != mau || !strings.Contains(isi, `"galat"`) {
			t.Errorf("%s: %d %s", jalur, kode, isi)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/occupation?cari=pa", "", ""); kode != 503 {
		t.Errorf("tanpa DB: %d", kode)
	}
}

// TestObjekItem - tiket 39: items pulang-pergi utuh (kunci persis kontrak ItemObjek, uang teks), items yang tidak
// dikirim = [] ; 400 ber-indeks; mata uang diperiksa ke daftar.
func TestObjekItem(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle).DenganPilihanItem(pilihanTiruan{})
	item := `{"itemTypeId":"UJI01","itemType":"UJI MESIN","note":"UJI KET","propertyYear":"2020","unit":"2","condition":"UJI",` +
		`"currency":"IDR","tsi":"1500000000.12345678","yearOfPlanting":"","noOfTree":"","areaHectar":"","remark":"UJI R",` +
		`"isAdjustable":true,"pctAdjust2":"","pctAdjustOther":"75.5"}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", `{"baris":[{"objectType":"UJI","items":[`+item+`]},{"objectType":"UJI"}]}`, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"items":[`+item+`]`) || !strings.HasSuffix(isi, `"items":[],"occupations":[],"fea":[],"lossRecords":[],"lossRatio":{"oneYearAmount":"0","oneYearPercent":"0","threeFiveYearAmount":"0","threeFiveYearPercent":"0"},"internalLossRecords":[]}]}`) {
		t.Fatalf("%d %s", kode, isi)
	}
	if len(d) != 2 || len(d[0].Items) != 1 || d[0].Items[0].TSI != "1500000000.12345678" || d[1].Items == nil {
		t.Errorf("model: %+v", d)
	}
	for nama, u := range map[string]struct {
		badan, pesan string
	}{
		"tsi angka JSON": {`{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","tsi":15}]}]}`, "bukan JSON"},
		"tsi koma":       {`{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","tsi":"1,5"}]}]}`, "baris[0].items[0].tsi"},
		"mata uang ITL":  {`{"baris":[{"objectType":"UJI","items":[{"currency":"ITL"}]}]}`, `baris[0].items[0].currency \"ITL\"`},
	} {
		if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", u.badan, "UJI-USER"); kode != 400 || !strings.Contains(isi, u.pesan) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
}

type pilihanTiruan struct{}

func (pilihanTiruan) DaftarJenisItem(context.Context) ([]models.JenisItem, error) {
	return []models.JenisItem{{Kode: "UJI01", Nama: "UJI MESIN", Keterangan: "UJI KET"}}, nil
}

func (pilihanTiruan) DaftarMataUang(context.Context) ([]string, error) {
	return []string{"IDR", "USD"}, nil
}

// TestPilihanItem - GET jenis-item-objek dan mata-uang (tiket 39): bentuk persis kontrak; tanpa identitas; 503.
func TestPilihanItem(t *testing.T) {
	svc := services.Baru(nil).DenganPilihanItem(pilihanTiruan{})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/jenis-item-objek", "", ""); kode != 200 || isi != `{"baris":[{"kode":"UJI01","nama":"UJI MESIN","keterangan":"UJI KET"}]}` {
		t.Errorf("jenis item: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/mata-uang", "", ""); kode != 200 || isi != `{"baris":["IDR","USD"]}` {
		t.Errorf("mata uang: %d %s", kode, isi)
	}
	for _, jalur := range []string{"/api/nbfacin/jenis-item-objek", "/api/nbfacin/mata-uang"} {
		if kode, _ := minta(t, services.Baru(nil), "GET", jalur, "", ""); kode != 503 {
			t.Errorf("%s tanpa DB: %d", jalur, kode)
		}
	}
}

// TestObjekOkupasi - tiket 40: occupations pulang-pergi utuh (kunci persis OkupasiObjek; pctLimit teks apa adanya
// termasuk koma dan spasi ujung, butir 68.1); occupations tidak dikirim = [].
func TestObjekOkupasi(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle)
	ok := `{"occupationId":"UJI01","occupationName":"UJI PABRIK","category":"III","constructionClass":"UJI KELAS","pctLimit":"70,000 "}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", `{"baris":[{"objectType":"UJI","occupations":[`+ok+`]},{"objectType":"UJI"}]}`, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"occupations":[`+ok+`]`) || !strings.HasSuffix(isi, `"occupations":[],"fea":[],"lossRecords":[],"lossRatio":{"oneYearAmount":"0","oneYearPercent":"0","threeFiveYearAmount":"0","threeFiveYearPercent":"0"},"internalLossRecords":[]}]}`) {
		t.Fatalf("%d %s", kode, isi)
	}
	if len(d) != 2 || len(d[0].Occupations) != 1 || d[0].Occupations[0].PctLimit != "70,000 " {
		t.Errorf("model: %+v", d)
	}
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek",
		`{"baris":[{"objectType":"UJI","occupations":[{"category":"`+strings.Repeat("I", 51)+`"}]}]}`, "UJI-USER"); kode != 400 ||
		!strings.Contains(isi, "baris[0].occupations[0].category") {
		t.Errorf("lebar: %d %s", kode, isi)
	}
}

// TestObjekFEA - tiket 41: fea pulang-pergi utuh (kunci persis BarisFEA); fea tidak dikirim = []; unit tidak sah 400.
func TestObjekFEA(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle)
	f := `{"apar":"2","sprinkler":"0","smokeDetector":"","hydrant":"3","privateTruckBrigade":"1","privateFireBrigade":"UJI",` +
		`"teamSopSafety":"","teamSopRiskManagement":"","info":"UJI INFO"}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", `{"baris":[{"objectType":"UJI","fea":[`+f+`]},{"objectType":"UJI"}]}`, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"fea":[`+f+`]`) || !strings.HasSuffix(isi, `"fea":[],"lossRecords":[],"lossRatio":{"oneYearAmount":"0","oneYearPercent":"0","threeFiveYearAmount":"0","threeFiveYearPercent":"0"},"internalLossRecords":[]}]}`) {
		t.Fatalf("%d %s", kode, isi)
	}
	if len(d) != 2 || len(d[0].FEA) != 1 || d[0].FEA[0].PrivateTruckBrigade != "1" {
		t.Errorf("model: %+v", d)
	}
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", `{"baris":[{"objectType":"UJI","fea":[{"apar":"-1"}]}]}`,
		"UJI-USER"); kode != 400 || !strings.Contains(isi, "baris[0].fea[0].apar") {
		t.Errorf("unit: %d %s", kode, isi)
	}
}

// TestObjekKerugian - tiket 42: lossRecords pulang-pergi (kunci persis CatatanKerugian); coinsName = tertanggung case;
// lossRatio dihitung server (badan PUT diabaikan); internalLossRecords selalu []; tanggal tak sah 400.
func TestObjekKerugian(t *testing.T) {
	var d []models.ObjekFire
	k := &models.Kasus{CaseID: "UJI-NB-1", InsuredName: "UJI TERTANGGUNG"}
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle).DenganPilihanItem(pilihanTiruan{}).
		DenganKasus(kasusTiruan{k: k})
	c := `{"dateOfLoss":"","coinsName":"UJI TERTANGGUNG","lossObject":"UJI OBJEK","currency":"IDR","amount":"400","claim":"100",` +
		`"preventionOfLoss":"0.5","causeOfLoss":"UJI SEBAB","remarks":"","detail":"UJI DETAIL"}`
	badan := `{"baris":[{"objectType":"UJI","lossRatio":{"oneYearAmount":"9"},"internalLossRecords":[{"location":"x"}],` +
		`"lossRecords":[` + strings.Replace(c, `"coinsName":"UJI TERTANGGUNG"`, `"coinsName":"UJI LAIN"`, 1) + `]}]}`
	kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", badan, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"lossRecords":[`+c+`]`) || !strings.Contains(isi, `"internalLossRecords":[]`) ||
		!strings.Contains(isi, `"lossRatio":{"oneYearAmount":"0","oneYearPercent":"0","threeFiveYearAmount":"0","threeFiveYearPercent":"0"}`) {
		t.Fatalf("%d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek",
		`{"baris":[{"objectType":"UJI","lossRecords":[{"currency":"IDR","dateOfLoss":"31-02-2026"}]}]}`, "UJI-USER"); kode != 400 ||
		!strings.Contains(isi, "baris[0].lossRecords[0].dateOfLoss") {
		t.Errorf("tanggal: %d %s", kode, isi)
	}
}
