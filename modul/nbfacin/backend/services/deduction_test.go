package services

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// deductionTiruan - PenyimpanDeduction dalam memori; "K404" = case tidak ada; tarif nil = TAX kosong.
type deductionTiruan struct {
	ada   models.KasusDeduction
	tulis []models.KasusDeduction
	tarif *models.TarifPajak
}

func (d *deductionTiruan) BacaDeduction(_ context.Context, id string) (models.KasusDeduction, error) {
	if id == "K404" {
		return models.KasusDeduction{}, repository.ErrKasusTidakAda
	}
	return d.ada, nil
}

func (d *deductionTiruan) TulisDeduction(_ context.Context, _ *db.Tx, _ string, k models.KasusDeduction) error {
	d.tulis = append(d.tulis, k)
	d.ada = k
	return nil
}

func (d *deductionTiruan) TarifPajak(context.Context) (models.TarifPajak, bool, error) {
	if d.tarif == nil {
		return models.TarifPajak{}, false, nil
	}
	return *d.tarif, true, nil
}

var tarifUji = models.TarifPajak{PPH: d("2"), PPN: d("11")}

// layananDeduction - objek kasusDuaItem PS 50: premi NR IDR 6 (4 + 2), USD 1.
func layananDeduction(dt *deductionTiruan) *Service {
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "USD")}}
	return Baru(nil).DenganDeduction(dt).DenganSpreading(src).DenganTransaksi(tanpaTx)
}

func barisUji(t *testing.T, b []models.BarisDeduction) string {
	t.Helper()
	var s []string
	for _, x := range b {
		s = append(s, strings.Join([]string{x.Currency, t8(x.Premium), t8(x.Commision), t8(x.Deduction2), t8(x.BrokerageFee), t8(x.PPh), t8(x.PPN)}, "/"))
	}
	return strings.Join(s, " ")
}

// TestHitungRIComm - SetRIComm_Act 8: Commision = 10 % × premi, Brokerage = 5 % × premi, PPh 2 % / PPN 11 % × Brokerage
// (PPnCheck true); PPnCheck false: dasar = Brokerage / 1,11 (8 desimal, cabang mati di Pega).
func TestHitungRIComm(t *testing.T) {
	b := []models.BarisDeduction{{Currency: "IDR", Premium: d("6")}, {Currency: "USD", Premium: d("1")}}
	hitungRIComm(b, d("10"), d("5"), tarifUji, true)
	if got := barisUji(t, b); got != "IDR/6/0.6//0.3/0.006/0.033 USD/1/0.1//0.05/0.001/0.0055" { // Deduction2 tidak disentuh (kosong)
		t.Fatalf("persen: %s", got)
	}
	hitungRIComm(b, d("10"), d("5"), tarifUji, false) // 0.3 / 1.11 = 0.27027027
	if t8(b[0].PPh) != "0.00540541" || t8(b[0].PPN) != "0.02972973" {
		t.Fatalf("PPnCheck false: %s / %s", t8(b[0].PPh), t8(b[0].PPN))
	}
}

func isianUji(c ...string) IsianDeduction {
	// c: commision, deduction2, brokerage untuk IDR lalu USD
	return IsianDeduction{PctPremiumDeduction: d("10"), PctBrokerageFee: d("5"), MataUang: []models.BarisDeduction{
		{Currency: "IDR", Commision: d(c[0]), Deduction2: d(c[1]), BrokerageFee: d(c[2])},
		{Currency: "USD", Commision: d(c[3]), Deduction2: d(c[4]), BrokerageFee: d(c[5])}}}
}

