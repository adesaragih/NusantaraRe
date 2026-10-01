package repository

// Baris kerja polis - `T_WORK_POLIS`, tiket 01 PremiumList Life.
//
// Untuk apa berkas ini: membaca dan menulis tahap serta status kerja sebuah
// polis. Tabelnya lahir di migrasi 050; dua kolom yang berarti - `POSITION`
// dan `STATUS_WORK` (bernama `STATUS` sampai 059) - menyimpan teks VERBATIM
// dari flow (`models.TahapPolis*`, `models.StatusPolis*`).
//
// ⛔ Sejak 059 (seragam `T_WORK_CLAIM`, keputusan work owner 01-10-2026)
// setiap ubah baris kasus menulis `TGL_UPDATE = SYSDATE` di pernyataan yang
// sama, dan keadaan ikut membaca pembuat, waktu, dan `COVER_KEY`.
//
// ⛔ Setiap query menyebut skemanya lewat `Qualify` (ADR-U-0033), nol
// `COMMIT` (ADR-U-0029).
//
// Dibaca sesudah: migrasi 050, 059.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

// ErrWorkPolisTidakAda - baris kerja polis yang diminta tidak ada.
var ErrWorkPolisTidakAda = errors.New("repository: baris kerja polis tidak ada")

// WorkPolis membaca dan menulis `T_WORK_POLIS`.
type WorkPolis struct{ db *db.DB }

// NewWorkPolis menyusunnya.
func NewWorkPolis(db *db.DB) *WorkPolis { return &WorkPolis{db: db} }

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
//	STATUS_WORK `pyWorkStatus`       -> Input Offer Life | Input Premium
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
	// Status - `pyWorkStatus` (kolom `STATUS_WORK`): nama assignment selama
	// berjalan, `Resolved-*` saat tertutup. Kosong berarti belum pernah
	// disetel (ADR-U-0027).
	Status string
	Lini   string
	// CoverKey - kasus induk (`T_WORK_POLIS.ID`); kosong sampai ada modul
	// yang terbukti mengisinya.
	CoverKey string
	// CreateOp, CreateOpName - pembuat kasus; kosong untuk baris yang lahir
	// sebelum 059 (CREATE_OP selalu, CREATE_OP_NAME bila tidak terisi ulang).
	CreateOp, CreateOpName string
	// TglCreate, TglUpdate - waktu lahir dan ubah terakhir; nol bila kosong.
	TglCreate, TglUpdate time.Time
}

// sqlKeadaanPolis merakit pembacaannya.
func sqlKeadaanPolis(tabel string) string {
	return fmt.Sprintf(
		`SELECT ID, LINI, POSITION, STATUS_WORK, COVER_KEY, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE FROM %s WHERE ID = :1`, tabel)
}

// sqlBenderaPolis membaca `FLAG_ONGOING_POLICY` saja - butir bq.
//
// ⛔ TERPISAH dari keadaan, dan itu sengaja (temuan tinjauan GILIRAN-14):
// kolomnya lahir di 057. Membacanya di `Keadaan` membuat SETIAP keputusan -
// `Reject`, `Decline` - gagal ORA-00904 selama 057 belum berjalan, padahal
// hanya `Confirm` di tahap penawaran yang membutuhkannya.
func sqlBenderaPolis(tabel string) string {
	return fmt.Sprintf(`SELECT FLAG_ONGOING_POLICY FROM %s WHERE ID = :1`, tabel)
}

// Bendera membaca `FLAG_ONGOING_POLICY` satu polis ("" bila kosong).
func (r *WorkPolis) Bendera(ctx context.Context, id string) (string, error) {
	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return "", err
	}
	q := sqlBenderaPolis(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var flag sql.NullString
	err = r.db.QueryRowContext(ctx, q, id).Scan(&flag)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrWorkPolisTidakAda
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca bendera polis: %w", err)
	}
	return flag.String, nil
}

// Keadaan membaca tahap dan status kerja satu polis.
func (r *WorkPolis) Keadaan(ctx context.Context, id string) (KeadaanPolis, error) {
	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return KeadaanPolis{}, err
	}
	q := sqlKeadaanPolis(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return KeadaanPolis{}, err
	}
	var pengenal, lini, posisi, status, induk, op, namaOp sql.NullString
	var lahir, ubah sql.NullTime
	err = r.db.QueryRowContext(ctx, q, id).Scan(&pengenal, &lini, &posisi, &status,
		&induk, &op, &namaOp, &lahir, &ubah)
	if errors.Is(err, sql.ErrNoRows) {
		return KeadaanPolis{}, ErrWorkPolisTidakAda
	}
	if err != nil {
		return KeadaanPolis{}, fmt.Errorf("repository: membaca kerja polis: %w", err)
	}
	return KeadaanPolis{
		ID: pengenal.String, Lini: lini.String,
		Position: posisi.String, Status: status.String,
		CoverKey: induk.String, CreateOp: op.String, CreateOpName: namaOp.String,
		TglCreate: lahir.Time, TglUpdate: ubah.Time,
	}, nil
}

