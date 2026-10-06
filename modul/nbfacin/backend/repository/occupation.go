package repository

// Saran Occupation sub-tab Surrounding Risk (tiket 38).
//
// `[terverifikasi]` `NB FacIn\Section\RiskAround.xml` sel 9 (dan sel Occupation sisi lain) =
// autocomplete atas data page `D_BrowseOccupationFacInFIRE` dengan parameter `TYPE = "FIRE"` (yang
// lain kosong), medan cari `.OldID` dan `.Name`, nilai tersimpan `.OldID`. RD
// `NB FacIn\ReportDefinition\BrowseOccupationFacInFIRE_RD.xml`: logika `A AND B AND D AND E AND
// (F OR C)` - filter berparameter KOSONG dibuang (pesan pxWarningMessage RD itu: "the filter will be
// dropped entirely"), sehingga yang tersisa hanya D `.Type = Param.TYPE`; tanpa DISTINCT,
// pyMaxRecords 500, urut Name ASC lalu OldID ASC. Kelas -> tabel: `RDBList\SearchOccupationIDSQL.xml`
// (kelas ASM-FW-GISFW-INT-OCCUPATION) membaca `FROM OCCUPATION`. Kolom DDL OCCUPATION: OLDID, NAME,
// TYPE (VARCHAR2(1000)).

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelOccupation - tabel warisan POOLDATA, baca saja.
	TabelOccupation = "OCCUPATION"
	// TipeOccupationFire - parameter TYPE data page sel 9.
	TipeOccupationFire = "FIRE"
	// BatasSaranOccupation - pyMaxRecords RD.
	BatasSaranOccupation = 500
)

// PembacaOccupation - saran Occupation FIRE.
type PembacaOccupation interface {
	CariOccupation(ctx context.Context, kata string) ([]models.BarisOccupation, error)
}

// OccupationOracle - PembacaOccupation atas Oracle.
type OccupationOracle struct{ db *db.DB }

// NewOccupationOracle merakit pembaca OCCUPATION.
func NewOccupationOracle(d *db.DB) *OccupationOracle { return &OccupationOracle{db: d} }

// awalSelectOccupation - awal pernyataan SELECT (kolom terpilih; (KDRISKEXPOSURE tiket 40: `SetDataOccupation` membaca `.KDRiskExposure`)).
const awalSelectOccupation = "SELECT OLDID, NAME, KDRISKEXPOSURE FROM "

// sqlCariOccupation - A131: OLDID ATAU NAME MENGANDUNG kata, tidak peka huruf (UPPER kedua sisi,
// pola PolaCari); TYPE dibandingkan persis. Urut RD (NAME, OLDID); paling banyak 500.
// :1 tipe, :2 dan :3 pola yang sama, :4 batas.
func sqlCariOccupation(occ string) string {
	return awalSelectOccupation + occ +
		` WHERE TYPE = :1 AND (UPPER(OLDID) LIKE :2 ESCAPE '\' OR UPPER(NAME) LIKE :3 ESCAPE '\')` +
		" ORDER BY NAME, OLDID FETCH FIRST :4 ROWS ONLY"
}

// sqlSemuaOccupation - tiket 40: kata kosong = seluruh okupasi FIRE (popup Choose Occupation memuat grid RD tanpa
// saringan kata - filter berparameter kosong dibuang). :1 tipe, :2 batas.
func sqlSemuaOccupation(occ string) string {
	return awalSelectOccupation + occ + " WHERE TYPE = :1 ORDER BY NAME, OLDID FETCH FIRST :2 ROWS ONLY"
}

// CariOccupation - lihat PembacaOccupation.
func (r *OccupationOracle) CariOccupation(ctx context.Context, kata string) ([]models.BarisOccupation, error) {
	q, err := r.db.Qualify(TabelOccupation)
	if err != nil {
		return nil, err
	}
	sqlq, arg := sqlSemuaOccupation(q), []any{TipeOccupationFire, BatasSaranOccupation}
	if pola := PolaCari(kata); pola != "" {
		sqlq, arg = sqlCariOccupation(q), []any{TipeOccupationFire, pola, pola, BatasSaranOccupation}
	}
	baris, err := r.db.QueryContext(ctx, sqlq, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelOccupation, err)
	}
	defer baris.Close()
	hasil := []models.BarisOccupation{}
	for baris.Next() {
		var oldID, nama, kd sql.NullString
		if err := baris.Scan(&oldID, &nama, &kd); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelOccupation, err)
		}
		hasil = append(hasil, models.BarisOccupation{OldID: oldID.String, Name: nama.String, KdRiskExposure: kd.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelOccupation, err)
	}
	return hasil, nil
}
