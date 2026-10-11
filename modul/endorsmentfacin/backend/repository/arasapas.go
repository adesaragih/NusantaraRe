package repository

// ⛔ BATAS LINTAS SKEMA - satu-satunya berkas modul ini yang membaca skema Arasapas.
//
// Gerbang 4 `SetErrorBatalEndorsement_Act` langkah 5.15 → `RDBList/SearcStatusBayarArasaps_SQL.xml`
// `[terverifikasi]`: `select * from ARASAPAS.DETAIL_INVOICE where inv_inv_no ={InputData.CARI18} and IVD_JR_ID ='5'`.
// `InputData.CARI18 = @replaceAll(.PolicyNo,".","")` (langkah 2). Gerbang hanya menguji `@LengthOfPageList(…)>0` -
// cukup ada/tidak; nol `SELECT *`. Kode jurnal dikirim sebagai penampung, bukan literal SQL.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
)

// tabelArasapas - skema lain, ditulis tetap seperti di rule (bukan `Qualify` - itu skema aplikasi).
const tabelArasapas = "ARASAPAS.DETAIL_INVOICE"

// KodeJurnalPembayaran - `IVD_JR_ID ='5'` rule asal. Arti kodenya `belum terverifikasi` di korpus EDM Fac In.
const KodeJurnalPembayaran = "5"

// ErrSkemaLainTakTerbaca - ARASAPAS tidak terbaca dari koneksi ini (ORA-00942 tabel/sinonim tidak ada, ORA-01031 hak
// baca). Gerbangnya "belum-diperiksa", bukan galat seluruh pemeriksaan.
var ErrSkemaLainTakTerbaca = errors.New("repository: " + tabelArasapas + " tidak terbaca dari koneksi ini")

// SqlAdaPembayaran - ada baris jurnal pembayaran untuk nomor invoice itu.
func SqlAdaPembayaran() string {
	return `SELECT 1 FROM ` + tabelArasapas + ` WHERE INV_INV_NO = :1 AND IVD_JR_ID = :2 FETCH FIRST 1 ROWS ONLY`
}

// NomorInvoice - `@replaceAll(.PolicyNo,".","")`.
func NomorInvoice(polis string) string { return strings.ReplaceAll(polis, ".", "") }

// AdaPembayaran - gerbang 4.
func (g *Gudang) AdaPembayaran(ctx context.Context, polis string) (bool, error) {
	q := SqlAdaPembayaran()
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var satu int
	err := g.db.QueryRowContext(ctx, q, NomorInvoice(polis), KodeJurnalPembayaran).Scan(&satu)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, galatArasapas(err)
	}
	return true, nil
}

// galatArasapas - ORA-00942 / ORA-01031 → ErrSkemaLainTakTerbaca; selain itu galat biasa.
func galatArasapas(err error) error {
	if s := err.Error(); strings.Contains(s, "ORA-00942") || strings.Contains(s, "ORA-01031") {
		return fmt.Errorf("%w: %v", ErrSkemaLainTakTerbaca, err)
	}
	return fmt.Errorf("repository: membaca pembayaran polis (SearcStatusBayarArasaps_SQL): %w", err)
}
