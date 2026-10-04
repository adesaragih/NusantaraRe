package repository

// Penulis `TREATYBUSINESS_LIFE` - tiruan `INSERTBUSINESS_LIFE` `[data DBA]`
// (upsert dikunci `ID`) dan `SaveTreatyBusinessAll_Life_SQL.xml` b84 (INSERT
// per sasaran): ID baru `'1'||lpad(to_Char(TREATYBUSINESS_LIFE_SEQ.nextval),6,'0')`,
// `TGLUPDATE = SYSDATE`.
//
// ⛔ `RIRATE` VARCHAR2(1000) TEKS apa adanya (R7). ⛔ Kunci induk tidak pernah
// diperbarui; salinan tahun/jenis ditulis dari induk (K4).

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func sqlSisipBusiness(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TREATYYEARID, TREATYYEAR, TREATYCONTRACTID, REINSTYPEID, REINSTYPENAME,
	       BIZCODE, BIZNAME, RIRATEID, RIRATE, USERID, TGLUPDATE)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, SYSDATE)`, t)
}

func sqlPerbaruiBusiness(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET TREATYYEAR = :1, REINSTYPEID = :2, REINSTYPENAME = :3, BIZCODE = :4, BIZNAME = :5,
	       RIRATEID = :6, RIRATE = :7, USERID = :8, TGLUPDATE = SYSDATE
	 WHERE ID = :9`, t)
}

// SisipBusiness menyisipkan business baru; ID dari sequence.
func (g *Gudang) SisipBusiness(ctx context.Context, tx *db.Tx, b models.Business) (string, error) {
	id, err := g.identitasBaru(ctx, tx, SeqBusiness)
	if err != nil {
		return "", err
	}
	q, err := g.siapkan(TabelBusiness, sqlSisipBusiness)
	if err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, id, b.TreatyYearID, b.TreatyYear, b.TreatyContractID, b.ReinsTypeID,
		b.ReinsTypeName, b.BizCode, b.BizName, b.RIRateID, b.RIRate, b.UserID)
	if err != nil {
		return "", fmt.Errorf("repository: inserting business: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "business insert")
}

// PerbaruiBusiness memperbarui satu business; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiBusiness(ctx context.Context, tx *db.Tx, b models.Business) error {
	q, err := g.siapkan(TabelBusiness, sqlPerbaruiBusiness)
	if err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, b.TreatyYear, b.ReinsTypeID, b.ReinsTypeName, b.BizCode, b.BizName,
		b.RIRateID, b.RIRate, b.UserID, b.ID)
	if err != nil {
		return fmt.Errorf("repository: updating business %s: %w", b.ID, err)
	}
	return adaSatu(hasil, "business", b.ID)
}
