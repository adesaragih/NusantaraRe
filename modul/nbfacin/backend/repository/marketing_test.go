package repository

import "testing"

// TestSQLMarketingOfficer - RD BrowseMarketingOfficer_RD: filter MOStatus = "1" (bind),
// urut ID DESC, pyMaxRecords 500 (bind); kolom ID (nilai) dan CLIENTNAME (teks).
func TestSQLMarketingOfficer(t *testing.T) {
	mau := "SELECT ID, CLIENTNAME FROM UJI.MARKETINGOFFICER WHERE MOSTATUS = :1 ORDER BY ID DESC FETCH FIRST :2 ROWS ONLY"
	if got := sqlMarketingOfficer("UJI.MARKETINGOFFICER"); got != mau {
		t.Errorf("SQL\n%q\nmau\n%q", got, mau)
	}
	if StatusMOAktif != "1" || BatasMarketingOfficer != 500 || TabelMarketingOfficer != "MARKETINGOFFICER" {
		t.Error("konstanta RD berubah")
	}
}
