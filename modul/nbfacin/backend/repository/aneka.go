package repository

// Tab Object keluarga Aneka (tiket 73 P1, migrasi 214). Lokasi = baris T_LOCATIONLIST / T_PROPERTY / T_RISKLOCATION yang SAMA
// dengan tab Object FIRE (186); di bawahnya okupasi Aneka T_OCCUPATIONLIST ber-PARENT_TABLE 'T_RISKLOCATION' (jalur rancangan
// `LocationList/Property/RiskLocation/OccupationList`, BUKAN okupasi FIRE ber-induk T_PROPERTY), objek T_ANEKALIST
// (PARENT_TABLE 'T_OCCUPATIONLIST', atau 'T_RISKLOCATION' untuk Golf), dan empat halaman anak satu-per-objek.
//
// GantiObjekAneka MENGGANTI UTUH seluruh lokasi case (juga sisa anak FIRE, lewat sqlHapusObjek tab FIRE yang tidak diubah).
// Coverage Aneka (T_COVERAGELIST ber-PARENT_TABLE 'T_ANEKALIST', tanpa FK - pola A140) TIDAK dihapus untuk objek yang tetap ada:
// barisnya dipindah induk ke ID objek baru menurut posisi (lokasi, okupasi / Golf, objek); coverage objek yang hilang dihapus
// beserta deductible dan spreading-nya. Seluruh nilai parameter terikat.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	TabelAnekaList        = "T_ANEKALIST"
	TabelMarineHull       = "T_MARINEHULL"
	TabelAviationHull     = "T_AVIATIONHULL"
	TabelVehicleHE        = "T_VEHICLEHE"
	TabelMemorandumAccept = "T_MEMORANDUMACCEPT"
	sequenceAnekaList     = "SEQ_T_ANEKALIST"
	sequenceMarineHull    = "SEQ_T_MARINEHULL"
	sequenceAviationHull  = "SEQ_T_AVIATIONHULL"
	sequenceVehicleHE     = "SEQ_T_VEHICLEHE"
	sequenceMemorandum    = "SEQ_T_MEMORANDUMACCEPT"
	// Isi PARENT_TABLE / SRC_PATH = jalur rancangan (loader flatten.go menulis yang sama).
	indukOkupasiAneka = TabelRiskLocation
	jalurOkupasiAneka = "LocationList/Property/RiskLocation/OccupationList"
	indukAnekaOkupasi = TabelOccupationList
	jalurAnekaOkupasi = "LocationList/Property/RiskLocation/OccupationList/AnekaList"
	indukAnekaGolf    = TabelRiskLocation
	jalurAnekaGolf    = "LocationList/Property/RiskLocation/AnekaList"
	indukCoverageAnek = TabelAnekaList
)

// PenyimpanObjekAneka - baca / ganti tab Object Aneka case.
type PenyimpanObjekAneka interface {
	// BacaObjekAneka - ErrKasusTidakAda bila case tidak ada / bukan LINI Fac In.
	BacaObjekAneka(ctx context.Context, id string) ([]models.LokasiAneka, error)
	// GantiObjekAneka - ganti seluruh lokasi case di dalam transaksi pemanggil (lihat kepala berkas); proRate = % Pro Rate CASE
	// (T_GENERAL_POLIS.PRO_RATE_PERCENT, migrasi 255; nil = NULL).
	GantiObjekAneka(ctx context.Context, tx *db.Tx, id string, lokasi []models.LokasiAneka, proRate *apd.Decimal) error
	// ProRateKasus - T_GENERAL_POLIS.PRO_RATE_PERCENT case (migrasi 255); baris / nilai tidak ada -> nil.
	ProRateKasus(ctx context.Context, id string) (*apd.Decimal, error)
	// PercentShareKasus - T_GENERAL_POLIS.PERCENT_SHARE case (tiket 73 P2, medan RNM coverage); baris / nilai tidak ada -> nil.
	PercentShareKasus(ctx context.Context, id string) (*apd.Decimal, error)
	// UbahPercentShare - T_GENERAL_POLIS.PERCENT_SHARE case di transaksi pemanggil (tiket 73 P3, Save tab Spreading Aneka; baris
	// General dipastikan ada).
	UbahPercentShare(ctx context.Context, tx *db.Tx, id string, ps *apd.Decimal) error
}

type tabelAneka struct {
	tabelObjek
	aneka, marine, aviasi, he, memo string
}

func (r *ObjekOracle) tabelAneka() (tabelAneka, error) {
	t, err := r.tabel()
	if err != nil {
		return tabelAneka{}, err
	}
	a := tabelAneka{tabelObjek: t}
	for _, x := range []struct {
		nama string
		ke   *string
	}{{TabelAnekaList, &a.aneka}, {TabelMarineHull, &a.marine}, {TabelAviationHull, &a.aviasi}, {TabelVehicleHE, &a.he},
		{TabelMemorandumAccept, &a.memo}} {
		q, err := r.db.Qualify(x.nama)
		if err != nil {
			return a, err
		}
		*x.ke = q
	}
	return a, nil
}

type kolomAneka = kolomData[models.ObjekAneka]

