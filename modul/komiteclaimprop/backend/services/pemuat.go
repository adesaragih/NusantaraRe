package services

// Untuk apa berkas ini: PEMUAT DATA LAMA (tiket 13) - sensus uji-kering kasus komite warisan untuk klaim yang ditunda
// pemuat Claim Prop. Nol tulisan: tangga lama tidak tersimpan dalam bentuk yang dapat dimuat tanpa mengarang
// (`models/lama.go`), maka `Jalankan` menolak terang sampai sumbernya diputuskan.

import (
	"context"
	"errors"

	"nusantarare/modul/komiteclaimprop/backend/models"
)

// SumberLama - bacaan data lama (repository.Gudang).
type SumberLama interface {
	BacaOSLama(ctx context.Context) ([]models.BarisOSLama, error)
	BacaRiwayatLama(ctx context.Context) ([]models.RiwayatLama, error)
	BacaKomitePega(ctx context.Context) ([]models.KomitePega, error)
}

// ErrPemuatTanpaSumberTangga - `-jalankan` ditolak: tangga lama belum punya sumber yang dapat dimuat (OQ).
var ErrPemuatTanpaSumberTangga = errors.New("services: tangga komite lama belum punya sumber yang dapat dimuat " +
	"tanpa mengarang (HISTORYAKSEPTASIPEGA tanpa tingkat / jabatan; BLOB DATAPEGA belum ada pengurainya) - OQ, " +
	"-jalankan ditolak")

// ErrPemuatProduksi - `-jalankan` ditolak di IS_PEGA_PROD=true.
var ErrPemuatProduksi = errors.New("services: pemuat data lama komite tidak dijalankan di IS_PEGA_PROD=true")

// SensusLama membaca ketiga sumber lalu mencacah (uji-kering). `pegaGagal` - galat bacaan DATAPEGA (hak baca dapat
// tidak ada di lingkungan lain): sensus tetap jalan dengan nol work object, galatnya dikembalikan untuk dilaporkan.
func SensusLama(ctx context.Context, s SumberLama) (hasil models.SensusLama, pegaGagal error, err error) {
	os, err := s.BacaOSLama(ctx)
	if err != nil {
		return models.SensusLama{}, nil, err
	}
	riw, err := s.BacaRiwayatLama(ctx)
	if err != nil {
		return models.SensusLama{}, nil, err
	}
	pega, pegaGagal := s.BacaKomitePega(ctx)
	if pegaGagal != nil {
		pega = nil
	}
	return models.Sensus(os, riw, pega), pegaGagal, nil
}

// JalankanLama - mode tulis: ditolak di produksi, dan ditolak di mana pun selama sumber tangga belum diputuskan.
func JalankanLama(produksi bool) error {
	if produksi {
		return ErrPemuatProduksi
	}
	return ErrPemuatTanpaSumberTangga
}
