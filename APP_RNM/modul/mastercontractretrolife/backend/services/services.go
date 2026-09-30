// Package services memuat aturan dagang modul Master Contract Retro Life.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers, dan tidak pernah mengimpor modul lain - yang
// bersama datang dari `inti/`.
//
// ⛔ Procedure `INSERT*_LIFE` TIDAK dipanggil (keputusan o, R4): logikanya
// ditiru di sini dan di repository - upsert dikunci `ID`, `ID` baru dari
// sequence, `TGLUPDATE = SYSDATE`, `USERID` pelaku. Transaksi dibuka di sini
// (`inti.Dasar.DalamTransaksi`), satu per permintaan, nol COMMIT di SQL.
package services

import (
	"context"
	"log"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/repository"
)

// Service adalah akar layanan Master Contract Retro Life.
type Service struct {
	*inti.Dasar
}

// DariDasar membuat Service di atas akar bersama yang disetel `cmd/api`.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// New membuat Service; db boleh nil (proses tanpa Oracle - rute menjawab 503).
func New(d *db.DB) *Service { return &Service{Dasar: inti.NewDasar(d)} }

// PunyaDatabase aman dipanggil pada nil (pola Treaty Contract Out).
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
	catat  func(string)
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan
// transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi, catat func(string)) *Layanan {
	if catat == nil {
		catat = func(string) {}
	}
	return &Layanan{gudang: g, tx: tx, catat: catat}
}

// LayananOracle menyusun Layanan di atas Oracle - satu-satunya penyusun yang
// dipakai handlers (handlers tidak mengimpor repository).
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi, func(baris string) { log.Print(baris) })
}