func TestHitungDeductionAksi(t *testing.T) {
	dt := &deductionTiruan{tarif: &tarifUji}
	svc := layananDeduction(dt)
	ctx := context.Background()
	got, err := svc.HitungDeduction(ctx, "NB-1", AksiPersen, 0, isianUji("0", "0", "0", "0", "0", "0"))
	if err != nil || barisUji(t, got.MataUang) != "IDR/6/0.6/0/0.3/0.006/0.033 USD/1/0.1/0/0.05/0.001/0.0055" || !got.PPnCheck ||
		t8(got.TarifPPH) != "2" || len(dt.tulis) != 0 {
		t.Fatalf("persen: %v %s", err, barisUji(t, got.MataUang))
	}
	// SetTotalPremium_Act: 1,2 / 6 × 100 = 20 % -> USD 0,2.
	got, err = svc.HitungDeduction(ctx, "NB-1", AksiCommision, 0, isianUji("1.2", "0", "0.3", "0.1", "0", "0.05"))
	if err != nil || t8(got.PctPremiumDeduction) != "20" || t8(got.MataUang[1].Commision) != "0.2" || t8(got.PctBrokerageFee) != "5" {
		t.Fatalf("commision: %v %s", err, barisUji(t, got.MataUang))
	}
	// ProtekDeduction2_Act: USD 0,5 / 1 = 50 % -> IDR 3.
	got, err = svc.HitungDeduction(ctx, "NB-1", AksiDeduction2, 1, isianUji("0", "0", "0", "0", "0.5", "0"))
	if err != nil || t8(got.MataUang[0].Deduction2) != "3" || t8(got.MataUang[1].Deduction2) != "0.5" {
		t.Fatalf("deduction2: %v %s", err, barisUji(t, got.MataUang))
	}
	if _, err := svc.HitungDeduction(ctx, "NB-1", AksiDeduction2, 1, isianUji("0", "0", "0", "0", "2", "0")); !errors.Is(err, ErrMasukanDeduction) ||
		!strings.HasSuffix(err.Error(), ": "+pesanDeduction2Lebih) {
		t.Fatalf("deduction2 > premi: %v", err)
	}
	// TotalBrokerageFee_Act: 0,6 / 6 × 100 = 10 % -> USD 0,1; PPh / PPN ulang.
	got, err = svc.HitungDeduction(ctx, "NB-1", AksiBrokerage, 0, isianUji("0", "0", "0.6", "0", "0", "0"))
	if err != nil || t8(got.PctBrokerageFee) != "10" || barisUji(t, got.MataUang) != "IDR/6/0/0/0.6/0.012/0.066 USD/1/0/0/0.1/0.002/0.011" {
		t.Fatalf("brokerage: %v %s", err, barisUji(t, got.MataUang))
	}
	for nama, x := range map[string]struct {
		aksi  string
		idx   int
		isian IsianDeduction
		mau   error
	}{
		"aksi":           {"lain", 0, isianUji("0", "0", "0", "0", "0", "0"), ErrMasukanDeduction},
		"idx":            {AksiCommision, 2, isianUji("0", "0", "0", "0", "0", "0"), ErrMasukanDeduction},
		"mata uang":      {AksiPersen, 0, IsianDeduction{MataUang: []models.BarisDeduction{{Currency: "IDR"}}}, ErrDeductionBerubah},
		"mata uang lain": {AksiPersen, 0, IsianDeduction{MataUang: []models.BarisDeduction{{Currency: "IDR"}, {Currency: "EUR"}}}, ErrDeductionBerubah},
	} {
		if _, err := svc.HitungDeduction(ctx, "NB-1", x.aksi, x.idx, x.isian); !errors.Is(err, x.mau) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	// Premi 0 (PS 0): persen dari baris tidak dapat dihitung (K50-5).
	nol := Baru(nil).DenganDeduction(dt).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("0", "USD")}})
	if _, err := nol.HitungDeduction(ctx, "NB-1", AksiCommision, 0, isianUji("1", "0", "0", "0", "0", "0")); !errors.Is(err, ErrMasukanDeduction) ||
		!strings.Contains(err.Error(), "Premium IDR 0") {
		t.Fatalf("premi 0: %v", err)
	}
}

