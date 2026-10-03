package repository

import (
	"strings"
	"testing"
)

// TestSQLAkun - teks SQL tiket 27: PERSIS lima kolom DDL T_M_ACCOUNT, tabel yang sudah
// dikualifikasi, saringan "mengandung" TIDAK peka huruf (UPPER di sisi kolom, butir 75)
// di tiga kolom tampil (A69) dengan ESCAPE, urutan deterministik (A72), paging
// OFFSET/FETCH; semua masukan lewat parameter terikat bernomor.
func TestSQLAkun(t *testing.T) {
	pilih := "SELECT ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME FROM UJI.T_M_ACCOUNT"
	saring := ` WHERE (UPPER(INSUREDID) LIKE :1 ESCAPE '\' OR UPPER(INSUREDNAME) LIKE :2 ESCAPE '\' OR UPPER(GROUPBUSINESS) LIKE :3 ESCAPE '\')`
	for _, u := range []struct{ got, mau string }{
		{sqlCariAkun("UJI.T_M_ACCOUNT", false), pilih + " ORDER BY INSUREDID, ID OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY"},
		{sqlCariAkun("UJI.T_M_ACCOUNT", true), pilih + saring + " ORDER BY INSUREDID, ID OFFSET :4 ROWS FETCH NEXT :5 ROWS ONLY"},
		{sqlCacahAkun("UJI.T_M_ACCOUNT", false), "SELECT COUNT(*) FROM UJI.T_M_ACCOUNT"},
		{sqlCacahAkun("UJI.T_M_ACCOUNT", true), "SELECT COUNT(*) FROM UJI.T_M_ACCOUNT" + saring},
	} {
		if u.got != u.mau {
			t.Errorf("SQL\n%q\nmau\n%q", u.got, u.mau)
		}
	}
	for _, s := range []string{sqlCariAkun("X", true), sqlCacahAkun("X", true)} {
		if strings.Contains(s, "SELECT *") || strings.Count(s, "UPPER(") != 3 || strings.Contains(s, "LOWER") {
			t.Errorf("SELECT *, atau bukan tepat tiga UPPER kolom (tidak peka huruf): %q", s)
		}
	}
	if TabelAkun != "T_M_ACCOUNT" {
		t.Error("nama tabel salah")
	}
}

// TestPolaCari - uji instrumen dengan jawaban yang diketahui: huruf BESAR (tidak peka
// huruf, butir 75 - "uji", "Uji", "UJI" berpola sama), `\` `%` `_` diloloskan sehingga
// masukan "50%" tidak menjadi wildcard, kosong/spasi = tanpa saringan.
func TestPolaCari(t *testing.T) {
	for masuk, mau := range map[string]string{
		"uji-a":    "%UJI-A%",
		"UJI-A":    "%UJI-A%",
		" Uji ":    "%UJI%",
		"50%":      `%50\%%`,
		"a_b":      `%A\_B%`,
		`c:\d`:     `%C:\\D%`,
		"":         "",
		"   ":      "",
		"%_":       `%\%\_%`,
		"uji'drop": "%UJI'DROP%", // kutip tunggal tetap data: nilainya parameter terikat
	} {
		if got := PolaCari(masuk); got != mau {
			t.Errorf("PolaCari(%q) = %q, mau %q", masuk, got, mau)
		}
	}
}
