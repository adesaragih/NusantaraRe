package repository

import (
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLRisk - tiket 36 (RD BrowseRisksAddress_RD): WHERE HANYA untuk saringan terisi, urut
// filter A..G, bind bernomor berurutan; tidak peka huruf (UPPER kolom + pola huruf besar) dengan
// ESCAPE; INNER JOIN RW pada kode pos; DISTINCT; 100 baris pertama urut sembilan kolom; paging.
func TestSQLRisk(t *testing.T) {
	dasar, arg := dasarRisk("UJI.A", "UJI.RW", models.SaringRisk{City: "jak%", Address: " uji_jalan "})
	mau := "SELECT DISTINCT a.ID, a.TITLE, a.ADDRESS, a.NATIONNAME, a.PROVINCENAME, a.CITYNAME, a.DISTRICTNAME, a.TERRITORYNAME," +
		" a.POSTALCODE FROM UJI.A a JOIN UJI.RW w ON w.ZIPCODE = a.POSTALCODE" +
		` WHERE UPPER(a.ADDRESS) LIKE :1 ESCAPE '\' AND UPPER(a.CITYNAME) LIKE :2 ESCAPE '\'`
	if dasar != mau {
		t.Errorf("SQL\n%q\nmau\n%q", dasar, mau)
	}
	if len(arg) != 2 || arg[0] != `%UJI\_JALAN%` || arg[1] != `%JAK\%%` {
		t.Errorf("argumen %v", arg)
	}
	semua, arg := dasarRisk("A", "RW", models.SaringRisk{Address: "a", ZipCode: "b", Country: "c", Province: "d", City: "e", District: "f", Territory: "g"})
	for i, k := range []string{"ADDRESS", "POSTALCODE", "NATIONNAME", "PROVINCENAME", "CITYNAME", "DISTRICTNAME", "TERRITORYNAME"} {
		if !strings.Contains(semua, "UPPER(a."+k+") LIKE :"+string(rune('1'+i))) {
			t.Errorf("saringan %s bukan :%d", k, i+1)
		}
	}
	if len(arg) != 7 || strings.Count(semua, " AND ") != 6 {
		t.Errorf("tujuh saringan: %d argumen, %q", len(arg), semua)
	}
	if kosong, arg := dasarRisk("A", "RW", models.SaringRisk{Territory: "   "}); strings.Contains(kosong, "WHERE") || len(arg) != 0 {
		t.Errorf("saringan spasi saja harus dilewati: %q", kosong)
	}
	halaman := sqlCariRisk(dasar, 2)
	if !strings.HasPrefix(halaman, "SELECT ID, TITLE, ADDRESS, NATIONNAME, PROVINCENAME, CITYNAME, DISTRICTNAME, TERRITORYNAME, POSTALCODE FROM (") ||
		!strings.HasSuffix(halaman, " ORDER BY "+kolomRisk+" FETCH FIRST 100 ROWS ONLY) ORDER BY ID, TITLE, ADDRESS, NATIONNAME, PROVINCENAME, CITYNAME,"+
			" DISTRICTNAME, TERRITORYNAME, POSTALCODE OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY") || strings.Contains(halaman, "SELECT *") {
		t.Errorf("halaman %q", halaman)
	}
	if got := sqlCacahRisk("X"); got != "SELECT COUNT(*) FROM (X FETCH FIRST 100 ROWS ONLY)" {
		t.Errorf("cacah %q", got)
	}
}
