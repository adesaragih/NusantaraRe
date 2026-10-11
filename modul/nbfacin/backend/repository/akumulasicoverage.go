package repository

// PUT /api/nbfacin/kasus/{caseId}/akumulasi (tiket 60): tahap uw hanya mengubah Accumulation Code / Address coverage
// (T_COVERAGELIST.ACCUMULATION_CODE / ACCUMULATION_DESCRIPTION, migrasi 195) - bukan ganti utuh objek seperti PUT …/objek.

import (
	"context"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// PenyimpanAkumulasiCoverage - tulis akumulasi coverage case.
type PenyimpanAkumulasiCoverage interface {
	// UbahAkumulasiCoverage - menyentuh TGL_UPDATE case lalu mengubah dua kolom akumulasi tiap coverage ber-ID (ID dibaca
	// pemanggil dari case yang SAMA di transaksi ini). ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	UbahAkumulasiCoverage(ctx context.Context, tx *db.Tx, id string, baris []models.AkumulasiCoverage) error
}

func sqlUbahAkumulasiCoverage(cov string) string {
	return "UPDATE " + cov + " SET ACCUMULATION_CODE = :1, ACCUMULATION_DESCRIPTION = :2 WHERE ID = :3"
}

// UbahAkumulasiCoverage - lihat PenyimpanAkumulasiCoverage. Teks kosong -> NULL (sama dengan PUT …/objek, kolomData).
func (r *ObjekOracle) UbahAkumulasiCoverage(ctx context.Context, tx *db.Tx, id string, baris []models.AkumulasiCoverage) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	for _, b := range baris {
		if _, err := jalankan(ctx, tx, sqlUbahAkumulasiCoverage(t.cov), "menulis akumulasi coverage",
			db.KosongJadiNil(b.AccumulationCode), db.KosongJadiNil(b.AccumulationDescription), b.ID); err != nil {
			return err
		}
	}
	return nil
}
