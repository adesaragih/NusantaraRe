package repository

// Untuk apa berkas ini: ACUAN - bacaan baca-saja yang dibutuhkan port activity (`models.Acuan`) dan pemilih layar Claim
// Fac In. Setiap kueri meniru satu RDB-List / Report Definition korpus `Claim Fac In/RDBList|ReportDefinition` (nama
// rule di komentar); alias CARI diluruskan; teks pencarian SELALU diikat (perbaikan prompt §5 butir 1). Pola disalin dari
// `modul/claimnonprop/backend/repository/acuan.go` (asal Claim Prop), bukan diimpor.
//
// ⚠️ Dokumen JSON warisan (JSON_POLIS.DATA_JSON, JSON_KLAIM.DATA_JSON) DIBACA - baca-saja.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimfacin/backend/models"
)

// Acuan - bacaan acuan Oracle (memenuhi `models.Acuan`).
type Acuan struct {
	g *Gudang
	// Produksi - padanan `When/IsPEGAPROD`: bacaan skema luar yang di korpus hanya berjalan di produksi.
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

func pola(cari string) string { return "%" + strings.ToUpper(strings.TrimSpace(cari)) + "%" }

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

// ---------------------------------------------------------------- polis

// sqlCariPolis = GetPolisForClaim_SQL (FACINPRODUCTION + M_BUSINESS + JSON_POLIS.PRODKE) dengan klausa SearchPolis_act
// langkah 1-4 sebagai teks TERIKAT (XML: `{ASIS:InputData.CARI1}` - teks pengguna disisipkan apa adanya).
func sqlCariPolis(fp, bisnis, jp, kolom string, persis bool) string {
	op := "LIKE :1 || '%'"
	if persis {
		op = "= :1"
	}
	return fmt.Sprintf(`SELECT DISTINCT a.NOPOLIS, a.INSUREDNAME, a.SOB, a.CEDINGCO, a.QQNAME,
		       TO_CHAR(a.BEGINDATE, 'YYYY-MM-DD'), TO_CHAR(a.ENDDATE, 'YYYY-MM-DD'),
		       (SELECT b.PRODKE FROM %s b WHERE b.IDPEGA = a.IDPEGA FETCH FIRST 1 ROWS ONLY),
		       (SELECT JSON_VALUE(c.JSONDATA, '$.Note' RETURNING VARCHAR2(400)) FROM %s c WHERE c.ID = a.BUSINESSCODE)
		  FROM %s a WHERE a.%s %s
		 FETCH FIRST %d ROWS ONLY`, jp, bisnis, fp, kolom, op, models.BatasCariPolis)
}

// CariPolis = SearchPolis_act + GetPolisForClaim_SQL.
func (a *Acuan) CariPolis(ctx context.Context, jenis, teks string) ([]models.BarisPolisCari, error) {
	kolom, persis := map[string]string{models.CariNoPolis: "NOPOLIS", models.CariCeding: "CEDINGCO",
		models.CariInsured: "INSUREDNAME", models.CariQQ: "QQNAME"}[jenis], jenis == models.CariNoPolis
	if kolom == "" || strings.TrimSpace(teks) == "" {
		return []models.BarisPolisCari{}, nil
	}
	fp, err := a.q("FACINPRODUCTION")
	if err != nil {
		return nil, err
	}
	bis, err := a.q("M_BUSINESS")
	if err != nil {
		return nil, err
	}
	jp, err := a.q("JSON_POLIS")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, sqlCariPolis(fp, bis, jp, kolom, persis), 9, strings.TrimSpace(teks))
	if err != nil {
		return nil, err
	}
	out := []models.BarisPolisCari{}
	for _, r := range rows {
		out = append(out, models.BarisPolisCari{PolicyNo: r[0], CustomerName: r[1], SourceOfBusinessName: r[2],
			CedingCoName: r[3], QQ: r[4], StartDateTime: r[5], EndDateTime: r[6], Prodke: r[7], BusinessName: r[8]})
	}
	return out, nil
}

