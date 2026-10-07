package handlers

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/services"
)

type scoringTiruanH struct {
	k     models.KasusScoring
	tulis int
}

func (s *scoringTiruanH) BacaScoring(context.Context, string) (models.KasusScoring, error) {
	return s.k, nil
}

func (s *scoringTiruanH) TulisScoring(context.Context, *db.Tx, string, models.Scoring) error {
	s.tulis++
	return nil
}

func layananScoringH(st *scoringTiruanH) *services.Service {
	d := []models.ObjekFire{{ObjectNo: "1", IsTopRisk: true, RiskLocation: "UJI ALAMAT",
		Occupations: []models.OkupasiObjek{{OccupationID: "UJI-OKP-4", OccupationName: "UJI PABRIK"}},
		Items:       []models.ItemObjek{{Currency: "IDR", TSI: angkaUji("1000")}}},
		{ObjectNo: "2", Occupations: []models.OkupasiObjek{{OccupationID: "UJI-OKP-1", OccupationName: "UJI KANTOR"}}}}
	return services.Baru(nil).DenganScoring(st).DenganKasus(kasusTiruan{k: kasusPeriode()}).DenganObjek(objekTiruan{&d}).
		DenganSpreading(&spreadingTiruanH{tersimpan: kasusSpreadingUji()}).DenganOccupation(occupationTiruan{}).
		DenganTransaksi(tanpaOracle)
}

// kunciScoring - kunci tingkat atas TampilanScoring, kunci kolom, dan 11 kunci hasil (kontrak tiket 54).
var (
	kunciScoring      = []string{"bolehUbah", "finalScore", "hasil", "insuredName", "kolom", "lokasiDominan", "mulai", "noteFinalScore", "pesan", "selesai", "status"}
	kunciKolomScoring = []string{"indexLocation", "occupationNote", "persenLossRatio", "pilihan", "skor", "topRiskLocation", "total"}
	kunciHasilScoring = []string{"acceptableNotApproval", "acceptableWithApproval", "dueLossRatio", "dueObjectConditions", "dueOccupation",
		"dueOthers", "notAcceptable", "opDirector", "opDivHead", "presDirector", "techDirector"}
)

