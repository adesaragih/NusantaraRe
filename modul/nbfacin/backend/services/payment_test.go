package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// paymentTiruan - PenyimpanPayment dalam memori; "K404" = case tidak ada.
type paymentTiruan struct {
	ada   models.KasusPayment
	tulis []models.SimpanPayment
}

func (p *paymentTiruan) BacaPayment(_ context.Context, id string) (models.KasusPayment, error) {
	if id == "K404" {
		return models.KasusPayment{}, repository.ErrKasusTidakAda
	}
	return p.ada, nil
}

func (p *paymentTiruan) TulisPayment(_ context.Context, _ *db.Tx, _ string, k models.SimpanPayment) error {
	for _, m := range k.MataUang {
		if err := periksaMask("Total Payment "+m.Baris.Currency, m.SumTotalPayment, m.NetPremium); err != nil {
			return err
		}
		for _, c := range m.Cicilan { // INSTALLMENT_PERCENTAGE = VARCHAR2, bukan angkaMasuk
			if err := periksaMask("cicilan "+m.Baris.Currency, c.Premium, c.RICommision, c.Deduction2, c.BrokerageFee, c.PPh, c.PPN,
				c.Discount, c.PaymentTotal); err != nil {
				return err
			}
		}
	}
	p.tulis = append(p.tulis, k)
	return nil
}

// layananPayment - kasusDuaItem PS 50 (premi IDR 6, USD 1), Premium Deduction tersimpan 10 % / 5 % (PPH 2, PPN 11),
// Deduction2 IDR 0.6: Net IDR = 6 − 0.6 − 0.6 − 0.3 + 0.006 + 0.033 = 4.539; USD = 1 − 0.1 − 0.05 + 0.001 + 0.0055 = 0.8565.
func layananPayment(p *paymentTiruan) *Service {
	dt := &deductionTiruan{tarif: &tarifUji, ada: models.KasusDeduction{RICommision: d("10"), PctBrokerageFee: d("5"), PPnCheck: "true",
		MataUang: []models.BarisDeduction{{Currency: "IDR", Deduction2: d("0.6")}}}}
	src := &spreadingTiruan{baca: []models.KasusSpreading{kasusDuaItem("50", "USD")}}
	return Baru(nil).DenganPayment(p).DenganDeduction(dt).DenganSpreading(src).DenganTransaksi(tanpaTx)
}

func cicilanUji(t *testing.T, c []models.CicilanPayment) string {
	t.Helper()
	var s []string
	for _, x := range c {
		s = append(s, strings.Join([]string{x.DueDate, t8(x.Persen), t8(x.Premium), t8(x.RICommision), t8(x.Deduction2), t8(x.BrokerageFee),
			t8(x.PPh), t8(x.PPN), t8(x.Discount), t8(x.PaymentTotal)}, "/"))
	}
	return strings.Join(s, " ")
}

func TestPersenRataDanJatuhTempo(t *testing.T) {
	r := persenRata(3)
	if r[0].String() != "33.33333333333333333333" || r[2].String() != "33.33333333333333333334" || jumlahPersen([]models.CicilanPayment{
		{Persen: r[0]}, {Persen: r[1]}, {Persen: r[2]}}).Cmp(seratus) != 0 {
		t.Fatalf("rata 3: %v", r)
	}
	if r := persenRata(4); t8(r[3]) != "25" {
		t.Fatalf("rata 4: %v", r)
	}
	// Calendar.add(MONTH): 31 Jan + 1 bulan = 28 Feb (bukan 3 Mar).
	if got := tambahBulan(time.Date(2026, 1, 31, 0, 0, 0, 0, WIB), 1).Format("20060102"); got != "20260228" {
		t.Fatalf("tambahBulan: %s", got)
	}
	for mulai, mau := range map[string]string{
		"20251130T050000.000 GMT": "20260228", // 30 Nov WIB + 3 bulan -> 28 Feb (dipotong)
		"20251130T170000.000 GMT": "20260301", // 1 Des WIB
		"":                        "",
	} {
		if got := jatuhTempoBawaan(mulai); got != mau {
			t.Errorf("jatuh tempo %q = %q, mau %q", mulai, got, mau)
		}
	}
}

