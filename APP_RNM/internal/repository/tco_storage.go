package repository

// Pembaca rantai penyimpanan lampiran - OQ-TCO-08 Treaty Contract Out.
//
// Untuk apa berkas ini: dua bacaan yang pelaksana penyimpanan nyata perlukan,
// di tabel warisan yang DIBACA SAJA dari berkas ini:
//
//	`T_FOLDER_IMAGE.APPNAME` - `Claim Life/RDBList/GetAppName_SQL.xml`
//	(`SELECT APPNAME ... FROM POOLDATA.T_FOLDER_IMAGE`), `App` panggilan
//	penyimpanan dan `APPNAME` token.
//	`GCP_IMAGE` token berlaku BESERTA kedaluwarsanya - `TokenBerlaku` bersama
//	hanya mengembalikan tokennya, sedangkan cache token modul (AC 60) perlu
//	tahu kapan ia habis. Penerbitan token baru tetap lewat `SimpanToken` bersama.
//
// ⛔ Nilai token TIDAK PERNAH masuk pesan galat maupun log.
//
// Dibaca sesudah: tokenstorage.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// MasterFolderImageTCO - tabel warisan asal `APPNAME` penyimpanan.
const MasterFolderImageTCO = "T_FOLDER_IMAGE"

// ErrAppStorageKosongTCO - `T_FOLDER_IMAGE` tidak memberi `APPNAME`.
var ErrAppStorageKosongTCO = errors.New("repository: APPNAME penyimpanan tidak ada di T_FOLDER_IMAGE")

func sqlAppStorageTCO(tabel string) string {
	return fmt.Sprintf(`SELECT APPNAME FROM %s WHERE APPNAME IS NOT NULL FETCH FIRST 1 ROWS ONLY`, tabel)
}

func sqlTokenStorageBerlakuTCO(tabel string) string {
	return fmt.Sprintf(`SELECT KODEAKSES, INPUTDATE FROM %s
		 WHERE APPNAME = :1 AND INPUTDATE > :2
		 ORDER BY INPUTDATE DESC
		 FETCH FIRST 1 ROWS ONLY`, tabel)
}

// AppStorageTCO membaca `APPNAME` penyimpanan saat jalan.
func (d *DB) AppStorageTCO(ctx context.Context) (string, error) {
	tabel, err := d.Qualify(MasterFolderImageTCO)
	if err != nil {
		return "", err
	}
	q := sqlAppStorageTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var app sql.NullString
	err = d.bacaTCO(ctx).QueryRowContext(ctx, q).Scan(&app)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && app.String == "") {
		return "", ErrAppStorageKosongTCO
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca APPNAME penyimpanan: %w", err)
	}
	return app.String, nil
}

// TokenStorageBerlakuTCO mencari token berlaku dan kedaluwarsanya.
// Mengembalikan token kosong bila tidak ada.
func (d *DB) TokenStorageBerlakuTCO(ctx context.Context, tx *Tx, appName string, saat time.Time) (string, time.Time, error) {
	if tx == nil {
		return "", time.Time{}, errors.New("repository: token penyimpanan menuntut transaksi")
	}
	tabel, err := d.Qualify("GCP_IMAGE")
	if err != nil {
		return "", time.Time{}, err
	}
	q := sqlTokenStorageBerlakuTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", time.Time{}, err
	}
	var (
		kode   sql.NullString
		sampai sql.NullTime
	)
	err = tx.tx.QueryRowContext(ctx, q, appName, saat).Scan(&kode, &sampai)
	if errors.Is(err, sql.ErrNoRows) {
		return "", time.Time{}, nil
	}
	if err != nil {
		// ⛔ Galat driver tidak diteruskan: pesannya dapat memuat nilai kolom kredensial.
		return "", time.Time{}, fmt.Errorf("repository: membaca token penyimpanan untuk %q", appName)
	}
	return kode.String, sampai.Time, nil
}
