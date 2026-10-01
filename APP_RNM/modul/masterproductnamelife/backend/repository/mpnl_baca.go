package repository

// Pembaca produk (paket 1) - grid daftar dan satu produk utuh.
//
//	grid      `BrowseProduct_Life` (`InboxProductName.xml` b76658): view `PRODUCT_LIFE`,
//	          tanpa saringan b682, urut `.ID ASC` b1096 - di sini `M_PRODUCT_LIFE`
//	          sendiri (P2), kunci yang sama dengan yang dibaca view
//	umum      `BrowseUnderwritingList` b61 `select JSONDATA … where ID = {ProductName.ID}`
//	          (`SetProductName` 3 b782 `·`, lalu `adoptJSONObject` 4 b959 `·`)
//	inward    `BrowseProductInward` saringan b840 `.PRODUCTID = Param.ProductID`
//	          (`SetProductNameInward` 2 b828 `·`, 3.1 b1048 `·` - baris TERAKHIR menang)
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

// MaksBarisGrid - `pyMaxRecords` RD `BrowseProduct_Life` b1078.
const MaksBarisGrid = 500

var (
	// ErrTidakAda - produk tidak ada di `M_PRODUCT_LIFE`.
	ErrTidakAda = errors.New("repository: product not found")
	// ErrIdentitasGanda - satu `ID` dipakai lebih dari satu baris (nol PK di
	// DEV): gagal terang, bukan memilih salah satunya diam-diam.
	ErrIdentitasGanda = errors.New("repository: the same ID is used by more than one row")
)

// sqlDaftarProduk - kelima kolom grid dibaca di Oracle (`JSON_VALUE`, seperti view
// `PRODUCT_LIFE`), bukan CLOB utuh: `CommentList` bertambah setiap simpan.
// `JSONDATA` dijaga constraint `IS JSON`, dan nilai yang ditulis modul ini dibatasi 4000 byte per kunci view
// (`LebarKunciView`), jadi `NULL ON ERROR` bawaan tidak menyembunyikan nilai tulisan modul ini; baris warisan
// yang melampauinya tampil kosong di grid, seperti di view `PRODUCT_LIFE`.
func sqlDaftarProduk(tabel string) string {
	// Batas baris = `pyMaxRecords` 500 `BrowseProduct_Life` b1078.
	return fmt.Sprintf(`SELECT ID, JSON_VALUE(JSONDATA, '$.CEDING'), JSON_VALUE(JSONDATA, '$.TREATYNUMBER'),
		JSON_VALUE(JSONDATA, '$.INWARDNAME'), JSON_VALUE(JSONDATA, '$.CREATEOP'), JSON_VALUE(JSONDATA, '$.UPDATEOP')
		FROM %s ORDER BY ID ASC FETCH FIRST %d ROWS ONLY`, tabel, MaksBarisGrid)
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
		var id, ceding, nomor, nama, buat, ubah sql.NullString
		if err := rows.Scan(&id, &ceding, &nomor, &nama, &buat, &ubah); err != nil {
			return nil, fmt.Errorf("repository: reading %s: %w", TabelProduk, err)
		}
		hasil = append(hasil, models.RingkasanProduk{ID: id.String, Ceding: ceding.String, TreatyNumber: nomor.String,
			InwardName: nama.String, CreateOp: buat.String, UpdateOp: ubah.String})
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
	// InwardMilikLain - `PRODUCTID` baris inward ber-`ID` = ID produk yang ternyata
	// milik produk LAIN (ID sequence inward warisan); kosong bila tidak ada.
	InwardMilikLain string
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
	} else {
		s.InwardMilikLain = inwardMilikLain(id, inward)
	}
	return s, nil
}

func productIDDari(b barisJSON) string {
	obj, err := uraiObjek(b.isi)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(teksDari(obj[kunciProductID]))
}

// pilihInward - baris inward produk: yang ber-`PRODUCTID` = produk (Pega,
// terakhir menang), lalu yang ber-`ID` = produk DAN `PRODUCTID`-nya kosong.
// ⛔ Baris ber-`ID` = produk yang `PRODUCTID`-nya menunjuk produk lain BUKAN
// milik produk ini (ID sequence inward warisan) - tidak pernah dipilih, dibaca, atau ditimpa.
func pilihInward(id string, baris []barisJSON) (barisJSON, bool) {
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

// inwardMilikLain - `PRODUCTID` baris ber-`ID` = produk yang milik produk lain.
func inwardMilikLain(id string, baris []barisJSON) string {
	for _, b := range baris {
		if pid := productIDDari(b); b.id == id && pid != "" && pid != id {
			return pid
		}
	}
	return ""
}

// AmbilProduk - satu produk utuh (tombol `View` b74798).
func (g *Gudang) AmbilProduk(ctx context.Context, tx *db.Tx, id string) (models.Produk, error) {
	s, err := g.AmbilSimpanan(ctx, tx, id, false)
	if err != nil {
		return models.Produk{}, err
	}
	return UraiProduk(s.ID, s.JSONUmum, s.IDInward, s.JSONInward)
}
