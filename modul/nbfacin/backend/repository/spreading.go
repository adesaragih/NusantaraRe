package repository

// Tab Spreading kasus FIRE (tiket 48): baca / tulis % Share RNM (T_GENERAL_POLIS.PERCENT_SHARE), TSI / premi Nusantara Re
// per coverage (T_COVERAGELIST.TSI_NUSANTARA_RE / PREMI_NUSANTARA_RE) dan T_SPREADINGLIST (migrasi 198), plus master
// treaty yang dibaca rantai Pega. Rumus TIDAK di sini (services/spreading.go).
//
// Master `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\RDBList\` dan `DDL\` (05-10-2026):
//   - GetTreatyName_SQL (PROPORTIONALARRG TREATY LIMIT × TREATYBUSINESS bizcode aktif × TREATYCONTRACT tanggal), PERSIS
//     kecuali placeholder ber-bind (bizcode dua kali, tanggal dua kali) dan kolom CARI4 / alias yang tidak dipakai.
//   - KAPASITAS_TREATY (Obj-Browse GetKapasitasTreaty: LimitBtmIDR <= TSI <= MaxLimitIDR, StartDate <= begin <= EndDate,
//     urut MaxLimitIDR naik, baris pertama). Pemetaan properti -> kolom (.LimitBtmIDR -> LIMIT_BOTTOM_IDR, .MaxLimitIDR ->
//     MAX_LIMIT_IDR, .TreatyNameQS -> TREATYNAME_QS, ...) `[dugaan]` menurut kemiripan nama (tidak ada rule class di korpus).
//   - GetCurrencyToIDR_SQL (TOIDR TREATYEXCHANGEYEARLY terbaru, urut STARTDATE teks turun) dan GetKursLimitSpreading_SQL
//     (TOIDR yang STARTDATE..ENDDATE-nya memuat Begin date; kedua cabang CASE asal identik -> satu subkueri).
//     ⚠️ STARTDATE / ENDDATE / TOIDR VARCHAR2(255): pembandingan teks seperti asal; bentuk tanggalnya `belum terverifikasi`.
//   - GetDataCurrencyByName_SQL: CURRENCY.ID ber-CURRENCY = nama (pxResults(1); di sini urut ID).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelSpreadingList - T_SPREADINGLIST (migrasi 198).
	TabelSpreadingList    = "T_SPREADINGLIST"
	sequenceSpreadingList = "SEQ_T_SPREADINGLIST"
	// indukSpreading, jalurSpreading - isi PARENT_TABLE / SRC_PATH (pola loader).
	indukSpreading = TabelCoverageList
	jalurSpreading = "LocationList/Property/PropertyItemList/CoverageList/SpreadingList"
	// Master treaty warisan (dibaca saja).
	TabelProportionalArrg     = "PROPORTIONALARRG"
	TabelReinsuranceType      = "REINSURANCETYPE"
	TabelTreatyBusiness       = "TREATYBUSINESS"
	TabelTreatyContract       = "TREATYCONTRACT"
	TabelKapasitasTreaty      = "KAPASITAS_TREATY"
	TabelTreatyExchangeYearly = "TREATYEXCHANGEYEARLY"
	TabelMataUang             = "CURRENCY"
	// mataUangKosong - CURRENCY_CODE bila item tanpa mata uang (DEFAULT rancangan).
	mataUangKosong = "UNKNOWN"
)