func TestIsiCicilan(t *testing.T) {
	m := models.MataUangPayment{Baris: models.BarisDeduction{Currency: "IDR", Premium: d("6"), Commision: d("0.6"), Deduction2: d("0.6"),
		BrokerageFee: d("0.3"), PPh: d("0.006"), PPN: d("0.033")}, Cicilan: []models.CicilanPayment{{Persen: d("50")}}}
	isiCicilan(&m)
	if t8(m.SumTotalPayment) != "4.539" || cicilanUji(t, m.Cicilan) != "/50/3/0.3/0.3/0.15/0.003/0.0165/0/2.2695" {
		t.Fatalf("Net %s cicilan %s", t8(m.SumTotalPayment), cicilanUji(t, m.Cicilan))
	}
}

func isianPaymentUji(n string, idr, usd []models.CicilanPayment) IsianPayment {
	return IsianPayment{Installment: n, MataUang: []IsianCicilanMataUang{{Currency: "IDR", Cicilan: idr}, {Currency: "USD", Cicilan: usd}}}
}

func TestHitungPayment(t *testing.T) {
	p := &paymentTiruan{}
	svc := layananPayment(p)
	ctx := context.Background()
	got, err := svc.HitungPayment(ctx, "NB-1", AksiInstallment, "", 0, IsianPayment{Installment: "2"})
	if err != nil || len(got.MataUang) != 2 || cicilanUji(t, got.MataUang[0].Cicilan) !=
		"20250402/50/3/0.3/0.3/0.15/0.003/0.0165/0/2.2695 20250402/50/3/0.3/0.3/0.15/0.003/0.0165/0/2.2695" ||
		t8(got.MataUang[1].SumTotalPayment) != "0.8565" || len(p.tulis) != 0 {
		t.Fatalf("installment: %v %+v", err, got)
	}
	c := func(ps ...string) []models.CicilanPayment {
		var out []models.CicilanPayment
		for _, x := range ps {
			out = append(out, models.CicilanPayment{Persen: d(x), DueDate: "20250402"})
		}
		return out
	}
	got, err = svc.HitungPayment(ctx, "NB-1", AksiPersenCicilan, "IDR", 0, isianPaymentUji("2", c("60", "30"), c("50", "50")))
	if err != nil || len(got.Pesan) != 1 || got.Pesan[0] != pesanPersenKurang100 || t8(got.MataUang[0].Cicilan[0].Premium) != "3.6" {
		t.Fatalf("persen 90: %v %v", err, got.Pesan)
	}
	if got, _ := svc.HitungPayment(ctx, "NB-1", AksiPersenCicilan, "IDR", 1, isianPaymentUji("2", c("60", "50"), c("50", "50"))); len(got.Pesan) != 1 ||
		got.Pesan[0] != pesanPersenLebih100 {
		t.Fatalf("persen 110: %v", got.Pesan)
	}
	for nama, x := range map[string]struct {
		aksi, cur string
		idx       int
		isian     IsianPayment
		mau       error
	}{
		"aksi":        {"lain", "", 0, IsianPayment{Installment: "1"}, ErrMasukanPayment},
		"installment": {AksiInstallment, "", 0, IsianPayment{Installment: "0"}, ErrMasukanPayment},
		"idx":         {AksiPersenCicilan, "IDR", 5, isianPaymentUji("2", c("50", "50"), c("50", "50")), ErrMasukanPayment},
		"mata uang":   {AksiPersenCicilan, "EUR", 0, isianPaymentUji("2", c("50", "50"), c("50", "50")), ErrPaymentBerubah},
		"himpunan":    {AksiPersenCicilan, "IDR", 0, IsianPayment{MataUang: []IsianCicilanMataUang{{Currency: "IDR"}}}, ErrPaymentBerubah},
	} {
		if _, err := svc.HitungPayment(ctx, "NB-1", x.aksi, x.cur, x.idx, x.isian); !errors.Is(err, x.mau) {
			t.Errorf("%s: %v", nama, err)
		}
	}
}

