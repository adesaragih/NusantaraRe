// Package services memuat aturan modul Treaty Group OJK (keputusan work owner 05-10-2026).
//
// Tabel warisan `POOLDATA.TREATYGROUPOJK` (16 baris DEV) tidak punya layar Pega di korpus: Pega hanya membacanya
// (`FetchTreatyGroupOJK`) untuk mengisi kolom OJK baris `TREATYGROUP`. Aturan modul ini keputusan work owner:
//   - Add dan Edit; TANPA hapus ("biarkan saja") - setiap OJK dipakai Treaty Group.
//   - Name dan Name (IDN) wajib, huruf besar; Name tidak boleh kembar.
//   - ID baru = nomor ID tertinggi + 1, dua digit (`17`); tabel tanpa sequence.
//   - Order No TIDAK di layar ("ORDERNO hide aja, isi sesuai max dari order no"): Add = Order No tertinggi + 1, Edit
//     membiarkannya.
//   - Salinan nama OJK di `TREATYGROUP` (OJKBUSINESSNAME, OJKBUSINESSNAMEIDN, ORDERNO) TIDAK ikut diubah ("jgn ada ubah
//     data") - modul ini hanya menulis `TREATYGROUPOJK`.
//   - Hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatygroupojk/backend/models"
	"nusantarare/modul/treatygroupojk/backend/repository"
)

// Service membawa akar bersama.
type Service struct{ *inti.Dasar }

// DariDasar membungkus akar bersama.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// PunyaDatabase - Oracle terpasang?
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// LayananOracle merakit layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi)
}

var _ Gudang = (*repository.Gudang)(nil)

type dbTx = db.Tx

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("services: treaty group OJK tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumAda - tabel TREATYGROUPOJK tidak ada di skema ini.
	ErrBelumAda = repository.ErrBelumAda
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

func dilarang(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrDilarang, fmt.Sprintf(format, a...))
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, kata string) ([]models.Ojk, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Ojk, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Ojk, error)
	// KunciNomorTertinggi - kunci seluruh baris, lalu nomor ID dan Order No tertinggi (di dalam transaksi Add).
	KunciNomorTertinggi(ctx context.Context, tx *dbTx) (id, order int, err error)
	Sisip(ctx context.Context, tx *dbTx, o models.Ojk) error
	Ubah(ctx context.Context, tx *dbTx, o models.Ojk) error
}

// Layanan adalah aturan modul di atas Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
}

// BaruLayanan membuat layanan.
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx} }

// Aktor - pelaku permintaan: akun login dan hak menunya (Full = bukan View only).
type Aktor struct {
	AkunID string
	Penuh  bool
}

func wajibPenuh(a Aktor) error {
	if a.AkunID == "" {
		return dilarang("Log in to save")
	}
	if !a.Penuh {
		return dilarang("Your access to the Treaty Group OJK menu is View only")
	}
	return nil
}
