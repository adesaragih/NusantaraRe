package repository

// Untuk apa berkas ini: ACUAN - bacaan baca-saja yang dibutuhkan port activity (`models.Acuan`) dan pemilih layar Claim
// Non Prop. Setiap kueri meniru satu RDB-List / Report Definition korpus `Claim Non Prop/RDBList` (nama rule di
// komentar); alias CARI diluruskan. Pola `modul/claimprop/backend/repository/acuan.go` (disalin, bukan diimpor).
//
// ⚠️ Dokumen JSON warisan (M_TREATY_IN.JSONDATA, JSON_KLAIM.DATA_JSON, OS_AKSEPTASI_KLAIM.DATA_JSON) DIBACA lewat
// JSON_VALUE / JSON_TABLE seperti rule-nya - baca-saja.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimnonprop/backend/models"
)

// Acuan - bacaan acuan Oracle (memenuhi `models.Acuan`).
type Acuan struct {
	g *Gudang
	// Produksi - padanan `When/IsPEGAPROD`: bacaan yang di korpus hanya berjalan di produksi (getStatusKonversi_SQL).
	Produksi bool
}

// AcuanDari menyusun acuan atas gudang.
func AcuanDari(g *Gudang, produksi bool) *Acuan { return &Acuan{g: g, Produksi: produksi} }

func (a *Acuan) q(objek string) (string, error) { return a.g.db.Qualify(objek) }

// satu menjalankan kueri satu kolom satu baris; tidak ada baris = "" tanpa galat.
func (a *Acuan) satu(ctx context.Context, q string, args ...any) (string, bool, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return "", false, err
	}
	var v sql.NullString
	err := a.g.db.QueryRowContext(ctx, q, args...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("repository: bacaan acuan: %w", err)
	}
	return rapikanDesimal(v.String), true, nil
}

