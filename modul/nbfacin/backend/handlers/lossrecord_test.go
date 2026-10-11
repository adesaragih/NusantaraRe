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

type lossRecordTiruanH struct{ tulis int }

func (l *lossRecordTiruanH) BacaLossRecord(context.Context, string) (models.LossRecord, error) {
	return models.LossRecord{Baris: []models.BarisLossRecord{}}, nil
}

func (l *lossRecordTiruanH) TulisLossRecord(context.Context, *db.Tx, string, models.LossRecord) error {
	l.tulis++
	return nil
}

type jabatanTiruanH map[string]string

func (j jabatanTiruanH) JabatanAkun(_ context.Context, akun string) (string, error) {
	return j[akun], nil
}

type deductionLossRecordH struct{}

func (deductionLossRecordH) BacaDeduction(context.Context, string) (models.KasusDeduction, error) {
	return models.KasusDeduction{RICommision: angkaUji("10")}, nil
}
func (deductionLossRecordH) TulisDeduction(context.Context, *db.Tx, string, models.KasusDeduction) error {
	return nil
}
func (deductionLossRecordH) TarifPajak(context.Context) (models.TarifPajak, bool, error) {
	return models.TarifPajak{}, false, nil
}

func layananLossRecordH(lt *lossRecordTiruanH) *services.Service {
	return services.Baru(nil).DenganKasus(kasusTiruan{k: kasusPeriode()}).DenganSpreading(&spreadingTiruanH{tersimpan: kasusSpreadingUji()}).
		DenganDeduction(deductionLossRecordH{}).DenganLossRecord(lt, jabatanTiruanH{"UJI-IT": "IT", "UJI-UW": "Underwriter"}).
		DenganKelasBisnis(kelasBisnisTiruan{}).DenganTransaksi(tanpaOracle)
}

// TestAksesLossRecordKasus - GET kasus: aksesLossRecord dan aksesKorespondensi (tiket 56, aturan sama) true hanya untuk
// akun ber-JOB_POSITION IT; tanpa sesi false.
func TestAksesLossRecordKasus(t *testing.T) {
	svc := layananLossRecordH(&lossRecordTiruanH{})
	for pelaku, mau := range map[string]string{"UJI-IT": "true", "UJI-UW": "false", "": "false"} {
		kode, isi := minta(t, svc, "GET", "/api/nbfacin/kasus/UJI-NB-1", "", pelaku)
		if kode != 200 || !strings.Contains(isi, `"aksesLossRecord":`+mau) || !strings.Contains(isi, `"aksesKorespondensi":`+mau) {
			t.Errorf("%q: %d %s", pelaku, kode, isi)
		}
	}
}

// TestLossRecordHandler - bentuk persis kontrak; hitung / simpan; badan ketat; 401 / 403 / 400.
func TestLossRecordHandler(t *testing.T) {
	lt := &lossRecordTiruanH{}
	svc := layananLossRecordH(lt)
	jalur := "/api/nbfacin/kasus/UJI-NB-1/loss-record"
	kode, isi := minta(t, svc, "GET", jalur, "", "UJI-IT")
	var j map[string]json.RawMessage
	if err := json.Unmarshal([]byte(isi), &j); err != nil || kode != 200 {
		t.Fatalf("%d %s", kode, isi)
	}
	var kunci []string
	for k := range j {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	if strings.Join(kunci, ",") != "baris,pesan,premiMb,premiOther,premiPar,premiPl,rasio,totalCob" ||
		!strings.Contains(isi, `"rasio":{"gross1":"","gross3":"","gross5":"","net1":"","net3":"","net5":""}`) {
		t.Errorf("bentuk: %s", isi)
	}
	badan := `{"premiPar":"1000","premiMb":"","premiPl":"","premiOther":"","baris":[{"classOfBusiness":"UJI COB","dateOfLoss":"15-06-2024","causeOfLoss":"uji","currency":"IDR","lossAmount":"10"}]}`
	kode, isi = minta(t, svc, "POST", jalur+"/hitung", badan, "UJI-IT")
	if kode != 200 || !strings.Contains(isi, `"baris":[{"uwYear":"2024","classOfBusiness":"UJI COB","dateOfLoss":"15-06-2024","causeOfLoss":"UJI","currency":"IDR","lossAmount":"10","lossAmountIdr":"10"}]`) ||
		!strings.Contains(isi, `"totalCob":[{"classOfBusiness":"UJI COB","currency":"IDR","lossAmount":"10","lossAmountIdr":"10"}]`) ||
		!strings.Contains(isi, `"gross1":"1"`) || lt.tulis != 0 {
		t.Fatalf("hitung: %d %s", kode, isi)
	}
	if kode, isi := minta(t, svc, "PUT", jalur, badan, "UJI-IT"); kode != 200 || lt.tulis != 1 {
		t.Fatalf("simpan: %d %s", kode, isi)
	}
	for nama, k := range map[string]struct {
		metode, badan, pelaku string
		kode                  int
	}{
		"tanpa identitas":      {"GET", "", "", 401},
		"bukan IT":             {"POST", badan, "UJI-UW", 403},
		"bukan IT simpan":      {"PUT", badan, "UJI-UW", 403},
		"kunci asing":          {"PUT", `{"premiPar":"","baris":[],"rasio":{}}`, "UJI-IT", 400},
		"angka":                {"PUT", `{"premiPar":"abc","baris":[]}`, "UJI-IT", 400},
		"tanggal":              {"PUT", `{"baris":[{"dateOfLoss":"2025-06-15"}]}`, "UJI-IT", 400},
		"bukan IT badan rusak": {"PUT", `{"rasio":{}}`, "UJI-UW", 403},
		"COB tidak di daftar":  {"PUT", `{"baris":[{"classOfBusiness":"UJI TAK ADA","dateOfLoss":"","causeOfLoss":"","currency":"IDR","lossAmount":"1"}]}`, "UJI-IT", 400},
	} {
		j := jalur
		if k.metode == "POST" {
			j += "/hitung"
		}
		if kode, isi := minta(t, svc, k.metode, j, k.badan, k.pelaku); kode != k.kode || lt.tulis != 1 {
			t.Errorf("%s: %d %s", nama, kode, isi)
		}
	}
}

// TestSaranKelasBisnisLossRecord - GET /api/nbfacin/class-of-business/saran?q= -> {baris:[{id,label}]}; tidak bertabrakan
// dengan GET /class-of-business (tiket 28).
func TestSaranKelasBisnisLossRecord(t *testing.T) {
	svc := layananLossRecordH(&lossRecordTiruanH{})
	if kode, isi := minta(t, svc, "GET", "/api/nbfacin/class-of-business/saran?q=api", "", ""); kode != 200 ||
		isi != `{"baris":[{"id":"UJI-B1","label":"UJI api"}]}` {
		t.Errorf("%d %s", kode, isi)
	}
	if kode, _ := minta(t, services.Baru(nil), "GET", "/api/nbfacin/class-of-business/saran", "", ""); kode != 503 {
		t.Errorf("503: %d", kode)
	}
}
