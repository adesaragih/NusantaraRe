package handlers

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// zonaTiruan - mencatat panggilan; rate UJI hanya untuk zona "UJI-ZONA-2".
type zonaTiruan struct{ catat *string }

func (z zonaTiruan) SaranZona(_ context.Context, oldID string, zona []string, kata string) ([]models.SaranAkumulasi, error) {
	*z.catat = oldID + "|" + strings.Join(zona, ",") + "|" + kata
	return []models.SaranAkumulasi{{ID: "Z2", Label: "UJI-ZONA-2"}}, nil
}

func (z zonaTiruan) SaranKonstruksiGempa(_ context.Context, kata string) ([]models.SaranAkumulasi, error) {
	*z.catat = "konstruksi|" + kata
	return []models.SaranAkumulasi{{ID: "UJI-KELAS-1", Label: "UJI-KELAS-1"}}, nil
}

func (z zonaTiruan) RateZona(_ context.Context, zona string) (*apd.Decimal, error) {
	*z.catat = "rate|" + zona
	if zona != "UJI-ZONA-2" {
		return nil, nil
	}
	d, _, err := apd.NewFromString("2.5")
	return d, err
}

// TestSaranZona - tiket 52: bentuk {baris:[{id,label}]}; zona menurut kode coverage (SetPropertyCoverageDT_FacIn);
// parameter dipangkas; 400 / 503.
func TestSaranZona(t *testing.T) {
	var catat string
	svc := services.Baru(nil).DenganZona(zonaTiruan{&catat})
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/zona?coverage=100829&oldId=%204.2%20PRGBI%20&q=%20uji%20", "", "")
	if kode != 200 || isi != `{"baris":[{"id":"Z2","label":"UJI-ZONA-2"}]}` || catat != "4.2 PRGBI|1,2,3,4,5,6|uji" {
		t.Fatalf("%d %s / %q", kode, isi, catat)
	}
	for kodeCov, mau := range map[string]string{"100835": "X,Y,Z", "100816": "X,Y,Z", "100821": "1,2,3,4,5,6", "100815": "", "": ""} {
		if kode, _ := minta(t, svc, "GET", "/api/nbfacin/zona?coverage="+kodeCov+"&oldId=UJI", "", ""); kode != 200 || catat != "UJI|"+mau+"|" {
			t.Errorf("coverage %q: %d %q", kodeCov, kode, catat)
		}
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/zona?q="+strings.Repeat("A", 256), "", ""); kode != 400 || !strings.Contains(isi, "q paling banyak") {
		t.Errorf("400: %d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/zona", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestSaranKonstruksiGempa - tiket 52: id = label = CONSTRUCTION; 400 / 503.
func TestSaranKonstruksiGempa(t *testing.T) {
	var catat string
	svc := services.Baru(nil).DenganZona(zonaTiruan{&catat})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/konstruksi-gempa?q=kelas", "", ""); kode != 200 ||
		isi != `{"baris":[{"id":"UJI-KELAS-1","label":"UJI-KELAS-1"}]}` || catat != "konstruksi|kelas" {
		t.Errorf("%d %s / %q", kode, isi, catat)
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/konstruksi-gempa?q="+strings.Repeat("A", 256), "", ""); kode != 400 {
		t.Errorf("400: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/konstruksi-gempa", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestHitungCoverageCallActRate - tiket 52 (lanjutan 06-10-2026, ikut Pega): callAct "rate" menjalankan CountPremiRate_ACT
// menurut OLDID harfiah - 4.2 PRGBI = ZONEOJK menurut Zone (tanpa baris / Zone kosong -> Rate KOSONG, premi 0); FLEXAS = max
// RATELOWER okupasi (kosong / NULL dilewati; tanpa nilai -> kosong); "4.1A CC" / "OTHERS" = 0.1; "4.3" = 0.5 / 0.45 menurut
// provinsi harfiah; OLDID lain / tanpa callAct -> Rate tetap. 400 callAct lain, province / occupationIds tanpa callAct,
// batas; 503 hanya bila pembaca memang perlu.
func TestHitungCoverageCallActRate(t *testing.T) {
	var catat string
	svc := services.Baru(nil).DenganKasus(kasusTiruan{k: kasusPeriode()}).DenganZona(zonaTiruan{&catat}).DenganOccupation(occupationTiruan{})
	jalur := "/api/nbfacin/kasus/UJI-NB-1/hitung-coverage"
	badan := func(oldID, zona, tambahan string) string {
		return `{"coverage":{"coverage":"100829","oldId":"` + oldID + `","coverageBasis":"1","rate":"1","indemnityPercentage":"100","zone":"` +
			zona + `","constructionEarthquake":"UJI-KELAS-1"},"tsi":"1000000000","mode":"percent"` + tambahan + "}"
	}
	const rate = `,"callAct":"rate"`
	hasil := func(kode int, isi string) map[string]any {
		t.Helper()
		var j map[string]any
		if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 {
			t.Fatalf("%d %s", kode, isi)
		}
		return j
	}
	for _, k := range []struct {
		nama, badan, rate, premi, catat string
	}{
		{"4.2 zona ada", badan("4.2 PRGBI", "UJI-ZONA-2", rate), "2.5", "2500000", "rate|UJI-ZONA-2"},
		{"4.2 zona tanpa baris", badan("4.2 PRGBI", "UJI-ZONA-LAIN", rate), "", "0", "rate|UJI-ZONA-LAIN"},
		{"4.2 zona kosong", badan("4.2 PRGBI", "", rate), "", "0", ""},
		{"4.2 zona spasi dicari apa adanya", badan("4.2 PRGBI", " ", rate), "", "0", "rate| "},
		{"FLEXAS max", badan("FLEXAS", "", rate+`,"occupationIds":["UJI-OKP-1","UJI-OKP-2","UJI-OKP-NULL","UJI-OKP-TAKADA","","UJI-OKP-1"]`), "0.75", "750000", ""},
		{"FLEXAS tanpa nilai", badan("FLEXAS", "", rate+`,"occupationIds":["UJI-OKP-NULL","UJI-OKP-TAKADA"]`), "", "0", ""},
		{"FLEXAS tanpa okupasi", badan("FLEXAS", "", rate), "", "0", ""},
		{"4.1A CC", badan("4.1A CC", "", rate), "0.1", "100000", ""},
		{"OTHERS", badan("OTHERS", "", rate), "0.1", "100000", ""},
		{"4.3 JAWA BARAT", badan("4.3", "", rate+`,"province":"JAWA BARAT"`), "0.5", "500000", ""},
		{"4.3 peka huruf", badan("4.3", "", rate+`,"province":"Jawa Barat"`), "0.45", "450000", ""},
		{"4.3 tanpa provinsi", badan("4.3", "", rate), "0.45", "450000", ""},
		{"OLDID lain", badan("4.2 prgbi", "UJI-ZONA-2", rate), "1", "1000000", ""},
		{"tanpa callAct", badan("4.2 PRGBI", "UJI-ZONA-2", ""), "1", "1000000", ""},
	} {
		catat = ""
		j := hasil(minta(t, svc, "POST", jalur, k.badan, ""))
		if j["rate"] != k.rate || j["premium"] != k.premi || catat != k.catat {
			t.Errorf("%s: rate %v premi %v catat %q", k.nama, j["rate"], j["premium"], catat)
		}
	}
	for nama, b := range map[string]string{
		"callAct lain":                 badan("4.2 PRGBI", "UJI-ZONA-2", `,"callAct":"premi"`),
		"province tanpa callAct":       badan("4.3", "", `,"province":"BANTEN"`),
		"okupasi tanpa callAct":        badan("FLEXAS", "", `,"occupationIds":["UJI-OKP-1"]`),
		"province terlalu panjang":     badan("4.3", "", rate+`,"province":"`+strings.Repeat("A", 256)+`"`),
		"okupasi terlalu banyak":       badan("FLEXAS", "", rate+`,"occupationIds":["`+strings.Repeat(`A","`, 100)+`A"]`),
		"kode okupasi terlalu panjang": badan("FLEXAS", "", rate+`,"occupationIds":["`+strings.Repeat("A", 1001)+`"]`),
	} {
		if kode, isi := minta(t, svc, "POST", jalur, b, ""); kode != 400 {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	dasar := services.Baru(nil).DenganKasus(kasusTiruan{k: kasusPeriode()})
	for nama, k := range map[string]struct {
		badan string
		kode  int
	}{
		"4.2 tanpa pembaca zona":       {badan("4.2 PRGBI", "UJI-ZONA-2", rate), 503},
		"FLEXAS tanpa pembaca okupasi": {badan("FLEXAS", "", rate+`,"occupationIds":["UJI-OKP-1"]`), 503},
		"4.2 zona kosong":              {badan("4.2 PRGBI", "", rate), 200},
		"FLEXAS tanpa okupasi":         {badan("FLEXAS", "", rate), 200},
		"4.3":                          {badan("4.3", "", rate), 200},
	} {
		if kode, isi := minta(t, dasar, "POST", jalur, k.badan, ""); kode != k.kode {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
}

// TestCoverageZonaSimpan - tiket 52: zone / constructionEarthquake pulang-pergi lewat PUT objek; lebar migrasi 203
// (zone 4000 byte sesudah 204, constructionEarthquake 500) -> 400 ber-jalur.
func TestCoverageZonaSimpan(t *testing.T) {
	var d []models.ObjekFire
	svc := services.Baru(nil).DenganObjek(objekTiruan{&d}).DenganTransaksi(tanpaOracle).DenganPilihanItem(pilihanTiruan{}).
		DenganKasus(kasusTiruan{k: kasusPeriode()})
	badan := func(isi string) string {
		return `{"baris":[{"objectType":"UJI","items":[{"currency":"IDR","tsi":"1000000000","coverages":[{"coverage":"100829","coverageBasis":"1",` + isi + `}]}]}]}`
	}
	zona := `"zone":"UJI-ZONA-2","constructionEarthquake":"UJI-KELAS-1"`
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", badan(zona), "UJI-USER"); kode != 200 || !strings.Contains(isi, zona+`}`) {
		t.Fatalf("%d %s", kode, isi)
	}
	if c := d[0].Items[0].Coverages[0]; c.Zone != "UJI-ZONA-2" || c.ConstructionEarthquake != "UJI-KELAS-1" {
		t.Errorf("model: %+v", c)
	}
	for medan, n := range map[string]int{"zone": 4000, "constructionEarthquake": 500} {
		if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", badan(`"`+medan+`":"`+strings.Repeat("A", n)+`"`), "UJI-USER"); kode != 200 {
			t.Errorf("%s %d byte: %d %s", medan, n, kode, isi)
		}
		if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/objek", badan(`"`+medan+`":"`+strings.Repeat("A", n+1)+`"`), "UJI-USER"); kode != 400 ||
			!strings.Contains(isi, "baris[0].items[0].coverages[0]."+medan+" paling banyak") {
			t.Errorf("%s %d byte: %d %s", medan, n+1, kode, isi)
		}
	}
}