// PenyimpanSpreading - data tab Spreading dan master treaty.
type PenyimpanSpreading interface {
	// BacaSpreading - % Share RNM, Begin date, dan pohon lokasi -> item -> coverage -> spreading case `id` (urut SEQ_NO).
	// tx nil = di luar transaksi. ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaSpreading(ctx context.Context, tx *db.Tx, id string) (models.KasusSpreading, error)
	// TulisSpreading - TSI / premi NR tiap coverage (ber-ID) dan SELURUH spreading case diganti; % Share RNM ditulis
	// bila ubahShare. ErrKasusTidakAda bila case tidak ada.
	TulisSpreading(ctx context.Context, tx *db.Tx, id string, k models.KasusSpreading, ubahShare bool) error
	// DaftarTreaty - GetTreatyName_SQL (bizcode, tanggal DD/MM/YYYY).
	DaftarTreaty(ctx context.Context, bizcode, tanggal string) ([]models.TreatySpreading, error)
	// KapasitasTreaty - baris pertama KAPASITAS_TREATY yang memuat tsi pada tanggal DD/MM/YYYY; ada=false bila tidak ada.
	KapasitasTreaty(ctx context.Context, tsi *apd.Decimal, tanggal string) (models.KapasitasTreaty, bool, error)
	// KursTerbaru - GetCurrencyToIDR_SQL (TOIDR terbaru mata uang ber-ID itu); nil bila tidak ada.
	KursTerbaru(ctx context.Context, idMataUang string) (*apd.Decimal, error)
	// KursPada - GetKursLimitSpreading_SQL (TOIDR yang periodenya memuat waktu, teks Pega); nil bila tidak ada.
	KursPada(ctx context.Context, idMataUang, waktu string) (*apd.Decimal, error)
	// IDMataUang - CURRENCY.ID ber-CURRENCY = nama; "" bila tidak ada.
	IDMataUang(ctx context.Context, nama string) (string, error)
	// CatatanJenisTreaty - REINSURANCETYPE.NOTE jenis treaty ber-ID itu (RDB-List PROPORTIONALARRG GetTreatyName,
	// dipakai CalcultePersentageSpeading_Act); "" bila tidak ada.
	CatatanJenisTreaty(ctx context.Context, id string) (string, error)
}

// SpreadingOracle - PenyimpanSpreading atas Oracle.
type SpreadingOracle struct{ db *db.DB }

// NewSpreadingOracle merakit penyimpan spreading.
func NewSpreadingOracle(d *db.DB) *SpreadingOracle { return &SpreadingOracle{db: d} }

// kueri - *db.DB atau *db.Tx.
type kueri interface {
	QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row
}

func (r *SpreadingOracle) pakai(tx *db.Tx) kueri {
	if tx != nil {
		return tx
	}
	return r.db
}

// tabelSpreading - nama berskema tabel yang dipakai.
type tabelSpreading struct{ work, general, loc, prop, risk, item, cov, spread string }

func (r *SpreadingOracle) tabel() (tabelSpreading, error) {
	var t tabelSpreading
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelWorkPolis, &t.work}, {TabelGeneralPolis, &t.general}, {TabelLocationList, &t.loc}, {TabelProperty, &t.prop},
		{TabelRiskLocation, &t.risk}, {TabelPropertyItemList, &t.item}, {TabelCoverageList, &t.cov}, {TabelSpreadingList, &t.spread}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return t, err
		}
		*x.ke = q
	}
	return t, nil
}

func sqlBacaGeneralSpreading(general string) string {
	return "SELECT START_DATE_TIME, " + angkaKeluar("PERCENT_SHARE") + " FROM " + general + " WHERE ID = :1"
}

// kolomPohonSpreading - lokasi (ID, nomor, nama, alamat risiko, top risk), item, coverage; urut sqlPohonSpreading.
var kolomPohonSpreading = []string{
	"TO_CHAR(l.ID)", "p.OBJECT_NO", "p.OBJECT_NAME", "r.ASM_ADDRESS", "p.IS_TOP_RISK",
	"TO_CHAR(i.ID)", "i.ITEM_TYPE", "i.CURRENCY", angkaKeluar("i.TSI_OBJECT_ITEM"), angkaKeluar("i.TOTAL_GROSS_PREMI"),
	"TO_CHAR(v.ID)", "v.OLDID", "v.COVERAGE_BASIS", "v.COVERAGE_NOTE", angkaKeluar("v.RATE"), angkaKeluar("v.PREMIUM"),
	angkaKeluar("v.DISCOUNT"), angkaKeluar("v.TSI"), angkaKeluar("v.TSI_LIABILITY"), angkaKeluar("v.LIMITOF_LIABILITY"),
	"v.FIRST_LOSS", angkaKeluar("v.TSI_NUSANTARA_RE"), angkaKeluar("v.PREMI_NUSANTARA_RE"),
}

// sqlPohonSpreading - satu baris per coverage (lokasi / item tanpa anak tetap muncul), urut lokasi, item, coverage.
func sqlPohonSpreading(t tabelSpreading) string {
	return "SELECT " + strings.Join(kolomPohonSpreading, ", ") + `
FROM ` + t.loc + ` l
LEFT JOIN ` + t.prop + ` p ON p.PARENT_ID = l.ID
LEFT JOIN ` + t.risk + ` r ON r.PARENT_ID = p.ID
LEFT JOIN ` + t.item + ` i ON i.PARENT_ID = p.ID
LEFT JOIN ` + t.cov + ` v ON v.PARENT_ID = i.ID AND ` + syaratIndukCoverage + `
WHERE l.PARENT_ID = :1
ORDER BY l.SEQ_NO, i.SEQ_NO, v.SEQ_NO`
}

