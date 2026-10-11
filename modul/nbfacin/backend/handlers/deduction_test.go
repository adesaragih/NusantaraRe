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

// deductionTiruanH - satu case "S1"; "K404" = case tidak ada; tarif PPH 2 / PPN 11.
type deductionTiruanH struct{ ada models.KasusDeduction }

func (d *deductionTiruanH) BacaDeduction(_ context.Context, id string) (models.KasusDeduction, error) {
	if id == "K404" {
		return models.KasusDeduction{}, repository.ErrKasusTidakAda
	}
	return d.ada, nil
}

func (d *deductionTiruanH) TulisDeduction(_ context.Context, _ *db.Tx, _ string, k models.KasusDeduction) error {
	d.ada = k
	return nil
}

func (d *deductionTiruanH) TarifPajak(context.Context) (models.TarifPajak, bool, error) {
	return models.TarifPajak{PPH: angkaUji("2"), PPN: angkaUji("11")}, true, nil
}

// TestDeductionKasus - tiket 50: GET bentuk TampilanDeduction persis (larik selalu [], angka teks), POST hitung tanpa
// simpan, PUT baca ulang, kunci asing / aksi / idx 400, 401, 404, 409, 503.
func TestDeductionKasus(t *testing.T) {
	dt := &deductionTiruanH{}
	sp := kasusSpreadingUji() // satu item IDR, premi NR 50 × 10 / 100 − 50 × 2 / 100 = 4
	sp.PercentShare = angkaUji("50")
	svc := services.Baru(nil).DenganDeduction(dt).DenganSpreading(&spreadingTiruanH{tersimpan: sp}).DenganTransaksi(tanpaOracleKlausa)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/S1/deduction", "", "")
	if kode != 200 || isi != `{"pctPremiumDeduction":"","pctBrokerageFee":"","ppnCheck":false,"tarifPph":"2","tarifPpn":"11",`+
		`"mataUang":[{"currency":"IDR","premium":"4","commision":"0","deduction2":"","brokerageFee":"0","pph":"0","ppn":"0"}],"pesan":[]}` {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	badan := `{"pctPremiumDeduction":"10","pctBrokerageFee":"5","ppnCheck":false,"mataUang":[{"currency":"IDR","commision":"","deduction2":"1","brokerageFee":""}]}`
	kode, isi = minta(t, svc, "POST", "/api/nbfacin/kasus/S1/deduction/hitung", `{"aksi":"persen","idx":0,`+badan[1:], "")
	if kode != 200 || !strings.Contains(isi, `"ppnCheck":true`) ||
		!strings.Contains(isi, `{"currency":"IDR","premium":"4","commision":"0.4","deduction2":"1","brokerageFee":"0.2","pph":"0.004","ppn":"0.022"}`) {
		t.Fatalf("hitung: %d %s", kode, isi)
	}
	kode, isi = minta(t, svc, "PUT", "/api/nbfacin/kasus/S1/deduction", badan, "UJI-USER")
	if kode != 200 || !strings.HasPrefix(isi, `{"pctPremiumDeduction":"10","pctBrokerageFee":"5","ppnCheck":true`) ||
		!strings.Contains(isi, `"deduction2":"1"`) {
		t.Fatalf("PUT: %d %s", kode, isi)
	}
	for nama, x := range map[string]struct {
		metode, jalur, badan, pelaku, mau string
		kode                              int
	}{
		"kunci asing":     {"PUT", "/api/nbfacin/kasus/S1/deduction", `{"mataUang":[],"lain":1}`, "UJI-USER", "badan permintaan", 400},
		"PUT ber-aksi":    {"PUT", "/api/nbfacin/kasus/S1/deduction", `{"aksi":"persen","idx":0,` + badan[1:], "UJI-USER", "badan permintaan", 400},
		"hitung tanpa":    {"POST", "/api/nbfacin/kasus/S1/deduction/hitung", badan, "", "aksi", 400},
		"aksi lain":       {"POST", "/api/nbfacin/kasus/S1/deduction/hitung", `{"aksi":"x","idx":0,` + badan[1:], "", "aksi harus", 400},
		"angka huruf":     {"PUT", "/api/nbfacin/kasus/S1/deduction", strings.Replace(badan, `"10"`, `"x"`, 1), "UJI-USER", "pctPremiumDeduction", 400},
		"deduction2 > p":  {"PUT", "/api/nbfacin/kasus/S1/deduction", strings.Replace(badan, `"deduction2":"1"`, `"deduction2":"5"`, 1), "UJI-USER", "Deduction 2 cannot be greater than the Premium", 400},
		"mata uang lain":  {"PUT", "/api/nbfacin/kasus/S1/deduction", strings.Replace(badan, `"IDR"`, `"USD"`, 1), "UJI-USER", "muat ulang", 409},
		"tanpa identitas": {"PUT", "/api/nbfacin/kasus/S1/deduction", badan, "", "", 401},
		"case tidak ada":  {"PUT", "/api/nbfacin/kasus/K404/deduction", badan, "UJI-USER", "", 404},
	} {
		if kode, isi := minta(t, svc, x.metode, x.jalur, x.badan, x.pelaku); kode != x.kode || !strings.Contains(isi, x.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/S1/deduction", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}

// TestDeductionKomaDesimal - `@replaceAll(RICommision, ",", ".")`: koma desimal % Premium Deduction diterima; medan lain tidak.
func TestDeductionKomaDesimal(t *testing.T) {
	sp := kasusSpreadingUji()
	sp.PercentShare = angkaUji("50")
	svc := services.Baru(nil).DenganDeduction(&deductionTiruanH{}).DenganSpreading(&spreadingTiruanH{tersimpan: sp}).DenganTransaksi(tanpaOracleKlausa)
	badan := `{"aksi":"persen","idx":0,"pctPremiumDeduction":"10,5","pctBrokerageFee":"5","ppnCheck":true,"mataUang":[{"currency":"IDR","commision":"","deduction2":"","brokerageFee":""}]}`
	if kode, isi := minta(t, svc, "POST", "/api/nbfacin/kasus/S1/deduction/hitung", badan, ""); kode != 200 ||
		!strings.Contains(isi, `"pctPremiumDeduction":"10.5"`) || !strings.Contains(isi, `"commision":"0.42"`) {
		t.Fatalf("koma: %d %s", kode, isi)
	}
	if kode, _ := minta(t, svc, "POST", "/api/nbfacin/kasus/S1/deduction/hitung", strings.Replace(badan, `"5"`, `"5,5"`, 1), ""); kode != 400 {
		t.Fatalf("koma di Pct Brokerage Fee: %d", kode)
	}
}