// banyak menjalankan kueri n kolom; setiap baris = []string.
func (a *Acuan) banyak(ctx context.Context, q string, n int, args ...any) ([][]string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := a.g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: bacaan acuan: %w", err)
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
			return nil, fmt.Errorf("repository: memindai acuan: %w", err)
		}
		s := make([]string, n)
		for i := range v {
			s[i] = rapikanDesimal(v[i].String)
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func dec(kol string) string { return fmt.Sprintf(db.FmtDesimal, kol) }

// NamaMataUang = GetCurrency (`CURRENCY FROM CURRENCY WHERE ID`).
func (a *Acuan) NamaMataUang(ctx context.Context, cur string) (string, error) {
	t, err := a.q("CURRENCY")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT CURRENCY FROM %s WHERE ID = :1`, t), cur)
	return v, err
}

// JenisReasXOL = GetReinsuranceTypeBYName_SQL (`ID, NOTE FROM REINSURANCETYPE WHERE TYPE = '4' AND FLAG = 'active' AND
// NOTE = nama`).
func (a *Acuan) JenisReasXOL(ctx context.Context, nama string) ([]models.BarisReinsType, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE TYPE = '4' AND FLAG = 'active' AND NOTE = :1
		ORDER BY ID`, t), 2, nama)
	if err != nil {
		return nil, err
	}
	var out []models.BarisReinsType
	for _, r := range rows {
		out = append(out, models.BarisReinsType{ID: r[0], Note: r[1]})
	}
	return out, nil
}

// ---------------------------------------------------------------- master treaty

// medanMaster - medan puncak JSON master Non Prop yang dibaca rule Claim Non Prop (urutan = kolom kueri).
var medanMaster = []string{"ProportionType", "Ceding", "CedingID", "LeadingReinsSource", "LeadingReinsSourceID",
	"Bordeaux", "BordereauxNote", "AccountingMode", "AccountingModeNonProp", "TeritorialScope", "Commencement",
	"Termination", "TreatyYear", "RNMShare", "EDMState", "StatusAkseptasi"}

// MasterTreaty = GetLimitsTreatyIn_SQL (`JSONDATA FROM M_TREATY_IN WHERE ID UNION ALL ... M_TREATY_IN_EDM`) +
// adoptJSONObject (SetValueClaimTNP_Act langkah 3-4): baris pertama yang ada (M_TREATY_IN lebih dulu).
func (a *Acuan) MasterTreaty(ctx context.Context, id string) (models.MasterTreaty, bool, error) {
	for _, tabel := range []string{"M_TREATY_IN", "M_TREATY_IN_EDM"} {
		m, ada, err := a.bacaMaster(ctx, tabel, id)
		if err != nil || ada {
			return m, ada, err
		}
	}
	return models.MasterTreaty{}, false, nil
}

func jv(p string) string {
	return fmt.Sprintf(`JSON_VALUE(JSONDATA, '$.%s' RETURNING VARCHAR2(4000))`, p)
}

func (a *Acuan) bacaMaster(ctx context.Context, objek, id string) (models.MasterTreaty, bool, error) {
	t, err := a.q(objek)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	var kol []string
	for _, m := range medanMaster {
		kol = append(kol, jv(m))
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, %s FROM %s WHERE ID = :1`, strings.Join(kol, ", "), t),
		len(medanMaster)+1, id)
	if err != nil || len(rows) == 0 {
		return models.MasterTreaty{}, false, err
	}
	r := rows[0]
	m := models.MasterTreaty{ID: r[0], ProportionType: r[1], Ceding: r[2], CedingID: r[3], LeadingReinsSource: r[4],
		LeadingReinsSourceID: r[5], Bordeaux: r[6], BordereauxNote: r[7], AccountingMode: r[8],
		AccountingModeNonProp: r[9], TeritorialScope: r[10], Commencement: tglPega(r[11]), Termination: tglPega(r[12]),
		TreatyYear: r[13], RNMShare: r[14], EDMState: r[15], StatusAkseptasi: r[16]}
	lim, err := a.banyak(ctx, sqlLimitsMaster(t), 14, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, l := range lim {
		m.Limits = append(m.Limits, models.LimitXOL{Layer: l[1], LayerType: l[2], LayerPart: l[3], LayerPartType: l[4],
			Currency: l[5], Currency2: l[6], Limit: l[7], Limit2: l[8], Deductible: l[9], Deductible2: l[10],
			ReinstatementPct: l[11], IsCombineMDP: l[12], NoRIPCalculation: l[13]})
	}
	grup, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.li, jt.tg FROM %s m, JSON_TABLE(m.JSONDATA, '$.Limits[*]' COLUMNS
		(li FOR ORDINALITY, NESTED PATH '$.TreatyGroupList[*]' COLUMNS (gi FOR ORDINALITY,
		tg VARCHAR2(400) PATH '$.TreatyGroup'))) jt WHERE m.ID = :1 AND jt.tg IS NOT NULL ORDER BY jt.li, jt.gi`, t), 2, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, g := range grup {
		if i := angkaBulat(g[0]) - 1; i >= 0 && i < len(m.Limits) {
			m.Limits[i].TreatyGroups = append(m.Limits[i].TreatyGroups, g[1])
		}
	}
	mdp, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.li, jt.cur, jt.val FROM %s m, JSON_TABLE(m.JSONDATA, '$.Limits[*]'
		COLUMNS (li FOR ORDINALITY, NESTED PATH '$.MDPList[*]' COLUMNS (mi FOR ORDINALITY,
		cur VARCHAR2(100) PATH '$.Currency', val VARCHAR2(100) PATH '$.Value'))) jt
		WHERE m.ID = :1 AND (jt.cur IS NOT NULL OR jt.val IS NOT NULL) ORDER BY jt.li, jt.mi`, t), 3, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, x := range mdp {
		if i := angkaBulat(x[0]) - 1; i >= 0 && i < len(m.Limits) {
			m.Limits[i].MDPList = append(m.Limits[i].MDPList, models.NilaiMataUang{Currency: x[1], Value: x[2]})
		}
	}
	ring, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.ly, jt.lt, jt.mdp, jt.mdp2, jt.cur, jt.cur2 FROM %s m,
		JSON_TABLE(m.JSONDATA, '$.LimitSummaryList[*]' COLUMNS (si FOR ORDINALITY, ly VARCHAR2(100) PATH '$.Layer',
		lt VARCHAR2(100) PATH '$.LayerType', mdp VARCHAR2(100) PATH '$.MDP', mdp2 VARCHAR2(100) PATH '$.MDP2',
		cur VARCHAR2(100) PATH '$.Currency', cur2 VARCHAR2(100) PATH '$.Currency2')) jt WHERE m.ID = :1 ORDER BY jt.si`, t),
		6, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, s := range ring {
		m.LimitSummaryList = append(m.LimitSummaryList, models.RingkasanLimit{Layer: s[0], LayerType: s[1], MDP: s[2],
			MDP2: s[3], Currency: s[4], Currency2: s[5]})
	}
	kurs, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.cur, jt.conv FROM %s m, JSON_TABLE(m.JSONDATA, '$.CurrencyList[*]'
		COLUMNS (ci FOR ORDINALITY, cur VARCHAR2(100) PATH '$.Currency', conv VARCHAR2(100) PATH '$.Conversion')) jt
		WHERE m.ID = :1 ORDER BY jt.ci`, t), 2, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, k := range kurs {
		m.CurrencyList = append(m.CurrencyList, models.KursMaster{Currency: k[0], Conversion: k[1]})
	}
	share, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.si, jt.tid, jt.tnm, jt.pct FROM %s m,
		JSON_TABLE(m.JSONDATA, '$.Share[*]' COLUMNS (si FOR ORDINALITY, tid VARCHAR2(100) PATH '$.SpreadingTypeIDXOL',
		tnm VARCHAR2(400) PATH '$.SpreadingTypeXOL', pct VARCHAR2(100) PATH '$.SpreadingTotalPctXOL')) jt
		WHERE m.ID = :1 ORDER BY jt.si`, t), 4, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, s := range share {
		m.Share = append(m.Share, models.ShareXOL{SpreadingTypeIDXOL: s[1], SpreadingTypeXOL: s[2],
			SpreadingTotalPctXOL: s[3]})
	}
	spr, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.si, jt.ri, jt.rn, jt.pct FROM %s m,
		JSON_TABLE(m.JSONDATA, '$.Share[*]' COLUMNS (si FOR ORDINALITY, NESTED PATH '$.SpreadingListXOL[*]' COLUMNS
		(xi FOR ORDINALITY, ri VARCHAR2(100) PATH '$.ReinsTypeID', rn VARCHAR2(400) PATH '$.ReinsTypeName',
		pct VARCHAR2(100) PATH '$.Pct'))) jt WHERE m.ID = :1 AND (jt.ri IS NOT NULL OR jt.rn IS NOT NULL)
		ORDER BY jt.si, jt.xi`, t), 4, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, s := range spr {
		if i := angkaBulat(s[0]) - 1; i >= 0 && i < len(m.Share) {
			m.Share[i].SpreadingListXOL = append(m.Share[i].SpreadingListXOL,
				models.SpreadingMaster{ReinsTypeID: s[1], ReinsTypeName: s[2], Pct: s[3]})
		}
	}
	return m, true, nil
}

// sqlLimitsMaster - `TreatyInMaster.Limits(n)` (layer XoL) satu master, urut indeks.
func sqlLimitsMaster(t string) string {
	kol := func(nama, p string) string { return fmt.Sprintf("%s VARCHAR2(400) PATH '$.%s'", nama, p) }
	return fmt.Sprintf(`SELECT jt.li, jt.ly, jt.lt, jt.lp, jt.lpt, jt.cur, jt.cur2, jt.lim, jt.lim2, jt.ded, jt.ded2,
		jt.rp, jt.cmb, jt.nrip FROM %s m, JSON_TABLE(m.JSONDATA, '$.Limits[*]' COLUMNS (li FOR ORDINALITY, %s, %s, %s, %s,
		%s, %s, %s, %s, %s, %s, %s, %s, %s)) jt WHERE m.ID = :1 ORDER BY jt.li`, t,
		kol("ly", "Layer"), kol("lt", "LayerType"), kol("lp", "LayerPart"), kol("lpt", "LayerPartType"),
		kol("cur", "Currency"), kol("cur2", "Currency2"), kol("lim", "Limit"), kol("lim2", "Limit2"),
		kol("ded", "Deductible"), kol("ded2", "Deductible2"), kol("rp", "ReinstatementPct"),
		kol("cmb", "IsCombineMDP"), kol("nrip", "NoRIPCalculation"))
}

