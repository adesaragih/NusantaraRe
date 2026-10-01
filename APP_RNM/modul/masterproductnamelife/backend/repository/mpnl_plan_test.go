package repository

import (
	"strings"
	"testing"
)

func TestSQLMasterPlanBrowseProductTypeLife(t *testing.T) {
	q := rata(sqlCariPlan("S.PRODUCT_TYPE_LIFE"))
	for _, w := range []string{"SELECT ID, COVERNAME, BUSINESS, BENEFIT FROM S.PRODUCT_TYPE_LIFE",
		"UPPER(COVERNAME) LIKE :1", "OR UPPER(BUSINESS) LIKE :2", "FETCH FIRST 500 ROWS ONLY"} {
		if !strings.Contains(q, w) {
			t.Errorf("SQL tanpa %q: %s", w, q)
		}
	}
	if !strings.Contains(rata(sqlAmbilPlan("S.PRODUCT_TYPE_LIFE")), "WHERE ID = :1") {
		t.Error("ambil plan dikunci ID")
	}
}
