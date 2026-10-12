package repository

// Uji SQL spreading retro (keputusan work owner 06-10-2026). TANPA Oracle.

import (
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/premiumlistlife/backend/models"
)

func TestSqlBacaBahanSpreading(t *testing.T) {
	for nama, k := range map[string]struct{ q, mau string }{
		"begin produk": {sqlBeginProduk("S.M_PRODUCTNAME_LIFE"), "SELECT TO_CHAR(BEGIN_DATE, 'YYYY-MM-DD') FROM S.M_PRODUCTNAME_LIFE WHERE ID = :1"},
		"tahun":        {sqlTahunRetro("S.TREATYYEAR_LIFE"), "SELECT ID, TO_CHAR(STARTDATE, 'YYYY-MM-DD'), TO_CHAR(ENDDATE, 'YYYY-MM-DD') FROM S.TREATYYEAR_LIFE"},
		"kontrak":      {sqlKontrakRetro("S.TREATYCONTRACT_LIFE"), "FROM S.TREATYCONTRACT_LIFE WHERE IDTREATYYEAR = :1"},
		"reinsurer":    {sqlReinsurerRetro("S.TREATYREINSURER_LIFE"), "SELECT TREATYCONTRACTID, REINSURERNAME, TO_CHAR(PCTSHARE, 'TM9'"},
		"business":     {sqlBusinessRetro("S.TREATYBUSINESS_LIFE"), "SELECT TREATYCONTRACTID, BIZNAME, RIRATEID, RIRATE FROM S.TREATYBUSINESS_LIFE WHERE TREATYYEARID = :1"},
	} {
		if err := db.PeriksaSQL(k.q); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(k.q, k.mau) {
			t.Errorf("%s:\n%s\nmau memuat %s", nama, k.q, k.mau)
		}
		if strings.Contains(k.q, "*") {
			t.Errorf("%s: SELECT *", nama)
		}
	}
	if q := sqlKontrakRetro("S.K"); !strings.Contains(q, "TO_CHAR(IDR, 'TM9'") || !strings.Contains(q, "TO_CHAR(TREATYSTARTDATE, 'YYYY-MM-DD')") {
		t.Errorf("kontrak tanpa limit / tanggal teks: %s", q)
	}
}

func TestSqlTulisDanBacaSpreading(t *testing.T) {
	qs, qr := sqlSisipSpreading("S.T_PREMIUM_LIST_SPREADING"), sqlSisipSpreadingRetro("S.T_PREMIUM_LIST_SPREADING_RETRO")
	for _, q := range []string{qs, qr, sqlHapusSpreadingRetro("S.R", "S.S", "S.D"), sqlHapusSpreading("S.S", "S.D"),
		sqlIdentitasPesertaSpreading("S.D"), sqlRincianSpreadingPeserta("S.S", "S.R", "S.Y"), sqlRekapRetro("S.D", "S.S", "S.R")} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Fatal(err)
		}
	}
	l := models.LapisanSpreading{TahunID: "T2", Bagian: apd.New(15000000005, -1),
		Kontrak: models.KontrakRetro{ReinsTypeID: "R1", ReinsTypeName: "QS", IDR: "3000000000", Mulai: "2025-01-01", Akhir: "2036-12-31"}}
	if n, args := strings.Count(qs, ":"), NilaiSisipSpreading("S1", "D1", "UJI-OP", l); n != len(args) {
		t.Errorf("sisip spreading: %d penanda, %d nilai", n, len(args))
	} else if args[5] != "1500000000.5000" || args[7] != nil || args[12] != "UJI-OP" {
		t.Errorf("nilai spreading = %v", args)
	}
	x := models.ReinsurerSpreading{Nama: "UJI-RE", PctShare: apd.New(60, 0), Komisi: apd.New(20, 0), Ovr: apd.New(10, 0),
		Amount: apd.New(1, 0), Rate: apd.New(25, -1), Gross: apd.New(123456789, -5), Net: apd.New(1, 0)}
	if n, args := strings.Count(qr, ":"), NilaiSisipSpreadingRetro("R1", "S1", "UJI-OP", l, x); n != len(args) {
		t.Errorf("sisip retro: %d penanda, %d nilai", n, len(args))
	} else if args[7] != "2.5000" || args[8] != "1234.5679" || args[12] != "2025-01-01" {
		t.Errorf("nilai retro = %v", args)
	}
	if !strings.Contains(sqlRekapRetro("S.D", "S.S", "S.R"), "GROUP BY r.REINSURER_NAME, d.CURRENCY") {
		t.Error("rekap tidak per reinsurer per mata uang")
	}
	for _, id := range []string{PengenalSpreading("NBLF-1", 3, 2), PengenalSpreadingRetro("NBLF-1", 3, 2, 1)} {
		if len(id) != 32 {
			t.Errorf("ID %q bukan 32 karakter", id)
		}
	}
	if PengenalSpreading("NBLF-1", 3, 2) == PengenalSpreading("NBLF-2", 3, 2) {
		t.Error("ID spreading antar polis bertabrakan")
	}
}
