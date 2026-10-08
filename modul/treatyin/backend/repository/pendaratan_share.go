package repository

// Tab Share Non-Prop — pohon dari tabel pendaratan, dan dua Report
// Definition spreading dari tabel master susunan treaty.
//
//	Share (T_TREATY_SHARE)
//	├ SpreadingListXOL (T_TREATY_SHARE_SPREADING)
//	├ DeductionList (T_TREATY_SHARE_DEDUCTION)
//	└ GrossPremiumList · NetPremiumList (T_TREATY_SHARE_AMOUNT, JENIS)
//	FacultativeShareList (T_TREATY_FAC_SHARE)
//	├ DeductionList (T_TREATY_FAC_SHARE_DEDUCTION)
//	└ GrossPremiumList · NetPremiumList (T_TREATY_FAC_SHARE_AMOUNT, JENIS)
//	ShareReins (T_TREATY_RETRO_SHARE)
//	ShareFacultativeReinsurers (T_TREATY_FAC_REINSURER)
//
// ⛔ Kunci dan kolomnya dari `PetaPendaratan` — peta yang MEMUAT tabel ini.
// Repository tidak menafsirkan; `services.ShareDariPendaratan` yang
// menafsirkan.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// tabelPohonShare — urutan baca.
var tabelPohonShare = []string{
	"T_TREATY_SHARE", "T_TREATY_SHARE_SPREADING", "T_TREATY_SHARE_DEDUCTION", "T_TREATY_SHARE_AMOUNT",
	"T_TREATY_FAC_SHARE", "T_TREATY_FAC_SHARE_DEDUCTION", "T_TREATY_FAC_SHARE_AMOUNT",
	"T_TREATY_RETRO_SHARE", "T_TREATY_FAC_REINSURER",
}

// BacaSharePendaratan membaca pohon tab Share satu kontrak.
func (g *Gudang) BacaSharePendaratan(ctx context.Context, masterID string) (models.SharePendaratan, error) {
	baris := make(map[string][]barisPendaratan, len(tabelPohonShare))
	for _, t := range tabelPohonShare {
		b, err := g.bacaEntri(ctx, t, masterID)
		if err != nil {
			return models.SharePendaratan{}, err
		}
		baris[t] = b
	}
	return RangkaiShare(baris), nil
}

// RangkaiShare merangkai pohon Share dari baris per tabel. Larik yang tak
// berbaris hadir sebagai larik KOSONG — layar membaca `.length`.
func RangkaiShare(baris map[string][]barisPendaratan) models.SharePendaratan {
	sp := models.SharePendaratan{
		Share:                      rangkaiBarisShare(baris, "T_TREATY_SHARE", "T_TREATY_SHARE_SPREADING", "T_TREATY_SHARE_DEDUCTION", "T_TREATY_SHARE_AMOUNT"),
		FacultativeShareList:       rangkaiBarisShare(baris, "T_TREATY_FAC_SHARE", "", "T_TREATY_FAC_SHARE_DEDUCTION", "T_TREATY_FAC_SHARE_AMOUNT"),
		ShareReins:                 rangkaiDatar(baris, "T_TREATY_RETRO_SHARE"),
		ShareFacultativeReinsurers: rangkaiDatar(baris, "T_TREATY_FAC_REINSURER"),
	}
	return sp
}

func rangkaiDatar(baris map[string][]barisPendaratan, tabel string) []map[string]any {
	p, _ := entriPeta(tabel)
	out := []map[string]any{}
	for _, b := range baris[tabel] {
		s := simpulPeta(b, p)
		s["ID"] = fmt.Sprint(b.ID)
		out = append(out, s)
	}
	return out
}

