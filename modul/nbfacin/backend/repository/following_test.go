package repository

import (
	"strings"
	"testing"
)

// TestSQLFollowing - tiket 70: BrowseDataViewPolis1 (kolom + subkueri RiSlipCeding, DISTINCT, NOENDORS IS NULL), klausa 4.x ber-bind
// (A70-1, nol literal nilai), halaman ber-bind :2 / :3, FOLLOWING (+ PERCENT_SHARE / PAY_RI_COMMISION).
func TestSQLFollowing(t *testing.T) {
	jenis := []string{CariFollowingNoPolis, CariFollowingIDPega, CariFollowingNoEndors, CariFollowingEDM, CariFollowingReff,
		CariFollowingInsured, CariFollowingSOB, CariFollowingQQ}
	mau := map[string]string{
		CariFollowingNoPolis:  "PC0.NOPOLIS = :1 AND PC0.NOENDORS IS NULL",
		CariFollowingIDPega:   "PC0.IDPEGA = :1 AND PC0.NOENDORS IS NULL",
		CariFollowingNoEndors: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM P WHERE NOENDORS = :1) AND PC0.NOENDORS IS NULL",
		CariFollowingEDM:      "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM P WHERE IDPEGA = :1) AND PC0.NOENDORS IS NULL",
		CariFollowingReff:     "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM P WHERE RISLIPCEDING LIKE :1) AND PC0.NOENDORS IS NULL",
		CariFollowingInsured:  "PC0.INSUREDNAME LIKE :1 AND PC0.NOENDORS IS NULL",
		CariFollowingSOB:      "PC0.SOB LIKE :1 AND PC0.NOENDORS IS NULL",
		CariFollowingQQ:       "PC0.QQNAME LIKE :1 AND PC0.NOENDORS IS NULL",
	}
	for _, j := range jenis {
		k, ada := klausaFollowing("P", j)
		if !ada || k != mau[j] {
			t.Errorf("%s: %q", j, k)
		}
		dasar := sqlDasarFollowing("P", k)
		if bindBerurutan(t, j+" jumlah", sqlJumlahFollowing(dasar)) != 1 || bindBerurutan(t, j+" halaman", sqlHalamanFollowing(dasar)) != 3 {
			t.Errorf("%s bind", j)
		}
	}
	if _, ada := klausaFollowing("P", "lain"); ada {
		t.Error("jenis lain")
	}
	dasar := sqlDasarFollowing("P", "X")
	if !strings.HasPrefix(dasar, "SELECT DISTINCT SUBSTR(PC0.IDPEGA, 19, 10) AS CASEID, PC0.NOPOLIS, PC0.BEGINDATE, PC0.ENDDATE, PC0.SOB, PC0.INSUREDNAME, PC0.QQNAME, (SELECT MAX(B.RISLIPCEDING) KEEP (DENSE_RANK LAST ORDER BY B.SPREAD_DATE) FROM P B WHERE B.NOPOLIS = PC0.NOPOLIS) AS RISLIP FROM P PC0 WHERE X") {
		t.Errorf("%q", dasar)
	}
	if h := sqlHalamanFollowing(dasar); !strings.Contains(h, "TO_CHAR(BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(ENDDATE, 'DD-MM-YYYY')") ||
		!strings.HasSuffix(h, "ORDER BY NOPOLIS, CASEID, BEGINDATE, ENDDATE, SOB, INSUREDNAME, QQNAME OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY") {
		t.Errorf("%q", h)
	}
	if len(AwalanIDPega) != 18 {
		t.Errorf("SUBSTR(IDPEGA, 19, ...) mensyaratkan awalan 18 karakter: %d", len(AwalanIDPega))
	}
	for nama, c := range map[string]struct {
		q    string
		bind int
	}{
		"ada":         {sqlAdaPolisFollowing("P"), 2},
		"ubah":        {sqlFollowing("G"), 2},
		"sisip":       {sqlSisipFollowing("G"), 2},
		"salin":       {sqlSalinFollowing("G"), 4},
		"sisip salin": {sqlSisipSalinFollowing("G"), 4},
	} {
		if bindBerurutan(t, nama, c.q) != c.bind {
			t.Errorf("%s: %q", nama, c.q)
		}
	}
	if !strings.HasPrefix(sqlSalinFollowing("G"), "UPDATE G SET FOLLOWING = :1, PERCENT_SHARE = TO_NUMBER(:2,") ||
		sqlFollowing("G") != "UPDATE G SET FOLLOWING = :1 WHERE ID = :2" ||
		sqlAdaPolisFollowing("P") != "SELECT COUNT(*) FROM P WHERE IDPEGA = :1 AND NOPOLIS = :2 AND NOENDORS IS NULL" {
		t.Error("teks SQL")
	}
}
