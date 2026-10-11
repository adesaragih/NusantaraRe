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
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TabelBisnis - nama tabel, ejaan DDL.
const TabelBisnis = "BUSINESS"

// PembacaKelasBisnis - sumber pilihan Class Of Business.
type PembacaKelasBisnis interface {
	// KelasBisnis - SEMUA baris ber-NOTE milik satu group business, tanpa paging.
	KelasBisnis(ctx context.Context, groupBusinessID string) ([]models.KelasBisnis, error)
	// SaranKelasBisnis - tiket 55 Class Of Business Loss Record: BrowseBusiness_RD dengan ketiga saringan (ID, Note,
	// Group) dikirim KOSONG oleh `Section\InputLossRecord_Sec.xml` -> semua grup; NOTE mengandung kata (tidak peka huruf).
	SaranKelasBisnis(ctx context.Context, kata string) ([]models.SaranAkumulasi, error)
	// NoteKelasBisnis - tiket 55 (Save Loss Record): untuk tiap kunci (NOTE dipangkas + huruf besar) ejaan BUSINESS.NOTE
	// tersimpan; kunci tanpa baris tidak ada di peta. Beberapa baris berkunci sama -> ID terkecil.
	NoteKelasBisnis(ctx context.Context, kunci []string) (map[string]string, error)
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

// sqlSaranKelasBisnis - urut ID DESC (pySortOrder 1 BrowseBusiness_RD), NOTE kosong dibuang (nilai yang ditulis = .Note),
// kata opsional, batas saran 50 (pola A124; RD 500). Dengan kata: NOTE yang SAMA PERSIS (dipangkas, tanpa beda huruf; :2 =
// UPPER(TRIM(kata))) lebih dulu - layar memutuskan "ada di daftar" dari jawaban ini, jadi yang persis sama tidak boleh
// terpotong batas 50 (tiket 55 K55-10).
func sqlSaranKelasBisnis(tabel string, denganKata bool) string {
	if !denganKata {
		return "SELECT ID, NOTE FROM " + tabel + " WHERE NOTE IS NOT NULL ORDER BY ID DESC FETCH FIRST :1 ROWS ONLY"
	}
	return "SELECT ID, NOTE FROM " + tabel + ` WHERE NOTE IS NOT NULL AND UPPER(NOTE) LIKE :1 ESCAPE '\'` +
		" ORDER BY CASE WHEN UPPER(TRIM(NOTE)) = :2 THEN 0 ELSE 1 END, ID DESC FETCH FIRST :3 ROWS ONLY"
}

// SaranKelasBisnis - lihat PembacaKelasBisnis.
func (r *KelasBisnisOracle) SaranKelasBisnis(ctx context.Context, kata string) ([]models.SaranAkumulasi, error) {
	q, err := r.db.Qualify(TabelBisnis)
	if err != nil {
		return nil, err
	}
	sqlq, arg := sqlSaranKelasBisnis(q, false), []any{BatasSaranAkumulasi}
	if pola := PolaCari(kata); pola != "" {
		sqlq, arg = sqlSaranKelasBisnis(q, true), []any{pola, strings.ToUpper(strings.TrimSpace(kata)), BatasSaranAkumulasi}
	}
	baris, err := r.db.QueryContext(ctx, sqlq, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelBisnis, err)
	}
	defer baris.Close()
	hasil := []models.SaranAkumulasi{}
	for baris.Next() {
		var id, note sql.NullString
		if err := baris.Scan(&id, &note); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
		}
		hasil = append(hasil, models.SaranAkumulasi{ID: id.String, Label: note.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
	}
	return hasil, nil
}

// sqlNoteKelasBisnis - :1..:n kunci UPPER(TRIM(NOTE)); urut ID (baris pertama per kunci dipakai).
func sqlNoteKelasBisnis(tabel string, n int) string {
	b := make([]string, n)
	for i := range b {
		b[i] = ":" + strconv.Itoa(i+1)
	}
	return "SELECT NOTE FROM " + tabel + " WHERE UPPER(TRIM(NOTE)) IN (" + strings.Join(b, ", ") + ") ORDER BY ID"
}

// NoteKelasBisnis - lihat PembacaKelasBisnis.
func (r *KelasBisnisOracle) NoteKelasBisnis(ctx context.Context, kunci []string) (map[string]string, error) {
	hasil := map[string]string{}
	if len(kunci) == 0 {
		return hasil, nil
	}
	q, err := r.db.Qualify(TabelBisnis)
	if err != nil {
		return nil, err
	}
	arg := make([]any, len(kunci))
	for i, k := range kunci {
		arg[i] = k
	}
	baris, err := r.db.QueryContext(ctx, sqlNoteKelasBisnis(q, len(kunci)), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelBisnis, err)
	}
	defer baris.Close()
	for baris.Next() {
		var note sql.NullString
		if err := baris.Scan(&note); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
		}
		k := strings.ToUpper(strings.TrimSpace(note.String))
		if _, ada := hasil[k]; !ada {
			hasil[k] = note.String
		}
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelBisnis, err)
	}
	return hasil, nil
}
