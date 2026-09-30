package repository

// Penulis `TREATYREINSURER_LIFE` - tiruan `INSERTREINSURER_LIFE` `[data DBA]`
// (upsert dikunci `ID`, ID baru `TREATYREINSURER_LIFE_SEQ`, `TGLUPDATE =
// SYSDATE`), dan pembaca total share per kontrak (tiket 11).
//
// ⛔ `TREATYYEARID`/`TREATYCONTRACTID` tidak pernah diperbarui - reinsurer
// tidak berpindah kontrak; salinan jenis ditulis dari kontrak (K4).

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func sqlSisipReinsurer(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, TREATYYEARID, TREATYCONTRACTID, REINSTYPEID, REINSTYPENAME,
	       REINSURERID, REINSURERNAME, USERID, TGLUPDATE, PCTSHARE, COMMISION, OVR_COMM)
	VALUES (:1, :2, :3, :4, :5, :6, :7, :8, SYSDATE, %s, %s, %s)`, t,
		angka(":9", ":10"), angka(":11", ":12"), angka(":13", ":14"))
}

func sqlPerbaruiReinsurer(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET REINSTYPEID = :1, REINSTYPENAME = :2, REINSURERID = :3, REINSURERNAME = :4, USERID = :5,
	       TGLUPDATE = SYSDATE, PCTSHARE = %s, COMMISION = %s, OVR_COMM = %s
	 WHERE ID = :12`, t, angka(":6", ":7"), angka(":8", ":9"), angka(":10", ":11"))
}

func argPersenReinsurer(args []any, r models.Reinsurer) []any {
	args = argAngka(args, r.PctShare)
	args = argAngka(args, r.Komisi)
	return argAngka(args, r.OvrComm)
}

// SisipReinsurer menyisipkan reinsurer baru; ID dari sequence.
func (g *Gudang) SisipReinsurer(ctx context.Context, tx *db.Tx, r models.Reinsurer) (string, error) {
	id, err := g.identitasBaru(ctx, tx, SeqReinsurer)
	if err != nil {
		return "", err
	}
	q, err := g.siapkan(TabelReinsurer, sqlSisipReinsurer)
	if err != nil {
		return "", err
	}
	args := argPersenReinsurer([]any{id, r.TreatyYearID, r.TreatyContractID, r.ReinsTypeID, r.ReinsTypeName,
		r.ReinsurerID, r.ReinsurerName, r.UserID}, r)
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return "", fmt.Errorf("repository: inserting reinsurer: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "reinsurer insert")
}

// PerbaruiReinsurer memperbarui satu reinsurer; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiReinsurer(ctx context.Context, tx *db.Tx, r models.Reinsurer) error {
	q, err := g.siapkan(TabelReinsurer, sqlPerbaruiReinsurer)
	if err != nil {
		return err
	}
	args := argPersenReinsurer([]any{r.ReinsTypeID, r.ReinsTypeName, r.ReinsurerID, r.ReinsurerName, r.UserID}, r)
	hasil, err := tx.ExecContext(ctx, q, append(args, r.ID)...)
	if err != nil {
		return fmt.Errorf("repository: updating reinsurer %s: %w", r.ID, err)
	}
	return adaSatu(hasil, "reinsurer", r.ID)
}

// sqlTotalSharePerKontrak - LEFT JOIN: kontrak tanpa reinsurer ikut (total NULL).
// Kunci gabung = dua kunci induk RD reinsurer (`TREATYYEARID`, `TREATYCONTRACTID`).
func sqlTotalSharePerKontrak(kontrak, reinsurer, tahun string, saringTahun bool) string {
	saring := ""
	if saringTahun {
		saring = "WHERE k.IDTREATYYEAR = :1"
	}
	return fmt.Sprintf(`SELECT k.ID, k.IDTREATYYEAR, y.TREATYYEAR, k.REINSTYPENAME,
	       TO_CHAR(SUM(r.PCTSHARE), 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')
	  FROM %s k
	  LEFT JOIN %s r ON r.TREATYCONTRACTID = k.ID AND r.TREATYYEARID = k.IDTREATYYEAR
	  LEFT JOIN %s y ON y.ID = k.IDTREATYYEAR
	 %s
	 GROUP BY k.ID, k.IDTREATYYEAR, y.TREATYYEAR, k.REINSTYPENAME
	 ORDER BY k.IDTREATYYEAR ASC, k.ID ASC`, kontrak, reinsurer, tahun, saring)
}

// TotalSharePerKontrak - lihat services.Gudang.
func (g *Gudang) TotalSharePerKontrak(ctx context.Context, tahunID string) ([]models.TotalShareKontrak, error) {
	var nama [3]string
	for i, o := range []string{TabelKontrak, TabelReinsurer, TabelTahun} {
		n, err := g.db.Qualify(o)
		if err != nil {
			return nil, err
		}
		nama[i] = n
	}
	q := sqlTotalSharePerKontrak(nama[0], nama[1], nama[2], tahunID != "")
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	var args []any
	if tahunID != "" {
		args = append(args, tahunID)
	}
	kolom := []string{"ID", "IDTREATYYEAR", "TREATYYEAR", "REINSTYPENAME", "TOTAL"}
	return daftar(ctx, g.db, q, kolom, func(b barisTeks) (models.TotalShareKontrak, error) {
		total, err := db.UraiDesimal(b.s("ID"), "SUM(PCTSHARE)", b["TOTAL"])
		return models.TotalShareKontrak{KontrakID: b.s("ID"), TahunID: b.s("IDTREATYYEAR"),
			TreatyYear: b.s("TREATYYEAR"), ReinsTypeName: b.s("REINSTYPENAME"), Total: total}, err
	}, "share total", args...)
}
