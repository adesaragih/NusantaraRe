package repository

// Popup Choose Accumulation Code (tiket 46) - pencarian akumulasi dan saran autocomplete. ACCUMULATION / PROVINCE /
// ACCUMULATEDTYPE / CZONE = tabel flat migrasi 196 (sebelumnya view atas JSON M_*; butir 103), nama dan kolom sama.
//
// Kelas -> tabel `[terverifikasi]`: ASM-FW-GISFW-Int-ACCUMULATION = view POOLDATA.ACCUMULATION (`RDBList\
// GetAccumulationProvince_SQL.xml` dan kawan-kawan "from accumulation"; DDL `DDL\ACCUMULATION.txt` 04-10-2026: SELECT
// DISTINCT atas JSON m_accumulation WHERE IsActive IS NULL, kolom tanpa tipe tertulis - dibaca sebagai teks). Properti
// `.PostalCode` = kolom ZIPCODE `[terverifikasi]` (`RDBList\SearchAccumulationbypersetase_SQL.xml` kelas yang sama:
// `ZIPCODE as "PostalCode"`, `NOTE as "Note"`). Int-RW = tabel RW (`RDBList\CheckZipCode_SQL.xml` "from rw").
// Int-CITY / Int-DISTRICT = view CITY / DISTRICT `[dugaan]` (nama dan seluruh kolom yang dibaca RD sama persis dengan
// DDL `DDL\CITY.txt` / `DISTRICT.txt`; tidak ada SQL kelas itu yang menyebut tabelnya).
//
// Jalur pencarian (`Activity\GetDataAccumulation_act.xml`, A172): RD `SearchRiskAccumulation_RD` atau RDB-List
// `GetAccumulationProvince_SQL` (kota) / `GetAccumulationDistrict_SQL` (kecamatan) / `GetSummaryRiskAccumPolis_Sql`
// (nomor polis) - pemilihan jalur di services.

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

const (
	// TabelAccumulation, TabelJSONPolis, TabelCity, TabelDistrict - objek warisan POOLDATA, baca saja.
	TabelAccumulation = "ACCUMULATION"
	TabelJSONPolis    = "JSON_POLIS"
	TabelCity         = "CITY"
	TabelDistrict     = "DISTRICT"
	// TabelNation, TabelProvince, TabelAccumulatedType, TabelCZone - saran tiket 46 (NATION tabel warisan; tiga lainnya
	// tabel flat migrasi 196, sebelumnya view atas JSON M_*).
	TabelNation = "NATION"
	// TabelCityInput, TabelDistrictInput, TabelMasterStatus - tabel milik modul Master Data (760 / 761) yang menyimpan
	// status aktif master yang dibaca saran ini (view CITY / DISTRICT tidak membawa status; NATION warisan tanpa kolom
	// status). StatusMasterAktif - '1' aktif (M-3).
	TabelCityInput     = "CITYINPUT"
	TabelDistrictInput = "DISTRICTINPUT"
	TabelMasterStatus  = "T_MASTER_STATUS"
	StatusMasterAktif  = "1"
	// syaratMasterAktif - tabel flat master berkolom STS_AKTIF sendiri.
	syaratMasterAktif    = "STS_AKTIF = '" + StatusMasterAktif + "'"
	TabelProvince        = "PROVINCE"
	TabelAccumulatedType = "ACCUMULATEDTYPE"
	TabelCZone           = "CZONE"
	// BatasAkumulasi - pyMaxRecords SearchRiskAccumulation_RD (juga dipakai jalur SQL, A174).
	BatasAkumulasi = 500
	// BatasSaranAkumulasi - saran autocomplete paling banyak 50 (pola A124; RD 1.500.000 / 150.000 / 10.000, A176).
	BatasSaranAkumulasi = 50
)

// PembacaAkumulasi - pencarian akumulasi dan saran autocomplete popup Choose Accumulation Code.
type PembacaAkumulasi interface {
	// CariAkumulasiRD - SearchRiskAccumulation_RD; saringan kosong dibuang.
	CariAkumulasiRD(ctx context.Context, s models.SaringAkumulasi) ([]models.BarisAkumulasi, error)
	// AkumulasiWilayah - GetAccumulationProvince_SQL (kota) / GetAccumulationDistrict_SQL (kecamatan).
	AkumulasiWilayah(ctx context.Context, kecamatan bool, id string) ([]models.BarisAkumulasi, error)
	// AkumulasiPolis - GetSummaryRiskAccumPolis_Sql + penyaringan ganda langkah 6.6.
	AkumulasiPolis(ctx context.Context, noPolis string) ([]models.BarisAkumulasi, error)
	// Saran - autocomplete kota / kecamatan / area.
	Saran(ctx context.Context, jenis, kata, induk string) ([]models.SaranAkumulasi, error)
}

