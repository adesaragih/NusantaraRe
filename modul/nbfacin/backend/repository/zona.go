package repository

// Zone dan 4.2 Construction form coverage FIRE (tiket 52) - saran autocomplete dan rate gempa menurut zona.
//
// Zone `[terverifikasi]`: autocomplete `Section\CoverageItem.xml` RD `BrowseZones_RD` (kelas ASM-FW-GISFW-Int-ZONES,
// nilai .Description, param ID / ACTIVEDATE dikirim kosong, COVERAGE = .OLDID). RD (`NB FacIn\ReportDefinition\
// BrowseZones_RD.xml`; salinan RNW sama, Endorsment beda isi dan tidak dipakai): pyFilterLogic "A AND B AND (C OR D OR E
// OR F OR G OR H) AND I" - A .ID = Param.ID, B .ActiveDate = Param.ACTIVEDATE, C..H .Zone = Zone.CARI1..6, I .OLDID =
// Param.COVERAGE; pyUseNullIfEmpty false (saringan kosong dibuang), pyCaseInsensitive true, tanpa urutan, tanpa
// DISTINCT, pyMaxRecords 500. Kelas -> view POOLDATA.ZONES (DDL `DDL\ZONES.txt`: ID, OLDID, ZONE, DESCRIPTION, RATE,
// DESCRIPTION2, ACTIVEDATE, TYPE atas JSON m_zones, kolom tanpa tipe tertulis - dibaca sebagai teks).
//
// 4.2 Construction `[terverifikasi]`: `RDBList\SearchClassConstructionReas_SQL.xml` = `select distinct construction from
// construction order by CONSTRUCTION` (view POOLDATA.CONSTRUCTION, DDL `DDL\CONSTRUCTION.txt`).
//
// Rate gempa `[terverifikasi]`: `Activity\CountPremiRate_ACT.xml` cabang Okupasi.CARI2 == "4.2 PRGBI" - Earthquake.CARI1
// = .Zone, RDB-List `RateEarthquake_Sql` (`SELECT RATE AS CARI1 FROM POOLDATA.ZONEOJK WHERE DESCRIPTION={Earthquake.CARI1}`,
// tanpa urutan), lalu `.Rate = InputEarthquake.pxResults(1).CARI1`. ZONEOJK = tabel (DDL `DDL\ZONEOJK.txt`: ID / ZONE
// VARCHAR2(20), DESCRIPTION VARCHAR2(1000), RATE NUMBER(20,3)).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelZones, TabelConstruction, TabelZoneOJK - objek warisan POOLDATA, baca saja.
	TabelZones        = "ZONES"
	TabelConstruction = "CONSTRUCTION"
	TabelZoneOJK      = "ZONEOJK"
)

// PembacaZona - saran Zone / 4.2 Construction dan rate ZONEOJK.
type PembacaZona interface {
	// SaranZona - BrowseZones_RD: OLDID = oldID (kosong -> dibuang), ZONE salah satu `zona` (kosong -> dibuang),
	// DESCRIPTION mengandung kata; tidak peka huruf.
	SaranZona(ctx context.Context, oldID string, zona []string, kata string) ([]models.SaranAkumulasi, error)
	// SaranKonstruksiGempa - SearchClassConstructionReas_SQL, CONSTRUCTION mengandung kata (tidak peka huruf).
	SaranKonstruksiGempa(ctx context.Context, kata string) ([]models.SaranAkumulasi, error)
	// RateZona - RateEarthquake_Sql: RATE baris pertama ZONEOJK.DESCRIPTION = zona; nil = tidak ada baris / RATE NULL.
	RateZona(ctx context.Context, zona string) (*apd.Decimal, error)
}

// ZonaOracle - PembacaZona atas Oracle.
type ZonaOracle struct{ db *db.DB }

// NewZonaOracle merakit pembaca ZONES / CONSTRUCTION / ZONEOJK.
func NewZonaOracle(d *db.DB) *ZonaOracle { return &ZonaOracle{db: d} }