// tglPega - "yyyyMMdd" master -> "2006-01-02".
func tglPega(s string) string {
	if t, ok := models.TanggalPega(s); ok {
		return models.TeksTanggal(t)
	}
	return s
}

// sqlDaftarMaster - popup ChooseMasterTNonProp: BrowseDtlTreatyNP (TREATYINDETAIL ∪ TREATYINDETAILEDM) atau
// BrowseDtlTreatyNPEDM (TREATYINDETAILEDM saja), `PROPORTIONTYPE = 'NonProportional'`, DISTINCT. Saringan per kolom
// (mengandung, tanpa beda huruf) dijalankan di server sebelum batas `models.BatasMaster` (pola Claim Prop).
func sqlDaftarMaster(detail, edm string, s models.SaringanMaster) (string, []any) {
	var args []any
	ph := func(v any) string { args = append(args, v); return fmt.Sprintf(":%d", len(args)) }
	saring := func() string {
		w := "PROPORTIONTYPE = " + ph(models.ProporsiMaster)
		for _, f := range []struct{ kol, nilai string }{
			{"TREATYID", s.TreatyID}, {"TREATYCONTRACTNAME", s.ContractName}, {"SOB", s.SOB}, {"CEDING", s.Ceding},
			{"TREATYGROUP", s.TreatyGroup}, {"CLASSOFBUSINESS", s.ClassOfBusiness}, {"TREATYYEAR", s.TreatyYear},
		} {
			if v := strings.TrimSpace(f.nilai); v != "" {
				w += " AND UPPER(" + f.kol + ") LIKE " + ph("%"+strings.ToUpper(v)+"%")
			}
		}
		return w
	}
	kol := `TREATYID, TREATYCONTRACTNAME, PROPORTIONTYPE, CLASSOFBUSINESSID, CLASSOFBUSINESS, SOB, CEDING, TREATYYEAR,
		TREATYGROUP`
	var bagian []string
	if s.Sumber != models.SumberMasterINEDM {
		bagian = append(bagian, fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, kol, detail, saring()))
	}
	bagian = append(bagian, fmt.Sprintf(`SELECT %s FROM %s WHERE %s`, kol, edm, saring()))
	return fmt.Sprintf(`SELECT DISTINCT %s FROM (%s) ORDER BY TREATYYEAR DESC, TREATYID DESC FETCH FIRST %d ROWS ONLY`,
		kol, strings.Join(bagian, " UNION ALL "), models.BatasMaster), args
}

// DaftarMaster = grid popup ChooseMasterTNonProp (BrowseDtlMasterTNP_Act: Param.TreatyIN IN / INEDM).
func (a *Acuan) DaftarMaster(ctx context.Context, s models.SaringanMaster) ([]models.BarisMaster, error) {
	detail, err := a.q("TREATYINDETAIL")
	if err != nil {
		return nil, err
	}
	edm, err := a.q("TREATYINDETAILEDM")
	if err != nil {
		return nil, err
	}
	q, args := sqlDaftarMaster(detail, edm, s)
	rows, err := a.banyak(ctx, q, 9, args...)
	if err != nil {
		return nil, err
	}
	out := []models.BarisMaster{}
	for _, r := range rows {
		out = append(out, models.BarisMaster{TreatyID: r[0], TreatyContractName: r[1], ProportionType: r[2],
			ClassOfBusinessID: r[3], ClassOfBusiness: r[4], SOB: r[5], Ceding: r[6], TreatyYear: r[7], TreatyGroup: r[8]})
	}
	return out, nil
}

// BarisMasterDari - satu baris popup master menurut TREATYID + CLASSOFBUSINESSID + TREATYGROUP (parameter tombol Choose
// dibaca ULANG di server, bukan dari layar).
func (a *Acuan) BarisMasterDari(ctx context.Context, sumber, treatyID, cobID, grup string) (models.BarisMaster, bool, error) {
	rows, err := a.DaftarMaster(ctx, models.SaringanMaster{Sumber: sumber, TreatyID: treatyID})
	if err != nil {
		return models.BarisMaster{}, false, err
	}
	for _, r := range rows {
		if r.TreatyID == treatyID && r.ClassOfBusinessID == cobID && r.TreatyGroup == grup {
			return r, true, nil
		}
	}
	return models.BarisMaster{}, false, nil
}

// ---------------------------------------------------------------- polis

// AdaPolisMaster = GetNopolis_SQL (`DISTINCT NOPOLIS FROM TREATYINPRODUCTION WHERE NOOFFER = :id`).
func (a *Acuan) AdaPolisMaster(ctx context.Context, idMaster string) (bool, error) {
	t, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return false, err
	}
	_, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT NOPOLIS FROM %s WHERE NOOFFER = :1 FETCH FIRST 1 ROWS ONLY`, t), idMaster)
	return ada, err
}

// DaftarPolis = GetDataPolisNonProp_SQL (`NOOFFER = IDMaster OR NOOFFER = 7 karakter pertama IDMaster`; alias CARI
// diluruskan menurut SQL - OQ-CNP-24 bawaan "ikut SQL").
func (a *Acuan) DaftarPolis(ctx context.Context, idMaster string) ([]models.BarisPolis, error) {
	t, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT IDPEGA, NOPOLIS, BUSINESSCODE, INSUREDNAME,
		TO_CHAR(BEGINDATE, 'YYYY-MM-DD'), TO_CHAR(ENDDATE, 'YYYY-MM-DD'), CEDINGCO, SOB, TREATYGROUP
		FROM %s WHERE NOOFFER = :1 OR NOOFFER = :2 ORDER BY NOPOLIS`, t), 9, idMaster, models.AwalanMaster(idMaster))
	if err != nil {
		return nil, err
	}
	out := []models.BarisPolis{}
	for _, r := range rows {
		out = append(out, models.BarisPolis{IDPega: r[0], PolicyNo: r[1], BusinessCode: r[2], InsuredName: r[3],
			BeginDate: r[4], EndDate: r[5], CedingCo: r[6], SOB: r[7], TreatyGroup: r[8]})
	}
	return out, nil
}

