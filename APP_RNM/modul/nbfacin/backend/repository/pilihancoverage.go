package repository

// Pilihan coverage tab Coverage FIRE (tiket 43).
//
// Popup Choose Coverage `[terverifikasi]`: RD `BrowseCoverageFacIn_RD` (kelas ASM-FW-GISFW-Int-COVERAGE_FACIN; tabel dari
// `RDBList\GetCoverageFacin.xml` "FROM COVERAGE_FACIN"): A `.Type = Param.Type` ("FIRE" dari `Section\ChooseCoverage.xml`),
// B `.BizCode = Param.BizCode` (dikirim KOSONG -> dibuang), C `.NamaCoverage` Contains `Param.Nama` (tidak peka huruf),
// D/E `.ACTIVESTATUS = 1` OR IS NULL; DISTINCT atas kolom laporan (ID, BizCode, NamaCoverage, Type, OLDID, ACTIVESTATUS);
// maks 500; urut NamaCoverage lalu OLDID. DDL `DDL\COVERAGE_FACIN.txt` (03-10-2026): VIEW (ID, BIZCODE, NAMACOVERAGE,
// TYPE, OLDID) atas JSON M_COVERAGE - TANPA ACTIVESTATUS, jadi saringan D/E tidak dapat dinyatakan (A160).
//
// Coverage otomatis `[terverifikasi]`: `Activity\AddCoverageAutoFire.xml` langkah 3 mengisi .Coverage berurutan 100815,
// 100828, 100829, 100825, 100840, lalu Obj-Browse kelas ASM-FW-GISFW-Int-COVERAGE (.ID = .Coverage) -> .CoverageNote =
// Name, .OLDID = OLDID. Kelas -> tabel COVERAGE (`RDBList\SearchCoverageIDSQL.xml` "FROM COVERAGE"); DDL
// `DDL\COVERAGE.txt`: ID / NAME VARCHAR2(4000), OLDID VARCHAR2(23).

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelCoverageFacIn - view warisan POOLDATA (baca saja); TabelCoverage - tabel warisan POOLDATA (baca saja).
	TabelCoverageFacIn = "COVERAGE_FACIN"
	TabelCoverage      = "COVERAGE"
	// TipeCoverageFire - Param.Type popup Choose Coverage. BatasCoverage - pyMaxRecords RD.
	TipeCoverageFire = "FIRE"
	BatasCoverage    = 500
)

// kodeCoverageOtomatis - AddCoverageAutoFire langkah 3, urutan korpus.
var kodeCoverageOtomatis = []string{"100815", "100828", "100829", "100825", "100840"}

// PembacaCoverage - pilihan coverage.
type PembacaCoverage interface {
	// CariCoverage - COVERAGE_FACIN FIRE; kata kosong = semua.
	CariCoverage(ctx context.Context, kata string) ([]models.BarisCoverage, error)
	// CoverageOtomatis - lima coverage AddCoverageAutoFire dari COVERAGE (urut korpus; kode tanpa baris -> nama kosong).
	CoverageOtomatis(ctx context.Context) ([]models.BarisCoverage, error)
}

// CoverageOracle - PembacaCoverage atas Oracle.
type CoverageOracle struct{ db *db.DB }

// NewCoverageOracle merakit pembaca COVERAGE_FACIN / COVERAGE.
func NewCoverageOracle(d *db.DB) *CoverageOracle { return &CoverageOracle{db: d} }

// sqlCariCoverage - :1 tipe, [:2 pola NAMACOVERAGE], batas terakhir. ID view lewat TO_CHAR (tipe M_COVERAGE.ID tidak
// ada di DDL `belum terverifikasi`).
func sqlCariCoverage(v string, denganKata bool) string {
	syarat, batas := "TYPE = :1", ":2"
	if denganKata {
		syarat, batas = `TYPE = :1 AND UPPER(NAMACOVERAGE) LIKE :2 ESCAPE '\'`, ":3"
	}
	return "SELECT ID, OLDID, NAMACOVERAGE FROM (SELECT DISTINCT TO_CHAR(ID) AS ID, BIZCODE, NAMACOVERAGE, TYPE, OLDID FROM " + v +
		" WHERE " + syarat + ") ORDER BY NAMACOVERAGE, OLDID, ID FETCH FIRST " + batas + " ROWS ONLY"
}

// sqlCoverageOtomatis - baris COVERAGE untuk lima kode (bind :1..:5), kolom urutan sama dengan sqlCariCoverage (ID,
// OLDID, nama); urut ID lalu NAME, OLDID (baris pertama per ID dipakai - Pega pxResults(1)).
func sqlCoverageOtomatis(c string) string {
	b := make([]string, len(kodeCoverageOtomatis))
	for i := range b {
		b[i] = ":" + strconv.Itoa(i+1)
	}
	return "SELECT ID, OLDID, NAME FROM " + c + " WHERE ID IN (" + strings.Join(b, ", ") + ") ORDER BY ID, NAME, OLDID"
}

func bacaBarisCoverage(baris *sql.Rows, tabel string) ([]models.BarisCoverage, error) {
	hasil := []models.BarisCoverage{}
	for baris.Next() {
		var id, a, b sql.NullString
		if err := baris.Scan(&id, &a, &b); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", tabel, err)
		}
		hasil = append(hasil, models.BarisCoverage{ID: id.String, OldID: a.String, Nama: b.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", tabel, err)
	}
	return hasil, nil
}

// CariCoverage - lihat PembacaCoverage.
func (r *CoverageOracle) CariCoverage(ctx context.Context, kata string) ([]models.BarisCoverage, error) {
	q, err := r.db.Qualify(TabelCoverageFacIn)
	if err != nil {
		return nil, err
	}
	sqlq, arg := sqlCariCoverage(q, false), []any{TipeCoverageFire, BatasCoverage}
	if pola := PolaCari(kata); pola != "" {
		sqlq, arg = sqlCariCoverage(q, true), []any{TipeCoverageFire, pola, BatasCoverage}
	}
	baris, err := r.db.QueryContext(ctx, sqlq, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCoverageFacIn, err)
	}
	defer baris.Close()
	return bacaBarisCoverage(baris, TabelCoverageFacIn)
}

// CoverageOtomatis - lihat PembacaCoverage.
func (r *CoverageOracle) CoverageOtomatis(ctx context.Context) ([]models.BarisCoverage, error) {
	q, err := r.db.Qualify(TabelCoverage)
	if err != nil {
		return nil, err
	}
	arg := make([]any, len(kodeCoverageOtomatis))
	for i, k := range kodeCoverageOtomatis {
		arg[i] = k
	}
	baris, err := r.db.QueryContext(ctx, sqlCoverageOtomatis(q), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCoverage, err)
	}
	defer baris.Close()
	mentah, err := bacaBarisCoverage(baris, TabelCoverage)
	if err != nil {
		return nil, err
	}
	pertama := map[string]models.BarisCoverage{}
	for _, m := range mentah {
		if _, ada := pertama[m.ID]; !ada {
			pertama[m.ID] = m
		}
	}
	hasil := make([]models.BarisCoverage, 0, len(kodeCoverageOtomatis))
	for _, k := range kodeCoverageOtomatis {
		b, ada := pertama[k]
		if !ada {
			b = models.BarisCoverage{ID: k}
		}
		hasil = append(hasil, b)
	}
	return hasil, nil
}