// sqlCoverageKasus - subkueri ID coverage berinduk item case :1.
func sqlCoverageKasus(t tabelSpreading) string {
	return "SELECT v.ID FROM " + t.cov + " v JOIN " + t.item + " i ON i.ID = v.PARENT_ID JOIN " + t.prop +
		" p ON p.ID = i.PARENT_ID JOIN " + t.loc + " l ON l.ID = p.PARENT_ID WHERE " + syaratIndukCoverage + " AND l.PARENT_ID = :1"
}

func sqlBacaBarisSpreading(t tabelSpreading) string {
	return "SELECT TO_CHAR(s.PARENT_ID), s.TREATY_TYPE, s.TREATY_NAME, " + angkaKeluar("s.SHARE_PERCENTAGE") + ", " +
		angkaKeluar("s.TSI_GROSS_SPREADED") + ", " + angkaKeluar("s.TSI_SPREADED") + ", " + angkaKeluar("s.CLAIM_ESTIMATION") +
		", " + angkaKeluar("s.PREMIUM_SPREADED") + " FROM " + t.spread + " s WHERE s.PARENT_ID IN (" + sqlCoverageKasus(t) +
		") ORDER BY s.PARENT_ID, s.SEQ_NO"
}

// sqlHapusSpreading - spreading seluruh coverage case :1 (juga dipakai tab Object sebelum coverage dihapus).
func sqlHapusSpreading(spread, cov, item, prop, loc string) string {
	return "DELETE FROM " + spread + " WHERE PARENT_ID IN (" + sqlCoverageKasus(tabelSpreading{cov: cov, item: item, prop: prop, loc: loc}) + ")"
}

func sqlUbahShare(general string) string {
	return "UPDATE " + general + " SET PERCENT_SHARE = " + angkaMasuk(":1") + " WHERE ID = :2"
}

func sqlUbahNR(cov string) string {
	return "UPDATE " + cov + " SET TSI_NUSANTARA_RE = " + angkaMasuk(":1") + ", PREMI_NUSANTARA_RE = " + angkaMasuk(":2") +
		" WHERE ID = :3"
}

func sqlSisipSpreading(spread string) string {
	return "INSERT INTO " + spread + " (ID, PARENT_ID, PARENT_TABLE, SRC_PATH, SEQ_NO, ROW_UID, TREATY_TYPE, TREATY_NAME," +
		" SHARE_PERCENTAGE, TSI_SPREADED, TSI_GROSS_SPREADED, PREMIUM_SPREADED, CLAIM_ESTIMATION, CURRENCY_CODE) VALUES (:1, :2," +
		" :3, :4, :5, :6, :7, :8, " + angkaMasuk(":9") + ", " + angkaMasuk(":10") + ", " + angkaMasuk(":11") + ", " +
		angkaMasuk(":12") + ", " + angkaMasuk(":13") + ", :14)"
}

// BacaSpreading - lihat PenyimpanSpreading.
func (r *SpreadingOracle) BacaSpreading(ctx context.Context, tx *db.Tx, id string) (models.KasusSpreading, error) {
	var k models.KasusSpreading
	t, err := r.tabel()
	if err != nil {
		return k, err
	}
	q := r.pakai(tx)
	var n int
	if err := q.QueryRowContext(ctx, sqlAdaKasus(t.work), id, LiniFacIn).Scan(&n); err != nil {
		return k, fmt.Errorf("repository: cek case NB: %w", err)
	}
	if n == 0 {
		return k, ErrKasusTidakAda
	}
	var mulai, share sql.NullString
	switch err := q.QueryRowContext(ctx, sqlBacaGeneralSpreading(t.general), id).Scan(&mulai, &share); {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		return k, fmt.Errorf("repository: baca %s: %w", TabelGeneralPolis, err)
	default:
		k.StartDateTime = mulai.String
		if k.PercentShare, err = bacaDesimal(TabelGeneralPolis+" "+id, "PERCENT_SHARE", &share); err != nil {
			return k, err
		}
	}
	if k.Lokasi, err = r.bacaPohon(ctx, q, t, id); err != nil {
		return k, err
	}
	return k, r.bacaBaris(ctx, q, t, id, k.Lokasi)
}