func TestSimpanPayment(t *testing.T) {
	p := &paymentTiruan{}
	svc := layananPayment(p)
	ctx := context.Background()
	c := func(ps ...string) []models.CicilanPayment {
		var out []models.CicilanPayment
		for _, x := range ps {
			out = append(out, models.CicilanPayment{Persen: d(x), DueDate: "20250402"})
		}
		return out
	}
	// Σ 100 dipertahankan (60 / 40); Σ 90 diratakan ulang (50 / 50).
	if _, err := svc.SimpanPayment(ctx, pelakuUji, "NB-1", isianPaymentUji("2", c("60", "40"), c("60", "30"))); err != nil {
		t.Fatal(err)
	}
	w := p.tulis[0]
	if w.Installment != 2 || t8(w.MataUang[0].Cicilan[0].Persen) != "60" || t8(w.MataUang[1].Cicilan[0].Persen) != "50" ||
		w.MataUang[1].Cicilan[1].InstallmentNo != 2 || t8(w.MataUang[0].SumTotalPayment) != "4.539" || t8(w.MataUang[0].NetPremium) != "4.539" {
		t.Fatalf("ditulis: %+v", w)
	}
	for nama, x := range map[string]struct {
		isian IsianPayment
		mau   string
	}{
		"installment kosong": {isianPaymentUji("", c("100"), c("100")), "installment wajib"},
		"jumlah beda":        {isianPaymentUji("2", c("100"), c("50", "50")), "jumlah cicilan IDR (1) tidak sama dengan Installment (2)"},
		"due date":           {isianPaymentUji("1", []models.CicilanPayment{{Persen: d("100")}}, c("100")), "cicilan IDR ke-1: Due Date wajib diisi"},
		"cicilan kosong":     {isianPaymentUji("1", nil, c("100")), pesanCicilanKosong},
	} {
		if _, err := svc.SimpanPayment(ctx, pelakuUji, "NB-1", x.isian); !errors.Is(err, ErrMasukanPayment) || !strings.Contains(err.Error(), x.mau) {
			t.Errorf("%s: %v", nama, err)
		}
	}
	if len(p.tulis) != 1 {
		t.Fatalf("tulis %d, mau 1 (proteksi menolak tanpa menyimpan)", len(p.tulis))
	}
}

// TestProteksiPayment - ProtectPremiPolicy_Act: List Payment kosong; premi 0; total tidak cocok hanya bila KETIGA selisih
// >= 0.001 (rantai OR Pega, K51-6).
func TestProteksiPayment(t *testing.T) {
	if m := proteksiPayment(nil, ""); len(m) != 1 || m[0] != pesanPaymentKosong {
		t.Fatalf("kosong: %v", m)
	}
	baris := models.BarisDeduction{Currency: "IDR", Premium: d("6"), Commision: d("0.6")}
	cocok := models.MataUangPayment{Baris: baris, NetPremium: d("5.4"), Cicilan: []models.CicilanPayment{{Premium: d("6"), RICommision: d("0.6"), PaymentTotal: d("5.4")}}}
	if m := proteksiPayment([]models.MataUangPayment{cocok}, ""); len(m) != 0 {
		t.Fatalf("cocok: %v", m)
	}
	satuBeda := cocok
	satuBeda.Cicilan = []models.CicilanPayment{{Premium: d("5"), RICommision: d("0.6"), PaymentTotal: d("5.4")}}
	if m := proteksiPayment([]models.MataUangPayment{satuBeda}, ""); len(m) != 0 {
		t.Fatalf("satu beda (verbatim lolos): %v", m)
	}
	semuaBeda := cocok
	semuaBeda.Cicilan = []models.CicilanPayment{{Premium: d("5"), RICommision: d("0.5"), PaymentTotal: d("4")}}
	if m := proteksiPayment([]models.MataUangPayment{semuaBeda}, ""); len(m) != 1 || m[0] != pesanTotalTakCocok {
		t.Fatalf("semua beda: %v", m)
	}
	nol := models.MataUangPayment{Baris: models.BarisDeduction{Currency: "IDR", Premium: d("0")}, NetPremium: d("0"),
		Cicilan: []models.CicilanPayment{{Premium: d("0"), RICommision: d("0"), PaymentTotal: d("0")}}}
	if m := proteksiPayment([]models.MataUangPayment{nol}, ""); len(m) != 1 || m[0] != pesanPremiNol {
		t.Fatalf("premi 0: %v", m)
	}
	// Work owner 09-10-2026 "Samakan dengan Pega": langkah 11 dilewati untuk Master Policy (kode "1"); kode lain tetap.
	for pt, ditahan := range map[string]bool{"1": false, "0": true, "2": true} {
		if m := proteksiPayment([]models.MataUangPayment{nol}, pt); (len(m) == 1 && m[0] == pesanPremiNol) != ditahan {
			t.Errorf("premi 0 PolicyType %q: %v", pt, m)
		}
	}
}

