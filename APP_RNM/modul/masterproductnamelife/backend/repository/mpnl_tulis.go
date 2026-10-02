package repository

// Penulis produk - tabel FLAT (tiket 01 bab 02-10-2026, K5; dulu `JSONDATA` kedua tabel warisan, P1).
//
//	SaveProductName_Act 8 b1625 `·` PRE=false  DATAPEGA ← @GetPageJSONString() halaman ProductName
//	                    9 b1833 `·` PRE=false  RDB SaveProductNameLIfe → PEGA_M_PRODUCT_LIFE (TIDAK dipanggil)
//	                   13 b2530 / 14 b2719 / 15 b2864                    halaman ProductNameInward (TIDAK dipanggil)
//
// ⛔ Prosedur ditiru: upsert dikunci `ID`, `ID` baru dari sequence (P6); SATU baris induk `M_PRODUCTNAME_LIFE` (sisi
// umum + sisi inward, Q1b) dan ketujuh anak ditulis ULANG di transaksi pemanggil (P4). Nol COMMIT. Nilai lebih dulu
// diseragamkan `NormalkanFlat`; yang tidak muat kolomnya = ErrNilaiTidakMuat (services menolaknya lebih dulu).
// ⛔ Kedua tabel JSON warisan TIDAK ditulis lagi - cadangan dan sumber alat pindah saja.

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// KunciProduk - produk utuh, baris induknya dikunci `FOR UPDATE` (di dalam simpan).
func (g *Gudang) KunciProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	return g.bacaFlat(ctx, tx, id, true)
}

// SisipProduk menerbitkan ID baru (sequence; ID terpakai dilewati) dan menulis produk baru di transaksi pemanggil.
// ID inward = `PRODUCTID` = ID produk (R14). Mengembalikan ID itu.
//
// Salinan (`Copy`, `p.SalinanDari`): services sudah menyalin medan milik server produk asal (`CopyProduct` menyalin
// halaman utuh); tabel flat tidak punya kunci tak dikelola untuk diwarisi (D2).
func (g *Gudang) SisipProduk(ctx context.Context, tx *db.Tx, p models.Produk) (string, error) {
	if p.SalinanDari != "" {
		if _, err := g.bacaFlat(ctx, tx, p.SalinanDari, false); err != nil {
			return "", err
		}
	}
	id, err := g.identitasBaru(ctx, tx)
	if err != nil {
		return "", err
	}
	p.ID = id
	p.Inward.ID, p.Inward.ProductID = id, id
	if err := g.tulisFlat(ctx, tx, p, true); err != nil {
		return "", err
	}
	return id, nil
}

// PerbaruiProduk menulis ulang produk yang ada: baris induk diperbarui (tepat satu baris), seluruh anak ditulis ulang.
func (g *Gudang) PerbaruiProduk(ctx context.Context, tx *db.Tx, p models.Produk) error {
	p.Inward.ID, p.Inward.ProductID = p.ID, p.ID
	return g.tulisFlat(ctx, tx, p, false)
}
