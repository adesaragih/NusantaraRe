package repository

import (
	"reflect"
	"strings"
	"testing"
)

// TestSQLSaranZona - BrowseZones_RD: saringan kosong dibuang (pyUseNullIfEmpty false), UPPER di kedua sisi
// (pyCaseInsensitive), zona sebagai satu kelompok OR (C..H), bind berurutan dengan batas terakhir.
func TestSQLSaranZona(t *testing.T) {
	q, arg := sqlSaranZona("UJI.ZONES", "4.2 PRGBI", []string{"1", "2"}, "%UJI%")
	mau := `SELECT ID, DESCRIPTION FROM UJI.ZONES WHERE UPPER(ZONE) IN (UPPER(:1), UPPER(:2)) AND UPPER(OLDID) = UPPER(:3) AND ` +
		`UPPER(DESCRIPTION) LIKE :4 ESCAPE '\' ORDER BY DESCRIPTION, ID FETCH FIRST :5 ROWS ONLY`
	if q != mau || !reflect.DeepEqual(arg, []any{"1", "2", "4.2 PRGBI", "%UJI%", BatasSaranAkumulasi}) {
		t.Errorf("penuh:\n%s\n%v", q, arg)
	}
	q, arg = sqlSaranZona("UJI.ZONES", "", nil, "")
	if q != "SELECT ID, DESCRIPTION FROM UJI.ZONES ORDER BY DESCRIPTION, ID FETCH FIRST :1 ROWS ONLY" || len(arg) != 1 {
		t.Errorf("kosong: %s %v", q, arg)
	}
	q, arg = sqlSaranZona("UJI.ZONES", "FLEXAS", nil, "")
	if !strings.Contains(q, " WHERE UPPER(OLDID) = UPPER(:1) ORDER BY") || !reflect.DeepEqual(arg, []any{"FLEXAS", BatasSaranAkumulasi}) {
		t.Errorf("oldId saja: %s %v", q, arg)
	}
}

// TestSQLKonstruksiDanRate - SearchClassConstructionReas_SQL (DISTINCT, urut CONSTRUCTION) dan RateEarthquake_Sql
// (DESCRIPTION = zona persis, RATE lewat angkaKeluar, satu baris).
func TestSQLKonstruksiDanRate(t *testing.T) {
	q, arg := sqlSaranKonstruksi("UJI.CONSTRUCTION", "")
	if q != "SELECT CONSTRUCTION AS KODE, CONSTRUCTION AS NAMA FROM (SELECT DISTINCT CONSTRUCTION FROM UJI.CONSTRUCTION WHERE CONSTRUCTION IS NOT NULL) "+
		"ORDER BY KODE FETCH FIRST :1 ROWS ONLY" || len(arg) != 1 {
		t.Errorf("tanpa kata: %s %v", q, arg)
	}
	q, arg = sqlSaranKonstruksi("UJI.CONSTRUCTION", "%KELAS%")
	if !strings.Contains(q, `WHERE UPPER(CONSTRUCTION) LIKE :1 ESCAPE '\') ORDER BY KODE FETCH FIRST :2 ROWS ONLY`) ||
		!reflect.DeepEqual(arg, []any{"%KELAS%", BatasSaranAkumulasi}) {
		t.Errorf("dengan kata: %s %v", q, arg)
	}
	// ORA-00960 (DEV 06-10-2026): kolom yang dipilih dua kali wajib beralias berbeda supaya ORDER BY tidak ambigu.
	for _, pola := range []string{"", "%KELAS%"} {
		if q, _ := sqlSaranKonstruksi("UJI.CONSTRUCTION", pola); strings.Contains(q, "SELECT CONSTRUCTION, CONSTRUCTION") {
			t.Errorf("kolom ganda tanpa alias (ORA-00960): %s", q)
		}
	}
	if q := sqlRateZona("UJI.ZONEOJK"); q != "SELECT "+angkaKeluar("RATE")+" FROM UJI.ZONEOJK WHERE DESCRIPTION = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY" {
		t.Errorf("rate: %s", q)
	}
}
