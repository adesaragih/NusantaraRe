package services

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// ErrKlaimTidakAda dikembalikan bila klaim yang diminta tidak ada.
var ErrKlaimTidakAda = errors.New("services: klaim tidak ada")

// KlaimLife merakit agregat klaim Life dari bacaan repository.
//
// Perakitan ada DI SINI, bukan di repository dan bukan di handlers.
type KlaimLife struct {
	repo *repository.KlaimLife
}

// KlaimLife mengembalikan layanan klaim Life, atau nil bila tanpa database.
func (s *Service) KlaimLife() *KlaimLife {
	if !s.PunyaDatabase() {
		return nil
	}
	return &KlaimLife{repo: repository.NewKlaimLife(s.db)}
}

// Ambil merakit satu klaim beserta seluruh peserta dan seluruh baris
// adjustment-nya.
//
// Tiket 01 AC-1: seluruh baris `AdjustmentList`, bukan hanya yang terakhir.
// Setiap baris tetap melekat pada pesertanya; meratakannya ke header menghapus
// informasi peserta pemilik dan mematahkan mesin status (ADR-U-0011).
func (k *KlaimLife) Ambil(ctx context.Context, id string) (*models.Klaim, error) {
	if k == nil || k.repo == nil {
		return nil, repository.ErrTanpaOracle
	}
	if id == "" {
		return nil, fmt.Errorf("%w: pengenal klaim kosong", ErrWajibIsi)
	}

	klaim, err := k.repo.AmbilHeader(ctx, id)
	if err != nil {
		return nil, err
	}
	if klaim == nil {
		return nil, ErrKlaimTidakAda
	}

	peserta, err := k.repo.AmbilPeserta(ctx, id)
	if err != nil {
		return nil, err
	}
	baris, err := k.repo.AmbilBaris(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range peserta {
		if b, ada := baris[peserta[i].ID]; ada {
			peserta[i].Baris = b
		}
	}
	klaim.Peserta = peserta
	return klaim, nil
}
