package repository

// Untuk apa berkas ini: GUDANG - work object Claim Non Prop di T_WORK_CLAIM (tabel bersama Claim Life; LINI 'NONPROP', TAHAP
// selalu terisi) dan baris induknya di T_GENERAL_CLAIM (shared PK).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimnonprop/backend/models"
)

// Galat repository.
var (
	// ErrKasusTidakAda - ID tidak ada, atau bukan kasus NONPROP.
	ErrKasusTidakAda = errors.New("repository: kasus Claim Non Prop tidak ada")
	// ErrTahapBerubah - tahap kasus berubah di antara pembacaan dan penulisan.
	ErrTahapBerubah = errors.New("repository: tahap kasus berubah sejak dibaca; muat ulang berkasnya")
)

// Gudang - seluruh SQL modul.
type Gudang struct {
	db *db.DB
}

// Baru menyusun gudang atas basis data.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// DB - basis data gudang (penomor bersama).
func (g *Gudang) DB() *db.DB { return g.db }

// lebarUrut - ID T_WORK_CLAIM `<awalan><LPAD(SEQ_WORK_CLAIM,6)>` (pola Claim Life `RakitPengenalWork`).
const lebarUrut = 6

// IDKasusBerikut - ID kasus baru dari SEQ_WORK_CLAIM (urutan bersama Claim Life - ruang nomor tidak bertabrakan karena
// awalannya berbeda).
func (g *Gudang) IDKasusBerikut(ctx context.Context, tx *db.Tx, awalan string) (string, error) {
	if awalan != models.AwalanKlaim && awalan != models.AwalanKomite {
		return "", fmt.Errorf("repository: awalan kasus %q tidak dikenal", awalan)
	}
	n, err := g.db.NomorBerikut(ctx, tx, "SEQ_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	n = strings.TrimSpace(n)
	if len(n) < lebarUrut {
		n = strings.Repeat("0", lebarUrut-len(n)) + n
	}
	return awalan + n, nil
}

// sqlSisipKasus - kelahiran baris T_WORK_CLAIM kasus baru.
func sqlSisipKasus(work string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, TAHAP, POSITION, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, work)
}

