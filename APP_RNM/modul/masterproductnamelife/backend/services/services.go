// Package services memuat aturan dagang modul Master Product Name Life.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers, dan tidak pernah mengimpor modul lain - yang
// bersama datang dari `inti/`.
//
// ⛔ Prosedur `PEGA_M_PRODUCT_LIFE` dan `PEGA_M_PRODUCT_INWARD_LIFE` TIDAK
// dipanggil (brief bab 1): logikanya ditiru di sini dan di repository -
// upsert dikunci `ID`, `ID` baru dari sequence, baris induk flat dan ketujuh anaknya dalam SATU
// transaksi (P4). Transaksi dibuka di sini (`inti.Dasar.DalamTransaksi`),
// satu per permintaan, nol COMMIT di SQL.
package services

import (
	"context"
	"log"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/layanan"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// Service adalah akar layanan Master Product Name Life.
type Service struct {
	*inti.Dasar
	// storage - penyimpanan nyata terpasang (`PELAKSANA_STORAGE=nyata`); nil = stub lokal.
	storage *storageNyata
}

// storageNyata - bahan token penyimpanan. ⛔ Garam tidak pernah dicetak, dicatat, atau masuk pesan galat.
type storageNyata struct {
	garam string
}

// DenganPenyimpananNyata memasang penyimpanan nyata (`mpnl_storage.go`) bila `nyata` - dipanggil sekali dari
// `modul.go` dengan `PELAKSANA_STORAGE == nyata` dan `STORAGE_TOKEN_SALT` (`inti/backend/config`, yang menolak
// `nyata` tanpa garam saat memuat). Selain itu stub lokal.
func (s *Service) DenganPenyimpananNyata(nyata bool, garam string) *Service {
	salin := *s
	salin.storage = nil
	if nyata {
		salin.storage = &storageNyata{garam: garam}
	}
	return &salin
}

// penyimpanan - penyimpanan berkas lampiran yang dipasang: alamat `M_LINK_SERVICE` dan token `GCP_IMAGE` dibaca
// saat jalan (tanpa Oracle gagal terang).
func (s *Service) penyimpanan() PenyimpananBerkas {
	if s.storage == nil {
		return PenyimpananLokal(s.UnggahanDir())
	}
	alamat := func(ctx context.Context, k layanan.KunciLayanan) (string, error) {
		if !s.PunyaDatabase() {
			return "", db.ErrTanpaOracle
		}
		return layanan.AlamatLayanan(ctx, layanan.ResolverLinkServiceOracle(s), k)
	}
	token := func(ctx context.Context, app string) (string, error) {
		if !s.PunyaDatabase() {
			return "", db.ErrTanpaOracle
		}
		var tok string
		err := s.DalamTransaksi(ctx, func(tx *db.Tx) error {
			var err error
			tok, err = tokenStorage(ctx, tx, layanan.NewPenyimpanToken(s.DB()), s.storage.garam, app, time.Now())
			return err
		})
		return tok, err
	}
	return PenyimpananGoogle(s.UnggahanDir(), alamat, token, nil, nil)
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
	// jam - `@CurrentDateTime()` (tanggal baris komentar, waktu lampiran).
	jam func() time.Time
	// berkas - penyimpanan berkas lampiran (stub lokal atau nyata; tiruan di uji).
	berkas PenyimpananBerkas
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan
// transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi, catat func(string)) *Layanan {
	if catat == nil {
		catat = func(string) {}
	}
	return &Layanan{gudang: g, tx: tx, catat: catat, jam: time.Now, berkas: penyimpananBelumDisetel{}}
}

// DenganPenyimpanan mengganti penyimpanan berkas (uji, atau stub lokal).
func (l *Layanan) DenganPenyimpanan(b PenyimpananBerkas) *Layanan {
	salinan := *l
	salinan.berkas = b
	return &salinan
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
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi, func(baris string) { log.Print(baris) }).
		DenganPenyimpanan(s.penyimpanan())
}