// DokumenPolis = GetCopyNBForClaim (`DATA_JSON FROM JSON_POLIS WHERE NOPOLIS AND PRODKE`), baris terbaru.
func (a *Acuan) DokumenPolis(ctx context.Context, nopolis, prodke string) ([]byte, bool, error) {
	t, err := a.q("JSON_POLIS")
	if err != nil {
		return nil, false, err
	}
	q := fmt.Sprintf(`SELECT DATA_JSON FROM %s WHERE NOPOLIS = :1 AND PRODKE = :2 AND DATA_JSON IS NOT NULL
		ORDER BY TGL_INPUT DESC NULLS LAST FETCH FIRST 1 ROWS ONLY`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, false, err
	}
	var v sql.NullString
	err = a.g.db.QueryRowContext(ctx, q, nopolis, prodke).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("repository: membaca JSON_POLIS: %w", err)
	}
	return []byte(v.String), v.Valid, nil
}

// AdaPolis - nomor polis + prodke ada di FACINPRODUCTION (pilihan tautan ViewPolis dibaca ulang di server).
func (a *Acuan) AdaPolis(ctx context.Context, nopolis, prodke string) (bool, error) {
	fp, err := a.q("FACINPRODUCTION")
	if err != nil {
		return false, err
	}
	jp, err := a.q("JSON_POLIS")
	if err != nil {
		return false, err
	}
	_, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT 1 FROM %s a WHERE a.NOPOLIS = :1 AND EXISTS (SELECT 1 FROM %s b
		WHERE b.IDPEGA = a.IDPEGA AND b.PRODKE = :2) FETCH FIRST 1 ROWS ONLY`, fp, jp), nopolis, prodke)
	return ada, err
}

// RiwayatKlaimPolis = CariHistoryClaim_SQL (JSON_KLAIM.DATA_JSON DateOfLoss / ClaimNo / CauseOfLoss / ClaimEstimate
// menurut NOPOLIS, urut TGL_INPUT) DITAMBAH klaim FACIN sistem baru (T_GENERAL_CLAIM ber-LINI FACIN) - JSON_KLAIM sistem
// baru tidak berisi DATA_JSON (keputusan work owner), pola Claim Prop.
func (a *Acuan) RiwayatKlaimPolis(ctx context.Context, nopolis string) ([]models.RiwayatKlaimPolis, error) {
	jk, err := a.q("JSON_KLAIM")
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
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT IDPEGA, NOKLAIM, DOL, SEBAB, NILAI FROM (
		SELECT a.IDPEGA, JSON_VALUE(a.DATA_JSON, '$.ClaimNo' RETURNING VARCHAR2(100)) NOKLAIM,
		       JSON_VALUE(a.DATA_JSON, '$.DateOfLoss' RETURNING VARCHAR2(100)) DOL,
		       JSON_VALUE(a.DATA_JSON, '$.CauseOfLoss' RETURNING VARCHAR2(1000)) SEBAB,
		       REPLACE(JSON_VALUE(a.DATA_JSON, '$.ClaimEstimate' RETURNING VARCHAR2(100)), ',', '.') NILAI, a.TGL_INPUT URUT
		  FROM %s a WHERE a.NOPOLIS = :1 AND a.DATA_JSON IS NOT NULL
		UNION ALL
		SELECT g.ID, g.CLAIM_NO, TO_CHAR(g.DATE_OF_LOSS, 'YYYY-MM-DD'), g.CAUSE_OF_LOSS, %s, w.TGL_CREATE
		  FROM %s g JOIN %s w ON w.ID = g.ID AND w.LINI = :2 AND w.TAHAP <> :3 WHERE g.POLICY_NO = :4)
		 ORDER BY URUT`, jk, dec("g.CLAIM_ESTIMATE"), gen, work), 5, nopolis, models.LiniFacIn, models.TahapKomite, nopolis)
	if err != nil {
		return nil, err
	}
	var out []models.RiwayatKlaimPolis
	for _, r := range rows {
		dol := r[2]
		if t, ok := models.UraiTanggal(dol); ok {
			dol = models.TeksTanggal(t)
		}
		out = append(out, models.RiwayatKlaimPolis{KunciKasus: r[0], ClaimNo: r[1], DateOfLoss: dol, CauseOfLoss: r[3],
			ClaimEstimate: r[4]})
	}
	return out, nil
}

