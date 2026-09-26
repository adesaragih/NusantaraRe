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

// JalankanMigrasi membentuk tabel di basis data yang dikonfigurasi.
//
// Ia dipanggil oleh `go run ./cmd/api -migrate` (target `make migrate`).
// Pelarinya menolak berjalan bila lingkungan menunjuk produksi Pega
// (ADR-U-0005), dan aman dijalankan berulang kali.
func (s *Service) JalankanMigrasi(ctx context.Context) (repository.LaporanMigrasi, error) {
	if !s.PunyaDatabase() {
		return repository.LaporanMigrasi{}, repository.ErrTanpaOracle
	}
	return s.db.JalankanMigrasi(ctx)
}

// BongkarMigrasi menjalankan jalur mundur tiap langkah yang TERCATAT selesai.
//
// ⛔ Ia MENGHAPUS tabel. Sampai 26-09-2026 satu-satunya pemanggilnya adalah
// skema uji, sehingga orang yang ingin membongkar skema uji sendiri terpaksa
// menyalin isi berkas *_down.sql ke sqlplus - dan itu melewati pengaman
// T_MIGRASI, yang hanya membongkar langkah yang benar-benar tercatat.
//
// Pemanggilnya WAJIB memagari lebih dulu lewat Config.PastikanSkemaUji.
// Lapisan ini tidak membaca environment sendiri.
func (s *Service) BongkarMigrasi(ctx context.Context) (repository.LaporanMigrasi, error) {
	if !s.PunyaDatabase() {
		return repository.LaporanMigrasi{}, repository.ErrTanpaOracle
	}
	return s.db.BongkarMigrasi(ctx)
}