// AkumulasiOracle - PembacaAkumulasi atas Oracle.
type AkumulasiOracle struct{ db *db.DB }

// NewAkumulasiOracle merakit pembaca ACCUMULATION / RW / CITY / DISTRICT / JSON_POLIS.
func NewAkumulasiOracle(d *db.DB) *AkumulasiOracle { return &AkumulasiOracle{db: d} }

// saringRD - saringan SearchRiskAccumulation_RD menurut urutan RD (A, B, D, E, F, G; C tidak dikirim layar, A173):
// kolom view, operator, dan nilainya. `Note` Contains tidak peka huruf (satu-satunya pyCaseInsensitive).
func saringRD(s models.SaringAkumulasi) (syarat []string, arg []any) {
	tambah := func(kolom, nilai string) {
		if nilai == "" {
			return
		}
		arg = append(arg, nilai)
		syarat = append(syarat, "a."+kolom+" = :"+strconv.Itoa(len(arg)))
	}
	syarat = append(syarat, "a.STS_AKTIF = '"+StatusMasterAktif+"'") // MD-5: baris aktif menu Master Data
	tambah("ID", s.ID)
	if pola := PolaCari(s.Note); pola != "" {
		arg = append(arg, pola)
		syarat = append(syarat, "UPPER(a.NOTE) LIKE :"+strconv.Itoa(len(arg))+` ESCAPE '\'`)
	}
	tambah("CZONE", s.CZone)
	tambah("KEYWORD", s.Keyword)
	tambah("ZIPCODE", s.PostalCode)
	tambah("PROVINCEID", s.ProvinceID)
	return syarat, arg
}

// sqlCariAkumulasi - INNER JOIN RW `.PostalCode = RW.ZipCode`, DISTINCT atas sebelas kolom laporan RD (PROVINCE =
// `.ProvinceName`), batas terakhir; urutan ID lalu NOTE (RD tanpa urutan, A174).
func sqlCariAkumulasi(acc, rw string, syarat []string, bindBatas int) string {
	where := ""
	if len(syarat) > 0 {
		where = " WHERE " + strings.Join(syarat, " AND ")
	}
	return "SELECT ID, ACCUMULATIONNAME, NOTE FROM (SELECT DISTINCT a.ACCUMULATION, a.CZONE, a.CZONEID, a.NOTE, a.KEYWORD, " +
		"a.SCOPEAREA, a.ID, a.ACCUMULATIONNAME, a.ZIPCODE, a.PROVINCE, a.PROVINCEID FROM " + acc + " a JOIN " + rw +
		" r ON r.ZIPCODE = a.ZIPCODE" + where + ") ORDER BY ID, NOTE, ACCUMULATIONNAME FETCH FIRST :" + strconv.Itoa(bindBatas) + " ROWS ONLY"
}

// sqlAkumulasiWilayah - `select id as CARI1, accumulationtype AS CARI2, note AS CARI3 from accumulation where ZIPCODE in
// (select zipcode from rw where cityid|DISTRICTID = ...)` persis (CARI2 = ACCUMULATIONTYPE, bukan ACCUMULATIONNAME);
// urutan dan batas A174.
func sqlAkumulasiWilayah(acc, rw string, kecamatan bool) string {
	kolom := "CITYID"
	if kecamatan {
		kolom = "DISTRICTID"
	}
	return "SELECT ID, ACCUMULATIONTYPE, NOTE FROM " + acc + " WHERE " + syaratMasterAktif + " AND ZIPCODE IN (SELECT ZIPCODE FROM " +
		rw + " WHERE " + kolom + " = :1) ORDER BY ID, NOTE FETCH FIRST :2 ROWS ONLY"
}