// KlaimDitolak = RD RejectedClaim_RD (CLAIMREJECTED `.INSKEY`).
func (a *Acuan) KlaimDitolak(ctx context.Context, kunci string) (bool, error) {
	t, err := a.q("CLAIMREJECTED")
	if err != nil {
		return false, err
	}
	_, ada, err := a.satu(ctx, fmt.Sprintf(`SELECT 1 FROM %s WHERE INSKEY = :1 FETCH FIRST 1 ROWS ONLY`, t), kunci)
	return ada, err
}

// ---------------------------------------------------------------- mata uang dan jenis reasuransi

// KursStandar = CurrencyStandard (`POOLDATA.GETCURRENCYSTANDARD(:cur, SYSDATE)`).
func (a *Acuan) KursStandar(ctx context.Context, cur string) (string, error) {
	f, err := a.q("GETCURRENCYSTANDARD")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT %s FROM DUAL`, dec(f+"(:1, SYSDATE)")), cur)
	return v, err
}

// NamaMataUang = GetCurrency (`CURRENCY FROM CURRENCY WHERE ID`).
func (a *Acuan) NamaMataUang(ctx context.Context, cur string) (string, error) {
	t, err := a.q("CURRENCY")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT CURRENCY FROM %s WHERE ID = :1`, t), cur)
	return v, err
}

// DaftarMataUang = BrowseCurrency_RD (`Currency != "ITL"`, urut `.Currency`).
func (a *Acuan) DaftarMataUang(ctx context.Context) ([]models.Pilihan, error) {
	t, err := a.q("CURRENCY")
	if err != nil {
		return nil, err
	}
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT ID, CURRENCY FROM %s WHERE NVL(CURRENCY, ' ') <> 'ITL'
		ORDER BY CURRENCY`, t), 2))
}

// NamaJenisReas = RD BrowseReinsuranceType_RD (`REINSURANCETYPE.NOTE` menurut ID).
func (a *Acuan) NamaJenisReas(ctx context.Context, id string) (string, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT NOTE FROM %s WHERE ID = :1`, t), id)
	return v, err
}

// DaftarJenisReas = BrowseReinsuranceType_RD (dropdown Treaty Type, urut NOTE).
func (a *Acuan) DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error) {
	t, err := a.q("REINSURANCETYPE")
	if err != nil {
		return nil, err
	}
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT ID, NOTE FROM %s ORDER BY NOTE`, t), 2))
}

// ---------------------------------------------------------------- lokasi, pelapor, adjuster

// Wilayah = BrowseRW_SQL (RW + CITY menurut ZIPCODE; alias diluruskan, pola Claim Prop).
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

// DaftarNegara = BrowseCountry_RD (ID, COUNTRY; saringan nama).
func (a *Acuan) DaftarNegara(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("COUNTRY")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT COUNTRY, COUNTRY, ID FROM %s WHERE UPPER(COUNTRY) LIKE :1
		ORDER BY COUNTRY FETCH FIRST 500 ROWS ONLY`, t), 3, pola(cari))
	return pilihanID(rows, err, "CountryID")
}

// pilihanID - pilihan bernilai nama dengan ID di kolom ketiga (isi autocomplete `.ID -> <jalur ID>`).
func pilihanID(rows [][]string, err error, prop string) ([]models.Pilihan, error) {
	if err != nil {
		return nil, err
	}
	out := []models.Pilihan{}
	for _, r := range rows {
		out = append(out, models.Pilihan{Nilai: r[0], Label: r[1], Tambahan: map[string]string{prop: r[2]}})
	}
	return out, nil
}

