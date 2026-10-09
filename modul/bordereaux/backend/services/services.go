// Package services memuat aturan modul Bordereaux (keputusan work owner 04-10-2026; XML Pega folder korpus
// `Bordereaux` patokannya, bug Pega diperbaiki di Go, format ID mengikuti XML, tanpa JSON).
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/bordereaux/backend/models"
	"nusantarare/modul/bordereaux/backend/repository"
)

// Service membawa akar bersama dan garam token penyimpanan (`config.StorageTokenSalt`).
type Service struct {
	*inti.Dasar
	garam string
}

// DariDasar membungkus akar bersama.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// DenganGaram menyetel garam token penyimpanan lampiran (`STORAGE_TOKEN_SALT`).
func (s *Service) DenganGaram(garam string) *Service {
	salin := *s
	salin.garam = garam
	return &salin
}

// PunyaDatabase - Oracle terpasang?
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// LayananOracle merakit layanan di atas Oracle; lampiran disimpan lewat penyimpanan bersama `inti/backend/penyimpanan`.
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi).DenganPenyimpanan(penyimpanan.Oracle(s.Dasar, s.garam))
}

var _ Gudang = (*repository.Gudang)(nil)

type dbTx = db.Tx

var (
	// ErrTidakAda - berkas tidak ada.
	ErrTidakAda = errors.New("services: berkas bordereaux tidak ada")
	// ErrMasukanTidakSah - permintaan ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak atas aksi ini; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrSudahDiproses - berkas sudah berpindah posisi/status sejak dibaca.
	ErrSudahDiproses = errors.New("services: berkas sudah diproses")
	// ErrBelumDimigrasi - migrasi 890 belum dijalankan.
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

// Gudang - semua yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, f models.Filter, offset, ukuran int) ([]models.Header, int, error)
	AmbilHeader(ctx context.Context, tx *dbTx, id string) (models.Header, error)
	AdaBerkas(ctx context.Context, tx *dbTx, id string) (bool, error)
	SisipHeader(ctx context.Context, tx *dbTx, h models.Header) error
	UbahHeader(ctx context.Context, tx *dbTx, h models.Header) error
	UbahStatus(ctx context.Context, tx *dbTx, id, posisiLama, statusLama, posisiBaru, statusBaru string) (bool, error)
	Detail(ctx context.Context, tx *dbTx, k models.KombinasiBdx, id string) ([]models.Baris, error)
	HapusDetail(ctx context.Context, tx *dbTx, k models.KombinasiBdx, id string) (int64, error)
	SisipDetail(ctx context.Context, tx *dbTx, k models.KombinasiBdx, id, idDetail string, b models.Baris) error
	HapusBerkas(ctx context.Context, tx *dbTx, id string) error
	Riwayat(ctx context.Context, id string) ([]models.Riwayat, error)
	SisipRiwayat(ctx context.Context, tx *dbTx, id, pic string, setuju bool, komentar string) error
	CariCedant(ctx context.Context, kata string, batas int) ([]models.Cedant, error)
	CariTreaty(ctx context.Context, cedingID string, batas int) ([]models.MasterTreaty, error)
	AmbilTreaty(ctx context.Context, id string) (models.MasterTreaty, error)
	Chart(ctx context.Context) ([]models.IrisanChart, error)
	// Copy Old Data.
	JSONLama(ctx context.Context, tx *dbTx, id string) ([]models.JSONLama, error)
	CacahDetail(ctx context.Context, k models.KombinasiBdx) (map[string]int, error)
	CacahRiwayat(ctx context.Context) (map[string]int, error)
	SisipRiwayatLama(ctx context.Context, tx *dbTx, id, tanggal, pic string, setuju bool, komentar string) error
	// Lampiran.
	KategoriLampiran(ctx context.Context, bdxID string) ([]models.KategoriLampiran, error)
	NamaKategoriLampiran(ctx context.Context, id string) (string, bool, error)
	DaftarLampiran(ctx context.Context, bdxID, kategoriID string) ([]models.Lampiran, error)
	AmbilLampiran(ctx context.Context, bdxID, kategoriID, id string) (models.Lampiran, bool, error)
	SisipLampiran(ctx context.Context, tx *dbTx, a models.Lampiran) error
	HapusLampiran(ctx context.Context, tx *dbTx, bdxID, id string) error
}

