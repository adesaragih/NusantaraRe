package repository

// Save Data bilah CSV objek General (tiket 66, `NB FacIn\Activity\InsertUploadFire_act.xml` + MappingJsonH2HNusare2_act +
// MappingCoverage): pencarian master PERSIS seperti Pega (bukan pencarian layar). Semua SELECT biasa - nol prosedur
// (ADR-0043). Tulisan master (risk address / akumulasi) memakai port tiket 37 / 46 yang sudah ada.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// TabelObjectItemType - POOLDATA.OBJECTITEMTYPE (DDL\OBJECTITEMTYPE.txt), baca saja.
const TabelObjectItemType = "OBJECTITEMTYPE"

// PembacaCsvObjek - pencarian master Save Data CSV dan penghapus CurrencyList case.
type PembacaCsvObjek interface {
	// RWKelurahan - RDBList\BrowseRW2_SQL: RW ber-zipcode itu dan NOTE = kelurahan (spasi diabaikan); baris pertama dipakai.
	RWKelurahan(ctx context.Context, zip, kelurahan string) ([]models.RWCsv, error)
	// RWKodePos - tiket 73 P6 (K73-58, CSV Aneka tanpa KELURAHAN): RW pertama ber-zipcode itu, urut NOTE / DISTRICTNAME /
	// CITYNAME `[keputusan agent]` (Pega Obj-Browse Int-RW `.ZipCode = POSTAL_CODE AND .Note = KELURAHAN` dengan Note kosong -
	// perilakunya `[belum terverifikasi]`); ok = false bila tidak ada.
	RWKodePos(ctx context.Context, zip string) (models.RWCsv, bool, error)
	// IDRiskAddress - Obj-Browse Int-RISKADDRESS `.ADDRESS = alamat AND .POSTALCODE = zip` (langkah 5.9) -> ID pertama; "" = tidak ada.
	IDRiskAddress(ctx context.Context, alamat, zip string) (string, error)
	// IDAkumulasi - Obj-Browse Int-ACCUMULATION `.NOTE = nama` (langkah 5.10) -> ID pertama; "" = tidak ada.
	IDAkumulasi(ctx context.Context, nama string) (string, error)
	// InisialNegara - Obj-Browse Int-NATION `.NOTE = nama` (langkah 5.10.2.4) -> NATIONINITIAL; "" = tidak ada.
	InisialNegara(ctx context.Context, nama string) (string, error)
	// IDProvinsi - PROVINCE ber-NOTE nama itu -> ID (Pega mengirim ProvinceName ke prosedur; port tiket 46 butuh ID); "" = tidak ada.
	IDProvinsi(ctx context.Context, nama string) (string, error)
	// MataUang - RDBList\GetDataCurrencyByName_SQL; ok false = tidak ada.
	MataUang(ctx context.Context, nama string) (models.MataUangCsv, bool, error)
	// IDJenisItemFire - RDBList\GetDataObjectItem (TYPE "FIRE", nama tanpa spasi, ISACTIVE '1') -> ID; "" = tidak ada.
	IDJenisItemFire(ctx context.Context, nama string) (string, error)
	// NamaOkupasiFire - ReportDefinition\BrowseOccupationFacInFIRE_RD (TYPE "FIRE", .OldID = kode; urut Name, OldID) -> NAME.
	NamaOkupasiFire(ctx context.Context, kode string) (string, error)
	// CoverageAda - ID yang ada di POOLDATA.COVERAGE (MappingCoverage memakai ID master ter-hardcode; D8).
	CoverageAda(ctx context.Context, id []string) (map[string]bool, error)
	// HapusCurrencyList - InsertUploadFire_act langkah 3 `Property-Remove .CurrencyList`: T_CURRENCYLIST case beserta
	// T_LISTINSTALLMENT anaknya (D5 iii, tanpa dibangun ulang).
	HapusCurrencyList(ctx context.Context, tx *db.Tx, id string) error
}

// CsvObjekOracle - PembacaCsvObjek atas Oracle.
type CsvObjekOracle struct{ db *db.DB }

