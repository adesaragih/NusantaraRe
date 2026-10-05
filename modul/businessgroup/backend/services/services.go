// Package services memuat aturan modul Business Group (keputusan work owner 05-10-2026).
//
// Tabel warisan `POOLDATA.BUSINESSGROUP` (45 baris DEV) tidak punya layar master Pega di korpus: Pega membacanya lewat
// `BrowseBusinessGroup_RD` (kelas `ASM-FW-GISFW-Int-BUSINESSGROUP`, saringan `.Note NotEndsWith "SYARIAH"`). Prosedur
// `UPSERT_BUSINESSGROUP` tidak dipanggil. Aturan modul ini:
//   - Add dan Edit; TANPA hapus ("biarkan saja") - grup bisnis dipakai BUSINESS dan TREATYGROUP.COAID.
//   - Grup berakhiran SYARIAH tidak tampil dan tidak dapat dibuka / diubah; nama baru berakhiran SYARIAH ditolak
//     ("tidak ada syariah").
//   - Treaty Group wajib (`TOPID` = `TREATYGROUP.ID`); TREATYNAME = salinan namanya. Edit yang tidak mengganti Treaty
//     Group membiarkan salinannya (juga TOPID lama yang tidak ada lagi di TREATYGROUP).
//   - Name (NOTE) wajib, huruf besar, tidak kembar; Alias Name huruf besar, kosong = Name.
//   - ID baru = situs aktif + `BUSINESSGROUP_SEQ` 4 digit; nomor yang sudah dipakai dilompati.
//   - Salinan nama di tabel lain (BUSINESS.BUSINESSGROUPNAME) TIDAK diubah ("jgn ada ubah data").
//   - Hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/businessgroup/backend/models"
	"nusantarare/modul/businessgroup/backend/repository"
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
	ErrTidakAda = errors.New("services: business group tidak ada")
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
	Daftar(ctx context.Context, kata, topID string) ([]models.BisnisGrup, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.BisnisGrup, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]string, error)
	DaftarTreaty(ctx context.Context) ([]models.TreatyGroup, error)
	AmbilTreaty(ctx context.Context, tx *dbTx, id string) (models.TreatyGroup, error)
	IDBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, b models.BisnisGrup) error
	Ubah(ctx context.Context, tx *dbTx, b models.BisnisGrup) error
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
		return dilarang("Your access to the Business Group menu is View only")
	}
	return nil
}
