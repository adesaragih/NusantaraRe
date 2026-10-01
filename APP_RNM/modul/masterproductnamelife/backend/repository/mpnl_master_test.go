package repository

// Teks SQL ketujuh pemilih master mengikuti RD Pega (PARITAS §4) - tanpa Oracle.

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/masterproductnamelife/backend/models"
)

func TestSQLMasterMengikutiRD(t *testing.T) {
	kasus := []struct {
		jenis models.JenisMaster
		wajib []string
		binds int
	}{
		// BrowseCedingCoLife_RD: B AND A AND C - ID Contains "L0" b570, ClientName Contains b587, StatusActive = 1 b607;
		// urut ClientName ASC b747; maks 10000 b723.
		{models.MasterCeding, []string{"FROM S.AGENT", "ID LIKE :1", "UPPER(CLIENTNAME) LIKE :2", "STATUSACTIVE = :3",
			"ORDER BY CLIENTNAME ASC", "FETCH FIRST 10000 ROWS ONLY"}, 3},
		{models.MasterSOB, []string{"FROM S.AGENT", "ID LIKE :1", "UPPER(CLIENTNAME) LIKE :2", "STATUSACTIVE = :3"}, 3},
		// BrowseClientNusaRe_RD: Name Contains b556, Name != "-" b570, Name IS NOT NULL b597; urut Name, BU_Note; maks 100000.
		{models.MasterPemegangPolis, []string{"FROM S.CLIENT", "UPPER(NAME) LIKE :1", "NAME <> :2", "NAME IS NOT NULL",
			"ORDER BY NAME ASC, BU_NOTE ASC", "FETCH FIRST 100000 ROWS ONLY"}, 2},
		// BrowseCurrencyLIFE_RD: Currency != "ITL" b541, Currency Contains b553; urut Currency ASC b757; maks 500.
		{models.MasterMataUang, []string{"FROM S.CURRENCY", "CURRENCY <> :1", "UPPER(CURRENCY) LIKE :2",
			"ORDER BY CURRENCY ASC", "FETCH FIRST 500 ROWS ONLY"}, 2},
		// BrowseRIRiskSummary: USEDBY Contains Param.SearchUsedby b555; urut ID ASC b682.
		{models.MasterRIRisk, []string{"FROM S.RIRISK_LIFE_SUMMARY", "UPPER(USEDBY) LIKE :1", "ORDER BY ID ASC",
			"FETCH FIRST 500 ROWS ONLY"}, 1},
		// BrowseCauseofLossLife_RD: CauseofLoss Contains b502; urut ID ASC b598.
		{models.MasterPenyebab, []string{"FROM S.CAUSEOFLOSS_LIFE", "UPPER(CAUSEOFLOSS) LIKE :1", "ORDER BY ID ASC"}, 1},
	}
	for _, k := range kasus {
		s, ada := sumberMaster[k.jenis]
		if !ada {
			t.Fatalf("%s: tanpa sumber", k.jenis)
		}
		q := rata(s.sqlCari("S." + s.objek))
		for _, w := range k.wajib {
			if !strings.Contains(q, w) {
				t.Errorf("%s: SQL tanpa %q:\n%s", k.jenis, w, q)
			}
		}
		if n := len(s.argCari("X")); n != k.binds {
			t.Errorf("%s: %d bind, mau %d", k.jenis, n, k.binds)
		}
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("%s: %v", k.jenis, err)
		}
		if !strings.Contains(rata(s.sqlAmbil("S."+s.objek)), "ID = :1") {
			t.Errorf("%s: pembaca satu nilai tidak dikunci ID", k.jenis)
		}
	}
	if _, ada := sumberMaster[models.MasterRIRate]; ada {
		t.Error("R/I Rate menunggu OQ-MPNL-03 - sumbernya tidak boleh dibaca")
	}
}

func TestSaringanMasterVerbatimRD(t *testing.T) {
	ag := sumberMaster[models.MasterCeding].argCari(PolaCari("uji"))
	if ag[0] != "%L0%" || ag[1] != "%UJI%" || ag[2] != "1" {
		t.Errorf("AGENT: %v", ag)
	}
	if a := sumberMaster[models.MasterPemegangPolis].argCari("A"); a[1] != "-" {
		t.Errorf("CLIENT Name != \"-\": %v", a)
	}
	if a := sumberMaster[models.MasterMataUang].argCari("A"); a[0] != "ITL" {
		t.Errorf("CURRENCY != ITL: %v", a)
	}
}

func TestPolaCariMeloloskanWildcard(t *testing.T) {
	if got := PolaCari(" a%b_c\\ "); got != `%A\%B\_C\\%` {
		t.Errorf("PolaCari = %q", got)
	}
	if PolaCari("") != "%%" {
		t.Error("kosong = semua (Contains '' Pega)")
	}
}

func TestGalatMasterMenyebutObjek(t *testing.T) {
	err := galatMaster("CLIENT", errors.New("ORA-00942: table or view does not exist"))
	if !errors.Is(err, ErrMasterTidakTerbaca) {
		t.Fatal("harus ErrMasterTidakTerbaca")
	}
	var l interface{ PesanLayar() string }
	if !errors.As(err, &l) || !strings.Contains(l.PesanLayar(), "CLIENT") || strings.Contains(l.PesanLayar(), "ORA-") {
		t.Errorf("kalimat layar menyebut objek, tanpa sebab Oracle: %v", err)
	}
}
