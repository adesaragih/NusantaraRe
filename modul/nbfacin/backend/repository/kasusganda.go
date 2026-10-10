package repository

// Kasus ganda Submit Aneka (tiket 73 P4) - RDBList\GetDataDoubleCase.xml L79-85 `[terverifikasi]` (`D:\migrasi\RNM\NB FacIn\`,
// `ASM-FW-GISFW-INT-POLICYJSON!ASM!GETDATADOUBLECASE`), dipanggil ProtectCoverage_Act 5.1.12 dengan ParamInProdCase.CARI10-15
// (L588-880). PERSIS kecuali placeholder ber-bind dan saringan `IDPEGA` case sendiri (K73-31 keputusan work owner 09-10-2026:
// case sendiri dikecualikan - menyimpang dari Pega). Kolom: DDL\FACINPRODUCTION.txt (BEGINDATE / ENDDATE DATE). Dibaca saja.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/modul/nbfacin/backend/models"
)

// sqlKasusGanda - :1 InsuredName tanpa "." / spasi, :2 SobName, :3 BusinessCode, :4 NoOfferSlip, :5 / :6 Begin / End
// dd/MM/yyyy, :7 IDPEGA case sendiri. Urutan baris tidak ditentukan (Pega memakai pxResults(1) tanpa ORDER BY).
func sqlKasusGanda(prod string) string {
	return "select distinct IDPEGA,NOPOLIS from " + prod + " where NOENDORS is null AND replace(replace(INSUREDNAME,'.',''),' ','') = :1" +
		" AND SOB = :2 AND BUSINESSCODE = :3 AND RISLIPCEDING = :4 AND to_char(BEGINDATE,'dd/MM/yyyy') = :5" +
		" AND to_char(ENDDATE,'dd/MM/yyyy') = :6 AND (IDPEGA IS NULL OR IDPEGA <> :7)"
}

// KasusGanda - lihat PenyimpanSubmit.
func (r *SubmitOracle) KasusGanda(ctx context.Context, c models.CariKasusGanda) ([]string, error) {
	prod, err := r.db.Qualify(TabelFacInProduction)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlKasusGanda(prod), c.InsuredName, c.SobName, c.BusinessCode, c.NoOfferSlip, c.Mulai,
		c.Akhir, c.IDPegaSendiri)
	if err != nil {
		return nil, fmt.Errorf("repository: GetDataDoubleCase: %w", err)
	}
	defer baris.Close()
	hasil := []string{}
	for baris.Next() {
		var id, polis sql.NullString
		if err := baris.Scan(&id, &polis); err != nil {
			return nil, fmt.Errorf("repository: GetDataDoubleCase: %w", err)
		}
		hasil = append(hasil, id.String)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: GetDataDoubleCase: %w", err)
	}
	return hasil, nil
}