// kolomAnekaData - kolom data T_ANEKALIST yang ditulis / dibaca (migrasi 214) sesudah kolom sistem; CURRENCY_CODE terpisah.
// Tidak ditulis: IS_OLD_DATA, LOCATION_ID, OCCUPATION_ID, SELECTED_*, FLAG_DELETE (EDM / loader).
var kolomAnekaData = []kolomAneka{
	kt("OBJECT_NAME", func(o *models.ObjekAneka) *string { return &o.ObjectName }),
	kt("SECTION", func(o *models.ObjekAneka) *string { return &o.Section }),
	kt("NO_KONTRAK", func(o *models.ObjekAneka) *string { return &o.NoKontrak }),
	kt("START_DATE", func(o *models.ObjekAneka) *string { return &o.StartDate }),
	kt("END_DATE", func(o *models.ObjekAneka) *string { return &o.EndDate }),
	kt("NO_KTP", func(o *models.ObjekAneka) *string { return &o.NoKtp }),
	kt("NO_NPWP", func(o *models.ObjekAneka) *string { return &o.NoNpwp }),
	ka("TSI", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.TSI }),
	ka("TSI_OLD", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.TSIOld }),
	kt("DRAFT_WORDING_ID", func(o *models.ObjekAneka) *string { return &o.DraftWordingID }),
	ka("MAINTENANCE_PERIOD", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.MaintenancePeriod }),
	kt("OTHERS", func(o *models.ObjekAneka) *string { return &o.Others }),
	kt("YEAR", func(o *models.ObjekAneka) *string { return &o.Year }),
	kt("QUANTITY", func(o *models.ObjekAneka) *string { return &o.Quantity }),
	kt("UNIT_QUANTITY", func(o *models.ObjekAneka) *string { return &o.UnitQuantity }),
	kt("SERIAL_NO", func(o *models.ObjekAneka) *string { return &o.SerialNo }),
	kt("BUILT_IN", func(o *models.ObjekAneka) *string { return &o.BuiltIn }),
	kt("FLAG", func(o *models.ObjekAneka) *string { return &o.Flag }),
	kt("GEOGRAPHICAL_LIMITS", func(o *models.ObjekAneka) *string { return &o.GeographicalLimits }),
	kt("INVOICE_NO", func(o *models.ObjekAneka) *string { return &o.InvoiceNo }),
	ka("NO_OF_TREE", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.NoOfTree }),
	ka("AREA_HECTAR", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.AreaHectar }),
	ka("LIMIT_IDEMNITY_DAY", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.LimitIdemnityDay }),
	ka("LIMIT_IDEMNITY_MONTH", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.LimitIdemnityMonth }),
	ka("LIMIT_OF_LIABILITY", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.LimitOfLiability }),
	kt("ALM_RISK_ID", func(o *models.ObjekAneka) *string { return &o.AlmRiskID }),
	kt("LOAN_ACCOUNT_NUMBER", func(o *models.ObjekAneka) *string { return &o.Briguna.LoanAccountNumber }),
	kt("LOAN_NAME", func(o *models.ObjekAneka) *string { return &o.Briguna.LoanName }),
	kt("JOB", func(o *models.ObjekAneka) *string { return &o.Briguna.Job }),
	kt("DEPOSIT_ACCOUNT_NO", func(o *models.ObjekAneka) *string { return &o.Briguna.DepositAccountNo }),
	kt("ID_CARD_NO", func(o *models.ObjekAneka) *string { return &o.Briguna.IDCardNo }),
	kt("INSTITUTION", func(o *models.ObjekAneka) *string { return &o.Briguna.Institution }),
	kt("LOAN_TYPE", func(o *models.ObjekAneka) *string { return &o.Briguna.LoanType }),
	kt("DATE_OF_BIRTH", func(o *models.ObjekAneka) *string { return &o.Briguna.DateOfBirth }),
	ka("SALARY", jenisAngka, func(o *models.ObjekAneka) **apd.Decimal { return &o.Briguna.Salary }),
}

var kolomMarineHull = []kolomData[models.MarineHull]{
	kt("TYPE_OF_VESSEL", func(m *models.MarineHull) *string { return &m.TypeOfVessel }),
	kt("NRTRT", func(m *models.MarineHull) *string { return &m.NRTRT }),
	kt("LENGTHM", func(m *models.MarineHull) *string { return &m.LengthM }),
	kt("CLASSIFICATION", func(m *models.MarineHull) *string { return &m.Classification }),
	kt("GRTRT", func(m *models.MarineHull) *string { return &m.GRTRT }),
	kt("BREADTHM", func(m *models.MarineHull) *string { return &m.BreadthM }),
	kt("CONSTRUCTION", func(m *models.MarineHull) *string { return &m.Construction }),
	kt("DWTRT", func(m *models.MarineHull) *string { return &m.DWTRT }),
	kt("DEPTHM", func(m *models.MarineHull) *string { return &m.DepthM }),
}

var kolomAviationHull = []kolomData[models.AviationHull]{
	kt("REGISTRATION_MARKS", func(a *models.AviationHull) *string { return &a.RegistrationMarks }),
	kt("MAX_PASSENGERS", func(a *models.AviationHull) *string { return &a.MaxPassengers }),
	kt("PILOTS", func(a *models.AviationHull) *string { return &a.Pilots }),
	kt("STANDARD_USES", func(a *models.AviationHull) *string { return &a.StandardUses }),
	kt("SPECIAL_USES", func(a *models.AviationHull) *string { return &a.SpecialUses }),
	kt("SPECIAL_RENTAL_USES", func(a *models.AviationHull) *string { return &a.SpecialRentalUses }),
}

var kolomVehicleHE = []kolomData[models.VehicleHE]{
	kt("OBJECT_NAME_HE", func(v *models.VehicleHE) *string { return &v.ObjectNameHE }),
	kt("OBJECT_ID_HE", func(v *models.VehicleHE) *string { return &v.ObjectIDHE }),
	kt("BRAND_NAME", func(v *models.VehicleHE) *string { return &v.BrandName }),
	kt("BRAND", func(v *models.VehicleHE) *string { return &v.Brand }),
	kt("TYPE_NAME", func(v *models.VehicleHE) *string { return &v.TypeName }),
	kt("TYPE", func(v *models.VehicleHE) *string { return &v.Type }),
	kt("CHASSIS_NUMBER", func(v *models.VehicleHE) *string { return &v.ChassisNumber }),
	kt("ENGINE_NUMBER", func(v *models.VehicleHE) *string { return &v.EngineNumber }),
}

