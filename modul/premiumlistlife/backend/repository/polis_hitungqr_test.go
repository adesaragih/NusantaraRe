package repository

// Uji SQL bahan hitung Type QR dan kolom tersimpan Calculate CSV
// (keputusan work owner 05-10-2026). TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestSqlBahanHitungQR(t *testing.T) {
	for nama, k := range map[string]struct{ q, mau string }{
		"kepala": {sqlKepalaUnggah("S.T_PREMIUM_LIST"),
			"SELECT TYPE, PRODUCT_NAME_ID, BUSINESS_NAME, PRO_RATE_TYPE, RI_SLIP_RNM FROM S.T_PREMIUM_LIST WHERE ID = :1"},
		// Batas produk dari M_PRODUCTNAME_LIFE, bukan PRODUCTINWARD_LIFE (05-10-2026).
		"batas": {sqlBatasProdukMPNL("S.M_PRODUCTNAME_LIFE"),
			"TO_CHAR(MAXSUMINSURED, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''') FROM S.M_PRODUCTNAME_LIFE WHERE ID = :1"},
		"produk": {sqlParamProdukQR("S.M_PRODUCTNAME_LIFE"), "RIRISKID FROM S.M_PRODUCTNAME_LIFE WHERE ID = :1"},
		"plan":   {sqlPlanProdukQR("S.P"), "SELECT NAME, RIRATEID FROM S.P WHERE PRODUCTID = :1 ORDER BY URUT"},
		"rate":   {sqlRateHitungQR("S.RATE_LIFE"), "SELECT ID, GENDER, AGE, CONTRACT, RATE FROM S.RATE_LIFE WHERE IDUSEDBY = :1"},
		"risk":   {sqlRiskHitungQR("S.RIRISK_LIFE"), "SELECT ID, YEAR, CONTRACT, RISK FROM S.RIRISK_LIFE WHERE IDUSEDBY = :1"},
	} {
		if err := db.PeriksaSQL(k.q); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(k.q, k.mau) {
			t.Errorf("%s:\n%s\nmau memuat %s", nama, k.q, k.mau)
		}
		if strings.Contains(k.q, "*") || strings.Contains(k.q, "JSONDATA") {
			t.Errorf("%s: SELECT * / JSONDATA", nama)
		}
	}
	// Master rate/risk untuk hitungan TIDAK dipotong batas tampilan.
	for _, q := range []string{sqlRateHitungQR("S.R"), sqlRiskHitungQR("S.R")} {
		if strings.Contains(q, "FETCH FIRST") {
			t.Errorf("master hitungan dipotong: %s", q)
		}
	}
	if q := sqlParamProdukQR("S.M"); !strings.Contains(q, "TO_CHAR(CEDINGLIMIT, 'TM9'") {
		t.Errorf("angka produk tidak dibaca sebagai teks desimal: %s", q)
	}
}

// TestSisipPesertaMenulisKolomHasilQR - kolom yang dihitung Calculate CSV
// benar-benar ikut INSERT (migrasi 052, nol perubahan skema).
func TestSisipPesertaMenulisKolomHasilQR(t *testing.T) {
	kolom, _, _ := ekspresiSisipPeserta()
	ada := map[string]bool{}
	for _, k := range strings.Split(kolom, ", ") {
		ada[k] = true
	}
	for _, k := range []string{"RATE", "RISK", "SUM_AT_RISK_GROSS", "CEDING_RETENTION", "SUM_REASURED",
		"SHARE_NUSANTARA_RE", "GROSS_PREMIUM", "DEDUCTION", "BROKERAGE_FEE", "NET_PREMIUM", "RI_ADMIN_FEE",
		"PERIOD_MM", "FACTOR", "EM_PERCENT"} {
		if !ada[k] {
			t.Errorf("SisipPeserta tidak menulis %s", k)
		}
	}
}
