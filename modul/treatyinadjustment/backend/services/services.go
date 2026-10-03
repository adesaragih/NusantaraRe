// Package services memegang aturan modul Treaty In Adjustment
// (`treatyinadjustment`).
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/repository"
)

// ErrKontrakTidakAda - kontrak yang diminta tidak ada.
var ErrKontrakTidakAda = errors.New("kontrak tidak ada")

// ErrIDTidakSah - pengenal kontrak bukan bilangan positif.
var ErrIDTidakSah = errors.New("pengenal kontrak tidak sah")

// Service membawa akar layanan modul ini.
type Service struct{ *inti.Dasar }

// DariDasar merakit Service di atas akar yang perakit sediakan.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// New membuat Service; db boleh nil (proses tanpa Oracle - rute menjawab 503).
func New(d *db.DB) *Service { return &Service{Dasar: inti.NewDasar(d)} }

// PunyaDatabase aman dipanggil pada nil.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// Gudang adalah seluruh sentuhan basis data yang Layanan butuhkan.
type Gudang interface {
	DaftarKontrak(ctx context.Context) ([]models.Kontrak, error)
	RantaiVersi(ctx context.Context, idKontrak int64) ([]models.Versi, error)
}

// Layanan memegang aturan modul ini di atas satu Gudang.
type Layanan struct{ gudang Gudang }

// LayananOracle merakit Layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan { return &Layanan{gudang: repository.Baru(s.DB())} }

// LayananDengan merakit Layanan di atas Gudang mana pun - dipakai uji.
func LayananDengan(g Gudang) *Layanan { return &Layanan{gudang: g} }

// DaftarKontrak membaca kepala seluruh kontrak.
func (l *Layanan) DaftarKontrak(ctx context.Context, p inti.Pelaku) ([]models.Kontrak, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.DaftarKontrak(ctx)
}

// RantaiVersi membaca rantai versi satu kontrak.
//
// Kontrak yang tidak ada menghasilkan ErrKontrakTidakAda, BUKAN daftar kosong.
// Keduanya berbeda arti: rantai kosong tidak mungkin - sebuah kontrak selalu
// punya sekurangnya versi pertamanya (tiket 14) - sehingga daftar kosong yang
// dikembalikan apa adanya akan terbaca di layar sebagai "kontrak ini tidak
// punya versi", padahal kontraknya sendiri yang tidak ada.
func (l *Layanan) RantaiVersi(ctx context.Context, p inti.Pelaku, idKontrak int64) ([]models.Versi, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if idKontrak <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrIDTidakSah, idKontrak)
	}
	versi, err := l.gudang.RantaiVersi(ctx, idKontrak)
	if err != nil {
		return nil, err
	}
	if len(versi) == 0 {
		return nil, fmt.Errorf("%w: %d", ErrKontrakTidakAda, idKontrak)
	}
	return versi, nil
}

// Pesan mengembalikan kalimat yang layak sampai ke layar.
func Pesan(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