// sqlBerkasPolis - berkas NB / EDM Treaty In terbaru untuk satu nomor polis (pola Claim Prop).
func sqlBerkasPolis(gen, kerja string) string {
	return fmt.Sprintf(`SELECT g.ID, g.PRODKE FROM %s g JOIN %s w ON w.ID = g.ID
		WHERE g.NOPOLIS = :1
		ORDER BY g.PRODKE DESC
		FETCH FIRST 1 ROWS ONLY`, gen, kerja)
}

// BerkasPolis = tombol View: modul dan ID berkas polis (bawaan OQ-CNP-13).
func (a *Acuan) BerkasPolis(ctx context.Context, nopolis string) (models.BerkasPolis, bool, error) {
	gen, err := a.q("T_GENERAL_POLIS_TREATY")
	if err != nil {
		return models.BerkasPolis{}, false, err
	}
	kerja, err := a.q("T_WORK_POLIS")
	if err != nil {
		return models.BerkasPolis{}, false, err
	}
	rows, err := a.banyak(ctx, sqlBerkasPolis(gen, kerja), 2, nopolis)
	if err != nil || len(rows) == 0 {
		return models.BerkasPolis{}, false, err
	}
	prodke, err := strconv.Atoi(strings.TrimSpace(rows[0][1]))
	if err != nil {
		return models.BerkasPolis{}, false, fmt.Errorf("repository: PRODKE berkas polis %q: %w", rows[0][0], err)
	}
	return models.BerkasPolis{Modul: models.ModulBerkasPolis(prodke), Kasus: rows[0][0]}, true, nil
}

// NamaTreatySpreading = GetTreatyName_SQL (PROPORTIONALARRG `TREATYDESCNAME = 'TREATY LIMIT'` menurut BusinessCode dan
// tanggal mulai master `Track.CARI33`), urut NOURUT jenis reasuransi. Mengembalikan REINSTYPENAME (CARI2).
func (a *Acuan) NamaTreatySpreading(ctx context.Context, bizCode, tanggalMulai string) ([]string, error) {
	pa, err := a.q("PROPORTIONALARRG")
	if err != nil {
		return nil, err
	}
	rt, err := a.q("REINSURANCETYPE")
	if err != nil {
		return nil, err
	}
	tb, err := a.q("TREATYBUSINESS")
	if err != nil {
		return nil, err
	}
	tc, err := a.q("TREATYCONTRACT")
	if err != nil {
		return nil, err
	}
	tgl := strings.ReplaceAll(tanggalMulai, "-", "")
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT a.REINSTYPENAME, (SELECT c.NOURUT FROM %s c WHERE c.ID = a.REINSTYPEID) URUT
		FROM %s a
		WHERE a.TREATYDESCNAME = 'TREATY LIMIT'
		  AND a.TREATYGROUPID IN (SELECT TREATYGROUPID FROM %s WHERE BIZCODE = :1 AND ISACTIVE = '1')
		  AND a.REINSTYPEID IN (SELECT REINSTYPEID FROM %s WHERE BIZCODE = :2 AND ISACTIVE = '1')
		  AND a.TREATYYEARID IN (SELECT IDTREATYYEAR FROM %s WHERE TO_DATE(:3, 'YYYYMMDD') BETWEEN TREATYSTARTDATE AND TREATYENDDATE)
		  AND a.REINSTYPEID IN (SELECT REINSTYPEID FROM %s WHERE TO_DATE(:4, 'YYYYMMDD') BETWEEN TREATYSTARTDATE AND TREATYENDDATE)
		ORDER BY URUT ASC`, rt, pa, tb, tb, tc, tc), 2, bizCode, bizCode, tgl, tgl)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, r := range rows {
		out = append(out, r[0])
	}
	return out, nil
}

// ---------------------------------------------------------------- lokasi, pelapor, adjuster

// Wilayah = BrowseRW_SQL (RW + CITY menurut ZIPCODE; alias diluruskan).
func (a *Acuan) Wilayah(ctx context.Context, kodePos string) ([]models.BarisWilayah, error) {
	rw, err := a.q("RW")
	if err != nil {
		return nil, err
	}
	city, err := a.q("CITY")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT b.ID, b.NOTE, b.DISTRICTID, b.DISTRICTNAME, c.ID, c.NOTE, c.PROVINCEID,
		c.PROVINCENAME FROM %s b LEFT JOIN %s c ON c.ID = b.CITYID WHERE b.ZIPCODE = :1`, rw, city), 8, kodePos)
	if err != nil {
		return nil, err
	}
	var out []models.BarisWilayah
	for _, r := range rows {
		out = append(out, models.BarisWilayah{RWID: r[0], RW: r[1], DistrictID: r[2], District: r[3], CityID: r[4],
			City: r[5], ProvinceID: r[6], Province: r[7]})
	}
	return out, nil
}

// AlamatAgen = GetAddressCeding (empat bagian alamat klien milik agen). XML membaca `M_CLIENT.JSONDATA.AddressList`;
// mengikuti keputusan work owner Claim Prop 08-10-2026 ("ubah jangan dari json, ambil dari client address") - dibaca
// dari tabel datar CLIENT_ADDRESS (PARITAS `[penyimpangan sadar]`, pola Claim Prop `AlamatKlien`).
func (a *Acuan) AlamatAgen(ctx context.Context, agentID string) ([4]string, error) {
	var out [4]string
	if strings.TrimSpace(agentID) == "" {
		return out, nil
	}
	ca, err := a.q("CLIENT_ADDRESS")
	if err != nil {
		return out, err
	}
	ag, err := a.q("AGENT")
	if err != nil {
		return out, err
	}
	rows, err := a.banyak(ctx, sqlAlamatKlien(ca, ag), 4, agentID)
	if err != nil || len(rows) == 0 {
		return out, err
	}
	copy(out[:], rows[0])
	return out, nil
}

