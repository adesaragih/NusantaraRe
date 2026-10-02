package repository

import (
	"strings"
	"testing"
)

// TestSQLAkun - teks SQL tiket 27: PERSIS lima kolom DDL T_M_ACCOUNT, tabel yang sudah
// dikualifikasi, saringan "mengandung" PEKA huruf (tanpa UPPER, keputusan work owner
// 02-10-2026) di tiga kolom tampil (A69)
// dengan ESCAPE, urutan deterministik (A72), paging OFFSET/FETCH; semua masukan lewat
// parameter terikat bernomor.
func TestSQLAkun(t *testing.T) {
	pilih := "SELECT ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME FROM UJI.T_M_ACCOUNT"
	saring := ` WHERE (INSUREDID LIKE :1 ESCAPE '\' OR INSUREDNAME LIKE :2 ESCAPE '\' OR GROUPBUSINESS LIKE :3 ESCAPE '\')`
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
		if strings.Contains(s, "SELECT *") || strings.Contains(s, "UPPER") || strings.Contains(s, "LOWER") {
			t.Errorf("SELECT * atau pengubah huruf (pencarian harus peka huruf): %q", s)
		}
	}
	if TabelAkun != "T_M_ACCOUNT" {
		t.Error("nama tabel salah")
	}
}

// TestPolaCari - uji instrumen dengan jawaban yang diketahui: huruf APA ADANYA (peka
// huruf, keputusan work owner 02-10-2026 - "uji" dan "UJI" pola berbeda), `\` `%` `_`
// diloloskan sehingga masukan "50%" tidak menjadi wildcard, kosong/spasi = tanpa saringan.
func TestPolaCari(t *testing.T) {
	for masuk, mau := range map[string]string{
		"uji-a":    "%uji-a%",
		"UJI-A":    "%UJI-A%",
		" Uji ":    "%Uji%",
		"50%":      `%50\%%`,
		"a_b":      `%a\_b%`,
		`c:\d`:     `%c:\\d%`,
		"":         "",
		"   ":      "",
		"%_":       `%\%\_%`,
		"uji'drop": "%uji'drop%", // kutip tunggal tetap data: nilainya parameter terikat
	} {
		if got := PolaCari(masuk); got != mau {
			t.Errorf("PolaCari(%q) = %q, mau %q", masuk, got, mau)
		}
	}
}
