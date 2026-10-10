package repository

// Lampiran berkas - tabel warisan `M_ATTACHMENTBORDEREAUX` (PK `ID`, katalog DEV 08-10-2026) dan master
// `M_KATEGORIBORDEREAUX` (dibaca saja), SQL rule Pega apa adanya:
//
//	GetKategoryDocBDX_SQL b80-90   kategori + COUNT lampiran berkas ini (LEFT JOIN); urutan ID ditambahkan - GROUP BY
//	                               Pega tidak berurutan
//	AttachDocumentBdx_SQL browse   lampiran satu kategori: `WHERE BDX_ID = … AND CATEGORY_ID = …`
//	AttachDocumentBdx_SQL save     `INSERT (ID, BDX_ID, CATEGORY_ID, CATEGORY, FILENAME, FILEMIMETYPE, USERNAME,
//	                               T_STORAGE_ID)`
//	AttachDocumentBdx_SQL delete   `WHERE BDX_ID = … AND ID = …`
//
// Objek berkasnya di `T_STORAGE_IMAGE` - milik `inti/backend/penyimpanan`, bukan modul ini.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/bordereaux/backend/models"
)

// Tabel lampiran warisan.
const (
	TabelLampiran         = "M_ATTACHMENTBORDEREAUX"
	TabelKategoriLampiran = "M_KATEGORIBORDEREAUX"
)

const kolomLampiran = `ID, BDX_ID, CATEGORY_ID, CATEGORY, FILENAME, FILEMIMETYPE, USERNAME, T_STORAGE_ID`

func sqlKategoriLampiran(kategori, lampiran string) string {
	return fmt.Sprintf(`SELECT B.ID, B.NOTE, COUNT(A.CATEGORY_ID)
		FROM %s B LEFT JOIN %s A ON B.ID = A.CATEGORY_ID AND A.BDX_ID = :1
		GROUP BY B.ID, B.NOTE ORDER BY B.ID`, kategori, lampiran)
}

func sqlNamaKategori(kategori string) string {
	return fmt.Sprintf(`SELECT NOTE FROM %s WHERE ID = :1`, kategori)
}

func sqlDaftarLampiran(lampiran string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE BDX_ID = :1 AND CATEGORY_ID = :2 ORDER BY ID`, kolomLampiran, lampiran)
}

func sqlAmbilLampiran(lampiran string) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE BDX_ID = :1 AND CATEGORY_ID = :2 AND ID = :3`, kolomLampiran, lampiran)
}

func sqlSisipLampiran(lampiran string) string {
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, lampiran, kolomLampiran)
}

func sqlHapusLampiran(lampiran string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE BDX_ID = :1 AND ID = :2`, lampiran)
}

func pindaiLampiran(p pemindai) (models.Lampiran, error) {
	v, err := pindaiTeks(p, 8)
	if err != nil {
		return models.Lampiran{}, err
	}
	return models.Lampiran{ID: v[0], BdxID: v[1], KategoriID: v[2], Kategori: v[3], FileName: v[4], Ekstensi: v[5],
		Username: v[6], StorageID: v[7]}, nil
}

// KategoriLampiran - grid `AttachmentsBdx` (`GetKategotyDocBdx`).
func (g *Gudang) KategoriLampiran(ctx context.Context, bdxID string) ([]models.KategoriLampiran, error) {
	k, err := g.nama(TabelKategoriLampiran)
	if err != nil {
		return nil, err
	}
	l, err := g.nama(TabelLampiran)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlKategoriLampiran(k, l), func(p pemindai) (models.KategoriLampiran, error) {
		v, err := pindaiTeks(p, 3)
		if err != nil {
			return models.KategoriLampiran{}, err
		}
		n, _ := strconv.Atoi(v[2])
		return models.KategoriLampiran{ID: v[0], Nama: v[1], Cacah: n}, nil
	}, bdxID)
}

// NamaKategoriLampiran - NOTE satu kategori; ada false = ID bukan kategori master.
func (g *Gudang) NamaKategoriLampiran(ctx context.Context, id string) (string, bool, error) {
	k, err := g.nama(TabelKategoriLampiran)
	if err != nil {
		return "", false, err
	}
	q := sqlNamaKategori(k)
	if err := db.PeriksaSQL(q); err != nil {
		return "", false, err
	}
	var nama sql.NullString
	err = g.db.QueryRowContext(ctx, q, id).Scan(&nama)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, bungkus(err, "membaca kategori lampiran")
	}
	return nama.String, true, nil
}

// DaftarLampiran - popup `AttachmentDetailBdx` (`getAttcachmentList`).
func (g *Gudang) DaftarLampiran(ctx context.Context, bdxID, kategoriID string) ([]models.Lampiran, error) {
	l, err := g.nama(TabelLampiran)
	if err != nil {
		return nil, err
	}
	return daftar(ctx, g.db, sqlDaftarLampiran(l), pindaiLampiran, bdxID, kategoriID)
}

// AmbilLampiran - satu lampiran berkas dan kategori itu; ada false = tidak ada.
func (g *Gudang) AmbilLampiran(ctx context.Context, bdxID, kategoriID, id string) (models.Lampiran, bool, error) {
	l, err := g.nama(TabelLampiran)
	if err != nil {
		return models.Lampiran{}, false, err
	}
	hs, err := daftar(ctx, g.db, sqlAmbilLampiran(l), pindaiLampiran, bdxID, kategoriID, id)
	if err != nil || len(hs) == 0 {
		return models.Lampiran{}, false, err
	}
	return hs[0], true, nil
}

// SisipLampiran - `AttachDocumentBdx_SQL` save; ID bentrok = `ErrIDTerpakai`.
func (g *Gudang) SisipLampiran(ctx context.Context, tx *db.Tx, a models.Lampiran) error {
	l, err := g.nama(TabelLampiran)
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, g.dari(tx), sqlSisipLampiran(l), "menyisipkan lampiran", a.ID, a.BdxID, a.KategoriID,
		db.KosongJadiNil(a.Kategori), a.FileName, db.KosongJadiNil(a.Ekstensi), a.Username, a.StorageID)
	return err
}

// HapusLampiran - `AttachDocumentBdx_SQL` delete.
func (g *Gudang) HapusLampiran(ctx context.Context, tx *db.Tx, bdxID, id string) error {
	l, err := g.nama(TabelLampiran)
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, g.dari(tx), sqlHapusLampiran(l), "menghapus lampiran", bdxID, id)
	return err
}
