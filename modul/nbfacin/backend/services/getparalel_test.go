package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Pembungkus pagar (pagarScoring, scoringparalel_test.go) untuk GET Payment / Deduction / Cedant / Submit (laporan work owner
// 12-10-2026: "mau, jangan halu").

type paymentPagar struct {
	*paymentTiruan
	p *pagarScoring
}

func (x paymentPagar) BacaPayment(ctx context.Context, id string) (models.KasusPayment, error) {
	if err := x.p.masuk("payment"); err != nil {
		return models.KasusPayment{}, err
	}
	return x.paymentTiruan.BacaPayment(ctx, id)
}

type deductionPagar struct {
	*deductionTiruan
	p *pagarScoring
}

func (x deductionPagar) BacaDeduction(ctx context.Context, id string) (models.KasusDeduction, error) {
	if err := x.p.masuk("deduction"); err != nil {
		return models.KasusDeduction{}, err
	}
	return x.deductionTiruan.BacaDeduction(ctx, id)
}

func (x deductionPagar) TarifPajak(ctx context.Context) (models.TarifPajak, bool, error) {
	if err := x.p.masuk("pajak"); err != nil {
		return models.TarifPajak{}, false, err
	}
	return x.deductionTiruan.TarifPajak(ctx)
}

type cedantPagar struct {
	*cedantTiruan
	p *pagarScoring
}

func (x cedantPagar) BacaCedant(ctx context.Context, id string) (models.KasusCedant, error) {
	if err := x.p.masuk("cedant"); err != nil {
		return models.KasusCedant{}, err
	}
	return x.cedantTiruan.BacaCedant(ctx, id)
}

type submitPagar struct {
	*submitTiruan
	p *pagarScoring
}

func (x submitPagar) RiwayatAkseptasi(ctx context.Context, id string) ([]models.RiwayatAkseptasi, error) {
	if err := x.p.masuk("riwayat"); err != nil {
		return nil, err
	}
	return x.submitTiruan.RiwayatAkseptasi(ctx, id)
}

func (x submitPagar) OpsiBanding(ctx context.Context, idPega, tim string) ([]string, error) {
	if err := x.p.masuk("banding"); err != nil {
		return nil, err
	}
	return x.submitTiruan.OpsiBanding(ctx, idPega, tim)
}

func (x submitPagar) SebabDecline(ctx context.Context) ([]models.BarisCauseOfDecline, error) {
	if err := x.p.masuk("sebab"); err != nil {
		return nil, err
	}
	return x.submitTiruan.SebabDecline(ctx)
}

// periksaParalel - hasil berpagar = hasil tanpa pagar, pagar tidak habis waktu; tiap nama galat diteruskan (errors.Is) tanpa
// menggantung.
func periksaParalel(t *testing.T, apa string, n int, tanpa func() (any, error), dengan func(*pagarScoring) (any, error), galat map[string]error) {
	t.Helper()
	mau, err := tanpa()
	if err != nil {
		t.Fatalf("%s tanpa pagar: %v", apa, err)
	}
	p := pagarBaru(n)
	got, err := dengan(p)
	if err != nil || p.habisWaktu || !reflect.DeepEqual(got, mau) {
		t.Fatalf("%s paralel: habis waktu %v galat %v\n got %+v\n mau %+v", apa, p.habisWaktu, err, got, mau)
	}
	for nama, g := range galat {
		p := pagarBaru(n)
		p.galat[nama] = g
		mulai := time.Now()
		if _, err := dengan(p); !errors.Is(err, g) && !(errors.Is(g, repository.ErrKasusTidakAda) && errors.Is(err, ErrKasusTidakAda)) ||
			time.Since(mulai) > 3*time.Second {
			t.Errorf("%s galat %s: %v (%v)", apa, nama, err, time.Since(mulai))
		}
	}
}

func TestGETParalel(t *testing.T) {
	ctx := context.Background()
	uji := errors.New("UJI galat baca")
	// Payment: BacaPayment ‖ (Deduction: BacaDeduction ‖ BacaSpreading ‖ TarifPajak) ‖ BacaSpreading = 5 bacaan bersamaan.
	periksaParalel(t, "Payment", 5, func() (any, error) {
		return layananPayment(&paymentTiruan{}).TampilanPayment(ctx, "NB-1")
	}, func(p *pagarScoring) (any, error) {
		svc := layananPayment(&paymentTiruan{})
		svc.DenganPayment(paymentPagar{svc.payment.(*paymentTiruan), p}).DenganDeduction(deductionPagar{svc.deduction.(*deductionTiruan), p}).
			DenganSpreading(spreadingPagar{svc.spreading.(*spreadingTiruan), p})
		return svc.TampilanPayment(ctx, "NB-1")
	}, map[string]error{"payment": repository.ErrKasusTidakAda, "deduction": uji, "pajak": uji, "spreading": uji})

	dt := func() *deductionTiruan {
		return &deductionTiruan{tarif: &tarifUji, ada: models.KasusDeduction{RICommision: d("10"), PctBrokerageFee: d("5")}}
	}
	periksaParalel(t, "Deduction", 3, func() (any, error) {
		return layananDeduction(dt()).TampilanDeduction(ctx, "NB-1")
	}, func(p *pagarScoring) (any, error) {
		svc := layananDeduction(dt())
		svc.DenganDeduction(deductionPagar{svc.deduction.(*deductionTiruan), p}).DenganSpreading(spreadingPagar{svc.spreading.(*spreadingTiruan), p})
		return svc.TampilanDeduction(ctx, "NB-1")
	}, map[string]error{"deduction": repository.ErrKasusTidakAda, "pajak": uji, "spreading": uji})

	ct := func() *cedantTiruan { return &cedantTiruan{ada: models.KasusCedant{SobName: "UJI SOB"}} }
	periksaParalel(t, "Cedant", 2, func() (any, error) {
		return layananCedant(ct()).TampilanCedant(ctx, "NB-1")
	}, func(p *pagarScoring) (any, error) {
		svc := layananCedant(ct())
		svc.DenganCedant(cedantPagar{svc.cedant.(*cedantTiruan), p}).DenganSpreading(spreadingPagar{svc.spreading.(*spreadingTiruan), p})
		return svc.TampilanCedant(ctx, "NB-1")
	}, map[string]error{"cedant": repository.ErrKasusTidakAda, "spreading": uji})

	periksaParalel(t, "Tampilan Submit", 3, func() (any, error) {
		return layananSubmit(submitUji()).TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1")
	}, func(p *pagarScoring) (any, error) {
		st := submitUji()
		return layananSubmit(st).DenganSubmit(submitPagar{st, p}).TampilanSubmit(ctx, pelakuMkt, "UJI-NB-1")
	}, map[string]error{"riwayat": uji, "banding": uji, "sebab": uji})
}

// kursHitung - sumberKurs yang menghitung panggilan; mata uang "GALAT" melempar galat (tidak di-cache).
type kursHitung struct{ id, kurs int }

func (k *kursHitung) IDMataUang(_ context.Context, nama string) (string, error) {
	k.id++
	if nama == "GALAT" {
		return "", errors.New("UJI galat kurs")
	}
	return "ID-" + nama, nil
}

func (k *kursHitung) KursPada(_ context.Context, id, _ string) (*apd.Decimal, error) {
	k.kurs++
	if id == "ID-TANPA" {
		return nil, nil
	}
	return d("16000"), nil
}

// TestKursMemo - cache kurs per permintaan Scoring: ulangan ID / kurs (termasuk "tidak ada") dibaca sekali; galat tidak di-cache.
func TestKursMemo(t *testing.T) {
	ctx := context.Background()
	src := &kursHitung{}
	m := kursMemoBaru(src)
	for i := 0; i < 3; i++ {
		if id, _ := m.IDMataUang(ctx, "USD"); id != "ID-USD" {
			t.Fatal(id)
		}
		if k, _ := m.KursPada(ctx, "ID-USD", "20250101"); k == nil || k.String() != "16000" {
			t.Fatal(k)
		}
		if k, err := m.KursPada(ctx, "ID-TANPA", "20250101"); k != nil || err != nil {
			t.Fatal(k, err)
		}
		_, _ = m.IDMataUang(ctx, "GALAT")
	}
	if _, _ = m.KursPada(ctx, "ID-USD", "20260101"); src.id != 1+3 || src.kurs != 3 {
		t.Errorf("panggilan id %d kurs %d (mau 4 / 3)", src.id, src.kurs)
	}
}

// TestRingkasanTotalParalel - Submit (tiket 76 ronde 2): dua cabang total (TSI jenjang ‖ Total RNM) membaca pohon spreading
// BERSAMAAN, hasil identik dengan tanpa pagar; galat satu cabang tidak membatalkan cabang lain (toleran, seperti sebelumnya).
func TestRingkasanTotalParalel(t *testing.T) {
	ctx := context.Background()
	baru := func() *Service {
		sp := spreadingSah()
		svc := layananSubmitUW(submitUji(), posisiJuniorUWB, sp)
		svc.DenganSpreading(&spreadingTiruan{baca: []models.KasusSpreading{sp}, idUang: map[string]string{"IDR": idMataUangIDR}})
		return svc
	}
	mau := baru().ringkasanTotalPolis(ctx, "UJI-NB-1")
	p := pagarBaru(2)
	svc := baru()
	svc.DenganSpreading(spreadingPagar{svc.spreading.(*spreadingTiruan), p})
	got := svc.ringkasanTotalPolis(ctx, "UJI-NB-1")
	if p.habisWaktu || !reflect.DeepEqual(got, mau) || got.TotalTSINusaRe == nil || got.TotalPremiNusaRe == nil {
		t.Fatalf("habis waktu %v\n got %+v\n mau %+v", p.habisWaktu, got, mau)
	}
}

// spreadingHitungKurs - spreadingTiruan yang menghitung bacaan kurs (IDMataUang + KursPada).
type spreadingHitungKurs struct {
	*spreadingTiruan
	n *int
}

func (s spreadingHitungKurs) KursPada(ctx context.Context, id, waktu string) (*apd.Decimal, error) {
	*s.n++
	return s.spreadingTiruan.KursPada(ctx, id, waktu)
}

func (s spreadingHitungKurs) IDMataUang(ctx context.Context, nama string) (string, error) {
	*s.n++
	return s.spreadingTiruan.IDMataUang(ctx, nama)
}

// TestScoringKursSekaliPerPermintaan - cache kurs per permintaan (kursMemo): kolom Top Risk dan Dominant = lokasi yang sama (IDR)
// -> keIDRScoring dua kali atas IDR; ID mata uang dibaca SEKALI (tanpa cache: dua kali); hasil sama dengan layanan tanpa penghitung.
func TestScoringKursSekaliPerPermintaan(t *testing.T) {
	ctx := context.Background()
	st := &scoringTiruan{k: models.KasusScoring{Position: repository.PosisiAwalCaseNB, FlagOnGoingPolicy: "0"}}
	masukan := [2]models.KolomScoring{{}, {IndexLocation: "2"}}
	mau, err := layananScoring(st).HitungScoring(ctx, "UJI-NB-1", masukan)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	svc := layananScoring(st)
	svc.DenganSpreading(spreadingHitungKurs{svc.spreading.(*spreadingTiruan), &n})
	got, err := svc.HitungScoring(ctx, "UJI-NB-1", masukan)
	if err != nil || !reflect.DeepEqual(got, mau) || n != 1 {
		t.Errorf("bacaan kurs %d (mau 1), galat %v", n, err)
	}
}
