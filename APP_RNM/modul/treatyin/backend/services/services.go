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

	// Tiket 14. Keduanya menulis dan membaca KONTRAK + VERSI_KONTRAK.
	// Pembuatan berjalan dalam SATU transaksi: kontrak tanpa versi pertamanya
	// bukan keadaan yang sah menurut tiket 14, dan dua pernyataan terpisah
	// dapat meninggalkannya bila yang kedua gagal.
	BuatKontrakDenganVersiPertama(ctx context.Context, k models.Kontrak, v models.VersiKontrak) (int64, int64, error)
	BacaKontrak(ctx context.Context, id int64) (models.KontrakDenganVersi, error)

	// Tiket 16, 17, 18, 19 - identitas kontrak.
	CariKontrakSerupa(ctx context.Context, k models.Kontrak) ([]int64, error)
	CariKontrakLewatNomorWarisan(ctx context.Context, nomor string) ([]models.Kontrak, error)
	PerbaruiKontrak(ctx context.Context, k models.Kontrak) error
	TambahVersi(ctx context.Context, idKontrak int64, v models.VersiKontrak) (int64, error)

	// Tiket 35 dan 36 - ketentuan proporsional.
	AdaQuotaSharePadaVersi(ctx context.Context, idVersi int64) (bool, error)
	CatatKetentuanProporsional(ctx context.Context, idVersi, idLayer, idKelompok int64, jenis, persenQS string, lines *int64) error

	// Tiket 32 - syarat berbeda tiap pemulihan limit. SELURUH daftar sekaligus:
	// nomor urut kembar hanya terlihat bila barisnya dilihat bersama.
	CatatPemulihanLimit(ctx context.Context, idLayer int64, baris []models.PemulihanLimit) error

	// Tiket 40 - versi DASAR dibaca lewat ID_VERSI_KONTRAK_DASAR, bukan dari
	// salinan. Nil tanpa galat berarti versinya yang pertama.
	BacaVersiDasar(ctx context.Context, idVersi int64) (*models.VersiKontrak, error)

	// Tiket 41 TIDAK menambah apa pun di sini: identitas bentuk lama
	// DITURUNKAN di services dari BacaKontrak. Menambah kueri tersendiri
	// berarti dua pembaca untuk kolom yang sama (INV-60).
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
