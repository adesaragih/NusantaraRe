package repository

// Kontrak On Retention untuk `OutwardList` (paket 9, RALAT R9/P3) - satu-satunya
// pembaca `TREATYCONTRACT_LIFE` / `TREATYYEAR_LIFE` di modul ini (R12), dibaca saja.
//
//	checkbox `On Retention` b47312 → onChange `GetReinsTypeOR_Life` b47488:
//	  1 b236 `·` PRE=false Property-Remove `ProductName.OutwardList`
//	  2 b383 `·` `Temp.CARI1/2 ← ProductNameInward.BEGIN/MATURE`
//	  3 b534 `·` PRE=false RDB-List `BrowseReinstypeOR_SQL`
//	  4 b731 `·` PRE=false ulang `Reinstype.pxResults`
//	  4.1 b770 `·` `OutwardList(<APPEND>)` ← REINSTYPEID, REINSTYPENAME, TREATYYEAR,
//	        TREATYCONTRACTID, UNDERWRITINGYEAR
//
//	BrowseReinstypeOR_SQL b84:
//	  SELECT tc.*, ty.* FROM pooldata.treatycontract_life tc
//	  JOIN pooldata.treatyyear_life ty ON ty.id = tc.idtreatyyear
//	  WHERE TO_DATE({Temp.CARI1}, 'DD/MM/YYYY') >= tc.treatystartdate
//	    AND TO_DATE({Temp.CARI2}, 'DD/MM/YYYY') <= tc.treatyenddate
//	    AND REINSTYPEID = '10200'
//
// ⚠️ `tc.*, ty.*` tidak memuat kolom bernama `TREATYCONTRACTID` (DDL kedua tabel,
// `mastercontractretrolife/docs/ddl-tables-from-dba.md`): di Pega `.TREATYCONTRACTID`
// baris hasil selalu kosong. Ditiru - `""` (OQ-MPNL-15 ditutup 01-10-2026: kolom itu tidak ada di DEV). Kolom dipilih menurut
// nama, bukan `*`; `ORDER BY tc.ID` ditambahkan agar urutan pasti (RDB tanpa urutan).

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Master kontrak treaty - dibaca saja.
const (
	MasterKontrakTreaty = "TREATYCONTRACT_LIFE"
	MasterTahunTreaty   = "TREATYYEAR_LIFE"
)

// ReinsTypeOR - `BrowseReinstypeOR_SQL` b84 `REINSTYPEID = '10200'` (jenis OR).
const ReinsTypeOR = "10200"

// Kolom yang disebut SQL On Retention (penjaga kata cadangan).
var (
	KolomKontrakTreaty = []string{"ID", "IDTREATYYEAR", "REINSTYPEID", "REINSTYPENAME", "TREATYSTARTDATE", "TREATYENDDATE"}
	KolomTahunTreaty   = []string{"ID", "TREATYYEAR", "UNDERWRITINGYEAR"}
)

func sqlReinstypeOR(kontrak, tahun string) string {
	return fmt.Sprintf(`SELECT tc.REINSTYPEID, tc.REINSTYPENAME, ty.TREATYYEAR, ty.UNDERWRITINGYEAR
		FROM %s tc JOIN %s ty ON ty.ID = tc.IDTREATYYEAR
		WHERE TO_DATE(:1, 'DD/MM/YYYY') >= tc.TREATYSTARTDATE
		  AND TO_DATE(:2, 'DD/MM/YYYY') <= tc.TREATYENDDATE
		  AND tc.REINSTYPEID = '%s'
		ORDER BY tc.ID ASC`, kontrak, tahun, ReinsTypeOR)
}

// DaftarReinstypeOR - baris `OutwardList` untuk periode `BEGIN`–`MATURE`
// (`dd/MM/yyyy`). Tanggal kosong = nol baris (di Pega `TO_DATE` atas teks kosong = NULL).
func (g *Gudang) DaftarReinstypeOR(ctx context.Context, tx *db.Tx, begin, mature string) ([]models.BarisOutward, error) {
	hasil := []models.BarisOutward{}
	if begin == "" || mature == "" {
		return hasil, nil
	}
	kontrak, err := g.db.Qualify(MasterKontrakTreaty)
	if err != nil {
		return nil, err
	}
	tahun, err := g.db.Qualify(MasterTahunTreaty)
	if err != nil {
		return nil, err
	}
	q := sqlReinstypeOR(kontrak, tahun)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.kueri(tx).QueryContext(ctx, q, begin, mature)
	if err != nil {
		return nil, galatMaster(MasterKontrakTreaty, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id, nama, tahunTreaty, uw sql.NullString
		if err := rows.Scan(&id, &nama, &tahunTreaty, &uw); err != nil {
			return nil, galatMaster(MasterKontrakTreaty, err)
		}
		hasil = append(hasil, models.BarisOutward{ReinsTypeID: id.String, ReinsTypeName: nama.String,
			TransactionYear: tahunTreaty.String, UnderwritingYear: uw.String})
	}
	if err := rows.Err(); err != nil {
		return nil, galatMaster(MasterKontrakTreaty, err)
	}
	return hasil, nil
}
