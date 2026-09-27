package repository

// Perekam jejak audit - butir am, A2.
//
// Untuk apa berkas ini: menulis SIAPA dan KAPAN untuk setiap transisi status
// dan setiap jalur balik (ADR-U-0007).
//
// ⛔ Ia dipanggil DI DALAM transaksi pemanggilnya, bukan sesudahnya. Jejak
// yang ditulis di transaksi terpisah dapat hilang sendirian, dan transisi
// tanpa jejak persis yang ADR-U-0007 larang.
//
// Dibaca sesudah: klaimlife.go.

import (
	"context"
	"fmt"
	"time"
)

// sqlSisipJejak merakit pernyataannya.
//
// Dipisah supaya bentuknya dapat diuji tanpa Oracle - kolom yang lupa di-bind
// tidak terlihat dari daftar nama mana pun.
func sqlSisipJejak(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
		(ID, ADJUSTMENT_ID, KLAIM_ID, DARI, KE, AKUN_ID, WAKTU)
		VALUES (:1,:2,:3,:4,:5,:6,:7)`, tabel)
}

// SisipJejak menulis satu catatan jejak.
//
// ⚠️ `dari` dan `ke` memuat DUA kosakata: kode status pada transisi baris, dan
// nama peran pada jalur balik tahap. Disengaja - keduanya "keadaan sebelum"
// dan "keadaan sesudah".
func (r *KlaimLife) SisipJejak(ctx context.Context, tx *Tx,
	adjustmentID, klaimID, dari, ke, akunID string, waktu time.Time) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_JEJAK")
	if err != nil {
		return err
	}
	id, err := NewPohonKlaim(r.db).nomorBerikut(ctx, tx, "SEQ_CLAIMLF_JEJAK")
	if err != nil {
		return err
	}
	q := sqlSisipJejak(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id,
		kosongJadiNil(adjustmentID), kosongJadiNil(klaimID),
		kosongJadiNil(dari), kosongJadiNil(ke), akunID, waktu)
	if err != nil {
		return fmt.Errorf("repository: merekam jejak: %w", err)
	}
	return pastikanSatuBaris(hasil, "perekaman jejak")
}
