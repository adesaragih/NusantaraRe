// Package repository memegang SQL modul Komite Claim Fac In. Nol aturan dagang.
//
// Tabel kasus komite dipakai bersama Komite Claim Life / Prop / Non Prop (`T_WORK_CLAIM`, `T_GENERAL_KOMITE`, tangga):
// setiap kueri menyaring `w.LINI = 'FACIN'` KETAT (pola Komite Claim Prop) dan awalan `KMT-` (fixture Claim Life juga
// memakai `KMT-` ber-LINI LIFE - saringan LINI yang menentukan). Tabel klaim (`T_GENERAL_CLAIM`, `T_CLAIM_ADJUSTMENT`)
// hanya DIBACA untuk kolom daftar kerja; tulisannya lewat kontrak `kontrak.KlaimFacInKomite`.
package repository

// Untuk apa berkas ini: GUDANG - kepala kasus komite (`T_GENERAL_KOMITE` + `T_WORK_CLAIM`): baca (dengan kunci untuk
// Submit), simpan tingkat / keputusan / usul / KOMITE_LOOP (KomitePost_Adjustment S14 / S24, KomitePost_Reject S16,
// ApprovalKomite_Act S8), tutup kasus (Komite_Flow END52 -> Resolved-Completed).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// Galat repository.
var (
	// ErrKasusTidakAda - ID tidak ada, bukan kasus komite FACIN, atau bukan KMT-.
	ErrKasusTidakAda = errors.New("repository: kasus Komite Claim Fac In tidak ada")
	// ErrKeputusanBersamaan - tangga / kepala kasus berubah sejak dibaca (dua klik, dua penyetuju).
	ErrKeputusanBersamaan = errors.New("repository: kasus komite berubah sejak dibaca; muat ulang kasusnya")
	// ErrTanpaTransaksi - penulis dipanggil di luar transaksi.
	ErrTanpaTransaksi = errors.New("repository: tulisan komite menuntut transaksi")
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

// sqlKepalaKasus - kepala satu kasus komite FACIN; `kunci` = FOR UPDATE (Submit berjalan di satu transaksi: dua Submit
// serentak diserialkan di sini). ADJUSTMENT_ID boleh NULL (TT3 / TT4, migrasi 641); TRANSFER_TYPE migrasi 642.
func sqlKepalaKasus(gen, work string, kunci bool) string {
	q := fmt.Sprintf(`SELECT g.ID, w.COVER_KEY, g.ADJUSTMENT_ID, g.TRANSFER_TYPE, g.KOMITE_LOOP, g.KOMITE_COUNT,
		       g.ACCEPT_STATUS, g.KOMITE_USUL_TUTUP, g.KOMITE_USUL_CADANG, w.TAHAP, w.STATUS_WORK, w.CREATE_OP,
		       w.CREATE_OP_NAME, %s, %s, g.KOMITE_CIRCUM_CAUSE_OF_LOSS, g.KOMITE_EXTENT_OF_LOSS,
		       g.KOMITE_LEGAL_LIABILITY
		  FROM %s g JOIN %s w ON w.ID = g.ID
		 WHERE g.ID = :1 AND w.LINI = :2 AND w.ID LIKE :3`, fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_CREATE"),
		fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_UPDATE"), gen, work)
	if kunci {
		q += ` FOR UPDATE OF g.KOMITE_COUNT`
	}
	return q
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
	err = g.barisAtau(ctx, tx, q, id, models.LiniFacIn, awalanLike()).Scan(&n[0], &n[1], &n[2], &n[3], &loop, &count,
		&n[4], &n[5], &n[6], &n[7], &n[8], &n[9], &n[10], &n[11], &n[12], &n[13], &n[14], &n[15])
	if errors.Is(err, sql.ErrNoRows) {
		return models.Kasus{}, fmt.Errorf("%w: %q", ErrKasusTidakAda, id)
	}
	if err != nil {
		return models.Kasus{}, fmt.Errorf("repository: membaca kasus komite: %w", err)
	}
	return models.Kasus{ID: n[0].String, KlaimID: n[1].String, AdjustmentID: n[2].String,
		TransferType: strings.TrimSpace(n[3].String), Loop: int(loop.Int64), Count: int(count.Int64),
		AcceptStatus: strings.TrimSpace(n[4].String), UsulTutup: strings.TrimSpace(n[5].String),
		UsulCadang: strings.TrimSpace(n[6].String), Tahap: n[7].String, StatusWork: n[8].String, PembuatID: n[9].String,
		PembuatNama: n[10].String, TglCreate: waktuDB(n[11].String), TglUpdate: waktuDB(n[12].String),
		Kronologi: n[13].String, Extent: n[14].String, Liability: n[15].String}, nil
}

// sqlSimpanKepala - KOMITE_COUNT (S14 / S24, Reject S16), KOMITE_LOOP (ApprovalKomite_Act S8 sesudah perluasan),
// `.AcceptStatus`, dua penanda usul (isian tingkat 1; kolom migrasi komiteclaimprop 680). Bersyarat `KOMITE_COUNT` lama.
func sqlSimpanKepala(gen string) string {
	return fmt.Sprintf(`UPDATE %s SET KOMITE_COUNT = :1, KOMITE_LOOP = :2, ACCEPT_STATUS = :3, KOMITE_USUL_TUTUP = :4,
		       KOMITE_USUL_CADANG = :5
		 WHERE ID = :6 AND KOMITE_COUNT = :7`, gen)
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
	tutup, cadang := kp.UsulTutup, kp.UsulCadang
	if tutup == "" {
		tutup = models.UsulTidak
	}
	if cadang == "" {
		cadang = models.UsulTidak
	}
	h, err := tx.ExecContext(ctx, q, kp.Count, kp.Loop, db.KosongJadiNil(kp.AcceptStatus), tutup, cadang, id, countLama)
	if err != nil {
		return fmt.Errorf("repository: menyimpan kepala kasus komite: %w", err)
	}
	return satuAtauBersamaan(h, "kepala kasus komite "+id)
}

// sqlTutupKasus - END52 `Komite_Flow` (decision IsKomiteLoop salah) -> Resolved-Completed. TAHAP dibiarkan (tidak pernah
// NULL); POSITION dikosongkan (tanpa pemegang).
func sqlTutupKasus(work string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, TGL_UPDATE = :2, POSITION = NULL
		 WHERE ID = :3 AND LINI = :4 AND STATUS_WORK IS NULL`, work)
}

// sqlSentuhKasus - TGL_UPDATE kasus komite yang masih berjalan (decision IsKomiteLoop -> assignment KomiteRouter);
// POSITION = KomiteID tingkat berikut (workbasket), seperti POSITION klaim = pemegang tahapnya.
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
	q, args := sqlSentuhKasus(work), []any{saat, db.KosongJadiNil(posisi), id, models.LiniFacIn}
	if selesai {
		q, args = sqlTutupKasus(work), []any{models.StatusSelesai, saat, id, models.LiniFacIn}
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
