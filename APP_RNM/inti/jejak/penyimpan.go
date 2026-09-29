package jejak

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

	"nusantarare/inti/db"
)

// PenyimpanJejak menulis `T_CLAIMLF_JEJAK`.
//
// Refactor bentuk B (30-09-2026): dulu menumpang di `KlaimLife`; dipindah
// apa adanya karena PremiumList dan Komite ikut merekam jejak.
type Penyimpan struct {
	db *db.DB
}

// NewPenyimpanJejak membuat penyimpan jejak.
func NewPenyimpan(db *db.DB) *Penyimpan { return &Penyimpan{db: db} }

// sqlSisipJejak merakit pernyataannya.
//
// Dipisah supaya bentuknya dapat diuji tanpa Oracle - kolom yang lupa di-bind
// tidak terlihat dari daftar nama mana pun.
func sqlSisipJejak(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
		(ID, ADJUSTMENT_ID, KLAIM_ID, DARI, KE, AKUN_ID, WAKTU, KOMENTAR)
		VALUES (:1,:2,:3,:4,:5,:6,:7,:8)`, tabel)
}

// SisipJejak menulis satu catatan jejak.
//
// ⚠️ `dari` dan `ke` memuat DUA kosakata: kode status pada transisi baris, dan
// nama peran pada jalur balik tahap. Disengaja - keduanya "keadaan sebelum"
// dan "keadaan sesudah". `komentar` (migrasi 021, OQ-M5) kosong = NULL.
func (r *Penyimpan) SisipJejak(ctx context.Context, tx *db.Tx,
	adjustmentID, klaimID, dari, ke, akunID string, waktu time.Time, komentar string) error {

	tabel, err := r.db.Qualify("T_CLAIMLF_JEJAK")
	if err != nil {
		return err
	}
	id, err := r.db.NomorBerikut(ctx, tx, "SEQ_CLAIMLF_JEJAK")
	if err != nil {
		return err
	}
	q := sqlSisipJejak(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id,
		db.KosongJadiNil(adjustmentID), db.KosongJadiNil(klaimID),
		db.KosongJadiNil(dari), db.KosongJadiNil(ke), akunID, waktu, db.KosongJadiNil(komentar))
	if err != nil {
		return fmt.Errorf("repository: merekam jejak: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "perekaman jejak")
}
