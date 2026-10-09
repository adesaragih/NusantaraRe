package penyimpanan

// Catatan objek di Oracle - tabel warisan `T_STORAGE_IMAGE` (tanpa PK) dan `T_FOLDER_IMAGE`, SQL rule kelas
// `T_STORAGE_IMAGE` apa adanya (pernah dipakai Master Product Name Life dan terbukti di DEV 03-10-2026):
//
//	Insert_T_Storage_SQL  INSERT (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE='standard')
//	GetLinkStorage_SQL    SELECT URLPUBLIC, APPFOLDER, EXPDATE, APPNAME, FILENAME WHERE imageid = …
//	Update_T_Storage_SQL  UPDATE URLPUBLIC, APPFOLDER, EXPDATE, TANGGAL_UPLOAD WHERE imageid = …
//	DeleteStorage_SQL     DELETE WHERE imageid = …
//	GetAppName_SQL        SELECT APPNAME FROM T_FOLDER_IMAGE

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"nusantarare/inti/backend/db"
)

// Tabel warisan kelas `T_STORAGE_IMAGE`.
const (
	TabelObjek  = "T_STORAGE_IMAGE"
	TabelFolder = "T_FOLDER_IMAGE"
)

func sqlSisipObjek(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)
		VALUES (:1, :2, :3, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, :6, 'standard')`, tabel)
}

// sqlAmbilObjek - `pxResults(1)`; EXPDATE dibaca dalam bentuk tulisnya sendiri, tanpa tafsir zona waktu oleh driver.
func sqlAmbilObjek(tabel string) string {
	return fmt.Sprintf(`SELECT URLPUBLIC, APPFOLDER, TO_CHAR(EXPDATE, 'DD/MM/YYYY HH24:MI:SS'), APPNAME, FILENAME
		FROM %s WHERE IMAGEID = :1 FETCH FIRST 1 ROWS ONLY`, tabel)
}

func sqlPerbaruiObjek(tabel string) string {
	return fmt.Sprintf(`UPDATE %s SET URLPUBLIC = :1, APPFOLDER = :2,
		EXPDATE = TO_DATE(:3, 'DD/MM/YYYY HH24:MI:SS'), TANGGAL_UPLOAD = TO_DATE(:4, 'MM/DD/YYYY HH24:MI:SS')
		WHERE IMAGEID = :5`, tabel)
}

func sqlHapusObjek(tabel string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IMAGEID = :1`, tabel)
}

func sqlNamaAplikasi(tabel string) string {
	return fmt.Sprintf(`SELECT APPNAME FROM %s FETCH FIRST 1 ROWS ONLY`, tabel)
}

type catatanOracle struct{ db *db.DB }

func (c catatanOracle) siapkan(objek string, susun func(string) string) (string, error) {
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

func (c catatanOracle) NamaAplikasi(ctx context.Context) (string, error) {
	q, err := c.siapkan(TabelFolder, sqlNamaAplikasi)
	if err != nil {
		return "", err
	}
	var app sql.NullString
	err = c.db.QueryRowContext(ctx, q).Scan(&app)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("penyimpanan: reading %s: %w", TabelFolder, err)
	}
	return app.String, nil
}

func (c catatanOracle) Ambil(ctx context.Context, imageID string) (Objek, bool, error) {
	if imageID == "" {
		return Objek{}, false, nil
	}
	q, err := c.siapkan(TabelObjek, sqlAmbilObjek)
	if err != nil {
		return Objek{}, false, err
	}
	var url, folder, exp, app, nama sql.NullString
	err = c.db.QueryRowContext(ctx, q, imageID).Scan(&url, &folder, &exp, &app, &nama)
	if errors.Is(err, sql.ErrNoRows) {
		return Objek{}, false, nil
	}
	if err != nil {
		return Objek{}, false, fmt.Errorf("penyimpanan: reading %s: %w", TabelObjek, err)
	}
	return Objek{ImageID: imageID, URLPublic: url.String, AppFolder: folder.String, Exp: exp.String, AppName: app.String,
		FileName: nama.String}, true, nil
}

func (c catatanOracle) Catat(ctx context.Context, tx *db.Tx, o Objek) error {
	if !tx.Terisi() {
		return errors.New("penyimpanan: recording a storage object requires a transaction")
	}
	q, err := c.siapkan(TabelObjek, sqlSisipObjek)
	if err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, o.ImageID, db.KosongJadiNil(o.URLPublic), o.AppFolder, db.KosongJadiNil(o.Exp),
		o.FileName, db.KosongJadiNil(o.AppName))
	if err != nil {
		return fmt.Errorf("penyimpanan: recording %s %s: %w", TabelObjek, o.ImageID, err)
	}
	return db.PastikanSatuBaris(hasil, TabelObjek)
}

// Perbarui - seperti Pega, SEMUA baris ber-IMAGEID itu diperbarui (`T_STORAGE_IMAGE` tanpa PK).
func (c catatanOracle) Perbarui(ctx context.Context, o Objek) error {
	q, err := c.siapkan(TabelObjek, sqlPerbaruiObjek)
	if err != nil {
		return err
	}
	if _, err := c.db.ExecContext(ctx, q, db.KosongJadiNil(o.URLPublic), db.KosongJadiNil(o.AppFolder),
		db.KosongJadiNil(o.Exp), db.KosongJadiNil(o.TanggalUpload), o.ImageID); err != nil {
		return fmt.Errorf("penyimpanan: updating %s %s: %w", TabelObjek, o.ImageID, err)
	}
	return nil
}

// Hapus - objek yang sudah tidak tercatat bukan galat.
func (c catatanOracle) Hapus(ctx context.Context, tx *db.Tx, imageID string) error {
	if !tx.Terisi() {
		return errors.New("penyimpanan: deleting a storage record requires a transaction")
	}
	q, err := c.siapkan(TabelObjek, sqlHapusObjek)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, imageID); err != nil {
		return fmt.Errorf("penyimpanan: deleting %s %s: %w", TabelObjek, imageID, err)
	}
	return nil
}
