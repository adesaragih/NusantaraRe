package repository

// Coverage Note tab Coverage Aneka (tiket 73 P2) `[terverifikasi]`: sel autocomplete `.CoverageNote`
// `NB FacIn\Section\InputCoverageAneka_FacIn.xml` L1709 -> RD `BrowseCoverageFacIn_RD` (sama dengan popup FIRE tiket 43, lihat
// pilihancoverage.go) dengan Type = "ANEKA" (L1891) dan BizCode = `pyWorkPage.OfferFacIn.QuotationData.BusinessCode` (L1898);
// Nama tidak dikirim (Contains diabaikan di Pega - di sini `q` dipakai sebagai kata pencarian autocomplete, pola FIRE); pilih:
// `.ID` -> `.Coverage`, NamaCoverage -> `.CoverageNote`.
// Sumber = tabel COVERAGE (keputusan work owner butir 96 untuk RD yang sama; kelas RD Int-COVERAGE_FACIN = view
// `DDL\COVERAGE_FACIN.txt` atas JSON M_COVERAGE yang TIDAK punya ACTIVESTATUS) - BizCode = BUSINESSCODE, peka huruf.

import (
	"context"
	"fmt"
	"strconv"

	"nusantarare/modul/nbfacin/backend/models"
)

// TipeCoverageAneka - Param.Type sel Coverage Note Aneka.
const TipeCoverageAneka = "ANEKA"

// sqlCariCoverageAneka - :1 tipe, :2 status aktif, [BUSINESSCODE], [pola NAME], batas terakhir; bentuk DISTINCT / urutan sama
// dengan sqlCariCoverage (FIRE).
func sqlCariCoverageAneka(c string, denganBisnis, denganKata bool) string {
	syarat, n := "TYPE = :1 AND (ACTIVESTATUS = :2 OR ACTIVESTATUS IS NULL)", 3
	if denganBisnis {
		syarat += " AND BUSINESSCODE = :" + strconv.Itoa(n)
		n++
	}
	if denganKata {
		syarat += " AND UPPER(NAME) LIKE :" + strconv.Itoa(n) + ` ESCAPE '\'`
		n++
	}
	return "SELECT ID, OLDID, NAME FROM (SELECT DISTINCT ID, BUSINESSCODE, NAME, TYPE, OLDID, ACTIVESTATUS FROM " + c +
		" WHERE " + syarat + ") ORDER BY NAME, OLDID, ID FETCH FIRST :" + strconv.Itoa(n) + " ROWS ONLY"
}

// CariCoverageAneka - lihat PembacaLookupAneka.
func (r *LookupAnekaOracle) CariCoverageAneka(ctx context.Context, kodeBisnis, kata string) ([]models.BarisCoverage, error) {
	q, err := r.db.Qualify(TabelCoverage)
	if err != nil {
		return nil, err
	}
	arg := []any{TipeCoverageAneka, StatusAktifCoverage}
	if kodeBisnis != "" {
		arg = append(arg, kodeBisnis)
	}
	pola := PolaCari(kata)
	if pola != "" {
		arg = append(arg, pola)
	}
	arg = append(arg, BatasCoverage)
	baris, err := r.db.QueryContext(ctx, sqlCariCoverageAneka(q, kodeBisnis != "", pola != ""), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCoverage, err)
	}
	defer baris.Close()
	return bacaBarisCoverage(baris, TabelCoverage)
}
