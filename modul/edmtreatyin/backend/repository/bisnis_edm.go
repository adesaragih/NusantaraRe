package repository

// Untuk apa berkas ini: GRID POPUP "Business And SOB List" (`Harness/BusinessAndSOBListEDM`,
// `Section/BusinessAndSOBListEDM` kelas Data-PolicyTreatyIn) - daftar REVISI MASTER treaty yang ditawarkan tombol
// Choose Business. Kolom kedua grid: (Choose) · ID Revision `.ID` · Previous ID `.OLDID` · Treaty Contract Name ·
// Proportion Type · SOB `.LeadingReinsSource` · Ceding · Start Date `.Commencement` · End Date `.Termination`.
//
//	S6  `.EDMType != 3`  pre-load `Activity/TreatyLoadMasterJoinEdmChooseBusiness`:
//	    bukan XOL Retro  RDB TreatyLoadMasterJoinEdmChooseBusiness - pooldata.treaty_in A UNION ALL
//	                     pooldata.treaty_in_edm B, `ID like '<OldData.NoOffer 7 karakter>%'`, order by ID
//	    XOL Retro        RDB TreatyLoadMasterJoinEdmChooseBusinessRetro - pooldata.treaty_out2,
//	                     `ID like '<10 karakter>%'`
//	S1  `.EDMType = 3`   RD BrowseTREATY_IN_EDM (kelas Int-treaty_in_edm = TREATY_IN_EDM; pyMaxRecords 500, urut
//	                     .ID DESC): A `.ID Contains OldMasterNo.CARI1` (DT SetOldMasterNo: 7 karakter), J
//	                     `.StatusAkseptasi = "Resolve Complete"` (case-insensitive), L `.EDMState = "3"`; filter
//	                     lain ber-parameter kosong tidak berlaku.
//
// ⚠️ `{ASIS:InputData.CARI1}` (teks SQL ditempel) diganti penampung bind - nilai pola dari halaman, bukan ketikan.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

const (
	tabelTreatyIn    = "TREATY_IN"
	tabelTreatyInEDM = "TREATY_IN_EDM"
	tabelTreatyOut2  = "TREATY_OUT2"
)

// kolomPopupEDM - kolom SELECT kedua RDB (urut XML) beserta nama properti barisnya.
var kolomPopupEDM = []struct{ kolom, properti string }{
	{"ID", "ID"}, {"OLDID", "OLDID"}, {"TREATYCONTRACTNAME", "TreatyContractName"}, {"PROPORTIONTYPE", "ProportionType"},
	{"LEADINGREINSSOURCE", "LeadingReinsSource"}, {"CEDING", "Ceding"}, {"COMMENCEMENT", "Commencement"},
	{"TERMINATION", "Termination"},
}

// SaringanPopupEDM - masukan grid popup.
type SaringanPopupEDM struct {
	// OldNoOffer - `PolicyTreatyIn.OldData.NoOffer`.
	OldNoOffer string
	// EDMType - `.EDMType` ("3" -> RD BrowseTREATY_IN_EDM).
	EDMType string
	// Retro - `pyWorkPage.PolicyTreatyIn.ClaimType == "XOL Retro"`.
	Retro bool
}

// awalanKarakter = `@substring(s, 0, n)` (Java String.substring: n karakter pertama; lebih pendek = galat Java ->
// di sini seluruh teks).
func awalanKarakter(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// pilihTanpaOldID - daftar SELECT dengan OLDID = NULL (`(select null from dual) as OLDID`).
func pilihTanpaOldID(alias string) string {
	p := "TO_CHAR(" + alias + ".ID), NULL"
	for _, k := range kolomPopupEDM[2:] {
		p += ", TO_CHAR(" + alias + "." + k.kolom + ")"
	}
	return p
}

// DaftarBisnisEDM membaca isi grid popup menurut `s`.
func (g *Gudang) DaftarBisnisEDM(ctx context.Context, s SaringanPopupEDM) ([]models.Baris, error) {
	pilih := ""
	for i, k := range kolomPopupEDM {
		if i > 0 {
			pilih += ", "
		}
		pilih += "TO_CHAR(%[1]s." + k.kolom + ")"
	}
	var q string
	var args []any
	switch {
	case s.EDMType == "3":
		t, err := g.nama(tabelTreatyInEDM)
		if err != nil {
			return nil, err
		}
		q = fmt.Sprintf(`SELECT %s FROM %s B
		  WHERE B.ID LIKE :1 AND UPPER(B.STATUSAKSEPTASI) = UPPER(:2) AND B.EDMSTATE = :3
		  ORDER BY B.ID DESC FETCH FIRST 500 ROWS ONLY`, fmt.Sprintf(pilih, "B"), t)
		args = []any{"%" + awalanKarakter(s.OldNoOffer, 7) + "%", "Resolve Complete", "3"}
	case s.Retro:
		t, err := g.nama(tabelTreatyOut2)
		if err != nil {
			return nil, err
		}
		// RDB ...Retro: `(select null from dual) as OLDID` - TREATY_OUT2 tanpa kolom OLDID (DEV ORA-00904)
		q = fmt.Sprintf(`SELECT %s FROM %s A WHERE A.ID LIKE :1 ORDER BY A.ID ASC`, pilihTanpaOldID("A"), t)
		args = []any{awalanKarakter(s.OldNoOffer, 10) + "%"}
	default:
		ti, err := g.nama(tabelTreatyIn)
		if err != nil {
			return nil, err
		}
		te, err := g.nama(tabelTreatyInEDM)
		if err != nil {
			return nil, err
		}
		// OLDID baris treaty_in = `(select null from dual)`; urut `order by ID asc` atas hasil union.
		pilihA := pilihTanpaOldID("A")
		q = fmt.Sprintf(`SELECT * FROM (
		  SELECT %s FROM %s A WHERE A.ID LIKE :1
		  UNION ALL
		  SELECT %s FROM %s B WHERE B.ID LIKE :2
		) ORDER BY 1 ASC`, pilihA, ti, fmt.Sprintf(pilih, "B"), te)
		pola := awalanKarakter(s.OldNoOffer, 7) + "%"
		args = []any{pola, pola}
	}
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar revisi master: %w", err)
	}
	defer rows.Close()
	out := []models.Baris{}
	for rows.Next() {
		nilai := make([]sql.NullString, len(kolomPopupEDM))
		tujuan := make([]any, len(nilai))
		for i := range nilai {
			tujuan[i] = &nilai[i]
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: membaca baris revisi master: %w", err)
		}
		b := models.Baris{}
		for i, k := range kolomPopupEDM {
			if v := teks(nilai[i]); v != "" {
				b[k.properti] = v
			}
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