// sqlAkumulasiPolis - GetSummaryRiskAccumPolis_Sql dipersempit ke yang dipakai langkah 6.6 (AccumulationCode -> CARI1,
// accumulationname -> CARI2, note -> CARI3; ganda dibuang): JSON_TABLE jalur `$.LocationList[*]` >
// `$.Property.PropertyItemList[*]` > `$.CoverageList[*]` `$.AccumulationCode` persis seperti SQL asal; subkueri dengan MAX
// (SQL asal bisa ORA-01427 bila ID ganda, A175); kode kosong dibuang (A175). SELECT biasa - nol prosedur (ADR-0043).
func sqlAkumulasiPolis(polis, acc string) string {
	return "SELECT j.ACCUMULATIONCODE, (SELECT MAX(x.ACCUMULATIONNAME) FROM " + acc + " x WHERE x.ID = j.ACCUMULATIONCODE), " +
		"(SELECT MAX(x.NOTE) FROM " + acc + " x WHERE x.ID = j.ACCUMULATIONCODE) FROM (SELECT DISTINCT jt.ACCUMULATIONCODE FROM " +
		polis + " p, JSON_TABLE(p.DATA_JSON, '$.LocationList[*]' COLUMNS (NESTED PATH '$.Property.PropertyItemList[*]' COLUMNS " +
		"(NESTED PATH '$.CoverageList[*]' COLUMNS (ACCUMULATIONCODE VARCHAR2(100) PATH '$.AccumulationCode')))) jt " +
		"WHERE p.NOPOLIS = :1 AND jt.ACCUMULATIONCODE IS NOT NULL AND EXISTS (SELECT 1 FROM " + acc + " x WHERE x.ID = jt.ACCUMULATIONCODE " +
		"AND x." + syaratMasterAktif + ")) j ORDER BY j.ACCUMULATIONCODE FETCH FIRST :2 ROWS ONLY"
}

// saranRD - satu autocomplete: tabel, kolom laporan RD (`unik` = pyGetDistinctRows), kolom id / label (medan cari,
// pyUseForSearch) / ekstra, saringan tetap ber-bind (`tetap = nilaiTetap`) atau tanpa bind (`syaratTetap`, teks SQL
// tetap), kolom induk, dan urutan RD (`urut`; kosong = label).
type saranRD struct {
	tabel             string
	kolom             string
	unik              bool
	id, label, ekstra string
	tetap, nilaiTetap string
	syaratTetap       string
	kolomInduk        string
	urut              string
	// aktif / tabelAktif - saringan baris AKTIF menu Master Data (MD-5, 04-10-2026): teks SQL tetap; `%s` diganti nama
	// berskema `tabelAktif` (kosong = tanpa tabel lain).
	aktif, tabelAktif string
}