// sqlAlamatKlien - satu baris CLIENT_ADDRESS klien milik agen: tipe 2 (Kantor) dulu, tipe 7 (Email) tidak pernah.
func sqlAlamatKlien(clientAddress, agent string) string {
	return fmt.Sprintf(`SELECT ASMADDRESS, RWNAME, DISTRICTNAME, CITYNAME FROM (
		SELECT a.ASMADDRESS, a.RWNAME, a.DISTRICTNAME, a.CITYNAME,
			ROW_NUMBER() OVER (ORDER BY CASE WHEN a.ASMADDRESSTYPE = '2' THEN 0 ELSE 1 END, a.PXCREATEDATETIME NULLS LAST, a.ROWID) URUT
		FROM %s a WHERE a.CLIENTID = (SELECT MAX(CLIENTID) FROM %s WHERE ID = :1) AND NVL(a.ASMADDRESSTYPE, '-') <> '7')
		WHERE URUT = 1`, clientAddress, agent)
}

// NamaAdjuster = BrowseAdjusterConsultant `.NAME` where `.ID`.
func (a *Acuan) NamaAdjuster(ctx context.Context, id string) (string, error) {
	t, err := a.q("ADJUSTERCONSULTANT")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT NAME FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY`, t), id)
	return v, err
}

// DaftarAdjuster = autocomplete Adjuster / Consultant ID (RD BrowseAdjusterConsultant: nilai `.ID`, tampil `.NAME`;
// 500).
func (a *Acuan) DaftarAdjuster(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("ADJUSTERCONSULTANT")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, NAME FROM %s WHERE UPPER(ID) LIKE :1 OR UPPER(NAME) LIKE :2
		ORDER BY ID FETCH FIRST 500 ROWS ONLY`, t), 2, pola, pola)
	return pilihan(rows, err)
}

// MarketingPolis = GetMObyNopol_SQL (MARKETINGOFFICER MOSTATUS '1' yang dipakai TREATYINPRODUCTION polis itu).
func (a *Acuan) MarketingPolis(ctx context.Context, nopolis string) ([6]string, bool, error) {
	var out [6]string
	mo, err := a.q("MARKETINGOFFICER")
	if err != nil {
		return out, false, err
	}
	tp, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return out, false, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT a.ID, a.CLIENTID, a.CLIENTNAME, a.TEAMGROUP, a.BRANCHDETAILID,
		a.BRANCHDETAILNAME FROM %s a WHERE a.MOSTATUS = '1'
		AND EXISTS (SELECT 1 FROM %s b WHERE b.MARKETINGOFFICERCODE = a.ID AND b.NOPOLIS = :1)
		FETCH FIRST 1 ROWS ONLY`, mo, tp), 6, nopolis)
	if err != nil || len(rows) == 0 {
		return out, false, err
	}
	copy(out[:], rows[0])
	return out, true, nil
}

// RiwayatKlaimPolis = CekHistoryClaimNonProp_SQL (JSON_KLAIM.DATA_JSON IDMaster / QuotationData.BusinessName /
// DateOfLoss menurut NOPOLIS, tidak di CLAIMREJECTED), DITAMBAH klaim NONPROP sistem baru (T_GENERAL_CLAIM) - JSON_KLAIM
// sistem baru tidak berisi DATA_JSON (keputusan work owner), tanpa tambahan ini klaim baru tak terlihat pembanding
// duplikat (pola Claim Prop).
func (a *Acuan) RiwayatKlaimPolis(ctx context.Context, nopolis string) ([]models.RiwayatKlaimPolis, error) {
	jk, err := a.q("JSON_KLAIM")
	if err != nil {
		return nil, err
	}
	rej, err := a.q("CLAIMREJECTED")
	if err != nil {
		return nil, err
	}
	gen, err := a.q("T_GENERAL_CLAIM")
	if err != nil {
		return nil, err
	}
	work, err := a.q("T_WORK_CLAIM")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT a.IDPEGA, JSON_VALUE(a.DATA_JSON, '$.IDMaster' RETURNING VARCHAR2(100)),
		JSON_VALUE(a.DATA_JSON, '$.QuotationData.BusinessName' RETURNING VARCHAR2(400)),
		JSON_VALUE(a.DATA_JSON, '$.DateOfLoss' RETURNING VARCHAR2(100))
		  FROM %s a WHERE a.NOPOLIS = :1 AND a.DATA_JSON IS NOT NULL
		   AND NOT EXISTS (SELECT 1 FROM %s r WHERE r.INSKEY = a.IDPEGA)
		UNION ALL
		SELECT g.ID, g.MASTER_ID, g.BUSINESS_NAME, TO_CHAR(g.DATE_OF_LOSS, 'YYYY-MM-DD')
		  FROM %s g JOIN %s w ON w.ID = g.ID AND w.LINI = :2 WHERE g.POLICY_NO = :3`, jk, rej, gen, work),
		4, nopolis, models.LiniNonProp, nopolis)
	if err != nil {
		return nil, err
	}
	var out []models.RiwayatKlaimPolis
	for _, r := range rows {
		dol := r[3]
		if t, ok := models.UraiTanggal(dol); ok {
			dol = models.TeksTanggal(t)
		}
		out = append(out, models.RiwayatKlaimPolis{KunciKasus: r[0], IDKasus: models.IDDariKunciPega(r[0]),
			IDMaster: r[1], Bisnis: r[2], DateOfLoss: dol})
	}
	return out, nil
}

