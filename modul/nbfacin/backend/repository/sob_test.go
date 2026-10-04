package repository

import (
	"strings"
	"testing"
)

// TestSQLSOB - tiket 33: ketiga syarat work owner selalu ada dan terikat (:1 '1', :2 'LIFE
// INSURANCE'), tipe NULL ikut (butir 79.1), kolom tipe-2 berejaan DDL AGENTTPYE2; cari
// "mengandung" tidak peka huruf atas ID / ClientID / nama dengan ESCAPE; urut ID DESC.
func TestSQLSOB(t *testing.T) {
	syarat := "STATUSACTIVE = :1 AND (AGENTTPYE2 IS NULL OR AGENTTPYE2 <> :2) AND CLIENTID IS NOT NULL"
	cari := ` AND (UPPER(ID) LIKE :3 ESCAPE '\' OR UPPER(CLIENTID) LIKE :4 ESCAPE '\' OR UPPER(CLIENTNAME) LIKE :5 ESCAPE '\')`
	for got, mau := range map[string]string{
		sqlCariSOB("UJI.AGENT", false):  "SELECT ID, CLIENTID, CLIENTNAME FROM UJI.AGENT WHERE " + syarat + " ORDER BY ID DESC OFFSET :3 ROWS FETCH NEXT :4 ROWS ONLY",
		sqlCariSOB("UJI.AGENT", true):   "SELECT ID, CLIENTID, CLIENTNAME FROM UJI.AGENT WHERE " + syarat + cari + " ORDER BY ID DESC OFFSET :6 ROWS FETCH NEXT :7 ROWS ONLY",
		sqlCacahSOB("UJI.AGENT", false): "SELECT COUNT(*) FROM UJI.AGENT WHERE " + syarat,
		sqlCacahSOB("UJI.AGENT", true):  "SELECT COUNT(*) FROM UJI.AGENT WHERE " + syarat + cari,
		sqlCekSOB("UJI.AGENT"):          "SELECT COUNT(*), MAX(CLIENTNAME) FROM UJI.AGENT WHERE " + syarat + " AND ID = :3",
	} {
		if got != mau {
			t.Errorf("SQL\n%q\nmau\n%q", got, mau)
		}
	}
	if StatusAgentAktif != "1" || TipeAgentDibuang != "LIFE INSURANCE" || TabelAgent != "AGENT" {
		t.Error("nilai syarat tiket 33 berubah")
	}
	if strings.Contains(sqlCariSOB("X", true), "SELECT *") || strings.Contains(sqlCariSOB("X", true), "AGENTTYPE2") {
		t.Error("SELECT * atau ejaan kolom bukan DDL (AGENTTYPE2)")
	}
}
