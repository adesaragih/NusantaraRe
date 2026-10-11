// Package repository memegang SQL modul Komite Claim Prop. Nol aturan dagang.
//
// Tabel kasus komite dipakai bersama Komite Claim Life (`T_WORK_CLAIM`, `T_GENERAL_KOMITE`, tangga): setiap kueri
// menyaring `w.LINI = 'PROP'` KETAT (tanpa `OR LINI IS NULL`, keputusan work owner 07-10-2026) dan awalan `TKMT-`.
// Tabel klaim Prop (`T_GENERAL_CLAIM`, `T_CLAIM_*`) hanya DIBACA untuk kolom daftar kerja (preseden
// `komiteclaimlife/.../komite_inbox.go`); tulisannya lewat kontrak `kontrak.KlaimTreatyKomite`.
package repository

// Untuk apa berkas ini: GUDANG - kepala kasus komite (`T_GENERAL_KOMITE` + `T_WORK_CLAIM`): baca (dengan kunci untuk
// Submit), simpan tingkat / keputusan / usul (`KomitePostAdjustment` S25 / S40, isian tingkat 1), tutup kasus
// (Decision `KomiteLoop` -> Resolved-Completed).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

// Galat repository.
var (
	// ErrKasusTidakAda - ID tidak ada, bukan kasus komite PROP, atau bukan TKMT-.
	ErrKasusTidakAda = errors.New("repository: kasus Komite Claim Prop tidak ada")
	// ErrKeputusanBersamaan - tangga / kepala kasus berubah sejak dibaca (dua klik, dua penyetuju).
	ErrKeputusanBersamaan = errors.New("repository: kasus komite berubah sejak dibaca; muat ulang kasusnya")
)

// Gudang - seluruh SQL modul.
type Gudang struct{ db *db.DB }

// Baru menyusun gudang atas basis data.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// barisAtau - QueryRow di dalam tx bila ada, selainnya di luar transaksi.
func (g *Gudang) barisAtau(ctx context.Context, tx *db.Tx, q string, args ...any) *sql.Row {
	if tx != nil {
		return tx.QueryRowContext(ctx, q, args...)
	}
	return g.db.QueryRowContext(ctx, q, args...)
}

// sqlKepalaKasus - kepala satu kasus komite PROP; `kunci` = FOR UPDATE (KomitePostAdjustment berjalan di satu
// transaksi: dua Submit serentak diserialkan di sini).
func sqlKepalaKasus(gen, work string, kunci bool) string {
	q := fmt.Sprintf(`SELECT g.ID, w.COVER_KEY, g.ADJUSTMENT_ID, g.KOMITE_LOOP, g.KOMITE_COUNT, g.ACCEPT_STATUS,
		       g.KOMITE_USUL_TUTUP, g.KOMITE_USUL_CADANG, g.KOMITE_SUBJECTIVITY, g.KOMITE_SUBJECTIVITY_NOTE, w.TAHAP, w.STATUS_WORK, w.CREATE_OP, w.CREATE_OP_NAME,
		       %s, %s
		  FROM %s g JOIN %s w ON w.ID = g.ID
		 WHERE g.ID = :1 AND w.LINI = :2 AND w.ID LIKE :3`, fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_CREATE"),
		fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_UPDATE"), gen, work)
	if kunci {
		q += ` FOR UPDATE OF g.KOMITE_COUNT`
	}
	return q
}

// awalanLike - pola LIKE awalan kasus komite.
func awalanLike() string { return models.AwalanKomite + "%" }

// waktuDB - teks "YYYY-MM-DD HH24:MI:SS" (jam Jakarta di Oracle) -> time.
func waktuDB(s string) time.Time {
	t, err := time.ParseInLocation(utils.TanggalWaktu, strings.TrimSpace(s), models.Jakarta)
	if err != nil {
		return time.Time{}
	}
	return t
}

// BacaKasus membaca kepala kasus `id` (tanpa tangga). `kunci` hanya di dalam transaksi.
func (g *Gudang) BacaKasus(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, error) {
	if kunci && tx == nil {
		return models.Kasus{}, errors.New("repository: kunci kasus komite menuntut transaksi")
	}
	gen, err := g.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return models.Kasus{}, err
	}
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return models.Kasus{}, err
	}
	q := sqlKepalaKasus(gen, work, kunci)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Kasus{}, err
	}
	var n [16]sql.NullString
	var loop, count sql.NullInt64
	err = g.barisAtau(ctx, tx, q, id, models.LiniProp, awalanLike()).Scan(&n[0], &n[1], &n[2], &loop, &count, &n[5],
		&n[6], &n[7], &n[14], &n[15], &n[8], &n[9], &n[10], &n[11], &n[12], &n[13])
	if errors.Is(err, sql.ErrNoRows) {
		return models.Kasus{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	if err != nil {
		return models.Kasus{}, fmt.Errorf("repository: membaca kasus komite: %w", err)
	}
	k := models.Kasus{ID: n[0].String, KlaimID: n[1].String, AdjustmentID: n[2].String, Loop: int(loop.Int64),
		Count: int(count.Int64), AcceptStatus: strings.TrimSpace(n[5].String), UsulTutup: strings.TrimSpace(n[6].String),
		UsulCadang: strings.TrimSpace(n[7].String), Subjectivity: strings.TrimSpace(n[14].String),
		SubjectivityNote: n[15].String, Tahap: n[8].String, StatusWork: n[9].String, PembuatID: n[10].String,
		PembuatNama: n[11].String, TglCreate: waktuDB(n[12].String), TglUpdate: waktuDB(n[13].String),
		TransferType: models.TransferAdjustment}
	if k.AdjustmentID == "" { // Close Without Payment (TT 4): kolom bersama migrasi komiteclaimfacin 642 / 643
		if err := g.bacaTeksTutup(ctx, tx, gen, &k); err != nil {
			return models.Kasus{}, err
		}
	}
	return k, nil
}

// sqlTeksTutup - jenis penyerahan dan teks Chronology kasus komite tanpa adjustment. Hanya dibaca untuk kasus itu,
// sehingga kasus TT 2 tidak bergantung kolom migrasi komiteclaimfacin 642 / 643.
func sqlTeksTutup(gen string) string {
	return fmt.Sprintf(`SELECT TRANSFER_TYPE, KOMITE_CIRCUM_CAUSE_OF_LOSS FROM %s WHERE ID = :1`, gen)
}

// bacaTeksTutup mengisi TransferType / Kronologi kasus tanpa adjustment; jenis selain TT 4 ditolak (Claim Prop tanpa
// penulis TT 3).
func (g *Gudang) bacaTeksTutup(ctx context.Context, tx *db.Tx, gen string, k *models.Kasus) error {
	q := sqlTeksTutup(gen)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	var tt, kron sql.NullString
	if err := g.barisAtau(ctx, tx, q, k.ID).Scan(&tt, &kron); err != nil {
		return fmt.Errorf("repository: membaca jenis penyerahan kasus komite: %w", err)
	}
	if strings.TrimSpace(tt.String) != models.TransferTutup {
		return fmt.Errorf("%w: kasus %q tanpa adjustment ber-TRANSFER_TYPE %q", ErrKasusTidakAda, k.ID, tt.String)
	}
	k.TransferType, k.Kronologi = models.TransferTutup, kron.String
	return nil
}

// sqlSimpanKepala - KomitePostAdjustment S25 (`KomiteCount := KomiteLoop` saat tolak) + S40 (`KomiteCount + 1`),
// `.AcceptStatus` (bahan `IsKomiteLoop`), dua penanda usul dan isian Subjectivity (isian tingkat 1, migrasi 680 / 682).
// Bersyarat `KOMITE_COUNT` lama.
func sqlSimpanKepala(gen string) string {
	return fmt.Sprintf(`UPDATE %s SET KOMITE_COUNT = :1, ACCEPT_STATUS = :2, KOMITE_USUL_TUTUP = :3,
		       KOMITE_USUL_CADANG = :4, KOMITE_SUBJECTIVITY = :5, KOMITE_SUBJECTIVITY_NOTE = :6
		 WHERE ID = :7 AND KOMITE_COUNT = :8`, gen)
}

// SimpanKepala menulis kepala kasus sesudah satu Submit.
func (g *Gudang) SimpanKepala(ctx context.Context, tx *db.Tx, id string, countLama int, kp models.Kepala) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	gen, err := g.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return err
	}
	q := sqlSimpanKepala(gen)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	subj := kp.Subjectivity
	if subj == "" {
		subj = models.UsulTidak
	}
	h, err := tx.ExecContext(ctx, q, kp.Count, db.KosongJadiNil(kp.AcceptStatus), kp.UsulTutup, kp.UsulCadang, subj,
		db.KosongJadiNil(kp.SubjectivityNote), id, countLama)
	if err != nil {
		return fmt.Errorf("repository: menyimpan kepala kasus komite: %w", err)
	}
	return satuAtauBersamaan(h, "kepala kasus komite "+id)
}

