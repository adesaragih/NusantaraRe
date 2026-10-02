package repository

// ⛔ BATAS LINTAS SISTEM - satu-satunya berkas modul ini yang membaca skema
// Arasapas (spec §13 penyimpangan sadar 4, AC 47). Kueri lintas skema tidak
// boleh tersebar; `TestArasapasHanyaDiSatuBerkas` menagihnya.
//
// Gerbang ke-5 `SetErrorBatalEndorsement_Act` 3.9 b2035 →
// `RDBList/SearcStatusBayarArasaps_SQL.xml` b58
// `select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no ={InputData.CARI18} and IVD_JR_ID ='5'`
// - `IVD_JR_ID = '5'` = pembayaran/pelunasan (spec §13 `[keputusan work owner]`).
// Nomor invoice = nomor polis tanpa titik (b464, `models.NomorInvoiceArasapas`).
//
// ⚠️ Yang dibutuhkan hanya ADA/TIDAK - nol `SELECT *`, satu baris paling banyak.
// Kode jurnal dikirim sebagai penampung, bukan literal.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

func sqlSudahDibayar() string {
	return `SELECT 1 FROM ` + tabelArasapas + ` WHERE INV_INV_NO = :1 AND IVD_JR_ID = :2 FETCH FIRST 1 ROWS ONLY`
}

// SudahDibayar - gerbang 5: ada baris pelunasan untuk nomor invoice itu.
func (g *Gudang) SudahDibayar(ctx context.Context, nomorInvoice string) (bool, error) {
	q := sqlSudahDibayar()
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var satu int
	err := g.db.QueryRowContext(ctx, q, nomorInvoice, models.KodeJurnalPelunasan).Scan(&satu)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrArasapasTakTerbaca, err)
	}
	return true, nil
}

// ErrArasapasTakTerbaca - pembacaan lintas skema gagal (hak baca, OQ-EDM-012).
var ErrArasapasTakTerbaca = fmt.Errorf("repository: pembayaran Arasapas tidak dapat dibaca")