var kolomMemorandum = []kolomData[models.MemorandumAccept]{
	kt("SURETY_TYPE", func(m *models.MemorandumAccept) *string { return &m.SuretyType }),
	kt("APPLICATION_NUMBER", func(m *models.MemorandumAccept) *string { return &m.ApplicationNumber }),
	kt("ACCUMULATION", func(m *models.MemorandumAccept) *string { return &m.Accumulation }),
	kt("COLLATERAL", func(m *models.MemorandumAccept) *string { return &m.Collateral }),
	kt("INSURED_SIUP", func(m *models.MemorandumAccept) *string { return &m.InsuredSiup }),
	kt("KEY_PERSON", func(m *models.MemorandumAccept) *string { return &m.KeyPerson }),
	kt("ORIGINAL_RATE", func(m *models.MemorandumAccept) *string { return &m.OriginalRate }),
	kt("PRINCIPALS_OUT_GO", func(m *models.MemorandumAccept) *string { return &m.PrincipalsOutGo }),
	kt("TYPE_OF_WORDING", func(m *models.MemorandumAccept) *string { return &m.TypeOfWording }),
	kt("BUILDING_STATUS", func(m *models.MemorandumAccept) *string { return &m.BuildingStatus }),
	kt("SUPPORTING_DOCUMENT", func(m *models.MemorandumAccept) *string { return &m.SupportingDocument }),
	kt("PRINCIPALS_DATA", func(m *models.MemorandumAccept) *string { return &m.PrincipalsData }),
	kt("CONCLUSIONS_PROPOSALS", func(m *models.MemorandumAccept) *string { return &m.ConclusionsProposals }),
	kt("PRINCIPALS_BACKGROUND", func(m *models.MemorandumAccept) *string { return &m.PrincipalsBackground }),
	kt("SUMMARY_FINANCIAL", func(m *models.MemorandumAccept) *string { return &m.SummaryFinancial }),
	kt("PRINCIPALS_GUARANTEE", func(m *models.MemorandumAccept) *string { return &m.PrincipalsGuarantee }),
}

// sqlRiskKasus / sqlOkupasiAnekaKasus - ID T_RISKLOCATION / okupasi Aneka case (bind `b`).
func sqlRiskKasus(t tabelAneka, b string) string {
	return "SELECT r.ID FROM " + t.risk + " r JOIN " + t.prop + " p ON p.ID = r.PARENT_ID JOIN " + t.loc +
		" l ON l.ID = p.PARENT_ID WHERE l.PARENT_ID = " + b
}

func sqlOkupasiAnekaKasus(t tabelAneka, b string) string {
	return "SELECT o.ID FROM " + t.okupasi + " o WHERE o.PARENT_TABLE = '" + indukOkupasiAneka + "' AND o.PARENT_ID IN (" +
		sqlRiskKasus(t, b) + ")"
}

// sqlAnekaKasus - ID T_ANEKALIST case (:1 dan :2 = ID case).
func sqlAnekaKasus(t tabelAneka) string {
	return "SELECT a.ID FROM " + t.aneka + " a WHERE (a.PARENT_TABLE = '" + indukAnekaOkupasi + "' AND a.PARENT_ID IN (" +
		sqlOkupasiAnekaKasus(t, ":1") + ")) OR (a.PARENT_TABLE = '" + indukAnekaGolf + "' AND a.PARENT_ID IN (" +
		sqlRiskKasus(t, ":2") + "))"
}

var kolomBacaLokasiAneka = []string{"p.OBJECT_NO", "p.IS_TOP_RISK", "p.ROAD_TYPE", "p.ROAD_NAME", "p.BUILDING_NO",
	"r.ASM_ZIP_CODE", "p.COUNTRY", "r.ASM_ADDRESS", "r.ASMRW", "r.ASM_CITY", "r.ASM_DISTRICT", "p.PROVINCE", "p.ALM_RISK_ID",
	"TO_CHAR(r.ID)"}

func sqlBacaLokasiAneka(t tabelAneka) string {
	return "SELECT " + strings.Join(kolomBacaLokasiAneka, ", ") + " FROM " + t.loc + " l LEFT JOIN " + t.prop +
		" p ON p.PARENT_ID = l.ID LEFT JOIN " + t.risk + " r ON r.PARENT_ID = p.ID WHERE l.PARENT_ID = :1 ORDER BY l.SEQ_NO"
}

var kolomBacaOkupasiAneka = []string{"TO_CHAR(o.PARENT_ID)", "TO_CHAR(o.ID)", "o.OCCUPATION_ID", "o.OCCUPATION_NAME"}

func sqlBacaOkupasiAneka(t tabelAneka) string {
	return "SELECT " + strings.Join(kolomBacaOkupasiAneka, ", ") + " FROM " + t.okupasi + " o WHERE o.ID IN (" +
		sqlOkupasiAnekaKasus(t, ":1") + ") ORDER BY o.PARENT_ID, o.SEQ_NO"
}

