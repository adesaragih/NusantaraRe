package repository

// Untuk apa berkas ini: ACUAN - bacaan baca-saja yang dibutuhkan port activity (`models.Acuan`) dan pemilih layar.
// Setiap kueri meniru satu RDB-List / Report Definition korpus (nama rule di komentar); alias berbohong diluruskan.
//
// ⚠️ Dokumen JSON warisan (M_TREATY_IN.JSONDATA, JSON_POLIS.DATA_JSON, M_CLIENT.JSONDATA, JSON_KLAIM.DATA_JSON,
// OS_AKSEPTASI_KLAIM.DATA_JSON) DIBACA lewat JSON_VALUE / JSON_TABLE persis seperti rule-nya - baca-saja; modul ini
// tidak pernah MENULIS JSON (keputusan work owner).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
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
			s[i] = v[i].String
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func dec(kol string) string { return fmt.Sprintf(db.FmtDesimal, kol) }

// KursStandar = CurrencyStandard (`POOLDATA.GETCURRENCYSTANDARD(:cur, SYSDATE)` - fungsi, bukan procedure).
func (a *Acuan) KursStandar(ctx context.Context, cur string) (string, error) {
	f, err := a.q("GETCURRENCYSTANDARD")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT %s FROM DUAL`, dec(f+"(:1, SYSDATE)")), cur)
	return v, err
}

// NamaMataUang = BrowseCurrency_RD `.Currency` where `.ID`.
func (a *Acuan) NamaMataUang(ctx context.Context, cur string) (string, error) {
	t, err := a.q("CURRENCY")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT CURRENCY FROM %s WHERE ID = :1`, t), cur)
	return v, err
}

// NamaJenisReasuransi = BrowseReinsuranceType_RD `.Note` where `.ID`.
func (a *Acuan) NamaJenisReasuransi(ctx context.Context, id string) (string, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT NOTE FROM %s WHERE ID = :1`, t), id)
	return v, err
}

// IDJenisReasuransi = BrowseReinsuranceType_RD `.ID` where `UPPER(.Note) = UPPER(:nama)` and `.Type` (baris pertama
// urut ID).
func (a *Acuan) IDJenisReasuransi(ctx context.Context, nama, tipe string) (string, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT ID FROM %s WHERE UPPER(NOTE) = UPPER(:1) AND TYPE = :2
		ORDER BY ID FETCH FIRST 1 ROWS ONLY`, t), nama, tipe)
	return v, err
}

// ---------------------------------------------------------------- master treaty

// MasterTreaty = GetLimitsTreatyIn_SQL (`SELECT JSONDATA FROM M_TREATY_IN WHERE ID UNION ALL ... M_TREATY_IN_EDM`) +
// adoptJSONObject: baris pertama yang ada (M_TREATY_IN lebih dulu).
func (a *Acuan) MasterTreaty(ctx context.Context, id string) (models.MasterTreaty, bool, error) {
	for _, tabel := range []string{"M_TREATY_IN", "M_TREATY_IN_EDM"} {
		m, ada, err := a.bacaMaster(ctx, tabel, id)
		if err != nil || ada {
			return m, ada, err
		}
	}
	return models.MasterTreaty{}, false, nil
}

var medanMaster = []string{"TreatyContractName", "ProportionType", "Ceding", "CedingID", "LeadingReinsSource",
	"LeadingReinsSourceID", "Bordeaux", "BordereauxNote", "AccountingMode", "TeritorialScope", "Commencement",
	"Termination", "TreatyYear", "RNMShareP", "StatusAkseptasi"}

