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

	"nusantarare/internal/models"
)

// ErrWorkPolisTidakAda - baris kerja polis yang diminta tidak ada.
var ErrWorkPolisTidakAda = errors.New("repository: baris kerja polis tidak ada")

// WorkPolis membaca dan menulis `T_WORK_POLIS`.
type WorkPolis struct{ db *DB }

// NewWorkPolis menyusunnya.
func NewWorkPolis(db *DB) *WorkPolis { return &WorkPolis{db: db} }

// KeadaanPolis adalah posisi layar dan status kerja sebuah polis.
//
// ⛔ RALAT 28-09-2026. Ronde pertama menulis bahwa `POSITION` menyimpan
// nama assignment. KELIRU, dan `Activity/ProtectAccept.xml` yang
// membantahnya: ia membandingkan `pyWorkPage.Position` dengan `"Offer"`
// b1207 dan `"Premium"` b2288 - bukan dengan nama assignment mana pun.
// Kedua nilai itu satu-satunya yang korpus pernah setel ke properti itu.
//
// Pembagian yang benar:
//
//	POSITION  `pyWorkPage.Position`  -> Offer | Premium   (posisi LAYAR)
//	STATUS    `pyWorkStatus`         -> Input Offer Life | Input Premium
//	                                    Detail | Input Premium Summary |
//	                                    Resolved-Rejected | Resolved-Completed
//
// ⚠️ Akibat kekeliruannya nyata: baris warisan menyimpan `Offer` di
// `POSITION`, dan pencarian dengan `"Input Offer Life"` menemukan NOL baris -
// kotak masuk kosong untuk pekerjaan yang benar-benar ada.
type KeadaanPolis struct {
	ID string
	// Position - `Offer` atau `Premium`; lihat models.Posisi*.
	Position string
	// Status - `pyWorkStatus`: nama assignment selama berjalan, `Resolved-*`
	// saat tertutup. Kosong berarti belum pernah disetel (ADR-U-0027).
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
// ⛔ Yang ditulis `STATUS`, bukan `POSITION` - lihat ralat di
// `KeadaanPolis`. Tahap adalah `pyWorkStatus`; `POSITION` menyimpan posisi
// layar (`Offer`/`Premium`) dan tidak berubah karena perpindahan tahap.
//
// ⛔ Syarat WHERE menyertakan STATUS LAMA, dan itu bukan kehati-hatian
// berlebih: baris dibaca di luar transaksi, jadi ia dapat berpindah di
// antara baca dan tulis. Tanpa syarat itu dua permintaan serentak sama-sama
// menang, dan yang kedua memindahkan kasus dari tahap yang sudah bukan
// tahapnya lagi. Pola yang sama dengan `PerbaruiStatusBaris` di Claim Life.
func sqlPindahTahapPolis(tabel string) string {
	return fmt.Sprintf(
		`UPDATE %s SET STATUS = :1
		  WHERE ID = :2 AND (STATUS = :3 OR (STATUS IS NULL AND :3 IS NULL))`, tabel)
}

// PindahTahap memindahkan polis ke tahap lain.
func (r *WorkPolis) PindahTahap(ctx context.Context, tx *Tx,
	id, tahapLama, tahapBaru string) error {

	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return err
	}
	q := sqlPindahTahapPolis(tabel)
	if err := PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.tx.ExecContext(ctx, q, tahapBaru, id, kosongJadiNil(tahapLama))
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
// ⛔ Syaratnya `STATUS` BUKAN salah satu status akhir - bukan `IS NULL`.
// Sejak ralat 28-09-2026, `STATUS` juga menyimpan nama tahap selama kasus
// berjalan, jadi `IS NULL` hanya benar untuk kasus yang belum pernah
// bertahap. Yang dijaga: kasus yang SUDAH tertutup tidak ditutup lagi dengan
// alasan yang berbeda - dan alasan penutupan itu jejak.
func sqlTutupPolis(tabel string) string {
	return fmt.Sprintf(
		`UPDATE %s SET STATUS = :1, POSITION = NULL
		  WHERE ID = :2 AND (STATUS IS NULL OR STATUS NOT IN (:3, :4))`, tabel)
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
	hasil, err := tx.tx.ExecContext(ctx, q, status, id,
		models.StatusPolisDitolak, models.StatusPolisSelesai)
	if err != nil {
		return fmt.Errorf("repository: menutup kasus polis: %w", err)
	}
	return pastikanSatuBaris(hasil, "penutupan kasus polis")
}
