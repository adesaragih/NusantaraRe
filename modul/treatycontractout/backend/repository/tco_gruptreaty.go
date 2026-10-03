package repository

// Master grup treaty - dibaca saja (tiket 03 Treaty Contract Out).
//
// `[terverifikasi]` Pemilih `Treaty Group` di form tahun treaty
// (`Section/InputTreatyContract.xml` b7776 `pyListSource reportdefinition`,
// b7792 `pySourceName BrowseTreatyGroup_RD`; `Section/InputDtlTreatyContact.xml`
// b6624) membaca `ReportDefinition/BrowseTreatyGroup_RD.xml` (kelas
// `ASM-FW-GISFW-Int-TREATYGROUP` b40): kolom `.ID` (sort DESC b587),
// `.TreatyGroupName`, `.TreatyGroupSOAName`, `.OJKBusinessID`, …; saringan
// berparameter `.ID = Param.ID` b546-b555 (kosong = seluruhnya).
//
// Tabel fisik `POOLDATA.TREATYGROUP` - kolom `ID`, `TREATYGROUPNAME`, `OLDID`
// disebut `NB Treaty In/RDBList/FetchTreatyGroupOLDID.xml`
// (`select oldid, treatygroupname from pooldata.treatygroup where ID = …`).
// Dua kolom yang dibaca modul ini keduanya terbukti SQL.
//
// ⛔ BACA SAJA (spec b107); penjaga TestTCOWarisanHanyaDibaca.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
)

// MasterGrupTreatyTCO adalah nama tabel master grup treaty.
const MasterGrupTreatyTCO = "TREATYGROUP"

// GrupTreatyTCO adalah satu baris master seperti dibaca modul ini.
type GrupTreatyTCO struct {
	ID              string
	TreatyGroupName string
}

// sqlGrupTreatyTCO - urutan `.ID DESC` (b587).
func sqlGrupTreatyTCO(tabel string) string {
	return fmt.Sprintf(`SELECT ID, TREATYGROUPNAME FROM %s ORDER BY ID DESC`, tabel)
}

// MasterGrupTreaty membaca master grup treaty.
type MasterGrupTreaty struct{ db *db.DB }

// NewMasterGrupTreaty menyusunnya.
func NewMasterGrupTreaty(db *db.DB) *MasterGrupTreaty { return &MasterGrupTreaty{db: db} }

// Daftar membaca seluruh grup treaty.
func (m *MasterGrupTreaty) Daftar(ctx context.Context) ([]GrupTreatyTCO, error) {
	tabel, err := m.db.Qualify(MasterGrupTreatyTCO)
	if err != nil {
		return nil, err
	}
	q := sqlGrupTreatyTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: reading treaty group master: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []GrupTreatyTCO
	for rows.Next() {
		var id, nama sql.NullString
		if err := rows.Scan(&id, &nama); err != nil {
			return nil, err
		}
		out = append(out, GrupTreatyTCO{ID: id.String, TreatyGroupName: nama.String})
	}
	return out, rows.Err()
}