// sqlPindahTahapPolis merakit perpindahan tahap.
//
// ⛔ Yang ditulis `STATUS_WORK`, bukan `POSITION` - lihat ralat di
// `KeadaanPolis`. Tahap adalah `pyWorkStatus`; `POSITION` menyimpan posisi
// layar (`Offer`/`Premium`) dan tidak berubah karena perpindahan tahap.
//
// ⛔ Syarat WHERE menyertakan STATUS_WORK LAMA, dan itu bukan kehati-hatian
// berlebih: baris dibaca di luar transaksi, jadi ia dapat berpindah di
// antara baca dan tulis. Tanpa syarat itu dua permintaan serentak sama-sama
// menang, dan yang kedua memindahkan kasus dari tahap yang sudah bukan
// tahapnya lagi. Pola yang sama dengan `PerbaruiStatusBaris` di Claim Life.
func sqlPindahTahapPolis(tabel string) string {
	return fmt.Sprintf(
		`UPDATE %s SET STATUS_WORK = :1, TGL_UPDATE = SYSDATE
		  WHERE ID = :2 AND (STATUS_WORK = :3 OR (STATUS_WORK IS NULL AND :3 IS NULL))`, tabel)
}

// PindahTahap memindahkan polis ke tahap lain.
func (r *WorkPolis) PindahTahap(ctx context.Context, tx *db.Tx,
	id, tahapLama, tahapBaru string) error {

	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return err
	}
	q := sqlPindahTahapPolis(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, tahapBaru, id, db.KosongJadiNil(tahapLama))
	if err != nil {
		return fmt.Errorf("repository: memindahkan tahap polis: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "perpindahan tahap polis")
}

// sqlTutupPolis merakit penutupan kasus.
//
// ⛔ `POSITION` DIKOSONGKAN saat ditutup - pola yang sama dengan
// `TutupKasus` Claim Life (butir bb): kotak masuk adalah worklist, dan kasus
// yang tertutup tidak boleh berdiri di antrean mana pun.
//
// ⛔ Syaratnya `STATUS_WORK` BUKAN salah satu status akhir - bukan `IS NULL`.
// Sejak ralat 28-09-2026, kolom itu juga menyimpan nama tahap selama kasus
// berjalan, jadi `IS NULL` hanya benar untuk kasus yang belum pernah
// bertahap. Yang dijaga: kasus yang SUDAH tertutup tidak ditutup lagi dengan
// alasan yang berbeda - dan alasan penutupan itu jejak.
//
// ⛔ DAN tahapnya harus tahap yang DIBACA (`:5`) - pola `sqlPindahTahapPolis`.
// Keadaan dibaca di luar transaksi; tanpa syarat ini kasus yang berpindah
// tahap di antara baca dan tulis tetap ditutup lewat konektor tahap lamanya
// (temuan /code-review giliran 10: `Submit` summary).
func sqlTutupPolis(tabel string) string {
	return fmt.Sprintf(
		`UPDATE %s SET STATUS_WORK = :1, POSITION = NULL, TGL_UPDATE = SYSDATE
		  WHERE ID = :2 AND (STATUS_WORK IS NULL OR STATUS_WORK NOT IN (:3, :4))
		    AND (STATUS_WORK = :5 OR (STATUS_WORK IS NULL AND :5 IS NULL))`, tabel)
}

// TutupKasus menutup kasus polis dengan status kerja akhirnya - hanya bila
// ia masih di `tahapLama`.
func (r *WorkPolis) TutupKasus(ctx context.Context, tx *db.Tx, id, tahapLama, status string) error {
	tabel, err := r.db.Qualify("T_WORK_POLIS")
	if err != nil {
		return err
	}
	q := sqlTutupPolis(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, status, id,
		models.StatusPolisDitolak, models.StatusPolisSelesai, db.KosongJadiNil(tahapLama))
	if err != nil {
		return fmt.Errorf("repository: menutup kasus polis: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penutupan kasus polis")
}