func (a *Acuan) bacaMaster(ctx context.Context, objek, id string) (models.MasterTreaty, bool, error) {
	t, err := a.q(objek)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	var kol []string
	for _, m := range medanMaster {
		kol = append(kol, fmt.Sprintf(`JSON_VALUE(JSONDATA, '$.%s' RETURNING VARCHAR2(4000))`, m))
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, %s FROM %s WHERE ID = :1`, strings.Join(kol, ", "), t),
		len(medanMaster)+1, id)
	if err != nil || len(rows) == 0 {
		return models.MasterTreaty{}, false, err
	}
	r := rows[0]
	m := models.MasterTreaty{ID: r[0], TreatyContractName: r[1], ProportionType: r[2], Ceding: r[3], CedingID: r[4],
		LeadingReinsSource: r[5], LeadingReinsSourceID: r[6], Bordeaux: r[7], BordereauxNote: r[8],
		AccountingMode: r[9], TeritorialScope: r[10], Commencement: tglPega(r[11]), Termination: tglPega(r[12]),
		TreatyYear: r[13], RNMShareP: r[14], StatusAkseptasi: r[15]}
	lim, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.li, jt.tt, jt.di, jt.tg, jt.rs FROM %s m,
		JSON_TABLE(m.JSONDATA, '$.Limits[*]' COLUMNS (li FOR ORDINALITY, tt VARCHAR2(200) PATH '$.TreatyType',
		  NESTED PATH '$.Detail[*]' COLUMNS (di FOR ORDINALITY, tg VARCHAR2(100) PATH '$.TreatyGroupID',
		  rs VARCHAR2(100) PATH '$.RNMShare'))) jt WHERE m.ID = :1 ORDER BY jt.li, jt.di`, t), 5, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	idx := map[string]int{}
	for _, l := range lim {
		li, ok := idx[l[0]]
		if !ok {
			m.Limits = append(m.Limits, models.LimitMaster{TreatyType: l[1]})
			li = len(m.Limits) - 1
			idx[l[0]] = li
		}
		if l[2] != "" {
			m.Limits[li].Detail = append(m.Limits[li].Detail, models.DetailLimit{TreatyGroupID: l[3], RNMShare: l[4]})
		}
	}
	cash, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.li, jt.di, jt.cur, jt.val FROM %s m,
		JSON_TABLE(m.JSONDATA, '$.Limits[*]' COLUMNS (li FOR ORDINALITY, NESTED PATH '$.Detail[*]' COLUMNS
		  (di FOR ORDINALITY, NESTED PATH '$.CashLossList[*]' COLUMNS (cur VARCHAR2(100) PATH '$.Currency',
		  val VARCHAR2(100) PATH '$.Value')))) jt WHERE m.ID = :1 AND (jt.cur IS NOT NULL OR jt.val IS NOT NULL)
		  ORDER BY jt.li, jt.di`, t), 4, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, c := range cash {
		if d := detail(&m, idx, c[0], c[1]); d != nil {
			d.CashLossList = append(d.CashLossList, models.CashLoss{Currency: c[2], Value: c[3]})
		}
	}
	spr, err := a.banyak(ctx, fmt.Sprintf(`SELECT jt.li, jt.di, jt.ri, jt.rn, jt.pc FROM %s m,
		JSON_TABLE(m.JSONDATA, '$.Limits[*]' COLUMNS (li FOR ORDINALITY, NESTED PATH '$.Detail[*]' COLUMNS
		  (di FOR ORDINALITY, NESTED PATH '$.SpreadingList[*]' COLUMNS (ri VARCHAR2(100) PATH '$.ReinsTypeID',
		  rn VARCHAR2(400) PATH '$.ReinsTypeName', pc VARCHAR2(100) PATH '$.Pct')))) jt
		 WHERE m.ID = :1 AND (jt.ri IS NOT NULL OR jt.rn IS NOT NULL) ORDER BY jt.li, jt.di`, t), 5, id)
	if err != nil {
		return models.MasterTreaty{}, false, err
	}
	for _, s := range spr {
		if d := detail(&m, idx, s[0], s[1]); d != nil {
			d.SpreadingList = append(d.SpreadingList, models.SpreadingMaster{ReinsTypeID: s[2], ReinsTypeName: s[3], Pct: s[4]})
		}
	}
	return m, true, nil
}

func detail(m *models.MasterTreaty, idx map[string]int, li, di string) *models.DetailLimit {
	l, ok := idx[li]
	if !ok {
		return nil
	}
	n := angkaBulat(di)
	if n < 1 || n > len(m.Limits[l].Detail) {
		return nil
	}
	return &m.Limits[l].Detail[n-1]
}

// tglPega - "yyyyMMdd" master -> "2006-01-02".
func tglPega(s string) string {
	if t, ok := models.TanggalPega(s); ok {
		return models.TeksTanggal(t)
	}
	return s
}

// sqlDaftarMaster - RD BrowseCLAIM_MASTER_TREATY: `.PROPORTIONTYPE = Param.TREATYTYPE` ("Proportional", parameter
// Section MasterTreatyInList), filter per kolom digabung AND (mengandung, tanpa beda huruf), terbaru dulu, batas
// `models.BatasMaster` (keputusan work owner 08-10-2026).
func sqlDaftarMaster(t string, s models.SaringanMaster) (string, []any) {
	args := []any{models.ProporsiMaster}
	w := "WHERE PROPORTIONTYPE = :1"
	for _, f := range []struct{ kol, nilai string }{
		{"TREATYID", s.TreatyID}, {"CLASSOFBUSINESS", s.ClassOfBusiness}, {"TREATYCONTRACTNAME", s.ContractName},
		{"SOB", s.SOB}, {"CEDING", s.InsuredName}, {"TREATYTYPE", s.TreatyType}, {"TREATYGROUP", s.TreatyGroup},
		{"TREATYYEAR", s.TreatyYear},
	} {
		if v := strings.TrimSpace(f.nilai); v != "" {
			args = append(args, "%"+strings.ToUpper(v)+"%")
			w += fmt.Sprintf(" AND UPPER(%s) LIKE :%d", f.kol, len(args))
		}
	}
	return fmt.Sprintf(`SELECT DISTINCT TREATYID, CLASSOFBUSINESS, CLASSOFBUSINESSID, TREATYCONTRACTNAME,
		SOB, CEDING, TREATYTYPE, PROPORTIONTYPE, TREATYGROUP, TREATYGROUPID, TREATYYEAR FROM %s
		%s
		ORDER BY TREATYYEAR DESC, TREATYID DESC
		FETCH FIRST %d ROWS ONLY`, t, w, models.BatasMaster), args
}

// DaftarMaster = grid "Data Master TreatyIn" (RD BrowseCLAIM_MASTER_TREATY, distinct; Section MasterTreatyInList
// mengirim TREATYTYPE = "Proportional" ke saringan `.PROPORTIONTYPE`). Saringan kolom dijalankan di server.
func (a *Acuan) DaftarMaster(ctx context.Context, s models.SaringanMaster) ([]models.BarisMaster, error) {
	t, err := a.q("CLAIM_MASTER_TREATY")
	if err != nil {
		return nil, err
	}
	q, args := sqlDaftarMaster(t, s)
	rows, err := a.banyak(ctx, q, 11, args...)
	if err != nil {
		return nil, err
	}
	out := []models.BarisMaster{}
	for _, r := range rows {
		out = append(out, models.BarisMaster{TreatyID: r[0], ClassOfBusiness: r[1], ClassOfBusinessID: r[2],
			TreatyContractName: r[3], SOB: r[4], Ceding: r[5], TreatyType: r[6], ProportionType: r[7], TreatyGroup: r[8],
			TreatyGroupID: r[9], TreatyYear: r[10]})
	}
	return out, nil
}

// BarisMasterDari - satu baris view master menurut TREATYID + TREATYGROUPID + CLASSOFBUSINESSID (parameter tombol
// Choose dibaca ULANG di server, bukan dari layar).
func (a *Acuan) BarisMasterDari(ctx context.Context, treatyID, grupID, cobID string) (models.BarisMaster, bool, error) {
	t, err := a.q("CLAIM_MASTER_TREATY")
	if err != nil {
		return models.BarisMaster{}, false, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT TREATYID, CLASSOFBUSINESS, CLASSOFBUSINESSID, TREATYCONTRACTNAME,
		SOB, CEDING, TREATYTYPE, PROPORTIONTYPE, TREATYGROUP, TREATYGROUPID, TREATYYEAR FROM %s
		WHERE TREATYID = :1 AND NVL(TREATYGROUPID, '-') = NVL(:2, '-') AND NVL(CLASSOFBUSINESSID, '-') = NVL(:3, '-')
		FETCH FIRST 1 ROWS ONLY`, t), 11, treatyID, teksAtauNil(grupID), teksAtauNil(cobID))
	if err != nil || len(rows) == 0 {
		return models.BarisMaster{}, false, err
	}
	r := rows[0]
	return models.BarisMaster{TreatyID: r[0], ClassOfBusiness: r[1], ClassOfBusinessID: r[2], TreatyContractName: r[3],
		SOB: r[4], Ceding: r[5], TreatyType: r[6], ProportionType: r[7], TreatyGroup: r[8], TreatyGroupID: r[9],
		TreatyYear: r[10]}, true, nil
}

