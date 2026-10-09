package repository

// Untuk apa berkas ini: BACAAN SUMBER DATA LAMA dan kelahiran kasus hasil pemuat (`services.Pemuat`,
// `backend/alat/pemuatlama`). Kedua sumber dibaca saja; tulisannya hanya T_WORK_CLAIM + T_GENERAL_CLAIM (sisanya
// `SimpanHalaman`, jalur yang sama dengan aplikasi).

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/claimprop/backend/models"
)

// sqlOSLama - seluruh baris OS_AKSEPTASI_KLAIM kasus Claim Prop lama. AcceptedNo dibaca dari DATA_JSON (kunci urutan
// AC 123); DATA_JSON baris itu sendiri JSON datar satu baris aksi.
func sqlOSLama(t string) string {
	return fmt.Sprintf(`SELECT CASEID, NOCLAIM, NOPOLIS, MASTERID, INSERT_OP, TO_CHAR(STS_REJECT),
		JSON_VALUE(DATA_JSON, '$.AcceptedNo' RETURNING VARCHAR2(100)), DATA_JSON, %s
		  FROM %s WHERE CASEID LIKE :1 ORDER BY CASEID, TANGGAL`, fmt.Sprintf(db.FmtTanggalOracle, "TANGGAL"), t)
}

// BacaOSLama membaca baris OS_AKSEPTASI_KLAIM kasus lama (CASEID `ASM-FW-GCNMFW-WORK CLMP-%`).
func (g *Gudang) BacaOSLama(ctx context.Context) ([]models.BarisOSLama, error) {
	t, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return nil, err
	}
	q := sqlOSLama(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, models.AwalanKunciLama+"%")
	if err != nil {
		return nil, fmt.Errorf("repository: membaca OS_AKSEPTASI_KLAIM lama: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.BarisOSLama
	for rows.Next() {
		var n [9]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &n[6], &n[7], &n[8]); err != nil {
			return nil, fmt.Errorf("repository: memindai OS_AKSEPTASI_KLAIM lama: %w", err)
		}
		b := models.BarisOSLama{CaseID: n[0].String, NoClaim: n[1].String, NoPolis: n[2].String, MasterID: n[3].String,
			InsertOp: n[4].String, StsReject: n[5].String, AcceptedNo: n[6].String, DataJSON: n[7].String}
		if t, err := utils.ParseTanggal(n[8].String); err == nil {
			b.Tanggal = t
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// sqlJSONKlaimLama - halaman JSON_KLAIM kasus lama.
func sqlJSONKlaimLama(t string) string {
	return fmt.Sprintf(`SELECT IDPEGA, DATA_JSON FROM %s WHERE IDPEGA LIKE :1`, t)
}

// BacaJSONKlaimLama membaca halaman JSON_KLAIM kasus lama: IDPEGA -> DATA_JSON.
func (g *Gudang) BacaJSONKlaimLama(ctx context.Context) (map[string]string, error) {
	t, err := g.db.Qualify("JSON_KLAIM")
	if err != nil {
		return nil, err
	}
	q := sqlJSONKlaimLama(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, models.AwalanKunciLama+"%")
	if err != nil {
		return nil, fmt.Errorf("repository: membaca JSON_KLAIM lama: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]string{}
	for rows.Next() {
		var id, dok sql.NullString
		if err := rows.Scan(&id, &dok); err != nil {
			return nil, fmt.Errorf("repository: memindai JSON_KLAIM lama: %w", err)
		}
		if dok.String != "" {
			out[id.String] = dok.String
		}
	}
	return out, rows.Err()
}

// sqlSisipKasusLama - kelahiran baris T_WORK_CLAIM kasus lama (STATUS_WORK ikut ditulis).
func sqlSisipKasusLama(work string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, TAHAP, POSITION, STATUS_WORK, CREATE_OP, CREATE_OP_NAME, TGL_CREATE,
		TGL_UPDATE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)`, work)
}

// SisipKasusLama melahirkan kasus hasil pemuat: ID pyID Pega, tahap / posisi / status dari baris berlaku, SUMBER
// 'PEGA' (AC 128-129: klaim hasil migrasi tidak dapat dihapus).
func (g *Gudang) SisipKasusLama(ctx context.Context, tx *db.Tx, k models.Kasus) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlSisipKasusLama(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, k.ID, models.LiniProp, k.Tahap, teksAtauNil(k.Posisi), teksAtauNil(k.StatusWork),
		teksAtauNil(k.PembuatID), teksAtauNil(k.PembuatNama), k.TglCreate, k.TglUpdate)
	if err != nil {
		return fmt.Errorf("repository: melahirkan kasus lama %s: %w", k.ID, err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran kasus lama"); err != nil {
		return err
	}
	return g.sisipInduk(ctx, tx, k.ID, models.SumberPega)
}
