package repository

// Untuk apa berkas ini: PEMBACAAN DATA KONTRAK DAN TABEL ACUAN warisan -
// padanan RD dan RDB-List yang dibaca layar dan aktivitas NB Treaty In
// (INVENTARIS-XML.md bab 8-9). Satu fungsi = satu rule; nama rulenya ditulis di
// atas tiap fungsi. Seluruhnya HANYA MEMBACA, kecuali tidak satu pun.
//
// ⭐ View kontrak `POOLDATA.TREATYINDETAILJOINEDM` (= TREATYINDETAIL UNION ALL
// TREATYINDETAILEDM) - `[keputusan work owner]` P29: data kontrak dari view
// relasional, BUKAN JSONDATA (AC 15, 16).
//
// ⚠️ TIPE KOLOM VIEW TIDAK DIKETAHUI dari korpus (naskah view diserahkan work
// owner, tipe kolomnya tidak). Supaya angka tidak lewat float dan tanggal tidak
// bergantung NLS sesi, tipe dibaca SEKALI dari katalog `SYS.ALL_TAB_COLUMNS`
// dan ekspresi SELECT dipilih per tipe: NUMBER -> TO_CHAR TM9 bertitik, DATE ->
// satu format tanggal, selainnya apa adanya (`COMMENCEMENT`/`TERMINATION`
// "dibaca apa adanya", PERTANYAAN-untuk-DBA P29).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// KolomRDDetail - 33 kolom RD `BrowseTreatyJoinEDM` / `BrowseTreatyInDetail`
// (urutan RD), ditambah dua kolom view yang dibaca halaman master.
var KolomRDDetail = []string{
	"ID", "TREATYID", "TREATYCONTRACTNAME", "PROPORTIONTYPE", "TREATYTYPE", "TREATYGROUP",
	"TREATYGROUPID", "CLASSOFBUSINESSID", "CLASSOFBUSINESS", "LIMITCURRENCY", "LIMITVALUE",
	"RETENTIONCURRENCY", "RETENTIONVALUE", "EPICURRENCY", "EPIVALUE", "MDPCURRENCY", "MDPVALUE",
	"NETPREMICURRENCY", "NETPREMIVALUE", "SOBID", "SOB", "CEDING", "CEDINGID", "TREATYYEAR",
	"INSTALLMENTNO", "LAYER", "LAYERPART", "LAYERTYPE", "LAYERPARTTYPE", "DEDUCTION1",
	"DEDUCTION2", "SHARECURRENCY", "SHAREVALUE",
}

// kolomMasterView - kolom view untuk halaman TreatyIn (`TerapkanMasterKontrak`).
var kolomMasterView = []string{"COMMENCEMENT", "TERMINATION"}

// batasDaftarDetail - baris terbanyak popup pilih bisnis.
const batasDaftarDetail = 500

// tipeKolom - cache tipe kolom per objek (nama berskema).
var tipeKolom sync.Map // map[string]map[string]string

// tipeKolomObjek membaca DATA_TYPE setiap kolom objek dari katalog Oracle.
func (g *Gudang) tipeKolomObjek(ctx context.Context, objek string) (map[string]string, error) {
	nama, err := g.nama(objek)
	if err != nil {
		return nil, err
	}
	if v, ada := tipeKolom.Load(nama); ada {
		return v.(map[string]string), nil
	}
	q := `SELECT COLUMN_NAME, DATA_TYPE FROM SYS.ALL_TAB_COLUMNS WHERE OWNER = :1 AND TABLE_NAME = :2`
	rows, err := g.db.QueryContext(ctx, q, g.db.Skema(), objek)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca tipe kolom %s: %w", nama, err)
	}
	defer rows.Close()
	tipe := map[string]string{}
	for rows.Next() {
		var k, t string
		if err := rows.Scan(&k, &t); err != nil {
			return nil, fmt.Errorf("repository: membaca tipe kolom %s: %w", nama, err)
		}
		tipe[k] = t
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(tipe) == 0 {
		return nil, fmt.Errorf("%w: objek %s tidak terbaca di katalog", ErrDataKontrakTidakAda, nama)
	}
	tipeKolom.Store(nama, tipe)
	return tipe, nil
}