// ---------------------------------------------------------------- polis

// AdaPolisMaster = GetNopolis_SQL (`SELECT DISTINCT NOPOLIS FROM TREATYINPRODUCTION WHERE NOOFFER = :id`).
func (a *Acuan) AdaPolisMaster(ctx context.Context, idMaster string) (bool, error) {
	t, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return false, err
	}
	_, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT NOPOLIS FROM %s WHERE NOOFFER = :1 FETCH FIRST 1 ROWS ONLY`, t), idMaster)
	return ada, err
}

// AdaPolisRealisasi = CekPolicyNumber_SQL (`SELECT IDPEGA FROM POOLDATA.TREATYINPRODUCTION WHERE NOPOLIS = :p`).
func (a *Acuan) AdaPolisRealisasi(ctx context.Context, nopolis string) (bool, error) {
	t, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return false, err
	}
	_, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT NOPOLIS FROM %s WHERE NOPOLIS = :1 FETCH FIRST 1 ROWS ONLY`, t), nopolis)
	return ada, err
}

// sqlBerkasPolis - berkas NB / EDM Treaty In terbaru untuk satu nomor polis: generasi PRODKE terbesar yang ada di
// T_WORK_POLIS (polis Pega lama yang belum disalin lewat Copy Old tidak punya berkas).
func sqlBerkasPolis(gen, kerja string) string {
	return fmt.Sprintf(`SELECT g.ID, g.PRODKE FROM %s g JOIN %s w ON w.ID = g.ID
		WHERE g.NOPOLIS = :1
		ORDER BY g.PRODKE DESC
		FETCH FIRST 1 ROWS ONLY`, gen, kerja)
}

// BerkasPolis = tombol View: modul dan ID berkas polis.
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