// sqlSisipInduk - kelahiran baris induk T_GENERAL_CLAIM (ID + SUMBER).
func sqlSisipInduk(gen string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, SUMBER) VALUES (:1, :2)`, gen)
}

// SisipKasus melahirkan work object klaim di tahap Outstanding Claim beserta baris induk T_GENERAL_CLAIM.
func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipKasus(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id, models.LiniNonProp, models.TahapOutstanding, teksAtauNil(pembuat),
		teksAtauNil(pembuat), teksAtauNil(namaPembuat), saat, saat)
	if err != nil {
		return fmt.Errorf("repository: melahirkan kasus Claim Non Prop: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran kasus Claim Non Prop"); err != nil {
		return err
	}
	return g.sisipInduk(ctx, tx, id, models.SumberGo)
}

// sisipInduk melahirkan baris induk T_GENERAL_CLAIM (kasus baru: SUMBER GO; kasus pemuat data lama: PEGA).
func (g *Gudang) sisipInduk(ctx context.Context, tx *db.Tx, id, sumber string) error {
	gen, err := g.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipInduk(gen)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id, sumber)
	if err != nil {
		return fmt.Errorf("repository: melahirkan induk klaim %s: %w", id, err)
	}
	return db.PastikanSatuBaris(hasil, "kelahiran induk klaim")
}

// sqlKeadaan - satu baris work object NONPROP.
func sqlKeadaan(work, gen string, kunci bool) string {
	q := fmt.Sprintf(`SELECT w.ID, w.TAHAP, w.POSITION, w.STATUS_WORK, w.CREATE_OP, w.CREATE_OP_NAME,
		       %s, %s, g.SUMBER
		  FROM %s w JOIN %s g ON g.ID = w.ID
		 WHERE w.ID = :1 AND w.LINI = :2`, fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_CREATE"),
		fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_UPDATE"), work, gen)
	if kunci {
		q += ` FOR UPDATE OF w.TAHAP`
	}
	return q
}

func (g *Gudang) keadaan(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, error) {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return models.Kasus{}, err
	}
	gen, err := g.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return models.Kasus{}, err
	}
	q := sqlKeadaan(work, gen, kunci)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Kasus{}, err
	}
	var n [9]sql.NullString
	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, q, id, models.LiniNonProp)
	} else {
		row = g.db.QueryRowContext(ctx, q, id, models.LiniNonProp)
	}
	if err := row.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8]); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Kasus{}, ErrKasusTidakAda
		}
		return models.Kasus{}, fmt.Errorf("repository: membaca kasus Claim Non Prop: %w", err)
	}
	k := models.Kasus{ID: n[0].String, Tahap: n[1].String, Posisi: n[2].String, StatusWork: n[3].String,
		PembuatID: n[4].String, PembuatNama: n[5].String, Sumber: n[8].String}
	if t, err := utils.ParseTanggal(n[6].String); err == nil {
		k.TglCreate = t
	}
	if t, err := utils.ParseTanggal(n[7].String); err == nil {
		k.TglUpdate = t
	}
	return k, nil
}

// Keadaan membaca satu kasus (tx nil = tanpa transaksi).
func (g *Gudang) Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error) {
	return g.keadaan(ctx, tx, id, false)
}

// KunciKasus mengunci baris kasus (`FOR UPDATE`) dan memastikan tahap dan status belum berubah.
func (g *Gudang) KunciKasus(ctx context.Context, tx *db.Tx, id, tahap string) (models.Kasus, error) {
	k, err := g.keadaan(ctx, tx, id, true)
	if err != nil {
		return models.Kasus{}, err
	}
	if k.Tahap != tahap || k.Tertutup() {
		return models.Kasus{}, ErrTahapBerubah
	}
	return k, nil
}

// sqlPindahTahap - pemindahan tahap kasus terbuka bertahap `lama`.
func sqlPindahTahap(work string) string {
	return fmt.Sprintf(`UPDATE %s SET TAHAP = :1, POSITION = :2, TGL_UPDATE = :3
		 WHERE ID = :4 AND LINI = :5 AND TAHAP = :6 AND STATUS_WORK IS NULL`, work)
}

// PindahTahap - finishAssignment: tahap lama -> tahap baru (POSITION = workbasket pemegang tahap baru).
func (g *Gudang) PindahTahap(ctx context.Context, tx *db.Tx, id, lama, baru, posisi string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlPindahTahap(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, baru, teksAtauNil(posisi), saat, id, models.LiniNonProp, lama)
	if err != nil {
		return fmt.Errorf("repository: memindah tahap kasus: %w", err)
	}
	if n, _ := hasil.RowsAffected(); n == 0 {
		return ErrTahapBerubah
	}
	return db.PastikanSatuBaris(hasil, "pemindahan tahap kasus")
}

// sqlTutupKasus - STATUS_WORK kasus terbuka bertahap `tahap`.
func sqlTutupKasus(work string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, TGL_UPDATE = :2
		 WHERE ID = :3 AND LINI = :4 AND TAHAP = :5 AND STATUS_WORK IS NULL`, work)
}

// TutupKasus - `ASMForceCaseClose` (CloseClaimTNonProp langkah 9): STATUS_WORK Resolved-Completed. TAHAP dibiarkan berisi
// tahap terakhir (aturan modul: TAHAP tidak pernah NULL).
func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id, tahap string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlTutupKasus(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, models.StatusSelesai, saat, id, models.LiniNonProp, tahap)
	if err != nil {
		return fmt.Errorf("repository: menutup kasus: %w", err)
	}
	if n, _ := hasil.RowsAffected(); n == 0 {
		return ErrTahapBerubah
	}
	return db.PastikanSatuBaris(hasil, "penutupan kasus")
}

// sqlSentuhKasus - TGL_UPDATE kasus.
func sqlSentuhKasus(work string) string {
	return fmt.Sprintf(`UPDATE %s SET TGL_UPDATE = :1 WHERE ID = :2 AND LINI = :3`, work)
}

