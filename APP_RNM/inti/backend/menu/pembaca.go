package menu

// Pembaca `M_NAV_MENU` - satu SELECT, baris aktif saja.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
)

// PembacaMenu membaca baris aktif `M_NAV_MENU`. Antarmuka supaya rute dapat
// diuji tanpa Oracle.
type PembacaMenu interface {
	Baca(ctx context.Context) ([]Baris, error)
}

// Pembaca membaca `M_NAV_MENU` dari Oracle.
type Pembaca struct {
	db *db.DB
}

// NewPembaca membuat pembaca menu.
func NewPembaca(d *db.DB) *Pembaca { return &Pembaca{db: d} }

// Nilai bendera satu-karakter `M_NAV_MENU` (STATUS_AKTIF, DIMIGRASI): '1' ya,
// '0' tidak - konvensi yang sama dengan data warisan.
const benderaYa = "1"

// sqlMenu merakit pernyataannya - baris STATUS_AKTIF '1' (bind `:1` =
// benderaYa), urut GROUPMENU, URUTAN (brief menu §3); ID sebagai penentu akhir
// supaya urutannya tetap.
func sqlMenu(tabel string) string {
	return fmt.Sprintf(`SELECT ID, PARENT_ID, KODE, LABEL, GROUPMENU, MODUL, URUTAN, DIMIGRASI FROM %s `+
		`WHERE STATUS_AKTIF = :1 ORDER BY GROUPMENU, URUTAN, ID`, tabel)
}

// Baca mengembalikan seluruh baris aktif, kelompok dan butir.
func (p *Pembaca) Baca(ctx context.Context) ([]Baris, error) {
	tabel, err := p.db.Qualify("M_NAV_MENU")
	if err != nil {
		return nil, err
	}
	q := sqlMenu(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := p.db.QueryContext(ctx, q, benderaYa)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca M_NAV_MENU: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []Baris
	for rows.Next() {
		var b Baris
		var induk sql.NullInt64
		var dimigrasi string
		if err := rows.Scan(&b.ID, &induk, &b.Kode, &b.Label, &b.Golongan, &b.Modul, &b.Urutan, &dimigrasi); err != nil {
			return nil, fmt.Errorf("repository: membaca baris M_NAV_MENU: %w", err)
		}
		b.IndukID = induk.Int64
		b.Dimigrasi = dimigrasi == benderaYa
		out = append(out, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca M_NAV_MENU: %w", err)
	}
	return out, nil
}
