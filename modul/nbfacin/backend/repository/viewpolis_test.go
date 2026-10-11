package repository

import (
	"reflect"
	"strings"
	"testing"
)

// TestSQLViewPolis - tiket 71 perluasan: klausa tiap jenis = CARI5 GetListViewPolicy_act 9.x / 10.1-10.7 / 11.1-11.6 / 12.1-12.3 dengan
// {ASIS} diganti bind (nol tanda kutip nilai), bind berurutan; jumlah = n bind, halaman = n + 2 (OFFSET / FETCH).
func TestSQLViewPolis(t *testing.T) {
	cases := []struct {
		jenis, klausa string
		tab           tabView
		nilai         string
		bind          []any
	}{
		{OfferFacIDPega, "PC0.IDPEGA IN (:1, :2)", tabOfferFac, "NB-1", []any{"NB-1", AwalanIDPega + "NB-1"}},
		{OfferFacNoOffer, "PC0.NOOFFER = :1", tabOfferFac, "RNM-1", []any{"RNM-1"}},
		{OfferFacReff, "PC0.RISLIPCEDING LIKE :1", tabOfferFac, "%X%", []any{"%X%"}},
		{OfferFacInsured, "PC0.INSUREDNAME LIKE :1", tabOfferFac, "X%", []any{"X%"}},
		{OfferFacSOB, "PC0.SOB LIKE :1", tabOfferFac, "X%", []any{"X%"}},
		{TreatyNoPolis, "PC0.NOPOLIS = :1 and PC0.NOENDORS is null", tabPolisTreaty, "P", []any{"P"}},
		{TreatyNoEndors, "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE NOENDORS = :1)  and PC0.NOENDORS is null", tabPolisTreaty, "E", []any{"E"}},
		{TreatyIDPega, " PC0.IDPEGA IN (:1, :2) and PC0.NOENDORS is null", tabPolisTreaty, "NB-1", []any{"NB-1", AwalanIDPega + "NB-1"}},
		{TreatyEDMT, "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE IDPEGA IN (:1, :2))  and PC0.NOENDORS is null", tabPolisTreaty,
			"EDMT-1", []any{"EDMT-1", AwalanIDPega + "EDMT-1"}},
		{TreatyMaster, "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE NOOFFER = :1)  and PC0.NOENDORS is null", tabPolisTreaty, "M", []any{"M"}},
		{TreatyInsured, "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE INSUREDNAME LIKE :1)  and PC0.NOENDORS is null", tabPolisTreaty, "%X%", []any{"%X%"}},
		{TreatySOB, "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE SOB LIKE :1)  and PC0.NOENDORS is null", tabPolisTreaty, "%X%", []any{"%X%"}},
		{LifeNoPolis, "NOPOLIS = :1 and NOENDORS is null", tabPolisLife, "P", []any{"P"}},
		{LifeDigitAkhir, "NOPOLIS like :1 and NOENDORS is null", tabPolisLife, "%00012", []any{"%00012"}},
		{LifeIDPega, " IDPEGA IN (:1, :2) and NOENDORS is null", tabPolisLife, "NBLF-1", []any{"NBLF-1", AwalanIDPega + "NBLF-1"}},
		{LifeNoEndors, "NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE NOENDORS  = :1) and NOENDORS is null", tabPolisLife, "E", []any{"E"}},
		{LifeEDML, "NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM T WHERE IDPEGA IN (:1, :2)) and NOENDORS is null", tabPolisLife, "EDML-1",
			[]any{"EDML-1", AwalanIDPega + "EDML-1"}},
		{LifeSOB, "(SOBNAME LIKE :1 OR CEDINGCONAME LIKE :2 OR POLICYHOLDERNAME LIKE :3)  AND NOENDORS IS NULL", tabPolisLife, "%X%",
			[]any{"%X%", "%X%", "%X%"}},
		{OfferLifeIDPega, " A.IDPEGA IN (:1, :2)", tabOfferLife, "NBLF-1", []any{"NBLF-1", AwalanIDPega + "NBLF-1"}},
		{OfferLifeOFR, " A.ID = SUBSTR(:1, 5)", tabOfferLife, "OFR-12", []any{"OFR-12"}},
		{OfferLifeCeding, " A.CEDINGCONAME like :1 OR  A.POLICYHOLDERNAME like :2", tabOfferLife, "%X%", []any{"%X%", "%X%"}},
	}
	for _, c := range cases {
		k, ada := klausaViewDari(c.jenis, "T")
		if !ada || k.teks != c.klausa || strings.Contains(k.teks, "'") {
			t.Errorf("%s: %q", c.jenis, k.teks)
		}
		if b := bindView(k, c.nilai); !reflect.DeepEqual(b, c.bind) {
			t.Errorf("%s bind %v", c.jenis, b)
		}
		dasar := c.tab.dasar("T", k.teks)
		n := len(c.bind)
		if bindBerurutan(t, c.jenis+" jumlah", sqlJumlahView(dasar)) != n || bindBerurutan(t, c.jenis+" halaman", sqlHalamanView(c.tab, dasar, n)) != n+2 {
			t.Errorf("%s jumlah bind", c.jenis)
		}
	}
	if len(cases) != 21 {
		t.Errorf("jenis %d", len(cases))
	}
	if _, ada := klausaViewDari("lain", "T"); ada {
		t.Error("jenis lain")
	}
}