// DaftarPolis = grid "Data Polis" (RDB SetPolicyTreatyProp: `NOOFFER = MasterID.CARI1 AND TREATYGROUP =
// ClaimData.TreatyGroupName`; `'0' AS "Prodke"`).
func (a *Acuan) DaftarPolis(ctx context.Context, noOffer, grup string) ([]models.BarisPolis, error) {
	t, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT NOPOLIS, NOOFFER, TREATYGROUP, SOB, UW_YEAR, QUARTER FROM %s
		WHERE NOOFFER = :1 AND TREATYGROUP = :2 ORDER BY NOPOLIS`, t), 6, noOffer, grup)
	if err != nil {
		return nil, err
	}
	out := []models.BarisPolis{}
	for _, r := range rows {
		out = append(out, models.BarisPolis{PolicyNo: r[0], NoOffer: r[1], TreatyGroup: r[2], SourceOfBusinessName: r[3],
			TreatyYear: r[4], Prodke: "0", Quarter: r[5]})
	}
	return out, nil
}

// TreatyGroupBisnis = GetTreatyGroupID (`DISTINCT TREATYGROUPID FROM TREATYBUSINESS WHERE BIZCODE AND TREATYYEAR AND
// ISACTIVE = '1'`), baris pertama.
func (a *Acuan) TreatyGroupBisnis(ctx context.Context, bizCode, treatyYear string) (string, error) {
	t, err := a.q("TREATYBUSINESS")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT DISTINCT TREATYGROUPID FROM %s
		WHERE BIZCODE = :1 AND TREATYYEAR = :2 AND ISACTIVE = '1' FETCH FIRST 1 ROWS ONLY`, t), bizCode, treatyYear)
	return v, err
}

// YearOfQuartal = GetYearofQuartal (`a.DATA_JSON.YearOfQuartal FROM JSON_POLIS a WHERE NOPOLIS AND PRODKE`).
// ⚠️ Polis realisasi sistem baru (NB Treaty In) menulis JSON_POLIS tanpa DATA_JSON - nilainya kosong (OQ-CP-09).
func (a *Acuan) YearOfQuartal(ctx context.Context, nopolis, prodke string) (string, error) {
	t, err := a.q("JSON_POLIS")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT JSON_VALUE(a.DATA_JSON, '$.YearOfQuartal' RETURNING VARCHAR2(100))
		FROM %s a WHERE a.NOPOLIS = :1 AND a.PRODKE = :2 FETCH FIRST 1 ROWS ONLY`, t), nopolis, prodke)
	return v, err
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

// AlamatKlien = GetAddressCeding (empat bagian alamat klien milik agen `agentID`). XML membaca
// `M_CLIENT.JSONDATA.AddressList.{ASMAddress,RWName,DistrictName,CityName}`; [keputusan work owner 08-10-2026] "ubah
// jangan dari json, ambil dari client address" - dibaca dari tabel datar CLIENT_ADDRESS (ASMADDRESS, RWNAME,
// DISTRICTNAME, CITYNAME). Satu baris: tipe 2 (Kantor) dulu, tipe 7 (Email) tidak pernah, lalu PXCREATEDATETIME; ROWID
// hanya pemutus seri terakhir. `[data DEV 08-10-2026]` agen ceding master treaty: setiap bagian alamat yang terisi di
// kedua sumber SAMA dengan JSON lama; 2 agen ber-JSON kosong kini mendapat alamat Kantor.
func (a *Acuan) AlamatKlien(ctx context.Context, agentID string) ([4]string, bool, error) {
	var out [4]string
	ca, err := a.q("CLIENT_ADDRESS")
	if err != nil {
		return out, false, err
	}
	ag, err := a.q("AGENT")
	if err != nil {
		return out, false, err
	}
	rows, err := a.banyak(ctx, sqlAlamatKlien(ca, ag), 4, agentID)
	if err != nil || len(rows) == 0 {
		return out, false, err
	}
	copy(out[:], rows[0])
	return out, true, nil
}

// sqlAlamatKlien - satu baris CLIENT_ADDRESS klien milik agen (lihat AlamatKlien).
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

// DaftarAdjuster = autocomplete Consultant / Adjuster (RD BrowseAdjusterConsultant, kolom ID*, NAME; 500). Saringan dan
// urutan menurut NAME saja (work owner 08-10-2026: "yang dropdown hanya dari namanya aja"); ID tetap nilai pilihan.
func (a *Acuan) DaftarAdjuster(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("ADJUSTERCONSULTANT")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, NAME FROM %s WHERE UPPER(NAME) LIKE :1
		ORDER BY NAME, ID FETCH FIRST 500 ROWS ONLY`, t), 2, pola)
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

// RiwayatKlaimPolis = CariHistoryClaim_SQL (JSON_KLAIM.DATA_JSON.DateOfLoss menurut NOPOLIS) + RejectedClaim_RD
// (CLAIMREJECTED menurut INSKEY), DITAMBAH klaim PROP sistem baru (T_GENERAL_CLAIM.DATE_OF_LOSS) - JSON_KLAIM sistem
// baru tidak berisi DATA_JSON (keputusan work owner), tanpa tambahan ini klaim baru tak terlihat pembanding duplikat.
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
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT a.IDPEGA, JSON_VALUE(a.DATA_JSON, '$.DateOfLoss' RETURNING VARCHAR2(100)),
		CASE WHEN EXISTS (SELECT 1 FROM %s r WHERE r.INSKEY = a.IDPEGA) THEN '1' ELSE '0' END
		  FROM %s a WHERE a.NOPOLIS = :1 AND a.DATA_JSON IS NOT NULL
		UNION ALL
		SELECT g.ID, TO_CHAR(g.DATE_OF_LOSS, 'YYYY-MM-DD HH24:MI:SS'), '0'
		  FROM %s g JOIN %s w ON w.ID = g.ID AND w.LINI = :2 WHERE g.POLICY_NO = :3`, rej, jk, gen, work),
		3, nopolis, models.LiniProp, nopolis)
	if err != nil {
		return nil, err
	}
	var out []models.RiwayatKlaimPolis
	for _, r := range rows {
		dol := r[1]
		if t, ok := models.UraiTanggal(dol); ok {
			dol = t.Format("2006-01-02 15:04:05")
		}
		out = append(out, models.RiwayatKlaimPolis{KunciKasus: r[0], IDKasus: models.IDDariKunciPega(r[0]),
			DateOfLoss: dol, Ditolak: r[2] == "1"})
	}
	return out, nil
}