// bacaPohon - lokasi -> item -> coverage dari baris rata sqlPohonSpreading.
func (r *SpreadingOracle) bacaPohon(ctx context.Context, q kueri, t tabelSpreading, id string) ([]models.LokasiSpreading, error) {
	baris, err := q.QueryContext(ctx, sqlPohonSpreading(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca pohon spreading: %w", err)
	}
	defer baris.Close()
	lokasi := []models.LokasiSpreading{}
	var lokTerakhir, itemTerakhir string
	for baris.Next() {
		v := make([]sql.NullString, len(kolomPohonSpreading))
		ptr := make([]any, len(v))
		for i := range v {
			ptr[i] = &v[i]
		}
		if err := baris.Scan(ptr...); err != nil {
			return nil, fmt.Errorf("repository: pohon spreading: %w", err)
		}
		if len(lokasi) == 0 || v[0].String != lokTerakhir {
			lokTerakhir, itemTerakhir = v[0].String, ""
			lokasi = append(lokasi, models.LokasiSpreading{ObjectNo: v[1].String, ObjectName: v[2].String, Location: v[3].String,
				IsTopRisk: v[4].String == teksBenar, Items: []models.ItemSpreading{}})
		}
		l := &lokasi[len(lokasi)-1]
		if !v[5].Valid {
			continue
		}
		if v[5].String != itemTerakhir {
			itemTerakhir = v[5].String
			it := models.ItemSpreading{ItemType: v[6].String, Currency: v[7].String, Coverages: []models.CoverageSpreading{}}
			if err := bacaDesimalKe(TabelPropertyItemList+" "+v[5].String, map[string]*sql.NullString{"tsi": &v[8], "gross": &v[9]},
				kolomDesimal{"tsi", "TSI_OBJECT_ITEM", &it.TSI}, kolomDesimal{"gross", "TOTAL_GROSS_PREMI", &it.TotalGrossPremi}); err != nil {
				return nil, err
			}
			l.Items = append(l.Items, it)
		}
		if !v[10].Valid {
			continue
		}
		c := models.CoverageSpreading{ID: v[10].String, OldID: v[11].String, CoverageBasis: v[12].String, CoverageNote: v[13].String,
			Spreading: []models.BarisSpreading{}}
		teks := map[string]*sql.NullString{"rate": &v[14], "premi": &v[15], "diskon": &v[16], "tsi": &v[17], "liab": &v[18],
			"lol": &v[19], "fl": &v[20], "tsinr": &v[21], "preminr": &v[22]}
		if err := bacaDesimalKe(TabelCoverageList+" "+c.ID, teks, kolomDesimal{"rate", "RATE", &c.Rate},
			kolomDesimal{"premi", "PREMIUM", &c.Premium}, kolomDesimal{"diskon", "DISCOUNT", &c.Discount},
			kolomDesimal{"tsi", "TSI", &c.TSI}, kolomDesimal{"liab", "TSI_LIABILITY", &c.TSILiability},
			kolomDesimal{"lol", "LIMITOF_LIABILITY", &c.LimitOfLiability}, kolomDesimal{"fl", "FIRST_LOSS", &c.FirstLoss},
			kolomDesimal{"tsinr", "TSI_NUSANTARA_RE", &c.TSINusantaraRe}, kolomDesimal{"preminr", "PREMI_NUSANTARA_RE", &c.PremiNusantaraRe}); err != nil {
			return nil, err
		}
		it := &l.Items[len(l.Items)-1]
		it.Coverages = append(it.Coverages, c)
	}
	return lokasi, baris.Err()
}

// bacaBaris - baris T_SPREADINGLIST dipasang ke coverage-nya menurut ID.
func (r *SpreadingOracle) bacaBaris(ctx context.Context, q kueri, t tabelSpreading, id string, lokasi []models.LokasiSpreading) error {
	cov := map[string]*models.CoverageSpreading{}
	for i := range lokasi {
		for j := range lokasi[i].Items {
			for m := range lokasi[i].Items[j].Coverages {
				c := &lokasi[i].Items[j].Coverages[m]
				cov[c.ID] = c
			}
		}
	}
	baris, err := q.QueryContext(ctx, sqlBacaBarisSpreading(t), id)
	if err != nil {
		return fmt.Errorf("repository: baca %s: %w", TabelSpreadingList, err)
	}
	defer baris.Close()
	for baris.Next() {
		var induk, tipe, nama, share, gross, tsi, klaim, premi sql.NullString
		if err := baris.Scan(&induk, &tipe, &nama, &share, &gross, &tsi, &klaim, &premi); err != nil {
			return fmt.Errorf("repository: %s: %w", TabelSpreadingList, err)
		}
		s := models.BarisSpreading{TreatyType: tipe.String, TreatyName: nama.String}
		teks := map[string]*sql.NullString{"share": &share, "gross": &gross, "tsi": &tsi, "klaim": &klaim, "premi": &premi}
		if err := bacaDesimalKe(TabelSpreadingList+" "+induk.String, teks, kolomDesimal{"share", "SHARE_PERCENTAGE", &s.SharePercentage},
			kolomDesimal{"gross", "TSI_GROSS_SPREADED", &s.TSIGrossSpreaded}, kolomDesimal{"tsi", "TSI_SPREADED", &s.TSISpreaded},
			kolomDesimal{"klaim", "CLAIM_ESTIMATION", &s.ClaimEstimation}, kolomDesimal{"premi", "PREMIUM_SPREADED", &s.PremiumSpreaded}); err != nil {
			return err
		}
		if c, ada := cov[induk.String]; ada {
			c.Spreading = append(c.Spreading, s)
		}
	}
	return baris.Err()
}

// TulisSpreading - lihat PenyimpanSpreading. Urutan: sentuh case (kunci + 404), pastikan General (+ % Share RNM bila
// ubahShare), TSI / premi NR tiap coverage, hapus lalu sisip ulang spreading (SEQ_NO 1..n per coverage).
func (r *SpreadingOracle) TulisSpreading(ctx context.Context, tx *db.Tx, id string, k models.KasusSpreading, ubahShare bool) error {
	t, err := r.tabel()
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, sqlSentuhCase(t.work), "menyentuh "+TabelWorkPolis, id, LiniFacIn)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrKasusTidakAda
	}
	if ubahShare {
		if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, sqlUbahShare(t.general), "menulis % Share RNM", ikatDesimal(k.PercentShare), id); err != nil {
			return err
		}
	}
	if _, err := jalankan(ctx, tx, sqlHapusSpreading(t.spread, t.cov, t.item, t.prop, t.loc), "menghapus spreading", id); err != nil {
		return err
	}
	for _, l := range k.Lokasi {
		for _, it := range l.Items {
			mataUang := it.Currency
			if mataUang == "" {
				mataUang = mataUangKosong
			}
			for _, c := range it.Coverages {
				if c.ID == "" {
					continue
				}
				if _, err := jalankan(ctx, tx, sqlUbahNR(t.cov), "menulis TSI / premi Nusantara Re", ikatDesimal(c.TSINusantaraRe),
					ikatDesimal(c.PremiNusantaraRe), c.ID); err != nil {
					return err
				}
				for n, s := range c.Spreading {
					idBaris, err := r.db.NomorBerikut(ctx, tx, sequenceSpreadingList)
					if err != nil {
						return err
					}
					uid, err := uidAcak()
					if err != nil {
						return err
					}
					if _, err := jalankan(ctx, tx, sqlSisipSpreading(t.spread), "menyisip "+TabelSpreadingList, idBaris, c.ID, indukSpreading,
						jalurSpreading, n+1, uid, db.KosongJadiNil(s.TreatyType), db.KosongJadiNil(s.TreatyName), ikatDesimal(s.SharePercentage),
						ikatDesimal(s.TSISpreaded), ikatDesimal(s.TSIGrossSpreaded), ikatDesimal(s.PremiumSpreaded),
						ikatDesimal(s.ClaimEstimation), mataUang); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// sqlDaftarTreaty - GetTreatyName_SQL (lihat kepala berkas); :1 / :2 bizcode, :3 / :4 tanggal DD/MM/YYYY.
func sqlDaftarTreaty(pa, rt, tb, tc string) string {
	return "SELECT (SELECT c.NOURUT FROM " + rt + " c WHERE c.ID = a.REINSTYPEID) AS CARI5, a.REINSTYPEID, a.REINSTYPENAME, a.RP FROM " +
		pa + ` a
WHERE a.TREATYDESCNAME = 'TREATY LIMIT'
  AND a.TREATYGROUPID IN (SELECT b.TREATYGROUPID FROM ` + tb + ` b WHERE b.BIZCODE = :1 AND b.ISACTIVE = '1')
  AND a.REINSTYPEID IN (SELECT b.REINSTYPEID FROM ` + tb + ` b WHERE b.BIZCODE = :2 AND b.ISACTIVE = '1')
  AND a.TREATYYEARID IN (SELECT t.IDTREATYYEAR FROM ` + tc + ` t WHERE TO_DATE(:3, 'DD/MM/RRRR') BETWEEN t.TREATYSTARTDATE AND t.TREATYENDDATE)
  AND a.REINSTYPEID IN (SELECT t.REINSTYPEID FROM ` + tc + ` t WHERE TO_DATE(:4, 'DD/MM/RRRR') BETWEEN t.TREATYSTARTDATE AND t.TREATYENDDATE)
ORDER BY CARI5 ASC`
}

// sqlKapasitasTreaty - Obj-Browse GetKapasitasTreaty: :1 / :2 TSI, :3 / :4 tanggal DD/MM/YYYY; baris pertama urut
// MAX_LIMIT_IDR naik (Sort A pada baris B).
func sqlKapasitasTreaty(kt string) string {
	return "SELECT TREATYNAME_QS, IDTREATY_QS, " + angkaKeluar("MAX_LIMIT_QS_IDR") + ", TREATYNAME_SPL, IDTREATY_SPL, " +
		angkaKeluar("MAX_LIMIT_SPL_IDR") + " FROM " + kt + " WHERE LIMIT_BOTTOM_IDR <= " + angkaMasuk(":1") + " AND MAX_LIMIT_IDR >= " +
		angkaMasuk(":2") + " AND STARTDATE <= TO_DATE(:3, 'DD/MM/YYYY') AND ENDDATE >= TO_DATE(:4, 'DD/MM/YYYY')" +
		" ORDER BY MAX_LIMIT_IDR ASC FETCH FIRST 1 ROWS ONLY"
}

// sqlKursTerbaru - GetCurrencyToIDR_SQL.
func sqlKursTerbaru(tey string) string {
	return "SELECT TOIDR FROM (SELECT x.TOIDR FROM " + tey + " x WHERE x.IDCURRENCY = :1 ORDER BY x.STARTDATE DESC) WHERE ROWNUM = 1"
}

// sqlKursPada - GetKursLimitSpreading_SQL (dua cabang CASE asal identik).
func sqlKursPada(tey string) string {
	return "SELECT TOIDR FROM (SELECT x.TOIDR FROM " + tey + " x WHERE x.IDCURRENCY = :1 AND :2 BETWEEN x.STARTDATE AND x.ENDDATE" +
		" ORDER BY x.STARTDATE DESC) WHERE ROWNUM = 1"
}

// sqlIDMataUang - GetDataCurrencyByName_SQL (ID saja).
func sqlIDMataUang(cur string) string {
	return "SELECT TO_CHAR(c.ID) FROM " + cur + " c WHERE c.CURRENCY = :1 ORDER BY c.ID FETCH FIRST 1 ROWS ONLY"
}

// DaftarTreaty - lihat PenyimpanSpreading.
func (r *SpreadingOracle) DaftarTreaty(ctx context.Context, bizcode, tanggal string) ([]models.TreatySpreading, error) {
	var nama [4]string
	for i, t := range []string{TabelProportionalArrg, TabelReinsuranceType, TabelTreatyBusiness, TabelTreatyContract} {
		q, err := r.db.Qualify(t)
		if err != nil {
			return nil, err
		}
		nama[i] = q
	}
	baris, err := r.db.QueryContext(ctx, sqlDaftarTreaty(nama[0], nama[1], nama[2], nama[3]), bizcode, bizcode, tanggal, tanggal)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelProportionalArrg, err)
	}
	defer baris.Close()
	hasil := []models.TreatySpreading{}
	for baris.Next() {
		var urut, id, nm, rp sql.NullString
		if err := baris.Scan(&urut, &id, &nm, &rp); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelProportionalArrg, err)
		}
		hasil = append(hasil, models.TreatySpreading{ID: id.String, Name: nm.String, Limit: rp.String})
	}
	return hasil, baris.Err()
}

