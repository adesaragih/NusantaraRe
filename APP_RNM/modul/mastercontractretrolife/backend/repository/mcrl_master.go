package repository

// Master rujukan layar - DIBACA SAJA (paket 1).
//
//	REINS TYPE      `BrowseReinsuranceTypeLimit_RD` b578 `.Flag = 1`, urut `.ID DESC` b765; dropdown nilai `.ID`, tampil `.Note`
//	REINSURER NAME  `BrowseCedingCoLife_RD` b556 `B AND A AND C`: `.ID Contains "L0"` b565,
//	                `.ClientName Contains Param.CedingCoLeader` b584 (kata yang diketik), `.StatusActive = 1` b601;
//	                urut `.ClientName ASC` b745; isi `REINSURERID ← .ID`
//	BUSINESS NAME   `BrowseBusinessLife_RD` b651 `.OLDID StartsWith "L"`, urut `.ID ASC` b1057; tampil `.Note`,
//	                isi `BIZCODE ← .ID`, kolom tampil `.OLDID`
//	R/I RATE, Rate List  TIDAK dibaca - menunggu persetujuan (OQ-MCRL-13, lihat mcrl_tabel.go)
//
// ⛔ Kolom fisik AGENT (`ID`, `CLIENTNAME`, `STATUSACTIVE`) dan BUSINESS
// (`ID`, `NOTE`) sama dengan yang dibaca Treaty Contract Out; `OLDID` dari
// nama properti RD. Pencarian autocomplete tidak peka huruf besar-kecil.
// ⛔ Nol INSERT/UPDATE/DELETE atas master (penjaga modul).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// BatasPilihan membatasi jumlah pilihan autocomplete per ketikan.
const BatasPilihan = 100

// ErrMasterTidakTerbaca - objek master tidak dapat dibaca (mis. tidak ada di
// skema); pesannya menyebut objeknya, tidak diam.
var ErrMasterTidakTerbaca = errors.New("repository: reference master cannot be read")

// Teks SQL master - fungsi murni atas nama tabel berkualifikasi (diuji tanpa Oracle).

func sqlJenisReasuransi(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE FLAG = :1 ORDER BY ID DESC`, t)
}

func sqlCariReinsurer(t string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s
	 WHERE ID LIKE :1 AND UPPER(CLIENTNAME) LIKE :2 ESCAPE '\' AND STATUSACTIVE = :3
	 ORDER BY CLIENTNAME ASC, ID ASC FETCH FIRST %d ROWS ONLY`, t, BatasPilihan)
}

func sqlAmbilReinsurer(t string) string {
	return fmt.Sprintf(`SELECT ID, CLIENTNAME FROM %s WHERE ID = :1`, t)
}

func sqlCariBusiness(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE, OLDID FROM %s
	 WHERE OLDID LIKE :1 AND UPPER(NOTE) LIKE :2 ESCAPE '\'
	 ORDER BY ID ASC FETCH FIRST %d ROWS ONLY`, t, BatasPilihan)
}

func sqlAmbilBusiness(t string) string {
	return fmt.Sprintf(`SELECT ID, NOTE, OLDID FROM %s WHERE ID = :1`, t)
}

// Argumen saringan tetap - nilai VERBATIM RD (flag dan status di models).
const (
	idReinsurerLife = "%L0%" // `BrowseCedingCoLife_RD` b565 Contains "L0"
	awalanBizLife   = "L%"   // `BrowseBusinessLife_RD` b651 StartsWith "L"
)

// PolaCari merakit pola LIKE ber-escape untuk "Contains" RD.
func PolaCari(kata string) string {
	k := strings.ToUpper(strings.TrimSpace(kata))
	k = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(k)
	return "%" + k + "%"
}

func (g *Gudang) bacaMaster(ctx context.Context, objek string, susun func(string) string, kolom []string,
	args ...any) ([]barisTeks, error) {

	q, err := g.siapkan(objek, susun)
	if err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrMasterTidakTerbaca, objek, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []barisTeks
	for rows.Next() {
		b, err := pindai(rows, kolom)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %v", ErrMasterTidakTerbaca, objek, err)
		}
		hasil = append(hasil, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrMasterTidakTerbaca, objek, err)
	}
	return hasil, nil
}

// JenisReasuransiLife - lima jenis life (`FLAG = '1'`, "1 for life").
func (g *Gudang) JenisReasuransiLife(ctx context.Context) ([]models.JenisReasuransi, error) {
	bb, err := g.bacaMaster(ctx, MasterJenisReasuransi, sqlJenisReasuransi, []string{"ID", "NOTE"},
		models.FlagJenisReasuransiLife)
	if err != nil {
		return nil, err
	}
	hasil := make([]models.JenisReasuransi, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.JenisReasuransi{ID: b.s("ID"), Note: b.s("NOTE")})
	}
	return hasil, nil
}

// CariMasterReinsurer - pilihan `REINSURER NAME` / `SECURITY REINSURER NAME`.
func (g *Gudang) CariMasterReinsurer(ctx context.Context, kata string) ([]models.MasterReinsurer, error) {
	bb, err := g.bacaMaster(ctx, MasterReinsurer, sqlCariReinsurer, []string{"ID", "CLIENTNAME"},
		idReinsurerLife, PolaCari(kata), models.StatusMasterReinsurerAktif)
	if err != nil {
		return nil, err
	}
	hasil := make([]models.MasterReinsurer, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.MasterReinsurer{ID: b.s("ID"), ClientName: b.s("CLIENTNAME")})
	}
	return hasil, nil
}

// AmbilMasterReinsurer - satu reinsurer master menurut ID (tanpa saringan
// aktif: baris lama yang reinsurernya kini nonaktif tetap dapat disunting).
func (g *Gudang) AmbilMasterReinsurer(ctx context.Context, id string) (models.MasterReinsurer, bool, error) {
	bb, err := g.bacaMaster(ctx, MasterReinsurer, sqlAmbilReinsurer, []string{"ID", "CLIENTNAME"}, id)
	if err != nil || len(bb) == 0 {
		return models.MasterReinsurer{}, false, err
	}
	return models.MasterReinsurer{ID: bb[0].s("ID"), ClientName: bb[0].s("CLIENTNAME")}, true, nil
}

// CariMasterBusiness - pilihan `BUSINESS NAME`.
func (g *Gudang) CariMasterBusiness(ctx context.Context, kata string) ([]models.MasterBusiness, error) {
	bb, err := g.bacaMaster(ctx, MasterBusiness, sqlCariBusiness, []string{"ID", "NOTE", "OLDID"},
		awalanBizLife, PolaCari(kata))
	if err != nil {
		return nil, err
	}
	hasil := make([]models.MasterBusiness, 0, len(bb))
	for _, b := range bb {
		hasil = append(hasil, models.MasterBusiness{ID: b.s("ID"), Note: b.s("NOTE"), OldID: b.s("OLDID")})
	}
	return hasil, nil
}

// AmbilMasterBusiness - satu business master menurut ID.
func (g *Gudang) AmbilMasterBusiness(ctx context.Context, id string) (models.MasterBusiness, bool, error) {
	bb, err := g.bacaMaster(ctx, MasterBusiness, sqlAmbilBusiness, []string{"ID", "NOTE", "OLDID"}, id)
	if err != nil || len(bb) == 0 {
		return models.MasterBusiness{}, false, err
	}
	return models.MasterBusiness{ID: bb[0].s("ID"), Note: bb[0].s("NOTE"), OldID: bb[0].s("OLDID")}, true, nil
}