// NewCsvObjekOracle merakit pembaca master Save Data CSV.
func NewCsvObjekOracle(d *db.DB) *CsvObjekOracle { return &CsvObjekOracle{db: d} }

func sqlRWKelurahan(rw string) string {
	return "SELECT NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FROM " + rw +
		" WHERE ZIPCODE = :1 AND REPLACE(NOTE, ' ', '') = REPLACE(:2, ' ', '')"
}

func sqlRWKodePos(rw string) string {
	return "SELECT NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FROM " + rw +
		" WHERE ZIPCODE = :1 ORDER BY NOTE, DISTRICTNAME, CITYNAME FETCH FIRST 1 ROWS ONLY"
}

func sqlCariIDRiskAddress(ra string) string {
	return "SELECT ID FROM " + ra + " WHERE ADDRESS = :1 AND POSTALCODE = :2 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlIDAkumulasi(acc string) string {
	return "SELECT ID FROM " + acc + " WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlInisialNegara(nation string) string {
	return "SELECT NATIONINITIAL FROM " + nation + " WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlIDProvinsi(prov string) string {
	return "SELECT ID FROM " + prov + " WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlMataUangCsv(cur string) string {
	return "SELECT ID, OLDID, CURRENCY FROM " + cur + " WHERE CURRENCY = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlJenisItemFire(oit string) string {
	return "SELECT ID FROM " + oit + " WHERE TYPE = 'FIRE' AND REPLACE(OBJECTITEMTYPE, ' ', '') = REPLACE(:1, ' ', '')" +
		" AND ISACTIVE = '1' ORDER BY ID FETCH FIRST 1 ROWS ONLY"
}

func sqlNamaOkupasiFire(occ string) string {
	return "SELECT NAME FROM " + occ + " WHERE TYPE = 'FIRE' AND OLDID = :1 ORDER BY NAME, OLDID FETCH FIRST 1 ROWS ONLY"
}

func sqlCoverageAda(cov string, n int) string {
	return "SELECT ID FROM " + cov + " WHERE ID IN (" + daftarBind(1, n) + ")"
}

func sqlHapusCicilanCase(inst, cur string) string {
	return "DELETE FROM " + inst + " WHERE PARENT_ID IN (SELECT ID FROM " + cur + " WHERE PARENT_ID = :1)"
}

func sqlHapusCurrencyCase(cur string) string { return "DELETE FROM " + cur + " WHERE PARENT_ID = :1" }

// daftarBind - ":a, :a+1, … " sebanyak n.
func daftarBind(awal, n int) string {
	b := make([]string, n)
	for i := range b {
		b[i] = fmt.Sprintf(":%d", awal+i)
	}
	return strings.Join(b, ", ")
}

// satuTeks - satu kolom teks baris pertama; tanpa baris -> "".
func (r *CsvObjekOracle) satuTeks(ctx context.Context, tabel string, susun func(string) string, arg ...any) (string, error) {
	q, err := r.db.Qualify(tabel)
	if err != nil {
		return "", err
	}
	var v sql.NullString
	err = r.db.QueryRowContext(ctx, susun(q), arg...).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: baca %s: %w", tabel, err)
	}
	return strings.TrimSpace(v.String), nil
}

// RWKelurahan - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) RWKelurahan(ctx context.Context, zip, kelurahan string) ([]models.RWCsv, error) {
	q, err := r.db.Qualify(TabelRW)
	if err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, sqlRWKelurahan(q), zip, kelurahan)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelRW, err)
	}
	defer baris.Close()
	var hasil []models.RWCsv
	for baris.Next() {
		var v [5]sql.NullString
		if err := baris.Scan(&v[0], &v[1], &v[2], &v[3], &v[4]); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelRW, err)
		}
		hasil = append(hasil, models.RWCsv{Note: v[0].String, District: v[1].String, City: v[2].String, Province: v[3].String,
			Nation: v[4].String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelRW, err)
	}
	return hasil, nil
}

