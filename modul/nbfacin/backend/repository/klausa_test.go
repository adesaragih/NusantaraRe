package repository

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSQLKlausa - pencarian = RetrieveClauseSQL (decode bahasa judul / isi, TYPE, Info LIKE '%' || q || '%' tanpa ESCAPE)
// + K47-7 (kata tidak peka huruf; isi Inggris jatuh ke TEXTINA bila TEXTENG kosong), halaman urut ID + jumlah argumen RetrieveArgumentNumberClauseSQL, bind :1..:6 urut kemunculan; argumen =
// SearchClauseArgFireSQL persis; ClauseList satu tabel (baris per argumen): baca urut klausa lalu argumen, hapus satu
// pernyataan, sisip 16 kolom.
func TestSQLKlausa(t *testing.T) {
	dasar := sqlCariKlausa("UJI.CLAUSE")
	if dasar != `SELECT * FROM (SELECT ID, DECODE(:1, '1', TITLEENG, '2', TITLEDUAL, '0', TITLEINA) AS "Title", `+
		`DECODE(:2, '1', CASE WHEN TRIM(TEXTENG) IS NULL THEN TEXTINA ELSE TEXTENG END, '2', TEXTDUAL, '0', TEXTINA) AS "Text", `+
		`INFO AS "Info" FROM UJI.CLAUSE WHERE TYPE = :3) bca WHERE "Title" IS NOT NULL AND UPPER("Info") LIKE '%' || UPPER(:4) || '%'` {
		t.Errorf("cari: %q", dasar)
	}
	if h := sqlHitungKlausa(dasar); h != "SELECT COUNT(*) FROM ("+dasar+")" {
		t.Errorf("hitung: %q", h)
	}
	if h := sqlHalamanKlausa(dasar, "UJI.A", "UJI.M"); h != `SELECT h.ID, h."Title", h."Text", h."Info", (SELECT COUNT(*) FROM UJI.A a `+
		`WHERE a.OLDID = (SELECT m.OLDID FROM UJI.M m WHERE m.ID = h.ID)) FROM (`+dasar+`) h ORDER BY h.ID OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY` {
		t.Errorf("halaman: %q", h)
	}
	// K47-6: tanpa M_ARGCLAUSEFIRE - kolom jumlah NULL, bind tetap :5 / :6; katalog ditanya per nama, pemilik tak peka huruf.
	if h := sqlHalamanKlausaTanpaArgumen(dasar); h != `SELECT h.ID, h."Title", h."Text", h."Info", NULL FROM (`+dasar+
		`) h ORDER BY h.ID OFFSET :5 ROWS FETCH NEXT :6 ROWS ONLY` || strings.Contains(h, "UJI.A") {
		t.Errorf("halaman tanpa argumen: %q", h)
	}
	if sqlAdaObjek != "SELECT COUNT(*) FROM SYS.ALL_OBJECTS WHERE UPPER(OWNER) = UPPER(:1) AND OBJECT_NAME = :2" {
		t.Error("katalog tabel argumen")
	}
	if q := sqlArgumenKlausa("UJI.A", "UJI.C"); q != "SELECT a.jsondata.ArgumentNumber, a.jsondata.Description, a.jsondata.DefaultValue "+
		"FROM UJI.A a WHERE a.OLDID = (SELECT c.OLDID FROM UJI.C c WHERE c.ID = :1) ORDER BY a.jsondata.ArgumentNumber" {
		t.Errorf("argumen: %q", q)
	}
	if q := sqlBacaKlausa("UJI.CL"); !strings.HasSuffix(q, "FROM UJI.CL WHERE PARENT_ID = :1 ORDER BY SEQ_NO, ARGUMENT_SEQ_NO NULLS FIRST") ||
		!strings.HasPrefix(q, "SELECT SEQ_NO, ARGUMENT_SEQ_NO, CLAUSE_CODE,") || !strings.Contains(q, "ARGUMENT_NUMBER, ARGUMENT_DESCRIPTION, ARGUMENT_VALUE FROM") {
		t.Errorf("baca: %q", q)
	}
	if h := sqlHapusKlausa("UJI.CL"); h != "DELETE FROM UJI.CL WHERE PARENT_ID = :1" {
		t.Errorf("hapus: %q", h)
	}
	if q := sqlSisipKlausa("UJI.CL"); strings.Count(q, ":") != 16 || !strings.Contains(q, "(ID, PARENT_ID, SEQ_NO, ARGUMENT_SEQ_NO, ROW_UID,") ||
		!strings.Contains(q, "ARGUMENT_VALUE) VALUES") {
		t.Errorf("sisip: %q", q)
	}
}

// TestSatuTabelKlausa - K47-4 satu baris per argumen: pecah lalu kumpulkan kembali = daftar semula (klausa tanpa argumen
// satu baris UrutArg 0, argumen berurutan, kolom klausa berulang), daftar kosong = larik kosong.
func TestSatuTabelKlausa(t *testing.T) {
	a := models.KlausaKasus{Code: "K1", Title: "UJI SATU", ArgumentCount: "2", Arguments: []models.ArgumenKlausa{
		{Number: "1", Description: "UJI A", Value: "x"}, {Number: "2", Description: "UJI B", Value: "y"}}}
	b := models.KlausaKasus{Code: "K2", Title: "UJI DUA", ArgumentCount: "0", Arguments: []models.ArgumenKlausa{}}
	c := models.KlausaKasus{Code: "K3", Title: "UJI TIGA", ArgumentCount: "1", Arguments: []models.ArgumenKlausa{{Number: "1", Value: "z"}}}
	baris := pecahKlausa([]models.KlausaKasus{a, b, c})
	var urut []string
	for _, r := range baris {
		urut = append(urut, strings.Join([]string{r.Klausa.Code, string(rune('0' + r.Urut)), string(rune('0' + r.UrutArg))}, "/"))
	}
	if strings.Join(urut, " ") != "K1/1/1 K1/1/2 K2/2/0 K3/3/1" {
		t.Errorf("pecah: %v", urut)
	}
	if got := kumpulkanKlausa(baris); !reflect.DeepEqual(got, []models.KlausaKasus{a, b, c}) {
		t.Errorf("kumpulkan: %+v", got)
	}
	if got := kumpulkanKlausa(nil); got == nil || len(got) != 0 {
		t.Errorf("kosong: %#v", got)
	}
}