// KodeLamaBisnis = GetOldIDBusiness (`BUSINESS.OLDID WHERE ID = BusinessCode`).
func (a *Acuan) KodeLamaBisnis(ctx context.Context, bizCode string) (string, error) {
	t, err := a.q("BUSINESS")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT OLDID FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY`, t), bizCode)
	return v, err
}

// TahunTreaty = TreatyYearTreatyin_SQL (`TREATYYEAR WHERE TREATYGROUPID AND :ymd BETWEEN STARTDATE AND ENDDATE`;
// STARTDATE/ENDDATE teks "yyyyMMdd" - katalog DEV 07-10-2026).
func (a *Acuan) TahunTreaty(ctx context.Context, grup, ymd string) (string, error) {
	t, err := a.q("TREATYYEAR")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT TREATYYEAR FROM %s WHERE TREATYGROUPID = :1
		AND :2 BETWEEN STARTDATE AND ENDDATE FETCH FIRST 1 ROWS ONLY`, t), grup, ymd)
	return v, err
}

// AnakSpreading - anak PROPORTIONALARRG (PARENTREINSTYPEID) dalam treaty group klaim, pada tahun arrangement TERBARU
// yang <= tahun treaty klaim: cadangan tabel bawah spreading bila SpreadingList master kosong (keputusan work owner
// 09-10-2026; aturan commit 24a6e315 - `[data DEV 08-10-2026]` TREATYYEAR arrangement = tahun kontrak treaty). Satu
// baris per REINSTYPEID: TGLUPDATE terbaru. Urut SPREADINGORDER lalu REINSTYPEID.
func (a *Acuan) AnakSpreading(ctx context.Context, induk, tahun, grup string) ([]models.AnakSpreading, error) {
	t, err := a.q("PROPORTIONALARRG")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, sqlAnakSpreading(t), 3, induk, grup, induk, grup, tahun)
	if err != nil {
		return nil, err
	}
	out := []models.AnakSpreading{}
	for _, r := range rows {
		out = append(out, models.AnakSpreading{ReinsTypeID: r[0], ReinsTypeName: r[1], Pct: r[2]})
	}
	return out, nil
}

func sqlAnakSpreading(t string) string {
	return fmt.Sprintf(`SELECT REINSTYPEID, REINSTYPENAME, PCT FROM (
		SELECT REINSTYPEID, REINSTYPENAME, PCT, SPREADINGORDER,
			ROW_NUMBER() OVER (PARTITION BY REINSTYPEID ORDER BY TGLUPDATE DESC NULLS LAST, ID DESC) URUT
		FROM %s WHERE PARENTREINSTYPEID = :1 AND TREATYGROUPID = :2
		  AND TREATYYEAR = (SELECT MAX(TREATYYEAR) FROM %s WHERE PARENTREINSTYPEID = :3 AND TREATYGROUPID = :4
		    AND TREATYYEAR <= :5))
		WHERE URUT = 1 ORDER BY SPREADINGORDER NULLS LAST, REINSTYPEID`, t, t)
}

// LimitPLA = GetLimitPLATreatyin (PROPORTIONALARRG `TREATYDESCID = '10001'`): kolom RP (alias CARI7) baris pertama.
func (a *Acuan) LimitPLA(ctx context.Context, tahun, grup, reins string) (string, bool, error) {
	t, err := a.q("PROPORTIONALARRG")
	if err != nil {
		return "", false, err
	}
	return a.satu(ctx, fmt.Sprintf(`SELECT RP FROM %s WHERE TREATYDESCID = '10001' AND TREATYYEAR = :1
		AND TREATYGROUPID = :2 AND REINSTYPEID = :3 FETCH FIRST 1 ROWS ONLY`, t), tahun, grup, reins)
}

