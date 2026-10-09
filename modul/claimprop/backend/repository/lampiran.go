package repository

// Untuk apa berkas ini: LAMPIRAN KLAIM - master kategori `T_KATEGORI_DOC_KLAIM` (migrasi 535/536) dan baris tabel
// warisan dokumen klaim (kelas `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`, `InsertDocument_Act` Obj-Save). Katalog DEV
// 08-10-2026: tabel warisan itu 14 kolom, PK ID; IDPEGA kasus lama berbentuk pzInsKey Pega ("ASM-FW-GCNMFW-WORK CLMP-n"),
// kasus sistem baru ID T_WORK_CLAIM - pembaca mencari kedua bentuk.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
)

// Tabel lampiran.
const (
	TabelKategoriDokumen = "T_KATEGORI_DOC_KLAIM"
	TabelDokumenKlaim    = "DOCUMENT_CLAIM" // tabel warisan Pega
)

// ErrIDDokumenTerpakai - ORA-00001: ID dokumen bentrok, pemanggil mencoba ID berikut.
var ErrIDDokumenTerpakai = errors.New("repository: ID dokumen klaim sudah terpakai")

// sqlKategoriLampiran - `AttachCategory.pxResults`: setiap kategori master TYPE_KLAIM `:3` dan cacah dokumen klaim
// (IDPEGA `:1` / `:2`) berkategori itu. ⛔ Oracle mengikat bind SQL menurut URUTAN KEMUNCULAN, bukan nomor `:n` -
// nomor mengikuti urutan teks (sempat `:2, :3 ... :1`: nol kategori terbaca di DEV, 09-10-2026).
func sqlKategoriLampiran(kat, dok string) string {
	return fmt.Sprintf(`SELECT k.ID, k.LABEL, COUNT(d.ID)
		  FROM %s k
		  LEFT JOIN %s d ON d.KATEGORI_1 = k.ID AND d.IDPEGA IN (:1, :2)
		 WHERE k.TYPE_KLAIM = :3
		 GROUP BY k.ID, k.LABEL
		 ORDER BY NLSSORT(k.ID, 'NLS_SORT=BINARY')`, kat, dok)
}

// sqlDaftarLampiran - dokumen klaim (IDPEGA `:1` / `:2`), terbaru dahulu.
func sqlDaftarLampiran(dok string) string {
	return fmt.Sprintf(`SELECT ID, NAMAFILE, KATEGORI_1, KATEGORI_2, MIME, TO_CHAR(TANGGAL, 'DD/MM/YYYY HH24:MI'),
		       PXCREATEOPERATOR, T_STORAGE_ID
		  FROM %s WHERE IDPEGA IN (:1, :2) ORDER BY TANGGAL DESC, ID DESC`, dok)
}

// sqlSisipDokumenKlaim - InsertDocument_Act S3 + S5 (kolom yang S3 isi dengan nilai, ditambah T_STORAGE_ID S4).
func sqlSisipDokumenKlaim(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TANGGAL, IDPEGA, NAMAFILE, MIME, KATEGORI_1, T_STORAGE_ID, PXCREATEOPERATOR)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

func (g *Gudang) tabelLampiran() (kat, dok string, err error) {
	if kat, err = g.db.Qualify(TabelKategoriDokumen); err != nil {
		return "", "", err
	}
	dok, err = g.db.Qualify(TabelDokumenKlaim)
	return kat, dok, err
}

// KategoriLampiran membaca `AttachCategory.pxResults` kasus klaim `id`.
func (g *Gudang) KategoriLampiran(ctx context.Context, id string) ([]models.KategoriLampiran, error) {
	kat, dok, err := g.tabelLampiran()
	if err != nil {
		return nil, err
	}
	q := sqlKategoriLampiran(kat, dok)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, models.KunciInstans(id), models.KunciPegaLama(id), models.TypeKlaimProp)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca kategori lampiran: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.KategoriLampiran{} // kosong = [] (bukan null) di JSON
	for rows.Next() {
		var k models.KategoriLampiran
		if err := rows.Scan(&k.ID, &k.Label, &k.CountAttach); err != nil {
			return nil, err
		}
		out = append(out, k)
	}
	return out, rows.Err()
}

