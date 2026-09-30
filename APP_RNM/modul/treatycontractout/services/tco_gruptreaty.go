package services

// Grup treaty untuk pemilih `Treaty Group` - tiket 03 Treaty Contract Out.
//
// Master `TREATYGROUP` dibaca saja (spec b107); kosong = kegagalan terlihat
// (ADR-0015), pola JenisReasuransiTreaty.

import (
	"context"
	"errors"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/treatycontractout/repository"
)

var (
	// ErrMasterGrupTreatyKosong - master grup treaty kosong atau tidak terbaca.
	ErrMasterGrupTreatyKosong = errors.New(
		"services: treaty group master TREATYGROUP is empty or unreadable")
	// ErrPembacaGrupTreatyBelumDisuntik - handler lupa memasang pembaca.
	ErrPembacaGrupTreatyBelumDisuntik = errors.New(
		"services: treaty group master reader is not injected")
)

// PembacaGrupTreatyTCO membaca master grup treaty.
type PembacaGrupTreatyTCO interface {
	Daftar(ctx context.Context) ([]repository.GrupTreatyTCO, error)
}

// GrupTreaty adalah satu pilihan untuk layar (`.ID`, `.TreatyGroupName`).
type GrupTreaty struct {
	ID              string `json:"id"`
	TreatyGroupName string `json:"treatyGroupName"`
}

type pembacaGrupTreatyBelumDisuntik struct{}

func (pembacaGrupTreatyBelumDisuntik) Daftar(context.Context) ([]repository.GrupTreatyTCO, error) {
	return nil, ErrPembacaGrupTreatyBelumDisuntik
}

type pembacaGrupTreatyOracle struct{ svc *Service }

func (p pembacaGrupTreatyOracle) Daftar(ctx context.Context) ([]repository.GrupTreatyTCO, error) {
	if !p.svc.PunyaDatabase() {
		return nil, db.ErrTanpaOracle
	}
	return repository.NewMasterGrupTreaty(p.svc.DB()).Daftar(ctx)
}

// PembacaGrupTreatyOracle adalah pembaca sungguhan, dipasang handler.
func PembacaGrupTreatyOracle(svc *Service) PembacaGrupTreatyTCO {
	return pembacaGrupTreatyOracle{svc: svc}
}

// GrupTreatyTreaty melayani daftar grup treaty.
type GrupTreatyTreaty struct {
	svc     *Service
	pembaca PembacaGrupTreatyTCO
}

// GrupTreaty menyusun layanannya dengan pembaca yang GAGAL TERANG.
func (s *Service) GrupTreaty() *GrupTreatyTreaty {
	return &GrupTreatyTreaty{svc: s, pembaca: pembacaGrupTreatyBelumDisuntik{}}
}

// DenganPembaca memasang pembaca master.
func (g *GrupTreatyTreaty) DenganPembaca(p PembacaGrupTreatyTCO) *GrupTreatyTreaty {
	salin := *g
	salin.pembaca = p
	return &salin
}

// Daftar mengembalikan seluruh grup treaty, urutan `.ID DESC` (RD b587).
func (g *GrupTreatyTreaty) Daftar(ctx context.Context, pelaku inti.Pelaku) ([]GrupTreaty, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return nil, err
	}
	baris, err := g.pembaca.Daftar(ctx)
	if err != nil {
		return nil, err
	}
	if len(baris) == 0 {
		return nil, ErrMasterGrupTreatyKosong
	}
	out := make([]GrupTreaty, 0, len(baris))
	for _, b := range baris {
		out = append(out, GrupTreaty{ID: b.ID, TreatyGroupName: b.TreatyGroupName})
	}
	return out, nil
}