// SpreadingPolis - spreading polis di TREATYINPRODUCTION (keputusan work owner 08-10-2026; AddSpreading_Act tidak
// diekspor): JN_REAS, PCT_SHARE_PREMI, CURR_ID unik per NOPOLIS. CurrencyID = CURRENCY.ID terkecil berkode CURR_ID
// (`[data DEV 08-10-2026]` satu kode kembar, TWD).
func (a *Acuan) SpreadingPolis(ctx context.Context, nopolis string) ([]models.SpreadingPolis, error) {
	p, err := a.q("TREATYINPRODUCTION")
	if err != nil {
		return nil, err
	}
	c, err := a.q("CURRENCY")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT p.JN_REAS, %s, (SELECT MIN(c.ID) FROM %s c WHERE c.CURRENCY = p.CURR_ID),
		p.CURR_ID FROM %s p WHERE p.NOPOLIS = :1 AND p.JN_REAS IS NOT NULL
		GROUP BY p.JN_REAS, p.PCT_SHARE_PREMI, p.CURR_ID ORDER BY p.JN_REAS, p.CURR_ID`, dec("p.PCT_SHARE_PREMI"), c, p),
		4, nopolis)
	if err != nil {
		return nil, err
	}
	out := []models.SpreadingPolis{}
	for _, r := range rows {
		out = append(out, models.SpreadingPolis{TreatyType: r[0], SharePercentage: rapikanDesimal(r[1]),
			CurrencyID: r[2], Currency: r[3]})
	}
	return out, nil
}

// DaftarRetro = GetListRetro_Sql (TREATYREINSURER `REINSTYPEID AND TREATYYEAR AND TREATYGROUPID`).
func (a *Acuan) DaftarRetro(ctx context.Context, reins, tahun, grup string) ([]models.Retro, error) {
	t, err := a.q("TREATYREINSURER")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT REINSURERID, NAME, %s, %s FROM %s
		WHERE REINSTYPEID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3 ORDER BY ID`, dec("PCTSHARE"), dec("RICOMM"), t),
		4, reins, tahun, grup)
	if err != nil {
		return nil, err
	}
	var out []models.Retro
	for _, r := range rows {
		out = append(out, models.Retro{ReinsurerID: r[0], ReinsurerName: r[1], PctShare: rapikanDesimal(r[2]),
			RiComm: rapikanDesimal(r[3])})
	}
	return out, nil
}

// ---------------------------------------------------------------- komite dan wewenang

// RosterKomite = RD FilterEmailKomiteWithLimit (`LIMIT_BOTTOM <= :nilai AND STS_KLAIM = :sts AND STS_AKTIF = '1'`,
// urut DEGREE). Batas atas TIDAK menyaring (AC 60).
func (a *Acuan) RosterKomite(ctx context.Context, nilai, sts string) ([]models.AnggotaKomite, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return nil, err
	}
	v, err := angkaAtauNil("LIMIT_BOTTOM", nilai)
	if err != nil {
		return nil, err
	}
	if v[0] == nil {
		v = []any{"0", int64(0)}
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(ID), OPERATOR_ID, EMAIL, JABATAN, DEGREE FROM %s
		WHERE LIMIT_BOTTOM <= (TO_NUMBER(:1) / POWER(10, :2)) AND STS_KLAIM = :3 AND STS_AKTIF = '1'
		ORDER BY DEGREE, ID`, t), 5, v[0], v[1], sts)
	if err != nil {
		return nil, err
	}
	var out []models.AnggotaKomite
	for _, r := range rows {
		out = append(out, models.AnggotaKomite{ID: r[0], OperatorID: r[1], Email: r[2], Jabatan: r[3], Degree: r[4]})
	}
	return out, nil
}

// TingkatPelaku - label tingkat wewenang pelaku (JABATAN baris roster aktif STS_KLAIM PROP; AC 54, 61). OPERATOR_ID
// roster = akun pelaku ATAU workbasket aktif yang dipegangnya (tangga komite PROP ke workbasket, migrasi 537,
// keputusan work owner 09-10-2026); beberapa baris = DEGREE terkecil.
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
	v, _, err := a.satu(ctx, sqlTingkatPelaku(t, lwb, wb), models.STSKlaimProp, operatorID, operatorID)
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

// OperatorKomiteAktif - KomiteID (OPERATOR_ID) roster aktif STS_KLAIM `sts` - nama workbasket sejak migrasi 537. Tabel
// komite inbox Claim Prop tampil bagi pemegang salah satunya (keputusan work owner 09-10-2026).
func (a *Acuan) OperatorKomiteAktif(ctx context.Context, sts string) ([]string, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT OPERATOR_ID FROM %s
		WHERE STS_KLAIM = :1 AND STS_AKTIF = '1' AND OPERATOR_ID IS NOT NULL ORDER BY OPERATOR_ID`, t), 1, sts)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r[0])
	}
	return out, nil
}

// DegreeDirekturUtama - baris roster Direktur Utama. `[dugaan]` GetLimitDirekturUtama_SQL mencarinya menurut NAMA orang
// (dibuang, AC 55); di DEV satu-satunya baris DEGREE 6 (STS_KLAIM kosong) - kunci DEGREE menunggu penegasan DBA
// (OQ-CP-10).
const DegreeDirekturUtama = "6"

// LimitDirekturUtama = GetLimitDirekturUtama_SQL (LIMIT_BOTTOM baris Direktur Utama; AC 56).
func (a *Acuan) LimitDirekturUtama(ctx context.Context) (string, bool, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return "", false, err
	}
	return a.satu(ctx, fmt.Sprintf(`SELECT %s FROM %s WHERE DEGREE = :1 FETCH FIRST 1 ROWS ONLY`, dec("LIMIT_BOTTOM"), t),
		DegreeDirekturUtama)
}

// ---------------------------------------------------------------- rekening dan premi

