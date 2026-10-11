package handlers

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// lokasiSpreadingAnekaUji - satu lokasi: okupasi OKU -> objek A (IDR, coverage C1 ber-baris OR 100 %), objek Golf G (USD).
func lokasiSpreadingAnekaUji() []models.LokasiAneka {
	return []models.LokasiAneka{{ObjectNo: "1", RiskLocation: "JL UJI", Okupasi: []models.OkupasiAneka{{OccupationName: "OKU",
		Objek: []models.ObjekAneka{{ObjectName: "A", Currency: "IDR", TSI: angkaUji("1000"), Coverages: []models.CoverageAneka{{OldID: "C1",
			CoverageNote: "UJI COV", TSI: angkaUji("1000"), Premium: angkaUji("10"), Discount: angkaUji("2"), LimitOfLiability: angkaUji("800"),
			Spreading: []models.BarisSpreading{{TreatyType: "10015", TreatyName: "FACOUT", SharePercentage: angkaUji("100")}}}}}}}},
		ObjekGolf: []models.ObjekAneka{{ObjectName: "G", Currency: "USD", TSI: angkaUji("200"), Coverages: []models.CoverageAneka{{OldID: "G1",
			TSI: angkaUji("200"), Premium: angkaUji("4")}}}}}}
}

// TestSpreadingAneka - tiket 73 P3: kunci kabel TampilanSpreadingAneka (kontrak 9d), hitung-share tidak menyimpan, PUT menerima
// tampilan utuh (kunci lain diabaikan) dan menulis pohon + % Share RNM; 400 / 401 / 404 / 409 / 503.
func TestSpreadingAneka(t *testing.T) {
	lokasi := lokasiSpreadingAnekaUji()
	src := &spreadingTiruanH{tersimpan: models.KasusSpreading{PercentShare: angkaUji("40"), StartDateTime: "20250101T170000.000 GMT"}}
	svc := services.Baru(nil).DenganObjekAneka(anekaTiruanH{d: &lokasi}).DenganSpreading(src).DenganTransaksi(tanpaOracleKlausa)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/spreading-aneka", "", "")
	if kode != 200 {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(isi), &v); err != nil {
		t.Fatal(err)
	}
	if k := kunciAneka(t, v); !reflect.DeepEqual(k, []string{"lokasi", "percentShare", "pesan", "ringkasanMataUang", "ringkasanTreaty",
		"tampilTotalLokasi", "template", "treaty"}) || v["percentShare"] != "40" || v["tampilTotalLokasi"] != false {
		t.Fatalf("tampilan: %v", k)
	}
	l := v["lokasi"].([]any)[0].(map[string]any)
	if k := kunciAneka(t, l); !reflect.DeepEqual(k, []string{"location", "objectNo", "objekGolf", "okupasi", "total"}) || l["location"] != "JL UJI" {
		t.Errorf("lokasi: %v", k)
	}
	ok := l["okupasi"].([]any)[0].(map[string]any)
	o := ok["objek"].([]any)[0].(map[string]any)
	if k := kunciAneka(t, ok); !reflect.DeepEqual(k, []string{"objek", "occupationName"}) || ok["occupationName"] != "OKU" {
		t.Errorf("okupasi: %v", k)
	}
	if k := kunciAneka(t, o); !reflect.DeepEqual(k, []string{"coverages", "currency", "objectName", "tsi", "tsiOld"}) || o["tsiOld"] != "" {
		t.Errorf("objek: %v", k)
	}
	c := o["coverages"].([]any)[0].(map[string]any)
	if k := kunciAneka(t, c); !reflect.DeepEqual(k, []string{"coverageNote", "discount", "limitOfLiability", "loading", "oldId", "premiNusantaraRe",
		"premium", "rate", "spreading", "tsi", "tsiNusantaraRe"}) || c["coverageNote"] != "UJI COV" {
		t.Errorf("coverage: %v", k)
	}
	b := c["spreading"].([]any)[0].(map[string]any)
	if k := kunciAneka(t, b); !reflect.DeepEqual(k, []string{"claimEstimation", "premiumSpreaded", "sharePercentage", "treatyName", "treatyType",
		"tsiGrossSpreaded", "tsiSpreaded"}) {
		t.Errorf("baris: %v", k)
	}
	if g := l["objekGolf"].([]any)[0].(map[string]any); len(g["coverages"].([]any)[0].(map[string]any)["spreading"].([]any)) != 0 {
		t.Errorf("Golf: %v", g)
	}
	kode, isi = minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/spreading-aneka/hitung-share", `{"percentShare":"50"}`, "")
	if kode != 200 || !strings.Contains(isi, `"tsiNusantaraRe":"500","premiNusantaraRe":"4"`) || !strings.Contains(isi, `"percentShare":"50"`) ||
		lokasi[0].Okupasi[0].Objek[0].Coverages[0].TSINusantaraRe != nil {
		t.Fatalf("hitung-share: %d %s", kode, isi)
	}
	// PUT: tampilan utuh + kunci tak dikenal; baris badan share 40.
	badan := strings.Replace(strings.Replace(isi, `"sharePercentage":"100"`, `"sharePercentage":"40","lain":1`, 1), `"percentShare":"50"`,
		`"percentShare":"50","kunciLain":true`, 1)
	kode, isi = minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/spreading-aneka", badan, "UJI-USER")
	if kode != 200 || lokasi[0].Okupasi[0].Objek[0].Coverages[0].Spreading[0].SharePercentage.Cmp(angkaUji("40")) != 0 ||
		lokasi[0].Okupasi[0].Objek[0].Coverages[0].TSINusantaraRe.Cmp(angkaUji("500")) != 0 {
		t.Fatalf("PUT: %d %s / %v", kode, isi, lokasi[0].Okupasi[0].Objek[0].Coverages[0])
	}
	for nama, u := range map[string]struct {
		metode, jalur, badan, pelaku, mau string
		kode                              int
	}{
		"tanpa percentShare":  {"PUT", "/spreading-aneka", `{"lokasi":[]}`, "UJI-USER", "percentShare", 400},
		"tanpa lokasi":        {"PUT", "/spreading-aneka", `{"percentShare":"50"}`, "UJI-USER", "lokasi", 400},
		"share huruf":         {"PUT", "/spreading-aneka", strings.Replace(badan, `"sharePercentage":"40"`, `"sharePercentage":"x"`, 1), "UJI-USER", "lokasi[0].okupasi[0].objek[0].coverages[0].spreading[0].sharePercentage", 400},
		"bentuk berubah":      {"PUT", "/spreading-aneka", `{"percentShare":"50","lokasi":[]}`, "UJI-USER", "muat ulang", 409},
		"tanpa identitas":     {"PUT", "/spreading-aneka", badan, "", "", 401},
		"salin identitas":     {"POST", "/spreading-aneka/salin", `{"percentShare":"50","template":[]}`, "", "", 401},
		"salin tanpa templat": {"POST", "/spreading-aneka/salin", `{"percentShare":"50","template":[]}`, "UJI-USER", "template", 400},
		"hitung-baris":        {"POST", "/spreading-aneka/hitung-baris", badan, "UJI-USER", `"percentShare":"50"`, 200},
	} {
		if kode, isi := minta(t, svc, u.metode, "/api/nbfacin/kasus/UJI-NB-1"+u.jalur, u.badan, u.pelaku); kode != u.kode || !strings.Contains(isi, u.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/K404/spreading-aneka", "", ""); kode != 404 {
		t.Errorf("case tidak ada: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1/spreading-aneka", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}