// DaftarProvinsi = BrowseProvince_RD (RW: distinct PROVINCENAME menurut negara, `PROVINCENAME Contains`).
//
// ⚠️ PERBAIKAN (PARITAS `[penyimpangan sadar]`): section mengirim parameter `ID2 = CountryID` dan membaca `.Note` /
// `.ID` yang tidak ada di kolom RD (RD hanya memilih PROVINCENAME, menyaring `Param.Nation`) - di Pega daftar tidak
// tersaring dan ProvinceID tidak pernah terisi. Di sini disaring nama negara dan ID provinsi ikut terisi (RW.PROVINCEID).
func (a *Acuan) DaftarProvinsi(ctx context.Context, negara, cari string) ([]models.Pilihan, error) {
	t, err := a.q("RW")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT PROVINCENAME, PROVINCENAME, MAX(PROVINCEID) FROM %s
		WHERE UPPER(NATION) = UPPER(:1) AND UPPER(PROVINCENAME) LIKE :2 GROUP BY PROVINCENAME ORDER BY PROVINCENAME
		FETCH FIRST 500 ROWS ONLY`, t), 3, negara, pola(cari))
	return pilihanID(rows, err, "ProvinceID")
}

// DaftarKota = BrowseCity_RD (RW: distinct CITYNAME menurut provinsi; perbaikan parameter sama dengan DaftarProvinsi).
func (a *Acuan) DaftarKota(ctx context.Context, provinsi, cari string) ([]models.Pilihan, error) {
	t, err := a.q("RW")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT CITYNAME, CITYNAME, MAX(CITYID) FROM %s
		WHERE PROVINCENAME = :1 AND UPPER(CITYNAME) LIKE :2 GROUP BY CITYNAME ORDER BY CITYNAME
		FETCH FIRST 500 ROWS ONLY`, t), 3, provinsi, pola(cari))
	return pilihanID(rows, err, "CityID")
}

// DaftarDistrik = BrowseDistrictInputC_RD (DISTRICT menurut CityID).
func (a *Acuan) DaftarDistrik(ctx context.Context, kotaID, cari string) ([]models.Pilihan, error) {
	t, err := a.q("DISTRICT")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT DISTRICTNAME, DISTRICTNAME, ID FROM %s WHERE CITYID = :1
		AND UPPER(DISTRICTNAME) LIKE :2 ORDER BY DISTRICTNAME FETCH FIRST 500 ROWS ONLY`, t), 3, kotaID, pola(cari))
	return pilihanID(rows, err, "DistrictID")
}

// DaftarRW = BrowseRWInput_RD (RW menurut DistrictID; isi ID -> RWID, ZipCode -> PostalCode).
func (a *Acuan) DaftarRW(ctx context.Context, distrikID, cari string) ([]models.Pilihan, error) {
	t, err := a.q("RW")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT NOTE, NOTE, ID, ZIPCODE FROM %s WHERE DISTRICTID = :1
		AND UPPER(NOTE) LIKE :2 ORDER BY NOTE FETCH FIRST 500 ROWS ONLY`, t), 4, distrikID, pola(cari))
	if err != nil {
		return nil, err
	}
	out := []models.Pilihan{}
	for _, r := range rows {
		out = append(out, models.Pilihan{Nilai: r[0], Label: r[1], Tambahan: map[string]string{"RWID": r[2],
			"PostalCode": r[3]}})
	}
	return out, nil
}

// AgenKlien = GetLeaderReport (`id FROM agent WHERE clientname = :nama`) - loop Pega menimpa: baris terakhir.
func (a *Acuan) AgenKlien(ctx context.Context, namaKlien string) (string, error) {
	t, err := a.q("AGENT")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT ID FROM %s WHERE CLIENTNAME = :1 ORDER BY ROWID DESC
		FETCH FIRST 1 ROWS ONLY`, t), namaKlien)
	return v, err
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

// DaftarAdjuster = autocomplete Adjuster / Consultant ID (RD BrowseAdjusterConsultant: nilai `.ID`, tampil `.NAME`).
func (a *Acuan) DaftarAdjuster(ctx context.Context, cari string) ([]models.Pilihan, error) {
	t, err := a.q("ADJUSTERCONSULTANT")
	if err != nil {
		return nil, err
	}
	p := pola(cari)
	return pilihan(a.banyak(ctx, fmt.Sprintf(`SELECT ID, NAME FROM %s WHERE UPPER(ID) LIKE :1 OR UPPER(NAME) LIKE :2
		ORDER BY ID FETCH FIRST 500 ROWS ONLY`, t), 2, p, p))
}

// MarketingOfficer = SetMOClaim_Act langkah 4 (Obj-Browse MARKETINGOFFICER menurut ID).
func (a *Acuan) MarketingOfficer(ctx context.Context, id string) ([6]string, bool, error) {
	var out [6]string
	t, err := a.q("MARKETINGOFFICER")
	if err != nil {
		return out, false, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, CLIENTID, CLIENTNAME, TEAMGROUP, BRANCHDETAILID, BRANCHDETAILNAME
		FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY`, t), 6, id)
	if err != nil || len(rows) == 0 {
		return out, false, err
	}
	copy(out[:], rows[0])
	return out, true, nil
}

