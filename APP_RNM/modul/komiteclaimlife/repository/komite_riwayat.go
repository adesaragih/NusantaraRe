package repository

// Pembaca jejak eskalasi kasus Komite - tiket 09 Komite Claim Life.
//
// Riwayat tangga dibaca dari `T_KOMITE_KOMITELIST` (`InboxKomite.Kasus`,
// urut `KOMITE_URUT`); yang TIDAK ada di sana adalah eskalasi - siapa
// memindahkan dan kapan. Itu tinggal di jejak `T_CLAIMLF_JEJAK` (tiket 03).
// Berkas ini membacanya kembali. Nol tulisan.

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
)

// JejakKomite adalah satu catatan jejak.
type JejakKomite struct {
	Dari   string
	Ke     string
	AkunID string
	Waktu  time.Time
	// Komentar - `KOMENTAR` (migrasi 021): komentar keputusan, dan komentar
	// ASLI tingkat yang ditimpa langkah 5.1 (OQ-K-05, GILIRAN-17).
	Komentar string
}

// sqlJejakEskalasi membaca jejak satu baris adjustment yang KE-nya berpola.
//
// ⛔ Berkunci `ADJUSTMENT_ID` (ADR-0011: riwayat melekat pada baris), dan
// polanya menyebut kasus komitenya - dua kasus atas baris yang sama (jalur
// baru sesudah Tolak) tidak saling meminjam eskalasi.
func sqlJejakEskalasi(tabel string) string {
	return fmt.Sprintf(`SELECT DARI, KE, AKUN_ID, WAKTU, KOMENTAR
	  FROM %s WHERE ADJUSTMENT_ID = :1 AND KE LIKE :2 ORDER BY WAKTU, ID`, tabel)
}

// JejakEskalasi membaca jejak eskalasi satu baris adjustment.
func (r *InboxKomite) JejakEskalasi(ctx context.Context, adjID, polaKe string) ([]JejakKomite, error) {
	tabel, err := r.db.Qualify("T_CLAIMLF_JEJAK")
	if err != nil {
		return nil, err
	}
	q := sqlJejakEskalasi(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, adjID, polaKe)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca jejak eskalasi: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []JejakKomite
	for rows.Next() {
		var j JejakKomite
		var dari, ke, akun, kom sql.NullString
		if err := rows.Scan(&dari, &ke, &akun, &j.Waktu, &kom); err != nil {
			return nil, fmt.Errorf("repository: memindai jejak eskalasi: %w", err)
		}
		j.Dari, j.Ke, j.AkunID, j.Komentar = dari.String, ke.String, akun.String, kom.String
		out = append(out, j)
	}
	return out, rows.Err()
}