// KodeLamaBisnis = GetDataBusiness_SQL (`OLDID FROM BUSINESS WHERE NOTE = nama bisnis`), baris pertama.
func (a *Acuan) KodeLamaBisnis(ctx context.Context, namaBisnis string) (string, error) {
	t, err := a.q("BUSINESS")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT OLDID FROM %s WHERE NOTE = :1 FETCH FIRST 1 ROWS ONLY`, t), namaBisnis)
	return v, err
}

// ---------------------------------------------------------------- komite dan wewenang

// RosterKomite = RD FilterEmailKomiteWithLimit (EMAILKOMITE `STS_KLAIM = 'NONPROP' AND STS_AKTIF = '1'`, urut DEGREE).
// Aturan tangga OQ-CNP-01 (ikut XML, `models.HanyaTingkat1`): `hanyaTingkat1` = hanya baris `LIMIT_BOTTOM <= 0`;
// selainnya semua baris aktif - kolom LIMIT tidak memilih tingkat.
func (a *Acuan) RosterKomite(ctx context.Context, hanyaTingkat1 bool) ([]models.AnggotaKomite, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return nil, err
	}
	saring := ""
	if hanyaTingkat1 {
		saring = " AND LIMIT_BOTTOM <= 0" // RD `LIMIT_BOTTOM <= Param.LIMIT_BOTTOM` (= 0): NULL tidak lolos, persis RD
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(ID), OPERATOR_ID, EMAIL, JABATAN, TO_CHAR(DEGREE) FROM %s
		WHERE STS_KLAIM = :1 AND STS_AKTIF = '1'%s ORDER BY DEGREE, ID`, t, saring), 5, models.STSKlaimNonProp)
	if err != nil {
		return nil, err
	}
	var out []models.AnggotaKomite
	for _, r := range rows {
		out = append(out, models.AnggotaKomite{ID: r[0], OperatorID: r[1], Email: r[2], Jabatan: r[3], Degree: r[4]})
	}
	return out, nil
}

// TingkatPelaku - label tingkat wewenang pelaku (JABATAN baris roster aktif STS_KLAIM NONPROP). OPERATOR_ID roster =
// akun pelaku ATAU workbasket aktif yang dipegangnya (roster ke workbasket, migrasi 6xx); beberapa baris = DEGREE
// terkecil (pola Claim Prop).
func (a *Acuan) TingkatPelaku(ctx context.Context, operatorID string) (string, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return "", err
	}
	lwb, err := a.q("M_LOGIN_GO_WORKBASKET")
	if err != nil {
		return "", err
	}
	wb, err := a.q("M_WORKBASKET")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, sqlTingkatPelaku(t, lwb, wb), models.STSKlaimNonProp, operatorID, operatorID)
	return v, err
}

// sqlTingkatPelaku - bind urut kemunculan: STS_KLAIM, akun (OPERATOR_ID), akun (pemegang workbasket).
func sqlTingkatPelaku(emk, lwb, wb string) string {
	return fmt.Sprintf(`SELECT e.JABATAN FROM %s e
		 WHERE e.STS_KLAIM = :1 AND e.STS_AKTIF = '1'
		   AND (UPPER(e.OPERATOR_ID) = UPPER(:2)
		        OR e.OPERATOR_ID IN (SELECT l.WORKBASKET_ID FROM %s l JOIN %s w ON w.WORKBASKET_ID = l.WORKBASKET_ID
		                              WHERE l.LOGIN_ID = :3 AND w.IS_ACTIVE = 1))
		 ORDER BY e.DEGREE, e.ID FETCH FIRST 1 ROWS ONLY`, emk, lwb, wb)
}

// EmailPelaku - `OperatorID.pyAddress` pelaku (ProteksiSendKomiteCNP_Act 19 / HitServiceToKasir_Act): M_LOGIN_GO.EMAIL.
func (a *Acuan) EmailPelaku(ctx context.Context, akun string) (string, error) {
	t, err := a.q("M_LOGIN_GO")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT EMAIL FROM %s WHERE LOGIN_ID = :1`, t), akun)
	return v, err
}

// NamaPelaku - nama tampilan akun (M_LOGIN_GO.NAME menurut LOGIN_ID); kosong = akun itu sendiri.
func (a *Acuan) NamaPelaku(ctx context.Context, akun string) (string, error) {
	t, err := a.q("M_LOGIN_GO")
	if err != nil {
		return akun, err
	}
	v, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT NAME FROM %s WHERE LOGIN_ID = :1`, t), akun)
	if err != nil {
		return "", err
	}
	if !ada || strings.TrimSpace(v) == "" {
		return akun, nil
	}
	return v, nil
}

// ---------------------------------------------------------------- rekening dan kasir

func (a *Acuan) rekening(ctx context.Context, where string, args ...any) ([]models.RekeningBank, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT CLIENTNAME, NAMEOFBANK, BRANCHOFBANK, ACCOUNTNO, SWIFTCODE, IDOFBANK,
		CURRENCYID, CURRENCY FROM %s WHERE %s ORDER BY ID`, t, where), 8, args...)
	if err != nil {
		return nil, err
	}
	out := []models.RekeningBank{}
	for _, r := range rows {
		out = append(out, models.RekeningBank{ClientName: r[0], NameOfBank: r[1], BranchOfBank: r[2], AccountNo: r[3],
			SwiftCode: r[4], IDOfBank: r[5], CurrencyID: r[6], Currency: r[7]})
	}
	return out, nil
}

// RekeningBank = GetDataBankAccount_sql (`CLIENTID AND CURRENCYID AND STS_CODE = '1'`).
func (a *Acuan) RekeningBank(ctx context.Context, klien, cur string) ([]models.RekeningBank, error) {
	return a.rekening(ctx, "CLIENTID = :1 AND CURRENCYID = :2 AND STS_CODE = '1'", klien, cur)
}

// RekeningBankKlien = GetDatabyClientName (`CLIENTID AND STS_CODE = '1'`).
func (a *Acuan) RekeningBankKlien(ctx context.Context, klien string) ([]models.RekeningBank, error) {
	return a.rekening(ctx, "CLIENTID = :1 AND STS_CODE = '1'", klien)
}

// RekeningBankNama = GetDatabyClientName2 (`CLIENTNAME = PayableTo AND CURRENCYID`) / GetDatabyClientName3
// (`CLIENTNAME = PayableTo`) - tanpa saringan STS_CODE, persis SQL.
func (a *Acuan) RekeningBankNama(ctx context.Context, nama, cur string) ([]models.RekeningBank, error) {
	if cur == "" {
		return a.rekening(ctx, "CLIENTNAME = :1", nama)
	}
	return a.rekening(ctx, "CLIENTNAME = :1 AND CURRENCYID = :2", nama, cur)
}

// DaftarKlien = GetClientName (`DISTINCT CLIENTNAME FROM BANKACCOUNT`) - dropdown Specify bila Payable To = 3; saringan
// mengandung dan batas 500 di server.
func (a *Acuan) DaftarKlien(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT CLIENTNAME AS NILAI, CLIENTNAME AS LABEL FROM %s
		WHERE CLIENTNAME IS NOT NULL AND UPPER(CLIENTNAME) LIKE :1 ORDER BY NILAI FETCH FIRST 500 ROWS ONLY`, t), 2, pola))
}

