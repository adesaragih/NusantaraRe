package repository

// Objek penyimpanan dokumen klaim di tabel BERSAMA `T_STORAGE_IMAGE` (keputusan work owner 11-10-2026: "untuk
// penyimpanan document samakan dengan penyimpanan document yang digunakan oleh product name" - dicatat di T_STORAGE_IMAGE
// seperti Master Product Name Life dan Pega). Bentuk SQL disalin dari masterproductnamelife repository/mpnl_lampiran.go
// (tanpa impor modul itu), yang meniru rule warisan:
//
//	Insert_T_Storage_SQL b85   INSERT T_STORAGE_IMAGE (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE='standard')
//	GetLinkStorage_SQL b85     SELECT URLPUBLIC, APPFOLDER, EXPDATE, APPNAME, FILENAME … WHERE imageid = …
//	Update_T_Storage_SQL b85   UPDATE … SET URLPUBLIC, APPFOLDER, EXPDATE, TANGGAL_UPLOAD WHERE imageid = …
//	DeleteStorage_SQL b85      DELETE … WHERE imageid = …
//	GetAppName_SQL b58         SELECT APPNAME FROM T_FOLDER_IMAGE
//
// ⛔ NOL COMMIT (rule warisan ber-commit tidak dipakai): transaksi milik pemanggil. ⚠️ T_CLAIMLF_STORAGE (migrasi 019,
// kartu stub lokal) DIPENSIUNKAN untuk dokumen baru dan hanya DIBACA untuk dokumen lama.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// Tabel bersama penyimpanan berkas (warisan, lintas aplikasi; APPNAME membedakan aplikasinya).
const (
	TabelObjekPenyimpanan  = "T_STORAGE_IMAGE"
	TabelFolderPenyimpanan = "T_FOLDER_IMAGE"
)

// ObjekPenyimpanan - satu baris T_STORAGE_IMAGE. Exp bentuk 'DD/MM/YYYY HH24:MI:SS'; TanggalUpload 'MM/DD/YYYY
// HH24:MI:SS' (bentuk `DateTime` jawaban layanan, Update_T_Storage_SQL). DurasiDetik bukan kolom (masukan layanan).
type ObjekPenyimpanan struct {
	ImageID       string
	URLPublic     string
	AppFolder     string
	Exp           string
	FileName      string
	AppName       string
	TanggalUpload string
	DurasiDetik   int
}

func sqlSisipObjek(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s (IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)
		VALUES (:1, :2, :3, TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS'), :5, :6, 'standard')`, tabel)
}

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

// penanyaSQL - transaksi bila ada, selain itu koneksi biasa.
type penanyaSQL interface {
	QueryRowContext(ctx context.Context, q string, arg ...any) *sql.Row
	ExecContext(ctx context.Context, q string, arg ...any) (sql.Result, error)
}

func (r *KlaimLife) penanya(tx *db.Tx) penanyaSQL {
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *KlaimLife) siapkanSQL(tabel string, rakit func(string) string) (string, error) {
	t, err := r.db.Qualify(tabel)
	if err != nil {
		return "", err
	}
	q := rakit(t)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	return q, nil
}

// SisipObjekPenyimpanan mencatat objek yang dijawab layanan penyimpanan (Insert_T_Storage_SQL).
func (r *KlaimLife) SisipObjekPenyimpanan(ctx context.Context, tx *db.Tx, o ObjekPenyimpanan) error {
	q, err := r.siapkanSQL(TabelObjekPenyimpanan, sqlSisipObjek)
	if err != nil {
		return err
	}
	if _, err := r.penanya(tx).ExecContext(ctx, q, o.ImageID, db.KosongJadiNil(o.URLPublic), db.KosongJadiNil(o.AppFolder),
		db.KosongJadiNil(o.Exp), db.KosongJadiNil(o.FileName), db.KosongJadiNil(o.AppName)); err != nil {
		return fmt.Errorf("repository: mencatat objek penyimpanan %s: %w", o.ImageID, err)
	}
	return nil
}

// AmbilObjekPenyimpanan membaca objek (GetLinkStorage_SQL); ada=false bila tidak ada baris.
func (r *KlaimLife) AmbilObjekPenyimpanan(ctx context.Context, tx *db.Tx, imageID string) (ObjekPenyimpanan, bool, error) {
	if strings.TrimSpace(imageID) == "" {
		return ObjekPenyimpanan{}, false, nil
	}
	q, err := r.siapkanSQL(TabelObjekPenyimpanan, sqlAmbilObjek)
	if err != nil {
		return ObjekPenyimpanan{}, false, err
	}
	var url, folder, exp, app, nama sql.NullString
	err = r.penanya(tx).QueryRowContext(ctx, q, imageID).Scan(&url, &folder, &exp, &app, &nama)
	if errors.Is(err, sql.ErrNoRows) {
		return ObjekPenyimpanan{}, false, nil
	}
	if err != nil {
		return ObjekPenyimpanan{}, false, fmt.Errorf("repository: membaca objek penyimpanan %s: %w", imageID, err)
	}
	return ObjekPenyimpanan{ImageID: imageID, URLPublic: url.String, AppFolder: folder.String, Exp: exp.String,
		AppName: app.String, FileName: nama.String}, true, nil
}

// PerbaruiObjekPenyimpanan menulis URL bertanda tangan yang diperbarui (Update_T_Storage_SQL).
func (r *KlaimLife) PerbaruiObjekPenyimpanan(ctx context.Context, tx *db.Tx, o ObjekPenyimpanan) error {
	q, err := r.siapkanSQL(TabelObjekPenyimpanan, sqlPerbaruiObjek)
	if err != nil {
		return err
	}
	if _, err := r.penanya(tx).ExecContext(ctx, q, db.KosongJadiNil(o.URLPublic), db.KosongJadiNil(o.AppFolder),
		db.KosongJadiNil(o.Exp), db.KosongJadiNil(o.TanggalUpload), o.ImageID); err != nil {
		return fmt.Errorf("repository: memperbarui objek penyimpanan %s: %w", o.ImageID, err)
	}
	return nil
}

// HapusObjekPenyimpanan membuang baris objek (DeleteStorage_SQL).
func (r *KlaimLife) HapusObjekPenyimpanan(ctx context.Context, tx *db.Tx, imageID string) error {
	q, err := r.siapkanSQL(TabelObjekPenyimpanan, sqlHapusObjek)
	if err != nil {
		return err
	}
	if _, err := r.penanya(tx).ExecContext(ctx, q, imageID); err != nil {
		return fmt.Errorf("repository: menghapus objek penyimpanan %s: %w", imageID, err)
	}
	return nil
}

// NamaAplikasiPenyimpanan - APPNAME penyimpanan (GetAppName_SQL, T_FOLDER_IMAGE); kosong bila tabelnya kosong.
func (r *KlaimLife) NamaAplikasiPenyimpanan(ctx context.Context, tx *db.Tx) (string, error) {
	q, err := r.siapkanSQL(TabelFolderPenyimpanan, sqlNamaAplikasi)
	if err != nil {
		return "", err
	}
	var app sql.NullString
	err = r.penanya(tx).QueryRowContext(ctx, q).Scan(&app)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca APPNAME penyimpanan: %w", err)
	}
	return strings.TrimSpace(app.String), nil
}