func (a *Acuan) rekening(ctx context.Context, where string, args ...any) ([]models.RekeningBank, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT CLIENTNAME, NAMEOFBANK, BRANCHOFBANK, ACCOUNTNO, SWIFTCODE, IDOFBANK,
		CURRENCYID FROM %s WHERE %s AND STS_CODE = '1' ORDER BY ID`, t, where), 7, args...)
	if err != nil {
		return nil, err
	}
	out := []models.RekeningBank{}
	for _, r := range rows {
		out = append(out, models.RekeningBank{ClientName: r[0], NameOfBank: r[1], BranchOfBank: r[2], AccountNo: r[3],
			SwiftCode: r[4], IDOfBank: r[5], CurrencyID: r[6]})
	}
	return out, nil
}

// RekeningBank = GetDataBankAccount_sql (`CLIENTID AND CURRENCYID AND STS_CODE = '1'`).
func (a *Acuan) RekeningBank(ctx context.Context, klien, cur string) ([]models.RekeningBank, error) {
	return a.rekening(ctx, "CLIENTID = :1 AND CURRENCYID = :2", klien, cur)
}

// RekeningBankMataUang = GetDataBankAccount2_sql (`CURRENCYID AND STS_CODE = '1'`).
func (a *Acuan) RekeningBankMataUang(ctx context.Context, cur string) ([]models.RekeningBank, error) {
	return a.rekening(ctx, "CURRENCYID = :1", cur)
}

// RekeningBankKlien = GetDatabyClientName (`CLIENTID AND STS_CODE = '1'`).
func (a *Acuan) RekeningBankKlien(ctx context.Context, klien string) ([]models.RekeningBank, error) {
	return a.rekening(ctx, "CLIENTID = :1", klien)
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

// SaldoPremi = CekLunasPremi_Sql (`SUM(IVD_TRANS_SIGN * IVD_TOTAL)` INVOICE x DETAIL_INVOICE - sinonim publik).
func (a *Acuan) SaldoPremi(ctx context.Context, invoice, cur string) (string, error) {
	// Sinonim publik INVOICE / DETAIL_INVOICE menunjuk ARASAPAS (katalog DEV 07-10-2026) - skema ditulis eksplisit.
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT %s FROM ARASAPAS.INVOICE a, ARASAPAS.DETAIL_INVOICE b
		WHERE a.INV_INV_NO = b.INV_INV_NO AND a.CURR_ENTRY_NO = b.CURR_ENTRY_NO AND a.INV_INV_NO = :1 AND a.INV_LKU_ID = :2`,
		dec("SUM(IVD_TRANS_SIGN * IVD_TOTAL)")), invoice, cur)
	if err != nil {
		return saldoPremiGagal(a.Produksi, err)
	}
	return v, nil
}

// saldoPremiGagal - `[keputusan work owner 09-10-2026, penyimpangan sadar]`: tabel ARASAPAS tidak ada di DEV (ORA-00942
// lewat ketiga cara baca, 09-10-2026), padahal Send to Committe selalu menjalankan CekPremiLunas_Act. Di luar
// produksi (`IsPEGAPROD` salah) tabel yang tidak terbaca = saldo kosong (premi dianggap lunas), dicatat di log; di
// produksi galatnya diteruskan. Galat lain tidak pernah ditelan.
func saldoPremiGagal(produksi bool, err error) (string, error) {
	if produksi || !strings.Contains(err.Error(), "ORA-00942") {
		return "", err
	}
	log.Printf("claimprop: CekLunasPremi_Sql dilewati di luar produksi (ARASAPAS.INVOICE tidak terbaca): %v", err)
	return "", nil
}

// AdaProteksiPremi = CekProteksiKlaim (OPENPROTEKSI_EDM `TYPE = '5' AND STS_AKSEP = '1' AND POLICY_NO`).
func (a *Acuan) AdaProteksiPremi(ctx context.Context, nopolis string) (bool, error) {
	t, err := a.q("OPENPROTEKSI_EDM")
	if err != nil {
		return false, err
	}
	_, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT POLICY_NO FROM %s WHERE TYPE = '5' AND STS_AKSEP = '1' AND POLICY_NO = :1
		FETCH FIRST 1 ROWS ONLY`, t), nopolis)
	return ada, err
}

// StatusKasir = GetStatusKasir_SQL (DIRECTTOKASIR_LOG.KET menurut NOAKSEPTASI; tanpa ORDER BY di korpus - di sini
// terbaru).
func (a *Acuan) StatusKasir(ctx context.Context, noAksep string) (string, bool, error) {
	t, err := a.q("DIRECTTOKASIR_LOG")
	if err != nil {
		return "", false, err
	}
	return a.satu(ctx, fmt.Sprintf(`SELECT KET FROM %s WHERE NOAKSEPTASI = :1 ORDER BY TGL_INPUT DESC NULLS LAST
		FETCH FIRST 1 ROWS ONLY`, t), noAksep)
}

// StatusKonversi = getStatusKonversi_SQL (`COUNT(1) FROM REINSURANCE.TRLOSS_DETAIL_T WHERE NO_AKSEP`) - HANYA di
// produksi (`IsPEGAPROD`); di luar produksi kosong (HitServiceToKasir_Act langkah 3 keluar). Skema REINSURANCE tidak
// terlihat dari akun DEV (OQ-CP-11, DBA).
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

// EmailCeding = GetEmailCeding_SQL (`GL.F_GET_EMAIL(:ceding) FROM DUAL`) - HANYA dibaca saat muatan kasir disusun di
// produksi. Skema GL tidak terlihat dari akun DEV (OQ-CP-11).
func (a *Acuan) EmailCeding(ctx context.Context, ceding string) (string, error) {
	if !a.Produksi {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT GL.F_GET_EMAIL(:1) FROM DUAL`, ceding)
	return v, err
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

