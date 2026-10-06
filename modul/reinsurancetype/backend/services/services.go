// Package services memuat aturan modul Reinsurance Type (keputusan work owner 05-10-2026).
//
// Tabel warisan `POOLDATA.REINSURANCETYPE` (130 baris DEV) tidak punya layar master Pega di korpus: Pega membacanya
// (`BrowseReinsuranceType_RD` - saringan ID, Note, Flag, Type; `GetReinsuranceTypeBYName_SQL` mencari ID lewat NAMA)
// dan menulisnya lewat prosedur `PEGA_REINSURANCETYPE` (tidak dipanggil di sini). Aturan modul ini:
//   - Add dan Edit; TANPA hapus - nonaktif = Flag `inactive`.
//   - Name (NOTE) wajib, huruf besar, tidak kembar; SOA Name huruf besar.
//   - Type 1 Own Retention / 2 Treaty Out / 3 Facultative / 4 Treaty In; Flag `active` / `inactive`; Group Type OR / QS /
//     RI / SPL atau kosong; nilai warisan lain (Flag `1` - dipakai Contract Retro Life - dan kosong) dibiarkan selama
//     tidak diubah.
//   - Code angka, kosong = `00`; No Urut angka atau kosong.
//   - ID baru = `1` + `M_REINSURANCETYPE_SEQ` 4 digit; TGLUPDATE format Pega GMT; USERID = akun.
//   - Hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/reinsurancetype/backend/models"
	"nusantarare/modul/reinsurancetype/backend/repository"
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
	ErrTidakAda = errors.New("services: reinsurance type tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumAda - tabel atau sequence warisan tidak ada di skema ini.
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
	Daftar(ctx context.Context, kata, tipe, flag string) ([]models.Jenis, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Jenis, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Jenis, error)
	IDBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, j models.Jenis) error
	Ubah(ctx context.Context, tx *dbTx, j models.Jenis) error
}

// Layanan adalah aturan modul di atas Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
	jam    func() time.Time
}

// BaruLayanan membuat layanan.
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx, jam: time.Now} }

// DenganJam mengganti jam (uji).
func (l *Layanan) DenganJam(jam func() time.Time) *Layanan {
	l.jam = jam
	return l
}

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
		return dilarang("Your access to the Reinsurance Type menu is View only")
	}
	return nil
}