// IDBankRekening = HitServiceToKasir_Act langkah 11-12 (Obj-Browse BANKACCOUNT `.IDOFBANK` menurut nama bank, cabang,
// nomor rekening).
func (a *Acuan) IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT IDOFBANK FROM %s WHERE NAMEOFBANK = :1 AND BRANCHOFBANK = :2
		AND ACCOUNTNO = :3 FETCH FIRST 1 ROWS ONLY`, t), bank, cabang, akun)
	return v, err
}

// StatusKasir = GetStatusKasir (DIRECTTOKASIR_LOG.KET menurut NOAKSEPTASI, terbaru).
func (a *Acuan) StatusKasir(ctx context.Context, noAksep string) (string, bool, error) {
	t, err := a.q("DIRECTTOKASIR_LOG")
	if err != nil {
		return "", false, err
	}
	return a.satu(ctx, fmt.Sprintf(`SELECT KET FROM %s WHERE NOAKSEPTASI = :1 ORDER BY TGL_INPUT DESC NULLS LAST
		FETCH FIRST 1 ROWS ONLY`, t), noAksep)
}

// StatusKonversi = getStatusKonversi_SQL (`COUNT(1) FROM REINSURANCE.TRLOSS_DETAIL_T WHERE NO_AKSEP`) - HANYA di
// produksi; di luar produksi kosong (HitServiceToKasir_Act langkah 3 keluar). Tidak terbaca dari akun DEV (OQ-CNP-32).
func (a *Acuan) StatusKonversi(ctx context.Context, noAksep string) (string, error) {
	if !a.Produksi || noAksep == "" {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT TO_CHAR(COUNT(1)) FROM REINSURANCE.TRLOSS_DETAIL_T WHERE NO_AKSEP = :1`, noAksep)
	if err != nil {
		return "", err
	}
	if angkaBulat(v) > 0 {
		return "1", nil
	}
	return v, nil
}

// EmailCeding = GetEmailCeding_SQL (`GL.F_GET_EMAIL(:ceding) FROM DUAL`) - HANYA di produksi (OQ-CNP-32).
func (a *Acuan) EmailCeding(ctx context.Context, ceding string) (string, error) {
	if !a.Produksi {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT GL.F_GET_EMAIL(:1) FROM DUAL`, ceding)
	return v, err
}

// ---------------------------------------------------------------- popup data

// BarisRiwayatMaster - satu baris CLAIMXOL2 (GetHistoryMasterID).
type BarisRiwayatMaster = models.BarisXOL2

// RiwayatMaster = GetHistoryMasterID (`CLAIMXOL2 WHERE CASEID IN (OS_AKSEPTASI_KLAIM MASTERID = :id | :id/R01 |
// :id/R02) AND TANGGAL = MAX(CLAIMXOL.TANGGAL) kasus itu`). ⚠️ View CLAIMXOL INVALID di DEV (OQ-CNP-28) - kueri gagal di
// DEV; teks SQL persis rule.
func (a *Acuan) RiwayatMaster(ctx context.Context, id [3]string) ([]models.BarisXOL2, error) {
	x2, err := a.q("CLAIMXOL2")
	if err != nil {
		return nil, err
	}
	x1, err := a.q("CLAIMXOL")
	if err != nil {
		return nil, err
	}
	os, err := a.q("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT a.CASEID, a."XOL", a."Currency", a."GrossAdjustment", a."CNPReinstatement",
		a."KursIDR" FROM %s a WHERE a.CASEID IN (SELECT CASEID FROM %s WHERE MASTERID = :1 OR MASTERID = :2 OR MASTERID = :3)
		AND (SELECT MAX(b.TANGGAL) FROM %s b WHERE b.CASEID = a.CASEID) = a.TANGGAL ORDER BY a.CASEID`, x2, os, x1),
		6, id[0], id[1], id[2])
	if err != nil {
		return nil, err
	}
	var out []models.BarisXOL2
	for _, r := range rows {
		out = append(out, models.BarisXOL2{CaseID: r[0], XOL: r[1], Currency: r[2], GrossAdjustment: r[3],
			CNPReinstatement: r[4], KursIDR: r[5]})
	}
	return out, nil
}

// BarisLampiranBayar - satu baris grid ViewAttachmentNP.
type BarisLampiranBayar struct {
	Kategori     string `json:"kategori"`
	NamaFile     string `json:"namaFile"`
	NoAksep      string `json:"noAksep"`
	NoPrekas     string `json:"noPrekas"`
	TanggalBayar string `json:"tanggal"`
}

// LampiranBayar - grid ViewAttachmentNP (lampiran invoice pembayaran). Report Definition GCNMGetInvoiceAttachments
// membaca Link-Attachment kategori Invoice kasus; langkah pembuat Link-Attachment-nya ber-REMARK
// (GetPayAttachmentAdj_Act 3.3.3-3.3.6) dan yang hidup menulis tabel warisan DOCUMENT_CLAIM (3.3.10
// InsertDocument_Act) - maka dibaca dari tabel warisan DOCUMENT_CLAIM (KATEGORI_1 Invoice, NOAKSEP terisi;
// `[penyimpangan sadar]` PARITAS). Kunci IDPEGA kedua bentuk.
func (a *Acuan) LampiranBayar(ctx context.Context, kasusID string) ([]BarisLampiranBayar, error) {
	t, err := a.q("DOCUMENT_CLAIM") // tabel warisan POOLDATA (14 kolom kelas Pega), bukan T_CLAIMLF_DOCUMENT
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT KATEGORI_1, NAMAFILE, NOAKSEP, NOPREKAS,
		TO_CHAR(NVL(PAYMENTDATE, TANGGAL), 'YYYY-MM-DD HH24:MI:SS') FROM %s
		WHERE IDPEGA IN (:1, :2) AND KATEGORI_1 = 'Invoice' AND NOAKSEP IS NOT NULL ORDER BY TANGGAL, ID`, t),
		5, models.KunciInstans(kasusID), models.KunciPegaLama(kasusID))
	if err != nil {
		return nil, err
	}
	out := []BarisLampiranBayar{}
	for _, r := range rows {
		out = append(out, BarisLampiranBayar{Kategori: r[0], NamaFile: r[1], NoAksep: r[2], NoPrekas: r[3], TanggalBayar: r[4]})
	}
	return out, nil
}

// ---------------------------------------------------------------- pilihan layar

func pilihan(rows [][]string, err error) ([]models.Pilihan, error) {
	if err != nil {
		return nil, err
	}
	out := []models.Pilihan{}
	for _, r := range rows {
		out = append(out, models.Pilihan{Nilai: r[0], Label: r[1]})
	}
	return out, nil
}

// DaftarMataUang = BrowseCurrency_RD (urut `.Currency`).
func (a *Acuan) DaftarMataUang(ctx context.Context) ([]models.Pilihan, error) {
	t, err := a.q("CURRENCY")
	if err != nil {
		return nil, err
	}
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT ID, CURRENCY FROM %s ORDER BY CURRENCY`, t), 2))
}

// DaftarProvinsi = BrowseProvince_RD (distinct PROVINCENAME, `Nation = "INDONESIA"`, `PROVINCENAME Contains`).
func (a *Acuan) DaftarProvinsi(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("RW")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT PROVINCENAME AS NILAI, PROVINCENAME AS LABEL FROM %s
		WHERE UPPER(NATION) = 'INDONESIA' AND UPPER(PROVINCENAME) LIKE :1 ORDER BY NILAI FETCH FIRST 500 ROWS ONLY`, t),
		2, pola))
}

