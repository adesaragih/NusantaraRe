package repository

// Untuk apa berkas ini: GUDANG - work object Claim Fac In di T_WORK_CLAIM (tabel bersama Claim Life; LINI 'FACIN', TAHAP
// = nama FlowAction assignment, selalu terisi) dan baris induknya di T_GENERAL_CLAIM (shared PK). ⛔ Lini disaring lewat
// T_WORK_CLAIM.LINI, tidak pernah lewat awalan ID (Claim Life juga ber-awalan `CLM-`).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimfacin/backend/models"
)

// Galat repository.
var (
	// ErrKasusTidakAda - ID tidak ada, atau bukan kasus FACIN.
	ErrKasusTidakAda = errors.New("repository: kasus Claim Fac In tidak ada")
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

// SisipKasus melahirkan work object klaim di tahap Input Register (Start -> Decision4 IsSPK salah: polis belum dipilih ->
// Assignment1) beserta baris induk T_GENERAL_CLAIM.
func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipKasus(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, id, models.LiniFacIn, models.TahapRegister, teksAtauNil(pembuat),
		teksAtauNil(pembuat), teksAtauNil(namaPembuat), saat, saat)
	if err != nil {
		return fmt.Errorf("repository: melahirkan kasus Claim Fac In: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran kasus Claim Fac In"); err != nil {
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

// sqlKeadaan - satu baris work object FACIN.
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
		row = tx.QueryRowContext(ctx, q, id, models.LiniFacIn)
	} else {
		row = g.db.QueryRowContext(ctx, q, id, models.LiniFacIn)
	}
	if err := row.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8]); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Kasus{}, ErrKasusTidakAda
		}
		return models.Kasus{}, fmt.Errorf("repository: membaca kasus Claim Fac In: %w", err)
	}
	if n[1].String == models.TahapKomite { // baris KMT- ber-LINI FACIN yang sama: bukan kasus klaim
		return models.Kasus{}, ErrKasusTidakAda
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
	hasil, err := tx.ExecContext(ctx, q, baru, teksAtauNil(posisi), saat, id, models.LiniFacIn, lama)
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

// TutupKasus - penutupan kasus (CloseClaim / SureRejectClaim): STATUS_WORK Resolved-Completed. TAHAP dibiarkan berisi
// tahap terakhir (aturan modul: TAHAP tidak pernah NULL).
func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id, tahap string, saat time.Time) error {
	return g.TutupKasusStatus(ctx, tx, id, tahap, models.StatusSelesai, saat)
}

// TutupKasusStatus - penutupan kasus berstatus `status` (`pxForceCaseClose` Komite Claim Fac In lewat kontrak:
// Resolved-Rejected KomitePost_Reject S17, Resolved-Completed KomitePost_CloseClaim S14).
func (g *Gudang) TutupKasusStatus(ctx context.Context, tx *db.Tx, id, tahap, status string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlTutupKasus(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, status, saat, id, models.LiniFacIn, tahap)
	if err != nil {
		return fmt.Errorf("repository: menutup kasus: %w", err)
	}
	if n, _ := hasil.RowsAffected(); n == 0 {
		return ErrTahapBerubah
	}
	return db.PastikanSatuBaris(hasil, "penutupan kasus")
}

