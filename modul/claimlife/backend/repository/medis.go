package repository

// Tulis tombol Save layar Medical Check - 11-10-2026: RECOMMENDATION per peserta (migrasi 003, VARCHAR2(1000)).
//
// Dibaca sesudah: isianoutstanding.go (SetelCatatanAdmin - pola yang sama).

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
)

// sqlSetelRekomendasi - RECOMMENDATION satu peserta aktif. Dibekukan bila peserta sudah diputus
// (`ClaimLifeDetailGCNM` b7323 pyDisabledWhen `.STS_REJECT=='1' || =='2'`) - syaratnya di DALAM pernyataan.
func sqlSetelRekomendasi(pes string) string {
	return fmt.Sprintf(`UPDATE %s SET RECOMMENDATION = :1
	 WHERE ID = :2 AND CLAIM_ID = :3 AND STS_HAPUS IS NULL AND (STS_REJECT IS NULL OR STS_REJECT NOT IN ('1', '2'))`, pes)
}

// SetelRekomendasi menulis RECOMMENDATION satu peserta di transaksi pemanggil; false bila peserta sudah diputus
// (atau bukan peserta aktif klaim ini).
func (r *KlaimLife) SetelRekomendasi(ctx context.Context, tx *db.Tx, klaimID, pesertaID, teks string) (bool, error) {
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return false, err
	}
	q := sqlSetelRekomendasi(pes)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(teks), pesertaID, klaimID)
	if err != nil {
		return false, fmt.Errorf("repository: menulis Recommendation peserta %s: %w", pesertaID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("repository: mencacah Recommendation: %w", err)
	}
	return n == 1, nil
}
