package repository

// Penulis `TREATYCONTRACT_LIFE` - tiruan `INSERTTREATYCONTRACT_LIFE` `[data DBA]`
// (upsert dikunci `ID`, ID baru `TREATYCONTRACT_LIFE_SEQ`, `TGLUPDATE =
// SYSDATE`) ditambah K4: salinan jenis ke reinsurer dan business.
//
// ⛔ Uang lewat `angka` (kebal NLS); `USD` dan `USD_SELISIH` boleh NULL.
// ⛔ `IDTREATYYEAR` TIDAK pernah diperbarui - kontrak tidak berpindah tahun.

import (
	"context"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

func sqlSisipKontrak(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, IDTREATYYEAR, REINSTYPEID, REINSTYPENAME, USERID, TGLUPDATE,
	       IDR, USD, B_IDR, B_USD, IDR_SELISIH, USD_SELISIH, TREATYSTARTDATE, TREATYENDDATE)
	VALUES (:1, :2, :3, :4, :5, SYSDATE, %s, %s, %s, %s, %s, %s, %s, %s)`, t,
		angka(":6", ":7"), angka(":8", ":9"), angka(":10", ":11"), angka(":12", ":13"),
		angka(":14", ":15"), angka(":16", ":17"), tgl(":18"), tgl(":19"))
}

func sqlPerbaruiKontrak(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET REINSTYPEID = :1, REINSTYPENAME = :2, USERID = :3, TGLUPDATE = SYSDATE,
	       IDR = %s, USD = %s, B_IDR = %s, B_USD = %s, IDR_SELISIH = %s, USD_SELISIH = %s,
	       TREATYSTARTDATE = %s, TREATYENDDATE = %s
	 WHERE ID = :18`, t, angka(":4", ":5"), angka(":6", ":7"), angka(":8", ":9"), angka(":10", ":11"),
		angka(":12", ":13"), angka(":14", ":15"), tgl(":16"), tgl(":17"))
}

// sqlSalinJenisKeAnak - K4: `REINSTYPEID`/`REINSTYPENAME` anak = kontrak.
func sqlSalinJenisKeAnak(t string) string {
	return fmt.Sprintf(`UPDATE %s
	   SET REINSTYPEID = :1, REINSTYPENAME = :2, USERID = :3, TGLUPDATE = SYSDATE
	 WHERE TREATYCONTRACTID = :4 AND (DECODE(REINSTYPEID, :5, 0, 1) = 1 OR DECODE(REINSTYPENAME, :6, 0, 1) = 1)`, t)
}

// argUang - enam kolom uang kontrak, urutan penampung SQL :6..:17 (sisip) / :4..:15 (ubah).
func argUang(args []any, k models.Kontrak) []any {
	args = argAngka(args, k.IDR)
	args = argAngka(args, k.USD)
	args = argAngka(args, k.BIDR)
	args = argAngka(args, k.BUSD)
	args = argAngka(args, k.IDRSelisih)
	return argAngka(args, k.USDSelisih)
}

// SisipKontrak menyisipkan kontrak baru dan mengembalikan ID dari sequence.
func (g *Gudang) SisipKontrak(ctx context.Context, tx *db.Tx, k models.Kontrak) (string, error) {
	id, err := g.identitasBaru(ctx, tx, SeqKontrak)
	if err != nil {
		return "", err
	}
	q, err := g.siapkan(TabelKontrak, sqlSisipKontrak)
	if err != nil {
		return "", err
	}
	args := argUang([]any{id, k.IDTreatyYear, k.ReinsTypeID, k.ReinsTypeName, k.UserID}, k)
	args = append(args, bindTanggal(k.TreatyStartDate), bindTanggal(k.TreatyEndDate))
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return "", fmt.Errorf("repository: inserting treaty contract: %w", err)
	}
	return id, db.PastikanSatuBaris(hasil, "treaty contract insert")
}

// PerbaruiKontrak memperbarui seluruh medan kontrak; nol baris = ErrTidakAda.
func (g *Gudang) PerbaruiKontrak(ctx context.Context, tx *db.Tx, k models.Kontrak) error {
	q, err := g.siapkan(TabelKontrak, sqlPerbaruiKontrak)
	if err != nil {
		return err
	}
	args := argUang([]any{k.ReinsTypeID, k.ReinsTypeName, k.UserID}, k)
	args = append(args, bindTanggal(k.TreatyStartDate), bindTanggal(k.TreatyEndDate), k.ID)
	hasil, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: updating treaty contract %s: %w", k.ID, err)
	}
	return adaSatu(hasil, "treaty contract", k.ID)
}

// SalinKontrakKeAnak - lihat services.Gudang.
func (g *Gudang) SalinKontrakKeAnak(ctx context.Context, tx *db.Tx, k models.Kontrak) (int64, error) {
	var hasil []interface{ RowsAffected() (int64, error) }
	for _, tabel := range []string{TabelReinsurer, TabelBusiness} {
		q, err := g.siapkan(tabel, sqlSalinJenisKeAnak)
		if err != nil {
			return 0, err
		}
		h, err := tx.ExecContext(ctx, q, k.ReinsTypeID, k.ReinsTypeName, k.UserID, k.ID, k.ReinsTypeID, k.ReinsTypeName)
		if err != nil {
			return 0, fmt.Errorf("repository: copying reinsurance type to %s: %w", tabel, err)
		}
		hasil = append(hasil, h)
	}
	return jumlahBaris(hasil...)
}
