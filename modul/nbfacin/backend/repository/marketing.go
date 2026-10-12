package repository

// Pilihan Marketing Name (sel 75 section Periode, `.QuotationData.MOID`) - tiket 31.
// `[terverifikasi]` dropdown sel 75: pyListSource reportdefinition `BrowseMarketingOfficer_RD`,
// nilai `.ID`, teks `.ClientName` (pyPrompt), pilihan kosong "Choose". RD
// (`NB FacIn\ReportDefinition\BrowseMarketingOfficer_RD.xml`, kelas
// ASM-FW-GISFW-Int-marketingofficer): filter A `.MOStatus = "1"` (Text), urut `.ID` DESC,
// pyMaxRecords 500. Tabel: DDL `D:\migrasi\RNM\DDL\MARKETINGOFFICER.txt` (POOLDATA.MARKETINGOFFICER;
// `PEGA_MARKETINGOFFICER.txt` adalah PROSEDUR penulis tabel itu, bukan tabel). Kolom yang dibaca:
// ID VARCHAR2(100), CLIENTNAME VARCHAR2(100), MOSTATUS VARCHAR2(100).

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelMarketingOfficer - tabel warisan POOLDATA, baca saja.
	TabelMarketingOfficer = "MARKETINGOFFICER"
	// StatusMOAktif - nilai filter RD `.MOStatus = "1"` (teks).
	StatusMOAktif = "1"
	// BatasMarketingOfficer - pyMaxRecords RD.
	BatasMarketingOfficer = 500
)

// PembacaMarketingOfficer - pilihan Marketing Name.
type PembacaMarketingOfficer interface {
	DaftarMarketingOfficer(ctx context.Context) ([]models.MarketingOfficer, error)
}

// MarketingOfficerOracle - PembacaMarketingOfficer atas Oracle.
type MarketingOfficerOracle struct{ db *db.DB }

// NewMarketingOfficerOracle merakit pembaca tabel marketing officer.
func NewMarketingOfficerOracle(d *db.DB) *MarketingOfficerOracle {
	return &MarketingOfficerOracle{db: d}
}

func sqlMarketingOfficer(tabel string) string {
	return "SELECT ID, CLIENTNAME FROM " + tabel + " WHERE MOSTATUS = :1 ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY"
}

// DaftarMarketingOfficer - lihat PembacaMarketingOfficer.
func (r *MarketingOfficerOracle) DaftarMarketingOfficer(ctx context.Context) ([]models.MarketingOfficer, error) {
	q, err := r.db.Qualify(TabelMarketingOfficer)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlMarketingOfficer(q), StatusMOAktif, BatasMarketingOfficer)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelMarketingOfficer, err)
	}
	defer baris.Close()
	hasil := []models.MarketingOfficer{}
	for baris.Next() {
		var id, nama sql.NullString
		if err := baris.Scan(&id, &nama); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelMarketingOfficer, err)
		}
		hasil = append(hasil, models.MarketingOfficer{ID: id.String, Nama: nama.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelMarketingOfficer, err)
	}
	return hasil, nil
}

// PembacaMarketingBawaan - saran awal Marketing Name layar InwardFacultative (keputusan work owner 12-10-2026: "buatin marketing name
// otomatis terpilih dari user yang buat, kalau usernya tidak ada di marketingofficer baru dikosongin"): ID MARKETINGOFFICER yang
// AKSES_LOGIN-nya = pembuat case (T_WORK_POLIS.CREATE_OP), urut ID turun. Opsional (type assertion) - tiruan lama tidak terdampak.
type PembacaMarketingBawaan interface {
	MarketingBawaan(ctx context.Context, caseID string) ([]string, error)
}

// sqlMarketingBawaan - :1 ID case, :2 LINI, :3 MOSTATUS aktif (saringan dropdown sqlMarketingOfficer); login dicocokkan
// UPPER(TRIM) seperti sqlTeamMarketing; dua baris cukup untuk mengenali ganda (ganda -> saran kosong, services).
func sqlMarketingBawaan(mo, work string) string {
	return "SELECT m.ID FROM " + mo + " m JOIN " + work + " w ON UPPER(TRIM(m.AKSES_LOGIN)) = UPPER(TRIM(w.CREATE_OP))" +
		" WHERE w.ID = :1 AND w.LINI = :2 AND m.MOSTATUS = :3 AND m.ID IS NOT NULL ORDER BY m.ID DESC FETCH FIRST 2 ROWS ONLY"
}

// MarketingBawaan - lihat PembacaMarketingBawaan.
func (r *MarketingOfficerOracle) MarketingBawaan(ctx context.Context, caseID string) ([]string, error) {
	mo, err := r.db.Qualify(TabelMarketingOfficer)
	if err != nil {
		return nil, err
	}
	work, err := r.db.Qualify(TabelWorkPolis)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlMarketingBawaan(mo, work), caseID, LiniFacIn, StatusMOAktif)
	if err != nil {
		return nil, fmt.Errorf("repository: marketing bawaan %s: %w", TabelMarketingOfficer, err)
	}
	defer baris.Close()
	var hasil []string
	for baris.Next() {
		var id sql.NullString
		if err := baris.Scan(&id); err != nil {
			return nil, fmt.Errorf("repository: marketing bawaan: %w", err)
		}
		hasil = append(hasil, id.String)
	}
	return hasil, baris.Err()
}
