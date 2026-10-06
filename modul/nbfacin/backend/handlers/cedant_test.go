package handlers

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
	"nusantarare/modul/nbfacin/backend/services"
)

// cedantTiruanH - satu case "S1"; "K404" = case tidak ada.
type cedantTiruanH struct{ ada models.KasusCedant }

func (c *cedantTiruanH) BacaCedant(_ context.Context, id string) (models.KasusCedant, error) {
	if id == "K404" {
		return models.KasusCedant{}, repository.ErrKasusTidakAda
	}
	return c.ada, nil
}

func (c *cedantTiruanH) TulisCedant(_ context.Context, _ *db.Tx, id string, k models.SimpanCedant) error {
	if id == "K404" {
		return repository.ErrKasusTidakAda
	}
	c.ada.ShareCedantType, c.ada.Cedant = k.ShareCedantType, k.Cedant
	return nil
}

// TestCedantKasus - tiket 49: GET bentuk TampilanCedant persis (larik selalu [], angka teks), PUT baca ulang, isian
// (kunci asing, share, tipe) 400, 401, 404, 503.
func TestCedantKasus(t *testing.T) {
	c := &cedantTiruanH{ada: models.KasusCedant{SobName: "UJI SOB", Cedant: []models.BarisCedant{},
		CedingUmum: []models.BarisCedant{{CedingCo: "A1", CedingCoName: "UJI A"}}}}
	sp := kasusSpreadingUji()
	sp.PercentShare = angkaUji("50")
	svc := services.Baru(nil).DenganCedant(c).DenganSpreading(&spreadingTiruanH{tersimpan: sp}).DenganTransaksi(tanpaOracleKlausa)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/S1/cedant", "", "")
	if kode != 200 || isi != `{"sobName":"UJI SOB","shareCedantType":"","percentShare":"50","totalTsiRnm":"500","totalPremiRnm":"4",`+
		`"wajib":false,"cedant":[],"cedingUmum":[{"cedingCo":"A1","cedingCoName":"UJI A","shareCeding":""}]}` {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	badan := `{"shareCedantType":"0","cedant":[{"cedingCo":"A1","cedingCoName":"UJI A","shareCeding":"60.5"}]}`
	kode, isi = minta(t, svc, "PUT", "/api/nbfacin/kasus/S1/cedant", badan, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `"shareCedantType":"0"`) ||
		!strings.Contains(isi, `"cedant":[{"cedingCo":"A1","cedingCoName":"UJI A","shareCeding":"60.5"}]`) {
		t.Fatalf("PUT: %d %s", kode, isi)
	}
	for nama, x := range map[string]struct {
		jalur, badan, pelaku, mau string
		kode                      int
	}{
		"kunci asing":     {"/api/nbfacin/kasus/S1/cedant", `{"shareCedantType":"0","cedant":[],"lain":1}`, "UJI-USER", "shareCedantType", 400},
		"tanpa cedant":    {"/api/nbfacin/kasus/S1/cedant", `{"shareCedantType":"0"}`, "UJI-USER", "cedant", 400},
		"tanpa tipe":      {"/api/nbfacin/kasus/S1/cedant", `{"cedant":[]}`, "UJI-USER", "shareCedantType", 400},
		"tipe 2":          {"/api/nbfacin/kasus/S1/cedant", `{"shareCedantType":"2","cedant":[]}`, "UJI-USER", "shareCedantType", 400},
		"share huruf":     {"/api/nbfacin/kasus/S1/cedant", strings.Replace(badan, `"60.5"`, `"x"`, 1), "UJI-USER", "cedant[0].shareCeding", 400},
		"share 101":       {"/api/nbfacin/kasus/S1/cedant", strings.Replace(badan, `"60.5"`, `"101"`, 1), "UJI-USER", "Value can't be more than 100 or less than 0!", 400},
		"tanpa identitas": {"/api/nbfacin/kasus/S1/cedant", badan, "", "", 401},
		"case tidak ada":  {"/api/nbfacin/kasus/K404/cedant", badan, "UJI-USER", "", 404},
	} {
		if kode, isi := minta(t, svc, "PUT", x.jalur, x.badan, x.pelaku); kode != x.kode || !strings.Contains(isi, x.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, svc, "GET", "/api/nbfacin/kasus/K404/cedant", "", ""); kode != 404 {
		t.Errorf("GET case tidak ada: %d", kode)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/S1/cedant", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}