// ---------------------------------------------------------------- cause of loss dan katastrofe

// BarisSebab - satu baris grid "List Cause of Loss".
type BarisSebab struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// DaftarSebab = grid "List Cause of Loss" (CauseofLoss_Section, RD BrowseCouseOfLoss_Business; `DESCRIPTION Contains`;
// urut DESCRIPTION; 500). Parameter BISNISID tidak dikirim section - tidak disaring (pola Claim Prop).
func (a *Acuan) DaftarSebab(ctx context.Context, cari string) ([]BarisSebab, error) {
	t, err := a.q("V_D_CAUSE_OF_LOSS_BUSINESS")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT DISTINCT D_COL_ID, DESCRIPTION FROM %s WHERE UPPER(DESCRIPTION) LIKE :1
		ORDER BY DESCRIPTION FETCH FIRST 500 ROWS ONLY`, t), 2, pola(cari))
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
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT D_COL_ID, DESCRIPTION FROM %s WHERE D_COL_ID = :1
		FETCH FIRST 1 ROWS ONLY`, t), 2, id)
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
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT ID, STSKATASTROFE, NONKATASTROFETYPE, NOTE FROM %s
		WHERE KLAIMTYPE = 'NON-LIFE' AND UPPER(NVL(NOTE, ' ')) LIKE :1 ORDER BY TGL_INPUT DESC NULLS LAST
		FETCH FIRST 500 ROWS ONLY`, t), 4, pola(cari))
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

// ---------------------------------------------------------------- limit treaty (CheckLimit_Act1)

// BatasTreaty = GetDataTreatyLimit_Sql (TREATYBUSINESS x PROPORTIONALARRG TREATYDESCID 10001, kontrak TREATYCONTRACT
// berlaku pada tanggal mulai polis), baris pertama.
func (a *Acuan) BatasTreaty(ctx context.Context, kodeBisnis, jenisTreaty, mulaiPolis string) (models.BatasTreaty, bool, error) {
	tb, err := a.q("TREATYBUSINESS")
	if err != nil {
		return models.BatasTreaty{}, false, err
	}
	pa, err := a.q("PROPORTIONALARRG")
	if err != nil {
		return models.BatasTreaty{}, false, err
	}
	tc, err := a.q("TREATYCONTRACT")
	if err != nil {
		return models.BatasTreaty{}, false, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT p.TREATYYEAR, p.TREATYGROUPID, p.RP
		  FROM %s b JOIN %s p ON p.TREATYYEARID = b.TREATYYEARID AND p.REINSTYPEID = b.REINSTYPEID
		 WHERE b.BIZCODE = :1 AND b.REINSTYPEID = :2 AND p.TREATYDESCID = '10001'
		   AND EXISTS (SELECT 1 FROM %s tc WHERE tc.IDTREATYYEAR = b.TREATYYEARID AND tc.REINSTYPEID = b.REINSTYPEID
		               AND TO_DATE(:3, 'YYYY-MM-DD') BETWEEN tc.TREATYSTARTDATE AND tc.TREATYENDDATE)
		 FETCH FIRST 1 ROWS ONLY`, tb, pa, tc), 3, kodeBisnis, jenisTreaty, tanggalAtauNil(mulaiPolis))
	if err != nil || len(rows) == 0 {
		return models.BatasTreaty{}, false, err
	}
	return models.BatasTreaty{TreatyYear: rows[0][0], TreatyGroupID: rows[0][1], Limit: rows[0][2]}, true, nil
}

