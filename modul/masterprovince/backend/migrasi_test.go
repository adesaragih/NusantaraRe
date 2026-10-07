package backend

// Migrasi 883 - TANPA Oracle: bentuk SQL-nya.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func langkah(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// 883: kolom PROVINCENAME selebar RW.PROVINCENAME; isi dari RW dengan pencocokan PERSIS kata work owner
// (rw.CITYNAME = cityinput.note, tanpa UPPER / TRIM), hanya bila TEPAT SATU nama provinsi, tanpa menimpa; mundur
// membuang kolomnya.
func TestMigrasiProvinceNameCity(t *testing.T) {
	p := langkah(t, "883_cityinput_provincename.sql")
	if len(p) != 3 {
		t.Fatalf("883: %d pernyataan, mau ADD, UPDATE, INSERT", len(p))
	}
	if tabel, kolom := migrasi.KolomAlterTambah(p[0]); tabel != "CITYINPUT" || !reflect.DeepEqual(kolom, []string{"PROVINCENAME"}) ||
		!strings.Contains(p[0], "PROVINCENAME VARCHAR2(4000)") {
		t.Errorf("883 kolom: %s %v", tabel, kolom)
	}
	u := p[1]
	for _, mau := range []string{"UPDATE {skema}.CITYINPUT c", "FROM {skema}.RW r WHERE r.CITYNAME = c.NOTE",
		"WHERE c.PROVINCENAME IS NULL", "COUNT(DISTINCT r.PROVINCENAME)", ") = 1"} {
		if !strings.Contains(u, mau) {
			t.Errorf("883 UPDATE tanpa %q:\n%s", mau, u)
		}
	}
	if strings.Count(u, "r.CITYNAME = c.NOTE") != 2 || strings.Contains(strings.ToUpper(u), "UPPER(") ||
		strings.Contains(strings.ToUpper(u), "TRIM(") || strings.Contains(u, "STS_AKTIF") {
		t.Errorf("883 UPDATE melonggarkan / menyaring di luar kata work owner:\n%s", u)
	}
	// Kota baru dari RW: kunci banding persis, ID digit + ROW_NUMBER, PROVINCENAME aturan yang sama, penanda mundur.
	s := p[2]
	for _, mau := range []string{"INSERT INTO {skema}.CITYINPUT (ID, NOTE, PROVINCENAME, CREATE_OP, TGL_CREATE)",
		"ROW_NUMBER() OVER (ORDER BY b.CITYNAME)", "REGEXP_LIKE(c.ID, '^[0-9]+$')", "CASE WHEN b.NPROV = 1 THEN b.PROV END",
		"COUNT(DISTINCT r.PROVINCENAME) AS NPROV", "'MIGRASI-883'", "WHERE NOT EXISTS (SELECT 1 FROM {skema}.CITYINPUT x WHERE x.NOTE = b.CITYNAME)",
		"WHERE TRIM(r.CITYNAME) IS NOT NULL", "GROUP BY r.CITYNAME"} {
		if !strings.Contains(s, mau) {
			t.Errorf("883 INSERT tanpa %q:\n%s", mau, s)
		}
	}
	if strings.Contains(strings.ToUpper(s), "UPPER(") || strings.Contains(s, "STS_AKTIF") || strings.Contains(s, "ZIPCODE") {
		t.Errorf("883 INSERT di luar kata work owner:\n%s", s)
	}
	d := langkah(t, "883_cityinput_provincename_down.sql")
	if len(d) != 2 || d[0] != "DELETE FROM {skema}.CITYINPUT WHERE CREATE_OP = 'MIGRASI-883'" ||
		d[1] != "ALTER TABLE {skema}.CITYINPUT DROP (PROVINCENAME)" {
		t.Errorf("883 mundur: %q", d)
	}
}
