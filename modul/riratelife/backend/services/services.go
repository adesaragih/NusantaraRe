// Package services memuat aturan modul R/I Rate Life (perintah work owner 05-10-2026; panduan XML section
// `InboxSummaryRIRate`, kelas `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`).
//
//   - Ringkasan (CRUD): daftar bersaring ID / R/I RATE NAME, 50 per halaman (`BrowseRateLifeSummary` b12645); Add dan
//     Edit lewat Save (`AddToListSummary_Act` b1833 - container Save/Cancel tersembunyi `1=2` b1529 di Pega, DITAMPILKAN
//     atas perintah CRUD work owner); Delete (`DeleteSummaryDetail` b12262) menghapus ringkasan BESERTA rate-nya.
//   - R/I RATE NAME (`.USEDBY`, wajib b1285) dipangkas dan TIDAK BOLEH KEMBAR tanpa beda huruf: upload mencocokkan
//     ringkasan lewat nama.
//   - ID baru dari sequence Oracle (K2) - nomor yang sudah terpakai dilewati (preseden `aggregate`).
//   - OPERATORID = akun login; MODIFIEDDATE = format Pega `YYYYMMDDTHHMMSS.mmm GMT` (ASUMSI A3).
//   - Rate Detail (`setIDUsedBy_Act` b11415 + harness `InboxRIRate` b11444, judul "Rate Detail" b11446): grid 20 per
//     halaman, tambah dan ubah satu baris (View Detail.xml, work owner 06-10-2026) - ratedetail.go.
//   - Upload CSV / View Upload / Simpan Upload (b2949, b3467, b4511): aturan dirancang dari label format b5305 (K3).
//   - Hak menu Full / View only: View only tanpa Add, Edit, Delete, Upload.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/riratelife/backend/models"
	"nusantarare/modul/riratelife/backend/repository"
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
	// ErrTidakAda - ID ringkasan tidak ada.
	ErrTidakAda = errors.New("services: ringkasan R/I rate tidak ada")
	// ErrRateTidakAda - baris rate tidak ada di ringkasan itu.
	ErrRateTidakAda = errors.New("services: baris rate tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumAda - tabel, view, atau sequence tidak ada di skema ini.
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

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji). `ringkasan` true = tabel /
// sequence ringkasan, false = rate.
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Ringkasan, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Ringkasan, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Ringkasan, error)
	IDBaru(ctx context.Context, tx *dbTx, ringkasan bool) (string, error)
	MaksID(ctx context.Context, tx *dbTx, ringkasan bool) (int, error)
	AdaID(ctx context.Context, tx *dbTx, ringkasan bool, id string) (bool, error)
	SisipRingkasan(ctx context.Context, tx *dbTx, r models.Ringkasan) error
	UbahRingkasan(ctx context.Context, tx *dbTx, r models.Ringkasan) error
	UbahNamaRate(ctx context.Context, tx *dbTx, idUsedBy, nama string) (int, error)
	HapusRingkasan(ctx context.Context, tx *dbTx, id string) error
	HapusRate(ctx context.Context, tx *dbTx, idUsedBy string) (int, error)
	JumlahRate(ctx context.Context, tx *dbTx, idUsedBy string) (int, error)
	DaftarRate(ctx context.Context, idUsedBy string, halaman int) ([]models.Rate, int, error)
	RateDari(ctx context.Context, tx *dbTx, ids []string) ([]models.Rate, error)
	SisipRate(ctx context.Context, tx *dbTx, r models.Rate) error
	UbahRate(ctx context.Context, tx *dbTx, r models.Rate) error
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
		return dilarang("Your access to the R/I Rate Life menu is View only")
	}
	return nil
}

// cobaIDMaksimum - nomor sequence yang sudah terpakai dilewati paling banyak sekian kali per ID.
const cobaIDMaksimum = 1000

// pemberiID - ID baru dari sequence; nomor di atas ID angka tertinggi tabel (dibaca sekali) pasti kosong, nomor di
// bawahnya diperiksa satu per satu dan dilewati bila terpakai.
type pemberiID struct {
	l         *Layanan
	ringkasan bool
	maks      int
	tahu      bool
}

func (p *pemberiID) baru(ctx context.Context, tx *dbTx) (string, error) {
	if !p.tahu {
		m, err := p.l.gudang.MaksID(ctx, tx, p.ringkasan)
		if err != nil {
			return "", err
		}
		p.maks, p.tahu = m, true
	}
	for i := 0; i < cobaIDMaksimum; i++ {
		id, err := p.l.gudang.IDBaru(ctx, tx, p.ringkasan)
		if err != nil {
			return "", err
		}
		if id == "" || len(id) > models.BatasID {
			return "", fmt.Errorf("services: nomor sequence %q di luar ID VARCHAR2(%d)", id, models.BatasID)
		}
		if n := angkaSaja(id); n > p.maks {
			return id, nil
		}
		ada, err := p.l.gudang.AdaID(ctx, tx, p.ringkasan, id)
		if err != nil {
			return "", err
		}
		if !ada {
			return id, nil
		}
	}
	return "", fmt.Errorf("services: tidak menemukan ID kosong sesudah %d nomor sequence", cobaIDMaksimum)
}

// angkaSaja - teks angka -> bilangan; selain angka = -1.
func angkaSaja(s string) int {
	if s == "" {
		return -1
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return -1
		}
		n = n*10 + int(c-'0')
	}
	return n
}