// TestTampilanDanSimpanDeduction - GET: nilai tersimpan dipakai bila premi sama, dihitung ulang dari persen bila premi
// berubah; Deduction2 tersimpan; Deduction2 > premi -> pesan; TAX kosong -> pesan. PUT: SetRIComm + Deduction2 apa adanya.
func TestTampilanDanSimpanDeduction(t *testing.T) {
	dt := &deductionTiruan{tarif: &tarifUji, ada: models.KasusDeduction{RICommision: d("10"), PctBrokerageFee: d("5"), PPnCheck: "true",
		MataUang: []models.BarisDeduction{
			{Currency: "IDR", Premium: d("6"), Commision: d("0.61"), Deduction2: d("1"), BrokerageFee: d("0.3"), PPh: d("0.006"), PPN: d("0.033")},
			{Currency: "USD", Premium: d("999"), Commision: d("99"), Deduction2: d("5"), BrokerageFee: d("9"), PPh: d("9"), PPN: d("9")}}}}
	svc := layananDeduction(dt)
	ctx := context.Background()
	got, err := svc.TampilanDeduction(ctx, "NB-1")
	if err != nil || barisUji(t, got.MataUang) != "IDR/6/0.61/1/0.3/0.006/0.033 USD/1/0.1/5/0.05/0.001/0.0055" || !got.PPnCheck ||
		len(got.Pesan) != 1 || got.Pesan[0] != pesanDeduction2Lebih {
		t.Fatalf("GET: %v %s %v", err, barisUji(t, got.MataUang), got.Pesan)
	}
	if _, err := svc.SimpanDeduction(ctx, pelakuUji, "NB-1", isianUji("0", "7", "0", "0", "0", "0")); !errors.Is(err, ErrMasukanDeduction) ||
		!strings.Contains(err.Error(), pesanDeduction2Lebih+" (IDR)") || len(dt.tulis) != 0 {
		t.Fatalf("PUT deduction2 > premi: %v", err)
	}
	if _, err := svc.SimpanDeduction(ctx, pelakuUji, "NB-1", isianUji("0", "3", "0", "0", "0.5", "0")); err != nil || len(dt.tulis) != 1 {
		t.Fatalf("PUT: %v", err)
	}
	w := dt.tulis[0]
	if t8(w.RICommision) != "10" || w.PPnCheck != "true" || barisUji(t, w.MataUang) != "IDR/6/0.6/3/0.3/0.006/0.033 USD/1/0.1/0.5/0.05/0.001/0.0055" {
		t.Fatalf("ditulis: %+v %s", w, barisUji(t, w.MataUang))
	}
	dt.tarif = nil
	got, err = svc.TampilanDeduction(ctx, "NB-1")
	if err != nil || got.TarifPPH != nil || len(got.Pesan) != 1 || !strings.Contains(got.Pesan[0], "TAX") {
		t.Fatalf("tanpa TAX: %v %v", err, got.Pesan)
	}
	if _, err := svc.TampilanDeduction(ctx, "K404"); !errors.Is(err, ErrKasusTidakAda) {
		t.Fatalf("404: %v", err)
	}
	if _, err := Baru(nil).TampilanDeduction(ctx, "NB-1"); !errors.Is(err, ErrDeductionTanpaDatabase) {
		t.Fatalf("503: %v", err)
	}
}

// TestDeductionMenyegarkanCicilan - K50-12: Save Premium Deduction menulis ulang cicilan tersimpan dengan nilai baru (persen
// tetap); tanpa cicilan tidak menulis apa pun ke Payment.
func TestDeductionMenyegarkanCicilan(t *testing.T) {
	dt := &deductionTiruan{tarif: &tarifUji}
	p := &paymentTiruan{ada: models.KasusPayment{Installment: "1", Cicilan: map[string][]models.CicilanPayment{
		"IDR": {{InstallmentNo: 1, DueDate: "20250402", Persen: d("100")}}, "USD": {{InstallmentNo: 1, DueDate: "20250402", Persen: d("100")}}}}}
	svc := Baru(nil).DenganDeduction(dt).DenganPayment(p).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "USD")}}).
		DenganTransaksi(tanpaTx)
	if _, err := svc.SimpanDeduction(context.Background(), pelakuUji, "NB-1", isianUji("0", "0", "0", "0", "0", "0")); err != nil {
		t.Fatal(err)
	}
	if len(p.tulis) != 1 || p.tulis[0].Installment != 1 || t8(p.tulis[0].MataUang[0].Cicilan[0].RICommision) != "0.6" ||
		t8(p.tulis[0].MataUang[0].SumTotalPayment) != "5.139" { // 6 − 0.6 − 0 − 0.3 + 0.006 + 0.033
		t.Fatalf("cicilan disegarkan: %+v", p.tulis)
	}
	kosong := &paymentTiruan{ada: models.KasusPayment{Cicilan: map[string][]models.CicilanPayment{}}}
	svc = Baru(nil).DenganDeduction(dt).DenganPayment(kosong).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "USD")}}).
		DenganTransaksi(tanpaTx)
	if _, err := svc.SimpanDeduction(context.Background(), pelakuUji, "NB-1", isianUji("0", "0", "0", "0", "0", "0")); err != nil || len(kosong.tulis) != 0 {
		t.Fatalf("tanpa cicilan: %v %d", err, len(kosong.tulis))
	}
}