// daftarSaran - jenis -> RD (`Section\SearchRiskAccumCov.xml`, `ReportDefinition\*`) `[terverifikasi]`:
//   - nation   BrowseNation_RD: cari .Note, ekstra .NationInitial (-> SearchAccumulation.SyariahStatus); saringan
//     `.ID = Param.ID OR .Note = Param.Note` dikirim kosong -> dibuang; tanpa DISTINCT. Tabel NATION (DDL NATION.txt).
//   - province BrowseProvince2_RD: cari .Note, id .ID (-> SearchAccumulation.ProvinceID), induk NationName =
//     SearchAccumulation.Nation (NAMA negara); tanpa DISTINCT.
//   - accumtype BrowseAccumulatedType_RD: cari .AccumulationType; `.AccumulationType = Param.AccType` dikirim kosong ->
//     dibuang; `.Note IS NOT NULL` tetap; tanpa DISTINCT.
//   - czone    BrowseCZoneIsNotNull_RD: cari .Code; `.GroupOf IS NOT NULL` tetap; GroupOfName / Code dikirim kosong ->
//     dibuang; urut .Description ASC (pySortOrder 1); tanpa DISTINCT.
//   - city     BrowseCityInput_RD: cari .Note, induk PROVINCEID = SearchAccumulation.ProvinceID, DISTINCT (ID, Note).
//   - district BrowseDistrictInputC_RD: cari .DistrictName, induk CityName = SearchAccumulation.City (NAMA kota),
//     DISTINCT (ID, CityID, DistrictName, CityName).
//   - area     BrowseRW_RD: cari .Note, `.STS_AKTIF = "1"`, DISTINCT tujuh kolom (NATIONNAME = kolom NATION, alias
//     `RDBList\BrowseRW2_SQL.xml`); TANPA saringan induk - layar mengirim param `DistrictName`, RD hanya mengenal
//     `District` (pyParameters: City, District, Province, Teritory, ZipCode) sehingga saringan D dibuang (A177). `.ID`
//     bukan kolom laporan RD -> id kosong.
var daftarSaran = map[string]saranRD{
	"city": {tabel: TabelCity, kolom: "ID, NOTE", unik: true, id: "ID", label: "NOTE", kolomInduk: "PROVINCEID",
		aktif: "ID IN (SELECT c.ID FROM %s c WHERE c.STS_AKTIF = '" + StatusMasterAktif + "')", tabelAktif: TabelCityInput},
	"district": {tabel: TabelDistrict, kolom: "ID, CITYID, DISTRICTNAME, CITYNAME", unik: true, id: "ID", label: "DISTRICTNAME",
		kolomInduk: "CITYNAME", aktif: "ID IN (SELECT d.ID FROM %s d WHERE d.STS_AKTIF = '" + StatusMasterAktif + "')",
		tabelAktif: TabelDistrictInput},
	"area": {tabel: TabelRW, kolom: "ZIPCODE, CZONE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION", unik: true, id: "''",
		label: "NOTE", ekstra: "ZIPCODE", tetap: "STS_AKTIF", nilaiTetap: StatusRWAktif},
	"nation": {tabel: TabelNation, kolom: "ID, NOTE, NATIONINITIAL", id: "ID", label: "NOTE", ekstra: "NATIONINITIAL",
		aktif: "NOT EXISTS (SELECT 1 FROM %s s WHERE s.NAMA_TABEL = '" + TabelNation + "' AND s.ID_BARIS = ID AND s.STS_AKTIF <> '" +
			StatusMasterAktif + "')", tabelAktif: TabelMasterStatus},
	"province": {tabel: TabelProvince, kolom: "ID, NATIONID, NOTE, NATIONNAME", id: "ID", label: "NOTE", kolomInduk: "NATIONNAME",
		aktif: syaratMasterAktif},
	"accumtype": {tabel: TabelAccumulatedType, kolom: "ID, ACCUMULATIONTYPE, KEYWORD, NOTE, TYPE", id: "ID", label: "ACCUMULATIONTYPE",
		syaratTetap: "NOTE IS NOT NULL", aktif: syaratMasterAktif},
	"czone": {tabel: TabelCZone, kolom: "DESCRIPTION, ID, GROUPOF, CODE, GROUPOFNAME", id: "ID", label: "CODE",
		syaratTetap: "GROUPOF IS NOT NULL", urut: "DESCRIPTION", aktif: syaratMasterAktif},
}

// sqlSaran - SELECT id, label, ekstra dari kolom laporan (DISTINCT bila RD-nya); saringan tetap / induk / kata (Contains tidak peka
// huruf atas kolom label, A176) hanya bila ada; urut label, ekstra, lalu id. Syarat dan nilai bind dibangun BERSAMA
// (pola saringRD) - urutan placeholder dan nilai tidak dapat bergeser.
func sqlSaran(t string, s saranRD, aktif, induk, pola string) (string, []any) {
	var syarat []string
	var arg []any
	bind := func(v any) string { arg = append(arg, v); return ":" + strconv.Itoa(len(arg)) }
	if s.tetap != "" {
		syarat = append(syarat, s.tetap+" = "+bind(s.nilaiTetap))
	}
	if s.syaratTetap != "" {
		syarat = append(syarat, s.syaratTetap)
	}
	if aktif != "" {
		syarat = append(syarat, aktif)
	}
	if induk != "" && s.kolomInduk != "" {
		syarat = append(syarat, s.kolomInduk+" = "+bind(induk))
	}
	if pola != "" {
		syarat = append(syarat, "UPPER("+s.label+") LIKE "+bind(pola)+` ESCAPE '\'`)
	}
	where := ""
	if len(syarat) > 0 {
		where = " WHERE " + strings.Join(syarat, " AND ")
	}
	ekstra := "''"
	if s.ekstra != "" {
		ekstra = s.ekstra
	}
	urut := s.label + ", " + ekstra + ", " + s.id
	if s.urut != "" {
		urut = s.urut + ", " + urut
	}
	pilih := "SELECT "
	if s.unik {
		pilih = "SELECT DISTINCT "
	}
	return "SELECT " + s.id + ", " + s.label + ", " + ekstra + " FROM (" + pilih + s.kolom + " FROM " + t + where +
		") ORDER BY " + urut + " FETCH FIRST " + bind(BatasSaranAkumulasi) + " ROWS ONLY", arg
}

