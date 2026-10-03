package repository

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// BacaFlatUji - pembaca produk flat untuk uji `db` paket luar (`mpnl_flat_db_test.go`).
func (g *Gudang) BacaFlatUji(ctx context.Context, id string) (models.Produk, error) {
	return g.bacaFlat(ctx, nil, id, false)
}

// TulisFlatUji - penulis produk flat (sisip) untuk uji `db` paket luar, di transaksi pemanggil.
func (g *Gudang) TulisFlatUji(ctx context.Context, tx *db.Tx, p models.Produk) error {
	return g.tulisFlat(ctx, tx, p, true)
}
