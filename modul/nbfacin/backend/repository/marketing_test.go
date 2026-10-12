package repository

import (
	"strings"
	"testing"
)

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

// TestSQLMarketingBawaan - saran Marketing Name (keputusan work owner 12-10-2026): MARKETINGOFFICER ber-AKSES_LOGIN = pembuat case
// (UPPER(TRIM) dua sisi, seperti sqlTeamMarketing - beda huruf / spasi tetap cocok), MOSTATUS aktif (bind :3 = StatusMOAktif, saringan
// dropdown - MOSTATUS lain diabaikan), urut ID turun seperti dropdown; dua baris cukup untuk mengenali ganda.
func TestSQLMarketingBawaan(t *testing.T) {
	if q := sqlMarketingBawaan("S.M", "S.W"); q != "SELECT m.ID FROM S.M m JOIN S.W w ON UPPER(TRIM(m.AKSES_LOGIN)) = UPPER(TRIM(w.CREATE_OP))"+
		" WHERE w.ID = :1 AND w.LINI = :2 AND m.MOSTATUS = :3 AND m.ID IS NOT NULL ORDER BY m.ID DESC FETCH FIRST 2 ROWS ONLY" ||
		bindBerurutan(t, "marketing bawaan", q) != 3 || StatusMOAktif != "1" ||
		!strings.Contains(sqlMarketingOfficer("S.M"), "MOSTATUS = :1 ORDER BY ID DESC") {
		t.Errorf("%q", q)
	}
	var _ PembacaMarketingBawaan = (*MarketingOfficerOracle)(nil)
}
