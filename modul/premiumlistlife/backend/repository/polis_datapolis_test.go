package repository

// Uji SQL data polis dan pencarian master Input Premium Detail. TANPA Oracle.

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

func TestSqlDataPolis(t *testing.T) {
	baca := sqlBacaDataPolis("S.T_PREMIUM_LIST")
	simpan := sqlSimpanDataPolis("S.T_PREMIUM_LIST")
	for _, q := range []string{baca, simpan} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Fatal(err)
		}
	}
	for _, k := range []string{"TO_CHAR(p.ANNUITY_INTEREST, 'TM9'", "TO_CHAR(p.WPC, 'YYYY", "WHERE p.ID = :1"} {
		if !strings.Contains(baca, k) {
			t.Errorf("baca tanpa %q", k)
		}
	}
	for _, k := range []string{"TYPE = :1", "PRO_RATE_TYPE = :11", "SECURITY_REINSURER = :20", "DATE_RECEIVED = :21", "WHERE ID = :22"} {
		if !strings.Contains(simpan, k) {
			t.Errorf("simpan tanpa %q", k)
		}
	}
	// Data penawaran (tanggal, Comment, dll.) dan nomor tidak tertimpa - KECUALI
	// DATE_RECEIVED, yang dapat diisi di layar ini (keputusan work owner 02-10-2026).
	for _, k := range []string{"DESCRIPTION", "TANGGAL_", "NO_POLIS", "WPC", "BUSINESS_CODE"} {
		if strings.Contains(simpan, k) {
			t.Errorf("simpan menimpa %s", k)
		}
	}
	if arg := argSimpanDataPolis("NBLF-1", models.IsianDataPolis{}); len(arg) != 22 || arg[21] != "NBLF-1" || arg[20] != nil || arg[14] != nil {
		t.Errorf("argumen %v", arg)
	}
}

func TestPencarianMasterPremiumDetail(t *testing.T) {
	for _, c := range []struct{ q, mau string }{
		{sqlCariProduk("S.PRODUCT_LIFE"), "WHERE CEDINGID = :1"},
		{sqlCariMarketing("S.MARKETINGOFFICER"), "MOSTATUS = :1"},
		{sqlCariRISlip("S.T_PREMIUM_LIST"), "ID = ID_PEGA AND (NO_POLIS LIKE :1 OR NO_POLIS LIKE :2)"},
	} {
		if err := db.PeriksaSQL(c.q); err != nil {
			t.Error(err)
		}
		if !strings.Contains(c.q, c.mau) {
			t.Errorf("tanpa %q:\n%s", c.mau, c.q)
		}
	}
	if PolaNomorRISlip[0] != "%RNML-Q%" || PolaNomorRISlip[1] != "%RNML-F%" {
		t.Error("pola GetPLandNopolis_sql berubah")
	}
}

func TestSqlBatasProdukDanPeserta(t *testing.T) {
	b := sqlBatasProduk("S.PRODUCTINWARD_LIFE")
	p := sqlPesertaBatas("S.T_PREMIUM_LIST_DETAIL")
	for _, q := range []string{b, p} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Fatal(err)
		}
	}
	if !strings.Contains(b, "WHERE ID = :1") || !strings.Contains(b, "TO_CHAR(MINAGE, 'TM9'") {
		t.Errorf("batas produk bukan GetRateProductLife:\n%s", b)
	}
	if !strings.Contains(p, "ORDER BY d.ID") {
		t.Error("urutan peserta bukan urutan grid - nomor \"at list\" menunjuk baris lain")
	}
}
