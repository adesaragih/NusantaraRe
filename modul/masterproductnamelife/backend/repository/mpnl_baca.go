package repository

// Pembaca produk - tabel FLAT (tiket 01 bab 02-10-2026, K5): grid daftar dan satu produk utuh.
//
//	grid      `BrowseProduct_Life` (`InboxProductName.xml` b76658): tanpa saringan b682, urut `.ID ASC` b1096,
//	          `pyMaxRecords` 500 b1078 - kelima kolom grid dibaca langsung dari induk `M_PRODUCTNAME_LIFE`
//	produk    induk + ketujuh anak (`mpnl_flat_sql.go` bacaFlat); dulu `BrowseUnderwritingList` b61 atas `JSONDATA`
//	          dan `BrowseProductInward` b840 atas baris inward - kini satu baris induk per produk
//
// Pembaca baris JSON warisan di bawah (barisJSON, pilihInward) hanya untuk alat pindah (`mpnl_pindah.go`).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// MaksBarisGrid - `pyMaxRecords` RD `BrowseProduct_Life` b1078.
const MaksBarisGrid = 500

var (
	// ErrTidakAda - produk tidak ada di tabel induk.
	ErrTidakAda = errors.New("repository: product not found")
	// ErrIdentitasGanda - satu `ID` dipakai lebih dari satu baris: gagal terang, bukan memilih salah satunya diam-diam
	// (tabel flat ber-PK; tetap dijaga untuk sumber JSON warisan yang tanpa PK).
	ErrIdentitasGanda = errors.New("repository: the same ID is used by more than one row")
)

// DaftarProduk - grid daftar, urut ID.
func (g *Gudang) DaftarProduk(ctx context.Context) ([]models.RingkasanProduk, error) {
	q, err := g.siapkan(TabelFlatInduk, sqlDaftarFlat)
	if err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s: %w", TabelFlatInduk, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.RingkasanProduk{}
	for rows.Next() {
		var id, ceding, nomor, nama, buat, ubah sql.NullString
		if err := rows.Scan(&id, &ceding, &nomor, &nama, &buat, &ubah); err != nil {
			return nil, fmt.Errorf("repository: reading %s: %w", TabelFlatInduk, err)
		}
		hasil = append(hasil, models.RingkasanProduk{ID: id.String, Ceding: ceding.String, TreatyNumber: nomor.String,
			InwardName: nama.String, CreateOp: buat.String, UpdateOp: ubah.String})
	}
	return hasil, rows.Err()
}

// AmbilProduk - satu produk utuh (tombol `View` b74798).
//
// ⛔ Di luar simpan (tx kosong) induk dan ketujuh anak dibaca di SATU transaksi baca-saja: satu potret data, bukan
// delapan pernyataan READ COMMITTED yang dapat menampilkan induk lama dengan anak baru dari simpan yang berjalan
// bersamaan (temuan /code-review 02-10-2026). Di dalam simpan, transaksi pemanggil dipakai.
func (g *Gudang) AmbilProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	if tx.Terisi() {
		return g.bacaFlat(ctx, tx, id, false)
	}
	baca, err := g.db.Mulai(ctx)
	if err != nil {
		return models.Produk{}, err
	}
	defer func() { _ = baca.Rollback() }()
	if _, err := baca.ExecContext(ctx, sqlTransaksiBacaSaja); err != nil {
		return models.Produk{}, fmt.Errorf("repository: opening a read-only snapshot: %w", err)
	}
	return g.bacaFlat(ctx, baca, id, false)
}

// sqlTransaksiBacaSaja - satu potret baca (Oracle: pernyataan PERTAMA transaksi).
const sqlTransaksiBacaSaja = `SET TRANSACTION READ ONLY`

// --- sumber JSON warisan (alat pindah) -------------------------------------------------

// barisJSON - satu baris (ID, JSONDATA) tabel warisan; `datar` = kolom datar `RIRISKID`, `RIRISK` `M_PRODUCT_LIFE`
// (`SaveProductNameLIfeFlat` b84) bila dibaca.
type barisJSON struct {
	id, isi  string
	datar    [2]string
	adaDatar bool
}

func (g *Gudang) bacaBaris(ctx context.Context, tx *db.Tx, q string, args ...any) ([]barisJSON, error) {
	rows, err := g.kueri(tx).QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var hasil []barisJSON
	for rows.Next() {
		var id, isi sql.NullString
		if err := rows.Scan(&id, &isi); err != nil {
			return nil, err
		}
		hasil = append(hasil, barisJSON{id: id.String, isi: isi.String})
	}
	return hasil, rows.Err()
}

func productIDDari(b barisJSON) string {
	obj, err := uraiObjek(b.isi)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(teksDari(obj[kunciProductID]))
}

// pilihInward - baris inward produk di sumber JSON: lebih dulu baris ber-`ID` = `PRODUCTID` = produk (baris yang
// dibaca pembaca hilir `WHERE ID = produk`), lalu yang ber-`PRODUCTID` = produk (Pega `SetProductNameInward` 3.1
// b1048, terakhir menang), lalu yang ber-`ID` = produk DAN `PRODUCTID`-nya kosong.
// ⛔ Baris ber-`ID` = produk yang `PRODUCTID`-nya menunjuk produk lain BUKAN milik produk ini - tidak pernah dipilih.
func pilihInward(id string, baris []barisJSON) (barisJSON, bool) {
	for _, b := range baris {
		if b.id == id && productIDDari(b) == id {
			return b, true
		}
	}
	var dipilih barisJSON
	ada := false
	for _, b := range baris {
		if productIDDari(b) == id {
			dipilih, ada = b, true
		}
	}
	if ada {
		return dipilih, true
	}
	for _, b := range baris {
		if b.id == id && productIDDari(b) == "" {
			return b, true
		}
	}
	return barisJSON{}, false
}
