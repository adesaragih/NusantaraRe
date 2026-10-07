package repository

// Untuk apa berkas ini: BATAS KOMITE di repository Claim Prop - penulis tangga kasus komite (`AddKomiteTreatyChild_ACT`
// langkah 22-25, `pxAddChildWork` ASM-FW-GCNMFW-Work-KomiteTreaty) dan pembacanya untuk grid "Committe Accept Status".
//
// ⛔ Satu-satunya berkas repository Claim Prop yang menyebut tabel tangga komite. Penjaga batas Claim Life
// `komite_statik_test.go` mengecualikan berkas ini (keputusan work owner 07-10-2026, opsi B). Keputusan anggota tetap
// DITULIS konteks Komite; berkas ini hanya menulis tangga awal (approval menunggu) dan membacanya.

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
)

// tabelTanggaKomite - tangga anggota kasus komite (tabel bersama Komite Claim Life).
const tabelTanggaKomite = "T_KOMITE_KOMITELIST"

// AnggotaTangga - satu anggota tangga komite yang ditulis.
type AnggotaTangga struct {
	Urut                       int
	OperatorID, Jabatan, Email string
}

// sqlSisipKasusKomite / sqlSisipKepalaKomite / sqlSisipAnggotaKomite - kelahiran kasus komite TKMT-.
func sqlSisipKasusKomite(work string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, COVER_KEY, LINI, TAHAP, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)
		VALUES (:1, :2, :3, :4, :5, :6, :7, :8)`, work)
}

func sqlSisipKepalaKomite(gen string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, ADJUSTMENT_ID, KOMITE_LOOP, KOMITE_COUNT) VALUES (:1, :2, :3, :4)`, gen)
}

func sqlSisipAnggotaKomite(list string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, DATA_KOMITE_ID, KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN, KOMITE_EMAIL,
		KOMITE_APPROVAL) VALUES (:1, :2, :3, :4, :5, :6, :7)`, list)
}

// sqlTanggaKomite - anggota tangga satu kasus komite, urut jenjang.
func sqlTanggaKomite(list string) string {
	return fmt.Sprintf(`SELECT ID, KOMITE_OPERATORID, KOMITE_JABATAN, KOMITE_APPROVAL, KOMITE_COMMENT, %s
		  FROM %s WHERE DATA_KOMITE_ID = :1 ORDER BY KOMITE_URUT, ID`, fmt.Sprintf(db.FmtTanggalOracle, "DATE_APPROVE"), list)
}

// sqlSetelKomiteAdjustment - penautan baris adjustment ke kasus komitenya (sekali; KOMITE_ID UNIQUE).
func sqlSetelKomiteAdjustment(adj string) string {
	return fmt.Sprintf(`UPDATE %s SET KOMITE_ID = :1 WHERE ID = :2 AND KOMITE_ID IS NULL`, adj)
}

// BuatKasusKomite = AddKomiteTreatyChild_ACT langkah 22-25: T_WORK_CLAIM TKMT- (COVER_KEY = klaim, LINI PROP, TAHAP
// KomiteTreaty_Flow), T_GENERAL_KOMITE (KOMITE_LOOP = cacah tangga, KOMITE_COUNT 1), tangga satu baris per anggota
// (approval menunggu). Mengembalikan ID kasus komite.
func (g *Gudang) BuatKasusKomite(ctx context.Context, tx *db.Tx, klaimID, adjID, pembuat, namaPembuat string,
	anggota []AnggotaTangga, saat time.Time) (string, error) {
	if len(anggota) == 0 {
		return "", fmt.Errorf("repository: tangga komite kosong; kasus tanpa anggota tidak dapat diputuskan siapa pun")
	}
	id, err := g.IDKasusBerikut(ctx, tx, models.AwalanKomite)
	if err != nil {
		return "", err
	}
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	q := sqlSisipKasusKomite(work)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, id, klaimID, models.LiniProp, models.TahapKomiteTreaty, teksAtauNil(pembuat),
		teksAtauNil(namaPembuat), saat, saat)
	if err != nil {
		return "", fmt.Errorf("repository: melahirkan kasus komite: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran kasus komite"); err != nil {
		return "", err
	}
	gen, err := g.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return "", err
	}
	q = sqlSisipKepalaKomite(gen)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	if hasil, err = tx.ExecContext(ctx, q, id, adjID, len(anggota), 1); err != nil {
		return "", fmt.Errorf("repository: melahirkan baris komite: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran baris komite"); err != nil {
		return "", err
	}
	list, err := g.db.Qualify(tabelTanggaKomite)
	if err != nil {
		return "", err
	}
	q = sqlSisipAnggotaKomite(list)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	for _, a := range anggota {
		rid, err := g.db.NomorBerikut(ctx, tx, "SEQ_KOMITE_KOMITELIST")
		if err != nil {
			return "", err
		}
		hasil, err := tx.ExecContext(ctx, q, rid, id, a.Urut, teksAtauNil(a.OperatorID), teksAtauNil(a.Jabatan),
			teksAtauNil(a.Email), models.ApprovalKomiteMenunggu)
		if err != nil {
			return "", fmt.Errorf("repository: menulis anggota tangga komite: %w", err)
		}
		if err := db.PastikanSatuBaris(hasil, "anggota tangga komite"); err != nil {
			return "", err
		}
	}
	return id, nil
}

// SetelKomiteAdjustment menautkan baris adjustment ke kasus komitenya (T_CLAIM_ADJUSTMENT.KOMITE_ID, UNIQUE).
func (g *Gudang) SetelKomiteAdjustment(ctx context.Context, tx *db.Tx, adjID, komiteID string) error {
	tabel, err := g.db.Qualify(models.TabelAdjustment.Nama)
	if err != nil {
		return err
	}
	q := sqlSetelKomiteAdjustment(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, komiteID, adjID)
	if err != nil {
		return fmt.Errorf("repository: menautkan adjustment ke komite: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penautan adjustment ke komite")
}

// TanggaKomite - anggota tangga kasus komite untuk grid "Committe Accept Status".
func (a *Acuan) TanggaKomite(ctx context.Context, komiteID string) ([]models.AnggotaKomite, error) {
	t, err := a.q(tabelTanggaKomite)
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, sqlTanggaKomite(t), 6, komiteID)
	if err != nil {
		return nil, err
	}
	out := []models.AnggotaKomite{}
	for _, r := range rows {
		out = append(out, models.AnggotaKomite{ID: r[0], OperatorID: r[1], Jabatan: r[2], Approval: r[3], Comment: r[4],
			TanggalSetuju: r[5]})
	}
	return out, nil
}
