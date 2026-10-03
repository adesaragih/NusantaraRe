package repository

// Jalur ketentuan proporsional — tiket 35 dan 36.

import (
	"context"
	"fmt"

	"github.com/cockroachdb/apd/v3"
)

// AdaQuotaSharePadaVersi — tiket 36.
//
// Lingkupnya VERSI, bukan layer: tiket 36 berbunyi "versi itu tidak punya baris
// QUOTA_SHARE", dan retensi yang menjadi pengali kapasitas surplus dicatat per
// versi. Gabungnya lewat LAYER sebab DETAIL_PROPORSIONAL menggantung di sana.
func (g *Gudang) AdaQuotaSharePadaVersi(ctx context.Context, idVersi int64) (bool, error) {
	dp, err := g.db.Qualify("DETAIL_PROPORSIONAL")
	if err != nil {
		return false, err
	}
	layer, err := g.db.Qualify("LAYER")
	if err != nil {
		return false, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s d JOIN %s l ON l.ID_LAYER = d.ID_LAYER
		WHERE l.ID_VERSI_KONTRAK = :1 AND d.JENIS_TREATY = 'QUOTA_SHARE'`, dp, layer)
	var n int
	if err := g.db.QueryRowContext(ctx, q, idVersi).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa baris quota share versi %d: %w", idVersi, err)
	}
	return n > 0, nil
}

// CatatKetentuanProporsional menyisipkan satu baris `DETAIL_PROPORSIONAL`.
//
// ⛔ Aturan tepat-satu (tiket 35) sudah ditegakkan di services dan TIDAK
// diulang di sini. Mengulangnya berarti dua tempat yang harus sepakat, dan
// INV-63 melarang rumus yang sama tertulis dua kali.
func (g *Gudang) CatatKetentuanProporsional(ctx context.Context, idVersi, idLayer, idKelompok int64, jenis, persenQS string, lines *int64) error {
	nama, err := g.db.Qualify("DETAIL_PROPORSIONAL")
	if err != nil {
		return err
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	selesai := false
	defer func() {
		if !selesai {
			_ = tx.Rollback()
		}
	}()
	id, err := nomor(ctx, g, tx, "SEQ_TRIN_DETAIL_PROPORSIONAL")
	if err != nil {
		return err
	}
	var koef any
	var skala int64
	if persenQS != "" {
		d, _, errD := apd.NewFromString(persenQS)
		if errD != nil {
			return fmt.Errorf("repository: persen quota share %q bukan desimal: %w", persenQS, errD)
		}
		koef, skala = pecahAngka(d)
	}
	q := fmt.Sprintf(`INSERT INTO %s (ID_DETAIL_PROPORSIONAL, ID_LAYER, ID_KELOMPOK_TREATY,
		JENIS_TREATY, PERSEN_QUOTA_SHARE, JUMLAH_LINES_SURPLUS)
		VALUES (:1, :2, :3, :4,
		        CASE WHEN :5 IS NULL THEN NULL ELSE (TO_NUMBER(:6) / POWER(10, :7)) END, :8)`, nama)
	if _, err := tx.ExecContext(ctx, q, id, idLayer, idKelompok, jenis, koef, koef, skala, lines); err != nil {
		return fmt.Errorf("repository: menyisipkan DETAIL_PROPORSIONAL: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository: menutup transaksi: %w", err)
	}
	selesai = true
	_ = idVersi // lingkup versi diperiksa di services (tiket 36), bukan di sini
	return nil
}
