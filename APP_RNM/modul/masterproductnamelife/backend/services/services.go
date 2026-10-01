// Package services memuat aturan dagang modul Master Product Name Life.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers, dan tidak pernah mengimpor modul lain - yang
// bersama datang dari `inti/`.
//
// ⛔ Prosedur `PEGA_M_PRODUCT_LIFE` dan `PEGA_M_PRODUCT_INWARD_LIFE` TIDAK
// dipanggil (brief bab 1): logikanya ditiru di sini dan di repository -
// upsert dikunci `ID`, `ID` baru dari sequence, kedua tabel dalam SATU
// transaksi (P4). Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`),
// satu per permintaan, nol COMMIT di SQL.
package services

import (
	"context"
	"log"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// Service adalah akar layanan Master Product Name Life.
type Service struct {
	*inti.Dasar
}

// DariDasar membuat Service di atas akar bersama yang disetel `cmd/api`.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// New membuat Service; db boleh nil (proses tanpa Oracle - rute menjawab 503).
func New(d *db.DB) *Service { return &Service{Dasar: inti.NewDasar(d)} }

// PunyaDatabase aman dipanggil pada nil.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
	catat  func(string)
	// jam - `@CurrentDateTime()` (tanggal baris komentar).
	jam func() time.Time
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
	salinan := *l
	salinan.jam = jam
	return &salinan
}

// LayananOracle menyusun Layanan di atas Oracle - satu-satunya penyusun yang
// dipakai handlers (handlers tidak mengimpor repository).
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi, func(baris string) { log.Print(baris) })
}
