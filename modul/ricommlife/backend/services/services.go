// Package services memuat aturan modul R/I Comm Life (perintah work owner 06-10-2026; panduan XML section
// `InboxSummaryRIComm` dan `InboxRIComm`; keputusan work owner 06-10-2026 butir 1-9, MODUL.md).
//
//   - Ringkasan (CRUD): daftar bersaring ID / R/I COMM NAME, 50 per halaman (`BrowseRICommSummary` b10081), ID menaik;
//     Add dan Edit lewat Save (`AddToListSummary_Act` b1823); Delete (`DeleteSummaryDetail` b9704) menghapus ringkasan
//     BESERTA rinciannya dalam satu transaksi.
//   - R/I COMM NAME (`.USEDBY`, wajib b1289) dipangkas dan TIDAK BOLEH KEMBAR tanpa beda huruf: upload mencocokkan
//     ringkasan lewat nama.
//   - ID baru (ringkasan DAN rincian) = site `M_SITE_DATABASE` || LPAD(sequence warisan, 6, '0') (butir 4).
//   - OPERATORID = akun login; MODIFIEDDATE = format Pega `yyyyMMdd'T'HHmmss.SSS 'GMT'` (butir 6).
//   - R/I COMM DETAIL (`setIDUsedBy_Act` b9059 + harness `InboxRIComm` b9096): tabel flat `RICOMM_LIFE`, tambah dan
//     ubah satu baris - detail.go.
//   - Upload CSV / View Upload / Simpan Upload (b3189, b3706, b4771): gaya riratelife (butir 7) - unggah.go.
//   - Transaksi milik Go: satu transaksi per simpan, prosedur PEGA_* tidak dipanggil (butir 5).
//   - Hak menu Full / View only: View only tanpa Add, Edit, Delete, Upload (butir 9).
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/ricommlife/backend/models"
	"nusantarare/modul/ricommlife/backend/repository"
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
	ErrTidakAda = errors.New("services: ringkasan R/I comm tidak ada")
	// ErrKomisiTidakAda - baris rincian tidak ada di ringkasan itu.
	ErrKomisiTidakAda = errors.New("services: baris R/I comm detail tidak ada")
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
// sequence ringkasan, false = rincian.
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Ringkasan, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Ringkasan, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Ringkasan, error)
	Situs(ctx context.Context, tx *dbTx) (string, error)
	NomorBaru(ctx context.Context, tx *dbTx, ringkasan bool) (string, error)
	AdaID(ctx context.Context, tx *dbTx, ringkasan bool, id string) (bool, error)
	SisipRingkasan(ctx context.Context, tx *dbTx, r models.Ringkasan) error
	UbahRingkasan(ctx context.Context, tx *dbTx, r models.Ringkasan) error
	HapusRingkasan(ctx context.Context, tx *dbTx, id string) error
	UbahNamaKomisi(ctx context.Context, tx *dbTx, idUsedBy, nama string) (int, error)
	HapusKomisi(ctx context.Context, tx *dbTx, idUsedBy string) (int, error)
	JumlahKomisi(ctx context.Context, tx *dbTx, idUsedBy string) (int, error)
	DaftarKomisi(ctx context.Context, idUsedBy string, halaman int) ([]models.Komisi, int, error)
	KomisiDari(ctx context.Context, tx *dbTx, ids []string) ([]models.Komisi, error)
	SisipKomisi(ctx context.Context, tx *dbTx, k models.Komisi) error
	UbahKomisi(ctx context.Context, tx *dbTx, k models.Komisi) error
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
		return dilarang("Your access to the R/I Comm Life menu is View only")
	}
	return nil
}

// cobaIDMaksimum - ID yang ternyata sudah terpakai dilewati paling banyak sekian kali per ID.
const cobaIDMaksimum = 100

// pemberiID - ID baru di dalam SATU transaksi: site dibaca sekali, nomor dari sequence warisan, dibentuk
// models.BentukID; ID yang sudah terpakai (mis. ditulis prosedur Pega) dilewati.
type pemberiID struct {
	l         *Layanan
	ringkasan bool
	situs     string
}

func (p *pemberiID) baru(ctx context.Context, tx *dbTx) (string, error) {
	if p.situs == "" {
		s, err := p.l.gudang.Situs(ctx, tx)
		if err != nil {
			return "", err
		}
		p.situs = s
	}
	for i := 0; i < cobaIDMaksimum; i++ {
		nomor, err := p.l.gudang.NomorBaru(ctx, tx, p.ringkasan)
		if err != nil {
			return "", err
		}
		id, err := models.BentukID(p.situs, nomor)
		if err != nil {
			return "", err
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