// sqlSaranZona - syarat dan bind dibangun BERSAMA (pola sqlSaran). UPPER di kedua sisi = pyCaseInsensitive. Urut
// DESCRIPTION lalu ID (RD tanpa urutan; pola A174), batas saran 50 (pola A124; RD 500).
func sqlSaranZona(t, oldID string, zona []string, pola string) (string, []any) {
	var syarat []string
	var arg []any
	bind := func(v any) string { arg = append(arg, v); return ":" + strconv.Itoa(len(arg)) }
	if len(zona) > 0 {
		b := make([]string, len(zona))
		for i, z := range zona {
			b[i] = "UPPER(" + bind(z) + ")"
		}
		syarat = append(syarat, "UPPER(ZONE) IN ("+strings.Join(b, ", ")+")")
	}
	if oldID != "" {
		syarat = append(syarat, "UPPER(OLDID) = UPPER("+bind(oldID)+")")
	}
	if pola != "" {
		syarat = append(syarat, "UPPER(DESCRIPTION) LIKE "+bind(pola)+` ESCAPE '\'`)
	}
	where := ""
	if len(syarat) > 0 {
		where = " WHERE " + strings.Join(syarat, " AND ")
	}
	return "SELECT ID, DESCRIPTION FROM " + t + where + " ORDER BY DESCRIPTION, ID FETCH FIRST " + bind(BatasSaranAkumulasi) +
		" ROWS ONLY", arg
}

// sqlSaranKonstruksi - DISTINCT CONSTRUCTION urut CONSTRUCTION (SQL asal), kata opsional, batas saran 50; kolom dipilih
// dua kali (id = label) dengan ALIAS berbeda: tanpa alias, `ORDER BY CONSTRUCTION` ambigu -> ORA-00960 (ditemukan work
// owner di DEV 06-10-2026, diperbaiki sesi 0f).
func sqlSaranKonstruksi(t, pola string) (string, []any) {
	if pola == "" {
		return "SELECT CONSTRUCTION AS KODE, CONSTRUCTION AS NAMA FROM (SELECT DISTINCT CONSTRUCTION FROM " + t +
			" WHERE CONSTRUCTION IS NOT NULL) ORDER BY KODE FETCH FIRST :1 ROWS ONLY", []any{BatasSaranAkumulasi}
	}
	return "SELECT CONSTRUCTION AS KODE, CONSTRUCTION AS NAMA FROM (SELECT DISTINCT CONSTRUCTION FROM " + t +
		` WHERE UPPER(CONSTRUCTION) LIKE :1 ESCAPE '\') ORDER BY KODE FETCH FIRST :2 ROWS ONLY`, []any{pola, BatasSaranAkumulasi}
}

// sqlRateZona - RateEarthquake_Sql; pxResults(1) dari kueri tanpa urutan -> baris pertama urut ID `[dugaan]` (urutan
// Pega = urutan basis data, tidak tertulis).
func sqlRateZona(t string) string {
	return "SELECT " + angkaKeluar("RATE") + " FROM " + t + " WHERE DESCRIPTION = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

// SaranZona - lihat PembacaZona.
func (r *ZonaOracle) SaranZona(ctx context.Context, oldID string, zona []string, kata string) ([]models.SaranAkumulasi, error) {
	t, err := r.db.Qualify(TabelZones)
	if err != nil {
		return nil, err
	}
	q, arg := sqlSaranZona(t, oldID, zona, PolaCari(kata))
	return r.bacaSaran(ctx, TabelZones, q, arg)
}

// SaranKonstruksiGempa - lihat PembacaZona.
func (r *ZonaOracle) SaranKonstruksiGempa(ctx context.Context, kata string) ([]models.SaranAkumulasi, error) {
	t, err := r.db.Qualify(TabelConstruction)
	if err != nil {
		return nil, err
	}
	q, arg := sqlSaranKonstruksi(t, PolaCari(kata))
	return r.bacaSaran(ctx, TabelConstruction, q, arg)
}

// RateZona - lihat PembacaZona.
func (r *ZonaOracle) RateZona(ctx context.Context, zona string) (*apd.Decimal, error) {
	t, err := r.db.Qualify(TabelZoneOJK)
	if err != nil {
		return nil, err
	}
	var rate sql.NullString
	err = r.db.QueryRowContext(ctx, sqlRateZona(t), zona).Scan(&rate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelZoneOJK, err)
	}
	return bacaDesimal(TabelZoneOJK, "RATE", &rate)
}

// bacaSaran - baris (id, label).
func (r *ZonaOracle) bacaSaran(ctx context.Context, tabel, q string, arg []any) ([]models.SaranAkumulasi, error) {
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", tabel, err)
	}
	defer baris.Close()
	hasil := []models.SaranAkumulasi{}
	for baris.Next() {
		var id, label sql.NullString
		if err := baris.Scan(&id, &label); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", tabel, err)
		}
		hasil = append(hasil, models.SaranAkumulasi{ID: id.String, Label: label.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", tabel, err)
	}
	return hasil, nil
}
