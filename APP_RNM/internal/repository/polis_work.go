package repository

// Baris kerja polis - `T_WORK_POLIS`, tiket 01 PremiumList Life.
//
// Untuk apa berkas ini: membaca dan menulis tahap serta status kerja sebuah
// polis. Tabelnya lahir di migrasi 050; kolomnya empat, dan keduanya yang
// berarti - `POSITION` dan `STATUS` - menyimpan teks VERBATIM dari flow
// (`models.TahapPolis*`, `models.StatusPolis*`).
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033), nol
// `COMMIT` (ADR-U-0029).
//
// Dibaca sesudah: migrasi 050.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrWorkPolisTidakAda - baris kerja polis yang diminta tidak ada.
var ErrWorkPolisTidakAda = errors.New("repository: baris kerja polis tidak ada")

// WorkPolis membaca dan menulis `T_WORK_POLIS`.
type WorkPolis struct{ db *DB }

// NewWorkPolis menyusunnya.
func NewWorkPolis(db *DB) *WorkPolis { return &WorkPolis{db: db} }

// KeadaanPolis adalah tahap dan status kerja sebuah polis.
type KeadaanPolis struct {
	ID string
	// Position VERBATIM nama assignment flow - lihat models.TahapPolis*.
	Position string
	// Status kosong berarti kasus BELUM ditutup (ADR-U-0027). Flow hanya
	// mengenal dua status akhir, dan nol status untuk kasus berjalan.
	Status string
	Lini   string
}

// sqlKeadaanPolis merakit pembacaannya.
func sqlKeadaanPolis(tabel string) string {
	return fmt.Sprintf(
		`SELECT ID, LINI, POSITION, STATUS FROM %s WHERE ID = :1`, tabel)
}

// Keadaan membaca tahap dan status kerja satu polis.
func (r *WorkPolis) Keadaan(ctx context.Context, id string) (KeadaanPolis, error) {
	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return KeadaanPolis{}, err
	}
	q := sqlKeadaanPolis(tabel)
	if err := PeriksaSQL(q); err != nil {
		return KeadaanPolis{}, err
	}
	var pengenal, lini, posisi, status sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, id).Scan(&pengenal, &lini, &posisi, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return KeadaanPolis{}, ErrWorkPolisTidakAda
	}
	if err != nil {
		return KeadaanPolis{}, fmt.Errorf("repository: membaca kerja polis: %w", err)
	}
	return KeadaanPolis{
		ID: pengenal.String, Lini: lini.String,
		Position: posisi.String, Status: status.String,
	}, nil
}

// sqlPindahTahapPolis merakit perpindahan tahap.
//
// ⛔ Syarat WHERE menyertakan POSITION LAMA, dan itu bukan kehati-hatian
// berlebih: baris dibaca di luar transaksi, jadi ia dapat berpindah di
// antara baca dan tulis. Tanpa syarat itu dua permintaan serentak sama-sama
// menang, dan yang kedua memindahkan kasus dari tahap yang sudah bukan
// tahapnya lagi. Pola yang sama dengan `PerbaruiStatusBaris` di Claim Life.
//
// ⚠️ `STATUS` ikut DIKOSONGKAN: perpindahan tahap berarti kasus berjalan
// lagi. Kasus yang pindah tetapi statusnya masih `Resolved-*` akan tampak
// tertutup dari satu sisi dan berjalan dari sisi lain.
func sqlPindahTahapPolis(tabel string) string {
	return fmt.Sprintf(
		`UPDATE %s SET POSITION = :1, STATUS = NULL
		  WHERE ID = :2 AND (POSITION = :3 OR (POSITION IS NULL AND :3 IS NULL))`, tabel)
}

// PindahTahap memindahkan polis ke tahap lain.
func (r *WorkPolis) PindahTahap(ctx context.Context, tx *Tx,
	id, posisiLama, posisiBaru string) error {

	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return err
	}
	q := sqlPindahTahapPolis(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, posisiBaru, id, kosongJadiNil(posisiLama))
	if err != nil {
		return fmt.Errorf("repository: memindahkan tahap polis: %w", err)
	}
	return pastikanSatuBaris(hasil, "perpindahan tahap polis")
}

// sqlTutupPolis merakit penutupan kasus.
//
// ⛔ `POSITION` DIKOSONGKAN saat ditutup - pola yang sama dengan
// `TutupKasus` Claim Life (butir bb): kotak masuk adalah worklist, dan kasus
// yang tertutup tidak boleh berdiri di antrean mana pun.
//
// ⛔ Dan `STATUS IS NULL` menjadi syarat: kasus yang SUDAH tertutup tidak
// ditutup lagi dengan status yang berbeda. Dua penutupan berturut-turut akan
// menimpa alasan penutupan yang pertama, dan alasan itu jejak.
func sqlTutupPolis(tabel string) string {
	return fmt.Sprintf(
		`UPDATE %s SET STATUS = :1, POSITION = NULL
		  WHERE ID = :2 AND STATUS IS NULL`, tabel)
}

// TutupKasus menutup kasus polis dengan status kerja akhirnya.
func (r *WorkPolis) TutupKasus(ctx context.Context, tx *Tx, id, status string) error {
	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return err
	}
	q := sqlTutupPolis(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("repository: menutup kasus polis: %w", err)
	}
	return pastikanSatuBaris(hasil, "penutupan kasus polis")
}