// tanggalAtauNil - "2006-01-02" dari teks tanggal halaman; tak terurai = NULL.
func tanggalAtauNil(s string) any {
	if t, ok := models.UraiTanggal(s); ok {
		return models.TeksTanggal(t)
	}
	return nil
}

// QuotaShare = GetQuotaShare (PROPORTIONALARRG TREATYDESCID 10001 menurut tahun, grup, induk).
func (a *Acuan) QuotaShare(ctx context.Context, tahun, grup, induk string) ([]models.BarisQuotaShare, error) {
	t, err := a.q("PROPORTIONALARRG")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT REINSTYPEID, REINSTYPENAME, PCT FROM %s
		WHERE TREATYYEAR = :1 AND TREATYGROUPID = :2 AND TREATYDESCID = '10001' AND PARENTREINSTYPEID = :3
		ORDER BY ID`, t), 3, tahun, grup, induk)
	if err != nil {
		return nil, err
	}
	var out []models.BarisQuotaShare
	for _, r := range rows {
		out = append(out, models.BarisQuotaShare{TreatyType: r[0], TreatyName: r[1], SharePercentage: r[2]})
	}
	return out, nil
}

// ReasuradurTreaty = GetListRetro_Sql (TREATYREINSURER menurut jenis, tahun, grup).
func (a *Acuan) ReasuradurTreaty(ctx context.Context, jenis, tahun, grup string) ([]models.ReasuradurTreaty, error) {
	t, err := a.q("TREATYREINSURER")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT REINSURERID, NAME, %s, %s FROM %s
		WHERE REINSTYPEID = :1 AND TREATYYEAR = :2 AND TREATYGROUPID = :3 ORDER BY ID`, dec("PCTSHARE"), dec("RICOMM"), t),
		4, jenis, tahun, grup)
	if err != nil {
		return nil, err
	}
	var out []models.ReasuradurTreaty
	for _, r := range rows {
		out = append(out, models.ReasuradurTreaty{ReinsurerID: r[0], ReinsurerName: r[1], PctShareAllObj: r[2],
			RiCommAllObj: r[3]})
	}
	return out, nil
}

// EstimasiKasusTerbuka = GetDataEstimation (`SUM(CLAIM_VALUE) REINSURANCE.TRLOSS_DETAIL_T WHERE NO_SPK = pyID AND
// NO_AKSEP IS NULL`) - HANYA di produksi (skema luar tidak terbaca dari akun DEV, pola OQ-CP-11).
func (a *Acuan) EstimasiKasusTerbuka(ctx context.Context, kasusID string) (string, error) {
	if !a.Produksi {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT TO_CHAR(SUM(CLAIM_VALUE), 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')
		FROM REINSURANCE.TRLOSS_DETAIL_T WHERE NO_SPK = :1 AND NO_AKSEP IS NULL`, kasusID)
	return v, err
}

// EstimasiPolisTerbuka = GetAllClaimWithSameNopolis_Sql - HANYA di produksi.
func (a *Acuan) EstimasiPolisTerbuka(ctx context.Context, nopolis string) (string, error) {
	if !a.Produksi {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT TO_CHAR(SUM(d.CLAIM_VALUE), 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')
		FROM REINSURANCE.TRLOSS_DETAIL_T d WHERE d.STS_AKSEP IS NULL
		AND EXISTS (SELECT 1 FROM REINSURANCE.TREATY_LOSS t WHERE t.NO_SPK = d.NO_SPK AND t.NO_POLIS = :1)`, nopolis)
	return v, err
}

// ---------------------------------------------------------------- komite dan wewenang

// RosterKomite = RD FilterEmailKomiteWithLimit (EMAILKOMITE `STS_KLAIM = 'FACIN' AND STS_AKTIF = '1'`, urut DEGREE).
func (a *Acuan) RosterKomite(ctx context.Context) ([]models.AnggotaKomite, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, fmt.Sprintf(`SELECT TO_CHAR(ID), OPERATOR_ID, EMAIL, JABATAN, TO_CHAR(DEGREE), %s, %s
		FROM %s WHERE STS_KLAIM = :1 AND STS_AKTIF = '1' ORDER BY DEGREE, ID`, dec("LIMIT_BOTTOM"), dec("LIMIT_TOP"), t),
		7, models.LiniFacIn)
	if err != nil {
		return nil, err
	}
	var out []models.AnggotaKomite
	for _, r := range rows {
		out = append(out, models.AnggotaKomite{ID: r[0], OperatorID: r[1], Email: r[2], Jabatan: r[3], Degree: r[4],
			LimitBottom: r[5], LimitTop: r[6]})
	}
	return out, nil
}

