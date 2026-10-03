package tiruan

// Copy Old (permintaan work owner 03-10-2026) - tiruan kedua tabel JSON warisan: produk lama menurut ID. Produk lama
// TETAP ada sesudah disalin (tabel JSON hanya dibaca); yang berubah hanya keberadaannya di tabel flat (`Produk`).

import (
	"context"
	"fmt"
	"sort"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/repository"
)

// produkLama - satu produk tabel JSON warisan: bentuk flat yang akan ditulis dan baris popup-nya.
type produkLama struct {
	produk models.Produk
	status models.ProdukLama
}

// IsiLama - fixture satu produk tabel JSON warisan (kodek `repository.UraiProduk` + `NormalkanFlat`, seperti alat
// pindah). Belum ada di tabel flat sampai disalin. Fixture yang keliru = panik.
func (g *Gudang) IsiLama(id, jsonUmum, jsonInward string) {
	idInward := ""
	if jsonInward != "" {
		idInward = id
	}
	p, err := repository.UraiProduk(id, jsonUmum, idInward, jsonInward)
	if err != nil {
		panic(fmt.Sprintf("tiruan: fixture lama %s: %v", id, err))
	}
	n, masalah := repository.NormalkanFlat(p)
	if len(masalah) > 0 {
		panic(fmt.Sprintf("tiruan: fixture lama %s: %+v", id, masalah))
	}
	g.lama[id] = &produkLama{produk: n, status: models.ProdukLama{ID: id, ProductName: n.Umum.ProductName,
		Ceding: n.Umum.Ceding, TreatyNumber: n.Umum.TreatyNumber, InwardName: n.Umum.InwardName, CreateOp: n.Umum.CreateOp,
		UpdateOp: n.Umum.UpdateOp, BolehDisalin: true, Alasan: []string{}, Catatan: []string{}}}
}

// TolakLama - produk lama yang ditolak rekonsiliasi (alasan menyebut kolom, tanpa nilai).
func (g *Gudang) TolakLama(id string, alasan ...string) {
	l := g.lama[id]
	l.status.BolehDisalin, l.status.Alasan = false, alasan
}

// SiapkanProdukLama - baris popup produk lama yang belum ada di tabel flat, urut ID, dan bentuk flat yang boleh disalin.
func (g *Gudang) SiapkanProdukLama(context.Context) ([]models.ProdukLama, map[string]models.Produk, error) {
	if g.GagalBaca != nil {
		return nil, nil, g.GagalBaca
	}
	ids := make([]string, 0, len(g.lama))
	for id := range g.lama {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	daftar := []models.ProdukLama{}
	siap := map[string]models.Produk{}
	for _, id := range ids {
		if _, diFlat := g.Produk[id]; diFlat {
			continue
		}
		l := g.lama[id]
		daftar = append(daftar, l.status)
		if l.status.BolehDisalin {
			siap[id] = salinProduk(l.produk)
		}
	}
	return daftar, siap, nil
}

// SalinProdukLama - tulis satu produk lama ke tabel flat (tiruan); ID yang sudah ada = repository.ErrSudahDiFlat.
func (g *Gudang) SalinProdukLama(_ context.Context, _ *db.Tx, p models.Produk) error {
	if err := g.GagalSalinLama[p.ID]; err != nil {
		return err
	}
	if _, ada := g.Produk[p.ID]; ada {
		return fmt.Errorf("%w: %s", repository.ErrSudahDiFlat, p.ID)
	}
	return g.tulis(p)
}
