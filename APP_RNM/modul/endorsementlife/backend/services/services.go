// Package services memuat aturan dagang modul Endorsement Life.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// mengimpor handlers dan tidak mengimpor modul lain - yang bersama datang dari
// `inti/`. Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`), satu per
// permintaan tulis, nol `COMMIT` di SQL.
//
// ⛔ Procedure Pega TIDAK dipanggil (`PEGA_M_LIFE_PREMIUM_SUMMARY`,
// `INSERTJSONPOLISLIFE`): logikanya ditiru di repository.
package services

import (
	"context"
	"log"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// Service adalah akar layanan Endorsement Life.
type Service struct {
	*inti.Dasar
}

// DariDasar membuat Service di atas akar bersama yang disetel `cmd/api`.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// New membuat Service; db boleh nil (proses tanpa Oracle - rute menjawab 503).
func New(d *db.DB) *Service { return &Service{Dasar: inti.NewDasar(d)} }

// PunyaDatabase aman dipanggil pada nil.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar != nil && s.Dasar.PunyaDatabase() }

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// Gudang adalah penyimpanan yang dibutuhkan layanan - Oracle
// (`repository.Gudang`) atau tiruan di memori (`tiruan.Gudang`, uji).
type Gudang interface {
	Inbox(ctx context.Context, halaman, ukuran int) ([]models.BarisInbox, int, error)
	AmbilKasus(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, error)
	DaftarPeserta(ctx context.Context, kasusID string, halaman, ukuran int) ([]models.Peserta, error)
	CacahPeserta(ctx context.Context, tx *db.Tx, kasusID string) (map[string]int, error)
	AdaRekap(ctx context.Context, tx *db.Tx, kasusID string) (bool, error)
	VersiBerjalan(ctx context.Context, tx *db.Tx, nomorPolis string, sebelumProdKe int) (models.Versi, bool, error)
}

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
	catat  func(string)
	jam    func() time.Time
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan
// transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi, catat func(string)) *Layanan {
	if catat == nil {
		catat = func(string) {}
	}
	return &Layanan{gudang: g, tx: tx, catat: catat, jam: time.Now}
}

// DenganJam mengganti jam layanan (uji).
func (l *Layanan) DenganJam(jam func() time.Time) *Layanan {
	l.jam = jam
	return l
}

// LayananOracle menyusun Layanan di atas Oracle - satu-satunya penyusun yang
// dipakai handlers (handlers tidak mengimpor repository).
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi, func(baris string) { log.Print(baris) })
}

// UkuranHalaman - baris per halaman grid (`pyGridPaginator`).
const UkuranHalaman = 20

// jepit merapikan nomor dan ukuran halaman.
func jepit(halaman, ukuran int) (int, int) {
	if halaman < 1 {
		halaman = 1
	}
	if ukuran < 1 || ukuran > 200 {
		ukuran = UkuranHalaman
	}
	return halaman, ukuran
}
