package repository

// Untuk apa berkas ini: BACAAN DATA LAMA (tiket 13) - BACA SAJA, untuk sensus `alat/pemuatlama`: baris
// OS_AKSEPTASI_KLAIM kasus CLMP warisan, HISTORYAKSEPTASIPEGA-nya (tanpa USERNAME: nama orang tidak dibaca), dan work
// object KomiteTreaty di DATAPEGA (tanpa BLOB). Nol pernyataan tulis.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimprop/backend/models"
)

func sqlOSLama(t string) string {
	return fmt.Sprintf(`SELECT CASEID, TO_CHAR(STS_REJECT), JSON_VALUE(DATA_JSON, '$.AcceptedNo' RETURNING VARCHAR2(100)),
		       %s
		  FROM %s WHERE CASEID LIKE :1 ORDER BY CASEID, TANGGAL`, fmt.Sprintf(db.FmtTanggalOracle, "TANGGAL"), t)
}

func sqlRiwayatLama(t string) string {
	return fmt.Sprintf(`SELECT ID_PEGA, ID_KOMITE, STATUS FROM %s WHERE ID_PEGA LIKE :1`, t)
}

// sqlKomitePega - work object KomiteTreaty warisan (skema DATAPEGA eksplisit, seperti Copy Old NB / EDM Treaty In).
const sqlKomitePega = `SELECT PZINSKEY, PXCOVERINSKEY, PYSTATUSWORK FROM DATAPEGA.PC_ASM_FW_GCNMFW_WORK
	 WHERE PXOBJCLASS = :1`

func (g *Gudang) bacaBaris(ctx context.Context, q string, n int, args ...any) ([][]string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: bacaan data lama komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out [][]string
	for rows.Next() {
		v := make([]sql.NullString, n)
		t := make([]any, n)
		for i := range v {
			t[i] = &v[i]
		}
		if err := rows.Scan(t...); err != nil {
			return nil, fmt.Errorf("repository: memindai data lama komite: %w", err)
		}
		b := make([]string, n)
		for i := range v {
			b[i] = strings.TrimSpace(v[i].String)
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// BacaOSLama - baris OS_AKSEPTASI_KLAIM kasus CLMP warisan.
func (g *Gudang) BacaOSLama(ctx context.Context) ([]models.BarisOSLama, error) {
	t, err := g.db.Qualify("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return nil, err
	}
	rows, err := g.bacaBaris(ctx, sqlOSLama(t), 4, models.AwalanKunciKlaimLama+"%")
	if err != nil {
		return nil, err
	}
	out := make([]models.BarisOSLama, 0, len(rows))
	for _, r := range rows {
		out = append(out, models.BarisOSLama{CaseID: r[0], StsReject: r[1], AcceptedNo: r[2], Tanggal: waktuDB(r[3])})
	}
	return out, nil
}

// BacaRiwayatLama - HISTORYAKSEPTASIPEGA kasus CLMP warisan.
func (g *Gudang) BacaRiwayatLama(ctx context.Context) ([]models.RiwayatLama, error) {
	t, err := g.db.Qualify("HISTORYAKSEPTASIPEGA")
	if err != nil {
		return nil, err
	}
	rows, err := g.bacaBaris(ctx, sqlRiwayatLama(t), 3, models.AwalanKunciKlaimLama+"%")
	if err != nil {
		return nil, err
	}
	out := make([]models.RiwayatLama, 0, len(rows))
	for _, r := range rows {
		out = append(out, models.RiwayatLama{IDPega: r[0], IDKomite: r[1], Status: r[2]})
	}
	return out, nil
}

// BacaKomitePega - work object KomiteTreaty warisan di DATAPEGA.
func (g *Gudang) BacaKomitePega(ctx context.Context) ([]models.KomitePega, error) {
	rows, err := g.bacaBaris(ctx, sqlKomitePega, 3, models.KelasKomite)
	if err != nil {
		return nil, err
	}
	out := make([]models.KomitePega, 0, len(rows))
	for _, r := range rows {
		out = append(out, models.KomitePega{Kunci: r[0], Cover: r[1], StatusWork: r[2]})
	}
	return out, nil
}
