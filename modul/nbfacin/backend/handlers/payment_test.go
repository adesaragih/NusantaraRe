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

// paymentTiruanH - satu case "S1"; "K404" = case tidak ada.
type paymentTiruanH struct{ ada models.KasusPayment }

func (p *paymentTiruanH) BacaPayment(_ context.Context, id string) (models.KasusPayment, error) {
	if id == "K404" {
		return models.KasusPayment{}, repository.ErrKasusTidakAda
	}
	return p.ada, nil
}

func (p *paymentTiruanH) TulisPayment(_ context.Context, _ *db.Tx, _ string, k models.SimpanPayment) error {
	p.ada = models.KasusPayment{Installment: "1", Cicilan: map[string][]models.CicilanPayment{}}
	for _, m := range k.MataUang {
		p.ada.Cicilan[m.Baris.Currency] = m.Cicilan
	}
	return nil
}

// TestPaymentKasus - tiket 51: GET bentuk TampilanPayment persis (larik selalu [], Due Date DD-MM-YYYY, persen apa adanya),
// POST hitung tanpa simpan, PUT baca ulang, kunci asing / aksi / tanggal 400, 401, 404, 409, 503.
func TestPaymentKasus(t *testing.T) {
	sp := kasusSpreadingUji() // satu item IDR, premi NR 4; Begin date 20250101T170000.000 GMT
	sp.PercentShare = angkaUji("50")
	pt := &paymentTiruanH{ada: models.KasusPayment{Cicilan: map[string][]models.CicilanPayment{}}}
	svc := services.Baru(nil).DenganPayment(pt).DenganDeduction(&deductionTiruanH{}).
		DenganSpreading(&spreadingTiruanH{tersimpan: sp}).DenganTransaksi(tanpaOracleKlausa)
	kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/S1/payment", "", "")
	if kode != 200 || isi != `{"installment":"","tarifPph":"2","tarifPpn":"11","mataUang":[{"currency":"IDR","premium":"4","sumTotalPayment":"4","cicilan":[]}],"pesan":[]}` {
		t.Fatalf("GET: %d %s", kode, isi)
	}
	kode, isi = minta(t, svc, "POST", "/api/nbfacin/kasus/S1/payment/hitung",
		`{"aksi":"installment","currency":"","idx":0,"installment":"3","mataUang":[{"currency":"IDR","cicilan":[]}]}`, "")
	if kode != 200 || !strings.Contains(isi, `"installmentNo":"1","dueDate":"02-04-2025","installmentPercentage":"33.33333333333333333333"`) ||
		!strings.Contains(isi, `"installmentPercentage":"33.33333333333333333334"`) {
		t.Fatalf("hitung installment: %d %s", kode, isi)
	}
	badan := `{"installment":"1","mataUang":[{"currency":"IDR","cicilan":[{"installmentNo":"1","dueDate":"02-04-2025","installmentPercentage":"100"}]}]}`
	kode, isi = minta(t, svc, "PUT", "/api/nbfacin/kasus/S1/payment", badan, "UJI-USER")
	if kode != 200 || !strings.Contains(isi, `{"installmentNo":"1","dueDate":"02-04-2025","installmentPercentage":"100","premium":"4"`) {
		t.Fatalf("PUT: %d %s", kode, isi)
	}
	for nama, x := range map[string]struct {
		metode, jalur, badan, pelaku, mau string
		kode                              int
	}{
		"kunci asing":     {"PUT", "/api/nbfacin/kasus/S1/payment", `{"installment":"1","mataUang":[],"lain":1}`, "UJI-USER", "badan permintaan", 400},
		"PUT ber-aksi":    {"PUT", "/api/nbfacin/kasus/S1/payment", `{"aksi":"persen","currency":"IDR","idx":0,` + badan[1:], "UJI-USER", "tidak memuat aksi", 400},
		"hitung tanpa":    {"POST", "/api/nbfacin/kasus/S1/payment/hitung", badan, "", "wajib memuat aksi", 400},
		"tanggal":         {"PUT", "/api/nbfacin/kasus/S1/payment", strings.Replace(badan, "02-04-2025", "2025-04-02", 1), "UJI-USER", "dueDate harus DD-MM-YYYY", 400},
		"persen":          {"PUT", "/api/nbfacin/kasus/S1/payment", strings.Replace(badan, `"100"`, `"1000"`, 1), "UJI-USER", "installmentPercentage", 400},
		"proteksi":        {"PUT", "/api/nbfacin/kasus/S1/payment", strings.Replace(badan, `"cicilan":[{"installmentNo":"1","dueDate":"02-04-2025","installmentPercentage":"100"}]`, `"cicilan":[]`, 1), "UJI-USER", "List Installment can't be null!", 400},
		"mata uang lain":  {"PUT", "/api/nbfacin/kasus/S1/payment", strings.Replace(badan, `"IDR"`, `"USD"`, 1), "UJI-USER", "muat ulang", 409},
		"tanpa identitas": {"PUT", "/api/nbfacin/kasus/S1/payment", badan, "", "", 401},
		"case tidak ada":  {"PUT", "/api/nbfacin/kasus/K404/payment", badan, "UJI-USER", "", 404},
	} {
		if kode, isi := minta(t, svc, x.metode, x.jalur, x.badan, x.pelaku); kode != x.kode || !strings.Contains(isi, x.mau) {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/kasus/S1/payment", "", ""); kode != 503 {
		t.Errorf("tanpa basis data: %d", kode)
	}
}