// TestSQLDasarViewPolis - SQL dasar = pyBrowseSQL keempat rule (kolom, DISTINCT, kurung WHERE), case id dua bentuk, tabel berkualifikasi
// (BrowseOfferLife subkueri C Pega tanpa skema - di sini ikut berkualifikasi), halaman TO_CHAR DD-MM-YYYY + urutan stabil.
func TestSQLDasarViewPolis(t *testing.T) {
	ci := "CASE WHEN SUBSTR(PC0.IDPEGA, 1, 18) = 'ASM-FW-GISFW-WORK ' THEN SUBSTR(PC0.IDPEGA, 19) ELSE PC0.IDPEGA END"
	for nama, c := range map[string]struct{ q, mau string }{
		"offer":      {sqlDasarOfferFac("S.FO", "K"), "SELECT DISTINCT " + ci + " AS CASEID, PC0.NOOFFER, PC0.BEGINDATE, PC0.ENDDATE, PC0.SOB, PC0.INSUREDNAME FROM S.FO PC0 WHERE (K)"},
		"treaty":     {sqlDasarPolisTreaty("S.TIP", "K"), "SELECT DISTINCT SUBSTR(" + ci + ", 1, 10) AS CASEID, PC0.NOPOLIS, PC0.BEGINDATE, PC0.ENDDATE, PC0.SOB, PC0.INSUREDNAME FROM S.TIP PC0 WHERE (K)"},
		"life":       {sqlDasarPolisLife("S.LIP", "K"), "SELECT SUBSTR(CASE WHEN SUBSTR(IDPEGA, 1, 18) = 'ASM-FW-GISFW-WORK ' THEN SUBSTR(IDPEGA, 19) ELSE IDPEGA END, 1, 10) AS CASEID, NOPOLIS, NOENDORS, CEDINGCONAME, POLICYHOLDERNAME AS QQ, TYPE AS TYPEPOLICY, SOBNAME, RISLIPRNM FROM S.LIP WHERE K"},
		"offer life": {sqlDasarOfferLife("S.J", "K"), "SELECT SUBSTR(CASE WHEN SUBSTR(A.IDPEGA, 1, 18) = 'ASM-FW-GISFW-WORK ' THEN SUBSTR(A.IDPEGA, 19) ELSE A.IDPEGA END, 1, 10) AS CASEID, 'OFR-' || D.ID AS NOOFFER, A.TGL_INPUT, A.CEDINGCONAME, A.POLICYHOLDERNAME AS QQ, A.TYPECEDING AS TYPEPOLICY, A.BUSINESSNAME, A.STATUS AS PRODKE FROM S.J A, (SELECT C.ID, C.IDPEGA FROM S.J C WHERE C.TGL_INPUT = (SELECT MIN(b.TGL_INPUT) FROM S.J b WHERE b.IDPEGA = C.IDPEGA)) D WHERE A.OLDID IS NULL AND A.TGL_INPUT = (SELECT MAX(b.TGL_INPUT) FROM S.J b WHERE b.IDPEGA = A.IDPEGA) AND A.IDPEGA = D.IDPEGA AND (K)"},
	} {
		if c.q != c.mau {
			t.Errorf("%s:\n%s\n%s", nama, c.q, c.mau)
		}
	}
	for _, c := range []struct {
		tab        tabView
		tabel      string
		awal, urut string
	}{
		{tabOfferFac, TabelFacInOffer, "SELECT CASEID, NOOFFER, INSUREDNAME, SOB, TO_CHAR(BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(ENDDATE, 'DD-MM-YYYY') FROM (",
			") ORDER BY NOOFFER, CASEID, BEGINDATE, ENDDATE, SOB, INSUREDNAME OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"},
		{tabPolisTreaty, TabelTreatyInProduction, "SELECT CASEID, NOPOLIS, INSUREDNAME, SOB, TO_CHAR(BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(ENDDATE, 'DD-MM-YYYY') FROM (",
			") ORDER BY NOPOLIS, CASEID, BEGINDATE, ENDDATE, SOB, INSUREDNAME OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"},
		{tabPolisLife, TabelLifeInProduction, "SELECT CASEID, NOPOLIS, NOENDORS, TYPEPOLICY, SOBNAME, CEDINGCONAME, QQ, RISLIPRNM FROM (",
			") ORDER BY NOPOLIS, NOENDORS NULLS FIRST, CASEID, TYPEPOLICY, SOBNAME, CEDINGCONAME, QQ, RISLIPRNM OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"},
		{tabOfferLife, TabelJSONOfferLife, "SELECT CASEID, NOOFFER, TYPEPOLICY, BUSINESSNAME, CEDINGCONAME, QQ, PRODKE FROM (",
			") ORDER BY TGL_INPUT DESC, CASEID, NOOFFER OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"},
	} {
		h := sqlHalamanView(c.tab, "D", 1)
		if h != c.awal+"D"+c.urut {
			t.Errorf("%s: %q", c.tab.nama, h)
		}
	}
	if TabelFacInOffer != "FACINOFFER" || TabelTreatyInProduction != "TREATYINPRODUCTION" || TabelLifeInProduction != "LIFEINPRODUCTION" ||
		TabelJSONOfferLife != "JSON_OFFER_LIFE" {
		t.Error("nama tabel POOLDATA")
	}
}

// TestSQLBerkasPolisTreaty - salinan sqlBerkasPolis Claim Non Prop (OQ-CNP-13): satu bind, PRODKE terbesar, JOIN T_WORK_POLIS.
func TestSQLBerkasPolisTreaty(t *testing.T) {
	q := sqlBerkasPolisTreaty("S.G", "S.W")
	if q != "SELECT g.ID, g.PRODKE FROM S.G g JOIN S.W w ON w.ID = g.ID WHERE g.NOPOLIS = :1 ORDER BY g.PRODKE DESC FETCH FIRST 1 ROWS ONLY" ||
		bindBerurutan(t, "berkas", q) != 1 {
		t.Errorf("%q", q)
	}
	if TabelGeneralPolisTreaty != "T_GENERAL_POLIS_TREATY" || TabelWorkPolis != "T_WORK_POLIS" {
		t.Error("nama tabel")
	}
}
