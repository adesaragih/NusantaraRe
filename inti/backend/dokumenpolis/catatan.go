package dokumenpolis

// Catatan di Oracle - tabel warisan `DOCUMENT_POLIS` (PK `ID`, katalog DEV 08-10-2026) dan master `CATEGORY_ATTACH_REAS`
// (dibaca saja; NOTE unik, 20 baris di DEV):
//
//	CategoryAttach_SQL        `select * from CATEGORY_ATTACH_REAS order by note`
//	InputParamUploadReas_act  Obj-Browse `DOCUMENT_POLIS` IDPEGA = pzInsKey, KATEGORI_2 = NOTE (dicacah `.CountAttach`)
//	InsertDocument_Act        Obj-Save NewDocument (ID, TANGGAL, IDPEGA, NAMAFILE, MIME, KATEGORI_1, KATEGORI_2,
//	                          INSKEY_LINK "", INSKEY_DATA "", pxCreateOperator, T_STORAGE_ID)
//	DeleteDocumentPolis_Act   Obj-Delete DOCUMENT_POLIS
//
// IDPEGA dibaca dalam BEBERAPA bentuk (`Kasus.Baca`): kunci sistem baru dan pzInsKey Pega lama.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
)

// Tabel warisan.
const (
	TabelDokumen  = "DOCUMENT_POLIS"
	TabelKategori = "CATEGORY_ATTACH_REAS"
)

// penampungIn - `:<mulai>, :<mulai+1>, …` sebanyak n.
func penampungIn(mulai, n int) string {
	p := make([]string, n)
	for i := range p {
		p[i] = ":" + strconv.Itoa(mulai+i)
	}
	return strings.Join(p, ", ")
}

func sqlKategori(kategori, dokumen string, n int) string {
	return fmt.Sprintf(`SELECT C.NOTE, COUNT(D.ID) FROM %s C
		LEFT JOIN %s D ON D.KATEGORI_2 = C.NOTE AND D.IDPEGA IN (%s)
		GROUP BY C.NOTE ORDER BY C.NOTE`, kategori, dokumen, penampungIn(1, n))
}

func sqlAdaKategori(kategori string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE NOTE = :1`, kategori)
}

const kolomDokumen = `ID, NAMAFILE, MIME, KATEGORI_2, TO_CHAR(TANGGAL, 'DD-MM-YYYY HH24:MI'), PXCREATEOPERATOR, T_STORAGE_ID`

// sqlDaftar - urut Upload Date: ID berformat jam 12-an (`hh`) tidak berurut waktu.
func sqlDaftar(dokumen string, n int) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDPEGA IN (%s) AND KATEGORI_2 = :%d ORDER BY TANGGAL, ID`,
		kolomDokumen, dokumen, penampungIn(1, n), n+1)
}

func sqlAmbil(dokumen string, n int) string {
	return fmt.Sprintf(`SELECT %s FROM %s WHERE IDPEGA IN (%s) AND ID = :%d`, kolomDokumen, dokumen, penampungIn(1, n), n+1)
}

// sqlSisip - TANGGAL = saat simpan (`@CurrentDateTime()`), INSKEY_LINK / INSKEY_DATA kosong (`""` = NULL).
func sqlSisip(dokumen string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, IDPEGA, TANGGAL, NAMAFILE, MIME, KATEGORI_1, KATEGORI_2, INSKEY_LINK, INSKEY_DATA,
		T_STORAGE_ID, PXCREATEOPERATOR) VALUES (:1, :2, SYSDATE, :3, :4, :5, :6, NULL, NULL, :7, :8)`, dokumen)
}

func sqlHapus(dokumen string, n int) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND IDPEGA IN (%s)`, dokumen, penampungIn(2, n))
}

type catatanOracle struct{ db *db.DB }

