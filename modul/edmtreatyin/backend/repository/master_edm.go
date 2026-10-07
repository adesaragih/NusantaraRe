package repository

// Untuk apa berkas ini: PEMBACA MASTER rantai `Activity/EDMChooseBusiness_Act` (`models.PembacaMasterEDM`) -
// padanan PERSIS empat RDB-List korpus EDM, BACA-SAJA JSONDATA (pengecualian K8 NB; nol penulisan JSON):
//
//	RDBList/BrowseTreatyInEDM (kelas Data-PolicyTreatyIn; SetValueEDM_Act 5):
//	  select JSONDATA as CLASSOFBUSINESS from pooldata.m_treaty_in_edm where OLDID = {OldData.NoOffer}
//	RDBList/BrowseTreatyIn (kelas Int-TREATY_IN, Access ASM; SetValueEDM_Act 6, SetTreatyIn_Act 4):
//	  select * from ( select JSONDATA ... from pooldata.M_TREATY_IN where ID={TreatyIn.ID}
//	                  union all select JSONDATA ... from pooldata.M_TREATY_IN_edm where ID={TreatyIn.ID} )
//	RDBList/BrowseTreatyInEDM_Int_treaty_in_edm (kelas Int-treaty_in_edm; SetTreatyInEDM_Act 3):
//	  select JSONDATA as ClassofBusiness from pooldata.M_TREATY_IN_EDM where ID={TreatyIn.ID}
//	RDBList/BrowseTreatyOutDetailEDM (kelas Int-TREATYOUTDETAIL; SetValueEDM_Act 7, XOL Retro):
//	  select JSONDATA as CLASSOFBUSINESS from pooldata.M_treaty_out where ID={OldData.NoOffer}
//	RDBList/GetDataCurrencyByName_SQL: select ID,... from CURRENCY where CURRENCY = {CARI8}
//
// Setiap baris diurai TERSENDIRI (`uraiMasterXOLPerBaris`); urutan baris = urutan RDB (tanpa ORDER BY di XML -
// urutan Oracle, `[dugaan]` tidak dijamin). Nol baris = daftar kosong (Pega: halaman tetap kosong), bukan galat.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

const tabelMasterTreatyOut = "M_TREATY_OUT"

// PembacaMasterEDM menyusun pembaca master atas gudang ini, terikat konteks permintaan.
func (g *Gudang) PembacaMasterEDM(ctx context.Context) models.PembacaMasterEDM {
	return pembacaMasterEDM{g: g, ctx: ctx}
}

type pembacaMasterEDM struct {
	g   *Gudang
	ctx context.Context
}

// sumberMaster - satu tabel master dan kolom kuncinya (klausa WHERE RDB-List).
type sumberMaster struct{ tabel, kunci string }

// baca - SATU-SATUNYA tempat kolom dokumen master disebut di berkas ini (penjaga F1
// `TestKolomDokumenHanyaDiMasterXOLDariJSON`); keluarannya selalu lewat `uraiMasterXOLPerBaris` (daftar medan
// tertutup, `TestPembacaMasterEDMHanyaMengembalikanHasilUrai`). Lebih dari satu sumber dirangkai `UNION ALL`
// berurutan dengan nilai kunci yang sama (RDB BrowseTreatyIn).
func (p pembacaMasterEDM) baca(apa, nilai string, sumber ...sumberMaster) ([]models.MasterXOL, error) {
	bagian := make([]string, 0, len(sumber))
	args := make([]any, 0, len(sumber))
	for i, s := range sumber {
		t, err := p.g.nama(s.tabel)
		if err != nil {
			return nil, err
		}
		bagian = append(bagian, fmt.Sprintf(`SELECT JSONDATA FROM %s WHERE %s = :%d`, t, s.kunci, i+1))
		args = append(args, nilai)
	}
	q := strings.Join(bagian, " UNION ALL ")
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := p.g.db.QueryContext(p.ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master treaty (%s): %w", apa, err)
	}
	defer rows.Close()
	var isi []string
	for rows.Next() {
		var v sql.NullString
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("repository: membaca master treaty (%s): %w", apa, err)
		}
		isi = append(isi, v.String)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca master treaty (%s): %w", apa, err)
	}
	return uraiMasterXOLPerBaris(isi)
}

// MasterMenurutOldID = RDB BrowseTreatyInEDM (kelas Data-PolicyTreatyIn).
func (p pembacaMasterEDM) MasterMenurutOldID(oldID string) ([]models.MasterXOL, error) {
	return p.baca("OLDID", oldID, sumberMaster{tabelMasterTreatyEDM, "OLDID"})
}

// MasterMenurutID = RDB BrowseTreatyIn (M_TREATY_IN ∪ M_TREATY_IN_EDM).
func (p pembacaMasterEDM) MasterMenurutID(id string) ([]models.MasterXOL, error) {
	return p.baca("ID", id, sumberMaster{tabelMasterTreaty, "ID"}, sumberMaster{tabelMasterTreatyEDM, "ID"})
}

// MasterEDMMenurutID = RDB BrowseTreatyInEDM (kelas Int-treaty_in_edm).
func (p pembacaMasterEDM) MasterEDMMenurutID(id string) ([]models.MasterXOL, error) {
	return p.baca("EDM ID", id, sumberMaster{tabelMasterTreatyEDM, "ID"})
}

// MasterOutMenurutID = RDB BrowseTreatyOutDetailEDM (M_TREATY_OUT, XOL Retro).
func (p pembacaMasterEDM) MasterOutMenurutID(id string) ([]models.MasterXOL, error) {
	return p.baca("treaty out", id, sumberMaster{tabelMasterTreatyOut, "ID"})
}

// IDMataUang = RDB GetDataCurrencyByName_SQL (`IDMataUangDariNama`).
func (p pembacaMasterEDM) IDMataUang(nama string) (string, error) {
	return p.g.IDMataUangDariNama(p.ctx, nama)
}
