package repository

// Tabel bisnis untuk daftar Class Of Business (tiket 28). Skema HANYA dari DDL
// `D:\migrasi\RNM\DDL\BUSINESS.txt` `[terverifikasi]`: ID, NOTE, BUSINESSGROUPID
// masing-masing VARCHAR2(4000 BYTE), tanpa NOT NULL. Saringan dari RD Pega
// BrowseBusiness_RD filter C `.BusinessGroupID = Param.Group`; NOTE yang ditampilkan
// (pyDisplayProperty `.Note` di section InputLossRecord_Sec, bukan di RD). `NOTE IS NOT
// NULL` dari brief, tidak ada di RD (A78).
//
// Hanya membaca dua kolom (+ satu kolom saringan); tidak ada `SELECT *`, nol tulisan.
// groupBusinessId selalu parameter terikat.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TabelBisnis - nama tabel, ejaan DDL.
const TabelBisnis = "BUSINESS"

// PembacaKelasBisnis - sumber pilihan Class Of Business.
type PembacaKelasBisnis interface {
	// KelasBisnis - SEMUA baris ber-NOTE milik satu group business, tanpa paging.
	KelasBisnis(ctx context.Context, groupBusinessID string) ([]models.KelasBisnis, error)
}

// KelasBisnisOracle - PembacaKelasBisnis atas Oracle.
type KelasBisnisOracle struct{ db *db.DB }

// NewKelasBisnisOracle merakit pembaca tabel bisnis.
func NewKelasBisnisOracle(d *db.DB) *KelasBisnisOracle { return &KelasBisnisOracle{db: d} }

// sqlKelasBisnis - A74: urut NOTE (tangkapan layar Pega alfabetis), BUKAN `.ID` DESC
// RD; A75: ID sebagai pemutus seri supaya urutan deterministik.
func sqlKelasBisnis(tabel string) string {
	return "SELECT ID, NOTE FROM " + tabel + " WHERE BUSINESSGROUPID = :1 AND NOTE IS NOT NULL ORDER BY NOTE, ID"
}

// KelasBisnis - lihat PembacaKelasBisnis.
func (r *KelasBisnisOracle) KelasBisnis(ctx context.Context, groupBusinessID string) ([]models.KelasBisnis, error) {
	q, err := r.db.Qualify(TabelBisnis)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlKelasBisnis(q), groupBusinessID)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelBisnis, err)
	}
	defer baris.Close()
	hasil := []models.KelasBisnis{}
	for baris.Next() {
		var id, note sql.NullString
		if err := baris.Scan(&id, &note); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
		}
		hasil = append(hasil, models.KelasBisnis{ID: id.String, Note: note.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
	}
	return hasil, nil
}