// DaftarJenisReas = BrowseReinsuranceType_RD (`.Type = Param.Type` bila diisi, urut `.Note`).
func (a *Acuan) DaftarJenisReas(ctx context.Context, tipe string) ([]models.Pilihan, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return nil, err
	}
	if tipe == "" {
		return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT ID, NOTE FROM %s ORDER BY NOTE`, t), 2))
	}
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT ID, NOTE FROM %s WHERE TYPE = :1 ORDER BY NOTE`, t), 2, tipe))
}

// DaftarProvinsi = BrowseProvince_RD (distinct PROVINCENAME, `Nation = "INDONESIA"`, `PROVINCENAME Contains`).
func (a *Acuan) DaftarProvinsi(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("RW")
	if err != nil {
		return nil, err
	}
	pola := "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%"
	// Alias wajib: PROVINCENAME dua kali tanpa alias membuat ORDER BY ambigu (ORA-00960, uji ujidev 07-10-2026).
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT PROVINCENAME AS NILAI, PROVINCENAME AS LABEL FROM %s
		WHERE UPPER(NATION) = 'INDONESIA' AND UPPER(PROVINCENAME) LIKE :1 ORDER BY NILAI FETCH FIRST 500 ROWS ONLY`, t),
		2, pola))
}

// BarisSebab - satu baris grid "List Cause of Loss".
type BarisSebab struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// DaftarSebab = grid "List Cause of Loss" (RD BrowseCouseOfLoss_Business: `BISNISID = Param.id` - tidak diisi pemanggil,
// tidak dipasang; `DESCRIPTION Contains Param.Name`; urut DESCRIPTION; 500).
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

// BarisKatastrofe - satu baris grid katastrofe (RD GetCatastrope_RD).
type BarisKatastrofe struct {
	ID                string `json:"id"`
	StsKatastrofe     string `json:"stsKatastrofe"`
	NonKatastrofeType string `json:"nonKatastrofeType"`
	Note              string `json:"note"`
}

// DaftarKatastrofe = RD GetCatastrope_RD (`KLAIMTYPE = "NON-LIFE"`; F1/F2 bernilai kosong tidak dipasang).
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

// BarisKatastrofeID - satu baris katastrofe menurut ID (tombol Choose dibaca ulang di server).
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

// BarisOS - satu baris ringkasan OS (DataOutstandingTreatyin / GetDataKlaimTreatyin).
type BarisRingkasanOS struct {
	Currency, CurrencyID, GrossValue, Value, NoClaim, NoPolis, AcceptedNo, Type, KursValue, Tanggal, TotalGross string
}

// RingkasanOS - baris OS_AKSEPTASI_KLAIM satu nomor klaim: kolom datar, jatuh ke JSON_VALUE DATA_JSON baris warisan.
// Urut `TANGGAL ASC, AcceptedNo DESC` (DataOutstandingTreatyin).
func (a *Acuan) RingkasanOS(ctx context.Context, noClaim string) ([]BarisRingkasanOS, error) {
	t, err := a.q("OS_AKSEPTASI_KLAIM")
	if err != nil {
		return nil, err
	}
	jv := func(p string) string {
		return fmt.Sprintf(`JSON_VALUE(a.DATA_JSON, '$.%s' RETURNING VARCHAR2(4000))`, p)
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT NVL(a.CURRENCY, %s), NVL(a.CURRENCYID, %s),
		NVL(%s, %s), NVL(%s, %s), a.NOCLAIM, a.NOPOLIS, NVL(a.ACCEPTEDNO, %s), NVL(a.TYPE, %s),
		NVL(%s, %s), TO_CHAR(a.TANGGAL, 'YYYY-MM-DD'), %s
		FROM %s a WHERE a.NOCLAIM = :1 ORDER BY a.TANGGAL ASC, NVL(a.ACCEPTEDNO, %s) DESC`,
		jv("Currency"), jv("CurrencyID"), dec("a.GROSSVALUE"), jv("GrossValue"), dec("a.VALUE"), jv("Value"),
		jv("AcceptedNo"), jv("Type"), dec("a.KURSVALUE"), jv("KursValue"), jv("TotalGross"), t, jv("AcceptedNo")), 11, noClaim)
	if err != nil {
		return nil, err
	}
	var out []BarisRingkasanOS
	for _, r := range rows {
		out = append(out, BarisRingkasanOS{Currency: r[0], CurrencyID: r[1], GrossValue: rapikanDesimal(r[2]),
			Value: rapikanDesimal(r[3]), NoClaim: r[4], NoPolis: r[5], AcceptedNo: r[6], Type: r[7],
			KursValue: rapikanDesimal(r[8]), Tanggal: r[9], TotalGross: r[10]})
	}
	return out, nil
}

// NamaPelaku - nama tampilan akun (M_LOGIN_GO.NAME menurut LOGIN_ID, pola NB Treaty In); kosong = akun itu sendiri.
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
