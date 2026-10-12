package repository

// Tulis tombol Save layar Claim Analis - 12-10-2026: NOTES (CLAIM ANALYST b15194) per peserta (migrasi 003,
// VARCHAR2(1000)).
//
// Dibaca sesudah: medis.go (SetelRekomendasi - pola yang sama).

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
)

// sqlSetelCatatanAnalis - NOTES satu peserta aktif klaim.
func sqlSetelCatatanAnalis(pes string) string {
	return fmt.Sprintf(`UPDATE %s SET NOTES = :1
	 WHERE ID = :2 AND CLAIM_ID = :3 AND STS_HAPUS IS NULL`, pes)
}

// SetelCatatanAnalis menulis NOTES satu peserta di transaksi pemanggil; false bila bukan peserta aktif klaim ini.
func (r *KlaimLife) SetelCatatanAnalis(ctx context.Context, tx *db.Tx, klaimID, pesertaID, teks string) (bool, error) {
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return false, err
	}
	q := sqlSetelCatatanAnalis(pes)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(teks), pesertaID, klaimID)
	if err != nil {
		return false, fmt.Errorf("repository: menulis catatan Claim Analyst peserta %s: %w", pesertaID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("repository: mencacah catatan Claim Analyst: %w", err)
	}
	return n == 1, nil
}
