package services

// Jalur baca produk (paket 1): grid daftar `InboxProductName` (mode daftar,
// wadah b71246) dan satu produk utuh (tombol `View` b74798 → `SetProductName`
// + `SetProductNameInward`).

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// Gudang adalah seluruh sentuhan basis data modul ini.
type Gudang interface {
	DaftarProduk(ctx context.Context) ([]models.RingkasanProduk, error)
	AmbilProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error)
	GudangMaster
	GudangTulis
	GudangLampiran
	GudangOutward
	GudangLama
}

var (
	// ErrProdukTidakAda - produk tidak ada (404); handler tidak mengimpor repository.
	ErrProdukTidakAda = errors.New("services: product not found")
	// ErrIdentitasGanda - satu ID dipakai beberapa baris: 500 berkalimat.
	ErrIdentitasGanda = repository.ErrIdentitasGanda
	// ErrJSONRusak - JSONDATA tabel warisan tidak terbaca (alat pindah): 500 berkalimat.
	ErrJSONRusak = repository.ErrJSONRusak
	// ErrIdentitasMelampauiLebar / ErrIdentitasBentrok - sequence harus ditinjau DBA (500 berkalimat).
	ErrIdentitasMelampauiLebar = repository.ErrIdentitasMelampauiLebar
	ErrIdentitasBentrok        = repository.ErrIdentitasBentrok
	// ErrNilaiTidakMuat - nilai tidak muat kolom flatnya (422). Services menolaknya lebih dulu berkalimat
	// (`periksaPanjang`); galat ini lapis kedua penulis.
	ErrNilaiTidakMuat = repository.ErrNilaiTidakMuat
)

// tidakAda menerjemahkan ErrTidakAda repository menjadi galat entitasnya.
func tidakAda(err, jadi error, id string) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return fmt.Errorf("%w: %s", jadi, id)
	}
	return err
}

// DaftarProduk - grid daftar (RD `BrowseProduct_Life`, urut `.ID ASC`).
func (l *Layanan) DaftarProduk(ctx context.Context, p inti.Pelaku) ([]models.RingkasanProduk, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	return l.gudang.DaftarProduk(ctx)
}

// AmbilProduk - satu produk utuh, kedua sisi dan ketujuh daftar.
func (l *Layanan) AmbilProduk(ctx context.Context, p inti.Pelaku, id string) (models.Produk, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Produk{}, err
	}
	pr, err := l.gudang.AmbilProduk(ctx, nil, id)
	return pr, tidakAda(err, ErrProdukTidakAda, id)
}
