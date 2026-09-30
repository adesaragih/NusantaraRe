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

	"nusantarare/inti/backend/db"
)

// MasterFolderImageTCO - tabel warisan asal `APPNAME` penyimpanan.
const MasterFolderImageTCO = "T_FOLDER_IMAGE"

// ErrAppStorageKosongTCO - `T_FOLDER_IMAGE` tidak memberi `APPNAME`.
var ErrAppStorageKosongTCO = errors.New("repository: storage APPNAME not in T_FOLDER_IMAGE")

func sqlAppStorageTCO(tabel string) string {
	return fmt.Sprintf(`SELECT APPNAME FROM %s WHERE APPNAME IS NOT NULL FETCH FIRST 1 ROWS ONLY`, tabel)
}

// sqlTokenStorageBerlakuTCO - sisa umur (detik) dihitung DI ORACLE terhadap
// waktu yang di-bind, bukan dengan membaca `INPUTDATE` ke Go: DATE tidak
// membawa zona, dan membandingkannya dengan jam aplikasi dapat meleset sebesar
// selisih zona. Tulis (`SimpanToken`) dan baca melewati konversi bind yang sama.
func sqlTokenStorageBerlakuTCO(tabel string) string {
	return fmt.Sprintf(`SELECT KODEAKSES, ROUND((CAST(INPUTDATE AS DATE) - CAST(:1 AS DATE)) * 86400) FROM %s
		 WHERE APPNAME = :2 AND INPUTDATE > :3
		 ORDER BY INPUTDATE DESC
		 FETCH FIRST 1 ROWS ONLY`, tabel)
}

// AppStorageTCO membaca `APPNAME` penyimpanan saat jalan.
func AppStorageTCO(ctx context.Context, d *db.DB) (string, error) {
	tabel, err := d.Qualify(MasterFolderImageTCO)
	if err != nil {
		return "", err
	}
	q := sqlAppStorageTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var app sql.NullString
	err = bacaTCO(ctx, d).QueryRowContext(ctx, q).Scan(&app)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && app.String == "") {
		return "", ErrAppStorageKosongTCO
	}
	if err != nil {
		return "", fmt.Errorf("repository: reading storage APPNAME: %w", err)
	}
	return app.String, nil
}

// TokenStorageBerlakuTCO mencari token yang masih berlaku LEBIH DARI
// `sisaMinimum` sesudah `saat`, beserta sisa umurnya. Token kosong = tidak ada.
func TokenStorageBerlakuTCO(ctx context.Context, d *db.DB, tx *db.Tx, appName string, saat time.Time,
	sisaMinimum time.Duration) (string, time.Duration, error) {
	if tx == nil {
		return "", 0, errors.New("repository: storage token requires a transaction")
	}
	tabel, err := d.Qualify("GCP_IMAGE")
	if err != nil {
		return "", 0, err
	}
	q := sqlTokenStorageBerlakuTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", 0, err
	}
	var (
		kode sql.NullString
		sisa sql.NullFloat64
	)
	err = tx.QueryRowContext(ctx, q, saat, appName, saat.Add(sisaMinimum)).Scan(&kode, &sisa)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, nil
	}
	if err != nil {
		// ⛔ Galat driver tidak diteruskan: pesannya dapat memuat nilai kolom kredensial.
		return "", 0, fmt.Errorf("repository: reading storage token for %q", appName)
	}
	return kode.String, time.Duration(sisa.Float64) * time.Second, nil
}