// RWKodePos - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) RWKodePos(ctx context.Context, zip string) (models.RWCsv, bool, error) {
	q, err := r.db.Qualify(TabelRW)
	if err != nil {
		return models.RWCsv{}, false, err
	}
	var v [5]sql.NullString
	err = r.db.QueryRowContext(ctx, sqlRWKodePos(q), zip).Scan(&v[0], &v[1], &v[2], &v[3], &v[4])
	if errors.Is(err, sql.ErrNoRows) {
		return models.RWCsv{}, false, nil
	}
	if err != nil {
		return models.RWCsv{}, false, fmt.Errorf("repository: baca %s: %w", TabelRW, err)
	}
	return models.RWCsv{Note: v[0].String, District: v[1].String, City: v[2].String, Province: v[3].String, Nation: v[4].String},
		true, nil
}

// IDRiskAddress - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) IDRiskAddress(ctx context.Context, alamat, zip string) (string, error) {
	return r.satuTeks(ctx, TabelRiskAddress, sqlCariIDRiskAddress, alamat, zip)
}

// IDAkumulasi - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) IDAkumulasi(ctx context.Context, nama string) (string, error) {
	return r.satuTeks(ctx, TabelAccumulation, sqlIDAkumulasi, nama)
}

// InisialNegara - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) InisialNegara(ctx context.Context, nama string) (string, error) {
	return r.satuTeks(ctx, TabelNation, sqlInisialNegara, nama)
}

// IDProvinsi - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) IDProvinsi(ctx context.Context, nama string) (string, error) {
	return r.satuTeks(ctx, TabelProvince, sqlIDProvinsi, nama)
}

// IDJenisItemFire - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) IDJenisItemFire(ctx context.Context, nama string) (string, error) {
	return r.satuTeks(ctx, TabelObjectItemType, sqlJenisItemFire, nama)
}

// NamaOkupasiFire - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) NamaOkupasiFire(ctx context.Context, kode string) (string, error) {
	return r.satuTeks(ctx, TabelOccupation, sqlNamaOkupasiFire, kode)
}

// MataUang - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) MataUang(ctx context.Context, nama string) (models.MataUangCsv, bool, error) {
	q, err := r.db.Qualify(TabelCurrency)
	if err != nil {
		return models.MataUangCsv{}, false, err
	}
	var id, old, cur sql.NullString
	err = r.db.QueryRowContext(ctx, sqlMataUangCsv(q), nama).Scan(&id, &old, &cur)
	if errors.Is(err, sql.ErrNoRows) {
		return models.MataUangCsv{}, false, nil
	}
	if err != nil {
		return models.MataUangCsv{}, false, fmt.Errorf("repository: baca %s: %w", TabelCurrency, err)
	}
	return models.MataUangCsv{ID: id.String, OldID: old.String, Currency: cur.String}, true, nil
}

// CoverageAda - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) CoverageAda(ctx context.Context, id []string) (map[string]bool, error) {
	hasil := map[string]bool{}
	if len(id) == 0 {
		return hasil, nil
	}
	q, err := r.db.Qualify(TabelCoverage)
	if err != nil {
		return nil, err
	}
	arg := make([]any, len(id))
	for i, v := range id {
		arg[i] = v
	}
	baris, err := r.db.QueryContext(ctx, sqlCoverageAda(q, len(id)), arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", TabelCoverage, err)
	}
	defer baris.Close()
	for baris.Next() {
		var v sql.NullString
		if err := baris.Scan(&v); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", TabelCoverage, err)
		}
		hasil[strings.TrimSpace(v.String)] = true
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", TabelCoverage, err)
	}
	return hasil, nil
}

// HapusCurrencyList - lihat PembacaCsvObjek.
func (r *CsvObjekOracle) HapusCurrencyList(ctx context.Context, tx *db.Tx, id string) error {
	cur, err := r.db.Qualify(TabelCurrencyList)
	if err != nil {
		return err
	}
	inst, err := r.db.Qualify(TabelListInstallment)
	if err != nil {
		return err
	}
	if _, err := jalankan(ctx, tx, sqlHapusCicilanCase(inst, cur), "menghapus "+TabelListInstallment, id); err != nil {
		return err
	}
	_, err = jalankan(ctx, tx, sqlHapusCurrencyCase(cur), "menghapus "+TabelCurrencyList, id)
	return err
}