// ekspresiTipe memilih ekspresi SELECT satu kolom menurut tipe Oracle-nya.
func ekspresiTipe(kolom, tipe string) string {
	switch {
	case tipe == "NUMBER" || tipe == "FLOAT" || strings.HasPrefix(tipe, "BINARY_"):
		return fmt.Sprintf(db.FmtDesimal, kolom)
	case tipe == "DATE" || strings.HasPrefix(tipe, "TIMESTAMP"):
		return fmt.Sprintf("TO_CHAR(%s, '%s')", kolom, fmtTanggal)
	}
	return kolom
}

// pilihKolom merakit daftar SELECT untuk `kolom`. Kolom yang TIDAK ada di objek
// adalah galat (AC 89: "nol medan yang dipakai tidak tersedia") - dibaca kosong
// diam-diam berarti layar menampilkan data kontrak yang tidak lengkap.
// (`CURRENCYID` yang dibaca `InputPolicyTreatyInDetail_preACT` langkah 3 memang
// bukan kolom RD maupun view, dan tidak ada di `kolom`.)
func pilihKolom(objek string, tipe map[string]string, kolom []string) ([]string, []string, error) {
	var eks, ada, hilang []string
	for _, k := range kolom {
		if t, ok := tipe[k]; ok {
			eks = append(eks, ekspresiTipe(k, t))
			ada = append(ada, k)
		} else {
			hilang = append(hilang, k)
		}
	}
	if len(hilang) > 0 {
		return nil, nil, fmt.Errorf("%w: kolom %s tidak ada di %s", ErrDataKontrakTidakAda, strings.Join(hilang, ", "), objek)
	}
	return eks, ada, nil
}

func pindaiKontrak(rows *sql.Rows, ada []string) (models.BarisKontrak, error) {
	nilai := make([]sql.NullString, len(ada))
	tujuan := make([]any, len(ada))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	if err := rows.Scan(tujuan...); err != nil {
		return nil, fmt.Errorf("repository: membaca baris kontrak: %w", err)
	}
	b := models.BarisKontrak{}
	for i, k := range ada {
		v := teks(nilai[i])
		if strings.HasPrefix(v, ".") || strings.HasPrefix(v, "-.") {
			v = rapikanDesimal(v)
		}
		b[k] = v
	}
	return b, nil
}

