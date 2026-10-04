package repository

// Bentuk SQL Accounts TANPA Oracle: kolom, bind, dan nol tulisan ke tabel yang dibaca saja.

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func rata(s string) string { return strings.Join(strings.Fields(s), " ") }

func TestSQLAccounts(t *testing.T) {
	for _, k := range []struct {
		nama  string
		sql   string
		wajib []string
	}{
		{"daftar", sqlDaftar("S.T_M_ACCOUNT"), []string{"FROM S.T_M_ACCOUNT WHERE (:1 IS NULL OR UPPER(INSUREDNAME) LIKE :2",
			"ORDER BY TO_NUMBER(REGEXP_SUBSTR(ID, '[0-9]+$')) DESC NULLS LAST, ID DESC", "OFFSET :6 ROWS FETCH NEXT :7 ROWS ONLY"}},
		{"hitung", sqlHitung("S.T_M_ACCOUNT"), []string{"SELECT COUNT(*) FROM S.T_M_ACCOUNT WHERE (:1 IS NULL"}},
		{"sisip", sqlSisip("S.T_M_ACCOUNT"), []string{"INSERT INTO S.T_M_ACCOUNT (ID, GROUPBUSINESSID, GROUPBUSINESS, INSUREDID, INSUREDNAME, DESCRIPTION, CREATEOP, CREATEDATE)",
			"VALUES (:1, :2, :3, :4, :5, :6, :7, SYSDATE)"}},
		{"pasangan lain", sqlPasanganLain("S.T_M_ACCOUNT"), []string{"WHERE INSUREDID = :1 AND GROUPBUSINESSID = :2 ORDER BY ID"}},
		{"nomor", sqlNomor("S.SEQ_T_M_ACCOUNT"), []string{"SELECT S.SEQ_T_M_ACCOUNT.NEXTVAL FROM DUAL"}},
		{"organisasi", sqlCariOrganisasi("S.CLIENT"), []string{"SELECT ID, IDVIEW, NAME FROM S.CLIENT WHERE FLAG = :1 AND NAME IS NOT NULL", "FETCH FIRST :5 ROWS ONLY"}},
		{"ambil organisasi", sqlAmbilOrganisasi("S.CLIENT"), []string{"WHERE ID = :1 AND FLAG = :2"}},
		{"group business", sqlDaftarGroupBusiness("S.BUSINESSGROUP"), []string{"SELECT ID, NOTE FROM S.BUSINESSGROUP WHERE NOTE IS NOT NULL ORDER BY UPPER(NOTE), ID"}},
	} {
		for _, w := range k.wajib {
			if !strings.Contains(rata(k.sql), w) {
				t.Errorf("%s: SQL tanpa %q:\n%s", k.nama, w, rata(k.sql))
			}
		}
		if err := db.PeriksaSQL(k.sql); err != nil {
			t.Errorf("%s: %v", k.nama, err)
		}
		if strings.Contains(strings.ToUpper(k.sql), "COMMIT") {
			t.Errorf("%s: COMMIT di SQL", k.nama)
		}
	}
	// Tabel yang dibaca saja tidak pernah ditulis.
	tulis := regexp.MustCompile(`(?i)\b(INSERT\s+INTO|UPDATE|DELETE|MERGE)\b`)
	for _, s := range []string{sqlCariOrganisasi("CLIENT"), sqlAmbilOrganisasi("CLIENT"), sqlDaftarGroupBusiness("BUSINESSGROUP"),
		sqlAmbilGroupBusiness("BUSINESSGROUP")} {
		if tulis.MatchString(s) {
			t.Errorf("menulis tabel baca-saja: %s", rata(s))
		}
	}
}

// Kata cari: tanpa beda huruf, karakter LIKE di-escape, kosong = NULL (tanpa saringan).
func TestPolaCari(t *testing.T) {
	if p := pola(` a%b_c\ `); p != `%A\%B\_C\\%` {
		t.Errorf("pola: %q", p)
	}
	if a := argCari("  "); a[0] != nil || len(a) != 5 {
		t.Errorf("kosong: %v", a)
	}
}
