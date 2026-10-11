package repository

import (
	"strings"
	"testing"
)

// TestKolomProduksiLife - kolom = daftar INSERT SaveLifeinProduction_SQL tanpa LAYER_1..4
// (keputusan work owner 06-10-2026); tiap kolom punya ekspresi.
func TestKolomProduksiLife(t *testing.T) {
	if len(KolomProduksiLife) != 32 {
		t.Errorf("kolom LIFEINPRODUCTION = %d, mau 32 (36 SaveLifeinProduction_SQL - LAYER_1..4)", len(KolomProduksiLife))
	}
	lihat := map[string]bool{}
	for _, k := range KolomProduksiLife {
		if lihat[k] {
			t.Errorf("kolom %s ganda", k)
		}
		lihat[k] = true
		if strings.HasPrefix(k, "LAYER_") || k == "JSONDATA" {
			t.Errorf("kolom %s tidak boleh diisi", k)
		}
		if nilaiProduksiLife[k] == "" {
			t.Errorf("kolom %s tanpa ekspresi nilai", k)
		}
	}
	if len(nilaiProduksiLife) != len(KolomProduksiLife) {
		t.Errorf("ekspresi %d ≠ kolom %d", len(nilaiProduksiLife), len(KolomProduksiLife))
	}
}

// TestSqlProduksiLife - DELETE + INSERT..SELECT dari baris utama kepala kasus, nol COMMIT.
func TestSqlProduksiLife(t *testing.T) {
	h := sqlHapusProduksiLife("SKEMAUJI.LIFEINPRODUCTION")
	if h != "DELETE FROM SKEMAUJI.LIFEINPRODUCTION WHERE IDPEGA = :1" {
		t.Errorf("hapus = %q", h)
	}
	q := sqlSisipProduksiLife("SKEMAUJI.LIFEINPRODUCTION", "SKEMAUJI.T_PREMIUM_LIST")
	for _, k := range []string{
		"INSERT INTO SKEMAUJI.LIFEINPRODUCTION (IDPEGA, NOPOLIS, BUSINESSCODE",
		"SELECT :1, :2, p.BUSINESS_CODE",
		"FROM SKEMAUJI.T_PREMIUM_LIST p WHERE p.ID = :3 AND p.ID = p.ID_PEGA",
		"TRUNC(p.DATE_RECEIVED)", "TRUNC(p.WPC)", "p.RI_SLIP_RNM", "p.CREATE_OP_NAME",
		"CASE p.TYPE_CEDING WHEN '1' THEN 'QS' WHEN '2' THEN 'SURPLUS' WHEN '3' THEN 'QS + SURPLUS' WHEN '4' THEN 'XOL' END",
		"PREMIUMREFUNDFACTOR, TGL_INPUT)", "p.PREMIUM_REFUND_FACTOR, SYSDATE FROM",
	} {
		if !strings.Contains(q, k) {
			t.Errorf("sisip tanpa %q:\n%s", k, q)
		}
	}
	for _, larangan := range []string{"SELECT *", "COMMIT", "LAYER_", "JSONDATA", "INSERTJSONOFFERLIFE", "NEXTVAL"} {
		if strings.Contains(strings.ToUpper(q+h), larangan) {
			t.Errorf("SQL memuat %q", larangan)
		}
	}
}
