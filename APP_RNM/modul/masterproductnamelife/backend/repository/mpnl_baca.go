package repository

// Pembaca produk (paket 1) - grid daftar dan satu produk utuh.
//
//	grid      `BrowseProduct_Life` (`InboxProductName.xml` b76658): view `PRODUCT_LIFE`,
//	          tanpa saringan b681, urut `.ID ASC` b1094 - di sini `M_PRODUCT_LIFE`
//	          sendiri (P2), kunci yang sama dengan yang dibaca view
//	umum      `BrowseUnderwritingList` b61 `select JSONDATA … where ID = {ProductName.ID}`
//	          (`SetProductName` 3 b780 `·`, lalu `adoptJSONObject` 4 b957 `·`)
//	inward    `BrowseProductInward` saringan b834 `.PRODUCTID = Param.ProductID`
//	          (`SetProductNameInward` 2 b826 `·`, 3.1 b1046 `·` - baris TERAKHIR menang)
//
// ⛔ Baris inward juga dicari lewat `ID` = ID produk: itulah kunci yang
// dipakai Claim Life (`GetProductName.xml` `WHERE ID = …ProductNameID`) dan
// yang ditulis modul ini (R14). Tanpa itu, baris yang `PRODUCTID`-nya kosong
// tidak terlihat dan simpan berikutnya menggandakan `ID`-nya (nol PK di DEV).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

var (
	// ErrTidakAda - produk tidak ada di `M_PRODUCT_LIFE`.
	ErrTidakAda = errors.New("repository: product not found")
	// ErrIdentitasGanda - satu `ID` dipakai lebih dari satu baris (nol PK di
	// DEV): gagal terang, bukan memilih salah satunya diam-diam.
	ErrIdentitasGanda = errors.New("repository: the same ID is used by more than one row")
)

func sqlDaftarProduk(tabel string) string {
	return fmt.Sprintf(`SELECT ID, JSONDATA FROM %s ORDER BY ID ASC`, tabel)
}

func sqlAmbilProduk(tabel string, kunci bool) string {
	q := fmt.Sprintf(`SELECT ID, JSONDATA FROM %s WHERE ID = :1`, tabel)
	if kunci {
		q += ` FOR UPDATE`
	}
	return q
}

// sqlAmbilInward - baris inward produk: `PRODUCTID` (Pega) ATAU `ID` (Claim Life).
func sqlAmbilInward(tabel string, kunci bool) string {
	q := fmt.Sprintf(`SELECT ID, JSONDATA FROM %s
		WHERE JSON_VALUE(JSONDATA, '$.PRODUCTID') = :1 OR ID = :2
		ORDER BY ID ASC`, tabel)
	if kunci {
		q += ` FOR UPDATE`
	}
	return q
}

// DaftarProduk - grid daftar, urut ID.
func (g *Gudang) DaftarProduk(ctx context.Context) ([]models.RingkasanProduk, error) {
	q, err := g.siapkan(TabelProduk, sqlDaftarProduk)
	if err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: reading %s: %w", TabelProduk, err)
	}
	defer func() { _ = rows.Close() }()
	hasil := []models.RingkasanProduk{}
	for rows.Next() {
		var id, isi sql.NullString
		if err := rows.Scan(&id, &isi); err != nil {
			return nil, fmt.Errorf("repository: reading %s: %w", TabelProduk, err)
		}
		r, err := RingkasanDari(id.String, isi.String)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, r)
	}
	return hasil, rows.Err()
}

// barisJSON - satu baris (ID, JSONDATA).
type barisJSON struct{ id, isi string }

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

// SimpananProduk - JSON mentah kedua sisi satu produk, apa adanya di tabel.
// Dipakai penulis (paket 3–4) untuk mempertahankan kunci yang tidak dikelola.
type SimpananProduk struct {
	ID, JSONUmum         string
	IDInward, JSONInward string
	AdaInward            bool
}

// AmbilSimpanan membaca JSON mentah satu produk; `kunci` = `FOR UPDATE`
// (dipakai di dalam transaksi simpan).
func (g *Gudang) AmbilSimpanan(ctx context.Context, tx *db.Tx, id string, kunci bool) (SimpananProduk, error) {
	q, err := g.siapkan(TabelProduk, func(t string) string { return sqlAmbilProduk(t, kunci) })
	if err != nil {
		return SimpananProduk{}, err
	}
	umum, err := g.bacaBaris(ctx, tx, q, id)
	if err != nil {
		return SimpananProduk{}, fmt.Errorf("repository: reading %s %s: %w", TabelProduk, id, err)
	}
	switch len(umum) {
	case 0:
		return SimpananProduk{}, fmt.Errorf("%w: %s", ErrTidakAda, id)
	case 1:
	default:
		return SimpananProduk{}, fmt.Errorf("%w: %s %s (%d rows)", ErrIdentitasGanda, TabelProduk, id, len(umum))
	}
	s := SimpananProduk{ID: umum[0].id, JSONUmum: umum[0].isi}
	qi, err := g.siapkan(TabelInward, func(t string) string { return sqlAmbilInward(t, kunci) })
	if err != nil {
		return SimpananProduk{}, err
	}
	inward, err := g.bacaBaris(ctx, tx, qi, id, id)
	if err != nil {
		return SimpananProduk{}, fmt.Errorf("repository: reading %s for product %s: %w", TabelInward, id, err)
	}
	if b, ada := pilihInward(id, inward); ada {
		s.IDInward, s.JSONInward, s.AdaInward = b.id, b.isi, true
	}
	return s, nil
}

// pilihInward - baris inward produk: yang ber-`PRODUCTID` = produk (Pega,
// terakhir menang), lalu yang ber-`ID` = produk.
func pilihInward(id string, baris []barisJSON) (barisJSON, bool) {
	var dipilih barisJSON
	ada := false
	for _, b := range baris {
		obj, err := uraiObjek(b.isi)
		if err == nil && strings.TrimSpace(teksDari(obj[kunciProductID])) == id {
			dipilih, ada = b, true
		}
	}
	if ada {
		return dipilih, true
	}
	for _, b := range baris {
		if b.id == id {
			return b, true
		}
	}
	return barisJSON{}, false
}

// AmbilProduk - satu produk utuh (tombol `View` b74753).
func (g *Gudang) AmbilProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	s, err := g.AmbilSimpanan(ctx, tx, id, false)
	if err != nil {
		return models.Produk{}, err
	}
	return UraiProduk(s.ID, s.JSONUmum, s.IDInward, s.JSONInward)
}
