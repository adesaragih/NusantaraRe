package layanan

// Penyimpanan token berkas - butir an, A2.
//
// `[data DBA]` `POOLDATA.GCP_IMAGE` - `APPNAME VARCHAR2(20)`,
// `KODEAKSES VARCHAR2(100)`, `USERINPUT VARCHAR2(50)`, `INPUTDATE DATE`.
//
// ⚠️ TULISAN KE TABEL WARISAN YANG DISENGAJA dan DISETUJUI: §6 brief
// mengizinkan `GCP_IMAGE` selain baris datar. Dicatat, bukan disembunyikan.
//
// ⛔ Nilai tokennya TIDAK PERNAH masuk pesan galat maupun log - ia kredensial.
//
// Dibaca sesudah: kategoridokumen.go.

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
)

// PenyimpanToken membaca dan menulis token penyimpanan di `GCP_IMAGE`.
//
// Refactor bentuk B (30-09-2026): dulu menumpang di `PohonKlaim`; dipindah
// apa adanya karena Treaty Contract Out ikut memakainya.
type PenyimpanToken struct {
	db *db.DB
}

// NewPenyimpanToken membuat penyimpan token.
func NewPenyimpanToken(db *db.DB) *PenyimpanToken { return &PenyimpanToken{db: db} }

// TokenBerlaku mencari token yang masih berlaku bagi sebuah aplikasi.
//
// `[data DBA]` procedure memilih `KODEAKSES` terbaru dengan
// `INPUTDATE > SYSDATE` - kedaluwarsa di MASA DEPAN berarti masih berlaku.
// Mengembalikan teks kosong bila tidak ada.
func (r *PenyimpanToken) TokenBerlaku(ctx context.Context, tx *db.Tx,
	appName string, saat time.Time) (string, error) {

	tabel, err := r.db.Qualify("GCP_IMAGE")
	if err != nil {
		return "", err
	}
	q := fmt.Sprintf(`SELECT KODEAKSES FROM %s
		 WHERE APPNAME = :1 AND INPUTDATE > :2
		 ORDER BY INPUTDATE DESC
		 FETCH FIRST 1 ROWS ONLY`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var kode sql.NullString
	err = tx.QueryRowContext(ctx, q, appName, saat).Scan(&kode)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		// ⛔ Galat driver TIDAK dibungkus apa adanya: pesannya dapat memuat
		// nilai kolom, dan kolom itu kredensial.
		return "", fmt.Errorf("repository: membaca token penyimpanan untuk %q", appName)
	}
	return kode.String, nil
}

// SimpanToken menulis token baru beserta kedaluwarsanya.
func (r *PenyimpanToken) SimpanToken(ctx context.Context, tx *db.Tx,
	appName, token, pengguna string, kedaluwarsa time.Time) error {

	tabel, err := r.db.Qualify("GCP_IMAGE")
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`INSERT INTO %s
		(APPNAME, KODEAKSES, USERINPUT, INPUTDATE)
		VALUES (:1,:2,:3,:4)`, tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, appName, token, pengguna, kedaluwarsa)
	if err != nil {
		return fmt.Errorf("repository: menyimpan token penyimpanan untuk %q", appName)
	}
	return db.PastikanSatuBaris(hasil, "penyimpanan token")
}
