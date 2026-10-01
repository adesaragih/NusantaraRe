package repository

// Teks SQL pembaca mengikuti RD/RDB Pega (PARITAS §2, §7) - tanpa Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func rata(s string) string { return strings.Join(strings.Fields(s), " ") }

func TestSQLBacaMengikutiRD(t *testing.T) {
	kasus := []struct {
		nama  string
		sql   string
		wajib []string
	}{
		{"grid BrowseProduct_Life b1094", sqlDaftarProduk("S.M_PRODUCT_LIFE"),
			[]string{"SELECT ID, JSONDATA FROM S.M_PRODUCT_LIFE", "ORDER BY ID ASC"}},
		{"umum BrowseUnderwritingList b61", sqlAmbilProduk("S.M_PRODUCT_LIFE", false),
			[]string{"FROM S.M_PRODUCT_LIFE WHERE ID = :1"}},
		{"inward BrowseProductInward b834 + kunci Claim Life", sqlAmbilInward("S.M_PRODUCTINWARD_LIFE", false),
			[]string{"JSON_VALUE(JSONDATA, '$.PRODUCTID') = :1 OR ID = :2", "ORDER BY ID ASC"}},
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
		if strings.Contains(k.sql, "FOR UPDATE") {
			t.Errorf("%s: pembaca biasa tidak mengunci", k.nama)
		}
	}
	if !strings.HasSuffix(rata(sqlAmbilProduk("S.T", true)), "FOR UPDATE") ||
		!strings.HasSuffix(rata(sqlAmbilInward("S.T", true)), "FOR UPDATE") {
		t.Error("pembaca di dalam simpan harus mengunci barisnya")
	}
}
