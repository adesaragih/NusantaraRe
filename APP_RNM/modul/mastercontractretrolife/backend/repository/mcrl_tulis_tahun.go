package repository

// Penulis `TREATYYEAR_LIFE` - tiruan `INSERTTREATYYEAR_LIFE` `[data DBA]`:
// upsert dikunci `ID`, ID baru dari `TREATYYEAR_LIFE_SEQ`, `TGLUPDATE =
// SYSDATE` (parameter cap waktu diabaikan procedure), `USERID` dari pemanggil.
// Upsert dipecah dua (Sisip/Perbarui): pemanggil yang tahu mana yang dimaksud.
//
// ⛔ Tanggal ditulis sebagai DATE tanpa jam (`TO_DATE(…, 'YYYY-MM-DD')`), sama
// dengan `To_date(CARI6, 'DD/MM/YYYY HH24:MI:SS')` dari `@FormatDateTime(…,
// "dd/MM/yyyy")` Pega (`SaveTreatyYearLife_Act` b823/b852).
// ⛔ Nol COMMIT: transaksinya milik services.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// bindTanggal membind tanggal sebagai teks YYYY-MM-DD; nol = NULL.
func bindTanggal(t time.Time) any { return db.KosongJadiNil(utils.FormatTanggal(t)) }

const tglOracle = `TO_DATE(%s, 'YYYY-MM-DD')`

func tgl(ph string) string { return fmt.Sprintf(tglOracle, ph) }

func sqlSisipTahun(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TREATYYEAR, UNDERWRITINGYEAR, USERID, TGLUPDATE, STARTDATE, ENDDATE)
	VALUES (:1, :2, :3, :4, SYSDATE, %s, %s)`, t, tgl(":5"), tgl(":6"))
}

func sqlPerbaruiTahun(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET TREATYYEAR = :1, UNDERWRITINGYEAR = :2, USERID = :3, TGLUPDATE = SYSDATE,
	       STARTDATE = %s, ENDDATE = %s
	 WHERE ID = :6`, t, tgl(":4"), tgl(":5"))
}

// sqlSalinTahunKeBusiness - K4: salinan `TREATYYEAR` business; hanya baris yang
// BERBEDA (`DECODE` menganggap NULL = NULL) yang disentuh.
func sqlSalinTahunKeBusiness(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET TREATYYEAR = :1, USERID = :2, TGLUPDATE = SYSDATE
	 WHERE TREATYYEARID = :3 AND DECODE(TREATYYEAR, :4, 0, 1) = 1`, t)
}

// sqlSalinTahunKeKontrak - R5: tanggal kontrak = salinan tanggal tahun.
func sqlSalinTahunKeKontrak(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET TREATYSTARTDATE = %s, TREATYENDDATE = %s, USERID = :3, TGLUPDATE = SYSDATE
	 WHERE IDTREATYYEAR = :4
	   AND (DECODE(TREATYSTARTDATE, %s, 0, 1) = 1 OR DECODE(TREATYENDDATE, %s, 0, 1) = 1)`,
		t, tgl(":1"), tgl(":2"), tgl(":5"), tgl(":6"))
}

// SisipTahun menyisipkan tahun baru dan mengembalikan ID dari sequence.
func (g *Gudang) SisipTahun(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (string, error) {
	id, err := g.identitasBaru(ctx, tx, SeqTahun)
	if err != nil {
		return "", err
	}
	q, err := g.siapkan(TabelTahun, sqlSisipTahun)
	if err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, id, t.TreatyYear, t.UnderwritingYear, t.UserID,
		bindTanggal(t.StartDate), bindTanggal(t.EndDate))
	if err != nil {
		return "", fmt.Errorf("repository: inserting treaty year: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "treaty year insert")
}

// PerbaruiTahun memperbarui seluruh medan satu tahun; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiTahun(ctx context.Context, tx *db.Tx, t models.TahunTreaty) error {
	q, err := g.siapkan(TabelTahun, sqlPerbaruiTahun)
	if err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, t.TreatyYear, t.UnderwritingYear, t.UserID,
		bindTanggal(t.StartDate), bindTanggal(t.EndDate), t.ID)
	if err != nil {
		return fmt.Errorf("repository: updating treaty year %s: %w", t.ID, err)
	}
	return adaSatu(hasil, "treaty year", t.ID)
}

// SalinTahunKeAnak - lihat services.Gudang.
func (g *Gudang) SalinTahunKeAnak(ctx context.Context, tx *db.Tx, t models.TahunTreaty) (int64, error) {
	qb, err := g.siapkan(TabelBusiness, sqlSalinTahunKeBusiness)
	if err != nil {
		return 0, err
	}
	hb, err := tx.ExecContext(ctx, qb, t.TreatyYear, t.UserID, t.ID, t.TreatyYear)
	if err != nil {
		return 0, fmt.Errorf("repository: copying treaty year to business rows: %w", err)
	}
	qk, err := g.siapkan(TabelKontrak, sqlSalinTahunKeKontrak)
	if err != nil {
		return 0, err
	}
	hk, err := tx.ExecContext(ctx, qk, bindTanggal(t.StartDate), bindTanggal(t.EndDate), t.UserID, t.ID,
		bindTanggal(t.StartDate), bindTanggal(t.EndDate))
	if err != nil {
		return 0, fmt.Errorf("repository: copying treaty year dates to contracts: %w", err)
	}
	return jumlahBaris(hb, hk)
}

// adaSatu - pernyataan tulis ber-ID menyentuh tepat satu baris; nol = ErrTidakAda.
func adaSatu(hasil interface{ RowsAffected() (int64, error) }, nama, id string) error {
	n, err := hasil.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: counting %s rows: %w", nama, err)
	}
	switch {
	case n == 0:
		return fmt.Errorf("%w: %s %s", ErrTidakAda, nama, id)
	case n > 1:
		// K3: `ID` tidak dijamin unik oleh DB (nol PK) - dua baris ber-ID sama
		// adalah keadaan data rusak, bukan pembaruan yang sah.
		return fmt.Errorf("repository: %s %s touched %d rows; ID must be unique: %w", nama, id, n, errIDGanda)
	}
	return nil
}

var errIDGanda = errors.New("duplicate ID")

func jumlahBaris(hasil ...interface{ RowsAffected() (int64, error) }) (int64, error) {
	var total int64
	for _, h := range hasil {
		n, err := h.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("repository: counting rows: %w", err)
		}
		total += n
	}
	return total, nil
}
