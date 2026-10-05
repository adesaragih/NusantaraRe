package handlers

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// spreadingTiruanH - satu case "S1" (satu lokasi, satu item IDR, satu coverage oldId C1); "K404" = case tidak ada.
type spreadingTiruanH struct {
	tersimpan models.KasusSpreading
	tulis     int
}

func angkaUji(s string) *apd.Decimal {
	v, err := utils.ParseDecimal(s)
	if err != nil {
		panic(err)
	}
	return v
}

func kasusSpreadingUji() models.KasusSpreading {
	return models.KasusSpreading{StartDateTime: "20250101T170000.000 GMT", Lokasi: []models.LokasiSpreading{{ObjectNo: "1",
		Items: []models.ItemSpreading{{ItemType: "BUILDING", Currency: "IDR", Coverages: []models.CoverageSpreading{{ID: "11", OldID: "C1",
			TSI: angkaUji("1200"), TSILiability: angkaUji("1000"), LimitOfLiability: angkaUji("800"), Premium: angkaUji("10"),
			Discount: angkaUji("2")}}}}}}}
}

func (s *spreadingTiruanH) BacaSpreading(_ context.Context, _ *db.Tx, id string) (models.KasusSpreading, error) {
	if id == "K404" {
		return models.KasusSpreading{}, repository.ErrKasusTidakAda
	}
	return s.tersimpan, nil
}

func (s *spreadingTiruanH) TulisSpreading(_ context.Context, _ *db.Tx, _ string, k models.KasusSpreading, ubah bool) error {
	s.tulis++
	if ubah {
		s.tersimpan.PercentShare = k.PercentShare
	}
	s.tersimpan.Lokasi = k.Lokasi
	return nil
}

func (s *spreadingTiruanH) DaftarTreaty(context.Context, string, string) ([]models.TreatySpreading, error) {
	return nil, nil
}

func (s *spreadingTiruanH) KapasitasTreaty(context.Context, *apd.Decimal, string) (models.KapasitasTreaty, bool, error) {
	return models.KapasitasTreaty{TreatyNameQS: "QS UJI", IDTreatyQS: "Q1", TreatyNameSPL: "SPL UJI", IDTreatySPL: "S1",
		MaxLimitQSIDR: angkaUji("300"), MaxLimitSPLIDR: angkaUji("400")}, true, nil
}

func (s *spreadingTiruanH) KursTerbaru(context.Context, string) (*apd.Decimal, error) {
	return nil, nil
}

func (s *spreadingTiruanH) KursPada(context.Context, string, string) (*apd.Decimal, error) {
	return nil, nil
}

func (s *spreadingTiruanH) CatatanJenisTreaty(_ context.Context, id string) (string, error) {
	if id == "Q1" {
		return "UJI QS", nil
	}
	return "", nil
}

func (s *spreadingTiruanH) IDMataUang(_ context.Context, nama string) (string, error) {
	if nama == "IDR" {
		return "10026", nil
	}
	return "", nil
}