// kolomBacaAneka - induk, lalu data T_ANEKALIST dan keempat anak (LEFT JOIN, UNIQUE PARENT_ID).
func kolomBacaAneka() []string {
	k := []string{"a.PARENT_TABLE", "TO_CHAR(a.PARENT_ID)", "TO_CHAR(a.ID)", "a.CURRENCY_CODE"}
	k = append(k, kolomBacaData(kolomAnekaData, "a")...)
	k = append(k, kolomBacaData(kolomMarineHull, "m")...)
	k = append(k, kolomBacaData(kolomAviationHull, "v")...)
	k = append(k, kolomBacaData(kolomVehicleHE, "h")...)
	return append(k, kolomBacaData(kolomMemorandum, "mm")...)
}

func sqlBacaAneka(t tabelAneka) string {
	return "SELECT " + strings.Join(kolomBacaAneka(), ", ") + " FROM " + t.aneka + " a LEFT JOIN " + t.marine +
		" m ON m.PARENT_ID = a.ID LEFT JOIN " + t.aviasi + " v ON v.PARENT_ID = a.ID LEFT JOIN " + t.he +
		" h ON h.PARENT_ID = a.ID LEFT JOIN " + t.memo + " mm ON mm.PARENT_ID = a.ID WHERE a.ID IN (" + sqlAnekaKasus(t) +
		") ORDER BY a.PARENT_TABLE, a.PARENT_ID, a.SEQ_NO"
}

// sqlKunciAnekaLama - posisi (lokasi, okupasi / 0 = Golf, objek) tiap objek Aneka case yang ber-coverage; :1 dan :2 = ID case.
func sqlKunciAnekaLama(t tabelAneka) string {
	berCov := " AND a.ID IN (SELECT c.PARENT_ID FROM " + t.cov + " c WHERE c.PARENT_TABLE = '" + indukCoverageAnek + "')"
	return "SELECT l.SEQ_NO, o.SEQ_NO, a.SEQ_NO, TO_CHAR(a.ID) FROM " + t.aneka + " a JOIN " + t.okupasi +
		" o ON o.ID = a.PARENT_ID JOIN " + t.risk + " r ON r.ID = o.PARENT_ID JOIN " + t.prop + " p ON p.ID = r.PARENT_ID JOIN " +
		t.loc + " l ON l.ID = p.PARENT_ID WHERE a.PARENT_TABLE = '" + indukAnekaOkupasi + "' AND o.PARENT_TABLE = '" +
		indukOkupasiAneka + "' AND l.PARENT_ID = :1" + berCov + " UNION ALL SELECT l.SEQ_NO, 0, a.SEQ_NO, TO_CHAR(a.ID) FROM " +
		t.aneka + " a JOIN " + t.risk + " r ON r.ID = a.PARENT_ID JOIN " + t.prop + " p ON p.ID = r.PARENT_ID JOIN " + t.loc +
		" l ON l.ID = p.PARENT_ID WHERE a.PARENT_TABLE = '" + indukAnekaGolf + "' AND l.PARENT_ID = :2" + berCov
}

// hapusAneka - satu pernyataan hapus dan jumlah bind ID case-nya (:1..:n, semuanya ID case).
type hapusAneka struct {
	q    string
	bind int
}

// sqlHapusAneka - anak lalu objek lalu okupasi Aneka case (sebelum sqlHapusObjek FIRE menghapus lokasi).
func sqlHapusAneka(t tabelAneka) []hapusAneka {
	var q []hapusAneka
	for _, anak := range []string{t.memo, t.he, t.aviasi, t.marine} {
		q = append(q, hapusAneka{"DELETE FROM " + anak + " WHERE PARENT_ID IN (" + sqlAnekaKasus(t) + ")", 2})
	}
	return append(q, hapusAneka{"DELETE FROM " + t.aneka + " WHERE ID IN (" + sqlAnekaKasus(t) + ")", 2},
		hapusAneka{"DELETE FROM " + t.okupasi + " o WHERE o.PARENT_TABLE = '" + indukOkupasiAneka + "' AND o.PARENT_ID IN (" +
			sqlRiskKasus(t, ":1") + ")", 1})
}

// sqlPindahCoverage - coverage Aneka objek lama :2 -> objek baru :1.
func sqlPindahCoverage(cov string) string {
	return "UPDATE " + cov + " SET PARENT_ID = :1 WHERE PARENT_TABLE = '" + indukCoverageAnek + "' AND PARENT_ID = :2"
}

// sqlHapusCoverageAneka - coverage objek lama :1 yang hilang: spreading, deductible, lalu coverage (FK).
func sqlHapusCoverageAneka(t tabelAneka) []string {
	cov := "SELECT c.ID FROM " + t.cov + " c WHERE c.PARENT_TABLE = '" + indukCoverageAnek + "' AND c.PARENT_ID = :1"
	return []string{"DELETE FROM " + t.spread + " WHERE PARENT_ID IN (" + cov + ")",
		"DELETE FROM " + t.ded + " WHERE PARENT_ID IN (" + cov + ")",
		"DELETE FROM " + t.cov + " WHERE PARENT_TABLE = '" + indukCoverageAnek + "' AND PARENT_ID = :1"}
}

func sqlSisipData(tabel string, sistem []string, nama, nilai []string) string {
	isi := make([]string, len(sistem))
	for i := range sistem {
		isi[i] = ":" + strconv.Itoa(i+1)
	}
	return "INSERT INTO " + tabel + " (" + strings.Join(append(append([]string{}, sistem...), nama...), ", ") + ") VALUES (" +
		strings.Join(append(isi, nilai...), ", ") + ")"
}