// JenisSaranDidukung - jenis saran yang dapat dilayani (sumber terverifikasi DDL).
func JenisSaranDidukung(jenis string) bool { _, ada := daftarSaran[jenis]; return ada }

// CariAkumulasiRD - lihat PembacaAkumulasi.
func (r *AkumulasiOracle) CariAkumulasiRD(ctx context.Context, s models.SaringAkumulasi) ([]models.BarisAkumulasi, error) {
	acc, rw, err := r.qualifyDua(TabelAccumulation, TabelRW)
	if err != nil {
		return nil, err
	}
	syarat, arg := saringRD(s)
	return r.bacaAkumulasi(ctx, TabelAccumulation, sqlCariAkumulasi(acc, rw, syarat, len(arg)+1), append(arg, BatasAkumulasi)...)
}

// AkumulasiWilayah - lihat PembacaAkumulasi.
func (r *AkumulasiOracle) AkumulasiWilayah(ctx context.Context, kecamatan bool, id string) ([]models.BarisAkumulasi, error) {
	acc, rw, err := r.qualifyDua(TabelAccumulation, TabelRW)
	if err != nil {
		return nil, err
	}
	return r.bacaAkumulasi(ctx, TabelAccumulation, sqlAkumulasiWilayah(acc, rw, kecamatan), id, BatasAkumulasi)
}

// AkumulasiPolis - lihat PembacaAkumulasi.
func (r *AkumulasiOracle) AkumulasiPolis(ctx context.Context, noPolis string) ([]models.BarisAkumulasi, error) {
	polis, acc, err := r.qualifyDua(TabelJSONPolis, TabelAccumulation)
	if err != nil {
		return nil, err
	}
	return r.bacaAkumulasi(ctx, TabelJSONPolis, sqlAkumulasiPolis(polis, acc), noPolis, BatasAkumulasi)
}

// Saran - lihat PembacaAkumulasi. Jenis tak dikenal -> galat (services memeriksanya lebih dulu).
func (r *AkumulasiOracle) Saran(ctx context.Context, jenis, kata, induk string) ([]models.SaranAkumulasi, error) {
	s, ada := daftarSaran[jenis]
	if !ada {
		return nil, fmt.Errorf("repository: saran akumulasi %q tidak dikenal", jenis)
	}
	t, err := r.db.Qualify(s.tabel)
	if err != nil {
		return nil, err
	}
	aktif := s.aktif
	if s.tabelAktif != "" {
		ta, err := r.db.Qualify(s.tabelAktif)
		if err != nil {
			return nil, err
		}
		aktif = fmt.Sprintf(s.aktif, ta)
	}
	q, arg := sqlSaran(t, s, aktif, induk, PolaCari(kata))
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", s.tabel, err)
	}
	defer baris.Close()
	hasil := []models.SaranAkumulasi{}
	for baris.Next() {
		var id, label, ekstra sql.NullString
		if err := baris.Scan(&id, &label, &ekstra); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", s.tabel, err)
		}
		hasil = append(hasil, models.SaranAkumulasi{ID: id.String, Label: label.String, Ekstra: ekstra.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", s.tabel, err)
	}
	return hasil, nil
}

// qualifyDua - nama berskema dua tabel sekaligus.
func (r *AkumulasiOracle) qualifyDua(a, b string) (string, string, error) {
	qa, err := r.db.Qualify(a)
	if err != nil {
		return "", "", err
	}
	qb, err := r.db.Qualify(b)
	return qa, qb, err
}

// bacaAkumulasi - baris (id, nama, note); `tabel` = sumber utama kueri untuk pesan galat.
func (r *AkumulasiOracle) bacaAkumulasi(ctx context.Context, tabel, q string, arg ...any) ([]models.BarisAkumulasi, error) {
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: baca %s: %w", tabel, err)
	}
	defer baris.Close()
	hasil := []models.BarisAkumulasi{}
	for baris.Next() {
		var id, nama, catatan sql.NullString
		if err := baris.Scan(&id, &nama, &catatan); err != nil {
			return nil, fmt.Errorf("repository: %s: %w", tabel, err)
		}
		hasil = append(hasil, models.BarisAkumulasi{ID: id.String, AccumulationName: nama.String, Note: catatan.String})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: %s: %w", tabel, err)
	}
	return hasil, nil
}
