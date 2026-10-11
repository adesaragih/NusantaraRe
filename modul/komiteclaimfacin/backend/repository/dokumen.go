package repository

// Untuk apa berkas ini: BARIS DOKUMEN KLAIM - `InsertDocument_Act` (kelas `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`, Obj-Save)
// untuk PDF akseptasi `PrintPDFAccep_MultiAksep_KMT` S27. Tabel warisan DOCUMENT_CLAIM (katalog DEV 08-10-2026, 14
// kolom, PK ID; baris AcceptanceNote berisi KATEGORI_2 / NOAKSEP / NOPREKAS / PAYMENTDATE NULL - parameter S27 kosong).
// Disalin dari Komite Claim Prop, bukan diimpor.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// TabelDokumenKlaim - tabel warisan kelas `ASM-FW-GCNMFW-Int-DOCUMENT_CLAIM`.
const TabelDokumenKlaim = "DOCUMENT_CLAIM" // tabel warisan Pega

// ErrIDDokumenTerpakai - ORA-00001: ID dokumen (`yyyyMMddhhmmssSSS`) bentrok, pemanggil mencoba ID berikut.
var ErrIDDokumenTerpakai = errors.New("repository: ID dokumen klaim sudah terpakai")

// sqlSisipDokumenKlaim - InsertDocument_Act S3 + S5 (kolom yang S3 isi dengan nilai, ditambah T_STORAGE_ID S4).
func sqlSisipDokumenKlaim(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TANGGAL, IDPEGA, NAMAFILE, MIME, KATEGORI_1, T_STORAGE_ID, PXCREATEOPERATOR)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, tabel)
}

// SisipDokumenKlaim menulis satu baris tabel warisan DOCUMENT_CLAIM di transaksi pemanggil.
func (g *Gudang) SisipDokumenKlaim(ctx context.Context, tx *db.Tx, d models.BarisDokumenKlaim) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	tabel, err := g.db.Qualify(TabelDokumenKlaim)
	if err != nil {
		return err
	}
	q := sqlSisipDokumenKlaim(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, d.ID, d.Tanggal, db.KosongJadiNil(d.IDPega), db.KosongJadiNil(d.NamaFile),
		db.KosongJadiNil(d.MIME), db.KosongJadiNil(d.Kategori1), db.KosongJadiNil(d.StorageID),
		db.KosongJadiNil(d.Operator))
	if err != nil {
		if strings.Contains(err.Error(), "ORA-00001") {
			return fmt.Errorf("%w: %v", ErrIDDokumenTerpakai, err)
		}
		return fmt.Errorf("repository: menyisipkan dokumen klaim: %w", err)
	}
	return db.PastikanSatuBaris(h, "dokumen klaim")
}
