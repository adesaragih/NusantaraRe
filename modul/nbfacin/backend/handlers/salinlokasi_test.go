package handlers

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

// TestSalinLokasi - tiket 73 P5 K73-40: kunci badan {percentShare, idxLokasi, template:[{treatyType, treatyName, sharePercentage,
// flagSpreading}]} dan jawaban (Tampilan salin, template []) untuk FIRE dan Aneka; 400 / 401.
func TestSalinLokasi(t *testing.T) {
	if k := kunciTagJSON(templateLokasiKabel{}); !reflect.DeepEqual(k, []string{"treatyType", "treatyName", "sharePercentage", "flagSpreading"}) {
		t.Errorf("template: %v", k)
	}
	if k := kunciTagJSON(badanSalinLokasi{}); !reflect.DeepEqual(k, []string{"percentShare", "idxLokasi", "template"}) {
		t.Errorf("badan: %v", k)
	}
	badan := `{"percentShare":"50","idxLokasi":1,"template":[{"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"100","flagSpreading":"SF-MR 2023"}]}`
	src := &spreadingTiruanH{tersimpan: kasusSpreadingUji()}
	svc := services.Baru(nil).DenganSpreading(src).DenganTransaksi(tanpaOracleKlausa)
	kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/S1/spreading/salin-lokasi", badan, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"template":[]`) || !strings.Contains(isi, `"treatyType":"10015","treatyName":"FACOUT","sharePercentage":"100"`) ||
		src.tulis != 0 {
		t.Fatalf("FIRE: %d %s", kode, isi)
	}
	lokasi := lokasiSpreadingAnekaUji()
	svcA := services.Baru(nil).DenganObjekAneka(anekaTiruanH{d: &lokasi}).DenganSpreading(&spreadingTiruanH{tersimpan: models.KasusSpreading{
		PercentShare: angkaUji("40"), StartDateTime: "20250101T170000.000 GMT"}})
	if kode, isi := minta(t, svcA, "POST", "/api/nbfacin/kasus/UJI-NB-1/spreading-aneka/salin-lokasi", badan, "UJI-USER"); kode != 200 ||
		!strings.Contains(isi, `"template":[]`) || !strings.Contains(isi, `"tampilTotalLokasi":false`) {
		t.Errorf("Aneka: %d %s", kode, isi)
	}
	for nama, c := range map[string]struct {
		badan, pelaku, mau string
		kode               int
	}{
		"tanpa idxLokasi": {`{"percentShare":"50","template":[]}`, "UJI-USER", "idxLokasi", 400},
		"idx di luar":     {strings.Replace(badan, `"idxLokasi":1`, `"idxLokasi":2`, 1), "UJI-USER", "idxLokasi 2", 400},
		"flag wajib":      {strings.Replace(badan, `"flagSpreading":"SF-MR 2023"`, `"flagSpreading":""`, 1), "UJI-USER", "flagSpreading wajib", 400},
		"share huruf":     {strings.Replace(badan, `"sharePercentage":"100"`, `"sharePercentage":"x"`, 1), "UJI-USER", "template[0].sharePercentage", 400},
		"tanpa identitas": {badan, "", "", 401},
	} {
		if kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/S1/spreading/salin-lokasi", c.badan, c.pelaku); kode != c.kode || !strings.Contains(isi, c.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
}