// TingkatPelaku - JABATAN baris roster FACIN aktif pelaku (akun pelaku ATAU workbasket aktif yang dipegangnya; DEGREE
// terkecil; pola Claim Prop).
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
	v, _, err := a.satu(ctx, sqlTingkatPelaku(t, lwb, wb), models.LiniFacIn, operatorID, operatorID)
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

// EmailPelaku - `OperatorID.pyAddress` pelaku: M_LOGIN_GO.EMAIL.
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

// ---------------------------------------------------------------- progres klaim (GetProgresClaim_ACT)

// ProgresKlaim = GetProgressClaim_SQL + GetSubProgressClaim_SQL: PROGRESSCLAIM kasus (urut TANGGAL) dengan baris
// SUBPROGRESSCLAIM ber-JENISPROGRES = POSITION di bawahnya (langkah 3.2.1).
func (a *Acuan) ProgresKlaim(ctx context.Context, kasusID string) ([]models.Baris, map[int][]models.Baris, error) {
	tp, err := a.q("PROGRESSCLAIM")
	if err != nil {
		return nil, nil, err
	}
	ts, err := a.q("SUBPROGRESSCLAIM")
	if err != nil {
		return nil, nil, err
	}
	pr, err := a.banyak(ctx, fmt.Sprintf(`SELECT IDPEGA, POSITION, TO_CHAR(TANGGAL, 'YYYY-MM-DD HH24:MI:SS'), STATUS
		FROM %s WHERE IDPEGA = :1 ORDER BY TANGGAL`, tp), 4, models.KunciInstans(kasusID))
	if err != nil {
		return nil, nil, err
	}
	sub, err := a.banyak(ctx, fmt.Sprintf(`SELECT IDPEGA, JENISPROGRES, TO_CHAR(TANGGALINPUT, 'YYYY-MM-DD HH24:MI:SS'),
		POSITION1, POSITION2, TO_CHAR(NEXT_FU, 'YYYY-MM-DD HH24:MI:SS'), USER_INPUT, STATUS
		FROM %s WHERE IDPROGRES = :1 ORDER BY TANGGALINPUT`, ts), 8, models.KunciInstans(kasusID))
	if err != nil {
		return nil, nil, err
	}
	var out []models.Baris
	anak := map[int][]models.Baris{}
	for i, r := range pr {
		out = append(out, models.Baris{"CaseID": r[0], "ProdKe": r[1], "AnalystTransferDate": r[2], "StatusClaim": r[3]})
		for _, s := range sub {
			if s[1] == r[1] {
				anak[i+1] = append(anak[i+1], models.Baris{"Notes": s[0], "Model": s[1], "ObjekTanggal": s[2],
					"RemarksPLA": s[3], "RemarksDLA": s[4], "TanggalFolloup": s[5], "TypeName": s[6], "Comment": s[7]})
			}
		}
	}
	return out, anak, nil
}

// ---------------------------------------------------------------- rekening, kasir (pola Claim Prop)

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