// TestDeductionTidakBergeser - K50-9: nilai yang diketik (Brokerage Fee 100000 atas premi 469550) tetap 100000 sesudah PUT
// walau persen tampil 8 desimal (21.29698648); mata uang lain ikut persen berpresisi penuh yang sama.
func TestDeductionTidakBergeser(t *testing.T) {
	k := models.KasusSpreading{PercentShare: d("100"), Lokasi: []models.LokasiSpreading{{Items: []models.ItemSpreading{
		{Currency: "IDR", Coverages: []models.CoverageSpreading{{OldID: "C1", TSILiability: d("1"), Premium: d("469550")}}},
		{Currency: "USD", Coverages: []models.CoverageSpreading{{OldID: "C2", TSILiability: d("1"), Premium: d("1000000")}}}}}}}
	dt := &deductionTiruan{tarif: &tarifUji}
	svc := Baru(nil).DenganDeduction(dt).DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{k}}).DenganTransaksi(tanpaTx)
	ctx := context.Background()
	isi := func(pct, idr, usd string) IsianDeduction {
		return IsianDeduction{PctPremiumDeduction: d("0"), PctBrokerageFee: d(pct), MataUang: []models.BarisDeduction{
			{Currency: "IDR", BrokerageFee: d(idr)}, {Currency: "USD", BrokerageFee: d(usd)}}}
	}
	h, err := svc.HitungDeduction(ctx, "NB-1", AksiBrokerage, 0, isi("0", "100000", "0"))
	if err != nil || t8(h.PctBrokerageFee) != "21.29698648" || t8(h.MataUang[0].BrokerageFee) != "100000" ||
		t8(h.MataUang[1].BrokerageFee) != "212969.86476414" {
		t.Fatalf("hitung: %v %s / %s / %s", err, t8(h.PctBrokerageFee), t8(h.MataUang[0].BrokerageFee), t8(h.MataUang[1].BrokerageFee))
	}
	if _, err := svc.SimpanDeduction(ctx, pelakuUji, "NB-1", isi("21.29698648", "100000", "212969.86476414")); err != nil {
		t.Fatal(err)
	}
	w := dt.tulis[0].MataUang
	if t8(w[0].BrokerageFee) != "100000" || t8(w[1].BrokerageFee) != "212969.86476414" || t8(dt.tulis[0].PctBrokerageFee) != "21.29698648" {
		t.Fatalf("tersimpan: %s / %s / %s", t8(w[0].BrokerageFee), t8(w[1].BrokerageFee), t8(dt.tulis[0].PctBrokerageFee))
	}
	// Persen diketik ulang (tidak cocok dengan grid): persen kiriman dipakai apa adanya.
	if _, err := svc.SimpanDeduction(ctx, pelakuUji, "NB-1", isi("10", "100000", "212969.86476414")); err != nil ||
		t8(dt.tulis[1].MataUang[0].BrokerageFee) != "46955" {
		t.Fatalf("persen baru: %v %s", err, t8(dt.tulis[1].MataUang[0].BrokerageFee))
	}
}
