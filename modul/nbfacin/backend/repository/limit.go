// Package repository memuat pembaca Oracle modul NB Fac In.
package repository

// Tabel limit akseptasi (tiket 20). Skema HANYA dari DDL
// `D:\migrasi\RNM\DDL\M_LIMIT_*.txt` `[terverifikasi]`: kolom yang dibaca adalah yang
// dibaca rule SQL Pega di `NB FacIn\RDBList\` (tiket 11, 12) - bentuk A `JABATAN`
// VARCHAR2(100), `TEAM_GROUP` VARCHAR2(10/20), `LIMIT_BOTTOM` dan `LIMIT_BOTTOM2`
// NUMBER(*,0); bentuk B `JABATAN`, `LIMITBOND_BOTTOM`, `LIMITCREDITCL_BOTTOM`,
// `LIMITCREDITNCL_BOTTOM` NUMBER(*,0) - bilangan bulat. Dibaca sebagai teks, diurai
// desimal (CLAUDE.md §7: tanpa float).
//
// ⛔ `NAMA` dan `LOGIN` TIDAK dibaca (CLAUDE.md §4 butir 10, K-025), dan tidak ada
// `SELECT *`. Hanya membaca; nol tulisan.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TabelBentukA - lima tabel bentuk A, ejaan DDL apa adanya (huruf ganda di akhir
// adalah nama tabel sungguhan).
var TabelBentukA = []string{
	"M_LIMIT_PROPERTYY", "M_LIMIT_PROPERTY_NON_PREFERREDD", "M_LIMIT_PROPERTY_PREFERRED_COMMERCIALL",
	"M_LIMIT_ENGINEERINGG", "M_LIMIT_NONPROPANDENGG",
}

// TabelBentukB - tabel limit bentuk B (Bond/Kredit/Trade).
const TabelBentukB = "M_LIMIT_FINANCIALINS"

// PembacaLimit - sumber tabel limit akseptasi.
type PembacaLimit interface {
	MuatLimit(ctx context.Context) ([]models.BarisLimitA, []models.BarisLimitB, error)
}

// LimitOracle - PembacaLimit atas Oracle.
type LimitOracle struct{ db *db.DB }

// NewLimitOracle merakit pembaca tabel limit.
func NewLimitOracle(d *db.DB) *LimitOracle { return &LimitOracle{db: d} }

func sqlBentukA(tabel string) string {
	return "SELECT JABATAN, TEAM_GROUP, LIMIT_BOTTOM, LIMIT_BOTTOM2 FROM " + tabel
}

func sqlBentukB(tabel string) string {
	return "SELECT JABATAN, LIMITBOND_BOTTOM, LIMITCREDITCL_BOTTOM, LIMITCREDITNCL_BOTTOM FROM " + tabel
}

// MuatLimit - seluruh baris keenam tabel. Urutan baris = urutan Oracle; tangga
// mengurutkan sendiri dan menolak seri yang tidak pasti (tiket 11, butir 44).
func (r *LimitOracle) MuatLimit(ctx context.Context) ([]models.BarisLimitA, []models.BarisLimitB, error) {
	var a []models.BarisLimitA
	for _, nama := range TabelBentukA {
		q, err := r.db.Qualify(nama)
		if err != nil {
			return nil, nil, err
		}
		baris, err := r.db.QueryContext(ctx, sqlBentukA(q))
		if err != nil {
			return nil, nil, fmt.Errorf("repository: baca %s: %w", nama, err)
		}
		for baris.Next() {
			var jab, tg, lb, lb2 sql.NullString
			if err := baris.Scan(&jab, &tg, &lb, &lb2); err != nil {
				baris.Close()
				return nil, nil, fmt.Errorf("repository: %s: %w", nama, err)
			}
			b, err := uraiBarisA(nama, jab, tg, lb, lb2)
			if err != nil {
				baris.Close()
				return nil, nil, err
			}
			a = append(a, b)
		}
		if err := baris.Err(); err != nil {
			baris.Close()
			return nil, nil, fmt.Errorf("repository: %s: %w", nama, err)
		}
		baris.Close()
	}
	q, err := r.db.Qualify(TabelBentukB)
	if err != nil {
		return nil, nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlBentukB(q))
	if err != nil {
		return nil, nil, fmt.Errorf("repository: baca %s: %w", TabelBentukB, err)
	}
	defer baris.Close()
	var b []models.BarisLimitB
	for baris.Next() {
		var jab, bond, cl, ncl sql.NullString
		if err := baris.Scan(&jab, &bond, &cl, &ncl); err != nil {
			return nil, nil, fmt.Errorf("repository: %s: %w", TabelBentukB, err)
		}
		f, err := uraiBarisB(jab, bond, cl, ncl)
		if err != nil {
			return nil, nil, err
		}
		b = append(b, f)
	}
	return a, b, baris.Err()
}

func uraiBarisA(tabel string, jab, tg, lb, lb2 sql.NullString) (models.BarisLimitA, error) {
	b := models.BarisLimitA{Tabel: tabel, Jabatan: jab.String, TeamGroup: tg.String}
	var err error
	if b.LimitBottom, err = db.UraiDesimal(tabel+" "+jab.String, "LIMIT_BOTTOM", lb); err != nil {
		return models.BarisLimitA{}, err
	}
	if b.LimitBottom2, err = db.UraiDesimal(tabel+" "+jab.String, "LIMIT_BOTTOM2", lb2); err != nil {
		return models.BarisLimitA{}, err
	}
	return b, nil
}

func uraiBarisB(jab, bond, cl, ncl sql.NullString) (models.BarisLimitB, error) {
	b := models.BarisLimitB{Jabatan: jab.String}
	id := TabelBentukB + " " + jab.String
	var err error
	if b.LimitBond, err = db.UraiDesimal(id, "LIMITBOND_BOTTOM", bond); err != nil {
		return models.BarisLimitB{}, err
	}
	if b.LimitCreditCL, err = db.UraiDesimal(id, "LIMITCREDITCL_BOTTOM", cl); err != nil {
		return models.BarisLimitB{}, err
	}
	if b.LimitCreditNCL, err = db.UraiDesimal(id, "LIMITCREDITNCL_BOTTOM", ncl); err != nil {
		return models.BarisLimitB{}, err
	}
	return b, nil
}