// IDBankRekening = HitServiceToKasir_Act 11-12 (Obj-Browse BANKACCOUNT `.IDOFBANK` menurut nama bank, cabang, nomor
// rekening).
func (a *Acuan) IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error) {
	t, err := a.q("BANKACCOUNT")
	if err != nil {
		return "", err
	}
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT IDOFBANK FROM %s WHERE NAMEOFBANK = :1 AND BRANCHOFBANK = :2
		AND ACCOUNTNO = :3 FETCH FIRST 1 ROWS ONLY`, t), bank, cabang, akun)
	return v, err
}

// StatusKasir = GetStatusKasir_SQL (DIRECTTOKASIR_LOG.KET menurut NOAKSEPTASI; tanpa ORDER BY di korpus - di sini
// terbaru, pola Claim Prop).
func (a *Acuan) StatusKasir(ctx context.Context, noAksep string) (string, bool, error) {
	t, err := a.q("DIRECTTOKASIR_LOG")
	if err != nil {
		return "", false, err
	}
	return a.satu(ctx, fmt.Sprintf(`SELECT KET FROM %s WHERE NOAKSEPTASI = :1 ORDER BY TGL_INPUT DESC NULLS LAST
		FETCH FIRST 1 ROWS ONLY`, t), noAksep)
}

// EmailCeding = GetEmailCeding_SQL (`GL.F_GET_EMAIL(:ceding) FROM DUAL`) - HANYA dibaca saat muatan kasir disusun di
// produksi (skema GL tidak terlihat dari akun DEV, pola Claim Prop).
func (a *Acuan) EmailCeding(ctx context.Context, ceding string) (string, error) {
	if !a.Produksi {
		return "", nil
	}
	v, _, err := a.satu(ctx, `SELECT GL.F_GET_EMAIL(:1) FROM DUAL`, ceding)
	return v, err
}

// ---------------------------------------------------------------- proteksi komite (pola Claim Prop)

// DegreeDirekturUtama - baris roster Direktur Utama. GetLimitDirekturUtama_SQL mencarinya menurut NAMA orang (dibuang,
// prompt §3 bawaan c); kunci DEGREE 6 mengikuti Claim Prop (OQ-CP-10).
const DegreeDirekturUtama = "6"

// LimitDirekturUtama = GetLimitDirekturUtama_SQL (SendPICProtect_Act 5): LIMIT_BOTTOM baris Direktur Utama. DEV 10-10-2026:
// satu baris DEGREE 6 (STS_KLAIM kosong, XML pun tanpa saringan lini); ORDER BY ID supaya tetap pasti bila bertambah.
func (a *Acuan) LimitDirekturUtama(ctx context.Context) (string, bool, error) {
	t, err := a.q("EMAILKOMITE")
	if err != nil {
		return "", false, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE DEGREE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY`, dec("LIMIT_BOTTOM"), t)
	return a.satu(ctx, q, DegreeDirekturUtama)
}

// SaldoPremi = CekLunasPremi_Sql (`SUM(IVD_TRANS_SIGN * IVD_TOTAL)` INVOICE x DETAIL_INVOICE - sinonim publik ARASAPAS).
// Di luar produksi tabel yang tidak terbaca (ORA-00942) = saldo kosong / premi dianggap lunas (pola Claim Prop, keputusan
// work owner 09-10-2026).
func (a *Acuan) SaldoPremi(ctx context.Context, invoice, cur string) (string, error) {
	v, _, err := a.satu(ctx, fmt.Sprintf(`SELECT %s FROM ARASAPAS.INVOICE a, ARASAPAS.DETAIL_INVOICE b
		WHERE a.INV_INV_NO = b.INV_INV_NO AND a.CURR_ENTRY_NO = b.CURR_ENTRY_NO AND a.INV_INV_NO = :1 AND a.INV_LKU_ID = :2`,
		dec("SUM(IVD_TRANS_SIGN * IVD_TOTAL)")), invoice, cur)
	if err != nil {
		if a.Produksi || !strings.Contains(err.Error(), "ORA-00942") {
			return "", err
		}
		log.Printf("claimfacin: CekLunasPremi_Sql dilewati di luar produksi (ARASAPAS.INVOICE tidak terbaca): %v", err)
		return "", nil
	}
	return v, nil
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
