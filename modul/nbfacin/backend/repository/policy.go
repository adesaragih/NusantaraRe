package repository

// Jalur Policy (tiket 69a) - bacaan IsFlagUW: RDBList\GetCountAccumulation.xml L86-89 VERBATIM `[terverifikasi]`
// (`D:\migrasi\RNM\NB FacIn\`), dipanggil Activity\CekLimitSpreading_Act.xml 4.1.1.2 dengan TempAcuum.CARI1 = AccumulationCode,
// CARI2 = StartDateTime dd/MM/yyyy WIB, CARI3 = pyID (FACINOFFER), CARI4 = pzInsKey (FACINPRODUCTION). Tabel POOLDATA dibaca saja.

import (
	"context"
	"database/sql"
	"fmt"
)

// TabelFacInProduction - POOLDATA.FACINPRODUCTION (ditulis SaveJsonPolicyFacIn_Act Pega; dibaca saja di sini).
const TabelFacInProduction = "FACINPRODUCTION"

// sqlHitungAkumulasi - GetCountAccumulation; bind posisi :1-:6 urut kemunculan (go-ora).
func sqlHitungAkumulasi(prod, offer string) string {
	return "select SUM(CARID1)AS CARID1 FROM (\n" +
		"select count(*) AS CARID1 from " + prod + " where ACCUMULATIONCODE = :1 AND TO_DATE (ENDDATE, 'dd/mm/rrrr') >= " +
		"TO_DATE (:2, 'dd/mm/rrrr') and IDPEGA != :3\n" +
		"UNION ALL\n" +
		"select count(*) AS CARID1 from " + offer + " where ACCUMULATIONCODE = :4 AND TO_DATE (ENDDATE, 'dd/mm/rrrr') >= " +
		"TO_DATE (:5, 'dd/mm/rrrr')and IDPEGA != :6)"
}

// HitungAkumulasiBerjalan - lihat PenyimpanSubmit.
func (r *SubmitOracle) HitungAkumulasiBerjalan(ctx context.Context, kode, tanggal, pzInsKey, pyID string) (int, error) {
	prod, err := r.db.Qualify(TabelFacInProduction)
	if err != nil {
		return 0, err
	}
	offer, err := r.db.Qualify(TabelFacInOffer)
	if err != nil {
		return 0, err
	}
	var n sql.NullInt64
	if err := r.db.QueryRowContext(ctx, sqlHitungAkumulasi(prod, offer), kode, tanggal, pzInsKey, kode, tanggal, pyID).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: GetCountAccumulation %s: %w", kode, err)
	}
	return int(n.Int64), nil
}