// Layanan adalah aturan modul di atas Gudang.
type Layanan struct {
	gudang   Gudang
	tx       Transaksi
	sekarang func() time.Time
	// berkas - penyimpanan lampiran (`lampiran.go`); nil = gagal terang.
	berkas PenyimpananBerkas
}

// BaruLayanan membuat layanan.
func BaruLayanan(g Gudang, tx Transaksi) *Layanan {
	return &Layanan{gudang: g, tx: tx, sekarang: sekarangJakarta}
}

// DenganJam mengganti jam (uji).
func (l *Layanan) DenganJam(f func() time.Time) *Layanan {
	salin := *l
	salin.sekarang = f
	return &salin
}

// Aktor - pelaku permintaan: akun, workbasket, apakah superadmin (pemegang menu Kelola User = pengganti `IT Developer`
// Pega), dan apakah menu Bordereaux-nya ber-hak PENUH (`M_LOGIN_GO_MENU.HAK`, keputusan work owner 04-10-2026).
type Aktor struct {
	AkunID     string
	Peran      []string
	Superadmin bool
	// Penuh - menu Bordereaux bukan View only: Input Data, Edit, Delete, dan Submit pembuat.
	Penuh bool
}

func (a Aktor) punya(peran string) bool {
	for _, p := range a.Peran {
		if p == peran {
			return true
		}
	}
	return false
}

// BolehBuat - tombol Add/Input Data: menu Bordereaux ber-hak PENUH (pengganti workbasket ReasBordereauxAdmin,
// keputusan work owner 04-10-2026). Superadmin pun harus PENUH (keputusan work owner 05-10-2026).
func (a Aktor) BolehBuat() bool {
	return a.AkunID != "" && a.Penuh
}

// Hak - aksi yang boleh dilakukan aktor atas satu berkas.
type Hak struct {
	Ubah     bool `json:"ubah"`
	Hapus    bool `json:"hapus"`
	Submit   bool `json:"submit"`
	Putuskan bool `json:"putuskan"`
	// Lampiran - Upload File / Delete lampiran (`BolehLampiran`); diisi saat berkas dibuka.
	Lampiran bool `json:"lampiran"`
}

// HakAtas menilai hak aktor atas satu berkas (`PortalBordereaux_Sec` kolom Edit/Delete, `InputBordereaux` tab Submit):
//   - Edit/Delete/Submit pembuat: berkas di tangan pembuatnya (POSITION = USER_INPUT), belum Resolve-Complete, menu
//     aktornya ber-hak PENUH, dan aktornya pembuat itu atau superadmin (mengambil alih berkas yang pembuatnya tanpa
//     akun). Superadmin ber-hak View only tidak mengambil alih (keputusan work owner 05-10-2026);
//   - Putuskan: POSITION Checker dan aktor pemegang ReasBordereauxChecker yang BUKAN pembuat berkasnya, atau POSITION
//     Supervisor dan aktor pemegang ReasBordereauxSupervisor.
func HakAtas(a Aktor, h models.Header) Hak {
	if a.AkunID == "" || h.StatusAksep == models.StatusResolveComplete {
		return Hak{}
	}
	var hak Hak
	diTanganPembuat := h.Position != "" && h.Position == h.UserInput
	if diTanganPembuat && a.Penuh && (a.AkunID == h.UserInput || a.Superadmin) {
		hak.Ubah, hak.Hapus, hak.Submit = true, true, true
	}
	switch h.Position {
	case models.PosisiChecker:
		hak.Putuskan = a.punya(models.PeranChecker) && a.AkunID != h.UserInput
	case models.PosisiSupervisor:
		hak.Putuskan = a.punya(models.PeranSupervisor)
	}
	return hak
}