// sqlSisipAneka - :1 ID, :2 PARENT_ID, :3 PARENT_TABLE, :4 SRC_PATH, :5 SEQ_NO, :6 ROW_UID, [:7 CURRENCY_CODE], lalu
// kolomAnekaData. Tanpa mata uang (mis. Aviation Hull C7: form tanpa Currency) kolom CURRENCY_CODE TIDAK disebut sehingga
// DEFAULT 'UNKNOWN' rancangan (K-069) yang mengisi - aplikasi tidak menulis 'UNKNOWN' (K-012); dibaca kembali kosong
// (mataUangBaca, pola A169 deductible).
func sqlSisipAneka(aneka string, denganMataUang bool) string {
	sistem := []string{"ID", "PARENT_ID", "PARENT_TABLE", "SRC_PATH", "SEQ_NO", "ROW_UID"}
	if denganMataUang {
		sistem = append(sistem, "CURRENCY_CODE")
	}
	nama, nilai := kolomSisipData(kolomAnekaData, len(sistem)+1)
	return sqlSisipData(aneka, sistem, nama, nilai)
}

func sqlSisipAnak[T any](tabel string, kolom []kolomData[T]) string {
	nama, nilai := kolomSisipData(kolom, 3)
	return sqlSisipData(tabel, []string{"ID", "PARENT_ID"}, nama, nilai)
}

// BacaObjekAneka - lihat PenyimpanObjekAneka.
func (r *ObjekOracle) BacaObjekAneka(ctx context.Context, id string) ([]models.LokasiAneka, error) {
	t, err := r.tabelAneka()
	if err != nil {
		return nil, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, sqlAdaKasus(t.work), id, LiniFacIn).Scan(&n); err != nil {
		return nil, fmt.Errorf("repository: cek case NB: %w", err)
	}
	if n == 0 {
		return nil, ErrKasusTidakAda
	}
	hasil := []models.LokasiAneka{}
	var idRisk []string
	baris, err := r.db.QueryContext(ctx, sqlBacaLokasiAneka(t), id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca lokasi Aneka: %w", err)
	}
	defer baris.Close()
	for baris.Next() {
		teks, err := pindaiBaris(baris, kolomBacaLokasiAneka)
		if err != nil {
			return nil, fmt.Errorf("repository: lokasi Aneka: %w", err)
		}
		v := func(k string) string { return teks[k].String }
		hasil = append(hasil, models.LokasiAneka{ObjectNo: v("p.OBJECT_NO"), IsTopRisk: v("p.IS_TOP_RISK") == teksBenar,
			RoadType: v("p.ROAD_TYPE"), RoadName: v("p.ROAD_NAME"), BuildingNo: v("p.BUILDING_NO"), ZipCode: v("r.ASM_ZIP_CODE"),
			Country: v("p.COUNTRY"), RiskLocation: v("r.ASM_ADDRESS"), Territory: v("r.ASMRW"), City: v("r.ASM_CITY"),
			District: v("r.ASM_DISTRICT"), Province: v("p.PROVINCE"), RiskAddressID: v("p.ALM_RISK_ID"),
			Okupasi: []models.OkupasiAneka{}, ObjekGolf: []models.ObjekAneka{}})
		idRisk = append(idRisk, v("TO_CHAR(r.ID)"))
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: lokasi Aneka: %w", err)
	}
	baris.Close()
	okupasi, idOkupasi, err := r.bacaOkupasiAneka(ctx, t, id)
	if err != nil {
		return nil, err
	}
	objek, err := r.bacaAneka(ctx, t, id)
	if err != nil {
		return nil, err
	}
	for i, risk := range idRisk {
		if risk == "" {
			continue
		}
		if d, ada := objek[indukAnekaGolf+"/"+risk]; ada {
			hasil[i].ObjekGolf = d
		}
		for j, o := range okupasi[risk] {
			if d, ada := objek[indukAnekaOkupasi+"/"+idOkupasi[risk][j]]; ada {
				o.Objek = d
			}
			hasil[i].Okupasi = append(hasil[i].Okupasi, o)
		}
	}
	return hasil, nil
}

// bacaOkupasiAneka - okupasi Aneka per ID T_RISKLOCATION (teks) dan ID okupasinya (urutan sama).
func (r *ObjekOracle) bacaOkupasiAneka(ctx context.Context, t tabelAneka, id string) (map[string][]models.OkupasiAneka,
	map[string][]string, error) {
	baris, err := r.db.QueryContext(ctx, sqlBacaOkupasiAneka(t), id)
	if err != nil {
		return nil, nil, fmt.Errorf("repository: baca okupasi Aneka: %w", err)
	}
	defer baris.Close()
	okupasi, ids := map[string][]models.OkupasiAneka{}, map[string][]string{}
	for baris.Next() {
		teks, err := pindaiBaris(baris, kolomBacaOkupasiAneka)
		if err != nil {
			return nil, nil, fmt.Errorf("repository: okupasi Aneka: %w", err)
		}
		induk := teks["TO_CHAR(o.PARENT_ID)"].String
		okupasi[induk] = append(okupasi[induk], models.OkupasiAneka{OccupationID: teks["o.OCCUPATION_ID"].String,
			OccupationName: teks["o.OCCUPATION_NAME"].String, Objek: []models.ObjekAneka{}})
		ids[induk] = append(ids[induk], teks["TO_CHAR(o.ID)"].String)
	}
	if err := baris.Err(); err != nil {
		return nil, nil, fmt.Errorf("repository: okupasi Aneka: %w", err)
	}
	return okupasi, ids, nil
}

