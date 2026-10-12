package services

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// pagarScoring - n bacaan wajib berjalan BERSAMAAN: tiap pembaca menunggu sampai n pembaca tiba (paling lama 2 detik). Alur
// berurutan membuat pembaca pertama menunggu sampai batas -> habisWaktu. Kedatangan lebih dari n tidak menunggu.
type pagarScoring struct {
	mu         sync.Mutex
	n, datang  int
	semua      chan struct{}
	habisWaktu bool
	galat      map[string]error
}

func pagarBaru(n ...int) *pagarScoring {
	jumlah := 4
	if len(n) > 0 {
		jumlah = n[0]
	}
	return &pagarScoring{n: jumlah, semua: make(chan struct{}), galat: map[string]error{}}
}

func (p *pagarScoring) masuk(nama string) error {
	p.mu.Lock()
	p.datang++
	if p.datang == p.n {
		close(p.semua)
	}
	p.mu.Unlock()
	select {
	case <-p.semua:
	case <-time.After(2 * time.Second):
		p.mu.Lock()
		p.habisWaktu = true
		p.mu.Unlock()
	}
	return p.galat[nama]
}

type scoringPagar struct {
	*scoringTiruan
	p *pagarScoring
}

func (s scoringPagar) BacaScoring(ctx context.Context, id string) (models.KasusScoring, error) {
	if err := s.p.masuk("scoring"); err != nil {
		return models.KasusScoring{}, err
	}
	return s.scoringTiruan.BacaScoring(ctx, id)
}

type kasusPagar struct {
	kasusTiruan
	p *pagarScoring
}

func (k kasusPagar) BacaKasus(ctx context.Context, id string) (models.Kasus, error) {
	if err := k.p.masuk("kasus"); err != nil {
		return models.Kasus{}, err
	}
	return k.kasusTiruan.BacaKasus(ctx, id)
}

type objekPagar struct {
	objekTiruan
	p *pagarScoring
}

func (o objekPagar) BacaObjek(ctx context.Context, id string) ([]models.ObjekFire, error) {
	if err := o.p.masuk("objek"); err != nil {
		return nil, err
	}
	return o.objekTiruan.BacaObjek(ctx, id)
}

type spreadingPagar struct {
	*spreadingTiruan
	p *pagarScoring
}

func (s spreadingPagar) BacaSpreading(ctx context.Context, tx *db.Tx, id string) (models.KasusSpreading, error) {
	if err := s.p.masuk("spreading"); err != nil {
		return models.KasusSpreading{}, err
	}
	return s.spreadingTiruan.BacaSpreading(ctx, tx, id)
}

// layananScoringPagar - layananScoring dengan keempat pembaca dibungkus pagar.
func layananScoringPagar(st *scoringTiruan, p *pagarScoring) *Service {
	svc := layananScoring(st)
	return svc.DenganScoring(scoringPagar{st, p}).DenganKasus(kasusPagar{svc.kasus.(kasusTiruan), p}).
		DenganObjek(objekPagar{svc.objek.(objekTiruan), p}).DenganSpreading(spreadingPagar{svc.spreading.(*spreadingTiruan), p})
}

// TestKonteksScoringParalel - laporan work owner 12-10-2026 (tab Scoring "Memuat…"): BacaScoring / BacaKasus / BacaObjek /
// BacaSpreading berjalan bersamaan (pagar), hasil TampilanScoringKasus IDENTIK dengan layanan tanpa pagar, galat salah satu pembaca
// diteruskan (ErrKasusTidakAda tetap lewat galatKasus) tanpa menggantung.
func TestKonteksScoringParalel(t *testing.T) {
	ctx := context.Background()
	st := &scoringTiruan{k: models.KasusScoring{Position: repository.PosisiAwalCaseNB, FlagOnGoingPolicy: "0"}}
	mau, err := layananScoring(st).TampilanScoringKasus(ctx, "UJI-NB-1")
	if err != nil {
		t.Fatal(err)
	}
	p := pagarBaru()
	got, err := layananScoringPagar(st, p).TampilanScoringKasus(ctx, "UJI-NB-1")
	if err != nil || p.habisWaktu || !reflect.DeepEqual(got, mau) {
		t.Fatalf("paralel: habis waktu %v, galat %v\n got %+v\n mau %+v", p.habisWaktu, err, got, mau)
	}
	galatUji := errors.New("UJI galat baca")
	for nama, c := range map[string]struct {
		galat error
		mau   error
	}{
		"scoring":   {repository.ErrKasusTidakAda, ErrKasusTidakAda},
		"kasus":     {galatUji, galatUji},
		"objek":     {repository.ErrKasusTidakAda, ErrKasusTidakAda},
		"spreading": {galatUji, galatUji},
	} {
		p := pagarBaru()
		p.galat[nama] = c.galat
		mulai := time.Now()
		_, err := layananScoringPagar(st, p).TampilanScoringKasus(ctx, "UJI-NB-1")
		if !errors.Is(err, c.mau) || time.Since(mulai) > 3*time.Second {
			t.Errorf("%s: %v (%v)", nama, err, time.Since(mulai))
		}
	}
}
