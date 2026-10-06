// Package services memuat aturan modul Treaty Exchange Yearly (keputusan work owner 05-10-2026).
//
// Tabel warisan `POOLDATA.TREATYEXCHANGEYEARLY` (140 baris DEV) tidak punya layar master Pega di korpus: Pega
// membacanya (FacIn `GetCurrencyToIDR_SQL` / `GetKursLimitSpreading_SQL` membandingkan STARTDATE/ENDDATE SEBAGAI TEKS;
// Treaty In `BrowseTreatyExchangeYearly_RD` mengambil baris pertama per tahun, mata uang, `Quarter = "0"`; Aggregate)
// dan menulisnya lewat prosedur `PEGA_TREATYEXCHANGE` (tidak dipanggil di sini). Aturan modul ini:
//   - Add dan Edit; TANPA hapus.
//   - Treaty Year (4 angka), Currency (view CURRENCY; IDCURRENCY + kode di CURRENCY), Start Date, End Date, To IDR
//     wajib; To USD boleh kosong; angka tanpa pemisah ribuan. Quarter 0 (tahunan) - 4, tampil di form.
//   - Tanggal ditulis format Pega yang benar `YYYYMMDDT000000.000 GMT` ("pake format seharusnya"); tanggal yang tidak
//     diubah saat Edit dibiarkan apa adanya.
//   - Treaty Year + Currency + Quarter tidak boleh kembar (baris kembar lama dibiarkan).
//   - ID baru = situs aktif + `TREATYEXCHANGE_SEQ` 4 digit; `M_TREATYEXCHANGE` (JSON) tidak disentuh. Baris dikenali
//     lewat ROWID karena ID warisan tidak unik.
//   - USERID = akun; DATEIU = waktu simpan; DATEIN = waktu Add (format Pega GMT).
//   - Hak menu Full / View only.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyexchangeyearly/backend/models"
	"nusantarare/modul/treatyexchangeyearly/backend/repository"
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
	// ErrTidakAda - baris tidak ada.
	ErrTidakAda = errors.New("services: kurs tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumAda - tabel, view, atau sequence warisan tidak ada di skema ini.
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
	Daftar(ctx context.Context, kata, tahun string) ([]models.Kurs, error)
	AmbilKunci(ctx context.Context, tx *dbTx, kunci string) (models.Kurs, error)
	AmbilID(ctx context.Context, tx *dbTx, id string) (models.Kurs, error)
	Kembar(ctx context.Context, tx *dbTx, tahun, idCurrency, quarter, kecualiKunci string) ([]string, error)
	DaftarMataUang(ctx context.Context) ([]models.MataUang, error)
	AmbilMataUang(ctx context.Context, tx *dbTx, id string) (models.MataUang, error)
	DaftarTahun(ctx context.Context) ([]string, error)
	IDBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, k models.Kurs) error
	Ubah(ctx context.Context, tx *dbTx, k models.Kurs) error
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
		return dilarang("Your access to the Treaty Exchange Yearly menu is View only")
	}
	return nil
}