func rangkaiBarisShare(baris map[string][]barisPendaratan, induk, tabelSebar, tabelDeduksi, tabelNilai string) []map[string]any {
	p, _ := entriPeta(induk)
	var sebar anakTabel
	if tabelSebar != "" {
		sebar = kelompokkan(tabelSebar, baris[tabelSebar])
	}
	deduksi := kelompokkan(tabelDeduksi, baris[tabelDeduksi])
	nilai := kelompokkan(tabelNilai, baris[tabelNilai])
	out := []map[string]any{}
	for _, b := range baris[induk] {
		s := simpulPeta(b, p)
		if tabelSebar != "" {
			s["SpreadingListXOL"] = isiAtauKosong(sebar.biasa[b.ID])
		} else {
			s["SpreadingListXOL"] = []map[string]any{}
		}
		s["DeductionList"] = isiAtauKosong(deduksi.biasa[b.ID])
		pasangLarikGabung(s, nilai, b.ID)
		out = append(out, s)
	}
	return out
}

// TabelSusunanTreaty / TabelTahunTreaty — kelas Pega
// `ASM-FW-GISFW-Int-PROPORTIONALARRG` dan `ASM-FW-GISFW-Int-TREATYYEAR`.
// Dibaca saja.
const (
	TabelSusunanTreaty = "PROPORTIONALARRG"
	TabelTahunTreaty   = "TREATYYEAR"
)

// SaringanIndukSpreading — `BrowseTreatyArrangement_ParentReinsMasterTrt`,
// logika `B AND C AND D AND (E OR A) AND F AND G AND H`, join `TrtYr`
// (`TREATYYEAR`) lewat `.TreatyYearID = TrtYr.ID`:
//
//	B .TreatyGroupID     = Param.TreatyGroupID
//	C .TreatyDescID      = "10001"   (semua pemanggil mengirim nilai ini)
//	D .ParentReinsTypeID = "00"
//	E .ReinsTypeName     Contains "TRT"   A .ReinsTypeName = "ORS"
//	F TrtYr.StartDate   <= Param.StartDate
//	G TrtYr.EndDate     >= Param.StartDate
//	H .ReinsTypeID      != Param.ReinsTypeID
//
// ⭐ KEDELAPAN filter TANPA `pyUseNullIfEmpty` — parameter KOSONG membuat
// filternya DILEWATI Pega (pembaca bersama `filter_diabaikan`). Itu bukan
// teori di sini: dropdown Reins Type spreading manual mengirim
// `TempSprd.TreatyGroupID`, yang tidak pernah diisi untuk baris Share
// Non-Prop — jadi B dilewati dan daftarnya induk SEMUA Treaty Group. Dan H
// hanya dikirim kedua dropdown (`"10246"` = `2025 XOL TRT`); pencarian di
// `FetchQSfromMasterXOL` / `SetSpreadingXOL` tidak mengirimnya.
//
// ⚠️ `NVL(…, '~')` pada H: `!=` Pega atas nilai kosong MELOLOSKAN baris,
// `<>` Oracle atas `NULL` menolaknya.
//
// ⚠️ Tanggal dibandingkan sebagai TEKS `YYYYMMDD` — bentuk `COMMENCEMENT`
// dan `TREATYYEAR.STARTDATE/ENDDATE` di basis data (terukur 7 Oktober 2026).
func SaringanIndukSpreading(treatyGroupID, tanggalMulai, kecuali string) (string, []any) {
	w := `p.TREATYDESCID = '10001' AND p.PARENTREINSTYPEID = '00' ` +
		`AND (INSTR(p.REINSTYPENAME, 'TRT') > 0 OR p.REINSTYPENAME = 'ORS')`
	arg := []any{}
	tambah := func(klausa string, nilai string) {
		arg = append(arg, nilai)
		w += fmt.Sprintf(" AND "+klausa, len(arg))
	}
	if treatyGroupID != "" {
		tambah("p.TREATYGROUPID = :%d", treatyGroupID)
	}
	if tanggalMulai != "" {
		tambah("t.STARTDATE <= :%d", tanggalMulai)
		tambah("t.ENDDATE >= :%d", tanggalMulai)
	}
	if kecuali != "" {
		tambah("NVL(p.REINSTYPEID, '~') <> :%d", kecuali)
	}
	return w, arg
}

// IndukDikecualikanDropdown — `ReinsTypeID` yang kedua dropdown spreading
// kirim sebagai parameter H (`2025 XOL TRT`).
const IndukDikecualikanDropdown = "10246"

