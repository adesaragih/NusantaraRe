// Package tiruan memuat gudang DI MEMORI modul Master Product Name Life -
// HANYA untuk uji services dan handlers (tidak diimpor kode produksi; penjaga
// modul `TestMPNLTiruanHanyaDiUji`).
//
// ⭐ Tiruan menyimpan `JSONDATA` MENTAH dan memakai kodek repository yang
// sungguhan (`repository.UraiProduk`, `RingkasanDari`), sehingga uji services
// dan handlers ikut menguji bentuk JSON Pega - bukan bentuk karangan tiruan.
// Fixture selalu berawalan `UJI-`.
package tiruan

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// Gudang adalah pengganti `repository.Gudang` di memori.
type Gudang struct {
	// Umum - `M_PRODUCT_LIFE`: ID → JSONDATA.
	Umum map[string]string
	// Inward - `M_PRODUCTINWARD_LIFE`: ID baris → JSONDATA.
	Inward map[string]string

	// Komit mencacah transaksi yang ditutup sukses.
	Komit int
	// GagalBaca - bila terisi, setiap pembacaan gagal dengan galat ini.
	GagalBaca error
}

// Baru menyusun gudang kosong.
func Baru() *Gudang {
	return &Gudang{Umum: map[string]string{}, Inward: map[string]string{}}
}

// Transaksi - tiruan `DalamTransaksi`: fn(nil); sukses = Komit++.
func (g *Gudang) Transaksi(_ context.Context, fn func(tx *db.Tx) error) error {
	if err := fn(nil); err != nil {
		return err
	}
	g.Komit++
	return nil
}

// DaftarProduk - urut ID.
func (g *Gudang) DaftarProduk(context.Context) ([]models.RingkasanProduk, error) {
	if g.GagalBaca != nil {
		return nil, g.GagalBaca
	}
	id := make([]string, 0, len(g.Umum))
	for k := range g.Umum {
		id = append(id, k)
	}
	sort.Strings(id)
	hasil := []models.RingkasanProduk{}
	for _, k := range id {
		r, err := repository.RingkasanDari(k, g.Umum[k])
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, r)
	}
	return hasil, nil
}

// AmbilProduk - kedua sisi; inward dicari lewat `PRODUCTID` lalu `ID`, seperti repository.
func (g *Gudang) AmbilProduk(_ context.Context, _ *db.Tx, id string) (models.Produk, error) {
	if g.GagalBaca != nil {
		return models.Produk{}, g.GagalBaca
	}
	umum, ada := g.Umum[id]
	if !ada {
		return models.Produk{}, fmt.Errorf("%w: %s", repository.ErrTidakAda, id)
	}
	idIn, isiIn := g.cariInward(id)
	return repository.UraiProduk(id, umum, idIn, isiIn)
}

func (g *Gudang) cariInward(id string) (string, string) {
	kunci := make([]string, 0, len(g.Inward))
	for k := range g.Inward {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	idPilih, isiPilih := "", ""
	for _, k := range kunci {
		p, err := repository.UraiProduk("", "", k, g.Inward[k])
		if err == nil && p.Inward.ProductID == id {
			idPilih, isiPilih = k, g.Inward[k]
		}
	}
	if idPilih != "" {
		return idPilih, isiPilih
	}
	if isi, ada := g.Inward[id]; ada {
		return id, isi
	}
	return "", ""
}

// ErrTiruan - galat buatan uji.
var ErrTiruan = errors.New("tiruan: injected failure")
