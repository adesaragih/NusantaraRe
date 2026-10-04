package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLCariAkumulasi - tiket 46: SearchRiskAccumulation_RD - INNER JOIN RW atas ZIPCODE, DISTINCT sebelas kolom laporan,
// saringan hanya yang diisi dengan bind berurutan, NOTE Contains tidak peka huruf, batas bind terakhir.
func TestSQLCariAkumulasi(t *testing.T) {
	syarat, arg := saringRD(models.SaringAkumulasi{})
	if len(syarat) != 0 || len(arg) != 0 {
		t.Fatalf("kosong: %v %v", syarat, arg)
	}
	if q := sqlCariAkumulasi("UJI.A", "UJI.R", nil, 1); q != "SELECT ID, ACCUMULATIONNAME, NOTE FROM (SELECT DISTINCT a.ACCUMULATION, a.CZONE, "+
		"a.CZONEID, a.NOTE, a.KEYWORD, a.SCOPEAREA, a.ID, a.ACCUMULATIONNAME, a.ZIPCODE, a.PROVINCE, a.PROVINCEID FROM UJI.A a JOIN UJI.R r "+
		"ON r.ZIPCODE = a.ZIPCODE) ORDER BY ID, NOTE, ACCUMULATIONNAME FETCH FIRST :1 ROWS ONLY" {
		t.Errorf("tanpa saringan: %q", q)
	}
	syarat, arg = saringRD(models.SaringAkumulasi{ID: "UJI-1", Note: "jl_a", CZone: "Z1", Keyword: "K", PostalCode: "12345", ProvinceID: "P1",
		PolicyNo: "diabaikan", SyariahStatus: "diabaikan", CityID: "diabaikan"})
	q := sqlCariAkumulasi("UJI.A", "UJI.R", syarat, len(arg)+1)
	mau := ` WHERE a.ID = :1 AND UPPER(a.NOTE) LIKE :2 ESCAPE '\' AND a.CZONE = :3 AND a.KEYWORD = :4 AND a.ZIPCODE = :5 AND a.PROVINCEID = :6) ` +
		`ORDER BY ID, NOTE, ACCUMULATIONNAME FETCH FIRST :7 ROWS ONLY`
	if !strings.HasSuffix(q, mau) || len(arg) != 6 || arg[1] != `%JL\_A%` || arg[4] != "12345" {
		t.Errorf("penuh: %q %v", q, arg)
	}
}

// TestSQLAkumulasiSQL - jalur RDB-List: kota / kecamatan persis SQL korpus (CARI2 = ACCUMULATIONTYPE), nomor polis lewat
// JSON_TABLE jalur korpus tanpa prosedur.
func TestSQLAkumulasiSQL(t *testing.T) {
	if q := sqlAkumulasiWilayah("UJI.A", "UJI.R", false); q != "SELECT ID, ACCUMULATIONTYPE, NOTE FROM UJI.A WHERE ZIPCODE IN (SELECT ZIPCODE FROM UJI.R WHERE CITYID = :1) ORDER BY ID, NOTE FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("kota: %q", q)
	}
	if q := sqlAkumulasiWilayah("UJI.A", "UJI.R", true); !strings.Contains(q, "WHERE DISTRICTID = :1)") {
		t.Errorf("kecamatan: %q", q)
	}
	q := sqlAkumulasiPolis("UJI.J", "UJI.A")
	for _, harus := range []string{"FROM UJI.J p, JSON_TABLE(p.DATA_JSON, '$.LocationList[*]'", "NESTED PATH '$.Property.PropertyItemList[*]'",
		"NESTED PATH '$.CoverageList[*]'", "PATH '$.AccumulationCode'", "WHERE p.NOPOLIS = :1 AND jt.ACCUMULATIONCODE IS NOT NULL",
		"SELECT MAX(x.ACCUMULATIONNAME) FROM UJI.A x WHERE x.ID = j.ACCUMULATIONCODE", "FETCH FIRST :2 ROWS ONLY"} {
		if !strings.Contains(q, harus) {
			t.Errorf("polis tanpa %q: %s", harus, q)
		}
	}
	if strings.Contains(strings.ToUpper(q), "CALL ") || strings.Contains(strings.ToUpper(q), "BEGIN") {
		t.Error("ADR-0043: nol prosedur")
	}
}

// TestSQLSaranAkumulasi - autocomplete: induk / kata / saringan tetap hanya bila ada, bind berurutan, kata atas kolom label.
func TestSQLSaranAkumulasi(t *testing.T) {
	for _, u := range []struct {
		jenis       string
		induk, kata bool
		mau         string
	}{
		{"city", true, true, "SELECT ID, NOTE, '' FROM (SELECT DISTINCT ID, NOTE FROM UJI.T WHERE PROVINCEID = :1 AND UPPER(NOTE) LIKE :2 ESCAPE '\\') ORDER BY NOTE, '', ID FETCH FIRST :3 ROWS ONLY"},
		{"city", false, false, "SELECT ID, NOTE, '' FROM (SELECT DISTINCT ID, NOTE FROM UJI.T) ORDER BY NOTE, '', ID FETCH FIRST :1 ROWS ONLY"},
		{"district", true, false, "SELECT ID, DISTRICTNAME, '' FROM (SELECT DISTINCT ID, CITYID, DISTRICTNAME, CITYNAME FROM UJI.T WHERE CITYNAME = :1) ORDER BY DISTRICTNAME, '', ID FETCH FIRST :2 ROWS ONLY"},
		{"area", true, true, "SELECT '', NOTE, ZIPCODE FROM (SELECT DISTINCT ZIPCODE, CZONE, NOTE, DISTRICTNAME, CITYNAME, PROVINCENAME, NATION FROM UJI.T WHERE STS_AKTIF = :1 AND UPPER(NOTE) LIKE :2 ESCAPE '\\') ORDER BY NOTE, ZIPCODE, '' FETCH FIRST :3 ROWS ONLY"},
	} {
		induk, pola := "", ""
		if u.induk {
			induk = "UJI-INDUK"
		}
		if u.kata {
			pola = "%UJI%"
		}
		q, arg := sqlSaran("UJI.T", daftarSaran[u.jenis], induk, pola)
		if arg[len(arg)-1] != BatasSaranAkumulasi || (u.kata && arg[len(arg)-2] != "%UJI%") || (u.jenis == "area" && arg[0] != "1") {
			t.Errorf("%s bind %v", u.jenis, arg)
		}
		if q != u.mau {
			t.Errorf("%s induk %v kata %v:\n%s\nmau\n%s", u.jenis, u.induk, u.kata, q, u.mau)
		}
	}
	for _, j := range []string{"city", "district", "area"} {
		if !JenisSaranDidukung(j) {
			t.Errorf("%s mestinya didukung", j)
		}
	}
	for _, j := range []string{"nation", "province", "accumtype", "czone", "x"} {
		if JenisSaranDidukung(j) {
			t.Errorf("%s mestinya belum didukung (DDL tidak ada)", j)
		}
	}
	if daftarSaran["city"].tabel != "CITY" || daftarSaran["district"].tabel != "DISTRICT" || daftarSaran["area"].tabel != "RW" ||
		TabelAccumulation != "ACCUMULATION" || TabelJSONPolis != "JSON_POLIS" || BatasAkumulasi != 500 || BatasSaranAkumulasi != 50 {
		t.Error("nama tabel DDL / batas berubah")
	}
}