// sqlTutupKomiteAnak - kasus komite KMT- klaim `COVER_KEY` yang masih terbuka; `kecuali` = tambahan `ID <> :6` (tanpa
// itu kosong: `ID <> ”` di Oracle = NULL, tidak pernah benar).
func sqlTutupKomiteAnak(work string, kecuali bool) string {
	q := fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, TGL_UPDATE = :2, POSITION = NULL
		 WHERE COVER_KEY = :3 AND LINI = :4 AND TAHAP = :5 AND STATUS_WORK IS NULL`, work)
	if kecuali {
		q += ` AND ID <> :6`
	}
	return q
}

// TutupKomiteAnak = `CloseAllSubCases=true` (`pxForceCaseClose` KomitePost_Reject S17 / KomitePost_CloseClaim S14,
// `ASMForceCaseClose` CloseClaim 12): kasus komite KMT- klaim `klaimID` yang masih terbuka ditutup berstatus `status`
// (`[inferensi]` status induk), kecuali `kecuali`. Nol baris = tidak ada sub-kasus terbuka (bukan galat).
func (g *Gudang) TutupKomiteAnak(ctx context.Context, tx *db.Tx, klaimID, kecuali, status string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlTutupKomiteAnak(work, kecuali != "")
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	args := []any{status, saat, klaimID, models.LiniFacIn, models.TahapKomite}
	if kecuali != "" {
		args = append(args, kecuali)
	}
	if _, err := tx.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("repository: menutup kasus komite anak: %w", err)
	}
	return nil
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
	hasil, err := tx.ExecContext(ctx, q, saat, id, models.LiniFacIn)
	if err != nil {
		return fmt.Errorf("repository: memperbarui kasus: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "pembaruan kasus")
}

// ---------------------------------------------------------------- daftar kerja

// SaringanKasus - daftar kerja halaman awal.
type SaringanKasus struct {
	// Tahap - satu tahap; kosong = semua tahap terbuka (dipersempit `TahapIn`).
	Tahap string
	// TahapIn - beberapa tahap (worklist pembuat: Input Register + Input Estimasi).
	TahapIn []string
	// Pembuat - CREATE_OP (worklist Assignment1 / Assignment7, ToCurrentOperator / ToWorklist). Kosong = semua.
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
	PolicyNo    string `json:"policyNo"`
	InsuredName string `json:"insuredName"`
	DateOfLoss  string `json:"dateOfLoss"`
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
	w := []string{"w.LINI = " + ph(models.LiniFacIn), "w.TAHAP <> " + ph(models.TahapKomite)}
	if s.Selesai {
		w = append(w, "w.STATUS_WORK = "+ph(models.StatusSelesai))
	} else {
		w = append(w, "w.STATUS_WORK IS NULL")
		if s.Tahap != "" {
			w = append(w, "w.TAHAP = "+ph(s.Tahap))
		}
		if len(s.TahapIn) > 0 {
			var in []string
			for _, t := range s.TahapIn {
				in = append(in, ph(t))
			}
			w = append(w, "w.TAHAP IN ("+strings.Join(in, ", ")+")")
		}
	}
	if s.Pembuat != "" { // tak peka huruf, sama dengan services.Pemegang (strings.EqualFold)
		w = append(w, "UPPER(w.CREATE_OP) = "+ph(strings.ToUpper(s.Pembuat)))
	}
	if c := strings.TrimSpace(s.Cari); c != "" {
		pola := "%" + strings.ToUpper(c) + "%"
		var atau []string
		for _, kol := range []string{"w.ID", "g.CLAIM_NO", "g.POLICY_NO", "g.INSURED_NAME"} {
			atau = append(atau, "UPPER("+kol+") LIKE "+ph(pola))
		}
		w = append(w, "("+strings.Join(atau, " OR ")+")")
	}
	batas := s.Batas
	if batas <= 0 || batas > 500 {
		batas = 500
	}
	q := fmt.Sprintf(`SELECT w.ID, w.TAHAP, w.STATUS_WORK, w.CREATE_OP_NAME, %s,
		       g.CLAIM_NO, g.POLICY_NO, g.INSURED_NAME, TO_CHAR(g.DATE_OF_LOSS, 'YYYY-MM-DD')
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
		return nil, fmt.Errorf("repository: membaca daftar kerja Claim Fac In: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []RingkasanKasus{}
	for rows.Next() {
		var n [9]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8]); err != nil {
			return nil, fmt.Errorf("repository: memindai daftar kerja Claim Fac In: %w", err)
		}
		out = append(out, RingkasanKasus{ID: n[0].String, Tahap: n[1].String, Label: models.LabelTahap[n[1].String],
			StatusWork: n[2].String, PembuatNama: n[3].String, TglCreate: n[4].String, NoClaim: n[5].String,
			PolicyNo: n[6].String, InsuredName: n[7].String, DateOfLoss: n[8].String})
	}
	return out, rows.Err()
}