// SentuhKasus memperbarui TGL_UPDATE (setiap aksi yang menyimpan halaman).
func (g *Gudang) SentuhKasus(ctx context.Context, tx *db.Tx, id string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlSentuhKasus(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, saat, id, models.LiniNonProp)
	if err != nil {
		return fmt.Errorf("repository: memperbarui kasus: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "pembaruan kasus")
}

// ---------------------------------------------------------------- daftar kerja

// SaringanKasus - daftar kerja halaman awal.
type SaringanKasus struct {
	// Tahap - TahapOutstanding (worklist pembuat) / TahapAcceptation (workbasket).
	Tahap string
	// Pembuat - CREATE_OP (worklist Assignment2, ToCurrentOperator). Kosong = semua.
	Pembuat string
	// Selesai - true = Resolved-Completed saja.
	Selesai bool
	Cari    string
	Batas   int
}

// RingkasanKasus - satu baris daftar kerja.
type RingkasanKasus struct {
	ID          string `json:"id"`
	Tahap       string `json:"tahap"`
	Label       string `json:"label"`
	StatusWork  string `json:"statusWork"`
	PembuatNama string `json:"pembuatNama"`
	TglCreate   string `json:"tglCreate"`
	NoClaim     string `json:"noClaim"`
	ClaimNoTemp string `json:"claimNoTemp"`
	PolicyNo    string `json:"policyNo"`
	TreatyName  string `json:"treatyName"`
	InsuredName string `json:"insuredName"`
}

// DaftarKasus membaca daftar kerja (terbaru dahulu).
func (g *Gudang) DaftarKasus(ctx context.Context, s SaringanKasus) ([]RingkasanKasus, error) {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return nil, err
	}
	gen, err := g.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return nil, err
	}
	var args []any
	ph := func(v any) string { args = append(args, v); return fmt.Sprintf(":%d", len(args)) }
	w := []string{"w.LINI = " + ph(models.LiniNonProp), "w.ID LIKE " + ph(models.AwalanKlaim+"%")}
	if s.Selesai {
		w = append(w, "w.STATUS_WORK = "+ph(models.StatusSelesai))
	} else {
		w = append(w, "w.STATUS_WORK IS NULL")
		if s.Tahap != "" {
			w = append(w, "w.TAHAP = "+ph(s.Tahap))
		}
	}
	if s.Pembuat != "" {
		w = append(w, "w.CREATE_OP = "+ph(s.Pembuat))
	}
	if c := strings.TrimSpace(s.Cari); c != "" {
		pola := "%" + strings.ToUpper(c) + "%"
		var atau []string
		for _, kol := range []string{"w.ID", "g.CLAIM_NO", "g.CLAIM_NO_TEMP", "g.POLICY_NO", "g.INSURED_NAME"} {
			atau = append(atau, "UPPER("+kol+") LIKE "+ph(pola))
		}
		w = append(w, "("+strings.Join(atau, " OR ")+")")
	}
	batas := s.Batas
	if batas <= 0 || batas > 500 {
		batas = 500
	}
	q := fmt.Sprintf(`SELECT w.ID, w.TAHAP, w.STATUS_WORK, w.CREATE_OP_NAME, %s,
		       g.CLAIM_NO, g.CLAIM_NO_TEMP, g.POLICY_NO, g.TREATY_NAME, g.INSURED_NAME
		  FROM %s w JOIN %s g ON g.ID = w.ID
		 WHERE %s
		 ORDER BY w.TGL_CREATE DESC, w.ID DESC
		 FETCH FIRST %d ROWS ONLY`, fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_CREATE"), work, gen,
		strings.Join(w, " AND "), batas)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar kerja Claim Non Prop: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []RingkasanKasus{}
	for rows.Next() {
		var n [10]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8], &n[9]); err != nil {
			return nil, fmt.Errorf("repository: memindai daftar kerja Claim Non Prop: %w", err)
		}
		out = append(out, RingkasanKasus{ID: n[0].String, Tahap: n[1].String, Label: models.LabelTahap[n[1].String],
			StatusWork: n[2].String, PembuatNama: n[3].String, TglCreate: n[4].String, NoClaim: n[5].String,
			ClaimNoTemp: n[6].String, PolicyNo: n[7].String, TreatyName: n[8].String, InsuredName: n[9].String})
	}
	return out, rows.Err()
}
