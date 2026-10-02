package repository

// Penulis `TREATYSECURITYREINSURER_LIFE` - tiruan `INSERTSECURITYREINSURER_LIFE`
// `[data DBA]` (upsert dikunci `ID`, ID baru `TREATYSECURITYREINSURER_LIFE_SEQ`,
// `TGLUPDATE = SYSDATE`). Kunci induk (`TREATYYEARID`, `TREATYCONTRACTID`,
// `TREATYREINSURERID`) ditulis saat lahir dan tidak pernah diperbarui.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func sqlSisipSecurity(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TREATYYEARID, TREATYCONTRACTID, TREATYREINSURERID, REINSURERID,
	       REINSURERNAME, USERID, TGLUPDATE, PCTSHARE)
	VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE, %s)`, t, angka(":8", ":9"))
}

func sqlPerbaruiSecurity(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET REINSURERID = :1, REINSURERNAME = :2, USERID = :3, TGLUPDATE = SYSDATE, PCTSHARE = %s
	 WHERE ID = :6`, t, angka(":4", ":5"))
}

// SisipSecurity menyisipkan security baru; ID dari sequence.
func (g *Gudang) SisipSecurity(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) (string, error) {
	id, err := g.identitasBaru(ctx, tx, SeqSecurity)
	if err != nil {
		return "", err
	}
	q, err := g.siapkan(TabelSecurity, sqlSisipSecurity)
	if err != nil {
		return "", err
	}
	args := argAngka([]any{id, s.TreatyYearID, s.TreatyContractID, s.TreatyReinsurerID, s.ReinsurerID,
		s.ReinsurerName, s.UserID}, s.PctShare)
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return "", fmt.Errorf("repository: inserting security reinsurer: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "security reinsurer insert")
}

// PerbaruiSecurity memperbarui satu security; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiSecurity(ctx context.Context, tx *db.Tx, s models.SecurityReinsurer) error {
	q, err := g.siapkan(TabelSecurity, sqlPerbaruiSecurity)
	if err != nil {
		return err
	}
	args := argAngka([]any{s.ReinsurerID, s.ReinsurerName, s.UserID}, s.PctShare)
	hasil, err := tx.ExecContext(ctx, q, append(args, s.ID)...)
	if err != nil {
		return fmt.Errorf("repository: updating security reinsurer %s: %w", s.ID, err)
	}
	return adaSatu(hasil, "security reinsurer", s.ID)
}