func kunciUrut(m map[string]json.RawMessage) []string {
	var k []string
	for x := range m {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}

// TestBacaScoringHandler - GET: bentuk persis kontrak, dua kolom selalu, 24 kunci pilihan / skor, daftar dominan.
func TestBacaScoringHandler(t *testing.T) {
	st := &scoringTiruanH{k: models.KasusScoring{Position: "Offer", FlagOnGoingPolicy: "0"}}
	kode, isi := minta(t, layananScoringH(st), "GET", "/api/nbfacin/kasus/UJI-NB-1/scoring", "", "")
	var j map[string]json.RawMessage
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 || strings.Join(kunciUrut(j), ",") != strings.Join(kunciScoring, ",") {
		t.Fatalf("%d %s", kode, isi)
	}
	var kolom []map[string]json.RawMessage
	var hasil map[string]json.RawMessage
	if json.Unmarshal(j["kolom"], &kolom) != nil || len(kolom) != 2 || json.Unmarshal(j["hasil"], &hasil) != nil ||
		strings.Join(kunciUrut(kolom[0]), ",") != strings.Join(kunciKolomScoring, ",") ||
		strings.Join(kunciUrut(hasil), ",") != strings.Join(kunciHasilScoring, ",") {
		t.Fatalf("kolom / hasil: %s", isi)
	}
	var pilihan map[string]int
	var skor map[string]string
	if json.Unmarshal(kolom[0]["pilihan"], &pilihan) != nil || json.Unmarshal(kolom[0]["skor"], &skor) != nil || len(pilihan) != 24 ||
		len(skor) != 24 || pilihan["occupation"] != 4 || pilihan["totalSumInsured"] != 1 || skor["occupation"] != "-5" ||
		string(kolom[0]["total"]) != `"-6"` || string(j["finalScore"]) != `"-6"` || string(j["noteFinalScore"]) != `"Poor"` ||
		string(j["bolehUbah"]) != "true" || string(j["status"]) != `"New Business"` ||
		string(j["lokasiDominan"]) != `[{"indexLocation":"1","label":"UJI PABRIK"},{"indexLocation":"2","label":"UJI KANTOR"}]` {
		t.Errorf("isi: %s", isi)
	}
	var skor2 map[string]string
	if json.Unmarshal(kolom[1]["skor"], &skor2) != nil || skor2["occupation"] != "" || string(kolom[1]["total"]) != `""` {
		t.Errorf("kolom dominan kosong: %s", kolom[1]["skor"])
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/UJI-NB-1/scoring", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}

// TestHitungSimpanScoringHandler - POST hitung / PUT: badan ketat, persen ber-jalur, 400 / 401 / 409, simpan sekali.
func TestHitungSimpanScoringHandler(t *testing.T) {
	st := &scoringTiruanH{k: models.KasusScoring{Position: "Offer", FlagOnGoingPolicy: "0"}}
	svc := layananScoringH(st)
	badan := `{"kolom":[{"indexLocation":"","pilihan":{"riot":1},"persenLossRatio":"20"},{"indexLocation":"2","pilihan":{},"persenLossRatio":""}]}`
	kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/scoring/hitung", badan, "")
	if kode != 200 || !strings.Contains(isi, `"lossRatio":5`) || !strings.Contains(isi, `"occupationNote":"UJI KANTOR"`) || st.tulis != 0 {
		t.Fatalf("hitung: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/scoring", badan, "UJI-USER"); kode != 200 || st.tulis != 1 {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	for nama, b := range map[string]string{
		"kunci asing": `{"kolom":[{"indexLocation":"","pilihan":{},"persenLossRatio":"","skor":{}}]}`,
		"tanpa kolom": `{"kolom":[]}`,
		"tiga kolom":  `{"kolom":[{},{},{}]}`,
		"persen":      `{"kolom":[{"persenLossRatio":"abc"}]}`,
		"pilihan":     `{"kolom":[{"pilihan":{"riot":9}}]}`,
		"sama":        `{"kolom":[{},{"indexLocation":"1"}]}`,
	} {
		if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/scoring", b, "UJI-USER"); kode != 400 || st.tulis != 1 {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, isi := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/scoring", `{"kolom":[{"persenLossRatio":"abc"}]}`, "UJI-USER"); !strings.Contains(isi, "kolom[0].persenLossRatio") {
		t.Errorf("jalur persen: %d %s", kode, isi)
	}
	if kode, _ := minta(t, svc, "PUT", "/api/nbfacin/kasus/UJI-NB-1/scoring", badan, ""); kode != 401 {
		t.Errorf("401: %d", kode)
	}
	st.k.Position = "Akseptasi"
	if kode, _ := minta(t, svc, "POST", "/api/nbfacin/kasus/UJI-NB-1/scoring/hitung", badan, ""); kode != 409 {
		t.Errorf("409: %d", kode)
	}
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1/scoring", "", ""); kode != 200 || !strings.Contains(isi, `"bolehUbah":false`) {
		t.Errorf("GET baca-saja: %d %s", kode, isi)
	}
}

// TestRiskScoringKasus - GET kasus membawa Scoring Risk TERSIMPAN (tiket 54, kolom "Risk Scoring" General): teks desimal
// apa adanya, tidak dihitung ulang.
func TestRiskScoringKasus(t *testing.T) {
	k := kasusPeriode()
	k.RiskScoring = models.RingkasanScoring{FinalScore: angkaUji("98.5"), NoteFinalScore: "Good"}
	kode, isi := minta(t, services.Baru(nil).DenganKasus(kasusTiruan{k: k}), "GET", "/api/nbfacin/kasus/UJI-NB-1", "", "")
	if kode != 200 || !strings.Contains(isi, `"riskScoring":{"finalScore":"98.5","noteFinalScore":"Good"}`) {
		t.Errorf("%d %s", kode, isi)
	}
}
