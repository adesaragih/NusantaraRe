package repository

// Isian Claim RNM baris adjustment sebelum Save to RNM - keputusan work owner 11-10-2026. Dipisah dari
// isianoutstanding.go: berkas ini TIDAK menyentuh tabel peserta polis M_LIFE_PREMIUM_DETAIL.

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
)

// sqlSetelClaimRNM - isian Claim RNM (`.CLAIM_GROSS` AdjustmentDetail b3012 = kolom CLAIM_AMOUNT) satu baris adjustment,
// keputusan work owner 11-10-2026: HANYA sebelum Save to RNM. Syaratnya di DALAM pernyataan: baris belum berstatus
// (STS_REJECT NULL), milik peserta aktif klaim itu, dan klaim belum bernomor - Save to RNM yang menyela tidak diikuti.
// Pengikatan menurut URUTAN kemunculan.
//
// Sejak 11-10-2026 juga data bank (Name of Bank, Branch of Bank, Account No) - isian TEKS sampai master BANKACCOUNT
// diketahui (keputusan work owner).
func sqlSetelClaimRNM(adj, pes, hdr string) string {
	// 11-10-2026 "semua read only kecuali name of bank": BRANCH_OF_BANK dan ACCOUNT_NO TIDAK ditulis lagi.
	return fmt.Sprintf(`UPDATE %s a SET a.CLAIM_AMOUNT = :1, a.NAME_OF_BANK = :2
	 WHERE a.ID = :3 AND a.STS_REJECT IS NULL
	   AND a.PREMIUM_LIST_DETAIL_ID IN (SELECT p.ID FROM %s p WHERE p.CLAIM_ID = :4 AND p.STS_HAPUS IS NULL)
	   AND EXISTS (SELECT 1 FROM %s h WHERE h.ID = :5 AND h.CLAIM_NO IS NULL)`, adj, pes, hdr)
}

// IsianBaris - Claim RNM dan Name of Bank satu baris adjustment (Branch of Bank / Account No baca saja, 11-10-2026).
type IsianBaris struct {
	Jumlah   string
	NamaBank string
}

// SetelClaimRNM menulis Claim RNM satu baris di transaksi pemanggil; false bila syaratnya tidak lagi terpenuhi.
func (r *KlaimLife) SetelClaimRNM(ctx context.Context, tx *db.Tx, klaimID, adjID string, is IsianBaris) (bool, error) {
	adj, err := r.db.Qualify("T_CLAIMLF_ADJUSTMENT")
	if err != nil {
		return false, err
	}
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return false, err
	}
	hdr, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return false, err
	}
	q := sqlSetelClaimRNM(adj, pes, hdr)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	// Claim RNM kosong = NULL (boleh pada Save grid; wajib ditegakkan Save to RNM, 11-10-2026).
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(is.Jumlah), db.KosongJadiNil(is.NamaBank), adjID, klaimID, klaimID)
	if err != nil {
		return false, fmt.Errorf("repository: menulis Claim RNM %s: %w", adjID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("repository: mencacah Claim RNM: %w", err)
	}
	return n == 1, nil
}

// TanggalPeserta adalah kelima tanggal yang diisi pemakai di rincian peserta sebelum Save to RNM (keputusan work owner
// 11-10-2026). nil = NULL (dikosongkan).
type TanggalPeserta struct {
	Kejadian       *time.Time // DATE_OF_LOSS
	TerimaKlaim    *time.Time // CLAIM_RECEIVED_DATE
	DokumenLengkap *time.Time // COMPLETE_DATE
	MaksData       *time.Time // MAX_DATA_RECEIVED (migrasi 028)
	Konfirmasi     *time.Time // CONFIRMATION_DATE
}

// sqlSetelTanggalPeserta - kelima tanggal satu peserta dalam SATU pernyataan (satu tombol Save). Hanya selama klaim
// belum bernomor (syarat di DALAM pernyataan, OQ-M1). Pengikatan menurut URUTAN kemunculan.
func sqlSetelTanggalPeserta(pes, hdr string) string {
	return fmt.Sprintf(`UPDATE %s p
	   SET p.DATE_OF_LOSS        = TO_DATE(:1, 'YYYY-MM-DD HH24:MI:SS'),
	       p.CLAIM_RECEIVED_DATE = TO_DATE(:2, 'YYYY-MM-DD HH24:MI:SS'),
	       p.COMPLETE_DATE       = TO_DATE(:3, 'YYYY-MM-DD HH24:MI:SS'),
	       p.MAX_DATA_RECEIVED   = TO_DATE(:4, 'YYYY-MM-DD HH24:MI:SS'),
	       p.CONFIRMATION_DATE   = TO_DATE(:5, 'YYYY-MM-DD HH24:MI:SS')
	 WHERE p.ID = :6 AND p.CLAIM_ID = :7 AND p.STS_HAPUS IS NULL
	   AND EXISTS (SELECT 1 FROM %s h WHERE h.ID = :8 AND h.CLAIM_NO IS NULL)`, pes, hdr)
}

// SetelTanggalPeserta menulis kelima tanggal di transaksi pemanggil; false bila klaim sudah bernomor / peserta tak ada.
func (r *KlaimLife) SetelTanggalPeserta(ctx context.Context, tx *db.Tx, klaimID, pesertaID string,
	t TanggalPeserta) (bool, error) {

	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return false, err
	}
	hdr, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return false, err
	}
	q := sqlSetelTanggalPeserta(pes, hdr)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	hasil, err := tx.ExecContext(ctx, q, argTanggal(t.Kejadian), argTanggal(t.TerimaKlaim), argTanggal(t.DokumenLengkap),
		argTanggal(t.MaksData), argTanggal(t.Konfirmasi), pesertaID, klaimID, klaimID)
	if err != nil {
		return false, fmt.Errorf("repository: menulis tanggal peserta %s: %w", pesertaID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("repository: mencacah tanggal peserta: %w", err)
	}
	return n == 1, nil
}
