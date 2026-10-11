package repository

// IDPEGA peserta klaim dari baris SUMBERNYA - keputusan work owner 10-10-2026 ("ambil ke tabel T_PREMIUM_LIST
// berdasarkan idpega peserta nya saja").
//
// Sejak migrasi 025 IDPEGA DISIMPAN di T_CLAIMLF_PREMIUMLIST_DETAIL; fungsi ini CADANGAN untuk baris yang kolomnya
// masih NULL (di luar jangkauan backfill 025), menurunkannya dari baris M_LIFE_PREMIUM_DETAIL yang disalin (SOURCE_ID),
// berkunci PL_NUMBER + CERTIFICATE_NO (index) + ID - pola yang sama dengan gandawarisan.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// sqlIdPegaSumber - IDPEGA satu baris sumber peserta.
func sqlIdPegaSumber(sumber string) string {
	return fmt.Sprintf(`SELECT IDPEGA FROM %s
	  WHERE PL_NUMBER = :1 AND CERTIFICATE_NO = :2 AND ID = :3
	  FETCH FIRST 1 ROWS ONLY`, sumber)
}

// IdPegaSumber - IDPEGA baris sumber k; "" bila barisnya tidak ada atau IDPEGA-nya kosong (dilaporkan pemanggil).
func (r *KlaimLife) IdPegaSumber(ctx context.Context, k KunciPesertaSumber) (string, error) {
	if strings.TrimSpace(k.SumberID) == "" {
		return "", nil
	}
	sumber, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return "", err
	}
	q := sqlIdPegaSumber(sumber)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	err = r.db.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca IDPEGA peserta sumber: %w", err)
	}
	return strings.TrimSpace(id.String), nil
}
