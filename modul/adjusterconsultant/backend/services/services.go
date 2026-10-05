// Package services memuat aturan modul Adjuster Consultant (keputusan work owner 05-10-2026).
//
// Sumber aturan - layar master Pega `MstAdjusterConsultant` (folder korpus Claim Fac In dan Claim Prop):
//   - `NewAdjustConsult_Act` / `EditMstConsultant_Act` / `SaveAdjusterConsultant_Act`: Save menulis NAME dan ADDRESS
//     huruf besar, TELPNO apa adanya, USERNAME = operator, EDITDATE = sekarang. ID baru = `GetIDConsultanAdj_SQL`
//     (ID situs aktif `M_SITE_DATABASE` + `ADJUSTERCONSULTANT_SEQ` 4 digit).
//   - Delete Pega membuka konfirmasi hapus TREATY GROUP (salin-tempel) - adjuster tidak pernah terhapus. Pengganti
//     (keputusan work owner): flag nonaktif `IS_ACTIVE` (migrasi 870), nol hapus permanen.
//
// Tambahan keputusan work owner: Name wajib; nama ganda DITOLAK (tanpa beda huruf dan spasi tepi, juga terhadap baris
// nonaktif); hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/adjusterconsultant/backend/models"
	"nusantarare/modul/adjusterconsultant/backend/repository"
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
	ErrTidakAda = errors.New("services: adjuster consultant tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumDimigrasi - migrasi 870 belum dijalankan.
	ErrBelumDimigrasi = repository.ErrBelumDimigrasi
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
	Daftar(ctx context.Context, kata, bendera string) ([]models.Adjuster, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Adjuster, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Adjuster, error)
	IDBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, a models.Adjuster, username string) error
	Ubah(ctx context.Context, tx *dbTx, a models.Adjuster, username string) error
	SetelAktif(ctx context.Context, tx *dbTx, id string, aktif bool, username string) error
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
		return dilarang("Your access to the Adjuster Consultant menu is View only")
	}
	return nil
}