// TestSpreadingKasus - tiket 48: GET (larik selalu [], angka teks), hitung-share / salin tidak menyimpan, PUT menerima
// TampilanSpreading utuh (kunci tambahan diabaikan), 400 / 401 / 404 / 409 / 503.
func TestSpreadingKasus(t *testing.T) {
	src := &spreadingTiruanH{tersimpan: kasusSpreadingUji()}
	svc := services.Baru(nil).DenganSpreading(src).DenganTransaksi(tanpaOracleKlausa)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/S1/spreading", "", "")
	if kode != 200 || !strings.HasPrefix(isi, `{"percentShare":"","treaty":[{"id":"10007","name":"ORS"},{"id":"10015","name":"FACOUT"}],"template":[],`) ||
		!strings.Contains(isi, `"spreading":[]`) || !strings.Contains(isi, `"total":[]`) || !strings.HasSuffix(isi, `"ringkasanMataUang":[],"pesan":[]}`) {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	kode, isi = minta(t, svc, "POST", "/api/nbfacin/kasus/S1/spreading/hitung-share", `{"percentShare":"50"}`, "")
	if kode != 200 || !strings.Contains(isi, `"tsiNusantaraRe":"500","premiNusantaraRe":"4"`) ||
		!strings.Contains(isi, `{"treatyType":"Q1","treatyName":"QS UJI","sharePercentage":"60","tsiGrossSpreaded":"360","tsiSpreaded":"300",`+
			`"claimEstimation":"240","premiumSpreaded":"2.4"}`) || !strings.Contains(isi, `"claimAmountIdr":"0"`) || src.tulis != 0 {
		t.Fatalf("hitung-share: %d %s (tulis %d)", kode, isi, src.tulis)
	}
	kode, isi = minta(t, svc, "POST", "/api/nbfacin/kasus/S1/spreading/salin",
		`{"percentShare":"50","template":[{"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"100"}]}`, "")
	if kode != 200 || !strings.Contains(isi, `"template":[{"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"100"}]`) ||
		!strings.Contains(isi, `"tsiSpreaded":"500"`) || src.tulis != 0 {
		t.Fatalf("salin: %d %s", kode, isi)
	}
	// PUT: badan = TampilanSpreading utuh dari layar, dengan medan tampil dan kunci yang tidak dikenal.
	badan := `{"percentShare":"50","treaty":[{"id":"10015","name":"FACOUT"}],"template":[],"pesan":["lama"],"ringkasanTreaty":[],` +
		`"ringkasanMataUang":[],"lokasi":[{"objectNo":"1","objectName":"","location":"","total":[{"currency":"IDR"}],` +
		`"items":[{"itemType":"BUILDING","currency":"IDR","tsi":"1","totalGrossPremi":"1","totalPremiumRnm":"9","lain":true,` +
		`"coverages":[{"oldId":"C1","coverageBasis":"1","rate":"0.1","premium":"10","discount":"2","tsi":"1200","tsiLiability":"1000",` +
		`"tsiNusantaraRe":"999","premiNusantaraRe":"999","spreading":[{"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"40",` +
		`"tsiGrossSpreaded":"1","tsiSpreaded":"1","claimEstimation":"1","premiumSpreaded":"1","isOldData":"x"}]}]}]}]}`
	kode, isi = minta(t, svc, "PUT", "/api/nbfacin/kasus/S1/spreading", badan, "UJI-USER")
	if kode != 200 || src.tulis != 1 || !strings.HasPrefix(isi, `{"percentShare":"50"`) ||
		!strings.Contains(isi, `"tsiNusantaraRe":"500"`) || !strings.Contains(isi, `"sharePercentage":"40","tsiGrossSpreaded":"240","tsiSpreaded":"200"`) {
		t.Fatalf("PUT: %d %s", kode, isi)
	}
	for nama, c := range map[string]struct {
		jalur, badan, pelaku, mau string
		kode                      int
	}{
		"tanpa percentShare":  {"/api/nbfacin/kasus/S1/spreading", `{"lokasi":[]}`, "UJI-USER", "percentShare", 400},
		"percentShare huruf":  {"/api/nbfacin/kasus/S1/spreading", `{"percentShare":"x","lokasi":[]}`, "UJI-USER", "percentShare", 400},
		"percentShare > 100":  {"/api/nbfacin/kasus/S1/spreading", `{"percentShare":"101","lokasi":[]}`, "UJI-USER", "percentShare", 400},
		"tanpa lokasi":        {"/api/nbfacin/kasus/S1/spreading", `{"percentShare":"50"}`, "UJI-USER", "lokasi", 400},
		"bentuk berubah":      {"/api/nbfacin/kasus/S1/spreading", `{"percentShare":"50","lokasi":[]}`, "UJI-USER", "muat ulang", 409},
		"share > 100":         {"/api/nbfacin/kasus/S1/spreading", strings.Replace(badan, `"sharePercentage":"40"`, `"sharePercentage":"140"`, 1), "UJI-USER", "paling besar 100", 400},
		"tanpa identitas":     {"/api/nbfacin/kasus/S1/spreading", badan, "", "", 401},
		"case tidak ada":      {"/api/nbfacin/kasus/K404/spreading", badan, "UJI-USER", "", 404},
		"salin tanpa templat": {"/api/nbfacin/kasus/S1/spreading/salin", `{"percentShare":"50","template":[]}`, "", "template", 400},
	} {
		metode := "PUT"
		if strings.HasSuffix(c.jalur, "/salin") {
			metode = "POST"
		}
		if kode, isi := minta(t, svc, metode, c.jalur, c.badan, c.pelaku); kode != c.kode || !strings.Contains(isi, c.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/K404/spreading", "", ""); kode != 404 {
		t.Errorf("GET case tidak ada: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/S1/spreading", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}
