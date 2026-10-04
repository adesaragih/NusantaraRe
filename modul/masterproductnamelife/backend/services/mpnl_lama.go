package services

// Copy Old - permintaan work owner 03-10-2026: tombol di samping `Add` membuka popup berisi produk tabel JSON warisan
// yang BELUM ada di tabel flat; yang dicentang disalin lewat `Process Copy`. Aturan salin SAMA dengan alat pindah
// (`repository.SiapkanProdukLama`); setiap produk disalin di transaksinya sendiri - satu produk yang gagal tidak
// membatalkan yang lain, dan hasilnya dilaporkan per produk. Tabel JSON hanya DIBACA; tulisan hanya ke tabel flat.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// GudangLama - sumber dan penulis Copy Old.
type GudangLama interface {
	// SiapkanProdukLama - baris popup (produk lama yang belum ada di tabel flat) dan bentuk flat yang boleh disalin.
	SiapkanProdukLama(ctx context.Context) ([]models.ProdukLama, map[string]models.Produk, error)
	// SalinProdukLama - tulis satu produk lama ke tabel flat di transaksi pemanggil, lalu baca ulang.
	SalinProdukLama(ctx context.Context, tx *db.Tx, p models.Produk) error
}

// MaksSalinLama - ID paling banyak dalam satu `Process Copy` (seluruh tabel lama DEV 02-10-2026: 196 produk).
const MaksSalinLama = 1000

// DaftarProdukLama - isi popup `Copy Old`.
func (l *Layanan) DaftarProdukLama(ctx context.Context, p inti.Pelaku) ([]models.ProdukLama, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	daftar, _, err := l.gudang.SiapkanProdukLama(ctx)
	return daftar, err
}

// SalinProdukLama - `Process Copy`: ID yang dicentang disalin satu per satu. Galat basis data satu produk dicatat log
// dan dilaporkan `gagal` untuk produk itu saja.
func (l *Layanan) SalinProdukLama(ctx context.Context, p inti.Pelaku, ids []string) (models.JawabanSalinLama, error) {
	j := models.JawabanSalinLama{Hasil: []models.HasilSalinLama{}}
	if err := inti.WajibIdentitas(p); err != nil {
		return j, err
	}
	var unik []string
	dilihat := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || dilihat[id] {
			continue
		}
		dilihat[id] = true
		unik = append(unik, id)
	}
	switch {
	case len(unik) == 0:
		return j, GalatValidasi{Pesan: []string{"select at least one old product to copy"}}
	case len(unik) > MaksSalinLama:
		return j, GalatValidasi{Pesan: []string{fmt.Sprintf("at most %d old products can be copied at once", MaksSalinLama)}}
	}
	daftar, siap, err := l.gudang.SiapkanProdukLama(ctx)
	if err != nil {
		return j, err
	}
	status := make(map[string]models.ProdukLama, len(daftar))
	for _, d := range daftar {
		status[d.ID] = d
	}
	for _, id := range unik {
		h := models.HasilSalinLama{ID: id, Pesan: []string{}}
		d, ada := status[id]
		switch {
		case !ada:
			// Tidak ada di antara produk lama yang belum disalin: sudah ada di tabel flat (mis. disalin pemakai lain).
			h.Status, h.Pesan = models.SalinSudahAda, []string{"already in the new tables, or not an old product"}
		case !d.BolehDisalin:
			h.Status, h.Pesan = models.SalinDitolak, d.Alasan
		default:
			err := l.tx(ctx, func(tx *db.Tx) error { return l.gudang.SalinProdukLama(ctx, tx, siap[id]) })
			switch {
			case err == nil:
				h.Status, h.Pesan = models.SalinDisalin, d.Catatan
				j.Disalin++
			case errors.Is(err, repository.ErrSudahDiFlat):
				h.Status, h.Pesan = models.SalinSudahAda, []string{"already in the new tables"}
			default:
				l.catat(fmt.Sprintf("master product name life: copying old product %s failed: %v", id, err))
				h.Status, h.Pesan = models.SalinGagal, []string{"database error - nothing was written for this product"}
			}
		}
		j.Hasil = append(j.Hasil, h)
	}
	l.catat(fmt.Sprintf("master product name life: %s copied %d of %d old product(s) from the JSON tables", p.AkunID, j.Disalin, len(unik)))
	return j, nil
}
