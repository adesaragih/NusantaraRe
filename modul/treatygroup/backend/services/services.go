// Package services memuat aturan modul Treaty Group (keputusan work owner 05-10-2026).
//
// Tabel warisan `POOLDATA.TREATYGROUP` (33 baris DEV) tidak punya layar master Pega di korpus: Pega membacanya lewat
// `BrowseTreatyGroup_RD` (kelas `ASM-FW-GISFW-Int-TREATYGROUP`) dan menulisnya lewat prosedur `PEGA_TREATYGROUP` (tidak
// dipanggil di sini). Aturan modul ini:
//   - Add dan Edit; TANPA hapus ("biarkan saja") - grup dipakai banyak tabel lain.
//   - OJK Business wajib, dipilih dari `TREATYGROUPOJK`; OJKBUSINESSID/NAME/NAMEIDN dan ORDERNO disalin dari baris OJK
//     itu (33 dari 33 baris DEV cocok). Edit yang tidak mengganti OJK membiarkan salinannya.
//   - Treaty Group Name wajib, huruf besar, tidak kembar; SOA Name boleh kosong, huruf besar.
//   - COAID (= `BUSINESSGROUP.ID`) opsi A: tidak diketik; ikut OJK yang dipilih ("ubah pas pilih OJK Business") -
//     COAID grup lain se-OJK yang paling sering. Edit yang tidak mengganti OJK membiarkannya.
//   - ID baru = situs aktif + `TREATYGROUP_SEQ` 4 digit; TGLUPDATE format Pega `YYYYMMDDTHHMMSS.mmm GMT`; USERID = akun.
//   - Salinan nama grup di tabel lain (BUSINESSGROUP.TREATYNAME, TREATYBUSINESS, TREATYINDETAIL, ...) TIDAK diubah
//     ("jgn ada ubah data").
//   - Hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatygroup/backend/models"
	"nusantarare/modul/treatygroup/backend/repository"
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
	ErrTidakAda = errors.New("services: treaty group tidak ada")
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
	Daftar(ctx context.Context, kata, ojk string) ([]models.Grup, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Grup, error)
	Anak(ctx context.Context, id string) ([]models.BisnisGrup, error)
	DaftarOjk(ctx context.Context) ([]models.Ojk, error)
	AmbilOjk(ctx context.Context, tx *dbTx, id string) (models.Ojk, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]string, error)
	CoaSeOjk(ctx context.Context, tx *dbTx, ojk string) (string, error)
	IDBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, r models.Grup) error
	Ubah(ctx context.Context, tx *dbTx, r models.Grup) error
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
		return dilarang("Your access to the Treaty Group menu is View only")
	}
	return nil
}