// BacaIndukSpreading — RD `ParentReinsMasterTrt`. `kecuali` = parameter H
// (kosong → filter dilewati).
//
// ⚠️ RD-nya TANPA urutan; di sini `ORDER BY p.ID` supaya "yang pertama"
// tidak bergantung pada urutan fisik tabel.
func (g *Gudang) BacaIndukSpreading(ctx context.Context, treatyGroupID, tanggalMulai, kecuali string) ([]models.SusunanSpreading, error) {
	arrg, err := g.db.Qualify(TabelSusunanTreaty)
	if err != nil {
		return nil, err
	}
	tahun, err := g.db.Qualify(TabelTahunTreaty)
	if err != nil {
		return nil, err
	}
	w, arg := SaringanIndukSpreading(treatyGroupID, tanggalMulai, kecuali)
	q := fmt.Sprintf(`SELECT p.REINSTYPEID, p.REINSTYPENAME, p.PARENTREINSTYPEID, p.TREATYYEARID, p.PCT, p.RP, p.USD, p.TREATYYEAR
		FROM %s p JOIN %s t ON t.ID = p.TREATYYEARID
		WHERE %s ORDER BY p.ID FETCH FIRST 500 ROWS ONLY`, arrg, tahun, w)
	return g.bacaSusunan(ctx, q, arg...)
}

