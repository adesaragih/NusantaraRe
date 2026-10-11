package repository

// Cacah bahan kesiapan Confirm (keputusan work owner 06-10-2026) - satu kueri,
// baca-saja: peserta, peserta Type QR ber-RETROCEDED_SHARE > 0 tanpa spreading
// retro, dan baris rekap summary.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
)

// CacahKesiapan - jumlah yang diperiksa kesiapan Confirm.
type CacahKesiapan struct {
	Peserta        int
	TanpaSpreading int
	Summary        int
}

// KesiapanPolis membaca cacah kesiapan Confirm.
type KesiapanPolis struct{ db *db.DB }

// NewKesiapanPolis menyusunnya.
func NewKesiapanPolis(db *db.DB) *KesiapanPolis { return &KesiapanPolis{db: db} }

func sqlCacahKesiapan(detail, spreading, summary string) string {
	return fmt.Sprintf(`SELECT
	   (SELECT COUNT(*) FROM %s d WHERE d.PREMIUM_LIST_ID = :1),
	   (SELECT COUNT(*) FROM %s d WHERE d.PREMIUM_LIST_ID = :2 AND NVL(d.RETROCEDED_SHARE, 0) > 0
	      AND NOT EXISTS (SELECT 1 FROM %s s WHERE s.DETAIL_ID = d.ID)),
	   (SELECT COUNT(*) FROM %s m WHERE m.PREMIUM_LIST_ID = :3)
	   FROM DUAL`, detail, detail, spreading, summary)
}

// Cacah membaca ketiga jumlah satu polis.
func (r *KesiapanPolis) Cacah(ctx context.Context, polisID string) (CacahKesiapan, error) {
	detail, err := r.db.Qualify("T_PREMIUM_LIST_DETAIL")
	if err != nil {
		return CacahKesiapan{}, err
	}
	spreading, err := r.db.Qualify(TabelSpreading)
	if err != nil {
		return CacahKesiapan{}, err
	}
	summary, err := r.db.Qualify("T_PREMIUM_LIST_SUMMARY")
	if err != nil {
		return CacahKesiapan{}, err
	}
	q := sqlCacahKesiapan(detail, spreading, summary)
	if err := db.PeriksaSQL(q); err != nil {
		return CacahKesiapan{}, err
	}
	var c CacahKesiapan
	if err := r.db.QueryRowContext(ctx, q, polisID, polisID, polisID).Scan(&c.Peserta, &c.TanpaSpreading, &c.Summary); err != nil {
		return CacahKesiapan{}, fmt.Errorf("repository: membaca kesiapan Confirm: %w", err)
	}
	return c, nil
}
