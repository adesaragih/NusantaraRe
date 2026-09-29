package repository

// Jejak audit Treaty Contract Out - `T_TREATYCO_JEJAK` (ADR-0007; AC 41).
//
// ⛔ Ditulis DI DALAM transaksi pemanggilnya, bukan sesudahnya: jejak yang
// ditulis di transaksi terpisah dapat hilang sendirian, dan penyimpanan
// tanpa jejak persis yang ADR-0007 larang.
//
// AKUN_ID adalah pengenal akun (ADR-U-0030), bukan nama orang.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Aksi jejak yang dikenal.
const (
	AksiJejakSimpan = "simpan"
	AksiJejakHapus  = "hapus"
	// Tiket 12 - lampiran tahun treaty.
	AksiJejakUnggah   = "unggah"
	AksiJejakUlangi   = "ulangi"
	AksiJejakMenyerah = "menyerah"
)

// CatatanJejakTCO adalah satu baris jejak seperti dibaca kembali.
type CatatanJejakTCO struct {
	ID         string
	Waktu      time.Time
	AkunID     string
	Tabel      string
	BarisID    string
	Aksi       string
	Keterangan string
}

func sqlSisipJejakTCO(tabel string) string {
	return fmt.Sprintf(`INSERT INTO %s
	       (ID, WAKTU, AKUN_ID, TABEL, BARIS_ID, AKSI, KETERANGAN)
	VALUES (:1, :2, :3, :4, :5, :6, :7)`, tabel)
}

func sqlBacaJejakTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, TO_CHAR(WAKTU, 'YYYY-MM-DD HH24:MI:SS'), AKUN_ID, TABEL, BARIS_ID, AKSI, KETERANGAN
	  FROM %s
	 WHERE TABEL = :1 AND BARIS_ID = :2
	 ORDER BY ID`, tabel)
}

// SisipJejakTCO menulis satu catatan jejak.
func (d *DB) SisipJejakTCO(ctx context.Context, tx *Tx,
	akunID, tabelDisentuh, barisID, aksi, keterangan string, waktu time.Time) error {

	if tx == nil {
		return errors.New("repository: jejak Treaty Contract Out menuntut transaksi")
	}
	tabel, err := d.Qualify(TabelJejakTCO)
	if err != nil {
		return err
	}
	id, err := d.IdentitasBerikutTCO(ctx, tx, SeqJejakTCO)
	if err != nil {
		return err
	}
	q := sqlSisipJejakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, id, waktu, akunID, tabelDisentuh,
		kosongJadiNil(barisID), aksi, kosongJadiNil(keterangan))
	if err != nil {
		return fmt.Errorf("repository: merekam jejak Treaty Contract Out: %w", err)
	}
	return pastikanSatuBaris(hasil, "perekaman jejak Treaty Contract Out")
}

// JejakTCO membaca jejak satu baris - dipakai uji dan layar riwayat kelak.
func (d *DB) JejakTCO(ctx context.Context, tabelDisentuh, barisID string) ([]CatatanJejakTCO, error) {
	tabel, err := d.Qualify(TabelJejakTCO)
	if err != nil {
		return nil, err
	}
	q := sqlBacaJejakTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := d.bacaTCO(ctx).QueryContext(ctx, q, tabelDisentuh, barisID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca jejak Treaty Contract Out: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []CatatanJejakTCO
	for rows.Next() {
		var n [7]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6]); err != nil {
			return nil, err
		}
		waktu, err := uraiTanggalTeks(n[1], "WAKTU")
		if err != nil {
			return nil, err
		}
		out = append(out, CatatanJejakTCO{ID: n[0].String, Waktu: waktu, AkunID: n[2].String,
			Tabel: n[3].String, BarisID: n[4].String, Aksi: n[5].String, Keterangan: n[6].String})
	}
	return out, rows.Err()
}