// BarisSebab - satu baris grid "List Cause of Loss".
type BarisSebab struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// DaftarSebab = grid "List Cause of Loss" (CauseofLoss_Section, RD BrowseCouseOfLoss_Business; `DESCRIPTION Contains`;
// urut DESCRIPTION; 500).
func (a *Acuan) DaftarSebab(ctx context.Context, cari string) ([]BarisSebab, error) {
	t, err := a.q("V_D_CAUSE_OF_LOSS_BUSINESS")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT D_COL_ID, DESCRIPTION FROM %s WHERE UPPER(DESCRIPTION) LIKE :1
		ORDER BY DESCRIPTION FETCH FIRST 500 ROWS ONLY`, t), 2, pola)
	if err != nil {
		return nil, err
	}
	out := []BarisSebab{}
	for _, r := range rows {
		out = append(out, BarisSebab{ID: r[0], Description: r[1]})
	}
	return out, nil
}

// BarisSebabID - satu baris cause of loss menurut D_COL_ID.
func (a *Acuan) BarisSebabID(ctx context.Context, id string) (BarisSebab, bool, error) {
	t, err := a.q("V_D_CAUSE_OF_LOSS_BUSINESS")
	if err != nil {
		return BarisSebab{}, false, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT D_COL_ID, DESCRIPTION FROM %s WHERE D_COL_ID = :1 FETCH FIRST 1 ROWS ONLY`, t), 2, id)
	if err != nil || len(rows) == 0 {
		return BarisSebab{}, false, err
	}
	return BarisSebab{ID: rows[0][0], Description: rows[0][1]}, true, nil
}

// BarisKatastrofe - satu baris grid katastrofe (RD GetCatastrope_RD).
type BarisKatastrofe struct {
	ID                string `json:"id"`
	StsKatastrofe     string `json:"stsKatastrofe"`
	NonKatastrofeType string `json:"nonKatastrofeType"`
	Note              string `json:"note"`
}

// DaftarKatastrofe = RD GetCatastrope_RD (`KLAIMTYPE = "NON-LIFE"`).
func (a *Acuan) DaftarKatastrofe(ctx context.Context, cari string) ([]BarisKatastrofe, error) {
	t, err := a.q("CATASTROPHE")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, STSKATASTROFE, NONKATASTROFETYPE, NOTE FROM %s
		WHERE KLAIMTYPE = 'NON-LIFE' AND UPPER(NVL(NOTE, ' ')) LIKE :1 ORDER BY TGL_INPUT DESC NULLS LAST
		FETCH FIRST 500 ROWS ONLY`, t), 4, pola)
	if err != nil {
		return nil, err
	}
	out := []BarisKatastrofe{}
	for _, r := range rows {
		out = append(out, BarisKatastrofe{ID: r[0], StsKatastrofe: r[1], NonKatastrofeType: r[2], Note: r[3]})
	}
	return out, nil
}

// BarisKatastrofeID - satu baris katastrofe menurut ID.
func (a *Acuan) BarisKatastrofeID(ctx context.Context, id string) (BarisKatastrofe, bool, error) {
	t, err := a.q("CATASTROPHE")
	if err != nil {
		return BarisKatastrofe{}, false, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, STSKATASTROFE, NONKATASTROFETYPE, NOTE FROM %s WHERE ID = :1`, t), 4, id)
	if err != nil || len(rows) == 0 {
		return BarisKatastrofe{}, false, err
	}
	r := rows[0]
	return BarisKatastrofe{ID: r[0], StsKatastrofe: r[1], NonKatastrofeType: r[2], Note: r[3]}, true, nil
}
