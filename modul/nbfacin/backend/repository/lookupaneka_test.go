package repository

import "testing"

// TestSQLLookupAneka - BrowseShip_RD: Remark != literal (bind) OR NULL, Name Contains tak peka huruf hanya bila diisi, ID DESC,
// batas; BrowseType_RD: TYPE = bind, ID DESC, batas; bind 1..n. Literal filter persis (peka huruf).
func TestSQLLookupAneka(t *testing.T) {
	if q := sqlCariKapal("S.SHIP", false); q != "SELECT ID, NAME, FLAGNAME, YEARMAKE1, NRT, GRT, DWT, CLASSIFICATION, CONSTRUCTION FROM S.SHIP WHERE (REMARK != :1 OR REMARK IS NULL) ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("tanpa kata: %q", q)
	}
	if q := sqlCariKapal("S.SHIP", true); q != `SELECT ID, NAME, FLAGNAME, YEARMAKE1, NRT, GRT, DWT, CLASSIFICATION, CONSTRUCTION FROM S.SHIP WHERE (REMARK != :1 OR REMARK IS NULL) AND UPPER(NAME) LIKE :2 ESCAPE '\' ORDER BY ID DESC FETCH FIRST :3 ROWS ONLY` ||
		bindTerbesar(t, q) != 3 {
		t.Errorf("dengan kata: %q", q)
	}
	if q := sqlTipeKendaraan("S.BD"); q != "SELECT ID, TYPENAME FROM S.BD WHERE TYPE = :1 ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY" {
		t.Errorf("tipe: %q", q)
	}
	if remarkTidakAktif != "Not Active" || tipeAneka != "ANEKA" || BatasLookupAneka != 500 || TabelShip != "SHIP" ||
		ViewBrandDetail != "BRANDDETAIL" {
		t.Error("literal RD / nama DDL berubah")
	}
}

// TestSQLOkupasiAneka - BrowseOccupationANEKA_RD: TYPE = bind, BUSINESSNAME = (peka huruf) dan NAME Contains (tak peka huruf)
// hanya bila diisi, ID DESC, batas; bind 1..n berurutan pada keempat kombinasi.
func TestSQLOkupasiAneka(t *testing.T) {
	for _, c := range []struct {
		bisnis, kata bool
		mau          string
	}{
		{false, false, "SELECT ID, NAME FROM S.O WHERE TYPE = :1 ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY"},
		{true, false, "SELECT ID, NAME FROM S.O WHERE TYPE = :1 AND BUSINESSNAME = :2 ORDER BY ID DESC FETCH FIRST :3 ROWS ONLY"},
		{false, true, `SELECT ID, NAME FROM S.O WHERE TYPE = :1 AND UPPER(NAME) LIKE :2 ESCAPE '\' ORDER BY ID DESC FETCH FIRST :3 ROWS ONLY`},
		{true, true, `SELECT ID, NAME FROM S.O WHERE TYPE = :1 AND BUSINESSNAME = :2 AND UPPER(NAME) LIKE :3 ESCAPE '\' ORDER BY ID DESC FETCH FIRST :4 ROWS ONLY`},
	} {
		if q := sqlOkupasiAneka("S.O", c.bisnis, c.kata); q != c.mau {
			t.Errorf("%v/%v: %q", c.bisnis, c.kata, q)
		}
	}
	if tipeOkupasiAneka != "ANEKA" || TabelOccupation != "OCCUPATION" {
		t.Error("literal")
	}
}

// TestSQLIndemnityPct - SearchIndemnityRate_SQL L84: INDEMNITY = :1 dan REPLACE(INDEMNITY_UNIT, ' ', ”) = :2 (padding CHAR(3) dipangkas
// seperti SQL Pega); nilai lewat angkaKeluar.
func TestSQLIndemnityPct(t *testing.T) {
	if q := sqlIndemnityPct("S.I"); q != "SELECT "+angkaKeluar("INDEMNITY_PCT")+" FROM S.I WHERE INDEMNITY = :1 AND REPLACE(INDEMNITY_UNIT, ' ', '') = :2" {
		t.Errorf("%q", q)
	}
	if TabelBIIndemnity != "M_BI_INDEMNITY" {
		t.Error("nama DDL")
	}
}