// TestPolicyTypeKasusPayment - Save Payment membaca Policy Type case lewat kodePolicyType (label lama "Master Policy" -> "1");
// tanpa pembaca kasus -> "" (tanpa gerbang); case tidak ada -> 404.
func TestPolicyTypeKasusPayment(t *testing.T) {
	ctx := context.Background()
	kasus := &models.Kasus{CaseID: "NB-1", General: models.General{PolicyType: "Master Policy"}}
	svc := Baru(nil).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"NB-1": kasus}})
	if pt, err := svc.policyTypeKasus(ctx, "NB-1"); err != nil || pt != "1" {
		t.Errorf("label: %q %v", pt, err)
	}
	if _, err := svc.policyTypeKasus(ctx, "NB-9"); !errors.Is(err, ErrKasusTidakAda) {
		t.Errorf("404: %v", err)
	}
	if pt, err := Baru(nil).policyTypeKasus(ctx, "NB-1"); err != nil || pt != "" {
		t.Errorf("tanpa pembaca: %q %v", pt, err)
	}
}

func TestTampilanPayment(t *testing.T) {
	p := &paymentTiruan{ada: models.KasusPayment{Installment: "1", Cicilan: map[string][]models.CicilanPayment{
		"IDR": {{InstallmentNo: 1, DueDate: "20250402", Persen: d("100")}}}}}
	got, err := layananPayment(p).TampilanPayment(context.Background(), "NB-1")
	if err != nil || got.Installment != "1" || len(got.MataUang) != 2 || len(got.MataUang[1].Cicilan) != 0 ||
		cicilanUji(t, got.MataUang[0].Cicilan) != "20250402/100/6/0.6/0.6/0.3/0.006/0.033/0/4.539" || t8(got.TarifPPH) != "2" {
		t.Fatalf("GET: %v %+v", err, got)
	}
	if _, err := layananPayment(p).TampilanPayment(context.Background(), "K404"); !errors.Is(err, ErrKasusTidakAda) {
		t.Fatalf("404: %v", err)
	}
	if _, err := Baru(nil).TampilanPayment(context.Background(), "NB-1"); !errors.Is(err, ErrPaymentTanpaDatabase) {
		t.Fatalf("503: %v", err)
	}
}

// TestSimpanPaymentMasterPolicy - Save Payment premi 0 (% Share 0): Individual Policy -> 400 memuat "Premi on paymnet can't be 0";
// Master Policy (label lama maupun kode) -> pesan itu TIDAK ada (langkah 11 dilewati, work owner 09-10-2026 "Samakan dengan Pega").
// Galat lain fixture (Total Payment negatif karena Deduction2) tetap ada di keduanya - hanya pesan premi 0 yang diuji.
func TestSimpanPaymentMasterPolicy(t *testing.T) {
	ctx := context.Background()
	c := func(ps ...string) []models.CicilanPayment {
		var out []models.CicilanPayment
		for _, x := range ps {
			out = append(out, models.CicilanPayment{Persen: d(x), DueDate: "20250402"})
		}
		return out
	}
	for pt, ditahan := range map[string]bool{"0": true, "Individual Policy": true, "1": false, "Master Policy": false} {
		p := &paymentTiruan{}
		kasus := &models.Kasus{CaseID: "NB-1", General: models.General{PolicyType: pt}}
		svc := layananPayment(p).DenganKasus(kasusTiruan{ada: map[string]*models.Kasus{"NB-1": kasus}})
		svc.spreading.(*spreadingTiruan).baca = []models.KasusSpreading{kasusDuaItem("0", "USD")}
		_, err := svc.SimpanPayment(ctx, pelakuUji, "NB-1", isianPaymentUji("2", c("60", "40"), c("60", "40")))
		if !errors.Is(err, ErrMasukanPayment) || strings.Contains(err.Error(), pesanPremiNol) != ditahan || len(p.tulis) != 0 {
			t.Errorf("%q: %v (tulis %d)", pt, err, len(p.tulis))
		}
	}
}
