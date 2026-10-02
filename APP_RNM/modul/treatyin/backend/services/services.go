// Package services memegang aturan modul Treaty In (`treatyin`).
//
// Arah ketergantungan: handlers -> services -> repository.
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

// ErrHimpunanTidakAda - himpunan acuan yang diminta bukan salah satu dari enam.
var ErrHimpunanTidakAda = errors.New("himpunan acuan tidak dikenal")

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
	DaftarAcuan(ctx context.Context, h models.Himpunan) ([]models.Acuan, error)
}

// Layanan memegang aturan modul ini di atas satu Gudang.
type Layanan struct{ gudang Gudang }

// LayananOracle merakit Layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan { return &Layanan{gudang: repository.Baru(s.DB())} }

// LayananDengan merakit Layanan di atas Gudang mana pun - dipakai uji.
func LayananDengan(g Gudang) *Layanan { return &Layanan{gudang: g} }

// himpunanSah adalah keenam himpunan acuan, DISEBUT satu per satu.
var himpunanSah = map[models.Himpunan]bool{
	models.HimpunanMataUang:        true,
	models.HimpunanJenisPotongan:   true,
	models.HimpunanKelasBisnis:     true,
	models.HimpunanKelompokTreaty:  true,
	models.HimpunanBahaya:          true,
	models.HimpunanJenisReasuransi: true,
}

// DaftarAcuan membaca isi satu himpunan acuan.
//
// Identitas pelaku diperiksa lebih dulu: tidak ada jalur baca tanpa identitas,
// termasuk untuk data acuan.
func (l *Layanan) DaftarAcuan(ctx context.Context, p inti.Pelaku, h models.Himpunan) ([]models.Acuan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if !himpunanSah[h] {
		return nil, fmt.Errorf("%w: %q", ErrHimpunanTidakAda, h)
	}
	return l.gudang.DaftarAcuan(ctx, h)
}

// Pesan mengembalikan kalimat yang layak sampai ke layar.
func Pesan(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