func (g *Gudang) BacaAnakSpreading(ctx context.Context, parentReinsTypeID, treatyYearID string) ([]models.SusunanSpreading, error) {
	if parentReinsTypeID == "" {
		return []models.SusunanSpreading{}, nil
	}
	arrg, err := g.db.Qualify(TabelSusunanTreaty)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT REINSTYPEID, REINSTYPENAME, PARENTREINSTYPEID, TREATYYEARID, PCT, RP, USD, TREATYYEAR
		FROM %s WHERE TREATYDESCID = '10001' AND PARENTREINSTYPEID = :1 AND NVL(TREATYYEARID, '~') = NVL(:2, '~')
		ORDER BY ID FETCH FIRST 500 ROWS ONLY`, arrg)
	return g.bacaSusunan(ctx, q, parentReinsTypeID, treatyYearID)
}

// SaringanAnakSpreadingProp — `BrowseTreatyArrangement_Limit_MstTrt_RD`
// seperti cabang PROPORSIONAL memanggilnya (`FetchQSfromMaster`,
// `SetSpreadName`), logika `A AND B AND C AND D AND E`:
//
//	A .TreatyYear        = Param.TreatyYear
//	B .TreatyGroupID     = Param.TreatyGroupID
//	C .TreatyDescID      = Param.TreatyDescID
//	D .ParentReinsTypeID = Param.ParentReinsTypeID
//	E .TreatyYearID      = Param.TreatyYearID
//
// ⭐ Kelima filter TANPA `pyUseNullIfEmpty` — parameter KOSONG membuat
// filternya DILEWATI (sama dengan `SaringanIndukSpreading`). Ini bukan
// teori: `TreatyInPropshare` memanggil `FetchQSfromMaster` TANPA
// `ParentReinsTypeID`, sehingga langkah 6.1 tidak pernah mengisi
// `Local.TreatyYear`/`Param.TreatyYearID` dan A serta E dilewati; dan
// `SetSpreadName` 3.2.5 mengosongkan A, B, C dengan sengaja.
func SaringanAnakSpreadingProp(treatyYear, treatyGroupID, treatyDescID, parent, treatyYearID string) (string, []any) {
	w := "1 = 1"
	arg := []any{}
	for _, f := range []struct{ kolom, nilai string }{
		{"TREATYYEAR", treatyYear},
		{"TREATYGROUPID", treatyGroupID},
		{"TREATYDESCID", treatyDescID},
		{"PARENTREINSTYPEID", parent},
		{"TREATYYEARID", treatyYearID},
	} {
		if f.nilai == "" {
			continue
		}
		arg = append(arg, f.nilai)
		w += fmt.Sprintf(" AND %s = :%d", f.kolom, len(arg))
	}
	return w, arg
}

// BacaAnakSpreadingProp — RD `Limit_MstTrt_RD` cabang Prop. Tanpa satu pun
// parameter yang terisi RD-nya membaca SELURUH susunan; di sini pagar:
// `ParentReinsTypeID` kosong → nol baris (setiap pemanggil Prop
// mengirimnya dari baris yang punya Spreading Type / Reins Type).
//
// ⚠️ RD-nya TANPA urutan; `ORDER BY ID` supaya urutan baris stabil.
func (g *Gudang) BacaAnakSpreadingProp(ctx context.Context, treatyYear, treatyGroupID, treatyDescID, parent, treatyYearID string) ([]models.SusunanSpreading, error) {
	if parent == "" {
		return []models.SusunanSpreading{}, nil
	}
	arrg, err := g.db.Qualify(TabelSusunanTreaty)
	if err != nil {
		return nil, err
	}
	w, arg := SaringanAnakSpreadingProp(treatyYear, treatyGroupID, treatyDescID, parent, treatyYearID)
	q := fmt.Sprintf(`SELECT REINSTYPEID, REINSTYPENAME, PARENTREINSTYPEID, TREATYYEARID, PCT, RP, USD, TREATYYEAR
		FROM %s WHERE %s ORDER BY ID FETCH FIRST 500 ROWS ONLY`, arrg, w)
	return g.bacaSusunan(ctx, q, arg...)
}

func (g *Gudang) bacaSusunan(ctx context.Context, q string, arg ...any) ([]models.SusunanSpreading, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", TabelSusunanTreaty, err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.SusunanSpreading{}
	for rows.Next() {
		var id, nm, induk, thn, pct, rp, usd, tahun sql.NullString
		if err := rows.Scan(&id, &nm, &induk, &thn, &pct, &rp, &usd, &tahun); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", TabelSusunanTreaty, err)
		}
		out = append(out, models.SusunanSpreading{
			ReinsTypeID: id.String, ReinsTypeName: nm.String, ParentReinsTypeID: induk.String,
			TreatyYearID: thn.String, Pct: pct.String, Rp: rp.String, Usd: usd.String, TreatyYear: tahun.String,
		})
	}
	return out, rows.Err()
}

// SaringanReasuradurShare — `BrowseAgentNusaRe_RD` seperti autocomplete
// `Reinsurer Name` / `Facultative Reinsurers` memanggilnya: A·B·G saja.
//
// ⚠️ BEDA dengan `SaringanAgenPega` (pemilih Ceding): autocomplete ini
// mengirim `ChildCount` TANPA nilai, dan filter D (`.ChildCount =
// Param.ChildCount`) tidak ber-`pyUseNullIfEmpty` — Pega MENGABAIKAN filter
// berparameter kosong. Jadi `CHILDCOUNT` tidak disaring di sini.
const SaringanReasuradurShare = `STATUSACTIVE = '1' AND NVL(AGENTTPYE2, '~') <> 'LIFE INSURANCE' AND CLIENTID IS NOT NULL`

// BacaDaftarReasuradurShare — isi autocomplete `Reinsurer Name`: nilai
// `.ClientName`, `pyAdditionalFields` `.ID → .ReinsID`; urut `.ID` menurun
// (sort RD). Filter E (`ClientName` Contains, tak peka huruf) dan batas 200
// baris dikerjakan layar atas kata yang diketik.
func (g *Gudang) BacaDaftarReasuradurShare(ctx context.Context) ([]models.PilihanWarisan, error) {
	nama, err := g.db.Qualify(TabelAgen)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(ID), CLIENTNAME FROM %s WHERE %s ORDER BY ID DESC`, nama, SaringanReasuradurShare)
	rows, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca reasuradur dari %s: %w", TabelAgen, err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.PilihanWarisan{}
	for rows.Next() {
		var id, nm sql.NullString
		if err := rows.Scan(&id, &nm); err != nil {
			return nil, fmt.Errorf("repository: membaca baris agen: %w", err)
		}
		out = append(out, models.PilihanWarisan{ID: id.String, Nama: nm.String})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	tandaiKembar(out)
	return out, nil
}