// KapasitasTreaty - lihat PenyimpanSpreading.
func (r *SpreadingOracle) KapasitasTreaty(ctx context.Context, tsi *apd.Decimal, tanggal string) (models.KapasitasTreaty, bool, error) {
	var k models.KapasitasTreaty
	kt, err := r.db.Qualify(TabelKapasitasTreaty)
	if err != nil {
		return k, false, err
	}
	nilai := ikatDesimal(tsi)
	var qs, spl sql.NullString
	var nqs, iqs, nspl, ispl sql.NullString
	err = r.db.QueryRowContext(ctx, sqlKapasitasTreaty(kt), nilai, nilai, tanggal, tanggal).Scan(&nqs, &iqs, &qs, &nspl, &ispl, &spl)
	if errors.Is(err, sql.ErrNoRows) {
		return k, false, nil
	}
	if err != nil {
		return k, false, fmt.Errorf("repository: baca %s: %w", TabelKapasitasTreaty, err)
	}
	k = models.KapasitasTreaty{TreatyNameQS: nqs.String, IDTreatyQS: iqs.String, TreatyNameSPL: nspl.String, IDTreatySPL: ispl.String}
	if err := bacaDesimalKe(TabelKapasitasTreaty, map[string]*sql.NullString{"qs": &qs, "spl": &spl},
		kolomDesimal{"qs", "MAX_LIMIT_QS_IDR", &k.MaxLimitQSIDR}, kolomDesimal{"spl", "MAX_LIMIT_SPL_IDR", &k.MaxLimitSPLIDR}); err != nil {
		return k, false, err
	}
	return k, true, nil
}