// ErrTanpaTransaksi - penulis dipanggil di luar transaksi.
var ErrTanpaTransaksi = errors.New("repository: tulisan komite menuntut transaksi")

func wajibTx(tx *db.Tx) error {
	if tx == nil {
		return ErrTanpaTransaksi
	}
	return nil
}

// satuAtauBersamaan - tepat satu baris, selainnya keadaan berubah sejak dibaca.
func satuAtauBersamaan(h sql.Result, nama string) error {
	n, err := h.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: membaca cacah %s: %w", nama, err)
	}
	if n != 1 {
		return fmt.Errorf("%w: %s (%d baris)", ErrKeputusanBersamaan, nama, n)
	}
	return nil
}

// sqlTutupKasus - Decision `KomiteLoop` salah -> End (Resolved-Completed). TAHAP dibiarkan (tidak pernah NULL);
// POSITION dikosongkan (tanpa pemegang).
func sqlTutupKasus(work string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, TGL_UPDATE = :2, POSITION = NULL
		 WHERE ID = :3 AND LINI = :4 AND STATUS_WORK IS NULL`, work)
}

// sqlSentuhKasus - TGL_UPDATE kasus komite yang masih berjalan (assignment kembali ke KomiteRouter); POSITION =
// KomiteID tingkat berikut (workbasket, keputusan work owner 09-10-2026), seperti POSITION klaim = pemegang tahapnya.
func sqlSentuhKasus(work string) string {
	return fmt.Sprintf(`UPDATE %s SET TGL_UPDATE = :1, POSITION = :2
		 WHERE ID = :3 AND LINI = :4 AND STATUS_WORK IS NULL`, work)
}

// TutupKasus menutup kasus komite (`selesai`) atau memperbarui TGL_UPDATE dan POSITION (`posisi`)-nya.
func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id string, selesai bool, posisi string, saat time.Time) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q, args := sqlSentuhKasus(work), []any{saat, db.KosongJadiNil(posisi), id, models.LiniProp}
	if selesai {
		q, args = sqlTutupKasus(work), []any{models.StatusSelesai, saat, id, models.LiniProp}
	}
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	h, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: menutup kasus komite: %w", err)
	}
	return satuAtauBersamaan(h, "work object komite "+id)
}