// bacaAneka - objek Aneka case per "PARENT_TABLE/PARENT_ID".
func (r *ObjekOracle) bacaAneka(ctx context.Context, t tabelAneka, id string) (map[string][]models.ObjekAneka, error) {
	// Coverage (dan deductible) dibaca lebih dulu - kursor objek belum terbuka.
	cov, err := r.bacaCoverageAneka(ctx, t, id)
	if err != nil {
		return nil, err
	}
	kunci := kolomBacaAneka()
	baris, err := r.db.QueryContext(ctx, sqlBacaAneka(t), id, id)
	if err != nil {
		return nil, fmt.Errorf("repository: baca objek Aneka: %w", err)
	}
	defer baris.Close()
	hasil := map[string][]models.ObjekAneka{}
	label := TabelAnekaList + " case " + id
	for baris.Next() {
		teks, err := pindaiBaris(baris, kunci)
		if err != nil {
			return nil, fmt.Errorf("repository: objek Aneka: %w", err)
		}
		o, err := isiData(kolomAnekaData, "a", label, teks)
		if err != nil {
			return nil, err
		}
		o.Currency = mataUangBaca(teks["a.CURRENCY_CODE"].String)
		if o.MarineHull, err = isiData(kolomMarineHull, "m", label, teks); err != nil {
			return nil, err
		}
		if o.AviationHull, err = isiData(kolomAviationHull, "v", label, teks); err != nil {
			return nil, err
		}
		if o.VehicleHE, err = isiData(kolomVehicleHE, "h", label, teks); err != nil {
			return nil, err
		}
		if o.Memorandum, err = isiData(kolomMemorandum, "mm", label, teks); err != nil {
			return nil, err
		}
		if o.Coverages = cov[teks["TO_CHAR(a.ID)"].String]; o.Coverages == nil {
			o.Coverages = []models.CoverageAneka{}
		}
		induk := teks["a.PARENT_TABLE"].String + "/" + teks["TO_CHAR(a.PARENT_ID)"].String
		hasil[induk] = append(hasil[induk], o)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: objek Aneka: %w", err)
	}
	return hasil, nil
}

// kunciAneka - posisi objek: lokasi (1..n), okupasi (1..n; 0 = Golf), objek (1..n).
func kunciAneka(lokasi, okupasi, objek int) string {
	return strconv.Itoa(lokasi) + "/" + strconv.Itoa(okupasi) + "/" + strconv.Itoa(objek)
}

// GantiObjekAneka - lihat PenyimpanObjekAneka. Urutan: sentuh case (kunci + 404), pastikan General, catat posisi objek lama
// yang ber-coverage, hapus Aneka lalu seluruh lokasi (sqlHapusObjek FIRE), sisip ulang, pindahkan / hapus coverage Aneka.
func (r *ObjekOracle) GantiObjekAneka(ctx context.Context, tx *db.Tx, id string, lokasi []models.LokasiAneka, proRate *apd.Decimal) error {
	t, err := r.tabelAneka()
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
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(t.general), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlUbahProRateKasus(t.general), "menulis % Pro Rate case", ikatDesimal(proRate), id); err != nil {
		return err
	}
	lama, err := kunciAnekaLama(ctx, tx, t, id)
	if err != nil {
		return err
	}
	for _, h := range sqlHapusAneka(t) {
		arg := make([]any, h.bind)
		for i := range arg {
			arg[i] = id
		}
		if _, err := jalankan(ctx, tx, h.q, "menghapus objek Aneka", arg...); err != nil {
			return err
		}
	}
	for _, q := range sqlHapusObjek(t.tabelObjek) {
		if _, err := jalankan(ctx, tx, q, "menghapus objek", id); err != nil {
			return err
		}
	}
	baru, ganti, err := sisipAnekaBerkelompok(ctx, oracleKelompok{d: r.db, tx: tx}, t, id, lokasi, uidAcak)
	if err != nil {
		return err
	}
	pindah, hapus := rencanaCoverageAneka(lama, baru, ganti)
	for _, p := range pindah {
		if _, err := jalankan(ctx, tx, sqlPindahCoverage(t.cov), "memindah coverage Aneka", p.baru, p.lama); err != nil {
			return err
		}
	}
	for _, idLama := range hapus {
		for _, q := range sqlHapusCoverageAneka(t) {
			if _, err := jalankan(ctx, tx, q, "menghapus coverage Aneka", idLama); err != nil {
				return err
			}
		}
	}
	return nil
}

// pindahCoverage - coverage Aneka objek lama -> objek baru di posisi yang sama.
type pindahCoverage struct{ lama, baru string }