// bacaKurs - satu TOIDR (teks VARCHAR2) -> desimal; tidak ada baris / kosong -> nil.
func (r *SpreadingOracle) bacaKurs(ctx context.Context, q string, arg ...any) (*apd.Decimal, error) {
	var v sql.NullString
	err := r.db.QueryRowContext(ctx, q, arg...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelTreatyExchangeYearly, err)
	}
	return bacaDesimal(TabelTreatyExchangeYearly, "TOIDR", &v)
}

// KursTerbaru - lihat PenyimpanSpreading.
func (r *SpreadingOracle) KursTerbaru(ctx context.Context, idMataUang string) (*apd.Decimal, error) {
	tey, err := r.db.Qualify(TabelTreatyExchangeYearly)
	if err != nil {
		return nil, err
	}
	return r.bacaKurs(ctx, sqlKursTerbaru(tey), idMataUang)
}

// KursPada - lihat PenyimpanSpreading.
func (r *SpreadingOracle) KursPada(ctx context.Context, idMataUang, waktu string) (*apd.Decimal, error) {
	tey, err := r.db.Qualify(TabelTreatyExchangeYearly)
	if err != nil {
		return nil, err
	}
	return r.bacaKurs(ctx, sqlKursPada(tey), idMataUang, waktu)
}

