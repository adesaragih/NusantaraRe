package repository

// % Indemnity coverage MBD (tiket 73 P2) `[terverifikasi]`: `NB FacIn\Activity\SetIndemnityRate_ACT.xml` langkah 2 RDB-List
// `NB FacIn\RDBList\SearchIndemnityRate_SQL.xml` L84 `select INDEMNITY_PCT from m_bi_indemnity where INDEMNITY =
// {SearchIndemnityRate.CARI1} and  replace(INDEMNITY_UNIT,' ','')={SearchIndemnityRate.CARI2}` - CARI1 = `.Indemnity`, CARI2 =
// `.Unit` (langkah 1). DDL `D:\migrasi\RNM\DDL\M_BI_INDEMNITY.txt` (ditaruh work owner 09-10-2026): INDEMNITY NUMBER(4,0),
// INDEMNITY_UNIT CHAR(3) (berpadding - dipangkas oleh REPLACE milik SQL Pega sendiri), INDEMNITY_PCT NUMBER(7,4). Isi tabel tidak
// ada di korpus. Baca saja; tanpa ORDER BY (RDB Pega juga tanpa urutan).

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cockroachdb/apd/v3"
)

// TabelBIIndemnity - tabel warisan POOLDATA (ejaan DDL).
const TabelBIIndemnity = "M_BI_INDEMNITY"

// sqlIndemnityPct - :1 Indemnity (bilangan bulat, sudah diperiksa services), :2 Unit.
func sqlIndemnityPct(tabel string) string {
	return "SELECT " + angkaKeluar("INDEMNITY_PCT") + " FROM " + tabel + " WHERE INDEMNITY = :1 AND REPLACE(INDEMNITY_UNIT, ' ', '') = :2"
}

// IndemnityPct - lihat PembacaLookupAneka: setiap baris cocok berurutan hasil kueri (NULL -> nil).
func (r *LookupAnekaOracle) IndemnityPct(ctx context.Context, indemnity, unit string) ([]*apd.Decimal, error) {
	q, err := r.db.Qualify(TabelBIIndemnity)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlIndemnityPct(q), indemnity, unit)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelBIIndemnity, err)
	}
	defer baris.Close()
	hasil := []*apd.Decimal{}
	for baris.Next() {
		var v sql.NullString
		if err := baris.Scan(&v); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelBIIndemnity, err)
		}
		d, err := bacaDesimal(TabelBIIndemnity, "INDEMNITY_PCT", &v)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, d)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelBIIndemnity, err)
	}
	return hasil, nil
}