// DaftarLampiran membaca dokumen klaim kasus `id`.
func (g *Gudang) DaftarLampiran(ctx context.Context, id string) ([]models.Lampiran, error) {
	_, dok, err := g.tabelLampiran()
	if err != nil {
		return nil, err
	}
	q := sqlDaftarLampiran(dok)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, models.KunciInstans(id), models.KunciPegaLama(id))
	if err != nil {
		return nil, fmt.Errorf("repository: membaca dokumen klaim: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.Lampiran{} // kosong = [] (bukan null) di JSON
	for rows.Next() {
		var a models.Lampiran
		var nama, kat, kat2, mime, tgl, op, st *string
		if err := rows.Scan(&a.ID, &nama, &kat, &kat2, &mime, &tgl, &op, &st); err != nil {
			return nil, err
		}
		a.NamaFile, a.Kategori, a.Note, a.MIME, a.Tanggal, a.Operator = teksNil(nama), teksNil(kat), teksNil(kat2),
			teksNil(mime), teksNil(tgl), teksNil(op)
		a.StorageID = teksNil(st)
		a.AdaObjek = strings.TrimSpace(a.StorageID) != ""
		out = append(out, a)
	}
	return out, rows.Err()
}

func teksNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// sqlHapusDokumenKlaim - DeleteDocument pola NB (`DeleteDocumentPolis_Act`): baris dokumen klaim ini saja.
func sqlHapusDokumenKlaim(dok string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND IDPEGA IN (:2, :3)`, dok)
}

// HapusDokumenKlaim menghapus satu baris tabel warisan dokumen klaim kasus `id` di transaksi pemanggil.
func (g *Gudang) HapusDokumenKlaim(ctx context.Context, tx *db.Tx, id, lid string) error {
	if tx == nil {
		return errors.New("repository: dokumen klaim dihapus di dalam transaksi")
	}
	tabel, err := g.db.Qualify(TabelDokumenKlaim)
	if err != nil {
		return err
	}
	q := sqlHapusDokumenKlaim(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, lid, models.KunciInstans(id), models.KunciPegaLama(id))
	if err != nil {
		return fmt.Errorf("repository: menghapus dokumen klaim: %w", err)
	}
	return db.PastikanSatuBaris(h, "dokumen klaim")
}

// sqlPindahKategoriDokumen - Change Category (layar Pega View File, screenshot work owner 09-10-2026): KATEGORI_1 satu
// baris dokumen klaim kasus ini.
func sqlPindahKategoriDokumen(dok string) string {
	return fmt.Sprintf(`UPDATE %s SET KATEGORI_1 = :1 WHERE ID = :2 AND IDPEGA IN (:3, :4)`, dok)
}

// PindahKategoriDokumen memindahkan satu dokumen klaim kasus `id` ke kategori `kategori` di transaksi pemanggil.
func (g *Gudang) PindahKategoriDokumen(ctx context.Context, tx *db.Tx, id, lid, kategori string) error {
	if tx == nil {
		return errors.New("repository: kategori dokumen klaim diubah di dalam transaksi")
	}
	tabel, err := g.db.Qualify(TabelDokumenKlaim)
	if err != nil {
		return err
	}
	q := sqlPindahKategoriDokumen(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, kategori, lid, models.KunciInstans(id), models.KunciPegaLama(id))
	if err != nil {
		return fmt.Errorf("repository: memindah kategori dokumen klaim: %w", err)
	}
	return db.PastikanSatuBaris(h, "dokumen klaim")
}

// SisipDokumenKlaim menulis satu baris tabel warisan dokumen klaim di transaksi pemanggil; ID bentrok =
// `ErrIDDokumenTerpakai`.
func (g *Gudang) SisipDokumenKlaim(ctx context.Context, tx *db.Tx, d models.BarisDokumenKlaim) error {
	if tx == nil {
		return errors.New("repository: dokumen klaim ditulis di dalam transaksi")
	}
	tabel, err := g.db.Qualify(TabelDokumenKlaim)
	if err != nil {
		return err
	}
	q := sqlSisipDokumenKlaim(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, d.ID, d.Tanggal, teksAtauNil(d.IDPega), teksAtauNil(d.NamaFile),
		teksAtauNil(d.MIME), teksAtauNil(d.Kategori1), teksAtauNil(d.StorageID), teksAtauNil(d.Operator))
	if err != nil {
		if strings.Contains(err.Error(), "ORA-00001") {
			return fmt.Errorf("%w: %v", ErrIDDokumenTerpakai, err)
		}
		return fmt.Errorf("repository: menyisipkan dokumen klaim: %w", err)
	}
	return db.PastikanSatuBaris(h, "dokumen klaim")
}