// DetailKontrak = RD `BrowseTreatyJoinEDM` dengan `Param.ID` (filter G
// `.ID = Param.ID`; parameter lain kosong -> diabaikan) - baris pertama.
// Dipakai pilih bisnis (`InputPolicyTreatyInDetail_preACT` langkah 1-3) dan
// pemuatan ulang halaman master saat kasus dibuka.
//
// ⛔ Nol baris = galat yang menghentikan proses (AC 36-38, spec §5.9) - Pega
// menelannya dan mengisi halaman dengan nilai kosong.
func (g *Gudang) DetailKontrak(ctx context.Context, id string) (models.BarisKontrak, error) {
	tipe, err := g.tipeKolomObjek(ctx, viewDetailGabung)
	if err != nil {
		return nil, err
	}
	nama, _ := g.nama(viewDetailGabung)
	eks, ada, err := pilihKolom(viewDetailGabung, tipe, append(append([]string{}, KolomRDDetail...), kolomMasterView...))
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1 FETCH FIRST 1 ROWS ONLY`, strings.Join(eks, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, id)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca data kontrak: %w", err)
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("repository: membaca data kontrak: %w", err)
		}
		return nil, fmt.Errorf("%w: ID %q", ErrDataKontrakTidakAda, id)
	}
	return pindaiKontrak(rows, ada)
}

// kolomKomisi - kolom view yang dibaca `TreatyInputPctCommSpreading` langkah
// 2-2.1.1.1 (`models.TreatyInputPctCommSpreading`).
var kolomKomisi = []string{"ID", models.KolomJenisTreaty, models.KolomGrupTreaty, models.KolomRIOGR, models.KolomRIONR}

// KomisiKontrak = pengganti relasional `FetchMasterTreatyIn` untuk
// `TreatyInputPctCommSpreading` (RDB `BrowseTreatyInJoinEDM`: JSON
// `M_TREATY_IN` UNION ALL `M_TREATY_IN_EDM` `where ID = PolicyTreatyIn.NoOffer`).
// Yang dibaca: baris view `TREATYINDETAILJOINEDM` (= TREATYINDETAIL UNION ALL
// TREATYINDETAILEDM) ber-`TREATYID = NoOffer` - satu baris per
// `Limits(n).Detail(m)` master - dengan kolom TREATYTYPE, TREATYGROUP, RIOGR,
// RIONR. Syarat langkah 2.1/2.1.1.1 diterapkan `models`, bukan di SQL.
//
// ⚠️ Urutan Limits/Detail dokumen JSON tidak ada di view: ORDER BY ID
// (lihat `models/komisi.go`). Nol baris = nol perubahan (kalang Pega atas
// halaman kosong), bukan galat.
func (g *Gudang) KomisiKontrak(ctx context.Context, treatyID string) ([]models.BarisKontrak, error) {
	if strings.TrimSpace(treatyID) == "" {
		return nil, nil // `where ID = NULL` - nol baris
	}
	tipe, err := g.tipeKolomObjek(ctx, viewDetailGabung)
	if err != nil {
		return nil, err
	}
	nama, _ := g.nama(viewDetailGabung)
	eks, ada, err := pilihKolom(viewDetailGabung, tipe, kolomKomisi)
	if err != nil {
		return nil, err
	}
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE TREATYID = :1 ORDER BY ID`, strings.Join(eks, ", "), nama)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, treatyID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca komisi kontrak: %w", err)
	}
	defer rows.Close()
	var out []models.BarisKontrak
	for rows.Next() {
		b, err := pindaiKontrak(rows, ada)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// SaringanDetail - parameter RD `BrowseTreatyInDetail`; kosong = diabaikan
// (aturan filter RD Pega).
type SaringanDetail struct {
	// TreatyID - filter I `.TREATYID Contains Param.TREATYID`.
	TreatyID string
}

// DaftarDetailKontrak = RD `BrowseTreatyInDetail` (kelas
// `ASM-FW-GISFW-Int-TREATYINDETAIL`, tabel TREATYINDETAIL) - grid popup
// `Section/BusinessAndSOBList`. Popup tidak mengirim parameter, sehingga
// filter A-L RD diabaikan; yang tersisa satu saringan teks nomor kontrak.
func (g *Gudang) DaftarDetailKontrak(ctx context.Context, s SaringanDetail) ([]models.BarisKontrak, error) {
	tipe, err := g.tipeKolomObjek(ctx, tabelDetail)
	if err != nil {
		return nil, err
	}
	nama, _ := g.nama(tabelDetail)
	eks, ada, err := pilihKolom(tabelDetail, tipe, KolomRDDetail)
	if err != nil {
		return nil, err
	}
	var q string
	var args []any
	if s.TreatyID != "" {
		q = fmt.Sprintf(`SELECT %s FROM %s WHERE UPPER(TREATYID) LIKE :1 ORDER BY TREATYID, ID FETCH FIRST %d ROWS ONLY`,
			strings.Join(eks, ", "), nama, batasDaftarDetail)
		args = append(args, "%"+strings.ToUpper(s.TreatyID)+"%")
	} else {
		q = fmt.Sprintf(`SELECT %s FROM %s ORDER BY TREATYID, ID FETCH FIRST %d ROWS ONLY`,
			strings.Join(eks, ", "), nama, batasDaftarDetail)
	}
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar kontrak: %w", err)
	}
	defer rows.Close()
	var out []models.BarisKontrak
	for rows.Next() {
		b, err := pindaiKontrak(rows, ada)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// ------------------------------------------------------------------ satu nilai

// satuTeks menjalankan kueri satu kolom satu baris; nol baris = "" (perilaku
// `pxResults(1).X` atas hasil kosong di Pega).
func (g *Gudang) satuTeks(ctx context.Context, apa, q string, args ...any) (string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var v sql.NullString
	err := g.db.QueryRowContext(ctx, q, args...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: %s: %w", apa, err)
	}
	return v.String, nil
}

// IDMataUangDariNama = RDB `GetCurrencyIDByName` (`SetTreatyCurrencyID`).
func (g *Gudang) IDMataUangDariNama(ctx context.Context, nama string) (string, error) {
	t, err := g.nama(tabelMataUang)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca ID mata uang",
		fmt.Sprintf(`SELECT TO_CHAR(ID) FROM %s WHERE CURRENCY = :1 FETCH FIRST 1 ROWS ONLY`, t), nama)
}

// NamaMataUang = RDB `GetCurrency` (`SetCurrency_act`).
func (g *Gudang) NamaMataUang(ctx context.Context, id string) (string, error) {
	t, err := g.nama(tabelMataUang)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca nama mata uang",
		fmt.Sprintf(`SELECT CURRENCY FROM %s WHERE TO_CHAR(ID) = :1 FETCH FIRST 1 ROWS ONLY`, t), id)
}

