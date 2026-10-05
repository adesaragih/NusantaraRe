// Package services memuat aturan modul Treaty Description (keputusan work owner 05-10-2026).
//
// Tabel warisan `POOLDATA.TREATYDESC` (13 baris DEV) = master jenis klausul treaty. Pega membacanya lewat
// `BrowseTreatyDesc_RD` (grid "For Non XOL" / "For XOL" Treaty Contract Out) dan menulisnya lewat prosedur
// `PEGA_TREATYDESC` (tidak dipanggil di sini); layar input Pega-nya tidak ada di korpus. Aturan modul ini:
//   - Add dan Edit; TANPA hapus - ID dipakai PROPORTIONALARRG dan rule Pega klaim (`TREATYDESCID='10001'`). Penonaktifan
//     lewat Status.
//   - ID baru = '1' + TREATY_DESCRIPTION_SEQ 4 digit (seperti PEGA_TREATYDESC); nomor yang sudah dipakai dilompati.
//     ID tidak pernah diubah.
//   - Description Name wajib, huruf besar, maks. 100 byte, tidak kembar (tanpa beda huruf, juga dengan baris nonaktif).
//   - Type: Non XOL ("0") / XOL ("1"); Status: Active ("1") / Inactive ("0"). STATUSAKTIF NULL (baris lama Pega)
//     dibaca Active; Edit menulisnya "1"/"0".
//   - Nama boleh diubah (keputusan work owner 05-10-2026), termasuk 10001; salinan TREATYDESCNAME di PROPORTIONALARRG
//     TIDAK ikut diubah.
//   - HANYA TREATYDESC yang dibaca dan ditulis ("hanya baca dari tabel yang saya kasih, jangan ber-experiment").
//   - Hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatydescription/backend/models"
	"nusantarare/modul/treatydescription/backend/repository"
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
	ErrTidakAda = errors.New("services: treaty description tidak ada")
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
	Daftar(ctx context.Context, kata, xol, status string) ([]models.Desc, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Desc, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]string, error)
	IDBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, r models.Desc) error
	Ubah(ctx context.Context, tx *dbTx, r models.Desc) error
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
		return dilarang("Your access to the Treaty Description menu is View only")
	}
	return nil
}