// sqlCatatanJenisTreaty - `SELECT ID AS CARI1, NOTE AS CARI2 FROM REINSURANCETYPE WHERE ID = {InputData.CARI3}` (NOTE
// saja). ID dibandingkan sebagai teks: treatyType datang dari layar, teks bukan-angka tidak boleh memicu ORA-01722.
func sqlCatatanJenisTreaty(rt string) string {
	return "SELECT NOTE FROM " + rt + " WHERE TO_CHAR(ID) = :1"
}

// CatatanJenisTreaty - lihat PenyimpanSpreading.
func (r *SpreadingOracle) CatatanJenisTreaty(ctx context.Context, id string) (string, error) {
	rt, err := r.db.Qualify(TabelReinsuranceType)
	if err != nil {
		return "", err
	}
	var note sql.NullString
	err = r.db.QueryRowContext(ctx, sqlCatatanJenisTreaty(rt), id).Scan(&note)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: baca %s: %w", TabelReinsuranceType, err)
	}
	return note.String, nil
}

// IDMataUang - lihat PenyimpanSpreading.
func (r *SpreadingOracle) IDMataUang(ctx context.Context, nama string) (string, error) {
	cur, err := r.db.Qualify(TabelMataUang)
	if err != nil {
		return "", err
	}
	var id sql.NullString
	err = r.db.QueryRowContext(ctx, sqlIDMataUang(cur), nama).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: baca %s: %w", TabelMataUang, err)
	}
	return id.String, nil
}