// OJKGrupTreaty = RD `BrowseTreatyGroup_RD` (`FetchTreatyGroupOJK`).
func (g *Gudang) OJKGrupTreaty(ctx context.Context, grupID string) (string, error) {
	t, err := g.nama(tabelGrupTreaty)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca OJK grup treaty",
		fmt.Sprintf(`SELECT TO_CHAR(OJKBUSINESSID) FROM %s WHERE TO_CHAR(ID) = :1 FETCH FIRST 1 ROWS ONLY`, t), grupID)
}

// OldIDGrupTreaty = RDB `FetchTreatyGroupOLDID` (`FetchTreatyGroupOldID`).
func (g *Gudang) OldIDGrupTreaty(ctx context.Context, grupID string) (string, error) {
	t, err := g.nama(tabelGrupTreaty)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca OLDID grup treaty",
		fmt.Sprintf(`SELECT TO_CHAR(OLDID) FROM %s WHERE TO_CHAR(ID) = :1 FETCH FIRST 1 ROWS ONLY`, t), grupID)
}

// KlienDariNama = RDB `GetClientID_SQL`:
// `replace(name,' ',”) = replace({SearchClient.CARI1},' ',”)`.
func (g *Gudang) KlienDariNama(ctx context.Context, nama string) (string, error) {
	t, err := g.nama(tabelKlien)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca ID klien",
		fmt.Sprintf(`SELECT TO_CHAR(ID) FROM %s WHERE REPLACE(NAME, ' ', '') = REPLACE(:1, ' ', '') FETCH FIRST 1 ROWS ONLY`, t), nama)
}

// StsPKPAgen = RD `BrowseClientName_RD` (`SetPPNPPH` langkah 1-3): filter
// D `.ID = Param.ID` (SourceOfBusiness) dan E `.StatusActive IS NULL`.
// `[data DBA - belum dikonfirmasi]` nama kolom STS_PKP / STATUSACTIVE tabel
// AGENT diambil dari nama properti RD kelas `ASM-FW-GISFW-Int-AGENT`.
func (g *Gudang) StsPKPAgen(ctx context.Context, sobID string) (string, error) {
	if sobID == "" {
		return "", nil
	}
	t, err := g.nama(tabelAgen)
	if err != nil {
		return "", err
	}
	return g.satuTeks(ctx, "membaca status PKP agen",
		fmt.Sprintf(`SELECT TO_CHAR(STS_PKP) FROM %s WHERE TO_CHAR(ID) = :1 AND STATUSACTIVE IS NULL FETCH FIRST 1 ROWS ONLY`, t), sobID)
}

