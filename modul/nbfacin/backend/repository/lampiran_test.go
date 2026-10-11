package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// TestSQLLampiran - SQL tiket 67: nol COMMIT / CALL, kolom DDL DOCUMENT_POLIS / T_STORAGE_IMAGE persis, hapus fisik dibatasi
// case, daftar dibatasi IDPEGA + 500 baris, kunci FOR UPDATE, pungutan per rujukan SKIP LOCKED.
func TestSQLLampiran(t *testing.T) {
	for nama, c := range map[string]struct {
		q   string
		ada []string
	}{
		"kategori": {sqlKategoriLampiran("S.C"), []string{"SELECT NOTE FROM S.C", "ORDER BY NOTE"}},
		"jumlah":   {sqlJumlahLampiran("S.D"), []string{"WHERE IDPEGA IN (:1, :2)", "GROUP BY KATEGORI_2"}},
		"daftar": {sqlDaftarLampiran("S.D", "S.O", "S.L", false), []string{"a.ID, a.NAMAFILE, a.MIME, a.KATEGORI_1, a.KATEGORI_2, " +
			"a.T_STORAGE_ID, a.PXCREATEOPERATOR", "TO_CHAR(a.TANGGAL, 'DD-MM-YYYY HH24:MI')", "s.IMAGEID = a.T_STORAGE_ID",
			"o.MODUL = :1 AND o.JENIS_EFEK = :2 AND o.RUJUKAN = a.ID", "a.KATEGORI_2 = :3 AND a.IDPEGA IN (:4, :5)", "FETCH FIRST 500 ROWS ONLY"}},
		"ambil": {sqlDaftarLampiran("S.D", "S.O", "S.L", true), []string{"a.ID = :3 AND a.IDPEGA IN (:4, :5)"}},
		"kunci": {sqlKunciLampiran("S.D", "S.O"), []string{"WHERE a.ID = :1 FOR UPDATE"}},
		"sisip": {sqlSisipLampiran("S.D"), []string{"(ID, IDPEGA, TANGGAL, NAMAFILE, MIME, KATEGORI_1, KATEGORI_2, INSKEY_LINK, " +
			"INSKEY_DATA, T_STORAGE_ID, PXCREATEOPERATOR) VALUES (:1, :2, SYSDATE, :3, :4, :5, :6, NULL, NULL, :7, :8)"}},
		"hapus":       {sqlHapusLampiran("S.D"), []string{"DELETE FROM S.D WHERE ID = :1 AND IDPEGA IN (:2, :3)"}},
		"link":        {sqlLinkStorage("S.O"), []string{"URLPUBLIC, APPFOLDER, TO_CHAR(EXPDATE, 'YYYYMMDDHH24MISS'), APPNAME, FILENAME", "IMAGEID = :1"}},
		"simpanObjek": {sqlSimpanObjekStorage("S.O"), []string{"(IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)", "TO_DATE(:4, 'DD/MM/YYYY HH24:MI:SS')", "'standard'"}},
		"perbarui":    {sqlPerbaruiObjekStorage("S.O"), []string{"TANGGAL_UPLOAD = TO_DATE(:4, 'MM/DD/YYYY HH24:MI:SS') WHERE IMAGEID = :5"}},
		"hapusObjek":  {sqlHapusObjekStorage("S.O"), []string{"DELETE FROM S.O WHERE IMAGEID = :1"}},
		"app":         {sqlAppStorage("S.F"), []string{"SELECT APPNAME FROM S.F"}},
		"token":       {sqlTokenStorageBerlaku("S.G"), []string{"WHERE APPNAME = :2 AND INPUTDATE > :3"}},
		"pungut": {sqlPungutEfekRujukan("S.L"), []string{"MODUL = :3 AND RUJUKAN = :4", "ORDER BY JADWAL_BERIKUT ASC) WHERE ROWNUM = 1)",
			"FROM S.L WHERE ID = (SELECT ID FROM (SELECT ID FROM S.L WHERE STATUS = :1", "FOR UPDATE SKIP LOCKED"}},
	} {
		if err := db.PeriksaSQL(c.q); err != nil || strings.Contains(strings.ToUpper(c.q), "CALL ") {
			t.Errorf("%s: %v", nama, err)
		}
		if strings.Contains(c.q, "FOR UPDATE") && strings.Contains(c.q, "FETCH FIRST") { // ORA-02014 (K67-8)
			t.Errorf("%s: FETCH FIRST bersama FOR UPDATE\n%s", nama, c.q)
		}
		for _, a := range c.ada {
			if !strings.Contains(c.q, a) {
				t.Errorf("%s: tanpa %q\n%s", nama, a, c.q)
			}
		}
	}
}
