// Package services memuat aturan modul R/I Risk (`ririsklife`; panduan XML section `InboxSummaryRIRisk` dan
// `InboxRIRisk`; keputusan work owner 08-10-2026 K1-K4, MODUL.md). Pola ricommlife; PERBEDAAN dari XML ditandai.
//
//   - Ringkasan: daftar bersaring ID / R/I RISK NAME, 50 per halaman (`BrowseRIRiskSummary` b7127, `pyPageSize` b10030),
//     ID menaik; Delete (`DeleteSummaryDetail` b9648) menghapus ringkasan BESERTA rinciannya dalam satu transaksi.
//     Add dan Edit lewat Save b1796 (`AddToListSummary_Act` b1819; Edit b8589 `EditListSummary_DT` b8616). Di XML
//     Save / Cancel ada di wadah TERSEMBUNYI (`pyContainerVisibleWhen` `1=2` b1511, CONDITION b1546) - DITAMPILKAN
//     atas keputusan work owner 08-10-2026 (RALAT R1 MODUL.md), sama dengan ricommlife.
//   - R/I RISK NAME (`.USEDBY`, wajib b1274) dipangkas dan TIDAK BOLEH KEMBAR tanpa beda huruf (indeks
//     IX_RIRISK_LIFE_SUMMARY_NAMA, migrasi inti 937); upload mencocokkan ringkasan lewat nama.
//   - ID baru: ringkasan = site `M_SITE_DATABASE` || LPAD(M_RIRISK_LIFE_SUMMARY_SEQ, 6, '0'); rincian = '1' ||
//     LPAD(M_RIRISK_LIFE_SEQ, 5, '0') - DUA rumus prosedur warisan (models.BentukID / BentukIDRincian).
//   - OPERATORID = akun login; MODIFIEDDATE = format Pega `yyyyMMdd'T'HHmmss.SSS 'GMT'` (seperti ricommlife).
//   - R/I RISK DETAIL (`setIDUsedBy_Act` b8886 + harness `InboxRIRisk` b8923): tabel `RIRISK_LIFE`, tambah dan ubah
//     satu baris (CONTRACT, YEAR, MONTH, RISK) - detail.go.
//   - Upload CSV / View Upload / Simpan Upload (b3156, b3676, b4690): kolom USEDBY, CONTRACT, YEAR, MONTH, RISK - unggah.go.
//   - Transaksi milik Go: satu transaksi per simpan, prosedur PEGA_* tidak dipanggil (dan INVALID sesudah 937/940).
//   - Hak menu Full / View only: View only tanpa Save, Edit, Delete, Upload, dan form detail.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/ririsklife/backend/models"
	"nusantarare/modul/ririsklife/backend/repository"
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
	ErrTidakAda = errors.New("services: ringkasan R/I risk tidak ada")
	// ErrRincianTidakAda - baris rincian tidak ada di ringkasan itu.
	ErrRincianTidakAda = errors.New("services: baris R/I risk detail tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini.
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
	UbahNamaRincian(ctx context.Context, tx *dbTx, idUsedBy, nama string) (int, error)
	HapusRincian(ctx context.Context, tx *dbTx, idUsedBy string) (int, error)
	JumlahRincian(ctx context.Context, tx *dbTx, idUsedBy string) (int, error)
	DaftarRincian(ctx context.Context, idUsedBy string, halaman int) ([]models.Rincian, int, error)
	RincianDari(ctx context.Context, tx *dbTx, ids []string) ([]models.Rincian, error)
	SisipRincian(ctx context.Context, tx *dbTx, k models.Rincian) error
	UbahRincian(ctx context.Context, tx *dbTx, k models.Rincian) error
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
		return dilarang("Your access to the R/I Risk menu is View only")
	}
	return nil
}

// cobaIDMaksimum - ID yang ternyata sudah terpakai dilewati paling banyak sekian kali per ID.
const cobaIDMaksimum = 100

// pemberiID - ID baru di dalam SATU transaksi: nomor dari sequence warisan; ringkasan = site (dibaca sekali) ||
// LPAD 6 (models.BentukID), rincian = '1' || LPAD 5 (models.BentukIDRincian, tanpa site); ID yang sudah terpakai
// (mis. ditulis prosedur Pega) dilewati.
type pemberiID struct {
	l         *Layanan
	ringkasan bool
	situs     string
}

func (p *pemberiID) baru(ctx context.Context, tx *dbTx) (string, error) {
	if p.ringkasan && p.situs == "" {
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
		var id string
		if p.ringkasan {
			id, err = models.BentukID(p.situs, nomor)
		} else {
			id, err = models.BentukIDRincian(nomor)
		}
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