// BisnisDariKunci = RDB `GetOldIDBusiness_SQL`:
// `where (ID = CARI2 OR NOTE = CARI2) AND GROUPPANEL IS NOT NULL`.
// ⚠️ RDB tanpa ORDER BY (baris pertama urutan basis data); di sini ORDER BY ID
// supaya hasilnya tetap - penyimpangan kecil, dicatat.
func (g *Gudang) BisnisDariKunci(ctx context.Context, kunci string) (models.BarisBisnis, error) {
	t, err := g.nama(tabelBisnis)
	if err != nil {
		return models.BarisBisnis{}, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(OLDID), TO_CHAR(GROUPPANEL), TO_CHAR(ID) FROM %s
	  WHERE (TO_CHAR(ID) = :1 OR NOTE = :2) AND GROUPPANEL IS NOT NULL ORDER BY ID FETCH FIRST 1 ROWS ONLY`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisBisnis{}, err
	}
	var a, b, c sql.NullString
	err = g.db.QueryRowContext(ctx, q, kunci, kunci).Scan(&a, &b, &c)
	if errors.Is(err, sql.ErrNoRows) {
		return models.BarisBisnis{}, nil
	}
	if err != nil {
		return models.BarisBisnis{}, fmt.Errorf("repository: membaca bisnis: %w", err)
	}
	return models.BarisBisnis{OldID: a.String, GroupPanel: b.String, ID: c.String}, nil
}

// MO = `CheckDataMkt` langkah 3 (`Obj-Browse` marketing officer, `.ID =
// QuotationData.MOID`).
func (g *Gudang) MO(ctx context.Context, id string) (models.BarisMO, error) {
	t, err := g.nama(tabelMO)
	if err != nil {
		return models.BarisMO{}, err
	}
	q := fmt.Sprintf(`SELECT TO_CHAR(ID), TO_CHAR(CLIENTID), CLIENTNAME, TO_CHAR(TEAMGROUP), TO_CHAR(BRANCHDETAILID), BRANCHDETAILNAME
	  FROM %s WHERE TO_CHAR(ID) = :1 FETCH FIRST 1 ROWS ONLY`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return models.BarisMO{}, err
	}
	var v [6]sql.NullString
	err = g.db.QueryRowContext(ctx, q, id).Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5])
	if errors.Is(err, sql.ErrNoRows) {
		return models.BarisMO{}, nil
	}
	if err != nil {
		return models.BarisMO{}, fmt.Errorf("repository: membaca marketing officer: %w", err)
	}
	return models.BarisMO{ID: v[0].String, ClientID: v[1].String, ClientName: v[2].String,
		TeamGroup: v[3].String, BranchDetailID: v[4].String, BranchDetailName: v[5].String}, nil
}

// ------------------------------------------------------------------ pilihan layar

func (g *Gudang) daftarPilihan(ctx context.Context, apa, q string, args ...any) ([]models.Pilihan, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: %s: %w", apa, err)
	}
	defer rows.Close()
	var out []models.Pilihan
	for rows.Next() {
		var n, l sql.NullString
		if err := rows.Scan(&n, &l); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", apa, err)
		}
		out = append(out, models.Pilihan{Nilai: n.String, Label: l.String})
	}
	return out, rows.Err()
}

// DaftarMataUang = RD `BrowseCurrencyTreatyIn_RD`: filter B `.Currency !=
// "ITL"` (AC 54); filter A ber-parameter kosong diabaikan.
func (g *Gudang) DaftarMataUang(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelMataUang)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar mata uang",
		fmt.Sprintf(`SELECT TO_CHAR(ID), CURRENCY FROM %s WHERE CURRENCY <> :1 ORDER BY CURRENCY`, t), "ITL")
}

// DaftarMO = RD `BrowseMarketingOfficer_RD`: filter A `.MOStatus = 1`.
func (g *Gudang) DaftarMO(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelMO)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar marketing officer",
		fmt.Sprintf(`SELECT TO_CHAR(ID), CLIENTNAME FROM %s WHERE MOSTATUS = :1 ORDER BY UPPER(CLIENTNAME), ID`, t), "1")
}