func (c catatanOracle) siapkan(objek string, susun func(string) string) (string, error) {
	if c.db == nil {
		return "", db.ErrTanpaOracle
	}
	tabel, err := c.db.Qualify(objek)
	if err != nil {
		return "", err
	}
	q := susun(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	return q, nil
}

func argumen(baca []string, lain ...any) []any {
	a := make([]any, 0, len(baca)+len(lain))
	for _, b := range baca {
		a = append(a, b)
	}
	return append(a, lain...)
}

func (c catatanOracle) Kategori(ctx context.Context, baca []string) ([]Kategori, error) {
	dok, err := c.siapkan(TabelDokumen, func(t string) string { return t })
	if err != nil {
		return nil, err
	}
	q, err := c.siapkan(TabelKategori, func(t string) string { return sqlKategori(t, dok, len(baca)) })
	if err != nil {
		return nil, err
	}
	rows, err := c.db.QueryContext(ctx, q, argumen(baca)...)
	if err != nil {
		return nil, fmt.Errorf("dokumenpolis: reading %s: %w", TabelKategori, err)
	}
	defer func() { _ = rows.Close() }()
	out := []Kategori{}
	for rows.Next() {
		var nama sql.NullString
		var n int
		if err := rows.Scan(&nama, &n); err != nil {
			return nil, fmt.Errorf("dokumenpolis: reading %s: %w", TabelKategori, err)
		}
		out = append(out, Kategori{Nama: nama.String, Cacah: n})
	}
	return out, rows.Err()
}

func (c catatanOracle) AdaKategori(ctx context.Context, nama string) (bool, error) {
	q, err := c.siapkan(TabelKategori, sqlAdaKategori)
	if err != nil {
		return false, err
	}
	var n int
	if err := c.db.QueryRowContext(ctx, q, nama).Scan(&n); err != nil {
		return false, fmt.Errorf("dokumenpolis: reading %s: %w", TabelKategori, err)
	}
	return n > 0, nil
}

func (c catatanOracle) baca(ctx context.Context, q string, args ...any) ([]Dokumen, error) {
	rows, err := c.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("dokumenpolis: reading %s: %w", TabelDokumen, err)
	}
	defer func() { _ = rows.Close() }()
	out := []Dokumen{}
	for rows.Next() {
		var v [7]sql.NullString
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5], &v[6]); err != nil {
			return nil, fmt.Errorf("dokumenpolis: reading %s: %w", TabelDokumen, err)
		}
		out = append(out, Dokumen{ID: v[0].String, NamaFile: v[1].String, Ekstensi: v[2].String, Kategori: v[3].String,
			Tanggal: v[4].String, Pengunggah: v[5].String, StorageID: v[6].String, AdaObjek: strings.TrimSpace(v[6].String) != ""})
	}
	return out, rows.Err()
}

func (c catatanOracle) Daftar(ctx context.Context, baca []string, kategori string) ([]Dokumen, error) {
	q, err := c.siapkan(TabelDokumen, func(t string) string { return sqlDaftar(t, len(baca)) })
	if err != nil {
		return nil, err
	}
	return c.baca(ctx, q, argumen(baca, kategori)...)
}

func (c catatanOracle) Ambil(ctx context.Context, baca []string, id string) (Dokumen, bool, error) {
	q, err := c.siapkan(TabelDokumen, func(t string) string { return sqlAmbil(t, len(baca)) })
	if err != nil {
		return Dokumen{}, false, err
	}
	d, err := c.baca(ctx, q, argumen(baca, id)...)
	if err != nil || len(d) == 0 {
		return Dokumen{}, false, err
	}
	return d[0], true, nil
}

func (c catatanOracle) Sisip(ctx context.Context, tx *db.Tx, d Dokumen, idPega string) error {
	if !tx.Terisi() {
		return errors.New("dokumenpolis: inserting a document requires a transaction")
	}
	q, err := c.siapkan(TabelDokumen, sqlSisip)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, q, d.ID, idPega, d.NamaFile, db.KosongJadiNil(d.Ekstensi), KategoriReas, d.Kategori,
		db.KosongJadiNil(d.StorageID), d.Pengunggah)
	if err != nil && strings.Contains(err.Error(), "ORA-00001") {
		return fmt.Errorf("%w: %v", ErrIDTerpakai, err)
	}
	if err != nil {
		return fmt.Errorf("dokumenpolis: inserting %s: %w", TabelDokumen, err)
	}
	return nil
}

func (c catatanOracle) Hapus(ctx context.Context, tx *db.Tx, baca []string, id string) error {
	if !tx.Terisi() {
		return errors.New("dokumenpolis: deleting a document requires a transaction")
	}
	q, err := c.siapkan(TabelDokumen, func(t string) string { return sqlHapus(t, len(baca)) })
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, append([]any{id}, argumen(baca)...)...); err != nil {
		return fmt.Errorf("dokumenpolis: deleting %s %s: %w", TabelDokumen, id, err)
	}
	return nil
}