// rencanaCoverageAneka - objek lama ber-coverage (posisi -> ID): posisi yang masih ada dan coverage-nya TIDAK dikirim (`ganti`)
// -> dipindah ke ID baru; selebihnya (posisi hilang, atau coverage dikirim dan sudah disisip baru) coverage lama dihapus. Urut
// posisi (deterministik).
func rencanaCoverageAneka(lama, baru map[string]string, ganti map[string]bool) ([]pindahCoverage, []string) {
	kunci := make([]string, 0, len(lama))
	for k := range lama {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	var pindah []pindahCoverage
	var hapus []string
	for _, k := range kunci {
		if idBaru, ada := baru[k]; ada && !ganti[k] {
			pindah = append(pindah, pindahCoverage{lama: lama[k], baru: idBaru})
			continue
		}
		hapus = append(hapus, lama[k])
	}
	return pindah, hapus
}

// kunciAnekaLama - posisi -> ID objek Aneka lama yang ber-coverage.
func kunciAnekaLama(ctx context.Context, tx *db.Tx, t tabelAneka, id string) (map[string]string, error) {
	q := sqlKunciAnekaLama(t)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := tx.QueryContext(ctx, q, id, id)
	if err != nil {
		return nil, fmt.Errorf("repository: posisi objek Aneka: %w", err)
	}
	defer baris.Close()
	hasil := map[string]string{}
	for baris.Next() {
		var l, o, a int
		var ida string
		if err := baris.Scan(&l, &o, &a, &ida); err != nil {
			return nil, fmt.Errorf("repository: posisi objek Aneka: %w", err)
		}
		hasil[kunciAneka(l, o, a)] = ida
	}
	return hasil, baris.Err()
}

// tulisAneka - traversal baris tab Object Aneka (pola tulisObjek); baru = posisi -> ID objek (lintasan terakhir).
type tulisAneka struct {
	t     tabelAneka
	p     penulisBaris
	uid   func() (string, error)
	baru  map[string]string
	ganti map[string]bool // posisi objek yang coverage-nya dikirim (diganti utuh, P2)
}

type langkahSisip struct {
	q, tabel string
	arg      []any
}

func (w *tulisAneka) jalankan(ctx context.Context, l []langkahSisip) error {
	for _, x := range l {
		if err := w.p.sisip(ctx, x.q, x.tabel, x.arg); err != nil {
			return err
		}
	}
	return nil
}

func (w *tulisAneka) sisipLokasi(ctx context.Context, id string, lokasi []models.LokasiAneka) error {
	t := w.t
	k := db.KosongJadiNil
	for i, o := range lokasi {
		ids := make([]string, 3)
		for j, seq := range []string{sequenceLocationList, sequenceProperty, sequenceRiskLocation} {
			var err error
			if ids[j], err = w.p.nomor(ctx, seq); err != nil {
				return err
			}
		}
		uid, err := w.uid()
		if err != nil {
			return err
		}
		if err := w.jalankan(ctx, []langkahSisip{
			{sqlSisipLokasi(t.loc), TabelLocationList, []any{ids[0], id, i + 1, uid, nil, nil, nil, nil}},
			{sqlSisipProperty(t.prop), TabelProperty, []any{ids[1], ids[0], k(o.ObjectNo), nil, nil, nil, teksBool(o.IsTopRisk),
				k(o.RoadType), k(o.RoadName), k(o.BuildingNo), k(o.Country), k(o.Province), k(o.RiskAddressID), nil, nil, nil, nil}},
			{sqlSisipRisk(t.risk), TabelRiskLocation, []any{ids[2], ids[1], k(o.RiskLocation), k(o.City), k(o.District),
				k(o.Territory), k(o.ZipCode)}},
		}); err != nil {
			return err
		}
		for j, ok := range o.Okupasi {
			idOk, err := w.p.nomor(ctx, sequenceOccupationList)
			if err != nil {
				return err
			}
			uid, err := w.uid()
			if err != nil {
				return err
			}
			if err := w.p.sisip(ctx, sqlSisipOkupasi(t.okupasi), TabelOccupationList, []any{idOk, ids[2], indukOkupasiAneka,
				jalurOkupasiAneka, j + 1, uid, k(ok.OccupationID), k(ok.OccupationName)}); err != nil {
				return err
			}
			if err := w.sisipObjek(ctx, idOk, indukAnekaOkupasi, jalurAnekaOkupasi, jalurCoverageAnekaOkupasi, i+1, j+1, ok.Objek); err != nil {
				return err
			}
		}
		if err := w.sisipObjek(ctx, ids[2], indukAnekaGolf, jalurAnekaGolf, jalurCoverageAnekaGolf, i+1, 0, o.ObjekGolf); err != nil {
			return err
		}
	}
	return nil
}

// sisipObjek - objek satu induk (SEQ_NO 1..n) dan halaman anak yang terisi.
func (w *tulisAneka) sisipObjek(ctx context.Context, induk, tabelInduk, jalur, jalurCov string, lokasi, okupasi int,
	objek []models.ObjekAneka) error {
	t := w.t
	for n, o := range objek {
		idA, err := w.p.nomor(ctx, sequenceAnekaList)
		if err != nil {
			return err
		}
		w.baru[kunciAneka(lokasi, okupasi, n+1)] = idA
		uid, err := w.uid()
		if err != nil {
			return err
		}
		arg := []any{idA, induk, tabelInduk, jalur, n + 1, uid}
		if o.Currency != "" {
			arg = append(arg, o.Currency)
		}
		arg = append(arg, argData(kolomAnekaData, o)...)
		if err := w.p.sisip(ctx, sqlSisipAneka(t.aneka, o.Currency != ""), TabelAnekaList, arg); err != nil {
			return err
		}
		for _, a := range []struct {
			kosong      bool
			seq, tabel  string
			q           string
			argumenData func() []any
		}{
			{o.MarineHull.Kosong(), sequenceMarineHull, TabelMarineHull, sqlSisipAnak(t.marine, kolomMarineHull),
				func() []any { return argData(kolomMarineHull, o.MarineHull) }},
			{o.AviationHull.Kosong(), sequenceAviationHull, TabelAviationHull, sqlSisipAnak(t.aviasi, kolomAviationHull),
				func() []any { return argData(kolomAviationHull, o.AviationHull) }},
			{o.VehicleHE.Kosong(), sequenceVehicleHE, TabelVehicleHE, sqlSisipAnak(t.he, kolomVehicleHE),
				func() []any { return argData(kolomVehicleHE, o.VehicleHE) }},
			{o.Memorandum.Kosong(), sequenceMemorandum, TabelMemorandumAccept, sqlSisipAnak(t.memo, kolomMemorandum),
				func() []any { return argData(kolomMemorandum, o.Memorandum) }},
		} {
			if a.kosong {
				continue
			}
			idAnak, err := w.p.nomor(ctx, a.seq)
			if err != nil {
				return err
			}
			if err := w.p.sisip(ctx, a.q, a.tabel, append([]any{idAnak, idA}, a.argumenData()...)); err != nil {
				return err
			}
		}
		if o.Coverages != nil {
			w.ganti[kunciAneka(lokasi, okupasi, n+1)] = true
			if err := w.sisipCoverageAneka(ctx, idA, jalurCov, o.Currency, o.Coverages); err != nil {
				return err
			}
		}
	}
	return nil
}

// sisipAnekaBerkelompok - pola sisipBerkelompok (tulisberkelompok.go): lintasan hitung, ID berkelompok, lintasan tulis,
// INSERT array-bind. Jawab posisi -> ID objek baru.
func sisipAnekaBerkelompok(ctx context.Context, e eksekutorKelompok, t tabelAneka, id string, lokasi []models.LokasiAneka,
	uid func() (string, error)) (map[string]string, map[string]bool, error) {
	hitung := &penghitungID{n: map[string]int{}}
	nolUID := func() (string, error) { return "", nil }
	if err := (&tulisAneka{t: t, p: hitung, uid: nolUID, baru: map[string]string{}, ganti: map[string]bool{}}).sisipLokasi(ctx, id,
		lokasi); err != nil {
		return nil, nil, err
	}
	kumpulan := map[string][]string{}
	var jumlah []int
	perJumlah := map[int][]string{}
	for _, s := range hitung.urut {
		n := hitung.n[s]
		if _, ada := perJumlah[n]; !ada {
			jumlah = append(jumlah, n)
		}
		perJumlah[n] = append(perJumlah[n], s)
	}
	for _, n := range jumlah {
		ids, err := e.nomorBerkelompok(ctx, perJumlah[n], n)
		if err != nil {
			return nil, nil, err
		}
		for i, s := range perJumlah[n] {
			kumpulan[s] = ids[i]
		}
	}
	kel := &penulisKelompok{id: kumpulan, indeks: map[string]*kelompokSisip{}}
	w := &tulisAneka{t: t, p: kel, uid: uid, baru: map[string]string{}, ganti: map[string]bool{}}
	if err := w.sisipLokasi(ctx, id, lokasi); err != nil {
		return nil, nil, err
	}
	for s, sisa := range kel.id {
		if len(sisa) != 0 {
			return nil, nil, fmt.Errorf("repository: %d ID %s tidak terpakai (lintasan hitung dan tulis berbeda)", len(sisa), s)
		}
	}
	for _, k := range kel.kelompok {
		for awal := 0; awal < len(k.baris); awal += ukuranKelompok {
			akhir := min(awal+ukuranKelompok, len(k.baris))
			if err := e.sisipArray(ctx, k.q, k.tabel, kolomArray(k.baris[awal:akhir]), akhir-awal); err != nil {
				return nil, nil, err
			}
		}
	}
	return w.baru, w.ganti, nil
}

// sqlPercentShareKasus - :1 ID case; nilai lewat angkaKeluar (NLS titik).
func sqlPercentShareKasus(general string) string {
	return "SELECT " + angkaKeluar("PERCENT_SHARE") + " FROM " + general + " WHERE ID = :1"
}

// PercentShareKasus - lihat PenyimpanObjekAneka.
func (r *ObjekOracle) PercentShareKasus(ctx context.Context, id string) (*apd.Decimal, error) {
	return r.desimalGeneral(ctx, id, sqlPercentShareKasus, "PERCENT_SHARE")
}

// sqlUbahProRateKasus - :1 % Pro Rate (angkaMasuk; NULL boleh), :2 ID case.
func sqlUbahProRateKasus(general string) string {
	return "UPDATE " + general + " SET PRO_RATE_PERCENT = " + angkaMasuk(":1") + " WHERE ID = :2"
}

// sqlProRateKasus - :1 ID case.
func sqlProRateKasus(general string) string {
	return "SELECT " + angkaKeluar("PRO_RATE_PERCENT") + " FROM " + general + " WHERE ID = :1"
}

// ProRateKasus - lihat PenyimpanObjekAneka.
// UbahPercentShare - lihat PenyimpanObjekAneka (SQL milik tab Spreading FIRE: sqlPastikanGeneral, sqlUbahShare).
func (r *ObjekOracle) UbahPercentShare(ctx context.Context, tx *db.Tx, id string, ps *apd.Decimal) error {
	q, err := r.db.Qualify(TabelGeneralPolis)
	if err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlPastikanGeneral(q), "memastikan "+TabelGeneralPolis, id, id); err != nil {
		return err
	}
	_, err = jalankan(ctx, tx, sqlUbahShare(q), "menulis % Share RNM", ikatDesimal(ps), id)
	return err
}

func (r *ObjekOracle) ProRateKasus(ctx context.Context, id string) (*apd.Decimal, error) {
	return r.desimalGeneral(ctx, id, sqlProRateKasus, "PRO_RATE_PERCENT")
}

// desimalGeneral - satu kolom desimal T_GENERAL_POLIS case; baris tidak ada -> nil.
func (r *ObjekOracle) desimalGeneral(ctx context.Context, id string, susun func(string) string, kolom string) (*apd.Decimal, error) {
	q, err := r.db.Qualify(TabelGeneralPolis)
	if err != nil {
		return nil, err
	}
	var v sql.NullString
	if err := r.db.QueryRowContext(ctx, susun(q), id).Scan(&v); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repository: baca %s.%s: %w", TabelGeneralPolis, kolom, err)
	}
	return bacaDesimal(TabelGeneralPolis+" "+id, kolom, &v)
}