// DaftarJenisSpreading = RDB `SelectSpreadingTreatyInProduction`
// (`InputPolicyTreatyInPre_Act` langkah 7-8 -> `ListSpreading`).
func (g *Gudang) DaftarJenisSpreading(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelJenisReas)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar spreading",
		fmt.Sprintf(`SELECT TO_CHAR(ID), NOTE FROM %s
		  WHERE NOTE LIKE :1 OR NOTE = :2 OR NOTE = :3 OR NOTE = :4 OR NOTE = :5`, t),
		"%TRT%", "ORS", "FAC-OUT", "QS (OR)", "QS (R/I)")
}

// DaftarJenisReas = RD `BrowseReinsuranceType_RD` (dropdown TreatyType layar
// atasan); keempat filternya ber-parameter, tidak diisi layar -> diabaikan.
func (g *Gudang) DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error) {
	t, err := g.nama(tabelJenisReas)
	if err != nil {
		return nil, err
	}
	return g.daftarPilihan(ctx, "membaca daftar jenis reasuransi",
		fmt.Sprintf(`SELECT TO_CHAR(ID), NOTE FROM %s ORDER BY ID`, t))
}

// ------------------------------------------------------------------ duplikat

// PolisSerupa = RDB `TreatyRealizationCheckDuplicate` atas
// `POOLDATA.TREATYINPRODUCTION`, dengan PEMETAAN PARAMETER APA ADANYA:
//
//	NOOFFER        <- PolicyTreatyIn.NoOffer
//	BEGINDATE      <- StartDate (YYYYMMDD)   ENDDATE <- EndDate
//	INSUREDNAME    <- PolicyTreatyIn.CedingCoName   (⚠️ bukan InsuredName)
//	UW_YEAR        <- TreatyYear
//	CURR_ID        <- PolicyTreatyIn.Currency       (⚠️ nama, bukan ID)
//	BALANCE_DUE_TO <- BalanceDueTo
//
// Kedua keanehan itu milik rule asli dan ditiru - menukarnya mengubah
// polis mana yang dianggap ganda.
func (g *Gudang) PolisSerupa(ctx context.Context, h *models.Halaman) ([]string, error) {
	t, err := g.nama(tabelProduksi)
	if err != nil {
		return nil, err
	}
	p := func(m string) string { return h.Ambil(models.HalamanPolis + "." + m) }
	awal, err1 := utils.ParseTanggal(strings.TrimSpace(p("StartDate")))
	akhir, err2 := utils.ParseTanggal(strings.TrimSpace(p("EndDate")))
	if err1 != nil || err2 != nil {
		return nil, nil // tanggal kosong: RDB membandingkan dengan NULL -> nol baris
	}
	// `BALANCE_DUE_TO = REPLACE({InputData.Totaltsi},',','.')` - Totaltsi kosong
	// menjadi NULL di Oracle -> nol baris (bukan dibandingkan dengan 0).
	if strings.TrimSpace(p("BalanceDueTo")) == "" {
		return nil, nil
	}
	saldo, err := models.AngkaTeks("BalanceDueTo", p("BalanceDueTo"))
	if err != nil {
		return nil, err
	}
	koef, skala := pecahAngka(saldo)
	q := fmt.Sprintf(`SELECT NOPOLIS FROM %s
	  WHERE NOOFFER = :1 AND BEGINDATE = TO_DATE(:2, 'YYYYMMDD') AND ENDDATE = TO_DATE(:3, 'YYYYMMDD')
	    AND INSUREDNAME = :4 AND UW_YEAR = :5 AND CURR_ID = :6
	    AND BALANCE_DUE_TO = (TO_NUMBER(:7) / POWER(10, :8))`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, p("NoOffer"), awal.Format("20060102"), akhir.Format("20060102"),
		p("CedingCoName"), p("TreatyYear"), p("Currency"), koef, skala)
	if err != nil {
		return nil, fmt.Errorf("repository: memeriksa polis serupa: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var n sql.NullString
		if err := rows.Scan(&n); err != nil {
			return nil, fmt.Errorf("repository: memeriksa polis serupa: %w", err)
		}
		out = append(out, n.String)
	}
	return out, rows.Err()
}
