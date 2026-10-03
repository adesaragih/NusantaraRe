package services

// Pencarian calon peserta klaim - tiket 02.
//
// Dibaca sesudah: repository/pesertapolis.go.

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimlife/backend/repository"
)

// PesertaLayanan membungkus pembaca peserta polis.
type PesertaLayanan struct{ svc *Service }

// Peserta menyusun layanan pencarian peserta.
func (s *Service) Peserta() *PesertaLayanan { return &PesertaLayanan{svc: s} }

// Cari mengembalikan calon peserta satu premium list.
//
// Penyaring peserta hidup dan batas hasil dipasang di repository, bukan di
// sini: AC 29 tiket 02 menuntut penyaringan terjadi di SATU tempat.
func (p *PesertaLayanan) Cari(ctx context.Context, nomorPremiList, sertifikat, nama string,
	batas int) ([]repository.CalonPeserta, error) {
	if !p.svc.PunyaDatabase() {
		return nil, db.ErrTanpaOracle
	}
	return repository.NewPesertaPolis(p.svc.DB()).
		Cari(ctx, nomorPremiList, sertifikat, nama, batas)
}
