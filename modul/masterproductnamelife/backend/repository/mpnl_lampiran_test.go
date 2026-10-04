package repository

// Teks SQL lampiran mengikuti RDB Pega (PARITAS §6) - tanpa Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSQLLampiranMengikutiRDB(t *testing.T) {
	kasus := []struct {
		nama  string
		sql   string
		wajib []string
	}{
		{"ID InsertAttachProdName_Sql b84", sqlWaktuID(),
			[]string{"TO_CHAR(SYSTIMESTAMP, 'YYYYMMDDHH24MISSFF3')"}},
		{"InsertAttachProdName_Sql b84", sqlSisipLampiran("S.M_ATTACHMENTPRODUCTNAME"),
			[]string{"INSERT INTO S.M_ATTACHMENTPRODUCTNAME (ID, TREATYID, CATEGORY, FILENAME, FILEMIMETYPE, DATA_JSON, USERNAME, T_STORAGE_ID)",
				"VALUES (:1, :2, :3, :4, :5, NULL, :6, :7)"}},
		{"GetAttachmentProdName_Sql b84", sqlDaftarLampiran("S.M_ATTACHMENTPRODUCTNAME"),
			[]string{"WHERE TREATYID = :1", "ORDER BY ID ASC"}},
		{"DeleteAttachProdName_Sql b85", sqlHapusLampiran("S.M_ATTACHMENTPRODUCTNAME"),
			[]string{"DELETE FROM S.M_ATTACHMENTPRODUCTNAME WHERE TREATYID = :1 AND ID = :2"}},
		{"Insert_T_Storage_SQL b85", sqlSisipObjek("S.T_STORAGE_IMAGE"),
			[]string{"(IMAGEID, URLPUBLIC, APPFOLDER, EXPDATE, FILENAME, APPNAME, STORAGE)", "VALUES (:1, NULL,", "'standard')"}},
		{"DeleteStorage_SQL b85", sqlHapusObjek("S.T_STORAGE_IMAGE"),
			[]string{"DELETE FROM S.T_STORAGE_IMAGE WHERE IMAGEID = :1"}},
		{"GetAppName_SQL b58", sqlNamaAplikasi("S.T_FOLDER_IMAGE"),
			[]string{"SELECT APPNAME FROM S.T_FOLDER_IMAGE"}},
		{"outbox pungut", sqlPungutUnggah("S.T_LOG_SERVICE_RNM"),
			[]string{"WHERE MODUL = :1 AND RUJUKAN = :2 AND STATUS = :3", "FOR UPDATE SKIP LOCKED"}},
		{"outbox status", sqlStatusUnggah("S.T_LOG_SERVICE_RNM"),
			[]string{"WHERE MODUL = :1 AND RUJUKAN = :2", "ORDER BY DIBUAT DESC"}},
	}
	for _, k := range kasus {
		for _, w := range k.wajib {
			if !strings.Contains(rata(k.sql), w) {
				t.Errorf("%s: SQL tanpa %q:\n%s", k.nama, w, rata(k.sql))
			}
		}
		if err := db.PeriksaSQL(k.sql); err != nil {
			t.Errorf("%s: %v", k.nama, err)
		}
	}
	// ORA-02014: batas baris tidak boleh bersama FOR UPDATE.
	if strings.Contains(sqlPungutUnggah("S.T"), "FETCH FIRST") {
		t.Error("pungut outbox memakai FETCH FIRST bersama FOR UPDATE")
	}
	// URLPUBLIC selalu NULL: alamat penyimpanan nyata tidak ditulis (OQ-MPNL-10).
	if strings.Count(sqlSisipObjek("S.T"), ":") != 5 {
		t.Errorf("sisip objek mengikat lebih dari 5 nilai: %s", rata(sqlSisipObjek("S.T")))
	}
}
